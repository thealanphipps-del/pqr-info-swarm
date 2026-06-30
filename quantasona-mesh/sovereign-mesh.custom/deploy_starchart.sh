#!/bin/bash
TARGET_DIR="/home/billing/sovereign_mesh"
mkdir -p $TARGET_DIR/proto $TARGET_DIR/grpc_node

# Decode and write files
echo $P1 | base64 -d > $TARGET_DIR/proto/sync.proto
echo $P2 | base64 -d > $TARGET_DIR/grpc_node/grpc_server.py
echo $P3 | base64 -d > $TARGET_DIR/grpc.go
echo $P4 | base64 -d > $TARGET_DIR/types.go
echo $P5 | base64 -d > $TARGET_DIR/sovereign.go
echo $P6 | base64 -d > $TARGET_DIR/radius.go
echo $P7 | base64 -d > $TARGET_DIR/grpc_node/web_server.py
echo $P8 | base64 -d > $TARGET_DIR/grpc_node/index.html
echo $P9 | base64 -d > $TARGET_DIR/grpc_node/mgsh_mcp.py

# Recompile and Restart
cd $TARGET_DIR
bash proto/compile_proto.sh
pkill -f python3
sleep 2
export ANTIGRAVITY_NODE_ID="GCP-OPS-mesh50"
nohup python3 -u grpc_node/grpc_server.py > ./grpc.log 2>&1 &
nohup python3 -u memory_bus/server.py > ./bus.log 2>&1 &
nohup python3 -u grpc_node/web_server.py > ./web.log 2>&1 &
