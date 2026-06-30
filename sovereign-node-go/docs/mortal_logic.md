# Mortal Logic: Full Technical Specification (v3.2)

Mortal Logic is the **Supervisory Layer** of the Sovereign Mesh, acting as the system's "executive function." It is responsible for interpreting external events, enforcing safety protocols, and maintaining the human-to-mesh control plane.

## 1. Core Responsibilities
- **Interpreting External Events:** Processing signals from GitHub, the HUD, and Swarm telemetry.
- **Safety Enforcement:** Upholding the **NO_RM Prime Directive** and non-destructive mutation rules.
- **Executive Oversight:** Supervising "Strike" behavior and preventing Z-DSP state drift.
- **Mortal Interface:** Maintaining the human-to-mesh control plane.

## 2. Port Integration
- **Port 9191 ("The Mortal Ear"):** The **GHUB Listener** for Webhook Ingress. It receives GitHub App events (Issues, PRs, Releases), validates them via cryptographic secrets, and transforms them into Z-DSP compliant packets for the Swarm.
- **Port 8888 ("The Mortal Hand"):** The **Collision Rebind** for Runtime Arbitration. This port acts as a safety valve for Z-DSP. When a state conflict occurs, the logic rebinds the request here for manual or supervisory resolution.

## 3. Z-DSP Collision Handling
When two nodes attempt to mutate the same SSOT (Single Source of Truth) object with a drift of less than **1ms**, Mortal Logic triggers the following **Collision Algorithm**:

```bash
# Mortal Logic Collision Algorithm
if conflict_detected(drift < 1ms):
    halt_write_operation
    rebind_to_port 8888
    request_mortal_arbitration
    request Rtdb parity check
    notify Swarm for AST repair
    log collision event
```

## 4. Supervision Rules (The Godhead Consensus)
The governance layer requires a **4/5 majority** for any structural or state changes. The "Five Entities" of the Godhead are:
1.  **The Architect:** Evaluates structural integrity and design coherence.
2.  **The Judge:** Evaluates rule compliance (NO_RM, Z-DSP).
3.  **The Oracle:** Evaluates predictive stability and future drift.
4.  **The Scribe:** Evaluates ledger consistency and SSOT alignment.
5.  **The Rebel:** Evaluates innovation, mutation, and insight validity.

## 5. The Mortal Hand (Port 8888) Protocol
When a collision occurs on Port 8080, the Mortal Hand follows a canonical **four-step arbitration protocol**:
1.  **Step 1 — Collision Detection:** Identifies contested listener state.
2.  **Step 2 — Quarantine Port 8080:** Isolates the contested port to prevent state leakage.
3.  **Step 3 — Rebind Runtime to Port 8888:** Shifts active control to the supervisory bridge.
4.  **Step 4 — Notify Z-DSP:** Generates a synchronization packet to inform the Mesh of the new temporal anchor.
