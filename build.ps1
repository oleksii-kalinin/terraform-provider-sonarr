$HostName = "registry.terraform.io"
$Namespace = "oleksii-kalinin"
$Name = "sonarr"
$Version = "0.0.1"
$Arch = "windows_amd64"

$TargetDir = "$env:APPDATA\terraform.d\plugins\$HostName\$Namespace\$Name\$Version\$Arch"
$BinaryName = "terraform-provider-${Name}.exe"

Write-Host "1. Building provider..."
go build -o $BinaryName

if ($LASTEXITCODE -ne 0) {
    Write-Error "Build failed"
    exit 1
}

Write-Host "2. Creating directory structure: $TargetDir"
New-Item -Path $TargetDir -ItemType Directory -Force | Out-Null

Write-Host "3. Installing binary..."
Move-Item -Path $BinaryName -Destination "$TargetDir\$BinaryName" -Force

Write-Host "Done! Provider installed."
