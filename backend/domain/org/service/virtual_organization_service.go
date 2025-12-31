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
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
)

// VirtualOrganizationService 虚拟组织服务
type VirtualOrganizationService struct {
	orgRepo    repository.VirtualOrganizationRepository
	memberRepo repository.VirtualOrgMemberRepository
	tagRepo    repository.VirtualOrgTagRepository
	empRepo    repository.EmployeeRepository
}

// NewVirtualOrganizationService 创建虚拟组织服务实例
func NewVirtualOrganizationService(
	orgRepo repository.VirtualOrganizationRepository,
	memberRepo repository.VirtualOrgMemberRepository,
	tagRepo repository.VirtualOrgTagRepository,
	empRepo repository.EmployeeRepository,
) *VirtualOrganizationService {
	return &VirtualOrganizationService{
		orgRepo:    orgRepo,
		memberRepo: memberRepo,
		tagRepo:    tagRepo,
		empRepo:    empRepo,
	}
}

// CreateVirtualOrg 创建虚拟组织
func (s *VirtualOrganizationService) CreateVirtualOrg(ctx context.Context, org *entity.VirtualOrganization) error {
	// 1. 验证参数
	if org.TenantID == "" || org.Name == "" || org.Type == "" || org.OwnerID == "" {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidParam"),
                )
	}

	// 2. 验证所有者存在
	emp, err := s.empRepo.GetByUserID(ctx, org.OwnerID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get owner"),
            )
	}
	if emp == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "EmployeeNotFound"),
                )
	}

	// 3. 生成虚拟组织ID
	org.VirtualOrgID = generateUUID()
	org.Status = entity.VirtualOrgStatusActive
	org.CreatedAt = time.Now().UnixMilli()
	org.UpdatedAt = time.Now().UnixMilli()

	// 4. 创建虚拟组织
	if err := s.orgRepo.Create(ctx, org); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "create virtual org"),
            )
	}

	// 5. 添加所有者为成员
	member := &entity.VirtualOrgMember{
		VirtualOrgID: org.VirtualOrgID,
		UserID:       org.OwnerID,
		Role:         entity.VirtualOrgRoleOwner,
		JoinedAt:     time.Now().UnixMilli(),
	}
	if err := s.memberRepo.AddMember(ctx, member); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "add owner as member"),
            )
	}

	return nil
}

// UpdateVirtualOrg 更新虚拟组织
func (s *VirtualOrganizationService) UpdateVirtualOrg(ctx context.Context, virtualOrgID string, org *entity.VirtualOrganization) error {
	// 1. 获取现有组织
	existing, err := s.orgRepo.GetByID(ctx, virtualOrgID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get virtual org"),
            )
	}
	if existing == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "VirtualOrgNotFound"),
                )
	}

	// 2. 更新
	org.VirtualOrgID = virtualOrgID
	org.UpdatedAt = time.Now().UnixMilli()
	if err := s.orgRepo.Update(ctx, org); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "update virtual org"),
            )
	}

	return nil
}

// DeleteVirtualOrg 删除虚拟组织
func (s *VirtualOrganizationService) DeleteVirtualOrg(ctx context.Context, virtualOrgID string) error {
	// 1. 获取现有组织
	existing, err := s.orgRepo.GetByID(ctx, virtualOrgID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get virtual org"),
            )
	}
	if existing == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "VirtualOrgNotFound"),
                )
	}

	// 2. 删除所有成员
	if err := s.memberRepo.RemoveMember(ctx, virtualOrgID, ""); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "remove members"),
            )
	}

	// 3. 删除组织
	if err := s.orgRepo.Delete(ctx, virtualOrgID); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "delete virtual org"),
            )
	}

	return nil
}

// GetVirtualOrg 获取虚拟组织
func (s *VirtualOrganizationService) GetVirtualOrg(ctx context.Context, id string) (*entity.VirtualOrganization, error) {
	if id == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get virtual org"),
            )
	}
	if org == nil {
		return nil, errorx.NewByErrorCode(errno.ErrVirtualOrgNotFound)
	}

	return org, nil
}

// ListVirtualOrgs 列出虚拟组织
func (s *VirtualOrganizationService) ListVirtualOrgs(ctx context.Context, tenantID string, orgType *entity.VirtualOrgType) ([]*entity.VirtualOrganization, error) {
	if tenantID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	var orgs []*entity.VirtualOrganization
	var err error

	if orgType != nil {
		orgs, err = s.orgRepo.SearchByType(ctx, tenantID, *orgType)
	} else {
		orgs, err = s.orgRepo.ListByTenantID(ctx, tenantID)
	}

	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "list virtual orgs"),
            )
	}

	return orgs, nil
}

// AddMember 添加成员
func (s *VirtualOrganizationService) AddMember(ctx context.Context, orgID, userID string, role entity.VirtualOrgMemberRole) error {
	// 1. 验证组织存在
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get virtual org"),
            )
	}
	if org == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "VirtualOrgNotFound"),
                )
	}

	// 2. 验证用户存在
	emp, err := s.empRepo.GetByUserID(ctx, userID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee"),
            )
	}
	if emp == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "EmployeeNotFound"),
                )
	}

	// 3. 检查是否已是成员
	existing, err := s.memberRepo.GetMember(ctx, orgID, userID)
	if err == nil && existing != nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "AlreadyMember"),
            )
	}

	// 4. 添加成员
	member := &entity.VirtualOrgMember{
		VirtualOrgID: orgID,
		UserID:       userID,
		Role:         role,
		JoinedAt:     time.Now().UnixMilli(),
	}
	if err := s.memberRepo.AddMember(ctx, member); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "add member"),
            )
	}

	return nil
}

// RemoveMember 移除成员
func (s *VirtualOrganizationService) RemoveMember(ctx context.Context, orgID, userID string) error {
	// 1. 验证组织存在
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get virtual org"),
            )
	}
	if org == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "VirtualOrgNotFound"),
                )
	}

	// 2. 不能移除所有者
	member, err := s.memberRepo.GetMember(ctx, orgID, userID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get member"),
            )
	}
	if member.Role == entity.VirtualOrgRoleOwner {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "CannotRemoveOwner"),
            )
	}

	// 3. 移除成员
	if err := s.memberRepo.RemoveMember(ctx, orgID, userID); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "remove member"),
            )
	}

	return nil
}

// UpdateMemberRole 更新成员角色
func (s *VirtualOrganizationService) UpdateMemberRole(ctx context.Context, orgID, userID string, role entity.VirtualOrgMemberRole) error {
	// 1. 验证成员存在
	member, err := s.memberRepo.GetMember(ctx, orgID, userID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get member"),
            )
	}
	if member == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "MemberNotFound"),
            )
	}

	// 2. 不能修改所有者角色
	if member.Role == entity.VirtualOrgRoleOwner {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "CannotChangeOwnerRole"),
            )
	}

	// 3. 更新角色
	if err := s.memberRepo.UpdateMemberRole(ctx, orgID, userID, role); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "update member role"),
            )
	}

	return nil
}

// ListMembers 列出成员
func (s *VirtualOrganizationService) ListMembers(ctx context.Context, orgID string) ([]*entity.VirtualOrgMember, error) {
	if orgID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	members, err := s.memberRepo.ListMembers(ctx, orgID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "list members"),
            )
	}

	return members, nil
}

// AddTag 添加标签
func (s *VirtualOrganizationService) AddTag(ctx context.Context, orgID, tag string) error {
	// 1. 验证组织存在
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get virtual org"),
            )
	}
	if org == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "VirtualOrgNotFound"),
                )
	}

	// 2. 添加标签
	tagEntity := &entity.VirtualOrgTag{
		VirtualOrgID: orgID,
		TagName:      tag,
	}
	if err := s.tagRepo.AddTag(ctx, tagEntity); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "add tag"),
            )
	}

	return nil
}

// RemoveTag 移除标签
func (s *VirtualOrganizationService) RemoveTag(ctx context.Context, orgID, tag string) error {
	if err := s.tagRepo.RemoveTag(ctx, orgID, tag); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "remove tag"),
            )
	}
	return nil
}

// ListTags 列出标签
func (s *VirtualOrganizationService) ListTags(ctx context.Context, orgID string) ([]*entity.VirtualOrgTag, error) {
	if orgID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	tags, err := s.tagRepo.ListTags(ctx, orgID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "list tags"),
            )
	}

	return tags, nil
}

// InheritPermissions 继承权限（简化实现）
func (s *VirtualOrganizationService) InheritPermissions(ctx context.Context, orgID string) ([]*entity.Permission, error) {
	// 简化实现：返回虚拟组织的基本权限
	// 实际应用中应该从权限系统查询并返回
	permissions := []*entity.Permission{
		{
			PermissionID:  "perm-1",
			ResourceID:    orgID,
			ResourceType:  "virtual_org",
			Action:        "read",
			Effect:        "allow",
		},
		{
			PermissionID:  "perm-2",
			ResourceID:    orgID,
			ResourceType:  "virtual_org",
			Action:        "write",
			Effect:        "allow",
		},
	}

	return permissions, nil
}
