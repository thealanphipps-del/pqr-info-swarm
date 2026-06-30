---
title: "sovereign-auto"
description: ""
date: "2026-05-21"
tags: []
---
---

## Table of Contents
- [🚀 SOVEREIGN-AUTO: Autonomous Orchestrator Guide](#2)
- [Overview](#2)
- [🛠️ Core Capabilities](#2)
- [1. Code Walking & Explanation (`walk`)](#2)
- [2. Autonomous Execution (`execute`)](#2)
- [3. Self-Healing & Repair (`fix`)](#2)
- [4. Integrity Checks (`check`)](#2)
- [🧬 Architectural Integration](#2)
- [📄 Manual Reference](#2)



# 🚀 SOVEREIGN-AUTO: Autonomous Orchestrator Guide

## Overview
`sovereign-auto` is the next-generation management interface for the Sovereign Mesh. Unlike traditional CLIs, it leverages the **Get Shit Done (GSD)** autonomous framework to perform complex code investigations and automated system repairs.

---
---

## 🛠️ Core Capabilities

### 1. Code Walking & Explanation (`walk`)
Recursively explores the codebase and uses the Mesh Engine's **7-Layer Pedigree** to explain the heritage and purpose of every file.
```bash
./sovereign-auto walk [path]
```

### 2. Autonomous Execution (`execute`)
Accepts high-level goals and spawns the `gsd-planner` and `gsd-executor` to fulfill them without human intervention.
```bash
./sovereign-auto execute "goal description"
```

### 3. Self-Healing & Repair (`fix`)
Identifies errors (or accepts them as input) and ignites the `gsd-debugger` to autonomously diagnose and apply fixes.
```bash
./sovereign-auto fix "Syntax error in ledger.go"
```

### 4. Integrity Checks (`check`)
Scans directories for syntax errors and architectural deviations from the Sovereign mandates.
```bash
./sovereign-auto check [path]
```

---
---

## 🧬 Architectural Integration
`sovereign-auto` is deeply integrated with the mesh:
*   **gRPC Control Plane:** Connects to `localhost:1111` for pedigree and state data.
*   **GSD Experts:** Directly orchestrates 33 specialized agents stored in `/agents`.
*   **Forensic Traceability:** All changes made via `execute` or `fix` are logged as PQR tickets in the ledger.

---
---

## 📄 Manual Reference
For detailed argument specifications, refer to the man page:
```bash
man -l sovereign-auto.1
```
