import urllib.request
import urllib.parse
import json
import sys

def run_test():
    import subprocess
    print("=== System Diagnostics ===")
    try:
        ps_out = subprocess.check_output(["ps", "-ef"]).decode("utf-8")
        pqr_pid = None
        for line in ps_out.splitlines():
            if "pqr" in line or "server" in line:
                print(line)
                parts = line.split()
                if len(parts) > 1 and "pqr-server" in parts[7]:
                    pqr_pid = parts[1]
        if pqr_pid:
            print(f"Reading environ for pid {pqr_pid}:")
            with open(f"/proc/{pqr_pid}/environ", "rb") as f:
                env_data = f.read().split(b"\x00")
                for item in env_data:
                    item_str = item.decode("utf-8", errors="ignore")
                    if any(x in item_str for x in ["DATABASE", "PORT", "CONFIG"]):
                        print("  ", item_str)
        
        # Try database connection
        print("Connecting to DB...")
        import psycopg2
        conn = psycopg2.connect("postgresql://root@localhost:26257/antigravity?sslmode=disable")
        cur = conn.cursor()
        cur.execute("SELECT 1;")
        print("DB connection success:", cur.fetchone())
        cur.close()
        conn.close()
    except Exception as ex:
        print("Failed to run diagnostics:", ex)
    print("==========================\n")

    base_url = "http://localhost:8080"
    agent_id = "test-agent-001"
    
    print("PQR Ticketing System - Agent Memory API Test")
    print("==========================================")
    print(f"Base URL: {base_url}")
    print(f"Agent ID: {agent_id}\n")
    
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

    # 1. Health check
    print("1. Health Check...")
    health = request("/REST/2.0/health")
    print(json.dumps(health, indent=2))
    print()

    # 2. Init schema
    print("2. Initialize Schema...")
    init_res = request("/REST/2.0/init", method="POST")
    print(json.dumps(init_res, indent=2))
    print()

    # 3. Create ticket
    print("3. Creating Ticket...")
    ticket_payload = {
        "Subject": "Agent Working Memory",
        "Queue": "processing",
        "Text": "Initial task content",
        "AgentID": agent_id,
        "Layer": 2,
        "Intent": {
            "task": "test",
            "priority": "high"
        }
    }
    ticket = request("/REST/2.0/ticket", method="POST", data=ticket_payload)
    print(json.dumps(ticket, indent=2))
    if not ticket or "id" not in ticket:
        print("Failed to create ticket")
        sys.exit(1)
    ticket_id = ticket["id"]
    print(f"Ticket ID: {ticket_id}\n")

    # 4. Store memory
    print("4. Storing Agent Memory...")
    mem_payload = {
        "memory_type": "context",
        "data": {
            "status": "processing",
            "items_processed": 5,
            "items_total": 10,
            "current_item": "data_point_5"
        },
        "relevance_score": 0.95
    }
    store_res = request(f"/REST/2.0/agent/{agent_id}/memory/{ticket_id}", method="POST", data=mem_payload)
    print(json.dumps(store_res, indent=2))
    print()

    # 5. Retrieve memory
    print("5. Retrieving Agent Memory...")
    retrieved = request(f"/REST/2.0/agent/{agent_id}/memory/{ticket_id}?type=context")
    print(json.dumps(retrieved, indent=2))
    print()

    # 6. Get ticket details
    print("6. Getting Ticket Details...")
    details = request(f"/REST/2.0/ticket/{ticket_id}")
    print(json.dumps(details, indent=2))
    print()

    # 7. Store knowledge memory
    print("7. Storing Knowledge Memory...")
    know_payload = {
        "memory_type": "knowledge",
        "data": {
            "patterns": ["pattern_a", "pattern_b"],
            "confidence": 0.87
        },
        "relevance_score": 0.85
    }
    store_know = request(f"/REST/2.0/agent/{agent_id}/memory/{ticket_id}", method="POST", data=know_payload)
    print(json.dumps(store_know, indent=2))
    print()

    # 8. Update ticket status
    print("8. Updating Ticket Status...")
    update_payload = {
        "Status": "PROCESSING",
        "Title": "Updated: Memory Storage Test"
    }
    update_res = request(f"/REST/2.0/ticket/{ticket_id}", method="PUT", data=update_payload)
    print(json.dumps(update_res, indent=2))
    print()

    # 9. Get audit trail
    print("9. Getting Audit Trail...")
    audit = request(f"/REST/2.0/ticket/{ticket_id}/audit")
    print(json.dumps(audit, indent=2))
    print()

    # 10. Get agent context
    print("10. Getting Agent Context...")
    ctx_res = request(f"/REST/2.0/agent/{agent_id}/context")
    print(json.dumps(ctx_res, indent=2))
    print()

    print("==========================================")
    print("Test Complete!")

if __name__ == "__main__":
    run_test()
