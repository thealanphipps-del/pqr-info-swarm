# 🛡️ PQR GitHub App Integration (Sovereign Sentinel)

The GitHub App acts as the primary "Forensic Proxy" for the Sovereign Node, allowing autonomous agents to manage repository state with granular control.

## 🛰️ Integration Stack
- **Library**: `google/go-github/v60`
- **Webhooks**: Routed via `ngrok` for local node exposure.
- **Identity**: Custom Bot (PQR-Sentinel) with read/write permissions.

## 🤖 Automation Triggers
1. **Vitality Slope Alerts**:
   - The app monitors `Issue` and `Commit` frequency.
   - If activity drops below the "Vitality Threshold," it triggers a swarm-wide status check.
2. **Fatality Purge**:
   - Automated cleanup of orphaned branches or failed experiment logs.
3. **Audit Forensic Hub Sync**:
   - Every GitHub event (Push/PR) is automatically mirrored as a **Fabric Ticket** in CockroachDB to maintain zero-divergence between Git and the Mesh.

## 🛠️ MCP Deployment (`mcp_pro_deploy.sh`)
```bash
#!/bin/bash
# One-click deployment for Termux/GoReleaser
echo "🚀 Deploying Sovereign MCP Server..."
go build -o mcp-server ./cmd/mcp
goreleaser release --snapshot --clean
echo "✅ MCP Node Active."
```

## 📂 File Offloading (`offload_sort.sh`)
```bash
#!/bin/bash
# Index and move files > 100MB to Google Drive
find . -size +100M -exec echo "Moving {} to Sovereign Archive..." \;
# Placeholder for rclone/google-drive-upload logic
```
