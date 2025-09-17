# 语言学习终端 (lang-cli)

这是一个支持多种语言的打字学习终端工具，用户可以在终端中进行不同语言的单词、短语、句子、文章等的打字练习。

## 功能特点

- 支持多种语言的打字练习
- 提供单词、短语、句子、文章等多种练习模式
- 支持资源的导入和管理
- 交互式终端界面，支持上下箭头选择和Tab补全

## 安装

```bash
go install github.com/daiweiwei/lang-cli@latest
```

## 使用方法

### 查看支持的语言

```bash
lang-cli lang ls
```

### 切换语言

```bash
lang-cli lang st japanese
```

### 单词练习

```bash
lang-cli practice words
```

### 短语练习

```bash
lang-cli practice phrases
```

### 句子练习

```bash
lang-cli practice sentences
```

### 文章练习

```bash
lang-cli practice articles
```

### 删除资源

```bash
lang-cli manage delete
```

### 导入资源

```bash
lang-cli manage import
```

## 配置

配置文件位于 `config/config.yaml`，可以修改以下配置：

- 支持的语言列表
- 当前选中的语言
- 单词练习配置
- 短语练习配置
- 句子练习配置
- 文章练习配置

## 许可证

MIT