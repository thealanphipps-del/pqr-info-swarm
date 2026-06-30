import json
import hashlib
import os

import argparse

parser = argparse.ArgumentParser()
parser.add_argument("--input_file", default="real_voice_sample.wav")
parser.add_argument("--output_prefix", default="real_voiceprint")
args = parser.parse_args()

print("Importing audio libraries...")
import librosa
import numpy as np

# Real Voice Sample
audio_file = args.input_file
output_file = f"{args.output_prefix}.json"

try:
    print(f"Loading {audio_file}...")
    # Load audio at 16kHz
    y, sr = librosa.load(audio_file, sr=16000)
    
    print("Extracting MFCCs...")
    # Extract MFCCs (Voice Print)
    mfcc = librosa.feature.mfcc(y=y, sr=sr, n_mfcc=13)
    
    # Generate mock epistemic deltas
    deltas = []
    
    print("Anonymizing (Double Blind) and structuring data...")
    # Process the audio in 32-frame windows (roughly ~0.5 second slices depending on hop length)
    # We will process the first 50 windows (25 seconds of audio) for testing to keep the artifact small
    for i in range(0, min(mfcc.shape[1], 32 * 50), 32): 
        window = mfcc[:, i:i+32]
        if window.shape[1] < 32: continue
        
        # Flatten and hash for anonymization (Double Blind)
        flat_embedding = window.flatten()
        hash_val = hashlib.sha256(flat_embedding.tobytes()).hexdigest()[:16]
        
        delta = {
            "sigmaId": f"biomarker/anon_{hash_val}",
            "semanticWeight": float(np.mean(flat_embedding)),
            "confidence": 0.95,
            "provenance": "anonymous-clinical-interview-001",
            "deltaType": "OBSERVATION",
            "relationType": "INTRODUCES",
            "sourceMic": "real-historical-audio-16kHz"
        }
        deltas.append(delta)

    output_file = "anonymized_real_voice_prints.json"
    with open(output_file, "w") as f:
        json.dump(deltas, f, indent=2)

    print(f"Successfully generated {len(deltas)} double-blind voice prints from REAL audio.")
    print(f"Saved to {output_file}")

except Exception as e:
    print(f"Error processing real voice audio: {e}")
