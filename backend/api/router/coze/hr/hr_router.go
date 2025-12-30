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

package hr

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/coze-dev/coze-studio/backend/api/handler/coze/hr"
	"github.com/coze-dev/coze-studio/backend/api/middleware"
)

// RegisterHRRoutes 注册HR生命周期管理的所有路由
//
// **功能**：注册合同管理、调岗管理、离职管理相关的所有API端点
//
// **路由结构**：
// - /api/hr/contracts/*     - 合同管理 (5个路由)
// - /api/hr/employees/*     - 员工调岗 (3个路由)
// - /api/hr/transfers/*     - 调岗记录 (1个路由)
// - /api/hr/resignations/*  - 离职管理 (6个路由)
//
// **中间件**：
// - 所有路由都应用 TenantIsolationMiddleware() 实现租户隔离
// - tenant_id从context中自动提取，无需手动传递
//
// **权限要求**：
// - 所有接口都需要认证
// - 合同创建和签署需要HR权限
// - 调岗和离职审批需要管理员权限
func RegisterHRRoutes(r *server.Hertz, handler *hr.HRLifecycleHandler) {
	// 创建根路由组，应用租户隔离中间件
	apiGroup := r.Group("/api", middleware.TenantIsolationMiddleware())
	{
		// ============================================================
		// HR管理 (15个路由)
		// ============================================================
		hrGroup := apiGroup.Group("/hr")
		{
			// ---------------------------- 合同管理 (5个路由) ----------------------------
			contracts := hrGroup.Group("/contracts")
			{
				contracts.POST("", handler.CreateContract)              // 创建合同
				contracts.GET("/:id", handler.GetContract)              // 获取合同详情
				contracts.POST("/:id/sign", handler.SignContract)       // 签署合同
			}

			employees := hrGroup.Group("/employees")
			{
				employees.GET("/:emp_id/contracts", handler.GetEmployeeContracts)   // 获取员工的所有合同
				employees.GET("/:emp_id/contracts/active", handler.GetActiveContract) // 获取员工的生效合同
			}

			// ---------------------------- 调岗管理 (3个路由) ----------------------------
			empGroup := hrGroup.Group("/employees")
			{
				empGroup.POST("/:emp_id/transfer", handler.TransferEmployee)    // 调岗
				empGroup.GET("/:emp_id/transfers", handler.GetEmployeeTransfers) // 获取员工的调岗记录
			}

			transfers := hrGroup.Group("/transfers")
			{
				transfers.GET("/:id", handler.GetTransfer) // 获取调岗记录详情
			}

			// ---------------------------- 离职管理 (6个路由) ----------------------------
			empResignGroup := hrGroup.Group("/employees")
			{
				empResignGroup.POST("/:emp_id/resign", handler.ResignEmployee)         // 员工离职
				empResignGroup.GET("/:emp_id/resignation", handler.GetEmployeeResignation) // 获取员工的离职记录
			}

			resignations := hrGroup.Group("/resignations")
			{
				resignations.GET("/:id", handler.GetResignation)               // 获取离职记录详情
				resignations.POST("/:id/approve", handler.ApproveResignation)  // 审批离职
				resignations.PUT("/:id/handover", handler.UpdateHandoverStatus) // 更新交接状态
				resignations.GET("/pending", handler.GetPendingResignations)   // 获取待审批的离职列表
			}

			// ---------------------------- 试用期管理 (1个路由) ----------------------------
			probation := hrGroup.Group("/employees")
			{
				probation.GET("/probation/upcoming", handler.GetUpcomingProbationEndings) // 获取即将结束试用期的员工列表
			}
		}
	}
}
