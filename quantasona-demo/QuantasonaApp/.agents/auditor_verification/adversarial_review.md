## Challenge Summary

**Overall risk assessment**: LOW

## Challenges

### [Medium] Challenge 1: Concurrency Safety of InMemoryFiveDStore

- **Assumption challenged**: The store maps `vertices` and `edges` are assumed to be accessed safely by the application, but they are implemented as plain `MutableMap`s.
- **Attack scenario**: Background telemetry bridge tasks collect Helium beacons concurrently and record them in the store while pathfinding queries run from other threads or the UI thread.
- **Blast radius**: A `ConcurrentModificationException` could be thrown or updates lost, leading to graph inconsistency or app crashes.
- **Mitigation**: Replace `mutableMapOf` with `ConcurrentHashMap` or synchronize map accesses using a Kotlin coroutines `Mutex`.

### [Low] Challenge 2: Dynamic List Mutation in HudTelemetryScreen Canvas

- **Assumption challenged**: The list of `neighbors` passed to `HudTelemetryScreen` is assumed to be stable during Canvas rendering.
- **Attack scenario**: The Helium telemetry pipeline updates the neighbor strength or inserts new nodes while the Compose Canvas is drawing, causing thread-safety issues if the list reference is mutated.
- **Blast radius**: Although `StateFlow` emits new lists, rapid updates could trigger unnecessary recompositions or UI flickering.
- **Mitigation**: Create an immutable snapshot copy of the list when rendering inside the `Canvas` block, and use a key-based `remember` wrapper.

### [Low] Challenge 3: Extreme Signal Strength Normalization Loss

- **Assumption challenged**: The RSSI values are assumed to fall within the standard operating range of `-100.0` to `-40.0`.
- **Attack scenario**: Beacons with extremely high power (e.g., `-30.0` dBm) or extremely low power (e.g., `-120.0` dBm) are registered.
- **Blast radius**: While the values are safely clamped to `1.0f` and `0.0f` using `coerceIn`, fine-grained topology resolution is lost at the boundaries.
- **Mitigation**: Introduce a dynamic calibration mapping that adjusts the minimum and maximum expected RSSI ranges based on historical signals.

## Stress Test Results

- **Empty Neighbors List** → Draw canvas with 0 nodes → Handled safely as `forEachIndexed` is not entered, avoiding division by zero → **PASS**
- **Concurrent addTrits Calls** → Call `addTrits` from 10 concurrent coroutines → `MutableStateFlow.update` runs atomically, state converges safely → **PASS**
- **Extreme RSSI values** → Pass RSSI of `-30` and `-110` → Clamped correctly to `1.0` and `0.0` by adapter → **PASS**

## Unchallenged Areas

- **VoiceRecorderManager** — Not challenged deeply as it delegates to platform-specific microphone recording, which was mock-backed during headless test suites.
- **FilecoinStorageEngine** — Not challenged as it uses an in-memory replica during test execution.
