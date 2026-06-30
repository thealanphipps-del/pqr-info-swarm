-- Migration to Sovereign Mesh v7.1
DROP TABLE IF EXISTS agent_memory_index CASCADE;
DROP TABLE IF EXISTS ticket_relationships CASCADE;
DROP TABLE IF EXISTS ticket_content CASCADE;
DROP TABLE IF EXISTS tickets CASCADE;
DROP TABLE IF EXISTS rtgo CASCADE;

-- Re-create with v7.1 Schema
CREATE TABLE tickets (
    ticket_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    layer_id INT NOT NULL CHECK (layer_id >= 1 AND layer_id <= 7),
    creator_agent_id STRING NOT NULL,
    status STRING DEFAULT 'PENDING',
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE ticket_relationships (
    parent_id UUID REFERENCES tickets(ticket_id),
    child_id UUID REFERENCES tickets(ticket_id),
    relationship_type STRING NOT NULL,
    PRIMARY KEY (parent_id, child_id)
);

CREATE TABLE ticket_content (
    ticket_id UUID REFERENCES tickets(ticket_id) PRIMARY KEY,
    intent_blob JSONB,
    state_vector FLOAT8[],
    consensus_score DECIMAL DEFAULT 0.0,
    raw_content BYTEA,
    summary_hash STRING,
    payload_hash STRING
);

CREATE TABLE agent_memory_index (
    agent_id STRING,
    ticket_id UUID REFERENCES tickets(ticket_id),
    context_depth INT DEFAULT 7,
    last_accessed TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (agent_id, ticket_id)
);

-- Initialize Genesis Block v7.1
INSERT INTO tickets (ticket_id, layer_id, creator_agent_id, status)
VALUES ('00000000-0000-0000-0000-000000000000', 1, 'GENESIS', 'IMMUTABLE');

INSERT INTO agent_memory_index (agent_id, ticket_id)
VALUES ('ORACLE', '00000000-0000-0000-0000-000000000000');
