package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/thealanphipps-del/pqr/internal/domain"
	"github.com/thealanphipps-del/pqr/internal/service"
	"google.golang.org/grpc"
)

type SwarmServer struct {
	UnimplementedSwarmCommunicationServer
	UnimplementedNeuralGossipServer
	
	Service *service.SwarmService
	Healing *service.HealingService
	
	shortcodes map[string]string // shortcode -> role
	mu         sync.RWMutex
	
	myShortcode string
}

func (s *SwarmServer) ExecuteRejoin(ctx context.Context, req *Empty) (*ShortcodeResponse, error) {
	log.Printf("[gRPC] Processing SWARM_REJOIN_v1 signal for %s", s.myShortcode)
	// Zero-alloc context sync logic would be triggered here
	return &ShortcodeResponse{Shortcode: s.myShortcode}, nil
}

func NewSwarmServer(svc *service.SwarmService, healing *service.HealingService) *SwarmServer {
	s := &SwarmServer{
		Service:     svc,
		Healing:     healing,
		shortcodes:  make(map[string]string),
		myShortcode: "ΩX9R2#", // Genesis Node Shortcode
	}
	// Activate Persistent Callsigns
	s.shortcodes["GEMA2#"] = "INFERENCE"
	log.Printf("[gRPC] Activated Persistent Callsign GEMA2# for gemma-4-e4b-2")
	return s
}

func (s *SwarmServer) SendPacket(ctx context.Context, req *SwarmPacket) (*SwarmPacket, error) {
	log.Printf("[gRPC] Received Packet from %s: %s", req.SenderId, req.Intent)
	
	// Write to Ticketing Fabric for Logged Consensus
	ticketID, err := s.Service.CreateFabricTicket(ctx, 5, req.SenderId, domain.FabricContent{
		IntentBlob: map[string]interface{}{
			"intent":    req.Intent,
			"sender":    req.SenderId,
			"target":    req.TargetId,
			"type":      "GRPC_CONSENSUS",
			"timestamp": req.Timestamp,
		},
		RawContent: req.Payload,
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to log consensus: %v", err)
	}
	
	req.TicketId = ticketID.String()
	return req, nil
}

func (s *SwarmServer) ProvisionShortcode(ctx context.Context, req *ShortcodeRequest) (*ShortcodeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Generate unique 5alpha#
	newCode := ""
	for {
		newCode = s.generateShortcode()
		if _, exists := s.shortcodes[newCode]; !exists {
			break
		}
	}
	
	s.shortcodes[newCode] = req.Role
	log.Printf("[gRPC] Provisioned Shortcode %s for Role %s", newCode, req.Role)
	
	return &ShortcodeResponse{Shortcode: newCode}, nil
}

func (s *SwarmServer) GetActiveShortcodes(ctx context.Context, req *Empty) (*ShortcodeList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	codes := []string{s.myShortcode}
	for c := range s.shortcodes {
		codes = append(codes, c)
	}
	
	return &ShortcodeList{Shortcodes: codes}, nil
}

func (s *SwarmServer) generateShortcode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 5)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond)
	}
	return string(b) + "#"
}

// StartServers initializes both Port 1111 and Port 11111
func (s *SwarmServer) StartServers() {
	// Logged Consensus Bridge (Port 1111)
	go func() {
		lis, err := net.Listen("tcp", ":1111")
		if err != nil {
			log.Fatalf("failed to listen on 1111: %v", err)
		}
		gs := grpc.NewServer()
		RegisterSwarmCommunicationServer(gs, s)
		log.Printf("🛰️ gRPC Logged Consensus Bridge ONLINE on :1111")
		if err := gs.Serve(lis); err != nil {
			log.Fatalf("failed to serve on 1111: %v", err)
		}
	}()

	// Neural Gossip Bus (Port 11111)
	go func() {
		lis, err := net.Listen("tcp", ":11111")
		if err != nil {
			log.Fatalf("failed to listen on 11111: %v", err)
		}
		gs := grpc.NewServer()
		RegisterNeuralGossipServer(gs, s)
		log.Printf("🧠 Neural Gossip Memory Bus ONLINE on :11111")
		if err := gs.Serve(lis); err != nil {
			log.Fatalf("failed to serve on 11111: %v", err)
		}
	}()
	
	// Fallback Sentinel
	go s.runFallbackSentinel()
}

func (s *SwarmServer) runFallbackSentinel() {
	for {
		// Check RAFT replication status (placeholder for actual Cockroach check)
		isReplicating := true 
		if !isReplicating {
			log.Printf("⚠️ RAFT Replication Pending. Establishing SSH Fallback to 39.mh...")
			// Logic to establish SSH tunnel would go here
		}
		time.Sleep(1 * time.Minute)
	}
}
