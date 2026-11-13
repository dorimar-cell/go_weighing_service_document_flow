param(
    [string]$TaskName = "Weighing DB Backup",
    [string]$Time = "02:00",
    [int]$KeepDays = 14,
    [switch]$WhatIf
)

$scriptDir = Split-Path -Parent $PSScriptRoot
$backupScript = Join-Path $scriptDir "tools\backup_db.ps1"

if (-not (Test-Path $backupScript)) {
    throw "Backup script not found: $backupScript"
}

$action = New-ScheduledTaskAction -Execute "powershell.exe" `
    -Argument "-ExecutionPolicy Bypass -File `"$backupScript`" -KeepDays $KeepDays"
$trigger = New-ScheduledTaskTrigger -Daily -At $Time
$settings = New-ScheduledTaskSettingsSet -StartWhenAvailable -DontStopOnIdleEnd

if ($WhatIf) {
    Write-Host "Task:    $TaskName"
    Write-Host "Action:  powershell.exe $($action.Arguments)"
    Write-Host "Trigger: daily at $Time"
    Write-Host "Keep:    $KeepDays days"
    return
}

Register-ScheduledTask -TaskName $TaskName -Action $action -Trigger $trigger `
    -Settings $settings -Description "Daily PostgreSQL backup for weighing_db" -Force | Out-Null

Write-Host "Scheduled task '$TaskName' registered: daily at $Time, keep $KeepDays days."
