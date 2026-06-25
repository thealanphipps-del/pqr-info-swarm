package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/thealanphipps-del/pqr/internal/infrastructure/db"
	"os/exec"
	"strings"
)

type SystemSnapshot struct {
	Timestamp     time.Time
	EngineState   []byte
	MemoryState   []byte
	ConfigState   []byte
	ProtoChecksum string
}

type EngineSnapshot struct {
	OSName           string            `json:"os_name"`
	HasWSL           bool              `json:"has_wsl"`
	DefaultTimeout   int               `json:"default_timeout_sec"`
	EmulatorDefaults map[string]string `json:"emulator_defaults"`
	VaultAddress     string            `json:"vault_address"`
	VaultRole        string            `json:"vault_role"`
}

type ConfigSnapshot struct {
	NodeID       string            `json:"node_id"`
	SwarmAddress string            `json:"swarm_address"`
	CockroachDSN string            `json:"cockroach_dsn"`
	Env          map[string]string `json:"env"`
}

type HostFingerprint struct {
	WindowsBuild string `json:"windows_build"`
	WSLKernel    string `json:"wsl_kernel"`
}

type GobackService struct {
	db *db.CockroachRepository
}

func collectHostFingerprint() *HostFingerprint {
	fp := &HostFingerprint{}
	
	// Get Windows Build
	out, err := exec.Command("cmd", "/c", "ver").Output()
	if err == nil {
		fp.WindowsBuild = strings.TrimSpace(string(out))
	}

	// Get WSL Kernel (if available)
	out, err = exec.Command("wsl", "uname", "-r").Output()
	if err == nil {
		fp.WSLKernel = strings.TrimSpace(string(out))
	}

	return fp
}

func NewGobackService(repo *db.CockroachRepository) *GobackService {
	return &GobackService{
		db: repo,
	}
}

func (g *GobackService) System(ts string) error {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return err
	}
	
	// Restore knowledge_journal up to ts
	if err := g.db.RewindKnowledge(context.Background(), t); err != nil {
		return fmt.Errorf("failed to rewind knowledge: %v", err)
	}
	
	// Restore fix memory up to ts
	if err := g.Fixes(ts); err != nil {
		return fmt.Errorf("failed to rewind fixes: %v", err)
	}

	// Restore journal entries up to ts
	if err := g.Chain(ts); err != nil {
		return fmt.Errorf("failed to rewind action journal: %v", err)
	}

	// Ultimately, the dashboard /diff handles the state comparison
	// But any system-wide side-effects should happen here.

	return nil
}

func (g *GobackService) Genesis(engineConfig *EngineSnapshot, nodeConfig *ConfigSnapshot, protoPath string) error {
	// 1. Hash proto file
	f, err := os.Open(protoPath)
	if err != nil {
		return fmt.Errorf("failed to open proto file for checksum: %v", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("failed to hash proto file: %v", err)
	}
	protoChecksum := hex.EncodeToString(h.Sum(nil))

	// 2. Serialize configurations
	engineBytes, err := json.Marshal(engineConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal engine config: %v", err)
	}
	configBytes, err := json.Marshal(nodeConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal node config: %v", err)
	}

	// 3. Serialize empty memory state (empty array)
	memoryBytes := []byte("[]")

	// 4. Create Genesis Snapshot
	ctx := context.Background()
	fp := collectHostFingerprint()
	fpBytes, _ := json.Marshal(fp)
	
	if err := g.db.CreateGenesisSnapshot(ctx, engineBytes, memoryBytes, configBytes, protoChecksum, fpBytes); err != nil {
		return fmt.Errorf("failed to insert genesis snapshot: %v", err)
	}

	// 5. Write Installation Boundary Journal Entry
	if err := g.db.LogAction(ctx, "installation", []byte(""), []byte("genesis")); err != nil {
		return fmt.Errorf("failed to log installation action: %v", err)
	}

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
