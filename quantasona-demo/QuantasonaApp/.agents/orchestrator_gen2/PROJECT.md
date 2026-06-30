# Project: Quantasona Sovereign Mesh Integration

## Architecture
- **Sovereign Mesh**: Runs client runtime (`MeshClient`/`MockMeshClient`) and CRDT graph engine (`MeshGraph`/`DefaultGraphQueryEngine`).
- **Data Flow**:
  - `HeliumClient` generates beacons (IOT signals) and rewards.
  - `HeliumMeshBridge` maps Helium beacons/rewards to Sovereign Mesh signals and reports them.
  - `DefaultDataRepository` updates trit balances and records signals.
  - `MainScreen` binds UI states and handles user navigation/actions.
- **Shared Interfaces**:
  - `FiveDVertexStore`/`FiveDEdgeStore` for graph persistence.
  - `GraphQueryEngine` for pathfinding and neighbor verification.

## Code Layout
- `app/src/main/java/com/example/quantasonaapp/MainActivity.kt`
- `app/src/main/java/com/example/quantasonaapp/Navigation.kt`
- `app/src/main/java/com/example/quantasonaapp/data/DataRepository.kt`
- `app/src/main/java/com/example/quantasonaapp/data/HeliumClient.kt`
- `app/src/main/java/com/example/quantasonaapp/data/MeshApp.kt`
- `app/src/main/java/com/example/quantasonaapp/data/MeshClient.kt`
- `app/src/main/java/com/example/quantasonaapp/data/MeshGraph.kt`
- `app/src/main/java/com/example/quantasonaapp/data/MeshModel.kt`
- `app/src/main/java/com/example/quantasonaapp/ui/main/MainScreen.kt`
- `app/src/test/java/com/example/quantasonaapp/ui/main/MainScreenViewModelTest.kt`

## Milestones
| # | Name | Scope | Dependencies | Status |
|---|------|-------|-------------|--------|
| 1 | E2E Test Suite Creation | Create robust opaque-box E2E test suite covering Tiers 1-4 | none | DONE |
| 2 | Compile & Unit Verification (R1) | Resolve compilation errors, ensure `./gradlew compileDebugKotlin` passes, and `MainScreenViewModelTest` passes | M1 | IN_PROGRESS |
| 3 | Compose UI Integration (R2) | Implement navigation tab container in `MainScreen.kt` to load `HpaAtlasScreen`, `GemMatchScreen`, `GeologyScannerScreen`, and `HudTelemetryScreen` | M2 | IN_PROGRESS |
| 4 | Helium & Telemetry Pipeline (R3) | Inject live beacons and rewards from `HeliumClient` into `MeshGraph` to update edge connection strengths | M3 | IN_PROGRESS |
| 5 | Adversarial Hardening (Tier 5) | Analyze coverage, run Challenger, Reviewer, and Forensic Audit checks to ensure CLEAN result | M4 | PLANNED |

## Interface Contracts
### HeliumClient ↔ MeshNodeClient
- `HeliumSignalAdapter.toSignalVertex` maps `HeliumBeacon` to `SignalVertex`.
- `HeliumRewardMapper.toRewardEvent` maps `HeliumReward` to `RewardEvent`.
- `HeliumMeshBridge` coordinates flow.
