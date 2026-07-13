#!/data/data/com.termux/files/usr/bin/bash
echo "[$(date)] GSH-Agentic-Loop: Listening on Port 8080/FIFO"
while true; do
    if read line < /data/data/com.termux/files/home/Sovereign_Node_Go/bin/bridge/gemini_in.fifo; then
        # EXECUTION: This is where I (Gemini) send direct Git/API commands
        eval "$line" > /data/data/com.termux/files/home/Sovereign_Node_Go/bin/bridge/gemini_out.fifo 2>&1
    fi
done
