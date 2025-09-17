package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/daiweiwei/lang-cli/internal/config"
	"github.com/daiweiwei/lang-cli/internal/lang"
	"github.com/daiweiwei/lang-cli/internal/manage"
	"github.com/daiweiwei/lang-cli/internal/practice"
	"github.com/daiweiwei/lang-cli/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lang-cli",
	Short: "一个多语言打字学习终端工具",
	Long: `这是一个支持多种语言的打字学习终端工具，
用户可以在终端中进行不同语言的单词、短语、句子、文章等的打字练习。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 显示主菜单
		fmt.Println("欢迎使用语言学习终端！")
		// 启动UI界面
		p := tea.NewProgram(ui.NewMainMenu(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("启动UI界面失败: %v\n", err)
			os.Exit(1)
		}
	},
}

// langCmd 表示lang子命令
var langCmd = &cobra.Command{
	Use:   "lang",
	Short: "语言管理模块",
	Long:  `语言管理模块，用于列出支持的语言和切换当前练习的语言。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("语言管理模块")
		// 这里将来会调用lang模块的功能
	},
}

// langLsCmd 表示lang ls子命令
var langLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出支持的语言",
	Long:  `列出支持的语言列表，当前选中的语言前面会有"·"标记。`,
	Run: func(cmd *cobra.Command, args []string) {
		lang.PrintLanguages()
	},
}

// langStCmd 表示lang st子命令
var langStCmd = &cobra.Command{
	Use:   "st [language]",
	Short: "切换练习语言",
	Long:  `切换当前练习的语言，例如：lang-cli lang st japanese。`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		language := args[0]
		if err := lang.SwitchLanguage(language); err != nil {
			fmt.Println(err)
		}
	},
	ValidArgs: lang.ListLanguages(),
}

// practiceCmd 表示practice子命令
var practiceCmd = &cobra.Command{
	Use:   "practice",
	Short: "练习模块",
	Long:  `练习模块，用于进行单词、短语、句子、文章等的打字练习。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("练习模块")
		// 这里将来会调用practice模块的功能
	},
}

// practiceWordsCmd 表示practice words子命令
var practiceWordsCmd = &cobra.Command{
	Use:   "words [file]",
	Short: "单词练习",
	Long:  `单词练习功能，从指定的单词列表文件中读取单词进行练习。`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 如果没有指定文件，则列出可用的单词列表文件
		if len(args) == 0 {
			files, err := practice.ListWordFiles()
			if err != nil {
				fmt.Println("获取单词列表文件失败:", err)
				return
			}

			if len(files) == 0 {
				fmt.Println("没有可用的单词列表文件，请先创建或导入单词列表。")
				return
			}

			fmt.Println("可用的单词列表文件:")
			for i, file := range files {
				fmt.Printf("%d. %s\n", i+1, file)
			}
			return
		}

		// 指定了文件，进行单词练习
		fileName := args[0]
		if err := practice.WordPractice(fileName); err != nil {
			fmt.Println("单词练习失败:", err)
		}
	},
}

// practicePhrasesCmd 表示practice phrases子命令
var practicePhrasesCmd = &cobra.Command{
	Use:   "phrases [file]",
	Short: "短语练习",
	Long:  `短语练习功能，从指定的短语列表文件中读取短语进行练习。`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 如果没有指定文件，则列出可用的短语列表文件
		if len(args) == 0 {
			files, err := practice.ListPhraseFiles()
			if err != nil {
				fmt.Println("获取短语列表文件失败:", err)
				return
			}

			if len(files) == 0 {
				fmt.Println("没有可用的短语列表文件，请先创建或导入短语列表。")
				return
			}

			fmt.Println("可用的短语列表文件:")
			for i, file := range files {
				fmt.Printf("%d. %s\n", i+1, file)
			}
			return
		}

		// 指定了文件，进行短语练习
		fileName := args[0]
		if err := practice.PhrasePractice(fileName); err != nil {
			fmt.Println("短语练习失败:", err)
		}
	},
}

// practiceSentencesCmd 表示practice sentences子命令
var practiceSentencesCmd = &cobra.Command{
	Use:   "sentences [file]",
	Short: "句子练习",
	Long:  `句子练习功能，从指定的句子列表文件中读取句子进行练习。`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 如果没有指定文件，则列出可用的句子列表文件
		if len(args) == 0 {
			files, err := practice.ListSentenceFiles()
			if err != nil {
				fmt.Println("获取句子列表文件失败:", err)
				return
			}

			if len(files) == 0 {
				fmt.Println("没有可用的句子列表文件，请先创建或导入句子列表。")
				return
			}

			fmt.Println("可用的句子列表文件:")
			for i, file := range files {
				fmt.Printf("%d. %s\n", i+1, file)
			}
			return
		}

		// 指定了文件，进行句子练习
		fileName := args[0]
		if err := practice.SentencePractice(fileName); err != nil {
			fmt.Println("句子练习失败:", err)
		}
	},
}

// practiceArticlesCmd 表示practice articles子命令
var practiceArticlesCmd = &cobra.Command{
	Use:   "articles [file]",
	Short: "文章练习",
	Long:  `文章练习功能，从指定的文章文件中读取文章进行练习。`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 如果没有指定文件，则列出可用的文章文件
		if len(args) == 0 {
			files, err := practice.ListArticleFiles()
			if err != nil {
				fmt.Println("获取文章文件失败:", err)
				return
			}

			if len(files) == 0 {
				fmt.Println("没有可用的文章文件，请先创建或导入文章。")
				return
			}

			fmt.Println("可用的文章文件:")
			for i, file := range files {
				fmt.Printf("%d. %s\n", i+1, file)
			}
			return
		}

		// 指定了文件，进行文章练习
		fileName := args[0]
		if err := practice.ArticlePractice(fileName); err != nil {
			fmt.Println("文章练习失败:", err)
		}
	},
}

// manageCmd 表示manage子命令
var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "资源管理模块",
	Long:  `资源管理模块，用于管理练习资源，包括删除和导入资源。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("资源管理模块")
		// 这里将来会调用manage模块的功能
	},
}

// manageDeleteCmd 表示manage delete子命令
var manageDeleteCmd = &cobra.Command{
	Use:   "delete [resourceType] [file]",
	Short: "删除资源",
	Long:  `删除资源，例如：lang-cli manage delete words daily.txt。`,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		// 获取资源类型
		resourceType := args[0]

		// 验证资源类型
		if !manage.ValidateResourceType(resourceType) {
			fmt.Printf("无效的资源类型: %s\n", resourceType)
			fmt.Println("有效的资源类型: words, phrases, sentences, articles")
			return
		}

		// 如果没有指定文件，则列出可用的资源文件
		if len(args) == 1 {
			files, err := manage.ListResourceFiles(resourceType)
			if err != nil {
				fmt.Printf("获取%s文件失败: %s\n", resourceType, err)
				return
			}

			if len(files) == 0 {
				fmt.Printf("没有可用的%s文件。\n", resourceType)
				return
			}

			fmt.Printf("可用的%s文件:\n", resourceType)
			for i, file := range files {
				fmt.Printf("%d. %s\n", i+1, file)
			}
			return
		}

		// 指定了文件，删除资源
		fileName := args[1]
		if err := manage.DeleteResource(resourceType, fileName); err != nil {
			fmt.Printf("删除%s文件失败: %s\n", resourceType, err)
		}
	},
}

// manageImportCmd 表示manage import子命令
var manageImportCmd = &cobra.Command{
	Use:   "import [resourceType] [file]",
	Short: "导入资源",
	Long:  `导入资源，例如：lang-cli manage import words /path/to/words.txt。`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// 获取资源类型和文件路径
		resourceType := args[0]
		filePath := args[1]

		// 验证资源类型
		if !manage.ValidateResourceType(resourceType) {
			fmt.Printf("无效的资源类型: %s\n", resourceType)
			fmt.Println("有效的资源类型: words, phrases, sentences, articles")
			return
		}

		// 导入资源
		if err := manage.ImportResource(resourceType, filePath); err != nil {
			fmt.Printf("导入%s文件失败: %s\n", resourceType, err)
		}
	},
}

func init() {
	// 加载配置文件
	if err := config.LoadConfig(); err != nil {
		fmt.Println("警告: 加载配置文件失败:", err)
		fmt.Println("将使用默认配置。")
	}

	// 添加子命令到根命令
	rootCmd.AddCommand(langCmd)
	rootCmd.AddCommand(practiceCmd)
	rootCmd.AddCommand(manageCmd)

	// 添加lang子命令
	langCmd.AddCommand(langLsCmd)
	langCmd.AddCommand(langStCmd)

	// 添加practice子命令
	practiceCmd.AddCommand(practiceWordsCmd)
	practiceCmd.AddCommand(practicePhrasesCmd)
	practiceCmd.AddCommand(practiceSentencesCmd)
	practiceCmd.AddCommand(practiceArticlesCmd)

	// 添加manage子命令
	manageCmd.AddCommand(manageDeleteCmd)
	manageCmd.AddCommand(manageImportCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
