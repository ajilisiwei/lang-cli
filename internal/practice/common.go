package practice

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/daiweiwei/lang-cli/internal/config"
)

// 资源类型
const (
	Words     = "words"
	Phrases   = "phrases"
	Sentences = "sentences"
	Articles  = "articles"
)

// 分隔符
const (
	Separator = " ->> "
)

// 初始化随机数生成器
func init() {
	rand.Seed(time.Now().UnixNano())
}

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

// ReadResourceFile 读取资源文件内容
func ReadResourceFile(resourceType string, fileName string) ([]string, error) {
	filePath := GetResourcePath(resourceType, fileName)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 创建空文件
		file, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("创建文件失败: %w", err)
		}
		defer file.Close()
		return []string{}, nil
	}

	// 读取文件内容
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	return lines, nil
}

// ParseLine 解析行内容，返回原文和翻译
func ParseLine(line string) (string, string) {
	parts := strings.Split(line, Separator)
	if len(parts) < 2 {
		return line, ""
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

// GetNextIndex 获取下一个索引
func GetNextIndex(currentIndex, total int, order string) int {
	if order == "sequential" {
		return (currentIndex + 1) % total
	}
	// 随机模式
	return rand.Intn(total)
}
