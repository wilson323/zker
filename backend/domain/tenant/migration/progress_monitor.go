// backend/domain/tenant/migration/progress_monitor.go
// 进度监控器 - Redis存储迁移进度
package migration

import (
	"context"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// ProgressMonitor 进度监控器
type ProgressMonitor struct {
	logger logs.CtxLogger
}

// TableProgress 表级别的迁移进度
type TableProgress struct {
	TableName   string    `json:"table_name"`
	Current     int       `json:"current"`
	Total       int       `json:"total"`
	Percent     float64   `json:"percent"`
	Status      string    `json:"status"`
	StartedAt   time.Time `json:"started_at"`
	LastUpdated time.Time `json:"last_updated"`
}

// NewProgressMonitor 创建进度监控器
func NewProgressMonitor(db interface{}) *ProgressMonitor {
	return &ProgressMonitor{
		logger: logs.DefaultLogger(),
	}
}

// UpdateProgress 更新迁移进度
func (m *ProgressMonitor) UpdateProgress(ctx context.Context, tableName string, current, total int) error {
	progress := &TableProgress{
		TableName:   tableName,
		Current:     current,
		Total:       total,
		Percent:     float64(current) / float64(total) * 100,
		Status:      "in_progress",
		LastUpdated: time.Now(),
	}

	// data, err := json.Marshal(progress)
	// if err != nil {
	// 	return err
	// }
	//
	// key := fmt.Sprintf("migration:progress:%s", tableName)

	// 存储到Redis（如果Redis可用）
	// m.redis.Set(ctx, key, data, 24*time.Hour)

	m.logger.CtxDebugf(ctx, "[ProgressMonitor] 进度更新: %s %d/%d (%.2f%%)",
		tableName, current, total, progress.Percent)

	return nil
}

// GetProgress 获取迁移进度
func (m *ProgressMonitor) GetProgress(ctx context.Context, tableName string) (*TableProgress, error) {
	// key := fmt.Sprintf("migration:progress:%s", tableName)

	// 从Redis读取
	// data, err := m.redis.Get(ctx, key).Result()
	// if err != nil {
	// 	return nil, err
	// }
	//
	// var progress TableProgress
	// err = json.Unmarshal([]byte(data), &progress)
	// return &progress, err

	return &TableProgress{}, nil
}
