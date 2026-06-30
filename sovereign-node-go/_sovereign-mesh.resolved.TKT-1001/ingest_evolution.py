import grpc
import sys
import os
import time

# Ensure we can import the generated proto files
sys.path.append("grpc_node")
import sync_pb2
import sync_pb2_grpc


def propose_mutation(host, port, key, value, reason):
    channel = grpc.insecure_channel(f"{host}:{port}")
    stub = sync_pb2_grpc.AgentSyncStub(channel)

    print(f"Ingesting into Time Machine: {key}...")
    try:
        res = stub.ProposeSwarmMutation(
            sync_pb2.MutationRequest(
                target_key=key,
                proposed_value=value,
                change_reason=reason,
                proposer_agent_id="EVOLUTION-INGEST-BOT",
            )
        )
        print(f"Status: {res.status} | Block Index: {res.block_index}")
        if not res.consensus_reached:
            print(f"Consensus Deficit: {res.consensus_ratio}")
    except Exception as e:
        print(f"Ingestion failed for {key}: {e}")


import argparse

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--host", default="136.113.240.237")
    parser.add_argument("--port", type=int, default=1111)
    args = parser.parse_args()

    host = args.host
    port = args.port

    # 1. Historical Context Ingestion
    propose_mutation(
        host,
        port,
        "HISTORICAL_ROOT",
        "STICKY-20260406-V1.2",
        "Tracing swarm ancestry to original V1.2 genesis state.",
    )
    propose_mutation(
        host,
        port,
        "STALL_MATRIX_ANALYSIS",
        "39_STALLED_TICKETS_RECOVERED",
        "Identified recurring authentication bottlenecks at node 201.mh. Applied self-healing logic.",
    )

    # 2. Economic Logic Ingestion
    propose_mutation(
        host,
        port,
        "SURFGO_SCARCITY_PROTOCOL",
        "HARD_CAP_27M_BURN_0.01_PERCENT",
        "Enforcing strict scarcity and supply contraction for the Sovereign economy.",
    )
    propose_mutation(
        host,
        port,
        "ARB_FEED_SYNC",
        "NUREMBERG-HELSINKI-DIRECT-TUNNEL",
        "Synchronized Nuremberg (0.mh) and Helsinki (39.mh) for sub-second arbitrage execution.",
    )

    # 3. Tooling & Architecture Evolution
    propose_mutation(
        host,
        port,
        "TOOL_MAP_L7",
        "AGENT_ARB_HUD_MCP_DEPLOYED",
        "Mapped the 7-layer tool suite from rtgo_public/bin into active agent procedural memory.",
    )
    propose_mutation(
        host,
        port,
        "OUROBOROS_SENTINEL_ACTIVE",
        "SELF_HEALING_V1.0",
        "Crystallizing the Ouroboros Sentinel as the primary process protection daemon.",
    )

    print("--- Jetweb Time Machine Ingestion Sequence Complete ---")
