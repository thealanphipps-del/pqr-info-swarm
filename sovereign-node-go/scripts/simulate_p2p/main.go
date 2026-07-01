package main

import (
	"context"
	"fmt"
	"sovereign-node-go/pkg/mesh"
)

func main() {
	fmt.Println("[P2P TEST] Simulating Gemini Handshake and Z-DSP Sync...")

	p2p := mesh.NewP2PManager("LOCAL-NODE-B")
	ctx := context.Background()

	// 1. Handshake
	remote, err := p2p.GeminiHandshake(ctx, "localhost:11111")
	if err != nil {
		fmt.Printf("[ERROR] Handshake failed: %v\n", err)
		return
	}
	fmt.Printf("[SUCCESS] Anchored to %s (Epoch %.1f)\n", remote.ID, remote.Epoch)

	// 2. Sync
	remoteState := map[string]string{
		"GENESIS_HASH": "0000000000000000000000000000000000000000000000000000000000000000",
		"ACTIVE_HASH":  "16bfe8ea48924fd48a8038f94937605d",
	}

	divergences, err := p2p.ZDSPReconcile(ctx, remoteState)
	if err != nil {
		fmt.Printf("[ERROR] Sync failed: %v\n", err)
		return
	}

	if len(divergences) == 0 {
		fmt.Println("[VERIFY] Node-B is now forensic-aligned with the Mesh.")
	}
}
