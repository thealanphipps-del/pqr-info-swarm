# Codebase Architecture

## Design Patterns & Core Subsystems

### 1. Zero-CGO Static Portability Constraint
All Go packages are compiled using `CGO_ENABLED=0` to guarantee absolute portability across Linux environments (including Termux). This requires strict avoidance of CGO dependencies (e.g. using pure-Go SQLite replacements instead of standard SQLite packages).

### 2. Lock-free Atomic Reasoning Queues
To achieve maximum concurrent throughput across the 8x8 Strategy Swarm, task nodes and ticket transitions leverage lock-free queues that avoid memory allocations and prevent garbage collection pressure.

### 3. Dynamic Database Routing Fallback
All ticketing modules (`pkg/tickets`) and fabric modules (`pkg/rtgo`) monitor PostgreSQL connections with a tight timeout limit. If connection fails, operations are dynamically translated and executed against a local pure-Go SQLite fallback engine (`local_ledger.db`), ensuring zero downtime.
