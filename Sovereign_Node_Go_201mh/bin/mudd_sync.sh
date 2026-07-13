#!/data/data/com.termux/files/usr/bin/bash
export TERMUX_BASE="/data/data/com.termux/files/home/Sovereign_Node_Go"
export BIN="/data/data/com.termux/files/usr/bin"
export HUB="root@204.168.138.60"
export SSH_KEY="$HOME/.ssh/id_rsa"

# Fetch list of recently modified .bak files from Helsinki
$BIN/ssh -i $SSH_KEY $HUB "find /root/Sovereign_Unified -name '*.bak' -mmin -60" > $TERMUX_BASE/logs/mudd_pull.list

while IFS= read -r remote_file; do
    if [ -n "$remote_file" ]; then
        clean_path=$($BIN/echo "$remote_file" | $BIN/sed -E 's/\.[0-9]+\.bak$//')
        local_path="${TERMUX_BASE}${clean_path#/root/Sovereign_Unified}"
        
        $BIN/mkdir -p "$($BIN/dirname "$local_path")"
        
        if $BIN/scp -i $SSH_KEY "$HUB:$remote_file" "$local_path"; then
            $BIN/echo "0 [SUCCESS] MUDD_APPLIED $local_path" >> $TERMUX_BASE/logs/mudd_sync.log
        fi
    fi
done < $TERMUX_BASE/logs/mudd_pull.list
