import urllib.request
import urllib.parse
import json
import sys
import time

def run_e2e_test():
    base_url = "http://localhost:8080"
    
    print("PQR Ticketing System - Full End-to-End Workflow Test")
    print("=====================================================")
    print(f"Base URL: {base_url}\n")
    
    # Helper for JSON requests
    def request(path, method="GET", data=None):
        url = f"{base_url}{path}"
        req = urllib.request.Request(url, method=method)
        if data is not None:
            req.add_header("Content-Type", "application/json")
            req_data = json.dumps(data).encode("utf-8")
        else:
            req_data = None
            
        try:
            with urllib.request.urlopen(req, data=req_data) as response:
                res_body = response.read().decode("utf-8")
                return json.loads(res_body) if res_body else {}
        except Exception as e:
            print(f"Request to {path} failed: {e}")
            if hasattr(e, 'read'):
                print(f"Error response: {e.read().decode('utf-8')}")
            return None

    # Step 1: Health Check and Schema Init
    print("Step 1: Initializing Schema & System Verification...")
    health = request("/REST/2.0/health")
    print(f"Health check status: {health.get('status')}, Version: {health.get('version')}")
    init_res = request("/REST/2.0/init", method="POST")
    print(f"Schema init response: {init_res.get('message')}\n")

    # Step 2: Agent A (Planner) creates the Parent task ticket
    print("Step 2: Planner Agent A creating parent ticket...")
    parent_payload = {
        "Subject": "Sovereign Node Upgrade",
        "Queue": "planning",
        "Text": "Upgrade all database schemas and verify healing triggers.",
        "AgentID": "planner-agent-001",
        "Layer": 3,
        "Intent": {
            "milestone": "RC2",
            "db_target": "CockroachDB"
        }
    }
    parent_ticket = request("/REST/2.0/ticket", method="POST", data=parent_payload)
    parent_id = parent_ticket["id"]
    print(f"Parent Ticket created with ID: {parent_id}\n")

    # Step 3: Agent A records its initialization context memory
    print("Step 3: Storing initial planning memory for Agent A...")
    mem_payload_a = {
        "memory_type": "context",
        "data": {
            "status": "initiated",
            "upgrade_steps": ["db_migration", "sentinel_validation", "healing_run"]
        },
        "relevance_score": 1.0
    }
    store_a_res = request(f"/REST/2.0/agent/planner-agent-001/memory/{parent_id}", method="POST", data=mem_payload_a)
    print(f"Agent A memory store response: {store_a_res.get('message')}\n")

    # Step 4: Agent B (Coder) creates a Subtask (Child) ticket
    print("Step 4: Coder Agent B creating child ticket...")
    child_payload = {
        "Subject": "Implement Database Migration Script",
        "Queue": "development",
        "Text": "Write SQL migrations and integration tests.",
        "AgentID": "coder-agent-001",
        "Layer": 2,
        "Intent": {
            "repo": "pqr-info-swarm",
            "priority": "high"
        }
    }
    child_ticket = request("/REST/2.0/ticket", method="POST", data=child_payload)
    child_id = child_ticket["id"]
    print(f"Child Ticket created with ID: {child_id}\n")

    # Step 5: Link the Child ticket to the Parent ticket (CONSEQUENCE/EVOLUTION)
    print(f"Step 5: Linking Child Ticket {child_id} to Parent Ticket {parent_id}...")
    link_res = request(f"/REST/2.0/ticket/{parent_id}/link/{child_id}", method="POST", data={
        "relationship_type": "CONSEQUENCE",
        "agent_id": "coder-agent-001"
    })
    print(f"Link response: {link_res.get('message') if link_res else 'Failed'}\n")

    # Step 6: Agent B records its code development state/knowledge
    print("Step 6: Storing coder progress memory...")
    mem_payload_b = {
        "memory_type": "knowledge",
        "data": {
            "migration_version": "v1.09",
            "files_modified": ["migrations.go", "server.go"]
        },
        "relevance_score": 0.9
    }
    store_b_res = request(f"/REST/2.0/agent/coder-agent-001/memory/{child_id}", method="POST", data=mem_payload_b)
    print(f"Agent B memory store response: {store_b_res.get('message')}\n")

    # Step 7: Agent B completes the subtask
    print("Step 7: Marking DB migration subtask as COMPLETED...")
    update_b_res = request(f"/REST/2.0/ticket/{child_id}", method="PUT", data={
        "Status": "COMPLETED",
        "Title": "Database Migration Script Complete"
    })
    print(f"Update response: {update_b_res.get('message')}\n")

    # Step 8: Agent C (QA/Monitor) scrapes a test log showing a failure and creates a Healing ticket
    print("Step 8: QA Agent C scraping log failure and initiating Self-Healing Loop...")
    mock_log = (
        "2026-06-28 14:42:15 FATAL [migration-runner] Migration v1.09 failed: column 'created_at' already exists\n"
        "at github.com/thealanphipps-del/pqr/internal/infrastructure/db.RunMigrations (migrations.go:112)"
    )
    healing_payload = {
        "issue": "Duplicate column 'created_at' in migrations",
        "logSnippet": mock_log
    }
    healing_ticket = request("/REST/2.0/healing/ticket", method="POST", data=healing_payload)
    healing_id = healing_ticket["id"]
    print(f"Self-Healing Ticket created with ID: {healing_id}\n")

    # Step 9: Link the Healing ticket to the Child ticket
    print(f"Step 9: Linking Self-Healing Ticket {healing_id} to child ticket {child_id}...")
    link_healing_res = request(f"/REST/2.0/ticket/{child_id}/link/{healing_id}", method="POST", data={
        "relationship_type": "EVOLUTION",
        "agent_id": "qa-agent-001"
    })
    print(f"Link response: {link_healing_res.get('message') if link_healing_res else 'Failed'}\n")

    # Step 10: Run diagnostic steps, record failure and resolve the healing ticket
    print("Step 10: Recording initial diagnostic failure & manually resolving the healing ticket...")
    
    # Record failure
    fail_payload = {
        "ticketID": healing_id,
        "failure": "SQL syntax check failed: Alter table query is not idempotent."
    }
    request("/REST/2.0/healing/failure", method="POST", data=fail_payload)
    
    # Resolve the ticket
    resolve_payload = {
        "ticketID": healing_id,
        "resolution": "Modified migration SQL to include 'IF NOT EXISTS' for column addition. Re-ran successfully.",
        "agentID": "qa-agent-001"
    }
    resolve_res = request("/REST/2.0/healing/resolve", method="POST", data=resolve_payload)
    print(f"Healing resolution response: {resolve_res.get('message')}\n")

    # Step 11: Agent A reviews audit trails for parent & children, then completes the parent task
    print("Step 11: Agent A verifying resolution and completing parent task...")
    parent_audit = request(f"/REST/2.0/ticket/{parent_id}/audit")
    child_audit = request(f"/REST/2.0/ticket/{child_id}/audit")
    print(f"Parent audit entries: {len(parent_audit.get('audit_trail', [])) if parent_audit else 0}")
    print(f"Child audit entries: {len(child_audit.get('audit_trail', [])) if child_audit else 0}")
    
    update_parent_res = request(f"/REST/2.0/ticket/{parent_id}", method="PUT", data={
        "Status": "COMPLETED",
        "Title": "Sovereign Node Upgrade (COMPLETED)"
    })
    print(f"Parent task completion response: {update_parent_res.get('message')}\n")

    # Step 12: Get final context windows for all involved agents
    print("Step 12: Fetching final context windows for all agents...")
    for agent in ["planner-agent-001", "coder-agent-001", "monitor-001"]:
        ctx_win = request(f"/REST/2.0/agent/{agent}/context")
        tickets_in_ctx = ctx_win.get("context_tickets") if ctx_win else []
        if tickets_in_ctx is None:
            tickets_in_ctx = []
        print(f" - {agent} has {len(tickets_in_ctx)} ticket(s) in active context")
        for t in tickets_in_ctx:
            print(f"   * [{t.get('status')}] {t.get('id')}: {t.get('intent', {}).get('title', t.get('content'))}")
            
    print("\n=====================================================")
    print("✓ Full End-to-End Workflow Test Completed Successfully!")

if __name__ == "__main__":
    run_e2e_test()
