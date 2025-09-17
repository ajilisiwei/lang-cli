package ui

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/daiweiwei/lang-cli/internal/config"
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
	// v0.2 新增：错误状态跟踪
	lastInputWrong bool   // 上次输入是否错误
	wrongInput     string // 错误的输入内容
	expectedText   string // 期望的正确文本
	// v0.2 新增：随机顺序支持
	practiceOrder []int // 练习顺序索引列表
	completedCount int   // 已完成的项目数量
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

	// 初始化练习顺序
	practiceOrder := make([]int, len(items))
	for i := range practiceOrder {
		practiceOrder[i] = i
	}

	// 根据配置决定是否打乱顺序
	orderMode := config.AppConfig.NextOneOrder
	if orderMode == "" {
		orderMode = "random" // 默认随机
	}
	if orderMode == "random" {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(practiceOrder), func(i, j int) {
			practiceOrder[i], practiceOrder[j] = practiceOrder[j], practiceOrder[i]
		})
	}

	return &PracticeSession{
		resourceType:  resourceType,
		fileName:      fileName,
		items:         items,
		currentIndex:  0,
		textInput:     ti,
		progress:      p,
		startTime:     time.Now(),
		state:         "practicing",
		practiceOrder: practiceOrder,
		completedCount: 0,
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
			// 无论在什么状态下，ESC键都应该退出练习并返回练习菜单
			practiceMenu := NewPracticeMenu()
			if m.width > 0 && m.height > 4 {
				updatedModel, _ := practiceMenu.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				return updatedModel, nil
			}
			return practiceMenu, nil

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

			if m.isInputCorrect(userInput, expectedInput) {
				// 输入正确
				m.correct++
				m.lastInputWrong = false
				m.wrongInput = ""
				m.expectedText = ""

				// 清空输入
				m.textInput.SetValue("")

				// 移动到下一项
				m.completedCount++
				if m.completedCount >= len(m.items) {
					// 练习结束
					m.state = "finished"
					m.endTime = time.Now()
					m.result = m.calculateResult()
				}
			} else {
				// 输入错误
				m.incorrect++
				m.lastInputWrong = true
				m.wrongInput = userInput
				m.expectedText = expectedInput

				// 清空输入，让用户重新输入
				m.textInput.SetValue("")
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
		progressValue := float64(m.completedCount) / float64(len(m.items))
		progressText := fmt.Sprintf("进度: %d/%d", m.completedCount+1, len(m.items))
		s.WriteString(RenderText(progressText) + "\n")
		s.WriteString(m.progress.ViewAs(progressValue) + "\n\n")

		// 显示当前项目
		currentItem := m.getCurrentItem()
		s.WriteString(RenderHighlight("当前项目:") + "\n")
		wrappedText := m.wrapText(currentItem, m.width-4) // 减去边距
		s.WriteString(RenderText(wrappedText) + "\n\n")

		// 显示输入框
		s.WriteString(RenderHighlight("请输入:") + "\n")
		s.WriteString(m.textInput.View() + "\n\n")

		// 显示错误信息（如果有）
		if m.lastInputWrong {
			s.WriteString(RenderError("❌ 输入错误！") + "\n")
			s.WriteString(m.renderWordLevelError() + "\n\n")
		}

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
	if m.completedCount < len(m.practiceOrder) {
		actualIndex := m.practiceOrder[m.completedCount]
		if actualIndex < len(m.items) {
			return m.items[actualIndex]
		}
	}
	return ""
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

// 检查输入是否正确
func (m PracticeSession) isInputCorrect(userInput, expectedInput string) bool {
	// 获取匹配模式
	matchMode := config.AppConfig.CorrectnessMatchMode
	if matchMode == "" {
		matchMode = "exact_match" // 默认值
	}

	switch matchMode {
	case "exact_match":
		// 完全匹配
		return userInput == expectedInput
	case "word_match":
		// 单词匹配，忽略大小写和标点符号
		return m.normalizeForWordMatch(userInput) == m.normalizeForWordMatch(expectedInput)
	default:
		// 默认使用完全匹配
		return userInput == expectedInput
	}
}

// 标准化文本用于单词匹配
func (m PracticeSession) normalizeForWordMatch(text string) string {
	// 转换为小写
	text = strings.ToLower(text)
	// 移除标点符号，只保留字母、数字和空格
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	text = reg.ReplaceAllString(text, "")
	// 移除多余的空格
	text = strings.TrimSpace(text)
	reg = regexp.MustCompile(`\s+`)
	text = reg.ReplaceAllString(text, " ")
	return text
}

// 渲染单词级别的错误高亮
func (m PracticeSession) renderWordLevelError() string {
	userWords := strings.Fields(m.wrongInput)
	expectedWords := strings.Fields(m.expectedText)
	
	var result strings.Builder
	result.WriteString(RenderError("你的输入: "))
	
	// 构建用户输入的高亮文本
	var userInputParts []string
	maxLen := len(userWords)
	if len(expectedWords) > maxLen {
		maxLen = len(expectedWords)
	}
	
	for i := 0; i < maxLen; i++ {
		if i < len(userWords) {
			userWord := userWords[i]
			expectedWord := ""
			if i < len(expectedWords) {
				expectedWord = expectedWords[i]
			}
			
			// 如果单词不匹配，用红色高亮
			if userWord != expectedWord {
				userInputParts = append(userInputParts, RenderError(userWord))
			} else {
				userInputParts = append(userInputParts, userWord)
			}
		}
	}
	
	// 将用户输入拼接并换行
	userInputText := strings.Join(userInputParts, " ")
	wrappedUserInput := m.wrapText(userInputText, m.width-4) // 减去边距
	result.WriteString(wrappedUserInput)
	
	result.WriteString("\n")
	
	// 正确答案也需要换行
	correctAnswerText := "正确答案: " + m.expectedText
	wrappedCorrectAnswer := m.wrapText(correctAnswerText, m.width-4) // 减去边距
	result.WriteString(RenderSuccess(wrappedCorrectAnswer))
	
	return result.String()
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

// wrapText 将文本按指定宽度换行
func (m PracticeSession) wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	currentLine := ""

	for _, word := range words {
		// 如果当前行为空，直接添加单词
		if currentLine == "" {
			currentLine = word
		} else {
			// 检查添加新单词后是否超过宽度
			testLine := currentLine + " " + word
			if len(testLine) <= width {
				currentLine = testLine
			} else {
				// 超过宽度，将当前行添加到结果中，开始新行
				lines = append(lines, currentLine)
				currentLine = word
			}
		}
	}

	// 添加最后一行
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}
