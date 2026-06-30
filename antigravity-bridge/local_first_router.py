#!/usr/bin/env python3
"""
ANTIGRAVITY LOCAL-FIRST ROUTER
KLO-PROTOCOL | ZERO-INFERENCE GROUNDING
Callsign: AURORA-R9 -> SENTRY-TWO -> GEMINI GATEWAY

Priority chain:
  1. Local Ollama (localhost:11434)
  2. S25 Training Agent Ollama (192.168.12.201:11434)
  3. S25 Bridge endpoint (192.168.12.201:8081)
  4. External Gemini API (last resort)
"""

import http.client
import json
import sys
import time
import os

# ─── Config ────────────────────────────────────────────────────────────────
LOCAL_OLLAMA    = ("127.0.0.1", 11434)
S25_OLLAMA      = ("192.168.12.201", 11434)
S25_BRIDGE      = ("192.168.12.201", 8081)
S25_CALLSIGN    = "AELLK"
DEFAULT_MODEL   = os.environ.get("LOCAL_MODEL", "gemma3:1b")
TIMEOUT         = 8  # seconds per tier

# ─── Helpers ────────────────────────────────────────────────────────────────
def _post_json(host, port, path, payload, timeout=TIMEOUT, extra_headers=None):
    """POST JSON payload and return parsed response or None on failure."""
    try:
        conn = http.client.HTTPConnection(host, port, timeout=timeout)
        headers = {"Content-Type": "application/json"}
        if extra_headers:
            headers.update(extra_headers)
        conn.request("POST", path, json.dumps(payload), headers)
        resp = conn.getresponse()
        data = resp.read().decode()
        conn.close()
        if resp.status == 200:
            return json.loads(data)
    except Exception as e:
        print(f"  [ROUTER] {host}:{port}{path} failed: {e}", file=sys.stderr)
    return None

def _get(host, port, path, timeout=TIMEOUT):
    """GET request, returns text or None."""
    try:
        conn = http.client.HTTPConnection(host, port, timeout=timeout)
        conn.request("GET", path)
        resp = conn.getresponse()
        data = resp.read().decode()
        conn.close()
        if resp.status == 200:
            return data
    except Exception as e:
        print(f"  [ROUTER] GET {host}:{port}{path} failed: {e}", file=sys.stderr)
    return None

def lmstudio_chat(host, port, prompt, model="google/gemma-4-e4b"):
    """Query an LM Studio chat completions endpoint."""
    payload = {
        "model": model,
        "messages": [{"role": "user", "content": prompt}],
        "stream": False,
        "temperature": 0.0
    }
    result = _post_json(host, port, "/v1/chat/completions", payload)
    if result and "choices" in result and len(result["choices"]) > 0:
        return result["choices"][0]["message"]["content"].strip()
    return None

def ollama_generate(host, port, prompt, model=DEFAULT_MODEL):
    """Query an Ollama /api/generate endpoint."""
    payload = {"model": model, "prompt": prompt, "stream": False}
    result = _post_json(host, port, "/api/generate", payload)
    if result and "response" in result:
        return result["response"].strip()
    return None

def s25_bridge_query(prompt):
    """Ask the S25 sovereign bridge endpoint."""
    import urllib.parse
    cmd = f"echo {urllib.parse.quote(prompt)}"
    result = _get(S25_BRIDGE[0], S25_BRIDGE[1], f"/api/bridge?cmd={urllib.parse.quote(prompt)}")
    return result.strip() if result else None

# ─── Main router ────────────────────────────────────────────────────────────
def route(prompt: str, context: str = "") -> dict:
    full_prompt = f"{context}\n\n{prompt}" if context else prompt
    tiers = [
        ("TIER-0 LM Studio (google/gemma-4-e4b)", lambda: lmstudio_chat("127.0.0.1", 1234, full_prompt, "google/gemma-4-e4b")),
        ("TIER-1 Local Ollama (smollm:135m)",      lambda: ollama_generate("127.0.0.1", 11434, full_prompt, "smollm:135m")),
        ("TIER-2 S25 Ollama (SENTRY-TWO)",  lambda: ollama_generate(*S25_OLLAMA,   full_prompt)),
        ("TIER-3 S25 Bridge",               lambda: s25_bridge_query(full_prompt)),
    ]

    for tier_name, fn in tiers:
        print(f"[ROUTER] Trying {tier_name}...", file=sys.stderr)
        t0 = time.time()
        response = fn()
        elapsed = round(time.time() - t0, 2)
        if response:
            print(f"[ROUTER] ✓ {tier_name} responded in {elapsed}s", file=sys.stderr)
            return {
                "response": response,
                "tier": tier_name,
                "latency_s": elapsed,
                "token_cost": 0
            }
        print(f"[ROUTER] ✗ {tier_name} miss ({elapsed}s)", file=sys.stderr)

    return {
        "response": None,
        "tier": "TIER-4 Gemini API (external)",
        "token_cost": "variable",
        "note": "All local tiers exhausted. Route to external Gemini API."
    }

# ─── CLI test harness ────────────────────────────────────────────────────────
if __name__ == "__main__":
    prompt = " ".join(sys.argv[1:]) if len(sys.argv) > 1 else "What is the status of the Aurora R9 sovereign mesh?"
    print(f"\n[ROUTER] Query: {prompt}\n")
    result = route(prompt)
    print(json.dumps(result, indent=2))
