package manage

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// DeleteResource 删除资源
func DeleteResource(resourceType, resourceIdentifier string) error {
	// 验证资源类型
	if !ValidateResourceType(resourceType) {
		return fmt.Errorf("无效的资源类型: %s", resourceType)
	}

	// 获取资源文件路径
	filePath := GetResourcePath(resourceType, resourceIdentifier)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	// 确认删除
	fmt.Printf("确认删除 %s 吗？(y/n): ", resourceIdentifier)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if strings.ToLower(input) != "y" {
		fmt.Println("取消删除。")
		return nil
	}

	// 删除文件
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	fmt.Printf("成功删除 %s\n", resourceIdentifier)
	return nil
}

// DeleteResourceWithoutConfirm 删除资源（不需要用户确认，用于UI界面）
func DeleteResourceWithoutConfirm(resourceType, resourceIdentifier string) error {
	// 验证资源类型
	if !ValidateResourceType(resourceType) {
		return fmt.Errorf("无效的资源类型: %s", resourceType)
	}

	// 获取资源文件路径
	filePath := GetResourcePath(resourceType, resourceIdentifier)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	// 直接删除文件
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// DeleteResourceForTest 删除资源（用于测试，不需要用户确认）
func DeleteResourceForTest(resourceType, resourceIdentifier string) error {
	// 验证资源类型
	if !ValidateResourceType(resourceType) {
		return fmt.Errorf("无效的资源类型: %s", resourceType)
	}

	// 获取资源文件路径
	filePath := GetResourcePath(resourceType, resourceIdentifier)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	// 删除文件
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// ListResourceFiles 列出指定类型的资源文件
func ListResourceFiles(resourceType string) ([]string, error) {
	// 验证资源类型
	if !ValidateResourceType(resourceType) {
		return nil, fmt.Errorf("无效的资源类型: %s", resourceType)
	}

	return GetResourceFiles(resourceType)
}
