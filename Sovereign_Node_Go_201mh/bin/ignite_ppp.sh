#!/data/data/com.termux/files/usr/bin/bash
# Clean slate
/data/data/com.termux/files/usr/bin/pkill -9 pppd 2>/dev/null

# Ignite PPP over SSH Tunnel
# Binding local IP 10.0.39.203 to Hub IP 10.0.39.1 via encrypted PTY
/data/data/com.termux/files/usr/bin/pppd pty "/data/data/com.termux/files/usr/bin/ssh -i /data/data/com.termux/files/home/.ssh/id_rsa root@204.168.138.60 pppd notty noauth" 10.0.39.203:10.0.39.1 noauth local nodefaultroute updetach >> /data/data/com.termux/files/home/Sovereign_Node_Go/logs/ppp_mesh_local.log 2>&1 &

# Wait for L3 handshake
/data/data/com.termux/files/usr/bin/sleep 8

# Verify Interface
if /data/data/com.termux/files/usr/bin/ifconfig ppp0 >> /data/data/com.termux/files/home/Sovereign_Node_Go/logs/ppp_mesh_local.log 2>&1; then
    echo "Success 0 PROOF OF EXECUTION"
    /data/data/com.termux/files/usr/bin/tail -n 15 /data/data/com.termux/files/home/Sovereign_Node_Go/logs/ppp_mesh_local.log
else
    echo "Error 1 Capture error RT ticket provide fix"
    /data/data/com.termux/files/usr/bin/tail -n 15 /data/data/com.termux/files/home/Sovereign_Node_Go/logs/ppp_mesh_local.log
fi
