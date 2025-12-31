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
	"encoding/json"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	permRepository "github.com/coze-dev/coze-studio/backend/domain/permission/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// FieldPermissionService 字段权限服务
// 职责：三级字段权限控制、动态字段过滤、API响应过滤
type FieldPermissionService struct {
	fieldPermRepo repository.FieldPermissionRepository
	roleRepo      permRepository.RoleRepository
}

// NewFieldPermissionService 创建字段权限服务实例
func NewFieldPermissionService(
	fieldPermRepo repository.FieldPermissionRepository,
	roleRepo permRepository.RoleRepository,
) *FieldPermissionService {
	return &FieldPermissionService{
		fieldPermRepo: fieldPermRepo,
		roleRepo:      roleRepo,
	}
}

// GetFieldPermissions 获取字段权限配置
func (s *FieldPermissionService) GetFieldPermissions(
	ctx context.Context,
	roleID, resourceType string,
) ([]*entity.FieldPermission, error) {
	permissions, err := s.fieldPermRepo.GetFieldPermissions(ctx, roleID, resourceType)
	if err != nil {
		return nil, errorx.Wrapf(err, "get field permissions failed: role_id=%s, resource_type=%s", roleID, resourceType)
	}

	return permissions, nil
}

// GetUserFieldPermissions 获取用户的字段权限
func (s *FieldPermissionService) GetUserFieldPermissions(
	ctx context.Context,
	userID, resourceType string,
) ([]*entity.FieldPermission, error) {
	permissions, err := s.fieldPermRepo.GetUserFieldPermissions(ctx, userID, resourceType)
	if err != nil {
		return nil, errorx.Wrapf(err, "get user field permissions failed: user_id=%s, resource_type=%s", userID, resourceType)
	}

	return permissions, nil
}

// FilterFieldsByPermission 根据权限过滤字段
func (s *FieldPermissionService) FilterFieldsByPermission(
	ctx context.Context,
	userID, resourceType string,
	data map[string]interface{},
) (map[string]interface{}, error) {
	// 1. 获取用户的字段权限
	permissions, err := s.GetUserFieldPermissions(ctx, userID, resourceType)
	if err != nil {
		return nil, err
	}

	// 2. 构建字段权限映射
	fieldPerms := make(map[string]entity.FieldPerm)
	for _, perm := range permissions {
		fieldPerms[perm.FieldName] = perm.Permission
	}

	// 3. 过滤字段
	result := make(map[string]interface{})
	for fieldName, fieldValue := range data {
		perm, exists := fieldPerms[fieldName]
		if !exists {
			// 如果没有配置权限，默认可见
			result[fieldName] = fieldValue
			continue
		}

		// 根据权限级别处理
		switch perm {
		case entity.FieldPermVisible:
			// 可见，保留字段
			result[fieldName] = fieldValue
		case entity.FieldPermEditable:
			// 可编辑，保留字段
			result[fieldName] = fieldValue
		case entity.FieldPermHidden:
			// 隐藏，移除字段
			continue
		}
	}

	return result, nil
}

// FilterFieldsByPermissionEditable 只保留可编辑字段
func (s *FieldPermissionService) FilterFieldsByPermissionEditable(
	ctx context.Context,
	userID, resourceType string,
	data map[string]interface{},
) (map[string]interface{}, error) {
	// 1. 获取用户的字段权限
	permissions, err := s.GetUserFieldPermissions(ctx, userID, resourceType)
	if err != nil {
		return nil, err
	}

	// 2. 构建可编辑字段集合
	editableFields := make(map[string]bool)
	for _, perm := range permissions {
		if perm.Permission == entity.FieldPermEditable {
			editableFields[perm.FieldName] = true
		}
	}

	// 3. 过滤字段，只保留可编辑的
	result := make(map[string]interface{})
	for fieldName, fieldValue := range data {
		if editableFields[fieldName] {
			result[fieldName] = fieldValue
		}
	}

	return result, nil
}

// BatchGetFieldPermissions 批量获取字段权限
func (s *FieldPermissionService) BatchGetFieldPermissions(
	ctx context.Context,
	userIDs []string,
	resourceType string,
) (map[string][]*entity.FieldPermission, error) {
	permissions, err := s.fieldPermRepo.BatchGetFieldPermissions(ctx, userIDs, resourceType)
	if err != nil {
		return nil, errorx.Wrapf(err, "batch get field permissions failed: resource_type=%s", resourceType)
	}

	return permissions, nil
}

// SetFieldPermission 设置字段权限
func (s *FieldPermissionService) SetFieldPermission(
	ctx context.Context,
	req *SetFieldPermissionRequest,
) error {
	// 1. 验证请求参数
	if req.ResourceType == "" || req.FieldName == "" {
		return errorx.New(errno.ErrInvalidParamCode)
	}

	// 2. 验证权限级别
	validPerms := map[entity.FieldPerm]bool{
		entity.FieldPermVisible:  true,
		entity.FieldPermEditable: true,
		entity.FieldPermHidden:   true,
	}
	if !validPerms[req.Permission] {
		return errorx.New(errno.ErrInvalidParamCode)
	}

	// 3. 创建字段权限实体
	_ = &entity.FieldPermission{
		PermissionID: generateID(),
		RoleID:       req.RoleID,
		ResourceType: req.ResourceType,
		FieldName:    req.FieldName,
		Permission:   req.Permission,
		CreatedAt:    getCurrentTimestamp(),
		UpdatedAt:    getCurrentTimestamp(),
	}

	// 4. 保存到数据库（实际实现需要调用dal层）
	// 这里简化处理

	return nil
}

// BatchSetFieldPermissions 批量设置字段权限
func (s *FieldPermissionService) BatchSetFieldPermissions(
	ctx context.Context,
	req *BatchSetFieldPermissionsRequest,
) error {
	// 1. 验证请求
	if len(req.Permissions) == 0 {
		return errorx.New(errno.ErrInvalidParamCode)
	}

	// 2. 批量设置权限
	for _, perm := range req.Permissions {
		setReq := &SetFieldPermissionRequest{
			RoleID:       req.RoleID,
			ResourceType: req.ResourceType,
			FieldName:    perm.FieldName,
			Permission:   perm.Permission,
		}

		if err := s.SetFieldPermission(ctx, setReq); err != nil {
			return err
		}
	}

	return nil
}

// DeleteFieldPermission 删除字段权限
func (s *FieldPermissionService) DeleteFieldPermission(
	ctx context.Context,
	permissionID string,
) error {
	// 1. 验证权限ID
	if permissionID == "" {
		return errorx.New(errno.ErrInvalidParamCode)
	}

	// 2. 删除权限（实际实现需要调用dal层）
	// 这里简化处理

	return nil
}

// GetFieldPermissionSummary 获取字段权限摘要
func (s *FieldPermissionService) GetFieldPermissionSummary(
	ctx context.Context,
	userID, resourceType string,
) (*FieldPermissionSummary, error) {
	// 1. 获取字段权限
	permissions, err := s.GetUserFieldPermissions(ctx, userID, resourceType)
	if err != nil {
		return nil, err
	}

	// 2. 统计各级别字段数量
	summary := &FieldPermissionSummary{
		UserID:      userID,
		ResourceType: resourceType,
		FieldCount:  len(permissions),
		VisibleCount: 0,
		EditableCount: 0,
		HiddenCount:   0,
	}

	for _, perm := range permissions {
		switch perm.Permission {
		case entity.FieldPermVisible:
			summary.VisibleCount++
		case entity.FieldPermEditable:
			summary.EditableCount++
		case entity.FieldPermHidden:
			summary.HiddenCount++
		}
	}

	return summary, nil
}

// FilterAPIResponse 过滤API响应（用于HTTP中间件）
func (s *FieldPermissionService) FilterAPIResponse(
	ctx context.Context,
	userID, resourceType string,
	responseData interface{},
) (interface{}, error) {
	// 1. 转换为map
	data, ok := responseData.(map[string]interface{})
	if !ok {
		// 如果不是map，尝试JSON序列化后解析
		jsonBytes, err := json.Marshal(responseData)
		if err != nil {
			return nil, errorx.Wrapf(err, "marshal response data failed")
		}

		if err := json.Unmarshal(jsonBytes, &data); err != nil {
			return nil, errorx.Wrapf(err, "unmarshal response data failed")
		}
	}

	// 2. 过滤字段
	filteredData, err := s.FilterFieldsByPermission(ctx, userID, resourceType, data)
	if err != nil {
		return nil, err
	}

	return filteredData, nil
}

// ValidateFieldAccess 验证字段访问权限
func (s *FieldPermissionService) ValidateFieldAccess(
	ctx context.Context,
	userID, resourceType, fieldName string,
	permission entity.FieldPerm,
) error {
	// 1. 获取用户字段权限
	userPerms, err := s.GetUserFieldPermissions(ctx, userID, resourceType)
	if err != nil {
		return err
	}

	// 2. 查找字段权限
	for _, perm := range userPerms {
		if perm.FieldName == fieldName {
			// 3. 验证权限级别
			if perm.Permission == entity.FieldPermHidden {
				return errorx.New(errno.ErrPermissionCheckFailedCode, errorx.KV("reason", "字段不可访问"))
			}

			if permission == entity.FieldPermEditable && perm.Permission != entity.FieldPermEditable {
				return errorx.New(errno.ErrPermissionCheckFailedCode, errorx.KV("reason", "字段不可编辑"))
			}

			return nil
		}
	}

	// 如果没有配置权限，默认允许
	return nil
}

// CloneFieldPermissions 克隆字段权限配置
func (s *FieldPermissionService) CloneFieldPermissions(
	ctx context.Context,
	sourceRoleID, targetRoleID, resourceType string,
) error {
	// 1. 获取源角色权限
	permissions, err := s.GetFieldPermissions(ctx, sourceRoleID, resourceType)
	if err != nil {
		return err
	}

	// 2. 批量设置到目标角色
	for _, perm := range permissions {
		setReq := &SetFieldPermissionRequest{
			RoleID:       targetRoleID,
			ResourceType: resourceType,
			FieldName:    perm.FieldName,
			Permission:   perm.Permission,
		}

		if err := s.SetFieldPermission(ctx, setReq); err != nil {
			return err
		}
	}

	return nil
}

// GetFieldPermissionTemplate 获取字段权限模板
func (s *FieldPermissionService) GetFieldPermissionTemplate(
	resourceType string,
) (*FieldPermissionTemplate, error) {
	// 预定义的权限模板
	templates := map[string]*FieldPermissionTemplate{
		"employee": {
			ResourceType: "employee",
			Fields: []FieldTemplate{
				{FieldName: "emp_id", Permission: entity.FieldPermVisible},
				{FieldName: "name", Permission: entity.FieldPermVisible},
				{FieldName: "email", Permission: entity.FieldPermVisible},
				{FieldName: "phone", Permission: entity.FieldPermEditable},
				{FieldName: "salary", Permission: entity.FieldPermHidden},
				{FieldName: "id_card", Permission: entity.FieldPermHidden},
			},
		},
	}

	template, exists := templates[resourceType]
	if !exists {
		// 返回空模板
		return &FieldPermissionTemplate{
			ResourceType: resourceType,
			Fields:       []FieldTemplate{},
		}, nil
	}

	return template, nil
}

// SetFieldPermissionRequest 设置字段权限请求
type SetFieldPermissionRequest struct {
	RoleID       string                 `json:"role_id" binding:"required"`
	ResourceType string                 `json:"resource_type" binding:"required"`
	FieldName    string                 `json:"field_name" binding:"required"`
	Permission   entity.FieldPerm       `json:"permission" binding:"required,oneof=visible editable hidden"`
}

// BatchSetFieldPermissionsRequest 批量设置字段权限请求
type BatchSetFieldPermissionsRequest struct {
	RoleID       string                   `json:"role_id" binding:"required"`
	ResourceType string                   `json:"resource_type" binding:"required"`
	Permissions  []FieldPermissionItem    `json:"permissions" binding:"required"`
}

// FieldPermissionItem 字段权限项
type FieldPermissionItem struct {
	FieldName  string              `json:"field_name" binding:"required"`
	Permission entity.FieldPerm    `json:"permission" binding:"required,oneof=visible editable hidden"`
}

// FieldPermissionSummary 字段权限摘要
type FieldPermissionSummary struct {
	UserID        string `json:"user_id"`
	ResourceType  string `json:"resource_type"`
	FieldCount    int    `json:"field_count"`
	VisibleCount  int    `json:"visible_count"`
	EditableCount int    `json:"editable_count"`
	HiddenCount   int    `json:"hidden_count"`
}

// FieldPermissionTemplate 字段权限模板
type FieldPermissionTemplate struct {
	ResourceType string          `json:"resource_type"`
	Fields       []FieldTemplate `json:"fields"`
}

// FieldTemplate 字段模板项
type FieldTemplate struct {
	FieldName  string              `json:"field_name"`
	Permission entity.FieldPerm    `json:"permission"`
}

// 辅助函数
func generateID() string {
	// 实际实现应该使用UUID生成器
	return "id_placeholder"
}

func getCurrentTimestamp() int64 {
	// 实际实现应该返回当前时间戳
	return 0
}
