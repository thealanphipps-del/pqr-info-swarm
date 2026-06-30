Start-Process powershell -Verb RunAs -ArgumentList "-ExecutionPolicy Bypass -File C:\temp\set_firewall.ps1" -Wait -WindowStyle Hidden
