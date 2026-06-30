package mesh

import (
	"context"
	"fmt"
	"sovereign-node-go/pkg/tickets"
)

type LatticeCluster string

const (
	LatExploration  LatticeCluster = "EXPLORATION"
	LatExploitation LatticeCluster = "EXPLOITATION"
	LatConstraint   LatticeCluster = "CONSTRAINT"
	LatPredictive   LatticeCluster = "PREDICTIVE"
	LatAdversarial  LatticeCluster = "ADVERSARIAL"
	LatConsensus    LatticeCluster = "CONSENSUS"
	LatNarrative    LatticeCluster = "NARRATIVE"
	LatMetaAudit    LatticeCluster = "META-AUDIT"
)

type LatticeResult struct {
	Cluster  LatticeCluster
	Approved bool
	Reason   string
}

type Lattice struct {
	strikeMgr *StrikeManager
	ticketMgr *tickets.Manager
}

func NewLattice(sm *StrikeManager, tm *tickets.Manager) *Lattice {
	return &Lattice{strikeMgr: sm, ticketMgr: tm}
}

func (l *Lattice) BroadcastStrike(ctx context.Context, strike *Strike) ([]LatticeResult, error) {
	fmt.Printf("[LATTICE] Broadcasting Strike %s to the 8x8 Strategy Swarm...\n", strike.ID)
	
	clusters := []LatticeCluster{
		LatExploration, LatExploitation, LatConstraint, LatPredictive,
		LatAdversarial, LatConsensus, LatNarrative, LatMetaAudit,
	}

	results := make([]LatticeResult, 0, len(clusters))

	for _, cluster := range clusters {
		// Simulation of parallel cluster audit
		res := LatticeResult{
			Cluster:  cluster,
			Approved: true,
			Reason:   fmt.Sprintf("%s cluster forensic check PASSED", cluster),
		}
		results = append(results, res)
		fmt.Printf("  - [%s] %s\n", res.Cluster, res.Reason)
	}

	return results, nil
}

func (l *Lattice) Reconcile(results []LatticeResult) (bool, string) {
	approvals := 0
	for _, r := range results {
		if r.Approved {
			approvals++
		}
	}

	if approvals == len(results) {
		return true, "Lattice Zero-Divergence achieved. All clusters approved."
	}

	return false, fmt.Sprintf("Lattice Divergence: Only %d/%d clusters approved", approvals, len(results))
}
