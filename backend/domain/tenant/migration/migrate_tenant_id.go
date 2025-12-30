// backend/domain/tenant/migration/migrate_tenant_id.go
package migration

import (
	"context"
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// TenantIDMigrator tenant_id迁移器
// 负责将业务表的user_id映射为tenant_id
type TenantIDMigrator struct {
	logger logs.CtxLogger
	db     *gorm.DB
}

// NewTenantIDMigrator 创建迁移器
func NewTenantIDMigrator(db *gorm.DB) *TenantIDMigrator {
	return &TenantIDMigrator{
		db:     db,
		logger: logs.DefaultLogger(),
	}
}

// MigrateTable 迁移单张表的tenant_id
//
// 参数:
//   - ctx: 上下文
//   - tableName: 表名，如 "bots"
//   - userIDField: 用户ID字段名，如 "creator_id"
//
// 迁移策略:
//  1. 分批查询未迁移数据（每批1000条）
//  2. 从users表查找tenant_id
//  3. 批量更新tenant_id
//  4. 批次之间休息100ms，避免锁表
func (m *TenantIDMigrator) MigrateTable(
	ctx context.Context,
	tableName string,
	userIDField string,
) error {
	m.logger.CtxInfof(ctx, "[Migration] starting migration for table %s (user_id_field=%s)",
		tableName, userIDField)

	startTime := time.Now()
	batchSize := 1000
	offset := 0
	totalMigrated := 0
	totalErrors := 0

	for {
		// 1. 查询一批未迁移的数据
		var results []map[string]interface{}
		err := m.db.Table(tableName).
			Select(fmt.Sprintf("%s as user_id, id", userIDField)).
			Where("tenant_id IS NULL").
			Where(fmt.Sprintf("%s IS NOT NULL", userIDField)).
			Limit(batchSize).
			Offset(offset).
			Find(&results).Error

		if err != nil {
			m.logger.CtxErrorf(ctx, "[Migration] failed to query batch %d: %v",
				offset/batchSize+1, err)
			return fmt.Errorf("查询数据失败: %w", err)
		}

		if len(results) == 0 {
			// 没有更多数据需要迁移
			break
		}

		// 2. 批量更新tenant_id
		batchMigrated := 0
		for _, row := range results {
			userID := fmt.Sprintf("%v", row["user_id"])
			id := fmt.Sprintf("%v", row["id"])

			// 从users表查找tenant_id
			var tenantID string
			err := m.db.Table("users").
				Select("tenant_id").
				Where("user_id = ?", userID).
				Pluck("tenant_id", &tenantID).Error

			if err != nil {
				m.logger.CtxErrorf(ctx, "[Migration] failed to get tenant_id for user %s: %v",
					userID, err)
				totalErrors++
				continue // 跳过这条数据，继续下一条
			}

			if tenantID == "" {
				// 用户没有租户，分配默认租户
				tenantID = "tenant_default"
				m.logger.CtxWarnf(ctx, "[Migration] user %s has no tenant, using default",
					userID)
			}

			// 更新tenant_id
			err = m.db.Table(tableName).
				Where("id = ?", id).
				Update("tenant_id", tenantID).Error

			if err != nil {
				m.logger.CtxErrorf(ctx, "[Migration] failed to update %s id %s: %v",
					tableName, id, err)
				totalErrors++
				continue
			}

			batchMigrated++
			totalMigrated++
		}

		// 3. 记录进度
		batchNum := offset/batchSize + 1
		m.logger.CtxInfof(ctx, "[Migration] table %s: batch %d, migrated %d records (total: %d, errors: %d)",
			tableName, batchNum, batchMigrated, totalMigrated, totalErrors)

		offset += batchSize

		// 4. 避免锁表，每批之间休息100ms
		time.Sleep(100 * time.Millisecond)
	}

	duration := time.Since(startTime)
	m.logger.CtxInfof(ctx, "[Migration] table %s migration completed: total %d records, errors %d, duration %s",
		tableName, totalMigrated, totalErrors, duration)

	return nil
}

// ValidateMigration 验证迁移结果
//
// 验证项:
//  1. 检查是否还有NULL的tenant_id
//  2. 检查外键完整性（所有tenant_id都存在于tenants表）
func (m *TenantIDMigrator) ValidateMigration(ctx context.Context, tableName string) error {
	m.logger.CtxInfof(ctx, "[Migration] validating table %s...", tableName)

	// 1. 检查是否还有NULL的tenant_id
	var nullCount int64
	err := m.db.Table(tableName).
		Where("tenant_id IS NULL").
		Count(&nullCount).Error

	if err != nil {
		return fmt.Errorf("验证NULL失败: %w", err)
	}

	if nullCount > 0 {
		return fmt.Errorf("表 %s 还有 %d 条记录的tenant_id为NULL", tableName, nullCount)
	}

	// 2. 检查外键完整性
	var invalidCount int64
	err = m.db.Table(tableName).
		Joins(fmt.Sprintf("LEFT JOIN tenants ON tenants.tenant_id = %s.tenant_id", tableName)).
		Where("tenants.tenant_id IS NULL").
		Count(&invalidCount).Error

	if err != nil {
		return fmt.Errorf("验证外键失败: %w", err)
	}

	if invalidCount > 0 {
		return fmt.Errorf("表 %s 有 %d 条记录的tenant_id不存在于tenants表", tableName, invalidCount)
	}

	m.logger.CtxInfof(ctx, "[Migration] table %s validation passed", tableName)
	return nil
}

// RollbackMigration 回滚迁移
//
// ⚠️ 警告: 此操作将清空已迁移的tenant_id数据
func (m *TenantIDMigrator) RollbackMigration(ctx context.Context, tableName string) error {
	m.logger.CtxWarnf(ctx, "[Migration] rolling back table %s...", tableName)

	result := m.db.Table(tableName).
		Where("tenant_id IS NOT NULL").
		Update("tenant_id", nil)

	if result.Error != nil {
		return fmt.Errorf("回滚失败: %w", result.Error)
	}

	m.logger.CtxInfof(ctx, "[Migration] table %s rollback completed: %d records cleared",
		tableName, result.RowsAffected)

	return nil
}

// GetMigrationStats 获取迁移统计信息
func (m *TenantIDMigrator) GetMigrationStats(
	ctx context.Context,
	tableName string,
) (*MigrationStats, error) {
	stats := &MigrationStats{}

	// 总记录数
	m.db.Table(tableName).Count(&stats.TotalRecords)

	// 已迁移记录数
	m.db.Table(tableName).Where("tenant_id IS NOT NULL").Count(&stats.MigratedRecords)

	// 未迁移记录数
	m.db.Table(tableName).Where("tenant_id IS NULL").Count(&stats.UnmigratedRecords)

	// 外键完整性检查
	m.db.Table(tableName).
		Joins(fmt.Sprintf("LEFT JOIN tenants ON tenants.tenant_id = %s.tenant_id", tableName)).
		Where("tenants.tenant_id IS NULL").
		Count(&stats.InvalidRecords)

	// 计算迁移进度
	if stats.TotalRecords > 0 {
		stats.Progress = float64(stats.MigratedRecords) / float64(stats.TotalRecords) * 100
	}

	return stats, nil
}

// MigrationStats 迁移统计信息
type MigrationStats struct {
	TotalRecords      int64   // 总记录数
	MigratedRecords   int64   // 已迁移记录数
	UnmigratedRecords int64   // 未迁移记录数
	InvalidRecords    int64   // 外键无效记录数
	Progress          float64 // 迁移进度（百分比）
}
