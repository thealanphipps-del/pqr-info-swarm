---
title: "AI_WORKFLOW"
description: ""
date: "2026-05-21"
tags: []
---
---

## Table of Contents
- [🤖 Sovereign Mesh AI Workflows](#2)
- [🛰️ Workflow A: With Gemini-CLI (Cloud-Augmented)](#2)
- [🛠️ Key Components](#2)
- [🔄 Operational Loop](#2)
- [✅ Best For:](#2)
- [🧬 Workflow B: Without Gemini-CLI (Sovereign-Local)](#2)
- [🛠️ Key Components](#2)
- [🔄 Operational Loop](#2)
- [✅ Best For:](#2)
- [📊 Comparison Matrix](#2)
- [🚀 Choosing Your Path](#2)



# 🤖 Sovereign Mesh AI Workflows

The Sovereign Mesh supports two distinct AI operational modes: **Cloud-Augmented (with Gemini-CLI)** and **Sovereign-Local (without Gemini-CLI)**.

---
---

## 🛰️ Workflow A: With Gemini-CLI (Cloud-Augmented)

This workflow utilizes the "Mothership" for high-level reasoning and complex sub-agent delegation.

### 🛠️ Key Components
*   **Primary Brain:** Gemini Pro/Flash (Cloud-based).
*   **Interface:** Standard `gemini-cli` environment.
*   **Context:** Large context window with global project instructions.
*   **Orchestration:** High-level tool use (grep, replace, shell).

### 🔄 Operational Loop
1.  **Directive:** User provides a goal via `gemini-cli`.
2.  **Delegation:** `gemini-cli` invokes specialized sub-agents (`codebase_investigator`, `generalist`).
3.  **Synthesis:** Cloud models perform deep analysis and plan-driven execution.
4.  **Verification:** Uses standard GSD skills to audit work.

### ✅ Best For:
*   Complex architectural refactoring.
*   Initial project bootstrapping.
*   High-level strategy and planning.

---
---

## 🧬 Workflow B: Without Gemini-CLI (Sovereign-Local)

This workflow is entirely private, GPU-bound, and runs natively on the local mesh infrastructure.

### 🛠️ Key Components
*   **Primary Brain:** Local Gemma Models (gemma-7b, codegemma).
*   **Interface:** `gemma-ui` (IDE + Chat) and `sovereign-auto` (CLI).
*   **Memory:** Permanent recall via CockroachDB `agentic_memories`.
*   **Tools:** Native MCP tool suite on Port 1111 (via `mgsh_mcp.py`).

### 🔄 Operational Loop
1.  **Directive:** User interacts via **Sovereign Gemma IDE** (localhost:3000) or `./sovereign-auto`.
2.  **Recall:** Gemma automatically retrieves historical context from the mesh ledger.
3.  **Execution:** `./sovereign-auto execute` triggers the local GSD-Executor.
4.  **Self-Healing:** `./sovereign-auto fix` ignites local repair cycles.
5.  **Backup:** If a query is too complex, Gemma uses `/backup` to consult the Mothership but remains the final arbiter.

### ✅ Best For:
*   Air-gapped development and high-security codebases.
*   Low-latency, GPU-accelerated local execution.
*   Continuous learning with perfect recall.

---
---

## 📊 Comparison Matrix

| Feature | Workflow A (Mothership) | Workflow B (Sovereign) |
| :--- | :--- | :--- |
| **Model** | Gemini Pro / Flash | Local Gemma / CodeGemma |
| **Latency** | Network Dependent | Sub-second (Local GPU) |
| **Privacy** | Cloud Processing | 100% Local / On-Prem |
| **Memory** | Session-bound | Permanent (CockroachDB) |
| **Interface** | `gemini-cli` | `gemma-ui` / `sovereign-auto` |
| **Failover** | High-level orchestration | Native SSH Fallback |

---
---

## 🚀 Choosing Your Path
*   **Use Workflow A** when you need maximum reasoning power and are connected to the Mothership.
*   **Use Workflow B** for daily autonomous maintenance, private state persistence, and native mesh integration.
