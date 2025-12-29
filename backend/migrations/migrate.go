/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MigrationConfig 迁移配置
type MigrationConfig struct {
	DSN        string // 数据库连接字符串
	ScriptDir  string // 迁移脚本目录
	DryRun     bool   // 是否为演练模式（不执行实际操作）
	SkipBackup bool   // 是否跳过备份警告
}

func main() {
	// 解析命令行参数
	config := parseFlags()

	// 验证配置
	if err := validateConfig(config); err != nil {
		fmt.Printf("❌ 配置验证失败: %v\n", err)
		os.Exit(1)
	}

	// 连接数据库
	db, err := connectDB(config.DSN)
	if err != nil {
		fmt.Printf("❌ 数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// 确认执行
	if !config.SkipBackup {
		fmt.Println("⚠️  警告：此操作将修改数据库结构！")
		fmt.Println("⚠️  请确保已备份数据库！")
		fmt.Print("是否继续？(yes/no): ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input != "yes" && input != "y" {
			fmt.Println("❌ 操作已取消")
			os.Exit(0)
		}
	}

	// 执行迁移
	if err := executeMigration(db, config); err != nil {
		fmt.Printf("❌ 迁移执行失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ 迁移执行成功！")
}

// parseFlags 解析命令行参数
func parseFlags() MigrationConfig {
	var config MigrationConfig

	flag.StringVar(&config.DSN, "dsn", "", "数据库连接字符串 (例: user:password@tcp(localhost:3306)/database)")
	flag.StringVar(&config.ScriptDir, "dir", "./migrations", "迁移脚本目录")
	flag.BoolVar(&config.DryRun, "dry-run", false, "演练模式（不执行实际操作）")
	flag.BoolVar(&config.SkipBackup, "skip-backup-warning", false, "跳过备份警告")

	flag.Parse()

	// 从环境变量读取 DSN（如果未指定）
	if config.DSN == "" {
		config.DSN = os.Getenv("MYSQL_DSN")
	}

	// 设置默认脚本目录
	if config.ScriptDir == "" {
		config.ScriptDir = "."
	}

	return config
}

// validateConfig 验证配置
func validateConfig(config MigrationConfig) error {
	if config.DSN == "" {
		return fmt.Errorf("未指定数据库连接字符串，请使用 -dsn 参数或设置 MYSQL_DSN 环境变量")
	}

	// 检查脚本目录是否存在
	if _, err := os.Stat(config.ScriptDir); os.IsNotExist(err) {
		return fmt.Errorf("脚本目录不存在: %s", config.ScriptDir)
	}

	return nil
}

// connectDB 连接数据库
func connectDB(dsn string) (*sqlx.DB, error) {
	fmt.Println("🔌 正在连接数据库...")

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接失败: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Ping 失败: %w", err)
	}

	fmt.Println("✅ 数据库连接成功")
	return db, nil
}

// executeMigration 执行迁移
func executeMigration(db *sqlx.DB, config MigrationConfig) error {
	// 查找所有 .sql 文件
	sqlFiles, err := findSQLFiles(config.ScriptDir)
	if err != nil {
		return fmt.Errorf("查找 SQL 文件失败: %w", err)
	}

	if len(sqlFiles) == 0 {
		return fmt.Errorf("未找到 SQL 迁移脚本")
	}

	fmt.Printf("📝 找到 %d 个迁移脚本\n", len(sqlFiles))

	// 执行每个脚本
	for i, sqlFile := range sqlFiles {
		fmt.Printf("\n[%d/%d] 正在执行: %s\n", i+1, len(sqlFiles), filepath.Base(sqlFile))

		if err := executeScript(db, sqlFile, config.DryRun); err != nil {
			return fmt.Errorf("执行脚本失败 %s: %w", sqlFile, err)
		}

		fmt.Printf("✅ [%d/%d] 完成: %s\n", i+1, len(sqlFiles), filepath.Base(sqlFile))
	}

	return nil
}

// findSQLFiles 查找所有 SQL 文件
func findSQLFiles(dir string) ([]string, error) {
	var sqlFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 只处理 .sql 文件
		if strings.HasSuffix(strings.ToLower(path), ".sql") {
			sqlFiles = append(sqlFiles, path)
		}

		return nil
	})

	return sqlFiles, err
}

// executeScript 执行单个 SQL 脚本
func executeScript(db *sqlx.DB, scriptPath string, dryRun bool) error {
	// 读取脚本内容
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	script := string(content)

	// 如果是演练模式，只打印不执行
	if dryRun {
		fmt.Printf("🔍 [DRY RUN] 将执行以下 SQL:\n")
		fmt.Printf("   文件: %s\n", scriptPath)
		fmt.Printf("   大小: %d 字节\n", len(content))
		return nil
	}

	// 开始事务
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}
	defer tx.Rollback() // 确保失败时回滚

	// 分割SQL语句并执行
	statements := splitSQL(script)

	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		fmt.Printf("   执行语句 %d/%d...\n", i+1, len(statements))

		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("执行语句失败 [%d]: %w\nSQL: %s", i+1, err, stmt)
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// splitSQL 分割SQL语句（简单实现）
func splitSQL(sql string) []string {
	// 按分号分割
	var statements []string
	var current strings.Builder
	inComment := false

	lines := strings.Split(sql, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 跳过注释行
		if strings.HasPrefix(trimmed, "--") {
			continue
		}

		// 检查多行注释
		if strings.HasPrefix(trimmed, "/*") {
			inComment = true
		}
		if strings.HasSuffix(trimmed, "*/") {
			inComment = false
			continue
		}
		if inComment {
			continue
		}

		// 累积当前语句
		current.WriteString(line)
		current.WriteString("\n")

		// 如果以分号结尾，表示一个完整语句
		if strings.HasSuffix(trimmed, ";") {
			statements = append(statements, current.String())
			current.Reset()
		}
	}

	// 处理最后一个语句（可能没有分号）
	if current.Len() > 0 {
		statements = append(statements, current.String())
	}

	return statements
}
