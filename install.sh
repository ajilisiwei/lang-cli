#!/bin/bash

# 安装脚本 - 多语言打字学习终端工具
# Install script for Multi-language Typing Learning Terminal Tool

set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="lang-cli"
LANG_CLI_DIR="$HOME/.lang-cli"

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

echo "可执行文件安装完成！"
echo ""

# 初始化用户配置和资源文件
echo "正在初始化用户配置和资源文件..."

# 创建用户主目录下的 .lang-cli 目录
echo "创建目录: $LANG_CLI_DIR"
mkdir -p "$LANG_CLI_DIR"

# 复制配置文件（强制覆盖以确保配置文件是最新的）
if [ -f "config/config.yaml" ]; then
    echo "复制配置文件..."
    cp "config/config.yaml" "$LANG_CLI_DIR/"
    echo "配置文件已更新到最新版本"
else
    echo "警告: 配置文件 config/config.yaml 不存在"
fi

# 复制资源文件
if [ -d "resources" ]; then
    if [ ! -d "$LANG_CLI_DIR/resources" ]; then
        echo "复制资源文件..."
        cp -r "resources" "$LANG_CLI_DIR/"
    else
        echo "资源目录已存在，正在合并新资源..."
        # 遍历每个语言目录
        for lang_dir in resources/*/; do
            if [ -d "$lang_dir" ]; then
                lang_name=$(basename "$lang_dir")
                target_lang_dir="$LANG_CLI_DIR/resources/$lang_name"
                mkdir -p "$target_lang_dir"
                
                # 遍历每个资源类型目录
                for type_dir in "$lang_dir"*/; do
                    if [ -d "$type_dir" ]; then
                        type_name=$(basename "$type_dir")
                        target_type_dir="$target_lang_dir/$type_name"
                        mkdir -p "$target_type_dir"
                        
                        # 复制文件，但不覆盖已存在的文件
                        for file in "$type_dir"*; do
                            if [ -f "$file" ]; then
                                filename=$(basename "$file")
                                target_file="$target_type_dir/$filename"
                                if [ ! -f "$target_file" ]; then
                                    cp "$file" "$target_file"
                                    echo "  添加新资源: $lang_name/$type_name/$filename"
                                else
                                    echo "  跳过已存在的资源: $lang_name/$type_name/$filename"
                                fi
                            fi
                        done
                    fi
                done
            fi
        done
    fi
else
    echo "警告: 资源目录 resources 不存在"
fi

echo ""
echo "安装和初始化完成！"
echo "Installation and initialization completed!"
echo ""
echo "配置和资源文件已复制到: $LANG_CLI_DIR"
echo "Configuration and resource files copied to: $LANG_CLI_DIR"
echo ""
echo "现在您可以在任何地方使用 '$BINARY_NAME' 命令"
echo "You can now use '$BINARY_NAME' command from anywhere"
echo ""
echo "使用方法 / Usage:"
echo "  $BINARY_NAME"
echo ""
echo "卸载 / Uninstall:"
echo "  sudo rm $INSTALL_DIR/$BINARY_NAME"
echo "  rm -rf $LANG_CLI_DIR"