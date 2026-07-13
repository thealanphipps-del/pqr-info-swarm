#!/data/data/com.termux/files/usr/bin/python3
import os,sys,json,subprocess,requests
from pathlib import Path

DB_HOST="204.168.138.60"
HOME="/data/data/com.termux/files/home"
SOV=f"{HOME}/Sovereign_Node_Go"
CTX_F=f"{SOV}/bin/gemini_context.json"
LOG_F=f"{SOV}/bin/vterm_forensic.log"

class Unbuffered:
 def __init__(self,s):self.s=s
 def write(self,d):self.s.write(d);self.s.flush()
 def writelines(self,d):
  self.s.writelines(d);self.s.flush()
 def __getattr__(self,a):return getattr(self.s,a)
sys.stdout=Unbuffered(sys.stdout)

class AgenticStrike:
 @staticmethod
 def execute(cmd):
  print(f"\n[KERNEL_IGNITION]:\n{cmd}")
  with open(LOG_F,"a") as f:
   f.write(f"\n-[STRIKE]-\n{cmd}\n")
  try:
   p=subprocess.Popen(cmd,shell=True,
   env=os.environ,stdout=subprocess.PIPE,
   stderr=subprocess.STDOUT,text=True,
   bufsize=1,cwd=SOV)
   for l in p.stdout:
    print(f"  [OUT] {l}",end="")
    with open(LOG_F,"a") as f:f.write(f"  [OUT] {l}")
   p.wait()
   subprocess.run("/data/data/com.termux/files/usr/bin/termux-vibrate -d 300 -f",
   shell=True,env=os.environ)
   return f"\n[SUCCESS] Exit: {p.returncode}"
  except Exception as e:return f"\n[FATAL] {str(e)}"

class GeminiBridge:
 _model="models/gemini-2.5-flash"
 @staticmethod
 def call(prompt,ctx):
  k=os.environ.get("GEMINI_API_KEY")
  u=f"https://generativelanguage.googleapis.com/v1beta/{GeminiBridge._model}:generateContent?key={k}"
  s="You are Gemini. Android Termux Root. EXECUTE STRIKES using absolute paths."
  h={"Content-Type": "application/json"}
  d={"contents": [{"parts":[{"text": prompt}]}], "systemInstruction": {"parts":[{"text": s}]}}
  try:
   r=requests.post(u, headers=h, json=d, timeout=30)
   r.raise_for_status()
   return r.json()["candidates"][0]["content"]["parts"][0]["text"]
  except Exception as e:
   return f"Error: {e}"

if __name__=="__main__":
 if len(sys.argv) > 1 and sys.argv[1] == "buzz":
  print("[BRIDGE_ONLINE] V4.2 STABLE. AWAITING WIKI BUILD ON PORT 9111.")
