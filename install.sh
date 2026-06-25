#!/bin/bash
set -e

echo "[SWEN Installer] Building SWEN executable..."
go build -o swend ./cmd/swend

echo "[SWEN Installer] Executing Genesis snapshot..."
./swend install

echo "[SWEN Installer] Installation complete! You can now start SWEN with './swend menu'."
