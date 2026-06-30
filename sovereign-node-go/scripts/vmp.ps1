param (
    [Parameter(Mandatory=$true)]
    [string]$Path
)

$archiveDir = "c:\Users\drphi\Jovian_Archives"
$logFile = "c:\Users\drphi\Sovereign_Node_Go\forensic_chain.log"

if (-not (Test-Path $Path)) {
    Write-Error "Path not found: $Path"
    exit 1
}

$timestamp = Get-Date -Format "yyyyMMdd_HHmmss_fff"
$fileName = Split-Path $Path -Leaf
$destination = Join-Path $archiveDir "$fileName.bak.$timestamp"

Write-Host "VMP: Moving $Path to $destination"
Move-Item -Path $Path -Destination $destination -Force

$logEntry = "[$(Get-Date -Format "yyyy-MM-dd HH:mm:ss.fff")] VMP_MOVE: Source=$Path, Destination=$destination"
Add-Content -Path $logFile -Value $logEntry

Write-Host "VMP: Forensic entry recorded."
