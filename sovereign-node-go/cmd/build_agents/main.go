package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sovereign-node-go/pkg/agent"
	"sovereign-node-go/pkg/github"
	"sovereign-node-go/pkg/tickets"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("[AGENTS] Initializing Agent Builds & Identities...")

	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/antigravity?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	client := github.NewClient("")
	ticketMgr, err := tickets.NewManager(client, "owner", "repo")
	if err != nil {
		log.Fatalf("Failed to initialize ticket manager: %v", err)
	}
	// We need to inject the db into ticket manager if possible. Wait, ticketMgr already connects to the db inside NewManager.
	// Oh wait, tickets.NewManager creates its own connection. Let's see how it works.
	// We'll just pass db to AgentManager.
	agentMgr := agent.NewAgentManager(db, ticketMgr)

	ctx := context.Background()

	err = agentMgr.InitSchema(ctx)
	if err != nil {
		log.Fatalf("Failed to init agent schema: %v", err)
	}

	identities, err := agentMgr.BuildAgents(ctx)
	if err != nil {
		log.Fatalf("Failed to build agents: %v", err)
	}

	fmt.Printf("[SUCCESS] Built and assigned identities to %d agents.\n", len(identities))

	// Simulate communication on Port 11111 for a few agents
	fmt.Println("[P2P] Simulating port 11111 communication...")
	for i := 0; i < 3 && i < len(identities); i++ {
		id := identities[i]
		fmt.Printf("--> [PORT 11111] %s: Handshake initialized from %s\n", id.Shortcode, id.AgentID)
	}
	fmt.Println("[P2P] Everything is tracked via Layer-2 Fabric Tickets.")
}
