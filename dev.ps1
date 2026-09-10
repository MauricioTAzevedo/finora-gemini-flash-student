param (
    [Parameter(Position=0)]
    [string]$Command = "help"
)

switch ($Command) {
    "dev" {
        Write-Host "Starting Finora API and Web..." -ForegroundColor Cyan
        Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd services/api; go run cmd/server/main.go"
        Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd apps/web; pnpm dev"
    }
    "test" {
        Write-Host "=== Running Go Ledger Tests ===" -ForegroundColor Green
        Push-Location services/api
        go test -v ./...
        Pop-Location

        Write-Host "=== Running Python AI Tests ===" -ForegroundColor Green
        Push-Location services/ai
        python -m pytest -v
        Pop-Location
    }
    "docker-up" {
        Write-Host "Starting Docker infrastructure..." -ForegroundColor Cyan
        docker compose up -d
    }
    "docker-down" {
        Write-Host "Stopping Docker infrastructure..." -ForegroundColor Cyan
        docker compose down
    }
    default {
        Write-Host "Finora Gemini Flash Student — Windows PowerShell Task Runner" -ForegroundColor Yellow
        Write-Host "  .\dev.ps1 dev         - Run API and Web apps"
        Write-Host "  .\dev.ps1 test        - Run all test suites"
        Write-Host "  .\dev.ps1 docker-up   - Start local containers"
        Write-Host "  .\dev.ps1 docker-down - Stop local containers"
    }
}
