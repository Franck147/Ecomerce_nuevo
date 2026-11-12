# Script para iniciar el Ecommerce profesional con barra lateral (Patrones Composite, Decorator y Singleton)

$ErrorActionPreference = "Stop"
$PG_BIN = "C:\Users\11AYUK1P2203A17\pgsql\bin"
$PG_DATA = "C:\Users\11AYUK1P2203A17\pgsql\data"
$PROJECT_ROOT = $PSScriptRoot

Write-Host "==> 1. Verificando PostgreSQL..." -ForegroundColor Cyan
$pgRunning = Get-Process -Name postgres -ErrorAction SilentlyContinue
if (-not $pgRunning) {
    Write-Host "Iniciando PostgreSQL..." -ForegroundColor Yellow
    Start-Process -FilePath "$PG_BIN\pg_ctl.exe" -ArgumentList "-D `"$PG_DATA`" -l `"$PG_DATA\..\server.log`" start" -WindowStyle Hidden
    Start-Sleep -Seconds 2
} else {
    Write-Host "PostgreSQL activo." -ForegroundColor Green
}

Write-Host "==> 2. Iniciando Backend API (:8081)..." -ForegroundColor Cyan
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$PROJECT_ROOT\backend'; .\api.exe"

Write-Host "==> 3. Iniciando Frontend Next.js (:3001)..." -ForegroundColor Cyan
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$PROJECT_ROOT\frontend'; npm run dev"

Write-Host "`n==> ECOMMERCE EN EJECUCIÓN:" -ForegroundColor Green
Write-Host "  - Productos:                   http://localhost:3001/products" -ForegroundColor White
Write-Host "  - Precios:                     http://localhost:3001/pricing" -ForegroundColor White
Write-Host "  - Backend REST API:            http://localhost:8081/health" -ForegroundColor White
