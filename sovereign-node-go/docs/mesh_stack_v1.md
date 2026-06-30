# Sovereign Mesh Protocol Stack (v1.0)

This document formalizes the consensus rules and protocol layers for the Sovereign Node mesh, as extracted from the forensic Copilot session. 

> [!NOTE]
> This document is partially superseded by the [MORTAL LOGIC -> Z-DSP INTERACTION MAP (CANONICAL v4.0)](file:///c:/Users/drphi/sovereign-node-go/docs/mortal_logic_zdsp_interaction.md).


## 1. The NO_RM Prime Directive (Non-Destructive Consensus)

The core tenet of the Sovereign Mesh is **Persistence through Non-Destruction**. 

### Rules:
1.  **Zero Deletion:** Files are never deleted; they are moved to the `Jovian_Archives` using the Versioned-Move Protocol (VMP).
2.  **Non-Destructive Merges:** Git merges must prioritize state preservation. Divergence is resolved by creating parallel forensic branches rather than overwriting history.
3.  **Hash-Anchoring:** Every mutation (file move, code change, state transition) must be anchored to a cryptographic hash.
4.  **Replayability:** Every state transition must be reproducible from the forensic logs.

## 2. Protocol Layers (Z-DSP, Starbirth, NO_RM)

The following table defines the operational boundaries and consensus layers for the Sovereign Node mesh.

| Component | Port | Purpose | Consensus Layer |
| :--- | :--- | :--- | :--- |
| **Kernel Ignition** | **9113** | Node boot confirmation & readiness | Starbirth |
| **Mesh Gateway (39.mh)** | **8080** | Primary consensus anchor join point | Z-DSP |
| **Collision Rebind** | **8888** | Fallback for Z-DSP drift/collision | Mortal Logic |
| **GHUB App Listener** | **9191** | Webhook receiver for IPC/Mortal Logic | Mortal Logic |
| **RTdb Local** | **5433** | Local NO_RM Postgres instance | NO_RM |
| **RTdb Tunnel** | **5432** | Secure SSH tunnel for remote RTdb sync | NO_RM |
| **Sovereign HUD** | **8082** | Local UI for Mortal Logic monitoring | Mortal Logic |
| **GMODEM VBR Rail** | **115200**| Variable bit-rate telemetry channel | Z-DSP |

## 3. Core Consensus Rules

### Z-DSP (Zero-Drift Sync)
- **Zero-Drift Enforcement:** All nodes must converge to a canonical state; no forks allowed.
- **Parent-Hash Requirement:** Every mutation must reference the previous state hash to prevent destructive edits.
- **Monotonic Logical Clock:** Rejects stale or regressive packets to ensure linear time.

### Starbirth (Autonomy Activation)
- **Autonomy Gate:** Mesh enters autonomous mode only when all nodes are anchored and consistent.
- **Kernel Ignition Rule:** Port 9113 must be bound and verified before node promotion.
- **Identity Profile (v3.1):** The canonical signature for a self-sovereign node, consisting of:
    1.  **KernelSignature:** Cryptographic proof that Port 9113 ignition succeeded.
    2.  **AnchorProof:** Evidence of a stable binding to the 39.mh mesh anchor.
    3.  **SutureState:** Confirmation that all RTdb orphaned transitions are resolved.
    4.  **DeltaLineage:** The full Z-DSP parent_hash chain, ensuring forensic continuity.

### Godhead Consensus (Canonical v4.0)
Promotion from "Strike" to "GOLD" requires a 4/5 majority from the Godhead entities (Oracle, Architect, Arbiter, Weaver, Finalizer). See [Interaction Map](file:///c:/Users/drphi/sovereign-node-go/docs/mortal_logic_zdsp_interaction.md) for details.


### NO_RM (Non-Destructive Consensus)
- **NO_RM Prime Directive:** Absolutely no deletions; only versioned-moves (VMP) are permitted.
- **Hash-Anchored Mutations:** Every state change is immutable and replayable from the forensic ledger.

## 4. Implementation Roadmap

- [x] **VMP Implementation:** Completed (`vmp.ps1`).
- [ ] **Hash-Anchoring:** Integration of Merkle tree logic into the Go binary.
- [ ] **Z-DSP Pulse:** Automated log syncing across mesh nodes on Port 8080.
- [ ] **Starbirth Ignition:** Automated binding of Port 9113 during boot.
