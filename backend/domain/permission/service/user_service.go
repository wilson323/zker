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
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
	"github.com/coze-dev/coze-studio/backend/domain/user/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// UserService 用户管理服务
type UserService struct {
	userRepo    repository.UserRepository
	userRoleRepo repository.UserRoleRepository
	roleRepo     repository.RoleRepository
	permChecker  *PermissionChecker
}

// NewUserService 创建用户管理服务实例
func NewUserService(
	userRepo repository.UserRepository,
	userRoleRepo repository.UserRoleRepository,
	roleRepo repository.RoleRepository,
	permChecker *PermissionChecker,
) *UserService {
	return &UserService{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
		permChecker:  permChecker,
	}
}

// ========== 请求/响应结构 ==========

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	TenantID    string
	Name        string
	Email       string
	UniqueName  string
	Description string
	IconURI     string
	Locale      string
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	UserID      int64
	Name        *string
	Description *string
	IconURI     *string
	Locale      *string
}

// UserListFilter 用户列表过滤器
type UserListFilter struct {
	TenantID  string
	Keyword   string
	RoleID    string
	Status    string
	PageToken string
	PageSize  int
}

// UserPermission 用户权限
type UserPermission struct {
	ResourceType entity.ResourceType `json:"resource_type"`
	Actions      []string            `json:"actions"`
	Scope        entity.DataPermissionScope `json:"scope"`
}

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	SuccessCount int      `json:"success_count"`
	FailureCount int      `json:"failure_count"`
	Errors       []string `json:"errors,omitempty"`
}

// ========== 1. 用户基础CRUD ==========

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*entity.User, error) {
	// 1. 检查邮箱是否已存在
	existing, err := s.userRepo.GetByEmail(ctx, req.TenantID, req.Email)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("reason", "failed to check email"),
		)
	}
	if existing != nil {
		return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("field", "email"),
			errorx.KV("value", req.Email),
			errorx.KV("reason", "already exists"),
		)
	}

	// 2. 检查UniqueName是否已存在
	existing, err = s.userRepo.GetByUniqueName(ctx, req.TenantID, req.UniqueName)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("reason", "failed to check unique_name"),
		)
	}
	if existing != nil {
		return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("field", "unique_name"),
			errorx.KV("value", req.UniqueName),
			errorx.KV("reason", "already exists"),
		)
	}

	// 3. 创建用户
	user := &entity.User{
		UserID:       generateUserID(),
		TenantID:     req.TenantID,
		Name:         req.Name,
		Email:        req.Email,
		UniqueName:   req.UniqueName,
		Description:  req.Description,
		IconURI:      req.IconURI,
		IconURL:      req.IconURI,
		Locale:       req.Locale,
		UserVerified: false,
		IsEnabled:    true,
		CreatedAt:    time.Now().UnixMilli(),
		UpdatedAt:    time.Now().UnixMilli(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("reason", "failed to create user"),
		)
	}

	// 4. 分配默认角色（tenant_member）
	if err := s.assignDefaultRole(ctx, user.UserID, req.TenantID); err != nil {
		// 记录错误但不阻止创建
		fmt.Printf("Warning: failed to assign default role: %v\n", err)
	}

	return user, nil
}

// assignDefaultRole 分配默认角色
func (s *UserService) assignDefaultRole(ctx context.Context, userID int64, tenantID string) error {
	// 查找租户的默认角色
	roles, _, err := s.roleRepo.List(ctx, &repository.RoleFilter{
		TenantID: tenantID,
		PageSize: 100,
	})
	if err != nil {
		return err
	}

	// 查找tenant_member角色
	var defaultRoleID string
	for _, role := range roles {
		if role.RoleCode == "tenant_member" {
			defaultRoleID = role.RoleID
			break
		}
	}

	if defaultRoleID == "" {
		return fmt.Errorf("default role not found")
	}

	// 分配角色
	userRole := &entity.UserRole{
		UserID:   fmt.Sprintf("%d", userID),
		TenantID: tenantID,
		RoleID:   defaultRoleID,
	}

	return s.userRoleRepo.Create(ctx, userRole)
}

// ListUsers 获取用户列表（支持分页、搜索）
func (s *UserService) ListUsers(ctx context.Context, filter *UserListFilter) ([]*entity.User, int64, error) {
	users, total, err := s.userRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "failed to list users"),
		)
	}

	return users, total, nil
}

// GetUser 获取用户详情
func (s *UserService) GetUser(ctx context.Context, userID int64) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("reason", "failed to get user"),
		)
	}

	if user == nil {
		return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("resource", "user"),
			errorx.KV("user_id", userID),
			errorx.KV("reason", "not found"),
		)
	}

	return user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, req *UpdateUserRequest) error {
	// 1. 获取现有用户
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("reason", "failed to get user"),
		)
	}
	if user == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("resource", "user"),
			errorx.KV("user_id", req.UserID),
			errorx.KV("reason", "not found"),
		)
	}

	// 2. 更新字段
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Description != nil {
		user.Description = *req.Description
	}
	if req.IconURI != nil {
		user.IconURI = *req.IconURI
		user.IconURL = *req.IconURI
	}
	if req.Locale != nil {
		user.Locale = *req.Locale
	}
	user.UpdatedAt = time.Now().UnixMilli()

	// 3. 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("reason", "failed to update user"),
		)
	}

	return nil
}

// DeleteUser 删除用户（软删除）
func (s *UserService) DeleteUser(ctx context.Context, userID int64) error {
	// 1. 删除用户角色关联
	if err := s.userRoleRepo.DeleteByUser(ctx, fmt.Sprintf("%d", userID), ""); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("reason", "failed to delete user roles"),
		)
	}

	// 2. 软删除用户
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("reason", "failed to delete user"),
		)
	}

	return nil
}

// ========== 2. 用户角色管理 ==========

// AssignRoles 分配角色给用户
func (s *UserService) AssignRoles(ctx context.Context, userID int64, tenantID string, roleIDs []string) error {
	// 1. 验证用户存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("reason", "failed to get user"),
		)
	}
	if user == nil || user.TenantID != tenantID {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("reason", "user not found"),
		)
	}

	// 2. 验证角色存在并分配
	for _, roleID := range roleIDs {
		// 检查角色是否存在
		role, err := s.roleRepo.GetByID(ctx, roleID)
		if err != nil {
			return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
				errorx.KV("reason", "failed to get role"),
			)
		}
		if role == nil || role.TenantID != tenantID {
			return errorx.New(errno.ErrPermissionInvalidParamCode,
				errorx.KV("resource", "role"),
				errorx.KV("role_id", roleID),
				errorx.KV("reason", "not found"),
			)
		}

		// 检查是否已经分配
		exists, err := s.userRoleRepo.Exists(ctx, fmt.Sprintf("%d", userID), tenantID, roleID)
		if err != nil {
			return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
				errorx.KV("reason", "failed to check role assignment"),
			)
		}
		if exists {
			continue // 已分配，跳过
		}

		// 创建用户角色关联
		userRole := &entity.UserRole{
			UserID:   fmt.Sprintf("%d", userID),
			TenantID: tenantID,
			RoleID:   roleID,
		}

		if err := s.userRoleRepo.Create(ctx, userRole); err != nil {
			return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
				errorx.KV("reason", "failed to assign role"),
			)
		}
	}

	return nil
}

// RemoveRole 移除用户角色
func (s *UserService) RemoveRole(ctx context.Context, userID int64, tenantID, roleID string) error {
	if err := s.userRoleRepo.Delete(ctx, fmt.Sprintf("%d", userID), tenantID, roleID); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("reason", "failed to revoke role"),
		)
	}
	return nil
}

// GetUserRoles 获取用户所有角色
func (s *UserService) GetUserRoles(ctx context.Context, userID int64, tenantID string) ([]*entity.Role, error) {
	roles, err := s.userRoleRepo.GetRolesByUser(ctx, fmt.Sprintf("%d", userID), tenantID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("reason", "failed to get user roles"),
		)
	}
	return roles, nil
}

// UpdateUserRoles 批量更新用户角色（先清除所有角色，再分配新角色）
func (s *UserService) UpdateUserRoles(ctx context.Context, userID int64, tenantID string, roleIDs []string) error {
	// 1. 删除用户所有角色
	if err := s.userRoleRepo.DeleteByUser(ctx, fmt.Sprintf("%d", userID), tenantID); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("reason", "failed to delete old roles"),
		)
	}

	// 2. 分配新角色
	if len(roleIDs) > 0 {
		if err := s.AssignRoles(ctx, userID, tenantID, roleIDs); err != nil {
			return err
		}
	}

	return nil
}

// ========== 3. 用户权限查询 ==========

// GetUserPermissions 获取用户所有权限
func (s *UserService) GetUserPermissions(ctx context.Context, userID int64, tenantID string) ([]*UserPermission, error) {
	// 1. 获取用户角色
	roles, err := s.GetUserRoles(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}

	// 2. 聚合所有角色的权限
	permissionsMap := make(map[string]*UserPermission)

	for _, role := range roles {
		// 获取角色的数据权限
		dataPerms, err := s.permChecker.dataPermRepo.GetByRole(ctx, role.RoleID)
		if err != nil {
			continue // 跳过错误，继续处理其他角色
		}

		for _, dp := range dataPerms {
			key := string(dp.ResourceType)
			if existing, ok := permissionsMap[key]; ok {
				// 合并权限（取最大权限）
				if dp.Scope > existing.Scope {
					existing.Scope = dp.Scope
				}
			} else {
				permissionsMap[key] = &UserPermission{
					ResourceType: dp.ResourceType,
					Actions:      []string{"read", "write", "delete"},
					Scope:        dp.Scope,
				}
			}
		}
	}

	// 3. 转换为列表
	permissions := make([]*UserPermission, 0, len(permissionsMap))
	for _, perm := range permissionsMap {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// GetUserEffectivePermissions 获取用户有效权限（含继承和临时授权）
func (s *UserService) GetUserEffectivePermissions(ctx context.Context, userID int64, tenantID string) ([]*UserPermission, error) {
	// 1. 获取基础权限
	permissions, err := s.GetUserPermissions(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}

	// 2. TODO: 叠加临时授权权限
	// 这里可以通过查询TemporaryGrantRepository获取临时授权，并叠加到基础权限上

	return permissions, nil
}

// ========== 4. 用户状态管理 ==========

// UpdateUserStatus 启用/禁用用户
func (s *UserService) UpdateUserStatus(ctx context.Context, userID int64, isEnabled bool) error {
	// 1. 获取现有用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("reason", "failed to get user"),
		)
	}
	if user == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("resource", "user"),
			errorx.KV("user_id", userID),
			errorx.KV("reason", "not found"),
		)
	}

	// 2. 更新状态
	user.IsEnabled = isEnabled
	user.UpdatedAt = time.Now().UnixMilli()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("reason", "failed to update user status"),
		)
	}

	return nil
}

// ResetPassword 重置用户密码
func (s *UserService) ResetPassword(ctx context.Context, userID int64, newPassword string) error {
	// 1. 验证密码强度
	if len(newPassword) < 8 {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("field", "new_password"),
			errorx.KV("reason", "password must be at least 8 characters"),
		)
	}

	// 2. 获取用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("reason", "failed to get user"),
		)
	}
	if user == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("resource", "user"),
			errorx.KV("user_id", userID),
			errorx.KV("reason", "not found"),
		)
	}

	// 3. TODO: 哈希密码并更新
	// 这里应该使用bcrypt等哈希算法对密码进行哈希处理
	// 然后调用userRepo.UpdatePassword(ctx, userID, hashedPassword)

	return nil
}

// UnlockUser 解锁用户账户
func (s *UserService) UnlockUser(ctx context.Context, userID int64) error {
	// TODO: 实现解锁逻辑
	// 1. 清除登录失败次数
	// 2. 清除锁定时间
	// 3. 更新用户状态

	return nil
}

// ========== 5. 用户批量操作 ==========

// BatchDeleteUsers 批量删除用户
func (s *UserService) BatchDeleteUsers(ctx context.Context, tenantID string, userIDs []int64) (*BatchOperationResult, error) {
	result := &BatchOperationResult{}

	for _, userID := range userIDs {
		// 验证用户存在且属于该租户
		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("user %d: %v", userID, err))
			continue
		}
		if user == nil || user.TenantID != tenantID {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("user %d: not found", userID))
			continue
		}

		// 删除用户
		if err := s.DeleteUser(ctx, userID); err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("user %d: %v", userID, err))
			continue
		}

		result.SuccessCount++
	}

	return result, nil
}

// BatchUpdateUsers 批量更新用户
func (s *UserService) BatchUpdateUsers(ctx context.Context, tenantID string, userIDs []int64, isEnabled *bool, roleIDs []string) (*BatchOperationResult, error) {
	result := &BatchOperationResult{}

	for _, userID := range userIDs {
		// 验证用户存在且属于该租户
		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("user %d: %v", userID, err))
			continue
		}
		if user == nil || user.TenantID != tenantID {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("user %d: not found", userID))
			continue
		}

		// 更新用户状态
		if isEnabled != nil {
			if err := s.UpdateUserStatus(ctx, userID, *isEnabled); err != nil {
				result.FailureCount++
				result.Errors = append(result.Errors, fmt.Sprintf("user %d: %v", userID, err))
				continue
			}
		}

		// 更新用户角色
		if roleIDs != nil {
			if err := s.UpdateUserRoles(ctx, userID, tenantID, roleIDs); err != nil {
				result.FailureCount++
				result.Errors = append(result.Errors, fmt.Sprintf("user %d: %v", userID, err))
				continue
			}
		}

		result.SuccessCount++
	}

	return result, nil
}

// ExportUsers 导出用户列表
func (s *UserService) ExportUsers(ctx context.Context, filter *UserListFilter) (string, error) {
	// 1. 查询用户列表
	users, _, err := s.userRepo.List(ctx, filter)
	if err != nil {
		return "", errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "failed to list users"),
		)
	}

	// 2. TODO: 生成Excel/CSV文件
	// 这里应该使用excelize等库生成Excel文件
	// 或者生成CSV文件

	// 3. TODO: 上传到存储服务
	// 这里应该上传到OSS或S3等存储服务
	// 返回可下载的URL

	return "", fmt.Errorf("export functionality not implemented yet")
}

// ========== 辅助函数 ==========

// generateUserID 生成用户ID
func generateUserID() int64 {
	return time.Now().UnixMilli()
}
