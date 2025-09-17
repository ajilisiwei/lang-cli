package practice

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/daiweiwei/lang-cli/internal/config"
)

// ArticlePractice 文章练习
func ArticlePractice(fileName string) error {
	// 读取文章内容
	lines, err := ReadResourceFile(Articles, fileName)
	if err != nil {
		return err
	}

	// 检查文章是否为空
	if len(lines) == 0 {
		fmt.Println("文章为空，请先添加文章内容。")
		return nil
	}

	fmt.Printf("开始练习文章: %s\n", fileName)
	fmt.Println("输入 'q' 退出练习。")

	// 获取配置
	showTranslation := config.AppConfig.Articles.ShowTranslation

	// 开始练习
	reader := bufio.NewReader(os.Stdin)

	for i := 0; i < len(lines); i++ {
		// 获取当前行
		line := lines[i]
		text, translation := ParseLine(line)

		// 显示当前行
		fmt.Printf("请输入第 %d 行: %s\n", i+1, text)

		// 读取用户输入
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// 检查是否退出
		if input == "q" {
			fmt.Println("练习结束。")
			return nil
		}

		// 检查输入是否正确
		if input == text {
			fmt.Println("正确！")
			if showTranslation && translation != "" {
				fmt.Printf("翻译: %s\n", translation)
			}
		} else {
			fmt.Println("错误，请重新输入。")
			i-- // 重新练习当前行
		}

		fmt.Println() // 空行分隔
	}

	fmt.Println("恭喜！文章练习完成。")
	return nil
}

// ListArticleFiles 列出文章文件
func ListArticleFiles() ([]string, error) {
	return GetResourceFiles(Articles)
}
