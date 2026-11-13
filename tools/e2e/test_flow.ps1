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

$pass = Post-Json "$server/passes" '{"number":"E2E-FLOW-001","transportNumber":"FLOW-A123BC77","passTypeId":3}'
"FLOW PASS: id=$($pass.id) enabled=$($pass.passEnabled)"

Set-Truck 42340
Start-Sleep -Seconds 8
$m = Get-MeasurementRaw
$body = '{"measurement":' + $m + ',"transportNumber":"","passId":' + $pass.id + ',"tareKg":15000}'
$flow = Post-Json "$server/sessions/flow" $body
"FLOW: id=$($flow.id) status=$($flow.status) brutto=$($flow.firstWeightKg) tara=$($flow.secondWeightKg) net=$($flow.netWeightKg) transport=$($flow.transportNumber) pass=$($flow.passNumber)"

$passAfter = Invoke-RestMethod "$server/passes/$($pass.id)"
"FLOW PASS AFTER: enabled=$($passAfter.passEnabled)"

$page = Invoke-WebRequest 'http://localhost:3002/weighing/operator' -UseBasicParsing
"PORTAL PAGE: $($page.StatusCode)"

Set-Truck 0 | Out-Null
