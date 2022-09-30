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

Write-Host "Downloading Mangal version $tag" -ForegroundColor DarkCyan

# download mangal to temp folder
$zip = "$env:TEMP\mangal.zip"
Invoke-WebRequest -Uri $url -OutFile $zip

# extract mangal at temp folder
Expand-Archive -Path $zip -DestinationPath $env:TEMP

# run mangal binary from the unzipped folder
$bin = "$env:TEMP\mangal\mangal.exe"
Start-Process $bin
