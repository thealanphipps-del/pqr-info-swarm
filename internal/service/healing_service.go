package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/thealanphipps-del/pqr/internal/domain"
)

type HealingService struct {
	repo domain.TicketRepository
	svc  *SwarmService
	ai   *AIService
}

func NewHealingService(repo domain.TicketRepository, svc *SwarmService, ai *AIService) *HealingService {
	return &HealingService{
		repo: repo,
		svc:  svc,
		ai:   ai,
	}
}

// StartBackgroundWorker kicks off the autonomous healing loop
func (h *HealingService) StartBackgroundWorker(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	log.Println("⚡ PQR Background Healing Worker: ONLINE")
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				h.processPendingTickets(ctx)
			}
		}
	}()
}

func (h *HealingService) processPendingTickets(ctx context.Context) {
	// Search for tickets in PENDING state that require healing (Layer 5-7)
	tickets, err := h.repo.Search(ctx, map[string]interface{}{"status": "PENDING"})
	if err != nil {
		return
	}
	
	for _, t := range tickets {
		if t.LayerID >= 5 {
			log.Printf("Autonomous Resolution: Found pending Layer %d ticket %s. Initiating healing loop...", t.LayerID, t.ID)
			h.ProcessHealingLoop(ctx, t.ID)
		}
	}
}


// ProcessHealingLoop advances a ticket through the self-healing escalation levels
func (h *HealingService) ProcessHealingLoop(ctx context.Context, ticketID uuid.UUID) error {
	ticket, content, err := h.repo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	if ticket.Status == "COMPLETED" || ticket.Status == "STALLED" {
		return nil
	}

	// Increment iteration
	ticket.Iteration++

	// Determine model and escalation level based on iteration
	var model string
	switch {
	case ticket.Iteration >= 1 && ticket.Iteration <= 3:
		ticket.EscalationLevel = 1
		model = "LM Studio / Gemma-4-e4b"
	case ticket.Iteration >= 4 && ticket.Iteration <= 6:
		ticket.EscalationLevel = 2
		model = "Gemini Flash/Lite"
	case ticket.Iteration >= 7 && ticket.Iteration <= 9:
		ticket.EscalationLevel = 3
		model = "Gemini Pro Lite / Thinking Model"
	case ticket.Iteration >= 10 && ticket.Iteration <= 11:
		ticket.EscalationLevel = 4
		model = "Gemini Pro"
	case ticket.Iteration >= 12:
		ticket.Status = "STALLED"
		log.Printf("Ticket %s STALLED after 11 iterations. Flagging for HUMAN INTERVENTION.", ticketID)
	}

	if ticket.Status != "STALLED" {
		log.Printf("Iteration %d: Escalating Ticket %s to Level %d using %s", ticket.Iteration, ticketID, ticket.EscalationLevel, model)
		
		// 1. Call the LLM to generate a resolution
		resolution, err := h.callModel(ctx, model, content)
		if err != nil {
			log.Printf("Healing Error: Failed to call model %s: %v", model, err)
			return err
		}

		// 2. Process the resolution (in a real system, this might involve code execution)
		log.Printf("Healing Solution Generated: %s", resolution)

		// Evolutionary Memory: Fetch previous failures and similar resolutions
		failedAttempts := content.FailedAttempts
		similar, _ := h.repo.FindSimilarResolutions(ctx, content.StateVector, 3)

		// Record audit
		entry := domain.AuditEntry{
			TicketID: ticket.ID,
			AgentID:  "healing-service",
			Action:   fmt.Sprintf("HEALING_ITERATION_%d", ticket.Iteration),
			NewValue: map[string]interface{}{
				"iteration":        ticket.Iteration,
				"escalation_level": ticket.EscalationLevel,
				"model_assigned":   model,
				"kb_context_count": len(similar),
				"failed_attempts":  len(failedAttempts),
			},
		}
		h.repo.AddAudit(ctx, entry)
	}

	return h.repo.Update(ctx, ticket.ID, ticket.Status, "")
}

// MarkResolved finalizes a ticket and adds it to the evolutionary knowledge base
func (h *HealingService) MarkResolved(ctx context.Context, ticketID uuid.UUID, resolution string, agentID string) error {
	return h.repo.(interface {
		UpdateExtended(context.Context, uuid.UUID, string, string, string, string) error
	}).UpdateExtended(ctx, ticketID, "COMPLETED", "", resolution, agentID)
}

// RecordFailure logs a failed attempt to avoid repetition in higher tiers
func (h *HealingService) RecordFailure(ctx context.Context, ticketID uuid.UUID, failure string) error {
	return h.repo.AddFailedAttempt(ctx, ticketID, failure)
}

// CreateHealingTicket starts a new self-healing process from a log issue
func (h *HealingService) CreateHealingTicket(ctx context.Context, issue string, logSnippet string) (uuid.UUID, error) {
	content := domain.FabricContent{
		IntentBlob: map[string]interface{}{
			"type":        "SELF_HEALING",
			"issue":       issue,
			"log_snippet": logSnippet,
		},
		RawContent: []byte(logSnippet),
	}

	// Healing tickets start at Layer 5 (Operational)
	return h.svc.CreateFabricTicket(ctx, 5, "monitor-001", content)
}
func (h *HealingService) callModel(ctx context.Context, modelName string, content *domain.FabricContent) (string, error) {
	prompt := fmt.Sprintf("System: You are a PQR Healing Agent. Resolve the following issue.\nContext: %v\nIssue: %s\nResolution:", 
		content.IntentBlob, string(content.RawContent))

	resp, node, err := h.ai.QuerySwarm(ctx, prompt)
	if err != nil {
		return "", err
	}
	log.Printf("[HEALING-AGENT] Solution obtained from %s", node)
	return resp, nil
}

// ExecuteDiagnostic uses the Swarm AI to analyze or execute a diagnostic command
func (h *HealingService) ExecuteDiagnostic(ctx context.Context, cmd string) (string, error) {
	prompt := fmt.Sprintf("System: You are a PQR Diagnostic Agent. Analyze and execute the following diagnostic command for the Sovereign Mesh.\nCommand: %s\nOutput:", cmd)
	
	log.Printf("[DIAGNOSTIC] Coaxing LM to process command: %s", cmd)
	resp, engine, err := h.ai.QuerySwarm(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("diagnostic failure: %v", err)
	}
	
	return fmt.Sprintf("[%s DIAGNOSTIC OUTPUT]\n%s", engine, resp), nil
}
