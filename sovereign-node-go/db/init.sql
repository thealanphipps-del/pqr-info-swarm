-- Initialize Genesis Block
INSERT INTO tickets (ticket_id, layer_id, creator_agent_id, status) 
VALUES ('00000000-0000-0000-0000-000000000000', 1, 'SYSTEM-GENESIS', 'IMMUTABLE')
ON CONFLICT (ticket_id) DO NOTHING;

INSERT INTO ticket_content (ticket_id, intent_blob, raw_content)
VALUES ('00000000-0000-0000-0000-000000000000', '{"name": "GENESIS_BLOCK", "rule": "ST-006"}', 'Genesis Block')
ON CONFLICT (ticket_id) DO NOTHING;

-- Initialize Oracle Agent Memory
INSERT INTO agent_memory_index (agent_id, ticket_id, context_depth)
VALUES ('ORACLE', '00000000-0000-0000-0000-000000000000', 7)
ON CONFLICT (agent_id, ticket_id) DO NOTHING;

