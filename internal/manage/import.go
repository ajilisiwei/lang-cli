package manage

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ajilisiwei/lang-cli/internal/practice"
)

// parseNewlineFormatForImport 解析换行符分隔格式的文件
func ParseNewlineFormatForImport(file *os.File) ([]string, error) {
	var result []string
	scanner := bufio.NewScanner(file)
	var currentPair []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行
		if line == "" {
			continue
		}

		currentPair = append(currentPair, line)

		// 当收集到两行时，组成一对
		if len(currentPair) == 2 {
			// 使用 ->> 作为分隔符
			result = append(result, currentPair[0]+" ->> "+currentPair[1])
			currentPair = nil
		}
	}

	// 如果还有剩余的单行，跳过它
	if len(currentPair) > 0 {
		// 可以选择记录警告或直接忽略
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// isValidSeparatorFormat 检查行是否包含有效的分隔符格式
func isValidSeparatorFormat(line string) bool {
	// 首先检查 " ->> " 分隔符
	if strings.Contains(line, " ->> ") {
		return true
	}

	// 然后检查制表符分隔符
	if strings.Contains(line, "\t") {
		return true
	}

	// 然后检查空格分隔符
	if strings.Contains(line, " ") {
		return true
	}

	// 定义其他支持的分隔符，按优先级顺序："/", ":", "："
	separators := []string{"/", ":", "："}

	// 检查是否包含其他分隔符
	for _, sep := range separators {
		if strings.Contains(line, sep) {
			return true
		}
	}

	return false
}

// ImportResource 导入资源
func ImportResource(resourceType, folderDir, sourcePath string) error {
	// 验证资源类型
	if !ValidateResourceType(resourceType) {
		return fmt.Errorf("无效的资源类型: %s", resourceType)
	}

	// 检查源文件是否存在
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("源文件不存在: %s", sourcePath)
	}

	// 检查文件扩展名
	if filepath.Ext(sourcePath) != ".txt" {
		return fmt.Errorf("文件必须是.txt格式: %s", sourcePath)
	}

	// 获取目标文件名
	fileName := filepath.Base(sourcePath)
	baseName := strings.TrimSuffix(fileName, ".txt")
	identifier := practice.BuildResourceIdentifier(folderDir, baseName)
	targetPath := GetResourcePath(resourceType, identifier)
	fmt.Printf("调试: 目标路径 = %s\n", targetPath)

	// 检查目标文件是否已存在
	if _, err := os.Stat(targetPath); err == nil {
		// 确认覆盖
		fmt.Printf("文件 %s 已存在，是否覆盖？(y/n): ", fileName)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.ToLower(input) != "y" {
			fmt.Println("取消导入。")
			return nil
		}
	}

	// 确认导入
	fmt.Printf("确认导入 %s 到 %s/%s 吗？(y/n): ", sourcePath, resourceType, practice.FormatResourceDisplayName(identifier))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if strings.ToLower(input) != "y" {
		fmt.Println("取消导入。")
		return nil
	}

	// 复制文件
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer sourceFile.Close()

	// 确保目标目录存在
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer targetFile.Close()

	// 直接复制文件，不进行格式转换
	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	fmt.Printf("成功导入 %s 到 %s/%s\n", fileName, resourceType, practice.FormatResourceDisplayName(identifier))
	return nil
}

// ImportResourceForTest 导入资源（用于测试，不需要用户确认）
// ImportResourceForTest 导入资源（用于测试，不需要用户确认）
func ImportResourceForTest(resourceType, folderDir, sourcePath string) error {
	// 验证资源类型
	if !ValidateResourceType(resourceType) {
		return fmt.Errorf("无效的资源类型: %s", resourceType)
	}

	// 检查源文件是否存在
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("源文件不存在: %s", sourcePath)
	}

	// 检查文件扩展名
	if filepath.Ext(sourcePath) != ".txt" {
		return fmt.Errorf("文件必须是.txt格式: %s", sourcePath)
	}

	// 获取目标文件名
	fileName := filepath.Base(sourcePath)
	baseName := strings.TrimSuffix(fileName, ".txt")
	identifier := practice.BuildResourceIdentifier(folderDir, baseName)
	targetPath := GetResourcePath(resourceType, identifier)

	// 复制文件
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer sourceFile.Close()

	// 确保目标目录存在
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer targetFile.Close()

	// 直接复制文件，不进行格式转换
	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	fmt.Printf("成功导入 %s 到 %s/%s\n", fileName, resourceType, practice.FormatResourceDisplayName(identifier))
	return nil
}
