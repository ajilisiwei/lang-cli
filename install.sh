#!/bin/bash

# 安装脚本 - 多语言打字学习终端工具
# Install script for Multi-language Typing Learning Terminal Tool

set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="lang-cli"

echo "多语言打字学习终端工具 - 安装脚本"
echo "Multi-language Typing Learning Terminal Tool - Install Script"
echo ""

# 检查是否存在可执行文件
if [ ! -f "$BINARY_NAME" ]; then
    echo "错误: 找不到可执行文件 '$BINARY_NAME'"
    echo "请先运行 './build.sh' 构建项目"
    echo ""
    echo "Error: Executable file '$BINARY_NAME' not found"
    echo "Please run './build.sh' to build the project first"
    exit 1
fi

# 检查安装目录是否存在
if [ ! -d "$INSTALL_DIR" ]; then
    echo "创建安装目录: $INSTALL_DIR"
    sudo mkdir -p "$INSTALL_DIR"
fi

# 复制可执行文件
echo "安装 $BINARY_NAME 到 $INSTALL_DIR..."
sudo cp "$BINARY_NAME" "$INSTALL_DIR/"

# 设置权限
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo "安装完成！"
echo "Installation completed!"
echo ""
echo "现在您可以在任何地方使用 '$BINARY_NAME' 命令"
echo "You can now use '$BINARY_NAME' command from anywhere"
echo ""
echo "使用方法 / Usage:"
echo "  $BINARY_NAME"
echo ""
echo "卸载 / Uninstall:"
echo "  sudo rm $INSTALL_DIR/$BINARY_NAME"