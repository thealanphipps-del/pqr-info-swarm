#!/data/data/com.termux/files/usr/bin/bash
# GSH Dashboard - Visual Heartbeat
FLOOR="982M"
STATUS=$(tail -n 1 $MESH_LOG)
termux-notification \
    --id GSH-DASH \
    --title "🦅 GSH-MESH | FLOOR: $FLOOR" \
    --content "LOG: $STATUS" \
    --priority high \
    --led-color 00FF00 \
    --ongoing
