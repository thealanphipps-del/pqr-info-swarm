#!/data/data/com.termux/files/usr/bin/bash
set -e

NODE_0_IP="46.224.84.64"
NODE_39_IP="204.168.138.60"

for port in 8081 8111 8112 5432 5433 9000 9111 9080; do
    /data/data/com.termux/files/usr/bin/fuser -k $port/tcp 2>/dev/null || true
done

/data/data/com.termux/files/usr/bin/ssh -i /data/data/com.termux/files/home/.ssh/id_rsa \
    -o ServerAliveInterval=30 \
    -o ServerAliveCountMax=5 \
    -L 8081:127.0.0.1:8081 \
    -L 8111:127.0.0.1:8111 \
    -L 8112:127.0.0.1:8112 \
    -L 5432:127.0.0.1:5432 \
    -L 5433:127.0.0.1:5433 \
    -L 9000:127.0.0.1:9000 \
    -L 9111:127.0.0.1:9111 \
    -L 9080:127.0.0.1:9080 \
    root@$NODE_39_IP
