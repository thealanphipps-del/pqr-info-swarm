package execution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/thealanphipps-del/pqr/internal/infrastructure/db"
)

type ExecutionEngine struct {
	osName      string
	credManager *RAMCredentialManager
	Memory      *db.CockroachRepository
}

type RAMCredentialManager struct {
	// Add fields for short-lived tokens, service account JSON, etc.
}

func NewExecutionEngine() *ExecutionEngine {
	return &ExecutionEngine{
		osName:      runtime.GOOS,
		credManager: &RAMCredentialManager{},
	}
}

// ExecuteWorkflow is the generic entrypoint for swarm-triggered tasks.
func (e *ExecutionEngine) ExecuteWorkflow(ctx context.Context, wf *WorkflowDescriptor) (*ExecutionResult, error) {
	// 0. Setup timeout
	timeout := 60 * time.Second
	if t, ok := wf.Metadata["timeout_sec"].(float64); ok {
		timeout = time.Duration(t) * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 1. First attempt
	res, err := e.runOnce(ctx, wf)
	if err == nil && res.Success {
		return res, nil
	}
	
	if errors.Is(err, context.DeadlineExceeded) {
		res = &ExecutionResult{
			Success: false,
			Logs: []string{
				toJSON(logf("error", "command timed out")),
			},
			Error: err,
		}
	}

	// 2. Lookup fix
	if e.Memory == nil {
		return res, err
	}

	sig := e.Memory.HashSignatureFromFields(string(wf.Kind), wf.Command, wf.TargetOS)
	fix, ok := e.Memory.QueryKnownFix(sig)
	if !ok {
		return res, err
	}

	// 3. Apply fix
	fixedCmd := fix.ApplyToCommand(wf.Command)
	wf.Command = fixedCmd

	// 4. Retry
	retryRes, retryErr := e.runOnce(ctx, wf)

	// 5. Update success rate
	if retryErr == nil && retryRes.Success {
		e.Memory.UpdateSuccessRate(sig, true)
	} else {
		e.Memory.UpdateSuccessRate(sig, false)
	}

	return retryRes, retryErr
}

func (e *ExecutionEngine) runOnce(ctx context.Context, wf *WorkflowDescriptor) (*ExecutionResult, error) {
	effectiveOS := e.osName
	if wf.TargetOS != "" {
		effectiveOS = wf.TargetOS
	}
	_ = effectiveOS // Will use this in Shell commands

	switch wf.Kind {
	case WorkflowKindEmulator:
		return e.LaunchEmulator(ctx, wf)
	case WorkflowKindGcpOp:
		return e.RunGcpCommand(ctx, wf)
	case WorkflowKindShell:
		return e.RunShellCommand(ctx, wf)
	default:
		return &ExecutionResult{
			Success: false,
			Logs:    []string{"unknown workflow kind"},
		}, nil
	}
}

func (e *ExecutionEngine) DiagnoseGCP(ctx context.Context) (*ExecutionResult, error) {
	// Check gcloud presence, auth status, project, etc.
	// This is where you can call `gcloud auth list`, `gcloud config list`, etc.
	return e.RunShellCommand(ctx, &WorkflowDescriptor{
		Kind:      WorkflowKindShell,
		Command:   "gcloud auth list",
		ShellHint: ShellHintAuto,
	})
}

func (e *ExecutionEngine) LaunchEmulator(ctx context.Context, wf *WorkflowDescriptor) (*ExecutionResult, error) {
	avdName, _ := wf.Metadata["avd_name"].(string)
	gpuMode, _ := wf.Metadata["gpu_mode"].(string)
	coldBoot, _ := wf.Metadata["cold_boot"].(bool)

	args := append([]string{}, wf.Args...)
	if avdName != "" {
		args = append(args, "-avd", avdName)
	}
	if coldBoot {
		args = append(args, "-no-snapshot-load")
	}
	if gpuMode != "" {
		args = append(args, "-gpu", gpuMode)
	}

	cmdName := wf.Command
	if cmdName == "" {
		cmdName = "emulator"
	}

	cmd := exec.CommandContext(ctx, cmdName, args...)
	out, err := cmd.CombinedOutput()
	
	logs := []string{
		toJSON(logf("info", "executing emulator command: "+cmdName)),
		toJSON(logf("info", "args: "+strings.Join(args, " "))),
		toJSON(logf("stdout", string(out))),
	}
	if err != nil {
		logs = append(logs, toJSON(logf("error", err.Error())))
	}

	return &ExecutionResult{
		Success: err == nil,
		Logs:    logs,
		Error:   err,
	}, nil
}

func (e *ExecutionEngine) RunGcpCommand(ctx context.Context, wf *WorkflowDescriptor) (*ExecutionResult, error) {
	// Use RAM-only credentials; never persist keys.
	// Example: `gcloud <args...>`
	cmd := exec.CommandContext(ctx, "gcloud", wf.Args...)
	out, err := cmd.CombinedOutput()
	
	logs := []string{
		toJSON(logf("info", "executing gcloud op")),
		toJSON(logf("info", "args: "+strings.Join(wf.Args, " "))),
		toJSON(logf("stdout", string(out))),
	}
	if err != nil {
		logs = append(logs, toJSON(logf("error", err.Error())))
	}

	return &ExecutionResult{
		Success: err == nil,
		Logs:    logs,
		Error:   err,
	}, nil
}

func hasWSL() bool {
	_, err := exec.LookPath("wsl")
	return err == nil
}

func (e *ExecutionEngine) RunShellCommand(ctx context.Context, wf *WorkflowDescriptor) (*ExecutionResult, error) {
	effectiveOS := e.osName
	if wf.TargetOS != "" {
		effectiveOS = wf.TargetOS
	}

	var cmd *exec.Cmd

	switch effectiveOS {
	case "wsl":
		if !hasWSL() {
			return &ExecutionResult{
				Success: false,
				Logs:    []string{toJSON(logf("error", "WSL not available on this host"))},
				Error:   fmt.Errorf("wsl not found"),
			}, nil
		}
		cmd = exec.CommandContext(ctx, "wsl", "bash", "-c", wf.Command)
	case "windows":
		cmd = exec.CommandContext(ctx, "powershell", "-Command", wf.Command)
	case "linux":
		cmd = exec.CommandContext(ctx, "bash", "-c", wf.Command)
	default:
		if e.osName == "windows" {
			cmd = exec.CommandContext(ctx, "powershell", "-Command", wf.Command)
		} else {
			cmd = exec.CommandContext(ctx, "bash", "-c", wf.Command)
		}
	}

	out, err := cmd.CombinedOutput()
	
	logs := []string{
		toJSON(logf("info", "executing command: "+wf.Command)),
		toJSON(logf("info", "targetOS: "+effectiveOS)),
		toJSON(logf("stdout", string(out))),
	}
	if err != nil {
		logs = append(logs, toJSON(logf("error", err.Error())))
	}

	return &ExecutionResult{
		Success: err == nil,
		Logs:    logs,
		Error:   err,
	}, nil
}

// --- Types used by the engine ---

type WorkflowKind string

const (
	WorkflowKindEmulator WorkflowKind = "emulator"
	WorkflowKindGcpOp    WorkflowKind = "gcp_op"
	WorkflowKindShell    WorkflowKind = "shell"
)

type ShellHint string

const (
	ShellHintAuto ShellHint = "auto"
)

type WorkflowDescriptor struct {
	Kind      WorkflowKind
	Command   string
	Args      []string
	TargetOS  string
	Metadata  map[string]interface{}
	ShellHint ShellHint
}

type ExecutionResult struct {
	Success bool
	Logs    []string
	Error   error
	Started time.Time
	Ended   time.Time
}

func (r *ExecutionResult) GetLogs() []string {
	return r.Logs
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

func logf(level, msg string) LogEntry {
	return LogEntry{
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Level:     level,
		Message:   msg,
	}
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
