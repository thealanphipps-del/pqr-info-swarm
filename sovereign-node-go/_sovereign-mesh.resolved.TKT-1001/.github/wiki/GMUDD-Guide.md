# GMUDD Guide – Self‑Guided Documentation

## Overview
`gmudd` (the **G**lobal **M**UDD interface) provides an interactive, color‑coded dashboard to explore the Sovereign Swarm mesh, manage agents, and now **integrates ticketing** directly from the UI.

---
### Menu Summary (current version)
| Key | Action |
|-----|--------|
| **1** | Look Around – diagnostic audit of the current node |
| **2** | Move – teleport to a directly connected tunnel |
| **3** | gRPC Control Ping – verify the mesh control plane |
| **4** | Audit Swarm Ancestry (Pedigree) |
| **5** | Query Swarm Immutable Ledger |
| **6** | Read Protocol Grimoire (docs) |
| **7** | Save Swarm Session State |
| **8** | View Agentic Memories (RT/PQR DB) |
| **9** | **Create Incomplete Task Tickets** – scans all code and `.md` files in the repo for `- [ ]` items and makes tickets |
| **10**| **Create Session History Tickets** – converts each room visited in `mudd_session.json` into a ticket |
| **11**| **Self‑Guided Documentation** – displays this guide inside gmudd |
| **12**| **Project File Tree Explorer** – shows a color-coded view of the project structure |
| **13**| **Browse Documentation (.md)** – interactive browser to read any markdown file in the project |
| **14**| **Mesh Divergence Audit** – lists all nodes and their unique/divergent files |
| **15**| **Reconcile Diff & Commit (Git)** – review changes and commit to mesh history |
| **16**| **Debug Build: Root Integrity Audit** – scan for all root-owned assets |
| **17**| **Virtualize Documentation** – materialize rooms from `.md` files |
| **0** | Exit the Swarm Matrix |

---
## Ticketing Features
### 1️⃣ Create Incomplete Task Tickets
- **What it does:** Runs `scripts/create_incomplete_task_tickets.py`.
- **Scope:** Scans all `.py`, `.go`, `.js`, `.ts`, `.tsx`, `.rs`, `.sh`, `.md`, `.json`, `.yaml`, `.proto`, and more.
- **Result:** Every unchecked checklist item (`- [ ]`) across the entire repository tree becomes a row in the `tickets` table (`Queue = "Incomplete-Tasks"`).
| **15**| **Reconcile Diff & Commit (Git)** – review changes and commit to mesh history |
| **16**| **Debug Build: Root Integrity Audit** – scan for all root-owned assets |
| **0** | Exit the Swarm Matrix |

---
## Navigation & Exploration
### 📁 Project File Tree Explorer (Option 12)
- **Debug Mode:** Now includes hidden files/folders (except `.git`).
- **Ownership Tracking:** Files owned by `root` are explicitly flagged with a `{C_RED}[ROOT]{C_RESET}` tag for security auditing.

### 📚 Documentation Browser (Option 13)
- **Full Scope:** Now includes hidden documentation (e.g., from `.agent` or `.github`) providing a complete view of swarm logic.

### 📉 Mesh Divergence Audit (Option 14)
- Scans the entire 9-node mesh (including Phased and Local nodes).
- Maps specific files and state deviations to their respective nodes based on forensic tickets and Jetweb logs.

### 💾 Reconcile Diff & Commit (Option 15)
- Provides a "forensic diff" of all uncommitted changes in the repository.
- Allows for direct `git commit` after reviewing the state of the mesh.

### 🔍 Debug Build: Root Integrity Audit (Option 16)
- A comprehensive deep-scan of the repository that identifies *every* asset owned by the `root` user.
- Useful for verifying that critical system-level binaries and configs are correctly permissioned.

### 🔮 Virtualize Documentation (Option 17)
- **Agent:** `DOC-PARSER-01`
- **What it does:** Scans the repository for virtually engineered room definitions embedded in markdown files.
- **Syntax:** Use HTML comments like `<!-- ROOM_START id="ID" ... -->` to define navigable space.
- **Result:** Documentation literally becomes part of the world. You can teleport to these rooms and interact with their descriptions.

### 🌐 Network Topology Audit (Option 18)
- **Agent:** `SCANNER-01`
- **What it does:** Performs a multi-node sweep (ifconfig/nmap) to visualize the actual network state across all 10 nodes.
- **Result:** Displays a real-time grid of all interfaces, IP bindings, and open ports, identifying potential security or connection issues in the swarm.

### 🆔 Agent Identity Audit (Option 19)
- **What it does:** Displays the registry of unique identifiers for all 33 GSD agents.
- **Shortcode Format:** `5alpha#XXX` (e.g., `5alpha#50c` for `gsd-advisor-researcher`).
- **Purpose:** Used for cross-node communication and global registry lookups, ensuring every agent is globally addressable.

### 2️⃣ Create Session History Tickets
- **What it does:** Runs `scripts/create_session_history_tickets.py`.
- **Result:** Each room ID stored in `mudd_session.json`’s `history` list generates a ticket (`Queue = "Session-History"`).
- **When to use:** To retroactively persist the path you traversed during a debugging or exploration session.

---
## Session Persistence & Resumption
- The file `mudd_session.json` now stores **additional flags** (e.g., `incomplete_task_tickets_created`, `session_history_tickets_created`).
- On start‑up `load_session()` reads the file and restores:
  - Current room location
  - Navigation history
  - Any boolean flags you set via `set_session_flag()`.
- **Resuming:** Simply run `gmudd` again; the UI will reflect the saved state and you can continue where you left off without re‑creating tickets.

---
## How to Use the New Commands
```bash
# Start the dashboard
$ gmudd
```
1. Choose **9** to generate tickets for every unchecked task.
2. Choose **10** to generate tickets for the rooms you visited.
3. Choose **11** to view this guide inside the console – handy when you’re already inside the UI.
4. After any ticket generation, the session file is automatically updated with a flag, so repeating the command will not duplicate tickets unless you clear the flag manually.

---
## Extending the Guide
- The guide lives at `docs/gmudd_guide.md`. Edit this file to add future menu items, screenshots, or workflow tips.
- The UI loads it on‑demand, so any changes are instantly visible the next time you press **11**.

---
## Screenshots (optional)
> *Add PNG/JPG screenshots in the `docs/` folder and reference them here using Markdown image syntax if you want visual aids.*

---
**Happy exploring!**
