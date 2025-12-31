package main

import (
	"fmt"
	

	"github.com/spf13/cobra"
	"github.com/coze-dev/coze-studio/backend/cmd/zker/api"
	"github.com/coze-dev/coze-studio/backend/cmd/zker/config"
	"github.com/coze-dev/coze-studio/backend/cmd/zker/utils"
)

var (
	apiEndpoint string
)

// loginCmd 登录命令
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录ZKER平台",
	Long:  `使用API Key或JWT Token登录ZKER平台。认证信息将保存到 ~/.zker/config.json`,
	Example: `  # 使用API Key登录
  zker login --api-key sk-xxxxxxxxxxxx

  # 使用JWT Token登录
  zker login --token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

  # 指定API端点
  zker login --api-key sk-xxxx --endpoint https://api.example.com`,
	RunE: runLogin,
}

var (
	apiKey string
	token  string
)

func init() {
	loginCmd.Flags().StringVar(&apiKey, "api-key", "", "API Key (sk-开头)")
	loginCmd.Flags().StringVar(&token, "token", "", "JWT Token")
	loginCmd.Flags().StringVar(&apiEndpoint, "endpoint", "", "API端点 (默认: http://localhost:8888)")
}

func runLogin(cmd *cobra.Command, args []string) error {
	// 验证参数
	if apiKey == "" && token == "" {
		return fmt.Errorf("必须提供 --api-key 或 --token 参数")
	}

	if apiKey != "" && token != "" {
		return fmt.Errorf("不能同时使用 --api-key 和 --token，请只提供一种认证方式")
	}

	// 获取API端点
	endpoint := apiEndpoint
	if endpoint == "" {
		cfg := config.GetGlobalConfig()
		endpoint = cfg.GetAPIEndpoint()
	}

	// 创建API客户端
	client := api.NewClient(endpoint, "")

	// 验证认证信息
	var authType, authValue string
	var err error

	if apiKey != "" {
		authType = "api-key"
		authValue = apiKey
		utils.PrintInfo("正在验证API Key...")

		err = client.ValidateAPIKey(apiKey)
		if err != nil {
			utils.PrintError("API Key验证失败: %v", err)
			return err
		}
	} else {
		authType = "token"
		authValue = token
		utils.PrintInfo("正在验证Token...")

		err = client.ValidateToken(token)
		if err != nil {
			utils.PrintError("Token验证失败: %v", err)
			return err
		}
	}

	// 保存配置
	cfg := &config.Config{
		AuthType:    authType,
		AuthValue:   authValue,
		APIEndpoint: endpoint,
	}

	if err := cfg.Save(); err != nil {
		utils.PrintError("保存配置失败: %v", err)
		return err
	}

	utils.PrintSuccess("登录成功！")
	fmt.Printf("  认证方式: %s\n", authType)
	fmt.Printf("  API端点: %s\n", endpoint)
	fmt.Printf("  配置文件: ~/.zker/config.json\n")

	return nil
}
