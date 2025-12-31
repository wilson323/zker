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

package coze

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/domain/digital_employee/entity"
	appDigitalEmployee "github.com/coze-dev/coze-studio/backend/application/digital_employee"
)

// getDigitalEmployeeSvc 获取数字员工服务容器
func getDigitalEmployeeSvc() *appDigitalEmployee.ApplicationService {
	return appDigitalEmployee.GetService()
}

// ==================== 员工画像管理接口 ====================

// CreateEmployeeProfile 创建员工画像
// @router /api/v1/digital-employees [POST]
func CreateEmployeeProfile(ctx context.Context, c *app.RequestContext) {
	svc := getDigitalEmployeeSvc()
	if svc == nil {
		InternalServerErrorResponse(ctx, c, nil)
		return
	}

	var err error
	var req entity.CreateProfileRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		InvalidParamRequestResponse(c, err.Error())
		return
	}

	profile, err := svc.DomainSVC.Profile.CreateProfile(ctx, &req)
	if err != nil {
		InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"code":    0,
		"message": "Employee profile created successfully",
		"data":    profile,
	})
}

// UpdateEmployeeProfile 更新员工画像
// @router /api/v1/digital-employees/:employee_id [PUT]
func UpdateEmployeeProfile(ctx context.Context, c *app.RequestContext) {
	svc := getDigitalEmployeeSvc()
	if svc == nil {
		InternalServerErrorResponse(ctx, c, nil)
		return
	}

	employeeID := c.Param("employee_id")
	if employeeID == "" {
		InvalidParamRequestResponse(c, "employee_id is required")
		return
	}

	var err error
	var req entity.UpdateProfileRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		InvalidParamRequestResponse(c, err.Error())
		return
	}

	req.EmployeeID = employeeID

	profile, err := svc.DomainSVC.Profile.UpdateProfile(ctx, &req)
	if err != nil {
		InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Employee profile updated successfully",
		"data":    profile,
	})
}

// GetEmployeeProfile 获取员工画像
// @router /api/v1/digital-employees/:employee_id [GET]
func GetEmployeeProfile(ctx context.Context, c *app.RequestContext) {
	svc := getDigitalEmployeeSvc()
	if svc == nil {
		InternalServerErrorResponse(ctx, c, nil)
		return
	}

	employeeID := c.Param("employee_id")
	if employeeID == "" {
		InvalidParamRequestResponse(c, "employee_id is required")
		return
	}

	profile, err := svc.DomainSVC.Profile.GetProfile(ctx, employeeID)
	if err != nil {
		InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Success",
		"data":    profile,
	})
}

// ListEmployeeProfiles 列出员工画像
// @router /api/v1/digital-employees [GET]
func ListEmployeeProfiles(ctx context.Context, c *app.RequestContext) {
	svc := getDigitalEmployeeSvc()
	if svc == nil {
		InternalServerErrorResponse(ctx, c, nil)
		return
	}

	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		InvalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 解析分页参数
	page := parseIntParam(c, "page", 1)
	pageSize := parseIntParam(c, "page_size", 20)

	req := &entity.ListProfilesRequest{
		TenantID: tenantID,
		Page:     page,
		PageSize: pageSize,
	}

	profiles, total, err := svc.DomainSVC.Profile.ListProfiles(ctx, req)
	if err != nil {
		InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Success",
		"data": map[string]interface{}{
			"profiles":  profiles,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// DeleteEmployeeProfile 删除员工画像
// @router /api/v1/digital-employees/:employee_id [DELETE]
func DeleteEmployeeProfile(ctx context.Context, c *app.RequestContext) {
	svc := getDigitalEmployeeSvc()
	if svc == nil {
		InternalServerErrorResponse(ctx, c, nil)
		return
	}

	employeeID := c.Param("employee_id")
	if employeeID == "" {
		InvalidParamRequestResponse(c, "employee_id is required")
		return
	}

	err := svc.DomainSVC.Profile.DeleteProfile(ctx, employeeID)
	if err != nil {
		InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Employee profile deleted successfully",
	})
}

// ==================== 任务分配管理接口 ====================

// AssignTask 分配任务给员工
// @router /api/v1/digital-employees/tasks/assign [POST]
func AssignTask(ctx context.Context, c *app.RequestContext) {
	svc := getDigitalEmployeeSvc()
	if svc == nil {
		InternalServerErrorResponse(ctx, c, nil)
		return
	}

	var err error
	var req entity.AssignTaskRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		InvalidParamRequestResponse(c, err.Error())
		return
	}

	assignment, err := svc.DomainSVC.Task.AssignTask(ctx, &req)
	if err != nil {
		InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"code":    0,
		"message": "Task assigned successfully",
		"data":    assignment,
	})
}

// AutoAssignTask 自动分配任务（智能技能匹配）
// @router /api/v1/digital-employees/tasks/auto-assign [POST]
func AutoAssignTask(ctx context.Context, c *app.RequestContext) {
	svc := getDigitalEmployeeSvc()
	if svc == nil {
		InternalServerErrorResponse(ctx, c, nil)
		return
	}

	var err error
	var req entity.AutoAssignTaskRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		InvalidParamRequestResponse(c, err.Error())
		return
	}

	assignment, err := svc.DomainSVC.Task.AutoAssignTask(ctx, &req)
	if err != nil {
		InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"code":    0,
		"message": "Task auto-assigned successfully",
		"data":    assignment,
	})
}

// CompleteTask 完成任务
// @router /api/v1/digital-employees/tasks/:assignment_id/complete [PUT]
func CompleteTask(ctx context.Context, c *app.RequestContext) {
	svc := getDigitalEmployeeSvc()
	if svc == nil {
		InternalServerErrorResponse(ctx, c, nil)
		return
	}

	assignmentID := c.Param("assignment_id")
	if assignmentID == "" {
		InvalidParamRequestResponse(c, "assignment_id is required")
		return
	}

	var err error
	var req entity.CompleteTaskRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		InvalidParamRequestResponse(c, err.Error())
		return
	}

	req.AssignmentID = assignmentID

	err = svc.DomainSVC.Task.CompleteTask(ctx, &req)
	if err != nil {
		InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Task completed successfully",
	})
}

// FailTask 标记任务失败
// @router /api/v1/digital-employees/tasks/:assignment_id/fail [PUT]
func FailTask(ctx context.Context, c *app.RequestContext) {
	svc := getDigitalEmployeeSvc()
	if svc == nil {
		InternalServerErrorResponse(ctx, c, nil)
		return
	}

	assignmentID := c.Param("assignment_id")
	if assignmentID == "" {
		InvalidParamRequestResponse(c, "assignment_id is required")
		return
	}

	var err error
	var req entity.FailTaskRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		InvalidParamRequestResponse(c, err.Error())
		return
	}

	req.AssignmentID = assignmentID

	err = svc.DomainSVC.Task.FailTask(ctx, &req)
	if err != nil {
		InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Task marked as failed",
	})
}

// GetEmployeeTasks 获取员工任务列表
// @router /api/v1/digital-employees/:employee_id/tasks [GET]
func GetEmployeeTasks(ctx context.Context, c *app.RequestContext) {
	svc := getDigitalEmployeeSvc()
	if svc == nil {
		InternalServerErrorResponse(ctx, c, nil)
		return
	}

	employeeID := c.Param("employee_id")
	if employeeID == "" {
		InvalidParamRequestResponse(c, "employee_id is required")
		return
	}

	// 解析分页参数
	page := parseIntParam(c, "page", 1)
	pageSize := parseIntParam(c, "page_size", 20)

	req := &entity.GetAssignmentsRequest{
		EmployeeID: employeeID,
		Page:       page,
		PageSize:   pageSize,
	}

	resp, err := svc.DomainSVC.Task.GetAssignments(ctx, req)
	if err != nil {
		InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Success",
		"data": map[string]interface{}{
			"assignments": resp.Assignments,
			"total":       resp.Total,
			"page":        resp.Page,
			"page_size":   resp.PageSize,
		},
	})
}

// ==================== 绩效统计接口 ====================

// GetEmployeePerformance 获取员工绩效
// @router /api/v1/digital-employees/:employee_id/performance [GET]
func GetEmployeePerformance(ctx context.Context, c *app.RequestContext) {
	svc := getDigitalEmployeeSvc()
	if svc == nil {
		InternalServerErrorResponse(ctx, c, nil)
		return
	}

	employeeID := c.Param("employee_id")
	if employeeID == "" {
		InvalidParamRequestResponse(c, "employee_id is required")
		return
	}

	period := c.DefaultQuery("period", "daily")
	date := c.Query("date")

	req := &entity.GetPerformanceRequest{
		EmployeeID: employeeID,
		Period:     entity.PerformancePeriod(period),
		Date:       &date,
	}

	performance, err := svc.DomainSVC.Performance.GetPerformance(ctx, req)
	if err != nil {
		InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Success",
		"data":    performance,
	})
}

// GetTeamPerformance 获取团队绩效
// @router /api/v1/digital-employees/performance/team [GET]
func GetTeamPerformance(ctx context.Context, c *app.RequestContext) {
	svc := getDigitalEmployeeSvc()
	if svc == nil {
		InternalServerErrorResponse(ctx, c, nil)
		return
	}

	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		InvalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	period := c.DefaultQuery("period", "daily")
	date := c.Query("date")

	// 解析分页参数
	page := parseIntParam(c, "page", 1)
	pageSize := parseIntParam(c, "page_size", 20)

	req := &entity.GetTeamPerformanceRequest{
		TenantID: tenantID,
		Period:   entity.PerformancePeriod(period),
		Date:     &date,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := svc.DomainSVC.Performance.GetTeamPerformance(ctx, req)
	if err != nil {
		InternalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Success",
		"data": map[string]interface{}{
			"performances": resp.Performances,
			"total":        resp.Total,
			"page":         resp.Page,
			"page_size":    resp.PageSize,
		},
	})
}

// parseIntParam 解析整数参数
func parseIntParam(c *app.RequestContext, key string, defaultValue int) int {
	strValue := c.Query(key)
	if strValue == "" {
		return defaultValue
	}

	var value int
	_, err := c.Query(key).Scan(&value)
	if err != nil {
		return defaultValue
	}
	return value
}
