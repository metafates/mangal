
# get latest release
$release_url = "https://api.github.com/repos/metafates/mangal/releases"
$tag = (Invoke-WebRequest -Uri $release_url -UseBasicParsing | ConvertFrom-Json)[0].tag_name
$loc = "$HOME\AppData\Local\mangal"
$url = ""
$arch = $env:PROCESSOR_ARCHITECTURE
$releases_api_url = "https://github.com/metafates/mangal/releases/download/$tag/mangal_${tag.substring(1)}_Windows"

if ($arch -eq "AMD64") {
    $url = "${releases_api_url}_x86_64.zip"
} elseif ($arch -eq "x86") {
    $url = "${releases_api_url}_i386.zip"
} elseif ($arch -eq "arm64") {
    $url = "${releases_api_url}_arm64.zip"
}

if (Test-Path -path $loc) {
    Remove-Item $loc -Recurse -Force
}

Write-Host "Downloading Mangal version $tag" -ForegroundColor DarkCyan
Invoke-WebRequest $url -outfile mangal.zip
Expand-Archive mangal.zip

.\mangal.exe

Remove-Item mangal* -Recurse -Force
