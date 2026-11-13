param(
    [int]$Days = 90,
    [string]$Container = "weighing_postgres",
    [string]$Database = "weighing_db",
    [string]$User = "weighing",
    [string]$PhotosDir = "",
    [switch]$DryRun
)

# Удаляет фото и записи weighting_photo старше $Days дней.
# Запускать вручную или из планировщика ОС; сервис не чистит архив сам.

if ([string]::IsNullOrWhiteSpace($PhotosDir)) {
    $PhotosDir = Join-Path (Split-Path -Parent $PSScriptRoot) "data/photos"
}

if (-not (Test-Path $PhotosDir)) {
    throw "Photos directory not found: $PhotosDir"
}

$query = "SELECT id, id_weighing, file_name FROM weighting_photo " +
    "WHERE created_date < now() - interval '$Days days' ORDER BY id"
$rows = docker exec $Container psql -U $User -d $Database -t -A -F "|" -c $query 2>&1
if ($LASTEXITCODE -ne 0) {
    throw "psql failed with exit code $LASTEXITCODE"
}

$items = @($rows | Where-Object { $_ -match '\|' })
if ($items.Count -eq 0) {
    Write-Host "No photos older than $Days days."
    exit 0
}

Write-Host "Photos to remove: $($items.Count) (older than $Days days)"
$deletedFiles = 0
$missingFiles = 0
$ids = New-Object System.Collections.Generic.List[string]

foreach ($row in $items) {
    $parts = $row.Split("|")
    if ($parts.Count -lt 3) {
        continue
    }
    $id = $parts[0].Trim()
    $sessionId = $parts[1].Trim()
    $fileName = $parts[2].Trim()
    $ids.Add($id)

    $path = Join-Path (Join-Path $PhotosDir $sessionId) $fileName
    if (Test-Path $path) {
        if ($DryRun) {
            Write-Host "Would delete: $path"
        }
        else {
            Remove-Item -Path $path -Force
        }
        $deletedFiles++
    }
    else {
        $missingFiles++
    }
}

Write-Host "Files: $deletedFiles, already missing: $missingFiles"

if ($DryRun) {
    Write-Host "Dry run: database rows are kept."
    exit 0
}

$idList = $ids -join ","
$delete = "DELETE FROM weighting_photo WHERE id IN ($idList)"
docker exec $Container psql -U $User -d $Database -c $delete | Out-Null
if ($LASTEXITCODE -ne 0) {
    throw "Failed to delete photo rows"
}

Write-Host "Removed rows: $($ids.Count)"
