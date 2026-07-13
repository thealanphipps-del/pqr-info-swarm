#!/data/data/com.termux/files/usr/bin/bash
LOG_FILE="/data/data/com.termux/files/home/Sovereign_Node_Go/logs/mesh_heal_repair.log"

# Function specifically for downstream WireGuard nodes via Helsinki
heal_wg_node() {
    local node_label=$1
    local wg_ip=$2
    local local_port=$3
    local gateway_ip="204.168.138.60"

    # We probe the mapped port directly via nc, avoiding SSH config parsing issues
    if /data/data/com.termux/files/usr/bin/nc -zv -w 3 127.0.0.1 "$local_port" >/dev/null 2>&1; then
        echo "[$(date -u)] $node_label ($wg_ip) ALIVE VIA PORT $local_port" >> "$LOG_FILE"
    else
        echo "[$(date -u)] $node_label ($wg_ip) UNREACHABLE. REBUILDING FORWARD..." >> "$LOG_FILE"
        
        # Kill the specific port forward process
        /data/data/com.termux/files/usr/bin/pkill -f "$local_port:$wg_ip:22" || true
        
        # Re-establish Local Port Forward
        /data/data/com.termux/files/usr/bin/ssh -F /dev/null -o IPQoS=throughput -o ControlMaster=no -o StrictHostKeyChecking=no -o ConnectTimeout=15 -o ServerAliveInterval=60 -i /data/data/com.termux/files/home/.ssh/id_rsa -L "$local_port":"$wg_ip":22 -f -N root@"$gateway_ip"
        
        sleep 5
        
        if /data/data/com.termux/files/usr/bin/nc -zv -w 3 127.0.0.1 "$local_port" >/dev/null 2>&1; then
            echo "[$(date -u)] $node_label HEAL VERIFIED: FORWARD ACTIVE" >> "$LOG_FILE"
        else
            echo "[$(date -u)] $node_label FATAL HEAL FAILURE: PORT $local_port STALLED" >> "$LOG_FILE"
        fi
    fi
}

# Probe the gateway directly
if /data/data/com.termux/files/usr/bin/ssh -F /dev/null -o IPQoS=throughput -o ControlMaster=no -o StrictHostKeyChecking=no -o ConnectTimeout=5 -o BatchMode=yes -i /data/data/com.termux/files/home/.ssh/id_rsa root@204.168.138.60 "exit 0" >/dev/null 2>&1; then
    echo "[$(date -u)] 39.mh GATEWAY ALIVE" >> "$LOG_FILE"
else
    echo "[$(date -u)] 39.mh GATEWAY FATAL OUTAGE" >> "$LOG_FILE"
fi

# Heal downstream nodes mapped to local ports
# heal_wg_node "201.mh" "10.8.0.2" "8081"
