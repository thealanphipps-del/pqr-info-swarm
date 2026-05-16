# PQR Swarm Starter Script
$projectDir = "C:\Users\drphi\pqr-info-swarm"

if (Test-Path $projectDir) {
    Set-Location $projectDir
    Write-Host "Initializing PQR Info Swarm Stack..." -ForegroundColor Cyan
    
    # Check if Docker is running using a command probe
    $dockerRunning = docker info 2>$null
    if (!$dockerRunning) {
        Write-Host "Docker Engine is not responding. Please ensure Docker Desktop is started." -ForegroundColor Red
        return
    }

    # Start the Docker Stack (Redundant Cluster Mode)
    Write-Host "Lifting the Sovereign Mesh (3-Node Redundant Cluster)..." -ForegroundColor Cyan
    docker-compose up -d --build
    if ($LASTEXITCODE -eq 0) {
        Write-Host "PQR Sovereign Cluster is ONLINE at pqr.info" -ForegroundColor Green
        Write-Host "Load Balancer (Nginx): http://localhost:3196" -ForegroundColor Green
        Write-Host "Monitoring Fabric: http://localhost:8081" -ForegroundColor Yellow
        Write-Host "Identity Vault: http://localhost:8200" -ForegroundColor Yellow
    } else {
        Write-Host "Failed to start PQR stack. Check docker-compose logs." -ForegroundColor Red
    }
} else {
    Write-Host "Error: Project directory not found." -ForegroundColor Red
}
