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

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
)

// OrgTemplateService 组织模板服务
// 职责：组织模板创建、模板应用、预置模板库、模板版本管理
type OrgTemplateService struct {
	db              *gorm.DB
	templateRepo    OrgTemplateRepository
	departmentRepo  repository.DepartmentRepository
	employeeRepo    repository.EmployeeRepository
}

// OrgTemplateRepository 组织模板仓储接口
type OrgTemplateRepository interface {
	Create(ctx context.Context, template *entity.OrgTemplate) error
	GetByID(ctx context.Context, templateID string) (*entity.OrgTemplate, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]*entity.OrgTemplate, error)
	GetBuiltinTemplates(ctx context.Context) ([]*entity.OrgTemplate, error)
	Update(ctx context.Context, template *entity.OrgTemplate) error
	Delete(ctx context.Context, templateID string) error
	List(ctx context.Context, filter *TemplateFilter) ([]*entity.OrgTemplate, int64, error)
}

// TemplateFilter 模板查询过滤器
type TemplateFilter struct {
	TenantID   string
	IsBuiltin  *bool
	IsEnabled  *bool
	Keyword    string
	PageToken  string
	PageSize   int
}

// NewOrgTemplateService 创建组织模板服务实例
func NewOrgTemplateService(
	db *gorm.DB,
	templateRepo OrgTemplateRepository,
	departmentRepo repository.DepartmentRepository,
	employeeRepo repository.EmployeeRepository,
) *OrgTemplateService {
	return &OrgTemplateService{
		db:             db,
		templateRepo:   templateRepo,
		departmentRepo: departmentRepo,
		employeeRepo:   employeeRepo,
	}
}

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	TenantID    string                 `json:"tenant_id" binding:"required"`
	Name        string                 `json:"name" binding:"required,min=1,max=200"`
	Description string                 `json:"description,omitempty"`
	Config      *entity.TemplateConfig `json:"config" binding:"required"`
	Version     string                 `json:"version" binding:"required"`
	CreatedBy   string                 `json:"created_by" binding:"required"`
}

// CreateTemplate 创建模板
func (s *OrgTemplateService) CreateTemplate(
	ctx context.Context,
	req *CreateTemplateRequest,
) (*entity.OrgTemplate, error) {
	// 1. 参数验证
	if req.TenantID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	// 2. 生成模板ID
	templateID := generateUUID()

	// 3. 创建模板实体
	now := time.Now().UnixMilli()
	template := &entity.OrgTemplate{
		TemplateID:   templateID,
		TenantID:     req.TenantID,
		Name:         req.Name,
		Description:  req.Description,
		Config:       req.Config,
		IsBuiltin:    false,
		Version:      req.Version,
		IsEnabled:    true,
		CreatedBy:    req.CreatedBy,
		UpdatedBy:    req.CreatedBy,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 4. 创建模板
	if err := s.templateRepo.Create(ctx, template); err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "create template"),
		)
	}

	return template, nil
}

// ApplyTemplateRequest 应用模板请求
type ApplyTemplateRequest struct {
	TenantID   string   `json:"tenant_id" binding:"required"`
	OrgID      string   `json:"org_id" binding:"required"`
	TemplateID string   `json:"template_id" binding:"required"`
	Options    *ApplyOptions `json:"options,omitempty"`
	AppliedBy  string   `json:"applied_by" binding:"required"`
}

// ApplyOptions 应用选项
type ApplyOptions struct {
	IncludeDepartments bool `json:"include_departments"`
	IncludeRoles       bool `json:"include_roles"`
	IncludePermissions bool `json:"include_permissions"`
	IncludeWorkflows   bool `json:"include_workflows"`
}

// ApplyTemplate 应用模板
func (s *OrgTemplateService) ApplyTemplate(
	ctx context.Context,
	req *ApplyTemplateRequest,
) error {
	// 1. 获取模板
	template, err := s.templateRepo.GetByID(ctx, req.TemplateID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get template"),
		)
	}
	if template == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("reason", "TemplateNotFound"),
		)
	}

	// 2. 验证模板启用状态
	if !template.IsEnabledTemplate() {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("reason", "TemplateDisabled"),
		)
	}

	// 3. 设置默认选项
	if req.Options == nil {
		req.Options = &ApplyOptions{
			IncludeDepartments: true,
			IncludeRoles:       true,
			IncludePermissions: true,
			IncludeWorkflows:   true,
		}
	}

	// 4. 在事务中应用模板
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 应用部门配置
		if req.Options.IncludeDepartments && template.Config.Departments != nil {
			for _, deptConfig := range template.Config.Departments {
				dept := &entity.Department{
					DeptID:      generateUUID(),
					TenantID:    req.TenantID,
					OrgID:       req.OrgID,
					DeptName:    deptConfig.DeptName,
					DeptCode:    deptConfig.DeptCode,
					ParentID:    deptConfig.ParentID,
					Description: deptConfig.Description,
					LeaderID:    deptConfig.LeaderID,
					SortOrder:   deptConfig.SortOrder,
					Status:      entity.OrgStatusActive,
					Level:       1,
					Path:        "/",
					CreatedAt:   time.Now().UnixMilli(),
					UpdatedAt:   time.Now().UnixMilli(),
				}

				if err := s.departmentRepo.Create(ctx, dept); err != nil {
					return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
						errorx.KV("operation", "create department from template"),
					)
				}
			}
		}

		// 应用角色配置（TODO: 需要角色仓储）
		if req.Options.IncludeRoles && template.Config.Roles != nil {
			// TODO: 实现角色应用
		}

		// 应用权限配置（TODO: 需要权限仓储）
		if req.Options.IncludePermissions && template.Config.Permissions != nil {
			// TODO: 实现权限应用
		}

		// 应用工作流配置（TODO: 需要工作流仓储）
		if req.Options.IncludeWorkflows && template.Config.Workflows != nil {
			// TODO: 实现工作流应用
		}

		return nil
	})
}

// GetTemplate 获取模板
func (s *OrgTemplateService) GetTemplate(
	ctx context.Context,
	templateID string,
) (*entity.OrgTemplate, error) {
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get template"),
		)
	}

	return template, nil
}

// ListTemplatesRequest 查询模板列表请求
type ListTemplatesRequest struct {
	TenantID  string  `json:"tenant_id,omitempty"`
	IsBuiltin *bool   `json:"is_builtin,omitempty"`
	IsEnabled *bool   `json:"is_enabled,omitempty"`
	Keyword   string  `json:"keyword,omitempty"`
	Page      int     `json:"page" binding:"min=1"`
	PageSize  int     `json:"page_size" binding:"min=1,max=100"`
}

// ListTemplatesResponse 查询模板列表响应
type ListTemplatesResponse struct {
	Templates []*entity.OrgTemplate `json:"templates"`
	Total     int64                `json:"total"`
	Page      int                  `json:"page"`
	PageSize  int                  `json:"page_size"`
}

// ListTemplates 查询模板列表
func (s *OrgTemplateService) ListTemplates(
	ctx context.Context,
	req *ListTemplatesRequest,
) (*ListTemplatesResponse, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 构建查询过滤器
	filter := &TemplateFilter{
		TenantID:  req.TenantID,
		IsBuiltin: req.IsBuiltin,
		IsEnabled: req.IsEnabled,
		Keyword:   req.Keyword,
		PageToken: "",
		PageSize:  req.PageSize,
	}

	// 查询模板
	templates, total, err := s.templateRepo.List(ctx, filter)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "list templates"),
		)
	}

	return &ListTemplatesResponse{
		Templates: templates,
		Total:     total,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}, nil
}

// GetBuiltinTemplates 获取预置模板
func (s *OrgTemplateService) GetBuiltinTemplates(
	ctx context.Context,
) ([]*entity.OrgTemplate, error) {
	templates, err := s.templateRepo.GetBuiltinTemplates(ctx)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get builtin templates"),
		)
	}

	return templates, nil
}

// UpdateTemplateRequest 更新模板请求
type UpdateTemplateRequest struct {
	TemplateID   string                 `json:"template_id" binding:"required"`
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Config       *entity.TemplateConfig `json:"config,omitempty"`
	Version      string                 `json:"version,omitempty"`
	IsEnabled    *bool                  `json:"is_enabled,omitempty"`
	UpdatedBy    string                 `json:"updated_by" binding:"required"`
}

// UpdateTemplate 更新模板
func (s *OrgTemplateService) UpdateTemplate(
	ctx context.Context,
	req *UpdateTemplateRequest,
) (*entity.OrgTemplate, error) {
	// 1. 获取模板
	template, err := s.templateRepo.GetByID(ctx, req.TemplateID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get template"),
		)
	}
	if template == nil {
		return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("reason", "TemplateNotFound"),
		)
	}

	// 2. 验证非预置模板
	if template.IsBuiltin {
		return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("reason", "CannotUpdateBuiltinTemplate"),
		)
	}

	// 3. 更新字段
	now := time.Now().UnixMilli()
	if req.Name != "" {
		template.Name = req.Name
	}
	if req.Description != "" {
		template.Description = req.Description
	}
	if req.Config != nil {
		template.Config = req.Config
	}
	if req.Version != "" {
		template.Version = req.Version
	}
	if req.IsEnabled != nil {
		template.IsEnabled = *req.IsEnabled
	}
	template.UpdatedBy = req.UpdatedBy
	template.UpdatedAt = now

	// 4. 保存更新
	if err := s.templateRepo.Update(ctx, template); err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "update template"),
		)
	}

	return template, nil
}

// DeleteTemplate 删除模板
func (s *OrgTemplateService) DeleteTemplate(
	ctx context.Context,
	templateID string,
) error {
	// 1. 获取模板
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get template"),
		)
	}
	if template == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("reason", "TemplateNotFound"),
		)
	}

	// 2. 验证非预置模板
	if template.IsBuiltin {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("reason", "CannotDeleteBuiltinTemplate"),
		)
	}

	// 3. 删除模板
	if err := s.templateRepo.Delete(ctx, templateID); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "delete template"),
		)
	}

	return nil
}

// CreateBuiltinTemplates 创建预置模板（初始化时调用）
func (s *OrgTemplateService) CreateBuiltinTemplates(ctx context.Context) error {
	// 预置模板1：科技公司模板
	techCompanyTemplate := &entity.OrgTemplate{
		TemplateID:   "builtin_tech_company_v1",
		TenantID:     "system",
		Name:         "科技公司标准模板",
		Description:  "适用于中小型科技公司的标准组织架构模板",
		Config: &entity.TemplateConfig{
			Departments: []*entity.DepartmentConfig{
				{DeptName: "研发部", DeptCode: "RD", SortOrder: 1},
				{DeptName: "产品部", DeptCode: "PD", SortOrder: 2},
				{DeptName: "市场部", DeptCode: "MKT", SortOrder: 3},
				{DeptName: "销售部", DeptCode: "SD", SortOrder: 4},
				{DeptName: "人力资源部", DeptCode: "HR", SortOrder: 5},
				{DeptName: "财务部", DeptCode: "FIN", SortOrder: 6},
			},
			Roles: []*entity.RoleConfig{
				{
					RoleName:     "系统管理员",
					RoleCode:     "admin",
					IsSystemRole: true,
				},
				{
					RoleName: "部门经理",
					RoleCode: "manager",
				},
				{
					RoleName: "普通员工",
					RoleCode: "employee",
				},
			},
		},
		IsBuiltin: true,
		Version:   "1.0",
		IsEnabled: true,
		CreatedBy: "system",
		UpdatedBy: "system",
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}

	// 检查是否已存在
	if _, err := s.templateRepo.GetByID(ctx, techCompanyTemplate.TemplateID); err == nil {
		return nil // 已存在，跳过
	}

	// 创建预置模板
	if err := s.templateRepo.Create(ctx, techCompanyTemplate); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "create builtin template"),
		)
	}

	return nil
}
