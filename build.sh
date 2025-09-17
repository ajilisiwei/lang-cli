#!/bin/bash

# 构建脚本 - 多语言打字学习终端工具
# Build script for Multi-language Typing Learning Terminal Tool

set -e

echo "开始构建 lang-cli..."
echo "Building lang-cli..."

# 清理之前的构建文件
if [ -f "lang-cli" ]; then
    echo "清理之前的构建文件..."
    rm -f lang-cli
fi

# 构建 macOS 版本
echo "构建 macOS 版本..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o lang-cli-darwin-arm64 cmd/lang-cli/main.go
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o lang-cli-darwin-amd64 cmd/lang-cli/main.go

# 创建通用二进制文件（如果需要）
if command -v lipo &> /dev/null; then
    echo "创建通用二进制文件..."
    lipo -create -output lang-cli lang-cli-darwin-arm64 lang-cli-darwin-amd64
    rm lang-cli-darwin-arm64 lang-cli-darwin-amd64
else
    # 如果没有 lipo，使用当前架构的版本
    if [ "$(uname -m)" = "arm64" ]; then
        mv lang-cli-darwin-arm64 lang-cli
        rm -f lang-cli-darwin-amd64
    else
        mv lang-cli-darwin-amd64 lang-cli
        rm -f lang-cli-darwin-arm64
    fi
fi

# 设置执行权限
chmod +x lang-cli

echo "构建完成！"
echo "可执行文件: ./lang-cli"
echo "文件大小: $(du -h lang-cli | cut -f1)"
echo "架构信息: $(file lang-cli)"
echo ""
echo "使用方法:"
echo "  ./lang-cli          # 直接运行"
echo "  sudo cp lang-cli /usr/local/bin/  # 安装到系统路径"