# pubspecmgr 一键安装脚本，适用于 Windows PowerShell
#
# 基本用法:
#   iwr -useb https://raw.githubusercontent.com/liasica/pubspecmgr/master/install.ps1 | iex
#
# 指定版本或自定义安装目录需要用 scriptblock 方式:
#   & ([scriptblock]::Create((iwr -useb https://raw.githubusercontent.com/liasica/pubspecmgr/master/install.ps1))) -Version 1.2.3 -InstallDir C:\tools\pubspecmgr
#
# 或使用环境变量:
#   $env:PUBSPECMGR_VERSION = '1.2.3'
#   iwr -useb https://raw.githubusercontent.com/liasica/pubspecmgr/master/install.ps1 | iex

[CmdletBinding()]
param(
    [string]$Version = $env:PUBSPECMGR_VERSION,
    [string]$InstallDir = $env:PUBSPECMGR_INSTALL_DIR
)

$ErrorActionPreference = 'Stop'

$Repo = 'liasica/pubspecmgr'
$BinName = 'pubspecmgr'

function Write-Step {
    param([string]$Message)
    Write-Host "==> $Message"
}

# 探测 CPU 架构
function Get-TargetArch {
    $arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture
    switch ($arch) {
        'X64'   { return 'amd64' }
        'Arm64' { return 'arm64' }
        default { throw "Unsupported architecture: $arch" }
    }
}

# 解析目标版本
function Resolve-TargetVersion {
    if ($Version) {
        return $Version.TrimStart('v')
    }

    $api = "https://api.github.com/repos/$Repo/releases/latest"
    try {
        $release = Invoke-RestMethod -Uri $api -UseBasicParsing
    }
    catch {
        throw "Failed to query GitHub API: $($_.Exception.Message)"
    }
    if (-not $release.tag_name) {
        throw 'GitHub API did not return tag_name'
    }
    return $release.tag_name.TrimStart('v')
}

# 解析安装目录
function Resolve-TargetDir {
    if ($InstallDir) {
        return $InstallDir
    }
    return Join-Path $env:LOCALAPPDATA 'Programs\pubspecmgr'
}

# 将目录添加到用户级 PATH 环境变量（若尚未存在）
function Add-UserPath {
    param([string]$Dir)
    $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    if (-not $userPath) { $userPath = '' }
    $parts = @($userPath -split ';' | Where-Object { $_ -ne '' })
    if ($parts -contains $Dir) {
        return $false
    }
    $newPath = if ($userPath) { "$userPath;$Dir" } else { $Dir }
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    return $true
}

function Invoke-Install {
    $arch = Get-TargetArch
    $ver = Resolve-TargetVersion
    $dir = Resolve-TargetDir
    $target = Join-Path $dir "$BinName.exe"

    Write-Step "Platform: windows-$arch"
    Write-Step "Version:  v$ver"
    Write-Step "Target:   $target"

    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Force -Path $dir | Out-Null
    }

    # GoReleaser 产物文件名不含 .exe 后缀，下载后落盘时再加
    $url = "https://github.com/$Repo/releases/download/v$ver/$BinName-windows-$arch"

    Write-Step "Downloading $url"
    try {
        Invoke-WebRequest -Uri $url -OutFile $target -UseBasicParsing
    }
    catch {
        throw "Download failed: $($_.Exception.Message)"
    }

    Write-Step "Installed: $target"

    # 冒烟测试
    try {
        & $target -v | Out-Null
    }
    catch {
        Write-Step "Warning: '$target -v' failed, but the binary is installed"
    }

    # PATH 处理
    $added = Add-UserPath -Dir $dir
    if ($added) {
        Write-Step "Added $dir to user PATH"
        Write-Step "Open a new terminal to use '$BinName'"
    }
    else {
        $sessionPath = @($env:PATH -split ';')
        if ($sessionPath -contains $dir) {
            Write-Step "Done. Run: $BinName -h"
        }
        else {
            Write-Step "$dir is already in user PATH but not in current session"
            Write-Step "Open a new terminal, or run: `$env:PATH += ';$dir'"
        }
    }
}

Invoke-Install
