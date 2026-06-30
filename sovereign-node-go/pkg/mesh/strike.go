package mesh

import (
	"context"
	"fmt"
	"sovereign-node-go/pkg/tickets"

	"github.com/google/uuid"
)

type ClusterID string

const (
	ClusterInitiator   ClusterID = "INITIATOR"
	ClusterResponder   ClusterID = "RESPONDER"
	ClusterCoordinator ClusterID = "COORDINATOR"
	ClusterObserver    ClusterID = "OBSERVER"
	ClusterCatalyst    ClusterID = "CATALYST"
	ClusterArbiter     ClusterID = "ARBITER"
	ClusterWeaver      ClusterID = "WEAVER"
	ClusterSovereign   ClusterID = "SOVEREIGN"
)

type StrikeStatus string

const (
	StrikeProposed StrikeStatus = "PROPOSED"
	StrikeAudited  StrikeStatus = "AUDITED"
	StrikeGold     StrikeStatus = "GOLD"
	StrikeRejected StrikeStatus = "REJECTED"
)

type Strike struct {
	ID            uuid.UUID
	ProposedBy    string // AgentID
	TargetLayer   int
	Content       tickets.FabricContent
	RequiredCluster ClusterID
	Status        StrikeStatus
	Audits        map[string]bool // AgentID -> Approved
	Votes         map[string]bool // Godhead Entity -> Approved
}

type StrikeManager struct {
	ticketMgr *tickets.Manager
	gov       *Governance
}

func NewStrikeManager(tm *tickets.Manager, g *Governance) *StrikeManager {
	return &StrikeManager{ticketMgr: tm, gov: g}
}

func (sm *StrikeManager) ProposeStrike(ctx context.Context, agentID string, layer int, cluster ClusterID, content tickets.FabricContent) (*Strike, error) {
	strikeID := uuid.New()
	strike := &Strike{
		ID:              strikeID,
		ProposedBy:      agentID,
		TargetLayer:     layer,
		Content:         content,
		RequiredCluster: cluster,
		Status:          StrikeProposed,
		Audits:          make(map[string]bool),
		Votes:           make(map[string]bool),
	}

	fmt.Printf("[STRIKE] Strike %s proposed by %s for Cluster %s\n", strikeID, agentID, cluster)
	return strike, nil
}

func (sm *StrikeManager) AuditStrike(strike *Strike, auditorID string, approved bool) error {
	// Verify if auditor belongs to the required cluster
	// (This would use a mapping from AgentID to ClusterID)
	strike.Audits[auditorID] = approved
	
	// If at least one approval from the required cluster, mark as AUDITED
	if approved {
		strike.Status = StrikeAudited
	}
	
	return nil
}

func (sm *StrikeManager) PromoteToGold(ctx context.Context, strike *Strike) error {
	if strike.Status != StrikeAudited {
		return fmt.Errorf("strike must be AUDITED before promotion")
	}

	// Verify Godhead Consensus (4/5)
	passed, err := sm.gov.VerifyConsensus(ctx, strike.Votes)
	if err != nil || !passed {
		strike.Status = StrikeRejected
		return fmt.Errorf("godhead consensus failed: %v", err)
	}

	// Promotion: Commit to DB
	ticketID, err := sm.ticketMgr.CreateFabricTicketV71(ctx, strike.TargetLayer, strike.ProposedBy, strike.Content)
	if err != nil {
		return err
	}

	strike.ID = ticketID
	strike.Status = StrikeGold
	fmt.Printf("[GOLD] Strike %s promoted to SSOT Ledger\n", ticketID)

	return nil
}
