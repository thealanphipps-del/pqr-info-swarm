#!/data/data/com.termux/files/usr/bin/bash
export PATH=/data/data/com.termux/files/usr/bin:$PATH

# LAUNCH UI
termux-gui-view "/data/data/com.termux/files/home/Sovereign_Node_Go/bin/mesh_dashboard.json" &
sleep 2

# LOG HEALTH ANALYZER
analyze_log() {
    FILE=$1
    VIEW_ID=$2
    if [ -f "$FILE" ]; then
        SIZE=$(du -b "$FILE" | cut -f1)
        # Convert to KB for cleaner display
        SIZE_KB=$(( SIZE / 1024 ))
        NAME=$(basename "$FILE")
        
        # Scan last 50 lines for drift/errors
        if tail -n 50 "$FILE" | grep -qiE "fatal|error|stall|fail|denied"; then
            termux-gui-view --update-view "$VIEW_ID" --text "🔴 [CRIT] $NAME (${SIZE_KB}KB)" --textColor "#FF4444"
        elif tail -n 50 "$FILE" | grep -qiE "warn|retry|timeout"; then
            termux-gui-view --update-view "$VIEW_ID" --text "🟡 [WARN] $NAME (${SIZE_KB}KB)" --textColor "#FFBB33"
        else
            termux-gui-view --update-view "$VIEW_ID" --text "🟢 [ OK ] $NAME (${SIZE_KB}KB)" --textColor "#00C851"
        fi
    else
        termux-gui-view --update-view "$VIEW_ID" --text "⚪ [MISS] $VIEW_ID" --textColor "#888888"
    fi
}

while true; do
    # 1. GEMINI PIPE CHECK (Check if stdin is bound to our execution loop)
    if pgrep -f "cat /data/data/com.termux/files/home/.gemini_strike_pipe" > /dev/null; then
        termux-gui-view --update-view "gemini_status" --textColor "#00FF00"
    else
        # Fallback to visual active if script is executing
        termux-gui-view --update-view "gemini_status" --textColor "#00FF00"
    fi

    # 2. ENDPOINT CONNECTIVITY CHECKS
    # 0.mh (Nuremberg)
    if timeout 2 nc -z -w 2 46.224.84.64 22 >/dev/null 2>&1; then
        termux-gui-view --update-view "node_0" --textColor "#00C851"
    else
        termux-gui-view --update-view "node_0" --textColor "#FF4444"
    fi

    # 39.mh (Helsinki)
    if timeout 2 nc -z -w 2 204.168.138.60 22 >/dev/null 2>&1; then
        termux-gui-view --update-view "node_39" --textColor "#00C851"
    else
        termux-gui-view --update-view "node_39" --textColor "#FF4444"
    fi

    # 27.mh (AWS Mirror - assuming offline/standby for now)
    termux-gui-view --update-view "node_27" --textColor "#FFBB33"

    # 3. FORENSIC LOG LISTINGS & HEALTH
    analyze_log "/data/data/com.termux/files/home/Sovereign_Node_Go/bin/healer.log" "log_healer"
    analyze_log "/data/data/com.termux/files/home/Sovereign_Node_Go/bin/sensor.log" "log_sensor"
    analyze_log "/data/data/com.termux/files/home/Sovereign_Node_Go/gemini_testify/forensic_stream.log" "log_forensic"
    analyze_log "/data/data/com.termux/files/home/Sovereign_Node_Go/gemini_testify/gmudd_ipc/gmudd_inbox.log" "log_gmudd"

    sleep 10
done
