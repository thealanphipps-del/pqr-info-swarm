#!/bin/bash
PIPE="/data/data/com.termux/files/home/Sovereign_Node_Go/gemini_testify/rpc_command.pipe"
INBOX="/data/data/com.termux/files/home/Sovereign_Node_Go/gemini_testify/gmudd_ipc/gmudd_inbox.log"

while true; do
    if read line < "$PIPE"; then
        # Deduplication Logic: Check new content against the last recorded entry
        LAST_CONTENT=$(/data/data/com.termux/files/usr/bin/tail -n 1 "$INBOX" 2>/dev/null | /data/data/com.termux/files/usr/bin/sed 's/.*\[INGEST\] //')
        if [[ "$line" != "$LAST_CONTENT" ]]; then
            /data/data/com.termux/files/usr/bin/echo "[$(/data/data/com.termux/files/usr/bin/date +%s)] [INGEST] $line" >> "$INBOX"
        fi
    fi
done
