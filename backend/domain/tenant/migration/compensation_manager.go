// backend/domain/tenant/migration/compensation_manager.go
// 失败补偿管理器 - 集成双写失败处理
package migration

import (
	"context"
	"encoding/json"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// CompensationManager 补偿管理器
type CompensationManager struct {
	db                  *gorm.DB
	logger              logs.CtxLogger
	compensationService *CompensationService
	maxRetryAttempts    int
	retryInterval       time.Duration
	enabled             bool
}

// NewCompensationManager 创建补偿管理器
func NewCompensationManager(db *gorm.DB) *CompensationManager {
	return &CompensationManager{
		db:                  db,
		logger:              logs.DefaultLogger(),
		compensationService: NewCompensationService(db),
		maxRetryAttempts:    3,
		retryInterval:       1 * time.Minute,
		enabled:             true,
	}
}

// Enable 启用补偿
func (m *CompensationManager) Enable() {
	m.enabled = true
	m.logger.CtxInfof(context.Background(), "[Compensation] 补偿机制已启用")
}

// Disable 禁用补偿
func (m *CompensationManager) Disable() {
	m.enabled = false
	m.logger.CtxInfof(context.Background(), "[Compensation] 补偿机制已禁用")
}

// RecordFailedWrite 记录失败的写入
func (m *CompensationManager) RecordFailedWrite(ctx context.Context, tableName, operation string, data interface{}, err error) {
	if !m.enabled {
		return
	}

	m.compensationService.RecordFailedWrite(ctx, tableName, operation, data, err)
}

// StartPeriodicCompensation 启动定时补偿任务
func (m *CompensationManager) StartPeriodicCompensation(ctx context.Context) {
	if !m.enabled {
		m.logger.CtxWarnf(ctx, "[Compensation] 补偿机制未启用，跳过定时任务")
		return
	}

	m.logger.CtxInfof(ctx, "[Compensation] 启动定时补偿任务")

	ticker := time.NewTicker(m.retryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			successCount, failedCount := m.compensationService.RetryFailedWrites(ctx)
			if successCount > 0 || failedCount > 0 {
				m.logger.CtxInfof(ctx, "[Compensation] 补偿完成: 成功=%d, 仍失败=%d", successCount, failedCount)
			}

			// 如果仍有失败的记录，增加告警
			if failedCount > 0 {
				m.logger.CtxWarnf(ctx, "[Compensation] 警告: %d 条记录补偿失败，可能需要人工介入", failedCount)
			}

		case <-ctx.Done():
			m.logger.CtxInfof(ctx, "[Compensation] 定时补偿任务已停止")
			return
		}
	}
}

// ManualCompensation 手动触发补偿
func (m *CompensationManager) ManualCompensation(ctx context.Context) (successCount, failedCount int) {
	m.logger.CtxInfof(ctx, "[Compensation] 手动触发补偿")
	return m.compensationService.RetryFailedWrites(ctx)
}

// GetFailedWritesCount 获取失败记录数
func (m *CompensationManager) GetFailedWritesCount(ctx context.Context) (int64, error) {
	var count int64
	err := m.db.Model(&FailedWrite{}).Where("retry_count < ?", m.maxRetryAttempts).Count(&count).Error
	return count, err
}

// GetFailedWrites 获取失败记录列表
func (m *CompensationManager) GetFailedWrites(ctx context.Context, limit int) ([]FailedWrite, error) {
	var failedWrites []FailedWrite
	err := m.db.Where("retry_count < ?", m.maxRetryAttempts).
		Order("created_at DESC").
		Limit(limit).
		Find(&failedWrites).Error
	return failedWrites, err
}

// ClearFailedWrites 清理已成功的失败记录
func (m *CompensationManager) ClearFailedWrites(ctx context.Context, beforeTime time.Time) (int64, error) {
	result := m.db.Where("created_at < ? AND retry_count >= ?", beforeTime, m.maxRetryAttempts).Delete(&FailedWrite{})
	if result.Error != nil {
		return 0, result.Error
	}

	m.logger.CtxInfof(ctx, "[Compensation] 清理失败记录: %d 条", result.RowsAffected)
	return result.RowsAffected, nil
}

// CompensationStats 补偿统计
type CompensationStats struct {
	TotalFailed        int64     `json:"total_failed"`
	PendingRetry       int64     `json:"pending_retry"`
	SuccessCompensated int64     `json:"success_compensated"`
	FailedPermanently  int64     `json:"failed_permanently"`
	LastCompensationAt time.Time `json:"last_compensation_at"`
}

// GetStats 获取补偿统计
func (m *CompensationManager) GetStats(ctx context.Context) (*CompensationStats, error) {
	stats := &CompensationStats{}

	// 总失败数
	m.db.Model(&FailedWrite{}).Count(&stats.TotalFailed)

	// 待重试数
	m.db.Model(&FailedWrite{}).Where("retry_count < ?", m.maxRetryAttempts).Count(&stats.PendingRetry)

	// 已补偿成功数（通过计算retry_count >= 3的记录）
	m.db.Model(&FailedWrite{}).Where("retry_count >= ?", m.maxRetryAttempts).Count(&stats.FailedPermanently)

	// 最后补偿时间
	var lastWrite FailedWrite
	m.db.Order("created_at DESC").First(&lastWrite)
	if lastWrite.ID > 0 {
		stats.LastCompensationAt = lastWrite.CreatedAt
	}

	return stats, nil
}

// ExportFailedWrites 导出失败记录到JSON
func (m *CompensationManager) ExportFailedWrites(ctx context.Context) ([]byte, error) {
	failedWrites, err := m.GetFailedWrites(ctx, 1000)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(failedWrites, "", "  ")
}
