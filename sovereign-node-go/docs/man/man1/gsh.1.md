# GSH(1) | Sovereign Node Manual

## NAME
gsh - Get Shit Done: The Unified Agentic Shell

## SYNOPSIS
**gsh** [OPTIONS] [COMMAND]

## DESCRIPTION
**gsh** is the primary entry point for agentic command orchestration within the Sovereign Node mesh. It acts as a forensic-aware wrapper around standard system shells, providing secure code injection and autonomous execution capabilities.

## OPTIONS
- **--secure-irs**
    Enable the Secure Injection Layer. This triggers the Forensic Guard and allows for atomic "Epoch" transitions via `secure_injection.py`.
- **--bridge-sync**
    Synchronize the local state with the Singularity HUD (Port 8081).
- **--version**
    Display version information. Current version: 4.2.

## ARCHITECTURE
GSH communicates via the **GSH Bridge** on Port 8081 (HTTP). It is designed to work across distributed nodes, including the Antigravity Laptop (Windows) and the Sovereign Master (S25 FE Phone).

## ENVIRONMENT
- **SOVEREIGN_NODE_ROOT**: Root directory of the mesh logic.
- **GSH_LOG_PATH**: Path to forensic command logs (default: `~/Jovian_Archives`).

## SEE ALSO
sovereign(1), secure-injection(8)
