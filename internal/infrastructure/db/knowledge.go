package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type KnowledgeEntry struct {
	ID        uuid.UUID
	Timestamp time.Time
	Source    string
	Path      string
	Content   string
	Tags      []string
}

func (c *CockroachRepository) InsertKnowledge(ctx context.Context, source, path, content string, tags []string) error {
	_, err := c.db.ExecContext(ctx, `
		INSERT INTO knowledge_journal (source, path, content, tags)
		VALUES ($1, $2, $3, $4)
	`, source, path, content, pq.Array(tags))
	return err
}

func (c *CockroachRepository) GetRecentKnowledge(ctx context.Context, limit int) ([]KnowledgeEntry, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT id, timestamp, source, path, content, tags
		FROM knowledge_journal
		ORDER BY timestamp DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []KnowledgeEntry
	for rows.Next() {
		var entry KnowledgeEntry
		if err := rows.Scan(&entry.ID, &entry.Timestamp, &entry.Source, &entry.Path, &entry.Content, pq.Array(&entry.Tags)); err != nil {
			return nil, err
		}
		out = append(out, entry)
	}
	return out, nil
}

func (c *CockroachRepository) SearchKnowledge(ctx context.Context, query string) ([]KnowledgeEntry, error) {
	// For now, doing a basic ILIKE text search on content and path.
	// In production with lots of data, this would use full-text indexing or vector embeddings.
	searchPattern := "%" + query + "%"
	rows, err := c.db.QueryContext(ctx, `
		SELECT id, timestamp, source, path, content, tags
		FROM knowledge_journal
		WHERE content ILIKE $1 OR path ILIKE $1
		ORDER BY timestamp DESC
		LIMIT 50
	`, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []KnowledgeEntry
	for rows.Next() {
		var entry KnowledgeEntry
		if err := rows.Scan(&entry.ID, &entry.Timestamp, &entry.Source, &entry.Path, &entry.Content, pq.Array(&entry.Tags)); err != nil {
			return nil, err
		}
		out = append(out, entry)
	}
	return out, nil
}

func (c *CockroachRepository) RewindKnowledge(ctx context.Context, ts time.Time) error {
	_, err := c.db.ExecContext(ctx, `
		DELETE FROM knowledge_journal WHERE timestamp > $1
	`, ts)
	return err
}
