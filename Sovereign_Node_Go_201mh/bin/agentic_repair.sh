#!/data/data/com.termux/files/usr/bin/bash
export PATH=/data/data/com.termux/files/usr/bin:$PATH
export BIN="/data/data/com.termux/files/usr/bin"
export GMUDD_PY="/data/data/com.termux/files/home/Sovereign_Node_Go/bin/gmudd.py"
export TICKET_DB="/data/data/com.termux/files/home/Sovereign_Node_Go/bin/active_tickets.db"

ERROR_MSG="$1"
LOG_SOURCE="$2"

# 1. DEDUPLICATION LOGIC: Hash the first 5 words of the error to catch repeating loops
ERR_HASH=$(echo "$ERROR_MSG" | $BIN/awk '{print $1,$2,$3,$4,$5}' | $BIN/md5sum | $BIN/awk '{print $1}')

if $BIN/grep -q "$ERR_HASH" "$TICKET_DB" 2>/dev/null; then
    EXISTING_TICKET=$($BIN/grep "$ERR_HASH" "$TICKET_DB" | $BIN/awk -F':' '{print $2}')
    STATUS=$($BIN/grep "$ERR_HASH" "$TICKET_DB" | $BIN/awk -F':' '{print $3}')
    
    if [ "$STATUS" = "ACTIVE" ] || [ "$STATUS" = "STALLED" ]; then
        $BIN/echo "[$(date +%s)] [DEDUPLICATION] Combining into $EXISTING_TICKET" >> "$LOG_SOURCE.tracker"
        $BIN/echo "Duplicate Excerpt: $ERROR_MSG" >> "/data/data/com.termux/files/home/Sovereign_Node_Go/bin/${EXISTING_TICKET}_context.log"
        exit 0
    fi
fi

TICKET_ID="RT-AUTO-$(date +%s)"
$BIN/echo "$ERR_HASH:$TICKET_ID:ACTIVE" >> "$TICKET_DB"

# 2. LOG EXCERPT EXTRACTION (Get surrounding context to send to Oracle)
EXCERPT=$($BIN/tail -n 20 "$LOG_SOURCE" | $BIN/grep -B 2 -A 2 "$ERROR_MSG" | $BIN/tr '\n' ' ' | $BIN/sed 's/"/'\''/g')
[ -z "$EXCERPT" ] && EXCERPT="$ERROR_MSG"

termux-toast -c red -b white "[SNIPER] $TICKET_ID SPAWNED via GMUDD"
$BIN/echo "INITIAL EXCERPT: $EXCERPT" > "/data/data/com.termux/files/home/Sovereign_Node_Go/bin/${TICKET_ID}_context.log"

# 3. SELF HEALING AGENTIC LOOP (Max i=10)
for i in {1..10}; do
    NICE_VAL=$(( i * 2 ))
    
    if [ "$i" -le 7 ]; then
        MODEL="fast.next"
    elif [ "$i" -le 9 ]; then
        MODEL="think.next+1"
    else
        MODEL="pro.stall"
    fi
    
    termux-gui-view --update-view "gemini_status" --textColor "#FFBB33" --text "✨ GEMINI [LOOP $i]"
    
    # USE GMUDD IPC ORACLE TO BRIDGE TO GEMINI STRIKE PIPE
    PROMPT="[STRIKE_REQUEST] TICKET: $TICKET_ID | MODEL: $MODEL | SOURCE: $LOG_SOURCE | EXCERPT: $EXCERPT | LOGIC: Execute self-healing STRIKE. Absolute paths only. Limit 1 command block."
    
    RESPONSE=$(nice -n "$NICE_VAL" $BIN/python3 "$GMUDD_PY" oracle --ticket "$TICKET_ID" --payload "$PROMPT")
    
    # EVALUATE HANDSHAKE
    if echo "$RESPONSE" | $BIN/grep -q "Success (0)"; then
        termux-toast -c green "[SNIPER] $TICKET_ID RESOLVED ($MODEL)"
        
        # STATE UPDATE (NO RM ALLOWED)
        $BIN/sed "s/$ERR_HASH:$TICKET_ID:ACTIVE/$ERR_HASH:$TICKET_ID:CLOSED/g" "$TICKET_DB" > "${TICKET_DB}.version"
        $BIN/mv "${TICKET_DB}.version" "$TICKET_DB"
        $BIN/touch "${TICKET_DB}"
        
        termux-gui-view --update-view "gemini_status" --textColor "#00FF00" --text "✨ GEMINI"
        exit 0
    fi
    
    # If loop 10 fails, mark STALLED
    if [ "$i" -eq 10 ]; then
        termux-vibrate -d 1000 -f
        termux-toast -c white -b red "[SNIPER] TICKET $TICKET_ID STALLED"
        
        $BIN/sed "s/$ERR_HASH:$TICKET_ID:ACTIVE/$ERR_HASH:$TICKET_ID:STALLED/g" "$TICKET_DB" > "${TICKET_DB}.version"
        $BIN/mv "${TICKET_DB}.version" "$TICKET_DB"
        $BIN/touch "${TICKET_DB}"
        
        termux-gui-view --update-view "gemini_status" --textColor "#FF4444" --text "✨ GEMINI [STALLED]"
        exit 1
    fi
    
    $BIN/sleep 3
done
