#!/bin/bash
export UNIFIED_ROOT="/root/Sovereign_Unified"

cat << 'EOF_GMUDD' > "$UNIFIED_ROOT/bin/gmudd.py"
#!/usr/bin/python3
import sys, os, time, requests

CORE_DIR = "/root/Sovereign_Unified/core"
IPC_LOG = "/root/Sovereign_Unified/gmudd_ipc/gmudd_inbox.log"

def call_oracle(prompt):
    k = os.environ.get("GEMINI_API_KEY")
    if not k: return "[ERROR] GEMINI_API_KEY missing."
    u = f"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key={k}"
    s = "Output ONLY valid Golang code for a Sovereign Node module. No markdown. Use package core."
    try:
        r = requests.post(u, json={
            "contents": [{"parts": [{"text": prompt}]}],
            "systemInstruction": {"parts": [{"text": s}]}
        }, timeout=30)
        return r.json()["candidates"][0]["content"]["parts"][0]["text"].strip()
    except: return "[ERROR] Oracle failed."

def main():
    if len(sys.argv) < 3:
        print("Usage: gmudd scribe [module_name] [prompt]")
        return
    
    cmd = sys.argv[1].lower()
    module = sys.argv[2].lower()
    
    if cmd == "scribe":
        target_path = os.path.join(CORE_DIR, f"{module}.go")
        prompt = " ".join(sys.argv[3:])
        print(f"[*] Consulting Oracle for module: {module}...")
        
        content = call_oracle(prompt)
        
        if "[ERROR]" not in content:
            with open(target_path, "w") as f:
                f.write(content)
            with open(IPC_LOG, "a") as f:
                f.write(f"[{int(time.time())}] [ORACLE_WRITE] Module {module} updated via Scribe\n")
            print(f"-> Success: {target_path} forge complete.")
        else:
            print(content)

if __name__ == "__main__": main()
EOF_GMUDD

chmod +x "$UNIFIED_ROOT/bin/gmudd.py"
echo "Success (0) [PROOF OF EXECUTION] - Oracle wired to Core"
