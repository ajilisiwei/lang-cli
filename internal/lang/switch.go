package lang

import (
	"fmt"

	"github.com/daiweiwei/lang-cli/internal/config"
)

// SwitchLanguage 切换当前练习的语言
func SwitchLanguage(language string) error {
	// 检查语言是否在支持的语言列表中
	languages := ListLanguages()
	var found bool
	for _, lang := range languages {
		if lang == language {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("不支持的语言: %s", language)
	}

	// 更新当前语言
	config.AppConfig.CurrentLanguage = language

	// 保存配置
	if err := config.SaveConfig(); err != nil {
		return err
	}

	fmt.Printf("已切换到语言: %s\n", language)
	return nil
}
