package lang

import (
	"os"
	"testing"

	"github.com/ajilisiwei/mllt-cli/internal/config"
)

func TestListLanguages(t *testing.T) {
	// 确保配置已加载
	if err := config.LoadConfig(); err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 测试列出语言
	languages := ListLanguages()
	if len(languages) == 0 {
		t.Error("语言列表为空")
	}

	// 检查是否包含默认语言
	found := false
	for _, lang := range languages {
		if lang == "english" || lang == "japanese" {
			found = true
			break
		}
	}

	if !found {
		t.Error("语言列表中没有找到默认语言")
	}
}

func TestSwitchLanguage(t *testing.T) {
	// 确保配置已加载
	if err := config.LoadConfig(); err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 保存原始语言设置
	originalLang := config.AppConfig.CurrentLanguage
	// 测试结束后恢复原始设置
	defer func() {
		config.AppConfig.CurrentLanguage = originalLang
		config.SaveConfig()
	}()

	// 测试切换语言
	testLang := "english"
	if originalLang == "english" {
		testLang = "japanese"
	}

	// 切换语言
	err := SwitchLanguage(testLang)
	if err != nil {
		t.Errorf("切换语言失败: %v", err)
	}

	// 验证语言已切换
	if config.AppConfig.CurrentLanguage != testLang {
		t.Errorf("语言切换失败，期望 %s，实际 %s", testLang, config.AppConfig.CurrentLanguage)
	}
}

func TestPrintLanguages(t *testing.T) {
	// 重定向标准输出
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 确保配置已加载
	if err := config.LoadConfig(); err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 调用函数
	PrintLanguages()

	// 恢复标准输出
	w.Close()
	os.Stdout = oldStdout

	// 读取输出
	var buf [1024]byte
	n, _ := r.Read(buf[:])
	output := string(buf[:n])

	// 验证输出包含语言列表
	if len(output) == 0 {
		t.Error("PrintLanguages没有输出内容")
	}
}
