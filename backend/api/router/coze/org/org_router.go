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
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/coze-dev/coze-studio/backend/api/handler/coze/org"
	"github.com/coze-dev/coze-studio/backend/api/middleware"
)

// OrgHandlers 组织管理相关的所有Handler集合
type OrgHandlers struct {
	Organization *org.OrganizationHandler
	Department   *org.DepartmentHandler
	Position     *org.PositionHandler
	// Employee和Directory使用函数式Handler，需要适配器
}

// EmployeeHandlerAdapter 员工Handler适配器（将函数式API转换为Handler方法）
type EmployeeHandlerAdapter struct{}

func (a *EmployeeHandlerAdapter) CreateEmployee(ctx context.Context, c *app.RequestContext) {
	org.CreateEmployee(ctx, c)
}

func (a *EmployeeHandlerAdapter) GetEmployee(ctx context.Context, c *app.RequestContext) {
	org.GetEmployee(ctx, c)
}

func (a *EmployeeHandlerAdapter) UpdateEmployee(ctx context.Context, c *app.RequestContext) {
	org.UpdateEmployee(ctx, c)
}

func (a *EmployeeHandlerAdapter) DeleteEmployee(ctx context.Context, c *app.RequestContext) {
	org.DeleteEmployee(ctx, c)
}

func (a *EmployeeHandlerAdapter) ActivateEmployee(ctx context.Context, c *app.RequestContext) {
	org.UpdateEmployeeStatus(ctx, c)
}

func (a *EmployeeHandlerAdapter) DeactivateEmployee(ctx context.Context, c *app.RequestContext) {
	org.UpdateEmployeeStatus(ctx, c)
}

func (a *EmployeeHandlerAdapter) GetEmployeeStatus(ctx context.Context, c *app.RequestContext) {
	org.GetEmployee(ctx, c)
}

func (a *EmployeeHandlerAdapter) ListEmployees(ctx context.Context, c *app.RequestContext) {
	org.ListEmployees(ctx, c)
}

func (a *EmployeeHandlerAdapter) GetEmployeeByCode(ctx context.Context, c *app.RequestContext) {
	org.GetEmployeeByCode(ctx, c)
}

func (a *EmployeeHandlerAdapter) SearchEmployees(ctx context.Context, c *app.RequestContext) {
	org.SearchEmployees(ctx, c)
}

func (a *EmployeeHandlerAdapter) GetEmployeesByDepartment(ctx context.Context, c *app.RequestContext) {
	org.GetEmployeesByDepartment(ctx, c)
}

func (a *EmployeeHandlerAdapter) GetEmployeesByPosition(ctx context.Context, c *app.RequestContext) {
	// 使用GetEmployeesByDepartment的类似逻辑
	org.GetEmployeesByDepartment(ctx, c)
}

// DirectoryHandlerAdapter 通讯录Handler适配器（将函数式API转换为Handler方法）
type DirectoryHandlerAdapter struct{}

func (a *DirectoryHandlerAdapter) GetDirectory(ctx context.Context, c *app.RequestContext) {
	org.GetOrganizationDirectory(ctx, c)
}

func (a *DirectoryHandlerAdapter) GetDirectoryTree(ctx context.Context, c *app.RequestContext) {
	org.GetOrganizationDirectory(ctx, c)
}

func (a *DirectoryHandlerAdapter) GetDirectoryByOrg(ctx context.Context, c *app.RequestContext) {
	org.GetOrganizationEmployees(ctx, c)
}

func (a *DirectoryHandlerAdapter) GetDirectoryByDept(ctx context.Context, c *app.RequestContext) {
	org.GetDepartmentEmployees(ctx, c)
}

func (a *DirectoryHandlerAdapter) ExportDirectory(ctx context.Context, c *app.RequestContext) {
	// TODO: 实现导出功能
	c.JSON(200, map[string]interface{}{"code": 0, "message": "TODO: Export directory"})
}

func (a *DirectoryHandlerAdapter) ImportDirectory(ctx context.Context, c *app.RequestContext) {
	// TODO: 实现导入功能
	c.JSON(200, map[string]interface{}{"code": 0, "message": "TODO: Import directory"})
}

// RegisterOrgRoutes 注册组织中心的所有路由
//
// **功能**：注册组织、部门、岗位、员工、通讯录相关的所有API端点
//
// **路由结构**：
// - /api/organizations/*    - 组织管理 (10个路由)
// - /api/org/departments/*  - 部门管理 (11个路由)
// - /api/org/positions/*    - 岗位管理 (9个路由)
// - /api/org/employees/*    - 员工管理 (12个路由)
// - /api/org/directory/*    - 通讯录 (6个路由)
//
// **中间件**：
// - 所有路由都应用 TenantIsolationMiddleware() 实现租户隔离
// - tenant_id从context中自动提取，无需手动传递
//
// **权限要求**：
// - 所有接口都需要认证
// - 部分接口需要管理员权限
func RegisterOrgRoutes(r *server.Hertz, handlers *OrgHandlers) {
	// 创建适配器
	employeeAdapter := &EmployeeHandlerAdapter{}
	directoryAdapter := &DirectoryHandlerAdapter{}

	// 创建根路由组，应用租户隔离中间件
	apiGroup := r.Group("/api", middleware.TenantIsolationMiddleware())
	{
		// ============================================================
		// 组织管理 (10个路由)
		// ============================================================
		organizations := apiGroup.Group("/organizations")
		{
			// 基础CRUD操作 (4个路由)
			organizations.POST("", handlers.Organization.CreateOrganization)       // 创建组织
			organizations.GET("/:id", handlers.Organization.GetOrganization)       // 获取组织详情
			organizations.PUT("/:id", handlers.Organization.UpdateOrganization)    // 更新组织
			organizations.DELETE("/:id", handlers.Organization.DeleteOrganization) // 删除组织

			// 组织树操作 (2个路由)
			organizations.GET("/tree", handlers.Organization.GetOrganizationTree) // 获取组织树
			organizations.GET("", handlers.Organization.ListOrganizations)        // 分页查询组织列表

			// 组织关系查询 (4个路由)
			organizations.POST("/:id/move", handlers.Organization.MoveOrganization)    // 移动组织
			organizations.GET("/:id/children", handlers.Organization.GetChildren)      // 获取子组织
			organizations.GET("/:id/ancestors", handlers.Organization.GetAncestors)    // 获取祖先组织
			organizations.GET("/:id/descendants", handlers.Organization.GetDescendants) // 获取后代组织
		}

		// ============================================================
		// 部门管理 (11个路由)
		// ============================================================
		departments := apiGroup.Group("/org")
		{
			deptGroup := departments.Group("/departments")
			{
				// 基础CRUD操作 (4个路由)
				deptGroup.POST("", handlers.Department.CreateDepartment)       // 创建部门
				deptGroup.GET("/:id", handlers.Department.GetDepartment)       // 获取部门详情
				deptGroup.PUT("/:id", handlers.Department.UpdateDepartment)    // 更新部门
				deptGroup.DELETE("/:id", handlers.Department.DeleteDepartment) // 删除部门

				// 部门树操作 (2个路由)
				deptGroup.GET("/tree", handlers.Department.GetDepartmentTree) // 获取部门树
				deptGroup.GET("", handlers.Department.ListDepartments)        // 分页查询部门列表

				// 部门关系操作 (4个路由)
				deptGroup.POST("/:id/move", handlers.Department.MoveDepartment)       // 移动部门
				deptGroup.GET("/:id/children", handlers.Department.GetChildren)       // 获取子部门
				deptGroup.GET("/:id/ancestors", handlers.Department.GetAncestors)     // 获取祖先部门
				deptGroup.GET("/:id/descendants", handlers.Department.GetDescendants) // 获取后代部门

				// 按组织查询部门 (1个路由)
				deptGroup.GET("/org/:org_id", handlers.Department.GetDepartmentsByOrg) // 获取组织的所有部门
			}
		}

		// ============================================================
		// 岗位管理 (9个路由)
		// ============================================================
		positions := apiGroup.Group("/org")
		{
			positionGroup := positions.Group("/positions")
			{
				// 基础CRUD操作 (4个路由)
				positionGroup.POST("", handlers.Position.CreatePosition)       // 创建岗位
				positionGroup.GET("/:id", handlers.Position.GetPosition)       // 获取岗位详情
				positionGroup.PUT("/:id", handlers.Position.UpdatePosition)    // 更新岗位
				positionGroup.DELETE("/:id", handlers.Position.DeletePosition) // 删除岗位

				// 岗位查询操作 (4个路由)
				positionGroup.GET("", handlers.Position.ListPositions)                 // 分页查询岗位列表
				positionGroup.GET("/code/:code", handlers.Position.GetPositionByCode) // 按岗位编码查询
				positionGroup.GET("/dept/:dept_id", handlers.Position.GetPositionsByDepartment) // 按部门查询岗位
				positionGroup.GET("/level/:level", handlers.Position.GetPositionsByLevel)       // 按职级查询岗位

				// 按类别查询岗位 (1个路由)
				positionGroup.GET("/category/:category", handlers.Position.GetPositionsByCategory) // 按类别查询岗位
			}
		}

		// ============================================================
		// 员工管理 (12个路由)
		// ============================================================
		employees := apiGroup.Group("/org")
		{
			empGroup := employees.Group("/employees")
			{
				// 基础CRUD操作 (4个路由)
				empGroup.POST("", employeeAdapter.CreateEmployee)       // 创建员工
				empGroup.GET("/:id", employeeAdapter.GetEmployee)       // 获取员工详情
				empGroup.PUT("/:id", employeeAdapter.UpdateEmployee)    // 更新员工
				empGroup.DELETE("/:id", employeeAdapter.DeleteEmployee) // 删除员工

				// 员工查询操作 (5个路由)
				empGroup.GET("", employeeAdapter.ListEmployees)                   // 分页查询员工列表
				empGroup.GET("/code/:code", employeeAdapter.GetEmployeeByCode)   // 按员工编号查询
				empGroup.GET("/search", employeeAdapter.SearchEmployees)          // 搜索员工
				empGroup.GET("/dept/:dept_id", employeeAdapter.GetEmployeesByDepartment) // 按部门查询员工
				empGroup.GET("/position/:position_id", employeeAdapter.GetEmployeesByPosition) // 按岗位查询员工

				// 员工状态管理 (3个路由)
				empGroup.POST("/:id/activate", employeeAdapter.ActivateEmployee)   // 激活员工
				empGroup.POST("/:id/deactivate", employeeAdapter.DeactivateEmployee) // 停用员工
				empGroup.GET("/:id/status", employeeAdapter.GetEmployeeStatus)      // 获取员工状态
			}
		}

		// ============================================================
		// 通讯录 (6个路由)
		// ============================================================
		directory := apiGroup.Group("/org")
		{
			dirGroup := directory.Group("/directory")
			{
				// 通讯录查询 (6个路由)
				dirGroup.GET("", directoryAdapter.GetDirectory)          // 获取通讯录列表
				dirGroup.GET("/tree", directoryAdapter.GetDirectoryTree) // 获取通讯录树
				dirGroup.GET("/org/:org_id", directoryAdapter.GetDirectoryByOrg) // 按组织查询通讯录
				dirGroup.GET("/dept/:dept_id", directoryAdapter.GetDirectoryByDept) // 按部门查询通讯录
				dirGroup.GET("/export", directoryAdapter.ExportDirectory) // 导出通讯录
				dirGroup.POST("/import", directoryAdapter.ImportDirectory) // 导入通讯录
			}
		}
	}
}
