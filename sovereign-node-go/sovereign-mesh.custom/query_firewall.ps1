Get-NetFirewallRule | Where-Object { $_.DisplayName -like "*WSL*" } | Select-Object DisplayName, Enabled, Direction, Action | Format-Table -AutoSize
