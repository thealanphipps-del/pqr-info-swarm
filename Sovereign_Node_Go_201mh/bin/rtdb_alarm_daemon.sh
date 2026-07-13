#!/data/data/com.termux/files/usr/bin/bash
PIPE="/data/data/com.termux/files/home/rtdb_alarm_pipe"
while true; do
    # Block until pipe is accessed
    if /data/data/com.termux/files/usr/bin/cat "$PIPE" > /dev/null; then
        # Check RTdb vitality on Helsinki Hub
        if ! /data/data/com.termux/files/usr/bin/timeout 3 /data/data/com.termux/files/usr/bin/ssh -i /data/data/com.termux/files/home/.ssh/id_rsa -o ConnectTimeout=2 root@204.168.138.60 'echo 1' >/dev/null 2>&1; then
            # RTdb Offline Trigger Vibrate
            /data/data/com.termux/files/usr/bin/termux-vibrate -d 1500 -f >/dev/null 2>&1
            echo "[$(/data/data/com.termux/files/usr/bin/date)] GSH_MESH ALARM TRIGGERED RTDB OFFLINE" >> /data/data/com.termux/files/home/Sovereign_Node_Go/logs/forensic_timeline.log
        fi
    fi
done
