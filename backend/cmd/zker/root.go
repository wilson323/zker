package main

import (
	"github.com/spf13/cobra"
)

// AddCommands 添加所有子命令
func AddCommands(rootCmd *cobra.Command) {
	// 添加所有命令
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(botCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(workflowCmd)
	rootCmd.AddCommand(knowledgeCmd)
}
