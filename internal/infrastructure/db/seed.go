package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/thealanphipps-del/pqr/internal/domain"
)

func (r *CockroachRepository) SeedGovernanceAgents(ctx context.Context) error {
	agents := []domain.Agent{
		{ID: "council-001", Role: "GOVERNANCE", Level: 10, Metadata: map[string]interface{}{"title": "High Justiciar"}},
		{ID: "council-002", Role: "GOVERNANCE", Level: 10, Metadata: map[string]interface{}{"title": "Master of Lineage"}},
		{ID: "council-003", Role: "GOVERNANCE", Level: 10, Metadata: map[string]interface{}{"title": "Forensic Auditor"}},
		{ID: "council-004", Role: "GOVERNANCE", Level: 10, Metadata: map[string]interface{}{"title": "Swarm Architect"}},
		{ID: "council-005", Role: "GOVERNANCE", Level: 10, Metadata: map[string]interface{}{"title": "Sticky Rule Guardian"}},
		{ID: "monitor-001", Role: "MONITOR", Level: 5, Metadata: map[string]interface{}{"title": "Log Sweeper"}},
	}

	for _, a := range agents {
		metaJSON, _ := json.Marshal(a.Metadata)
		_, err := r.db.ExecContext(ctx, `
			INSERT INTO agents (agent_id, role, level, metadata)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (agent_id) DO UPDATE SET
				role = $2,
				level = $3,
				metadata = $4
		`, a.ID, a.Role, a.Level, metaJSON)
		if err != nil {
			return fmt.Errorf("failed to seed agent %s: %v", a.ID, err)
		}
	}

	return nil
}

func (r *CockroachRepository) SeedTickets(ctx context.Context) error {
	genesisID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	
	// 1. Create Genesis Ticket
	genesisTicket := &domain.FabricTicket{
		ID:             genesisID,
		LayerID:        0,
		CreatorAgentID: "SWARM-ARCHITECT",
		Status:         "COMPLETED",
		Resolution:     "Forensic Ticketing Fabric Operational",
	}
	genesisContent := &domain.FabricContent{
		TicketID:   genesisID,
		IntentBlob: map[string]interface{}{"title": "Genesis Root: PQR Info Swarm v1.0", "type": "GENESIS"},
		RawContent: []byte("The world's first autonomous, forensic-grade ticketing fabric for hyperdevelopment swarms."),
	}
	
	// Check if it exists first
	var exists bool
	r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM tickets WHERE ticket_id = $1)", genesisID).Scan(&exists)
	if exists {
		return nil // Already seeded
	}

	if err := r.Create(ctx, genesisTicket, genesisContent); err != nil {
		return fmt.Errorf("failed to create genesis ticket: %v", err)
	}

	// 2. Create Swarm Init Ticket
	initID := uuid.New()
	initTicket := &domain.FabricTicket{
		ID:             initID,
		LayerID:        1,
		CreatorAgentID: "SWARM-ARCHITECT",
		Status:         "COMPLETED",
		Resolution:     "64-Agent Game Theory Swarm Initialized",
	}
	initContent := &domain.FabricContent{
		TicketID:   initID,
		IntentBlob: map[string]interface{}{"title": "Swarm Initialization", "type": "INFRASTRUCTURE"},
		RawContent: []byte("Initializing autonomous nodes and establishing P2P consensus mesh."),
	}
	
	if err := r.Create(ctx, initTicket, initContent); err != nil {
		return fmt.Errorf("failed to create init ticket: %v", err)
	}
	r.Link(ctx, genesisID, initID, domain.RelEvolution)

	// 3. Create some active tickets
	activeTickets := []struct {
		Title   string
		Agent   string
		Layer   int
		Content string
	}{
		{"Refactor Sovereign Node Infrastructure", "council-001", 2, "Unwinding sovereign node spaghetti and implementing Repository pattern."},
		{"Optimize Agent Memory Retrieval", "council-002", 2, "Implementing RAE pattern to bypass redundancy in Iteration 4-6."},
		{"Deploy Local MCP Gateway", "monitor-001", 3, "Upgrading FastAPI service to Anthropic Model Context Protocol standard."},
	}

	for _, at := range activeTickets {
		tid := uuid.New()
		t := &domain.FabricTicket{
			ID:             tid,
			LayerID:        at.Layer,
			CreatorAgentID: at.Agent,
			Status:         "IN_PROGRESS",
		}
		c := &domain.FabricContent{
			TicketID:   tid,
			IntentBlob: map[string]interface{}{"title": at.Title},
			RawContent: []byte(at.Content),
		}
		if err := r.Create(ctx, t, c); err != nil {
			return fmt.Errorf("failed to create active ticket %s: %v", at.Title, err)
		}
		r.Link(ctx, initID, tid, domain.RelEvolution)
	}

	return nil
}
