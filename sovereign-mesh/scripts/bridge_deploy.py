import requests
import json

BRIDGE_URL = "http://192.168.12.201:8081/api/bridge"

def provision_node(node_id, ip):
    print(f"Provisioning node {node_id} via RTGO Bridge...")
    payload = {
        "command": "start_mesh_services",
        "node_id": node_id,
        "listen_ip": "0.0.0.0",
        "grpc_port": 1111,
        "web_port": 8085
    }
    try:
        r = requests.post(BRIDGE_URL, json=payload)
        print(f"Response: {r.status_code} - {r.text}")
    except Exception as e:
        print(f"Failed to bridge {node_id}: {e}")

if __name__ == "__main__":
    nodes = {"50.MH": "34.42.122.68", "201.MH": "89.167.91.81"}
    for nid, ip in nodes.items():
        provision_node(nid, ip)
