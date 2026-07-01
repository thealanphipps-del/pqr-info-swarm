import urllib.request
import json
import sys

url = "http://192.168.12.110:4111/v1/chat/completions"
prompt = """
Hello Gemma! We are running a memory test. 
Can you recall our previous conversation regarding Alan's request for the Cloudflare API token and the HashiCorp vault? 
Please summarize what you told me about how to securely retrieve tokens.
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
        print("Gemma responded:\n", response_data['choices'][0]['message']['content'])
except Exception as e:
    print(f"Error testing Gemma memory: {e}")
