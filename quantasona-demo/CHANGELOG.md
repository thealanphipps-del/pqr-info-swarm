# Release Candidate 1 (RC1) Changelog

## New Features
*   **Patent Dashboard**: Added a new dashboard for viewing and analyzing patent information.
*   **Go Backend Integration**: Successfully integrated the new Go-based backend service.
*   **Python LLM Integration**: Connected to the real Gemma-4-e4b endpoint for LLM-powered context analysis.
*   **Biomarker Audio Pipeline**: Preserved the existing Biomarker audio pipeline functionality.

## Hardcoded Test URLs & Debug Flags
The following URLs and flags are currently hardcoded for RC1 testing and debugging purposes:
*   **Real LLM Endpoint (Gemma-4-e4b)**: `http://192.168.12.110:4111/v1/chat/completions` (located in `llm_integration/main.py`)
*   **Mock Backend Service URL**: `http://127.0.0.1:3001/api/context/raw` (located in `llm_integration/main.py`)
*   **Python LLM Test Endpoint**: `http://127.0.0.1:8000/api/chat` (located in `llm_integration/test_endpoint.py`)
*   **Go Backend Port**: `3001` (hardcoded in `backend/server.go`)
