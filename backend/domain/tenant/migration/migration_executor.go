// backend/domain/tenant/migration/migration_executor.go
// 迁移执行器 - 实际执行数据迁移的主工具
package migration

import (
	"context"
	"fmt"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// MigrationExecutor 迁移执行器
type MigrationExecutor struct {
	db             *gorm.DB
	batchProcessor *BatchProcessor
	logger         logs.CtxLogger
}

// NewMigrationExecutor 创建迁移执行器
func NewMigrationExecutor(db *gorm.DB, batchSize, maxConcurrency int, defaultTenant string) *MigrationExecutor {
	return &MigrationExecutor{
		db:             db,
		batchProcessor: NewBatchProcessor(db, batchSize, maxConcurrency, defaultTenant),
		logger:         logs.DefaultLogger(),
	}
}

// ExecuteMigrationPlan 执行迁移计划
func (e *MigrationExecutor) ExecuteMigrationPlan(ctx context.Context, plan *MigrationPlan) (*MigrationExecutionReport, error) {
	e.logger.CtxInfof(ctx, "[Executor] ========== 开始执行迁移计划 ==========")
	e.logger.CtxInfof(ctx, "[Executor] 计划名称: %s", plan.PlanName)
	e.logger.CtxInfof(ctx, "[Executor] 涉及表数: %d", len(plan.Tables))
	e.logger.CtxInfof(ctx, "[Executor] 批次大小: %d", plan.BatchSize)
	e.logger.CtxInfof(ctx, "[Executor] 并发度: %d", plan.MaxConcurrency)

	report := &MigrationExecutionReport{
		PlanName:     plan.PlanName,
		StartedAt:    time.Now(),
		TableReports: make(map[string]*BatchProcessResult),
	}

	// 1. 预检查
	if err := e.preCheck(ctx, plan); err != nil {
		return nil, fmt.Errorf("预检查失败: %w", err)
	}

	// 2. 执行迁移
	for _, tableName := range plan.Tables {
		e.logger.CtxInfof(ctx, "[Executor] 开始迁移表: %s", tableName)

		result, err := e.batchProcessor.ProcessTable(ctx, tableName)
		if err != nil {
			e.logger.CtxErrorf(ctx, "[Executor] 迁移表 %s 失败: %v", tableName, err)
			report.FailedTables = append(report.FailedTables, tableName)
			continue
		}

		report.TableReports[tableName] = result
		report.SuccessTables = append(report.SuccessTables, tableName)

		e.logger.CtxInfof(ctx, "[Executor] 表 %s 迁移完成: %d/%d (%.2f%%)",
			tableName, result.MigratedRecords, result.TotalRecords,
			float64(result.MigratedRecords)/float64(result.TotalRecords)*100)
	}

	// 3. 后验证
	if err := e.postCheck(ctx, plan); err != nil {
		e.logger.CtxErrorf(ctx, "[Executor] 后验证失败: %v", err)
		report.ValidationError = err.Error()
	}

	report.FinishedAt = time.Now()
	report.DurationSeconds = report.FinishedAt.Sub(report.StartedAt).Seconds()

	// 4. 生成汇总
	e.generateSummary(report)

	e.logFinalReport(ctx, report)

	return report, nil
}

// preCheck 预检查
func (e *MigrationExecutor) preCheck(ctx context.Context, plan *MigrationPlan) error {
	e.logger.CtxInfof(ctx, "[Executor] ========== 预检查 ==========")

	// 1. 检查表是否存在
	for _, tableName := range plan.Tables {
		if !e.db.Migrator().HasTable(tableName) {
			return fmt.Errorf("表不存在: %s", tableName)
		}
		e.logger.CtxInfof(ctx, "[Executor] ✓ 表存在: %s", tableName)
	}

	// 2. 检查tenant_id字段是否存在
	for _, tableName := range plan.Tables {
		if !e.db.Migrator().HasColumn(tableName, "tenant_id") {
			return fmt.Errorf("字段不存在: %s.tenant_id", tableName)
		}
		e.logger.CtxInfof(ctx, "[Executor] ✓ 字段存在: %s.tenant_id", tableName)
	}

	// 3. 检查是否有索引
	for _, tableName := range plan.Tables {
		var indexes []struct {
			IndexName string `gorm:"column:Key_name"`
		}
		e.db.Table("information_schema.STATISTICS").
			Select("Key_name").
			Where("table_schema = DATABASE() AND table_name = ? AND key_name LIKE ?", tableName, "%tenant%").
			Scan(&indexes)

		if len(indexes) == 0 {
			e.logger.CtxWarnf(ctx, "[Executor] ⚠ 警告: %s.tenant_id 缺少索引", tableName)
		} else {
			e.logger.CtxInfof(ctx, "[Executor] ✓ 索引存在: %s.tenant_id (%d个)", tableName, len(indexes))
		}
	}

	e.logger.CtxInfof(ctx, "[Executor] ========== 预检查通过 ==========")
	return nil
}

// postCheck 后验证
func (e *MigrationExecutor) postCheck(ctx context.Context, plan *MigrationPlan) error {
	e.logger.CtxInfof(ctx, "[Executor] ========== 后验证 ==========")

	// 1. 检查是否有NULL的tenant_id
	for _, tableName := range plan.Tables {
		var nullCount int64
		e.db.Table(tableName).
			Where("tenant_id IS NULL OR tenant_id = ''").
			Count(&nullCount)

		if nullCount > 0 {
			e.logger.CtxWarnf(ctx, "[Executor] ⚠ 警告: %s 有 %d 条记录的tenant_id为空", tableName, nullCount)
		} else {
			e.logger.CtxInfof(ctx, "[Executor] ✓ 无空值: %s.tenant_id", tableName)
		}
	}

	// 2. 检查数据一致性（记录数）
	for _, tableName := range plan.Tables {
		var totalCount int64
		e.db.Table(tableName).Count(&totalCount)

		var notNullCount int64
		e.db.Table(tableName).
			Where("tenant_id IS NOT NULL AND tenant_id != ''").
			Count(&notNullCount)

		completeness := float64(notNullCount) / float64(totalCount) * 100
		e.logger.CtxInfof(ctx, "[Executor] ✓ 数据完整性: %s (%.2f%%)", tableName, completeness)

		if completeness < 99.9 {
			return fmt.Errorf("数据完整性不足: %s (%.2f%% < 99.9%%)", tableName, completeness)
		}
	}

	// 3. 抽样验证
	for _, tableName := range plan.Tables {
		var samples []struct {
			ID       int64  `gorm:"column:id"`
			TenantID string `gorm:"column:tenant_id"`
		}
		e.db.Table(tableName).
			Limit(10).
			Scan(&samples)

		validCount := 0
		for _, sample := range samples {
			if sample.TenantID != "" && sample.TenantID != "default_tenant" {
				validCount++
			}
		}

		e.logger.CtxInfof(ctx, "[Executor] ✓ 抽样验证: %s (有效: %d/10)", tableName, validCount)
	}

	e.logger.CtxInfof(ctx, "[Executor] ========== 后验证通过 ==========")
	return nil
}

// generateSummary 生成汇总
func (e *MigrationExecutor) generateSummary(report *MigrationExecutionReport) {
	totalTables := len(report.TableReports)
	totalRecords := int64(0)
	totalMigrated := int64(0)
	totalFailed := int64(0)

	for _, result := range report.TableReports {
		totalRecords += result.TotalRecords
		totalMigrated += result.MigratedRecords
		totalFailed += result.FailedRecords
	}

	report.TotalTables = totalTables
	report.SuccessTablesCount = len(report.SuccessTables)
	report.FailedTablesCount = len(report.FailedTables)
	report.TotalRecords = totalRecords
	report.TotalMigrated = totalMigrated
	report.TotalFailed = totalFailed
	report.SuccessRate = float64(totalMigrated) / float64(totalMigrated+totalFailed) * 100
}

// logFinalReport 输出最终报告
func (e *MigrationExecutor) logFinalReport(ctx context.Context, report *MigrationExecutionReport) {
	e.logger.CtxInfof(ctx, "[Executor] ========== 迁移执行报告 ==========")
	e.logger.CtxInfof(ctx, "[Executor] 计划名称: %s", report.PlanName)
	e.logger.CtxInfof(ctx, "[Executor] 开始时间: %s", report.StartedAt.Format("2006-01-02 15:04:05"))
	e.logger.CtxInfof(ctx, "[Executor] 完成时间: %s", report.FinishedAt.Format("2006-01-02 15:04:05"))
	e.logger.CtxInfof(ctx, "[Executor] 总耗时: %.2f 秒", report.DurationSeconds)
	e.logger.CtxInfof(ctx, "[Executor]")
	e.logger.CtxInfof(ctx, "[Executor] 涉及表数: %d", report.TotalTables)
	e.logger.CtxInfof(ctx, "[Executor] 成功表数: %d", report.SuccessTablesCount)
	e.logger.CtxInfof(ctx, "[Executor] 失败表数: %d", report.FailedTablesCount)
	e.logger.CtxInfof(ctx, "[Executor]")
	e.logger.CtxInfof(ctx, "[Executor] 总记录数: %d", report.TotalRecords)
	e.logger.CtxInfof(ctx, "[Executor] 已迁移: %d", report.TotalMigrated)
	e.logger.CtxInfof(ctx, "[Executor] 失败: %d", report.TotalFailed)
	e.logger.CtxInfof(ctx, "[Executor] 成功率: %.2f%%", report.SuccessRate)
	e.logger.CtxInfof(ctx, "[Executor] 平均速率: %.2f 条/秒", float64(report.TotalMigrated)/report.DurationSeconds)
	e.logger.CtxInfof(ctx, "[Executor] ======================================")

	if report.ValidationError != "" {
		e.logger.CtxWarnf(ctx, "[Executor] 验证错误: %s", report.ValidationError)
	}
}

// MigrationPlan 迁移计划
type MigrationPlan struct {
	PlanName       string
	Tables         []string
	BatchSize      int
	MaxConcurrency int
	DefaultTenant  string
}

// MigrationExecutionReport 迁移执行报告
type MigrationExecutionReport struct {
	PlanName           string                         `json:"plan_name"`
	StartedAt          time.Time                      `json:"started_at"`
	FinishedAt         time.Time                      `json:"finished_at"`
	DurationSeconds    float64                        `json:"duration_seconds"`
	TableReports       map[string]*BatchProcessResult `json:"table_reports"`
	SuccessTables      []string                       `json:"success_tables"`
	FailedTables       []string                       `json:"failed_tables"`
	TotalTables        int                            `json:"total_tables"`
	SuccessTablesCount int                            `json:"success_tables_count"`
	FailedTablesCount  int                            `json:"failed_tables_count"`
	TotalRecords       int64                          `json:"total_records"`
	TotalMigrated      int64                          `json:"total_migrated"`
	TotalFailed        int64                          `json:"total_failed"`
	SuccessRate        float64                        `json:"success_rate"`
	ValidationError    string                         `json:"validation_error,omitempty"`
}

// DefaultMigrationPlan 默认迁移计划
func DefaultMigrationPlan() *MigrationPlan {
	return &MigrationPlan{
		PlanName: "tenant_id_migration_phase2",
		Tables: []string{
			"bots",
			"bot_configs",
			"conversations",
			"messages",
			"knowledge_bases",
			"knowledge_chunks",
			"workflows",
			"workflow_executions",
			"single_agent_draft",
			"published_bots",
		},
		BatchSize:      1000,
		MaxConcurrency: 10,
		DefaultTenant:  "default_tenant",
	}
}
