#!/bin/bash
set -e

# # 获取最新版本号
# echo "[info] Fetching latest linux-command version..."
# VERSION=$(curl -sI https://unpkg.com/linux-command | grep "location" | awk -F'@|/dist/data.json' '{print $2}')
# echo "[info] Latest version: $VERSION"

# # 生成Go版本号文件（覆盖写入）
# echo -n "$VERSION" > ./cmd/version-info
# cat ./version >>> ./cmd/version-info

# echo "[info] Downloading linux-command@$VERSION data.json..."
# wget -q "https://unpkg.com/linux-command@$VERSION/dist/data.json" -O ./cmd/linux-command.json

echo "[info] Building binaries..."
GOOS=linux go build -o lsx_linux_amd64 .
GOOS=darwin go build -o lsx_darwin_amd64 .
GOOS=windows go build -o lsx_windows_amd64.exe .

# 仅对非 Windows 可执行文件加可执行权限
chmod +x lsx_linux_amd64 lsx_darwin_amd64 || true
cp lsx_darwin_amd64 ~/.local/bin/lsx
ls | grep lsx
