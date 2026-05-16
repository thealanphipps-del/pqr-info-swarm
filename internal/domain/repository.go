package domain

import (
	"context"

	"github.com/google/uuid"
)

// TicketRepository defines the interface for ticket persistence
type TicketRepository interface {
	CreateTicket(ctx context.Context, ticket *FabricTicket, content *FabricContent) error
	GetByID(ctx context.Context, id uuid.UUID) (*FabricTicket, *FabricContent, error)
	Update(ctx context.Context, id uuid.UUID, status string, title string) error
	Link(ctx context.Context, parentID, childID uuid.UUID, relType RelationshipType) error
	ListRecent(ctx context.Context, limit int) ([]FabricTicket, error)
	GetAuditTrail(ctx context.Context, id uuid.UUID) ([]AuditEntry, error)
	AddAudit(ctx context.Context, entry AuditEntry) error
	AddFailedAttempt(ctx context.Context, id uuid.UUID, attempt string) error
	FindSimilarResolutions(ctx context.Context, vector []float64, limit int) ([]FabricTicket, error)
	Search(ctx context.Context, criteria map[string]interface{}) ([]FabricTicket, error)
	IncrementMetric(ctx context.Context, key string, amount float64) error
	GetMetric(ctx context.Context, key string) (float64, float64, error)
}

// AgentMemoryRepository defines the interface for agent-specific memory storage
type AgentMemoryRepository interface {
	Store(ctx context.Context, agentID string, ticketID uuid.UUID, memType string, data map[string]interface{}, score float64) error
	Get(ctx context.Context, agentID string, ticketID uuid.UUID, memType string) (map[string]interface{}, error)
	GetContext(ctx context.Context, agentID string, limit int) ([]FabricTicket, error)
}

// UserRepository defines the interface for user persistence and authentication
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}
