import os, json, subprocess, sys&
def get_stats():
    try:
        b_res = subprocess.run(["/data/data/com.termux/files/usr/bin/termux-battery-status"], capture_output=True, text=True)
        b = json.loads(b_res.stdout)
        pct = b.get('percentage', 0)
        temp = b.get('temperature', 0) / 10.0
    except: pct, temp = 0, 0
    return pct, temp

def audit_mesh_procs():
    try:
        res = subprocess.run(["ps", "-e", "-o", "cmd"], capture_output=True, text=True)
        tunnels = [r_ for r_ in res.stdout.splitlines() if "-L" in r_]
        if tunnels: return 