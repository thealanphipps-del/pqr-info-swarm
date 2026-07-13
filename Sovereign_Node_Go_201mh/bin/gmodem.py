#!/data/data/com.termux/files/usr/bin/python3
import os, sys, requests
if len(sys.argv) < 3:
    print("Usage: gmodem <prompt> <output_path>"); sys.exit(1)
prompt, path = sys.argv[1], sys.argv[2]
key = os.environ.get("GEMINI_API_KEY")
if not key: print("[-] Missing GEMINI_API_KEY"); sys.exit(1)
print(f"[j] GMODEM Fetching payload to {path}...")
url = f"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key={key}"
sys_p = "You are GMODEM. Output ONLY raw code in a markdown code block (e.g. ``bash or ```python). NO PROSE."
pay = {"contents": [{"parts": [{"text": f"{sys_p}\n\nUSER: {prompt}"}]}]}
try:
    r = requests.post(url, json=pay, timeout=30)
    t = r.json()["candidates"][0]["content"]["parts"][0]["text"]
    if "```" in t:
        code = t.split("```")[1].split("\n", 1)[1].rsplit("```", 1)[0].strip()
        with open(path, "w") as f: f.write(code)
        os.chmod(path, 0o755)
        print(f"[+] GMODEM SUCCESS: Wrote to {path}")
        os.system("/data/data/com.termux/files/usr/bin/termux-vibrate -d 300 -f")
    else: print("[-] No code block found.")
except Exception as e: print(f"[-] GMODEM FATAL: {e}")