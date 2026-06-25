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

func (c *CockroachRepository) CreateGenesisSnapshot(ctx context.Context, engineBytes, memoryBytes, configBytes []byte, protoChecksum string, hostFingerprint []byte) error {
	_, err := c.db.ExecContext(ctx, `
		INSERT INTO system_snapshots (snapshot_time, engine_state, memory_state, config_state, proto_checksum, is_genesis, host_fingerprint)
		VALUES ($1, $2, $3, $4, $5, TRUE, $6)
	`, time.Now(), engineBytes, memoryBytes, configBytes, protoChecksum, hostFingerprint)
	return err
}

type SystemSnapshotRecord struct {
	Timestamp       time.Time
	EngineState     []byte
	MemoryState     []byte
	ConfigState     []byte
	ProtoChecksum   string
	IsGenesis       bool
	HostFingerprint []byte
}

type ErrorHistoryRecord struct {
	SignatureHash  string
	SnapshotTime   time.Time
	ErrorContext   []byte
	SynthesizedFix []byte
	SuccessRate    float64
}

func (c *CockroachRepository) GetSnapshots(ctx context.Context) ([]SystemSnapshotRecord, error) {
	rows, err := c.db.QueryContext(ctx, "SELECT snapshot_time, engine_state, memory_state, config_state, proto_checksum, is_genesis, host_fingerprint FROM system_snapshots ORDER BY snapshot_time DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []SystemSnapshotRecord
	for rows.Next() {
		var r SystemSnapshotRecord
		if err := rows.Scan(&r.Timestamp, &r.EngineState, &r.MemoryState, &r.ConfigState, &r.ProtoChecksum, &r.IsGenesis, &r.HostFingerprint); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func (c *CockroachRepository) GetSnapshot(ctx context.Context, ts time.Time) (*SystemSnapshotRecord, error) {
	var r SystemSnapshotRecord
	err := c.db.QueryRowContext(ctx, "SELECT snapshot_time, engine_state, memory_state, config_state, proto_checksum, is_genesis, host_fingerprint FROM system_snapshots WHERE snapshot_time = $1", ts).
		Scan(&r.Timestamp, &r.EngineState, &r.MemoryState, &r.ConfigState, &r.ProtoChecksum, &r.IsGenesis, &r.HostFingerprint)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *CockroachRepository) GetJournal(ctx context.Context) ([]JournalEntry, error) {
	rows, err := c.db.QueryContext(ctx, "SELECT id, timestamp, action, before, after FROM action_journal ORDER BY timestamp DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []JournalEntry
	for rows.Next() {
		var r JournalEntry
		if err := rows.Scan(&r.ID, &r.Timestamp, &r.Action, &r.Before, &r.After); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func (c *CockroachRepository) GetErrorHistory(ctx context.Context) ([]ErrorHistoryRecord, error) {
	rows, err := c.db.QueryContext(ctx, "SELECT signature_hash, snapshot_time, error_context, synthesized_fix, success_rate FROM error_solution_history ORDER BY snapshot_time DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ErrorHistoryRecord
	for rows.Next() {
		var r ErrorHistoryRecord
		if err := rows.Scan(&r.SignatureHash, &r.SnapshotTime, &r.ErrorContext, &r.SynthesizedFix, &r.SuccessRate); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
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
