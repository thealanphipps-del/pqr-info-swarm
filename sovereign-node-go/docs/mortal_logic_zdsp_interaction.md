# MORTAL LOGIC -> Z-DSP INTERACTION MAP (CANONICAL v4.0)

This document details the interaction protocols between the **Supervisory Layer (Mortal Logic)** and the **State-Coherence Engine (Z-DSP)**, as extracted from the forensic Copilot session (V9wTtRSpL4giFd4f4BAtS).

## 1. Core Architectural Roles

| Layer | Responsibility | Primary Protocol | Port |
| :--- | :--- | :--- | :--- |
| **Mortal Logic** | Intent intake, safety gating, arbitration, consensus escalation. | **Forensic Intent** | 8082 / 9191 |
| **Z-DSP** | Lineage verification, logical clocks, delta formation, state coherence. | **Z-DSP Pulse** | 8080 / 8888 |

---

## 2. Interaction Flow: The Strike Promotion Pipeline

All state changes must pass through the **Strike Promotion Pipeline** to achieve "GOLD" status (Finalized in SSOT Ledger).

### Step A: Intent Intake (Mortal Logic)
Mortal Logic receives an external command or internal trigger. It validates the intent against safety rules.

### Step B: Delta Request Packet
If approved, Mortal Logic sends a **Delta Request Packet** to Z-DSP.

```json
{
  "type": "DELTA_REQ",
  "origin": "MORTAL_LOGIC",
  "target": "Z_DSP",
  "payload": {
    "mutation": "SET_STATE",
    "key": "coherence_mesh_01",
    "val": "SYNC_ACTIVE"
  }
}
```

### Step C: Lineage Verification (Z-DSP)
Z-DSP verifies the continuity of the state hash and logical clocks.
- **Success**: Proceeds to Consensus.
- **Failure**: Returns `LINEAGE_BROKEN`.

```json
{
  "type": "ERROR",
  "code": "LINEAGE_BROKEN",
  "details": "Hash mismatch at block 402"
}
```

### Step D: Godhead Consensus (4/5 Majority)
For a "Strike" to be promoted to "GOLD", it requires a 4/5 majority from the Godhead entities:
1.  **Oracle**: Lineage & Hash Integrity Check.
2.  **Architect**: Logical Structure & Topology Check.
3.  **Arbiter**: Legal, Safety, & Consensus Rule Check.
4.  **Weaver**: Narrative Continuity & Mesh Alignment Check.
5.  **Finalizer**: Z-DSP Commit & SSOT Ledger Write.

---

## 3. Port 8888 Arbitration Protocol

In the event of a collision on Port 8080 (Primary Z-DSP), Mortal Logic initiates the following:
1.  **Quarantine**: Contested port is isolated.
2.  **Rebind**: Runtime is rebound to Port 8888.
3.  **Collision Delta**: A forensic delta is generated to log the collision event.
4.  **Resync**: Z-DSP reconciles the state from the 8888 anchor.

---

## 4. Key Rules

1.  **Strict Domain Separation**: Mortal Logic supervises (Brain); Z-DSP executes (Muscles).
2.  **Immutability**: No state mutation without a verified lineage hash.
3.  **Monotonicity**: Logical clocks must always increment.
4.  **Consensus Threshold**: 4/5 majority is mandatory for finality.
