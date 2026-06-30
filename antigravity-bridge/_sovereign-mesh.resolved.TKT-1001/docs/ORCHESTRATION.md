---
title: "ORCHESTRATION"
description: ""
date: "2026-05-21"
tags: []
---
---

## Table of Contents
- [🌌 Sovereign Mesh Orchestration Suite](#2)
- [🏗️ Orchestration Layers](#2)
- [1. System Layer: `mesh_control.sh`](#2)
- [2. Interactive Layer: `gmudd`](#2)
- [3. Native API Layer: `sovereign-cli`](#2)
- [4. Autonomous Layer: `sovereign-auto`](#2)
- [🤖 AI Workflows](#2)
- [🛠️ Feature Matrix](#2)
- [🧬 Integrated Workflows](#2)
- [Scenario A: Initializing the Swarm](#2)
- [Scenario B: Autonomous Feature Implementation](#2)
- [Scenario C: Forensic Recovery](#2)



# 🌌 Sovereign Mesh Orchestration Suite

This document defines the hierarchy and functional roles of the Sovereign Mesh orchestration tools.

---
---

## 🏗️ Orchestration Layers

### 1. System Layer: `mesh_control.sh`
The foundational script for infrastructure lifecycle and low-level mesh operations.
*   **Role:** Daemon management, network auditing, and forensic state mutation.
*   **Key Command:** `./mesh_control.sh start`
*   **Manual:** `man -l mesh_control.8`

### 2. Interactive Layer: `gmudd`
The human-centric "Global MUDD" dashboard.
*   **Role:** Real-time exploration, interactive navigation, and visual mesh auditing.
*   **Key Command:** `gmudd`
*   **Manual:** `man -l gmudd.8`

### 3. Native API Layer: `sovereign-cli`
The compiled Go interface for direct gRPC Control Bus interaction.
*   **Role:** Programmatic access to strike protocols, neural training, and citizenship.
*   **Key Command:** `./sovereign-cli ping`
*   **Manual:** `man -l sovereign-cli.1`

### 4. Autonomous Layer: `sovereign-auto`
The high-level autonomous agent and code-walking orchestrator.
*   **Role:** Code exploration, automatic goal execution, and self-healing.
*   **Key Command:** `./sovereign-auto walk .`
*   **Manual:** `man -l sovereign-auto.1`

---
---

## 🤖 AI Workflows
The mesh supports dual operational modes depending on whether the cloud Mothership is utilized.
*   **Workflow A (Cloud-Augmented):** Uses `gemini-cli` and cloud models for deep architectural refactoring.
*   **Workflow B (Sovereign-Local):** Uses `gemma-ui`, `./sovereign-auto`, and local GPU-bound models for private, persistent execution.
*   **Detailed Guide:** `docs/AI_WORKFLOW.md`
*   **Manual:** `man -l ai-workflow.7`

---
---

## 🛠️ Feature Matrix

| Feature | `mesh_control` | `gmudd` | `sovereign-cli` | `sovereign-auto` |
| :--- | :---: | :---: | :---: | :---: |
| **Start/Stop Daemons** | ✅ | ❌ | ❌ | ❌ |
| **Code Walking** | ❌ | ❌ | ❌ | ✅ |
| **Goal Execution** | ❌ | ❌ | ❌ | ✅ |
| **Self-Healing** | ❌ | ❌ | ❌ | ✅ |
| **State Forensics** | ✅ | ✅ | ❌ | ❌ |
| **Neural Training** | ❌ | ❌ | ✅ | ❌ |
| **Interactive Map** | ❌ | ✅ | ❌ | ❌ |

---
---

## 🧬 Integrated Workflows

### Scenario A: Initializing the Swarm
1.  Run `./mesh_control.sh start` to ignite the infrastructure.
2.  Use `gmudd` to verify all 10 nodes are online and phased.

### Scenario B: Autonomous Feature Implementation
1.  Run `./sovereign-auto execute "Add telemetry logging to dex.go"`
2.  The orchestrator will plan, implement, and verify the change autonomously.

### Scenario C: Forensic Recovery
1.  Identify a bad decision index in `gmudd` (Option 5).
2.  Use `./mesh_control.sh timemachine` to refactor the timeline history.
