package manage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ajilisiwei/mllt-cli/internal/config"
	"github.com/ajilisiwei/mllt-cli/internal/practice"
)

const testFolderName = "custom_folder"

// 测试前的准备工作
func setupTest(t *testing.T) {
	// 确保配置已加载
	if err := config.LoadConfig(); err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 确保资源目录存在
	resourceDirs := []string{
		filepath.Join("resources", config.AppConfig.CurrentLanguage, practice.Words),
		filepath.Join("resources", config.AppConfig.CurrentLanguage, practice.Phrases),
		filepath.Join("resources", config.AppConfig.CurrentLanguage, practice.Sentences),
		filepath.Join("resources", config.AppConfig.CurrentLanguage, practice.Articles),
	}

	for _, dir := range resourceDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatalf("创建目录失败: %v", err)
			}
		}
	}

	// 创建测试文件
	testFiles := map[string][]string{
		practice.Words: {"test_manage_words.txt", "apple ->> 苹果\nbanana ->> 香蕉\norange ->> 橙子"},
	}

	for resourceType, fileInfo := range testFiles {
		filePath := filepath.Join("resources", config.AppConfig.CurrentLanguage, resourceType, fileInfo[0])
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}
		defer file.Close()

		if _, err := file.WriteString(fileInfo[1]); err != nil {
			t.Fatalf("写入测试文件失败: %v", err)
		}
	}

	// 创建导入测试文件
	importFilePath := filepath.Join(os.TempDir(), "import_test_words.txt")
	importFile, err := os.Create(importFilePath)
	if err != nil {
		t.Fatalf("创建导入测试文件失败: %v", err)
	}
	defer importFile.Close()

	if _, err := importFile.WriteString("grape ->> 葡萄\npeach ->> 桃子\n"); err != nil {
		t.Fatalf("写入导入测试文件失败: %v", err)
	}
}

// 测试后的清理工作
func cleanupTest(t *testing.T) {
	// 删除测试文件
	testFiles := []string{
		filepath.Join("resources", config.AppConfig.CurrentLanguage, practice.Words, "test_manage_words.txt"),
		filepath.Join("resources", "user-data", config.AppConfig.CurrentLanguage, practice.Words, testFolderName, "import_test_words.txt"),
	}

	for _, file := range testFiles {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			t.Logf("删除测试文件失败: %v", err)
		}
	}

	userFolderPath := filepath.Join("resources", "user-data", config.AppConfig.CurrentLanguage, practice.Words, testFolderName)
	if err := os.Remove(userFolderPath); err != nil && !os.IsNotExist(err) {
		t.Logf("删除测试文件夹失败: %v", err)
	}

	// 删除导入测试文件
	importFilePath := filepath.Join(os.TempDir(), "import_test_words.txt")
	if err := os.Remove(importFilePath); err != nil && !os.IsNotExist(err) {
		t.Logf("删除导入测试文件失败: %v", err)
	}
}

func TestValidateResourceType(t *testing.T) {
	// 测试验证资源类型
	tests := []struct {
		resourceType string
		expected     bool
	}{
		{practice.Words, true},
		{practice.Phrases, true},
		{practice.Sentences, true},
		{practice.Articles, true},
		{"invalid", false},
	}

	for _, test := range tests {
		result := ValidateResourceType(test.resourceType)
		if result != test.expected {
			t.Errorf("验证资源类型失败，资源类型 %s，期望 %v，实际 %v", test.resourceType, test.expected, result)
		}
	}
}

func TestListResourceFiles(t *testing.T) {
	// 设置测试环境
	setupTest(t)
	defer cleanupTest(t)

	// 测试列出资源文件
	files, err := ListResourceFiles(practice.Words)
	if err != nil {
		t.Errorf("列出资源文件失败: %v", err)
	}

	// 验证文件列表包含测试文件
	found := false
	for _, file := range files {
		if file == "test_manage_words" {
			found = true
			break
		}
	}

	if !found {
		t.Error("资源文件列表中没有找到测试文件")
	}
}

func TestDeleteResource(t *testing.T) {
	// 设置测试环境
	setupTest(t)
	defer cleanupTest(t)

	// 测试删除资源，使用测试专用的删除函数，不需要用户确认
	err := DeleteResourceForTest(practice.Words, "test_manage_words.txt")
	if err != nil {
		t.Errorf("删除资源失败: %v", err)
	}

	// 验证文件已删除
	filePath := filepath.Join("resources", config.AppConfig.CurrentLanguage, practice.Words, "test_manage_words.txt")
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("文件未被删除")
	}
}

func TestImportResource(t *testing.T) {
	// 设置测试环境
	setupTest(t)
	defer cleanupTest(t)

	// 测试导入资源
	importFilePath := filepath.Join(os.TempDir(), "import_test_words.txt")
	// 使用测试专用的导入函数，不需要用户确认
	err := ImportResourceForTest(practice.Words, testFolderName, importFilePath)
	if err != nil {
		t.Errorf("导入资源失败: %v", err)
	}

	// 验证文件已导入
	filePath := filepath.Join("resources", "user-data", config.AppConfig.CurrentLanguage, practice.Words, testFolderName, "import_test_words.txt")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("文件未被导入")
	}

	// 验证文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("读取导入文件失败: %v", err)
	}

	expectedContent := "grape ->> 葡萄\npeach ->> 桃子\n"
	if string(content) != expectedContent {
		t.Errorf("导入文件内容不匹配，期望 %s，实际 %s", expectedContent, string(content))
	}
}
