!/bin/bash

VERSION=$(curl -sI https://unpkg.com/linux-command | grep "location" | awk -F'@|/dist/data.json' '{print $2}')

# 生成Go版本号文件
echo -n "$VERSION" > ./cmd/version
cat ./version >> ./cmd/version

wget "https://unpkg.com/linux-command@$VERSION/dist/data.json" -O ./cmd/linux-command.json

GOOS=linux go build -o lsx_linux_amd64 .
GOOS=darwin go build -o lsx_darwin_amd64 .
GOOS=windows go build -o lsx_windows_amd64.exe .

chmod +x lsx_linux_amd64 lsx_darwin_amd64 lsx_windows_amd64.exe

ls | grep lsx
