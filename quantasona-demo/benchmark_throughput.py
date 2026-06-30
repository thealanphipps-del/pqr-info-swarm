import json
import time
import sqlite3
import random
import os
import argparse

# Configuration
INPUT_FILE = "anonymized_real_voice_prints.json"
DB_FILE = "benchmark.db"

def run_benchmark(multiplier):
    if os.path.exists(DB_FILE):
        os.remove(DB_FILE)

    print(f"Loading base dataset from {INPUT_FILE}...")
    with open(INPUT_FILE, "r") as f:
        base_data = json.load(f)

    print(f"Base dataset size: {len(base_data)}")
    
    # 1. Amplifying the dataset
    print(f"Amplifying dataset {multiplier}x...")
    amplified_data = []
    for i in range(multiplier):
        for item in base_data:
            new_item = item.copy()
            # Modify sigmaId slightly to ensure uniqueness for indexing
            new_item["sigmaId"] = f"{item['sigmaId']}_{i}"
            amplified_data.append(new_item)
            
    total_records = len(amplified_data)
    print(f"Amplified dataset size: {total_records}")
    
    # Connect to SQLite
    conn = sqlite3.connect(DB_FILE)
    cursor = conn.cursor()
    
    # Create table
    cursor.execute('''
    CREATE TABLE deltas (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sigmaId TEXT,
        semanticWeight REAL,
        confidence REAL,
        provenance TEXT,
        deltaType TEXT,
        relationType TEXT,
        sourceMic TEXT
    )
    ''')
    conn.commit()

    # 2. Ingestion
    print("\nStarting ingestion...")
    start_time = time.time()
    
    # Insert in batches
    batch_size = 1000
    for i in range(0, total_records, batch_size):
        batch = amplified_data[i:i+batch_size]
        cursor.executemany('''
            INSERT INTO deltas (sigmaId, semanticWeight, confidence, provenance, deltaType, relationType, sourceMic)
            VALUES (:sigmaId, :semanticWeight, :confidence, :provenance, :deltaType, :relationType, :sourceMic)
        ''', batch)
    conn.commit()
    
    ingest_time = time.time() - start_time
    print(f"Ingestion Time: {ingest_time:.4f} seconds")
    if ingest_time > 0:
        print(f"Ingestion Throughput: {total_records / ingest_time:.2f} records/second")
    
    # 3. Indexing
    print("\nStarting indexing...")
    start_time = time.time()
    cursor.execute('CREATE INDEX idx_sigmaId ON deltas(sigmaId)')
    conn.commit()
    index_time = time.time() - start_time
    print(f"Indexing Time: {index_time:.4f} seconds")
    
    # 4. Querying / Hashing lookups
    print("\nStarting querying (hashing lookups)...")
    # Pick random subset of sigmaIds to query
    query_count = min(10000, total_records)
    random_samples = random.sample(amplified_data, query_count)
    query_ids = [item['sigmaId'] for item in random_samples]
    
    start_time = time.time()
    found = 0
    for q_id in query_ids:
        cursor.execute('SELECT * FROM deltas WHERE sigmaId = ?', (q_id,))
        if cursor.fetchone():
            found += 1
            
    query_time = time.time() - start_time
    print(f"Query Time for {query_count} lookups: {query_time:.4f} seconds")
    if query_count > 0:
        print(f"Query Latency: {(query_time / query_count) * 1000:.4f} ms/query")
        if query_time > 0:
            print(f"Query Throughput: {query_count / query_time:.2f} queries/second")
    
    conn.close()
    
    # Print Bottleneck Analysis
    print("\n--- Bottleneck Analysis ---")
    if ingest_time > 5.0:
        print("- Bottleneck detected: Database write operations are slow under this load.")
        print("  Recommendation: Consider bulk inserts with larger batches, disabling synchronous disk writes, or switching to an in-memory/NoSQL datastore.")
    else:
        print("- Database write operations handled the load without significant bottlenecks.")
        
    if query_count > 0 and (query_time / query_count) * 1000 > 1.0:
        print("- Bottleneck detected: Hashing lookups / Query latency is high (>1ms per query).")
        print("  Recommendation: Implement a caching layer (e.g., Redis) or in-memory bloom filters to speed up index lookups under load.")
    else:
        print("- Hashing lookups (Indexed queries) were fast and did not bottleneck the system.")

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--multiplier", type=int, default=25000, help="How many times to multiply the base dataset")
    args = parser.parse_args()
    
    run_benchmark(args.multiplier)
