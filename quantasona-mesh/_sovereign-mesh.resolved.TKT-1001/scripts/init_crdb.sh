#!/bin/bash
# Initialize CockroachDB across WSL (Ubuntu) and Windows 11 host
# NOTE: This script assumes cockroach binary is in your PATH.

# 1. WSL Node (Primary)
mkdir -p ~/crdb_data/node1
cockroach start --insecure --store=~/crdb_data/node1 --listen-addr=localhost:26257 --http-addr=localhost:8080 --background

# 2. Join instruction for Windows host
echo "--- COPY THIS COMMAND TO WINDOWS CMD/POWERSHELL ---"
echo "mkdir C:\crdb_data\node2"
echo "cockroach start --insecure --store=C:\crdb_data\node2 --listen-addr=0.0.0.0:26258 --http-addr=0.0.0.0:8081 --join=127.0.0.1:26257"
echo "----------------------------------------------------"

# 3. Finalize Cluster
echo "Initializing cluster..."
cockroach init --insecure --host=localhost:26257
