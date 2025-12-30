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
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	orgservice "github.com/coze-dev/coze-studio/backend/domain/org/service"
)

// PositionHandler 岗位管理处理器
// 职责：处理HTTP请求，调用领域服务，返回HTTP响应
// 遵循SOLID原则：
// - 单一职责：仅处理HTTP层逻辑
// - 依赖倒置：依赖PositionService接口而非具体实现
// - 开闭原则：通过接口扩展，无需修改现有代码
type PositionHandler struct {
	positionService *orgservice.PositionService
}

// NewPositionHandler 创建岗位处理器实例
// 依赖注入：通过构造函数注入PositionService
func NewPositionHandler(positionService *orgservice.PositionService) *PositionHandler {
	return &PositionHandler{
		positionService: positionService,
	}
}

// ==================== 基础CRUD操作 ====================

// CreatePosition 创建岗位
// @router /api/positions [POST]
// @summary 创建岗位
// @description 创建新岗位，支持指定所属部门、职级、类别等信息
// @tags 岗位管理
// @accept json
// @produce json
// @param request body org.CreatePositionRequest true "创建岗位请求"
// @success 200 {object} org.CreatePositionResponse
// @failure 400 {object} coze.data
// @failure 500 {object} coze.data
func (h *PositionHandler) CreatePosition(ctx context.Context, c *app.RequestContext) {
	// 1. 参数绑定与验证
	var req org.CreatePositionRequest
	if err := c.BindAndValidate(&req); err != nil {
		coze.InvalidParamRequestResponse(c, err.Error())
		return
	}

	// 2. 转换为领域服务请求
	serviceReq := &orgservice.CreatePositionRequest{
		TenantID:         req.TenantID,
		DeptID:           req.DeptID,
		PositionName:     req.PositionName,
		PositionCode:     req.PositionCode,
		Level:            req.Level,
		Category:         req.Category,
		Responsibilities: req.Responsibilities,
		Requirements:     req.Requirements,
		SortOrder:        req.SortOrder,
	}

	// 3. 调用领域服务
	positionEntity, err := h.positionService.CreatePosition(ctx, serviceReq)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 4. 转换并返回响应
	c.JSON(http.StatusOK, &org.CreatePositionResponse{
		Code:    0,
		Message: "success",
		Data:    toPositionData(positionEntity),
	})
}

// GetPosition 获取岗位详情
// @router /api/positions/:id [GET]
// @summary 获取岗位详情
// @description 根据岗位ID获取岗位详细信息
// @tags 岗位管理
// @accept json
// @produce json
// @param id path string true "岗位ID"
// @success 200 {object} org.GetPositionResponse
// @failure 400 {object} coze.data
// @failure 404 {object} coze.data
// @failure 500 {object} coze.data
func (h *PositionHandler) GetPosition(ctx context.Context, c *app.RequestContext) {
	// 1. 获取路径参数
	positionID := c.Param("id")
	if positionID == "" {
		coze.InvalidParamRequestResponse(c, "position id is required")
		return
	}

	// 2. 调用领域服务
	positionEntity, err := h.positionService.GetPosition(ctx, positionID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 3. 转换并返回响应
	c.JSON(http.StatusOK, &org.GetPositionResponse{
		Code:    0,
		Message: "success",
		Data:    toPositionData(positionEntity),
	})
}

// UpdatePosition 更新岗位信息
// @router /api/positions/:id [PUT]
// @summary 更新岗位信息
// @description 更新岗位的基本信息、职责、要求等
// @tags 岗位管理
// @accept json
// @produce json
// @param id path string true "岗位ID"
// @param request body org.UpdatePositionRequest true "更新岗位请求"
// @success 200 {object} org.UpdatePositionResponse
// @failure 400 {object} coze.data
// @failure 404 {object} coze.data
// @failure 500 {object} coze.data
func (h *PositionHandler) UpdatePosition(ctx context.Context, c *app.RequestContext) {
	// 1. 获取路径参数
	positionID := c.Param("id")
	if positionID == "" {
		coze.InvalidParamRequestResponse(c, "position id is required")
		return
	}

	// 2. 参数绑定与验证
	var req org.UpdatePositionRequest
	if err := c.BindAndValidate(&req); err != nil {
		coze.InvalidParamRequestResponse(c, err.Error())
		return
	}

	// 3. 转换为领域服务请求
	serviceReq := &orgservice.UpdatePositionRequest{
		PositionID:       positionID,
		PositionName:     req.PositionName,
		Level:            req.Level,
		Category:         req.Category,
		Responsibilities: req.Responsibilities,
		Requirements:     req.Requirements,
		SortOrder:        req.SortOrder,
		Status:           req.Status,
	}

	// 4. 调用领域服务
	positionEntity, err := h.positionService.UpdatePosition(ctx, serviceReq)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 5. 转换并返回响应
	c.JSON(http.StatusOK, &org.UpdatePositionResponse{
		Code:    0,
		Message: "success",
		Data:    toPositionData(positionEntity),
	})
}

// DeletePosition 删除岗位
// @router /api/positions/:id [DELETE]
// @summary 删除岗位
// @description 软删除岗位，如果岗位下有员工则无法删除
// @tags 岗位管理
// @accept json
// @produce json
// @param id path string true "岗位ID"
// @success 200 {object} org.DeletePositionResponse
// @failure 400 {object} coze.data
// @failure 403 {object} coze.data
// @failure 404 {object} coze.data
// @failure 500 {object} coze.data
func (h *PositionHandler) DeletePosition(ctx context.Context, c *app.RequestContext) {
	// 1. 获取路径参数
	positionID := c.Param("id")
	if positionID == "" {
		coze.InvalidParamRequestResponse(c, "position id is required")
		return
	}

	// 2. 调用领域服务
	err := h.positionService.DeletePosition(ctx, positionID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 3. 返回成功响应
	c.JSON(http.StatusOK, &org.DeletePositionResponse{
		Code:    0,
		Message: "success",
	})
}

// ==================== 查询操作 ====================

// ListPositions 分页查询岗位列表
// @router /api/positions [GET]
// @summary 分页查询岗位列表
// @description 支持按部门、职级、类别、状态等条件筛选岗位
// @tags 岗位管理
// @accept json
// @produce json
// @param tenant_id query string true "租户ID"
// @param dept_id query string false "部门ID"
// @param level query int false "职级"
// @param category query string false "岗位类别"
// @param status query string false "状态" Enums(active, inactive, frozen)
// @param keyword query string false "关键词（岗位名称或编码）"
// @param page query int false "页码" default(1)
// @param page_size query int false "每页数量" default(20)
// @param sort_by query string false "排序字段" Enums(created_at, updated_at, sort_order, position_name)
// @param sort_order_dir query string false "排序方向" Enums(asc, desc)
// @success 200 {object} org.ListPositionsResponse
// @failure 400 {object} coze.data
// @failure 500 {object} coze.data
func (h *PositionHandler) ListPositions(ctx context.Context, c *app.RequestContext) {
	// 1. 获取查询参数
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		coze.InvalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 2. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 3. 解析查询参数
	var level *int
	if levelStr := c.Query("level"); levelStr != "" {
		if l, err := strconv.Atoi(levelStr); err == nil {
			level = &l
		}
	}

	// 4. 构建查询过滤器
	filter := &repository.PositionFilter{
		TenantID:  tenantID,
		DeptID:    c.Query("dept_id"),
		Level:     level,
		Category:  c.Query("category"),
		Status:    c.Query("status"),
		Keyword:   c.Query("keyword"),
		Page:      page,
		PageSize:  pageSize,
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order_dir"),
	}

	// 5. 调用领域服务
	positions, total, err := h.positionService.ListPositions(ctx, filter)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 6. 转换并返回响应
	c.JSON(http.StatusOK, &org.ListPositionsResponse{
		Code:     0,
		Message:  "success",
		Data:     toPositionDataList(positions),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetPositionByCode 根据编码获取岗位
// @router /api/positions/by-code/:code [GET]
// @summary 根据编码获取岗位
// @description 通过岗位编码获取岗位信息，编码在租户内唯一
// @tags 岗位管理
// @accept json
// @produce json
// @param code path string true "岗位编码"
// @param tenant_id query string true "租户ID"
// @success 200 {object} org.GetPositionByCodeResponse
// @failure 400 {object} coze.data
// @failure 404 {object} coze.data
// @failure 500 {object} coze.data
func (h *PositionHandler) GetPositionByCode(ctx context.Context, c *app.RequestContext) {
	// 1. 获取路径参数
	code := c.Param("code")
	if code == "" {
		coze.InvalidParamRequestResponse(c, "position code is required")
		return
	}

	// 2. 获取租户ID
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		coze.InvalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 3. 调用领域服务
	positionEntity, err := h.positionService.GetPositionByCode(ctx, tenantID, code)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 4. 转换并返回响应
	c.JSON(http.StatusOK, &org.GetPositionByCodeResponse{
		Code:    0,
		Message: "success",
		Data:    toPositionData(positionEntity),
	})
}

// GetPositionsByDepartment 获取部门的所有岗位
// @router /api/positions/by-dept/:dept_id [GET]
// @summary 获取部门的所有岗位
// @description 获取指定部门下的所有岗位列表
// @tags 岗位管理
// @accept json
// @produce json
// @param dept_id path string true "部门ID"
// @success 200 {object} org.GetPositionsByDepartmentResponse
// @failure 400 {object} coze.data
// @failure 500 {object} coze.data
func (h *PositionHandler) GetPositionsByDepartment(ctx context.Context, c *app.RequestContext) {
	// 1. 获取路径参数
	deptID := c.Param("dept_id")
	if deptID == "" {
		coze.InvalidParamRequestResponse(c, "department id is required")
		return
	}

	// 2. 调用领域服务
	positions, err := h.positionService.GetPositionsByDepartment(ctx, deptID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 3. 转换并返回响应
	c.JSON(http.StatusOK, &org.GetPositionsByDepartmentResponse{
		Code:    0,
		Message: "success",
		Data:    toPositionDataList(positions),
	})
}

// GetPositionsByLevel 按职级获取岗位
// @router /api/positions/by-level/:level [GET]
// @summary 按职级获取岗位
// @description 获取指定职级的所有岗位
// @tags 岗位管理
// @accept json
// @produce json
// @param level path int true "职级"
// @param tenant_id query string true "租户ID"
// @success 200 {object} org.GetPositionsByLevelResponse
// @failure 400 {object} coze.data
// @failure 500 {object} coze.data
func (h *PositionHandler) GetPositionsByLevel(ctx context.Context, c *app.RequestContext) {
	// 1. 获取路径参数
	levelStr := c.Param("level")
	if levelStr == "" {
		coze.InvalidParamRequestResponse(c, "level is required")
		return
	}

	// 2. 解析职级
	level, err := strconv.Atoi(levelStr)
	if err != nil || level <= 0 {
		coze.InvalidParamRequestResponse(c, "invalid level format")
		return
	}

	// 3. 获取租户ID
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		coze.InvalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 4. 调用领域服务
	positions, err := h.positionService.GetPositionsByLevel(ctx, tenantID, level)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 5. 转换并返回响应
	c.JSON(http.StatusOK, &org.GetPositionsByLevelResponse{
		Code:    0,
		Message: "success",
		Data:    toPositionDataList(positions),
	})
}

// GetPositionsByCategory 按类别获取岗位
// @router /api/positions/by-category/:category [GET]
// @summary 按类别获取岗位
// @description 获取指定类别的所有岗位（如：技术岗、管理岗、职能岗）
// @tags 岗位管理
// @accept json
// @produce json
// @param category path string true "岗位类别"
// @param tenant_id query string true "租户ID"
// @success 200 {object} org.GetPositionsByCategoryResponse
// @failure 400 {object} coze.data
// @failure 500 {object} coze.data
func (h *PositionHandler) GetPositionsByCategory(ctx context.Context, c *app.RequestContext) {
	// 1. 获取路径参数
	category := c.Param("category")
	if category == "" {
		coze.InvalidParamRequestResponse(c, "category is required")
		return
	}

	// 2. 获取租户ID
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		coze.InvalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 3. 调用领域服务
	positions, err := h.positionService.GetPositionsByCategory(ctx, tenantID, category)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 4. 转换并返回响应
	c.JSON(http.StatusOK, &org.GetPositionsByCategoryResponse{
		Code:    0,
		Message: "success",
		Data:    toPositionDataList(positions),
	})
}

// ==================== 辅助函数 ====================

// toPositionData 将领域实体转换为API数据模型
func toPositionData(entity *orgentity.Position) *org.PositionData {
	if entity == nil {
		return nil
	}

	return &org.PositionData{
		PositionID:       entity.PositionID,
		TenantID:         entity.TenantID,
		DeptID:           entity.DeptID,
		PositionName:     entity.PositionName,
		PositionCode:     entity.PositionCode,
		Level:            entity.Level,
		Category:         entity.Category,
		Responsibilities: entity.Responsibilities,
		Requirements:     entity.Requirements,
		Status:           string(entity.Status),
		SortOrder:        entity.SortOrder,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
	}
}

// toPositionDataList 批量转换领域实体列表
func toPositionDataList(entities []*orgentity.Position) []*org.PositionData {
	if entities == nil {
		return nil
	}

	dataList := make([]*org.PositionData, 0, len(entities))
	for _, entity := range entities {
		dataList = append(dataList, toPositionData(entity))
	}
	return dataList
}
