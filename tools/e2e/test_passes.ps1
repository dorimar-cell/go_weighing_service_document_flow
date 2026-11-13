$ErrorActionPreference = 'Stop'
$server = 'http://localhost:3002/spa/api/weighing'
$agent = 'http://127.0.0.1:8095'

function Wait-Health([string]$url) {
    for ($i = 0; $i -lt 60; $i++) {
        try { $h = Invoke-RestMethod "$url/health" -TimeoutSec 2; if ($h.status -eq 'ok') { return } } catch { Start-Sleep -Milliseconds 500 }
    }
    throw "health timeout: $url"
}

function Set-Truck([double]$kg) {
    Invoke-RestMethod -Uri "$agent/api/simulate/truck" -Method Post `
        -Headers @{ 'X-Station-Token' = 'dev-token' } `
        -ContentType 'application/json' -Body "{`"weightKg`":$kg}" | Out-Null
}

function Get-MeasurementRaw { (Invoke-WebRequest "$agent/api/measurement" -UseBasicParsing).Content }

function Post-Json([string]$url, [string]$json) {
    Invoke-RestMethod -Uri $url -Method Post -ContentType 'application/json; charset=utf-8' `
        -Body ([System.Text.Encoding]::UTF8.GetBytes($json))
}

Wait-Health $agent

$pass = Post-Json "$server/passes" '{"number":"E2E-P-001","transportNumber":"PASS-A123BC77","passTypeId":3}'
"PASS CREATED: id=$($pass.id) number=$($pass.number) type=$($pass.passTypeName) enabled=$($pass.passEnabled)"

$active = Invoke-RestMethod "$server/passes?status=active"
"PASSES ACTIVE: total=$($active.total)"

Set-Truck 42340
Start-Sleep -Seconds 8
$m1 = Get-MeasurementRaw
$body = '{"measurement":' + $m1 + ',"transportNumber":"","passId":' + $pass.id + '}'
$session = Post-Json "$server/sessions" $body
"SESSION BY PASS: id=$($session.id) transport=$($session.transportNumber) pass=$($session.passNumber)"

Set-Truck 0
Start-Sleep -Seconds 3
Set-Truck 18340
Start-Sleep -Seconds 8
$m2 = Get-MeasurementRaw
$closed = Post-Json "$server/sessions/$($session.id)/second" ('{"measurement":' + $m2 + '}')
"SESSION CLOSED: status=$($closed.status) net=$($closed.netWeightKg) pass=$($closed.passNumber)"

$activeAfter = Invoke-RestMethod "$server/passes?status=active"
$passAfter = Invoke-RestMethod "$server/passes/$($pass.id)"
"PASS AFTER CLOSE: enabled=$($passAfter.passEnabled) active_total=$($activeAfter.total)"

$page = Invoke-WebRequest 'http://localhost:3002/weighing/operator' -UseBasicParsing
"PORTAL PAGE: $($page.StatusCode)"

Set-Truck 0 | Out-Null
