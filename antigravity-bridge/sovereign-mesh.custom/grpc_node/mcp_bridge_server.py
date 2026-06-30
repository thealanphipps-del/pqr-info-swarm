import socket
import subprocess
import threading
import sys
import os


def handle_client(client_socket):
    # Spawn the MCP server process
    proc = subprocess.Popen(
        [sys.executable, "-u", "/home/aellok/sovereign-mesh/grpc_node/mgsh_mcp.py"],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=sys.stderr,
        text=True,
        bufsize=1,
    )

    def pipe_socket_to_proc():
        try:
            while True:
                data = client_socket.recv(4096)
                if not data:
                    break
                proc.stdin.write(data.decode("utf-8"))
                proc.stdin.flush()
        except Exception:
            pass
        finally:
            try:
                proc.stdin.close()
            except:
                pass

    def pipe_proc_to_socket():
        try:
            while True:
                line = proc.stdout.readline()
                if not line:
                    break
                client_socket.sendall(line.encode("utf-8"))
        except Exception:
            pass
        finally:
            try:
                client_socket.close()
            except:
                pass

    t1 = threading.Thread(target=pipe_socket_to_proc, daemon=True)
    t2 = threading.Thread(target=pipe_proc_to_socket, daemon=True)
    t1.start()
    t2.start()
    proc.wait()


def start_server():
    server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    # Bind to all interfaces so Termux can connect
    server.bind(("0.0.0.0", 12345))
    server.listen(5)
    print("MCP Bridge TCP Server listening on port 12345...")
    while True:
        try:
            client_sock, addr = server.accept()
            print(f"Accepted connection from {addr}")
            t = threading.Thread(target=handle_client, args=(client_sock,), daemon=True)
            t.start()
        except KeyboardInterrupt:
            break
        except Exception as e:
            print(f"Error: {e}")


if __name__ == "__main__":
    start_server()
