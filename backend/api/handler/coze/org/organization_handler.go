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

package org

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/api/handler/coze"
	"github.com/coze-dev/coze-studio/backend/api/model/org"
	orgentity "github.com/coze-dev/coze-studio/backend/domain/org/entity"
	orgservice "github.com/coze-dev/coze-studio/backend/domain/org/service"
)

// OrganizationHandler 组织管理处理器
// 职责：处理HTTP请求，调用领域服务，返回HTTP响应
// 遵循SOLID原则：单一职责（仅处理HTTP层），依赖倒置（依赖service接口）
type OrganizationHandler struct {
	orgService *orgservice.OrganizationService
}

// NewOrganizationHandler 创建组织处理器实例
func NewOrganizationHandler(orgService *orgservice.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: orgService,
	}
}

// ==================== 基础CRUD操作 ====================

// CreateOrganization 创建组织
// @router /api/organizations [POST]
func (h *OrganizationHandler) CreateOrganization(ctx context.Context, c *app.RequestContext) {
	var req org.CreateOrganizationRequest
	if err := c.BindAndValidate(&req); err != nil {
		coze.InvalidParamRequestResponse(c, err.Error())
		return
	}

	// 转换为领域服务请求
	serviceReq := &orgservice.CreateOrganizationRequest{
		TenantID:     req.TenantID,
		OrgName:      req.OrgName,
		OrgType:      orgentity.OrganizationType(req.OrgType),
		ParentID:     req.ParentID,
		OrgCode:      req.OrgCode,
		LeaderID:     req.LeaderID,
		Description:  req.Description,
		SortOrder:    req.SortOrder,
	}

	// 调用领域服务
	orgEntity, err := h.orgService.CreateOrganization(ctx, serviceReq)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换响应
	c.JSON(http.StatusOK, &org.CreateOrganizationResponse{
		Code:    0,
		Message: "success",
		Data:    toOrganizationData(orgEntity),
	})
}

// GetOrganization 获取组织详情
// @router /api/organizations/:id [GET]
func (h *OrganizationHandler) GetOrganization(ctx context.Context, c *app.RequestContext) {
	orgID := c.Param("id")
	if orgID == "" {
		coze.InvalidParamRequestResponse(c, "organization id is required")
		return
	}

	// 调用领域服务
	orgEntity, err := h.orgService.GetOrganization(ctx, orgID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换响应
	c.JSON(http.StatusOK, &org.GetOrganizationResponse{
		Code:    0,
		Message: "success",
		Data:    toOrganizationData(orgEntity),
	})
}

// UpdateOrganization 更新组织信息
// @router /api/organizations/:id [PUT]
func (h *OrganizationHandler) UpdateOrganization(ctx context.Context, c *app.RequestContext) {
	orgID := c.Param("id")
	if orgID == "" {
		coze.InvalidParamRequestResponse(c, "organization id is required")
		return
	}

	var req org.UpdateOrganizationRequest
	if err := c.BindAndValidate(&req); err != nil {
		coze.InvalidParamRequestResponse(c, err.Error())
		return
	}

	// 转换为领域服务请求
	serviceReq := &orgservice.UpdateOrganizationRequest{
		OrgID:       orgID,
		OrgName:     req.OrgName,
		LeaderID:    req.LeaderID,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      orgentity.OrganizationStatus(req.Status),
	}

	// 调用领域服务
	orgEntity, err := h.orgService.UpdateOrganization(ctx, serviceReq)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换响应
	c.JSON(http.StatusOK, &org.UpdateOrganizationResponse{
		Code:    0,
		Message: "success",
		Data:    toOrganizationData(orgEntity),
	})
}

// DeleteOrganization 删除组织
// @router /api/organizations/:id [DELETE]
func (h *OrganizationHandler) DeleteOrganization(ctx context.Context, c *app.RequestContext) {
	orgID := c.Param("id")
	if orgID == "" {
		coze.InvalidParamRequestResponse(c, "organization id is required")
		return
	}

	// 调用领域服务
	if err := h.orgService.DeleteOrganization(ctx, orgID); err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, &org.DeleteOrganizationResponse{
		Code:    0,
		Message: "organization deleted successfully",
	})
}

// ==================== 组织树操作 ====================

// GetOrganizationTree 获取组织树
// @router /api/organizations/tree [GET]
func (h *OrganizationHandler) GetOrganizationTree(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		coze.InvalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 调用领域服务
	orgTree, err := h.orgService.GetOrganizationTree(ctx, tenantID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换响应
	data := make([]*org.OrganizationData, 0, len(orgTree))
	for _, orgEntity := range orgTree {
		data = append(data, toOrganizationData(orgEntity))
	}

	c.JSON(http.StatusOK, &org.GetOrganizationTreeResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// ListOrganizations 分页查询组织列表
// @router /api/organizations [GET]
func (h *OrganizationHandler) ListOrganizations(ctx context.Context, c *app.RequestContext) {
	// 解析查询参数
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		coze.InvalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	orgType := c.Query("org_type")
	status := c.Query("status")
	parentID := c.Query("parent_id")
	keyword := c.Query("keyword")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 构建过滤条件
	filter := &orgservice.OrganizationFilter{
		TenantID:  tenantID,
		OrgType:   orgType,
		Status:    status,
		ParentID:  parentID,
		Keyword:   keyword,
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	// 调用领域服务
	orgs, total, err := h.orgService.ListOrganizations(ctx, filter)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换响应
	data := make([]*org.OrganizationData, 0, len(orgs))
	for _, orgEntity := range orgs {
		data = append(data, toOrganizationData(orgEntity))
	}

	c.JSON(http.StatusOK, &org.ListOrganizationsResponse{
		Code:     0,
		Message:  "success",
		Data:     data,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// MoveOrganization 移动组织
// @router /api/organizations/:id/move [POST]
func (h *OrganizationHandler) MoveOrganization(ctx context.Context, c *app.RequestContext) {
	orgID := c.Param("id")
	if orgID == "" {
		coze.InvalidParamRequestResponse(c, "organization id is required")
		return
	}

	var req org.MoveOrganizationRequest
	if err := c.BindAndValidate(&req); err != nil {
		coze.InvalidParamRequestResponse(c, err.Error())
		return
	}

	// 转换为领域服务请求
	serviceReq := &orgservice.MoveOrganizationRequest{
		OrgID:       orgID,
		NewParentID: req.NewParentID,
	}

	// 调用领域服务
	if err := h.orgService.MoveOrganization(ctx, serviceReq); err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, &org.MoveOrganizationResponse{
		Code:    0,
		Message: "organization moved successfully",
	})
}

// ==================== 组织关系查询 ====================

// GetChildren 获取子组织列表
// @router /api/organizations/:id/children [GET]
func (h *OrganizationHandler) GetChildren(ctx context.Context, c *app.RequestContext) {
	parentID := c.Param("id")
	if parentID == "" {
		coze.InvalidParamRequestResponse(c, "organization id is required")
		return
	}

	// 调用领域服务
	children, err := h.orgService.GetChildren(ctx, parentID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换响应
	data := make([]*org.OrganizationData, 0, len(children))
	for _, orgEntity := range children {
		data = append(data, toOrganizationData(orgEntity))
	}

	c.JSON(http.StatusOK, &org.GetChildrenResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// GetAncestors 获取组织的所有祖先
// @router /api/organizations/:id/ancestors [GET]
func (h *OrganizationHandler) GetAncestors(ctx context.Context, c *app.RequestContext) {
	orgID := c.Param("id")
	if orgID == "" {
		coze.InvalidParamRequestResponse(c, "organization id is required")
		return
	}

	// 调用领域服务
	ancestors, err := h.orgService.GetAncestors(ctx, orgID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换响应
	data := make([]*org.OrganizationData, 0, len(ancestors))
	for _, orgEntity := range ancestors {
		data = append(data, toOrganizationData(orgEntity))
	}

	c.JSON(http.StatusOK, &org.GetAncestorsResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// GetDescendants 获取组织的所有后代
// @router /api/organizations/:id/descendants [GET]
func (h *OrganizationHandler) GetDescendants(ctx context.Context, c *app.RequestContext) {
	orgID := c.Param("id")
	if orgID == "" {
		coze.InvalidParamRequestResponse(c, "organization id is required")
		return
	}

	// 调用领域服务
	descendants, err := h.orgService.GetDescendants(ctx, orgID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换响应
	data := make([]*org.OrganizationData, 0, len(descendants))
	for _, orgEntity := range descendants {
		data = append(data, toOrganizationData(orgEntity))
	}

	c.JSON(http.StatusOK, &org.GetDescendantsResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// ==================== 辅助函数 ====================

// toOrganizationData 将领域实体转换为API数据模型
func toOrganizationData(entity *orgentity.Organization) *org.OrganizationData {
	if entity == nil {
		return nil
	}

	return &org.OrganizationData{
		OrgID:       entity.OrgID,
		TenantID:    entity.TenantID,
		OrgName:     entity.OrgName,
		OrgType:     string(entity.OrgType),
		ParentID:    entity.ParentID,
		OrgCode:     entity.OrgCode,
		Level:       entity.Level,
		Path:        entity.Path,
		SortOrder:   entity.SortOrder,
		Status:      string(entity.Status),
		Description: entity.Description,
		LeaderID:    entity.LeaderID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

// toOrganizationDataList 批量转换领域实体为API数据模型
func toOrganizationDataList(entities []*orgentity.Organization) []*org.OrganizationData {
	if entities == nil {
		return nil
	}

	data := make([]*org.OrganizationData, 0, len(entities))
	for _, entity := range entities {
		data = append(data, toOrganizationData(entity))
	}
	return data
}
