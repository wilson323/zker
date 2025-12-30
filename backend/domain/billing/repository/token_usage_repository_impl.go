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
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
)

// ============================================================
// TokenUsageRepository GORM实现
// ============================================================

// tokenUsageRepository Token使用仓储实现
type tokenUsageRepository struct {
	db *gorm.DB
}

// NewTokenUsageRepository 创建Token使用仓储实例
func NewTokenUsageRepository(db *gorm.DB) TokenUsageRepository {
	return &tokenUsageRepository{db: db}
}

// Create 创建Token使用记录
func (r *tokenUsageRepository) Create(ctx context.Context, usage *entity.TokenUsage) error {
	return r.db.WithContext(ctx).Create(usage).Error
}

// GetByID 根据ID获取Token使用记录
func (r *tokenUsageRepository) GetByID(ctx context.Context, id string) (*entity.TokenUsage, error) {
	var usage entity.TokenUsage
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&usage).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &usage, nil
}

// ListByTenant 分页查询租户的Token使用记录
func (r *tokenUsageRepository) ListByTenant(
	ctx context.Context,
	tenantID string,
	limit, offset int,
) ([]*entity.TokenUsage, error) {
	var usages []*entity.TokenUsage
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&usages).Error

	return usages, err
}

// GetMonthlyUsage 获取租户指定模型的月度使用量（Token数）
func (r *tokenUsageRepository) GetMonthlyUsage(
	ctx context.Context,
	tenantID, modelID string,
) (int, error) {
	var totalTokens int64

	// 获取当前月份的第一天和最后一天
	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastDay := firstDay.AddDate(0, 1, -1)

	err := r.db.WithContext(ctx).
		Model(&entity.TokenUsage{}).
		Where("tenant_id = ?", tenantID).
		Where("model_id = ?", modelID).
		Where("created_at >= ? AND created_at <= ?", firstDay, lastDay).
		Select("COALESCE(SUM(total_tokens), 0)").
		Scan(&totalTokens).Error

	return int(totalTokens), err
}

// GetUsageByModel 获取租户各模型的使用统计
func (r *tokenUsageRepository) GetUsageByModel(
	ctx context.Context,
	tenantID string,
) ([]*ModelUsage, error) {
	type UsageResult struct {
		ModelID      string
		TotalTokens  int64
		CostUSD      float64
		RequestCount int64
	}

	var results []UsageResult
	err := r.db.WithContext(ctx).
		Model(&entity.TokenUsage{}).
		Select("model_id, COALESCE(SUM(total_tokens), 0) as total_tokens, COALESCE(SUM(cost_usd), 0) as cost_usd, COUNT(*) as request_count").
		Where("tenant_id = ?", tenantID).
		Group("model_id").
		Order("total_tokens DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 转换为返回类型
	usages := make([]*ModelUsage, 0, len(results))
	for _, r := range results {
		usages = append(usages, &ModelUsage{
			ModelID:      r.ModelID,
			TotalTokens:  int(r.TotalTokens),
			CostUSD:      r.CostUSD,
			RequestCount: int(r.RequestCount),
		})
	}

	return usages, nil
}

// ============================================================
// TokenBudgetRepository GORM实现
// ============================================================

// tokenBudgetRepository Token预算仓储实现
type tokenBudgetRepository struct {
	db *gorm.DB
}

// NewTokenBudgetRepository 创建Token预算仓储实例
func NewTokenBudgetRepository(db *gorm.DB) TokenBudgetRepository {
	return &tokenBudgetRepository{db: db}
}

// Create 创建Token预算配置
func (r *tokenBudgetRepository) Create(ctx context.Context, budget *entity.TokenBudget) error {
	return r.db.WithContext(ctx).Create(budget).Error
}

// GetByTenantID 获取租户的预算配置
func (r *tokenBudgetRepository) GetByTenantID(
	ctx context.Context,
	tenantID string,
) (*entity.TokenBudget, error) {
	var budget entity.TokenBudget
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		First(&budget).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

// Update 更新Token预算配置
func (r *tokenBudgetRepository) Update(ctx context.Context, budget *entity.TokenBudget) error {
	return r.db.WithContext(ctx).
		Model(&entity.TokenBudget{}).
		Where("tenant_id = ?", budget.TenantID).
		Updates(budget).Error
}

// GetByTenantAndModel 获取租户指定模型的预算配置
func (r *tokenBudgetRepository) GetByTenantAndModel(
	ctx context.Context,
	tenantID, modelID, period string,
) (*entity.TokenBudget, error) {
	var budget entity.TokenBudget

	// 注意：由于entity.TokenBudget没有model_id字段，
	// 这里我们获取租户的整体预算，实际应用中可能需要修改entity结构以支持多模型预算
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		First(&budget).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// 如果指定了reset_period且不匹配，返回nil
	if period != "" && budget.ResetPeriod != period {
		return nil, nil
	}

	return &budget, nil
}
