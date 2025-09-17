package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/daiweiwei/lang-cli/internal/practice"
)

// PracticeSession 练习会话模型
type PracticeSession struct {
	resourceType string
	fileName     string
	items        []string        // 练习项目
	currentIndex int             // 当前项目索引
	textInput    textinput.Model // 文本输入
	progress     progress.Model  // 进度条
	width        int             // 窗口宽度
	height       int             // 窗口高度
	startTime    time.Time       // 开始时间
	endTime      time.Time       // 结束时间
	correct      int             // 正确数量
	incorrect    int             // 错误数量
	state        string          // 状态："practicing", "finished"
	result       string          // 结果信息
	quitting     bool
}

// 创建新的练习会话
func NewPracticeSession(resourceType, fileName string) *PracticeSession {
	// 读取资源文件
	items, err := practice.ReadResourceFile(resourceType, fileName)
	if err != nil || len(items) == 0 {
		// 处理错误或空文件
		items = []string{"没有可用的练习内容"}
	}

	// 创建文本输入
	ti := textinput.New()
	ti.Placeholder = "输入这里..."
	ti.Focus()
	ti.Width = 40

	// 创建进度条
	p := progress.New(progress.WithDefaultGradient())

	return &PracticeSession{
		resourceType: resourceType,
		fileName:     fileName,
		items:        items,
		currentIndex: 0,
		textInput:    ti,
		progress:     p,
		startTime:    time.Now(),
		state:        "practicing",
	}
}

// Init 初始化模型
func (m PracticeSession) Init() tea.Cmd {
	return textinput.Blink
}

// Update 更新模型
func (m PracticeSession) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 20
		m.textInput.Width = msg.Width - 20
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "esc":
			if m.state == "finished" {
				// 返回练习菜单并传递窗口大小
				practiceMenu := NewPracticeMenu()
				if m.width > 0 && m.height > 4 {
					updatedModel, _ := practiceMenu.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
					return updatedModel, nil
				}
				return practiceMenu, nil
			}
			// 确认是否要退出练习
			return m, nil

		case "enter":
			if m.state == "finished" {
				// 返回练习菜单并传递窗口大小
				practiceMenu := NewPracticeMenu()
				if m.width > 0 && m.height > 4 {
					updatedModel, _ := practiceMenu.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
					return updatedModel, nil
				}
				return practiceMenu, nil
			}

			// 检查输入是否正确
			userInput := m.textInput.Value()
			currentItem := m.getCurrentItem()
			expectedInput := m.getExpectedInput(currentItem)

			if userInput == expectedInput {
				m.correct++
			} else {
				m.incorrect++
			}

			// 清空输入
			m.textInput.SetValue("")

			// 移动到下一项
			m.currentIndex++
			if m.currentIndex >= len(m.items) {
				// 练习结束
				m.state = "finished"
				m.endTime = time.Now()
				m.result = m.calculateResult()
			}

			return m, nil
		}

		// 如果练习已结束，不处理其他按键
		if m.state == "finished" {
			return m, nil
		}
	}

	// 处理文本输入
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View 渲染视图
func (m PracticeSession) View() string {
	if m.quitting {
		return "练习已中断！"
	}

	var s strings.Builder

	// 标题
	title := fmt.Sprintf("%s练习 - %s", getResourceTypeTitle(m.resourceType), m.fileName)
	s.WriteString(RenderTitle(title) + "\n\n")

	if m.state == "practicing" {
		// 显示当前进度
		progressValue := float64(m.currentIndex) / float64(len(m.items))
		progressText := fmt.Sprintf("进度: %d/%d", m.currentIndex+1, len(m.items))
		s.WriteString(RenderText(progressText) + "\n")
		s.WriteString(m.progress.ViewAs(progressValue) + "\n\n")

		// 显示当前项目
		currentItem := m.getCurrentItem()
		s.WriteString(RenderHighlight("当前项目:") + "\n")
		s.WriteString(RenderText(currentItem) + "\n\n")

		// 显示输入框
		s.WriteString(RenderHighlight("请输入:") + "\n")
		s.WriteString(m.textInput.View() + "\n\n")

		// 显示提示
		s.WriteString(RenderText("按 Enter 提交，按 Esc 退出练习") + "\n")
	} else if m.state == "finished" {
		// 显示练习结果
		s.WriteString(RenderSuccess("练习完成！") + "\n\n")
		s.WriteString(RenderText(m.result) + "\n\n")
		s.WriteString(RenderText("按 Enter 或 Esc 返回练习菜单") + "\n")
	}

	return s.String()
}

// 获取当前项目
func (m PracticeSession) getCurrentItem() string {
	if m.currentIndex < 0 || m.currentIndex >= len(m.items) {
		return ""
	}
	return m.items[m.currentIndex]
}

// 获取期望输入
func (m PracticeSession) getExpectedInput(item string) string {
	// 对于单词和短语，可能有翻译，需要分离
	if m.resourceType == practice.Words || m.resourceType == practice.Phrases {
		parts := strings.Split(item, practice.Separator)
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	return item
}

// 计算练习结果
func (m PracticeSession) calculateResult() string {
	// 计算练习时间
	duration := m.endTime.Sub(m.startTime)
	durationStr := fmt.Sprintf("%d分%d秒", int(duration.Minutes()), int(duration.Seconds())%60)

	// 计算正确率
	total := m.correct + m.incorrect
	accuracy := 0.0
	if total > 0 {
		accuracy = float64(m.correct) / float64(total) * 100
	}

	// 计算速度（每分钟字符数）
	totalChars := 0
	for _, item := range m.items {
		expectedInput := m.getExpectedInput(item)
		totalChars += len(expectedInput)
	}

	cpm := 0.0
	if duration.Minutes() > 0 {
		cpm = float64(totalChars) / duration.Minutes()
	}

	return fmt.Sprintf(
		"练习时间: %s\n正确数量: %d\n错误数量: %d\n正确率: %.1f%%\n打字速度: %.1f CPM (每分钟字符数)",
		durationStr, m.correct, m.incorrect, accuracy, cpm,
	)
}
