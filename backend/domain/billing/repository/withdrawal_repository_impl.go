/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
)

// ============================================================
// 提现管理Repository实现
// ============================================================

// withdrawalRequestRepositoryImpl 提现申请仓储实现
type withdrawalRequestRepositoryImpl struct {
	db *gorm.DB
}

// NewWithdrawalRequestRepository 创建提现申请仓储实例
func NewWithdrawalRequestRepository(db *gorm.DB) WithdrawalRequestRepository {
	return &withdrawalRequestRepositoryImpl{db: db}
}

// Create 创建提现申请
func (r *withdrawalRequestRepositoryImpl) Create(ctx context.Context, withdrawal *entity.WithdrawalRequest) error {
	return r.db.WithContext(ctx).Create(withdrawal).Error
}

// GetByID 根据ID获取提现申请
func (r *withdrawalRequestRepositoryImpl) GetByID(ctx context.Context, id uint64) (*entity.WithdrawalRequest, error) {
	var withdrawal entity.WithdrawalRequest
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&withdrawal).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &withdrawal, nil
}

// GetByWithdrawalNumber 根据提现单号获取申请
func (r *withdrawalRequestRepositoryImpl) GetByWithdrawalNumber(ctx context.Context, withdrawalNumber string) (*entity.WithdrawalRequest, error) {
	var withdrawal entity.WithdrawalRequest
	err := r.db.WithContext(ctx).
		Where("withdrawal_number = ? AND deleted_at IS NULL", withdrawalNumber).
		First(&withdrawal).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &withdrawal, nil
}

// Update 更新提现申请
func (r *withdrawalRequestRepositoryImpl) Update(ctx context.Context, withdrawal *entity.WithdrawalRequest) error {
	return r.db.WithContext(ctx).Save(withdrawal).Error
}

// Delete 删除提现申请（软删除）
func (r *withdrawalRequestRepositoryImpl) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.WithdrawalRequest{}, id).Error
}

// List 分页查询提现申请
func (r *withdrawalRequestRepositoryImpl) List(ctx context.Context, filter *WithdrawalRequestFilter) ([]*entity.WithdrawalRequest, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.WithdrawalRequest{})

	// 租户过滤
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	// 状态过滤
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	// 提现方式过滤
	if filter.WithdrawalMethod != nil {
		query = query.Where("withdrawal_method = ?", *filter.WithdrawalMethod)
	}

	// 日期范围过滤
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	// 金额范围过滤
	if filter.MinAmount != nil {
		query = query.Where("amount >= ?", *filter.MinAmount)
	}
	if filter.MaxAmount != nil {
		query = query.Where("amount <= ?", *filter.MaxAmount)
	}

	// 审核人过滤
	if filter.ReviewerID != nil {
		query = query.Where("reviewed_by = ?", *filter.ReviewerID)
	}

	// 提现单号过滤
	if filter.WithdrawalNumber != nil {
		query = query.Where("withdrawal_number LIKE ?", "%"+*filter.WithdrawalNumber+"%")
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var withdrawals []*entity.WithdrawalRequest
	offset := (filter.PageToken - 1) * filter.PageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&withdrawals).Error

	return withdrawals, total, err
}

// GetByTenant 获取租户的提现申请
func (r *withdrawalRequestRepositoryImpl) GetByTenant(
	ctx context.Context,
	tenantID string,
	page, pageSize int,
) ([]*entity.WithdrawalRequest, int64, error) {
	var withdrawals []*entity.WithdrawalRequest
	var total int64

	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&withdrawals).Error

	return withdrawals, total, err
}

// GetPendingWithdrawals 获取待处理的提现申请
func (r *withdrawalRequestRepositoryImpl) GetPendingWithdrawals(ctx context.Context, limit int) ([]*entity.WithdrawalRequest, error) {
	var withdrawals []*entity.WithdrawalRequest
	err := r.db.WithContext(ctx).
		Where("status IN ? AND deleted_at IS NULL", []string{
			entity.WithdrawalStatusPending,
			entity.WithdrawalStatusApproved,
		}).
		Order("created_at ASC").
		Limit(limit).
		Find(&withdrawals).Error
	return withdrawals, err
}

// CountByTenantAndDate 统计租户在指定日期的提现数量
func (r *withdrawalRequestRepositoryImpl) CountByTenantAndDate(ctx context.Context, tenantID string, date string) (int, error) {
	var count int64
	startDate := date + " 00:00:00"
	endDate := date + " 23:59:59"

	err := r.db.WithContext(ctx).
		Model(&entity.WithdrawalRequest{}).
		Where("tenant_id = ? AND created_at >= ? AND created_at <= ? AND deleted_at IS NULL",
			tenantID, startDate, endDate).
		Count(&count).Error

	return int(count), err
}

// SumAmountByTenantAndDate 统计租户在指定日期的提现金额
func (r *withdrawalRequestRepositoryImpl) SumAmountByTenantAndDate(ctx context.Context, tenantID string, date string) (float64, error) {
	var sum float64
	startDate := date + " 00:00:00"
	endDate := date + " 23:59:59"

	err := r.db.WithContext(ctx).
		Model(&entity.WithdrawalRequest{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("tenant_id = ? AND created_at >= ? AND created_at <= ? AND deleted_at IS NULL",
			tenantID, startDate, endDate).
		Scan(&sum).Error

	return sum, err
}

// GetStatisticsByTenant 获取租户提现统计
func (r *withdrawalRequestRepositoryImpl) GetStatisticsByTenant(
	ctx context.Context,
	tenantID string,
	startDate, endDate string,
) (*WithdrawalStatistics, error) {
	stats := &WithdrawalStatistics{}

	startDateTime := startDate + " 00:00:00"
	endDateTime := endDate + " 23:59:59"

	// 总数和总金额
	err := r.db.WithContext(ctx).
		Model(&entity.WithdrawalRequest{}).
		Where("tenant_id = ? AND created_at >= ? AND created_at <= ? AND deleted_at IS NULL",
			tenantID, startDateTime, endDateTime).
		Select(`
			COUNT(*) as total_count,
			COALESCE(SUM(amount), 0) as total_amount,
			COALESCE(SUM(fee), 0) as total_fee
		`).
		Scan(stats).Error
	if err != nil {
		return nil, err
	}

	// 成功统计
	var successStats struct {
		Count int     `json:"count"`
		Amount float64 `json:"amount"`
	}
	err = r.db.WithContext(ctx).
		Model(&entity.WithdrawalRequest{}).
		Where("tenant_id = ? AND created_at >= ? AND created_at <= ? AND status = ? AND deleted_at IS NULL",
			tenantID, startDateTime, endDateTime, entity.WithdrawalStatusSuccess).
		Select("COUNT(*) as count, COALESCE(SUM(amount), 0) as amount").
		Scan(&successStats).Error
	if err == nil {
		stats.SuccessCount = successStats.Count
		stats.SuccessAmount = successStats.Amount
	}

	// 待处理统计
	var pendingStats struct {
		Count int     `json:"count"`
		Amount float64 `json:"amount"`
	}
	err = r.db.WithContext(ctx).
		Model(&entity.WithdrawalRequest{}).
		Where("tenant_id = ? AND created_at >= ? AND created_at <= ? AND status IN ? AND deleted_at IS NULL",
			tenantID, startDateTime, endDateTime, []string{
				entity.WithdrawalStatusPending,
				entity.WithdrawalStatusReviewing,
			}).
		Select("COUNT(*) as count, COALESCE(SUM(amount), 0) as amount").
		Scan(&pendingStats).Error
	if err == nil {
		stats.PendingCount = pendingStats.Count
		stats.PendingAmount = pendingStats.Amount
	}

	// 失败统计
	var failedStats struct {
		Count int     `json:"count"`
		Amount float64 `json:"amount"`
	}
	err = r.db.WithContext(ctx).
		Model(&entity.WithdrawalRequest{}).
		Where("tenant_id = ? AND created_at >= ? AND created_at <= ? AND status = ? AND deleted_at IS NULL",
			tenantID, startDateTime, endDateTime, entity.WithdrawalStatusFailed).
		Select("COUNT(*) as count, COALESCE(SUM(amount), 0) as amount").
		Scan(&failedStats).Error
	if err == nil {
		stats.FailedCount = failedStats.Count
		stats.FailedAmount = failedStats.Amount
	}

	return stats, nil
}

// BatchUpdateStatus 批量更新状态
func (r *withdrawalRequestRepositoryImpl) BatchUpdateStatus(ctx context.Context, ids []uint64, status string) error {
	return r.db.WithContext(ctx).
		Model(&entity.WithdrawalRequest{}).
		Where("id IN ?", ids).
		Update("status", status).Error
}

// ============================================================

// withdrawalLimitRepositoryImpl 提现限额仓储实现
type withdrawalLimitRepositoryImpl struct {
	db *gorm.DB
}

// NewWithdrawalLimitRepository 创建提现限额仓储实例
func NewWithdrawalLimitRepository(db *gorm.DB) WithdrawalLimitRepository {
	return &withdrawalLimitRepositoryImpl{db: db}
}

// Create 创建限额配置
func (r *withdrawalLimitRepositoryImpl) Create(ctx context.Context, limit *entity.WithdrawalLimit) error {
	return r.db.WithContext(ctx).Create(limit).Error
}

// GetByID 根据ID获取限额配置
func (r *withdrawalLimitRepositoryImpl) GetByID(ctx context.Context, id uint64) (*entity.WithdrawalLimit, error) {
	var limit entity.WithdrawalLimit
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&limit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &limit, nil
}

// Update 更新限额配置
func (r *withdrawalLimitRepositoryImpl) Update(ctx context.Context, limit *entity.WithdrawalLimit) error {
	return r.db.WithContext(ctx).Save(limit).Error
}

// Delete 删除限额配置（软删除）
func (r *withdrawalLimitRepositoryImpl) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.WithdrawalLimit{}, id).Error
}

// List 列出所有限额配置
func (r *withdrawalLimitRepositoryImpl) List(ctx context.Context, enabledOnly bool) ([]*entity.WithdrawalLimit, error) {
	var limits []*entity.WithdrawalLimit
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL")
	if enabledOnly {
		query = query.Where("enabled = ?", true)
	}
	err := query.Find(&limits).Error
	return limits, err
}

// GetByMethod 获取指定提现方式的限额配置
func (r *withdrawalLimitRepositoryImpl) GetByMethod(ctx context.Context, withdrawalMethod string) ([]*entity.WithdrawalLimit, error) {
	var limits []*entity.WithdrawalLimit
	err := r.db.WithContext(ctx).
		Where("withdrawal_method = ? OR withdrawal_method IS NULL", withdrawalMethod).
		Where("enabled = ? AND deleted_at IS NULL", true).
		Find(&limits).Error
	return limits, err
}

// GetEffectiveLimit 获取有效的限额配置
func (r *withdrawalLimitRepositoryImpl) GetEffectiveLimit(ctx context.Context, withdrawalMethod string, amount float64) (*entity.WithdrawalLimit, error) {
	var limit entity.WithdrawalLimit
	err := r.db.WithContext(ctx).
		Where("(withdrawal_method = ? OR withdrawal_method IS NULL)", withdrawalMethod).
		Where("enabled = ? AND min_amount <= ? AND (max_amount >= ? OR max_amount = 0)",
			true, amount, amount).
		Where("deleted_at IS NULL").
		Order("withdrawal_method DESC, id ASC").
		First(&limit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &limit, nil
}

// ============================================================

// withdrawalRecordRepositoryImpl 提现流水仓储实现
type withdrawalRecordRepositoryImpl struct {
	db *gorm.DB
}

// NewWithdrawalRecordRepository 创建提现流水仓储实例
func NewWithdrawalRecordRepository(db *gorm.DB) WithdrawalRecordRepository {
	return &withdrawalRecordRepositoryImpl{db: db}
}

// Create 创建流水记录
func (r *withdrawalRecordRepositoryImpl) Create(ctx context.Context, record *entity.WithdrawalRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

// GetByID 根据ID获取流水记录
func (r *withdrawalRecordRepositoryImpl) GetByID(ctx context.Context, id uint64) (*entity.WithdrawalRecord, error) {
	var record entity.WithdrawalRecord
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// GetByWithdrawalID 获取提现的所有流水记录
func (r *withdrawalRecordRepositoryImpl) GetByWithdrawalID(ctx context.Context, withdrawalID uint64) ([]*entity.WithdrawalRecord, error) {
	var records []*entity.WithdrawalRecord
	err := r.db.WithContext(ctx).
		Where("withdrawal_id = ?", withdrawalID).
		Order("created_at ASC").
		Find(&records).Error
	return records, err
}

// List 分页查询流水记录
func (r *withdrawalRecordRepositoryImpl) List(ctx context.Context, filter *WithdrawalRecordFilter) ([]*entity.WithdrawalRecord, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.WithdrawalRecord{})

	// 租户过滤
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	// 提现ID过滤
	if filter.WithdrawalID != nil {
		query = query.Where("withdrawal_id = ?", *filter.WithdrawalID)
	}

	// 交易类型过滤
	if filter.TransactionType != nil {
		query = query.Where("transaction_type = ?", *filter.TransactionType)
	}

	// 日期范围过滤
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var records []*entity.WithdrawalRecord
	offset := (filter.PageToken - 1) * filter.PageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&records).Error

	return records, total, err
}

// BatchCreate 批量创建流水记录
func (r *withdrawalRecordRepositoryImpl) BatchCreate(ctx context.Context, records []*entity.WithdrawalRecord) error {
	if len(records) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(records, 100).Error
}

// GetDailyStatistics 获取每日统计
func (r *withdrawalRecordRepositoryImpl) GetDailyStatistics(
	ctx context.Context,
	tenantID string,
	startDate, endDate string,
) ([]*DailyWithdrawalStatistics, error) {
	var stats []*DailyWithdrawalStatistics

	query := `
		SELECT
			transaction_date as date,
			COUNT(*) as transaction_count,
			COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) as total_amount,
			COUNT(CASE WHEN transaction_type = 'withdrawal' THEN 1 END) as success_count,
			COALESCE(SUM(CASE WHEN transaction_type = 'withdrawal' AND amount > 0 THEN amount ELSE 0 END), 0) as success_amount,
			COUNT(CASE WHEN transaction_type = 'refund' THEN 1 END) as failed_count,
			COALESCE(SUM(CASE WHEN transaction_type = 'refund' AND amount > 0 THEN amount ELSE 0 END), 0) as failed_amount
		FROM withdrawal_records
		WHERE tenant_id = ? AND transaction_date >= ? AND transaction_date <= ?
		GROUP BY transaction_date
		ORDER BY transaction_date DESC
	`

	err := r.db.WithContext(ctx).Raw(query, tenantID, startDate, endDate).Scan(&stats).Error
	return stats, err
}

// GetBalanceHistory 获取余额历史
func (r *withdrawalRecordRepositoryImpl) GetBalanceHistory(
	ctx context.Context,
	tenantID string,
	startDate, endDate string,
) ([]*BalanceHistory, error) {
	var history []*BalanceHistory

	query := `
		SELECT
			transaction_date as date,
			MIN(CASE WHEN rn = 1 THEN balance_before END) as opening_balance,
			MAX(balance_after) as closing_balance,
			COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) as total_inflow,
			COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0) as total_outflow
		FROM (
			SELECT *,
				ROW_NUMBER() OVER (PARTITION BY transaction_date ORDER BY created_at ASC) as rn
			FROM withdrawal_records
			WHERE tenant_id = ? AND transaction_date >= ? AND transaction_date <= ?
		) AS sub
		GROUP BY transaction_date
		ORDER BY transaction_date DESC
	`

	err := r.db.WithContext(ctx).Raw(query, tenantID, startDate, endDate).Scan(&history).Error
	return history, err
}

// ============================================================

// withdrawalReconciliationRepositoryImpl 提现对账仓储实现
type withdrawalReconciliationRepositoryImpl struct {
	db *gorm.DB
}

// NewWithdrawalReconciliationRepository 创建提现对账仓储实例
func NewWithdrawalReconciliationRepository(db *gorm.DB) WithdrawalReconciliationRepository {
	return &withdrawalReconciliationRepositoryImpl{db: db}
}

// Create 创建对账记录
func (r *withdrawalReconciliationRepositoryImpl) Create(ctx context.Context, reconciliation *entity.WithdrawalReconciliation) error {
	return r.db.WithContext(ctx).Create(reconciliation).Error
}

// GetByID 根据ID获取对账记录
func (r *withdrawalReconciliationRepositoryImpl) GetByID(ctx context.Context, id uint64) (*entity.WithdrawalReconciliation, error) {
	var reconciliation entity.WithdrawalReconciliation
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&reconciliation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &reconciliation, nil
}

// GetByBatch 根据批次号获取对账记录
func (r *withdrawalReconciliationRepositoryImpl) GetByBatch(ctx context.Context, batch string) (*entity.WithdrawalReconciliation, error) {
	var reconciliation entity.WithdrawalReconciliation
	err := r.db.WithContext(ctx).
		Where("reconciliation_batch = ? AND deleted_at IS NULL", batch).
		First(&reconciliation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &reconciliation, nil
}

// Update 更新对账记录
func (r *withdrawalReconciliationRepositoryImpl) Update(ctx context.Context, reconciliation *entity.WithdrawalReconciliation) error {
	return r.db.WithContext(ctx).Save(reconciliation).Error
}

// List 分页查询对账记录
func (r *withdrawalReconciliationRepositoryImpl) List(ctx context.Context, filter *WithdrawalReconciliationFilter) ([]*entity.WithdrawalReconciliation, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.WithdrawalReconciliation{})

	// 日期范围过滤
	if filter.StartDate != nil {
		query = query.Where("reconciliation_date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("reconciliation_date <= ?", *filter.EndDate)
	}

	// 状态过滤
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	// 批次号过滤
	if filter.ReconciliationBatch != nil {
		query = query.Where("reconciliation_batch LIKE ?", "%"+*filter.ReconciliationBatch+"%")
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var reconciliations []*entity.WithdrawalReconciliation
	offset := (filter.PageToken - 1) * filter.PageSize
	err := query.Order("reconciliation_date DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&reconciliations).Error

	return reconciliations, total, err
}

// GetLatest 获取最新的对账记录
func (r *withdrawalReconciliationRepositoryImpl) GetLatest(ctx context.Context, limit int) ([]*entity.WithdrawalReconciliation, error) {
	var reconciliations []*entity.WithdrawalReconciliation
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("reconciliation_date DESC").
		Limit(limit).
		Find(&reconciliations).Error
	return reconciliations, err
}
