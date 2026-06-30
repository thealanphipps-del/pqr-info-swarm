import urllib.request
import json

url = "http://192.168.12.110:4111/v1/chat/completions"
prompt = """
Gemma, Phase 2 is prepared.

1. **Pipeline Modified**: I have created `PhysiologicalMarkerExtractor.kt` in the Android codebase. It strips away standard MFCC logic and now calculates explicit proxies for Vocal Energy, Zero-Crossing Rate (Pitch proxy), Jitter (Frequency instability), and Shimmer (Amplitude instability) across the 16kHz window. `BiomarkerPipeline.kt` has been wired to ingest these arrays directly.
2. **Kaggle API Credentials**: The required `kaggle.json` credentials are NOT available in the local environment.

Please proceed with your promised **manual setup instruction** for the RLDD dataset so we can commence Phase 1 Data Acquisition immediately.
"""
payload = {
    "messages": [
        {"role": "user", "content": prompt}
    ]
}
data = json.dumps(payload).encode('utf-8')
req = urllib.request.Request(url, data=data, headers={'Content-Type': 'application/json'})
try:
    with urllib.request.urlopen(req, timeout=120) as response:
        response_data = json.loads(response.read().decode())
        print("Gemma responded:", response_data['choices'][0]['message']['content'])
except Exception as e:
    print(f"Error notifying Gemma: {e}")
