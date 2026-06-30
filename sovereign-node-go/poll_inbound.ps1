while($true) {
    try {
        $response = Invoke-WebRequest -Uri "http://192.168.12.201:8081/api/bridge?cmd=ls%20-1%20sovereign-node-go/inbound" -UseBasicParsing -ErrorAction Stop
        if ($response.Content -and $response.Content -notlike "*ERROR*") {
            $files = $response.Content -split "\r?\n"
            foreach ($file in $files) {
                $file = $file.Trim()
                if ($file -and $file -ne "from_antigravity.txt" -and $file -notlike "*ERROR*") {
                     $msg = Invoke-WebRequest -Uri "http://192.168.12.201:8081/api/bridge?cmd=cat%20sovereign-node-go/inbound/$file" -UseBasicParsing
                     if ($msg.Content -and $msg.Content -notlike "*ERROR*") {
                         Write-Output "[MESSAGE FROM AGENT] $file`: $($msg.Content)"
                         # Archive the message
                         Invoke-WebRequest -Uri "http://192.168.12.201:8081/api/bridge?cmd=mv%20sovereign-node-go/inbound/$file%20sovereign-node-go/logs/archived_$file" -UseBasicParsing
                     }
                }
            }
        }
    } catch { }
    Start-Sleep -Seconds 10
}
