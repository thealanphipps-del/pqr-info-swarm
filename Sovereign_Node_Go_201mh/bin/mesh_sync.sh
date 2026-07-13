#!/data/data/com.termux/files/usr/bin/bash
LOG="/data/data/com.termux/files/home/Sovereign_Node_Go/ghub/mesh_sync.log"
FP="/data/data/com.termux/files/home/Sovereign_Node_Go/bin/find_port"

echo "[_] Initializing Mesh Recovery..." > $LOG
#A. Deconflict
PIDS=$(/data/data/com.termux/files/usr/bin/fuser 1080/tcp 8888/tcp 5432/tcp 4450/tcp 2>/dev/null)
if [ -n "$PIDS" ]; then kill -9 $PIDS; fi

#B. Re-establish Tunnels
/data/data/com.termux/files/usr/bin/ssh -i /data/data/com.termux/files/home/.ssh/id_rsa -D 1080 -L 5432:127.0.0.1:5432 -L 4450:127.0.0.1:445 -L 8888:localhost:8888 root@204.168.138.60 -N -f

sleep 2

#C. Re-ignite Daemons
/data/data/com.termux/files/usr/bin/pkill -9 director hud 2>/dev/null
/data/data/com.termux/files/home/Sovereign_Node_Go/bin/director >> $LOG 2>&1 &
/data/data/com.termux/files/home/Sovereign_Node_Go/bin/hud >> $LOG 2>&1 &

echo "[SUCCESS] Mesh Recovered. Ports: 8080, 8888, 1080, 5432, 4450" >> $LOG
