package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/coze-dev/coze-studio/backend/cmd/zker/api"
	"github.com/coze-dev/coze-studio/backend/cmd/zker/config"
	"github.com/coze-dev/coze-studio/backend/cmd/zker/utils"
)

// botCmd Bot管理命令
var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "Bot管理",
	Long:  `管理ZKER Bot，包括创建、部署、查看和删除操作`,
}

var (
	botTemplate    string
	botDescription string
	deployEnv      string
)

// listBotsCmd 列出所有Bot
var listBotsCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有Bot",
	Long:  `列出当前账户下的所有Bot`,
	Example: `  # 列出所有Bot
  zker bot list

  # 使用JSON格式输出
  zker bot list --output json`,
	RunE: runBotList,
}

var outputFormat string

// createBotCmd 创建Bot
var createBotCmd = &cobra.Command{
	Use:   "create [bot-name]",
	Short: "创建Bot",
	Long:  `创建一个新的Bot`,
	Example: `  # 创建聊天机器人
  zker bot create my-bot --template chatbot

  # 创建助手
  zker bot create assistant --template assistant --description "我的助手"

  # 创建工作流Bot
  zker bot create workflow-bot --template workflow`,
	Args: cobra.ExactArgs(1),
	RunE: runBotCreate,
}

// deployBotCmd 部署Bot
var deployBotCmd = &cobra.Command{
	Use:   "deploy [bot-name]",
	Short: "部署Bot",
	Long:  `部署Bot到指定环境`,
	Example: `  # 部署到开发环境
  zker bot deploy my-bot --env dev

  # 部署到生产环境
  zker bot deploy my-bot --env production`,
	Args: cobra.ExactArgs(1),
	RunE: runBotDeploy,
}

// infoBotCmd 查看Bot详情
var infoBotCmd = &cobra.Command{
	Use:   "info [bot-name]",
	Short: "查看Bot详情",
	Long:  `查看Bot的详细信息`,
	Example: `  zker bot info my-bot`,
	Args: cobra.ExactArgs(1),
	RunE: runBotInfo,
}

// deleteBotCmd 删除Bot
var deleteBotCmd = &cobra.Command{
	Use:   "delete [bot-name]",
	Short: "删除Bot",
	Long:  `删除指定的Bot（此操作不可恢复）`,
	Example: `  zker bot delete my-bot`,
	Args: cobra.ExactArgs(1),
	RunE: runBotDelete,
}

func init() {
	botCmd.AddCommand(listBotsCmd)
	botCmd.AddCommand(createBotCmd)
	botCmd.AddCommand(deployBotCmd)
	botCmd.AddCommand(infoBotCmd)
	botCmd.AddCommand(deleteBotCmd)

	// list命令参数
	listBotsCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "输出格式 (table/json)")

	// create命令参数
	createBotCmd.Flags().StringVar(&botTemplate, "template", "chatbot", "Bot模板 (chatbot/assistant/workflow)")
	createBotCmd.Flags().StringVarP(&botDescription, "description", "d", "", "Bot描述")

	// deploy命令参数
	deployBotCmd.Flags().StringVarP(&deployEnv, "env", "e", "dev", "部署环境 (dev/staging/production)")
}

// getAPIClient 获取API客户端
func getAPIClient() *api.Client {
	cfg := config.GetGlobalConfig()
	if !cfg.IsAuthenticated() {
		utils.PrintError("未登录，请先运行: zker login")
		return nil
	}

	return api.NewClient(cfg.GetAPIEndpoint(), cfg.GetAuthToken())
}

func runBotList(cmd *cobra.Command, args []string) error {
	client := getAPIClient()
	if client == nil {
		return fmt.Errorf("获取API客户端失败")
	}

	utils.PrintInfo("正在获取Bot列表...")

	bots, err := client.ListBots()
	if err != nil {
		utils.PrintError("获取Bot列表失败: %v", err)
		return err
	}

	if outputFormat == "json" {
		return utils.PrintJSON(bots)
	}

	utils.PrintSuccess("共找到 %d 个Bot\n", len(bots))
	utils.PrintBotTable(convertBots(bots))
	return nil
}

func runBotCreate(cmd *cobra.Command, args []string) error {
	botName := args[0]

	client := getAPIClient()
	if client == nil {
		return fmt.Errorf("获取API客户端失败")
	}

	utils.PrintInfo("正在创建Bot: %s", botName)

	req := &api.CreateBotRequest{
		Name:        botName,
		Template:    botTemplate,
		Description: botDescription,
	}

	bot, err := client.CreateBot(req)
	if err != nil {
		utils.PrintError("创建Bot失败: %v", err)
		return err
	}

	utils.PrintSuccess("Bot创建成功！")
	fmt.Printf("  Bot ID: %s\n", bot.ID)
	fmt.Printf("  Bot名称: %s\n", bot.Name)
	fmt.Printf("  模板: %s\n", botTemplate)
	if bot.Description != "" {
		fmt.Printf("  描述: %s\n", bot.Description)
	}

	fmt.Printf("\n下一步:\n")
	fmt.Printf("  部署Bot: zker bot deploy %s --env dev\n", botName)
	fmt.Printf("  查看日志: zker logs %s --follow\n", botName)

	return nil
}

func runBotDeploy(cmd *cobra.Command, args []string) error {
	botName := args[0]

	client := getAPIClient()
	if client == nil {
		return fmt.Errorf("获取API客户端失败")
	}

	utils.PrintInfo("正在部署Bot: %s (环境: %s)", botName, deployEnv)

	// TODO: 通过botName获取botID
	botID := botName // 临时使用botName作为botID

	deployment, err := client.DeployBot(botID, deployEnv)
	if err != nil {
		utils.PrintError("部署Bot失败: %v", err)
		return err
	}

	utils.PrintSuccess("Bot部署成功！")
	fmt.Printf("  Bot名称: %s\n", botName)
	fmt.Printf("  环境: %s\n", deployEnv)
	fmt.Printf("  部署ID: %s\n", deployment.ID)
	fmt.Printf("  状态: %s\n", deployment.Status)
	if deployment.URL != "" {
		fmt.Printf("  访问地址: %s\n", deployment.URL)
	}

	fmt.Printf("\n查看部署状态:\n")
	fmt.Printf("  zker bot info %s\n", botName)

	return nil
}

func runBotInfo(cmd *cobra.Command, args []string) error {
	botName := args[0]

	client := getAPIClient()
	if client == nil {
		return fmt.Errorf("获取API客户端失败")
	}

	utils.PrintInfo("正在获取Bot信息: %s", botName)

	// TODO: 通过botName获取botID
	botID := botName

	bot, err := client.GetBot(botID)
	if err != nil {
		utils.PrintError("获取Bot信息失败: %v", err)
		return err
	}

	utils.PrintSuccess("Bot信息:")
	fmt.Printf("  ID: %s\n", bot.ID)
	fmt.Printf("  名称: %s\n", bot.Name)
	fmt.Printf("  描述: %s\n", bot.Description)
	fmt.Printf("  状态: %s\n", bot.Status)
	fmt.Printf("  创建时间: %s\n", bot.CreatedAt)
	fmt.Printf("  更新时间: %s\n", bot.UpdatedAt)

	return nil
}

func runBotDelete(cmd *cobra.Command, args []string) error {
	botName := args[0]

	// 确认删除
	fmt.Printf("警告: 您即将删除Bot '%s'，此操作不可恢复！\n", botName)
	fmt.Printf("是否继续? (yes/no): ")

	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "yes" && confirm != "y" {
		utils.PrintInfo("操作已取消")
		return nil
	}

	client := getAPIClient()
	if client == nil {
		return fmt.Errorf("获取API客户端失败")
	}

	utils.PrintInfo("正在删除Bot: %s", botName)

	// TODO: 通过botName获取botID
	botID := botName

	err := client.DeleteBot(botID)
	if err != nil {
		utils.PrintError("删除Bot失败: %v", err)
		return err
	}

	utils.PrintSuccess("Bot '%s' 已删除", botName)
	return nil
}

// convertBots 转换Bot切片为interface{}切片
func convertBots(bots []api.Bot) []interface{} {
	result := make([]interface{}, len(bots))
	for i, bot := range bots {
		result[i] = bot
	}
	return result
}
