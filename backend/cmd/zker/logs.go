package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/coze-dev/coze-studio/backend/cmd/zker/api"
	"github.com/coze-dev/coze-studio/backend/cmd/zker/config"
	"github.com/coze-dev/coze-studio/backend/cmd/zker/utils"
)

// logsCmd 日志查询命令
var logsCmd = &cobra.Command{
	Use:   "logs [bot-name]",
	Short: "查看日志",
	Long:  `查看Bot的运行日志`,
	Example: `  # 实时查看日志
  zker logs my-bot --follow

  # 查看最近100行
  zker logs my-bot --tail 100

  # 过滤错误日志
  zker logs my-bot --level error

  # 查看最近1小时的日志
  zker logs my-bot --since "1h"`,
	Args: cobra.ExactArgs(1),
	RunE: runLogs,
}

var (
	logFollow bool
	logTail   int
	logLevel  string
	logSince  string
)

func init() {
	logsCmd.Flags().BoolVarP(&logFollow, "follow", "f", false, "实时跟踪日志（类似tail -f）")
	logsCmd.Flags().IntVarP(&logTail, "tail", "n", 100, "显示最后N行日志")
	logsCmd.Flags().StringVarP(&logLevel, "level", "l", "", "过滤日志级别 (debug/info/warn/error)")
	logsCmd.Flags().StringVarP(&logSince, "since", "s", "", "时间范围 (e.g., '1h', '30m', '1d')")
}

func runLogs(cmd *cobra.Command, args []string) error {
	botName := args[0]

	client := getLogClient()
	if client == nil {
		return fmt.Errorf("获取API客户端失败")
	}

	// 解析时间范围
	var sinceTime time.Time
	var err error

	if logSince != "" {
		sinceTime, err = parseSince(logSince)
		if err != nil {
			return fmt.Errorf("解析时间范围失败: %w", err)
		}
	}

	if logFollow {
		// 实时日志流
		utils.PrintInfo("正在实时跟踪日志 (Ctrl+C 退出)...\n")
		return streamLogs(client, botName, logLevel)
	}

	// 查询历史日志
	utils.PrintInfo("正在获取日志...")

	query := &api.LogQuery{
		Tail:   logTail,
		Level:  logLevel,
		Since:  sinceTime,
		Follow: false,
	}

	logs, err := client.GetLogs(botName, query)
	if err != nil {
		utils.PrintError("获取日志失败: %v", err)
		return err
	}

	if len(logs) == 0 {
		utils.PrintInfo("暂无日志")
		return nil
	}

	utils.PrintSuccess("共获取 %d 条日志\n", len(logs))
	utils.PrintLogs(convertLogs(logs))

	return nil
}

// getLogClient 获取日志API客户端
func getLogClient() *api.Client {
	cfg := config.GetGlobalConfig()
	if !cfg.IsAuthenticated() {
		utils.PrintError("未登录，请先运行: zker login")
		return nil
	}

	return api.NewClient(cfg.GetAPIEndpoint(), cfg.GetAuthToken())
}

// streamLogs 实时流式日志
func streamLogs(client *api.Client, botName, level string) error {
	err := client.StreamLogs(botName, level, func(entry api.LogEntry) {
		// 格式化输出日志
		fmt.Printf("[%s] [%s] %s\n",
			entry.Timestamp,
			entry.Level,
			entry.Message,
		)
	})

	if err != nil {
		utils.PrintError("流式日志失败: %v", err)
		return err
	}

	return nil
}

// parseSince 解析时间范围字符串
func parseSince(since string) (time.Time, error) {
	now := time.Now()

	var duration time.Duration
	var err error

	// 解析时间范围
	switch {
	case len(since) > 0:
		// 尝试解析为duration
		duration, err = time.ParseDuration(since)
		if err != nil {
			return time.Time{}, fmt.Errorf("无效的时间格式: %s (支持格式: 1h, 30m, 1d)", since)
		}
	default:
		return time.Time{}, fmt.Errorf("时间范围不能为空")
	}

	return now.Add(-duration), nil
}

// convertLogs 转换Log切片为interface{}切片
func convertLogs(logs []api.LogEntry) []interface{} {
	result := make([]interface{}, len(logs))
	for i, log := range logs {
		result[i] = log
	}
	return result
}
