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

package service

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"gorm.io/gorm"
)

// TenantService 租户服务接口
// 提供租户基本查询功能，供中间件使用
type TenantService interface {
	// GetTenant 获取租户信息
	GetTenant(ctx context.Context, tenantID string) (*entity.Tenant, error)
}

// tenantServiceImpl 租户服务实现
type tenantServiceImpl struct {
	db *gorm.DB
}

// NewTenantService 创建租户服务实例
func NewTenantService(db *gorm.DB) TenantService {
	return &tenantServiceImpl{db: db}
}

// GetTenant 获取租户信息
func (s *tenantServiceImpl) GetTenant(ctx context.Context, tenantID string) (*entity.Tenant, error) {
	var tenant entity.Tenant
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		First(&tenant).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &tenant, nil
}

// GetDB 获取数据库连接（供其他服务使用）
func (s *tenantServiceImpl) GetDB() *gorm.DB {
	return s.db
}
