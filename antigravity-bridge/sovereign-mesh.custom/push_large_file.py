import grpc
import sync_pb2
import sync_pb2_grpc
import base64
import os
import sys

# Ensure we can import sync_pb2
sys.path.append("grpc_node")


def push_file_chunks(host, port, local_path, remote_path):
    with open(local_path, "rb") as f:
        content = base64.b64encode(f.read()).decode("utf-8")

    channel = grpc.insecure_channel(f"{host}:{port}")
    stub = sync_pb2_grpc.AgentSyncStub(channel)

    # Clear remote file first
    print(f"Clearing remote {remote_path}...")
    stub.RemoteExecute(
        sync_pb2.CommandPayload(command="bash", args=["-c", f"> {remote_path}"])
    )

    chunk_size = 5000  # Small enough to avoid "Argument list too long"
    for i in range(0, len(content), chunk_size):
        chunk = content[i : i + chunk_size]
        print(f"Pushing chunk {i//chunk_size + 1}...")
        res = stub.RemoteExecute(
            sync_pb2.CommandPayload(
                command="bash", args=["-c", f"echo {chunk} >> {remote_path}.tmp"]
            )
        )
        if res.exit_code != 0:
            print(f"Error: {res.stderr}")
            return

    print("Finalizing file...")
    stub.RemoteExecute(
        sync_pb2.CommandPayload(
            command="bash",
            args=[
                "-c",
                f"cat {remote_path}.tmp | base64 -d > {remote_path} && rm {remote_path}.tmp",
            ],
        )
    )
    print("Success.")


if __name__ == "__main__":
    host = "136.113.240.237"
    port = 1111
    push_file_chunks(
        host,
        port,
        "grpc_node/index.html",
        "/home/billing/sovereign-mesh/grpc_node/index.html",
    )

    print("Triggering restart...")
    channel = grpc.insecure_channel(f"{host}:{port}")
    stub = sync_pb2_grpc.AgentSyncStub(channel)
    stub.RemoteExecute(sync_pb2.CommandPayload(command="pkill", args=["-f", "python3"]))
    print("Remote server updated.")
