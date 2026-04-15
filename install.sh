#!/bin/sh
# pubspecmgr 一键安装脚本，适用于 macOS 和 Linux
#
# 基本用法:
#   curl -fsSL https://raw.githubusercontent.com/liasica/pubspecmgr/master/install.sh | sh
#
# 指定版本:
#   curl -fsSL https://raw.githubusercontent.com/liasica/pubspecmgr/master/install.sh | VERSION=1.2.3 sh
#
# 自定义安装目录:
#   curl -fsSL https://raw.githubusercontent.com/liasica/pubspecmgr/master/install.sh | INSTALL_DIR=$HOME/bin sh

set -eu

REPO="liasica/pubspecmgr"
BIN_NAME="pubspecmgr"

log() {
    printf '==> %s\n' "$*"
}

err() {
    printf 'Error: %s\n' "$*" >&2
    exit 1
}

# 探测操作系统
detect_os() {
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$os" in
        darwin|linux) printf '%s' "$os" ;;
        *) err "Unsupported OS: $os" ;;
    esac
}

# 探测 CPU 架构
detect_arch() {
    arch=$(uname -m)
    case "$arch" in
        x86_64|amd64) printf '%s' "amd64" ;;
        aarch64|arm64) printf '%s' "arm64" ;;
        *) err "Unsupported architecture: $arch" ;;
    esac
}

# 解析目标版本，优先使用 VERSION 环境变量，否则查询 GitHub API 的 latest
resolve_version() {
    if [ -n "${VERSION:-}" ]; then
        printf '%s' "$VERSION" | sed 's/^v//'
        return
    fi

    api_url="https://api.github.com/repos/${REPO}/releases/latest"
    latest=$(curl -fsSL "$api_url" \
        | grep '"tag_name"' \
        | head -1 \
        | sed 's/.*"tag_name"[[:space:]]*:[[:space:]]*"v\{0,1\}\([^"]*\)".*/\1/')

    [ -n "$latest" ] || err "Failed to resolve latest version from GitHub API"
    printf '%s' "$latest"
}

# 解析安装目录，优先级: INSTALL_DIR > /usr/local/bin (可写) > $HOME/.local/bin
resolve_install_dir() {
    if [ -n "${INSTALL_DIR:-}" ]; then
        printf '%s' "$INSTALL_DIR"
        return
    fi

    if [ -w "/usr/local/bin" ] 2>/dev/null; then
        printf '%s' "/usr/local/bin"
    else
        printf '%s' "$HOME/.local/bin"
    fi
}

main() {
    command -v curl >/dev/null 2>&1 || err "curl is required but not found"

    os=$(detect_os)
    arch=$(detect_arch)
    version=$(resolve_version)
    install_dir=$(resolve_install_dir)
    target="${install_dir}/${BIN_NAME}"

    log "Platform: ${os}-${arch}"
    log "Version:  v${version}"
    log "Target:   ${target}"

    mkdir -p "$install_dir" || err "Failed to create $install_dir"
    [ -w "$install_dir" ] || err "Install directory is not writable: $install_dir (set INSTALL_DIR to override)"

    url="https://github.com/${REPO}/releases/download/v${version}/${BIN_NAME}-${os}-${arch}"
    tmpfile=$(mktemp)
    # shellcheck disable=SC2064
    trap "rm -f '$tmpfile'" EXIT INT TERM

    log "Downloading ${url}"
    curl -fL --progress-bar -o "$tmpfile" "$url" || err "Download failed"

    mv "$tmpfile" "$target"
    chmod +x "$target"

    log "Installed: ${target}"

    # 冒烟测试
    if ! "$target" -v >/dev/null 2>&1; then
        log "Warning: '${target} -v' failed, but the binary is installed"
    fi

    # PATH 检查
    case ":$PATH:" in
        *":${install_dir}:"*)
            log "Done. Run: ${BIN_NAME} -h"
            ;;
        *)
            log "Note: ${install_dir} is not in your PATH"
            log "Add the following line to your shell profile:"
            log "  export PATH=\"${install_dir}:\$PATH\""
            ;;
    esac
}

main "$@"
