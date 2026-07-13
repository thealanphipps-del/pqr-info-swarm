#!/data/data/com.termux/files/usr/bin/bash
export PATH=/data/data/com.termux/files/usr/bin:$PATH
export BIN="/data/data/com.termux/files/usr/bin"
export HUD_JSON="/data/data/com.termux/files/home/Sovereign_Node_Go/bin/mesh_dashboard_v2.json"

# INITIALIZE HUD
termux-gui-view "$HUD_JSON" &
sleep 2

while true; do
    # 1. SCAN MESH CONNECTIVITY
    # 0.mh Nuremberg
    if timeout 1 nc -z 46.224.84.64 22 2>/dev/null; then
        termux-gui-view --update-view "node_0" --textColor "#00C851"
    else
        termux-gui-view --update-view "node_0" --textColor "#FF4444"
    fi
    
    # 39.mh Helsinki
    if timeout 1 nc -z 204.168.138.60 22 2>/dev/null; then
        termux-gui-view --update-view "node_39" --textColor "#00C851"
    else
        termux-gui-view --update-view "node_39" --textColor "#FF4444"
    fi

    # 2. SCAN LOG HEALTH AND RE-RENDER LIST
    # Clear and rebuild list logic via JSON update
    LOGS=$(ls /data/data/com.termux/files/home/Sovereign_Node_Go/bin/*.log /data/data/com.termux/files/home/Sovereign_Node_Go/gemini_testify/*.log 2>/dev/null)
    
    # Simple list update via toast for now while termux-gui-view dynamic insertion matures
    for log in $LOGS; do
        NAME=$(basename "$log")
        if tail -n 50 "$log" | grep -qiE "fatal|error|panic"; then
            COLOR="#FF4444" # CRITICAL
            termux-vibrate -d 30
        elif tail -n 50 "$log" | grep -qiE "warn|retry"; then
            COLOR="#FFBB33" # WARNING
        else
            COLOR="#00C851" # NOMINAL
        fi
        # Updating placeholder view IDs if they exist in layout
        # Note: In production we use a ListView, here we simulate status updates
        termux-gui-view --update-view "log_list" --text "SCANNING: $NAME" --textColor "$COLOR"
    done

    # 3. GEMINI HEARTBEAT
    if pgrep -f "cat /data/data/com.termux/files/home/.gemini_strike_pipe" > /dev/null; then
        termux-gui-view --update-view "gemini_icon" --textColor "#00FF00"
    else
        termux-gui-view --update-view "gemini_icon" --textColor "#888888"
    fi

    sleep 15
done
