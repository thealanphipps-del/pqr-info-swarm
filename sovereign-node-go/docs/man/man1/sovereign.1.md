# SOVEREIGN(1) | Sovereign Node Manual

## NAME
sovereign - Core Execution Engine and Mesh Consensus Node

## SYNOPSIS
**sovereign** [OPTIONS]

## DESCRIPTION
**sovereign** is the central execution engine of the Sovereign Node. It manages real-time trading logic, gRPC mesh communication, and consensus across the Starburst Monolith.

## FEATURES
- **Trading Engine**: Monitors RSI signals (Target: 28.5) and executes buy/sell signals based on mesh consensus.
- **Mesh Consensus**: Re-anchors to peers using gRPC on Port 1111.
- **Asset Lock**: Maintains a hardcoded baseline (Current: $814.68).

## NETWORK
- **gRPC Port**: 1111 (Engine communication)
- **HTTP Port**: 8080 (MCP Pipeline)

## STATE
Current status is tracked in the `rt_ledger.log` and synchronized to the cloud database (CockroachDB).

## SEE ALSO
gsh(1), singularity-hud(1)
