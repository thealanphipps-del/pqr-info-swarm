#!/usr/bin/env python3
import http.server
import socketserver
import subprocess
import re

PORT = 4111

class StatsDirectoryHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-Type", "text/html")
        self.end_headers()

        # Run ss -tuln to find listening TCP ports
        ports = []
        try:
            output = subprocess.check_output(["ss", "-tuln"], text=True)
            for line in output.splitlines():
                if "LISTEN" in line:
                    # Match port after : or ]:
                    match = re.search(r'(?::|\]:)(\d+)\s+', line)
                    if match:
                        port_num = int(match.group(1))
                        if port_num not in ports and port_num != PORT:
                            ports.append(port_num)
        except Exception as e:
            ports = [8085, 9990, 1111, 11111, 1113] # Default fallback

        # Sort ports numerically
        ports.sort()

        # Identify known services
        known_services = {
            8085: ("Swarm Web Dashboard / HUD", "Primary web UI for relational database and sovereign mesh status"),
            9990: ("Local HTTP File Server", "Dynamic project asset directory host"),
            1111: ("gRPC Control Bus Engine", "Protobuf-based orchestrator sync endpoint (gRPC API)"),
            11111: ("High-Speed RAM Bus Server", "Zero-copy shared memory synchronization endpoint"),
            1113: ("Stealth Cloud Run Proxy Listener", "Inbound telemetry channel proxy"),
            22: ("SSH Daemon", "Local terminal gateway"),
            53: ("DNS Server", "Local domain binding resolver")
        }

        html_content = f"""<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Sovereign Mesh Services Directory</title>
    <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600&display=swap" rel="stylesheet">
    <style>
        body {{
            font-family: 'Outfit', sans-serif;
            background: linear-gradient(135deg, #0f172a 0%, #1e1b4b 100%);
            color: #f8fafc;
            min-height: 100vh;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            margin: 0;
            padding: 20px;
        }}
        .container {{
            background: rgba(255, 255, 255, 0.03);
            backdrop-filter: blur(16px);
            border: 1px solid rgba(255, 255, 255, 0.08);
            border-radius: 24px;
            padding: 40px;
            max-width: 650px;
            width: 100%;
            box-shadow: 0 20px 50px rgba(0, 0, 0, 0.3);
        }}
        h1 {{
            font-size: 2.2rem;
            font-weight: 600;
            margin-top: 0;
            margin-bottom: 8px;
            background: linear-gradient(to right, #38bdf8, #818cf8);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            text-align: center;
        }}
        .subtitle {{
            color: #94a3b8;
            font-size: 1rem;
            text-align: center;
            margin-bottom: 30px;
        }}
        .grid {{
            display: flex;
            flex-direction: column;
            gap: 16px;
        }}
        .card {{
            background: rgba(255, 255, 255, 0.02);
            border: 1px solid rgba(255, 255, 255, 0.05);
            border-radius: 16px;
            padding: 16px 20px;
            transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
            display: flex;
            align-items: center;
            justify-content: space-between;
        }}
        .card:hover {{
            background: rgba(255, 255, 255, 0.06);
            border-color: rgba(99, 102, 241, 0.4);
            transform: translateY(-2px);
        }}
        .details {{
            display: flex;
            flex-direction: column;
            gap: 4px;
        }}
        .port-tag {{
            font-size: 0.8rem;
            font-weight: 600;
            background: rgba(99, 102, 241, 0.15);
            color: #a5b4fc;
            padding: 4px 10px;
            border-radius: 9999px;
            border: 1px solid rgba(99, 102, 241, 0.3);
            display: inline-block;
            margin-bottom: 4px;
            width: fit-content;
        }}
        .title {{
            font-size: 1.1rem;
            font-weight: 600;
            color: #f1f5f9;
        }}
        .desc {{
            font-size: 0.85rem;
            color: #94a3b8;
        }}
        .btn {{
            background: linear-gradient(135deg, #4f46e5 0%, #3730a3 100%);
            color: #ffffff;
            text-decoration: none;
            padding: 10px 18px;
            border-radius: 12px;
            font-weight: 600;
            font-size: 0.9rem;
            transition: all 0.2s;
            box-shadow: 0 4px 12px rgba(79, 70, 229, 0.3);
        }}
        .btn:hover {{
            opacity: 0.95;
            box-shadow: 0 6px 16px rgba(79, 70, 229, 0.4);
        }}
        .footer {{
            text-align: center;
            margin-top: 30px;
            font-size: 0.8rem;
            color: #64748b;
        }}
    </style>
</head>
<body>
    <div class="container">
        <h1>Sovereign Swarm</h1>
        <div class="subtitle">Active Ports & Management Interfaces on 192.168.12.169</div>
        <div class="grid">
"""

        for p in ports:
            name, desc = known_services.get(p, ("Unknown Port Connection", "No additional service descriptions available."))

            # Form link
            link_url = f"http://192.168.12.169:{p}"
            btn_text = "Navigate"

            if p == 1111 or p == 11111:
                link_url = "#"
                btn_text = "Raw Protocol"

            html_content += f"""
            <div class="card">
                <div class="details">
                    <span class="port-tag">PORT {p}</span>
                    <span class="title">{name}</span>
                    <span class="desc">{desc}</span>
                </div>
                <a href="{link_url}" class="btn" {"style='pointer-events: none; opacity: 0.4;'" if btn_text == "Raw Protocol" else ""}>{btn_text}</a>
            </div>
            """

        html_content += """
        </div>
        <div class="footer">Sovereign Mesh Control Hub &bull; 2026</div>
    </div>
</body>
</html>
"""
        self.wfile.write(html_content.encode("utf-8"))

def start_server():
    import threading
    import subprocess
    
    bind_ips = ["127.0.0.1"]
    try:
        for ip in subprocess.check_output(["hostname", "-I"]).decode().split():
            if ip.startswith("192.168.12."):
                bind_ips.append(ip)
    except Exception:
        pass
        
    socketserver.TCPServer.allow_reuse_address = True

    def serve_on(ip):
        try:
            with socketserver.TCPServer((ip, PORT), StatsDirectoryHandler) as httpd:
                print(f"Stats Directory Server active on {ip}:{PORT}")
                httpd.serve_forever()
        except Exception as e:
            print(f"Error serving on {ip}:{PORT}: {e}")

    for ip in set(bind_ips):
        t = threading.Thread(target=serve_on, args=(ip,), daemon=True)
        t.start()
        
    import time
    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        print("\nStopping stats directory server...")

if __name__ == "__main__":
    start_server()
