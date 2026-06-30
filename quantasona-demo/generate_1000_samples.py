import json
import numpy as np
import hashlib

def generate_baseline_samples(w_base, num_samples=1000, sigma_pitch=0.02, sigma_channel=0.001):
    D = len(w_base)
    baseline_samples = np.zeros((num_samples, D))
    
    # In MFCCs, the first few coefficients (0-3 out of 13) carry the bulk of pitch/spectral envelope
    # For a flattened 13x32 vector, we can treat indices corresponding to the first 4 MFCCs across all 32 frames as D_pitch
    d_pitch_indices = [i for i in range(D) if (i % 13) < 4]
    
    for i in range(num_samples):
        p_i = np.zeros(D)
        
        # 1. Simulate Pitch Variation
        p_i[d_pitch_indices] = np.random.normal(0, sigma_pitch, len(d_pitch_indices)) * w_base[d_pitch_indices]
        
        # 2. Simulate Channel Noise
        p_i += np.random.normal(0, sigma_channel, D) * w_base
        
        # 3. Simulate Speaking Rate Fluctuation
        scaling_factor = np.random.uniform(0.98, 1.02)
        p_i *= scaling_factor
        
        # Output Sample
        w_i = w_base + p_i
        baseline_samples[i] = w_i
        
    return baseline_samples

if __name__ == "__main__":
    # Simulate a realistic W_base for a voice print (416-dimensional vector for 13 MFCCs x 32 frames)
    # Using a normal distribution to mock an authentic seed state
    np.random.seed(42)
    w_base = np.random.normal(0, 1.5, 416)
    
    print("Generating 1000 Perturbative Samples based on Gemma's logic...")
    samples = generate_baseline_samples(w_base, num_samples=1000)
    
    # Validation step as requested by Gemma
    mu_baseline = np.mean(samples, axis=0)
    covariance_baseline = np.cov(samples, rowvar=False)
    
    # Save the 1000 samples to an artifact for the user
    # To save space, we just store the semanticWeight and sigmaId like the pipeline does
    deltas = []
    for idx, w_i in enumerate(samples):
        hash_val = hashlib.sha256(w_i.tobytes()).hexdigest()[:16]
        delta = {
            "sigmaId": f"biomarker/baseline_{hash_val}",
            "semanticWeight": float(np.mean(w_i)),
            "confidence": 1.0, # Ground truth baseline
            "deltaType": "BASELINE",
            "sampleIndex": idx
        }
        deltas.append(delta)
        
    with open("baseline_ml_weights_1000.json", "w") as f:
        json.dump(deltas, f, indent=2)
        
    # Check deviations
    mean_deviation = np.mean(np.abs(mu_baseline - w_base))
    
    print(f"Completed 1000 samples.")
    print(f"Mean deviation from W_base: {mean_deviation:.6f}")
    print(f"Covariance matrix shape: {covariance_baseline.shape}")
