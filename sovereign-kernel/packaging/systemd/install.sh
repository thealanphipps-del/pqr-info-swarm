#!/bin/bash
set -e

echo "=> Installing Sovereign Kernel (OS-SPARK) Systemd Service..."

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root (sudo)"
  exit 1
fi

mkdir -p /opt/sovereign-kernel
mkdir -p /etc/sovereign-kernel

# Copy binary
cp ../../sovereign-kernel /opt/sovereign-kernel/

# Copy service
cp sovereign-kernel.service /etc/systemd/system/

systemctl daemon-reload
systemctl enable sovereign-kernel
systemctl start sovereign-kernel

echo "=> OS-SPARK Kernel installed and running via Systemd."
systemctl status sovereign-kernel --no-pager
