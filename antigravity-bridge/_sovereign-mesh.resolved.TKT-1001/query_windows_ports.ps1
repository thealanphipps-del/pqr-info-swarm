Get-NetTCPConnection -State Listen | Select-Object LocalAddress, LocalPort | Format-Table -AutoSize
