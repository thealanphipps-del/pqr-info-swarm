package domain

import (
	"context"

	"github.com/google/uuid"
)

// TicketRepository defines the interface for ticket persistence
type TicketRepository interface {
	Create(ctx context.Context, ticket *FabricTicket, content *FabricContent) error
	GetByID(ctx context.Context, id uuid.UUID) (*FabricTicket, *FabricContent, error)
	Update(ctx context.Context, id uuid.UUID, status string, title string) error
	Link(ctx context.Context, parentID, childID uuid.UUID, relType RelationshipType) error
	ListRecent(ctx context.Context, limit int) ([]FabricTicket, error)
	GetAuditTrail(ctx context.Context, id uuid.UUID) ([]AuditEntry, error)
	AddAudit(ctx context.Context, entry AuditEntry) error
	AddFailedAttempt(ctx context.Context, id uuid.UUID, attempt string) error
	FindSimilarResolutions(ctx context.Context, vector []float64, limit int) ([]FabricTicket, error)
	Search(ctx context.Context, criteria map[string]interface{}) ([]FabricTicket, error)
}

// AgentMemoryRepository defines the interface for agent-specific memory storage
type AgentMemoryRepository interface {
	Store(ctx context.Context, agentID string, ticketID uuid.UUID, memType string, data map[string]interface{}, score float64) error
	Get(ctx context.Context, agentID string, ticketID uuid.UUID, memType string) (map[string]interface{}, error)
	GetContext(ctx context.Context, agentID string, limit int) ([]FabricTicket, error)
}
