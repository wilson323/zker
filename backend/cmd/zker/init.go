package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/coze-dev/coze-studio/backend/cmd/zker/utils"
)

// initCmd 初始化项目命令
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "初始化Bot项目",
	Long:  `创建一个新的ZKER Bot项目，包含标准的目录结构和配置文件`,
	Example: `  # 创建新项目
  zker init my-bot

  # 进入项目目录
  cd my-bot`,
	Args: cobra.ExactArgs(1),
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	// 检查目录是否已存在
	if _, err := os.Stat(projectName); err == nil {
		return fmt.Errorf("目录 '%s' 已存在", projectName)
	}

	utils.PrintInfo("正在创建项目: %s", projectName)

	// 创建项目目录
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("创建项目目录失败: %w", err)
	}

	// 创建目录结构
	dirs := []string{
		"bots",
		"workflows",
		"knowledge",
		"config",
		"logs",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(projectName, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
	}

	// 创建配置文件
	if err := createZkerConfig(projectName); err != nil {
		return err
	}

	// 创建示例Bot配置
	if err := createSampleBot(projectName); err != nil {
		return err
	}

	// 创建.gitignore
	if err := createGitignore(projectName); err != nil {
		return err
	}

	// 创建README
	if err := createReadme(projectName); err != nil {
		return err
	}

	utils.PrintSuccess("项目初始化成功！")
	fmt.Printf("\n项目名称: %s\n", projectName)
	fmt.Printf("\n目录结构:\n")
	fmt.Printf("  %s/\n", projectName)
	fmt.Printf("  ├── bots/          # Bot配置文件\n")
	fmt.Printf("  ├── workflows/     # 工作流配置\n")
	fmt.Printf("  ├── knowledge/     # 知识库文件\n")
	fmt.Printf("  ├── config/        # 配置文件\n")
	fmt.Printf("  ├── logs/          # 日志文件\n")
	fmt.Printf("  ├── zker.yaml      # 项目配置\n")
	fmt.Printf("  ├── .gitignore\n")
	fmt.Printf("  └── README.md\n")
	fmt.Printf("\n下一步:\n")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Printf("  zker login --api-key sk-your-api-key\n")
	fmt.Printf("  zker bot create my-first-bot\n")

	return nil
}

// createZkerConfig 创建zker.yaml配置文件
func createZkerConfig(projectName string) error {
	configContent := "# ZKER Bot项目配置\n" +
		fmt.Sprintf("project: %s\n", projectName) +
		"version: 1.0.0\n\n" +
		"# Bot配置目录\n" +
		"bots:\n" +
		"  directory: ./bots\n\n" +
		"# 工作流配置目录\n" +
		"workflows:\n" +
		"  directory: ./workflows\n\n" +
		"# 知识库目录\n" +
		"knowledge:\n" +
		"  directory: ./knowledge\n\n" +
		"# 日志目录\n" +
		"logs:\n" +
		"  directory: ./logs\n\n" +
		"# API配置\n" +
		"api:\n" +
		"  endpoint: http://localhost:8888\n" +
		"  timeout: 30s\n\n" +
		"# 部署配置\n" +
		"deployment:\n" +
		"  default_env: dev\n" +
		"  environments:\n" +
		"    - name: dev\n" +
		"      endpoint: http://localhost:8888\n" +
		"    - name: staging\n" +
		"      endpoint: https://staging.api.example.com\n" +
		"    - name: production\n" +
		"      endpoint: https://api.example.com\n"

	configPath := filepath.Join(projectName, "zker.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("创建配置文件失败: %w", err)
	}

	return nil
}

// createSampleBot 创建示例Bot配置
func createSampleBot(projectName string) error {
	sampleBot := "# 示例Bot配置\n" +
		"name: hello-bot\n" +
		"version: 1.0.0\n" +
		"description: 我的第一个ZKER Bot\n\n" +
		"# Bot类型: chatbot, assistant, workflow\n" +
		"type: chatbot\n\n" +
		"# LLM配置\n" +
		"llm:\n" +
		"  provider: openai\n" +
		"  model: gpt-4\n" +
		"  temperature: 0.7\n" +
		"  max_tokens: 2000\n\n" +
		"# 提示词\n" +
		"prompt: |\n" +
		"  你是一个友好的助手，可以帮助用户解答问题。\n" +
		"  请用简洁、专业的语言回答。\n\n" +
		"# 知识库配置\n" +
		"knowledge:\n" +
		"  enabled: false\n" +
		"  bases: []\n\n" +
		"# 工具配置\n" +
		"tools:\n" +
		"  - name: search\n" +
		"    enabled: false\n"

	botPath := filepath.Join(projectName, "bots", "hello-bot.yaml")
	if err := os.WriteFile(botPath, []byte(sampleBot), 0644); err != nil {
		return fmt.Errorf("创建示例Bot失败: %w", err)
	}

	return nil
}

// createGitignore 创建.gitignore文件
func createGitignore(projectName string) error {
	gitignore := "# ZKER CLI\n" +
		".zker/\n\n" +
		"# 日志文件\n" +
		"logs/\n" +
		"*.log\n\n" +
		"# IDE\n" +
		".idea/\n" +
		".vscode/\n" +
		"*.swp\n" +
		"*.swo\n" +
		"*~\n\n" +
		"# OS\n" +
		".DS_Store\n" +
		"Thumbs.db\n\n" +
		"# 配置文件（包含敏感信息）\n" +
		"config/secrets.yaml\n"

	gitignorePath := filepath.Join(projectName, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignore), 0644); err != nil {
		return fmt.Errorf("创建.gitignore失败: %w", err)
	}

	return nil
}

// createReadme 创建README.md文件
func createReadme(projectName string) error {
	readmeContent := fmt.Sprintf("# %s\n\n", projectName) +
		"这是一个使用ZKER CLI创建的AI Bot项目。\n\n" +
		"## 快速开始\n\n" +
		"### 1. 安装ZKER CLI\n\n" +
		"从源码安装: `go install github.com/coze-dev/coze-studio/backend/cmd/zker@latest`\n\n" +
		"### 2. 登录ZKER平台\n\n" +
		"`zker login --api-key sk-your-api-key`\n\n" +
		"### 3. 创建Bot\n\n" +
		"- 创建新Bot: `zker bot create my-bot --template chatbot`\n" +
		"- 部署Bot: `zker bot deploy my-bot --env dev`\n" +
		"- 查看日志: `zker logs my-bot --follow`\n\n" +
		"## 项目结构\n\n" +
		"- `bots/` - Bot配置文件\n" +
		"- `workflows/` - 工作流配置\n" +
		"- `knowledge/` - 知识库文件\n" +
		"- `config/` - 配置文件\n" +
		"- `logs/` - 日志文件\n\n" +
		"## 更多信息\n\n" +
		"查看 ZKER文档 了解更多。\n"

	readmePath := filepath.Join(projectName, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("创建README失败: %w", err)
	}

	return nil
}
