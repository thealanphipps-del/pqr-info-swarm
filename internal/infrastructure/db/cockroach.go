package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/thealanphipps-del/pqr/internal/domain"
	"github.com/lib/pq"
)

type CockroachRepository struct {
	db *sql.DB
}

func NewCockroachRepository(connStr string) (*CockroachRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cockroachdb: %v", err)
	}
	return &CockroachRepository{db: db}, nil
}

func (r *CockroachRepository) Create(ctx context.Context, t *domain.FabricTicket, c *domain.FabricContent) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO tickets (ticket_id, layer_id, creator_agent_id, status, iteration, escalation_level)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, t.ID, t.LayerID, t.CreatorAgentID, t.Status, t.Iteration, t.EscalationLevel)
	if err != nil {
		return err
	}

	intentJSON, _ := json.Marshal(c.IntentBlob)
	_, err = tx.ExecContext(ctx, `
		INSERT INTO ticket_content (ticket_id, intent_blob, consensus_score, raw_content, summary_hash, payload_hash, state_vector)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, t.ID, intentJSON, c.ConsensusScore, c.RawContent, c.SummaryHash, c.PayloadHash, pq.Array(c.StateVector))
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *CockroachRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.FabricTicket, *domain.FabricContent, error) {
	var t domain.FabricTicket
	var c domain.FabricContent
	var intentJSON []byte
	var failedAttemptsJSON []byte

	err := r.db.QueryRowContext(ctx, `
		SELECT t.ticket_id, t.layer_id, t.creator_agent_id, t.status, t.iteration, t.escalation_level, t.resolution, t.resolved_by, t.created_at,
		       c.intent_blob, c.raw_content, c.consensus_score, c.summary_hash, c.payload_hash, c.failed_attempts
		FROM tickets t
		LEFT JOIN ticket_content c ON t.ticket_id = c.ticket_id
		WHERE t.ticket_id = $1
	`, id).Scan(&t.ID, &t.LayerID, &t.CreatorAgentID, &t.Status, &t.Iteration, &t.EscalationLevel, &t.Resolution, &t.ResolvedBy, &t.CreatedAt, 
		&intentJSON, &c.RawContent, &c.ConsensusScore, &c.SummaryHash, &c.PayloadHash, &failedAttemptsJSON)

	if err != nil {
		return nil, nil, err
	}

	if intentJSON != nil {
		json.Unmarshal(intentJSON, &c.IntentBlob)
	}
	if failedAttemptsJSON != nil {
		json.Unmarshal(failedAttemptsJSON, &c.FailedAttempts)
	}
	c.TicketID = t.ID

	return &t, &c, nil
}

func (r *CockroachRepository) Update(ctx context.Context, id uuid.UUID, status string, title string) error {
	// We'll expand this to handle Resolution if status is COMPLETED
	// But for now let's just make it generic
	return r.UpdateExtended(ctx, id, status, title, "", "")
}

func (r *CockroachRepository) UpdateExtended(ctx context.Context, id uuid.UUID, status string, title string, resolution string, resolvedBy string) error {
	if status != "" {
		_, err := r.db.ExecContext(ctx, `UPDATE tickets SET status = $1 WHERE ticket_id = $2`, status, id)
		if err != nil {
			return err
		}
	}

	if resolution != "" {
		_, err := r.db.ExecContext(ctx, `UPDATE tickets SET resolution = $1, resolved_by = $2 WHERE ticket_id = $3`, resolution, resolvedBy, id)
		if err != nil {
			return err
		}
	}

	if title != "" {
		_, err := r.db.ExecContext(ctx, `
			UPDATE ticket_content 
			SET intent_blob = intent_blob || jsonb_build_object('title', $1::STRING)
			WHERE ticket_id = $2
		`, title, id)
		if err != nil {
			return err
		}
	}

	if resolvedBy != "" && resolution == "" { // Using resolvedBy for reassignment if resolution is empty
		_, err := r.db.ExecContext(ctx, `UPDATE tickets SET creator_agent_id = $1 WHERE ticket_id = $2`, resolvedBy, id)
		if err != nil {
			return err
		}
	}

	return nil
}


func (r *CockroachRepository) AddFailedAttempt(ctx context.Context, id uuid.UUID, attempt string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE ticket_content 
		SET failed_attempts = failed_attempts || jsonb_build_array($1::STRING)
		WHERE ticket_id = $2
	`, attempt, id)
	return err
}

func (r *CockroachRepository) FindSimilarResolutions(ctx context.Context, vector []float64, limit int) ([]domain.FabricTicket, error) {
	// This would use vector similarity in a real CockroachDB instance with vector support
	// For now we'll just return the last successful resolutions
	rows, err := r.db.QueryContext(ctx, `
		SELECT t.ticket_id, t.resolution, t.resolved_by
		FROM tickets t
		WHERE t.status = 'COMPLETED' AND t.resolution IS NOT NULL
		ORDER BY t.created_at DESC LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.FabricTicket
	for rows.Next() {
		var t domain.FabricTicket
		if err := rows.Scan(&t.ID, &t.Resolution, &t.ResolvedBy); err != nil {
			return nil, err
		}
		results = append(results, t)
	}
	return results, nil
}

func (r *CockroachRepository) Link(ctx context.Context, parentID, childID uuid.UUID, relType domain.RelationshipType) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO ticket_relationships (parent_id, child_id, relationship_type)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`, parentID, childID, string(relType))
	return err
}

func (r *CockroachRepository) ListRecent(ctx context.Context, limit int) ([]domain.FabricTicket, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT ticket_id, layer_id, creator_agent_id, status, created_at 
		FROM tickets ORDER BY created_at DESC LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []domain.FabricTicket
	for rows.Next() {
		var t domain.FabricTicket
		if err := rows.Scan(&t.ID, &t.LayerID, &t.CreatorAgentID, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}

func (r *CockroachRepository) Store(ctx context.Context, agentID string, ticketID uuid.UUID, memType string, data map[string]interface{}, score float64) error {
	memJSON, _ := json.Marshal(data)
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO agent_memory (agent_id, ticket_id, memory_type, memory_data, relevance_score)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (agent_id, ticket_id, memory_type) DO UPDATE SET
			memory_data = $4,
			relevance_score = $5,
			accessed_at = CURRENT_TIMESTAMP
	`, agentID, ticketID, memType, memJSON, score)
	return err
}

func (r *CockroachRepository) Get(ctx context.Context, agentID string, ticketID uuid.UUID, memType string) (map[string]interface{}, error) {
	var memJSON []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT memory_data FROM agent_memory
		WHERE agent_id = $1 AND ticket_id = $2 AND memory_type = $3
	`, agentID, ticketID, memType).Scan(&memJSON)
	
	if err != nil {
		return nil, err
	}
	
	var data map[string]interface{}
	json.Unmarshal(memJSON, &data)
	return data, nil
}

func (r *CockroachRepository) GetContext(ctx context.Context, agentID string, limit int) ([]domain.FabricTicket, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT DISTINCT t.ticket_id, t.layer_id, t.creator_agent_id, t.status, t.created_at
		FROM tickets t
		INNER JOIN agent_memory am ON t.ticket_id = am.ticket_id
		WHERE am.agent_id = $1
		ORDER BY am.relevance_score DESC, t.created_at DESC
		LIMIT $2
	`, agentID, limit)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tickets []domain.FabricTicket
	for rows.Next() {
		var t domain.FabricTicket
		if err := rows.Scan(&t.ID, &t.LayerID, &t.CreatorAgentID, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}

func (r *CockroachRepository) GetAuditTrail(ctx context.Context, id uuid.UUID) ([]domain.AuditEntry, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, agent_id, action, old_value, new_value, created_at
		FROM ticket_audit
		WHERE ticket_id = $1
		ORDER BY created_at DESC
	`, id)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var trail []domain.AuditEntry
	for rows.Next() {
		var entry domain.AuditEntry
		var oldVal, newVal sql.NullString
		
		if err := rows.Scan(&entry.ID, &entry.AgentID, &entry.Action, &oldVal, &newVal, &entry.CreatedAt); err != nil {
			return nil, err
		}
		
		if oldVal.Valid {
			json.Unmarshal([]byte(oldVal.String), &entry.OldValue)
		}
		if newVal.Valid {
			json.Unmarshal([]byte(newVal.String), &entry.NewValue)
		}
		entry.TicketID = id
		trail = append(trail, entry)
	}
	return trail, nil
}

func (r *CockroachRepository) AddAudit(ctx context.Context, entry domain.AuditEntry) error {
	oldJSON, _ := json.Marshal(entry.OldValue)
	newJSON, _ := json.Marshal(entry.NewValue)
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO ticket_audit (ticket_id, agent_id, action, old_value, new_value)
		VALUES ($1, $2, $3, $4, $5)
	`, entry.TicketID, entry.AgentID, entry.Action, oldJSON, newJSON)
	return err
}

func (r *CockroachRepository) InitSchema(ctx context.Context) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS tickets (
			ticket_id UUID PRIMARY KEY,
			layer_id INT NOT NULL,
			creator_agent_id STRING NOT NULL,
			status STRING NOT NULL DEFAULT 'PENDING',
			iteration INT NOT NULL DEFAULT 0,
			escalation_level INT NOT NULL DEFAULT 0,
			resolution STRING,
			resolved_by STRING,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_status (status),
			INDEX idx_creator (creator_agent_id),
			INDEX idx_layer (layer_id)
		)`,
		`CREATE TABLE IF NOT EXISTS agents (
			agent_id STRING PRIMARY KEY,
			role STRING NOT NULL,
			level INT NOT NULL DEFAULT 0,
			metadata JSONB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS ticket_content (
			ticket_id UUID PRIMARY KEY REFERENCES tickets(ticket_id) ON DELETE CASCADE,
			intent_blob JSONB,
			consensus_score DECIMAL,
			raw_content BYTES,
			summary_hash STRING,
			payload_hash STRING,
			state_vector FLOAT8[],
			failed_attempts JSONB DEFAULT '[]',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS ticket_relationships (
			parent_id UUID NOT NULL REFERENCES tickets(ticket_id) ON DELETE CASCADE,
			child_id UUID NOT NULL REFERENCES tickets(ticket_id) ON DELETE CASCADE,
			relationship_type STRING NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (parent_id, child_id, relationship_type),
			INDEX idx_parent (parent_id),
			INDEX idx_child (child_id)
		)`,
		`CREATE TABLE IF NOT EXISTS agent_memory (
			agent_id STRING NOT NULL,
			ticket_id UUID NOT NULL REFERENCES tickets(ticket_id) ON DELETE CASCADE,
			memory_type STRING NOT NULL,
			memory_data JSONB,
			relevance_score DECIMAL,
			accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (agent_id, ticket_id, memory_type),
			INDEX idx_agent (agent_id),
			INDEX idx_relevance (agent_id, relevance_score DESC)
		)`,
		`CREATE TABLE IF NOT EXISTS ticket_audit (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			ticket_id UUID NOT NULL REFERENCES tickets(ticket_id) ON DELETE CASCADE,
			agent_id STRING NOT NULL,
			action STRING NOT NULL,
			old_value JSONB,
			new_value JSONB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_ticket (ticket_id),
			INDEX idx_agent (agent_id),
			INDEX idx_created (created_at DESC)
		)`,
	}

	for _, table := range tables {
		if _, err := r.db.ExecContext(ctx, table); err != nil {
			return err
		}
	}
	return nil
}

func (r *CockroachRepository) Search(ctx context.Context, criteria map[string]interface{}) ([]domain.FabricTicket, error) {
	query := `SELECT ticket_id, layer_id, creator_agent_id, status, created_at FROM tickets WHERE 1=1`
	var args []interface{}
	i := 1
	for k, v := range criteria {
		// Basic sanitization: check against allowed column names
		allowedColumns := map[string]bool{"status": true, "layer_id": true, "creator_agent_id": true}
		columnName := k
		if k == "layer" { columnName = "layer_id" } // Map friendly name
		
		if allowedColumns[columnName] {
			query += fmt.Sprintf(" AND %s = $%d", columnName, i)
			args = append(args, v)
			i++
		}
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []domain.FabricTicket
	for rows.Next() {
		var t domain.FabricTicket
		if err := rows.Scan(&t.ID, &t.LayerID, &t.CreatorAgentID, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}

