#!/usr/bin/env python3
import os
import hashlib
import json

AGENT_DIR = ".agent/agents"
SHORTCODE_MAP = {}

def generate_shortcode(agent_name):
    # Deterministic generation: hash of agent name -> 5alpha#XXX
    h = hashlib.md5(agent_name.encode()).hexdigest()
    # Take first 3 chars of hash, map to 0-9a-z range
    code = h[:3]
    return f"5alpha#{code}"

def provision_agents():
    print("🆔 PROVISIONING AGENT IDENTITIES")
    agents = sorted([f for f in os.listdir(AGENT_DIR) if f.endswith(".md")])
    
    for a in agents:
        agent_name = a.replace(".md", "")
        code = generate_shortcode(agent_name)
        SHORTCODE_MAP[agent_name] = code
        print(f"  {code:<12} | {agent_name}")
        
        # Inject shortcode into the agent's MD file
        fpath = os.path.join(AGENT_DIR, a)
        with open(fpath, "r") as f:
            lines = f.readlines()
            
        # Ensure header block
        if not lines[0].startswith("---"):
            lines.insert(0, "---\n")
            lines.insert(1, f"shortcode: {code}\n")
            lines.insert(2, "---\n")
        else:
            found = False
            for i in range(len(lines)):
                if "shortcode:" in lines[i]:
                    lines[i] = f"shortcode: {code}\n"
                    found = True
                    break
            if not found:
                lines.insert(1, f"shortcode: {code}\n")
                
        with open(fpath, "w") as f:
            f.writelines(lines)
            
    with open("agent_identity_map.json", "w") as f:
        json.dump(SHORTCODE_MAP, f, indent=4)
        
    print("\n✅ All 33 agents provisioned with 5alpha# identities.")

if __name__ == "__main__":
    provision_agents()
