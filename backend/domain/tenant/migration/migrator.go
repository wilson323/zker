// backend/domain/tenant/migration/migrator.go
// 租户数据迁移器 - 批量迁移历史数据，分配tenant_id
package migration

import (
	"context"
	"fmt"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// TenantMigrator 租户迁移器
type TenantMigrator struct {
	db            *gorm.DB
	batchSize     int    // 批处理大小
	defaultTenant string // 默认租户ID
	monitor       *ProgressMonitor
	checkpoint    *CheckpointManager
	logger        logs.CtxLogger
}

// MigrationResult 迁移结果
type MigrationResult struct {
	TableName       string    `json:"table_name"`
	TotalRecords    int64     `json:"total_records"`
	MigratedRecords int64     `json:"migrated_records"`
	FailedRecords   int64     `json:"failed_records"`
	ProgressPercent float64   `json:"progress_percent"`
	StartedAt       time.Time `json:"started_at"`
	FinishedAt      time.Time `json:"finished_at"`
	DurationSeconds int64     `json:"duration_seconds"`
	Status          string    `json:"status"` // pending, in_progress, completed, failed
}

// NewTenantMigrator 创建迁移器
func NewTenantMigrator(db *gorm.DB, batchSize int, defaultTenant string) *TenantMigrator {
	return &TenantMigrator{
		db:            db,
		batchSize:     batchSize,
		defaultTenant: defaultTenant,
		monitor:       NewProgressMonitor(db),
		checkpoint:    NewCheckpointManager(db),
		logger:        logs.DefaultLogger(),
	}
}

// MigrateTable 迁移单张表
func (m *TenantMigrator) MigrateTable(ctx context.Context, tableName string) (*MigrationResult, error) {
	m.logger.CtxInfof(ctx, "[Migrator] 开始迁移表: %s", tableName)

	result := &MigrationResult{
		TableName: tableName,
		Status:    "in_progress",
		StartedAt: time.Now(),
	}

	// 1. 获取总记录数
	totalRecords, err := m.getTotalRecords(ctx, tableName)
	if err != nil {
		return nil, fmt.Errorf("获取总记录数失败: %w", err)
	}
	result.TotalRecords = totalRecords

	m.logger.CtxInfof(ctx, "[Migrator] 表 %s 共有 %d 条记录", tableName, totalRecords)

	// 2. 加载检查点（支持断点续传）
	lastID, err := m.checkpoint.LoadCheckpoint(ctx, tableName)
	if err != nil {
		m.logger.CtxWarnf(ctx, "[Migrator] 加载检查点失败，从头开始: %v", err)
		lastID = 0
	}

	// 3. 批量迁移
	migrated := int64(0)
	failed := int64(0)

	for {
		// 查询一批记录
		records, err := m.getBatch(ctx, tableName, lastID, m.batchSize)
		if err != nil {
			m.logger.CtxErrorf(ctx, "[Migrator] 查询批次失败: %v", err)
			result.Status = "failed"
			return result, err
		}

		if len(records) == 0 {
			break // 所有记录已迁移完成
		}

		// 更新这批记录的 tenant_id
		batchMigrated, batchFailed := m.updateBatch(ctx, tableName, records)

		migrated += batchMigrated
		failed += batchFailed
		result.MigratedRecords = migrated
		result.FailedRecords = failed
		result.ProgressPercent = float64(migrated) / float64(totalRecords) * 100

		// 更新进度
		if err := m.monitor.UpdateProgress(ctx, tableName, int(migrated), int(totalRecords)); err != nil {
			m.logger.CtxWarnf(ctx, "[Migrator] 更新进度失败: %v", err)
		}

		// 保存检查点
		if err := m.checkpoint.SaveCheckpoint(ctx, tableName, lastID); err != nil {
			m.logger.CtxWarnf(ctx, "[Migrator] 保存检查点失败: %v", err)
		}

		// 记录日志
		m.logger.CtxInfof(ctx, "[Migrator] 表 %s 进度: %d/%d (%.2f%%)",
			tableName, migrated, totalRecords, result.ProgressPercent)

		// 更新 lastID
		lastID = records[len(records)-1].ID
	}

	// 4. 迁移完成
	result.Status = "completed"
	result.FinishedAt = time.Now()
	result.DurationSeconds = int64(time.Since(result.StartedAt).Seconds())

	m.logger.CtxInfof(ctx, "[Migrator] 表 %s 迁移完成: 成功 %d, 失败 %d, 耗时 %d秒",
		tableName, migrated, failed, result.DurationSeconds)

	return result, nil
}

// getTotalRecords 获取表总记录数
func (m *TenantMigrator) getTotalRecords(ctx context.Context, tableName string) (int64, error) {
	var count int64
	err := m.db.Table(tableName).Count(&count).Error
	return count, err
}

// getBatch 获取一批记录
func (m *TenantMigrator) getBatch(ctx context.Context, tableName string, lastID int64, batchSize int) ([]Record, error) {
	var records []Record

	err := m.db.Table(tableName).
		Where("id > ?", lastID).
		Order("id ASC").
		Limit(batchSize).
		Find(&records).Error

	return records, err
}

// updateBatch 更新一批记录的 tenant_id
func (m *TenantMigrator) updateBatch(ctx context.Context, tableName string, records []Record) (int64, int64) {
	success := int64(0)
	failed := int64(0)

	for _, record := range records {
		// 分配 tenant_id
		tenantID := m.assignTenantID(ctx, record)

		// 更新记录
		err := m.db.Table(tableName).
			Where("id = ?", record.ID).
			Update("tenant_id", tenantID).Error

		if err != nil {
			failed++
			m.logger.CtxErrorf(ctx, "[Migrator] 更新记录失败: id=%d, error=%v", record.ID, err)
		} else {
			success++
		}
	}

	return success, failed
}

// assignTenantID 为记录分配 tenant_id
func (m *TenantMigrator) assignTenantID(ctx context.Context, record Record) string {
	// 根据业务规则分配 tenant_id
	// 1. 如果记录有 user_id，查询用户所属租户
	// 2. 否则使用默认租户

	// 简化实现：使用默认租户
	return m.defaultTenant
}

// Record 通用记录结构
type Record struct {
	ID     int64  `gorm:"column:id"`
	UserID *int64 `gorm:"column:user_id"`
}
