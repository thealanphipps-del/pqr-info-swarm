# E2E Test Infra: Quantasona Sovereign Mesh

## Test Philosophy
- Opaque-box, requirement-driven. No dependency on internal view implementation details.
- Runs via JUnit on the JVM for fast, deterministic execution during build.

## Feature Inventory
| # | Feature | Source (requirement) | Tier 1 | Tier 2 | Tier 3 | Tier 4 |
|---|---------|---------------------|:------:|:------:|:------:|:------:|
| 1 | HpaAtlasScreen | ORIGINAL_REQUEST R2 | 5 | 5 | ✓ | ✓ |
| 2 | GemMatchScreen | ORIGINAL_REQUEST R2 | 5 | 5 | ✓ | ✓ |
| 3 | GeologyScannerScreen | ORIGINAL_REQUEST R2 | 5 | 5 | ✓ | ✓ |
| 4 | HudTelemetryScreen & Helium | ORIGINAL_REQUEST R3 | 5 | 5 | ✓ | ✓ |

## Test Architecture
- Test Runner: Gradle command `./gradlew test --tests com.example.quantasonaapp.E2ETestSuite`
- Test Case Format: Kotlin JUnit test cases validating state transitions, boundary conditions, cross-feature interaction, and real-world workloads.

## Coverage Thresholds
- Tier 1: 20 tests (5 per feature)
- Tier 2: 20 tests (5 per feature)
- Tier 3: 4 tests (pairwise cross-feature)
- Tier 4: 5 tests (real-world workloads)
- Total: 49 tests
