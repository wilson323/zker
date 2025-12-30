// backend/domain/tenant/migration/dual_write.go
package migration

import (
	"context"
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// DualWriteVerifier 双写验证器
// 用于验证新旧字段的数据一致性
type DualWriteVerifier struct {
	db     *gorm.DB
	logger logs.CtxLogger
}

// NewDualWriteVerifier 创建双写验证器
func NewDualWriteVerifier(db *gorm.DB) *DualWriteVerifier {
	return &DualWriteVerifier{
		db:     db,
		logger: logs.DefaultLogger(),
	}
}

// VerifyDualWrite 验证双写数据一致性
//
// 参数:
//   - ctx: 上下文
//   - tableName: 表名
//   - duration: 检查最近N小时的数据
//
// 验证逻辑:
//  1. 查询指定时间段内更新的数据
//  2. 对比tenant_id字段与从users表查到的值
//  3. 计算一致性比例
//
// 要求: 数据一致性 ≥ 99.9%
func (v *DualWriteVerifier) VerifyDualWrite(
	ctx context.Context,
	tableName string,
	duration time.Duration,
) error {
	v.logger.CtxInfof(ctx, "[DualWrite] starting dual write verification for table %s (last %s)",
		tableName, duration)

	startTime := time.Now()
	checkTime := time.Now().Add(-duration)

	// 1. 查询指定时间段内更新的数据
	type Record struct {
		ID       string
		UserID   string
		TenantID string
	}

	var records []Record
	err := v.db.Table(tableName).
		Select("id, creator_id as user_id, tenant_id").
		Where("updated_at >= ?", checkTime).
		Find(&records).Error

	if err != nil {
		return fmt.Errorf("查询数据失败: %w", err)
	}

	if len(records) == 0 {
		v.logger.CtxWarnf(ctx, "[DualWrite] no records found in the last %s", duration)
		return nil
	}

	// 2. 对比新旧字段数据
	inconsistencyCount := 0
	consistentCount := 0

	for _, record := range records {
		if record.UserID == "" {
			// 没有user_id字段，跳过
			continue
		}

		// 从users表查找正确的tenant_id
		var correctTenantID string
		err := v.db.Table("users").
			Select("tenant_id").
			Where("user_id = ?", record.UserID).
			Pluck("tenant_id", &correctTenantID).Error

		if err != nil {
			v.logger.CtxErrorf(ctx, "[DualWrite] failed to get tenant_id for user %s: %v",
				record.UserID, err)
			inconsistencyCount++
			continue
		}

		// 对比tenant_id
		if record.TenantID != correctTenantID {
			v.logger.CtxErrorf(ctx,
				"[DualWrite] data inconsistency found in record %s: expected=%s, actual=%s",
				record.ID, correctTenantID, record.TenantID)
			inconsistencyCount++
		} else {
			consistentCount++
		}
	}

	// 3. 计算一致性
	totalChecked := consistentCount + inconsistencyCount
	consistencyRate := float64(consistentCount) / float64(totalChecked) * 100

	elapsed := time.Since(startTime)
	v.logger.CtxInfof(ctx,
		"[DualWrite] verification completed: total=%d, consistent=%d, inconsistent=%d, consistency=%.2f%%, duration=%s",
		totalChecked, consistentCount, inconsistencyCount, consistencyRate, elapsed)

	// 4. 检查一致性是否达标
	if consistencyRate < 99.9 {
		return fmt.Errorf("数据一致性低于99.9%%: %.2f%% (%d/%d)",
			consistencyRate, consistentCount, totalChecked)
	}

	v.logger.CtxInfof(ctx, "[DualWrite] ✓ Data consistency passed: %.2f%%", consistencyRate)
	return nil
}

// StartDualWriteMode 启动双写模式
//
// 双写模式说明:
//   - 同时读写新旧字段
//   - 确保数据一致性
//   - 用于迁移验证期间
func (v *DualWriteVerifier) StartDualWriteMode(
	ctx context.Context,
	tableName string,
) error {
	v.logger.CtxInfof(ctx, "[DualWrite] starting dual write mode for table %s", tableName)

	// 注意：实际的双写逻辑需要在应用层实现
	// 这里只是记录日志，提醒开启双写配置

	v.logger.CtxInfof(ctx,
		"[DualWrite] ✓ Dual write mode enabled for table %s. "+
			"Please ensure application-level dual-write is enabled.",
		tableName)

	return nil
}

// StopDualWriteMode 停止双写模式
func (v *DualWriteVerifier) StopDualWriteMode(
	ctx context.Context,
	tableName string,
) error {
	v.logger.CtxInfof(ctx, "[DualWrite] stopping dual write mode for table %s", tableName)

	v.logger.CtxInfof(ctx,
		"[DualWrite] ✓ Dual write mode disabled for table %s. "+
			"Please ensure application-level dual-write is disabled.",
		tableName)

	return nil
}

// MonitorDualWrite 监控双写数据一致性
//
// 定期检查数据一致性，持续监控
func (v *DualWriteVerifier) MonitorDualWrite(
	ctx context.Context,
	tableName string,
	interval time.Duration,
	duration time.Duration,
) error {
	v.logger.CtxInfof(ctx,
		"[DualWrite] starting dual write monitoring: table=%s, interval=%s, duration=%s",
		tableName, interval, duration)

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

			v.logger.CtxInfof(ctx, "[DualWrite] running consistency check #%d...", checkCount)

			err := v.VerifyDualWrite(ctx, tableName, 1*time.Hour)
			if err != nil {
				v.logger.CtxErrorf(ctx, "[DualWrite] consistency check #%d failed: %v",
					checkCount, err)
				// 不返回错误，继续监控
			}
		}
	}

	v.logger.CtxInfof(ctx, "[DualWrite] monitoring completed: %d checks performed", checkCount)
	return nil
}

// DualWriteStats 双写统计信息
type DualWriteStats struct {
	TableName           string    // 表名
	TotalRecords        int64     // 总记录数
	ConsistentRecords   int64     // 一致的记录数
	InconsistentRecords int64     // 不一致的记录数
	ConsistencyRate     float64   // 一致性比例（百分比）
	LastCheckTime       time.Time // 最后检查时间
}

// GetDualWriteStats 获取双写统计信息
func (v *DualWriteVerifier) GetDualWriteStats(
	ctx context.Context,
	tableName string,
	duration time.Duration,
) (*DualWriteStats, error) {
	stats := &DualWriteStats{
		TableName:     tableName,
		LastCheckTime: time.Now(),
	}

	checkTime := time.Now().Add(-duration)

	// 查询总记录数
	err := v.db.Table(tableName).
		Where("updated_at >= ?", checkTime).
		Count(&stats.TotalRecords).Error

	if err != nil {
		return nil, err
	}

	// 查询一致的记录数（通过JOIN验证）
	err = v.db.Raw(`
		SELECT COUNT(*)
		FROM `+tableName+` b
		INNER JOIN users u ON b.creator_id = u.user_id
		WHERE b.updated_at >= ? AND b.tenant_id = u.tenant_id
	`, checkTime).Scan(&stats.ConsistentRecords).Error

	if err != nil {
		return nil, err
	}

	// 不一致的记录数
	stats.InconsistentRecords = stats.TotalRecords - stats.ConsistentRecords

	// 一致性比例
	if stats.TotalRecords > 0 {
		stats.ConsistencyRate = float64(stats.ConsistentRecords) / float64(stats.TotalRecords) * 100
	}

	return stats, nil
}
