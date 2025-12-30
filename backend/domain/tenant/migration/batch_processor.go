// backend/domain/tenant/migration/batch_processor.go
// 批处理器 - 高效批量迁移工具
package migration

import (
	"context"
	"fmt"
	"sync"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// BatchProcessor 批处理器
type BatchProcessor struct {
	db             *gorm.DB
	batchSize      int
	maxConcurrency int
	defaultTenant  string
	logger         logs.CtxLogger
	checkpoint     *CheckpointManager
	monitor        *HybridProgressMonitor
}

// NewBatchProcessor 创建批处理器
func NewBatchProcessor(db *gorm.DB, batchSize, maxConcurrency int, defaultTenant string) *BatchProcessor {
	return &BatchProcessor{
		db:             db,
		batchSize:      batchSize,
		maxConcurrency: maxConcurrency,
		defaultTenant:  defaultTenant,
		logger:         logs.DefaultLogger(),
		checkpoint:     NewCheckpointManager(db),
		monitor:        NewHybridProgressMonitor(nil, db),
	}
}

// BatchProcessResult 批处理结果
type BatchProcessResult struct {
	TableName        string
	TotalRecords     int64
	MigratedRecords  int64
	FailedRecords    int64
	SuccessRate      float64
	DurationSeconds  float64
	RecordsPerSecond float64
	StartedAt        time.Time
	FinishedAt       time.Time
}

// ProcessTable 处理单张表的迁移
func (p *BatchProcessor) ProcessTable(ctx context.Context, tableName string) (*BatchProcessResult, error) {
	p.logger.CtxInfof(ctx, "[BatchProcessor] 开始迁移表: %s (批次大小: %d, 并发度: %d)",
		tableName, p.batchSize, p.maxConcurrency)

	result := &BatchProcessResult{
		TableName: tableName,
		StartedAt: time.Now(),
	}

	// 1. 获取总记录数
	totalRecords, err := p.getTotalRecords(ctx, tableName)
	if err != nil {
		return nil, fmt.Errorf("获取总记录数失败: %w", err)
	}
	result.TotalRecords = totalRecords

	p.logger.CtxInfof(ctx, "[BatchProcessor] 表 %s 共有 %d 条记录", tableName, totalRecords)

	// 2. 加载检查点
	lastID, err := p.checkpoint.LoadCheckpoint(ctx, tableName)
	if err != nil {
		lastID = 0
		p.logger.CtxWarnf(ctx, "[BatchProcessor] 加载检查点失败，从头开始: %v", err)
	} else if lastID > 0 {
		p.logger.CtxInfof(ctx, "[BatchProcessor] 从检查点恢复: lastID=%d", lastID)
	}

	// 3. 并发处理批次
	var wg sync.WaitGroup
	var mu sync.Mutex
	semaphore := make(chan struct{}, p.maxConcurrency)

	migrated := int64(0)
	failed := int64(0)
	currentLastID := lastID

	// 性能统计
	startTime := time.Now()

	for {
		// 获取一批记录ID
		ids, err := p.getBatchIDs(ctx, tableName, currentLastID, p.batchSize*10)
		if err != nil {
			p.logger.CtxErrorf(ctx, "[BatchProcessor] 获取批次ID失败: %v", err)
			break
		}

		if len(ids) == 0 {
			break // 没有更多记录
		}

		// 分批处理（每批batchSize条）
		for i := 0; i < len(ids); i += p.batchSize {
			end := i + p.batchSize
			if end > len(ids) {
				end = len(ids)
			}
			batch := ids[i:end]

			// 使用信号量控制并发
			semaphore <- struct{}{}
			wg.Add(1)

			go func(batchIDs []int64, batchNum int) {
				defer func() {
					<-semaphore
					wg.Done()
				}()

				// 更新这批记录
				batchMigrated, batchFailed := p.updateBatch(ctx, tableName, batchIDs)

				mu.Lock()
				migrated += batchMigrated
				failed += batchFailed
				result.MigratedRecords = migrated
				result.FailedRecords = failed
				result.SuccessRate = float64(migrated) / float64(migrated+failed) * 100

				// 更新进度
				progress := int(migrated) + int(lastID)
				p.monitor.UpdateProgress(ctx, tableName, progress, int(totalRecords))
				p.monitor.RecordPerformance(ctx, tableName, batchMigrated, time.Since(startTime))

				// 保存检查点
				if len(batchIDs) > 0 {
					lastBatchID := batchIDs[len(batchIDs)-1]
					p.checkpoint.SaveCheckpoint(ctx, tableName, lastBatchID)
				}

				// 记录日志
				if batchNum%10 == 0 {
					p.logger.CtxInfof(ctx, "[BatchProcessor] 表 %s: 已迁移 %d/%d (%.2f%%), 失败 %d",
						tableName, migrated, totalRecords,
						float64(migrated)/float64(totalRecords)*100, failed)
				}
				mu.Unlock()
			}(batch, i/p.batchSize)
		}

		currentLastID = ids[len(ids)-1]
	}

	wg.Wait()

	// 4. 完成
	result.FinishedAt = time.Now()
	result.DurationSeconds = result.FinishedAt.Sub(result.StartedAt).Seconds()
	result.RecordsPerSecond = float64(result.MigratedRecords) / result.DurationSeconds

	p.logBatchResult(ctx, result)

	// TODO: 清理检查点 - 需要在CheckpointManager中添加ClearCheckpoint方法
	// p.checkpoint.ClearCheckpoint(ctx, tableName)

	return result, nil
}

// getTotalRecords 获取表总记录数
func (p *BatchProcessor) getTotalRecords(ctx context.Context, tableName string) (int64, error) {
	var count int64
	err := p.db.Table(tableName).Count(&count).Error
	return count, err
}

// getBatchIDs 获取一批记录ID
func (p *BatchProcessor) getBatchIDs(ctx context.Context, tableName string, lastID int64, limit int) ([]int64, error) {
	var ids []int64
	err := p.db.Table(tableName).
		Select("id").
		Where("id > ?", lastID).
		Order("id ASC").
		Limit(limit).
		Pluck("id", &ids).Error
	return ids, err
}

// updateBatch 更新一批记录
func (p *BatchProcessor) updateBatch(ctx context.Context, tableName string, ids []int64) (int64, int64) {
	success := int64(0)
	failed := int64(0)

	// 批量更新tenant_id
	for _, id := range ids {
		tenantID := p.assignTenantID(ctx, tableName, id)

		err := p.db.Table(tableName).
			Where("id = ?", id).
			Update("tenant_id", tenantID).Error

		if err != nil {
			failed++
			p.logger.CtxErrorf(ctx, "[BatchProcessor] 更新失败: 表=%s, id=%d, 错误=%v",
				tableName, id, err)
		} else {
			success++
		}
	}

	return success, failed
}

// assignTenantID 分配tenant_id
func (p *BatchProcessor) assignTenantID(ctx context.Context, tableName string, id int64) string {
	// TODO: 根据实际业务逻辑分配tenant_id
	// 目前简化为返回默认租户ID
	// 实际应该根据user_id查询关联的tenant_id

	// 示例: 从users表查询
	if tableName == "bots" || tableName == "conversations" {
		// 根据user_id查询tenant_id
		// 这里简化处理，实际应该JOIN查询
		return p.defaultTenant
	}

	return p.defaultTenant
}

// logBatchResult 输出批处理结果
func (p *BatchProcessor) logBatchResult(ctx context.Context, result *BatchProcessResult) {
	p.logger.CtxInfof(ctx, "[BatchProcessor] ========== 批处理完成 (%s) ==========", result.TableName)
	p.logger.CtxInfof(ctx, "[BatchProcessor] 总记录数: %d", result.TotalRecords)
	p.logger.CtxInfof(ctx, "[BatchProcessor] 已迁移: %d", result.MigratedRecords)
	p.logger.CtxInfof(ctx, "[BatchProcessor] 失败: %d", result.FailedRecords)
	p.logger.CtxInfof(ctx, "[BatchProcessor] 成功率: %.2f%%", result.SuccessRate)
	p.logger.CtxInfof(ctx, "[BatchProcessor] 总耗时: %.2f 秒", result.DurationSeconds)
	p.logger.CtxInfof(ctx, "[BatchProcessor] 迁移速率: %.2f 条/秒", result.RecordsPerSecond)
	p.logger.CtxInfof(ctx, "[BatchProcessor] ===========================================")
}

// ProcessAllTables 处理所有表
func (p *BatchProcessor) ProcessAllTables(ctx context.Context, tables []string) map[string]*BatchProcessResult {
	results := make(map[string]*BatchProcessResult)

	for _, tableName := range tables {
		result, err := p.ProcessTable(ctx, tableName)
		if err != nil {
			p.logger.CtxErrorf(ctx, "[BatchProcessor] 迁移表 %s 失败: %v", tableName, err)
			continue
		}
		results[tableName] = result
	}

	p.logAllTablesSummary(ctx, results)

	return results
}

// logAllTablesSummary 输出所有表的汇总
func (p *BatchProcessor) logAllTablesSummary(ctx context.Context, results map[string]*BatchProcessResult) {
	p.logger.CtxInfof(ctx, "[BatchProcessor] ========== 所有表迁移汇总 ==========")

	totalRecords := int64(0)
	totalMigrated := int64(0)
	totalFailed := int64(0)

	for tableName, result := range results {
		totalRecords += result.TotalRecords
		totalMigrated += result.MigratedRecords
		totalFailed += result.FailedRecords

		p.logger.CtxInfof(ctx, "[BatchProcessor] %s: %d/%d (%.2f%%)",
			tableName, result.MigratedRecords, result.TotalRecords,
			float64(result.MigratedRecords)/float64(result.TotalRecords)*100)
	}

	successRate := float64(totalMigrated) / float64(totalMigrated+totalFailed) * 100

	p.logger.CtxInfof(ctx, "[BatchProcessor] 总计: 迁移 %d, 失败 %d, 成功率 %.2f%%",
		totalMigrated, totalFailed, successRate)
	p.logger.CtxInfof(ctx, "[BatchProcessor] ===========================================")
}
