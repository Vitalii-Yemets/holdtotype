param(
    [string]$Model = "small",
    [string]$WhisperVersion = "v1.7.4",
    [string]$Version = "",
    [switch]$StopRunning
)
$ErrorActionPreference = "Stop"
Push-Location $PSScriptRoot
try {
    $running = Get-Process -Name holdtotype, whisper-server -ErrorAction SilentlyContinue
    if ($running) {
        if (-not $StopRunning) {
            Write-Host "Запущены процессы: $(($running.Name | Sort-Object -Unique) -join ', ') — они блокируют файлы в dist\." -ForegroundColor Yellow
            $answer = Read-Host "Завершить их и продолжить сборку? [y/N]"
            if ($answer -notmatch '^[yYдД]') {
                Write-Host "Закройте HoldToType (иконка в трее -> Выход) и запустите сборку снова." -ForegroundColor Red
                exit 1
            }
        }
        $running | Stop-Process -Force
        Start-Sleep -Milliseconds 500
    }

    docker build `
        --file build/Dockerfile `
        --target dist `
        --build-arg WHISPER_MODEL=$Model `
        --build-arg WHISPER_CPP_VERSION=$WhisperVersion `
        --build-arg APP_VERSION=$Version `
        --output type=local,dest=dist `
        .

    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
    Write-Host ""
    Write-Host "Сборка завершена. Запуск: dist\holdtotype.exe" -ForegroundColor Green
} finally {
    Pop-Location
}
