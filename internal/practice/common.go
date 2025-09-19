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

	// 对于非单词类型的资源，尝试解析换行符分隔格式
	if resourceType != Words {
		if lines, err := parseNewlineFormat(file); err == nil && len(lines) > 0 {
			return lines, nil
		}
		// 如果解析换行符格式失败，回退到原来的逐行解析
		file.Seek(0, 0) // 重置文件指针
	}

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

// parseNewlineFormat 解析换行符分隔格式的资源文件
// 格式：原文\n翻译\n\n（空行分隔不同的条目）
func parseNewlineFormat(file *os.File) ([]string, error) {
	var lines []string
	var currentOriginal string
	var currentTranslation string
	var lineCount int
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineCount++
		
		if line == "" {
			// 遇到空行，如果有当前条目，则保存
			if currentOriginal != "" {
				if currentTranslation != "" {
					// 有翻译的情况
					lines = append(lines, currentOriginal+" ->> "+currentTranslation)
				} else {
					// 没有翻译的情况
					lines = append(lines, currentOriginal)
				}
				currentOriginal = ""
				currentTranslation = ""
			}
		} else if currentOriginal == "" {
			// 第一行是原文
			currentOriginal = line
		} else if currentTranslation == "" {
			// 第二行是翻译
			currentTranslation = line
		} else {
			// 如果已经有原文和翻译，但还有内容，说明格式不对
			// 回退到普通格式解析
			return nil, fmt.Errorf("格式不符合换行符分隔规则")
		}
	}
	
	// 处理文件末尾的条目
	if currentOriginal != "" {
		if currentTranslation != "" {
			lines = append(lines, currentOriginal+" ->> "+currentTranslation)
		} else {
			lines = append(lines, currentOriginal)
		}
	}
	
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	
	// 如果没有解析到任何内容，或者格式不符合要求，返回错误
	if len(lines) == 0 {
		return nil, fmt.Errorf("没有找到符合换行符分隔格式的内容")
	}
	
	return lines, nil
}

// ParseLine 解析行内容，返回原文和翻译
// 支持多种分隔符，按优先级顺序：" ->> ", ":", "：", "/", 空格
func ParseLine(line string) (string, string) {
	// 定义分隔符优先级列表
	separators := []string{" ->> ", ":", "：", "/"}
	
	// 按优先级尝试分隔符
	for _, sep := range separators {
		if strings.Contains(line, sep) {
			parts := strings.SplitN(line, sep, 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			}
		}
	}
	
	// 如果没有找到任何分隔符，使用空格分隔
	// 第一个空格之前的文本作为原文，之后的作为翻译
	spaceIndex := strings.Index(line, " ")
	if spaceIndex > 0 {
		return strings.TrimSpace(line[:spaceIndex]), strings.TrimSpace(line[spaceIndex+1:])
	}
	
	// 如果没有空格，整行作为原文，翻译为空
	return line, ""
}

// GetNextIndex 获取下一个索引
func GetNextIndex(currentIndex, total int, order string) int {
	if order == "sequential" {
		return (currentIndex + 1) % total
	}
	// 随机模式
	return rand.Intn(total)
}
