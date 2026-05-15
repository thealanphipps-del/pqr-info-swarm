# PQR Swarm Starter Script
$projectDir = "C:\Users\drphi\pqr-info-swarm"

if (Test-Path $projectDir) {
    Set-Location $projectDir
    Write-Host "Initializing PQR Info Swarm Stack..." -ForegroundColor Cyan
    
    # Check if Docker is running
    $dockerCheck = Get-Process docker -ErrorAction SilentlyContinue
    if ($null -eq $dockerCheck) {
        Write-Host "Docker Desktop is not running. Please start it first." -ForegroundColor Red
        return
    }

    docker-compose up -d
    if ($LASTEXITCODE -eq 0) {
        Write-Host "PQR Stack is Online" -ForegroundColor Green
    } else {
        Write-Host "Failed to start Docker stack." -ForegroundColor Red
    }
} else {
    Write-Host ("Error: Project directory not found at " + $projectDir) -ForegroundColor Red
}
