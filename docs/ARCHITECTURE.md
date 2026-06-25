# SWEN Architecture (v1.0.0 Substrate)

**SWEN (Swarm Execution Node)** is a self-healing, cross-OS, metadata-driven execution organism. It acts as the physical actuator for the overarching swarm intelligence, running locally on a host machine while maintaining its own deterministic error memory and temporal reversibility.

## Core Capabilities

### 1. Cross-OS Execution Boundary
SWEN abstracts the host OS from the swarm's intent. Using the `target_os` instruction, SWEN dynamically bridges execution boundaries:
- `windows`: Natively executes via PowerShell.
- `linux`: Natively executes via Bash.
- `wsl`: Transparently bridges Linux payloads onto a Windows host via `wsl bash -c`.

### 2. Metadata-Driven Execution
SWEN treats metadata as first-class declarative configuration. 
- **Emulator Controllers**: Workflows of kind `emulator` consume `avd_name`, `cold_boot`, and `gpu_mode` dynamically from the swarm to construct command paths without hardcoded bash scripts.
- **Timeouts**: Every execution is strictly bounded by `timeout_sec`, mapping hangs and stalls directly to a `context.DeadlineExceeded` error, which feeds into the learning loop.

### 3. Self-Healing Learning Loop
SWEN observes its own failures and learns to fix them autonomously.
- **Deterministic Hashing**: Every workflow failure is hashed deterministically (`Kind` + `Command` + `OS`).
- **Fix Synthesis Engine**: The engine parses `stderr` and structured JSON logs to synthesize immediate patches (e.g., prepending `gcloud auth login &&`, injecting `sudo`, or echoing instructions for missing AVDs).
- **Intelligent Retry**: Before executing any command, SWEN checks its `error_solution_memory`. If a fix exists, it applies it dynamically via a strict DSL (`replace`, `prepend`, `append`) and measures the resulting success rate.

### 4. The Goback System (Time Machine)
SWEN possesses complete temporal reversibility across three distinct layers, providing surgical rollback capabilities.
- **Layer 1: Immutable Snapshots**: `system_snapshots` captures the entire node state (engine config, CockroachDB memory dump, proto checksums) allowing full-node restoration via `swend goback --to <timestamp>`.
- **Layer 2: Change Journaling**: `action_journal` tracks structural configuration and schema updates, enabling sequential undo (`swend goback --last` or `--chain`).
- **Layer 3: Memory-Level Rewind**: `error_solution_history` maintains a ledger of all synthesized fixes. `RewindFix` allows SWEN to unlearn a bad fix and revert to a historically successful state.

## Subsystems

- **gRPC Dialer**: `swarm_client.go` manages the high-throughput `OpenExecutionStream` to the overarching swarm.
- **CockroachDB**: Provides the transactional, replicated memory layer for error signatures, action journals, and system snapshots.
- **Vault Authenticator**: Secures short-lived credentials for GCP operations and swarm authentication.

*This document captures the v1.0.0 stable baseline before entering live swarm integration.*
