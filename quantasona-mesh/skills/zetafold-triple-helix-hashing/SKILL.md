---
name: zetafold-triple-helix-hashing
description: >-
  Generates a 27x3 Triple Helix manifold matrix, a cryptographic ZetaFold hash, and a 3-word identifier for an agent. It fuses physical protein topology (AlphaFold) with chemical interaction data (ChEMBL) to seed cybernetic AlphaGo agents.
---

# ZetaFold Triple Helix Hashing

## Overview
Generates the core cryptographic identity and initial 5D topological mapping weights for Sovereign Mesh agents. The output is a `[27][3]` grid matrix mapping AlphaFold (structural), ChEMBL (chemical), and AlphaGo (cybernetic) weights.

## Dependencies
This skill acts as an orchestrator. You MUST use the following existing skills to gather the raw data first:
- `alphafold-database-fetch-and-analyze`
- `chembl-database`

## Workflow

### 1. Fetch Structural Data (x-axis)
Use the `alphafold-database-fetch-and-analyze` skill to download the structure and metadata for the provided UniProt ID. Save the resulting AlphaFold JSON metadata file.
*If the API returns a 403 Forbidden, remember to set the `SCIENCE_SKILLS_USER_AGENT` environment variable!*

### 2. Fetch Chemical Data (y-axis)
Use the `chembl-database` skill to search for bioactivity data associated with the provided ChEMBL Target ID (e.g., `uv run scripts/chembl_api.py activity --filter target_chembl_id=CHEMBL203 --limit 50 --output /tmp/chembl_data.json`).

### 3. Generate the Manifold
Pass the output JSON files from steps 1 and 2 into this skill's generative script to construct the 27x3 Triple Helix mapping.

## Utility Scripts

The core computation is handled by the `generate_manifold.py` script.

**Usage:**
```bash
uv run scripts/generate_manifold.py generate \
  --alphafold-json /path/to/alphafold_metadata.json \
  --chembl-json /path/to/chembl_data.json \
  --node-id 4 \
  --output /tmp/triple_helix.json
```

**Missing Data Handling:**
If you could not find ChEMBL data, you may omit the `--chembl-json` argument. The script will automatically default the y-axis to `0.0` and lower the final `confidence_score` by multiplying it by 0.7. However, `--alphafold-json` is strictly required.

## Common Mistakes
- **Skipping Dependencies:** Do not attempt to parse the AlphaFold or ChEMBL APIs directly. Always use the upstream skills to produce the local JSON files first.
- **Missing `--output`:** The script writes to a file, not stdout. Always specify an output file.
