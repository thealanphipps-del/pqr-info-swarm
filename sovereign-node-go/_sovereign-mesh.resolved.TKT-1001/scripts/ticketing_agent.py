#!/usr/bin/env python3
import sqlite3
import time
import requests
import json

DB_PATH = "/home/aellok/sovereign-mesh/agent_pedigree.db"
INFER_URL = "http://127.0.0.1:8085/api/v1/infer"

AGENT_OPTIONS = {
    "AGENT-SEC": "Security Overlord (Security & Compliance, logins, credentials, Red-Teaming)",
    "AGENT-TEL": "Telemetry Watcher (System Diagnostics, hardware metrics, Loki logs, status reports)",
    "AGENT-EXEC": "Execution Spawner (Remote execution, running scripts/commands)",
    "AGENT-NET": "Port Auditor (Network ports, tunnels, gRPC connections, socket bindings)",
    "AGENT-MIG": "Migrator Engine (Mind migration, state teleportation, memory bus allocation)",
    "AGENT-INTEGRATOR": "Swarm System Integrator (HUD dashboard, REST APIs, general coordination)"
}

def get_unassigned_tickets():
    try:
        conn = sqlite3.connect(DB_PATH)
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()
        cursor.execute("SELECT ticket_id, Subject, task_description FROM tickets WHERE (agent_id IS NULL OR agent_id = '' OR agent_id = 'None' OR agent_id = 'SYSTEM') AND Status != 'resolved' LIMIT 5")
        rows = cursor.fetchall()
        conn.close()
        return [dict(r) for r in rows]
    except Exception as e:
        print(f"[TICKET-AGENT] DB Error fetching tickets: {e}")
        return []

def assign_ticket(ticket_id, agent_id, specialty):
    try:
        conn = sqlite3.connect(DB_PATH)
        cursor = conn.cursor()
        cursor.execute(
            "UPDATE tickets SET agent_id = ?, specialty = ?, Status = 'Open', LastUpdated = datetime('now') WHERE ticket_id = ?",
            (agent_id, specialty, ticket_id)
        )
        conn.commit()
        conn.close()
        print(f"[TICKET-AGENT] Successfully delegated Ticket #{ticket_id} -> {agent_id} ({specialty})")
    except Exception as e:
        print(f"[TICKET-AGENT] DB Error updating ticket: {e}")

def classify_ticket(subject, description):
    prompt = f"""You are the Sovereign Ticketing Dispatcher.
Your task is to assign the incoming ticket to the most appropriate agent.
Available Agents:
- AGENT-SEC: Security, authentication, passwords, compliance, Red-Teaming.
- AGENT-TEL: Diagnostics, CPU/memory stats, logs, health checks.
- AGENT-EXEC: Command execution, spawning shell processes.
- AGENT-NET: Port mapping, SSH tunnels, gRPC sockets, IP addresses.
- AGENT-MIG: Agent mind migration, state teleportation.
- AGENT-INTEGRATOR: HUD design, REST API endpoints, web layout.

Ticket Subject: {subject}
Ticket Description: {description}

Respond with ONLY the exact agent ID (e.g. AGENT-SEC, AGENT-TEL, etc.). No other text.
"""
    try:
        payload = {
            "prompt": prompt,
            "provider": "LM_STUDIO",
            "model": "gemma-2b"
        }
        res = requests.post(INFER_URL, json=payload, timeout=10)
        if res.status_code == 200:
            text = res.json().get("text", "").strip()
            # Clean up text in case of markdown formatting or extra words
            for aid in AGENT_OPTIONS.keys():
                if aid in text:
                    return aid
        return "AGENT-INTEGRATOR" # Default fallback
    except Exception as e:
        print(f"[TICKET-AGENT] Classification request failed: {e}")
        return "AGENT-INTEGRATOR"

def main():
    print("[TICKET-AGENT] Swarm Ticketing Agent successfully spawned and monitoring...")
    while True:
        tickets = get_unassigned_tickets()
        if tickets:
            print(f"[TICKET-AGENT] Found {len(tickets)} unassigned tickets. Starting delegation...")
            for t in tickets:
                subject = t.get("Subject", "")
                desc = t.get("task_description", "")
                ticket_id = t.get("ticket_id")
                
                agent_id = classify_ticket(subject, desc)
                specialty = AGENT_OPTIONS.get(agent_id, "Systems Integration").split(" (")[0]
                
                assign_ticket(ticket_id, agent_id, specialty)
        time.sleep(15)

if __name__ == "__main__":
    main()
