param(
    [Parameter(Mandatory = $true)]
    [string]$File,
    [string]$Container = "weighing_postgres",
    [string]$User = "weighing",
    [string]$TempDb = "",
    [switch]$Keep
)

# Проверка резервной копии: разворачивает дамп во временную БД,
# считает ключевые таблицы и удаляет временную БД. Рабочую БД не трогает.

if (-not (Test-Path $File)) {
    throw "Backup file not found: $File"
}
if ([string]::IsNullOrWhiteSpace($TempDb)) {
    $TempDb = "weighing_verify_" + (Get-Date -Format "yyyyMMdd_HHmmss")
}

$created = $false
$failed = $false

try {
    Write-Host "Creating temporary database $TempDb..."
    docker exec $Container createdb -U $User $TempDb
    if ($LASTEXITCODE -ne 0) {
        throw "createdb failed with exit code $LASTEXITCODE"
    }
    $created = $true

    Write-Host "Restoring $File into $TempDb..."
    $command = "docker exec -i $Container pg_restore -U $User -d $TempDb --no-owner < `"$File`""
    cmd /c $command
    if ($LASTEXITCODE -ne 0) {
        throw "pg_restore failed with exit code $LASTEXITCODE"
    }

    $query = "SELECT " +
        "(SELECT count(*) FROM weighing) AS sessions, " +
        "(SELECT count(*) FROM weighing_measurement) AS measurements, " +
        "(SELECT count(*) FROM dir_users) AS users, " +
        "(SELECT count(*) FROM schema_migrations) AS migrations"
    $counts = docker exec $Container psql -U $User -d $TempDb -t -A -F "|" -c $query
    if ($LASTEXITCODE -ne 0) {
        throw "Verification query failed with exit code $LASTEXITCODE"
    }

    $parts = ($counts | Select-Object -First 1).Split("|")
    if ($parts.Count -lt 4) {
        throw "Unexpected verification output: $counts"
    }
    if ([int]$parts[3] -le 0) {
        throw "No applied migrations in restored database"
    }

    Write-Host "Restore verified:"
    Write-Host "  sessions=$($parts[0]) measurements=$($parts[1]) users=$($parts[2]) migrations=$($parts[3])"
}
catch {
    $failed = $true
    Write-Host "VERIFY FAILED: $($_.Exception.Message)" -ForegroundColor Red
}
finally {
    if ($created) {
        if ($Keep) {
            Write-Host "Temporary database kept: $TempDb"
        }
        else {
            Write-Host "Dropping temporary database $TempDb..."
            docker exec $Container dropdb -U $User --if-exists $TempDb | Out-Null
        }
    }
}

if ($failed) {
    exit 1
}

Write-Host "Backup restore check passed."
