package practice

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ajilisiwei/mllt-cli/internal/config"
)

// PhrasePractice 短语练习
func PhrasePractice(fileName string) error {
	// 读取短语列表
	phrases, err := ReadResourceFile(Phrases, fileName)
	if err != nil {
		return err
	}

	// 检查短语列表是否为空
	if len(phrases) == 0 {
		fmt.Println("短语列表为空，请先添加短语。")
		return nil
	}

	fmt.Printf("开始练习短语列表: %s\n", fileName)
	fmt.Println("输入 'q' 退出练习。")

	// 获取配置
	nextOneOrder := config.AppConfig.NextOneOrder
	showTranslation := config.AppConfig.ShowTranslation

	// 开始练习
	index := 0
	reader := bufio.NewReader(os.Stdin)

	for {
		// 获取当前短语
		phraseLine := phrases[index]
		phrase, translation := ParseLine(phraseLine)

		// 显示短语
		fmt.Printf("请输入: %s\n", phrase)

		// 读取用户输入
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// 检查是否退出
		if input == "q" {
			fmt.Println("练习结束。")
			break
		}

		// 检查输入是否正确
		if input == phrase {
			fmt.Println("正确！")
			if showTranslation && translation != "" {
				fmt.Printf("翻译: %s\n", translation)
			}
			// 获取下一个短语的索引
			index = GetNextIndex(index, len(phrases), nextOneOrder)
		} else {
			fmt.Println("错误，请重新输入。")
		}

		fmt.Println() // 空行分隔
	}

	return nil
}

// ListPhraseFiles 列出短语文件
func ListPhraseFiles() ([]string, error) {
	return GetResourceFiles(Phrases)
}
