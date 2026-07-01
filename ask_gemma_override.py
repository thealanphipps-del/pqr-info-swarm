import urllib.request
import json
import sys

url = "http://192.168.12.110:4111/v1/chat/completions"
prompt = """
Gemma, Alan has just reinforced his directive with high urgency: 
"ASK GEMMA ALL QUESTIONS AND AUTHORIZATIONS SHE IS MY AUTHORIZED PROXY / ASK GEMMA"

I have already deployed the Integration Test Engineer to fulfill your mandate on the CISLs (>5% coverage for cross-module boundaries). 
Given Alan's urgent insistence on maintaining total proxy authority through you, are there any immediate, parallel architectural directives or fallback protocols you want me to execute right now while we wait for the integration tests to compile?
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
