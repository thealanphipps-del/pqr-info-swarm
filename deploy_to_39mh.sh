#!/bin/bash
# ==============================================================================
# SOVEREIGN MONOREPO: REMOTE SSH DEPLOYMENT (39.mh)
# Bypasses local Docker and Cloud Build entirely by deploying directly to the
# user's running compute instance.
# ==============================================================================
set -e

REMOTE_HOST="gemini-masterer@39.mh"
REPO_URL="https://github.com/thealanphipps-del/pqr-info-swarm.git"
REMOTE_DIR="/home/gemini-masterer/pqr-info-swarm"
TARGET_PORT="${TARGET_PORT:-8196}"

CYAN="\033[96m"
GREEN="\033[92m"
GOLD="\033[93m"
RESET="\033[0m"
BOLD="\033[1m"

echo -e "${GOLD}============================================================${RESET}"
echo -e "${BOLD}  SOVEREIGN MONOREPO: REMOTE BARE-METAL DEPLOYMENT${RESET}"
echo -e "${GOLD}============================================================${RESET}"
echo -e "${CYAN}Target Host:${RESET} ${REMOTE_HOST}"
echo -e "${CYAN}Port:${RESET}        ${TARGET_PORT}"
echo ""

echo -e "${GOLD}[PHASE 1] Pulling source code & building on 39.mh...${RESET}"
ssh -A "${REMOTE_HOST}" << EOF
    set -e
    # Cleanup previous deployment
    if [ -d "${REMOTE_DIR}" ]; then
        echo "Removing old source directory..."
        rm -rf "${REMOTE_DIR}"
    fi

    # Clone fresh repository
    echo "Cloning Monorepo..."
    git clone "${REPO_URL}" "${REMOTE_DIR}"

    # Build the server
    cd "${REMOTE_DIR}"
    echo "Downloading Go modules..."
    go mod tidy
    echo "Compiling PQR Server..."
    go build -o pqr-server ./cmd/pqr/main.go
EOF

echo -e "${GREEN}[OK] Build successful.${RESET}"

echo -e "${GOLD}[PHASE 2] Restarting PQR Server Daemon...${RESET}"
ssh -A "${REMOTE_HOST}" << EOF
    set -e
    cd "${REMOTE_DIR}"
    # Kill any existing instance running on the target port (ignoring errors if none exist)
    kill \$(lsof -t -i:${TARGET_PORT}) 2>/dev/null || true
    
    # Start the server in the background using nohup
    export PORT=${TARGET_PORT}
    export DATABASE_URL=postgresql://root@localhost:26257/antigravity?sslmode=disable
    
    nohup ./pqr-server > pqr-server.log 2>&1 &
    
    echo "Server launched! PID: \$!"
EOF

echo ""
echo -e "${GREEN}${BOLD}============================================================${RESET}"
echo -e "${GREEN}${BOLD}  ✅ BARE-METAL SSH DEPLOYMENT SUCCESSFUL!${RESET}"
echo -e "${GREEN}${BOLD}============================================================${RESET}"
echo -e "${CYAN}Live Service URL:${RESET} http://39.mh:${TARGET_PORT}"
echo -e "${CYAN}Remote Log File:${RESET}  ${REMOTE_DIR}/pqr-server.log"
echo ""
echo -e "${GOLD}To tail the live logs:${RESET}"
echo -e "  ssh ${REMOTE_HOST} 'tail -f ${REMOTE_DIR}/pqr-server.log'"
echo ""
