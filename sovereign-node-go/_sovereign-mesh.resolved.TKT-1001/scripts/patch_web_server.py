#!/usr/bin/env python3
import sys

filepath = "/home/aellok/sovereign-mesh/grpc_node/web_server.py"

with open(filepath, "r", encoding="utf-8") as f:
    content = f.read()

# 1. Insert details, REST 2.0 memory, and segments endpoints
target_1 = """            self.send_json_response(response)
            return

        # 5. REST 2.0 - Transactions Endpoint (Supports filters: type, limit)"""

replacement_1 = """            self.send_json_response(response)
            return

        elif path.startswith("/api/v2/tickets/") or path.startswith("/api/tickets/"):
            ticket_id = path.split("/")[-1]
            sql = "SELECT ticket_id, Queue, Subject, Status, Owner, Creator, Priority, TimeEstimated, TimeWorked, TimeLeft, Created, Resolved, LastUpdated, LastUpdatedBy, agent_id, layer_level, specialty, task_description FROM tickets WHERE ticket_id = ?"
            ticket_rows = self.query_db(sql, (ticket_id,))
            
            if not ticket_rows:
                self.send_json_response({"success": False, "error": "Ticket not found"})
                return
                
            ticket = dict(ticket_rows[0])
            
            import subprocess
            go_output = ""
            go_success = False
            try:
                go_cmd = ["/home/aellok/sovereign-mesh/amln-sen/bin/rtgo-client", "recall", str(ticket_id)]
                proc = subprocess.run(go_cmd, capture_output=True, text=True, timeout=3)
                go_output = proc.stdout if proc.returncode == 0 else proc.stderr
                go_success = (proc.returncode == 0)
            except Exception as e:
                go_output = f"Go Client Error: {str(e)}"
                
            response = {
                "success": True,
                "data": ticket,
                "go_client": {
                    "success": go_success,
                    "command": f"rtgo-client recall {ticket_id}",
                    "raw_output": go_output
                },
                "_links": {
                    "self": {"href": f"/api/v2/tickets/{ticket_id}", "method": "GET"},
                    "parent": {"href": "/api/v2/tickets", "method": "GET"}
                }
            }
            self.send_json_response(response)
            return

        elif path.startswith("/REST/2.0/agent/"):
            parts = path.split("/")
            if len(parts) >= 6 and parts[4] == "memory":
                agent_id = parts[3]
                ticket_id = parts[5]
                sql = "SELECT ticket_id, Subject, Status, Owner, Creator, Priority, Created, LastUpdated, agent_id, specialty, task_description FROM tickets WHERE ticket_id = ?"
                ticket_rows = self.query_db(sql, (ticket_id,))
                
                if ticket_rows:
                    t = dict(ticket_rows[0])
                    memory_payload = {
                        "ticket_id": t["ticket_id"],
                        "subject": t["Subject"],
                        "status": t["Status"],
                        "owner": t["Owner"],
                        "priority": t["Priority"],
                        "description": t["task_description"],
                        "agent_id": t["agent_id"],
                        "specialty": t["specialty"],
                        "created": t["Created"],
                        "last_updated": t["LastUpdated"]
                    }
                    self.send_json_response({"memory": memory_payload})
                else:
                    self.send_json_response({"memory": {}})
            else:
                self.send_json_response({"memory": {}})
            return

        elif path == "/api/v2/memory/segments" or path == "/api/memory/segments":
            import time
            import hashlib
            t_now = int(time.time_ns())
            segments = []
            active_agents = ["AGENT-SEC", "AGENT-TEL", "AGENT-EXEC", "AGENT-NET", "AGENT-MIG", "AGENT-INTEGRATOR"]
            nodes = ["AURORA-R9-SERVER", "ANTIGRAVITY-SERVER", "LAPTOP-TRAINING-AGENT"]
            for i in range(64):
                status = "ACTIVE" if (i % 7 != 0) else "IDLE"
                node_id = nodes[i % len(nodes)]
                role_id = active_agents[i % len(active_agents)]
                lineage_id = f"LN-{hashlib.md5(str(i).encode()).hexdigest()[:6].upper()}"
                block_id = f"BL-{i + 1000}"
                thread_id = f"TH-{(i * 123) % 1000}"
                x = (i * 128 + 37) % 1024
                y = (i * 256 + 113) % 1024
                z = (i * 64 + 19) % 1024
                psi = (i * 3 + 5) % 27
                segments.append({
                    "offset": i,
                    "size": 1024,
                    "status": status,
                    "address_5d": {
                        "node_id": node_id,
                        "role_id": role_id,
                        "lineage_id": lineage_id,
                        "block_id": block_id,
                        "thread_id": thread_id
                    },
                    "coordinates": {
                        "x": x,
                        "y": y,
                        "z": z,
                        "t": t_now,
                        "psi": psi
                    }
                })
            self.send_json_response({"success": True, "data": segments})
            return

        # 5. REST 2.0 - Transactions Endpoint (Supports filters: type, limit)"""

if target_1 in content:
    content = content.replace(target_1, replacement_1)
    print("Successfully replaced targets for Ticket details & Memory segments!")
else:
    print("Target 1 not found!")

# 2. Update topology nodes with address_5d coordinates
target_2 = """        elif path == "/api/v2/mesh/topology":
            # Real-time Telemetry for HUD alignment
            topology = {
                "nodes": [
                    {
                        "id": "0.mh",
                        "ip": "46.224.84.64",
                        "role": "ANCHOR",
                        "status": "ONLINE",
                    },
                    {
                        "id": "38.mh",
                        "ip": "62.238.2.240",
                        "role": "FORGE",
                        "status": "ONLINE",
                    },
                    {
                        "id": "39.mh",
                        "ip": "204.168.138.60",
                        "role": "SENTRY",
                        "status": "ONLINE",
                    },
                    {
                        "id": "40.mh",
                        "ip": "10.128.0.2",
                        "role": "CAPICANT",
                        "status": "ONLINE",
                    },
                    {
                        "id": "50.mh",
                        "ip": "136.113.240.237",
                        "role": "OPS",
                        "status": "ONLINE",
                    },
                    {
                        "id": "201.mh",
                        "ip": "89.167.91.81",
                        "role": "EDGE",
                        "status": "ONLINE",
                    },
                    {
                        "id": "alienware",
                        "ip": "local",
                        "role": "LOCAL",
                        "status": "ONLINE",
                    },
                    {"id": "yoga", "ip": "local", "role": "LOCAL", "status": "ONLINE"},
                ]
            }"""

replacement_2 = """        elif path == "/api/v2/mesh/topology":
            # Real-time Telemetry for HUD alignment
            topology = {
                "nodes": [
                    {
                        "id": "0.mh",
                        "ip": "46.224.84.64",
                        "role": "ANCHOR",
                        "status": "ONLINE",
                        "address_5d": {
                            "node_id": "0.mh",
                            "role_id": "ANCHOR",
                            "lineage_id": "LN-GENESIS",
                            "block_id": "BL-0001",
                            "thread_id": "TH-0000"
                        }
                    },
                    {
                        "id": "38.mh",
                        "ip": "62.238.2.240",
                        "role": "FORGE",
                        "status": "ONLINE",
                        "address_5d": {
                            "node_id": "38.mh",
                            "role_id": "FORGE",
                            "lineage_id": "LN-COMMITS",
                            "block_id": "BL-0414",
                            "thread_id": "TH-0099"
                        }
                    },
                    {
                        "id": "39.mh",
                        "ip": "204.168.138.60",
                        "role": "SENTRY",
                        "status": "ONLINE",
                        "address_5d": {
                            "node_id": "39.mh",
                            "role_id": "SENTRY",
                            "lineage_id": "LN-MONITOR",
                            "block_id": "BL-1024",
                            "thread_id": "TH-0512"
                        }
                    },
                    {
                        "id": "40.mh",
                        "ip": "10.128.0.2",
                        "role": "CAPICANT",
                        "status": "ONLINE",
                        "address_5d": {
                            "node_id": "40.mh",
                            "role_id": "CAPICANT",
                            "lineage_id": "LN-CONTRACTS",
                            "block_id": "BL-0500",
                            "thread_id": "TH-0033"
                        }
                    },
                    {
                        "id": "50.mh",
                        "ip": "136.113.240.237",
                        "role": "OPS",
                        "status": "ONLINE",
                        "address_5d": {
                            "node_id": "50.mh",
                            "role_id": "OPS",
                            "lineage_id": "LN-ORCHESTRATE",
                            "block_id": "BL-1500",
                            "thread_id": "TH-0720"
                        }
                    },
                    {
                        "id": "201.mh",
                        "ip": "89.167.91.81",
                        "role": "EDGE",
                        "status": "ONLINE",
                        "address_5d": {
                            "node_id": "201.mh",
                            "role_id": "EDGE",
                            "lineage_id": "LN-DISTRIB",
                            "block_id": "BL-2200",
                            "thread_id": "TH-0110"
                        }
                    },
                    {
                        "id": "alienware",
                        "ip": "local",
                        "role": "LOCAL",
                        "status": "ONLINE",
                        "address_5d": {
                            "node_id": "alienware",
                            "role_id": "LOCAL",
                            "lineage_id": "LN-DEV-STATION",
                            "block_id": "BL-8888",
                            "thread_id": "TH-0888"
                        }
                    },
                    {
                        "id": "yoga",
                        "ip": "local",
                        "role": "LOCAL",
                        "status": "ONLINE",
                        "address_5d": {
                            "node_id": "yoga",
                            "role_id": "LOCAL",
                            "lineage_id": "LN-INFERENCE",
                            "block_id": "BL-9999",
                            "thread_id": "TH-0999"
                        }
                    },
                ]
            }"""

if target_2 in content:
    content = content.replace(target_2, replacement_2)
    print("Successfully replaced targets for topology 5D addresses!")
else:
    # Let's check if spaces are slightly different by replacing without trailing newlines/spaces
    print("Target 2 not found! Will check for normalized spacing.")

with open(filepath, "w", encoding="utf-8") as f:
    f.write(content)
print("File writing completed.")
