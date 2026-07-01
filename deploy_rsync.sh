#!/bin/bash
set -e

REMOTE_HOST=$1
if [ -z "$REMOTE_HOST" ]; then
    echo "Usage: $0 <user@host>"
    exit 1
fi

TARGET_PORT="8196"

echo "Deploying to $REMOTE_HOST ..."
echo "Syncing local Monorepo via rsync..."
rsync -avz --exclude '.git' --exclude 'node_modules' --exclude '.next' --exclude 'dist' -e "ssh -o BatchMode=yes -i ~/.ssh/id_ed25519" /home/aellok/pqr-info-swarm/ ${REMOTE_HOST}:/home/aellok/pqr-info-swarm/

ssh -o BatchMode=yes -i ~/.ssh/id_ed25519 -A "${REMOTE_HOST}" << EOF
set -e
export PATH=\$PATH:/home/aellok/go/bin
if [ ! -d "/home/aellok/go/bin" ]; then
    echo "Installing Go locally..."
    cd /home/aellok
    wget -q https://go.dev/dl/go1.22.4.linux-amd64.tar.gz
    mkdir -p /home/aellok/go
    tar -C /home/aellok/go -xzf go1.22.4.linux-amd64.tar.gz --strip-components=1
    rm go1.22.4.linux-amd64.tar.gz
fi
cd /home/aellok/pqr-info-swarm
go mod tidy
go build -o pqr-server ./cmd/pqr/main.go
kill \$(lsof -t -i:${TARGET_PORT}) 2>/dev/null || true
export PORT=${TARGET_PORT}
export DATABASE_URL=postgresql://root@localhost:26257/antigravity?sslmode=disable
nohup ./pqr-server > pqr-server.log 2>&1 &
echo "Deployment successful on $REMOTE_HOST"
EOF
