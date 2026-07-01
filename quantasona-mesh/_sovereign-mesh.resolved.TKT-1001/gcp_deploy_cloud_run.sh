#!/bin/bash
# ==============================================================================
# SOVEREIGN MONOREPO: GOOGLE CLOUD RUN DEPLOYMENT SCRIPT (REFACTORED)
# ==============================================================================
set -e

# --- CONFIG (Abstracted Targets) ---
# Default to environment variables, fallback to hardcoded for the current target
PROJECT_ID="${GCP_PROJECT_ID:-model-loader-495607-m2}"
REGION="${GCP_REGION:-us-central1}"
REPOSITORY="${GCP_ARTIFACT_REPO:-sovereign-neural-repo}"
SERVICE_NAME="${GCP_SERVICE_NAME:-pqr-server-pool}"
TARGET_PORT="${TARGET_PORT:-8196}"
GPU_TYPE="${GPU_TYPE:-nvidia-l4}"
GPU_COUNT="${GPU_COUNT:-0}" # Defaulting to 0 since PQR Server doesn't strictly need it, but configurable

# --- VERSIONING INTEGRATION ---
GIT_SHA=$(git rev-parse --short HEAD 2>/dev/null || echo "manual")
BUILD_TIMESTAMP=$(date +%Y%m%d-%H%M%S)
TAG="v-${GIT_SHA}-${BUILD_TIMESTAMP}"
IMAGE_NAME="pqr-server-monorepo"
FULL_IMAGE_URI="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${IMAGE_NAME}:${TAG}"

CYAN="\033[96m"
GREEN="\033[92m"
GOLD="\033[93m"
RESET="\033[0m"
BOLD="\033[1m"

echo -e "${GOLD}============================================================${RESET}"
echo -e "${BOLD}  SOVEREIGN MONOREPO: CLOUD RUN PIPELINE${RESET}"
echo -e "${GOLD}============================================================${RESET}"
echo -e "${CYAN}Target Project:${RESET} ${PROJECT_ID}"
echo -e "${CYAN}Build Tag:${RESET}      ${TAG}"
echo -e "${CYAN}Image URI:${RESET}      ${FULL_IMAGE_URI}"
echo -e "${CYAN}Service:${RESET}        ${SERVICE_NAME}"
echo ""

# --- PHASE 1: BUILD ARTIFACT ---
echo -e "${GOLD}[PHASE 1] Building Artifact...${RESET}"
# Ensure we are in the root of the monorepo
cd /home/aellok/pqr-info-swarm || exit 1

# APIs should already be enabled for this project

# Build and push the container natively in GCP Cloud Build (bypassing local Docker)
gcloud builds submit --tag "${FULL_IMAGE_URI}" . --project="${PROJECT_ID}"
echo -e "${GREEN}[OK] Cloud Build complete and image pushed to Artifact Registry.${RESET}"

# --- PHASE 2: DEPLOY MANIFEST ---
echo -e "${GOLD}[PHASE 2] Deploying to Cloud Run...${RESET}"

# Configure GPU flags only if count > 0
GPU_ARGS=""
if [ "${GPU_COUNT}" -gt 0 ]; then
    GPU_ARGS="--gpu=${GPU_COUNT} --gpu-type=${GPU_TYPE}"
fi

# Execute Deployment
gcloud beta run deploy "${SERVICE_NAME}" \
    --image="${FULL_IMAGE_URI}" \
    --region="${REGION}" \
    --project="${PROJECT_ID}" \
    --cpu=2 \
    --memory=4Gi \
    $GPU_ARGS \
    --min-instances=0 \
    --max-instances=5 \
    --port="${TARGET_PORT}" \
    --no-allow-unauthenticated \
    --set-env-vars="GIT_COMMIT=${GIT_SHA},DEPLOYMENT_TIME=${BUILD_TIMESTAMP}" \
    --quiet

# --- PHASE 3: VERIFICATION ---
DEPLOY_URL=$(gcloud run services describe "${SERVICE_NAME}" --region="${REGION}" --project="${PROJECT_ID}" --format="value(status.url)")

echo ""
echo -e "${GREEN}${BOLD}============================================================${RESET}"
echo -e "${GREEN}${BOLD}  ✅ MONOREPO DEPLOYMENT SUCCESSFUL!${RESET}"
echo -e "${GREEN}${BOLD}============================================================${RESET}"
echo -e "${CYAN}Service URL:${RESET}   ${DEPLOY_URL}"
echo -e "${CYAN}Active Tag:${RESET}    ${TAG}"
echo -e "${CYAN}Root Image:${RESET}    ${FULL_IMAGE_URI}"
echo ""
echo -e "${GOLD}To view live logs and verify traffic routing:${RESET}"
echo -e "  gcloud run services describe ${SERVICE_NAME} --region=${REGION} --project=${PROJECT_ID}"
echo ""
