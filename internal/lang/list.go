package lang

import (
	"fmt"

	"github.com/ajilisiwei/mllt-cli/internal/config"
)

// ListLanguages 列出支持的语言列表
func ListLanguages() []string {
	return config.AppConfig.Languages
}

// PrintLanguages 打印支持的语言列表
func PrintLanguages() {
	languages := ListLanguages()
	currentLanguage := config.AppConfig.CurrentLanguage

	fmt.Println("支持的语言列表:")
	for _, lang := range languages {
		prefix := "  "
		if lang == currentLanguage {
			prefix = "✔ "
		}
		fmt.Printf("%s%s\n", prefix, lang)
	}
}

// GetCurrentLanguage 获取当前设置的语言
func GetCurrentLanguage() string {
	return config.AppConfig.CurrentLanguage
}
