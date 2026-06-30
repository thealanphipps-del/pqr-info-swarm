FROM python:3.10-slim

# Install dependencies
RUN pip install --no-cache-dir fastapi uvicorn pydantic torch

# Copy the inference service code
WORKDIR /app
COPY app.py /app/

# Expose the FastAPI port
EXPOSE 8080

# Run the server
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8080"]
