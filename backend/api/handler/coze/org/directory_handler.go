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
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/api/handler/coze"
	"github.com/coze-dev/coze-studio/backend/pkg/contextutil"
	"github.com/coze-dev/coze-studio/backend/domain/org/service"
)

var (
	// directoryService 通讯录服务实例（在应用启动时注入）
	directoryService *service.DirectoryService
)

// InitDirectoryHandler 初始化通讯录Handler依赖
//
// **使用场景**：应用启动时调用，注入DirectoryService
//
// **示例**：
//   org.InitDirectoryHandler(directoryService)
func InitDirectoryHandler(ds *service.DirectoryService) {
	directoryService = ds
}

// ==================== 组织目录接口 ====================

// GetOrganizationDirectory 获取组织架构目录
// @router /api/directory/organization [GET]
// @param tenant_id query string false "租户ID（可选，默认从context获取）"
// @success 200 {object} object
func GetOrganizationDirectory(ctx context.Context, c *app.RequestContext) {
	if directoryService == nil {
		coze.InternalServerErrorResponse(ctx, c, ErrServiceNotInitialized)
		return
	}

	// 获取tenant_id（优先从query参数，否则从context）
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		tenantID = contextutil.GetTenantIDFromContext(ctx)
	}

	if tenantID == "" {
		coze.InvalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 调用服务层
	result, err := directoryService.GetOrganizationDirectory(ctx, tenantID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    result,
	})
}

// GetDepartmentDirectory 获取部门目录
// @router /api/directory/department [GET]
// @param tenant_id query string false "租户ID（可选，默认从context获取）"
// @success 200 {object} object
func GetDepartmentDirectory(ctx context.Context, c *app.RequestContext) {
	if directoryService == nil {
		coze.InternalServerErrorResponse(ctx, c, ErrServiceNotInitialized)
		return
	}

	// 获取tenant_id（优先从query参数，否则从context）
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		tenantID = contextutil.GetTenantIDFromContext(ctx)
	}

	if tenantID == "" {
		coze.InvalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 调用服务层
	result, err := directoryService.GetDepartmentDirectory(ctx, tenantID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    result,
	})
}

// ==================== 员工列表接口 ====================

// GetDepartmentEmployees 获取部门员工列表
// @router /api/directory/departments/:dept_id/employees [GET]
// @param dept_id path string true "部门ID"
// @success 200 {object} object
func GetDepartmentEmployees(ctx context.Context, c *app.RequestContext) {
	if directoryService == nil {
		coze.InternalServerErrorResponse(ctx, c, ErrServiceNotInitialized)
		return
	}

	deptID := c.Param("dept_id")
	if deptID == "" {
		coze.InvalidParamRequestResponse(c, "dept_id is required")
		return
	}

	// 调用服务层
	employees, err := directoryService.GetDepartmentEmployees(ctx, deptID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"employees": employees,
			"total":     len(employees),
		},
	})
}

// GetOrganizationEmployees 获取组织员工列表
// @router /api/directory/organizations/:org_id/employees [GET]
// @param org_id path string true "组织ID"
// @success 200 {object} object
func GetOrganizationEmployees(ctx context.Context, c *app.RequestContext) {
	if directoryService == nil {
		coze.InternalServerErrorResponse(ctx, c, ErrServiceNotInitialized)
		return
	}

	orgID := c.Param("org_id")
	if orgID == "" {
		coze.InvalidParamRequestResponse(c, "org_id is required")
		return
	}

	// 调用服务层
	employees, err := directoryService.GetOrganizationEmployees(ctx, orgID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"employees": employees,
			"total":     len(employees),
		},
	})
}

// ==================== 员工搜索接口 ====================

// SearchEmployeesRequest 搜索员工请求参数
type SearchEmployeesRequest struct {
	TenantID string `json:"tenant_id" query:"tenant_id" binding:"omitempty"`
	Keyword  string `json:"keyword" query:"keyword" binding:"required,min=1,max=50"`
	Limit    int    `json:"limit" query:"limit" binding:"omitempty,min=1,max=100"`
}

// SearchEmployees 搜索员工
// @router /api/directory/search [GET]
// @param keyword query string true "搜索关键词（姓名/工号/邮箱/手机号）"
// @param limit query int false "返回数量限制（默认20，最大100）"
// @param tenant_id query string false "租户ID（可选，默认从context获取）"
// @success 200 {object} object
func SearchEmployees(ctx context.Context, c *app.RequestContext) {
	if directoryService == nil {
		coze.InternalServerErrorResponse(ctx, c, ErrServiceNotInitialized)
		return
	}

	var req SearchEmployeesRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		coze.InvalidParamRequestResponse(c, err.Error())
		return
	}

	// 获取tenant_id（优先从请求参数，否则从context）
	tenantID := req.TenantID
	if tenantID == "" {
		tenantID = contextutil.GetTenantIDFromContext(ctx)
	}

	if tenantID == "" {
		coze.InvalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 设置默认limit
	if req.Limit <= 0 {
		req.Limit = 20
	}

	// 调用服务层
	employees, err := directoryService.SearchEmployees(ctx, tenantID, req.Keyword, req.Limit)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"employees": employees,
			"total":     len(employees),
			"keyword":   req.Keyword,
		},
	})
}

// ==================== 员工详情接口 ====================

// GetEmployeeByCode 根据工号获取员工信息
// @router /api/directory/employee/by-code/:code [GET]
// @param code path string true "工号"
// @param tenant_id query string false "租户ID（可选，默认从context获取）"
// @success 200 {object} object
func GetEmployeeByCode(ctx context.Context, c *app.RequestContext) {
	if directoryService == nil {
		coze.InternalServerErrorResponse(ctx, c, ErrServiceNotInitialized)
		return
	}

	code := c.Param("code")
	if code == "" {
		coze.InvalidParamRequestResponse(c, "code is required")
		return
	}

	// 获取tenant_id（优先从query参数，否则从context）
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		tenantID = contextutil.GetTenantIDFromContext(ctx)
	}

	if tenantID == "" {
		coze.InvalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 调用服务层
	employee, err := directoryService.GetEmployeeByCode(ctx, tenantID, code)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    employee,
	})
}

// ==================== 错误定义 ====================

// ErrServiceNotInitialized 服务未初始化错误
var ErrServiceNotInitialized = fmt.Errorf("directory service not initialized")
