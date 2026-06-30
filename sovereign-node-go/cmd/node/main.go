package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"sovereign-node-go/pkg/mesh"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Sovereign Node Go: Initializing...")

	// Phase 2: Mesh Re-Anchoring
	// Attempt to re-anchor to a known peer if needed
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Example: Target the bridge IP for mesh re-anchoring
	go func() {
		err := mesh.ReAnchor(ctx, "204.168.138.60:1111")
		if err != nil {
			fmt.Printf("[WARNING] Mesh re-anchor failed: %v\n", err)
		}
	}()

	// Default to port 1111 as the engine port
	port := ":1111"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	fmt.Printf("Sovereign Mesh Listener active on port %s\n", port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
