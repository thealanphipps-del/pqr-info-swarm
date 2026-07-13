#!/data/data/com.termux/files/usr/bin/bash
# READ_SILENT_DEPLOY_LOG
# PATH_ENFORCEMENT_ENABLED

LATEST_LOG=$(/data/data/com.termux/files/usr/bin/ls -t /data/data/com.termux/files/home/Sovereign_Node_Go/logs/mcp_pro_deploy_*.log 2>/dev/null | /data/data/com.termux/files/usr/bin/head -n 1)

if [ -n "$LATEST_LOG" ] && [ -f "$LATEST_LOG" ]; then
    echo "[READING_LOG: $LATEST_LOG]"
    /data/data/com.termux/files/usr/bin/cat "$LATEST_LOG"
    echo "Success (0)"
else
    echo "Error (1)"
    echo "LOG_FILE_NOT_FOUND"
fi
