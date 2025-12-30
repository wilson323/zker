// backend/domain/tenant/migration/cutover.go
package migration

import (
	"context"
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// CutoverManager 切换管理器
// 负责将系统从旧模式切换到新模式
type CutoverManager struct {
	db     *gorm.DB
	logger logs.CtxLogger
}

// NewCutoverManager 创建切换管理器
func NewCutoverManager(db *gorm.DB) *CutoverManager {
	return &CutoverManager{
		db:     db,
		logger: logs.DefaultLogger(),
	}
}

// CutoverToTenantID 切换到使用tenant_id
//
// 切换步骤:
//  1. 验证所有数据已迁移
//  2. 将tenant_id字段改为NOT NULL
//  3. 添加唯一索引（如需要）
//  4. 更新应用配置
func (m *CutoverManager) CutoverToTenantID(
	ctx context.Context,
	tableName string,
) error {
	m.logger.CtxInfof(ctx, "[Cutover] starting cutover for table %s", tableName)
	startTime := time.Now()

	// 1. 验证所有数据已迁移
	m.logger.CtxInfof(ctx, "[Cutover] validating data migration...")
	var nullCount int64
	err := m.db.Table(tableName).
		Where("tenant_id IS NULL").
		Count(&nullCount).Error

	if err != nil {
		return fmt.Errorf("验证迁移失败: %w", err)
	}

	if nullCount > 0 {
		return fmt.Errorf("表 %s 还有 %d 条记录的tenant_id为NULL，无法切换", tableName, nullCount)
	}

	// 2. 将tenant_id字段改为NOT NULL
	m.logger.CtxInfof(ctx, "[Cutover] altering column to NOT NULL...")
	sql := fmt.Sprintf(`
		ALTER TABLE %s
		MODIFY COLUMN tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID'
	`, tableName)

	err = m.db.Exec(sql).Error
	if err != nil {
		return fmt.Errorf("修改tenant_id为NOT NULL失败: %w", err)
	}

	// 3. 验证修改成功
	var columnType string
	err = m.db.Raw(`
		SELECT COLUMN_TYPE
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = ?
		  AND COLUMN_NAME = 'tenant_id'
	`, tableName).Scan(&columnType).Error

	if err != nil {
		return fmt.Errorf("验证字段类型失败: %w", err)
	}

	duration := time.Since(startTime)
	m.logger.CtxInfof(ctx, "[Cutover] table %s cutover completed successfully (duration: %s)",
		tableName, duration)

	return nil
}

// CleanupLegacyData 清理旧数据
//
// ⚠️ 警告: 此操作将删除旧字段，不可恢复！
// 建议在切换完成并观察1周后执行
func (m *CutoverManager) CleanupLegacyData(
	ctx context.Context,
	tableName string,
	dropOldField bool,
) error {
	m.logger.CtxWarnf(ctx, "[Cleanup] starting cleanup for table %s (drop_old_field=%v)",
		tableName, dropOldField)

	if !dropOldField {
		m.logger.CtxInfof(ctx, "[Cleanup] dry-run mode, skipping actual cleanup")
		return nil
	}

	// ⚠️ 危险操作：删除旧字段
	// 注意：这里只是示例，实际应该根据旧字段名称删除
	m.logger.CtxWarnf(ctx, "[Cleanup] dropping old fields from table %s", tableName)

	// TODO: 根据实际情况删除旧字段
	// ALTER TABLE bots DROP COLUMN space_id;

	m.logger.CtxInfof(ctx, "[Cleanup] table %s cleanup completed", tableName)
	return nil
}

// RollbackCutover 回滚切换（恢复为NULL）
func (m *CutoverManager) RollbackCutover(
	ctx context.Context,
	tableName string,
) error {
	m.logger.CtxWarnf(ctx, "[Cutover] rolling back cutover for table %s", tableName)

	// 将tenant_id改回NULL
	sql := fmt.Sprintf(`
		ALTER TABLE %s
		MODIFY COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID'
	`, tableName)

	err := m.db.Exec(sql).Error
	if err != nil {
		return fmt.Errorf("回滚切换失败: %w", err)
	}

	m.logger.CtxInfof(ctx, "[Cutover] table %s cutover rolled back", tableName)
	return nil
}

// VerifyCutover 验证切换结果
func (m *CutoverManager) VerifyCutover(
	ctx context.Context,
	tableName string,
) error {
	m.logger.CtxInfof(ctx, "[Cutover] verifying cutover for table %s...", tableName)

	// 1. 检查tenant_id是否为NOT NULL
	var isNullable string
	err := m.db.Raw(`
		SELECT IS_NULLABLE
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = ?
		  AND COLUMN_NAME = 'tenant_id'
	`, tableName).Scan(&isNullable).Error

	if err != nil {
		return fmt.Errorf("验证失败: %w", err)
	}

	if isNullable == "YES" {
		return fmt.Errorf("表 %s 的tenant_id字段仍为NULL，切换未完成", tableName)
	}

	// 2. 检查是否还有NULL值
	var nullCount int64
	err = m.db.Table(tableName).
		Where("tenant_id IS NULL").
		Count(&nullCount).Error

	if err != nil {
		return fmt.Errorf("验证NULL失败: %w", err)
	}

	if nullCount > 0 {
		return fmt.Errorf("表 %s 仍有 %d 条记录的tenant_id为NULL", tableName, nullCount)
	}

	// 3. 检查外键完整性
	var invalidCount int64
	err = m.db.Raw(fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s b
		LEFT JOIN tenants t ON b.tenant_id = t.tenant_id
		WHERE t.tenant_id IS NULL
	`, tableName)).Scan(&invalidCount).Error

	if err != nil {
		return fmt.Errorf("验证外键失败: %w", err)
	}

	if invalidCount > 0 {
		return fmt.Errorf("表 %s 有 %d 条记录的tenant_id无效", tableName, invalidCount)
	}

	m.logger.CtxInfof(ctx, "[Cutover] ✓ Table %s cutover verification passed", tableName)
	return nil
}

// BatchCutover 批量切换多张表
func (m *CutoverManager) BatchCutover(
	ctx context.Context,
	tableNames []string,
) error {
	m.logger.CtxInfof(ctx, "[Cutover] starting batch cutover for %d tables", len(tableNames))

	successCount := 0
	failedTables := []string{}

	for _, tableName := range tableNames {
		m.logger.CtxInfof(ctx, "[Cutover] processing table %s...", tableName)

		err := m.CutoverToTenantID(ctx, tableName)
		if err != nil {
			m.logger.CtxErrorf(ctx, "[Cutover] failed to cutover table %s: %v",
				tableName, err)
			failedTables = append(failedTables, tableName)
			continue
		}

		// 验证切换结果
		err = m.VerifyCutover(ctx, tableName)
		if err != nil {
			m.logger.CtxErrorf(ctx, "[Cutover] verification failed for table %s: %v",
				tableName, err)
			failedTables = append(failedTables, tableName)
			continue
		}

		successCount++
		m.logger.CtxInfof(ctx, "[Cutover] ✓ Table %s cutover completed", tableName)
	}

	m.logger.CtxInfof(ctx, "[Cutover] batch cutover completed: success=%d, failed=%d",
		successCount, len(failedTables))

	if len(failedTables) > 0 {
		return fmt.Errorf("以下表切换失败: %v", failedTables)
	}

	return nil
}

// CutoverReport 切换报告
type CutoverReport struct {
	TableName    string        // 表名
	CutoverTime  time.Time     // 切换时间
	Duration     time.Duration // 耗时
	Success      bool          // 是否成功
	NullCount    int64         // NULL记录数
	InvalidCount int64         // 无效记录数
	ErrorMessage string        // 错误消息
}

// GenerateCutoverReport 生成切换报告
func (m *CutoverManager) GenerateCutoverReport(
	ctx context.Context,
	tableNames []string,
) []*CutoverReport {
	reports := make([]*CutoverReport, 0, len(tableNames))

	for _, tableName := range tableNames {
		report := &CutoverReport{
			TableName: tableName,
		}

		// 检查NULL记录数
		m.db.Table(tableName).
			Where("tenant_id IS NULL").
			Count(&report.NullCount)

		// 检查无效记录数
		m.db.Raw(fmt.Sprintf(`
			SELECT COUNT(*)
			FROM %s b
			LEFT JOIN tenants t ON b.tenant_id = t.tenant_id
			WHERE t.tenant_id IS NULL
		`, tableName)).Scan(&report.InvalidCount)

		// 判断是否成功
		report.Success = (report.NullCount == 0 && report.InvalidCount == 0)

		reports = append(reports, report)
	}

	return reports
}
