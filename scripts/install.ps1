# WAS NOT TESTED

$release_url = "https://api.github.com/repos/metafates/mangal/releases"
$tag = (Invoke-WebRequest -Uri $release_url -UseBasicParsing | ConvertFrom-Json)[0].tag_name
$version = $tag.substring(1)
$loc = "$HOME\AppData\Local\mangal"
$url = ""
$arch = $env:PROCESSOR_ARCHITECTURE
$releases_api_url = "https://github.com/metafates/mangal/releases/download/$tag/mangal_${version}_Windows"

if ($arch -eq "AMD64")
{
    $url = "${releases_api_url}_x86_64.zip"
}
elseif ($arch -eq "x86")
{
    $url = "${releases_api_url}_i386.zip"
}
elseif ($arch -eq "arm64")
{
    $url = "${releases_api_url}_arm64.zip"
}

if (Test-Path -path $loc)
{
    Remove-Item $loc -Recurse -Force
}

Write-Host "Installing mangal version $tag" -ForegroundColor DarkCyan

Invoke-WebRequest $url -outfile mangal.zip

Expand-Archive mangal.zip

New-Item -ItemType "directory" -Path $loc

Move-Item -Path mangal\mangal.exe -Destination $loc

Remove-Item mangal* -Recurse -Force

[System.Environment]::SetEnvironmentVariable("Path", $Env:Path + ";$loc", [System.EnvironmentVariableTarget]::User)

if (Test-Path -path $loc)
{
    Write-Host "Mangal version $tag installed successfully" -ForegroundColor Green
}
else
{
    Write-Host "Download failed" -ForegroundColor Red
    Write-Host "Please try again later" -ForegroundColor Red
}
