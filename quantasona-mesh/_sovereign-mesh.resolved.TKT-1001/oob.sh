#!/bin/bash
# Sovereign Mesh - Out-of-Band (OOB) Ignition & Self-Provisioning Script
set -e

BASE_DIR="/home/billing/sovereign_mesh"
CORE_DIR="$BASE_DIR/mgsh_core"
VENV_DIR="$HOME/.sovereign_venv"

echo "🐍 [OOB] Initializing Sovereign Autonomous Environment..."

# 1. Self-Provisioning: Ensure directory structures exist
mkdir -p "$BASE_DIR" "$CORE_DIR"

# 2. Virtual Environment Check
if [ ! -d "$VENV_DIR" ]; then
    echo "📦 [PROVISION] Creating virtual environment..."
    python3 -m venv "$VENV_DIR"
fi
source "$VENV_DIR/bin/activate"

# 3. Tool Dependency Check
if ! python3 -c "import grpc" 2>/dev/null; then
    echo "🛠️ [PROVISION] Installing gRPC dependencies..."
    pip install --quiet grpcio grpcio-tools
fi

# 4. Correct log paths and Compile Protos
echo "📜 [OOB] Setting up local log paths and compiling protos..."
export ANTIGRAVITY_NODE_ID="GCP-OPS-mesh50"

# Compile protos inside the venv
python3 -m grpc_tools.protoc \
    -Iproto \
    --python_out=grpc_node \
    --grpc_python_out=grpc_node \
    proto/sync.proto

GRPC_LOG="$BASE_DIR/grpc.log"
BUS_LOG="$BASE_DIR/bus.log"
WEB_LOG="$BASE_DIR/web.log"

# 5. Ignite Core Services
echo "🚀 [OOB] Igniting Mesh Swarm..."
pkill -f python3 || true
sleep 1

nohup python3 -u "$BASE_DIR/grpc_node/grpc_server.py" > "$GRPC_LOG" 2>&1 &
nohup python3 -u "$BASE_DIR/memory_bus/server.py" > "$BUS_LOG" 2>&1 &
nohup python3 -u "$BASE_DIR/grpc_node/web_server.py" > "$WEB_LOG" 2>&1 &

echo "✅ [OOB] Swarm activated. Listening on 1111, 11111, 8080."
