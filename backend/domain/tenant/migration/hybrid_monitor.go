// backend/domain/tenant/migration/hybrid_monitor.go
// 混合监控器 - Redis降级 + 性能监控指标
package migration

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// HybridProgressMonitor 混合进度监控器（Redis降级）
type HybridProgressMonitor struct {
	redis    *redis.Client
	db       *gorm.DB
	logger   logs.CtxLogger
	fallback bool
	stats    *PerformanceStats
}

// PerformanceStats 性能统计
type PerformanceStats struct {
	TotalMigrated     int64
	TotalDuration     time.Duration
	AvgRate           float64 // 记录/秒
	CurrentRate       float64 // 当前速率
	PeakRate          float64 // 峰值速率
	DBQueryCount      int64
	RedisQueryCount   int64
	RedisFailureCount int64
}

// NewHybridProgressMonitor 创建混合监控器
func NewHybridProgressMonitor(redisClient *redis.Client, db *gorm.DB) *HybridProgressMonitor {
	return &HybridProgressMonitor{
		redis:    redisClient,
		db:       db,
		logger:   logs.DefaultLogger(),
		fallback: false,
		stats:    &PerformanceStats{},
	}
}

// UpdateProgress 更新进度（Redis优先，数据库降级）
func (m *HybridProgressMonitor) UpdateProgress(ctx context.Context, tableName string, current, total int) error {
	progress := &MigrationProgress{
		CurrentTable:    tableName,
		MigratedRecords: int64(current),
		TotalRecords:    int64(total),
		ProgressPercent: float64(current) / float64(total) * 100,
		Status:          "in_progress",
	}

	data, err := json.Marshal(progress)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("migration:progress:%s", tableName)

	// 优先使用Redis
	if m.redis != nil && !m.fallback {
		err := m.redis.Set(ctx, key, data, 24*time.Hour).Err()
		if err == nil {
			m.stats.RedisQueryCount++
			return nil
		}

		// Redis失败，降级到数据库
		m.logger.CtxWarnf(ctx, "[Monitor] Redis不可用，降级到数据库")
		m.fallback = true
		m.stats.RedisFailureCount++
	}

	// 降级：存储到数据库
	return m.saveToDatabase(ctx, progress)
}

// saveToDatabase 保存到数据库
func (m *HybridProgressMonitor) saveToDatabase(ctx context.Context, progress *MigrationProgress) error {
	// 使用内存表或临时表存储进度
	return m.db.Table("migration_progress").Save(progress).Error
}

// GetProgress 获取进度（自动选择数据源）
func (m *HybridProgressMonitor) GetProgress(ctx context.Context, tableName string) (*MigrationProgress, error) {
	key := fmt.Sprintf("migration:progress:%s", tableName)

	// 优先从Redis读取
	if !m.fallback && m.redis != nil {
		data, err := m.redis.Get(ctx, key).Result()
		if err == nil {
			var progress MigrationProgress
			if err := json.Unmarshal([]byte(data), &progress); err == nil {
				m.stats.RedisQueryCount++
				return &progress, nil
			}
		}
	}

	// 降级：从数据库读取
	var progress MigrationProgress
	err := m.db.Table("migration_progress").
		Where("table_name = ?", tableName).
		First(&progress).Error

	m.stats.DBQueryCount++
	return &progress, err
}

// RecordPerformance 记录性能指标
func (m *HybridProgressMonitor) RecordPerformance(ctx context.Context, tableName string, migrated int64, duration time.Duration) {
	// 计算当前速率
	rate := float64(migrated) / duration.Seconds()

	m.stats.TotalMigrated += migrated
	m.stats.TotalDuration += duration
	m.stats.CurrentRate = rate

	// 更新平均速率
	if m.stats.TotalDuration.Seconds() > 0 {
		m.stats.AvgRate = float64(m.stats.TotalMigrated) / m.stats.TotalDuration.Seconds()
	}

	// 更新峰值速率
	if rate > m.stats.PeakRate {
		m.stats.PeakRate = rate
	}

	// 记录到Prometheus（如果集成）
	// prometheus.RecordMigrationRate(tableName, rate)
}

// GetPerformanceStats 获取性能统计
func (m *HybridProgressMonitor) GetPerformanceStats() *PerformanceStats {
	return m.stats
}

// LogPerformanceStats 输出性能统计
func (m *HybridProgressMonitor) LogPerformanceStats(ctx context.Context, tableName string) {
	stats := m.stats

	m.logger.CtxInfof(ctx, "[Monitor] ========== 性能统计 (%s) ==========", tableName)
	m.logger.CtxInfof(ctx, "[Monitor] 已迁移: %d 条", stats.TotalMigrated)
	m.logger.CtxInfof(ctx, "[Monitor] 总耗时: %dms", stats.TotalDuration.Milliseconds())
	m.logger.CtxInfof(ctx, "[Monitor] 平均速率: %.2f 条/秒", stats.AvgRate)
	m.logger.CtxInfof(ctx, "[Monitor] 当前速率: %.2f 条/秒", stats.CurrentRate)
	m.logger.CtxInfof(ctx, "[Monitor] 峰值速率: %.2f 条/秒", stats.PeakRate)
	m.logger.CtxInfof(ctx, "[Monitor] 数据库查询: %d 次", stats.DBQueryCount)
	m.logger.CtxInfof(ctx, "[Monitor] Redis查询: %d 次", stats.RedisQueryCount)
	m.logger.CtxInfof(ctx, "[Monitor] Redis失败: %d 次", stats.RedisFailureCount)
	m.logger.CtxInfof(ctx, "[Monitor] ========================================")
}

// IsRedisAvailable 检查Redis是否可用
func (m *HybridProgressMonitor) IsRedisAvailable(ctx context.Context) bool {
	if m.redis == nil {
		return false
	}

	err := m.redis.Ping(ctx).Err()
	return err == nil
}
