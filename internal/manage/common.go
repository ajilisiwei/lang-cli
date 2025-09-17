package manage

import (
	"fmt"
	"os"
	"path/filepath"

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
	return filepath.Join("resources", currentLanguage, resourceType, fileName)
}

// GetResourceFiles 获取指定类型的资源文件列表
func GetResourceFiles(resourceType string) ([]string, error) {
	currentLanguage := config.AppConfig.CurrentLanguage
	resourcePath := filepath.Join("resources", currentLanguage, resourceType)

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

	// 提取文件名
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
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
