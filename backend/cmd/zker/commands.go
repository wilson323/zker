package main

import (
	
)

// addCommands 添加所有子命令
func addCommands() {
	// 添加所有命令
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(botCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(workflowCmd)
	rootCmd.AddCommand(knowledgeCmd)
}
