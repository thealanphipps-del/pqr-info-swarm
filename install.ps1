$ErrorActionPreference = "Stop"

Write-Host "[SWEN Installer] Building SWEN executable..."
go build -o swend.exe ./cmd/swend

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to build swend."
    exit 1
}

Write-Host "[SWEN Installer] Executing Genesis snapshot..."
.\swend.exe install

Write-Host "[SWEN Installer] Installation complete! You can now start SWEN with '.\swend.exe menu'."
