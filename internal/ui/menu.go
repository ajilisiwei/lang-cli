package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// 菜单项
type MenuItem struct {
	title       string
	description string
	action      func() (tea.Model, error)
}

// 实现list.Item接口
func (i MenuItem) Title() string       { return i.title }
func (i MenuItem) Description() string { return i.description }
func (i MenuItem) FilterValue() string { return i.title }

// MainMenu 主菜单模型
type MainMenu struct {
	list     list.Model
	choice   *MenuItem
	quitting bool
}

// 创建新的主菜单
func NewMainMenu() *MainMenu {
	items := []list.Item{
		MenuItem{
			title:       "语言管理",
			description: "列出支持的语言和切换当前练习的语言",
			action:      func() (tea.Model, error) { return NewLanguageMenu(), nil },
		},
		MenuItem{
			title:       "练习",
			description: "进行单词、短语、句子、文章等的打字练习",
			action:      func() (tea.Model, error) { return NewPracticeMenu(), nil },
		},
		MenuItem{
			title:       "资源管理",
			description: "管理练习资源，包括删除和导入资源",
			action:      func() (tea.Model, error) { return NewManageMenu(), nil },
		},
		MenuItem{
			title:       "退出",
			description: "退出程序",
			action:      func() (tea.Model, error) { return nil, fmt.Errorf("quit") },
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "多语言打字学习终端工具"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	return &MainMenu{
		list: l,
	}
}

// Init 初始化模型
func (m MainMenu) Init() tea.Cmd {
	return nil
}

// Update 更新模型
func (m MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4) // 减去标题和状态栏的高度
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(MenuItem)
			if ok {
				m.choice = &i
				if i.action != nil {
					newModel, err := i.action()
					if err != nil {
						if err.Error() == "quit" {
							return m, tea.Quit
						}
						// 处理错误
						return m, nil
					}
					// 传递当前窗口大小给新模型
					width, height := m.list.Width(), m.list.Height()+4
					if width > 0 && height > 4 {
						updatedModel, _ := newModel.Update(tea.WindowSizeMsg{Width: width, Height: height})
						return updatedModel, nil
					}
					return newModel, nil
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View 渲染视图
func (m MainMenu) View() string {
	if m.quitting {
		return "再见！"
	}

	return m.list.View()
}
