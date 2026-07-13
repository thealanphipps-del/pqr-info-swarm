#!/data/data/com.termux/files/usr/bin/bash
export PATH=/data/data/com.termux/files/usr/bin:$PATH

SESSION_NAME="GSH_SPLIT"

# PURGE EXISTING ORPHANED SESSIONS TO PREVENT NESTED ATTACHMENT ERRORS
tmux kill-session -t $SESSION_NAME 2>/dev/null || true

# INITIALIZE DETACHED SESSION
tmux new-session -d -s $SESSION_NAME

# SPLIT WINDOW VERTICALLY (BOTTOM PANE AT 40% HEIGHT)
tmux split-window -v -p 40

# DISPATCH LIVE FORENSIC TAIL TO BOTTOM PANE
tmux send-keys -t $SESSION_NAME:0.1 "tail -n 10 -F /data/data/com.termux/files/home/Sovereign_Node_Go/bin/sensor.log /data/data/com.termux/files/home/Sovereign_Node_Go/bin/healer.log /data/data/com.termux/files/home/Sovereign_Node_Go/gemini_testify/forensic_stream.log" C-m

# LOCK FOCUS TO TOP PANE FOR MASTERER INPUT
tmux select-pane -t $SESSION_NAME:0.0

# CLEAR TERMINAL AND ATTACH
clear
tmux attach-session -t $SESSION_NAME
