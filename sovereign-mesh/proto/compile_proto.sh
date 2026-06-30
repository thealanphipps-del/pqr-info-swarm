#!/bin/bash
echo -e "\033[93m[PROTO] Compiling sync.proto...\033[0m"

# Go generation
pushd proto > /dev/null
protoc -I. \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    sync.proto mesh_proto.proto
popd > /dev/null

# Removed legacy nested path hack - Go workspace uses source_relative

# Python generation
PYTHON_BIN="python3"
if [ -x "../.venv/bin/python3" ]; then
    PYTHON_BIN="../.venv/bin/python3"
elif [ -x ".venv/bin/python3" ]; then
    PYTHON_BIN=".venv/bin/python3"
fi

$PYTHON_BIN -m grpc_tools.protoc \
    -Iproto \
    --python_out=grpc_node \
    --grpc_python_out=grpc_node \
    proto/sync.proto proto/mesh_proto.proto

echo -e "\033[92m[PROTO] Compilation successful!\033[0m"
