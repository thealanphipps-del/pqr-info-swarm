#!/bin/bash

# Sovereign Node: Automated File Offloading Script
# Target: Files > 100MB to Google Drive

SOURCE_DIR="./data"
DEST_DIR="/sdcard/GoogleDrive/Sovereign_Offload"

echo "[INFO] Scanning for files larger than 100MB in $SOURCE_DIR..."

mkdir -p "$DEST_DIR"

find "$SOURCE_DIR" -type f -size +100M | while read -r file; do
    echo "[OFFLOAD] Moving $file to $DEST_DIR"
    mv "$file" "$DEST_DIR/"
done

echo "[SUCCESS] Offloading complete."
