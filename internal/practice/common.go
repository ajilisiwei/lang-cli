package practice

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
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
// 优先返回用户资源路径，若不存在则回退到基础资源目录
func GetResourcePath(resourceType string, fileName string) string {
	currentLanguage := config.AppConfig.CurrentLanguage
	if !strings.HasSuffix(fileName, ".txt") {
		fileName = fileName + ".txt"
	}

	userPath := filepath.Join(getUserDataBaseDir(), currentLanguage, resourceType, fileName)
	if _, err := os.Stat(userPath); err == nil {
		return userPath
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

func getUserDataBaseDir() string {
	if isTestEnvironment() {
		return filepath.Join("resources", "user-data")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("resources", "user-data")
	}
	return filepath.Join(homeDir, ".lang-cli", "user-data")
}

// GetUserDataDir 返回用于存储用户练习数据的目录
func GetUserDataDir() string {
	return getUserDataBaseDir()
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
	baseDir := filepath.Join(getResourceBaseDir(), currentLanguage, resourceType)
	userDir := filepath.Join(getUserDataBaseDir(), currentLanguage, resourceType)

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("创建基础资源目录失败: %w", err)
	}
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, fmt.Errorf("创建用户资源目录失败: %w", err)
	}

	fileSet := make(map[string]struct{})

	addFiles := func(dir string) error {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if !strings.HasSuffix(name, ".txt") {
				continue
			}
			trimmed := strings.TrimSuffix(name, ".txt")
			fileSet[trimmed] = struct{}{}
		}
		return nil
	}

	if err := addFiles(baseDir); err != nil {
		return nil, fmt.Errorf("读取基础资源目录失败: %w", err)
	}
	if err := addFiles(userDir); err != nil {
		return nil, fmt.Errorf("读取用户资源目录失败: %w", err)
	}

	var fileNames []string
	for name := range fileSet {
		fileNames = append(fileNames, name)
	}

	sort.Strings(fileNames)

	return fileNames, nil
}

// ReadResourceFile 读取资源文件内容
func ReadResourceFile(resourceType string, fileName string) ([]string, error) {
	filePath, err := resolveReadableResourcePath(resourceType, fileName)
	if err != nil {
		return nil, err
	}

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

// WriteResourceFile 将内容写入资源文件（覆盖写入）
func WriteResourceFile(resourceType string, fileName string, lines []string) error {
	filePath := getWritableResourcePath(resourceType, fileName)

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("创建资源目录失败: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("打开资源文件失败: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("写入资源文件失败: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("刷新资源文件失败: %w", err)
	}

	return nil
}

func resolveReadableResourcePath(resourceType, fileName string) (string, error) {
	currentLanguage := config.AppConfig.CurrentLanguage
	if !strings.HasSuffix(fileName, ".txt") {
		fileName = fileName + ".txt"
	}

	userPath := filepath.Join(getUserDataBaseDir(), currentLanguage, resourceType, fileName)
	if _, err := os.Stat(userPath); err == nil {
		return userPath, nil
	}

	basePath := filepath.Join(getResourceBaseDir(), currentLanguage, resourceType, fileName)
	if _, err := os.Stat(basePath); err == nil {
		return basePath, nil
	}

	if err := os.MkdirAll(filepath.Dir(userPath), 0755); err != nil {
		return "", fmt.Errorf("创建用户资源目录失败: %w", err)
	}

	file, err := os.Create(userPath)
	if err != nil {
		return "", fmt.Errorf("创建资源文件失败: %w", err)
	}
	file.Close()

	return userPath, nil
}

func getWritableResourcePath(resourceType, fileName string) string {
	currentLanguage := config.AppConfig.CurrentLanguage
	if !strings.HasSuffix(fileName, ".txt") {
		fileName = fileName + ".txt"
	}

	return filepath.Join(getUserDataBaseDir(), currentLanguage, resourceType, fileName)
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
			if containsInlineSeparator(line) {
				return nil, fmt.Errorf("检测到行内分隔符，回退到逐行解析")
			}
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

func containsInlineSeparator(line string) bool {
	if strings.Contains(line, " ->> ") {
		return true
	}
	return false
}

// ParseLine 解析行内容，返回原文和翻译
// 支持多种分隔符，按优先级顺序：" ->> ", 制表符, 空格, "/", ":", "："
func ParseLine(line string) (string, string) {
	// 首先检查 " ->> " 分隔符
	if strings.Contains(line, " ->> ") {
		parts := strings.SplitN(line, " ->> ", 2)
		if len(parts) == 2 {
			return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		}
	}

	// 然后检查制表符分隔符
	// 第一个制表符之前的文本作为原文，之后的作为翻译
	tabIndex := strings.Index(line, "\t")
	if tabIndex > 0 {
		return strings.TrimSpace(line[:tabIndex]), strings.TrimSpace(line[tabIndex+1:])
	}

	// 然后检查空格分隔符
	// 第一个空格之前的文本作为原文，之后的作为翻译
	spaceIndex := strings.Index(line, " ")
	if spaceIndex > 0 {
		return strings.TrimSpace(line[:spaceIndex]), strings.TrimSpace(line[spaceIndex+1:])
	}

	// 定义其他分隔符优先级列表
	separators := []string{"/", ":", "："}

	// 按优先级尝试其他分隔符
	for _, sep := range separators {
		if strings.Contains(line, sep) {
			parts := strings.SplitN(line, sep, 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			}
		}
	}

	// 如果没有找到任何分隔符，整行作为原文，翻译为空
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
