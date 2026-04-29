#!/bin/sh
set -eu
REPO="wjames2000/skill-guard"
VERSION="${1:-latest}"
ARCH=$(uname -m)
case $ARCH in x86_64) ARCH="amd64" ;; aarch64|arm64) ARCH="arm64" ;; *) echo "不支持的架构: $ARCH"; exit 1 ;; esac
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
BINARY="skill-guard-${OS}-${ARCH}"
[ "$OS" = "windows" ] && BINARY="${BINARY}.exe"
URL="https://github.com/${REPO}/releases/${VERSION}/download/${BINARY}"
DEST="/usr/local/bin/skill-guard"
echo "下载 ${URL}..."
if command -v curl >/dev/null 2>&1; then curl -sfL "$URL" -o "$BINARY"
elif command -v wget >/dev/null 2>&1; then wget -q "$URL" -O "$BINARY"
else echo "需要 curl 或 wget"; exit 1; fi
chmod +x "$BINARY"
sudo mv "$BINARY" "$DEST"
echo "安装完成: ${DEST}"
