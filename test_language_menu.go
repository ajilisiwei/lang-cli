package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/daiweiwei/lang-cli/internal/config"
	"github.com/daiweiwei/lang-cli/internal/ui"
)

func main() {
	// 加载配置文件
	if err := config.LoadConfig(); err != nil {
		fmt.Println("警告: 加载配置文件失败:", err)
		fmt.Println("将使用默认配置。")
	}

	// 直接启动语言管理菜单
	fmt.Println("启动语言管理菜单测试...")
	p := tea.NewProgram(ui.NewLanguageMenu(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("启动语言管理菜单失败: %v\n", err)
		os.Exit(1)
	}
}