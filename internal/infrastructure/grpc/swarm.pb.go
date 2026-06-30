// Code generated manually for compilation purposes (no protoc available).
// It provides the minimal set of types/interfaces required by internal/infrastructure/grpc/server.go.

package grpc

import (
	context "context"
	reflect "reflect"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Messages

type SwarmPacket struct {
	SenderId  string `protobuf:"bytes,1,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	TargetId  string `protobuf:"bytes,2,opt,name=target_id,json=targetId,proto3" json:"target_id,omitempty"`
	Intent    string `protobuf:"bytes,3,opt,name=intent,proto3" json:"intent,omitempty"`
	Payload   []byte `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`
	Timestamp int64  `protobuf:"varint,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	TicketId  string `protobuf:"bytes,6,opt,name=ticket_id,json=ticketId,proto3" json:"ticket_id,omitempty"`
}

type StrikeRequest struct {
	TicketId     string `protobuf:"bytes,1,opt,name=ticket_id,json=ticketId,proto3" json:"ticket_id,omitempty"`
	LogicPayload string `protobuf:"bytes,2,opt,name=logic_payload,json=logicPayload,proto3" json:"logic_payload,omitempty"`
	Signature     []byte `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
}

type StrikeResponse struct {
	ExitCode   int32  `protobuf:"varint,1,opt,name=exit_code,json=exitCode,proto3" json:"exit_code,omitempty"`
	LogOutput  string `protobuf:"bytes,2,opt,name=log_output,json=logOutput,proto3" json:"log_output,omitempty"`
	ProofHash  string `protobuf:"bytes,3,opt,name=proof_hash,json=proofHash,proto3" json:"proof_hash,omitempty"`
}

type BlackholeUpdate struct {
	IdentifierHash []byte `protobuf:"bytes,1,opt,name=identifier_hash,json=identifierHash,proto3" json:"identifier_hash,omitempty"`
	IsRevoked      bool   `protobuf:"varint,2,opt,name=is_revoked,json=isRevoked,proto3" json:"is_revoked,omitempty"`
}

type TelemetryRequest struct{ NodeId string `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"` }

type TelemetryData struct {
	VitalityScore float64 `protobuf:"fixed64,1,opt,name=vitality_score,json=vitalityScore,proto3" json:"vitality_score,omitempty"`
	Slope         float64 `protobuf:"fixed64,2,opt,name=slope,proto3" json:"slope,omitempty"`
}

type MemoryPage struct {
	PageId string `protobuf:"bytes,1,opt,name=page_id,json=pageId,proto3" json:"page_id,omitempty"`
	Data    []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Offset  uint64 `protobuf:"varint,3,opt,name=offset,proto3" json:"offset,omitempty"`
}

type ShortcodeRequest struct{ Role string `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"` }
type ShortcodeResponse struct{ Shortcode string `protobuf:"bytes,1,opt,name=shortcode,proto3" json:"shortcode,omitempty"` }
type ShortcodeList struct{ Shortcodes []string `protobuf:"bytes,1,rep,name=shortcodes,proto3" json:"shortcodes,omitempty"` }
type Empty struct{}

// GOSSIP SERVICE

type NeuralGossipClient interface {
	StreamVitality(ctx context.Context, in *TelemetryRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ *TelemetryData ], error)
	MemoryPageSwap(ctx context.Context, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ *MemoryPage ], error)
}

type NeuralGossipServer interface {
	StreamVitality(*TelemetryRequest, grpc.ServerStream) error
	MemoryPageSwap(grpc.ServerStream) error
	mustEmbedUnimplementedNeuralGossipServer()
}

type UnimplementedNeuralGossipServer struct{}

func (UnimplementedNeuralGossipServer) StreamVitality(*TelemetryRequest, grpc.ServerStream) error {
	return status.Errorf(codes.Unimplemented, "method StreamVitality not implemented")
}
func (UnimplementedNeuralGossipServer) MemoryPageSwap(grpc.ServerStream) error {
	return status.Errorf(codes.Unimplemented, "method MemoryPageSwap not implemented")
}
func (UnimplementedNeuralGossipServer) mustEmbedUnimplementedNeuralGossipServer() {}

// SWARM SERVICE

type SwarmCommunicationClient interface {
	SendPacket(ctx context.Context, in *SwarmPacket, opts ...grpc.CallOption) (*SwarmPacket, error)
	StreamDeliberation(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStream, error)
	ProvisionShortcode(ctx context.Context, in *ShortcodeRequest, opts ...grpc.CallOption) (*ShortcodeResponse, error)
	GetActiveShortcodes(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ShortcodeList, error)
	ExecuteStrike(ctx context.Context, in *StrikeRequest, opts ...grpc.CallOption) (*StrikeResponse, error)
	SyncBlackhole(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStream, error)
}

type SwarmCommunicationServer interface {
	SendPacket(context.Context, *SwarmPacket) (*SwarmPacket, error)
	StreamDeliberation(grpc.ServerStream) error
	ProvisionShortcode(context.Context, *ShortcodeRequest) (*ShortcodeResponse, error)
	GetActiveShortcodes(context.Context, *Empty) (*ShortcodeList, error)
	ExecuteStrike(context.Context, *StrikeRequest) (*StrikeResponse, error)
	SyncBlackhole(grpc.ServerStream) error
	mustEmbedUnimplementedSwarmCommunicationServer()
}

type UnimplementedSwarmCommunicationServer struct{}

func (UnimplementedSwarmCommunicationServer) SendPacket(context.Context, *SwarmPacket) (*SwarmPacket, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendPacket not implemented")
}
func (UnimplementedSwarmCommunicationServer) StreamDeliberation(grpc.ServerStream) error {
	return status.Errorf(codes.Unimplemented, "method StreamDeliberation not implemented")
}
func (UnimplementedSwarmCommunicationServer) ProvisionShortcode(context.Context, *ShortcodeRequest) (*ShortcodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProvisionShortcode not implemented")
}
func (UnimplementedSwarmCommunicationServer) GetActiveShortcodes(context.Context, *Empty) (*ShortcodeList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetActiveShortcodes not implemented")
}
func (UnimplementedSwarmCommunicationServer) ExecuteStrike(context.Context, *StrikeRequest) (*StrikeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteStrike not implemented")
}
func (UnimplementedSwarmCommunicationServer) SyncBlackhole(grpc.ServerStream) error {
	return status.Errorf(codes.Unimplemented, "method SyncBlackhole not implemented")
}
func (UnimplementedSwarmCommunicationServer) mustEmbedUnimplementedSwarmCommunicationServer() {}

// Registration helpers (minimal)

func RegisterSwarmCommunicationServer(s grpc.ServiceRegistrar, srv SwarmCommunicationServer) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "swarm.SwarmCommunication",
		HandlerType: (*SwarmCommunicationServer)(nil),
		Methods: []grpc.MethodDesc{
			{MethodName: "SendPacket", Handler: unimplementedMethodHandler},
			{MethodName: "ProvisionShortcode", Handler: unimplementedMethodHandler},
			{MethodName: "GetActiveShortcodes", Handler: unimplementedMethodHandler},
			{MethodName: "ExecuteStrike", Handler: unimplementedMethodHandler},
		},
		Streams: []grpc.StreamDesc{
			{StreamName: "StreamDeliberation", Handler: unimplementedStreamHandler, ServerStreams: true},
			{StreamName: "SyncBlackhole", Handler: unimplementedStreamHandler, ServerStreams: true},
		},
		Metadata: "proto/swarm.proto",
	}, srv)
}

func RegisterNeuralGossipServer(s grpc.ServiceRegistrar, srv NeuralGossipServer) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "swarm.NeuralGossip",
		HandlerType: (*NeuralGossipServer)(nil),
		Methods: []grpc.MethodDesc{},
		Streams: []grpc.StreamDesc{
			{StreamName: "StreamVitality", Handler: unimplementedStreamHandler, ServerStreams: true},
			{StreamName: "MemoryPageSwap", Handler: unimplementedStreamHandler, ServerStreams: true},
		},
		Metadata: "proto/swarm.proto",
	}, srv)
}

func unimplementedMethodHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented stub")
}
func unimplementedStreamHandler(srv interface{}, stream grpc.ServerStream) error {
	return status.Error(codes.Unimplemented, "unimplemented stub")
}

var _ = reflect.TypeOf
