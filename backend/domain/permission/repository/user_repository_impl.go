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

	"gorm.io/gorm"

	userentity "github.com/coze-dev/coze-studio/backend/domain/user/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

// Create 创建用户
func (r *UserRepositoryImpl) Create(ctx context.Context, user *userentity.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return errorx.WrapByCode(err, errno.DatabaseErrorCode,
			errorx.KV("operation", "create user"),
		)
	}
	return nil
}

// GetByID 根据ID获取用户
func (r *UserRepositoryImpl) GetByID(ctx context.Context, userID int64) (*userentity.User, error) {
	var user userentity.User
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.DatabaseErrorCode,
			errorx.KV("operation", "get user by id"),
		)
	}

	return &user, nil
}

// GetByEmail 根据租户ID和邮箱获取用户
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, tenantID, email string) (*userentity.User, error) {
	var user userentity.User
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND email = ? AND deleted_at IS NULL", tenantID, email).
		First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.DatabaseErrorCode,
			errorx.KV("operation", "get user by email"),
		)
	}

	return &user, nil
}

// GetByUniqueName 根据租户ID和唯一名称获取用户
func (r *UserRepositoryImpl) GetByUniqueName(ctx context.Context, tenantID, uniqueName string) (*userentity.User, error) {
	var user userentity.User
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND unique_name = ? AND deleted_at IS NULL", tenantID, uniqueName).
		First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.DatabaseErrorCode,
			errorx.KV("operation", "get user by unique_name"),
		)
	}

	return &user, nil
}

// Update 更新用户
func (r *UserRepositoryImpl) Update(ctx context.Context, user *userentity.User) error {
	err := r.db.WithContext(ctx).
		Model(&userentity.User{}).
		Where("user_id = ? AND tenant_id = ? AND deleted_at IS NULL", user.UserID, user.TenantID).
		Updates(user).Error

	if err != nil {
		return errorx.WrapByCode(err, errno.DatabaseErrorCode,
			errorx.KV("operation", "update user"),
		)
	}

	return nil
}

// Delete 软删除用户
func (r *UserRepositoryImpl) Delete(ctx context.Context, userID int64) error {
	now := int64(0) // TODO: 设置为当前时间戳
	err := r.db.WithContext(ctx).
		Model(&userentity.User{}).
		Where("user_id = ?", userID).
		Update("deleted_at", now).Error

	if err != nil {
		return errorx.WrapByCode(err, errno.DatabaseErrorCode,
			errorx.KV("operation", "delete user"),
		)
	}

	return nil
}

// List 分页查询用户列表
func (r *UserRepositoryImpl) List(ctx context.Context, filter *UserFilter) ([]*userentity.User, int64, error) {
	query := r.db.WithContext(ctx).Model(&userentity.User{})

	// 租户隔离
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	// 未删除
	query = query.Where("deleted_at IS NULL")

	// 搜索关键词（姓名或邮箱）
	if filter.Keyword != "" {
		query = query.Where("name LIKE ? OR email LIKE ? OR unique_name LIKE ?",
			"%"+filter.Keyword+"%",
			"%"+filter.Keyword+"%",
			"%"+filter.Keyword+"%")
	}

	// 按状态筛选
	if filter.Status == "active" {
		query = query.Where("is_enabled = ?", true)
	} else if filter.Status == "inactive" {
		query = query.Where("is_enabled = ?", false)
	}

	// 按角色筛选
	if filter.RoleID != "" {
		// 关联查询user_roles表
		query = query.Joins("INNER JOIN user_roles ON user_roles.user_id = CAST(users.user_id AS CHAR)").
			Where("user_roles.role_id = ?", filter.RoleID)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errorx.WrapByCode(err, errno.DatabaseErrorCode,
			errorx.KV("operation", "count users"),
		)
	}

	// 分页
	if filter.PageSize > 0 {
		offset := 0
		if filter.PageToken != "" {
			// TODO: 解析page_token获取offset
			offset = 0
		}
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 按创建时间倒序
	query = query.Order("created_at DESC")

	// 查询
	var users []*userentity.User
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, errorx.WrapByCode(err, errno.DatabaseErrorCode,
			errorx.KV("operation", "list users"),
		)
	}

	return users, total, nil
}

// GetByIDs 根据ID列表批量获取用户
func (r *UserRepositoryImpl) GetByIDs(ctx context.Context, userIDs []int64) ([]*userentity.User, error) {
	if len(userIDs) == 0 {
		return []*userentity.User{}, nil
	}

	var users []*userentity.User
	err := r.db.WithContext(ctx).
		Where("user_id IN ? AND deleted_at IS NULL", userIDs).
		Find(&users).Error

	if err != nil {
		return nil, errorx.WrapByCode(err, errno.DatabaseErrorCode,
			errorx.KV("operation", "get users by ids"),
		)
	}

	return users, nil
}
