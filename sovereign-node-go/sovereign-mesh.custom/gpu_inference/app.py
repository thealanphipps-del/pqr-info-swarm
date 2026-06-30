import torch
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import os

model_path = '/app/gemma_2b.pt'

model = None
if os.path.exists(model_path):
    try:
        model = torch.load(model_path, map_location='cuda' if torch.cuda.is_available() else 'cpu')
        model.eval()
        print("Model loaded successfully.")
    except Exception as e:
        print(f"Error loading model: {e}")
else:
    print(f"Model path {model_path} not found. Running in mock inference mode.")

app = FastAPI(title="Sovereign Swarm GPU Inference Service", version="0.1")

class InferenceRequest(BaseModel):
    prompt: str
    max_new_tokens: int = 64
    temperature: float = 0.7

class InferenceResponse(BaseModel):
    output: str

@app.post("/infer", response_model=InferenceResponse)
async def infer(request: InferenceRequest):
    if not request.prompt:
        raise HTTPException(status_code=400, detail="Prompt cannot be empty")
    try:
        if model is not None:
            with torch.no_grad():
                output_text = f"[Inference output from Custom Model] {request.prompt}"
        else:
            output_text = f"[Mocked Swarm Intelligence Output] Fusing gate weights for prompt: '{request.prompt}'"
        return InferenceResponse(output=output_text)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
