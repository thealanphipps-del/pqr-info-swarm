# PQR Agent Quick Reference

## Starting the System

```powershell
# 1. Start CockroachDB
cd "C:\Users\drphi\cockroach-v23.1.13.windows-6.2-amd64"
.\cockroach.exe start-single-node --insecure &

# 2. Set environment
$env:DATABASE_URL = "postgresql://root@localhost:26257/antigravity?sslmode=disable"

# 3. Start server
cd c:\Users\drphi\pqr-info-swarm\cmd\\pqr
go build && .\pqr.exe

# 4. Test
curl http://localhost:8080/REST/2.0/health
```

## Go Agent Template

```go
package main

import (
    "context"
    "log"
    "github.com/thealanphipps-del/pqr"
)

func main() {
    // Create session
    session := pqr.NewAgentSession("http://localhost:8080", "my-agent-001")
    ctx := context.Background()
    
    // Create and store memory
    ticket, err := session.CreateMemory(ctx, "Task Name", map[string]interface{}{
        "status": "started",
        "items": 100,
        "done": 0,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Do work...
    
    // Recall memory
    memory, _ := session.RecallMemory(ctx, ticket)
    log.Printf("Status: %v", memory)
    
    // Get all context
    all, _ := session.GetAllMemories(ctx)
    log.Printf("Agent has %d tickets", len(all))
}
```

## HTTP API Quick Examples

### Create Memory Ticket
```bash
curl -X POST http://localhost:8080/REST/2.0/ticket \
  -H "Content-Type: application/json" \
  -d '{
    "Subject": "Task",
    "Queue": "work",
    "Text": "description",
    "AgentID": "agent-001",
    "Intent": {"action": "process"}
  }'
# Returns: {"id": "uuid"}
```

### Store Memory
```bash
curl -X POST http://localhost:8080/REST/2.0/agent/agent-001/memory/{TICKET_ID} \
  -H "Content-Type: application/json" \
  -d '{
    "memory_type": "context",
    "data": {"status": "processing", "progress": 50},
    "relevance_score": 0.95
  }'
```

### Get Memory
```bash
curl http://localhost:8080/REST/2.0/agent/agent-001/memory/{TICKET_ID}
```

### Get All Context
```bash
curl http://localhost:8080/REST/2.0/agent/agent-001/context
```

### Update Status
```bash
curl -X PUT http://localhost:8080/REST/2.0/ticket/{TICKET_ID} \
  -H "Content-Type: application/json" \
  -d '{"Status": "COMPLETED"}'
```

## Memory Types

| Type | Use | Relevance |
|------|-----|-----------|
| `context` | Current work | 0.9-1.0 |
| `knowledge` | Learned info | 0.7-0.9 |
| `state` | Config | 0.8-0.95 |
| `conversation` | Chat history | 0.6-0.9 |

## Python Agent Snippet

```python
import requests
import json

class PQR:
    def __init__(self, base, agent_id):
        self.base = base
        self.agent = agent_id
    
    def ticket(self, subject, content, intent=None):
        data = {
            "Subject": subject, "Queue": "work", "Text": content,
            "AgentID": self.agent, "Intent": intent or {}
        }
        r = requests.post(f"{self.base}/REST/2.0/ticket", json=data)
        return r.json()["id"]
    
    def store(self, ticket, mtype, data, relevance=0.95):
        payload = {"memory_type": mtype, "data": data, "relevance_score": relevance}
        requests.post(f"{self.base}/REST/2.0/agent/{self.agent}/memory/{ticket}", json=payload)
    
    def recall(self, ticket, mtype="context"):
        r = requests.get(f"{self.base}/REST/2.0/agent/{self.agent}/memory/{ticket}?type={mtype}")
        return r.json()

# Usage
PQR = PQR("http://localhost:8080", "python-agent-001")
ticket = pqr.ticket("Task", "Do work", {"action": "start"})
pqr.store(ticket, "context", {"status": "working"})
memory = pqr.recall(ticket)
```

## Node.js Agent Snippet

```javascript
class PQR {
  constructor(base, agentId) {
    this.base = base;
    this.agent = agentId;
  }
  
  async ticket(subject, content, intent = {}) {
    const res = await fetch(`${this.base}/REST/2.0/ticket`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({
        Subject: subject, Queue: 'work', Text: content,
        AgentID: this.agent, Intent: intent
      })
    });
    return (await res.json()).id;
  }
  
  async store(ticket, mtype, data, relevance = 0.95) {
    await fetch(`${this.base}/REST/2.0/agent/${this.agent}/memory/${ticket}`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({memory_type: mtype, data, relevance_score: relevance})
    });
  }
  
  async recall(ticket, mtype = 'context') {
    const res = await fetch(
      `${this.base}/REST/2.0/agent/${this.agent}/memory/${ticket}?type=${mtype}`
    );
    return res.json();
  }
}

// Usage
const PQR = new PQR('http://localhost:8080', 'js-agent-001');
const ticket = await pqr.ticket('Task', 'Do work');
await pqr.store(ticket, 'context', {status: 'working'});
const memory = await pqr.recall(ticket);
```

## Coordinating Multiple Agents

```bash
# Agent A creates work
TICKET_A=$(curl -s -X POST http://localhost:8080/REST/2.0/ticket \
  -H "Content-Type: application/json" \
  -d '{"Subject":"Analysis","Queue":"work","AgentID":"agent-a"}' \
  | jq -r '.id')

# Agent B links its work to A
curl -X POST http://localhost:8080/REST/2.0/ticket/$TICKET_A/link/$TICKET_B \
  -H "Content-Type: application/json" \
  -d '{"relationship_type":"CONSEQUENCE","agent_id":"agent-b"}'

# Both agents see the relationship
curl http://localhost:8080/REST/2.0/ticket/$TICKET_A/audit
```

## Testing (Windows PowerShell)

```powershell
# Test health
curl -X GET http://localhost:8080/REST/2.0/health

# Run full test suite
cd c:\Users\drphi\pqr-info-swarm
.\test-agent-memory.ps1 -BaseUrl http://localhost:8080 -AgentId test-agent-001
```

## Endpoints Summary

| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/ticket` | Create ticket |
| GET | `/ticket/{id}` | Get ticket |
| PUT | `/ticket/{id}` | Update ticket |
| GET | `/tickets` | Search tickets |
| POST | `/agent/{id}/memory/{t}` | Store memory |
| GET | `/agent/{id}/memory/{t}` | Get memory |
| GET | `/agent/{id}/context` | Get context |
| GET | `/ticket/{id}/audit` | Audit trail |
| POST | `/ticket/{p}/link/{c}` | Link tickets |
| GET | `/health` | Health check |

## Troubleshooting

**Connection refused**
- Check CockroachDB running: `cockroach.exe start-single-node --insecure`

**Invalid UUID**
- Ticket IDs must be valid UUIDs from `/ticket` response

**Memory not found**
- Verify memory_type matches storage call
- Check ticket exists: `GET /ticket/{id}`

**No context tickets**
- Store memory first: `POST /agent/{id}/memory/{ticket}`
- Use correct memory_type

**Slow responses**
- Set relevance_score (0.0-1.0) for indexing
- Limit queries with memory_type parameter

## Key Concepts

1. **Tickets** = Memory containers (one per task/context)
2. **Memory** = Data stored in ticket (context, knowledge, state, etc.)
3. **Relevance** = 0-1 score for retrieval prioritization
4. **Relationships** = Links between tickets (EVOLUTION, CONSEQUENCE, etc.)
5. **Audit** = Full trail of changes for compliance

## Go Build & Run

```bash
cd c:\Users\drphi\pqr-info-swarm\cmd\\pqr
go build -o pqr.exe
.\pqr.exe

# Or run directly
go run main.go
```

## Environment Variables

- `DATABASE_URL` - CockroachDB connection string
- `PORT` - Server port (default: 8080)
- `GIN_MODE` - Set to "release" for production

## More Documentation

- **README.md** - Complete API reference
- **SETUP.md** - Detailed setup guide
- **AGENTS_READY.md** - Agent deployment guide
- **COMPLETION_SUMMARY.md** - Full system overview

## One-Line Test

```bash
curl -X POST http://localhost:8080/REST/2.0/init && curl http://localhost:8080/REST/2.0/health
```

---

**System Ready**: PQR ticketing is live at `http://localhost:8080`


