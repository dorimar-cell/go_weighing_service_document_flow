$ErrorActionPreference = 'Stop'
$server = 'http://127.0.0.1:8091/spa/api/weighing'
$agent = 'http://127.0.0.1:8095'

function Wait-Health([string]$url) {
    for ($i = 0; $i -lt 60; $i++) {
        try {
            $h = Invoke-RestMethod "$url/health" -TimeoutSec 2
            if ($h.status -eq 'ok') { return $h }
        } catch { Start-Sleep -Milliseconds 500 }
    }
    throw "health timeout: $url"
}

function Set-Truck([double]$kg) {
    Invoke-RestMethod -Uri "$agent/api/simulate/truck" -Method Post `
        -Headers @{ 'X-Station-Token' = 'dev-token' } `
        -ContentType 'application/json' -Body "{`"weightKg`":$kg}" | Out-Null
}

function Get-MeasurementRaw {
    return (Invoke-WebRequest "$agent/api/measurement" -UseBasicParsing).Content
}

function Post-Json([string]$url, [string]$json) {
    return Invoke-RestMethod -Uri $url -Method Post `
        -ContentType 'application/json; charset=utf-8' `
        -Body ([System.Text.Encoding]::UTF8.GetBytes($json))
}

Wait-Health $agent | Out-Null

Set-Truck 42340
Start-Sleep -Seconds 8
$m1 = Get-MeasurementRaw
$body1 = '{"measurement":' + $m1 + ',"transportNumber":"A123BC77","trailerNumber":"EA1234-77","driverName":"Ivanov Ivan"}'
$session1 = Post-Json "$server/sessions" $body1
"FIRST: id=$($session1.id) status=$($session1.status) weight=$($session1.firstWeightKg) transport=$($session1.transportNumber)"

$again = Post-Json "$server/sessions" $body1
"IDEMPOTENT: id=$($again.id) same=$($again.id -eq $session1.id)"

Set-Truck 0
Start-Sleep -Seconds 3
Set-Truck 18340
Start-Sleep -Seconds 8
$m2 = Get-MeasurementRaw
$body2 = '{"measurement":' + $m2 + '}'
$closed = Post-Json "$server/sessions/$($session1.id)/second" $body2
"SECOND: status=$($closed.status) first=$($closed.firstWeightKg) second=$($closed.secondWeightKg) net=$($closed.netWeightKg)"

$open = Invoke-RestMethod "$server/sessions?status=open"
$all = Invoke-RestMethod "$server/sessions?status=all"
"LIST: open=$($open.total) all=$($all.total)"

# Отмена: новое измерение -> новый рейс -> cancel
Set-Truck 0
Start-Sleep -Seconds 3
Set-Truck 25000
Start-Sleep -Seconds 8
$m3 = Get-MeasurementRaw
$body3 = '{"measurement":' + $m3 + ',"transportNumber":"B456EK50"}'
$session3 = Post-Json "$server/sessions" $body3
Invoke-RestMethod -Uri "$server/sessions/$($session3.id)/cancel" -Method Post | Out-Null
$cancelled = Invoke-RestMethod "$server/sessions/$($session3.id)"
"CANCEL: id=$($cancelled.id) status=$($cancelled.status)"

# Неверная подпись
$bad = $m3 -replace '"checksum":"[a-f0-9]+"', '"checksum":"0000000000000000000000000000000000000000000000000000000000000000"'
try {
    Post-Json "$server/sessions" ('{"measurement":' + $bad + ',"transportNumber":"C789OP99"}') | Out-Null
    "BAD SIGNATURE: unexpectedly accepted"
} catch {
    "BAD SIGNATURE: rejected status=$([int]$_.Exception.Response.StatusCode)"
}

Set-Truck 0 | Out-Null
