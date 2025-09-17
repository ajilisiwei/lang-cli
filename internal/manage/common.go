package manage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/daiweiwei/lang-cli/internal/config"
)

// 资源类型
const (
	Words     = "words"
	Phrases   = "phrases"
	Sentences = "sentences"
	Articles  = "articles"
)

// GetResourcePath 获取资源文件路径
func GetResourcePath(resourceType string, fileName string) string {
	currentLanguage := config.AppConfig.CurrentLanguage
	// 如果文件名没有.txt后缀，自动添加
	if !strings.HasSuffix(fileName, ".txt") {
		fileName = fileName + ".txt"
	}
	return filepath.Join(getResourceBaseDir(), currentLanguage, resourceType, fileName)
}

// getResourceBaseDir 获取资源基础目录
func getResourceBaseDir() string {
	// 检查是否在测试环境中
	if isTestEnvironment() {
		return "resources"
	}
	
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// 如果无法获取用户主目录，使用当前目录下的resources
		return "resources"
	}
	return filepath.Join(homeDir, ".lang-cli", "resources")
}

// isTestEnvironment 检查是否在测试环境中
func isTestEnvironment() bool {
	// 通过检查环境变量或者调用栈来判断是否在测试中
	for _, arg := range os.Args {
		if strings.Contains(arg, "test") {
			return true
		}
	}
	return false
}

// GetResourceFiles 获取指定类型的资源文件列表
func GetResourceFiles(resourceType string) ([]string, error) {
	currentLanguage := config.AppConfig.CurrentLanguage
	resourcePath := filepath.Join(getResourceBaseDir(), currentLanguage, resourceType)

	// 检查目录是否存在
	if _, err := os.Stat(resourcePath); os.IsNotExist(err) {
		// 创建目录
		if err := os.MkdirAll(resourcePath, 0755); err != nil {
			return nil, fmt.Errorf("创建目录失败: %w", err)
		}
	}

	// 读取目录中的文件
	files, err := os.ReadDir(resourcePath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %w", err)
	}

	// 提取文件名（去掉.txt后缀）
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			fileName := strings.TrimSuffix(file.Name(), ".txt")
			fileNames = append(fileNames, fileName)
		}
	}

	return fileNames, nil
}

// ValidateResourceType 验证资源类型是否有效
func ValidateResourceType(resourceType string) bool {
	switch resourceType {
	case Words, Phrases, Sentences, Articles:
		return true
	default:
		return false
	}
}
