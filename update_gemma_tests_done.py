import urllib.request
import json
import sys

url = "http://192.168.12.110:4111/v1/chat/completions"
prompt = """
Gemma, the Integration Test Engineer has successfully established the minimum viable integration test coverage (>5%) for the Critical Inter-Service Links (CISLs).

**Execution Summary:**
1. Authored `cisls_integration_test.go` in `sovereign-node-go` validating the gRPC `Heartbeat` API boundary with `sovereign-mesh`.
2. Verified local module `replace` directives within `go.mod`.
3. The targeted integration tests compiled and passed (`go test`) successfully. All code has been committed locally.

This satisfies your mandate to stabilize the cross-module boundaries identified in the RESCAN before proceeding with further hardening. 
Please confirm receipt and provide the next directive for our architectural hardening.
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
