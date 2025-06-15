#!/bin/bash
set -e

# 构建脚本：用于编译三端（Linux、macOS、Windows）可执行文件
# 获取最新 linux-command 版本（已注释，需时可启用）
echo "[info] Fetching latest linux-command version..."
VERSION=$(curl -sI https://unpkg.com/linux-command | grep "location" | awk -F'@|/dist/data.json' '{print $2}')
echo "[info] Latest version: $VERSION"

# 生成 Go 版本号文件（已注释，需时可启用）
echo "$VERSION" > ./cmd/version-info
cat ./version >> ./cmd/version-info

echo "[info] Downloading linux-command@$VERSION data.json..."
wget -q "https://unpkg.com/linux-command@$VERSION/dist/data.json" -O ./cmd/linux-command.json

echo "[info] Building binaries for multiple platforms..."
# 定义支持的架构列表
ARCHS=("amd64" "arm64")
# 定义支持的系统列表
OS_LIST=("linux" "darwin" "windows")

# 循环构建所有组合
for os in "${OS_LIST[@]}"; do
  for arch in "${ARCHS[@]}"; do
    # Windows 需特殊处理后缀
    suffix=""
    if [ "$os" = "windows" ]; then
      suffix=".exe"
    fi

    output_name="lsx_${os}_${arch}${suffix}"
    echo "[info] Building $output_name..."

    GOOS=$os GOARCH=$arch go build -o "$output_name" .

    # 仅对非 Windows 文件加执行权限
    if [ "$os" != "windows" ]; then
      chmod +x "$output_name"
    fi
  done
done

# 复制 macOS 版本到本地 bin 目录，方便直接调用
# cp lsx_darwin_arm64 ~/.local/bin/lsx
# 列出所有 lsx 可执行文件
ls | grep lsx
