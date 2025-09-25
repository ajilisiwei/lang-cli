package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/daiweiwei/lang-cli/internal/practice"
)

// PracticeMenu 练习菜单模型
type PracticeMenu struct {
	list     list.Model
	quitting bool
}

// 创建新的练习菜单
func NewPracticeMenu() *PracticeMenu {
	items := []list.Item{
		MenuItem{
			title:       "单词练习",
			description: "进行单词打字练习",
			action:      func() (tea.Model, error) { return NewResourceSelectionMenu(practice.Words), nil },
		},
		MenuItem{
			title:       "短语练习",
			description: "进行短语打字练习",
			action:      func() (tea.Model, error) { return NewResourceSelectionMenu(practice.Phrases), nil },
		},
		MenuItem{
			title:       "句子练习",
			description: "进行句子打字练习",
			action:      func() (tea.Model, error) { return NewResourceSelectionMenu(practice.Sentences), nil },
		},
		MenuItem{
			title:       "文章练习",
			description: "进行文章打字练习",
			action:      func() (tea.Model, error) { return NewResourceSelectionMenu(practice.Articles), nil },
		},
		MenuItem{
			title:       "返回主菜单",
			description: "返回到主菜单",
			action:      func() (tea.Model, error) { return NewMainMenu(), nil },
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "练习菜单"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	return &PracticeMenu{
		list: l,
	}
}

// Init 初始化模型
func (m PracticeMenu) Init() tea.Cmd {
	return nil
}

// Update 更新模型
func (m PracticeMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			i, ok := m.list.SelectedItem().(MenuItem)
			if ok && i.action != nil {
				newModel, err := i.action()
				if err != nil {
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

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View 渲染视图
func (m PracticeMenu) View() string {
	if m.quitting {
		return "再见！"
	}

	return m.list.View()
}

// ResourceSelectionMenu 资源选择菜单
type ResourceSelectionMenu struct {
	list         list.Model
	resourceType string
	quitting     bool
}

// 创建新的资源选择菜单
func NewResourceSelectionMenu(resourceType string) *ResourceSelectionMenu {
	// 获取资源文件列表
	files, err := practice.GetResourceFiles(resourceType)
	if err != nil {
		files = []string{}
	}

	// 创建菜单项
	items := []list.Item{}
	for _, file := range files {
		fileName := file // 创建副本以避免闭包问题
		items = append(items, MenuItem{
			title:       fileName,
			description: "选择" + fileName + "进行练习",
			action: func() (tea.Model, error) {
				return NewPracticeSession(resourceType, fileName), nil
			},
		})
	}

	// 添加返回选项
	items = append(items, MenuItem{
		title:       "返回练习菜单",
		description: "返回到练习菜单",
		action:      func() (tea.Model, error) { return NewPracticeMenu(), nil },
	})

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = getResourceTypeTitle(resourceType) + "选择"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	return &ResourceSelectionMenu{
		list:         l,
		resourceType: resourceType,
	}
}

// 获取资源类型标题
func getResourceTypeTitle(resourceType string) string {
	switch resourceType {
	case practice.Words:
		return "单词"
	case practice.Phrases:
		return "短语"
	case practice.Sentences:
		return "句子"
	case practice.Articles:
		return "文章"
	default:
		return "资源"
	}
}

// Init 初始化模型
func (m ResourceSelectionMenu) Init() tea.Cmd {
	return nil
}

// Update 更新模型
func (m ResourceSelectionMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			practiceMenu := NewPracticeMenu()
			width, height := m.list.Width(), m.list.Height()+4
			if width > 0 && height > 4 {
				updatedModel, _ := practiceMenu.Update(tea.WindowSizeMsg{Width: width, Height: height})
				return updatedModel, nil
			}
			return practiceMenu, nil

		case "enter":
			i, ok := m.list.SelectedItem().(MenuItem)
			if ok && i.action != nil {
				newModel, err := i.action()
				if err != nil {
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

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View 渲染视图
func (m ResourceSelectionMenu) View() string {
	if m.quitting {
		return "再见！"
	}

	return m.list.View()
}
