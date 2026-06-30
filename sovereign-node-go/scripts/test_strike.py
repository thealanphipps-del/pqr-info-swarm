import requests
import json

def test_propose_strike():
    url = "http://localhost:11111/strike/propose"
    payload = {
        "agent_id": "GATE-01",
        "layer": 3,
        "cluster": "INITIATOR",
        "content": {
            "intent_blob": {
                "objective": "Expand Mesh Connectivity",
                "target_node": "Node-X"
            },
            "consensus_score": 0.0,
            "raw_content": "Proposal to expand the mesh by adding Node-X as a forensic anchor.",
        }
    }

    print(f"[TEST] Sending Strike Proposal to {url}...")
    try:
        response = requests.post(url, json=payload, timeout=10)
        response.raise_for_status()
        print("[SUCCESS] Strike Proposed:")
        print(json.dumps(response.json(), indent=2))
    except Exception as e:
        print(f"[ERROR] Failed to propose strike: {e}")

if __name__ == "__main__":
    test_propose_strike()
