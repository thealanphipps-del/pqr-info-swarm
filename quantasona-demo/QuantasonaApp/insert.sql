INSERT INTO system_manifest (key, value, updated_at) VALUES (
  'network_topology',
  '{
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
  }'::JSONB,
  CURRENT_TIMESTAMP
) ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP;

INSERT INTO system_manifest (key, value, updated_at) VALUES (
  'config_manifest',
  '{
    "hosts": {
      "antigravity-aurora-r9": {
        "ip": "127.0.0.1",
        "os": "Windows 11 / WSL 2",
        "role": "Development Host & Execution Node (Antigravity IDE)"
      },
      "0.mh": {
        "ip": "46.224.84.64",
        "role": "Sovereign Mesh Node / PQR Root Authority"
      },
      "38.mh": {
        "ip": "62.238.2.240",
        "role": "Sovereign Mesh Node / 5D Router Peer Discovery"
      },
      "39.mh": {
        "ip": "204.168.138.60",
        "role": "Sovereign Mesh Node / SSH Fallback Tunnel"
      },
      "201.mh": {
        "ip": "89.167.91.81",
        "role": "Sovereign Mesh Node"
      }
    },
    "devices": {
      "R3GYB02VRQL": {
        "type": "Android Physical Device",
        "status": "authorized",
        "app": "QuantasonaApp"
      }
    },
    "database": {
      "engine": "CockroachDB",
      "host": "localhost",
      "port": 26257,
      "database": "antigravity",
      "url": "postgresql://root@localhost:26257/antigravity?sslmode=disable"
    },
    "services": {
      "grpc_logged_consensus": 1111,
      "grpc_neural_gossip": 11111,
      "nginx_http": 3196,
      "nginx_https": 443,
      "swend_api_rest": 8196,
      "vault_api": 8200
    },
    "tokens_and_keys": {
      "vault_token": "pqr-vault-token",
      "sovereign_seal": "㉗ (U+3257)"
    }
  }'::JSONB,
  CURRENT_TIMESTAMP
) ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP;
