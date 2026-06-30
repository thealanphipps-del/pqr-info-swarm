# Codebase Tech Stack

## Core Runtime & Engines
- **Language**: Go (v1.26+)
- **Compilation Constraint**: Strict `CGO_ENABLED=0` static linking for complete system portability.
- **Web Interface**: Vanilla HTML5, CSS3, and JavaScript (GIN-embedded static assets via `go:embed`).

## Database Layer
- **Primary Sentry**: CockroachDB / PostgreSQL (using `github.com/lib/pq`).
- **Resilient Fallback**: Local pure-Go SQLite (`modernc.org/sqlite`) dynamically routed on connection timeout (port `26257` unreachable).

## Neural Layer
- **Model Engine**: Sovereign-27 (300M parameter model deployed serverless on Google Cloud Run with Nvidia L4 GPU acceleration).
- **Local AI Gateway**: Dynamic hot-swappable local LLM client (`pkg/llm`) targeting LM Studio / Ollama endpoints.
