# Build OTA bundle.zip on Windows (PowerShell)
$ErrorActionPreference = "Stop"
$Root = Split-Path (Split-Path $PSScriptRoot -Parent) -Parent
$BundleDir = Join-Path $Root "ota\bundle"
$OutDir = Join-Path $Root "ota\public\v1"
$ZipPath = Join-Path $OutDir "bundle.zip"
$ManifestPath = Join-Path $Root "ota\public\manifest.json"

New-Item -ItemType Directory -Force -Path $OutDir | Out-Null
if (Test-Path $ZipPath) { Remove-Item $ZipPath }

# Zip contents of bundle/ (pages, css, js folders)
Compress-Archive -Path (Join-Path $BundleDir "*") -DestinationPath $ZipPath -Force

$hash = (Get-FileHash -Path $ZipPath -Algorithm SHA256).Hash.ToLower()
Write-Host "SHA256: $hash"

$manifest = Get-Content $ManifestPath -Raw | ConvertFrom-Json
$manifest.sha256 = $hash
$manifest | ConvertTo-Json -Depth 5 | Set-Content $ManifestPath -Encoding UTF8

Write-Host "Updated manifest.json with sha256"
Write-Host "Upload folder ota\public\ to https://run.8me.in/ota/"
