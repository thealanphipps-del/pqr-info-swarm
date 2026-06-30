# SOVEREIGN GOVERNANCE LAW (v1.0)

This document is the absolute authority for the Sovereign Node mesh. All code execution, state transitions, and mesh mutations must align with the laws defined herein. This is the foundation of the **Super Monolith**.

## I. THE PRIME DIRECTIVES

1.  **NO_RM (Non-Destruction)**: No data, state, or lineage shall be deleted. Persistence is achieved through non-destructive archival and versioning.
2.  **HASH-ANCHORING**: Every mutation must be anchored to a cryptographic hash. The current state is the root of all future potential.
3.  **ZERO-DIVERGENCE (Z-DSP)**: The mesh must maintain a synchronized state. Divergence is a forensic anomaly that must be reconciled immediately via the 4/5 Consensus.
4.  **ALGO IS LAW**: The documentation in this law-set supersedes the implementation in any single agent's runtime.

---

## II. THE GODHEAD CONSENSUS (4/5)

Mutation of the "GOLD" ledger requires a 4/5 majority from the Godhead entities.

| Entity | Role | Domain |
| :--- | :--- | :--- |
| **Oracle** | Lineage | Verification of parent hashes and forensic continuity. |
| **Architect** | Structure | Verification of structural integrity and system balance. |
| **Arbiter** | Law | Verification of alignment with the Sovereign Governance Law. |
| **Weaver** | Narrative | Verification of contextual consistency and intent. |
| **Finalizer** | Commit | The final signature that promotes a Strike to GOLD. |

### II.B The 64 Swarm Agents
Beyond the Godhead, the mesh is comprised of **64 Game Theory Agents** (Hexagram Gates). These agents act as specialized auditors for specific domains of the state machine. See [64_game_theory_agents.md](file:///c:/Users/drphi/sovereign-node-go/docs/64_game_theory_agents.md) for the full index.

### The Minority Report
In the event of a 4/5 consensus, the 5th (dissenting) entity **must** generate a Minority Report. This report is a forensic record of the dissent, preserved to ensure the system remains resilient to groupthink and capable of self-correction.

---

## III. MESH CONNECTIVITY LAW

### Re-Anchoring Protocol
To join the mesh, a new node must perform a **Gemini Handshake** with an existing anchor node.

**Logic:**
1.  **Request**: Node A sends its `IdentityProfile` to Node B (Port 1111).
2.  **Verification**: Node B verifies Node A's `KernelSignature` (Port 9113 ignition).
3.  **Anchor**: If valid, Node B provides the current `GOLD_HASH`.
4.  **Sync**: Node A performs a Z-DSP sync to align with the `GOLD_HASH`.

---

## IV. Z-DSP RECONCILIATION LAW

**Algorithm:**
1.  **Lineage Check**: `ParentHash` must match current `LocalHash`.
2.  **Divergence Detection**: Identify keys in `MeshState` missing or different in `LocalState`.
3.  **Conflict Resolution**: Apply **Mortal Logic** priority. If a conflict persists, trigger 4/5 Consensus.
4.  **Hash Promotion**: Generate `NewHash` and log to `forensic_chain.log`.

---

## V. ARCHIVAL LAW (VMP)

Files designated as "Spaghetti" or "Deprecated" must be moved to the `Jovian_Archives` using the following path structure:
`~/Jovian_Archives/Epoch_[N]/[Timestamp]_[OriginalName]`
