#!/bin/bash
# Sovereign MCP One-Click Setup
echo "Initializing MCP Pro Deployment..."
go build -o pqr-mcp ./cmd/mcp
echo "Deploying to Termux environment..."
