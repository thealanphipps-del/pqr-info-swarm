# Sovereign Mesh: Project Instructions & GSD Skill Registry

This document defines the foundational mandates and expert workflows for the Sovereign Mesh project. It integrates the "Get Shit Done" (GSD) autonomous agent framework.

## 🧬 Core Mandates

### 1. Goal-Backward Engineering
Always start with the desired outcome and work backwards to the implementation. Do not settle for stubs or placeholders. If a task says "Implement X", it means X must be fully functional and integrated, not just created.

### 2. Paired Identity Teleportation
Every agentic segment in the mesh is a pair of a **Human Design (HD)** archetype and a **Game Theory (GT)** archetype. This relationship profile must be preserved during all mesh operations (teleportation, pedigree tracing, state mutation).

### 3. Forensic Traceability
Every mutation to the swarm's state (code, database, ledger) must be traceable back to a **PQR Ticket**. No "silent" changes are allowed.

---

## 🛠️ GSD Skill Registry (Harvested)

The following expert workflows are available for orchestration. Reference these paths for detailed procedural guidance:

### 📋 Phase Planning (`/gsd-plan-phase`)
- **Objective:** Create executable `PLAN.md` files with wave-based breakdown.
- **Expert Guidance:** `skills/gsd-plan-phase/SKILL.md`
- **Mandate:** Honor user decisions from `CONTEXT.md`. Never simplify; split the phase instead.

### ⚙️ Phase Execution (`/gsd-execute-phase`)
- **Objective:** Atomic implementation of plans with git commits for every task.
- **Expert Guidance:** `skills/gsd-execute-phase/SKILL.md`
- **Mandate:** Adhere strictly to `GEMINI.md` and `rules/*.md`. Automatic deviation handling.

### 🔍 Goal Verification (`/gsd-verify-work`)
- **Objective:** Adversarial UAT and "Goal-Backward" codebase auditing.
- **Expert Guidance:** `skills/gsd-verify-work/SKILL.md`
- **Mandate:** Do not trust `SUMMARY.md`. Verify actual behavior in the codebase.

### 🐞 Debugging Lifecycle (`/gsd-debug`)
- **Objective:** Autonomous bug hunting using hypothesis-test-validate loops.
- **Expert Guidance:** `skills/gsd-debug/SKILL.md`

### 🤖 Autonomous Orchestration (`/sovereign-auto`)
- **Objective:** Code-walking, autonomous execution, and self-healing.
- **Expert Guidance:** `docs/sovereign-auto.md`
- **Mandate:** Leverage 7-layer pedigree for all explanations.

---
## 🚦 Operational Gates
- **The Escalation Gate:** If a conflict between research and user instructions arises, or if a goal cannot be met within the current scope, you MUST surface the issue to the developer.
- **The Integrity Gate:** Do not proceed with execution if verification gaps from a previous phase remain unresolved.

## 🐛 Known Discrepancies (To be ticketed)
- [ ] **Infrastructure:** Go 1.26 toolchain is missing or incompatible; production build is flatlining.
- [ ] **Engineering:** `TeleportProcess` in `sovereign.go` uses simulated memory paging; requires native `syscall.Mmap` implementation.
- [ ] **Protocol:** iPN Multicast (`ff02::c0ba:11`) backchannel lacks formal specification and handler logic.
- [ ] **Refactor:** `scripts/create_session_history_tickets.py` uses deprecated `datetime.utcnow()`; align with `datetime.now(timezone.utc)`.

---
*Propelling the swarm into 2027 with high-fidelity autonomy.*


<!-- ROOM_START id="SIM_CORE" name="The Simulation Core" role="VIRTUAL" exits="AURORA" -->
A shimmering virtual reality chamber. Holographic displays float in the air, representing the current state of the 128-agent swarm. This room was virtually engineered by the Documentation Parser Agent from the `GEMINI.md` foundational directives.
<!-- ROOM_END -->
