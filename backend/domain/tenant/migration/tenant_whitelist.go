// backend/domain/tenant/migration/tenant_whitelist.go
// 租户白名单 - 灰度发布租户控制
package migration

import (
	"context"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// TenantWhitelist 租户白名单
type TenantWhitelist struct {
	db     *gorm.DB
	logger logs.CtxLogger
}

// WhitelistEntry 白名单条目
type WhitelistEntry struct {
	ID        int64     `gorm:"primaryKey"`
	TenantID  string    `gorm:"column:tenant_id;type:varchar(36);not null;uniqueIndex"`
	Reason    string    `gorm:"column:reason;type:varchar(500)"`
	Enabled   bool      `gorm:"column:enabled;default:true"`
	CreatedBy string    `gorm:"column:created_by;type:varchar(100)"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// NewTenantWhitelist 创建租户白名单
func NewTenantWhitelist(db *gorm.DB) *TenantWhitelist {
	wl := &TenantWhitelist{
		db:     db,
		logger: logs.DefaultLogger(),
	}

	// 自动创建表
	db.AutoMigrate(&WhitelistEntry{})

	return wl
}

// Add 添加租户到白名单
func (w *TenantWhitelist) Add(ctx context.Context, tenantID, reason, createdBy string) error {
	entry := &WhitelistEntry{
		TenantID:  tenantID,
		Reason:    reason,
		Enabled:   true,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := w.db.Create(entry).Error; err != nil {
		w.logger.CtxErrorf(ctx, "[Whitelist] 添加失败: tenant_id=%s, 错误=%v", tenantID, err)
		return err
	}

	w.logger.CtxInfof(ctx, "[Whitelist] 租户已添加到白名单: tenant_id=%s, 原因=%s", tenantID, reason)
	return nil
}

// Remove 从白名单移除租户
func (w *TenantWhitelist) Remove(ctx context.Context, tenantID string) error {
	result := w.db.Where("tenant_id = ?", tenantID).Delete(&WhitelistEntry{})
	if result.Error != nil {
		w.logger.CtxErrorf(ctx, "[Whitelist] 移除失败: tenant_id=%s, 错误=%v", tenantID, result.Error)
		return result.Error
	}

	if result.RowsAffected > 0 {
		w.logger.CtxInfof(ctx, "[Whitelist] 租户已从白名单移除: tenant_id=%s", tenantID)
	}

	return nil
}

// Enable 启用租户
func (w *TenantWhitelist) Enable(ctx context.Context, tenantID string) error {
	result := w.db.Model(&WhitelistEntry{}).
		Where("tenant_id = ?", tenantID).
		Update("enabled", true)

	if result.Error != nil {
		return result.Error
	}

	w.logger.CtxInfof(ctx, "[Whitelist] 租户已启用: tenant_id=%s", tenantID)
	return nil
}

// Disable 禁用租户
func (w *TenantWhitelist) Disable(ctx context.Context, tenantID string) error {
	result := w.db.Model(&WhitelistEntry{}).
		Where("tenant_id = ?", tenantID).
		Update("enabled", false)

	if result.Error != nil {
		return result.Error
	}

	w.logger.CtxInfof(ctx, "[Whitelist] 租户已禁用: tenant_id=%s", tenantID)
	return nil
}

// IsEnabled 检查租户是否启用
func (w *TenantWhitelist) IsEnabled(tenantID string) bool {
	var count int64
	w.db.Model(&WhitelistEntry{}).
		Where("tenant_id = ? AND enabled = ?", tenantID, true).
		Count(&count)

	return count > 0
}

// GetAll 获取所有白名单租户
func (w *TenantWhitelist) GetAll(ctx context.Context) ([]WhitelistEntry, error) {
	var entries []WhitelistEntry
	err := w.db.Order("created_at DESC").Find(&entries).Error
	return entries, err
}

// GetCount 获取白名单租户数
func (w *TenantWhitelist) GetCount(ctx context.Context) (int64, error) {
	var count int64
	err := w.db.Model(&WhitelistEntry{}).Where("enabled = ?", true).Count(&count).Error
	return count, err
}

// GetEnabledTenants 获取启用的租户列表
func (w *TenantWhitelist) GetEnabledTenants(ctx context.Context) ([]string, error) {
	var tenantIDs []string
	err := w.db.Model(&WhitelistEntry{}).
		Where("enabled = ?", true).
		Pluck("tenant_id", &tenantIDs).Error

	return tenantIDs, err
}

// BatchAdd 批量添加租户到白名单
func (w *TenantWhitelist) BatchAdd(ctx context.Context, tenantIDs []string, reason, createdBy string) (int, int) {
	successCount := 0
	failedCount := 0

	for _, tenantID := range tenantIDs {
		if err := w.Add(ctx, tenantID, reason, createdBy); err != nil {
			failedCount++
		} else {
			successCount++
		}
	}

	w.logger.CtxInfof(ctx, "[Whitelist] 批量添加完成: 成功=%d, 失败=%d", successCount, failedCount)

	return successCount, failedCount
}

// Clear 清空白名单
func (w *TenantWhitelist) Clear(ctx context.Context) error {
	result := w.db.Where("1 = 1").Delete(&WhitelistEntry{})
	if result.Error != nil {
		return result.Error
	}

	w.logger.CtxInfof(ctx, "[Whitelist] 白名单已清空: 删除 %d 条记录", result.RowsAffected)
	return nil
}
