#!/bin/bash
# ==============================================================================
# DEPLOYMENT INTEGRATION TEST SUITE
# Validates the abstracted GCP deploy scripts in a dry-run / mocked environment.
# ==============================================================================
set -e

echo "Running Deployment Integration Tests..."

# Check if the primary deployment script is executable
DEPLOY_SCRIPT="./quantasona-mesh/_sovereign-mesh.resolved.TKT-1001/gcp_deploy_cloud_run.sh"

if [ ! -f "$DEPLOY_SCRIPT" ]; then
    echo "ERROR: Deployment script not found!"
    exit 1
fi

chmod +x "$DEPLOY_SCRIPT"

# Mock the gcloud command for the test suite
export PATH="$(pwd)/ci/mock_bin:$PATH"
mkdir -p ./ci/mock_bin
cat << 'EOF' > ./ci/mock_bin/gcloud
#!/bin/bash
echo "[MOCK GCLOUD] Execution: gcloud $@"
if [[ "$*" == *"run deploy"* ]]; then
    echo "Deploying to mocked service..."
elif [[ "$*" == *"run services describe"* ]]; then
    echo "https://mock-service-url.run.app"
fi
exit 0
EOF
chmod +x ./ci/mock_bin/gcloud

cat << 'EOF' > ./ci/mock_bin/docker
#!/bin/bash
echo "[MOCK DOCKER] Execution: docker $@"
exit 0
EOF
chmod +x ./ci/mock_bin/docker

echo "Executing deployment script in test environment..."
# Run the deployment script and capture output
output=$($DEPLOY_SCRIPT 2>&1) || (echo "Deployment script failed in CI." && exit 1)

# Validate Versioning logic was used
if echo "$output" | grep -q "Build Tag:      v-"; then
    echo "✅ Versioning logic verified."
else
    echo "❌ Versioning logic failed!"
    exit 1
fi

if echo "$output" | grep -q "\[MOCK GCLOUD\] Execution: gcloud beta run deploy pqr-server-pool"; then
    echo "✅ Abstraction and deployment execution verified."
else
    echo "❌ Abstraction execution failed!"
    exit 1
fi

echo "========================================="
echo "✅ DEPLOYMENT INTEGRATION TESTS PASSED"
echo "========================================="

# Clean up
rm -rf ./ci/mock_bin
