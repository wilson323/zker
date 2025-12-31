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

package developer

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	devservice "github.com/coze-dev/coze-studio/backend/domain/developer/service"
	"github.com/coze-dev/coze-studio/backend/domain/developer/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	berrno "github.com/coze-dev/coze-studio/backend/types/errno"
)

// ProjectManagementHandler 项目管理Handler
type ProjectManagementHandler struct {
	projectService devservice.ProjectManagementService
}

// NewProjectManagementHandler 创建项目管理Handler实例
func NewProjectManagementHandler(projectService devservice.ProjectManagementService) *ProjectManagementHandler {
	return &ProjectManagementHandler{
		projectService: projectService,
	}
}

// ListProjectsRequest 列出项目请求
type ListProjectsRequest struct {
	ProjectType   string `json:"project_type" query:"project_type"`
	ProjectStatus string `json:"project_status" query:"project_status"`
	DeveloperID   string `json:"developer_id" query:"developer_id"`
	PageToken     string `json:"page_token" query:"page_token"`
	PageSize      int    `json:"page_size" query:"page_size"`
}

// ListProjectsResponse 列出项目响应
type ListProjectsResponse struct {
	Projects  []*ProjectListItem `json:"projects"`
	Total     int64              `json:"total"`
	PageToken string             `json:"page_token,omitempty"`
}

// ProjectListItem 项目列表项
type ProjectListItem struct {
	ProjectID    string             `json:"project_id"`
	ProjectName  string             `json:"project_name"`
	ProjectType  entity.ProjectType `json:"project_type"`
	Description  string             `json:"description"`
	Status       entity.ProjectStatus `json:"status"`
	CreatedAt    int64              `json:"created_at"`
	UpdatedAt    int64              `json:"updated_at"`
}

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	ProjectName  string                `json:"project_name" validate:"required,max=100"`
	ProjectType  entity.ProjectType    `json:"project_type" validate:"required"`
	Description  string                `json:"description"`
	DeveloperID  string                `json:"developer_id" validate:"required"`
	Config       *entity.ProjectConfig `json:"config"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	ProjectName  *string               `json:"project_name" validate:"omitempty,max=100"`
	Description  *string               `json:"description"`
	Config       *entity.ProjectConfig `json:"config"`
}

// PublishProjectRequest 发布项目请求
type PublishProjectRequest struct {
	ProjectID string `json:"project_id" validate:"required"`
}

// ArchiveProjectRequest 归档项目请求
type ArchiveProjectRequest struct {
	ProjectID string `json:"project_id" validate:"required"`
}

// ListProjects 列出项目
// @router /api/developer/projects [GET]
func (h *ProjectManagementHandler) ListProjects(ctx context.Context, c *app.RequestContext) {
	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 解析请求参数
	var req ListProjectsRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 设置默认分页参数
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 构建过滤器
	filter := &devservice.ProjectListFilter{
		TenantID:      tenantID,
		DeveloperID:   req.DeveloperID,
		PageToken:     req.PageToken,
		PageSize:      req.PageSize,
	}

	if req.ProjectType != "" {
		filter.ProjectType = entity.ProjectType(req.ProjectType)
	}
	if req.ProjectStatus != "" {
		filter.ProjectStatus = entity.ProjectStatus(req.ProjectStatus)
	}

	// 查询项目列表
	projects, total, err := h.projectService.ListProjects(ctx, filter)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	items := make([]*ProjectListItem, 0, len(projects))
	for _, p := range projects {
		items = append(items, &ProjectListItem{
			ProjectID:    p.ProjectID,
			ProjectName:  p.ProjectName,
			ProjectType:  p.ProjectType,
			Description:  p.Description,
			Status:       p.Status,
			CreatedAt:    p.CreatedAt,
			UpdatedAt:    p.UpdatedAt,
		})
	}

	httputil.BuildSuccessResp(c, &ListProjectsResponse{
		Projects:  items,
		Total:     total,
		PageToken: req.PageToken,
	})
}

// CreateProject 创建项目
// @router /api/developer/projects [POST]
func (h *ProjectManagementHandler) CreateProject(ctx context.Context, c *app.RequestContext) {
	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 解析请求参数
	var req CreateProjectRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 构建创建请求
	createReq := &devservice.CreateProjectRequest{
		TenantID:    tenantID,
		DeveloperID: req.DeveloperID,
		ProjectName: req.ProjectName,
		ProjectType: req.ProjectType,
		Description: req.Description,
		Config:      req.Config,
	}

	// 创建项目
	project, err := h.projectService.CreateProject(ctx, createReq)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, project)
}

// GetProject 获取项目详情
// @router /api/developer/projects/:id [GET]
func (h *ProjectManagementHandler) GetProject(ctx context.Context, c *app.RequestContext) {
	projectID := c.Param("id")
	if projectID == "" {
		httputil.BadRequest(c, "project_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 查询项目
	project, err := h.projectService.GetProject(ctx, projectID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	// 验证租户隔离
	if project == nil || project.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrProjectNotFound.Int32Code(), berrno.ErrProjectNotFound.Message(), berrno.ErrProjectNotFound.MessageZH(), nil)
		return
	}

	httputil.BuildSuccessResp(c, project)
}

// UpdateProject 更新项目
// @router /api/developer/projects/:id [PUT]
func (h *ProjectManagementHandler) UpdateProject(ctx context.Context, c *app.RequestContext) {
	projectID := c.Param("id")
	if projectID == "" {
		httputil.BadRequest(c, "project_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 解析请求参数
	var req UpdateProjectRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 验证项目存在且属于该租户
	project, err := h.projectService.GetProject(ctx, projectID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if project == nil || project.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrProjectNotFound.Int32Code(), berrno.ErrProjectNotFound.Message(), berrno.ErrProjectNotFound.MessageZH(), nil)
		return
	}

	// 检查项目状态（生产环境不可修改）
	if project.IsProduction() {
		httputil.BuildErrorResp(c, berrno.ErrProjectInProduction.Int32Code(), berrno.ErrProjectInProduction.Message(), berrno.ErrProjectInProduction.MessageZH(), nil)
		return
	}

	// 构建更新请求
	updateReq := &devservice.UpdateProjectRequest{
		ProjectID:   projectID,
		ProjectName: req.ProjectName,
		Description: req.Description,
		Config:      req.Config,
	}

	// 更新项目
	if err := h.projectService.UpdateProject(ctx, updateReq); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"project_id": projectID,
		"message":    "Project updated successfully",
	})
}

// DeleteProject 删除项目
// @router /api/developer/projects/:id [DELETE]
func (h *ProjectManagementHandler) DeleteProject(ctx context.Context, c *app.RequestContext) {
	projectID := c.Param("id")
	if projectID == "" {
		httputil.BadRequest(c, "project_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 验证项目存在且属于该租户
	project, err := h.projectService.GetProject(ctx, projectID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if project == nil || project.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrProjectNotFound.Int32Code(), berrno.ErrProjectNotFound.Message(), berrno.ErrProjectNotFound.MessageZH(), nil)
		return
	}

	// 检查项目状态（生产环境不可删除）
	if project.IsProduction() {
		httputil.BuildErrorResp(c, berrno.ErrProjectInProduction.Int32Code(), berrno.ErrProjectInProduction.Message(), berrno.ErrProjectInProduction.MessageZH(), nil)
		return
	}

	// 删除项目
	if err := h.projectService.DeleteProject(ctx, projectID); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Project deleted successfully",
	})
}

// PublishProject 发布项目到生产环境
// @router /api/developer/projects/:id/publish [POST]
func (h *ProjectManagementHandler) PublishProject(ctx context.Context, c *app.RequestContext) {
	projectID := c.Param("id")
	if projectID == "" {
		httputil.BadRequest(c, "project_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 验证项目存在且属于该租户
	project, err := h.projectService.GetProject(ctx, projectID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if project == nil || project.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrProjectNotFound.Int32Code(), berrno.ErrProjectNotFound.Message(), berrno.ErrProjectNotFound.MessageZH(), nil)
		return
	}

	// 发布项目
	if err := h.projectService.PublishProject(ctx, projectID); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"project_id": projectID,
		"status":     entity.ProjectStatusProduction,
		"message":    "Project published successfully",
	})
}

// ArchiveProject 归档项目
// @router /api/developer/projects/:id/archive [POST]
func (h *ProjectManagementHandler) ArchiveProject(ctx context.Context, c *app.RequestContext) {
	projectID := c.Param("id")
	if projectID == "" {
		httputil.BadRequest(c, "project_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 验证项目存在且属于该租户
	project, err := h.projectService.GetProject(ctx, projectID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if project == nil || project.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrProjectNotFound.Int32Code(), berrno.ErrProjectNotFound.Message(), berrno.ErrProjectNotFound.MessageZH(), nil)
		return
	}

	// 检查项目状态（生产环境不可归档）
	if project.IsProduction() {
		httputil.BuildErrorResp(c, berrno.ErrProjectInProduction.Int32Code(), berrno.ErrProjectInProduction.Message(), berrno.ErrProjectInProduction.MessageZH(), nil)
		return
	}

	// 归档项目
	if err := h.projectService.ArchiveProject(ctx, projectID); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"project_id": projectID,
		"status":     entity.ProjectStatusArchived,
		"message":    "Project archived successfully",
	})
}

// GetProjectStats 获取项目统计信息
// @router /api/developer/projects/:id/stats [GET]
func (h *ProjectManagementHandler) GetProjectStats(ctx context.Context, c *app.RequestContext) {
	projectID := c.Param("id")
	if projectID == "" {
		httputil.BadRequest(c, "project_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 验证项目存在且属于该租户
	project, err := h.projectService.GetProject(ctx, projectID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if project == nil || project.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrProjectNotFound.Int32Code(), berrno.ErrProjectNotFound.Message(), berrno.ErrProjectNotFound.MessageZH(), nil)
		return
	}

	// 获取统计信息
	stats, err := h.projectService.GetProjectStats(ctx, projectID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, stats)
}

// GetDeveloperProjects 获取开发者的项目列表
// @router /api/developer/developers/:id/projects [GET]
func (h *ProjectManagementHandler) GetDeveloperProjects(ctx context.Context, c *app.RequestContext) {
	developerID := c.Param("id")
	if developerID == "" {
		httputil.BadRequest(c, "developer_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 查询项目列表
	projects, err := h.projectService.GetDeveloperProjects(ctx, developerID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	// 过滤租户（确保租户隔离）
	filteredProjects := make([]*entity.Project, 0)
	for _, p := range projects {
		if p.TenantID == tenantID {
			filteredProjects = append(filteredProjects, p)
		}
	}

	httputil.BuildSuccessResp(c, filteredProjects)
}

// GetTenantProjects 获取租户的项目列表（内部接口）
// @router /api/developer/tenant/projects [GET]
func (h *ProjectManagementHandler) GetTenantProjects(ctx context.Context, c *app.RequestContext) {
	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 查询项目列表
	projects, err := h.projectService.GetTenantProjects(ctx, tenantID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, projects)
}
