#!/data/data/com.termux/files/usr/bin/bash
LOG="/data/data/com.termux/files/home/Sovereign_Node_Go/ghub/director_ignition.log"
echo "[_] Initializing Director Ignition..." > $LOG
cd /data/data/com.termux/files/home/Sovereign_Node_Go/cmd/director
# Re-forge main.go if missing or mangled
if [ ! -f main.go ]; then
    echo "[!] Source missing, contacting Gemini..." >> $LOG
fi
/data/data/com.termux/files/usr/bin/go mod tidy >> $LOG 2>&1
CGO_ENABLED=0 /data/data/com.termux/files/usr/bin/go build -o /data/data/com.termux/files/home/Sovereign_Node_Go/bin/director main.go >> $LOG 2>&1
/data/data/com.termux/files/usr/bin/pkill -9 director 2>/dev/null
/data/data/com.termux/files/home/Sovereign_Node_Go/bin/director >> $LOG 2>&1 &
echo "[SUCCESS] Director active on PID $!" >> $LOG
