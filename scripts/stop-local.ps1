$ErrorActionPreference = "Continue"

$repoRoot = Split-Path -Parent $PSScriptRoot
$runtimeDir = Join-Path $repoRoot ".local"

foreach ($serviceName in @("frontend", "backend")) {
    $pidFile = Join-Path $runtimeDir "$serviceName.pid"
    if (-not (Test-Path -LiteralPath $pidFile)) {
        continue
    }

    $storedPid = Get-Content -LiteralPath $pidFile -ErrorAction SilentlyContinue
    if ($storedPid -and (Get-Process -Id ([int]$storedPid) -ErrorAction SilentlyContinue)) {
        & taskkill.exe /PID $storedPid /T /F | Out-Null
        Write-Host "$serviceName detenido."
    }

    Remove-Item -LiteralPath $pidFile -Force -ErrorAction SilentlyContinue
}

docker compose --project-directory $repoRoot -f (Join-Path $repoRoot "compose.yaml") stop postgres
Write-Host "PostgreSQL detenido. Los datos locales se conservaron."
