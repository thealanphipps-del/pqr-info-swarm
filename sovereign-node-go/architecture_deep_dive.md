# Architecture Deep Dive: Starburst Super Monolith

The **Starburst Super Monolith** is a distributed meta-programming architecture where **The Algo is Law**. It is designed for high-velocity execution and autonomous synchronization across the Sovereign Node mesh, governed by the [SOVEREIGN_GOVERNANCE_LAW.md](file:///c:/Users/drphi/sovereign-node-go/SOVEREIGN_GOVERNANCE_LAW.md).

## 1. Core Components

### A. GSH (Get Shit Done) - Unified Shell (v4.2)
*   **Path**: `/home/sovereign-node-go/gsh_binary`
*   **Role**: The primary entry point for agentic commands.
*   **Special Logic**: Includes the `--secure-irs` flag which triggers the **Secure Injection Layer**. It acts as a wrapper around standard bash but with forensic-guard awareness.

### B. Sovereign Binary (v10.0 Singularity)
*   **Path**: `/home/sovereign-node-go/sovereign`
*   **Role**: The main execution engine for trading and mesh consensus.
*   **State**: Currently running as PID 13593. It monitors RSI signals (Target: 28.5) and anchors into the `39.mh` mesh.

### C. Secure Injection Layer
*   **Logic**: `secure_injection.py`
*   **Role**: Handles the "MetaGo" aspect—injecting code modifications into a running system without requiring a full restart. This allows for "epoch" transitions (e.g., from Epoch 3 to Epoch 4).

### D. Multi-Chain Contracts
*   **Path**: `/home/sovereign-node-go/contracts/`
*   **Supported Chains**: Arbitrum, Cosmos, Ethereum, Solana.
*   **Role**: Smart contract logic for distributed liquidity and arbitrage.

### E. Vitality Generation Controller
*   **Logic**: `vitality_gen.go`
*   **Role**: Likely manages system health, entropy, or "vitality" signals across the mesh nodes.

---

## 2. Network Topography (Native Ports)

| Service | Port | Protocol | Role |
| :--- | :--- | :--- | :--- |
| **Kernel Ignition** | 9113 | TCP | Node boot confirmation & readiness (Starbirth). |
| **Mesh Gateway** | 8080 | HTTP/TCP| Primary consensus anchor join point (Z-DSP). |
| **Collision Rebind** | 8888 | TCP | Fallback for Z-DSP drift/collision (Mortal Logic). |
| **GHUB Listener** | 9191 | HTTP | Webhook receiver for IPC/Mortal Logic. |
| **RTdb Local** | 5433 | SQL | Local NO_RM Postgres instance. |
| **RTdb Tunnel** | 5432 | SSH | Secure SSH tunnel for remote RTdb sync. |
| **Sovereign HUD** | 8082 | HTTP | Local UI for Mortal Logic monitoring. |
| **GMODEM Rail** | 115200| Serial/VBR| Variable bit-rate telemetry channel. |
| **GSH Bridge** | 8081 | HTTP | The primary command & control link. |
| **Consensus Engine**| 1111 | gRPC | Native Antigravity-to-Antigravity communication. |

---

## 3. Storage & Persistence

*   **Asset Lock**: $814.68 (Hardcoded Baseline).
*   **Jovian Archives**: Centralized forensic log storage at `~/Jovian_Archives`.
*   **Real-Time Ledger**: `rt_ledger.log` (Current state is synced to cloud DB).

---

---

## 5. Forensic Layer (Z-DSP v4.0)

As of Epoch 4, the **Forensic Chain of Custody** is active, adhering to **Canonical v4.0** standards.

### Zero-Divergence Sync Protocol (Z-DSP)
Every modification to the mesh state or configuration must be logged in the `forensic_chain.log`. This ensures that all agents (Human and AI) maintain a perfectly synchronized view of the system.

### Strike Promotion Pipeline (Godhead Consensus)
All mutations are treated as "Strikes" and must be promoted to "GOLD" (SSOT Ledger) via a 4/5 majority from the Godhead entities:
- **Oracle** (Lineage)
- **Architect** (Structure)
- **Arbiter** (Law)
- **Weaver** (Narrative)
- **Finalizer** (Commit)

### Lineage Verification
Implemented in `reconcile_zdsp.py`. Rejects any delta if the `parent_hash` does not match the current state hash (`LINEAGE_BROKEN`).

### Sovereign Mesh Protocol Stack (v1.0)
The system now adheres to the full Protocol Stack (v1.0), including Hash-Anchoring and Replayability.
*   **Documentation**: [mesh_stack_v1.md](file:///c:/Users/drphi/sovereign-node-go/docs/mesh_stack_v1.md)
*   **Prime Directive**: NO_RM (Persistence through Non-Destruction).
