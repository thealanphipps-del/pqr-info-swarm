#!/usr/bin/env python3
import grpc
import sys
import os
import json
import time

# Add grpc_node to path for stubs
sys.path.append(os.path.join(os.path.dirname(__file__), "..", "grpc_node"))
import mesh_proto_pb2 as mesh_proto
import mesh_proto_pb2_grpc as mesh_proto_grpc

# Load agent identity map
AGENT_MAP_PATH = os.path.join(os.path.dirname(__file__), "..", "agent_identity_map.json")

def birth_agents():
    print("🧬 SOVEREIGN SWARM: GENESIS HEARTBEAT INITIATION")
    
    if not os.path.exists(AGENT_MAP_PATH):
        print(f"Error: Agent map not found at {AGENT_MAP_PATH}")
        return

    with open(AGENT_MAP_PATH, "r") as f:
        identities = json.load(f)

    # Dial the Go mesh_server on Port 1113
    channel = grpc.insecure_channel('localhost:1113')
    stub = mesh_proto_grpc.SovereignMeshStub(channel)

    print("Entering heartbeat loop... (Press Ctrl+C to stop)")
    while True:
        count = 0
        for agent_name, shortcode in identities.items():
            # Every 5th agent is a VALIDATOR to satisfy infra requirement
            node_class = "VALIDATOR" if count % 5 == 0 else "AGENT"
            
            try:
                res = stub.Heartbeat(mesh_proto.AgentHeartbeat(
                    agent_id=agent_name,
                    address="localhost",
                    status="idle",
                    intelligence_level=7,
                    node_class=node_class
                ))
                if res.success:
                    count += 1
            except Exception as e:
                print(f"  ❌ gRPC error for {agent_name}: {e}")
            
        print(f"[{time.strftime('%H:%M:%S')}] Swarm Vitality: {count}/33 agents heartbeating.")
        time.sleep(30) # Send heartbeats every 30 seconds

if __name__ == "__main__":
    birth_agents()
