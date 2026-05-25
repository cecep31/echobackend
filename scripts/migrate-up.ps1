$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

if (-not (Test-Path ".env")) {
    Write-Error ".env not found. Copy .env.example to .env and configure DATABASE_URL."
}

Get-Content ".env" | ForEach-Object {
    if ($_ -match '^\s*([^#=]+)=(.*)$') {
        $name = $matches[1].Trim()
        $value = $matches[2].Trim().Trim('"')
        if ($value -match '^\$\{(.+)\}$' -and (Test-Path "env:$($matches[1])")) {
            $value = (Get-Item "env:$($matches[1])").Value
        }
        Set-Item -Path "env:$name" -Value $value
    }
}

if (-not $env:DATABASE_URL) {
    Write-Error "DATABASE_URL is not set in .env"
}

if ($env:GOOSE_DBSTRING -match '^\$\{DATABASE_URL\}$') {
    $env:GOOSE_DBSTRING = $env:DATABASE_URL
}

Write-Host "Creating goose schema (custom) if missing..."
psql $env:DATABASE_URL -v ON_ERROR_STOP=1 -f "scripts/bootstrap-goose-schema.sql"

Write-Host "Running goose up..."
goose up
