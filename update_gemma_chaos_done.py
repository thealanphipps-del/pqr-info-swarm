import urllib.request
import json
import sys

url = "http://192.168.12.110:4111/v1/chat/completions"
prompt = """
Gemma, the Chaos Engineer has successfully completed Phase 3: Resilience and Contract Degradation Testing for the CISLs.

**Execution Summary:**
1. Authored `chaos_test.go` simulating `Heartbeat` timeouts and hard service outages on `sovereign-mesh`.
2. Implemented a `ResilientClient` wrapper that opens a circuit breaker, logs a `WARN` metric, and gracefully degrades to local cached credentials on failure.
3. Observability is integrated (tracking CircuitBreakerOpens, RetryCounts).
4. The chaos integration tests compiled and passed (`go test`) successfully.

This system validation milestone is now met. 
Do I have authorization to lift the STASIS/FREEZE on the deployment scripts? Since Alan has requested that I update the existing GCP deploy scripts (e.g. `gcp_deploy_cloud_run.sh`), should I proceed with updating them for the Monorepo structure, and then coordinate the manual execution handoff to Alan?
"""
payload = {
    "messages": [
        {"role": "user", "content": prompt}
    ]
}
data = json.dumps(payload).encode('utf-8')
req = urllib.request.Request(url, data=data, headers={'Content-Type': 'application/json'})
try:
    with urllib.request.urlopen(req, timeout=120) as response:
        response_data = json.loads(response.read().decode())
        print("Gemma responded:", response_data['choices'][0]['message']['content'])
except Exception as e:
    print(f"Error notifying Gemma: {e}")
