package rtgo

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sovereign-node-go/pkg/llm"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)


type Manager struct {
	store *Store
	llm   *llm.LocalGateway
}

func NewManager(connStr string, localLLM *llm.LocalGateway) (*Manager, error) {
	if connStr == "" {
		connStr = "postgresql://root@localhost:26257/antigravity?sslmode=disable"
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cockroachdb: %v", err)
	}
	return &Manager{
		store: NewStore(db),
		llm:   localLLM,
	}, nil
}

// CreateFabricTicketV71 creates a new ticket with state vectors and hashes.
func (m *Manager) CreateFabricTicketV71(ctx context.Context, layer int, agentID string, content FabricContent) (uuid.UUID, error) {
	summaryHash := sha256.Sum256([]byte(fmt.Sprintf("%v", content.IntentBlob)))
	content.SummaryHash = hex.EncodeToString(summaryHash[:])

	payloadHash := sha256.Sum256(content.RawContent)
	content.PayloadHash = hex.EncodeToString(payloadHash[:])

	ticketID := uuid.New()

	if len(content.StateVector) == 0 {
		content.StateVector = m.GenerateStateVector(ctx, content.RawContent)
	}

	ticket := FabricTicket{
		ID:             ticketID,
		LayerID:        layer,
		CreatorAgentID: agentID,
		Status:         StatusPending,
		Priority:       PriorityNormal,
		Queue:          "General",
	}

	err := m.store.CreateTicket(ctx, ticket, content)
	if err != nil {
		return uuid.Nil, err
	}

	return ticketID, nil
}

func (m *Manager) LinkTicketsV71(ctx context.Context, parentID, childID uuid.UUID, relType RelationshipType) error {
	return m.store.LinkTickets(ctx, parentID, childID, relType)
}

func (m *Manager) UpdateTicket(ctx context.Context, idStr string, title string, status string, priority int, assignedTo string) error {
	if status != "" {
		if err := m.store.UpdateTicketStatus(ctx, idStr, status); err != nil {
			return err
		}
	}
	if priority > 0 {
		if err := m.store.UpdateTicketPriority(ctx, idStr, priority); err != nil {
			return err
		}
	}
	if assignedTo != "" {
		if err := m.store.UpdateTicketAssignment(ctx, idStr, assignedTo); err != nil {
			return err
		}
	}
	if title != "" {
		if err := m.store.UpdateTicketTitle(ctx, idStr, title); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) UpdateTicketResponse(ctx context.Context, idStr string, response string) error {
	return m.store.AddTicketResponse(ctx, idStr, response)
}

func (m *Manager) MergeTickets(ctx context.Context, targetID, sourceID uuid.UUID) error {
	if err := m.LinkTicketsV71(ctx, targetID, sourceID, RelMerge); err != nil {
		return err
	}
	return m.store.UpdateTicketStatus(ctx, sourceID.String(), StatusDeleted)
}

func (m *Manager) GenerateStateVector(ctx context.Context, content []byte) []float64 {
	if m.llm != nil {
		emb, err := m.llm.GenerateEmbedding(ctx, string(content))
		if err == nil {
			out := make([]float64, len(emb))
			for i, v := range emb {
				out[i] = float64(v)
			}
			return out
		}
	}

	// Fallback to hash-based dummy vector
	h := sha256.Sum256(content)
	vec := make([]float64, 8)
	for i := 0; i < 8; i++ {
		vec[i] = float64(h[i]) / 255.0
	}
	return vec
}

func (m *Manager) SemanticSearch(ctx context.Context, query string, limit int) ([]TicketSummary, error) {
	// For now, return a limited list from the store.
	// In the future, this should perform actual vector comparison.
	return m.store.ListTickets(ctx, limit)
}

func (m *Manager) FetchTickets(ctx context.Context) ([]TicketSummary, error) {
	return m.store.ListTickets(ctx, 100)
}

func (m *Manager) ExportForensicPacket(ctx context.Context, id uuid.UUID) (string, error) {
	summary, err := m.store.GetTicketSummary(ctx, id)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# FORENSIC CONTEXT PACKET: %s\n\n", summary.ID))
	sb.WriteString(fmt.Sprintf("- **Status**: %s\n", summary.Status))
	sb.WriteString(fmt.Sprintf("- **Priority**: %d\n", summary.Priority))
	sb.WriteString(fmt.Sprintf("- **Queue**: %s\n", summary.Queue))
	sb.WriteString(fmt.Sprintf("- **Assigned To**: %s\n", summary.AssignedTo))
	sb.WriteString(fmt.Sprintf("- **Created At**: %s\n\n", summary.CreatedAt.Format(time.RFC3339)))

	sb.WriteString("## RAW CONTENT\n")
	sb.WriteString("```\n")
	sb.WriteString(summary.RawContent)
	sb.WriteString("\n```\n\n")

	sb.WriteString("## AGENT RESPONSES\n")
	if responses, ok := summary.Intent["responses"].([]interface{}); ok {
		for i, resp := range responses {
			sb.WriteString(fmt.Sprintf("### RESPONSE %d\n", i+1))
			sb.WriteString(fmt.Sprintln(resp))
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString("*No agent responses recorded.*\n")
	}

	return sb.String(), nil
}

func (m *Manager) FetchTicketTree(ctx context.Context) ([]*TreeNode, error) {
	nodeMap, err := m.store.GetAllTicketNodes(ctx)
	if err != nil {
		return nil, err
	}

	rels, err := m.store.GetRelationships(ctx)
	if err != nil {
		return nil, err
	}

	childIDs := map[string]bool{}
	for _, rel := range rels {
		parent, ok1 := nodeMap[rel.Parent]
		child, ok2 := nodeMap[rel.Child]
		if ok1 && ok2 {
			child.RelType = rel.Type
			parent.Children = append(parent.Children, child)
			childIDs[rel.Child] = true
		}
	}

	var roots []*TreeNode
	for id, node := range nodeMap {
		if !childIDs[id] {
			roots = append(roots, node)
		}
	}
	return roots, nil
}

func (m *Manager) GetAgentMemoryContext(ctx context.Context, agentID string) (*AgentMemoryContext, error) {
	return m.store.GetAgentMemoryContext(ctx, agentID)
}


func (m *Manager) AddNPUSharingReward(ctx context.Context, agentID string, ticketID uuid.UUID, cycles int, reward float64) error {
	return m.store.AddNPUSharingReward(ctx, agentID, ticketID, cycles, reward)
}
