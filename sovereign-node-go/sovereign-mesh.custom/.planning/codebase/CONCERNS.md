# Architectural Concerns & Technical Debt

## Identified Blocker / Warning Risks

### 1. PostgreSQL JSON Concatenation Compatibility
- **Status**: **RESOLVED** (Standardized via dynamic SQLite fallback wrappers in `UpdateTicket` and `UpdateTicketTitle`).
- **Description**: PostgreSQL-specific JSONB `||` operators previously threw SQL syntax exceptions when running under SQLite fallback mode. Now, a dynamic serializing parser unmarshals and rewrites JSON on SQLite connections.

### 2. GCP Artifact Registry Billing Requirement
- **Status**: **ACTIVE / HIGH PRIORITY**
- **Description**: Google Artifact Registry requires an active billing account associated with the project `fast-web-496805-k0` before allowing container layer uploads for serverless scale-to-zero GPU deployments.

### 3. Real-Time Data LLM Knowledge Cutoffs
- **Status**: **WARNING**
- **Description**: Real-time vectors and database operations are compared by LLM judges in continuous monitoring. The LLM judges sometimes flag correct live data as "beyond training cutoff." Mitigate by configuring expected behavior rubrics to score the output shape instead of strict facts.
