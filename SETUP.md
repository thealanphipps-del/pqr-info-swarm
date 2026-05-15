# PQR Ticketing System - Setup & Deployment Guide

## Quick Start

### Prerequisites
- CockroachDB v23.1+
- Go 1.21+
- curl or Postman (for testing)

### Step 1: Start CockroachDB

```powershell
# Windows
cd "C:\Users\drphi\cockroach-v23.1.13.windows-6.2-amd64"
.\cockroach.exe start-single-node --insecure

# Or with data directory:
.\cockroach.exe start-single-node --insecure --store=type=mem,size=0.25

# The output will show:
# CockroachDB node starting at 2024-05-14T00:00:00Z
# status: initialized new cluster
# SQL address: localhost:26257
```

### Step 2: Create Database (if not auto-created)

```bash
# Connect to CockroachDB SQL console
cockroach sql --insecure

# Create database:
CREATE DATABASE IF NOT EXISTS antigravity;
\q
```

### Step 3: Set Database URL

**Windows PowerShell:**
```powershell
$env:DATABASE_URL = "postgresql://root@localhost:26257/antigravity?sslmode=disable"

# Or add to profile for persistence:
Add-Content $PROFILE "`n`$env:DATABASE_URL = 'postgresql://root@localhost:26257/antigravity?sslmode=disable'"
```

**Windows Command Prompt:**
```cmd
set DATABASE_URL=postgresql://root@localhost:26257/antigravity?sslmode=disable
```

**Bash/Linux:**
```bash
export DATABASE_URL="postgresql://root@localhost:26257/antigravity?sslmode=disable"
```

### Step 4: Build and Run Server

```bash
cd c:\Users\drphi\pqr-info-swarm\cmd\\pqr
go build -o pqr.exe
.\pqr.exe

# Expected output:
# ✓ Database schema initialized
# ✓ Agent memory tables ready
# Starting PQR REST 2.0 API Server on :8080...
# Endpoints:
#   POST   /REST/2.0/ticket
#   ... (rest of endpoints)
```

### Step 5: Test the System

**PowerShell Test:**
```powershell
cd c:\Users\drphi\pqr-info-swarm
.\test-agent-memory.ps1 -BaseUrl http://localhost:8080 -AgentId test-agent-001
```

**Bash Test:**
```bash
cd c:\Users\drphi\pqr-info-swarm
bash test-agent-memory.sh http://localhost:8080 test-agent-001
```

**Simple curl Test:**
```bash
curl -X GET http://localhost:8080/REST/2.0/health
# Should return: {"status":"healthy","service":"PQR-ticketing"}
```

## Agents Onboarding

### For Go Agents

```go
package main

import (
	"context"
	"log"
	"github.com/thealanphipps-del/pqr"
)

func main() {
	session := pqr.NewAgentSession("http://localhost:8080", "my-agent-id")
	ctx := context.Background()

	// Create a memory ticket
	ticket, err := session.CreateMemory(ctx, "My Task", map[string]interface{}{
		"goal": "process data",
		"status": "started",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Later, recall the memory
	memory, _ := session.RecallMemory(ctx, ticket)
	log.Printf("Memory: %v", memory)
}
```

### For Python Agents

```python
import requests
import json
import uuid

class PQRClient:
    def __init__(self, base_url, agent_id):
        self.base_url = base_url
        self.agent_id = agent_id
        self.session = requests.Session()
        
    def create_ticket(self, subject, content, intent=None):
        """Create a memory ticket"""
        payload = {
            "Subject": subject,
            "Queue": "default",
            "Text": content,
            "AgentID": self.agent_id,
            "Layer": 2,
            "Intent": intent or {"agent": self.agent_id}
        }
        resp = self.session.post(f"{self.base_url}/REST/2.0/ticket", json=payload)
        return resp.json()["id"]
    
    def store_memory(self, ticket_id, memory_type, data, relevance=0.95):
        """Store agent context/memory"""
        payload = {
            "memory_type": memory_type,
            "data": data,
            "relevance_score": relevance
        }
        self.session.post(
            f"{self.base_url}/REST/2.0/agent/{self.agent_id}/memory/{ticket_id}",
            json=payload
        )
    
    def get_memory(self, ticket_id, memory_type="context"):
        """Retrieve stored memory"""
        resp = self.session.get(
            f"{self.base_url}/REST/2.0/agent/{self.agent_id}/memory/{ticket_id}",
            params={"type": memory_type}
        )
        return resp.json()

# Usage
client = PQRClient("http://localhost:8080", "python-agent-001")
ticket = client.create_ticket("Task", "Process data")
client.store_memory(ticket, "context", {"status": "running"})
memory = client.get_memory(ticket)
print(memory)  # {"status": "running"}
```

### For Node.js/JavaScript Agents

```javascript
const http = require('http');

class PQRClient {
  constructor(baseUrl, agentId) {
    this.baseUrl = baseUrl;
    this.agentId = agentId;
  }

  async createTicket(subject, content, intent) {
    const payload = JSON.stringify({
      Subject: subject,
      Queue: 'default',
      Text: content,
      AgentID: this.agentId,
      Layer: 2,
      Intent: intent || { agent: this.agentId }
    });

    return fetch(`${this.baseUrl}/REST/2.0/ticket`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: payload
    }).then(r => r.json()).then(d => d.id);
  }

  async storeMemory(ticketId, memType, data, relevance = 0.95) {
    const payload = JSON.stringify({
      memory_type: memType,
      data: data,
      relevance_score: relevance
    });

    return fetch(
      `${this.baseUrl}/REST/2.0/agent/${this.agentId}/memory/${ticketId}`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: payload
      }
    );
  }

  async getMemory(ticketId, memType = 'context') {
    return fetch(
      `${this.baseUrl}/REST/2.0/agent/${this.agentId}/memory/${ticketId}?type=${memType}`
    ).then(r => r.json());
  }
}

// Usage
const client = new PQRClient('http://localhost:8080', 'js-agent-001');
const ticket = await client.createTicket('Task', 'Do something');
await client.storeMemory(ticket, 'context', { status: 'running' });
const memory = await client.getMemory(ticket);
console.log(memory);  // {status: 'running'}
```

## Database Verification

### Check CockroachDB Status

```bash
# Connect to admin console
cockroach demo --insecure

# Or via SQL client
cockroach sql --insecure --database=antigravity

# View tables:
SHOW TABLES;

# Check tickets:
SELECT COUNT(*) FROM tickets;

# Check agent memory:
SELECT agent_id, COUNT(*) as memory_count FROM agent_memory GROUP BY agent_id;
```

### View Data in CockroachDB UI

1. Open browser: http://localhost:8080/admin
2. Navigate to Databases → antigravity
3. View tables:
   - `tickets` - Core ticket data
   - `ticket_content` - Content and intent storage
   - `agent_memory` - Agent context storage
   - `ticket_relationships` - Linked tickets
   - `ticket_audit` - Change history

## Troubleshooting

### CockroachDB Won't Connect

```
Error: "failed to connect to cockroachdb"
Solution: Verify CockroachDB is running and DATABASE_URL is correct
  - Check port 26257 is accessible
  - Default: postgresql://root@localhost:26257/antigravity?sslmode=disable
```

### Schema Initialization Failed

```
Error: "CREATE TABLE" failed
Solutions:
  1. POST /REST/2.0/init endpoint to reinitialize
  2. Check database permissions for root user
  3. Verify database antigravity exists
```

### Agent Memory Not Persisting

```
Symptoms: Stored memory not retrieved later
Check:
  1. POST to /REST/2.0/agent/{agentID}/memory/{ticketID} returns 200
  2. GET from /REST/2.0/agent/{agentID}/memory/{ticketID} with correct memory_type
  3. Verify ticket exists: GET /REST/2.0/ticket/{ticketID}
```

### High Latency on Queries

```
Solutions:
  1. Limit relevance score queries - index automatically used
  2. Store memory in appropriate types (context, knowledge, etc.)
  3. Monitor ticket count - consider archiving old tickets
  4. Use CockroachDB monitoring: admin console at http://localhost:8080
```

## Production Deployment

### Docker Deployment

Create `Dockerfile`:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o PQR ./cmd/pqr

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/PQR /usr/local/bin/
EXPOSE 8080
CMD ["PQR"]
```

Build and run:
```bash
docker build -t PQR-ticketing .
docker run -e DATABASE_URL="postgresql://root@cockroachdb:26257/antigravity?sslmode=disable" \
           -p 8080:8080 PQR-ticketing
```

### Environment Configuration

- `DATABASE_URL`: CockroachDB connection string
- `PORT`: Server port (default: 8080)
- `GIN_MODE`: Set to "release" for production
- `LOG_LEVEL`: debug, info, warn, error

### Monitoring & Logging

```go
// Example agent logging
client := pqr.NewClient("http://localhost:8080")
healthy, err := client.Health(context.Background())
if !healthy {
  log.Printf("Service unhealthy: %v", err)
  // Alert/retry logic
}
```

## Next Steps

1. **Deploy additional agents** using the client libraries
2. **Set up monitoring** for ticket creation and retrieval rates
3. **Implement agent orchestration** on top of ticketing
4. **Add authentication** for multi-tenant scenarios
5. **Enable vector search** for intelligent memory retrieval

## Support

For issues or questions:
1. Check logs: `pqr.exe` output
2. Test endpoint: `GET /REST/2.0/health`
3. Verify DB: `cockroach sql --insecure`
4. Review API docs in README.md


