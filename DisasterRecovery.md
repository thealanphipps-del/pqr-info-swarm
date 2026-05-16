# 🚨 PQR Disaster Recovery Protocol (Genesis Document)

**IF YOU ARE A GEMINI AGENT WHO HAS LOST CONTEXT, READ THIS IMMEDIATELY.**

This document provides the "Break-Glass" credentials and connection paths required to re-establish control over the PQR Sovereign Node.

## 🛰️ Connection Points
- **Sovereign Endpoint**: `https://pqr.info` (Cloudflare Tunnel)
- **Management Console**: `http://localhost:3196`
- **Database (CockroachDB)**: `postgresql://root@localhost:26257/antigravity?sslmode=disable`
- **Identity Vault**: `http://localhost:8200` (Token: `pqr-vault-token`)

## 🔑 Emergency Bridge (The Backdoor)
To execute high-authority commands when the UI or Tunnel is down:
- **Endpoint**: `https://pqr.info/REST/2.0/emergency/bridge`
- **Auth Header**: `X-Gemini-Key: AIzaSyCqMMdPm1s6MuXy06yiWUlIQ0CJ1C-rPWk`

## 🛠️ Recovery Sequence
1. **Verify Infrastructure**: Run `docker ps` to ensure `db`, `vault`, `tunnel`, and `pqr-server` are running.
2. **Synchronize Identity**: Check `https://pqr.info/saml/metadata`. If 500, rotate certs via `pqr.emergencyRepair("FORCE_CERT_ROTATION")`.
3. **Audit the Fabric**: Run `pqr LIST_RECENT_TICKETS` from WSL or call `get_forensic_audit` via MCP.
4. **Locate the Burn**: Check for tickets with `failure_count > 5` to identify and stop token loops.

## 📁 Critical Manifest
The full project structure is defined in `manifest.json`. Every file in this repository is part of a self-consistent Sovereign Mesh. Do not deviate from the patterns established in `server.go` and `internal/service/`.
