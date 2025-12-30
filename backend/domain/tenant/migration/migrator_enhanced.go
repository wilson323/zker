// backend/domain/tenant/migration/migrator_enhanced.go
// 租户数据迁移器增强版 - 并发控制和性能优化
package migration

import (
	"context"
	"fmt"
	"sync"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// TenantMigratorEnhanced 增强版迁移器（并发优化）
type TenantMigratorEnhanced struct {
	db             *gorm.DB
	batchSize      int
	defaultTenant  string
	maxConcurrency int // 最大并发数
	monitor        *ProgressMonitor
	checkpoint     *CheckpointManager
	logger         logs.CtxLogger
	semaphore      chan struct{} // 并发信号量
	stats          *MigrationStats
	totalRecords   int64     // 总记录数缓存
	migrated       int64     // 已迁移记录数缓存
	startedAt      time.Time // 开始时间
}

// NewTenantMigratorEnhanced 创建增强版迁移器
func NewTenantMigratorEnhanced(db *gorm.DB, batchSize int, defaultTenant string, maxConcurrency int) *TenantMigratorEnhanced {
	return &TenantMigratorEnhanced{
		db:             db,
		batchSize:      batchSize,
		defaultTenant:  defaultTenant,
		maxConcurrency: maxConcurrency,
		monitor:        NewProgressMonitor(db),
		checkpoint:     NewCheckpointManager(db),
		logger:         logs.DefaultLogger(),
		semaphore:      make(chan struct{}, maxConcurrency),
		stats:          &MigrationStats{},
		startedAt:      time.Now(),
	}
}

// MigrateTableEnhanced 增强版表迁移（并发优化）
func (m *TenantMigratorEnhanced) MigrateTableEnhanced(ctx context.Context, tableName string) (*MigrationResult, error) {
	m.logger.CtxInfof(ctx, "[Migrator] 开始迁移表: %s (并发度: %d)", tableName, m.maxConcurrency)

	result := &MigrationResult{
		TableName: tableName,
		Status:    "in_progress",
		StartedAt: time.Now(),
	}
	m.startedAt = time.Now()

	// 1. 获取总记录数
	totalRecords, err := m.getTotalRecords(ctx, tableName)
	if err != nil {
		return nil, fmt.Errorf("获取总记录数失败: %w", err)
	}
	result.TotalRecords = totalRecords
	m.totalRecords = totalRecords
	m.migrated = 0

	m.logger.CtxInfof(ctx, "[Migrator] 表 %s 共有 %d 条记录", tableName, totalRecords)

	// 2. 加载检查点
	lastID, err := m.checkpoint.LoadCheckpoint(ctx, tableName)
	if err != nil {
		lastID = 0
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	migrated := int64(0)
	failed := int64(0)

	// 3. 批量并发迁移
	for {
		// 获取一批记录ID
		ids, err := m.getBatchIDs(ctx, tableName, lastID, m.batchSize*10)
		if err != nil || len(ids) == 0 {
			break
		}

		// 分批处理（每批batchSize条）
		for i := 0; i < len(ids); i += m.batchSize {
			end := i + m.batchSize
			if end > len(ids) {
				end = len(ids)
			}
			batch := ids[i:end]

			// 使用信号量控制并发
			m.semaphore <- struct{}{}
			wg.Add(1)

			go func(batchIDs []int64) {
				defer func() {
					<-m.semaphore
					wg.Done()
				}()

				// 更新这批记录
				batchMigrated, batchFailed := m.updateBatch(ctx, tableName, batchIDs)

				mu.Lock()
				migrated += batchMigrated
				failed += batchFailed
				result.MigratedRecords = migrated
				result.FailedRecords = failed
				result.ProgressPercent = float64(migrated) / float64(totalRecords) * 100

				// 动态调整批次大小（基于性能）
				m.adjustBatchSize(ctx, migrated, result.TotalRecords)
				mu.Unlock()

				// 更新进度
				m.monitor.UpdateProgress(ctx, tableName, int(migrated), int(totalRecords))

				// 保存检查点
				m.checkpoint.SaveCheckpoint(ctx, tableName, batchIDs[len(batchIDs)-1])
			}(batch)
		}

		lastID = ids[len(ids)-1]
	}

	wg.Wait()

	// 4. 完成
	result.Status = "completed"
	result.FinishedAt = time.Now()
	result.DurationSeconds = int64(time.Since(result.StartedAt).Seconds())

	m.logStats(ctx, tableName, result)

	return result, nil
}

// getBatchIDs 获取一批记录ID
func (m *TenantMigratorEnhanced) getBatchIDs(ctx context.Context, tableName string, lastID int64, limit int) ([]int64, error) {
	var ids []int64
	err := m.db.Table(tableName).
		Select("id").
		Where("id > ?", lastID).
		Order("id ASC").
		Limit(limit).
		Pluck("id", &ids).Error
	return ids, err
}

// getTotalRecords 获取表总记录数
func (m *TenantMigratorEnhanced) getTotalRecords(ctx context.Context, tableName string) (int64, error) {
	var count int64
	err := m.db.Table(tableName).Count(&count).Error
	return count, err
}

// updateBatch 更新一批记录
func (m *TenantMigratorEnhanced) updateBatch(ctx context.Context, tableName string, ids []int64) (int64, int64) {
	success := int64(0)
	failed := int64(0)

	for _, id := range ids {
		tenantID := m.assignTenantID(ctx, id)

		err := m.db.Table(tableName).
			Where("id = ?", id).
			Update("tenant_id", tenantID).Error

		if err != nil {
			failed++
		} else {
			success++
		}
	}

	return success, failed
}

// assignTenantID 分配tenant_id
func (m *TenantMigratorEnhanced) assignTenantID(ctx context.Context, id int64) string {
	// 根据user_id查询租户（简化实现）
	return m.defaultTenant
}

// adjustBatchSize 动态调整批次大小
func (m *TenantMigratorEnhanced) adjustBatchSize(ctx context.Context, migrated, total int64) {
	elapsed := time.Since(m.startedAt)
	if elapsed.Seconds() == 0 {
		return
	}

	// 计算当前速率
	rate := float64(migrated) / elapsed.Seconds()

	// 动态调整
	if rate > 10000 && m.batchSize < 2000 {
		m.batchSize = 2000 // 加速
		m.logger.CtxInfof(ctx, "[Migrator] 性能良好，批次大小调整为: %d", m.batchSize)
	} else if rate < 1000 && m.batchSize > 500 {
		m.batchSize = 500 // 减速
		m.logger.CtxInfof(ctx, "[Migrator] 性能较低，批次大小调整为: %d", m.batchSize)
	}
}

// logStats 输出统计信息
func (m *TenantMigratorEnhanced) logStats(ctx context.Context, tableName string, result *MigrationResult) {
	duration := time.Since(m.startedAt)
	rate := float64(result.MigratedRecords) / duration.Seconds()

	m.logger.CtxInfof(ctx, "[Migrator] ========== 迁移统计 (%s) ==========", tableName)
	m.logger.CtxInfof(ctx, "[Migrator] 总记录数: %d", result.TotalRecords)
	m.logger.CtxInfof(ctx, "[Migrator] 已迁移: %d", result.MigratedRecords)
	m.logger.CtxInfof(ctx, "[Migrator] 失败: %d", result.FailedRecords)
	m.logger.CtxInfof(ctx, "[Migrator] 进度: %.2f%%", result.ProgressPercent)
	m.logger.CtxInfof(ctx, "[Migrator] 总耗时: %dms", duration.Milliseconds())
	m.logger.CtxInfof(ctx, "[Migrator] 迁移速率: %.2f 条/秒", rate)
	m.logger.CtxInfof(ctx, "[Migrator] ======================================")
}
