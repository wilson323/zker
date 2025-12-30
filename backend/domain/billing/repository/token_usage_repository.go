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

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
)

// TokenUsageRepository Token使用仓储接口
type TokenUsageRepository interface {
	Create(ctx context.Context, usage *entity.TokenUsage) error
	GetByID(ctx context.Context, id string) (*entity.TokenUsage, error)
	ListByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*entity.TokenUsage, error)

	// GetMonthlyUsage 获取租户指定模型的月度使用量（Token数）
	GetMonthlyUsage(ctx context.Context, tenantID, modelID string) (int, error)

	// GetUsageByModel 获取租户各模型的使用统计
	GetUsageByModel(ctx context.Context, tenantID string) ([]*ModelUsage, error)
}

// TokenBudgetRepository Token预算仓储接口
type TokenBudgetRepository interface {
	Create(ctx context.Context, budget *entity.TokenBudget) error
	GetByTenantID(ctx context.Context, tenantID string) (*entity.TokenBudget, error)
	Update(ctx context.Context, budget *entity.TokenBudget) error

	// GetByTenantAndModel 获取租户指定模型的预算配置
	GetByTenantAndModel(ctx context.Context, tenantID, modelID, period string) (*entity.TokenBudget, error)
}

// ModelUsage 模型使用统计（用于GetUsageByModel返回）
type ModelUsage struct {
	ModelID      string
	TotalTokens  int
	CostUSD      float64
	RequestCount int
}
