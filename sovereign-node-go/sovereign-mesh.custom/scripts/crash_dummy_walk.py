#!/usr/bin/env python3
import os
import sys
import json
import pathlib

# Import ROOMS from mudd_interface
REPO_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
sys.path.append(REPO_ROOT)
sys.path.append(os.path.join(REPO_ROOT, "grpc_node"))

try:
    from mudd_interface import ROOMS
except ImportError as e:
    print(f"Failed to import ROOMS from mudd_interface.py: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

def audit_grid():
    print("🛡️ CRASH DUMMY AGENT 'DUMMY-01' ACTIVATED")
    print(f"Walking the mesh grid of {len(ROOMS)} nodes...\n")
    
    findings = []
    
    for node_id, data in ROOMS.items():
        print(f"Checking Node: {node_id} ({data['name']})")
        
        # 1. Check Tunnels (Exits)
        for exit_id in data['exits']:
            if exit_id not in ROOMS:
                finding = f"[MISSING LINK] Node {node_id} points to non-existent exit: {exit_id}"
                print(f"  ❌ {finding}")
                findings.append(finding)
            else:
                # Check for bidirectional tunnel (recommended but not strictly required for some)
                if node_id not in ROOMS[exit_id]['exits']:
                    finding = f"[ASYMMETRIC TUNNEL] Node {exit_id} lacks return path to {node_id}"
                    print(f"  ⚠️  {finding}")
                    # findings.append(finding) # Non-blocking

        # 2. Check Divergent Files
        for f in data.get('divergent_files', []):
            fpath = os.path.join(REPO_ROOT, f)
            if not os.path.exists(fpath):
                finding = f"[MISSING FILE] Node {node_id} references missing divergent file: {f}"
                print(f"  ❌ {finding}")
                findings.append(finding)
            else:
                print(f"  ✅ File verified: {f}")

        # 3. Check Virtual Sources
        if data.get('virtual'):
            source_path = os.path.join(REPO_ROOT, data['source'])
            if not os.path.exists(source_path):
                finding = f"[MISSING SOURCE] Virtual Node {node_id} source missing: {data['source']}"
                print(f"  ❌ {finding}")
                findings.append(finding)

    print("\n" + "="*50)
    if not findings:
        print("🎉 GRID WALK COMPLETE: 100% Integrity. All links and files accounted for.")
    else:
        print(f"🚨 GRID WALK COMPLETE: Found {len(findings)} discrepancies.")
        for f in findings:
            # Simple way to flag for ticketing script
            print(f"- [ ] {f}")
            
if __name__ == "__main__":
    audit_grid()
