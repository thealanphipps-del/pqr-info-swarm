#!/usr/bin/env python3
import os
import sys
import re

# Repo Setup
REPO_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
sys.path.append(REPO_ROOT)
from mudd_interface import ROOMS

def research_and_initialize():
    print("🕵️ RESEARCH AGENT 'DOC-WALKER-01' DEPLOYED")
    print("Performing forensic grid initialization...\n")
    
    # Locate all rooms from docs
    discovered_rooms = {}
    pattern = re.compile(r'<!-- ROOM_START id="([^"]+)" name="([^"]+)" role="([^"]+)" exits="([^"]*)" -->(.*?)<!-- ROOM_END -->', re.DOTALL)
    
    for root, dirs, files in os.walk(REPO_ROOT):
        dirs[:] = [d for d in dirs if d not in ('.git', 'node_modules', 'venv', '.venv')]
        for f in files:
            if f.endswith(".md"):
                fpath = os.path.join(root, f)
                with open(fpath, "r", encoding="utf-8") as md:
                    matches = pattern.finditer(md.read())
                    for m in matches:
                        r_id, r_name, r_role, r_exits, r_desc = m.groups()
                        discovered_rooms[r_id] = {
                            "name": r_name, "role": r_role, "exits": [e.strip() for e in r_exits.split(",") if e.strip()],
                            "desc": r_desc.strip(), "source": os.path.relpath(fpath, REPO_ROOT)
                        }

    # Initialize each discovered room
    initialized = 0
    for r_id, meta in discovered_rooms.items():
        print(f"Initializing {r_id}: {meta['name']}...")
        
        # Cross-reference exits and divergent files
        # Divergent file mapping logic (heuristic based on role)
        divergent = []
        if "ANCHOR" in meta['role']: divergent = ["genesis_ledger.go", "price_feed_sync.py"]
        elif "FORGE" in meta['role']: divergent = ["tournament.go", "forge_engine/src/main.rs"]
        elif "SENTRY" in meta['role']: divergent = ["shadow_controller.go", "indicator_logic.py", "avg_price_query.sql"]
        
        # Commit room to runtime memory
        ROOMS[r_id] = {
            "name": meta['name'], "ip": "VIRTUAL-NODE", "role": meta['role'],
            "desc": meta['desc'], "exits": meta['exits'], "virtual": True,
            "source": meta['source'], "divergent_files": divergent
        }
        
        # Log to Forensic Audit
        with open(f"{REPO_ROOT}/logs/grid_init.log", "a") as log:
            log.write(f"Initialized Room: {r_id} (Source: {meta['source']})\n")
        initialized += 1

    print(f"\n✅ Initialization complete. {initialized} rooms fully initialized.")

if __name__ == "__main__":
    os.makedirs(f"{REPO_ROOT}/logs", exist_ok=True)
    research_and_initialize()
