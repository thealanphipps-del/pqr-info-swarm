import urllib.request
import urllib.parse
import json
import sys

def run_test():
    base_url = "http://localhost:8080"
    agent_id = "healing-test-agent"
    
    print("PQR Ticketing System - Log Scraping & Self-Healing Loop Test")
    print("==========================================================")
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

    # 1. Health check
    print("1. Checking Server Health...")
    health = request("/REST/2.0/health")
    print(json.dumps(health, indent=2))
    print()

    # 2. Simulate Log Scraping: Create a self-healing ticket from a mock database connection failure log snippet
    print("2. Simulating Log Scraping & Creating Healing Ticket...")
    mock_log = (
        "2026-06-28 14:38:00 ERROR [db-pool-1] connection refused: dial tcp 127.0.0.1:26257: connect: connection refused\n"
        "at github.com/thealanphipps-del/pqr/internal/infrastructure/db.Connect (db.go:45)\n"
        "at github.com/thealanphipps-del/pqr/internal/service/healing_service.go:50"
    )
    issue_desc = "CockroachDB connection refused on port 26257"
    
    payload = {
        "issue": issue_desc,
        "logSnippet": mock_log
    }
    
    ticket = request("/REST/2.0/healing/ticket", method="POST", data=payload)
    print(json.dumps(ticket, indent=2))
    if not ticket or "id" not in ticket:
        print("Failed to create self-healing ticket")
        sys.exit(1)
    ticket_id = ticket["id"]
    print(f"Healing Ticket ID: {ticket_id}\n")

    # 3. Retrieve ticket details to check status and layer
    print("3. Retrieving Healing Ticket Details...")
    details = request(f"/REST/2.0/ticket/{ticket_id}")
    print(json.dumps(details, indent=2))
    print()

    # 4. Process a healing loop iteration (Escalation iteration 1)
    print("4. Processing Healing Loop Iteration 1...")
    iter1_res = request(f"/REST/2.0/healing/iterate/{ticket_id}", method="POST")
    print(json.dumps(iter1_res, indent=2))
    print()

    # 5. Record a healing failure to simulate a failed resolution attempt
    print("5. Recording Healing Attempt Failure...")
    fail_payload = {
        "ticketID": ticket_id,
        "failure": "Attempted to ping localhost:26257 but connection refused. Diagnostic script output: exit status 1"
    }
    fail_res = request("/REST/2.0/healing/failure", method="POST", data=fail_payload)
    print(json.dumps(fail_res, indent=2))
    print()

    # 6. Process another iteration to escalate to Level 2 (Iter 2)
    print("6. Processing Healing Loop Iteration 2...")
    iter2_res = request(f"/REST/2.0/healing/iterate/{ticket_id}", method="POST")
    print(json.dumps(iter2_res, indent=2))
    print()

    # 7. Resolve the healing ticket
    print("7. Resolving Healing Ticket...")
    resolve_payload = {
        "ticketID": ticket_id,
        "resolution": "Re-started CockroachDB service using systemctl. Port 26257 is now accepting connections.",
        "agentID": agent_id
    }
    resolve_res = request("/REST/2.0/healing/resolve", method="POST", data=resolve_payload)
    print(json.dumps(resolve_res, indent=2))
    print()

    # 8. Retrieve updated ticket details
    print("8. Checking Final Healing Ticket Status...")
    final_details = request(f"/REST/2.0/ticket/{ticket_id}")
    print(json.dumps(final_details, indent=2))
    print()

    # 9. Get the Audit Trail of the self-healing loop
    print("9. Getting Audit Trail for Self-Healing Ticket...")
    audit = request(f"/REST/2.0/ticket/{ticket_id}/audit")
    print(json.dumps(audit, indent=2))
    print()

    print("==========================================================")
    print("Self-Healing Test Complete!")

if __name__ == "__main__":
    run_test()
