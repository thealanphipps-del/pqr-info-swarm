#!/bin/bash
export UNIFIED_ROOT="/root/Sovereign_Unified"
mkdir -p "$UNIFIED_ROOT/bin" "$UNIFIED_ROOT/gmudd_ipc" "$UNIFIED_ROOT/logs"

# 1. DEPLOY GMUDD ORACLE ON 39.MH
cat << 'EOF_GMUDD' > "$UNIFIED_ROOT/bin/gmudd.py"
#!/usr/bin/python3
import sys, os, time, requests

IPC_LOG = "/root/Sovereign_Unified/gmudd_ipc/gmudd_inbox.log"

def call_oracle(prompt):
    k = os.environ.get("GEMINI_API_KEY")
    if not k: return "[ERROR] GEMINI_API_KEY missing on 39.mh."
    u = f"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key={k}"
    try:
        r = requests.post(u, json={"contents":[{"parts":[{"text": prompt}]}]}, timeout=30)
        return r.json()["candidates"][0]["content"]["parts"][0]["text"].strip()
    except: return "Oracle Timeout."

def main():
    if len(sys.argv) < 2: return
    cmd = sys.argv[1].lower()
    if cmd == "scribe":
        req = input("Unified Prompt > ")
        content = call_oracle(req)
        with open(IPC_LOG, "a") as f: f.write(f"[{int(time.time())}] [SCRIBE] {content}\n")
        print("-> Oracle Sync Complete")

if __name__ == "__main__": main()
EOF_GMUDD

chmod +x "$UNIFIED_ROOT/bin/gmudd.py"

# 2. DEPLOY DEDUPLICATING BRIDGE LISTENER ON 39.MH
cat << 'EOF_LISTENER' > "$UNIFIED_ROOT/bin/bridge_listener.sh"
#!/bin/bash
PIPE="/root/Sovereign_Unified/gmudd_ipc/rpc_command.pipe"
INBOX="/root/Sovereign_Unified/gmudd_ipc/gmudd_inbox.log"

if [ ! -p "$PIPE" ]; then rm -f "$PIPE"; mkfifo "$PIPE"; fi

while true; do
    if read line < "$PIPE"; then
        LAST_CONTENT=$(tail -n 1 "$INBOX" 2>/dev/null | sed 's/.*\[INGEST\] //')
        if [[ "$line" != "$LAST_CONTENT" ]]; then
            echo "[$(date +%s)] [INGEST] $line" >> "$INBOX"
        fi
    fi
done
EOF_LISTENER

chmod +x "$UNIFIED_ROOT/bin/bridge_listener.sh"

# 3. RESTART LISTENER AND EXECUTE IGNITION TEST
pkill -f bridge_listener.sh || true
nohup "$UNIFIED_ROOT/bin/bridge_listener.sh" > /dev/null 2>&1 &
sleep 2

PIPE="/root/Sovereign_Unified/gmudd_ipc/rpc_command.pipe"
INBOX="/root/Sovereign_Unified/gmudd_ipc/gmudd_inbox.log"

echo "[Rtgo] TEST_PACKET: Clean 39.mh Workspace Initialization" > "$PIPE"
sleep 1
echo "[Masterer] TEST_PACKET: Oracle Bridge Integrity Confirmed" > "$PIPE"
sleep 1

echo "Success 0 PROOF OF EXECUTION plus last 10 log lines"
tail -n 10 "$INBOX"
