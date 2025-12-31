package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// 颜色定义
var (
	successColor = color.New(color.FgGreen).SprintFunc()
	errorColor   = color.New(color.FgRed).SprintFunc()
	warnColor    = color.New(color.FgYellow).SprintFunc()
	infoColor    = color.New(color.FgCyan).SprintFunc()
	boldColor    = color.New(color.Bold).SprintFunc()
)

// PrintSuccess 打印成功消息
func PrintSuccess(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", successColor("✅"), fmt.Sprintf(format, args...))
}

// PrintError 打印错误消息
func PrintError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s %s\n", errorColor("❌"), fmt.Sprintf(format, args...))
}

// PrintWarning 打印警告消息
func PrintWarning(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", warnColor("⚠️"), fmt.Sprintf(format, args...))
}

// PrintInfo 打印信息消息
func PrintInfo(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", infoColor("ℹ️"), fmt.Sprintf(format, args...))
}

// PrintBotTable 打印Bot列表表格
func PrintBotTable(bots []interface{}) {
	if len(bots) == 0 {
		PrintInfo("暂无Bot")
		return
	}

	// 简化表格实现
	fmt.Println("\nBot ID      名称              描述            状态      创建时间")
	fmt.Println("-------     --------------     --------------  --------  ---------")

	for range bots {
		// TODO: 根据实际Bot结构解析
		fmt.Println("bot-001     示例Bot          这是一个示例     active    2025-01-03")
	}
	fmt.Println()
}

// PrintLogs 打印日志
func PrintLogs(logs []interface{}) {
	for range logs {
		// TODO: 根据实际日志结构解析
		fmt.Println("[2025-01-03 12:00:00] [INFO] Sample log message")
	}
}

// PrintJSON 打印JSON格式输出
func PrintJSON(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))
	return nil
}

// TruncateString 截断字符串
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// FormatStatus 格式化状态显示
func FormatStatus(status string) string {
	switch strings.ToLower(status) {
	case "active", "running", "online":
		return successColor(status)
	case "inactive", "stopped", "offline":
		return color.New(color.FgHiBlack).Sprint(status)
	case "error", "failed":
		return errorColor(status)
	case "warning", "pending":
		return warnColor(status)
	default:
		return status
	}
}
