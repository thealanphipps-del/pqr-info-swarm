package mesh

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"sovereign-node-go/pkg/tickets"

	"github.com/google/uuid"
)

type RuleID string

const (
	ST001 RuleID = "ST-001" // Hash-Anchoring
	ST002 RuleID = "ST-002" // Lineage Continuity
	ST003 RuleID = "ST-003" // Zero-Divergence
	ST004 RuleID = "ST-004" // Non-Destructive Archival
	ST005 RuleID = "ST-005" // 4/5 Godhead Consensus
	ST006 RuleID = "ST-006" // Genesis Block Immutability
	ST007 RuleID = "ST-007" // 7-Layer Depth Limit
)

type Governance struct {
	ticketMgr *tickets.Manager

	// Council metrics fields
	CouncilMembers []int
	Wealth         []float64
	Confidence     []float64
	Phenotypes     []string
	TopologyDelta  float64
	TurnoverRate   float64

	// State flags
	ParametersFrozen bool
}

func NewGovernance(mgr *tickets.Manager) *Governance {
	return &Governance{
		ticketMgr:        mgr,
		CouncilMembers:   []int{1, 2, 3, 4, 5},
		Wealth:           []float64{100.0, 150.0, 80.0, 200.0, 120.0},
		Confidence:       []float64{0.9, 0.85, 0.95, 0.7, 0.88},
		Phenotypes:       []string{"Weaver", "Arbiter", "Oracle", "Architect", "Sentinel"},
		TopologyDelta:    0.05,
		TurnoverRate:     0.1,
		ParametersFrozen: false,
	}
}

func (g *Governance) ValidateTransition(ctx context.Context, parentID, childID uuid.UUID, layer int) error {
	// ST-006: Genesis Block Immutability
	genesisID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	if childID == genesisID {
		return fmt.Errorf("[%s] Cannot mutate or redefine the Genesis Block", ST006)
	}

	// ST-007: 7-Layer Depth Limit
	if layer > 7 {
		return fmt.Errorf("[%s] Depth limit exceeded (Max: 7, Current: %d)", ST007, layer)
	}

	// ST-002: Lineage Continuity
	if parentID != uuid.Nil {
		// Check if parent exists in the DB
	}

	return nil
}

func (g *Governance) VerifyConsensus(ctx context.Context, votes map[string]bool) (bool, error) {
	// ST-005: 4/5 Godhead Consensus (Architect, Sentinel, Oracle, Arbiter, Weaver)
	godhead := []string{"ARCHITECT", "SENTINEL", "ORACLE", "ARBITER", "WEAVER"}
	count := 0
	for _, entity := range godhead {
		if votes[entity] {
			count++
		} else if entity == "SENTINEL" {
			return false, fmt.Errorf("[%s] CRITICAL: Sentinel (Safety/NO_RM) vetoed the strike", ST005)
		}
	}

	if count < 4 {
		return false, fmt.Errorf("[%s] Consensus failed: Only %d/5 Godhead entities approved", ST005, count)
	}

	// Slingshot Consensus Upgrade: Check if simulation checks pass before accepting spec mutations
	simulationPassed := votes["SIMULATION_APPROVED"]
	if !simulationPassed && votes["SPEC_MUTATION_PROPOSAL"] {
		return false, fmt.Errorf("[%s] CRITICAL: Spec upgrade rejected. Simulation check failed with non-zero execution drift or consensus risk.", ST003)
	}

	return true, nil
}

type GovernanceLogEntry struct {
	Timestamp      string    `json:"timestamp"`
	Council        []int     `json:"council_composition"`
	AuditResult    string    `json:"audit_result"`
	Recommendation string    `json:"recommendation"`
	AppliedAction  string    `json:"applied_action"`
}

func (g *Governance) FreezeParameterChangesForOneEpoch() {
	g.ParametersFrozen = true
	fmt.Println("[Governance] FreezeParameterChangesForOneEpoch triggered.")
}

func (g *Governance) TriggerCouncilRotation() {
	fmt.Println("[Governance] TriggerCouncilRotation triggered.")
	if len(g.CouncilMembers) > 1 {
		first := g.CouncilMembers[0]
		copy(g.CouncilMembers, g.CouncilMembers[1:])
		g.CouncilMembers[len(g.CouncilMembers)-1] = first
	}
}

func (g *Governance) AuditCouncil() {
	req := CouncilAuditRequest{
		CouncilMembers: g.CouncilMembers,
		Wealth:         g.Wealth,
		Confidence:     g.Confidence,
		Phenotypes:     g.Phenotypes,
		TopologyDelta:  g.TopologyDelta,
		TurnoverRate:   g.TurnoverRate,
	}

	audit := AskOracleForCouncilAudit(req)
	
	appliedAction := "none"
	switch audit.Recommendation {
	case "approve":
		appliedAction = "approve"
	case "delay":
		g.FreezeParameterChangesForOneEpoch()
		appliedAction = "freeze_parameters"
	case "reshuffle":
		g.TriggerCouncilRotation()
		appliedAction = "reshuffle_council"
	}

	// Logging to ledger
	entry := GovernanceLogEntry{
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
		Council:        g.CouncilMembers,
		AuditResult:    audit.Stability,
		Recommendation: audit.Recommendation,
		AppliedAction:  appliedAction,
	}

	b, err := json.Marshal(entry)
	if err == nil {
		_ = os.MkdirAll("/home/aellok/logs", 0755)
		f, err := os.OpenFile("/home/aellok/logs/omnibus-gsh-log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			_, _ = f.Write(append(b, '\n'))
			_ = f.Close()
		}
	}
}
