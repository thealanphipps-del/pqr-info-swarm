package rtgo

import (
	"time"

	"github.com/google/uuid"
)

// RelationshipType defines the nature of the link between two tickets.
type RelationshipType string

const (
	RelEvolution   RelationshipType = "EVOLUTION"   // Direct successor
	RelConsequence RelationshipType = "CONSEQUENCE" // Indirect result
	RelContext     RelationshipType = "CONTEXT"     // Background information
	RelGenesis     RelationshipType = "GENESIS"     // Root of a new chain
	RelMerge       RelationshipType = "MERGE"       // Combination of two tickets
	RelSwap        RelationshipType = "SWAP"        // Replacement/Referrer swap
)

// Ticket Status constants
const (
	StatusNew        = "NEW"
	StatusOpen       = "OPEN"
	StatusPending    = "PENDING"
	StatusHITLRequired = "HITL_REQUIRED"
	StatusInProgress = "IN_PROGRESS"
	StatusStalled    = "STALLED"
	StatusStuck      = "STUCK"
	StatusResolved   = "RESOLVED"
	StatusPromoted   = "PROMOTED"
	StatusDeleted    = "DELETED"
)

// Ticket Priority constants
const (
	PriorityLow      = 10
	PriorityNormal   = 20
	PriorityHigh     = 30
	PriorityCritical = 40
)

// FabricTicket represents the base record of a ticket in the Sovereign Node.
type FabricTicket struct {
	ID             uuid.UUID
	LayerID        int
	CreatorAgentID string
	Status         string
	Priority       int
	Queue          string
	AssignedTo     string
	IsSticky       bool
	ReferrerCode   string
	CreatedAt      time.Time
}

// FabricContent holds the rich data associated with a ticket, including intent and embeddings.
type FabricContent struct {
	TicketID       uuid.UUID
	IntentBlob     map[string]interface{}
	StateVector    []float64
	ConsensusScore float64
	RawContent     []byte
	SummaryHash    string
	PayloadHash    string
}

// TicketSummary is the wire-format for the SSE stream and REST list endpoints.
type TicketSummary struct {
	ID           string                 `json:"id"`
	Layer        int                    `json:"layer"`
	Creator      string                 `json:"creator"`
	Status       string                 `json:"status"`
	Priority     int                    `json:"priority"`
	Queue        string                 `json:"queue"`
	AssignedTo   string                 `json:"assigned_to"`
	IsSticky     bool                   `json:"is_sticky"`
	ReferrerCode string                 `json:"referrer_code"`
	Title        string                 `json:"title"`
	CreatedAt    time.Time              `json:"created_at"`
	Intent       map[string]interface{} `json:"intent,omitempty"`
	RawContent   string                 `json:"raw_content,omitempty"`
}

// TreeNode represents a single node in the GEDCOM-style ticket genealogy.
type TreeNode struct {
	ID           string                 `json:"id"`
	Layer        int                    `json:"layer"`
	Creator      string                 `json:"creator"`
	Status       string                 `json:"status"`
	Priority     int                    `json:"priority"`
	Queue        string                 `json:"queue"`
	AssignedTo   string                 `json:"assigned_to"`
	IsSticky     bool                   `json:"is_sticky"`
	ReferrerCode string                 `json:"referrer_code"`
	Title        string                 `json:"title"`
	CreatedAt    time.Time              `json:"created_at"`
	RelType      string                 `json:"rel_type,omitempty"`
	Intent       map[string]interface{} `json:"intent,omitempty"`
	Children     []*TreeNode            `json:"children,omitempty"`
}
