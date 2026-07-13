#!/data/data/com.termux/files/usr/bin/bash

# --- CORE PATHS ---
BIN="/data/data/com.termux/files/usr/bin"
BASE_DIR="/data/data/com.termux/files/home"
ARCHIVE_DIR="$BASE_DIR/Sovereign_Node_Go/offworld_staging"
RCLONE_CONF="$BASE_DIR/.config/rclone/rclone.conf"

# --- SYSTEM HEALTH MODULE ---
sys_health() {
    BATT_JSON=$($BIN/termux-battery-status)
    BATT_PCT=$($BIN/echo "$BATT_JSON" | $BIN/jq '.percentage')
    BATT_TEMP=$($BIN/echo "$BATT_JSON" | $BIN/jq '.temperature')
    LOAD=$($BIN/uptime | $BIN/awk -F'load average:' '{print $2}')
    MEM_FREE=$($BIN/free -m | $BIN/awk '/^Mem:/{print $4 "MB"}')
    
    $BIN/termux-dialog confirm -t "System Vitality" -i "Battery: ${BATT_PCT}% \nTemp: ${BATT_TEMP}C \nLoad:${LOAD} \nMem Free: ${MEM_FREE}"
}

# --- RCLONE GOOGLE PHOTOS SYNC ---
configure_photos() {
    $BIN/termux-dialog confirm -t "Google Photos Bridge" -i "Initiating rclone config for alan@w-isp.net. A browser will open for OAuth."
    $BIN/rclone config create gphotos_alan google photos --non-interactive
    $BIN/termux-dialog confirm -t "Rclone" -i "Run 'rclone config reconnect gphotos_alan:' to authenticate via browser."
}

# --- TAKEOUT INTERACTIVE BROWSER ---
takeout_browser() {
    # Find zip archives and extracted folders
    PAYLOADS=$($BIN/find $BASE_DIR -maxdepth 4 -iname "*takeout*" -type d,f | $BIN/tr '\n' ',')
    SELECTION=$($BIN/termux-dialog radio -t "Select Takeout Payload" -v "$PAYLOADS" | $BIN/jq -r '.text')
    
    if [ "$SELECTION" != "null" ] && [ -n "$SELECTION" ]; then
        if [[ "$SELECTION" == *.zip ]]; then
            $BIN/termux-dialog confirm -t "Archive Detected" -i "Extracting $SELECTION to staging..."
            DEST="$ARCHIVE_DIR/extract_$(basename "$SELECTION" .zip)"
            $BIN/mkdir -p "$DEST"
            $BIN/unzip -q "$SELECTION" -d "$DEST"
            SELECTION="$DEST"
        fi
        
        # Navigate extracted directory
        while true; do
            FILES=$($BIN/ls -1 "$SELECTION" | $BIN/tr '\n' ',')
            NAV=$($BIN/termux-dialog radio -t "Browsing: $SELECTION" -v "..,${FILES}" | $BIN/jq -r '.text')
            
            if [ "$NAV" == "null" ] || [ -z "$NAV" ]; then break; fi
            
            if [ "$NAV" == ".." ]; then
                SELECTION=$($BIN/dirname "$SELECTION")
            elif [ -d "$SELECTION/$NAV" ]; then
                SELECTION="$SELECTION/$NAV"
            elif [ -f "$SELECTION/$NAV" ]; then
                # Dispatch to Android Intent (PDF viewing/signing, Audio, Video)
                $BIN/termux-open "$SELECTION/$NAV"
            fi
        done
    fi
}

# --- FORENSIC ORGANIZER & DUPLICATE SWEEP ---
organize_sweep() {
    TARGET_DIR=$($BIN/termux-dialog text -t "Enter Target Directory for Organization" -i "$BASE_DIR/downloads" | $BIN/jq -r '.text')
    if [ -d "$TARGET_DIR" ]; then
        $BIN/termux-toast -s "Sweeping Duplicates..."
        $BIN/fdupes -rdN "$TARGET_DIR" > /dev/null 2>&1
        
        $BIN/termux-toast -s "Organizing Documents..."
        $BIN/mkdir -p "$TARGET_DIR/PDFs" "$TARGET_DIR/Media" "$TARGET_DIR/Archives"
        $BIN/find "$TARGET_DIR" -maxdepth 1 -iname "*.pdf" -exec $BIN/mv {} "$TARGET_DIR/PDFs/" \;
        $BIN/find "$TARGET_DIR" -maxdepth 1 -iname "*.mp3" -o -iname "*.wav" -o -iname "*.jpg" -o -iname "*.png" -exec $BIN/mv {} "$TARGET_DIR/Media/" \;
        $BIN/find "$TARGET_DIR" -maxdepth 1 -iname "*.zip" -o -iname "*.tar.gz" -exec $BIN/mv {} "$TARGET_DIR/Archives/" \;
        
        $BIN/termux-dialog confirm -t "Organization Complete" -i "Duplicates purged. Files sorted into PDFs, Media, and Archives."
    fi
}

# --- MAIN LOOP ---
while true; do
    CHOICE=$($BIN/termux-dialog radio -t "Jovian Overseer (39.mh Bridge)" -v "1. System Vitality,2. Google Photos Config (alan@w-isp.net),3. Interactive Takeout Browser,4. Duplicate Sweep & Organizer,5. Exit" | $BIN/jq -r '.text')
    
    case "$CHOICE" in
        "1. System Vitality") sys_health ;;
        "2. Google Photos Config (alan@w-isp.net)") configure_photos ;;
        "3. Interactive Takeout Browser") takeout_browser ;;
        "4. Duplicate Sweep & Organizer") organize_sweep ;;
        "5. Exit"|"null"|"") break ;;
    esac
done
