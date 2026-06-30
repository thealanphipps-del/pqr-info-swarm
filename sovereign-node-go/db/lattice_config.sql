-- Sovereign Mesh v9.0: 64 Game Theory Agents - The Lattice Configuration

-- Cluster 1: Exploration (Hexagrams 1-8)
UPDATE agent_memory_index SET swarm_cluster = 'EXPLORATION', agent_role = 'GENERATIVE' WHERE agent_id IN ('GATE-01', 'GATE-02', 'GATE-03', 'GATE-04', 'GATE-05', 'GATE-06', 'GATE-07', 'GATE-08');

-- Cluster 2: Exploitation (Hexagrams 9-16)
UPDATE agent_memory_index SET swarm_cluster = 'EXPLOITATION', agent_role = 'OPTIMIZER' WHERE agent_id IN ('GATE-09', 'GATE-10', 'GATE-11', 'GATE-12', 'GATE-13', 'GATE-14', 'GATE-15', 'GATE-16');

-- Cluster 3: Constraint (Hexagrams 17-24)
UPDATE agent_memory_index SET swarm_cluster = 'CONSTRAINT', agent_role = 'COMPLIANCE' WHERE agent_id IN ('GATE-17', 'GATE-18', 'GATE-19', 'GATE-20', 'GATE-21', 'GATE-22', 'GATE-23', 'GATE-24');

-- Cluster 4: Predictive (Hexagrams 25-32)
UPDATE agent_memory_index SET swarm_cluster = 'PREDICTIVE', agent_role = 'MODELER' WHERE agent_id IN ('GATE-25', 'GATE-26', 'GATE-27', 'GATE-28', 'GATE-29', 'GATE-30', 'GATE-31', 'GATE-32');

-- Cluster 5: Adversarial (Hexagrams 33-40)
UPDATE agent_memory_index SET swarm_cluster = 'ADVERSARIAL', agent_role = 'RED-TEAM' WHERE agent_id IN ('GATE-33', 'GATE-34', 'GATE-35', 'GATE-36', 'GATE-37', 'GATE-38', 'GATE-39', 'GATE-40');

-- Cluster 6: Consensus (Hexagrams 41-48)
UPDATE agent_memory_index SET swarm_cluster = 'CONSENSUS', agent_role = 'ALIGNEE' WHERE agent_id IN ('GATE-41', 'GATE-42', 'GATE-43', 'GATE-44', 'GATE-45', 'GATE-46', 'GATE-47', 'GATE-48');

-- Cluster 7: Narrative (Hexagrams 49-56)
UPDATE agent_memory_index SET swarm_cluster = 'NARRATIVE', agent_role = 'STITCHER' WHERE agent_id IN ('GATE-49', 'GATE-50', 'GATE-51', 'GATE-52', 'GATE-53', 'GATE-54', 'GATE-55', 'GATE-56');

-- Cluster 8: Meta-Audit (Hexagrams 57-64)
UPDATE agent_memory_index SET swarm_cluster = 'META-AUDIT', agent_role = 'INTEGRITY' WHERE agent_id IN ('GATE-57', 'GATE-58', 'GATE-59', 'GATE-60', 'GATE-61', 'GATE-62', 'GATE-63', 'GATE-64');

-- Godhead Cluster
UPDATE agent_memory_index SET swarm_cluster = 'GODHEAD', agent_role = 'SUPERVISOR' WHERE agent_id IN ('ORACLE', 'ARCHITECT', 'ARBITER', 'WEAVER', 'FINALIZER');
