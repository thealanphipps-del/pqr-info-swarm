# SENTINEL.ps1 - PQR Sovereign Guardian
# Monitors the PQR Mesh from the Windows Host side.

$ErrorActionPreference = "Continue"
$projectDir = "C:\Users\drphi\pqr-info-swarm"
$healthUrl = "http://localhost:3196/REST/2.0/health"
$logFile = "$projectDir\sentinel.log"
$signalsDir = "$projectDir\signals"
$triggerFile = "$signalsDir\RESTART_TRIGGER"

if (!(Test-Path $signalsDir)) { New-Item -ItemType Directory -Path $signalsDir | Out-Null }

function Write-Log($msg, $color = "White") {
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $formattedMsg = "[$timestamp] $msg"
    Write-Host $formattedMsg -ForegroundColor $color
    try {
        $formattedMsg | Out-File -FilePath $logFile -Append -ErrorAction SilentlyContinue
    } catch {}
}

Clear-Host
Write-Log "--------------------------------------------------------" -Color Cyan
Write-Log "      PQR SENTINEL: Sovereign Guardian v1.0.0           " -Color Cyan
Write-Log "      Autonomous Host-Side Monitoring Active            " -Color Cyan
Write-Log "--------------------------------------------------------" -Color Cyan
Write-Log "Project Root: $projectDir" -Color Gray
Write-Log "Health Check: $healthUrl" -Color Gray
Write-Log "Sentinel Log: $logFile" -Color Gray
Write-Log "--------------------------------------------------------" -Color Cyan

Set-Location $projectDir

while ($true) {
    # 1. Check Docker Engine Status
    $dockerCheck = docker info 2>$null
    if ($LASTEXITCODE -ne 0) {
        Write-Log "CRITICAL: Docker Engine is OFFLINE or not responding." -Color Red
        Write-Log "Please ensure Docker Desktop is running on Windows." -Color Yellow
        Start-Sleep -Seconds 10
        continue
    }

    # 2. Check PQR Server Health
    $isHealthy = $false
    try {
        $response = Invoke-RestMethod -Uri $healthUrl -Method Get -TimeoutSec 5
        if ($response.status -eq "healthy") {
            # Sub-heartbeat: We could check which replica served this if we added a header
        } else {
            Write-Log "WARNING: PQR Mesh reported DEGRADED status: $($response.status)" -Color Yellow
        }
    } catch {
        Write-Log "ALERT: PQR Mesh is UNREACHABLE (HTTP Timeout/Error)." -Color Red
        Write-Log "Attempting recovery of Gateway and Server Cluster..." -Color Cyan
        
        # Re-lift the gateway and servers
        docker-compose up -d gateway pqr-server
        
        Start-Sleep -Seconds 20 # Allow cluster to synchronize
    }

    # 3. Check for Agent-Driven Signals
    if (Test-Path $triggerFile) {
        Write-Log "SIGNAL DETECTED: Agent-driven restart trigger (RESTART_TRIGGER)." -Color Magenta
        Write-Log "Executing full stack reconciliation (up -d --build)..." -Color Cyan
        
        docker-compose up -d --build
        
        if ($LASTEXITCODE -eq 0) {
            Write-Log "SUCCESS: Stack rebuilt and lifted." -Color Green
            Remove-Item $triggerFile -Force
        } else {
            Write-Log "ERROR: Stack rebuild failed. Check Docker logs." -Color Red
        }
    }

    # 4. Check for Orphaned Containers
    # (Future expansion: check for high CPU/Memory and alert)

    Start-Sleep -Seconds 15
}
