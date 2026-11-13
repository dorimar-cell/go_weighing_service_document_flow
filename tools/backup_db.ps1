param(
    [string]$Container = "weighing_postgres",
    [string]$Database = "weighing_db",
    [string]$User = "weighing",
    [string]$OutDir = "",
    [int]$KeepDays = 14
)

if ([string]::IsNullOrWhiteSpace($OutDir)) {
    $OutDir = Join-Path (Split-Path -Parent $PSScriptRoot) "backups"
}

New-Item -ItemType Directory -Force -Path $OutDir | Out-Null

$stamp = Get-Date -Format "yyyyMMdd-HHmmss"
$file = Join-Path $OutDir "$Database-$stamp.dump"

# cmd /c нужен, чтобы не испортить бинарный поток дампа перенаправлением PowerShell.
$command = "docker exec $Container pg_dump -U $User -d $Database -Fc > `"$file`""
cmd /c $command
if ($LASTEXITCODE -ne 0) {
    throw "pg_dump failed with exit code $LASTEXITCODE"
}

$info = Get-Item -Path $file
if ($info.Length -le 0) {
    throw "Backup file is empty: $file"
}

Write-Host "Backup created: $file ($([math]::Round($info.Length / 1MB, 2)) MB)"

$limit = (Get-Date).AddDays(-$KeepDays)
Get-ChildItem -Path $OutDir -Filter "$Database-*.dump" |
    Where-Object { $_.LastWriteTime -lt $limit } |
    ForEach-Object {
        Remove-Item -Path $_.FullName -Force
        Write-Host "Removed old backup: $($_.Name)"
    }
