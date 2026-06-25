package service

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/thealanphipps-del/pqr/internal/infrastructure/db"
)

type KnowledgeIngester struct {
	db *db.CockroachRepository
}

func NewKnowledgeIngester(repo *db.CockroachRepository) *KnowledgeIngester {
	return &KnowledgeIngester{
		db: repo,
	}
}

func (k *KnowledgeIngester) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		// Run immediately on start
		k.ingestKnowledgeBase(ctx)

		for {
			select {
			case <-ticker.C:
				k.ingestKnowledgeBase(ctx)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (k *KnowledgeIngester) ingestKnowledgeBase(ctx context.Context) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("[swend:knowledge] failed to get user home dir: %v", err)
		return
	}

	targets := []struct {
		Source string
		Path   string
	}{
		{"antigravity", filepath.Join(home, ".antigravity")},
		{"gemini", filepath.Join(home, ".gemini")},
	}

	for _, target := range targets {
		if _, err := os.Stat(target.Path); os.IsNotExist(err) {
			continue // Skip if directory does not exist
		}
		
		err := filepath.WalkDir(target.Path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil // skip errors
			}
			if d.IsDir() {
				return nil
			}

			// We only want text/json/markdown files
			ext := strings.ToLower(filepath.Ext(path))
			if ext != ".md" && ext != ".json" && ext != ".txt" && ext != ".jsonl" {
				return nil
			}

			// Read file
			contentBytes, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			content := string(contentBytes)

			// Simple check if it's already ingested
			// In a robust system, we would check file hashes to prevent duplicates.
			// For Phase 4 MVP, we assume we just do a simplistic search or 
			// use a WHERE NOT EXISTS logic, but here we just insert it.
			// To prevent infinite duplication, let's skip files that might already be in DB based on path.
			// Note: This is an unoptimized approach for the skeleton.
			existing, searchErr := k.db.SearchKnowledge(ctx, path)
			if searchErr == nil && len(existing) > 0 {
				return nil // Already ingested something matching this path
			}

			tags := []string{ext}

			if err := k.db.InsertKnowledge(ctx, target.Source, path, content, tags); err != nil {
				log.Printf("[swend:knowledge] failed to insert %s: %v", path, err)
			}
			return nil
		})

		if err != nil {
			log.Printf("[swend:knowledge] error walking %s: %v", target.Path, err)
		}
	}
}
