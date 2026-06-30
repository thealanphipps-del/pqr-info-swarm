import json
import psycopg2

db_url = "postgresql://root@localhost:26257/antigravity?sslmode=disable"

topology_data = {
  "topology": {
    "localhost": {
      "port_mappings": [
        {"port": 1111, "service": "gRPC Swarm Consensus", "protocol": "TCP", "status": "bound"},
        {"port": 11111, "service": "gRPC Neural Gossip", "protocol": "TCP", "status": "bound"},
        {"port": 26257, "service": "CockroachDB SQL", "protocol": "TCP", "status": "bound"},
        {"port": 8080, "service": "CockroachDB Console", "protocol": "TCP", "status": "bound"},
        {"port": 8200, "service": "HashiCorp Vault", "protocol": "TCP", "status": "bound"},
        {"port": 8196, "service": "SWEND REST API v2", "protocol": "TCP", "status": "bound"},
        {"port": 3196, "service": "Nginx Gateway (HTTP)", "protocol": "TCP", "status": "bound"},
        {"port": 443, "service": "Nginx Gateway (HTTPS)", "protocol": "TCP", "status": "bound"}
      ]
    },
    "external_peers": {
      "38.mh": {"ip": "62.238.2.240", "services": ["gRPC Swarm Consensus", "gRPC Neural Gossip"]},
      "39.mh": {"ip": "204.168.138.60", "services": ["SSH Tunnel Target"]},
      "0.mh": {"ip": "46.224.84.64", "services": ["PQR Root Registry"]}
    }
  }
}

try:
    conn = psycopg2.connect(db_url)
    cur = conn.cursor()
    
    cur.execute("""
        INSERT INTO system_manifest (key, value, updated_at)
        VALUES (%s, %s, CURRENT_TIMESTAMP)
        ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP
    """, ("network_topology", json.dumps(topology_data)))
    conn.commit()
    print("Successfully saved network topology to database.")
    
    cur.close()
    conn.close()
except Exception as e:
    print("Database error:", e)
