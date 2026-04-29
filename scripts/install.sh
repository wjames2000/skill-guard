#!/bin/sh
# skill-guard 一键安装脚本
# 用法: curl -sfL https://github.com/wjames2000/skill-guard/releases/latest/download/install.sh | sh
# 或: curl -sfL https://github.com/wjames2000/skill-guard/releases/latest/download/install.sh | sh -s v0.1.0

set -eu

REPO="wjames2000/skill-guard"

# 版本优先级: 命令行参数 > SKILLGUARD_VERSION 环境变量 > latest
VERSION="${1:-${SKILLGUARD_VERSION:-latest}}"

# --- 检测系统架构 ---
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "错误: 不支持的架构 '$ARCH'。仅支持 amd64 和 arm64。"
    exit 1
    ;;
esac

# --- 检测操作系统 ---
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux|darwin) ;;
  *)
    echo "错误: 不支持的操作系统 '$OS'。仅支持 Linux 和 macOS。"
    exit 1
    ;;
esac

BINARY="skill-guard-${OS}-${ARCH}"

# --- 构建下载 URL ---
if [ "$VERSION" = "latest" ]; then
  URL="https://github.com/${REPO}/releases/latest/download/${BINARY}"
else
  URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY}"
fi

# --- 选择安装路径 ---
if command -v sudo >/dev/null 2>&1; then
  DEST="/usr/local/bin/skill-guard"
  INSTALL_WITH_SUDO=true
else
  DEST="${HOME}/.local/bin/skill-guard"
  INSTALL_WITH_SUDO=false
  mkdir -p "${HOME}/.local/bin"
fi

# --- 下载二进制文件 ---
echo ">>> 下载 skill-guard (${OS}/${ARCH})..."
echo ">>> 目标路径: ${DEST}"

DOWNLOAD_OK=false
if command -v curl >/dev/null 2>&1; then
  curl -sfL "$URL" -o "${BINARY}" && DOWNLOAD_OK=true
elif command -v wget >/dev/null 2>&1; then
  wget -q "$URL" -O "${BINARY}" && DOWNLOAD_OK=true
else
  echo "错误: 需要 curl 或 wget 来下载文件。请先安装其中之一。"
  exit 1
fi

if [ "$DOWNLOAD_OK" != "true" ]; then
  echo "错误: 下载失败！"
  echo ""
  echo "请检查:"
  echo "  1. 网络连接是否正常"
  echo "  2. URL 是否可以访问: ${URL}"
  echo "  3. 版本 '${VERSION}' 是否存在"
  echo ""
  echo "你也可以手动下载:"
  echo "  https://github.com/${REPO}/releases/latest"
  exit 1
fi

# --- 安装 ---
chmod +x "${BINARY}"

if [ "$INSTALL_WITH_SUDO" = "true" ]; then
  sudo mv "${BINARY}" "${DEST}"
else
  mv "${BINARY}" "${DEST}"
fi

echo ""
echo "✓ 安装完成: ${DEST}"
echo ""

# --- 验证安装 ---
if command -v skill-guard >/dev/null 2>&1; then
  echo "运行以下命令开始使用:"
  echo "  skill-guard --help"
else
  echo "请将 ${DEST%/*} 添加到 PATH 环境变量中:"
  echo "  export PATH=\"${DEST%/*}:\$PATH\""
  echo ""
  echo "或将这一行添加到 ~/.bashrc 或 ~/.zshrc 中:"
  echo "  echo 'export PATH=\"${DEST%/*}:\$PATH\"' >> ~/.zshrc"
fi
