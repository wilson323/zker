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

	"github.com/coze-dev/coze-studio/backend/domain/developer/entity"
	"github.com/coze-dev/coze-studio/backend/domain/developer/repository"
)

// ProjectManagementService 项目管理服务接口
type ProjectManagementService interface {
	// CreateProject 创建项目
	CreateProject(ctx context.Context, req *CreateProjectRequest) (*entity.Project, error)

	// GetProject 获取项目详情
	GetProject(ctx context.Context, projectID string) (*entity.Project, error)

	// ListProjects 查询项目列表
	ListProjects(ctx context.Context, filter *ProjectListFilter) ([]*entity.Project, int64, error)

	// UpdateProject 更新项目信息
	UpdateProject(ctx context.Context, req *UpdateProjectRequest) error

	// DeleteProject 删除项目
	DeleteProject(ctx context.Context, projectID string) error

	// ArchiveProject 归档项目
	ArchiveProject(ctx context.Context, projectID string) error

	// PublishProject 发布项目到生产环境
	PublishProject(ctx context.Context, projectID string) error

	// GetDeveloperProjects 获取开发者的项目列表
	GetDeveloperProjects(ctx context.Context, developerID string) ([]*entity.Project, error)

	// GetTenantProjects 获取租户的项目列表
	GetTenantProjects(ctx context.Context, tenantID string) ([]*entity.Project, error)

	// UpdateProjectConfig 更新项目配置
	UpdateProjectConfig(ctx context.Context, projectID string, config *entity.ProjectConfig) error

	// GetProjectStats 获取项目统计信息
	GetProjectStats(ctx context.Context, projectID string) (*ProjectStats, error)
}

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	TenantID     string                  `json:"tenant_id" validate:"required"`
	DeveloperID  string                  `json:"developer_id" validate:"required"`
	ProjectName  string                  `json:"project_name" validate:"required,max=100"`
	ProjectType  entity.ProjectType      `json:"project_type" validate:"required"`
	Description  string                  `json:"description"`
	Config       *entity.ProjectConfig   `json:"config"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	ProjectID   string                `json:"project_id" validate:"required"`
	ProjectName *string               `json:"project_name"`
	Description *string               `json:"description"`
	Config      *entity.ProjectConfig `json:"config"`
}

// ProjectListFilter 项目列表过滤器
type ProjectListFilter struct {
	TenantID        string                 `validate:"required"`
	DeveloperID     string
	ProjectType     entity.ProjectType
	ProjectStatus   entity.ProjectStatus
	PageToken       string
	PageSize        int
}

// ProjectStats 项目统计信息
type ProjectStats struct {
	ProjectID       string `json:"project_id"`
	TotalAPIKeys    int64  `json:"total_api_keys"`
	ActiveAPIKeys   int64  `json:"active_api_keys"`
	TotalWebhooks   int64  `json:"total_webhooks"`
	ActiveWebhooks  int64  `json:"active_webhooks"`
	TotalSDKs       int64  `json:"total_sdks"`
	PublishedSDKs   int64  `json:"published_sdks"`
}

// projectManagementService 项目管理服务实现
type projectManagementService struct {
	projectRepo repository.ProjectRepository
	apiKeyRepo  repository.APIKeyRepository
	webhookRepo repository.WebhookRepository
	sdkRepo     repository.SDKRepository
}

// NewProjectManagementService 创建项目管理服务实例
func NewProjectManagementService(
	projectRepo repository.ProjectRepository,
	apiKeyRepo repository.APIKeyRepository,
	webhookRepo repository.WebhookRepository,
	sdkRepo repository.SDKRepository,
) ProjectManagementService {
	return &projectManagementService{
		projectRepo: projectRepo,
		apiKeyRepo:  apiKeyRepo,
		webhookRepo: webhookRepo,
		sdkRepo:     sdkRepo,
	}
}

// CreateProject 创建项目
func (s *projectManagementService) CreateProject(ctx context.Context, req *CreateProjectRequest) (*entity.Project, error) {
	// TODO: 验证开发者是否存在
	// TODO: 验证租户配额是否充足
	// TODO: 验证项目名称唯一性

	project := &entity.Project{
		ProjectID:    generateID("proj"),
		TenantID:     req.TenantID,
		DeveloperID:  req.DeveloperID,
		ProjectName:  req.ProjectName,
		ProjectType:  req.ProjectType,
		Description:  req.Description,
		Status:       entity.ProjectStatusDevelopment,
	}

	// 序列化配置
	if req.Config != nil {
		// TODO: 序列化配置为JSON
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// GetProject 获取项目详情
func (s *projectManagementService) GetProject(ctx context.Context, projectID string) (*entity.Project, error) {
	return s.projectRepo.GetByID(ctx, projectID)
}

// ListProjects 查询项目列表
func (s *projectManagementService) ListProjects(ctx context.Context, filter *ProjectListFilter) ([]*entity.Project, int64, error) {
	repoFilter := &repository.ProjectFilter{
		TenantID:      filter.TenantID,
		DeveloperID:   filter.DeveloperID,
		ProjectType:   filter.ProjectType,
		ProjectStatus: filter.ProjectStatus,
		PageToken:     filter.PageToken,
		PageSize:      filter.PageSize,
	}

	return s.projectRepo.List(ctx, repoFilter)
}

// UpdateProject 更新项目信息
func (s *projectManagementService) UpdateProject(ctx context.Context, req *UpdateProjectRequest) error {
	project, err := s.projectRepo.GetByID(ctx, req.ProjectID)
	if err != nil {
		return err
	}
	if project == nil {
		// TODO: 返回项目不存在错误
		return err
	}

	// 更新字段
	if req.ProjectName != nil {
		project.ProjectName = *req.ProjectName
	}
	if req.Description != nil {
		project.Description = *req.Description
	}
	if req.Config != nil {
		// TODO: 序列化配置
	}

	return s.projectRepo.Update(ctx, project)
}

// DeleteProject 删除项目
func (s *projectManagementService) DeleteProject(ctx context.Context, projectID string) error {
	// TODO: 检查项目是否有关联资源（API Keys、Webhooks等）
	// TODO: 级联删除或拒绝删除
	return s.projectRepo.Delete(ctx, projectID)
}

// ArchiveProject 归档项目
func (s *projectManagementService) ArchiveProject(ctx context.Context, projectID string) error {
	return s.projectRepo.UpdateStatus(ctx, projectID, entity.ProjectStatusArchived)
}

// PublishProject 发布项目到生产环境
func (s *projectManagementService) PublishProject(ctx context.Context, projectID string) error {
	return s.projectRepo.UpdateStatus(ctx, projectID, entity.ProjectStatusProduction)
}

// GetDeveloperProjects 获取开发者的项目列表
func (s *projectManagementService) GetDeveloperProjects(ctx context.Context, developerID string) ([]*entity.Project, error) {
	return s.projectRepo.GetByDeveloperID(ctx, developerID)
}

// GetTenantProjects 获取租户的项目列表
func (s *projectManagementService) GetTenantProjects(ctx context.Context, tenantID string) ([]*entity.Project, error) {
	return s.projectRepo.GetByTenantID(ctx, tenantID)
}

// UpdateProjectConfig 更新项目配置
func (s *projectManagementService) UpdateProjectConfig(ctx context.Context, projectID string, config *entity.ProjectConfig) error {
	// TODO: 序列化配置
	return s.projectRepo.UpdateConfig(ctx, projectID, "")
}

// GetProjectStats 获取项目统计信息
func (s *projectManagementService) GetProjectStats(ctx context.Context, projectID string) (*ProjectStats, error) {
	stats := &ProjectStats{ProjectID: projectID}

	// 统计API Keys
	totalKeys, _ := s.apiKeyRepo.CountByProjectID(ctx, projectID)
	stats.TotalAPIKeys = totalKeys

	// 统计Webhooks
	totalWebhooks, _ := s.webhookRepo.CountByProjectID(ctx, projectID)
	stats.TotalWebhooks = totalWebhooks

	// TODO: 统计活跃的API Keys和Webhooks
	// TODO: 统计SDKs

	return stats, nil
}

// generateID 生成唯一ID（简化实现）
func generateID(prefix string) string {
	// TODO: 使用UUID生成器
	return prefix + "_" + randomString(32)
}

// randomString 生成随机字符串
func randomString(length int) string {
	// TODO: 实现随机字符串生成
	return "random"
}
