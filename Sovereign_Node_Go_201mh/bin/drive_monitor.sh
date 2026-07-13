#!/data/data/com.termux/files/usr/bin/bash
while true; do
    SPACE_AVAIL=$(/data/data/com.termux/files/usr/bin/df -h /data | /data/data/com.termux/files/usr/bin/awk 'NR==2 {print $4}')
    /data/data/com.termux/files/usr/bin/termux-toast -s -b black -c white "Drive Space: $SPACE_AVAIL"
    /data/data/com.termux/files/usr/bin/sleep 300
done
