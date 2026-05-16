# 🛰️ Swarm Operations Guide

Practical procedures for managing and monitoring the PQR Sovereign Node.

## 🚀 Lifecycle Management
The entire mesh is managed via Docker Compose.

### Start the Node
```powershell
.\start_pqr.ps1
```

### Stop the Node
```powershell
docker-compose down
```

### Restart a Specific Service
```bash
docker-compose restart pqr-server
docker-compose restart tunnel
```

## 🔍 Forensic Monitoring
Monitoring is the key to understanding agent evolution.

### View Real-time Logs
```bash
# All services
docker-compose logs -f

# Just the server
docker logs -f pqr-info-swarm-pqr-server-1

# Just the tunnel
docker logs -f pqr-info-swarm-tunnel-1
```

### Check Database Health (Console)
Visit [http://localhost:8081](http://localhost:8081) to view the CockroachDB console.

## 🧬 Autonomous Healing Verification
To see what the agents are currently "thinking":
1. Open the [HUD](http://localhost:3196/hud).
2. Look for tickets with Layer 7 or 10.
3. Check the `IntentBlob` for the reasoning behind an autonomous resolution.

## 🧪 Testing Connectivity
Verify the tunnel and access bypass headers:
```powershell
$headers = @{ 
  "CF-Access-Client-Id" = "c98ca7026f54305b05cd24975a3ce6d2.access";
  "CF-Access-Client-Secret" = "ebf3177d992adb0c3db7b088fb5b9e3d83e96649fb9bc5b86a25301af5c8e744"
}
Invoke-RestMethod -Uri "https://pqr.info/REST/2.0/health" -Headers $headers
```
