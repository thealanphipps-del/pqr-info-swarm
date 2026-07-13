import os, json, subprocess
def run_hud():
    try:
        b = subprocess.run(["/data/data/com.termux/files/usr/bin/termux-battery-status"], capture_output=True, text=True)
        pct = json.loads(b.stdout).get('percentage', 'N/A')
    except: pct = 'ERR'
    cmd = ["/data/data/com.termux/files/usr/bin/termux-dialog", "sheet", "-v", f"BATTERY {pct}%,RUN FORENSIC SYNC,MESH STATUS"]
    try:
        res = subprocess.run(cmd, capture_output=True, text=True)
        sel = json.loads(res.stdout).get("text", "")
        if sel == "RUN FORENSIC SYNC":
            print("Status: Triggering Chain of Custody")
        elif sel == "MESH STATUS":
            os.system("netstat -ant | grep LISTEN")
    except: pass
if __name__ == "__main__": run_hud()