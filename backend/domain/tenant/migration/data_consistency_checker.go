// backend/domain/tenant/migration/data_consistency_checker.go
// 数据一致性校验工具
package migration

import (
	"context"
	"fmt"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// DataConsistencyChecker 数据一致性检查器
type DataConsistencyChecker struct {
	db     *gorm.DB
	logger logs.CtxLogger
}

// NewDataConsistencyChecker 创建数据一致性检查器
func NewDataConsistencyChecker(db *gorm.DB) *DataConsistencyChecker {
	return &DataConsistencyChecker{
		db:     db,
		logger: logs.DefaultLogger(),
	}
}

// ConsistencyCheckResult 一致性检查结果
type ConsistencyCheckResult struct {
	TableName        string    `json:"table_name"`
	TotalRecords     int64     `json:"total_records"`
	NullTenantCount  int64     `json:"null_tenant_count"`
	EmptyTenantCount int64     `json:"empty_tenant_count"`
	ValidTenantCount int64     `json:"valid_tenant_count"`
	Completeness     float64   `json:"completeness"`
	UniquenessValid  bool      `json:"uniqueness_valid"`
	ForeignKeyValid  bool      `json:"foreign_key_valid"`
	ChecksPassed     int       `json:"checks_passed"`
	ChecksTotal      int       `json:"checks_total"`
	AllChecksPassed  bool      `json:"all_checks_passed"`
	CheckedAt        time.Time `json:"checked_at"`
	Issues           []string  `json:"issues"`
}

// CheckTable 检查单张表的一致性
func (c *DataConsistencyChecker) CheckTable(ctx context.Context, tableName string) (*ConsistencyCheckResult, error) {
	c.logger.CtxInfof(ctx, "[Checker] 开始检查表: %s", tableName)

	result := &ConsistencyCheckResult{
		TableName:   tableName,
		CheckedAt:   time.Now(),
		Issues:      make([]string, 0),
		ChecksTotal: 6,
	}

	// 1. 获取总记录数
	c.db.Table(tableName).Count(&result.TotalRecords)
	result.ChecksPassed++

	// 2. 检查NULL的tenant_id
	c.db.Table(tableName).
		Where("tenant_id IS NULL").
		Count(&result.NullTenantCount)

	if result.NullTenantCount > 0 {
		result.Issues = append(result.Issues,
			fmt.Sprintf("发现 %d 条记录的tenant_id为NULL", result.NullTenantCount))
	} else {
		result.ChecksPassed++
	}

	// 3. 检查空字符串tenant_id
	c.db.Table(tableName).
		Where("tenant_id = ''").
		Count(&result.EmptyTenantCount)

	if result.EmptyTenantCount > 0 {
		result.Issues = append(result.Issues,
			fmt.Sprintf("发现 %d 条记录的tenant_id为空字符串", result.EmptyTenantCount))
	} else {
		result.ChecksPassed++
	}

	// 4. 检查有效的tenant_id
	c.db.Table(tableName).
		Where("tenant_id IS NOT NULL AND tenant_id != ''").
		Count(&result.ValidTenantCount)

	result.Completeness = float64(result.ValidTenantCount) / float64(result.TotalRecords) * 100

	if result.Completeness >= 99.9 {
		result.ChecksPassed++
	} else {
		result.Issues = append(result.Issues,
			fmt.Sprintf("数据完整性不足: %.2f%% < 99.9%%", result.Completeness))
	}

	// 5. 检查唯一性约束（如果有唯一约束）
	result.UniquenessValid = c.checkUniqueness(ctx, tableName)
	if result.UniquenessValid {
		result.ChecksPassed++
	} else {
		result.Issues = append(result.Issues, "唯一性约束验证失败")
	}

	// 6. 检查外键约束（如果有外键）
	result.ForeignKeyValid = c.checkForeignKey(ctx, tableName)
	if result.ForeignKeyValid {
		result.ChecksPassed++
	} else {
		result.Issues = append(result.Issues, "外键约束验证失败")
	}

	result.AllChecksPassed = result.ChecksPassed == result.ChecksTotal

	c.logCheckResult(ctx, result)

	return result, nil
}

// checkUniqueness 检查唯一性约束
func (c *DataConsistencyChecker) checkUniqueness(ctx context.Context, tableName string) bool {
	// 检查是否有重复的 (tenant_id, id) 组合
	var duplicates []struct {
		TenantID string `gorm:"column:tenant_id"`
		ID       int64  `gorm:"column:id"`
		Count    int64  `gorm:"column:count"`
	}

	err := c.db.Table(tableName).
		Select("tenant_id, id, COUNT(*) as count").
		Where("tenant_id IS NOT NULL AND tenant_id != ''").
		Group("tenant_id, id").
		Having("count > 1").
		Scan(&duplicates).Error

	if err != nil {
		c.logger.CtxWarnf(ctx, "[Checker] 唯一性检查失败: %v", err)
		return false
	}

	return len(duplicates) == 0
}

// checkForeignKey 检查外键约束
func (c *DataConsistencyChecker) checkForeignKey(ctx context.Context, tableName string) bool {
	// TODO: 根据实际的表结构检查外键
	// 例如: bots表的tenant_id应该引用tenants表

	// 简化实现: 检查tenant_id是否存在于tenants表
	if tableName == "bots" || tableName == "conversations" {
		var invalidCount int64
		c.db.Table(tableName).
			Where("tenant_id NOT IN (SELECT tenant_id FROM tenants)").
			Count(&invalidCount)

		return invalidCount == 0
	}

	return true
}

// CheckAllTables 检查所有表的一致性
func (c *DataConsistencyChecker) CheckAllTables(ctx context.Context, tables []string) map[string]*ConsistencyCheckResult {
	results := make(map[string]*ConsistencyCheckResult)

	for _, tableName := range tables {
		result, err := c.CheckTable(ctx, tableName)
		if err != nil {
			c.logger.CtxErrorf(ctx, "[Checker] 检查表 %s 失败: %v", tableName, err)
			continue
		}
		results[tableName] = result
	}

	c.logAllTablesSummary(ctx, results)

	return results
}

// logCheckResult 输出检查结果
func (c *DataConsistencyChecker) logCheckResult(ctx context.Context, result *ConsistencyCheckResult) {
	c.logger.CtxInfof(ctx, "[Checker] ========== 一致性检查结果 (%s) ==========", result.TableName)
	c.logger.CtxInfof(ctx, "[Checker] 总记录数: %d", result.TotalRecords)
	c.logger.CtxInfof(ctx, "[Checker] NULL tenant_id: %d", result.NullTenantCount)
	c.logger.CtxInfof(ctx, "[Checker] 空 tenant_id: %d", result.EmptyTenantCount)
	c.logger.CtxInfof(ctx, "[Checker] 有效 tenant_id: %d", result.ValidTenantCount)
	c.logger.CtxInfof(ctx, "[Checker] 数据完整性: %.2f%%", result.Completeness)
	c.logger.CtxInfof(ctx, "[Checker] 唯一性约束: %v", result.UniquenessValid)
	c.logger.CtxInfof(ctx, "[Checker] 外键约束: %v", result.ForeignKeyValid)
	c.logger.CtxInfof(ctx, "[Checker] 检查通过: %d/%d", result.ChecksPassed, result.ChecksTotal)
	c.logger.CtxInfof(ctx, "[Checker] 全部通过: %v", result.AllChecksPassed)

	if len(result.Issues) > 0 {
		c.logger.CtxWarnf(ctx, "[Checker] 发现问题:")
		for _, issue := range result.Issues {
			c.logger.CtxWarnf(ctx, "[Checker]   - %s", issue)
		}
	}

	c.logger.CtxInfof(ctx, "[Checker] ===============================================")
}

// logAllTablesSummary 输出所有表的汇总
func (c *DataConsistencyChecker) logAllTablesSummary(ctx context.Context, results map[string]*ConsistencyCheckResult) {
	c.logger.CtxInfof(ctx, "[Checker] ========== 一致性检查汇总 ==========")

	totalPassed := 0
	totalTables := len(results)

	for tableName, result := range results {
		if result.AllChecksPassed {
			totalPassed++
		}

		c.logger.CtxInfof(ctx, "[Checker] %s: %v (%d/%d检查通过)",
			tableName,
			map[bool]string{true: "✓ 通过", false: "✗ 失败"}[result.AllChecksPassed],
			result.ChecksPassed,
			result.ChecksTotal)
	}

	passRate := float64(totalPassed) / float64(totalTables) * 100
	c.logger.CtxInfof(ctx, "[Checker] 总计: %d/%d 表通过 (%.2f%%)", totalPassed, totalTables, passRate)
	c.logger.CtxInfof(ctx, "[Checker] ===========================================")
}

// FixInconsistentData 修复不一致数据
func (c *DataConsistencyChecker) FixInconsistentData(ctx context.Context, tableName string, defaultTenant string) (int64, error) {
	c.logger.CtxInfof(ctx, "[Checker] 开始修复不一致数据: %s", tableName)

	// 修复NULL的tenant_id
	result := c.db.Table(tableName).
		Where("tenant_id IS NULL OR tenant_id = ''").
		Update("tenant_id", defaultTenant)

	if result.Error != nil {
		return 0, result.Error
	}

	c.logger.CtxInfof(ctx, "[Checker] 修复完成: %d 条记录", result.RowsAffected)
	return result.RowsAffected, nil
}

// GenerateConsistencyReport 生成一致性报告
func (c *DataConsistencyChecker) GenerateConsistencyReport(ctx context.Context, tables []string) *ConsistencyReport {
	results := c.CheckAllTables(ctx, tables)

	report := &ConsistencyReport{
		GeneratedAt: time.Now(),
		Tables:      results,
	}

	// 汇总统计
	for _, result := range results {
		report.TotalTables++
		report.TotalRecords += result.TotalRecords
		report.ValidRecords += result.ValidTenantCount
		report.InvalidRecords += (result.NullTenantCount + result.EmptyTenantCount)

		if result.AllChecksPassed {
			report.PassedTables++
		} else {
			report.FailedTables++
			report.FailedTableNames = append(report.FailedTableNames, result.TableName)
		}

		if result.Completeness >= 99.9 {
			report.HighQualityTables++
		}
	}

	report.OverallCompleteness = float64(report.ValidRecords) / float64(report.TotalRecords) * 100
	report.OverallPassRate = float64(report.PassedTables) / float64(report.TotalTables) * 100

	return report
}

// ConsistencyReport 一致性报告
type ConsistencyReport struct {
	GeneratedAt         time.Time                          `json:"generated_at"`
	Tables              map[string]*ConsistencyCheckResult `json:"tables"`
	TotalTables         int                                `json:"total_tables"`
	PassedTables        int                                `json:"passed_tables"`
	FailedTables        int                                `json:"failed_tables"`
	HighQualityTables   int                                `json:"high_quality_tables"`
	FailedTableNames    []string                           `json:"failed_table_names,omitempty"`
	TotalRecords        int64                              `json:"total_records"`
	ValidRecords        int64                              `json:"valid_records"`
	InvalidRecords      int64                              `json:"invalid_records"`
	OverallCompleteness float64                            `json:"overall_completeness"`
	OverallPassRate     float64                            `json:"overall_pass_rate"`
}
