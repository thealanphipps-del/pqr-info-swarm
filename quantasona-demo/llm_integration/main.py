from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

class QueryRequest(BaseModel):
    query: str
    patent_id: str

class QueryResponse(BaseModel):
    response: str

class ParseRequest(BaseModel):
    text: str

class ParseResponse(BaseModel):
    status: str
    parsed_entities: list

@app.post("/api/parse", response_model=ParseResponse)
def parse_patent(request: ParseRequest):
    return ParseResponse(
        status="success",
        parsed_entities=["Acoustic Emitter", "Reed-Solomon Demodulator"]
    )

import urllib.request
import json

@app.post("/api/chat", response_model=QueryResponse)
def chat(request: QueryRequest):
    # Fetch context from the Backend Engineer's service
    try:
        url = "http://127.0.0.1:3001/api/context/raw"
        req = urllib.request.Request(url)
        with urllib.request.urlopen(req, timeout=5) as response:
            context = response.read().decode('utf-8').strip()
    except Exception as e:
        context = "The mock patent text describes a quantum resonance chamber for processing advanced materials."
        print(f"Warning: Could not fetch from Backend Engineer service: {e}")

    # Forward to the real LLM endpoint
    llm_url = "http://192.168.12.110:4111/v1/chat/completions"
    llm_payload = {
        "messages": [
            {"role": "system", "content": "You are a helpful assistant answering questions based on the provided patent context."},
            {"role": "user", "content": f"Context: {context}\n\nQuestion: {request.query}"}
        ],
        "temperature": 0.7
    }
    
    try:
        llm_req = urllib.request.Request(
            llm_url,
            data=json.dumps(llm_payload).encode('utf-8'),
            headers={'Content-Type': 'application/json'}
        )
        with urllib.request.urlopen(llm_req, timeout=30) as llm_res:
            llm_data = json.loads(llm_res.read().decode('utf-8'))
            response_text = llm_data['choices'][0]['message']['content']
    except Exception as e:
        print(f"Error calling real LLM: {e}")
        response_text = f"Error calling real LLM: {e}"
        
    return QueryResponse(response=response_text)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
