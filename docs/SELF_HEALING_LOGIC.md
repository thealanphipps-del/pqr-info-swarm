# 🧬 Sovereign Self-Healing Logic

The PQR Swarm uses a "Model-Aware" healing strategy. Every system failure is treated as a "Ticket" that must evolve toward a resolution.

## 📊 The Escalation Hierarchy
Tickets are categorized by **Layers** (1-11):
- **Layers 1-3**: Routine maintenance, data cleanup, and log rotation. (Resolved by local Gemma:2b).
- **Layers 4-6**: Service connectivity, database performance, and model availability. (Resolved by LM Studio / Gemma:9b).
- **Layers 7-9**: Identity conflicts, Tunnel stability, and Security breaches. (Resolved by Gemini Pro).
- **Layers 10-11**: Governance, Strategic Drift, and "Genesis" events. (Requires High Justiciar / Manual Oversight).

## 🧠 Evolutionary Memory
When a ticket is resolved, the successful "Resolution Pattern" is stored in the **Ticketing Fabric**. 
- **Learning**: Future agents at lower layers check the knowledge base before initiating a new reasoning loop.
- **Forensics**: Every attempt (success or failure) is logged with a context vector, allowing for full auditing of the AI's decision-making process.

## 🛰️ Monitoring Loops
1. **Connectivity**: Probes `pqr.info/saml/metadata` (1-min interval).
2. **Identity**: Probes Vault for certificate health (1-min interval).
3. **Strategic**: Probes Gemini Pro for "Sovereign Alignment" (5-min interval).
