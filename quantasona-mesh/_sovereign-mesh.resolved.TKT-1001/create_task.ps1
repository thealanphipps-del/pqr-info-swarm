schtasks.exe /create /tn "SetWSLFirewall" /tr "powershell.exe -ExecutionPolicy Bypass -File C:\temp\set_firewall.ps1" /sc once /st 23:59 /sd 12/31/2026 /ru "SYSTEM"
schtasks.exe /run /tn "SetWSLFirewall"
