#!/bin/bash
# Sovereign Mesh: Secure Transport Layer Bridge
# Establishes persistent tunnels through 39.MH Sentry Host

JUMP_HOST="39.mh"
# Ports: 1111 (gRPC), 11111 (RAM Bus), 8080 (Web)
LOCAL_PORTS=(1111 11111 8085)

echo -e "\033[96m=== MESH TRANSPORT LAYER INITIALIZATION ===\033[0m"

for port in "${LOCAL_PORTS[@]}"; do
    if ! ss -tuln | grep -q ":$port "; then
        echo -e "🔗 Bridging Port $port through $JUMP_HOST..."
        ssh -f -N -L $port:localhost:$port $JUMP_HOST
    else
        echo -e "✅ Port $port already bridged."
    fi
done

echo -e "\033[92mMesh transport layer secured.\033[0m"
