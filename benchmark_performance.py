import urllib.request
import urllib.parse
import json
import time
import statistics

def run_benchmark():
    base_url = "http://localhost:8080"
    agent_id = "benchmark-agent"
    iterations = 25
    
    print("PQR Ticketing System - Benchmark Performance Test")
    print("==================================================")
    print(f"Base URL: {base_url}")
    print(f"Iterations per operation: {iterations}\n")
    
    # Helper for JSON requests with built-in sleep to avoid rate limits
    def request(path, method="GET", data=None):
        time.sleep(0.08) # Sleep 80ms to ensure we stay below 12.5 req/sec
        url = f"{base_url}{path}"
        req = urllib.request.Request(url, method=method)
        if data is not None:
            req.add_header("Content-Type", "application/json")
            req_data = json.dumps(data).encode("utf-8")
        else:
            req_data = None
            
        start = time.perf_counter()
        try:
            with urllib.request.urlopen(req, data=req_data) as response:
                response.read()
                latency = (time.perf_counter() - start) * 1000 # ms
                return latency, True
        except Exception as e:
            latency = (time.perf_counter() - start) * 1000 # ms
            return latency, False

    # 1. Health check benchmark
    print("Benchmarking Health Check (GET /REST/2.0/health)...")
    health_latencies = []
    for _ in range(iterations):
        latency, success = request("/REST/2.0/health")
        if success:
            health_latencies.append(latency)
    
    # 2. Ticket creation benchmark
    print("Benchmarking Ticket Creation (POST /REST/2.0/ticket)...")
    ticket_latencies = []
    ticket_ids = []
    for i in range(iterations):
        time.sleep(0.08)
        payload = {
            "Subject": f"Benchmark Task {i}",
            "Queue": "benchmark",
            "Text": "Testing latency of ticket creation endpoint under load.",
            "AgentID": agent_id,
            "Layer": 1,
            "Intent": {"iteration": i}
        }
        url = f"{base_url}/REST/2.0/ticket"
        req = urllib.request.Request(url, method="POST")
        req.add_header("Content-Type", "application/json")
        req_data = json.dumps(payload).encode("utf-8")
        
        start = time.perf_counter()
        try:
            with urllib.request.urlopen(req, data=req_data) as response:
                body = json.loads(response.read().decode('utf-8'))
                latency = (time.perf_counter() - start) * 1000
                ticket_latencies.append(latency)
                ticket_ids.append(body["id"])
        except Exception as e:
            print(f"Ticket creation failed at index {i}: {e}")
            if hasattr(e, 'read'):
                print(e.read().decode('utf-8'))

    if not ticket_ids:
        print("Failed to create tickets for benchmark.")
        return

    # 3. Store memory benchmark
    print("Benchmarking Memory Storage (POST /REST/2.0/agent/:id/memory/:t)...")
    store_latencies = []
    for i, t_id in enumerate(ticket_ids):
        payload = {
            "memory_type": "context",
            "data": {"status": "active", "value": i},
            "relevance_score": 0.95
        }
        latency, success = request(f"/REST/2.0/agent/{agent_id}/memory/{t_id}", method="POST", data=payload)
        if success:
            store_latencies.append(latency)

    # 4. Get memory benchmark
    print("Benchmarking Memory Retrieval (GET /REST/2.0/agent/:id/memory/:t)...")
    get_latencies = []
    for t_id in ticket_ids:
        latency, success = request(f"/REST/2.0/agent/{agent_id}/memory/{t_id}?type=context")
        if success:
            get_latencies.append(latency)

    # 5. Update ticket benchmark
    print("Benchmarking Ticket Update (PUT /REST/2.0/ticket/:t)...")
    update_latencies = []
    for i, t_id in enumerate(ticket_ids):
        payload = {
            "Status": "COMPLETED",
            "Title": f"Benchmark Task {i} (Done)"
        }
        latency, success = request(f"/REST/2.0/ticket/{t_id}", method="PUT", data=payload)
        if success:
            update_latencies.append(latency)

    # 6. Context retrieval benchmark
    print("Benchmarking Context Retrieval (GET /REST/2.0/agent/:id/context)...")
    context_latencies = []
    for _ in range(iterations):
        latency, success = request(f"/REST/2.0/agent/{agent_id}/context")
        if success:
            context_latencies.append(latency)

    # 7. Audit trail query benchmark
    print("Benchmarking Audit Trail (GET /REST/2.0/ticket/:t/audit)...")
    audit_latencies = []
    for t_id in ticket_ids:
        latency, success = request(f"/REST/2.0/ticket/{t_id}/audit")
        if success:
            audit_latencies.append(latency)

    # Report helper
    def print_stats(name, latencies):
        if not latencies:
            print(f"{name:<45} | No data")
            return
        avg = statistics.mean(latencies)
        p95 = statistics.quantiles(latencies, n=20)[18] if len(latencies) >= 20 else max(latencies)
        min_val = min(latencies)
        max_val = max(latencies)
        print(f"{name:<45} | Min: {min_val:6.2f}ms | Max: {max_val:6.2f}ms | Avg: {avg:6.2f}ms | p95: {p95:6.2f}ms")

    print("\n==========================================================================")
    print("PQR API Performance Benchmark Results")
    print("==========================================================================")
    print_stats("GET  /health", health_latencies)
    print_stats("POST /ticket (Create Ticket)", ticket_latencies)
    print_stats("POST /agent/:id/memory/:t (Store)", store_latencies)
    print_stats("GET  /agent/:id/memory/:t (Get)", get_latencies)
    print_stats("PUT  /ticket/:t (Update Status)", update_latencies)
    print_stats("GET  /agent/:id/context (List Context)", context_latencies)
    print_stats("GET  /ticket/:t/audit (Audit Trail)", audit_latencies)
    print("==========================================================================")

if __name__ == "__main__":
    run_benchmark()
