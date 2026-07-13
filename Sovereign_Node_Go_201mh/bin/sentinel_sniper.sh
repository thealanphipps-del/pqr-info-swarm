#!/data/data/com.termux/files/usr/bin/bash
export PATH=/data/data/com.termux/files/usr/bin:$PATH
export REPAIR_SCRIPT="/data/data/com.termux/files/home/Sovereign_Node_Go/bin/agentic_repair.sh"
export SAFE_TMP="/data/data/com.termux/files/usr/tmp"

# TAIL MULTIPLE LOGS IN REAL-TIME, FILTERING FOR CRITICAL FLAGS
tail -n 0 -F /data/data/com.termux/files/home/Sovereign_Node_Go/bin/*.log /data/data/com.termux/files/home/Sovereign_Node_Go/gemini_testify/*.log 2>/dev/null | \
while read -r line; do
    if echo "$line" | grep -qiE "fatal|error|panic|denied|exception"; then
        
        # HASH THE ERROR TO PREVENT FORK BOMBS ON REPEATED IDENTICAL LOG LINES
        HASH=$(echo "$line" | md5sum | awk '{print $1}')
        LOCK_FILE="$SAFE_TMP/sniper_lock_$HASH"
        
        if [ ! -f "$LOCK_FILE" ]; then
            touch "$LOCK_FILE"
            
            # HAPTIC & GUI ALERT
            termux-vibrate -d 500 -f
            
            # ISOLATE LOGIC AND SPAWN REPAIR THREAD
            nohup "$REPAIR_SCRIPT" "$line" "mesh_log_stream" > "$SAFE_TMP/repair_output_$HASH.log" 2>&1 &
            
            # LOCK CLEARANCE DELAY (NO RM COMMAND ALLOWED)
            (sleep 120 && mv "$LOCK_FILE" "$LOCK_FILE.expired") &
        fi
    fi
done
