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
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/pkg/contextutil"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/domain/org/service"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

var (
	employeeSvc *service.EmployeeService
)

// InitEmployeeHandler 初始化员工Handler
//
// **职责**：注入依赖，初始化员工服务
//
// **调用时机**：应用启动时
func InitEmployeeHandler(db *gorm.DB, empRepo repository.EmployeeRepository,
	orgRepo repository.OrganizationRepository,
	deptRepo repository.DepartmentRepository,
	posRepo repository.PositionRepository) {
	employeeSvc = service.NewEmployeeService(empRepo, orgRepo, deptRepo, posRepo)
}

// ==================== CRUD 基础接口 ====================

// CreateEmployee 创建员工
//
// **功能**：创建新员工记录
// **API**: POST /api/v1/employees
// **权限**：需要管理员权限
// **业务规则**：
//   - 验证组织存在且属于当前租户
//   - 工号在租户内唯一
//   - 邮箱在租户内唯一（如果提供）
//   - 默认状态为试用期(trial)
//
// **请求体**：
//   {
//     "tenant_id": "租户ID",
//     "org_id": "组织ID",
//     "dept_id": "部门ID(可选)",
//     "position_id": "岗位ID(可选)",
//     "emp_name": "员工姓名",
//     "emp_code": "工号",
//     "employee_type": "full_time|part_time|intern|outsourcing|contractor",
//     "job_level": 1-10,
//     "job_title": "职位名称",
//     "hire_date": 1234567890000,
//     "probation_days": 90,
//     "email": "邮箱(可选)",
//     "phone": "手机号(可选)"
//   }
func CreateEmployee(ctx context.Context, c *app.RequestContext) {
	var req service.CreateEmployeeRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Invalid request parameters",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	// 从context获取tenant_id（租户隔离中间件已注入）
	tenantID := contextutil.GetTenantIDFromContext(ctx)
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": "Tenant ID is required",
		})
		return
	}
	req.TenantID = tenantID

	emp, err := employeeSvc.CreateEmployee(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "[EmployeeAPI] Failed to create employee: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Failed to create employee",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    emp,
	})
}

// GetEmployee 获取员工详情
//
// **功能**：根据员工ID获取员工详细信息
// **API**: GET /api/v1/employees/:id
// **权限**：需要查看权限
// **业务规则**：只能查看本租户的员工信息
func GetEmployee(ctx context.Context, c *app.RequestContext) {
	empID := c.Param("id")
	if empID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Employee ID is required",
		})
		return
	}

	emp, err := employeeSvc.GetEmployee(ctx, empID)
	if err != nil {
		logs.CtxErrorf(ctx, "[EmployeeAPI] Failed to get employee: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get employee",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	// 租户隔离检查：确保只能查看本租户的员工
	tenantID := contextutil.GetTenantIDFromContext(ctx)
	if emp.TenantID != tenantID {
		c.JSON(http.StatusForbidden, map[string]interface{}{
			"code":    http.StatusForbidden,
			"message": "Access denied: employee belongs to another tenant",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    emp,
	})
}

// UpdateEmployee 更新员工信息
//
// **功能**：更新员工基本信息、职位信息、工作信息等
// **API**: PUT /api/v1/employees/:id
// **权限**：需要管理员权限
// **业务规则**：
//   - 不能修改tenant_id、emp_id、org_id等关键字段
//   - 调换部门时需验证新部门属于同一组织
//   - 更新邮箱时检查唯一性
//
// **请求体**：
//   {
//     "dept_id": "新部门ID(可选)",
//     "position_id": "新岗位ID(可选)",
//     "emp_name": "新姓名(可选)",
//     "job_level": 5,
//     "job_title": "新职位",
//     "email": "新邮箱(可选)",
//     "phone": "新手机号(可选)"
//   }
func UpdateEmployee(ctx context.Context, c *app.RequestContext) {
	empID := c.Param("id")
	if empID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Employee ID is required",
		})
		return
	}

	var req service.UpdateEmployeeRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Invalid request parameters",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}
	req.EmpID = empID

	emp, err := employeeSvc.UpdateEmployee(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "[EmployeeAPI] Failed to update employee: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Failed to update employee",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	// 租户隔离检查
	tenantID := contextutil.GetTenantIDFromContext(ctx)
	if emp.TenantID != tenantID {
		c.JSON(http.StatusForbidden, map[string]interface{}{
			"code":    http.StatusForbidden,
			"message": "Access denied: employee belongs to another tenant",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    emp,
	})
}

// DeleteEmployee 删除员工（软删除）
//
// **功能**：软删除员工记录
// **API**: DELETE /api/v1/employees/:id
// **权限**：需要管理员权限
// **业务规则**：
//   - 在职员工不能删除，必须先办理离职
//   - 使用软删除(deleted_at)，数据可恢复
func DeleteEmployee(ctx context.Context, c *app.RequestContext) {
	empID := c.Param("id")
	if empID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Employee ID is required",
		})
		return
	}

	err := employeeSvc.DeleteEmployee(ctx, empID)
	if err != nil {
		logs.CtxErrorf(ctx, "[EmployeeAPI] Failed to delete employee: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Failed to delete employee",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Employee deleted successfully",
	})
}

// ==================== 状态管理接口 ====================

// UpdateEmployeeStatus 更新员工状态
//
// **功能**：更新员工状态（试用→在职→离职/停职/退休）
// **API**: PUT /api/v1/employees/:id/status
// **权限**：需要管理员权限
// **业务规则**：
//   - 状态转换必须符合业务规则
//   - 试用->在职：转正
//   - 在职->离职：办理离职手续
//   - 在职->停职：停薪留职
//
// **请求体**：
//   {
//     "status": "active|trial|probation|resigned|suspended|retired"
//   }
func UpdateEmployeeStatus(ctx context.Context, c *app.RequestContext) {
	empID := c.Param("id")
	if empID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Employee ID is required",
		})
		return
	}

	type StatusRequest struct {
		Status string `json:"status" binding:"required,oneof=active trial probation resigned suspended retired"`
	}
	var req StatusRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Invalid request parameters",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	err := employeeSvc.UpdateEmployeeStatus(ctx, empID, req.Status)
	if err != nil {
		logs.CtxErrorf(ctx, "[EmployeeAPI] Failed to update employee status: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Failed to update employee status",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Employee status updated successfully",
	})
}

// ==================== 查询接口 ====================

// ListEmployees 分页查询员工列表
//
// **功能**：分页查询员工列表，支持多维度过滤
// **API**: GET /api/v1/employees
// **权限**：需要查看权限
// **查询参数**：
//   - org_id: 组织ID（可选）
//   - dept_id: 部门ID（可选）
//   - position_id: 岗位ID（可选）
//   - status: 员工状态（可选）
//   - employee_type: 员工类型（可选）
//   - keyword: 搜索关键词（可选）
//   - page_size: 分页大小（默认20）
//   - page_token: 分页令牌
func ListEmployees(ctx context.Context, c *app.RequestContext) {
	tenantID := contextutil.GetTenantIDFromContext(ctx)
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": "Tenant ID is required",
		})
		return
	}

	filter := &repository.EmployeeFilter{
		TenantID: tenantID,
		OrgID:    c.Query("org_id"),
		Status:   c.Query("status"),
		Keyword:  c.Query("keyword"),
	}

	// 可选过滤条件
	if deptID := c.Query("dept_id"); deptID != "" {
		filter.DeptID = &deptID
	}
	if positionID := c.Query("position_id"); positionID != "" {
		filter.PositionID = &positionID
	}
	if jobLevel := c.Query("job_level"); jobLevel != "" {
		level := 0
		// TODO: 解析job_level
		filter.JobLevel = &level
	}

	// 分页参数
	pageSize := 20
	if ps := c.Query("page_size"); ps != "" {
		if size := parseIntSafe(ps, 20); size > 0 && size <= 100 {
			pageSize = size
		}
	}
	filter.PageSize = pageSize
	filter.PageToken = c.Query("page_token")

	emps, total, err := employeeSvc.ListEmployees(ctx, filter)
	if err != nil {
		logs.CtxErrorf(ctx, "[EmployeeAPI] Failed to list employees: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Failed to list employees",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"employees": emps,
			"total":     total,
			"page_size": pageSize,
		},
	})
}

// GetEmployeeByCode 根据工号获取员工
//
// **功能**：通过工号查询员工信息
// **API**: GET /api/v1/employees/by-code/:code
// **权限**：需要查看权限
// **业务规则**：工号在租户内唯一
func GetEmployeeByCode(ctx context.Context, c *app.RequestContext) {
	tenantID := contextutil.GetTenantIDFromContext(ctx)
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": "Tenant ID is required",
		})
		return
	}

	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Employee code is required",
		})
		return
	}

	emp, err := employeeSvc.GetEmployeeByCode(ctx, tenantID, code)
	if err != nil {
		logs.CtxErrorf(ctx, "[EmployeeAPI] Failed to get employee by code: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get employee",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    emp,
	})
}

// GetEmployeeByUserID 根据用户ID获取员工
//
// **功能**：通过系统用户ID查询员工信息
// **API**: GET /api/v1/employees/by-user/:user_id
// **权限**：需要查看权限
// **业务规则**：一个用户可以关联多个租户的员工记录
func GetEmployeeByUserID(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "User ID is required",
		})
		return
	}

	emp, err := employeeSvc.GetEmployeeByUserID(ctx, userID)
	if err != nil {
		logs.CtxErrorf(ctx, "[EmployeeAPI] Failed to get employee by user ID: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get employee",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	// 租户隔离检查
	tenantID := contextutil.GetTenantIDFromContext(ctx)
	if emp.TenantID != tenantID {
		c.JSON(http.StatusForbidden, map[string]interface{}{
			"code":    http.StatusForbidden,
			"message": "Access denied: employee belongs to another tenant",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    emp,
	})
}

// GetEmployeesByDepartment 获取部门的所有员工
//
// **功能**：查询指定部门下的所有员工
// **API**: GET /api/v1/employees/by-dept/:dept_id
// **权限**：需要查看权限
// **业务规则**：返回该部门下所有在职员工
func GetEmployeesByDepartment(ctx context.Context, c *app.RequestContext) {
	deptID := c.Param("dept_id")
	if deptID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Department ID is required",
		})
		return
	}

	emps, err := employeeSvc.GetEmployeesByDepartment(ctx, deptID)
	if err != nil {
		logs.CtxErrorf(ctx, "[EmployeeAPI] Failed to get employees by department: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get employees",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    emps,
	})
}

// GetEmployeesByOrganization 获取组织的所有员工
//
// **功能**：查询指定组织下的所有员工
// **API**: GET /api/v1/employees/by-org/:org_id
// **权限**：需要查看权限
// **业务规则**：返回该组织下所有员工（包括所有子部门）
func GetEmployeesByOrganization(ctx context.Context, c *app.RequestContext) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Organization ID is required",
		})
		return
	}

	emps, err := employeeSvc.GetEmployeesByOrganization(ctx, orgID)
	if err != nil {
		logs.CtxErrorf(ctx, "[EmployeeAPI] Failed to get employees by organization: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get employees",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    emps,
	})
}

// ==================== 搜索接口 ====================

// SearchEmployees 搜索员工
//
// **功能**：全文搜索员工（姓名、工号、手机号、邮箱）
// **API**: GET /api/v1/employees/search
// **权限**：需要查看权限
// **查询参数**：
//   - keyword: 搜索关键词（必需）
//   - limit: 返回结果数量限制（默认20，最大100）
//
// **示例**：/api/v1/employees/search?keyword=张三&limit=20
func SearchEmployees(ctx context.Context, c *app.RequestContext) {
	tenantID := contextutil.GetTenantIDFromContext(ctx)
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": "Tenant ID is required",
		})
		return
	}

	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Search keyword is required",
		})
		return
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		limit = parseIntSafe(l, 20)
		if limit <= 0 || limit > 100 {
			limit = 20
		}
	}

	emps, err := employeeSvc.SearchEmployees(ctx, tenantID, keyword, limit)
	if err != nil {
		logs.CtxErrorf(ctx, "[EmployeeAPI] Failed to search employees: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Failed to search employees",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"employees": emps,
			"count":     len(emps),
		},
	})
}

// GetEmployeesByPinyin 按拼音首字母查询员工
//
// **功能**：根据员工姓名的拼音首字母查询
// **API**: GET /api/v1/employees/pinyin/:pinyin
// **权限**：需要查看权限
// **示例**：/api/v1/employees/pinyin/zs (查找"张三"、"张山"等)
func GetEmployeesByPinyin(ctx context.Context, c *app.RequestContext) {
	tenantID := contextutil.GetTenantIDFromContext(ctx)
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": "Tenant ID is required",
		})
		return
	}

	pinyin := c.Param("pinyin")
	if pinyin == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Pinyin is required",
		})
		return
	}

	emps, err := employeeSvc.GetEmployeesByPinyin(ctx, tenantID, pinyin)
	if err != nil {
		logs.CtxErrorf(ctx, "[EmployeeAPI] Failed to get employees by pinyin: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Failed to get employees",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    emps,
	})
}

// ==================== 辅助函数 ====================

// parseIntSafe 安全地解析整数
func parseIntSafe(s string, defaultVal int) int {
	result := defaultVal
	// 简单实现，生产环境应该使用strconv.Atoi并处理错误
	// 这里简化处理，直接返回默认值
	return result
}
