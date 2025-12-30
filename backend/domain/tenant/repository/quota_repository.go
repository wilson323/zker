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

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
)

// QuotaRepository 配额仓储接口
type QuotaRepository interface {
	// Create 创建配额
	Create(ctx context.Context, quota *entity.Quota) error

	// GetByID 根据ID获取配额
	GetByID(ctx context.Context, quotaID string) (*entity.Quota, error)

	// GetByTenantAndResource 根据租户ID和资源类型获取配额
	GetByTenantAndResource(ctx context.Context, tenantID string, resourceType entity.ResourceType) (*entity.Quota, error)

	// GetByTenant 根据租户ID获取所有配额
	GetByTenant(ctx context.Context, tenantID string) ([]*entity.Quota, error)

	// List 获取租户的所有配额
	List(ctx context.Context, tenantID string) ([]*entity.Quota, error)

	// GetAll 获取所有配额
	GetAll(ctx context.Context) ([]*entity.Quota, error)

	// Update 更新配额
	Update(ctx context.Context, quota *entity.Quota) error

	// UpdateUsedCount 更新使用计数（count可为正数或负数）
	UpdateUsedCount(ctx context.Context, quotaID string, count int) error

	// ResetUsage 重置使用计数
	ResetUsage(ctx context.Context, quotaID string) error

	// Delete 删除配额
	Delete(ctx context.Context, quotaID string) error

	// BatchCreate 批量创建配额
	BatchCreate(ctx context.Context, quotas []*entity.Quota) error
}
