package statistics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ajilisiwei/mllt-cli/internal/practice"
)

// SessionRecord 记录一次练习的统计数据
type SessionRecord struct {
	Timestamp       time.Time `json:"timestamp"`
	ResourceType    string    `json:"resource_type"`
	FileName        string    `json:"file_name"`
	Total           int       `json:"total"`
	Correct         int       `json:"correct"`
	Incorrect       int       `json:"incorrect"`
	Accuracy        float64   `json:"accuracy"`
	DurationSeconds int64     `json:"duration_seconds"`
	OrderMode       string    `json:"order_mode"`
	Completed       bool      `json:"completed"`
}

// DailySummary 汇总某一天的统计数据
type DailySummary struct {
	Date         string
	SessionCount int
	Total        int
	Correct      int
	Incorrect    int
	Accuracy     float64
}

func statsDir() (string, error) {
	base := practice.GetUserDataDir()
	dir := filepath.Join(base, "statistics")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

// LogSession 记录一次练习结果
func LogSession(record SessionRecord) error {
	dir, err := statsDir()
	if err != nil {
		return err
	}

	date := record.Timestamp.Local().Format("2006-01-02")
	path := filepath.Join(dir, date+".json")

	var records []SessionRecord
	if data, err := os.ReadFile(path); err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &records); err != nil {
			return fmt.Errorf("解析统计文件失败: %w", err)
		}
	}

	records = append(records, record)

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetDailySummaries 返回按日期汇总的统计信息
func GetDailySummaries() ([]DailySummary, error) {
	dir, err := statsDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	summaries := make([]DailySummary, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if filepath.Ext(name) != ".json" {
			continue
		}

		date := strings.TrimSuffix(name, ".json")
		records, err := GetSessionsByDate(date)
		if err != nil {
			return nil, err
		}

		summary := DailySummary{Date: date}
		for _, record := range records {
			summary.SessionCount++
			summary.Total += record.Total
			summary.Correct += record.Correct
			summary.Incorrect += record.Incorrect
		}

		if summary.Total > 0 {
			summary.Accuracy = float64(summary.Correct) / float64(summary.Total) * 100
		}

		summaries = append(summaries, summary)
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Date > summaries[j].Date
	})

	return summaries, nil
}

// GetSessionsByDate 返回指定日期的所有练习记录
func GetSessionsByDate(date string) ([]SessionRecord, error) {
	dir, err := statsDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, date+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []SessionRecord{}, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return []SessionRecord{}, nil
	}

	var records []SessionRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("解析统计文件失败: %w", err)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp.After(records[j].Timestamp)
	})

	return records, nil
}
