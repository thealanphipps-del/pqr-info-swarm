#!/data/data/com.termux/files/usr/bin/python
import os
import json
import subprocess
import requests
from pathlib import Path

CONTEXT_FILE = "/data/data/com.termux/files/home/Sovereign_Node_Go/bin/gemini_context.json"
LOG_FILE = "/data/data/com.termux/files/home/Sovereign_Node_Go/bin/vterm_forensic.log"

def load_context():
    if Path(CONTEXT_FILE).exists():
        with open(CONTEXT_FILE, 'r') as f:
            return json.load(f)
    return {}

def log_action(action, data):
    with open(LOG_FILE, 'a') as f:
        f.write(f"[{action}] {data}\n")

def call_gemini(prompt, context):
    api_key = os.environ.get("GEMINI_API_KEY")
    if not api_key:
        
    
    url = f"https://generativelanguage.googleapis.com/v1beta/models/gemini-3.1-pro:generateContent?key={api_key}"
    headers = {'Content-Type': 'application/json'}
    payload = {
        "contents": [{"parts": [{"text": f"SYSTEM CONTEXT: {json.dumps(context)}\n\nUSER PROMPT: {prompt}"}]}]
    }
    
    try:
        resp = requests.post(url, headers=headers, json=payload)
        resp.raise_for_status()
        return resp.json()['candidates'][0]['content']['parts'][0]['text']
    except Exception as e:
        return f"API ERROR: {str(e)}"

def execute_shell(command):
    print(f"\n[EXECUTING LOCAL COMMAND IN USERSPACE]: {command}")
    try:
        result = subprocess.run(command, shell=True, check=True, capture_output=True, text=True)
        log_action("EXEC_SUCCESS", result.stdout)
        return f"Success (0)\n{result.stdout[-500:]}"
    except subprocess.CalledProcessError as e:
        log_action("EXEC_FAIL", e.stderr)
        return f"Error (1)\n{e.stderr[-500:]}"

def main():
    print("=== S25 FE LOCAL VIRTUAL TERMINAL IGNITION ===")
    print("--- WIREGUARD DOWN STATE ACKNOWLEDGED ---")
    context = load_context()
    
    while True:
        try:
            user_in = input("\nroot@S25-FE-VTERM:~# ")
            if user_in.strip().lower() in ['exit', 'quit']:
                break
            
            log_action("USER_INPUT", user_in)
            response = call_gemini(user_in, context)
            print(f"\n[GEMINI]:\n{response}")
            log_action("GEMINI_RESPONSE", response)
            
            if "```bash" in response:
                blocks = response.split("```bash")
                for block in blocks[1:]:
                    cmd = block.split("```")[0].strip()
                    exec_res = execute_shell(cmd)
                    print(exec_res)
                    
        except KeyboardInterrupt:
            break
        except Exception as e:
            print(f"FATAL ERROR: {str(e)}")
            break

if __name__ == '__main__':
    main()
