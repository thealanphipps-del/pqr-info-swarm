#!/usr/bin/env python3
import json
import socket
import threading
import time
from http.server import BaseHTTPRequestHandler, HTTPServer

# ----------------------------------------------------------------------
# Config
# ----------------------------------------------------------------------

CLUSTER_ID = "sovereign-alpha"
GENESIS_EPOCH = 1
BIND_HOST = "127.0.0.3"

PQR_BUS_PORT = 11111      # agent protocol, MPUDP, HyperFrames
SRRP_PORT = 11112         # SRRP transport tunnels (mocked)
MESHQUORUM_PORT = 1111    # governance truth / epoch sync
PBA_PORT = 1112           # Port Binding Agent

SUBSTRATE = {}            # in-memory mock: object_id -> record


# ----------------------------------------------------------------------
# Utilities
# ----------------------------------------------------------------------

def send_json(sock: socket.socket, obj: dict):
    data = json.dumps(obj).encode("utf-8")
    length = len(data).to_bytes(4, "big")
    sock.sendall(length + data)


def recv_json(sock: socket.socket) -> dict:
    length_bytes = sock.recv(4)
    if not length_bytes:
        return None
    length = int.from_bytes(length_bytes, "big")
    buf = b""
    while len(buf) < length:
        chunk = sock.recv(length - len(buf))
        if not chunk:
            break
        buf += chunk
    return json.loads(buf.decode("utf-8"))


# ----------------------------------------------------------------------
# Port Binding Agent (PBA) on 1112
# ----------------------------------------------------------------------

def pba_server():
    srv = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    srv.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    srv.bind((BIND_HOST, PBA_PORT))
    srv.listen(16)
    print("[PBA] listening on", PBA_PORT)

    while True:
        conn, addr = srv.accept()
        threading.Thread(target=handle_pba_client, args=(conn,), daemon=True).start()


def handle_pba_client(conn: socket.socket):
    try:
        req = recv_json(conn)
        if not req:
            return
        agent_id = req.get("agent_id", "unknown")
        print(f"[PBA] registration from {agent_id}")

        resp = {
            "agent_id": agent_id,
            "assigned_ports": {
                "pqr_bus": PQR_BUS_PORT,
                "srrp": SRRP_PORT,
                "control": MESHQUORUM_PORT,
            },
            "epoch": GENESIS_EPOCH,
        }
        send_json(conn, resp)
        print(f"[PBA] assigned ports to {agent_id}")
    finally:
        conn.close()


# ----------------------------------------------------------------------
# MeshQuorum (epoch sync) on 1111
# ----------------------------------------------------------------------

AGENT_ACKS = set()


class MeshQuorumHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        length = int(self.headers.get("Content-Length", "0"))
        body = self.rfile.read(length)
        msg = json.loads(body.decode("utf-8"))

        if msg.get("type") == "epoch_announce":
            print(f"[MeshQuorum] epoch announce: {msg}")
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b'{"status":"ok"}')
            return

        if msg.get("type") == "epoch_ack":
            agent_id = msg.get("agent_id")
            epoch = msg.get("epoch")
            print(f"[MeshQuorum] epoch ACK from {agent_id} for epoch {epoch}")
            AGENT_ACKS.add(agent_id)
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b'{"status":"ok"}')
            return

        self.send_response(400)
        self.end_headers()
        self.wfile.write(b'{"status":"error"}')


def meshquorum_server():
    httpd = HTTPServer((BIND_HOST, MESHQUORUM_PORT), MeshQuorumHandler)
    print("[MeshQuorum] listening on", MESHQUORUM_PORT)
    httpd.serve_forever()


# ----------------------------------------------------------------------
# PQR Bus on 11111 (HyperFrame, MPUDP, etc.)
# ----------------------------------------------------------------------

def pqr_bus_server():
    srv = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    srv.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    srv.bind((BIND_HOST, PQR_BUS_PORT))
    srv.listen(16)
    print("[PQR Bus] listening on", PQR_BUS_PORT)

    while True:
        conn, addr = srv.accept()
        threading.Thread(target=handle_pqr_bus_client, args=(conn,), daemon=True).start()


def handle_pqr_bus_client(conn: socket.socket):
    try:
        frame = recv_json(conn)
        if not frame:
            return
        msg_type = frame.get("msg_type")
        if msg_type == "hyperframe":
            handle_hyperframe(frame)
        else:
            print(f"[PQR Bus] unknown msg_type: {msg_type}")
    finally:
        conn.close()


def handle_hyperframe(frame: dict):
    hf = frame.get("hyperframe", {})
    object_id = hf.get("object_id")
    payload = hf.get("payload")
    epoch = hf.get("epoch")
    origin5d = hf.get("origin5d")
    vertex = hf.get("vertex")

    print(f"[PQR Bus] received HyperFrame object_id={object_id} epoch={epoch}")
    SUBSTRATE[object_id] = {
        "payload": payload,
        "epoch": epoch,
        "origin5d": origin5d,
        "vertex": vertex,
        "lineage": frame.get("lineage_hash"),
    }


# ----------------------------------------------------------------------
# SRRP transport (mock) on 11112
# ----------------------------------------------------------------------

def srrp_server():
    srv = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    srv.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    srv.bind((BIND_HOST, SRRP_PORT))
    srv.listen(16)
    print("[SRRP] listening on", SRRP_PORT)

    while True:
        conn, addr = srv.accept()
        threading.Thread(target=handle_srrp_client, args=(conn,), daemon=True).start()


def handle_srrp_client(conn: socket.socket):
    try:
        req = recv_json(conn)
        if not req:
            return
        print(f"[SRRP] mock route request: {req}")
        resp = {"status": "ok", "route": ["kernel-root"]}
        send_json(conn, resp)
    finally:
        conn.close()


# ----------------------------------------------------------------------
# Test agent: registration, epoch ACK, HyperFrame send, reconstruction
# ----------------------------------------------------------------------

def test_agent():
    agent_id = "agent-test-1"

    # 1) Register with PBA
    pba_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    pba_sock.connect((BIND_HOST, PBA_PORT))
    send_json(pba_sock, {
        "agent_id": agent_id,
        "role": "mesh-node",
        "capabilities": ["mpudp", "midi", "hyperframe", "srrp"],
    })
    pba_resp = recv_json(pba_sock)
    pba_sock.close()
    print("[TestAgent] PBA response:", pba_resp)

    # 2) Epoch ACK to MeshQuorum
    mq_payload = {
        "type": "epoch_ack",
        "agent_id": agent_id,
        "epoch": GENESIS_EPOCH,
        "status": "ready",
    }
    mq_data = json.dumps(mq_payload).encode("utf-8")
    mq_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    mq_sock.connect((BIND_HOST, MESHQUORUM_PORT))
    # Simple HTTP POST
    req = (
        "POST /epoch HTTP/1.1\r\n"
        f"Host: {BIND_HOST}:{MESHQUORUM_PORT}\r\n"
        "Content-Type: application/json\r\n"
        f"Content-Length: {len(mq_data)}\r\n"
        "\r\n"
    ).encode("utf-8") + mq_data
    mq_sock.sendall(req)
    mq_sock.close()
    print("[TestAgent] sent epoch ACK")

    time.sleep(0.5)

    # 3) Send HyperFrame over PQR Bus
    hf = {
        "frame_id": "frame-1",
        "object_id": "test-1",
        "epoch": GENESIS_EPOCH,
        "payload": "hello sovereign",
        "origin5d": {"x": 1, "y": 2, "z": 3, "phi": 4, "lambda": 5},
        "vertex": {"v0": 1, "v1": 2, "v2": 3, "v3": 4},
    }
    frame = {
        "version": 1,
        "msg_type": "hyperframe",
        "agent_id": agent_id,
        "epoch": GENESIS_EPOCH,
        "origin5d": hf["origin5d"],
        "vertex": hf["vertex"],
        "hyperframe": hf,
        "lineage_hash": "lineage-test-1",
    }

    bus_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    bus_sock.connect((BIND_HOST, PQR_BUS_PORT))
    send_json(bus_sock, frame)
    bus_sock.close()
    print("[TestAgent] sent HyperFrame test-1")

    time.sleep(0.5)

    # 4) Reconstruction check
    record = SUBSTRATE.get("test-1")
    print("[TestAgent] reconstruction record:", record)
    assert record is not None, "Substrate did not store test-1"
    assert record["payload"] == "hello sovereign", "Payload mismatch"
    assert record["epoch"] == GENESIS_EPOCH, "Epoch mismatch"
    print("[TestAgent] ✅ Sovereign OS round-trip verified for test-1")


# ----------------------------------------------------------------------
# Kernel orchestration
# ----------------------------------------------------------------------

def kernel_announce_epoch():
    payload = {
        "type": "epoch_announce",
        "epoch": GENESIS_EPOCH,
        "cluster_id": CLUSTER_ID,
        "timestamp": time.time(),
    }
    data = json.dumps(payload).encode("utf-8")
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((BIND_HOST, MESHQUORUM_PORT))
    req = (
        "POST /epoch HTTP/1.1\r\n"
        f"Host: {BIND_HOST}:{MESHQUORUM_PORT}\r\n"
        "Content-Type: application/json\r\n"
        f"Content-Length: {len(data)}\r\n"
        "\r\n"
    ).encode("utf-8") + data
    sock.sendall(req)
    sock.close()
    print("[Kernel] epoch announce sent")


def main():
    # Start servers
    threading.Thread(target=pba_server, daemon=True).start()
    threading.Thread(target=meshquorum_server, daemon=True).start()
    threading.Thread(target=pqr_bus_server, daemon=True).start()
    threading.Thread(target=srrp_server, daemon=True).start()

    time.sleep(0.5)

    # Kernel announces genesis epoch
    kernel_announce_epoch()

    # Run test agent
    test_agent()

    # Keep process alive a bit for logs
    time.sleep(1)


if __name__ == "__main__":
    main()
