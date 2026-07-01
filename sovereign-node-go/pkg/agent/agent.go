package agent

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"sovereign-node-go/pkg/tickets"

	"github.com/google/uuid"
)

type AgentIdentity struct {
	AgentID   string
	Shortcode string
}

type AgentManager struct {
	db        *sql.DB
	ticketMgr *tickets.Manager
}

func NewAgentManager(db *sql.DB, ticketMgr *tickets.Manager) *AgentManager {
	return &AgentManager{
		db:        db,
		ticketMgr: ticketMgr,
	}
}

func (m *AgentManager) InitSchema(ctx context.Context) error {
	_, err := m.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS agent_identity (
			agent_id STRING PRIMARY KEY,
			shortcode STRING UNIQUE
		);
	`)
	return err
}

func generateShortcode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

func (m *AgentManager) BuildAgents(ctx context.Context) ([]AgentIdentity, error) {
	// 1. Fetch all distinct agents
	rows, err := m.db.QueryContext(ctx, `SELECT DISTINCT agent_id FROM agent_memory_index`)
	if err != nil {
		return nil, fmt.Errorf("failed to query agents: %v", err)
	}
	defer rows.Close()

	var agents []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		agents = append(agents, id)
	}

	var identities []AgentIdentity
	genesisID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	// 2. Assign shortcodes and track via tickets
	for _, agentID := range agents {
		// Check if already assigned
		var existingCode string
		err := m.db.QueryRowContext(ctx, `SELECT shortcode FROM agent_identity WHERE agent_id = $1`, agentID).Scan(&existingCode)
		if err == nil {
			identities = append(identities, AgentIdentity{AgentID: agentID, Shortcode: existingCode})
			continue // Already assigned
		} else if err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed checking identity: %v", err)
		}

		shortcode := generateShortcode()

		// Store identity
		_, err = m.db.ExecContext(ctx, `INSERT INTO agent_identity (agent_id, shortcode) VALUES ($1, $2)`, agentID, shortcode)
		if err != nil {
			log.Printf("Failed to insert shortcode for %s: %v", agentID, err)
			continue
		}

		identities = append(identities, AgentIdentity{AgentID: agentID, Shortcode: shortcode})

		// Track assignment in Ticket
		content := tickets.FabricContent{
			IntentBlob: map[string]interface{}{
				"action":    "ASSIGN_IDENTITY",
				"agent_id":  agentID,
				"shortcode": shortcode,
			},
			ConsensusScore: 1.0,
			RawContent:     []byte(fmt.Sprintf("Assigned 5-alphanumeric shortcode %s to %s", shortcode, agentID)),
		}

		ticketID, err := m.ticketMgr.CreateFabricTicketV71(ctx, 2, agentID, content)
		if err != nil {
			log.Printf("Failed to track ticket for %s: %v", agentID, err)
			continue
		}

		err = m.ticketMgr.LinkTicketsV71(ctx, genesisID, ticketID, tickets.RelEvolution)
		if err != nil {
			log.Printf("Failed to link ticket for %s: %v", agentID, err)
		}
	}

	return identities, nil
}
