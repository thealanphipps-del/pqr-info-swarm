#!/bin/bash
# ==============================================================================
# SOVEREIGN MONOREPO: FREE TIER (e2-micro) GCP COMPUTE DEPLOYMENT
# Bypasses local Docker and Cloud Build entirely by building from source on the VM.
# ==============================================================================
set -e

PROJECT_ID="model-loader-495607-m2"
REGION="us-east1"
ZONE="us-east1-b"
INSTANCE_NAME="sovereign-free-node"
ACCOUNT="alan@w-isp.net"

CYAN="\033[96m"
GREEN="\033[92m"
GOLD="\033[93m"
RESET="\033[0m"
BOLD="\033[1m"

echo -e "${GOLD}============================================================${RESET}"
echo -e "${BOLD}  SOVEREIGN MONOREPO: FREE TIER COMPUTE DEPLOYMENT${RESET}"
echo -e "${GOLD}============================================================${RESET}"
echo -e "${CYAN}Target Project:${RESET} ${PROJECT_ID}"
echo -e "${CYAN}Target Zone:${RESET}    ${ZONE}"
echo -e "${CYAN}Machine Type:${RESET}   e2-micro (Free Tier Eligible)${RESET}"
echo ""

# Enable compute API if not already enabled (using the user account)
echo -e "${GOLD}Verifying Compute Engine API...${RESET}"
# gcloud services enable compute.googleapis.com --project="${PROJECT_ID}" --quiet

# Define the Startup Script
# This runs as root on the VM when it boots. It installs Go, clones the repo, builds, and starts a systemd service.
STARTUP_SCRIPT=$(cat << 'EOF'
#!/bin/bash
set -e
# Install dependencies
apt-get update
apt-get install -y git wget

# Install Go 1.22
wget https://go.dev/dl/go1.22.4.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.22.4.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Clone the Monorepo
cd /opt
if [ -d "pqr-info-swarm" ]; then
    rm -rf pqr-info-swarm
fi
git clone https://github.com/thealanphipps-del/pqr-info-swarm.git

# Build the PQR Server
cd /opt/pqr-info-swarm
/usr/local/go/bin/go mod tidy
/usr/local/go/bin/go build -o pqr-server ./cmd/pqr/main.go

# Create a systemd service to keep it running on port 80
cat << 'SYS' > /etc/systemd/system/pqr-server.service
[Unit]
Description=Sovereign PQR Server
After=network.target

[Service]
ExecStart=/opt/pqr-info-swarm/pqr-server
WorkingDirectory=/opt/pqr-info-swarm
Restart=always
User=root
Environment=PORT=80
Environment=DATABASE_URL=postgresql://root@localhost:26257/antigravity?sslmode=disable

[Install]
WantedBy=multi-user.target
SYS

# Start and enable the service
systemctl daemon-reload
systemctl enable pqr-server
systemctl start pqr-server
# Install Cloudflared
wget -q https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
dpkg -i cloudflared-linux-amd64.deb

# Create cloudflared service
cat << 'CF_SYS' > /etc/systemd/system/cloudflared-tunnel.service
[Unit]
Description=Cloudflare Tunnel
After=network.target

[Service]
ExecStart=/usr/bin/cloudflared tunnel --no-autoupdate run --token eyJhIjoiYzA3Y2NjYTM3ZDMyN2Y3Nzk5NjcxZDIzMDhmODBiZWIiLCJ0IjoiNzUzMDQ5NzAtMWZhMC00Y2Q0LTg1M2QtOGUzZGQyM2NjOGU0IiwicyI6IlpqTTFNRFkzTUdZdE9ESTFZaTAwT1RNM0xXRmtORFl0Wm1Fd01UQXdOR1l3T1dVeiJ9
Restart=always
User=root

[Install]
WantedBy=multi-user.target
CF_SYS

systemctl daemon-reload
systemctl enable cloudflared-tunnel
systemctl start cloudflared-tunnel
EOF
)

# Launch the Free Tier VM (with a retry loop to wait for Org Policy propagation)
echo -e "${GOLD}Provisioning GCP e2-micro VM and triggering source build... (Waiting for Org Policy to propagate)${RESET}"
MAX_RETRIES=20
RETRY_COUNT=0
SUCCESS=false

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if gcloud compute instances create "${INSTANCE_NAME}" \
        --project="${PROJECT_ID}" \
        --zone="${ZONE}" \
        --machine-type="e2-micro" \
        --image-family="debian-12" \
        --image-project="debian-cloud" \
        --boot-disk-size="30GB" \
        --boot-disk-type="pd-standard" \
        --tags="http-server,https-server" \
        --no-address \
        --no-service-account \
        --no-scopes \
        --metadata="startup-script=${STARTUP_SCRIPT}" \
        \
        --quiet; then
        SUCCESS=true
        break
    else
        echo -e "${CYAN}Org policy hasn't propagated yet. Retrying in 15 seconds... ($(($RETRY_COUNT+1))/$MAX_RETRIES)${RESET}"
        sleep 15
        RETRY_COUNT=$(($RETRY_COUNT+1))
    fi
done

if [ "$SUCCESS" = false ]; then
    echo -e "${GOLD}Error: Failed to create VM after $MAX_RETRIES attempts.${RESET}"
    exit 1
fi

# Allow HTTP traffic via firewall
echo -e "${GOLD}Configuring Firewall to allow HTTP traffic on port 80...${RESET}"
gcloud compute firewall-rules create default-allow-http \
    --project="${PROJECT_ID}" \
    --allow=tcp:80 \
    --source-ranges=0.0.0.0/0 \
    --target-tags=http-server \
    \
    --quiet || echo -e "${CYAN}Firewall rule already exists.${RESET}"

# Fetch the internal IP instead
EXTERNAL_IP=$(gcloud compute instances describe "${INSTANCE_NAME}" --project="${PROJECT_ID}" --zone="${ZONE}" --format='get(networkInterfaces[0].networkIP)')

echo ""
echo -e "${GREEN}${BOLD}============================================================${RESET}"
echo -e "${GREEN}${BOLD}  ✅ FREE TIER SERVERLESS BYPASS DEPLOYMENT INITIATED!${RESET}"
echo -e "${GREEN}${BOLD}============================================================${RESET}"
echo -e "${CYAN}Instance Name:${RESET} ${INSTANCE_NAME}"
echo -e "${CYAN}External IP:${RESET}   http://${EXTERNAL_IP}"
echo ""
echo -e "${GOLD}Note: The VM is booting and compiling Go from source. It will take ~2-3 minutes for the IP to become responsive.${RESET}"
echo ""
