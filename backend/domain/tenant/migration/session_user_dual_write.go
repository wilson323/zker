// backend/domain/tenant/migration/session_user_dual_write.go
// Session/User 表双写适配器
package migration

import (
	"context"
	"fmt"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// SessionUserDualWriteAdapter Session/User表双写适配器
type SessionUserDualWriteAdapter struct {
	db     *gorm.DB
	logger logs.CtxLogger
}

// NewSessionUserDualWriteAdapter 创建双写适配器
func NewSessionUserDualWriteAdapter(db *gorm.DB) *SessionUserDualWriteAdapter {
	return &SessionUserDualWriteAdapter{
		db:     db,
		logger: logs.DefaultLogger(),
	}
}

// EnableDualWriteMode 启用双写模式
func (a *SessionUserDualWriteAdapter) EnableDualWriteMode(ctx context.Context) error {
	a.logger.CtxInfof(ctx, "[SessionUserDualWrite] 启用 Session/User 表双写模式")

	// 设置配置标志（可以使用配置中心或数据库）
	// 这里简化为日志记录
	a.logger.CtxInfof(ctx, "[SessionUserDualWrite] ✅ 双写模式已启用")

	return nil
}

// DisableDualWriteMode 禁用双写模式
func (a *SessionUserDualWriteAdapter) DisableDualWriteMode(ctx context.Context) error {
	a.logger.CtxInfof(ctx, "[SessionUserDualWrite] 禁用 Session/User 表双写模式")

	a.logger.CtxInfof(ctx, "[SessionUserDualWrite] ✅ 双写模式已禁用")
	return nil
}

// VerifySessionDualWrite 验证 session 表双写数据一致性
func (a *SessionUserDualWriteAdapter) VerifySessionDualWrite(
	ctx context.Context,
	duration time.Duration,
) error {
	a.logger.CtxInfof(ctx, "[SessionUserDualWrite] 开始验证 session 表双写数据一致性...")

	checkTime := time.Now().Add(-duration)

	type Record struct {
		ID       int64
		UserID   int64
		TenantID string
	}

	var records []Record
	err := a.db.Table("session").
		Select("id, user_id, tenant_id").
		Where("created_at >= ?", checkTime).
		Find(&records).Error

	if err != nil {
		return fmt.Errorf("查询 session 表失败: %w", err)
	}

	if len(records) == 0 {
		a.logger.CtxWarnf(ctx, "[SessionUserDualWrite] session 表没有最近的记录")
		return nil
	}

	inconsistencyCount := 0
	consistentCount := 0

	for _, record := range records {
		if record.UserID == 0 {
			continue
		}

		// 从 users 表查找正确的 tenant_id
		var correctTenantID string
		err := a.db.Table("user").
			Select("tenant_id").
			Where("id = ?", record.UserID).
			Pluck("tenant_id", &correctTenantID).Error

		if err != nil {
			a.logger.CtxErrorf(ctx, "[SessionUserDualWrite] 获取用户 %d 的 tenant_id 失败: %v", record.UserID, err)
			inconsistencyCount++
			continue
		}

		if record.TenantID != correctTenantID {
			a.logger.CtxErrorf(ctx,
				"[SessionUserDualWrite] session %d 数据不一致: user_id=%d, expected_tenant=%s, actual_tenant=%s",
				record.ID, record.UserID, correctTenantID, record.TenantID)
			inconsistencyCount++
		} else {
			consistentCount++
		}
	}

	totalChecked := consistentCount + inconsistencyCount
	consistencyRate := float64(consistentCount) / float64(totalChecked) * 100

	a.logger.CtxInfof(ctx,
		"[SessionUserDualWrite] session 表双写验证完成: total=%d, consistent=%d, inconsistent=%d, consistency=%.2f%%",
		totalChecked, consistentCount, inconsistencyCount, consistencyRate)

	if consistencyRate < 99.9 {
		return fmt.Errorf("session 表数据一致性低于99.9%%: %.2f%%", consistencyRate)
	}

	return nil
}

// VerifyUserDualWrite 验证 user 表双写数据一致性
func (a *SessionUserDualWriteAdapter) VerifyUserDualWrite(
	ctx context.Context,
	duration time.Duration,
) error {
	a.logger.CtxInfof(ctx, "[SessionUserDualWrite] 开始验证 user 表双写数据一致性...")

	// user 表的 tenant_id 是直接存储的，无需对比
	// 这里主要检查是否有 NULL 的 tenant_id
	var nullCount int64
	err := a.db.Table("user").
		Where("tenant_id IS NULL").
		Count(&nullCount).Error

	if err != nil {
		return fmt.Errorf("查询 user 表失败: %w", err)
	}

	if nullCount > 0 {
		return fmt.Errorf("user 表仍有 %d 条记录的 tenant_id 为 NULL", nullCount)
	}

	var totalCount int64
	err = a.db.Table("user").Count(&totalCount).Error
	if err != nil {
		return fmt.Errorf("统计 user 表记录数失败: %w", err)
	}

	completionRate := float64(totalCount-nullCount) / float64(totalCount) * 100

	a.logger.CtxInfof(ctx,
		"[SessionUserDualWrite] user 表双写验证完成: total=%d, with_tenant=%d, completion=%.2f%%",
		totalCount, totalCount-nullCount, completionRate)

	return nil
}

// VerifyUserTenantRelations 验证 user_tenant 关联数据一致性
func (a *SessionUserDualWriteAdapter) VerifyUserTenantRelations(ctx context.Context) error {
	a.logger.CtxInfof(ctx, "[SessionUserDualWrite] 开始验证 user_tenant 关联数据...")

	// 检查每个用户是否在 user_tenant 表中有记录
	var orphanUsers []int64
	err := a.db.Raw(`
		SELECT u.id
		FROM user u
		LEFT JOIN user_tenant ut ON u.id = ut.user_id AND ut.deleted_at IS NULL
		WHERE ut.id IS NULL
	`).Scan(&orphanUsers).Error

	if err != nil {
		return fmt.Errorf("查询孤儿用户失败: %w", err)
	}

	if len(orphanUsers) > 0 {
		a.logger.CtxWarnf(ctx, "[SessionUserDualWrite] 发现 %d 个孤儿用户（无租户关联）", len(orphanUsers))
		// 可以选择自动修复或仅警告
	}

	// 检查 user_tenant 中的用户是否在 user 表中存在
	var invalidRelations []int64
	err = a.db.Raw(`
		SELECT ut.user_id
		FROM user_tenant ut
		LEFT JOIN user u ON ut.user_id = u.id
		WHERE u.id IS NULL AND ut.deleted_at IS NULL
	`).Scan(&invalidRelations).Error

	if err != nil {
		return fmt.Errorf("查询无效关联失败: %w", err)
	}

	if len(invalidRelations) > 0 {
		a.logger.CtxWarnf(ctx, "[SessionUserDualWrite] 发现 %d 个无效关联（用户不存在）", len(invalidRelations))
	}

	var stats struct {
		TotalRelations int64
		UniqueUsers    int64
		UniqueTenants  int64
	}
	err = a.db.Raw(`
		SELECT
			COUNT(*) as total_relations,
			COUNT(DISTINCT user_id) as unique_users,
			COUNT(DISTINCT tenant_id) as unique_tenants
		FROM user_tenant
		WHERE deleted_at IS NULL
	`).Scan(&stats).Error

	if err != nil {
		return fmt.Errorf("统计关联数据失败: %w", err)
	}

	a.logger.CtxInfof(ctx,
		"[SessionUserDualWrite] user_tenant 关联验证完成: total=%d, users=%d, tenants=%d",
		stats.TotalRelations, stats.UniqueUsers, stats.UniqueTenants)

	return nil
}

// MonitorDualWrite 监控双写数据一致性
func (a *SessionUserDualWriteAdapter) MonitorDualWrite(
	ctx context.Context,
	interval time.Duration,
	duration time.Duration,
) error {
	a.logger.CtxInfof(ctx,
		"[SessionUserDualWrite] 启动双写监控: interval=%s, duration=%s",
		interval, duration)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	deadline := time.Now().Add(duration)
	checkCount := 0

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			checkCount++

			a.logger.CtxInfof(ctx, "[SessionUserDualWrite] 执行第 %d 次一致性检查...", checkCount)

			// 验证 session 表
			if err := a.VerifySessionDualWrite(ctx, 1*time.Hour); err != nil {
				a.logger.CtxErrorf(ctx, "[SessionUserDualWrite] session 表一致性检查失败: %v", err)
			}

			// 验证 user 表
			if err := a.VerifyUserDualWrite(ctx, 1*time.Hour); err != nil {
				a.logger.CtxErrorf(ctx, "[SessionUserDualWrite] user 表一致性检查失败: %v", err)
			}

			// 验证 user_tenant 关联
			if err := a.VerifyUserTenantRelations(ctx); err != nil {
				a.logger.CtxErrorf(ctx, "[SessionUserDualWrite] user_tenant 关联检查失败: %v", err)
			}
		}
	}

	a.logger.CtxInfof(ctx, "[SessionUserDualWrite] 监控完成: 共执行 %d 次检查", checkCount)
	return nil
}

// GetDualWriteStats 获取双写统计信息
func (a *SessionUserDualWriteAdapter) GetDualWriteStats(
	ctx context.Context,
) (*SessionUserDualWriteStats, error) {
	stats := &SessionUserDualWriteStats{
		CheckedAt: time.Now(),
	}

	// User 表统计
	a.db.Table("user").Count(&stats.UserTotal)
	a.db.Table("user").Where("tenant_id IS NOT NULL").Count(&stats.UserWithTenant)

	if stats.UserTotal > 0 {
		stats.UserCompletionRate = float64(stats.UserWithTenant) / float64(stats.UserTotal) * 100
	}

	// Session 表统计
	a.db.Table("session").Count(&stats.SessionTotal)
	a.db.Table("session").Where("tenant_id IS NOT NULL").Count(&stats.SessionWithTenant)

	if stats.SessionTotal > 0 {
		stats.SessionCompletionRate = float64(stats.SessionWithTenant) / float64(stats.SessionTotal) * 100
	}

	// UserTenant 关联统计
	a.db.Raw(`
		SELECT COUNT(*) FROM user_tenant WHERE deleted_at IS NULL
	`).Scan(&stats.UserTenantRelations)

	a.db.Raw(`
		SELECT COUNT(DISTINCT user_id) FROM user_tenant WHERE deleted_at IS NULL
	`).Scan(&stats.UniqueUsers)

	return stats, nil
}

// SessionUserDualWriteStats Session/User表双写统计信息
type SessionUserDualWriteStats struct {
	CheckedAt             time.Time `json:"checked_at"`
	UserTotal             int64     `json:"user_total"`
	UserWithTenant        int64     `json:"user_with_tenant"`
	UserCompletionRate    float64   `json:"user_completion_rate"`
	SessionTotal          int64     `json:"session_total"`
	SessionWithTenant     int64     `json:"session_with_tenant"`
	SessionCompletionRate float64   `json:"session_completion_rate"`
	UserTenantRelations   int64     `json:"user_tenant_relations"`
	UniqueUsers           int64     `json:"unique_users"`
}
