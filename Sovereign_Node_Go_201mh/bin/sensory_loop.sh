#!/data/data/com.termux/files/usr/bin/bash
export PATH=/data/data/com.termux/files/usr/bin:$PATH

# INITIALIZE HUD
termux-gui-view "/data/data/com.termux/files/home/Sovereign_Node_Go/bin/alpaca_mudd_multi_v2.json" &
sleep 2

while true; do
    # GRAB TELEMETRY
    BATT_PCT=$(termux-battery-status | jq '.percentage')
    BATT_PCT=${BATT_PCT:-100}
    WIFI_SSID=$(termux-wifi-connectioninfo | jq -r '.ssid')

    # UPDATE GUI
    termux-gui-view --update-view "sys_battery" --text "SYS_BATTERY: ${BATT_PCT}%"
    termux-gui-view --update-view "sys_wlan" --text "SYS_WLAN: ${WIFI_SSID}"

    # HAPTIC LOGIC: VIBRATE ON LOW BATTERY (< 15%) OR NO WIFI
    if [ "$BATT_PCT" -lt 15 ]; then
        termux-vibrate -d 1000 -f
        echo "[HAPTIC] FATAL: LOW BATTERY" >> /data/data/com.termux/files/home/Sovereign_Node_Go/bin/sensor.log
    fi

    if [[ "$WIFI_SSID" == "<unknown ssid>" || -z "$WIFI_SSID" ]]; then
        termux-vibrate -d 300 -f
        echo "[HAPTIC] WARN: WLAN DISCONNECTED" >> /data/data/com.termux/files/home/Sovereign_Node_Go/bin/sensor.log
    fi

    # SUBTLE PULSE FOR TICK CONFIRMATION
    termux-vibrate -d 15
    echo "[HEARTBEAT] Telemetry synced | BATT: ${BATT_PCT}% | WLAN: ${WIFI_SSID}" >> /data/data/com.termux/files/home/Sovereign_Node_Go/bin/sensor.log

    sleep 10
done
