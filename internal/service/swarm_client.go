package service

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/thealanphipps-del/pqr/internal/execution"
	"github.com/thealanphipps-del/pqr/internal/infrastructure/db"
	pb "github.com/thealanphipps-del/pqr/proto"
)

type SwarmClient struct {
	grpc pb.SwarmCommunicationClient
	conn *grpc.ClientConn
}

type AgentInfo struct {
	AgentID   string
	Shortcode string
}

type CapabilityManifest struct {
	AgentType    string   `json:"agent_type"`
	OS           string   `json:"os"`
	Capabilities []string `json:"capabilities"`
}

func NewSwarmClient() *SwarmClient {
	addr := os.Getenv("PQR_SWARM_ADDR")
	if addr == "" {
		addr = "localhost:8196" // Changed to 8196 to match the pqr server default port
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // use insecure credentials for dev
	)
	if err != nil {
		log.Printf("[WARNING] failed to dial swarm gRPC at %s: %v", addr, err)
	}

	client := pb.NewSwarmCommunicationClient(conn)

	return &SwarmClient{
		grpc: client,
		conn: conn,
	}
}

func (c *SwarmClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *SwarmClient) Register(manifest CapabilityManifest) (*AgentInfo, error) {
	req := &pb.ShortcodeRequest{
		Role: manifest.AgentType,
	}

	resp, err := c.grpc.ProvisionShortcode(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &AgentInfo{
		AgentID:   "swen-node-1", // Would come from resp if proto updated
		Shortcode: resp.Shortcode,
	}, nil
}

// StartExecutionStream routes physical execution commands to the SWEN.
func StartExecutionStream(
	client *SwarmClient,
	engine *execution.ExecutionEngine,
	memory *db.CockroachRepository,
) {
	ctx := context.Background()
	stream, err := client.grpc.OpenExecutionStream(ctx)
	if err != nil {
		log.Printf("failed to open execution stream: %v", err)
		return
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("execution stream recv error: %v", err)
			return
		}

		wf := &execution.WorkflowDescriptor{
			Kind:     execution.WorkflowKind(msg.Kind),
			Command:  msg.Command,
			Args:     msg.Args,
			TargetOS: msg.TargetOs,
			Metadata: decodeMetadata(msg.MetadataJson),
		}

		sig := memory.HashSignatureFromFields(
			msg.Kind,
			msg.Command,
			msg.TargetOs,
		)
		if fix, ok := memory.QueryKnownFix(sig); ok {
			wf.Command = fix.ApplyToCommand(wf.Command)
		}

		res, execErr := engine.ExecuteWorkflow(ctx, wf)

		if execErr != nil || !res.Success {
			memory.LogErrorSolution(sig, msg, res, execErr)
		}

		if err := stream.Send(&pb.ExecutionResult{
			Success: res.Success,
			Logs:    res.Logs,
			Error:   errorToString(execErr),
		}); err != nil {
			log.Printf("execution stream send error: %v", err)
			return
		}
	}
}

func decodeMetadata(raw string) map[string]interface{} {
	if raw == "" {
		return map[string]interface{}{}
	}
	var out map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return map[string]interface{}{}
	}
	return out
}

func errorToString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
