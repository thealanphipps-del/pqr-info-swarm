<#
.SYNOPSIS
Builds the Windows Installer (MSI) for the OS-SPARK Sovereign Kernel.
#>

$ErrorActionPreference = "Stop"

Write-Host "=> Initializing OS-SPARK Windows Build Pipeline" -ForegroundColor Cyan

$BuildDir = ".\build-msi-temp"
$OutDir = ".\dist"

if (!(Test-Path $OutDir)) { New-Item -ItemType Directory -Path $OutDir | Out-Null }
if (Test-Path $BuildDir) { Remove-Item -Recurse -Force $BuildDir }
New-Item -ItemType Directory -Path $BuildDir | Out-Null

Write-Host "=> Compiling Go Binary for Windows (GOOS=windows GOARCH=amd64)"
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o "$BuildDir\sovereign-kernel.exe" ..\..\main.go

Write-Host "=> Generating WiX Toolset XML Configurations"
$WxsContent = @"
<?xml version='1.0' encoding='windows-1252'?>
<Wix xmlns='http://schemas.microsoft.com/wix/2006/wi'>
  <Product Name='OS-SPARK Sovereign Kernel' Manufacturer='Sovereign Mesh'
           Id='*' UpgradeCode='12345678-1234-1234-1234-123456789012'
           Language='1033' Codepage='1252' Version='1.0.0'>
    <Package Id='*' Keywords='Installer' Description='OS-SPARK Routing Daemon'
             Comments='Sovereign Mesh' Manufacturer='Sovereign Mesh'
             InstallerVersion='100' Languages='1033' Compressed='yes' SummaryCodepage='1252' />
    <Media Id='1' Cabinet='Sample.cab' EmbedCab='yes' DiskPrompt='CD-ROM #1' />
    <Property Id='DiskPrompt' Value='OS-SPARK Installation [1]' />
    
    <Directory Id='TARGETDIR' Name='SourceDir'>
      <Directory Id='ProgramFilesFolder' Name='PFiles'>
        <Directory Id='SovereignMesh' Name='SovereignMesh'>
          <Directory Id='INSTALLDIR' Name='OS-SPARK'>
            <Component Id='MainExecutable' Guid='87654321-4321-4321-4321-210987654321'>
              <File Id='SovereignKernelEXE' Name='sovereign-kernel.exe' DiskId='1' Source='$BuildDir\sovereign-kernel.exe' KeyPath='yes'/>
              <ServiceInstall Id='ServiceInstaller' Type='ownProcess' Name='SovereignKernel' DisplayName='OS-SPARK Continuity Kernel' Description='Multi-dimensional routing daemon' Start='auto' Account='LocalSystem' ErrorControl='normal'/>
              <ServiceControl Id='StartService' Stop='both' Remove='uninstall' Name='SovereignKernel' Wait='yes' />
            </Component>
          </Directory>
        </Directory>
      </Directory>
    </Directory>
    
    <Feature Id='Complete' Level='1'>
      <ComponentRef Id='MainExecutable' />
    </Feature>
  </Product>
</Wix>
"@

Set-Content -Path "$BuildDir\sovereign-kernel.wxs" -Value $WxsContent

Write-Host "=> To compile the MSI, install WiX Toolset and run:" -ForegroundColor Yellow
Write-Host "   candle.exe $BuildDir\sovereign-kernel.wxs"
Write-Host "   light.exe sovereign-kernel.wixobj -out $OutDir\sovereign-kernel.msi"

Write-Host "=> OS-SPARK Windows Build Pipeline Staged." -ForegroundColor Green
