param(
    [string]$AppPath = "C:\WeightingApp\WeightingApp.exe",
    [switch]$Test
)

# Prints the SQL Server connection target (server/database/user) for WeightingApp.
# MUST be run on the station where WeightingApp is installed, under the same Windows user
# that runs the application: the connection string is encrypted with DPAPI for that user/machine.
# The password is never printed.

$ErrorActionPreference = "Stop"
Add-Type -AssemblyName System.Configuration

function Fail([string]$message) {
    Write-Host "ERROR: $message" -ForegroundColor Red
    exit 1
}

if (-not (Test-Path -LiteralPath $AppPath)) {
    Fail "Application not found: $AppPath. Pass the real path via -AppPath."
}

$configPath = "$AppPath.config"
if (-not (Test-Path -LiteralPath $configPath)) {
    Fail "Config not found: $configPath"
}

Write-Host "Machine : $env:COMPUTERNAME"
Write-Host "User    : $env:USERDOMAIN\$env:USERNAME"
Write-Host "Config  : $configPath"
Write-Host ""

try {
    $config = [System.Configuration.ConfigurationManager]::OpenExeConfiguration($AppPath)
    $section = $config.GetSection("connectionStrings")
}
catch {
    Fail "Cannot decrypt the connection string. Run this script on the original station under the same Windows user. Details: $($_.Exception.Message)"
}

$found = $false
foreach ($entry in $section.ConnectionStrings) {
    if ($entry.Name -ne "SQLContext") {
        continue
    }

    $found = $true
    $builder = New-Object System.Data.SqlClient.SqlConnectionStringBuilder($entry.ConnectionString)

    Write-Host "Connection      : $($entry.Name)"
    Write-Host "Server          : $($builder.DataSource)"
    Write-Host "Database        : $($builder.InitialCatalog)"
    Write-Host "User            : $($builder.UserID)"
    Write-Host "IntegratedAuth  : $($builder.IntegratedSecurity)"
    Write-Host "Encrypt         : $($builder.Encrypt)"
    Write-Host ""

    if ($Test) {
        Write-Host "Testing connection..."
        $connection = New-Object System.Data.SqlClient.SqlConnection($entry.ConnectionString)
        $connection.Open()
        $command = $connection.CreateCommand()
        $command.CommandText = "SELECT @@VERSION AS Version, DB_NAME() AS DbName, SUSER_SNAME() AS LoginName"
        $reader = $command.ExecuteReader()
        if ($reader.Read()) {
            Write-Host "OK. Version   : $($reader['Version'])"
            Write-Host "OK. Database  : $($reader['DbName'])"
            Write-Host "OK. Login     : $($reader['LoginName'])"
        }
        $reader.Close()
        $connection.Close()
    }
}

if (-not $found) {
    Fail "Connection 'SQLContext' not found in $configPath"
}
