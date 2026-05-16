package domain

import (
	"time"

	"github.com/google/uuid"
)

// RelationshipType defines how two tickets are connected
type RelationshipType string

const (
	RelEvolution   RelationshipType = "EVOLUTION"
	RelConsequence RelationshipType = "CONSEQUENCE"
	RelContext     RelationshipType = "CONTEXT"
	RelGenesis     RelationshipType = "GENESIS"
)

// FabricTicket represents the metadata and state of a single memory unit
type FabricTicket struct {
	ID              uuid.UUID `json:"id"`
	LayerID         int       `json:"layer_id"`
	CreatorAgentID  string    `json:"creator_agent_id"`
	Status          string    `json:"status"`
	Iteration       int       `json:"iteration"`
	EscalationLevel int       `json:"escalation_level"`
	Resolution      string    `json:"resolution,omitempty"`
	ResolvedBy      string    `json:"resolved_by,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// FabricContent holds the actual payload and forensic hashes of a ticket
type FabricContent struct {
	TicketID       uuid.UUID              `json:"ticket_id"`
	IntentBlob     map[string]interface{} `json:"intent_blob"`
	StateVector    []float64              `json:"state_vector"`
	ConsensusScore float64                `json:"consensus_score"`
	RawContent     []byte                 `json:"raw_content"`
	SummaryHash    string                 `json:"summary_hash"`
	PayloadHash    string                 `json:"payload_hash"`
	FailedAttempts []string               `json:"failed_attempts,omitempty"`
}

// Agent represents an autonomous entity in the swarm
type Agent struct {
	ID       string `json:"id"`
	Role     string `json:"role"`
	Level    int    `json:"level"`
	Metadata map[string]interface{} `json:"metadata"`
}

// AuditEntry records a forensic trail of changes to a ticket
type AuditEntry struct {
	ID        uuid.UUID              `json:"id"`
	TicketID  uuid.UUID              `json:"ticket_id"`
	AgentID   string                 `json:"agent_id"`
	Action    string                 `json:"action"`
	OldValue  map[string]interface{} `json:"old_value,omitempty"`
	NewValue  map[string]interface{} `json:"new_value,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

// User represents a human or system user for authentication
type User struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}
