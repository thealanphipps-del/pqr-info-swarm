# RFC-004: Merkle Proof Routing

- **Status:** Proposed
- **Target Version:** v0.4.0
- **Author:** Antigravity Engineering

## Summary
Introduce partial chunk verification via Merkle proofs to eliminate full recomputation during mismatch detection.

## Specification
Nodes must:
1. Compute the MerkleRoot for each task block.
2. Return a list of per-chunk hashes.

On mismatch:
1. The requester requests a specific chunk proof.
2. Only that chunk is recomputed.

## Consensus Rule Update
- Root must match quorum.
- Chunk mismatch triggers targeted recompute.
