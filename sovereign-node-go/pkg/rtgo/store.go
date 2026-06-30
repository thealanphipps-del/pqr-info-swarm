package rtgo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// Store handles all direct interactions with CockroachDB.
type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// CreateTicket inserts the base ticket record and its content.
func (s *Store) CreateTicket(ctx context.Context, t FabricTicket, c FabricContent) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO tickets (ticket_id, layer_id, creator_agent_id, status, priority, queue, assigned_to, is_sticky, referrer_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, t.ID, t.LayerID, t.CreatorAgentID, t.Status, t.Priority, t.Queue, t.AssignedTo, t.IsSticky, t.ReferrerCode)
	if err != nil {
		return fmt.Errorf("failed to insert ticket: %v", err)
	}

	intentJSON, _ := json.Marshal(c.IntentBlob)
	_, err = tx.ExecContext(ctx, `
		INSERT INTO ticket_content (ticket_id, intent_blob, consensus_score, raw_content, summary_hash, payload_hash, state_vector)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, t.ID, intentJSON, c.ConsensusScore, c.RawContent, c.SummaryHash, c.PayloadHash, c.StateVector)
	if err != nil {
		return fmt.Errorf("failed to insert ticket content: %v", err)
	}

	return tx.Commit()
}

// LinkTickets establishes a relationship between two tickets.
func (s *Store) LinkTickets(ctx context.Context, parentID, childID uuid.UUID, relType RelationshipType) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO ticket_relationships (parent_id, child_id, relationship_type)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`, parentID, childID, string(relType))
	return err
}

// UpdateTicketStatus updates the status of a ticket.
func (s *Store) UpdateTicketStatus(ctx context.Context, id string, status string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE tickets SET status = $1 WHERE ticket_id = $2`, status, id)
	return err
}

// UpdateTicketPriority updates the priority of a ticket.
func (s *Store) UpdateTicketPriority(ctx context.Context, id string, priority int) error {
	_, err := s.db.ExecContext(ctx, `UPDATE tickets SET priority = $1 WHERE ticket_id = $2`, priority, id)
	return err
}

// UpdateTicketAssignment updates the assigned agent for a ticket.
func (s *Store) UpdateTicketAssignment(ctx context.Context, id string, assignedTo string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE tickets SET assigned_to = $1 WHERE ticket_id = $2`, assignedTo, id)
	return err
}

// UpdateTicketTitle updates the title within the intent blob.
func (s *Store) UpdateTicketTitle(ctx context.Context, id string, title string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE ticket_content 
		SET intent_blob = intent_blob || jsonb_build_object('title', $1::STRING)
		WHERE ticket_id = $2
	`, title, id)
	return err
}

// AddTicketResponse appends a response to the intent blob's responses array.
func (s *Store) AddTicketResponse(ctx context.Context, id string, response string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE ticket_content 
		SET intent_blob = jsonb_set(
			COALESCE(intent_blob, '{}'::jsonb), 
			'{responses}', 
			(COALESCE(intent_blob->'responses', '[]'::jsonb) || jsonb_build_array($1::STRING)), 
			true
		)
		WHERE ticket_id = $2
	`, response, id)
	return err
}

// GetTicketSummary fetches a single ticket summary.
func (s *Store) GetTicketSummary(ctx context.Context, id uuid.UUID) (TicketSummary, error) {
	var ts TicketSummary
	var intentJSON []byte
	var rawContent []byte
	err := s.db.QueryRowContext(ctx, `
		SELECT t.ticket_id, t.layer_id, t.creator_agent_id, t.status, t.priority, t.queue, t.assigned_to, t.is_sticky, t.referrer_code, t.created_at, c.intent_blob, c.raw_content
		FROM tickets t
		LEFT JOIN ticket_content c ON t.ticket_id = c.ticket_id
		WHERE t.ticket_id = $1
	`, id).Scan(&ts.ID, &ts.Layer, &ts.Creator, &ts.Status, &ts.Priority, &ts.Queue, &ts.AssignedTo, &ts.IsSticky, &ts.ReferrerCode, &ts.CreatedAt, &intentJSON, &rawContent)
	
	if err != nil {
		return ts, err
	}
	if intentJSON != nil {
		json.Unmarshal(intentJSON, &ts.Intent)
	}
	if rawContent != nil {
		ts.RawContent = string(rawContent)
	}
	return ts, nil
}

// ListTickets returns a list of the most recent tickets.
func (s *Store) ListTickets(ctx context.Context, limit int) ([]TicketSummary, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT t.ticket_id, t.layer_id, t.creator_agent_id, t.status, t.priority, t.queue, t.assigned_to, t.is_sticky, t.referrer_code, t.created_at, c.intent_blob, c.raw_content
		FROM tickets t
		LEFT JOIN ticket_content c ON t.ticket_id = c.ticket_id
		ORDER BY t.created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TicketSummary
	for rows.Next() {
		var ts TicketSummary
		var intentJSON []byte
		var rawContent []byte
		if err := rows.Scan(&ts.ID, &ts.Layer, &ts.Creator, &ts.Status, &ts.Priority, &ts.Queue, &ts.AssignedTo, &ts.IsSticky, &ts.ReferrerCode, &ts.CreatedAt, &intentJSON, &rawContent); err != nil {
			continue
		}
		if rawContent != nil {
			ts.RawContent = string(rawContent)
		}
		if intentJSON != nil {
			json.Unmarshal(intentJSON, &ts.Intent)
			// Helper to set title from intent
			if t, ok := ts.Intent["title"].(string); ok {
				ts.Title = t
			} else if t, ok := ts.Intent["subject"].(string); ok {
				ts.Title = t
			}
		}
		if ts.Title == "" {
			ts.Title = "Ticket " + ts.ID[:8]
		}
		out = append(out, ts)
	}
	return out, nil
}

// GetAllTicketNodes fetches all tickets and their relationships for tree building.
func (s *Store) GetAllTicketNodes(ctx context.Context) (map[string]*TreeNode, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT t.ticket_id, t.layer_id, t.creator_agent_id, t.status, t.created_at, c.intent_blob
		FROM tickets t
		LEFT JOIN ticket_content c ON t.ticket_id = c.ticket_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodeMap := map[string]*TreeNode{}
	for rows.Next() {
		var n TreeNode
		var intentJSON []byte
		if err := rows.Scan(&n.ID, &n.Layer, &n.Creator, &n.Status, &n.CreatedAt, &intentJSON); err != nil {
			continue
		}
		if intentJSON != nil {
			json.Unmarshal(intentJSON, &n.Intent)
			if t, ok := n.Intent["title"].(string); ok {
				n.Title = t
			} else if t, ok := n.Intent["subject"].(string); ok {
				n.Title = t
			}
		}
		if n.Title == "" {
			n.Title = "Ticket " + n.ID[:8]
		}
		nodeMap[n.ID] = &n
	}
	return nodeMap, nil
}

// GetRelationships fetches all ticket relationships.
func (s *Store) GetRelationships(ctx context.Context) ([]struct{Parent, Child, Type string}, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT parent_id, child_id, relationship_type FROM ticket_relationships
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []struct{Parent, Child, Type string}
	for rows.Next() {
		var r struct{Parent, Child, Type string}
		if err := rows.Scan(&r.Parent, &r.Child, &r.Type); err != nil {
			continue
		}
		out = append(out, r)
	}
	return out, nil
}

// AgentMemoryContext represents the loaded memory and configuration for an agent.
type AgentMemoryContext struct {
	AgentID      string
	SwarmCluster string
	AgentRole    string
	Memories     []struct {
		TicketID   string
		Layer      int
		Status     string
		Title      string
		RawContent string
	}
	LearningCases []struct {
		CaseID        string
		CaseHash      string
		ParentTitle   string
		ParentProblem string
		ChildSolution string
	}
}

// GetAgentMemoryContext retrieves the agent's metadata and their associated tickets.
func (s *Store) GetAgentMemoryContext(ctx context.Context, agentID string) (*AgentMemoryContext, error) {
	// Query containing upgrade v2 columns
	rows, err := s.db.QueryContext(ctx, `
		SELECT COALESCE(m.swarm_cluster, '') as swarm_cluster, 
		       COALESCE(m.agent_role, '') as agent_role, 
		       t.ticket_id, t.layer_id, t.status, c.intent_blob, c.raw_content
		FROM agent_memory_index m
		JOIN tickets t ON m.ticket_id = t.ticket_id
		LEFT JOIN ticket_content c ON t.ticket_id = c.ticket_id
		WHERE m.agent_id = $1
	`, agentID)
	if err != nil {
		// Fallback query if v2 columns are not present
		rowsFallback, errFallback := s.db.QueryContext(ctx, `
			SELECT '' as swarm_cluster, 
			       '' as agent_role, 
			       t.ticket_id, t.layer_id, t.status, c.intent_blob, c.raw_content
			FROM agent_memory_index m
			JOIN tickets t ON m.ticket_id = t.ticket_id
			LEFT JOIN ticket_content c ON t.ticket_id = c.ticket_id
			WHERE m.agent_id = $1
		`, agentID)
		if errFallback != nil {
			return nil, fmt.Errorf("failed to fetch agent memory: %v (fallback: %v)", err, errFallback)
		}
		rows = rowsFallback
	}
	defer rows.Close()

	amc := &AgentMemoryContext{AgentID: agentID}
	for rows.Next() {
		var cluster, role, ticketID, status string
		var layer int
		var intentJSON, rawContent []byte
		if err := rows.Scan(&cluster, &role, &ticketID, &layer, &status, &intentJSON, &rawContent); err != nil {
			continue
		}
		amc.SwarmCluster = cluster
		amc.AgentRole = role

		title := ""
		if intentJSON != nil {
			var intent map[string]interface{}
			if err := json.Unmarshal(intentJSON, &intent); err == nil {
				if t, ok := intent["title"].(string); ok {
					title = t
				} else if t, ok := intent["subject"].(string); ok {
					title = t
				}
			}
		}
		if title == "" {
			title = "Ticket " + ticketID[:8]
		}

		amc.Memories = append(amc.Memories, struct {
			TicketID   string
			Layer      int
			Status     string
			Title      string
			RawContent string
		}{
			TicketID:   ticketID,
			Layer:      layer,
			Status:     status,
			Title:      title,
			RawContent: string(rawContent),
		})
	}

	// Fetch historical self-healing resolution patterns for active learning with deterministic sorting
	rowsLC, errLC := s.db.QueryContext(ctx, `
		SELECT r.parent_id::STRING as case_id,
		       hash(r.parent_id)::STRING as case_hash,
		       COALESCE(pc.intent_blob->>'title', 'Parent Case') as parent_title,
		       COALESCE(CAST(pc.raw_content AS STRING), '') as parent_problem,
		       COALESCE(CAST(cc.raw_content AS STRING), '') as child_solution
		FROM ticket_relationships r
		JOIN tickets p ON r.parent_id = p.ticket_id
		LEFT JOIN ticket_content pc ON p.ticket_id = pc.ticket_id
		JOIN tickets c ON r.child_id = c.ticket_id
		LEFT JOIN ticket_content cc ON c.ticket_id = cc.ticket_id
		WHERE c.status = 'PROMOTED' AND r.relationship_type = 'CONSEQUENCE'
		ORDER BY hash(r.parent_id) ASC
		LIMIT 10
	`)
	if errLC == nil {
		defer rowsLC.Close()
		for rowsLC.Next() {
			var caseItem struct {
				CaseID        string
				CaseHash      string
				ParentTitle   string
				ParentProblem string
				ChildSolution string
			}
			if errScan := rowsLC.Scan(&caseItem.CaseID, &caseItem.CaseHash, &caseItem.ParentTitle, &caseItem.ParentProblem, &caseItem.ChildSolution); errScan == nil {
				amc.LearningCases = append(amc.LearningCases, caseItem)
			}
		}
	}

	return amc, nil
}


// AddNPUSharingReward logs a reward event to the npu_sharing_rewards table.
func (s *Store) AddNPUSharingReward(ctx context.Context, agentID string, ticketID uuid.UUID, cycles int, reward float64) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO npu_sharing_rewards (agent_id, ticket_id, npu_compute_cycles, cobalt_chrome_rewarded)
		VALUES ($1, $2, $3, $4)
	`, agentID, ticketID, cycles, reward)
	return err
}
