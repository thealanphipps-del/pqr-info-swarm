-- Sovereign Mesh: RTGO Memory v2 (v9.0 Upgrade)

-- 1. Metadata Layer
CREATE TABLE IF NOT EXISTS ticket_metadata (
    ticket_id UUID REFERENCES tickets(ticket_id) PRIMARY KEY,
    external_ref_id STRING, -- e.g., 'GH-123'
    source_system STRING, -- e.g., 'GITHUB', 'S25', 'RTGO_SSH'
    audit_trail JSONB,
    forensic_signature STRING
);

-- 2. Non-Causal References (Lattice Links)
CREATE TABLE IF NOT EXISTS ticket_references (
    from_ticket_id UUID REFERENCES tickets(ticket_id),
    to_ticket_id UUID REFERENCES tickets(ticket_id),
    ref_type STRING NOT NULL, -- 'ANALOGY', 'CONTRADICTION', 'DEPENDENCY'
    PRIMARY KEY (from_ticket_id, to_ticket_id)
);

-- 3. Swarm Role Configuration
ALTER TABLE agent_memory_index ADD COLUMN IF NOT EXISTS swarm_cluster STRING;
ALTER TABLE agent_memory_index ADD COLUMN IF NOT EXISTS agent_role STRING;

-- 4. Create composite index for graph traversal performance
CREATE INDEX IF NOT EXISTS idx_relationships_parent_child ON ticket_relationships (parent_id, child_id);
CREATE INDEX IF NOT EXISTS idx_relationships_child_parent ON ticket_relationships (child_id, parent_id);
