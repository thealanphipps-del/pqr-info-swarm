package service

import (
	"context"
	"time"

	"github.com/thealanphipps-del/pqr/internal/infrastructure/db"
)

type SystemSnapshot struct {
	Timestamp     time.Time
	EngineState   []byte
	MemoryState   []byte
	ConfigState   []byte
	ProtoChecksum string
}

type GobackService struct {
	db *db.CockroachRepository
}

func NewGobackService(repo *db.CockroachRepository) *GobackService {
	return &GobackService{
		db: repo,
	}
}

func (g *GobackService) System(ts string) error {
	// Restore system_snapshots up to ts
	return nil
}

func (g *GobackService) Last() error {
	// Undo last journal entry
	return g.db.UndoLast(context.Background(), func(action string, state []byte) error {
		// apply state based on action type
		return nil
	})
}

func (g *GobackService) Chain(ts string) error {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return err
	}
	// Undo chain of entries
	return g.db.UndoChain(context.Background(), t, func(action string, state []byte) error {
		// apply state based on action type
		return nil
	})
}

func (g *GobackService) Fixes(ts string) error {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return err
	}
	// Restore fix memory
	return g.db.RewindAllFixes(t)
}
