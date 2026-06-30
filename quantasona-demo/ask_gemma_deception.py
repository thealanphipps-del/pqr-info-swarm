import urllib.request
import json
import sys

url = "http://192.168.12.110:4111/v1/chat/completions"
prompt = """
Gemma, Alan has issued a new major directive:
"search for cd quality wav files or no loss compression formats of real people talking preferably telling lies and truth mixed together we are going to evaluate them for nutrient deficency and pathogens and/or parasites you can use reality tv shows as a last resort after exhausting all sources first"

I have conducted a preliminary web search and found the following academic datasets for deception detection:
1. **Real-life Deception Detection Dataset (RLDD)**: Courtroom trial videos (truths and lies from defendants/witnesses). Audio can be extracted to lossless FLAC/WAV.
2. **DOLOS Dataset**: 1,675 clips from reality-TV gameshows featuring truths/lies. (Alan mentioned reality TV as a last resort, so this is a fallback).
3. **Bag-of-Lies (BoL)**: Multimodal dataset often used for spectral analysis.

Given Alan's objective to evaluate these audio signatures for "nutrient deficiency and pathogens and/or parasites" using the BiometricPipeline, which dataset do you authorize me to source? 
As his proxy, please provide explicit instructions on how I should proceed with acquiring and processing this data, or if you need me to adjust the search parameters.
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
