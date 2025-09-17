package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 备份原始配置文件
	originalConfig, err := os.ReadFile("../../config/config.yaml")
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("读取原始配置文件失败: %v", err)
	}

	// 测试结束后恢复原始配置文件
	defer func() {
		if originalConfig != nil {
			if err := os.WriteFile("../../config/config.yaml", originalConfig, 0644); err != nil {
				t.Logf("恢复原始配置文件失败: %v", err)
			}
		}
	}()

	// 测试加载配置
	err = LoadConfig()
	if err != nil {
		t.Errorf("加载配置失败: %v", err)
	}

	// 验证配置已加载
	if AppConfig.CurrentLanguage == "" {
		t.Error("当前语言为空")
	}
}

func TestSaveConfig(t *testing.T) {
	// 备份原始配置文件
	originalConfig, err := os.ReadFile("../../config/config.yaml")
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("读取原始配置文件失败: %v", err)
	}

	// 测试结束后恢复原始配置文件
	defer func() {
		if originalConfig != nil {
			if err := os.WriteFile("../../config/config.yaml", originalConfig, 0644); err != nil {
				t.Logf("恢复原始配置文件失败: %v", err)
			}
		}
	}()

	// 确保配置已加载
	if err := LoadConfig(); err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 修改配置
	originalLang := AppConfig.CurrentLanguage
	testLang := "english"
	if originalLang == "english" {
		testLang = "japanese"
	}
	AppConfig.CurrentLanguage = testLang

	// 保存配置
	err = SaveConfig()
	if err != nil {
		t.Errorf("保存配置失败: %v", err)
	}

	// 重新加载配置
	AppConfig.CurrentLanguage = originalLang
	err = LoadConfig()
	if err != nil {
		t.Errorf("重新加载配置失败: %v", err)
	}

	// 验证配置已保存
	if AppConfig.CurrentLanguage != testLang {
		t.Errorf("配置保存失败，期望 %s，实际 %s", testLang, AppConfig.CurrentLanguage)
	}

	// 恢复原始配置
	AppConfig.CurrentLanguage = originalLang
	SaveConfig()
}
