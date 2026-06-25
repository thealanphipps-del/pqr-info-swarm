package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type JournalEntry struct {
	ID        uuid.UUID
	Timestamp time.Time
	Action    string
	Before    []byte
	After     []byte
}

func (c *CockroachRepository) LogAction(ctx context.Context, action string, before, after []byte) error {
	_, err := c.db.ExecContext(ctx, `
		INSERT INTO action_journal (action, before, after)
		VALUES ($1, $2, $3)
	`, action, before, after)
	return err
}

func (c *CockroachRepository) UndoLast(ctx context.Context, applyFunc func(action string, state []byte) error) error {
	var entry JournalEntry
	err := c.db.QueryRowContext(ctx, `
		SELECT id, action, before 
		FROM action_journal 
		ORDER BY timestamp DESC 
		LIMIT 1
	`).Scan(&entry.ID, &entry.Action, &entry.Before)
	
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	if err := applyFunc(entry.Action, entry.Before); err != nil {
		return err
	}

	_, err = c.db.ExecContext(ctx, "DELETE FROM action_journal WHERE id = $1", entry.ID)
	return err
}

func (c *CockroachRepository) UndoChain(ctx context.Context, start time.Time, applyFunc func(action string, state []byte) error) error {
	rows, err := c.db.QueryContext(ctx, `
		SELECT id, action, before 
		FROM action_journal 
		WHERE timestamp >= $1 
		ORDER BY timestamp DESC
	`, start)
	if err != nil {
		return err
	}
	defer rows.Close()

	var toDelete []uuid.UUID

	for rows.Next() {
		var entry JournalEntry
		if err := rows.Scan(&entry.ID, &entry.Action, &entry.Before); err != nil {
			return err
		}
		if err := applyFunc(entry.Action, entry.Before); err != nil {
			return err
		}
		toDelete = append(toDelete, entry.ID)
	}

	for _, id := range toDelete {
		_, _ = c.db.ExecContext(ctx, "DELETE FROM action_journal WHERE id = $1", id)
	}
	
	return nil
}
