"Starting set_firewall.ps1" | Out-File C:\temp\log.txt
Set-NetFirewallRule -DisplayName "Ubuntu Pro for WSL" -RemoteAddress @("192.168.12.0/24", "LocalSubnet") 2>&1 | Out-File -Append C:\temp\log.txt
"Finished set_firewall.ps1" | Out-File -Append C:\temp\log.txt
