// backend/domain/tenant/migration/cli.go
// 租户数据迁移CLI工具
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/migration"
)

func main() {
	app := &cli.App{
		Name:  "tenant-migration",
		Usage: "ZKER 租户数据迁移工具",
		Commands: []*cli.Command{
			{
				Name:  "migrate",
				Usage: "执行tenant_id迁移",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "tables",
						Usage:    "要迁移的表名，逗号分隔",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "mode",
						Usage:   "迁移模式: double_write, direct",
						Value:   "double_write",
					},
					&cli.StringFlag{
						Name:    "default-tenant",
						Usage:   "默认租户ID",
						Value:   "system_tenant",
					},
					&cli.IntFlag{
						Name:    "batch-size",
						Usage:   "批处理大小",
						Value:   1000,
					},
					&cli.IntFlag{
						Name:    "workers",
						Usage:   "并发worker数量",
						Value:   5,
					},
					&cli.StringFlag{
						Name:    "dsn",
						Usage:   "数据库连接字符串",
						EnvVars: []string{"DATABASE_DSN"},
					},
				},
				Action: migrateCmd,
			},
			{
				Name:  "progress",
				Usage: "查询迁移进度",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "migration-id",
						Usage:    "迁移任务ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "dsn",
						Usage:   "数据库连接字符串",
						EnvVars: []string{"DATABASE_DSN"},
					},
				},
				Action: progressCmd,
			},
			{
				Name:  "rollback",
				Usage: "回滚迁移",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "migration-id",
						Usage:    "迁移任务ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "dsn",
						Usage:   "数据库连接字符串",
						EnvVars: []string{"DATABASE_DSN"},
					},
				},
				Action: rollbackCmd,
			},
			{
				Name:  "validate",
				Usage: "验证迁移结果",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "migration-id",
						Usage:    "迁移任务ID",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "dsn",
						Usage:   "数据库连接字符串",
						EnvVars: []string{"DATABASE_DSN"},
					},
				},
				Action: validateCmd,
			},
			{
				Name:  "list",
				Usage: "列出所有迁移任务",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "dsn",
						Usage:   "数据库连接字符串",
						EnvVars: []string{"DATABASE_DSN"},
					},
				},
				Action: listCmd,
			},
			{
				Name:  "verify-dual-write",
				Usage: "验证双写数据一致性",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "table",
						Usage:    "表名",
						Required: true,
					},
					&cli.DurationFlag{
						Name:    "duration",
						Usage:   "检查最近N小时的数据",
						Value:   1 * time.Hour,
					},
					&cli.StringFlag{
						Name:    "dsn",
						Usage:   "数据库连接字符串",
						EnvVars: []string{"DATABASE_DSN"},
					},
				},
				Action: verifyDualWriteCmd,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		hlog.Fatalf("运行失败: %v", err)
	}
}

// ================================================================================
// 命令处理函数
// ================================================================================

// migrateCmd 执行迁移
func migrateCmd(c *cli.Context) error {
	// 1. 连接数据库
	db, err := connectDB(c.String("dsn"))
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 2. 创建迁移服务
	service := migration.NewMigrationService(db)

	// 3. 解析表名
	tables := parseStringList(c.String("tables"))

	// 4. 执行迁移
	config := &migration.MigrationConfig{
		TableNames:    tables,
		Mode:          c.String("mode"),
		DefaultTenant: c.String("default-tenant"),
		BatchSize:     c.Int("batch-size"),
		Workers:       c.Int("workers"),
	}

	ctx := context.Background()
	if err := service.StartMigration(ctx, config); err != nil {
		return fmt.Errorf("启动迁移失败: %w", err)
	}

	fmt.Printf("✅ 迁移任务已启动: %s\n", config.MigrationID)
	fmt.Printf("   表数量: %d\n", len(tables))
	fmt.Printf("   迁移模式: %s\n", config.Mode)
	fmt.Printf("   批处理大小: %d\n", config.BatchSize)
	fmt.Printf("\n使用以下命令查询进度:\n")
	fmt.Printf("   tenant-migration progress --migration-id %s\n", config.MigrationID)

	return nil
}

// progressCmd 查询迁移进度
func progressCmd(c *cli.Context) error {
	db, err := connectDB(c.String("dsn"))
	if err != nil {
		return err
	}

	service := migration.NewMigrationService(db)
	migrationID := c.String("migration-id")

	ctx := context.Background()
	progress, err := service.GetProgress(ctx, migrationID)
	if err != nil {
		return fmt.Errorf("查询进度失败: %w", err)
	}

	fmt.Printf("迁移任务: %s\n", progress.MigrationID)
	fmt.Printf("状态: %s\n", progress.Status)
	fmt.Printf("进度: %.2f%%\n", progress.ProgressPercent)
	fmt.Printf("当前表: %s\n", progress.CurrentTable)
	fmt.Printf("完成表: %d/%d\n", progress.CompletedTables, progress.TotalTables)
	fmt.Printf("记录数: %d/%d\n", progress.MigratedRecords, progress.TotalRecords)
	fmt.Printf("失败数: %d\n", progress.FailedRecords)
	fmt.Printf("开始时间: %s\n", progress.StartedAt.Format("2006-01-02 15:04:05"))
	if progress.FinishedAt != nil {
		fmt.Printf("完成时间: %s\n", progress.FinishedAt.Format("2006-01-02 15:04:05"))
	}
	if progress.Error != "" {
		fmt.Printf("错误: %s\n", progress.Error)
	}

	return nil
}

// rollbackCmd 回滚迁移
func rollbackCmd(c *cli.Context) error {
	db, err := connectDB(c.String("dsn"))
	if err != nil {
		return err
	}

	service := migration.NewMigrationService(db)
	migrationID := c.String("migration-id")

	ctx := context.Background()
	if err := service.Rollback(ctx, migrationID); err != nil {
		return fmt.Errorf("回滚失败: %w", err)
	}

	fmt.Printf("✅ 迁移已回滚: %s\n", migrationID)
	return nil
}

// validateCmd 验证迁移结果
func validateCmd(c *cli.Context) error {
	db, err := connectDB(c.String("dsn"))
	if err != nil {
		return err
	}

	service := migration.NewMigrationService(db)
	migrationID := c.String("migration-id")

	ctx := context.Background()
	if err := service.Validate(ctx, migrationID); err != nil {
		return fmt.Errorf("验证失败: %w", err)
	}

	fmt.Printf("✅ 迁移验证通过: %s\n", migrationID)
	return nil
}

// listCmd 列出所有迁移任务
func listCmd(c *cli.Context) error {
	db, err := connectDB(c.String("dsn"))
	if err != nil {
		return err
	}

	service := migration.NewMigrationService(db)

	ctx := context.Background()
	records, err := service.ListMigrations(ctx)
	if err != nil {
		return fmt.Errorf("查询迁移记录失败: %w", err)
	}

	if len(records) == 0 {
		fmt.Println("没有找到迁移记录")
		return nil
	}

	fmt.Printf("找到 %d 条迁移记录:\n\n", len(records))
	for _, record := range records {
		fmt.Printf("ID: %s\n", record.MigrationID)
		fmt.Printf("  状态: %s\n", record.Status)
		fmt.Printf("  进度: %.2f%%\n", record.ProgressPercent)
		fmt.Printf("  开始时间: %s\n", record.StartedAt.Format("2006-01-02 15:04:05"))
		if record.FinishedAt != nil {
			fmt.Printf("  完成时间: %s\n", record.FinishedAt.Format("2006-01-02 15:04:05"))
		}
		if record.Error != "" {
			fmt.Printf("  错误: %s\n", record.Error)
		}
		fmt.Println()
	}

	return nil
}

// verifyDualWriteCmd 验证双写数据一致性
func verifyDualWriteCmd(c *cli.Context) error {
	db, err := connectDB(c.String("dsn"))
	if err != nil {
		return err
	}

	verifier := migration.NewDualWriteVerifier(db)
	tableName := c.String("table")
	duration := c.Duration("duration")

	ctx := context.Background()
	if err := verifier.VerifyDualWrite(ctx, tableName, duration); err != nil {
		return fmt.Errorf("双写验证失败: %w", err)
	}

	fmt.Printf("✅ 双写数据一致性验证通过: %s\n", tableName)
	return nil
}

// ================================================================================
// 辅助函数
// ================================================================================

// connectDB 连接数据库
func connectDB(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("数据库连接字符串不能为空")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}

// parseStringList 解析逗号分隔的字符串列表
func parseStringList(s string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	current := ""
	for _, ch := range s {
		if ch == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}

	return result
}
