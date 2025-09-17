package manage

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ImportResource 导入资源
func ImportResource(resourceType, sourcePath string) error {
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
	targetPath := GetResourcePath(resourceType, fileName)

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
	fmt.Printf("确认导入 %s 到 %s 吗？(y/n): ", sourcePath, resourceType)
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

	// 验证文件格式（如果是单词、短语或句子，需要检查格式）
	if resourceType == Words || resourceType == Phrases || resourceType == Sentences {
		scanner := bufio.NewScanner(sourceFile)
		writer := bufio.NewWriter(targetFile)

		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			// 跳过空行
			if strings.TrimSpace(line) == "" {
				continue
			}

			// 检查格式
			if !strings.Contains(line, " ->> ") {
				fmt.Printf("警告: 第 %d 行格式不正确，应为 '内容 ->> 翻译'\n", lineNum)
				// 继续处理，不中断导入
			}

			// 写入行
			fmt.Fprintln(writer, line)
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("读取源文件失败: %w", err)
		}

		if err := writer.Flush(); err != nil {
			return fmt.Errorf("写入目标文件失败: %w", err)
		}
	} else {
		// 对于文章，直接复制
		if _, err := io.Copy(targetFile, sourceFile); err != nil {
			return fmt.Errorf("复制文件失败: %w", err)
		}
	}

	fmt.Printf("成功导入 %s 到 %s\n", fileName, resourceType)
	return nil
}

// ImportResourceForTest 导入资源（用于测试，不需要用户确认）
// ImportResourceForTest 导入资源（用于测试，不需要用户确认）
func ImportResourceForTest(resourceType, sourcePath string) error {
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
	targetPath := GetResourcePath(resourceType, fileName)

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

	// 验证文件格式（如果是单词、短语或句子，需要检查格式）
	if resourceType == Words || resourceType == Phrases || resourceType == Sentences {
		scanner := bufio.NewScanner(sourceFile)
		writer := bufio.NewWriter(targetFile)

		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			// 跳过空行
			if strings.TrimSpace(line) == "" {
				continue
			}

			// 检查格式
			if !strings.Contains(line, " ->> ") {
				fmt.Printf("警告: 第 %d 行格式不正确，应为 '内容 ->> 翻译'\n", lineNum)
				// 继续处理，不中断导入
			}

			// 写入行
			fmt.Fprintln(writer, line)
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("读取源文件失败: %w", err)
		}

		if err := writer.Flush(); err != nil {
			return fmt.Errorf("写入目标文件失败: %w", err)
		}
	} else {
		// 对于文章，直接复制
		if _, err := io.Copy(targetFile, sourceFile); err != nil {
			return fmt.Errorf("复制文件失败: %w", err)
		}
	}

	return nil
}
