package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/thealanphipps-del/pqr/internal/domain"
)

type SwarmService struct {
	repo domain.TicketRepository
	mem  domain.AgentMemoryRepository
}

func NewSwarmService(repo domain.TicketRepository, mem domain.AgentMemoryRepository) *SwarmService {
	return &SwarmService{
		repo: repo,
		mem:  mem,
	}
}

func (s *SwarmService) CreateFabricTicket(ctx context.Context, layer int, agentID string, content domain.FabricContent) (uuid.UUID, error) {
	summaryHash := sha256.Sum256([]byte(fmt.Sprintf("%v", content.IntentBlob)))
	content.SummaryHash = hex.EncodeToString(summaryHash[:])

	payloadHash := sha256.Sum256(content.RawContent)
	content.PayloadHash = hex.EncodeToString(payloadHash[:])

	ticketID := uuid.New()
	content.TicketID = ticketID

	if len(content.StateVector) == 0 {
		content.StateVector = s.GenerateStateVector(content.RawContent)
	}

	ticket := &domain.FabricTicket{
		ID:             ticketID,
		LayerID:        layer,
		CreatorAgentID: agentID,
		Status:         "PENDING",
	}

	if err := s.repo.Create(ctx, ticket, &content); err != nil {
		return uuid.Nil, err
	}

	return ticketID, nil
}

func (s *SwarmService) LinkTicketsWithAudit(ctx context.Context, parentID, childID uuid.UUID, relType domain.RelationshipType, agentID string) error {
	if err := s.repo.Link(ctx, parentID, childID, relType); err != nil {
		return err
	}

	entry := domain.AuditEntry{
		TicketID: parentID,
		AgentID:  agentID,
		Action:   fmt.Sprintf("LINK_%s", relType),
		NewValue: map[string]interface{}{"child_id": childID.String()},
	}

	return s.repo.AddAudit(ctx, entry)
}

func (s *SwarmService) GenerateStateVector(content []byte) []float64 {
	h := sha256.Sum256(content)
	vec := make([]float64, 8)
	for i := 0; i < 8; i++ {
		vec[i] = float64(h[i]) / 255.0
	}
	return vec
}

func (s *SwarmService) GetRecentTickets(ctx context.Context, limit int) ([]domain.FabricTicket, error) {
	return s.repo.ListRecent(ctx, limit)
}

func (s *SwarmService) UpdateTicket(ctx context.Context, id uuid.UUID, status string, title string) error {
	return s.repo.Update(ctx, id, status, title)
}

func (s *SwarmService) UpdateExtended(ctx context.Context, id uuid.UUID, status, title, resolution, creator string) error {
	// Cast repository to check for extended update support
	if r, ok := s.repo.(interface {
		UpdateExtended(context.Context, uuid.UUID, string, string, string, string) error
	}); ok {
		return r.UpdateExtended(ctx, id, status, title, resolution, creator)
	}
	// Fallback to basic update
	return s.repo.Update(ctx, id, status, title)
}


func (s *SwarmService) GetTicketWithContent(ctx context.Context, id uuid.UUID) (*domain.FabricTicket, *domain.FabricContent, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SwarmService) StoreAgentMemory(ctx context.Context, agentID string, ticketID uuid.UUID, memType string, data map[string]interface{}, score float64) error {
	return s.mem.Store(ctx, agentID, ticketID, memType, data, score)
}

func (s *SwarmService) GetAgentMemory(ctx context.Context, agentID string, ticketID uuid.UUID, memType string) (map[string]interface{}, error) {
	return s.mem.Get(ctx, agentID, ticketID, memType)
}

func (s *SwarmService) GetAgentContext(ctx context.Context, agentID string, limit int) ([]domain.FabricTicket, error) {
	return s.mem.GetContext(ctx, agentID, limit)
}

func (s *SwarmService) GetAuditTrail(ctx context.Context, id uuid.UUID) ([]domain.AuditEntry, error) {
	return s.repo.GetAuditTrail(ctx, id)
}

func (s *SwarmService) InitSchema(ctx context.Context) error {
	// In a real repo-based app, we'd cast or use a specific init method
	if r, ok := s.repo.(interface{ InitSchema(context.Context) error }); ok {
		return r.InitSchema(ctx)
	}
	return nil
}
