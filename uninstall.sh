#!/bin/bash

# 卸载脚本 - 多语言打字学习终端工具
# Uninstall script for Multi-language Typing Learning Terminal Tool

set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="lang-cli"
LANG_CLI_DIR="$HOME/.lang-cli"

echo "多语言打字学习终端工具 - 卸载脚本"
echo "Multi-language Typing Learning Terminal Tool - Uninstall Script"
echo ""

# 检查是否以管理员权限运行（用于删除系统目录中的文件）
if [ "$EUID" -ne 0 ]; then
    echo "注意: 需要管理员权限来删除系统目录中的文件"
    echo "Note: Administrator privileges required to remove files from system directories"
    echo ""
fi

# 询问用户确认
echo "此操作将完全删除 lang-cli 及其所有数据，包括："
echo "This operation will completely remove lang-cli and all its data, including:"
echo "  - 可执行文件 / Executable file: $INSTALL_DIR/$BINARY_NAME"
echo "  - 用户数据目录 / User data directory: $LANG_CLI_DIR"
echo "    (包含配置文件、资源文件、用户数据等)"
echo "    (Contains configuration files, resource files, user data, etc.)"
echo ""
read -p "确认卸载？(y/N) / Confirm uninstall? (y/N): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "卸载已取消"
    echo "Uninstall cancelled"
    exit 0
fi

echo ""
echo "开始卸载..."
echo "Starting uninstall..."
echo ""

# 删除可执行文件
if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
    echo "删除可执行文件: $INSTALL_DIR/$BINARY_NAME"
    echo "Removing executable file: $INSTALL_DIR/$BINARY_NAME"
    if sudo rm -f "$INSTALL_DIR/$BINARY_NAME"; then
        echo "✓ 可执行文件删除成功"
        echo "✓ Executable file removed successfully"
    else
        echo "✗ 删除可执行文件失败，请手动删除"
        echo "✗ Failed to remove executable file, please remove manually"
    fi
else
    echo "可执行文件不存在: $INSTALL_DIR/$BINARY_NAME"
    echo "Executable file not found: $INSTALL_DIR/$BINARY_NAME"
fi

echo ""

# 删除用户数据目录
if [ -d "$LANG_CLI_DIR" ]; then
    echo "删除用户数据目录: $LANG_CLI_DIR"
    echo "Removing user data directory: $LANG_CLI_DIR"
    
    # 显示将要删除的内容
    echo "目录内容 / Directory contents:"
    ls -la "$LANG_CLI_DIR" 2>/dev/null || echo "无法列出目录内容 / Cannot list directory contents"
    echo ""
    
    # 再次确认删除用户数据
    read -p "确认删除用户数据目录？(y/N) / Confirm deletion of user data directory? (y/N): " -n 1 -r
    echo ""
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        if rm -rf "$LANG_CLI_DIR"; then
            echo "✓ 用户数据目录删除成功"
            echo "✓ User data directory removed successfully"
        else
            echo "✗ 删除用户数据目录失败"
            echo "✗ Failed to remove user data directory"
        fi
    else
        echo "保留用户数据目录: $LANG_CLI_DIR"
        echo "User data directory preserved: $LANG_CLI_DIR"
    fi
else
    echo "用户数据目录不存在: $LANG_CLI_DIR"
    echo "User data directory not found: $LANG_CLI_DIR"
fi

echo ""

# 检查卸载结果
UNINSTALL_SUCCESS=true

if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
    echo "⚠️  可执行文件仍然存在: $INSTALL_DIR/$BINARY_NAME"
    echo "⚠️  Executable file still exists: $INSTALL_DIR/$BINARY_NAME"
    UNINSTALL_SUCCESS=false
fi

if [ -d "$LANG_CLI_DIR" ]; then
    echo "⚠️  用户数据目录仍然存在: $LANG_CLI_DIR"
    echo "⚠️  User data directory still exists: $LANG_CLI_DIR"
fi

if [ "$UNINSTALL_SUCCESS" = true ]; then
    echo "🎉 卸载完成！"
    echo "🎉 Uninstall completed!"
    echo ""
    echo "lang-cli 已从您的系统中完全移除"
    echo "lang-cli has been completely removed from your system"
else
    echo "⚠️  卸载未完全成功，请手动删除剩余文件"
    echo "⚠️  Uninstall not completely successful, please manually remove remaining files"
    echo ""
    echo "手动删除命令 / Manual removal commands:"
    if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        echo "  sudo rm -f $INSTALL_DIR/$BINARY_NAME"
    fi
    if [ -d "$LANG_CLI_DIR" ]; then
        echo "  rm -rf $LANG_CLI_DIR"
    fi
fi

echo ""
echo "感谢使用 lang-cli！"
echo "Thank you for using lang-cli!"