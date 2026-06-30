package main

import (
	"context"
	"fmt"
	"log"
	"sovereign-node-go/pkg/github"
	"sovereign-node-go/pkg/tickets"

	"github.com/google/uuid"
)

func main() {
	fmt.Println("[VERIFY] Starting Sovereign Mesh v7.1 Fabric Verification...")

	// 1. Initialize Manager
	// We'll use a dummy github client for now as we're testing the DB part
	client := github.NewClient("") 
	mgr, err := tickets.NewManager(client, "owner", "repo")
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize manager: %v", err)
	}

	ctx := context.Background()

	// 2. Create a Layer 2 Ticket (Child of Genesis)
	genesisID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	agentID := "GATE-01"

	content := tickets.FabricContent{
		IntentBlob: map[string]interface{}{
			"action": "BOOT_SEQUENCE",
			"params": map[string]string{"mode": "FORENSIC"},
		},
		ConsensusScore: 0.95,
		RawContent:     []byte("System boot sequence initiated by Gate-01."),
	}

	ticketID, err := mgr.CreateFabricTicketV71(ctx, 2, agentID, content)
	if err != nil {
		log.Fatalf("[FATAL] Failed to create ticket: %v", err)
	}
	fmt.Printf("[SUCCESS] Created Ticket: %s\n", ticketID)

	// 3. Link to Genesis
	err = mgr.LinkTicketsV71(ctx, genesisID, ticketID, tickets.RelEvolution)
	if err != nil {
		log.Fatalf("[FATAL] Failed to link ticket: %v", err)
	}
	fmt.Println("[SUCCESS] Linked to Genesis Block.")

	// 4. Update Agent Memory
	err = mgr.UpdateAgentMemory(ctx, agentID, ticketID)
	if err != nil {
		log.Fatalf("[FATAL] Failed to update agent memory: %v", err)
	}
	fmt.Printf("[SUCCESS] Updated Agent %s memory.\n", agentID)

	// 5. Retrieve Context (Anti-Hallucination Check)
	snapshot, err := mgr.GetContextV71(ctx, agentID)
	if err != nil {
		log.Fatalf("[FATAL] Failed to retrieve context: %v", err)
	}

	fmt.Println("\n--- Context Snapshot ---")
	fmt.Printf("Active Ticket: %s (Layer %d)\n", snapshot.ActiveTicket.ID, snapshot.ActiveTicket.LayerID)
	fmt.Printf("Content: %s\n", string(snapshot.Content.RawContent))
	fmt.Printf("Ancestors: %d\n", len(snapshot.Ancestors))
	for _, a := range snapshot.Ancestors {
		fmt.Printf("  - Ancestor: %s (Layer %d)\n", a.ID, a.LayerID)
	}
	fmt.Println("------------------------")

	fmt.Println("[VERIFY] Verification Complete.")
}
