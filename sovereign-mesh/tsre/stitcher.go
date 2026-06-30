package tsre

import (
	"fmt"
	"sync"
	"time"

	pb "github.com/pqr-info/sovereign-mesh/proto"
)

// ShardStitcher takes resolved mutations from the TCR and splices them into the Scatter Platter.
type ShardStitcher struct {
	workers int
	taskCh  chan []*pb.QStateMutation
	wg      sync.WaitGroup
}

func NewShardStitcher(workerCount int) *ShardStitcher {
	return &ShardStitcher{
		workers: workerCount,
		taskCh:  make(chan []*pb.QStateMutation, 500),
	}
}

// Start boots the worker pool
func (s *ShardStitcher) Start() {
	fmt.Printf("[STITCHER] Booting Scatter Platter Shard-Stitcher with %d discrete thread-workers\n", s.workers)
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.workerLoop(i)
	}
}

func (s *ShardStitcher) workerLoop(id int) {
	defer s.wg.Done()
	for batch := range s.taskCh {
		s.stitchToShard(id, batch)
	}
}

// Dispatch sends a batch of resolved mutations to the thread pool
func (s *ShardStitcher) Dispatch(resolved []*pb.QStateMutation) {
	s.taskCh <- resolved
}

func (s *ShardStitcher) stitchToShard(workerID int, resolved []*pb.QStateMutation) {
	if len(resolved) == 0 {
		return
	}

	start := time.Now()
	// Simulate cryptographic localized hashing into the Scatter Platter blockchain
	time.Sleep(time.Duration(len(resolved)*5) * time.Millisecond)

	fmt.Printf("[STITCHER-%d] Successfully stitched %d mutations into the Scatter Platter Shard (Latency: %s)\n",
		workerID, len(resolved), time.Since(start))
	
	// Gossip Protocol trigger goes here to propagate the shard update to the P2P Mesh
}

// Stop gracefully shuts down the workers
func (s *ShardStitcher) Stop() {
	close(s.taskCh)
	s.wg.Wait()
}
