package tsre

import (
	"fmt"
	"sync"
	"time"

	pb "github.com/pqr-info/sovereign-mesh/proto"
)

// TCRManager handles the Temporal Conflict Resolution Algorithm (Step 2)
// It operates as a DAG state manager, reconciling Q-State paradoxes before they hit the ledger.
type TCRManager struct {
	mu           sync.RWMutex
	anchorStates map[string]*pb.QStateMutation // Anchor states keyed by 5D Coordinate hash
}

func NewTCRManager() *TCRManager {
	return &TCRManager{
		anchorStates: make(map[string]*pb.QStateMutation),
	}
}

// ReconcileBatch applies a tensor-based merging function to a batch of mutations.
func (tcr *TCRManager) ReconcileBatch(batch *pb.QStateBatch) []*pb.QStateMutation {
	tcr.mu.Lock()
	defer tcr.mu.Unlock()

	var resolvedMutations []*pb.QStateMutation

	for _, mut := range batch.Mutations {
		coordKey := hashCoordinateVector(mut.CoordinateVector)
		
		existingMut, exists := tcr.anchorStates[coordKey]
		if !exists {
			// No collision, anchor the state
			tcr.anchorStates[coordKey] = mut
			resolvedMutations = append(resolvedMutations, mut)
			continue
		}

		// COLLISION DETECTED
		// We have two mutations attempting to alter the exact same 5D spatial coordinate.
		// Resolve using the Human-in-the-Loop Replay Timestamp (Most recent human override wins).
		if mut.HitlReplayTimestamp > existingMut.HitlReplayTimestamp {
			fmt.Printf("[TCR] Resolved Paradox at %s in favor of Sequence %s (HITL Override)\n", coordKey, mut.SequenceId)
			tcr.anchorStates[coordKey] = mut
			resolvedMutations = append(resolvedMutations, mut)
		} else {
			fmt.Printf("[TCR] Rejected Chronos-Spike at %s (Outdated Entropy Hash)\n", coordKey)
			// The incoming mutation is discarded
		}
	}

	return resolvedMutations
}

// Hash coordinate vector mathematically to create a 5D positional key
func hashCoordinateVector(coords []float32) string {
	if len(coords) < 5 {
		return "0.0.0.0.0" // Default fallback
	}
	return fmt.Sprintf("%.2f-%.2f-%.2f-%.2f-%.2f", coords[0], coords[1], coords[2], coords[3], coords[4])
}

// StartReconciliationWorker listens for batches from the TSRE engine and processes them
func (tcr *TCRManager) StartReconciliationWorker(batchCh <-chan *pb.QStateBatch, stitcher *ShardStitcher) {
	go func() {
		fmt.Println("[TCR] Temporal Conflict Resolution DAG Online.")
		for batch := range batchCh {
			start := time.Now()
			resolved := tcr.ReconcileBatch(batch)
			latency := time.Since(start)
			fmt.Printf("[TCR] Batch %s Reconciled: %d incoming -> %d resolved (Latency: %s)\n", 
				batch.EpochId[:8], len(batch.Mutations), len(resolved), latency)
			
			// Dispatch to the Scatter Platter Shard-Stitcher
			if stitcher != nil && len(resolved) > 0 {
				stitcher.Dispatch(resolved)
			}
		}
	}()
}
