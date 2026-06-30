import urllib.request
import urllib.parse
import json
import sys
import time

def run_escalation_test():
    base_url = "http://localhost:8080"
    
    print("PQR Ticketing System - Complex Problem Escalation Self-Healing Test")
    print("====================================================================")
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
            # For iteration requests, we expect a 500 error because AI nodes are offline.
            # We want to capture the status code and response body to verify the system logic still progressed.
            if hasattr(e, 'code'):
                err_body = e.read().decode('utf-8')
                try:
                    return {"http_code": e.code, "body": json.loads(err_body)}
                except:
                    return {"http_code": e.code, "body": err_body}
            print(f"Request to {path} failed: {e}")
            return None

    # 1. Create a complex healing ticket
    print("1. Creating a complex self-healing ticket...")
    healing_payload = {
        "issue": "Cascade Failure: Mesh Tunnel Out of Sync with Vault Credentials",
        "logSnippet": "CRITICAL [tunnel-sync] SAML authentication failed: invalid vault signature key path /etc/certs/vault.pem"
    }
    ticket = request("/REST/2.0/healing/ticket", method="POST", data=healing_payload)
    if not ticket or "id" not in ticket:
        print("Failed to create self-healing ticket")
        sys.exit(1)
    ticket_id = ticket["id"]
    print(f"Created Ticket ID: {ticket_id}\n")

    # 2. Iterate and escalate through all levels (up to 12 iterations to reach STALLED state)
    print("2. Simulating healing iterations & tracking model/level escalation...")
    for iter_num in range(1, 13):
        print(f"\n--- Running Iteration {iter_num} ---")
        
        # Call the iterate endpoint
        res = request(f"/REST/2.0/healing/iterate/{ticket_id}", method="POST")
        
        # Note: AI nodes offline will return HTTP 500. However, the database transaction
        # increments the iteration counter first! Let's get the ticket details to verify.
        details = request(f"/REST/2.0/ticket/{ticket_id}")
        
        status = details.get("status")
        # In a real system, the iteration increment and status check are stored in the db.
        # Let's inspect the status and details returned.
        print(f"HTTP Endpoint Response: {res}")
        print(f"Ticket Status: {status}")
        print(f"Current DB Details: {json.dumps(details, indent=2)}")
        
        # Record a failure at each iteration to simulate the problem remaining unresolved
        fail_payload = {
            "ticketID": ticket_id,
            "failure": f"Failed resolution attempt at iteration {iter_num} due to persistent config error."
        }
        request("/REST/2.0/healing/failure", method="POST", data=fail_payload)
        time.sleep(0.15)

    # 3. Retrieve final ticket details and audit trail
    print("\n====================================================================")
    print("3. Fetching Final Ticket Status & Audit Trail...")
    final_details = request(f"/REST/2.0/ticket/{ticket_id}")
    print(f"Final Ticket Details: {json.dumps(final_details, indent=2)}")
    
    audit = request(f"/REST/2.0/ticket/{ticket_id}/audit")
    print(f"Final Audit Trail: {json.dumps(audit, indent=2)}")
    
    print("\n====================================================================")
    if final_details.get("status") == "STALLED":
        print("✓ Complex Escalation Loop Test Successful: Ticket correctly transitioned to STALLED after 12 iterations!")
    else:
        print("✗ Test Failed: Ticket status is not STALLED.")

if __name__ == "__main__":
    run_escalation_test()
