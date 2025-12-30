// backend/cmd/migrate/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/migration"
	"github.com/coze-dev/coze-studio/backend/infra/database"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// 表迁移配置
type TableConfig struct {
	TableName    string // 表名
	UserIDField  string // 用户ID字段名
	Description  string // 描述
}

// 所有需要迁移的表
var tablesToMigrate = []TableConfig{
	{"bots", "creator_id", "Bot表"},
	{"bot_configs", "", "Bot配置表（通过bot_id关联）"},
	{"conversations", "creator_id", "对话表"},
	{"messages", "", "消息表（通过conversation_id关联）"},
	{"knowledge_bases", "creator_id", "知识库表"},
	{"knowledge_chunks", "", "知识库分块表（通过kb_id关联）"},
	{"workflows", "creator_id", "工作流表"},
	{"workflow_executions", "", "工作流执行记录表（通过workflow_id关联）"},
	{"single_agent_draft", "user_id", "单Agent草稿表"},
	{"published_bots", "", "已发布Bot表（通过bot_id关联）"},
}

func main() {
	// 命令行参数
	command := flag.String("cmd", "help", "命令: start, validate, rollback, stats")
	table := flag.String("table", "", "表名（不指定则迁移所有表）")
	batchSize := flag.Int("batch", 1000, "批次大小")
	flag.Parse()

	// 初始化日志
	logger.Init("migration")

	// 初始化数据库
	db, err := database.InitDB()
	if err != nil {
		logger.Errorf("Failed to init database: %v", err)
		os.Exit(1)
	}

	// 创建迁移器
	migrator := migration.NewTenantIDMigrator(db)

	ctx := context.Background()

	// 执行命令
	switch *command {
	case "start":
		err = doMigrate(ctx, migrator, *table)
	case "validate":
		err = doValidate(ctx, migrator, *table)
	case "rollback":
		err = doRollback(ctx, migrator, *table)
	case "stats":
		err = doStats(ctx, migrator, *table)
	default:
		printHelp()
	}

	if err != nil {
		logger.Errorf("Command failed: %v", err)
		os.Exit(1)
	}
}

// doMigrate 执行迁移
func doMigrate(ctx context.Context, migrator *migration.TenantIDMigrator, tableName string) error {
	logger.Infof("=== Starting tenant_id migration ===")

	var tables []TableConfig
	if tableName != "" {
		// 迁移单张表
		found := false
		for _, t := range tablesToMigrate {
			if t.TableName == tableName {
				tables = []TableConfig{t}
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("table %s not found in migration list", tableName)
		}
	} else {
		// 迁移所有表
		tables = tablesToMigrate
	}

	// 逐表迁移
	for i, tableConfig := range tables {
		logger.Infof("========================================")
		logger.Infof("[%d/%d] Migrating table: %s (%s)",
			i+1, len(tables), tableConfig.TableName, tableConfig.Description)
		logger.Infof("========================================")

		if tableConfig.UserIDField == "" {
			logger.Warnf("Table %s has no user_id_field, skipping direct migration",
				tableConfig.TableName)
			// 这些表通过关联表迁移，不需要直接处理
			continue
		}

		startTime := time.Now()
		err := migrator.MigrateTable(ctx, tableConfig.TableName, tableConfig.UserIDField)
		duration := time.Since(startTime)

		if err != nil {
			logger.Errorf("Failed to migrate table %s: %v (duration: %s)",
				tableConfig.TableName, err, duration)
			// 继续迁移下一张表，不中断
			continue
		}

		logger.Infof("✓ Table %s migration completed successfully (duration: %s)",
			tableConfig.TableName, duration)

		// 验证迁移结果
		err = migrator.ValidateMigration(ctx, tableConfig.TableName)
		if err != nil {
			logger.Errorf("Validation failed for table %s: %v",
				tableConfig.TableName, err)
		} else {
			logger.Infof("✓ Validation passed for table %s", tableConfig.TableName)
		}
	}

	logger.Infof("=== Migration completed ===")
	return nil
}

// doValidate 验证迁移结果
func doValidate(ctx context.Context, migrator *migration.TenantIDMigrator, tableName string) error {
	logger.Infof("=== Validating tenant_id migration ===")

	var tables []TableConfig
	if tableName != "" {
		for _, t := range tablesToMigrate {
			if t.TableName == tableName {
				tables = []TableConfig{t}
				break
			}
		}
	} else {
		tables = tablesToMigrate
	}

	validationPassed := true
	for _, tableConfig := range tables {
		logger.Infof("Validating table: %s...", tableConfig.TableName)

		err := migrator.ValidateMigration(ctx, tableConfig.TableName)
		if err != nil {
			logger.Errorf("✗ Table %s validation failed: %v",
				tableConfig.TableName, err)
			validationPassed = false
		} else {
			logger.Infof("✓ Table %s validation passed", tableConfig.TableName)
		}
	}

	if validationPassed {
		logger.Infof("=== All validations passed ===")
	} else {
		logger.Errorf("=== Some validations failed ===")
	}

	return nil
}

// doRollback 回滚迁移
func doRollback(ctx context.Context, migrator *migration.TenantIDMigrator, tableName string) error {
	logger.Warnf("=== Rolling back tenant_id migration ===")

	fmt.Print("⚠️  WARNING: This will clear all migrated tenant_id data!\n")
	fmt.Print("Are you sure? (type 'yes' to continue): ")

	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "yes" {
		logger.Infof("Rollback cancelled")
		return nil
	}

	var tables []TableConfig
	if tableName != "" {
		for _, t := range tablesToMigrate {
			if t.TableName == tableName {
				tables = []TableConfig{t}
				break
			}
		}
	} else {
		tables = tablesToMigrate
	}

	for _, tableConfig := range tables {
		logger.Infof("Rolling back table: %s...", tableConfig.TableName)

		err := migrator.RollbackMigration(ctx, tableConfig.TableName)
		if err != nil {
			logger.Errorf("Failed to rollback table %s: %v",
				tableConfig.TableName, err)
			continue
		}

		logger.Infof("✓ Table %s rollback completed", tableConfig.TableName)
	}

	logger.Infof("=== Rollback completed ===")
	return nil
}

// doStats 显示迁移统计
func doStats(ctx context.Context, migrator *migration.TenantIDMigrator, tableName string) error {
	logger.Infof("=== Migration Statistics ===")

	var tables []TableConfig
	if tableName != "" {
		for _, t := range tablesToMigrate {
			if t.TableName == tableName {
				tables = []TableConfig{t}
				break
			}
		}
	} else {
		tables = tablesToMigrate
	}

	for _, tableConfig := range tables {
		stats, err := migrator.GetMigrationStats(ctx, tableConfig.TableName)
		if err != nil {
			logger.Errorf("Failed to get stats for table %s: %v",
				tableConfig.TableName, err)
			continue
		}

		fmt.Printf("\n--- Table: %s (%s) ---\n", tableConfig.TableName, tableConfig.Description)
		fmt.Printf("Total Records:      %d\n", stats.TotalRecords)
		fmt.Printf("Migrated Records:   %d\n", stats.MigratedRecords)
		fmt.Printf("Unmigrated Records: %d\n", stats.UnmigratedRecords)
		fmt.Printf("Invalid Records:    %d\n", stats.InvalidRecords)
		fmt.Printf("Progress:           %.2f%%\n", stats.Progress)
		fmt.Printf("\n")
	}

	return nil
}

// printHelp 打印帮助信息
func printHelp() {
	fmt.Printf(`
tenant_id Migration Tool

Usage:
  migrate -cmd <command> [options]

Commands:
  start      Start migration (migrate all tables or specific table)
  validate   Validate migration results
  rollback   Rollback migration (clear tenant_id data)
  stats      Show migration statistics

Options:
  -table     Specify table name (default: all tables)
  -batch     Batch size for migration (default: 1000)

Examples:
  # Migrate all tables
  migrate -cmd=start

  # Migrate specific table
  migrate -cmd=start -table=bots

  # Validate migration
  migrate -cmd=validate

  # Show statistics
  migrate -cmd=stats

  # Rollback migration
  migrate -cmd=rollback

Tables to migrate:
`)
	for _, t := range tablesToMigrate {
		fmt.Printf("  - %s (%s)\n", t.TableName, t.Description)
	}
}
