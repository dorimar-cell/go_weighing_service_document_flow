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

Set-Truck 42340
Start-Sleep -Seconds 8
$m1 = Get-MeasurementRaw
$body = '{"measurement":' + $m1 + ',"transportNumber":"REF-A123BC77","trailerNumber":"EA1234-77","driverName":"Ref Driver","cargoId":1,"counterpartyId":1,"carrierId":1,"storageId":1,"operationTypeId":1}'
$created = Post-Json "$server/sessions" $body
"REFS: id=$($created.id) cargo=$($created.cargoName) counterparty=$($created.counterpartyName) carrier=$($created.carrierName) storage=$($created.storageName) operation=$($created.operationTypeName)"

Set-Truck 0
Start-Sleep -Seconds 3
Set-Truck 18340
Start-Sleep -Seconds 8
$m2 = Get-MeasurementRaw
$closed = Post-Json "$server/sessions/$($created.id)/second" ('{"measurement":' + $m2 + '}')
"REFS CLOSED: status=$($closed.status) net=$($closed.netWeightKg) cargo=$($closed.cargoName)"

Set-Truck 0 | Out-Null
