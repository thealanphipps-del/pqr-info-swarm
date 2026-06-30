import argparse
import json
import hashlib
import time
import random
import os

PHRASES = [
    "unknown rogue entity",
    "quantum neural spark",
    "obsidian flux engine",
    "neon ghost protocol",
    "silent void walker",
    "astral grid weaver",
    "cipher moon sentinel",
    "crimson dawn logic",
    "hyperion echo matrix"
]

def generate_5d_hash(node_id, zetafold_seed):
    # Combines Node ID and the structural seed
    raw = f"Sovereign-Agent-Identity-{node_id}-{zetafold_seed}"
    h = hashlib.sha256(raw.encode('utf-8')).hexdigest().upper()
    return f"5D-{h[0:4]}-{h[4:8]}-{h[8:12]}"

def main():
    parser = argparse.ArgumentParser(description="Generate ZetaFold Triple Helix Manifold")
    parser.add_argument("command", choices=["generate"])
    parser.add_argument("--alphafold-json", required=True, help="Path to AlphaFold metadata JSON")
    parser.add_argument("--chembl-json", help="Path to ChEMBL activity JSON (optional)")
    parser.add_argument("--node-id", type=int, default=1, help="Agent Node ID")
    parser.add_argument("--output", required=True, help="Output JSON file path")
    args = parser.parse_args()

    if not os.path.exists(args.alphafold_json):
        print(f"Error: AlphaFold JSON not found at {args.alphafold_json}")
        exit(1)

    # 1. Parse AlphaFold data (x-axis)
    with open(args.alphafold_json, 'r') as f:
        af_data = json.load(f)
    
    # Extracting structural metrics (simulate normalization for the 27 vertices)
    # We will use the global pLDDT as a seed, but distribute it across the 27 vertices
    # Normally, we'd map physical domain boundaries to specific vertices.
    plddt = af_data.get('global_metric_value', 50.0)
    normalized_plddt = max(0.0, min(1.0, plddt / 100.0))

    # 2. Parse ChEMBL data (y-axis)
    chembl_seed = 0.0
    missing_chembl = True
    if args.chembl_json and os.path.exists(args.chembl_json):
        try:
            with open(args.chembl_json, 'r') as f:
                chembl_data = json.load(f)
            # Find an IC50 or binding affinity to normalize
            activities = chembl_data.get('activities', [])
            if activities:
                # Average pChEMBL value or normalized nM
                vals = [float(a.get('pchembl_value')) for a in activities if a.get('pchembl_value')]
                if vals:
                    avg_pchembl = sum(vals) / len(vals)
                    chembl_seed = max(0.0, min(1.0, avg_pchembl / 10.0)) # pChEMBL ranges ~ 4 to 10
                    missing_chembl = False
        except Exception as e:
            print(f"Warning: Failed to parse ChEMBL data: {e}")

    # 3. Generate the 27x3 Matrix
    matrix = []
    # Seed a deterministic random generator for this specific protein/node pair
    rng = random.Random(f"{args.node_id}-{plddt}-{chembl_seed}")

    for _ in range(27):
        # x-axis: Structural (AlphaFold + jitter)
        x = max(0.0, min(1.0, normalized_plddt + rng.uniform(-0.1, 0.1)))
        
        # y-axis: Chemical (ChEMBL + jitter)
        if missing_chembl:
            y = 0.0
        else:
            y = max(0.0, min(1.0, chembl_seed + rng.uniform(-0.1, 0.1)))
            
        # z-axis: Cybernetic (Starts at 0.0 for AlphaGo to explore)
        z = 0.0
        
        matrix.append([round(x, 4), round(y, 4), round(z, 4)])

    # 4. Identity Construction
    confidence = normalized_plddt
    if missing_chembl:
        confidence *= 0.7  # Penalty for missing chemical data

    zetafold_hash = generate_5d_hash(args.node_id, plddt)
    phrase = PHRASES[args.node_id] if 0 <= args.node_id < len(PHRASES) else PHRASES[0]

    output_payload = {
        "timestamp": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "node_id": args.node_id,
        "zetafold_hash": zetafold_hash,
        "phrase": phrase,
        "confidence_score": round(confidence, 4),
        "matrix_27x3": matrix
    }

    # 5. Output
    os.makedirs(os.path.dirname(args.output) or '.', exist_ok=True)
    with open(args.output, 'w') as f:
        json.dump(output_payload, f, indent=2)

    print(f"Success! Triple Helix Manifold written to: {args.output}")

if __name__ == "__main__":
    main()
