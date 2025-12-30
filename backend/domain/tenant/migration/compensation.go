// backend/domain/tenant/migration/compensation.go
// 失败补偿机制 - 记录和重试失败的写入
package migration

import (
	"context"
	"encoding/json"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// CompensationService 补偿服务
type CompensationService struct {
	db     *gorm.DB
	logger logs.CtxLogger
}

// FailedWrite 失败的写入记录
type FailedWrite struct {
	ID         int64     `gorm:"primaryKey"`
	TableName  string    `gorm:"column:table_name;type:varchar(100)"`
	Operation  string    `gorm:"column:operation;type:varchar(20)"` // INSERT/UPDATE
	Data       string    `gorm:"column:data;type:text"`             // JSON序列化
	RetryCount int       `gorm:"column:retry_count;type:int"`
	LastError  string    `gorm:"column:last_error;type:text"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

// NewCompensationService 创建补偿服务
func NewCompensationService(db *gorm.DB) *CompensationService {
	return &CompensationService{
		db:     db,
		logger: logs.DefaultLogger(),
	}
}

// RecordFailedWrite 记录失败的写入
func (s *CompensationService) RecordFailedWrite(ctx context.Context, tableName, operation string, data interface{}, err error) {
	dataBytes, _ := json.Marshal(data)

	failedWrite := &FailedWrite{
		TableName:  tableName,
		Operation:  operation,
		Data:       string(dataBytes),
		RetryCount: 0,
		LastError:  err.Error(),
	}

	s.db.Create(failedWrite)
	s.logger.CtxWarnf(ctx, "[Compensation] 记录失败写入: 表=%s, 操作=%s, 错误=%v", tableName, operation, err)
}

// RetryFailedWrites 重试失败的写入（后台任务）
func (s *CompensationService) RetryFailedWrites(ctx context.Context) (int, int) {
	s.logger.CtxInfof(ctx, "[Compensation] 开始补偿失败写入")

	var failedWrites []FailedWrite
	s.db.Where("retry_count < 3").Find(&failedWrites)

	successCount := 0
	stillFailedCount := 0

	for _, fw := range failedWrites {
		// 重试写入
		if err := s.retryWrite(ctx, &fw); err == nil {
			s.db.Delete(&fw) // 成功后删除
			successCount++
		} else {
			s.db.Model(&fw).Updates(map[string]interface{}{
				"retry_count": fw.RetryCount + 1,
				"last_error":  err.Error(),
			})
			stillFailedCount++
		}
	}

	s.logger.CtxInfof(ctx, "[Compensation] 补偿完成: 成功=%d, 仍失败=%d", successCount, stillFailedCount)
	return successCount, stillFailedCount
}

// retryWrite 重试单条写入
func (s *CompensationService) retryWrite(ctx context.Context, fw *FailedWrite) error {
	// 反序列化数据
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(fw.Data), &data); err != nil {
		return err
	}

	// 执行操作（根据operation类型）
	switch fw.Operation {
	case "INSERT":
		return s.db.Table(fw.TableName).Create(data).Error
	case "UPDATE":
		return s.db.Table(fw.TableName).Save(data).Error
	default:
		return s.db.Table(fw.TableName).Create(data).Error
	}
}

// StartPeriodicCompensation 启动定时补偿任务（每分钟执行一次）
func (s *CompensationService) StartPeriodicCompensation(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.RetryFailedWrites(ctx)
		case <-ctx.Done():
			return
		}
	}
}
