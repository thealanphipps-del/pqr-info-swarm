-- Sovereign Mesh: Ticketing Fabric Schema v7.1
-- Designed for 7-Layer Genealogy and Anti-Hallucination Architecture

-- 1. Atomic Tickets
CREATE TABLE IF NOT EXISTS tickets (
    ticket_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    layer_id INT NOT NULL CHECK (layer_id >= 1 AND layer_id <= 7),
    creator_agent_id STRING NOT NULL, -- Hexagram ID (e.g., 'GATE-01')
    status STRING DEFAULT 'PENDING',
    created_at TIMESTAMPTZ DEFAULT now()
);

-- 2. Causal Relationship Graph
CREATE TABLE IF NOT EXISTS ticket_relationships (
    parent_id UUID REFERENCES tickets(ticket_id),
    child_id UUID REFERENCES tickets(ticket_id),
    relationship_type STRING NOT NULL, -- 'EVOLUTION', 'CONSEQUENCE', 'CONTEXT', 'GENESIS'
    PRIMARY KEY (parent_id, child_id)
);

-- 3. Content-Addressable Storage (Content-Heavy)
CREATE TABLE IF NOT EXISTS ticket_content (
    ticket_id UUID REFERENCES tickets(ticket_id) PRIMARY KEY,
    intent_blob JSONB,
    state_vector FLOAT8[], -- For semantic search/embedding
    consensus_score DECIMAL DEFAULT 0.0,
    raw_content BYTEA,
    summary_hash STRING, -- SHA-256
    payload_hash STRING -- SHA-256
);

-- 4. Agent Memory & Context Index
CREATE TABLE IF NOT EXISTS agent_memory_index (
    agent_id STRING,
    ticket_id UUID REFERENCES tickets(ticket_id),
    context_depth INT DEFAULT 7,
    last_accessed TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (agent_id, ticket_id)
);

-- 5. Sticky Rule ST-006: Genesis Block Initialization
-- Initialized in init.sql
