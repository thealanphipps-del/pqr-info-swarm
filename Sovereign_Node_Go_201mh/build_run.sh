#!/bin/bash
set -e
APP_DIR="/data/data/com.termux/files/home/Sovereign_Node_Go"
cd "$APP_DIR"
echo "[*] Compiling Sovereign Mother App (CGO_ENABLED=0)..."
CGO_ENABLED=0 /usr/bin/go build -o mother_main .
echo "[*] Launching Binary..."
./mother_main
