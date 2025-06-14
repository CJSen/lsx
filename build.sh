#!/bin/bash
set -e

# # 构建脚本：用于编译三端（Linux、macOS、Windows）可执行文件
# # 获取最新 linux-command 版本（已注释，需时可启用）
# echo "[info] Fetching latest linux-command version..."
# VERSION=$(curl -sI https://unpkg.com/linux-command | grep "location" | awk -F'@|/dist/data.json' '{print $2}')
# echo "[info] Latest version: $VERSION"

# # 生成 Go 版本号文件（已注释，需时可启用）
# echo -n "$VERSION" > ./cmd/version-info
# cat ./version >>> ./cmd/version-info

# echo "[info] Downloading linux-command@$VERSION data.json..."
# wget -q "https://unpkg.com/linux-command@$VERSION/dist/data.json" -O ./cmd/linux-command.json

echo "[info] Building binaries..."
# 分别为 Linux、macOS、Windows 构建可执行文件
GOOS=linux go build -o lsx_linux_amd64 .
GOOS=darwin go build -o lsx_darwin_amd64 .
GOOS=windows go build -o lsx_windows_amd64.exe .

# 仅对非 Windows 可执行文件加执行权限
chmod +x lsx_linux_amd64 lsx_darwin_amd64 || true
# 复制 macOS 版本到本地 bin 目录，方便直接调用
cp lsx_darwin_amd64 ~/.local/bin/lsx
# 列出所有 lsx 可执行文件
ls | grep lsx
