package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/coze-dev/coze-studio/backend/cmd/zker/config"
	
)

var (
	// Version 版本号
	Version = "1.0.0"
	// BuildTime 构建时间
	BuildTime = "unknown"
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "zker",
	Short: "ZKER AI智能体工作台 - CLI工具",
	Long: `ZKER CLI - 企业级AI智能体开发平台命令行工具

快速开始:
  zker init my-bot      # 初始化Bot项目
  zker login            # 登录账户
  zker bot list         # 列出所有Bot
  zker bot deploy       # 部署Bot
  zker logs my-bot      # 查看日志

文档:
  https://github.com/coze-dev/coze-studio/docs`,

	Version: "1.0.0",
}

func main() {
	// 设置版本信息
	rootCmd.Version = fmt.Sprintf("%s (built: %s)", Version, BuildTime)

	// 添加子命令
	addCommands()

	// 初始化配置
	cobra.OnInitialize(initConfig)

	// 禁用默认的completion命令
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func initConfig() {
	// 初始化全局配置
	if err := config.InitGlobalConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize config: %v\n", err)
	}
}
