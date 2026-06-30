# Antigravity JSON Bridge

This project provides a lightweight, JSON-based communication layer for the **Sovereign AI Infrastructure**. It was designed as a "handy" client/server setup to facilitate communication between nodes (like the Alienware Aurora R9 and the Laptop node) even when full gRPC or Linux environments aren't yet active.

## Components

- **`client.py`**: A premium CLI client that loads a JSON payload (e.g., your TechProfile) and transmits it to a target node using the **Event Horizon** protocol handshake.
- **`server.py`**: A listener node that receives transmissions, validates the payload, and integrates the state.
- **`payloads/`**: Directory for storing JSON payloads.
- **`config.json`**: Configure target IP, Port, and Node ID here.

## Quick Start

### 1. Start the Node Listener
Run this on the machine intended to receive data:
```powershell
python server.py
```

### 2. Configure the Client
Edit `config.json` and set the `server_ip` to the IP address of the machine running the listener.

### 3. Execute Handshake
Run the client to send the profile:
```powershell
python client.py
```

## Protocol Details
- **Path**: `/bridge`
- **Method**: `POST`
- **Headers**:
  - `X-Antigravity-Node`: Identifies the source node.
  - `X-Protocol-Version`: 1.0
- **Aesthetics**: Full ANSI color support for Sovereign-themed console output.
