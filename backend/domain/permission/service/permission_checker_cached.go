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
	"fmt"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
	"github.com/coze-dev/coze-studio/backend/infra/cache"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// CachedPermissionChecker 带缓存的权限检查服务
type CachedPermissionChecker struct {
	*PermissionChecker                    // 嵌入原始权限检查器
	permCache          cache.PermissionCache // 权限缓存
}

// NewCachedPermissionChecker 创建带缓存的权限检查服务实例
func NewCachedPermissionChecker(
	db *gorm.DB,
	roleRepo repository.RoleRepository,
	dataPermRepo repository.DataPermissionRepository,
	fieldPermRepo repository.FieldPermissionRepository,
	userRoleRepo repository.UserRoleRepository,
	userDeptRepo repository.UserDepartmentRepository,
	departmentRepo repository.DepartmentRepository,
	permCache cache.PermissionCache,
) *CachedPermissionChecker {
	baseChecker := &PermissionChecker{
		db:             db,
		roleRepo:       roleRepo,
		dataPermRepo:   dataPermRepo,
		fieldPermRepo:  fieldPermRepo,
		userRoleRepo:   userRoleRepo,
		userDeptRepo:   userDeptRepo,
		departmentRepo: departmentRepo,
	}

	return &CachedPermissionChecker{
		PermissionChecker: baseChecker,
		permCache:         permCache,
	}
}

// CheckDataPermission 检查数据权限（带缓存）
// 性能优化：
//   - L1: 本地缓存（1分钟）- 响应时间 < 1ms
//   - L2: Redis缓存（5分钟）- 响应时间 < 5ms
//   - L3: 数据库查询 - 响应时间 20ms
// 预期缓存命中率：95%+
func (p *CachedPermissionChecker) CheckDataPermission(
	ctx context.Context,
	tenantID, userID string,
	resourceType entity.ResourceType,
	action string,
	resourceID string,
) (bool, error) {
	// 1. 参数验证
	if tenantID == "" {
		return false, errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("field", "tenant_id"),
                errorx.KV("reason", "required field is missing"),
            )
	}
	if userID == "" {
		return false, errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("field", "user_id"),
                errorx.KV("reason", "required field is missing"),
            )
	}
	if resourceType == "" {
		return false, errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("field", "resource_type"),
                errorx.KV("reason", "required field is missing"),
            )
	}

	// 2. 构建缓存键
	// 格式: {tenant_id}:{user_id}:{resource_type}:{action}:{resource_id}
	cacheKey := fmt.Sprintf("%s:%s:%s:%s:%s", tenantID, userID, resourceType, action, resourceID)

	// 3. 使用缓存检查权限
	allowed, err := p.permCache.CheckDataPermission(ctx, cacheKey, func() (bool, error) {
		// 缓存未命中时调用原始的数据库查询逻辑
		return p.PermissionChecker.CheckDataPermission(ctx, tenantID, userID, resourceType, action, resourceID)
	})

	if err != nil {
		return false, err
	}

	return allowed, nil
}

// GetFieldPermissions 获取字段权限（带缓存）
func (p *CachedPermissionChecker) GetFieldPermissions(
	ctx context.Context,
	tenantID, userID string,
	resourceType string,
) (map[string]string, error) {
	// 1. 参数验证
	if tenantID == "" || userID == "" || resourceType == "" {
		return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "tenant_id, user_id and resource_type are required"),
            )
	}

	// 2. 构建缓存键
	cacheKey := fmt.Sprintf("%s:%s:%s", tenantID, userID, resourceType)

	// 3. 使用缓存获取字段权限
	fieldPerms, err := p.permCache.GetFieldPermissions(ctx, cacheKey, func() (map[string]string, error) {
		// 缓存未命中时调用原始的数据库查询逻辑
		return p.PermissionChecker.GetFieldPermissions(ctx, tenantID, userID, resourceType)
	})

	if err != nil {
		return nil, err
	}

	return fieldPerms, nil
}

// InvalidateUserPermissions 使用户权限缓存失效
// 调用时机：
//   1. 用户角色变更时
//   2. 角色权限配置变更时
//   3. 数据权限配置变更时
func (p *CachedPermissionChecker) InvalidateUserPermissions(ctx context.Context, userID string) error {
	return p.permCache.InvalidateUserPermissions(ctx, userID)
}

// InvalidateDataPermission 使特定数据权限缓存失效
// 调用时机：
//   1. 资源所有者变更时
//   2. 资源部门变更时
func (p *CachedPermissionChecker) InvalidateDataPermission(
	ctx context.Context,
	tenantID, userID string,
	resourceType entity.ResourceType,
	action string,
	resourceID string,
) error {
	cacheKey := fmt.Sprintf("%s:%s:%s:%s:%s", tenantID, userID, resourceType, action, resourceID)
	return p.permCache.InvalidateDataPermission(ctx, cacheKey)
}

// GetCacheStats 获取缓存统计信息
func (p *CachedPermissionChecker) GetCacheStats() cache.PermissionCacheStats {
	return p.permCache.GetCacheStats()
}
