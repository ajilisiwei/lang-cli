package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 表示应用程序的配置
type Config struct {
	// 支持的语言列表
	Languages []string `mapstructure:"languages"`
	// 当前选中的语言
	CurrentLanguage string `mapstructure:"current_language"`
	// 单词练习配置
	Words WordsConfig `mapstructure:"words"`
	// 短语练习配置
	Phrases PhrasesConfig `mapstructure:"phrases"`
	// 句子练习配置
	Sentences SentencesConfig `mapstructure:"sentences"`
	// 文章练习配置
	Articles ArticlesConfig `mapstructure:"articles"`
	// v0.2 新增：正确性匹配模式，可选值：exact_match（完全匹配）、word_match（单词匹配）
	CorrectnessMatchMode string `mapstructure:"correctness_match_mode"`
	// v0.2 新增：全局下一个资源的出现顺序，可选值：random（随机）、sequential（顺序）
	NextOneOrder string `mapstructure:"next_one_order"`
}

// WordsConfig 表示单词练习的配置
type WordsConfig struct {
	// 下一个单词的出现顺序，可选值：random（随机）、sequential（顺序）
	NextOneOrder string `mapstructure:"next_one_order"`
	// 输入正确后是否显示翻译
	ShowTranslation bool `mapstructure:"show_translation"`
}

// PhrasesConfig 表示短语练习的配置
type PhrasesConfig struct {
	// 下一个短语的出现顺序，可选值：random（随机）、sequential（顺序）
	NextOneOrder string `mapstructure:"next_one_order"`
	// 输入正确后是否显示翻译
	ShowTranslation bool `mapstructure:"show_translation"`
}

// SentencesConfig 表示句子练习的配置
type SentencesConfig struct {
	// 下一个句子的出现顺序，可选值：random（随机）、sequential（顺序）
	NextOneOrder string `mapstructure:"next_one_order"`
	// 输入正确后是否显示翻译
	ShowTranslation bool `mapstructure:"show_translation"`
}

// ArticlesConfig 表示文章练习的配置
type ArticlesConfig struct {
	// 输入正确后是否显示翻译
	ShowTranslation bool `mapstructure:"show_translation"`
}

// 全局配置实例
var AppConfig Config

// LoadConfig 加载配置文件
// LoadConfig 加载配置文件
func LoadConfig() error {
	// 设置配置文件的名称和路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 获取用户主目录下的配置路径
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configDir := filepath.Join(homeDir, ".lang-cli")
		viper.AddConfigPath(configDir)
	}

	// 获取当前执行文件的目录
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取执行文件路径失败: %w", err)
	}
	execDir := filepath.Dir(execPath)

	// 添加配置文件的搜索路径
	viper.AddConfigPath(filepath.Join(execDir, "config"))
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	// 添加测试环境下的配置文件路径
	viper.AddConfigPath("../../config")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 将配置文件的内容解析到结构体中
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}

// SaveConfig 保存配置到文件
func SaveConfig() error {
	// 将结构体的内容写入到配置文件中
	for k, v := range map[string]interface{}{
		"languages":               AppConfig.Languages,
		"current_language":        AppConfig.CurrentLanguage,
		"words":                   AppConfig.Words,
		"phrases":                 AppConfig.Phrases,
		"sentences":               AppConfig.Sentences,
		"articles":                AppConfig.Articles,
		"correctness_match_mode":  AppConfig.CorrectnessMatchMode,
		"next_one_order":          AppConfig.NextOneOrder,
	} {
		viper.Set(k, v)
	}

	// 保存配置文件
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("保存配置文件失败: %w", err)
	}

	return nil
}
