$ErrorActionPreference = 'Stop'
$server = 'http://localhost:3002/spa/api/weighing'
$agent = 'http://127.0.0.1:8095'

function Wait-Health([string]$url) {
    for ($i = 0; $i -lt 60; $i++) {
        try {
            $h = Invoke-RestMethod "$url/health" -TimeoutSec 2
            if ($h.status -eq 'ok') { return }
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

Wait-Health $agent

Set-Truck 42340
Start-Sleep -Seconds 8
$m1 = Get-MeasurementRaw
$body1 = '{"measurement":' + $m1 + ',"transportNumber":"E2E-A123BC77","trailerNumber":"EA1234-77","driverName":"E2E Ivanov"}'
$created = Post-Json "$server/sessions" $body1
"PORTAL FIRST: id=$($created.id) status=$($created.status) weight=$($created.firstWeightKg)"

$open = Invoke-RestMethod "$server/sessions?status=open"
"PORTAL OPEN: total=$($open.total)"

Set-Truck 0
Start-Sleep -Seconds 3
Set-Truck 18340
Start-Sleep -Seconds 8
$m2 = Get-MeasurementRaw
$closed = Post-Json "$server/sessions/$($created.id)/second" ('{"measurement":' + $m2 + '}')
"PORTAL SECOND: status=$($closed.status) net=$($closed.netWeightKg)"

$openAfter = Invoke-RestMethod "$server/sessions?status=open"
$allAfter = Invoke-RestMethod "$server/sessions?status=all"
"PORTAL LIST: open=$($openAfter.total) all=$($allAfter.total)"

$page = Invoke-WebRequest 'http://localhost:3002/weighing/operator' -UseBasicParsing
"PORTAL PAGE: $($page.StatusCode)"

Set-Truck 0 | Out-Null
