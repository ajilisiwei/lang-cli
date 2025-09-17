package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// 定义样式常量
var (
	// 标题样式
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Margin(1, 0)

	// 普通文本样式
	TextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDDDDD"))

	// 高亮文本样式
	HighlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	// 错误文本样式
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	// 成功文本样式
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	// 按钮样式
	ButtonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#888888")).
			Padding(0, 3).
			Margin(0, 1)

	// 选中按钮样式
	ActiveButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 3).
				Margin(0, 1)
)

// RenderTitle 渲染标题
func RenderTitle(title string) string {
	return TitleStyle.Render(title)
}

// RenderText 渲染普通文本
func RenderText(text string) string {
	return TextStyle.Render(text)
}

// RenderHighlight 渲染高亮文本
func RenderHighlight(text string) string {
	return HighlightStyle.Render(text)
}

// RenderError 渲染错误文本
func RenderError(text string) string {
	return ErrorStyle.Render(text)
}

// RenderSuccess 渲染成功文本
func RenderSuccess(text string) string {
	return SuccessStyle.Render(text)
}

// RenderButton 渲染按钮
func RenderButton(text string, active bool) string {
	if active {
		return ActiveButtonStyle.Render(text)
	}
	return ButtonStyle.Render(text)
}

// CenterText 居中显示文本
func CenterText(text string, width int) string {
	if width <= 0 {
		return text
	}

	padding := (width - lipgloss.Width(text)) / 2
	if padding <= 0 {
		return text
	}

	return strings.Repeat(" ", padding) + text
}
