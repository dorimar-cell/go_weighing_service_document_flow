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
$session = Post-Json "$server/sessions" ('{"measurement":' + $m1 + ',"transportNumber":"PHOTO-A123BC77"}')
"SESSION: id=$($session.id)"

$photo1 = Post-Json "http://localhost:3002/api/weighing-agent/api/cameras/1/capture" '{}'
"CAMERA: file=$($photo1.fileName) type=$($photo1.contentType) base64len=$($photo1.base64.Length)"
$upload1 = Post-Json "$server/sessions/$($session.id)/photos" ('{"base64":"' + $photo1.base64 + '","contentType":"' + $photo1.contentType + '","firstWeight":true}')
"PHOTO1: id=$($upload1.id) url=$($upload1.url)"

Set-Truck 0
Start-Sleep -Seconds 3
Set-Truck 18340
Start-Sleep -Seconds 8
$m2 = Get-MeasurementRaw
$closed = Post-Json "$server/sessions/$($session.id)/second" ('{"measurement":' + $m2 + '}')
$photo2 = Post-Json "http://localhost:3002/api/weighing-agent/api/cameras/1/capture" '{}'
$upload2 = Post-Json "$server/sessions/$($session.id)/photos" ('{"base64":"' + $photo2.base64 + '","contentType":"' + $photo2.contentType + '","firstWeight":false}')
"PHOTO2: id=$($upload2.id) firstWeight=$($upload2.firstWeight)"

$list = Invoke-RestMethod "$server/sessions/$($session.id)/photos"
"PHOTOS: total=$($list.items.Count)"
$file = Invoke-WebRequest "http://localhost:3002$($upload1.url)" -UseBasicParsing
$bytes = $file.Content
"FILE: status=$($file.StatusCode) type=$($file.Headers['Content-Type']) jpeg=$($bytes[0] -eq 0xFF -and $bytes[1] -eq 0xD8 -and $bytes[2] -eq 0xFF) bytes=$($bytes.Length)"

$page = Invoke-WebRequest 'http://localhost:3002/weighing/journal' -UseBasicParsing
"PAGE: $($page.StatusCode)"

Set-Truck 0 | Out-Null
