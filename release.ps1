# PQR Continuous Development Release Script (WSL Optimized)
# Usage: ./release.ps1 "Commit message"

param (
    [string]$Message = "RC-1 Deployment Update"
)

# 1. Detect Version from server.go
$VersionLine = Get-Content "server.go" | Select-String 'const Version = "(.*)"'
$Version = $VersionLine.Matches.Groups[1].Value

Write-Host "Preparing Release $Version"

# 2. Git Workflow via WSL
Write-Host "Syncing with GitHub (via WSL)..."
wsl git add .
$FullMessage = "[$Version] $Message"
wsl git commit -m "$FullMessage"
wsl git push origin main

Write-Host "Release $Version Pushed Successfully via WSL"
