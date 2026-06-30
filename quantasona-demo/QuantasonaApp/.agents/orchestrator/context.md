# Context

## Environment
- OS: Windows
- Project Root: `C:\Users\theal\QuantasonaApp`
- Gradle project with Jetpack Compose, Kotlin, navigation3 runtime

## Architecture & Codebase Discovery
- `MainActivity` launches `MainNavigation`.
- `Navigation.kt` manages navigation3 backstack starting with `Main`, displaying `MainScreen`.
- `MainScreen` lists four tabs: "HPA Atlas", "Gem Match", "Geology", "Node HUD".
- Data model contains `HpaGene`, `Addr5D`, `SilkRoadPacket`, `SignalVertex`, `RewardEvent`.
- `MockMeshClient` implements `MeshNodeClient`.
- `InMemoryFiveDStore` implements `FiveDVertexStore` and `FiveDEdgeStore` using CRDT merging.
- `DefaultGraphQueryEngine` query neighbors and find paths.
- `InsulinLattice` bootstraps healthy, T1D, and T2D states into graph stores.
- `HeliumClient` generates beacons and rewards.
- `HeliumMeshBridge` pipes beacons into Sovereign Mesh node clients.
- `MeshApp` defines the scope, contexts, and runtimes for running mesh actions.
