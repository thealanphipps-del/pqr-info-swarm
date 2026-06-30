import json
import random
import statistics

def load_deltas(filepath):
    with open(filepath, 'r') as f:
        return json.load(f)

def generate_positive_controls(deltas):
    # Historical records with known biographical overlaps. True matches.
    records = []
    for d in deltas:
        records.append({
            "sigmaId": d["sigmaId"],
            "semanticWeight": d["semanticWeight"] + random.uniform(-0.5, 0.5), # Slight noise
            "source": "historical-audio"
        })
    return records

def generate_cross_modality(deltas):
    # Records from different source modalities. True matches, but more noise.
    records = []
    for d in deltas:
        records.append({
            "sigmaId": d["sigmaId"],
            "semanticWeight": d["semanticWeight"] + random.uniform(-2.0, 2.0), # Higher noise
            "source": "written-interview"
        })
    return records

def generate_negative_controls(num_records):
    # Randomly selected, non-matching records
    records = []
    for _ in range(num_records):
        records.append({
            "sigmaId": f"biomarker/anon_{random.randint(1000000, 9999999)}",
            "semanticWeight": random.uniform(-50.0, 50.0),
            "source": "random-database"
        })
    return records

def match_algorithm(delta, database, threshold=1.0):
    matches = []
    for record in database:
        if record["sigmaId"] == delta["sigmaId"]:
            diff = abs(record["semanticWeight"] - delta["semanticWeight"])
            if diff <= threshold:
                matches.append((record, diff))
    return matches

def run_test(deltas, db, expected_positive=True, threshold=1.0):
    true_positives = 0
    false_positives = 0
    false_negatives = 0
    true_negatives = 0
    
    correlation_strengths = [] # Store similarities for true positives (1 / (1 + diff))
    
    for delta in deltas:
        matches = match_algorithm(delta, db, threshold)
        if expected_positive:
            if len(matches) > 0:
                true_positives += 1
                for m, diff in matches:
                    correlation_strengths.append(1.0 / (1.0 + diff))
            else:
                false_negatives += 1
        else:
            if len(matches) > 0:
                false_positives += len(matches)
            else:
                true_negatives += len(db)
                
    fpr = 0.0
    if not expected_positive and (false_positives + true_negatives) > 0:
         fpr = false_positives / (false_positives + true_negatives)
         
    accuracy = 0.0
    if expected_positive and len(deltas) > 0:
         accuracy = true_positives / len(deltas)
         
    avg_corr = statistics.mean(correlation_strengths) if correlation_strengths else 0.0
    
    return {
        "true_positives": true_positives,
        "false_positives": false_positives,
        "false_negatives": false_negatives,
        "true_negatives": true_negatives,
        "fpr": fpr,
        "accuracy": accuracy,
        "correlation_strength": avg_corr
    }

def main():
    deltas = load_deltas('anonymized_real_voice_prints.json')
    print(f"Loaded {len(deltas)} deltas.")
    
    # 1. Positive Controls
    pos_db = generate_positive_controls(deltas)
    pos_results = run_test(deltas, pos_db, expected_positive=True, threshold=1.0)
    print("--- 1. Positive Controls ---")
    print(f"Accuracy (TPR): {pos_results['accuracy']:.2f}")
    print(f"Correlation Strength: {pos_results['correlation_strength']:.2f}")
    
    # 2. Cross-Modality
    cross_db = generate_cross_modality(deltas)
    cross_results = run_test(deltas, cross_db, expected_positive=True, threshold=2.5)
    print("--- 2. Cross-Modality Test ---")
    print(f"Accuracy (TPR): {cross_results['accuracy']:.2f}")
    print(f"Correlation Strength: {cross_results['correlation_strength']:.2f}")
    
    # 3. Negative Controls
    neg_db = generate_negative_controls(1000)
    neg_results = run_test(deltas, neg_db, expected_positive=False, threshold=1.0)
    print("--- 3. Negative Controls ---")
    print(f"False Positive Rate (FPR): {neg_results['fpr']:.4f}")
    print(f"False Positives: {neg_results['false_positives']}")
    
if __name__ == '__main__':
    main()
