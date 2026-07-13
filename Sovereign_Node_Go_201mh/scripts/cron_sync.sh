#!/data/data/com.termux/files/usr/bin/bash
TIMESTAMP=$(/data/data/com.termux/files/usr/bin/date +"%Y-%m-%d %H:%M:%S")
cd /data/data/com.termux/files/home/Sovereign_Node_Go/ || exit 1

/data/data/com.termux/files/usr/bin/git fetch origin master
LOCAL=$(/data/data/com.termux/files/usr/bin/git rev-parse HEAD)
REMOTE=$(/data/data/com.termux/files/usr/bin/git rev-parse origin/master)

if [ "$LOCAL" != "$REMOTE" ]; then
    echo "[$TIMESTAMP] INFO: State mismatch detected. Syncing from Git DB..." >> /data/data/com.termux/files/home/Sovereign_Node_Go/logs/cron_sync.log
    /data/data/com.termux/files/usr/bin/git reset --hard origin/master
    
    export CGO_ENABLED=0
    if /data/data/com.termux/files/usr/bin/go build -o build/ignition ./cmd/main.go; then
        echo "[$TIMESTAMP] INFO: Build SUCCESS. Rotating binary." >> /data/data/com.termux/files/home/Sovereign_Node_Go/logs/cron_sync.log
        /data/data/com.termux/files/usr/bin/fuser -k 8080/tcp
        nohup ./build/ignition > ./logs/runtime.log 2>&1 &
    else
        echo "[$TIMESTAMP] ERROR: Build FAILED on sync." >> /data/data/com.termux/files/home/Sovereign_Node_Go/logs/cron_sync.log
    fi
fi
