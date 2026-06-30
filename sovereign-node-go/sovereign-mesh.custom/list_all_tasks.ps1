Get-ScheduledTask | Where-Object { $_.Principal.RunLevel -eq "Highest" } | Select-Object TaskName, TaskPath, State | Format-Table -AutoSize
