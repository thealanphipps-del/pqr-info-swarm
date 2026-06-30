#!/bin/bash
cd /home/billing/sovereign_mesh
source ~/.sovereign_venv/bin/activate

# 1. Start Web Server with SSL
pkill -f web_server.py || true
nohup python3 -u grpc_node/web_server.py > ./web.log 2>&1 &

# 2. Start Cloudflare Tunnel
TOKEN="cfut_UBEZeWO99lOIPXbwpzHb6pLY03f0K4cvfzKgWWxJdab06f92"

echo "[TUNNEL] Cleaning up old binaries..."
pkill -f cloudflared || true
rm -f /tmp/cloudflared

echo "[TUNNEL] Downloading ARM64 binary..."
wget -q https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-arm64 -O /tmp/cloudflared
chmod +x /tmp/cloudflared

echo "[TUNNEL] Starting tunnel..."
nohup /tmp/cloudflared tunnel run --token "$TOKEN" > ./pqr_tunnel.log 2>&1 &

echo "✅ [GCP] Public services ignited: https://pqr.info (ARM64 Native)"
