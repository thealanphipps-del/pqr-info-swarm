#!/data/data/com.termux/files/usr/bin/python3
import sys, os, time, requests

WIKI_DIR = "/data/data/com.termux/files/home/Sovereign_Node_Go/wiki"
CONTENT_FILE = f"{WIKI_DIR}/content.txt"
IPC_LOG = "/data/data/com.termux/files/home/Sovereign_Node_Go/gemini_testify/gmudd_ipc/gmudd_inbox.log"

def call_oracle(prompt):
    k = os.environ.get("GEMINI_API_KEY")
    if not k: return "[ERROR] GEMINI_API_KEY missing."
    u = f"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key={k}"
    try:
        r = requests.post(u, json={"contents":[{"parts":[{"text": prompt}]}]}, timeout=30)
        return r.json()["candidates"][0]["content"]["parts"][0]["text"].strip()
    except: return "Oracle Timeout."

def main():
    if len(sys.argv) < 2: return
    cmd = sys.argv[1].lower()
    if cmd == "scribe":
        req = input("Wiki Prompt > ")
        content = call_oracle(req)
        with open(CONTENT_FILE, "w") as f: f.write(content)
        with open(IPC_LOG, "a") as f: f.write(f"[{int(time.time())}] [SCRIBE] Content Updated\n")
        print(f"-> Wiki Updated: {CONTENT_FILE}")

if __name__ == "__main__": main()
