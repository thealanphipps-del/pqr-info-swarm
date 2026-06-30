#!/usr/bin/env python3
import os
import sys
import json
import subprocess

# Mesh nodes and their IPs from mudd_interface
# We'll use the IP list directly for network auditing
MESH_NODES = {
    "0.MH": "46.224.84.64",
    "38.MH": "62.238.2.240",
    "39.MH": "204.168.138.60",
    "40.MH": "10.128.0.2",
    "50.MH": "136.113.240.237",
    "201.MH": "89.167.91.81",
    "7.MH": "35.243.111.42",
    "8.MH": "34.93.120.15",
    "9.MH": "34.123.50.201",
    "AURORA": "127.0.0.1"
}

def scan_network():
    print("🌐 NETWORK-SENTINEL AGENT 'SCANNER-01' ACTIVATED")
    manifest = {"timestamp": "", "nodes": {}}
    
    for node_id, ip in MESH_NODES.items():
        print(f"Scanning node: {node_id} ({ip})...")
        node_data = {"ip": ip, "interfaces": [], "open_ports": []}
        
        # 1. Interface Audit (ifconfig/ip addr)
        res_if = subprocess.run(["./mesh_control.sh", "exec", node_id, "ip -brief addr show"], capture_output=True, text=True)
        node_data["interfaces"] = res_if.stdout.splitlines()
        
        # 2. Port Audit (nmap/ss)
        res_ports = subprocess.run(["./mesh_control.sh", "exec", node_id, "ss -tuln"], capture_output=True, text=True)
        node_data["open_ports"] = res_ports.stdout.splitlines()
        
        manifest["nodes"][node_id] = node_data
        
    # Export manifest
    with open("mesh_wide_manifest.json", "w") as f:
        json.dump(manifest, f, indent=4)
        
    print("\n✅ Manifest generated: mesh_wide_manifest.json")

if __name__ == "__main__":
    scan_network()
