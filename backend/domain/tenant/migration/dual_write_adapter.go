// backend/domain/tenant/migration/dual_write_adapter.go
// 双写适配器 - 同时写入旧表和新表
package migration

import (
	"context"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// DualWriteAdapter 双写适配器
type DualWriteAdapter struct {
	legacyDB *gorm.DB // 旧库（不含tenant_id）
	newDB    *gorm.DB // 新库（含tenant_id）
	enabled  bool     // 是否启用双写
	logger   logs.CtxLogger
}

// NewDualWriteAdapter 创建双写适配器
func NewDualWriteAdapter(legacyDB, newDB *gorm.DB) *DualWriteAdapter {
	return &DualWriteAdapter{
		legacyDB: legacyDB,
		newDB:    newDB,
		enabled:  false, // 默认关闭，通过配置启用
		logger:   logs.DefaultLogger(),
	}
}

// Enable 启用双写
func (d *DualWriteAdapter) Enable(ctx context.Context) {
	d.enabled = true
	d.logger.CtxInfof(ctx, "[DualWriteAdapter] 双写已启用")
}

// Disable 禁用双写
func (d *DualWriteAdapter) Disable(ctx context.Context) {
	d.enabled = false
	d.logger.CtxInfof(ctx, "[DualWriteAdapter] 双写已禁用")
}

// CreateBotWithDualWrite 双写创建Bot
func (d *DualWriteAdapter) CreateBotWithDualWrite(ctx context.Context, bot interface{}) error {
	// 1. 写入旧库（同步，主路径）
	if err := d.writeToLegacy(ctx, bot); err != nil {
		return err
	}

	// 2. 异步写入新库
	if d.enabled {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := d.writeToNew(ctx, bot); err != nil {
				d.logger.CtxErrorf(ctx, "[DualWriteAdapter] 写入新库失败: %v", err)
				// 记录失败日志，后续补偿
			}
		}()
	}

	return nil
}

// CreateWithDualWrite 通用的双写创建方法
func (d *DualWriteAdapter) CreateWithDualWrite(ctx context.Context, tableName string, data interface{}) error {
	// 1. 写入旧库（同步，主路径）
	if err := d.legacyDB.WithContext(ctx).Table(tableName).Create(data).Error; err != nil {
		return err
	}

	// 2. 异步写入新库
	if d.enabled {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := d.newDB.WithContext(ctx).Table(tableName).Create(data).Error; err != nil {
				d.logger.CtxErrorf(ctx, "[DualWriteAdapter] 写入新库失败: %v", err)
			}
		}()
	}

	return nil
}

// writeToLegacy 写入旧库
func (d *DualWriteAdapter) writeToLegacy(ctx context.Context, bot interface{}) error {
	return d.legacyDB.WithContext(ctx).Create(bot).Error
}

// writeToNew 写入新库
func (d *DualWriteAdapter) writeToNew(ctx context.Context, bot interface{}) error {
	return d.newDB.WithContext(ctx).Create(bot).Error
}

// UpdateBotWithDualWrite 双写更新Bot
func (d *DualWriteAdapter) UpdateBotWithDualWrite(ctx context.Context, bot interface{}) error {
	// 1. 更新旧库
	if err := d.updateLegacy(ctx, bot); err != nil {
		return err
	}

	// 2. 异步更新新库
	if d.enabled {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := d.updateNew(ctx, bot); err != nil {
				d.logger.CtxErrorf(ctx, "[DualWriteAdapter] 更新新库失败: %v", err)
			}
		}()
	}

	return nil
}

// updateLegacy 更新旧库
func (d *DualWriteAdapter) updateLegacy(ctx context.Context, bot interface{}) error {
	return d.legacyDB.WithContext(ctx).Save(bot).Error
}

// updateNew 更新新库
func (d *DualWriteAdapter) updateNew(ctx context.Context, bot interface{}) error {
	return d.newDB.WithContext(ctx).Save(bot).Error
}

// DeleteBotWithDualWrite 双写删除Bot（软删除）
func (d *DualWriteAdapter) DeleteBotWithDualWrite(ctx context.Context, bot interface{}) error {
	// 1. 软删除旧库（设置deleted_at）
	if err := d.softDeleteLegacy(ctx, bot); err != nil {
		return err
	}

	// 2. 异步软删除新库
	if d.enabled {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := d.softDeleteNew(ctx, bot); err != nil {
				d.logger.CtxErrorf(ctx, "[DualWriteAdapter] 软删除新库失败: %v", err)
			}
		}()
	}

	return nil
}

// softDeleteLegacy 软删除旧库
func (d *DualWriteAdapter) softDeleteLegacy(ctx context.Context, bot interface{}) error {
	return d.legacyDB.WithContext(ctx).Table("bots").Delete(bot).Error
}

// softDeleteNew 软删除新库
func (d *DualWriteAdapter) softDeleteNew(ctx context.Context, bot interface{}) error {
	return d.newDB.WithContext(ctx).Table("bots").Delete(bot).Error
}
