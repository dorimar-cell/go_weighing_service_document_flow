param(
    [Parameter(Mandatory = $true)]
    [string]$File,
    [string]$Container = "weighing_postgres",
    [string]$Database = "weighing_db",
    [string]$User = "weighing"
)

if (-not (Test-Path $File)) {
    throw "Backup file not found: $File"
}

Write-Host "Restoring $File into $Database ($Container)..."
Write-Host "WARNING: existing objects will be dropped (--clean --if-exists)."

$command = "docker exec -i $Container pg_restore -U $User -d $Database --clean --if-exists --no-owner < `"$File`""
cmd /c $command
if ($LASTEXITCODE -ne 0) {
    throw "pg_restore failed with exit code $LASTEXITCODE"
}

Write-Host "Restore completed."
