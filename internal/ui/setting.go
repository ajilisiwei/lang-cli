package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/daiweiwei/lang-cli/internal/config"
)

// SettingMenuItem 设置菜单项
type SettingMenuItem struct {
	title       string
	description string
	action      func() (tea.Model, error)
}

// 实现list.Item接口
func (i SettingMenuItem) Title() string       { return i.title }
func (i SettingMenuItem) Description() string { return i.description }
func (i SettingMenuItem) FilterValue() string { return i.title }

// SettingMenu 设置菜单模型
type SettingMenu struct {
	list     list.Model
	quitting bool
}

// 创建新的设置菜单
func NewSettingMenu() *SettingMenu {
	// 创建菜单项
	items := []list.Item{
		SettingMenuItem{
			title:       "匹配模式设置",
			description: "设置正确性匹配模式（完全匹配/单词匹配）",
			action: func() (tea.Model, error) {
				return NewMatchModeMenu(), nil
			},
		},
		SettingMenuItem{
			title:       "顺序设置",
			description: "设置练习资源的出现顺序（随机/顺序）",
			action: func() (tea.Model, error) {
				return NewOrderMenu(), nil
			},
		},
		MenuItem{
			title:       "返回主菜单",
			description: "返回到主菜单",
			action:      func() (tea.Model, error) { return NewMainMenu(), nil },
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "设置"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	return &SettingMenu{
		list: l,
	}
}

// Init 初始化模型
func (m SettingMenu) Init() tea.Cmd {
	return nil
}

// Update 更新模型
func (m SettingMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			switch i := m.list.SelectedItem().(type) {
			case SettingMenuItem:
				if i.action != nil {
					newModel, err := i.action()
					if err != nil {
						// 处理错误
						return m, nil
					}
					// 传递当前窗口大小
					width, height := m.list.Width(), m.list.Height()+4
					if width > 0 && height > 4 {
						updatedModel, _ := newModel.Update(tea.WindowSizeMsg{Width: width, Height: height})
						return updatedModel, nil
					}
					return newModel, nil
				}
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
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View 渲染视图
func (m SettingMenu) View() string {
	if m.quitting {
		return "再见！"
	}

	return m.list.View()
}

// MatchModeMenu 匹配模式菜单模型
type MatchModeMenu struct {
	list     list.Model
	quitting bool
}

// MatchModeMenuItem 匹配模式菜单项
type MatchModeMenuItem struct {
	mode        string
	title       string
	description string
	isCurrent   bool
}

// 实现list.Item接口
func (i MatchModeMenuItem) Title() string {
	title := i.title
	if i.isCurrent {
		title = "· " + title
	}
	return title
}
func (i MatchModeMenuItem) Description() string { return i.description }
func (i MatchModeMenuItem) FilterValue() string { return i.title }

// 创建新的匹配模式菜单
func NewMatchModeMenu() *MatchModeMenu {
	currentMode := config.AppConfig.CorrectnessMatchMode
	if currentMode == "" {
		currentMode = "exact_match" // 默认值
	}

	// 创建菜单项
	items := []list.Item{
		MatchModeMenuItem{
			mode:        "exact_match",
			title:       "完全匹配",
			description: "用户输入必须完全匹配才能被认为是正确的",
			isCurrent:   currentMode == "exact_match",
		},
		MatchModeMenuItem{
			mode:        "word_match",
			title:       "单词匹配",
			description: "只要包含单词就被认为是正确的，忽略大小写和标点符号",
			isCurrent:   currentMode == "word_match",
		},
		MenuItem{
			title:       "返回设置菜单",
			description: "返回到设置菜单",
			action:      func() (tea.Model, error) { return NewSettingMenu(), nil },
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "匹配模式设置"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	return &MatchModeMenu{
		list: l,
	}
}

// Init 初始化模型
func (m MatchModeMenu) Init() tea.Cmd {
	return nil
}

// Update 更新模型
func (m MatchModeMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			switch i := m.list.SelectedItem().(type) {
			case MatchModeMenuItem:
				// 设置匹配模式
				config.AppConfig.CorrectnessMatchMode = i.mode
				config.SaveConfig()
				// 刷新菜单
				newModel := NewMatchModeMenu()
				// 传递当前窗口大小
				width, height := m.list.Width(), m.list.Height()+4
				if width > 0 && height > 4 {
					updatedModel, _ := newModel.Update(tea.WindowSizeMsg{Width: width, Height: height})
					return updatedModel, nil
				}
				return newModel, nil
			case MenuItem:
				if i.action != nil {
					newModel, err := i.action()
					if err != nil {
						return m, nil
					}
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
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View 渲染视图
func (m MatchModeMenu) View() string {
	if m.quitting {
		return "再见！"
	}

	return m.list.View()
}

// OrderMenu 顺序菜单模型
type OrderMenu struct {
	list     list.Model
	quitting bool
}

// OrderMenuItem 顺序菜单项
type OrderMenuItem struct {
	order       string
	title       string
	description string
	isCurrent   bool
}

// 实现list.Item接口
func (i OrderMenuItem) Title() string {
	title := i.title
	if i.isCurrent {
		title = "· " + title
	}
	return title
}
func (i OrderMenuItem) Description() string { return i.description }
func (i OrderMenuItem) FilterValue() string { return i.title }

// 创建新的顺序菜单
func NewOrderMenu() *OrderMenu {
	currentOrder := config.AppConfig.NextOneOrder
	if currentOrder == "" {
		currentOrder = "random" // 默认值
	}

	// 创建菜单项
	items := []list.Item{
		OrderMenuItem{
			order:       "random",
			title:       "随机顺序",
			description: "练习时资源的出现顺序是随机的",
			isCurrent:   currentOrder == "random",
		},
		OrderMenuItem{
			order:       "sequential",
			title:       "顺序出现",
			description: "练习时资源的出现顺序是顺序的",
			isCurrent:   currentOrder == "sequential",
		},
		MenuItem{
			title:       "返回设置菜单",
			description: "返回到设置菜单",
			action:      func() (tea.Model, error) { return NewSettingMenu(), nil },
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "顺序设置"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	return &OrderMenu{
		list: l,
	}
}

// Init 初始化模型
func (m OrderMenu) Init() tea.Cmd {
	return nil
}

// Update 更新模型
func (m OrderMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			switch i := m.list.SelectedItem().(type) {
			case OrderMenuItem:
				// 设置顺序
				config.AppConfig.NextOneOrder = i.order
				config.SaveConfig()
				// 刷新菜单
				newModel := NewOrderMenu()
				// 传递当前窗口大小
				width, height := m.list.Width(), m.list.Height()+4
				if width > 0 && height > 4 {
					updatedModel, _ := newModel.Update(tea.WindowSizeMsg{Width: width, Height: height})
					return updatedModel, nil
				}
				return newModel, nil
			case MenuItem:
				if i.action != nil {
					newModel, err := i.action()
					if err != nil {
						return m, nil
					}
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
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View 渲染视图
func (m OrderMenu) View() string {
	if m.quitting {
		return "再见！"
	}

	return m.list.View()
}