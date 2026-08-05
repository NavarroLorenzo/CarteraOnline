$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$backendDir = Join-Path $repoRoot "backend"
$frontendDir = Join-Path $repoRoot "frontend"
$runtimeDir = Join-Path $repoRoot ".local"
$logsDir = Join-Path $runtimeDir "logs"

function Assert-Command([string]$command, [string]$help) {
    if (-not (Get-Command $command -ErrorAction SilentlyContinue)) {
        throw "$command no está disponible. $help"
    }
}

function Test-ProcessFromPidFile([string]$pidFile) {
    if (-not (Test-Path -LiteralPath $pidFile)) {
        return $false
    }

    $storedPid = Get-Content -LiteralPath $pidFile -ErrorAction SilentlyContinue
    if (-not $storedPid) {
        return $false
    }

    return $null -ne (Get-Process -Id ([int]$storedPid) -ErrorAction SilentlyContinue)
}

function Wait-Url([string]$url, [int]$attempts) {
    for ($attempt = 1; $attempt -le $attempts; $attempt++) {
        $previousPreference = $ErrorActionPreference
        $ErrorActionPreference = "Continue"
        & curl.exe --silent --output NUL --fail --max-time 2 $url
        $curlExitCode = $LASTEXITCODE
        $ErrorActionPreference = $previousPreference

        if ($curlExitCode -eq 0) {
            return $true
        }

        Start-Sleep -Seconds 1
    }

    return $false
}

function Start-DetachedCommand(
    [string]$workingDirectory,
    [string]$command,
    [string]$pidFile
) {
    $payload = "Set-Location -LiteralPath '$($workingDirectory.Replace("'", "''"))'; $command"
    $encodedPayload = [Convert]::ToBase64String([Text.Encoding]::Unicode.GetBytes($payload))
    $commandLine = "powershell.exe -NoProfile -WindowStyle Hidden -EncodedCommand $encodedPayload"
    $result = Invoke-CimMethod -ClassName Win32_Process -MethodName Create -Arguments @{
        CommandLine = $commandLine
        CurrentDirectory = $workingDirectory
    }

    if ($result.ReturnValue -ne 0) {
        throw "No se pudo iniciar el proceso local. Código: $($result.ReturnValue)"
    }

    Set-Content -LiteralPath $pidFile -Value $result.ProcessId
}

Assert-Command "docker" "Instalá o iniciá Docker Desktop."
Assert-Command "go" "Instalá la versión de Go indicada en backend/go.mod."
Assert-Command "node" "Instalá Node.js 20 o superior."
Assert-Command "npm.cmd" "La instalación de Node.js debe incluir npm."

$previousErrorActionPreference = $ErrorActionPreference
$ErrorActionPreference = "Continue"
docker info *> $null
$dockerInfoExitCode = $LASTEXITCODE
$ErrorActionPreference = $previousErrorActionPreference
if ($dockerInfoExitCode -ne 0) {
    throw "Docker Desktop está instalado pero no está iniciado. Abrilo y volvé a ejecutar este comando."
}

New-Item -ItemType Directory -Path $logsDir -Force | Out-Null

Write-Host "Iniciando PostgreSQL local..." -ForegroundColor Cyan
docker compose --project-directory $repoRoot -f (Join-Path $repoRoot "compose.yaml") up -d postgres
if ($LASTEXITCODE -ne 0) {
    throw "No se pudo iniciar PostgreSQL."
}

$databaseReady = $false
for ($attempt = 1; $attempt -le 60; $attempt++) {
    $health = docker inspect --format "{{.State.Health.Status}}" cenz-postgres 2>$null
    if ($health -eq "healthy") {
        $databaseReady = $true
        break
    }
    Start-Sleep -Seconds 1
}
if (-not $databaseReady) {
    throw "PostgreSQL no quedó listo. Revisá: docker logs cenz-postgres"
}

if (-not (Test-Path -LiteralPath (Join-Path $frontendDir "node_modules"))) {
    Write-Host "Instalando dependencias del frontend..." -ForegroundColor Cyan
    Push-Location $frontendDir
    try {
        & npm.cmd ci
        if ($LASTEXITCODE -ne 0) {
            throw "npm ci terminó con error."
        }
    } finally {
        Pop-Location
    }
}

$backendPidFile = Join-Path $runtimeDir "backend.pid"
$frontendPidFile = Join-Path $runtimeDir "frontend.pid"

if (-not (Test-ProcessFromPidFile $backendPidFile)) {
    Write-Host "Iniciando backend..." -ForegroundColor Cyan
    $backendOutLog = Join-Path $logsDir "backend.out.log"
    $backendErrLog = Join-Path $logsDir "backend.err.log"
    $backendCommand = "& go run ./cmd/api 1> '$($backendOutLog.Replace("'", "''"))' 2> '$($backendErrLog.Replace("'", "''"))'"
    Start-DetachedCommand $backendDir $backendCommand $backendPidFile
}

if (-not (Wait-Url "http://localhost:8080/api/health" 60)) {
    throw "El backend no respondió. Revisá .local/logs/backend.err.log"
}

if (-not (Test-ProcessFromPidFile $frontendPidFile)) {
    Write-Host "Iniciando frontend..." -ForegroundColor Cyan
    $frontendOutLog = Join-Path $logsDir "frontend.out.log"
    $frontendErrLog = Join-Path $logsDir "frontend.err.log"
    $frontendCommand = "& npm.cmd run dev 1> '$($frontendOutLog.Replace("'", "''"))' 2> '$($frontendErrLog.Replace("'", "''"))'"
    Start-DetachedCommand $frontendDir $frontendCommand $frontendPidFile
}

if (-not (Wait-Url "http://localhost:5173" 60)) {
    throw "El frontend no respondió. Revisá .local/logs/frontend.err.log"
}

Write-Host ""
Write-Host "Cenz local está funcionando." -ForegroundColor Green
Write-Host "Frontend: http://localhost:5173"
Write-Host "API:      http://localhost:8080/api/health"
Write-Host "Logs:     $logsDir"
Write-Host "Detener:  scripts\stop-local.cmd"
