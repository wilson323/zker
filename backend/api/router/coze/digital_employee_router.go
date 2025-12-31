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
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// registerDigitalEmployeeRoutes 注册数字员工管理路由
func registerDigitalEmployeeRoutes(r *server.Hertz) {
	// 员工画像管理接口（需要认证和租户上下文）
	_employees := r.Group("/api/v1/digital-employees")
	{
		// 员工画像CRUD
		_employees.POST("", CreateEmployeeProfile)                         // 创建员工画像
		_employees.PUT("/:employee_id", UpdateEmployeeProfile)             // 更新员工画像
		_employees.GET("/:employee_id", GetEmployeeProfile)                // 获取员工画像
		_employees.GET("", ListEmployeeProfiles)                           // 列出员工画像
		_employees.DELETE("/:employee_id", DeleteEmployeeProfile)          // 删除员工画像

		// 任务分配管理
		_tasks := _employees.Group("/tasks")
		{
			_tasks.POST("/assign", AssignTask)                              // 分配任务给员工
			_tasks.POST("/auto-assign", AutoAssignTask)                     // 自动分配任务
			_tasks.PUT("/:assignment_id/complete", CompleteTask)            // 完成任务
			_tasks.PUT("/:assignment_id/fail", FailTask)                    // 标记任务失败
		}

		// 员工任务列表
		_employees.GET("/:employee_id/tasks", GetEmployeeTasks)            // 获取员工任务列表

		// 绩效统计
		_employees.GET("/:employee_id/performance", GetEmployeePerformance) // 获取员工绩效
	}

	// 团队绩效接口
	_performance := r.Group("/api/v1/digital-employees/performance")
	{
		_performance.GET("/team", GetTeamPerformance)                      // 获取团队绩效
	}
}
