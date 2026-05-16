# 🚨 Gemini Emergency Bridge Protocol

The Emergency Bridge is a "Break-Glass" management channel that bypasses the primary Cloudflare Tunnel and UI, allowing for direct command injection from Gemini.

## 🔑 Authentication
The bridge is secured via the **X-Gemini-Key**.
- **Key**: `AIzaSyCqMMdPm1s6MuXy06yiWUlIQ0CJ1C-rPWk`
- **Location**: Stored in Vault at `secret/pqr/emergency_bridge`.

## 🛠️ Execution
Commands are sent via POST to `/REST/2.0/emergency/bridge`.

### Example: System Health Check
```bash
curl -X POST https://pqr.info/REST/2.0/emergency/bridge \
     -H "X-Gemini-Key: <YOUR_KEY>" \
     -H "Content-Type: application/json" \
     -d '{"command": "GET_SYSTEM_HEALTH"}'
```

### Supported Commands
- `GET_SYSTEM_HEALTH`: Returns node status, version, and uptime.
- `LIST_RECENT_TICKETS`: Retrieves the last 10 entries from the Ticketing Fabric.
- `TRIGGER_HEALING`: Manually initiates a healing loop for a specific issue.

## 🧠 Proactive Heartbeat
The Swarm maintains a permanent heartbeat to Gemini. Every 5 minutes, the node sends a **Sovereign Snapshot** to Gemini Pro for strategic review. If Gemini detects a critical failure, it can respond with an autonomous override command via this bridge.
