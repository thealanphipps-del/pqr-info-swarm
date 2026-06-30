package tsre

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	pb "github.com/pqr-info/sovereign-mesh/proto"
)

// TemporalEngine represents the Chronos-Q Temporal State Reconciliation Engine.
// It buffers incoming MIDI-transcoded 5D mutations and batches them into Quantum Epochs.
type TemporalEngine struct {
	mu           sync.Mutex
	mutationCh   chan *pb.QStateMutation
	batchCh      chan *pb.QStateBatch
	currentBatch *pb.QStateBatch
	epochTicker  *time.Ticker
	ctx          context.Context
	cancel       context.CancelFunc
	tcr          *TCRManager
	stitcher     *ShardStitcher
}

func NewTemporalEngine(epochDuration time.Duration) *TemporalEngine {
	ctx, cancel := context.WithCancel(context.Background())
	engine := &TemporalEngine{
		mutationCh:  make(chan *pb.QStateMutation, 10000), 
		batchCh:     make(chan *pb.QStateBatch, 100),
		epochTicker: time.NewTicker(epochDuration),
		ctx:         ctx,
		cancel:      cancel,
		tcr:         NewTCRManager(),
		stitcher:    NewShardStitcher(4), // 4 concurrent stitcher workers
	}
	engine.stitcher.Start()
	engine.tcr.StartReconciliationWorker(engine.batchCh, engine.stitcher)
	return engine
}

// Start begins the event sourcing pipeline.
func (e *TemporalEngine) Start() {
	e.initNewBatch()
	go e.runPipeline()
	fmt.Println("[TSRE] Chronos-Q Temporal Engine Online.")
}

// Stop halts the engine gracefully.
func (e *TemporalEngine) Stop() {
	e.cancel()
	e.epochTicker.Stop()
	close(e.mutationCh)
	close(e.batchCh)
	e.stitcher.Stop()
}

// IngestMutation pushes a live mix adjustment into the buffer.
func (e *TemporalEngine) IngestMutation(mut *pb.QStateMutation) {
	select {
	case e.mutationCh <- mut:
	case <-e.ctx.Done():
	default:
		fmt.Printf("[TSRE-WARNING] Mutation buffer full, dropping sequence %s\n", mut.SequenceId)
	}
}

func (e *TemporalEngine) runPipeline() {
	for {
		select {
		case <-e.ctx.Done():
			return
		case mut := <-e.mutationCh:
			e.mu.Lock()
			if e.currentBatch != nil {
				e.currentBatch.Mutations = append(e.currentBatch.Mutations, mut)
			}
			e.mu.Unlock()
		case <-e.epochTicker.C:
			e.flushEpoch()
		}
	}
}

func (e *TemporalEngine) initNewBatch() {
	e.currentBatch = &pb.QStateBatch{
		EpochId:       uuid.New().String(),
		WindowStartMs: time.Now().UnixMilli(),
		Mutations:     make([]*pb.QStateMutation, 0),
	}
}

func (e *TemporalEngine) flushEpoch() {
	e.mu.Lock()
	batch := e.currentBatch
	batch.WindowEndMs = time.Now().UnixMilli()
	e.initNewBatch()
	e.mu.Unlock()

	if len(batch.Mutations) > 0 {
		e.processBatch(batch)
	}
}

func (e *TemporalEngine) processBatch(batch *pb.QStateBatch) {
	fmt.Printf("[TSRE] Epoch %s flushed with %d mutations spanning %dms\n", 
		batch.EpochId[:8], len(batch.Mutations), batch.WindowEndMs-batch.WindowStartMs)
	e.batchCh <- batch
}
