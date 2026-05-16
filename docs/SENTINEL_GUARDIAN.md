# 🛡️ PQR Sentinel: Sovereign Guardian

The **Sentinel** is a host-side guardian agent running on Windows that ensures the PQR Sovereign Mesh remains operational, even if the Docker Engine or specific containers fail.

## 🛰️ Core Functions
1.  **Engine Monitoring**: Periodically checks if the Docker Engine is responsive.
2.  **Health Verification**: Polls the PQR REST 2.0 API health endpoint (`/health`).
3.  **Auto-Recovery**: Automatically restarts the `pqr-server` if it becomes unreachable.
4.  **Agent-Driven Signal Handling**: Watches for "triggers" from within the mesh to execute host-side operations (like a full rebuild).

## 🚀 Deployment
To start the Sentinel, run the following in a Windows PowerShell terminal:

```powershell
.\SENTINEL.ps1
```

It is recommended to keep this terminal visible or run it as a background scheduled task.

## 🧬 Inter-Process Communication (Agent-to-Host)
Agents inside the Docker containers can "signal" the host-side Sentinel by interacting with the shared `signals/` directory.

### Triggering a Full Rebuild
If an agent determines that the entire stack needs a hard reset (e.g., after a major code change or critical failure), it can create a file:

**Inside Container:**
```bash
touch /app/signals/RESTART_TRIGGER
```

**Sentinel Reaction:**
1. Detects the file.
2. Logs the request.
3. Executes `docker-compose up -d --build`.
4. Clears the trigger.

## 📋 Monitoring
The Sentinel logs all its observations and recovery actions to:
`C:\Users\drphi\pqr-info-swarm\sentinel.log`

---
**Status**: 🟢 Active
**Role**: Host-side Watchdog
**System**: PQR Sovereign Mesh
