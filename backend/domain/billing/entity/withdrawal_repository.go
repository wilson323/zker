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
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
)

// ============================================================
// 提现管理Repository接口
// ============================================================

// WithdrawalRequestRepository 提现申请仓储接口
type WithdrawalRequestRepository interface {
	// Create 创建提现申请
	Create(ctx context.Context, withdrawal *entity.WithdrawalRequest) error

	// GetByID 根据ID获取提现申请
	GetByID(ctx context.Context, id uint64) (*entity.WithdrawalRequest, error)

	// GetByWithdrawalNumber 根据提现单号获取申请
	GetByWithdrawalNumber(ctx context.Context, withdrawalNumber string) (*entity.WithdrawalRequest, error)

	// Update 更新提现申请
	Update(ctx context.Context, withdrawal *entity.WithdrawalRequest) error

	// Delete 删除提现申请（软删除）
	Delete(ctx context.Context, id uint64) error

	// List 分页查询提现申请
	List(ctx context.Context, filter *WithdrawalRequestFilter) ([]*entity.WithdrawalRequest, int64, error)

	// GetByTenant 获取租户的提现申请
	GetByTenant(
		ctx context.Context,
		tenantID string,
		page, pageSize int,
	) ([]*entity.WithdrawalRequest, int64, error)

	// GetPendingWithdrawals 获取待处理的提现申请
	GetPendingWithdrawals(ctx context.Context, limit int) ([]*entity.WithdrawalRequest, error)

	// CountByTenantAndDate 统计租户在指定日期的提现数量
	CountByTenantAndDate(ctx context.Context, tenantID string, date string) (int, error)

	// SumAmountByTenantAndDate 统计租户在指定日期的提现金额
	SumAmountByTenantAndDate(ctx context.Context, tenantID string, date string) (float64, error)

	// GetStatisticsByTenant 获取租户提现统计
	GetStatisticsByTenant(
		ctx context.Context,
		tenantID string,
		startDate, endDate string,
	) (*WithdrawalStatistics, error)

	// BatchUpdateStatus 批量更新状态
	BatchUpdateStatus(ctx context.Context, ids []uint64, status string) error
}

// WithdrawalRequestFilter 提现申请查询过滤器
type WithdrawalRequestFilter struct {
	TenantID         string
	Status           *string
	WithdrawalMethod *string
	StartDate        *time.Time
	EndDate          *time.Time
	MinAmount        *float64
	MaxAmount        *float64
	ReviewerID       *uint64
	WithdrawalNumber *string
	PageToken        string
	PageSize         int
}

// WithdrawalStatistics 提现统计
type WithdrawalStatistics struct {
	TotalCount   int     `json:"total_count"`
	TotalAmount  float64 `json:"total_amount"`
	SuccessCount int     `json:"success_count"`
	SuccessAmount float64 `json:"success_amount"`
	PendingCount int     `json:"pending_count"`
	PendingAmount float64 `json:"pending_amount"`
	FailedCount  int     `json:"failed_count"`
	FailedAmount float64 `json:"failed_amount"`
	TotalFee     float64 `json:"total_fee"`
}

// ============================================================

// WithdrawalLimitRepository 提现限额仓储接口
type WithdrawalLimitRepository interface {
	// Create 创建限额配置
	Create(ctx context.Context, limit *entity.WithdrawalLimit) error

	// GetByID 根据ID获取限额配置
	GetByID(ctx context.Context, id uint64) (*entity.WithdrawalLimit, error)

	// Update 更新限额配置
	Update(ctx context.Context, limit *entity.WithdrawalLimit) error

	// Delete 删除限额配置（软删除）
	Delete(ctx context.Context, id uint64) error

	// List 列出所有启用的限额配置
	List(ctx context.Context, enabledOnly bool) ([]*entity.WithdrawalLimit, error)

	// GetByMethod 获取指定提现方式的限额配置
	GetByMethod(ctx context.Context, withdrawalMethod string) ([]*entity.WithdrawalLimit, error)

	// GetEffectiveLimit 获取有效的限额配置
	GetEffectiveLimit(ctx context.Context, withdrawalMethod string, amount float64) (*entity.WithdrawalLimit, error)
}

// ============================================================

// WithdrawalRecordRepository 提现流水仓储接口
type WithdrawalRecordRepository interface {
	// Create 创建流水记录
	Create(ctx context.Context, record *entity.WithdrawalRecord) error

	// GetByID 根据ID获取流水记录
	GetByID(ctx context.Context, id uint64) (*entity.WithdrawalRecord, error)

	// GetByWithdrawalID 获取提现的所有流水记录
	GetByWithdrawalID(ctx context.Context, withdrawalID uint64) ([]*entity.WithdrawalRecord, error)

	// List 分页查询流水记录
	List(ctx context.Context, filter *WithdrawalRecordFilter) ([]*entity.WithdrawalRecord, int64, error)

	// BatchCreate 批量创建流水记录
	BatchCreate(ctx context.Context, records []*entity.WithdrawalRecord) error

	// GetDailyStatistics 获取每日统计
	GetDailyStatistics(
		ctx context.Context,
		tenantID string,
		startDate, endDate string,
	) ([]*DailyWithdrawalStatistics, error)

	// GetBalanceHistory 获取余额历史
	GetBalanceHistory(
		ctx context.Context,
		tenantID string,
		startDate, endDate string,
	) ([]*BalanceHistory, error)
}

// WithdrawalRecordFilter 流水记录查询过滤器
type WithdrawalRecordFilter struct {
	TenantID        string
	WithdrawalID    *uint64
	TransactionType *string
	StartDate       *time.Time
	EndDate         *time.Time
	PageToken       string
	PageSize        int
}

// DailyWithdrawalStatistics 每日提现统计
type DailyWithdrawalStatistics struct {
	Date          string  `json:"date"`
	TransactionCount int   `json:"transaction_count"`
	TotalAmount   float64 `json:"total_amount"`
	SuccessCount  int     `json:"success_count"`
	SuccessAmount float64 `json:"success_amount"`
	FailedCount   int     `json:"failed_count"`
	FailedAmount  float64 `json:"failed_amount"`
}

// BalanceHistory 余额历史
type BalanceHistory struct {
	Date         string  `json:"date"`
	OpeningBalance float64 `json:"opening_balance"`
	ClosingBalance float64 `json:"closing_balance"`
	TotalInflow  float64 `json:"total_inflow"`
	TotalOutflow float64 `json:"total_outflow"`
}

// ============================================================

// WithdrawalReconciliationRepository 提现对账仓储接口
type WithdrawalReconciliationRepository interface {
	// Create 创建对账记录
	Create(ctx context.Context, reconciliation *entity.WithdrawalReconciliation) error

	// GetByID 根据ID获取对账记录
	GetByID(ctx context.Context, id uint64) (*entity.WithdrawalReconciliation, error)

	// GetByBatch 根据批次号获取对账记录
	GetByBatch(ctx context.Context, batch string) (*entity.WithdrawalReconciliation, error)

	// Update 更新对账记录
	Update(ctx context.Context, reconciliation *entity.WithdrawalReconciliation) error

	// List 分页查询对账记录
	List(ctx context.Context, filter *WithdrawalReconciliationFilter) ([]*entity.WithdrawalReconciliation, int64, error)

	// GetLatest 获取最新的对账记录
	GetLatest(ctx context.Context, limit int) ([]*entity.WithdrawalReconciliation, error)
}

// WithdrawalReconciliationFilter 对账记录查询过滤器
type WithdrawalReconciliationFilter struct {
	StartDate           *time.Time
	EndDate             *time.Time
	Status              *string
	ReconciliationBatch *string
	PageToken           string
	PageSize            int
}
