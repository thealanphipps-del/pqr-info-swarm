import grpc
import sync_pb2
import sync_pb2_grpc
import base64
import os
import sys

# Ensure we can import sync_pb2
sys.path.append("grpc_node")


def push_file(host, port, local_path, remote_path):
    with open(local_path, "rb") as f:
        content = base64.b64encode(f.read()).decode("utf-8")

    channel = grpc.insecure_channel(f"{host}:{port}")
    stub = sync_pb2_grpc.AgentSyncStub(channel)

    print(f"Pushing {local_path} -> {remote_path}...")
    # Use chunks if file is large (like index.html)
    chunk_size = 5000
    if len(content) > chunk_size:
        stub.RemoteExecute(
            sync_pb2.CommandPayload(command="bash", args=["-c", f"> {remote_path}.tmp"])
        )
        for i in range(0, len(content), chunk_size):
            chunk = content[i : i + chunk_size]
            stub.RemoteExecute(
                sync_pb2.CommandPayload(
                    command="bash", args=["-c", f"echo '{chunk}' >> {remote_path}.tmp"]
                )
            )
        stub.RemoteExecute(
            sync_pb2.CommandPayload(
                command="bash",
                args=[
                    "-c",
                    f"cat {remote_path}.tmp | base64 -d > {remote_path} && rm {remote_path}.tmp",
                ],
            )
        )
    else:
        res = stub.RemoteExecute(
            sync_pb2.CommandPayload(
                command="bash",
                args=["-c", f"echo '{content}' | base64 -d > {remote_path}"],
            )
        )
        if res.exit_code != 0:
            print(f"Error pushing {local_path}: {res.stderr}")


if __name__ == "__main__":
    # Target the reverse tunnel via 39.mh -> local:1112
    host = "localhost"
    port = 1112
    files = [
        ("proto/sync.proto", "/home/billing/sovereign-mesh/proto/sync.proto"),
        (
            "grpc_node/grpc_server.py",
            "/home/billing/sovereign-mesh/grpc_node/grpc_server.py",
        ),
        ("grpc.go", "/home/billing/sovereign-mesh/grpc.go"),
        ("dao.go", "/home/billing/sovereign-mesh/dao.go"),
        ("types.go", "/home/billing/sovereign-mesh/types.go"),
        ("sovereign.go", "/home/billing/sovereign-mesh/sovereign.go"),
        ("radius.go", "/home/billing/sovereign-mesh/radius.go"),
        ("dns.go", "/home/billing/sovereign-mesh/dns.go"),
        ("cloud.go", "/home/billing/sovereign-mesh/cloud.go"),
        (
            "grpc_node/web_server.py",
            "/home/billing/sovereign-mesh/grpc_node/web_server.py",
        ),
        ("grpc_node/index.html", "/home/billing/sovereign-mesh/grpc_node/index.html"),
        ("grpc_node/mgsh_mcp.py", "/home/billing/sovereign-mesh/grpc_node/mgsh_mcp.py"),
        ("oob.sh", "/home/billing/sovereign-mesh/oob.sh"),
    ]

    for local, remote in files:
        push_file(host, port, local, remote)

    print("Triggering recompile and restart...")
    channel = grpc.insecure_channel(f"{host}:{port}")
    stub = sync_pb2_grpc.AgentSyncStub(channel)
    stub.RemoteExecute(
        sync_pb2.CommandPayload(
            command="bash", args=["/home/billing/sovereign-mesh/proto/compile_proto.sh"]
        )
    )
    stub.RemoteExecute(sync_pb2.CommandPayload(command="pkill", args=["-f", "python3"]))
    print("Remote server updated and restarting.")
