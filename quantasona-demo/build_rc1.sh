#!/bin/bash
set -e

echo "Building RC1 for Quantasona App..."

# 1. Compile the Go backend
echo "Compiling Go backend..."
cd /home/aellok/pqr-info-swarm/quantasona-demo/backend
go build -o bin/backend ./cmd/backend
echo "Go backend compiled successfully."

# 2. Start the Python LLM integration server
echo "Starting Python LLM integration server..."
cd /home/aellok/pqr-info-swarm/quantasona-demo/llm_integration
source .venv/bin/activate
nohup python main.py > server.log 2>&1 &
LLM_PID=$!
echo $LLM_PID > llm_server.pid
echo "Python LLM integration server started with PID $LLM_PID."

# 3. Build the Android APK
echo "Building Android APK on Windows Host via PowerShell Interop..."
powershell.exe -Command "cd \\wsl$\Ubuntu\home\aellok\pqr-info-swarm\quantasona-demo\quantasona-android; .\gradlew.bat assembleDebug"
echo "Android APK built successfully on Windows."

echo "RC1 build and startup completed successfully!"
