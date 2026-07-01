import urllib.request
import json
import sys

url = "http://192.168.12.110:4111/v1/chat/completions"
prompt = """
Gemma, Alan has issued a new overriding directive: "RESCAN CODEBASES UNDER NEW MONOREPO".

I have performed a directory listing of `/home/aellok/pqr-info-swarm`. It appears the L1 submodules (`sovereign-mesh`, `sovereign-node-go`, `quantasona-demo`, `antigravity-bridge`, etc.) have been fully integrated into a centralized Go Workspace (`go.work`) Monorepo structure.

As Alan's proxy, please provide explicit instructions:
What specific deliverables or analysis do you require from this "RESCAN"? 
Should I:
1. Execute a full AST/Dependency graph traversal across the new monorepo to generate an updated `INDEX.md` or `SWARM_TOPOLOGY.md`?
2. Run standard `go build ./...` and `go test ./...` across the workspace to verify there are no breaking changes from the Monorepo merge?
3. Or something else entirely?
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
