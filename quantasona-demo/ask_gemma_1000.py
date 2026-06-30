import urllib.request
import json
import sys

url = "http://192.168.12.110:4111/v1/chat/completions"
prompt = """
Gemma, Alan has just requested: "run 1000 samples for baseline ml weights".

As his proxy, I am seeking your authorization and architectural guidance on how to proceed.
Should I simulate 1000 patient voice prints using statistical variance (e.g., Gaussian noise around the real audio sample's weights) to establish a baseline distribution, or is there a specific ML weighting script you want me to write to calibrate the BiometricPipeline?

Please provide the explicit parameters or the script logic for generating and aggregating these 1000 baseline samples.
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
