package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/daiweiwei/lang-cli/internal/lang"
)

// LanguageMenuItem 语言菜单项
type LanguageMenuItem struct {
	language    string
	description string
	isCurrent   bool
}

// 实现list.Item接口
func (i LanguageMenuItem) Title() string {
	title := i.language
	if i.isCurrent {
		title = "✔ " + title
	}
	return title
}

func (i LanguageMenuItem) Description() string { return i.description }
func (i LanguageMenuItem) FilterValue() string { return i.language }

// LanguageMenu 语言菜单模型
type LanguageMenu struct {
	list     list.Model
	quitting bool
}

// 创建新的语言菜单
func NewLanguageMenu() *LanguageMenu {
	// 获取语言列表
	languages := lang.ListLanguages()
	currentLang := lang.GetCurrentLanguage()

	// 创建菜单项
	items := []list.Item{}
	for _, language := range languages {
		items = append(items, LanguageMenuItem{
			language:    language,
			description: "切换到" + language + "语言",
			isCurrent:   language == currentLang,
		})
	}

	// 添加返回选项
	items = append(items, MenuItem{
		title:       "返回主菜单",
		description: "返回到主菜单",
		action:      func() (tea.Model, error) { return NewMainMenu(), nil },
	})

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "语言管理"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	return &LanguageMenu{
		list: l,
	}
}

// Init 初始化模型
func (m LanguageMenu) Init() tea.Cmd {
	return nil
}

// Update 更新模型
func (m LanguageMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		case "esc":
			mainMenu := NewMainMenu()
			width, height := m.list.Width(), m.list.Height()+4
			if width > 0 && height > 4 {
				updatedModel, _ := mainMenu.Update(tea.WindowSizeMsg{Width: width, Height: height})
				return updatedModel, nil
			}
			return mainMenu, nil

		case "enter":
			switch i := m.list.SelectedItem().(type) {
			case MenuItem:
				if i.action != nil {
					newModel, err := i.action()
					if err != nil {
						// 处理错误
						return m, nil
					}
					// 如果返回的是主菜单，需要传递当前窗口大小
					if mainMenu, ok := newModel.(*MainMenu); ok {
						// 获取当前窗口大小并设置
						width, height := m.list.Width(), m.list.Height()+4
						if width > 0 && height > 4 {
							updatedModel, _ := mainMenu.Update(tea.WindowSizeMsg{Width: width, Height: height})
							return updatedModel, nil
						}
						return mainMenu, nil
					}
					return newModel, nil
				}
			case LanguageMenuItem:
				// 切换语言
				lang.SwitchLanguage(i.language)
				// 刷新菜单
				newModel := NewLanguageMenu()
				// 传递当前窗口大小
				width, height := m.list.Width(), m.list.Height()+4
				if width > 0 && height > 4 {
					updatedModel, _ := newModel.Update(tea.WindowSizeMsg{Width: width, Height: height})
					return updatedModel, nil
				}
				return newModel, nil
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View 渲染视图
func (m LanguageMenu) View() string {
	if m.quitting {
		return "再见！"
	}

	return m.list.View()
}
