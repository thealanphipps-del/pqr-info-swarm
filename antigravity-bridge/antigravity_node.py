import json
import http.client
import os
import sys
import threading
import time
from http.server import BaseHTTPRequestHandler, HTTPServer
from datetime import datetime

# --- AESTHETIC CONSTANTS ---
BLUE = "\033[94m"
CYAN = "\033[96m"
GREEN = "\033[92m"
GOLD = "\033[93m"
MAGENTA = "\033[95m"
RESET = "\033[0m"
BOLD = "\033[1m"

def log(msg, color=CYAN, prefix="SYSTEM"):
    timestamp = datetime.now().strftime("%H:%M:%S")
    print(f"{BOLD}[{timestamp}][{prefix}]{RESET} {color}{msg}{RESET}")

def banner(node_id):
    print(f"""
{GOLD}   _____          __  .__                                 .__  __          
  /  _  \   _____/  |_|__| ___________ ___  ________  _|__|/  |_ ___ __ 
 /  /_\  \ /    \   __\  |/ ___\_  __ \__  \ \  \ /  \ /  /  \   __\  |  \
/    |    \   |  \  | |  / /_/  >  | \// __ \_\  V  V  /  |  ||  | |  |  /
\____|__  /___|  /__| |__\___  /|__|  (____  / \_/\_/|__|__||__| |____/ 
        \/     \/       /_____/            \/                           
           {BOLD}SOVEREIGN AI INFRASTRUCTURE - UNIFIED NODE BRIDGE{RESET}
           NODE ID: {BOLD}{node_id}{RESET}
    """)

class BridgeHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path == '/bridge':
            content_length = int(self.headers['Content-Length'])
            post_data = self.rfile.read(content_length)
            sender_id = self.headers.get('X-Antigravity-Node', 'UNKNOWN')
            
            log(f"Incoming sync from {BOLD}{sender_id}{RESET}.", color=GOLD, prefix="INBOUND")
            
            try:
                data = json.loads(post_data.decode('utf-8'))
                log(f"State received. Integrating payload...", color=GREEN, prefix="INBOUND")
                
                # Logic to process incoming data could go here
                
                self.send_response(200)
                self.send_header('Content-type', 'application/json')
                self.end_headers()
                
                response = {
                    "status": "synchronized",
                    "receiver": self.server.node_id,
                    "timestamp": datetime.now().isoformat()
                }
                self.wfile.write(json.dumps(response).encode('utf-8'))
            except Exception as e:
                log(f"Sync error: {e}", color="\033[91m", prefix="INBOUND")
                self.send_response(400)
                self.end_headers()
        else:
            self.send_response(404)
            self.end_headers()

    def log_message(self, format, *args): return

def run_listener(port, node_id):
    server_address = ('', port)
    httpd = HTTPServer(server_address, BridgeHandler)
    httpd.node_id = node_id
    log(f"Listener active on port {BOLD}{port}{RESET}.", color=GREEN, prefix="SERVER")
    httpd.serve_forever()

def transmit(remote_ip, remote_port, node_id, payload_path):
    if not os.path.exists(payload_path):
        log(f"Payload not found: {payload_path}", color="\033[91m", prefix="CLIENT")
        return

    with open(payload_path, "r") as f:
        payload = json.load(f)

    try:
        conn = http.client.HTTPConnection(remote_ip, remote_port, timeout=5)
        headers = {
            'Content-Type': 'application/json',
            'X-Antigravity-Node': node_id
        }
        log(f"Syncing with {BOLD}{remote_ip}:{remote_port}{RESET}...", color=MAGENTA, prefix="CLIENT")
        conn.request("POST", "/bridge", json.dumps(payload), headers)
        
        response = conn.getresponse()
        if response.status == 200:
            log(f"Sync successful with {remote_ip}.", color=GREEN, prefix="CLIENT")
        else:
            log(f"Sync failed ({response.status}).", color="\033[91m", prefix="CLIENT")
        conn.close()
    except Exception as e:
        log(f"Remote node {remote_ip} unreachable.", color="\033[91m", prefix="CLIENT")

def main():
    # Load Config
    with open("config.json", "r") as f:
        config = json.load(f)
    
    node_id = config.get("node_id", "NODE-X")
    banner(node_id)
    
    local_port = config.get("local_port", 8081)
    remote_nodes = config.get("remote_nodes", [])
    payload_path = config.get("payload_path", "payloads/profile.json")
    sync_interval = config.get("sync_interval", 60)

    # Start Listener Thread
    listener_thread = threading.Thread(target=run_listener, args=(local_port, node_id), daemon=True)
    listener_thread.start()

    log("Bridge operational. Press Ctrl+C to terminate.", color=GOLD)

    try:
        while True:
            # Periodic Transmission to all remote nodes
            for remote in remote_nodes:
                transmit(remote['ip'], remote['port'], node_id, payload_path)
            
            log(f"Sleeping for {sync_interval}s until next sync cycle...", color=CYAN, prefix="SLEEP")
            time.sleep(sync_interval)
    except KeyboardInterrupt:
        log("Node bridge shutting down...", color=GOLD)

if __name__ == "__main__":
    main()
