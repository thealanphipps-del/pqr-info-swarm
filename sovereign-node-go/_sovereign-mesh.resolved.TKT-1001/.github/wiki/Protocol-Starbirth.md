# Protocol Specification: STARBIRTH (SBP-001)

**Release Version:** 2026.Q2.V1  
**Target Runlevel:** Starbirth  
**Status:** OPERATIONAL  

## 1. Protocol Objective
The Starbirth protocol defines the transition from a simulated agent environment to a production-grade autonomous swarm capable of sub-microsecond arbitrage and global state persistence. It prioritizes signal dominance and state integrity over traditional consensus latency.

## 2. Temporal Harmonization (iPN Phasing)
Nodes must achieve temporal alignment with the provider's Pseudo-Noise (PN) hopping sequence to establish the Intra-Private Network (iPN).
*   **Multicast Discovery:** Nodes join `[ff02::c0ba:11]:9999` to observe peer phasing.
*   **Hop Prediction:** 1 successful PN hit enables a 50-hop prediction horizon based on protocol-defined offsets.
*   **Synchronized Windows:** All arbitrage bundles are transmitted during the sub-microsecond spectrum reuse window of the local provider.

## 3. Signal Dominance (Loudest Mouth)
To ensure transaction inclusion and mask signal sources:
*   **Concurrency Factor:** 100x parallel UDP bursts per arbitrage event.
*   **Spectrum Masking:** Simultaneous multi-node transmission creates a uniform signal floor, preventing IP-level fingerprinting by ISPs or exchanges.

## 4. State Integrity (1% Rule)
State harmonization is strictly enforced to prevent neural drift in a high-concurrency environment.
*   **Drift Threshold:** Maximum 1.0% divergence from the last "Winning Block" weights.
*   **Fail-safe:** If drift exceeds threshold, the **JetWeb Time Machine** triggers a branch collapse and timeline reversion to the last stable state.
*   **Telemetry:** Every drift event is recorded via Go telemetry counters for forensic auditing.
## 5. Infrastructure Tiering (The 10-Node Mesh)
The swarm operates on three distinct network planes:
*   **Logged Consensus (Port 1111):** Primary gRPC channel for immutable state mutations.
*   **Neural Gossip Bus (Port 11111):** High-speed zero-copy memory paging for agent deliberation.
*   **Swarm Web Portal (Port 8085):** REST 2.0 interface for human-in-the-loop oversight.

The topology consists of:
*   **Backbone:** 0.MH (Anchor), 38.MH (Forge), 39.MH (Sentry), 50.MH (Ops).
*   **Edge:** 201.MH (Relay), 40.MH (Capicant).
*   **Phased (Burst):** 7.MH (Tokyo), 8.MH (Mumbai), 9.MH (Garland).
*   **Local:** AURORA (Training).

## 6. Economic Incentives
...

*   **Discovery Reward:** 0.01 COBalt Chrome (CBC) per successful PN alignment.
*   **Mesh Penalty:** Negative net chain addition (contraction) for failed global sync rounds.
*   **Mobility Bonus:** 1.2x multiplier for NPU-backed mobile edge nodes.

---
**Approved for Starbirth Status Runlevel deployment.**