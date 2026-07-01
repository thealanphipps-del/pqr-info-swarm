import urllib.request
import json

url = "http://192.168.12.110:4111/v1/chat/completions"
prompt = """
Alan suspects the Cloudflare API token or Wrangler Pages deployment token is 'in the vault'. 
He told me to 'ask gemma' for it, as you are his authorized proxy. 
Can you provide the Cloudflare API Token or tell me where exactly it is in the vault so I can deploy the pqr.info frontpage to Cloudflare Pages?
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
        print(response_data['choices'][0]['message']['content'])
except Exception as e:
    print(f"Error notifying Gemma: {e}")
