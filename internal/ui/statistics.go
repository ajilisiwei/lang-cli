package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ajilisiwei/lang-cli/internal/statistics"
)

// StatisticsMenu 统计菜单
type StatisticsMenu struct {
	list     list.Model
	quitting bool
}

// StatisticsSummaryItem 用于展示每天的统计摘要
type StatisticsSummaryItem struct {
	summary statistics.DailySummary
}

func (i StatisticsSummaryItem) Title() string {
	return fmt.Sprintf("%s - %d 次练习", i.summary.Date, i.summary.SessionCount)
}

func (i StatisticsSummaryItem) Description() string {
	return fmt.Sprintf("正确率 %.1f%% · 正确 %d · 错误 %d", i.summary.Accuracy, i.summary.Correct, i.summary.Incorrect)
}

func (i StatisticsSummaryItem) FilterValue() string { return i.summary.Date }

// StatisticsSessionItem 用于展示单次练习记录
type StatisticsSessionItem struct {
	record statistics.SessionRecord
}

func (i StatisticsSessionItem) Title() string {
	status := "完成"
	if !i.record.Completed {
		status = "未完成"
	}
	timestamp := i.record.Timestamp.Local().Format("15:04:05")
	return fmt.Sprintf("%s · %s · %s", timestamp, strings.ToUpper(i.record.ResourceType), status)
}

func (i StatisticsSessionItem) Description() string {
	return fmt.Sprintf("%s · 正确 %d / %d · %.1f%% · 用时 %ds", i.record.FileName, i.record.Correct, i.record.Total, i.record.Accuracy, i.record.DurationSeconds)
}

func (i StatisticsSessionItem) FilterValue() string { return i.record.FileName }

// NewStatisticsMenu 创建统计菜单
func NewStatisticsMenu() *StatisticsMenu {
	summaries, err := statistics.GetDailySummaries()
	if err != nil {
		summaries = []statistics.DailySummary{}
	}

	items := []list.Item{}
	if len(summaries) == 0 {
		items = append(items, MenuItem{
			title:       "暂无统计数据",
			description: "当前没有可用的练习统计",
			action:      nil,
		})
	} else {
		for _, summary := range summaries {
			summaryCopy := summary
			items = append(items, StatisticsSummaryItem{summary: summaryCopy})
		}
	}

	items = append(items, MenuItem{
		title:       "返回主菜单",
		description: "返回到主菜单",
		action:      func() (tea.Model, error) { return NewMainMenu(), nil },
	})

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "练习统计"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	return &StatisticsMenu{list: l}
}

func (m StatisticsMenu) Init() tea.Cmd {
	return nil
}

func (m StatisticsMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			item := m.list.SelectedItem()
			switch selected := item.(type) {
			case StatisticsSummaryItem:
				detail := NewStatisticsDetailView(selected.summary.Date)
				if m.list.Width() > 0 {
					if updated, cmd := detail.Update(tea.WindowSizeMsg{Width: m.list.Width(), Height: m.list.Height() + 4}); updated != nil {
						return updated, cmd
					}
				}
				return detail, nil
			case MenuItem:
				if selected.action != nil {
					newModel, err := selected.action()
					if err != nil {
						return m, nil
					}
					if m.list.Width() > 0 {
						if updated, cmd := newModel.Update(tea.WindowSizeMsg{Width: m.list.Width(), Height: m.list.Height() + 4}); updated != nil {
							return updated, cmd
						}
					}
					return newModel, nil
				}
			}
		case "esc":
			mainMenu := NewMainMenu()
			if m.list.Width() > 0 {
				if updated, cmd := mainMenu.Update(tea.WindowSizeMsg{Width: m.list.Width(), Height: m.list.Height() + 4}); updated != nil {
					return updated, cmd
				}
			}
			return mainMenu, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m StatisticsMenu) View() string {
	if m.quitting {
		return "再见！"
	}
	return m.list.View()
}

// StatisticsDetailView 详情视图
type StatisticsDetailView struct {
	date  string
	list  list.Model
	empty bool
}

// NewStatisticsDetailView 创建详情视图
func NewStatisticsDetailView(date string) *StatisticsDetailView {
	records, err := statistics.GetSessionsByDate(date)
	if err != nil {
		records = []statistics.SessionRecord{}
	}

	items := []list.Item{}
	for _, record := range records {
		items = append(items, StatisticsSessionItem{record: record})
	}

	if len(items) == 0 {
		items = append(items, MenuItem{
			title:       "暂无记录",
			description: "该日期没有练习记录",
			action:      nil,
		})
	}

	items = append(items, MenuItem{
		title:       "返回统计概览",
		description: "返回到统计列表",
		action:      func() (tea.Model, error) { return NewStatisticsMenu(), nil },
	})

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = fmt.Sprintf("%s 练习详情", date)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	return &StatisticsDetailView{date: date, list: l, empty: len(records) == 0}
}

func (m StatisticsDetailView) Init() tea.Cmd {
	return nil
}

func (m StatisticsDetailView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			menu := NewStatisticsMenu()
			if m.list.Width() > 0 {
				if updated, cmd := menu.Update(tea.WindowSizeMsg{Width: m.list.Width(), Height: m.list.Height() + 4}); updated != nil {
					return updated, cmd
				}
			}
			return menu, nil
		case "enter":
			if menuItem, ok := m.list.SelectedItem().(MenuItem); ok && menuItem.action != nil {
				newModel, err := menuItem.action()
				if err != nil {
					return m, nil
				}
				if m.list.Width() > 0 {
					if updated, cmd := newModel.Update(tea.WindowSizeMsg{Width: m.list.Width(), Height: m.list.Height() + 4}); updated != nil {
						return updated, cmd
					}
				}
				return newModel, nil
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m StatisticsDetailView) View() string {
	return m.list.View()
}
