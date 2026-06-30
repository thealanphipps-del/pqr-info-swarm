#!/bin/bash
# Sovereign Mesh - Cloudflare Tunnel Orchestrator for pqr.info
set -e

# Load token from environment or absolute path
TOKEN="eyJhIjoiYzA3Y2NjYTM3ZDMyN2Y3Nzk5NjcxZDIzMDhmODBiZWIiLCJ0IjoiNzUzMDQ5NzAtMWZhMC00Y2Q0LTg1M2QtOGUzZGQyM2NjOGU0IiwicyI6IlpqTTFNRFkzTUdZdE9ESTFZaTAwT1RNM0xXRmtORFl0Wm1Fd01UQXdOR1l3T1dVeiJ9"

echo "🛡️ [TUNNEL] Initiating Cloudflare Argo Tunnel for pqr.info..."

# Ensure cloudflared is available
if [ ! -f "/tmp/cloudflared" ]; then
    wget -q https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64 -O /tmp/cloudflared
    chmod +x /tmp/cloudflared
fi

# Run the tunnel in background
# This will map the tunnel to our local SSL-enabled web server on port 8080
nohup /tmp/cloudflared tunnel --token "$TOKEN" > ./pqr_tunnel.log 2>&1 &

echo "✅ [TUNNEL] pqr.info tunnel ignited in background. Monitoring ./pqr_tunnel.log"
