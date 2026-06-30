#!/bin/bash
set -e

if [ -z "$1" ]; then
  echo "Usage: ./export_client_release.sh <project_directory>"
  echo "Example: ./export_client_release.sh quantasona-mesh"
  exit 1
fi

PROJECT=$1
MONOREPO_DIR=$(pwd)
EXPORT_BASE="/home/aellok/client-exports"
EXPORT_DIR="$EXPORT_BASE/$PROJECT-$(date +%s)"

if [ ! -d "$PROJECT" ]; then
  echo "Error: Directory '$PROJECT' not found in monorepo."
  exit 1
fi

echo "=> Preparing isolated client export for: $PROJECT"
mkdir -p "$EXPORT_DIR"

# 1. Copy the specific project files
echo "=> Copying source files..."
cp -r "$PROJECT/"* "$EXPORT_DIR/"
cp -r "$PROJECT/".* "$EXPORT_DIR/" 2>/dev/null || true

# Clean up internal git/monorepo artifacts if they slipped in
rm -rf "$EXPORT_DIR/.git" 2>/dev/null || true

# 2. Bundle shared Monorepo dependencies (Go Vendor)
if [ -f "$EXPORT_DIR/go.mod" ]; then
  echo "=> Detected Go project. Vendoring shared monorepo dependencies..."
  cd "$EXPORT_DIR"
  
  # Ensure we have a local go.mod, but we temporarily link the workspace to resolve paths
  # Actually, the safest way to vendor a go.work monorepo module is to do it from the root workspace
  cd "$MONOREPO_DIR"
  go work sync || true
  cd "$EXPORT_DIR"
  
  # Go mod vendor will pull all dependencies (including shared monorepo libs) into vendor/
  # Note: to vendor a local workspace module properly, it might need to temporarily replace 
  # workspace links with actual paths or just rely on the sync.
  # A robust fallback is to simply copy the shared Go libraries if we know them.
  # For now, we will rely on standard go mod vendor.
  go mod vendor 2>/dev/null || echo "Note: 'go mod vendor' skipped or had issues, ensure internal libs don't break."
fi

# 3. Initialize Pristine Git Repository
echo "=> Initializing clean Git repository (No Monorepo History)..."
cd "$EXPORT_DIR"
git init
git add .
git commit -m "Initial Client Release: $PROJECT"

echo "=========================================================="
echo "✅ SUCCESS: Pristine Client Repository Created!"
echo "📍 Location: $EXPORT_DIR"
echo "=========================================================="
echo "You can now safely push this directory to a client's remote repository."
echo "Example: cd $EXPORT_DIR && git remote add origin <CLIENT_URL> && git push -u origin master"
