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

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// OrganizationService 组织管理服务
// 职责：组织的CRUD操作、层级管理、业务规则验证
type OrganizationService struct {
	orgRepo  repository.OrganizationRepository
	treeRepo repository.OrganizationTreeRepository
	db       *gorm.DB
}

// NewOrganizationService 创建组织服务实例
func NewOrganizationService(
	orgRepo repository.OrganizationRepository,
	treeRepo repository.OrganizationTreeRepository,
	db *gorm.DB,
) *OrganizationService {
	return &OrganizationService{
		orgRepo:  orgRepo,
		treeRepo: treeRepo,
		db:       db,
	}
}

// CreateOrganizationRequest 创建组织请求
type CreateOrganizationRequest struct {
	TenantID     string                       `json:"tenant_id" binding:"required"`
	OrgName      string                       `json:"org_name" binding:"required,min=1,max=200"`
	OrgType      entity.OrganizationType      `json:"org_type" binding:"required,oneof=company division department project"`
	ParentID     *string                      `json:"parent_id,omitempty"`
	OrgCode      string                       `json:"org_code" binding:"required,min=1,max=50"`
	LeaderID     *string                      `json:"leader_id,omitempty"`
	Description string                       `json:"description,omitempty"`
	SortOrder    int                          `json:"sort_order,omitempty"`
}

// UpdateOrganizationRequest 更新组织请求
type UpdateOrganizationRequest struct {
	OrgID        string                       `json:"org_id" binding:"required"`
	OrgName      string                       `json:"org_name" binding:"omitempty,min=1,max=200"`
	LeaderID     *string                      `json:"leader_id,omitempty"`
	Description  string                       `json:"description,omitempty"`
	SortOrder    int                          `json:"sort_order,omitempty"`
	Status       entity.OrganizationStatus    `json:"status" binding:"omitempty,oneof=active inactive frozen"`
}

// MoveOrganizationRequest 移动组织请求
type MoveOrganizationRequest struct {
	OrgID       string  `json:"org_id" binding:"required"`
	NewParentID *string `json:"new_parent_id,omitempty"`
}

// CreateOrganization 创建组织
func (s *OrganizationService) CreateOrganization(ctx context.Context, req *CreateOrganizationRequest) (*entity.Organization, error) {
	// 1. 验证租户唯一性
	orgs, err := s.orgRepo.GetByTenantID(ctx, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to check tenant organizations: %w", err)
	}

	// 2. 验证组织类型规则
	if req.OrgType == entity.OrgTypeCompany && len(orgs) > 0 {
		return nil, errno.ErrTenantAlreadyHasOrg
	}

	// 3. 如果有父组织，验证父组织存在且不是公司类型
	var parent *entity.Organization
	var level int
	var path string

	if req.ParentID != nil && *req.ParentID != "" {
		parent, err = s.orgRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent organization: %w", err)
		}
		if parent == nil {
			return nil, errno.ErrParentOrgNotFound
		}
		if parent.TenantID != req.TenantID {
			return nil, errno.ErrTenantMismatch
		}
		if parent.OrgType == entity.OrgTypeDepartment {
			return nil, errno.ErrInvalidParentOrg
		}

		level = parent.Level + 1
		path = parent.Path + "/"
	} else {
		if req.OrgType != entity.OrgTypeCompany {
			return nil, errno.ErrRootOrgMustBeCompany
		}
		level = 1
		path = "/"
	}

	// 4. 验证编码唯一性
	exists, err := s.orgRepo.ExistsByCode(ctx, req.TenantID, req.OrgCode, "")
	if err != nil {
		return nil, fmt.Errorf("failed to check org code: %w", err)
	}
	if exists {
		return nil, errno.ErrOrgCodeAlreadyExists
	}

	// 5. 生成组织ID
	orgID := generateUUID()

	// 6. 创建组织实体
	org := &entity.Organization{
		OrgID:        orgID,
		TenantID:     req.TenantID,
		OrgName:      req.OrgName,
		OrgType:      req.OrgType,
		ParentID:     req.ParentID,
		OrgCode:      req.OrgCode,
		Level:        level,
		Path:         path + orgID,
		SortOrder:    req.SortOrder,
		Status:       entity.OrgStatusActive,
		Description:  req.Description,
		LeaderID:     req.LeaderID,
		CreatedAt:    time.Now().UnixMilli(),
		UpdatedAt:    time.Now().UnixMilli(),
	}

	// 7. 使用事务创建组织和路径
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 创建组织
		if err := s.orgRepo.Create(ctx, org); err != nil {
			return fmt.Errorf("failed to create organization: %w", err)
		}

		// 创建闭包表路径
		if err := s.createOrganizationPaths(ctx, org, parent); err != nil {
			return fmt.Errorf("failed to create organization paths: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return org, nil
}

// createOrganizationPaths 创建组织的所有闭包表路径
func (s *OrganizationService) createOrganizationPaths(ctx context.Context, org *entity.Organization, parent *entity.Organization) error {
	paths := make([]*entity.OrganizationTree, 0, org.Level)

	// 添加自己到自己的路径（深度0）
	paths = append(paths, &entity.OrganizationTree{
		TenantID:     org.TenantID,
		AncestorID:   org.OrgID,
		DescendantID: org.OrgID,
		Depth:        0,
	})

	// 如果有父组织，添加所有祖先的路径
	if parent != nil {
		ancestorPaths, err := s.treeRepo.GetAncestors(ctx, parent.OrgID)
		if err != nil {
			return fmt.Errorf("failed to get ancestor paths: %w", err)
		}

		for _, ancestorPath := range ancestorPaths {
			paths = append(paths, &entity.OrganizationTree{
				TenantID:     org.TenantID,
				AncestorID:   ancestorPath.AncestorID,
				DescendantID: org.OrgID,
				Depth:        ancestorPath.Depth + 1,
			})
		}
	}

	// 批量创建路径
	return s.treeRepo.CreatePaths(ctx, paths)
}

// GetOrganization 获取组织详情
func (s *OrganizationService) GetOrganization(ctx context.Context, orgID string) (*entity.Organization, error) {
	if orgID == "" {
		return nil, errno.ErrInvalidParam
	}

	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	if org == nil {
		return nil, errno.ErrOrgNotFound
	}

	return org, nil
}

// UpdateOrganization 更新组织
func (s *OrganizationService) UpdateOrganization(ctx context.Context, req *UpdateOrganizationRequest) (*entity.Organization, error) {
	// 1. 获取组织
	org, err := s.orgRepo.GetByID(ctx, req.OrgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	if org == nil {
		return nil, errno.ErrOrgNotFound
	}

	// 2. 更新字段
	if req.OrgName != "" {
		org.OrgName = req.OrgName
	}
	if req.LeaderID != nil {
		org.LeaderID = req.LeaderID
	}
	if req.Description != "" {
		org.Description = req.Description
	}
	if req.SortOrder != 0 {
		org.SortOrder = req.SortOrder
	}
	if req.Status != "" {
		// 验证状态变更规则
		if req.Status == entity.OrgStatusFrozen && org.OrgType == entity.OrgTypeCompany {
			return nil, errno.ErrCannotFreezeRootOrg
		}
		org.Status = req.Status
	}

	// 3. 保存更新
	if err := s.orgRepo.Update(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	return org, nil
}

// DeleteOrganization 删除组织
func (s *OrganizationService) DeleteOrganization(ctx context.Context, orgID string) error {
	if orgID == "" {
		return errno.ErrInvalidParam
	}

	// 1. 获取组织
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return fmt.Errorf("failed to get organization: %w", err)
	}
	if org == nil {
		return errno.ErrOrgNotFound
	}

	// 2. 检查是否有子组织
	children, err := s.orgRepo.GetChildren(ctx, orgID)
	if err != nil {
		return fmt.Errorf("failed to check children: %w", err)
	}
	if len(children) > 0 {
		return errno.ErrOrgHasChildren
	}

	// 3. 公司类型不能删除
	if org.OrgType == entity.OrgTypeCompany {
		return errno.ErrCannotDeleteRootOrg
	}

	// 4. 使用事务删除组织和路径
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 软删除组织
		if err := s.orgRepo.Delete(ctx, orgID); err != nil {
			return fmt.Errorf("failed to delete organization: %w", err)
		}

		// 删除闭包表路径
		if err := s.treeRepo.DeletePaths(ctx, orgID); err != nil {
			return fmt.Errorf("failed to delete organization paths: %w", err)
		}

		return nil
	})
}

// MoveOrganization 移动组织
func (s *OrganizationService) MoveOrganization(ctx context.Context, req *MoveOrganizationRequest) error {
	if req.OrgID == "" {
		return errno.ErrInvalidParam
	}

	// 1. 获取要移动的组织
	org, err := s.orgRepo.GetByID(ctx, req.OrgID)
	if err != nil {
		return fmt.Errorf("failed to get organization: %w", err)
	}
	if org == nil {
		return errno.ErrOrgNotFound
	}

	// 2. 验证不能移动到自己
	if req.NewParentID != nil && *req.NewParentID == org.OrgID {
		return errno.ErrCannotMoveToSelf
	}

	// 3. 验证公司类型不能移动
	if org.OrgType == entity.OrgTypeCompany {
		return errno.ErrCannotMoveRootOrg
	}

	// 4. 获取新旧父组织
	var newParent *entity.Organization
	var newLevel int
	var newPath string

	if req.NewParentID != nil && *req.NewParentID != "" {
		newParent, err = s.orgRepo.GetByID(ctx, *req.NewParentID)
		if err != nil {
			return fmt.Errorf("failed to get new parent: %w", err)
		}
		if newParent == nil {
			return errno.ErrParentOrgNotFound
		}
		if newParent.TenantID != org.TenantID {
			return errno.ErrTenantMismatch
		}
		if newParent.OrgType == entity.OrgTypeDepartment {
			return nil, errno.ErrInvalidParentOrg
		}

		// 验证不能移动到自己的后代
		isDescendant, err := s.isDescendant(ctx, org.OrgID, newParent.OrgID)
		if err != nil {
			return fmt.Errorf("failed to check descendant: %w", err)
		}
		if isDescendant {
			return errno.ErrCannotMoveToDescendant
		}

		newLevel = newParent.Level + 1
		newPath = newParent.Path + "/" + org.OrgID
	} else {
		// 移动到根级别
		newLevel = 1
		newPath = "/" + org.OrgID
	}

	// 5. 使用事务移动组织和更新所有后代的路径
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 更新组织的层级和路径
		org.ParentID = req.NewParentID
		org.Level = newLevel
		org.Path = newPath
		org.UpdatedAt = time.Now().UnixMilli()

		if err := s.orgRepo.Update(ctx, org); err != nil {
			return fmt.Errorf("failed to update organization: %w", err)
		}

		// 更新闭包表路径
		if err := s.treeRepo.MoveSubtree(ctx, org.OrgID, *req.NewParentID); err != nil {
			return fmt.Errorf("failed to move subtree: %w", err)
		}

		// TODO: 更新所有后代的层级和路径
		// 这里需要递归更新所有子组织

		return nil
	})
}

// isDescendant 检查是否是后代节点
func (s *OrganizationService) isDescendant(ctx context.Context, ancestorID, descendantID string) (bool, error) {
	paths, err := s.treeRepo.GetAncestors(ctx, descendantID)
	if err != nil {
		return false, err
	}

	for _, path := range paths {
		if path.AncestorID == ancestorID {
			return true, nil
		}
	}

	return false, nil
}

// GetOrganizationTree 获取组织树
func (s *OrganizationService) GetOrganizationTree(ctx context.Context, tenantID string) ([]*entity.Organization, error) {
	if tenantID == "" {
		return nil, errno.ErrInvalidParam
	}

	orgs, err := s.orgRepo.GetTree(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization tree: %w", err)
	}

	// 构建树形结构
	return s.buildOrgTree(orgs, ""), nil
}

// buildOrgTree 构建组织树
func (s *OrganizationService) buildOrgTree(orgs []*entity.Organization, parentID string) []*entity.Organization {
	var tree []*entity.Organization

	for _, org := range orgs {
		isChild := false
		if parentID == "" {
			// 根节点
			isChild = org.ParentID == nil || *org.ParentID == ""
		} else {
			// 子节点
			isChild = org.ParentID != nil && *org.ParentID == parentID
		}

		if isChild {
			org.Children = s.buildOrgTree(orgs, org.OrgID)
			tree = append(tree, org)
		}
	}

	return tree
}

// OrganizationFilter 组织查询过滤器（扩展版，支持排序）
type OrganizationFilter struct {
	TenantID  string
	OrgType   string // entity.OrganizationType
	Status    string // entity.OrganizationStatus
	ParentID  string
	Keyword   string
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
}

// ListOrganizations 分页查询组织列表
func (s *OrganizationService) ListOrganizations(ctx context.Context, filter *OrganizationFilter) ([]*entity.Organization, int64, error) {
	if filter.TenantID == "" {
		return nil, 0, errno.ErrInvalidParam
	}

	// 转换为仓储层过滤器
	repoFilter := &repository.OrganizationFilter{
		TenantID:  filter.TenantID,
		OrgType:   entity.OrganizationType(filter.OrgType),
		Status:    entity.OrganizationStatus(filter.Status),
		Keyword:   filter.Keyword,
		PageSize:  filter.PageSize,
		PageToken: fmt.Sprintf("%d", filter.Page),
	}

	// 如果有ParentID，设置过滤
	if filter.ParentID != "" {
		repoFilter.ParentID = &filter.ParentID
	}

	orgs, total, err := s.orgRepo.List(ctx, repoFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list organizations: %w", err)
	}

	return orgs, total, nil
}

// GetChildren 获取子组织
func (s *OrganizationService) GetChildren(ctx context.Context, parentID string) ([]*entity.Organization, error) {
	if parentID == "" {
		return nil, errno.ErrInvalidParam
	}

	orgs, err := s.orgRepo.GetChildren(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get children: %w", err)
	}

	return orgs, nil
}

// GetAncestors 获取组织的所有祖先
func (s *OrganizationService) GetAncestors(ctx context.Context, orgID string) ([]*entity.Organization, error) {
	if orgID == "" {
		return nil, errno.ErrInvalidParam
	}

	paths, err := s.treeRepo.GetAncestors(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ancestor paths: %w", err)
	}

	// 转换为组织实体
	var ancestors []*entity.Organization
	for _, path := range paths {
		if path.Depth > 0 { // 排除自己
			org, err := s.orgRepo.GetByID(ctx, path.AncestorID)
			if err != nil {
				return nil, fmt.Errorf("failed to get ancestor organization: %w", err)
			}
			if org != nil {
				ancestors = append(ancestors, org)
			}
		}
	}

	return ancestors, nil
}

// GetDescendants 获取组织的所有后代
func (s *OrganizationService) GetDescendants(ctx context.Context, orgID string) ([]*entity.Organization, error) {
	if orgID == "" {
		return nil, errno.ErrInvalidParam
	}

	paths, err := s.treeRepo.GetDescendants(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get descendant paths: %w", err)
	}

	// 转换为组织实体
	var descendants []*entity.Organization
	for _, path := range paths {
		if path.Depth > 0 { // 排除自己
			org, err := s.orgRepo.GetByID(ctx, path.DescendantID)
			if err != nil {
				return nil, fmt.Errorf("failed to get descendant organization: %w", err)
			}
			if org != nil {
				descendants = append(descendants, org)
			}
		}
	}

	return descendants, nil
}

// generateUUID 生成UUID
func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
