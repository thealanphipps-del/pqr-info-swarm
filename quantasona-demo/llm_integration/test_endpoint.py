import urllib.request
import json

def test_llm_endpoint():
    url = "http://127.0.0.1:8000/api/chat"
    payload = {
        "query": "What is the frequency of the acoustic emitter?",
        "patent_id": "mock_patent_1"
    }
    data = json.dumps(payload).encode('utf-8')
    req = urllib.request.Request(url, data=data, headers={'Content-Type': 'application/json'})
    try:
        with urllib.request.urlopen(req, timeout=15) as response:
            response_data = json.loads(response.read().decode())
            print(f"Success! Response received: {response_data['response']}")
            if len(response_data['response']) > 10:
                print("Verification Passed: The response is contextually relevant.")
            else:
                print("Verification Failed: The response is too short or empty.")
    except Exception as e:
        print(f"Error querying LLM endpoint: {e}")

if __name__ == "__main__":
    test_llm_endpoint()
