package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIsInDevelopmentMode 测试开发模式检测
func TestIsInDevelopmentMode(t *testing.T) {
	// 在项目根目录下应该检测为开发模式
	if !isInDevelopmentMode() {
		t.Error("在项目根目录下应该检测为开发模式")
	}

	// 创建临时目录测试生产模式
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	// 切换到临时目录
	os.Chdir(tempDir)
	if isInDevelopmentMode() {
		t.Error("在没有go.mod的目录下应该检测为生产模式")
	}
}

// TestConfigLoadingInDifferentModes 测试不同模式下的配置加载
func TestConfigLoadingInDifferentModes(t *testing.T) {
	// 保存原始工作目录
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	// 测试开发模式
	t.Run("开发模式", func(t *testing.T) {
		// 确保在项目根目录
		os.Chdir(originalWd)
		err := LoadConfig()
		if err != nil {
			t.Fatalf("开发模式下加载配置失败: %v", err)
		}
		// 验证配置是否正确加载
		if AppConfig.CurrentLanguage == "" {
			t.Error("配置未正确加载")
		}
	})

	// 测试生产模式（模拟）
	t.Run("生产模式模拟", func(t *testing.T) {
		// 创建临时目录和配置文件
		tempDir := t.TempDir()
		homeDir := filepath.Join(tempDir, "home")
		configDir := filepath.Join(homeDir, ".mllt-cli")
		os.MkdirAll(configDir, 0755)

		// 创建测试配置文件
		configContent := `current_language: english
languages:
  - english
  - japanese
words:
  show_translation: true
  next_one_order: ""
phrases:
  show_translation: true
  next_one_order: ""
sentences:
  show_translation: true
  next_one_order: ""
articles:
  show_translation: true
correctness_match_mode: exact_match
next_one_order: random`
		configFile := filepath.Join(configDir, "config.yaml")
		os.WriteFile(configFile, []byte(configContent), 0644)

		// 设置临时的HOME环境变量
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", homeDir)
		defer os.Setenv("HOME", originalHome)

		// 切换到没有go.mod的目录
		workDir := filepath.Join(tempDir, "work")
		os.MkdirAll(workDir, 0755)
		os.Chdir(workDir)

		// 加载配置
		err := LoadConfig()
		if err != nil {
			t.Fatalf("生产模式下加载配置失败: %v", err)
		}

		// 验证配置是否正确加载
		if AppConfig.CurrentLanguage != "english" {
			t.Errorf("期望当前语言为english，实际为: %s", AppConfig.CurrentLanguage)
		}
	})
}