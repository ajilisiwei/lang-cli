package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/daiweiwei/lang-cli/internal/manage"
	"github.com/daiweiwei/lang-cli/internal/practice"
)

// ManageMenu 资源管理菜单模型
type ManageMenu struct {
	list     list.Model
	quitting bool
}

// 创建新的资源管理菜单
func NewManageMenu() *ManageMenu {
	items := []list.Item{
		MenuItem{
			title:       "删除资源",
			description: "删除练习资源",
			action:      func() (tea.Model, error) { return NewResourceTypeMenu("delete"), nil },
		},
		MenuItem{
			title:       "导入资源",
			description: "导入练习资源",
			action:      func() (tea.Model, error) { return NewResourceTypeMenu("import"), nil },
		},
		MenuItem{
			title:       "返回主菜单",
			description: "返回到主菜单",
			action:      func() (tea.Model, error) { return NewMainMenu(), nil },
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "资源管理"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	return &ManageMenu{
		list: l,
	}
}

// Init 初始化模型
func (m ManageMenu) Init() tea.Cmd {
	return nil
}

// Update 更新模型
func (m ManageMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			// 返回到主菜单
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
func (m ManageMenu) View() string {
	if m.quitting {
		return "再见！"
	}

	return m.list.View()
}

// ResourceTypeMenu 资源类型菜单模型
type ResourceTypeMenu struct {
	list     list.Model
	action   string // "delete" 或 "import"
	quitting bool
}

// 创建新的资源类型菜单
func NewResourceTypeMenu(action string) *ResourceTypeMenu {
	items := []list.Item{
		MenuItem{
			title:       "单词",
			description: getActionTitle(action) + "单词资源",
			action: func() (tea.Model, error) {
				if action == "delete" {
					return NewManageResourceMenu(practice.Words, action), nil
				} else {
					return NewImportView(practice.Words), nil
				}
			},
		},
		MenuItem{
			title:       "短语",
			description: getActionTitle(action) + "短语资源",
			action: func() (tea.Model, error) {
				if action == "delete" {
					return NewManageResourceMenu(practice.Phrases, action), nil
				} else {
					return NewImportView(practice.Phrases), nil
				}
			},
		},
		MenuItem{
			title:       "句子",
			description: getActionTitle(action) + "句子资源",
			action: func() (tea.Model, error) {
				if action == "delete" {
					return NewManageResourceMenu(practice.Sentences, action), nil
				} else {
					return NewImportView(practice.Sentences), nil
				}
			},
		},
		MenuItem{
			title:       "文章",
			description: getActionTitle(action) + "文章资源",
			action: func() (tea.Model, error) {
				if action == "delete" {
					return NewManageResourceMenu(practice.Articles, action), nil
				} else {
					return NewImportView(practice.Articles), nil
				}
			},
		},
		MenuItem{
			title:       "返回资源管理菜单",
			description: "返回到资源管理菜单",
			action:      func() (tea.Model, error) { return NewManageMenu(), nil },
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = getActionTitle(action) + "资源"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	return &ResourceTypeMenu{
		list:   l,
		action: action,
	}
}

// 获取操作标题
func getActionTitle(action string) string {
	switch action {
	case "delete":
		return "删除"
	case "import":
		return "导入"
	default:
		return "管理"
	}
}

// Init 初始化模型
func (m ResourceTypeMenu) Init() tea.Cmd {
	return nil
}

// Update 更新模型
func (m ResourceTypeMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			// 返回到管理菜单
			manageMenu := NewManageMenu()
			width, height := m.list.Width(), m.list.Height()+4
			if width > 0 && height > 4 {
				updatedModel, _ := manageMenu.Update(tea.WindowSizeMsg{Width: width, Height: height})
				return updatedModel, nil
			}
			return manageMenu, nil

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
func (m ResourceTypeMenu) View() string {
	if m.quitting {
		return "再见！"
	}

	return m.list.View()
}

// ManageResourceMenu 资源管理菜单模型
type ManageResourceMenu struct {
	list         list.Model
	resourceType string
	action       string
	quitting     bool
}

// 创建新的资源管理菜单
func NewManageResourceMenu(resourceType, action string) *ManageResourceMenu {
	// 获取资源文件列表
	files, err := manage.ListResourceFiles(resourceType)
	if err != nil {
		files = []string{}
	}

	// 创建菜单项
	items := []list.Item{}
	for _, file := range files {
		fileName := file // 创建副本以避免闭包问题
		items = append(items, MenuItem{
			title:       fileName,
			description: getActionTitle(action) + fileName,
			action: func() (tea.Model, error) {
				if action == "delete" {
					// 显示删除确认对话框
					return NewDeleteConfirmView(resourceType, fileName), nil
				}
				return nil, nil
			},
		})
	}

	// 添加返回选项
	items = append(items, MenuItem{
		title:       "返回资源类型菜单",
		description: "返回到资源类型菜单",
		action:      func() (tea.Model, error) { return NewResourceTypeMenu(action), nil },
	})

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = getActionTitle(action) + getResourceTypeTitle(resourceType)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	return &ManageResourceMenu{
		list:         l,
		resourceType: resourceType,
		action:       action,
	}
}

// Init 初始化模型
func (m ManageResourceMenu) Init() tea.Cmd {
	return nil
}

// Update 更新模型
func (m ManageResourceMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			// 返回到资源类型菜单
			resourceTypeMenu := NewResourceTypeMenu(m.action)
			width, height := m.list.Width(), m.list.Height()+4
			if width > 0 && height > 4 {
				updatedModel, _ := resourceTypeMenu.Update(tea.WindowSizeMsg{Width: width, Height: height})
				return updatedModel, nil
			}
			return resourceTypeMenu, nil

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
func (m ManageResourceMenu) View() string {
	if m.quitting {
		return "再见！"
	}

	return m.list.View()
}

// ImportView 导入视图模型
type ImportView struct {
	resourceType string
	textInput    textinput.Model
	message      string
	width        int
	height       int
	quitting     bool
}

// 创建新的导入视图
func NewImportView(resourceType string) *ImportView {
	// 创建文本输入
	ti := textinput.New()
	ti.Placeholder = "输入文件路径..."
	ti.Focus()
	ti.Width = 40

	return &ImportView{
		resourceType: resourceType,
		textInput:    ti,
		message:      "",
	}
}

// Init 初始化模型
func (m ImportView) Init() tea.Cmd {
	return textinput.Blink
}

// Update 更新模型
func (m ImportView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textInput.Width = msg.Width - 20
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "esc":
			// 返回到资源类型菜单
			newModel := NewResourceTypeMenu("import")
			// 传递当前窗口大小给新模型
			if m.width > 0 && m.height > 0 {
				updatedModel, _ := newModel.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				return updatedModel, nil
			}
			return newModel, nil

		case "enter":
			// 获取文件路径
			filePath := m.textInput.Value()
			if filePath == "" {
				m.message = "请输入文件路径"
				return m, nil
			}

			// 导入资源 - 使用不需要用户确认的版本
			err := manage.ImportResourceForTest(m.resourceType, filePath)
			if err != nil {
				m.message = "导入失败: " + err.Error()
				return m, nil
			} else {
				// 导入成功，自动返回到资源类型菜单
				newModel := NewResourceTypeMenu("import")
				// 传递当前窗口大小给新模型
				if m.width > 0 && m.height > 0 {
					updatedModel, _ := newModel.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
					return updatedModel, nil
				}
				return newModel, nil
			}
		}
	}

	// 处理文本输入
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View 渲染视图
func (m ImportView) View() string {
	if m.quitting {
		return "导入已取消！"
	}

	var s strings.Builder

	// 标题
	title := "导入" + getResourceTypeTitle(m.resourceType)
	s.WriteString(RenderTitle(title) + "\n\n")

	// 显示输入框
	s.WriteString(RenderHighlight("请输入文件路径:") + "\n")
	s.WriteString(m.textInput.View() + "\n\n")

	// 显示消息
	if m.message != "" {
		if strings.HasPrefix(m.message, "导入失败") {
			s.WriteString(RenderError(m.message) + "\n\n")
		} else {
			s.WriteString(RenderSuccess(m.message) + "\n\n")
		}
	}

	// 显示提示
	s.WriteString(RenderText("按 Enter 导入，按 Esc 返回") + "\n")

	return s.String()
}

// DeleteConfirmView 删除确认视图模型
type DeleteConfirmView struct {
	resourceType string
	fileName     string
	message      string
	width        int
	height       int
	quitting     bool
}

// 创建新的删除确认视图
func NewDeleteConfirmView(resourceType, fileName string) *DeleteConfirmView {
	return &DeleteConfirmView{
		resourceType: resourceType,
		fileName:     fileName,
	}
}

// Init 初始化模型
func (m DeleteConfirmView) Init() tea.Cmd {
	return nil
}

// Update 更新模型
func (m DeleteConfirmView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "y", "Y":
			// 确认删除
			err := manage.DeleteResourceWithoutConfirm(m.resourceType, m.fileName)
			if err != nil {
				m.message = "删除失败: " + err.Error()
				return m, nil
			}
			// 删除成功，返回资源列表
			newModel := NewManageResourceMenu(m.resourceType, "delete")
			// 传递当前窗口大小
			if m.width > 0 && m.height > 0 {
				updatedModel, _ := newModel.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				return updatedModel, nil
			}
			return newModel, nil

		case "n", "N", "esc":
			// 取消删除，返回资源列表
			newModel := NewManageResourceMenu(m.resourceType, "delete")
			// 传递当前窗口大小
			if m.width > 0 && m.height > 0 {
				updatedModel, _ := newModel.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				return updatedModel, nil
			}
			return newModel, nil
		}
	}

	return m, nil
}

// View 渲染视图
func (m DeleteConfirmView) View() string {
	if m.quitting {
		return "再见！"
	}

	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(TitleStyle.Render("确认删除"))
	s.WriteString("\n\n")
	s.WriteString("确认删除以下资源吗？\n\n")
	s.WriteString("资源类型: " + getResourceTypeTitle(m.resourceType) + "\n")
	s.WriteString("文件名: " + m.fileName + "\n\n")
	s.WriteString("按 Y 确认删除，按 N 或 Esc 取消\n")

	if m.message != "" {
		s.WriteString("\n")
		s.WriteString(m.message)
	}

	return s.String()
}
