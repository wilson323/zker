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

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/domain/org/service"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// DepartmentHandler 部门管理处理器
// 职责：处理部门相关的HTTP请求，进行参数验证，调用领域服务
type DepartmentHandler struct {
	deptService *service.DepartmentService
}

// NewDepartmentHandler 创建部门处理器
func NewDepartmentHandler(deptService *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		deptService: deptService,
	}
}

// CreateDepartmentRequest 创建部门请求
type CreateDepartmentRequest struct {
	TenantID     string  `json:"tenant_id" binding:"required"`
	OrgID        string  `json:"org_id" binding:"required"`
	ParentID     *string `json:"parent_id,omitempty"`
	DeptName     string  `json:"dept_name" binding:"required,min=1,max=200"`
	DeptCode     string  `json:"dept_code" binding:"required,min=1,max=50"`
	LeaderID     *string `json:"leader_id,omitempty"`
	ParentLeader *string `json:"parent_leader,omitempty"`
	Description  string  `json:"description,omitempty"`
	SortOrder    int     `json:"sort_order,omitempty"`
}

// UpdateDepartmentRequest 更新部门请求
type UpdateDepartmentRequest struct {
	DeptID       string  `json:"dept_id" binding:"required"`
	DeptName     string  `json:"dept_name" binding:"omitempty,min=1,max=200"`
	LeaderID     *string `json:"leader_id,omitempty"`
	ParentLeader *string `json:"parent_leader,omitempty"`
	Description  string  `json:"description,omitempty"`
	SortOrder    int     `json:"sort_order,omitempty"`
	Status       string  `json:"status" binding:"omitempty,oneof=active inactive frozen"`
}

// MoveDepartmentRequest 移动部门请求
type MoveDepartmentRequest struct {
	DeptID      string  `json:"dept_id" binding:"required"`
	NewParentID *string `json:"new_parent_id,omitempty"`
}

// ListDepartmentsRequest 查询部门列表请求
type ListDepartmentsRequest struct {
	TenantID string `query:"tenant_id" binding:"required"`
	OrgID    string `query:"org_id,omitempty"`
	Status   string `query:"status" binding:"omitempty,oneof=active inactive frozen"`
	ParentID string `query:"parent_id,omitempty"`
	Level    int    `query:"level,omitempty"`
	Keyword  string `query:"keyword,omitempty"`
	Page     int    `query:"page" binding:"omitempty,min=1"`
	PageSize int    `query:"page_size" binding:"omitempty,min=1,max=100"`
}

// CreateDepartment 创建部门
// @router /api/org/departments [POST]
func (h *DepartmentHandler) CreateDepartment(ctx context.Context, c *app.RequestContext) {
	var req CreateDepartmentRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, Fail(errno.InvalidParams.Code(), err.Error()))
		return
	}

	// 转换为服务层请求
	svcReq := &service.CreateDepartmentRequest{
		TenantID:     req.TenantID,
		OrgID:        req.OrgID,
		ParentID:     req.ParentID,
		DeptName:     req.DeptName,
		DeptCode:     req.DeptCode,
		LeaderID:     req.LeaderID,
		ParentLeader: req.ParentLeader,
		Description:  req.Description,
		SortOrder:    req.SortOrder,
	}

	dept, err := h.deptService.CreateDepartment(ctx, svcReq)
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, Success(dept))
}

// GetDepartment 获取部门详情
// @router /api/org/departments/:id [GET]
func (h *DepartmentHandler) GetDepartment(ctx context.Context, c *app.RequestContext) {
	deptID := c.Param("id")
	if deptID == "" {
		c.JSON(http.StatusBadRequest, Fail(errno.InvalidParams.Code(), "department id is required"))
		return
	}

	dept, err := h.deptService.GetDepartment(ctx, deptID)
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, Success(dept))
}

// UpdateDepartment 更新部门
// @router /api/org/departments/:id [PUT]
func (h *DepartmentHandler) UpdateDepartment(ctx context.Context, c *app.RequestContext) {
	deptID := c.Param("id")
	if deptID == "" {
		c.JSON(http.StatusBadRequest, Fail(errno.InvalidParams.Code(), "department id is required"))
		return
	}

	var req UpdateDepartmentRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, Fail(errno.InvalidParams.Code(), err.Error()))
		return
	}

	// 转换为服务层请求
	svcReq := &service.UpdateDepartmentRequest{
		DeptID:       deptID,
		DeptName:     req.DeptName,
		LeaderID:     req.LeaderID,
		ParentLeader: req.ParentLeader,
		Description:  req.Description,
		SortOrder:    req.SortOrder,
		Status:       req.Status,
	}

	dept, err := h.deptService.UpdateDepartment(ctx, svcReq)
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, Success(dept))
}

// DeleteDepartment 删除部门
// @router /api/org/departments/:id [DELETE]
func (h *DepartmentHandler) DeleteDepartment(ctx context.Context, c *app.RequestContext) {
	deptID := c.Param("id")
	if deptID == "" {
		c.JSON(http.StatusBadRequest, Fail(errno.InvalidParams.Code(), "department id is required"))
		return
	}

	if err := h.deptService.DeleteDepartment(ctx, deptID); err != nil {
		h.handleError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, Success(nil))
}

// GetDepartmentTree 获取部门树
// @router /api/org/departments/tree [GET]
func (h *DepartmentHandler) GetDepartmentTree(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, Fail(errno.InvalidParams.Code(), "tenant_id is required"))
		return
	}

	tree, err := h.deptService.GetDepartmentTree(ctx, tenantID)
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, Success(tree))
}

// ListDepartments 分页查询部门列表
// @router /api/org/departments [GET]
func (h *DepartmentHandler) ListDepartments(ctx context.Context, c *app.RequestContext) {
	var req ListDepartmentsRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, Fail(errno.InvalidParams.Code(), err.Error()))
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 转换为仓储过滤器
	filter := &repository.DepartmentFilter{
		TenantID: req.TenantID,
		OrgID:    req.OrgID,
		ParentID: parseStringPtr(req.ParentID),
		Level:    parseIntPtr(req.Level),
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	if req.Status != "" {
		filter.Status = req.Status
	}

	depts, total, err := h.deptService.ListDepartments(ctx, filter)
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, Success(map[string]interface{}{
		"departments": depts,
		"total":       total,
		"page":        req.Page,
		"page_size":   req.PageSize,
	}))
}

// MoveDepartment 移动部门
// @router /api/org/departments/:id/move [POST]
func (h *DepartmentHandler) MoveDepartment(ctx context.Context, c *app.RequestContext) {
	deptID := c.Param("id")
	if deptID == "" {
		c.JSON(http.StatusBadRequest, Fail(errno.InvalidParams.Code(), "department id is required"))
		return
	}

	var req MoveDepartmentRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, Fail(errno.InvalidParams.Code(), err.Error()))
		return
	}

	// 转换为服务层请求
	svcReq := &service.MoveDepartmentRequest{
		DeptID:      deptID,
		NewParentID: req.NewParentID,
	}

	if err := h.deptService.MoveDepartment(ctx, svcReq); err != nil {
		h.handleError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, Success(nil))
}

// GetChildren 获取子部门
// @router /api/org/departments/:id/children [GET]
func (h *DepartmentHandler) GetChildren(ctx context.Context, c *app.RequestContext) {
	parentID := c.Param("id")
	if parentID == "" {
		c.JSON(http.StatusBadRequest, Fail(errno.InvalidParams.Code(), "department id is required"))
		return
	}

	children, err := h.deptService.GetChildren(ctx, parentID)
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, Success(children))
}

// GetAncestors 获取部门的所有祖先
// @router /api/org/departments/:id/ancestors [GET]
func (h *DepartmentHandler) GetAncestors(ctx context.Context, c *app.RequestContext) {
	deptID := c.Param("id")
	if deptID == "" {
		c.JSON(http.StatusBadRequest, Fail(errno.InvalidParams.Code(), "department id is required"))
		return
	}

	ancestors, err := h.deptService.GetAncestors(ctx, deptID)
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, Success(ancestors))
}

// GetDescendants 获取部门的所有后代
// @router /api/org/departments/:id/descendants [GET]
func (h *DepartmentHandler) GetDescendants(ctx context.Context, c *app.RequestContext) {
	deptID := c.Param("id")
	if deptID == "" {
		c.JSON(http.StatusBadRequest, Fail(errno.InvalidParams.Code(), "department id is required"))
		return
	}

	descendants, err := h.deptService.GetDescendants(ctx, deptID)
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, Success(descendants))
}

// GetDepartmentsByOrg 获取组织的所有部门
// @router /api/org/departments/by-org/:org_id [GET]
func (h *DepartmentHandler) GetDepartmentsByOrg(ctx context.Context, c *app.RequestContext) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, Fail(errno.InvalidParams.Code(), "organization id is required"))
		return
	}

	depts, err := h.deptService.GetDepartmentsByOrg(ctx, orgID)
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, Success(depts))
}

// handleError 统一错误处理
func (h *DepartmentHandler) handleError(ctx context.Context, c *app.RequestContext, err error) {
	if codeErr, ok := err.(errno.ErrorCode); ok {
		c.JSON(codeErr.HTTPStatus(), Fail(codeErr.Code(), codeErr.Message()))
		return
	}
	c.JSON(http.StatusInternalServerError, Fail(errno.InternalError.Code(), errno.InternalError.Message()))
}

// Success 成功响应
func Success(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    data,
	}
}

// Fail 失败响应
func Fail(code int32, message string) map[string]interface{} {
	return map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    nil,
	}
}

// parseStringPtr 解析字符串指针
func parseStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// parseIntPtr 解析整型指针
func parseIntPtr(i int) *int {
	if i <= 0 {
		return nil
	}
	return &i
}
