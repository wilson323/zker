// backend/domain/tenant/migration/checkpoint_manager.go
// 检查点管理器 - 支持断点续传
package migration

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// CheckpointManager 检查点管理器
type CheckpointManager struct {
	db *gorm.DB
}

// MigrationCheckpoint 迁移检查点
type MigrationCheckpoint struct {
	ID        int64     `gorm:"primaryKey"`
	TableName string    `gorm:"column:table_name"`
	LastID    int64     `gorm:"column:last_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// NewCheckpointManager 创建检查点管理器
func NewCheckpointManager(db *gorm.DB) *CheckpointManager {
	return &CheckpointManager{db: db}
}

// SaveCheckpoint 保存检查点
func (c *CheckpointManager) SaveCheckpoint(ctx context.Context, tableName string, lastID int64) error {
	checkpoint := &MigrationCheckpoint{
		TableName: tableName,
		LastID:    lastID,
		UpdatedAt: time.Now(),
	}

	// 使用 ON DUPLICATE KEY UPDATE
	return c.db.Save(checkpoint).Error
}

// LoadCheckpoint 加载检查点
func (c *CheckpointManager) LoadCheckpoint(ctx context.Context, tableName string) (int64, error) {
	var checkpoint MigrationCheckpoint
	err := c.db.Where("table_name = ?", tableName).First(&checkpoint).Error
	if err != nil {
		return 0, err // 没有检查点，从头开始
	}
	return checkpoint.LastID, nil
}

// ClearCheckpoint 清除检查点
func (c *CheckpointManager) ClearCheckpoint(ctx context.Context, tableName string) error {
	return c.db.Where("table_name = ?", tableName).Delete(&MigrationCheckpoint{}).Error
}
