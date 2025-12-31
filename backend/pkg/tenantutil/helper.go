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

package tenantutil

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/api/middleware"
	tenantentity "github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
)

var (
	// ErrTenantNotFound 租户未找到错误
	ErrTenantNotFound = errors.New("tenant_id not found in context")
	// ErrTenantAccessDenied 租户访问被拒绝
	ErrTenantAccessDenied = errors.New("tenant access denied")
)

// WithTenantFilter 为GORM查询添加租户过滤
//
// 该helper函数确保所有查询都包含tenant_id过滤，防止数据泄露
//
// 参数:
//   - db: GORM DB实例
//   - ctx: 上下文（必须包含tenant_id）
//
// 返回:
//   - *gorm.DB: 添加了tenant_id过滤的查询构建器
//
// 示例:
//
//	query := tenantutil.WithTenantFilter(db, ctx)
//	query.Find(&bots)
func WithTenantFilter(db *gorm.DB, ctx context.Context) *gorm.DB {
	tenantID, err := GetTenantID(ctx)
	if err != nil {
		// 如果无法获取tenant_id，返回一个会产生空结果的查询
		// 这是为了安全：宁可不查询，也不能泄露数据
		return db.Session(&gorm.Session{}).Where("1 = 0")
	}

	return db.Where("tenant_id = ?", tenantID)
}

// WithTenantFilterAndDeleted 为GORM查询添加租户过滤和软删除过滤
//
// 该helper函数确保查询包含：
// 1. tenant_id过滤（防止跨租户数据泄露）
// 2. deleted_at IS NULL过滤（只查询未删除数据）
//
// 参数:
//   - db: GORM DB实例
//   - ctx: 上下文（必须包含tenant_id）
//
// 返回:
//   - *gorm.DB: 添加了过滤条件的查询构建器
//
// 示例:
//
//	query := tenantutil.WithTenantFilterAndDeleted(db, ctx)
//	query.Find(&bots)
func WithTenantFilterAndDeleted(db *gorm.DB, ctx context.Context) *gorm.DB {
	return WithTenantFilter(db, ctx).Where("deleted_at IS NULL")
}

// GetTenantID 从上下文中获取租户ID
//
// 该函数从上下文中提取tenant_id，支持多种上下文实现
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - string: 租户ID
//   - error: 如果无法获取租户ID则返回错误
func GetTenantID(ctx context.Context) (string, error) {
	// 尝试从标准context中获取（中间件设置）
	if tenantID, ok := ctx.Value(middleware.TenantContextKey).(*tenantentity.Tenant); ok && tenantID != nil {
		return tenantID.TenantID, nil
	}

	// 尝试从context直接获取（兼容性）
	if tenantID, ok := ctx.Value("tenant_id").(string); ok && tenantID != "" {
		return tenantID, nil
	}

	return "", ErrTenantNotFound
}

// MustGetTenantID 从上下文中获取租户ID（panic如果失败）
//
// 该函数用于确保tenant_id必须存在的场景
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - string: 租户ID
//
// Panic:
//   - 如果无法获取租户ID则panic
func MustGetTenantID(ctx context.Context) string {
	tenantID, err := GetTenantID(ctx)
	if err != nil {
		panic(err)
	}
	return tenantID
}

// WithTenantFilterByID 使用指定租户ID创建过滤查询
//
// 该函数用于已经拥有tenantID的场景，避免重复查询context
//
// 参数:
//   - db: GORM DB实例
//   - tenantID: 租户ID
//
// 返回:
//   - *gorm.DB: 添加了租户过滤的查询构建器
//
// 示例:
//
//	tenantID := "tenant-123"
//	query := tenantutil.WithTenantFilterByID(db, tenantID)
//	query.Find(&bots)
func WithTenantFilterByID(db *gorm.DB, tenantID string) *gorm.DB {
	if tenantID == "" {
		// 返回空结果查询
		return db.Session(&gorm.Session{}).Where("1 = 0")
	}

	return db.Where("tenant_id = ?", tenantID)
}

// WithTenantFilterByIDAndDeleted 使用指定租户ID创建过滤查询（包含软删除过滤）
//
// 参数:
//   - db: GORM DB实例
//   - tenantID: 租户ID
//
// 返回:
//   - *gorm.DB: 添加了过滤条件的查询构建器
func WithTenantFilterByIDAndDeleted(db *gorm.DB, tenantID string) *gorm.DB {
	return WithTenantFilterByID(db, tenantID).Where("deleted_at IS NULL")
}

// CheckTenantAccess 检查资源是否属于当前租户
//
// 该函数用于验证用户是否有权访问指定资源
//
// 参数:
//   - ctx: 上下文
//   - resourceTenantID: 资源的租户ID
//
// 返回:
//   - error: 如果租户不匹配则返回权限错误
func CheckTenantAccess(ctx context.Context, resourceTenantID string) error {
	currentTenantID, err := GetTenantID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tenant_id from context: %w", err)
	}

	if currentTenantID != resourceTenantID {
		return fmt.Errorf("%w: current tenant %s cannot access resource of tenant %s",
			ErrTenantAccessDenied, currentTenantID, resourceTenantID)
	}

	return nil
}

// IsTenantOwned 检查资源是否属于指定租户
//
// 参数:
//   - ctx: 上下文
//   - resourceTenantID: 资源的租户ID
//
// 返回:
//   - bool: 如果资源属于当前租户则返回true
func IsTenantOwned(ctx context.Context, resourceTenantID string) bool {
	return CheckTenantAccess(ctx, resourceTenantID) == nil
}
