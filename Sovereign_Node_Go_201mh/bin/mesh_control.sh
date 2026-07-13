#!/data/data/com.termux/files/usr/bin/bash
export PATH=/data/data/com.termux/files/usr/bin:$PATH

UI_DATA=$(termux-gui-dialog --title "GSH: MESH CONTROL CENTER" \
  label "System State Configuration:" \
  switch "Autonomous Healer" true \
  switch "Sentinel Sniper" true \
  checkbox "Enable Trace Forensics" true \
  spinner 3 "Anchor: MA-25" "Anchor: Fibonacci" "Anchor: RSI" \
  radio 3 "Policy: Passive" "Policy: Active" "Policy: Aggressive" 2 \
  label "Helsinki Auth Token (39.mh):" \
  text textPassword)

if [ $? -eq 2 ]; then
  termux-toast -c white -b red "Config Update Aborted"
  exit 1
fi

echo "[$(date +%s)] [CONTROL_UI] New Params: $UI_DATA" >> /data/data/com.termux/files/home/Sovereign_Node_Go/gemini_testify/gmudd_ipc/gmudd_inbox.log

termux-vibrate -d 100
termux-toast -c white -b blue "[GEMINI] FIX_APPLIED: MESH_PARAMETERS_LOCKED"
exit 0
