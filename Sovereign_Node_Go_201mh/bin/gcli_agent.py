#!/data/data/com.termux/files/usr/bin/python
import sys, os, urllib.request, json, socket

ROOT = "/data/data/com.termux/files/home/Sovereign_Node_Go"
ENV = f"{ROOT}/.env"

def get_key():
    try:
        with open(ENV) as f:
            for line in f:
                if "GEMINI_API_KEY" in line: return line.split("=")[1].strip().strip('"')
    except: return None

key = get_key()
input_data = sys.stdin.read().strip()

if not key or not input_data:
    sys.exit(0)

# Forensic DNS Check
try:
    socket.gethostbyname("generativelanguage.googleapis.com")
except socket.gaierror:
    print("Error (1) DNS FAILURE: Cannot resolve Gemini API. Check 39.mh Mesh Status.")
    sys.exit(1)

payload = {"contents": [{"parts": [{"text": input_data}]}]}
url = f"https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-pro-latest:generateContent?key={key}"

req = urllib.request.Request(url, data=json.dumps(payload).encode(), headers={'Content-Type': 'application/json'})

try:
    with urllib.request.urlopen(req, timeout=15) as res:
        print(json.loads(res.read().decode())['candidates'][0]['content']['parts'][0]['text'])
except Exception as e:
    print(f"Error (1) API Request Failed: {e}")
