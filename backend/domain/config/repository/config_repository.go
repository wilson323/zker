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

	"github.com/coze-dev/coze-studio/backend/domain/config/entity"
)

// ConfigRepository 配置仓储接口
type ConfigRepository interface {
	// Create 创建配置项
	Create(ctx context.Context, config *entity.ConfigItem) error

	// GetByID 根据ID获取配置
	GetByID(ctx context.Context, configID string) (*entity.ConfigItem, error)

	// GetByKey 根据租户和key获取配置
	GetByKey(ctx context.Context, tenantID, key string) (*entity.ConfigItem, error)

	// ListByTenant 获取租户所有配置
	ListByTenant(ctx context.Context, tenantID string) ([]*entity.ConfigItem, error)

	// Update 更新配置
	Update(ctx context.Context, config *entity.ConfigItem) error

	// Delete 删除配置
	Delete(ctx context.Context, configID string) error

	// GetByPrefix 根据前缀获取配置列表
	GetByPrefix(ctx context.Context, tenantID, prefix string) ([]*entity.ConfigItem, error)
}

// ConfigHistoryRepository 配置历史仓储接口
type ConfigHistoryRepository interface {
	// Create 创建历史记录
	Create(ctx context.Context, history *entity.ConfigHistory) error

	// ListByConfigID 获取配置的所有历史记录
	ListByConfigID(ctx context.Context, configID string, limit, offset int) ([]*entity.ConfigHistory, error)

	// GetLatestVersion 获取最新版本号
	GetLatestVersion(ctx context.Context, configID string) (int, error)

	// ListByTenant 获取租户的所有配置变更历史
	ListByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*entity.ConfigHistory, error)
}
