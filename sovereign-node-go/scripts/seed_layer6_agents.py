#!/usr/bin/env python3
"""
Seed Layer 6 (Shared Resource Bindings) agents into the Sovereign Swarm
via the REST 2.0 /api/v2/agents endpoint on port 8080.
"""

import urllib.request
import json
import time

BASE_URL = "http://127.0.0.1:8080/api/v2/agents"

# Layer 6 — Shared Resource Bindings
# These agents govern cross-domain resource allocation, memory brokering,
# inter-layer I/O arbitration and the GMODEM VBR rail interfaces.
LAYER_6_AGENTS = [
    {
        "agent_id": "AGENT-L6-MEMBROKER",
        "name": "Memory Broker",
        "parent_agent_id": "AGENT-0",
        "layer_level": 6,
        "specialty": "Shared Memory Management",
        "subspecialty": "POSIX /dev/shm allocation & zero-copy RAM-page brokering"
    },
    {
        "agent_id": "AGENT-L6-NETBIND",
        "name": "Network Binding Arbiter",
        "parent_agent_id": "AGENT-0",
        "layer_level": 6,
        "specialty": "Socket & Port Lifecycle",
        "subspecialty": "Port-conflict arbitration, ephemeral binding audit, and DNS rebind detection"
    },
    {
        "agent_id": "AGENT-L6-VAULTGATE",
        "name": "Vault Gateway Sentinel",
        "parent_agent_id": "AGENT-0",
        "layer_level": 6,
        "specialty": "Secret & Credential Routing",
        "subspecialty": "HashiCorp Vault token renewal, CF Access credential relay, and zero-trust pipeline gating"
    },
    {
        "agent_id": "AGENT-L6-DBPOOL",
        "name": "Database Pool Orchestrator",
        "parent_agent_id": "AGENT-0",
        "layer_level": 6,
        "specialty": "Connection Pool Management",
        "subspecialty": "CockroachDB & SQLite connection pooling, query throttle, and dead-lock resolution"
    },
    {
        "agent_id": "AGENT-L6-MODEM",
        "name": "GMODEM VBR Rail Controller",
        "parent_agent_id": "AGENT-0",
        "layer_level": 6,
        "specialty": "Telemetry Channel Management",
        "subspecialty": "Variable bit-rate serial telemetry scheduling on 115200-baud GMODEM rail (Z-DSP layer)"
    },
    {
        "agent_id": "AGENT-L6-FSMOUNT",
        "name": "Filesystem Mount Warden",
        "parent_agent_id": "AGENT-0",
        "layer_level": 6,
        "specialty": "Volume & Mount Governance",
        "subspecialty": "NTFS/WSL mount arbitration, NO_RM Jovian archive VMP enforcement, and bind-mount lifecycle"
    }
]

def seed_agents():
    succeeded = 0
    for ag in LAYER_6_AGENTS:
        payload = json.dumps(ag).encode("utf-8")
        req = urllib.request.Request(
            BASE_URL,
            data=payload,
            headers={"Content-Type": "application/json", "User-Agent": "Sovereign-Seed/2.0"},
            method="POST"
        )
        try:
            with urllib.request.urlopen(req, timeout=5) as resp:
                body = json.loads(resp.read().decode("utf-8"))
                if body.get("success"):
                    print(f"  ✅  Seeded: {ag['agent_id']} ({ag['name']})")
                    succeeded += 1
                else:
                    print(f"  ⚠️  Already exists / conflict: {ag['agent_id']} — {body}")
        except Exception as e:
            print(f"  ❌  Failed to seed {ag['agent_id']}: {e}")
        time.sleep(0.3)

    print(f"\n  Total seeded: {succeeded}/{len(LAYER_6_AGENTS)} Layer-6 Shared Resource Binding agents")

if __name__ == "__main__":
    print("\n🧬  Seeding Layer 6 (Shared Resource Bindings) agents into Sovereign Swarm...\n")
    seed_agents()
    print("\n✅  Layer 6 seeding complete. Refresh the Agent Manager in the dashboard.\n")
