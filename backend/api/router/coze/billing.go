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
	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterBillingRoutes 注册计费相关路由
func RegisterBillingRoutes(r *server.Hertz) {
	root := r.Group("/")
	{
		_api := root.Group("/api")
		{
			_billing := _api.Group("/billing")
			{
				// ==================== 实时计费路由 ====================
				_realtime := _billing.Group("/realtime")
				{
					_realtime.POST("/charge", Charge)           // 实时扣费
				}

				// ==================== 余额查询路由 ====================
				_balance := _billing.Group("/tenant")
				{
					_balance.GET("/:tenant_id/balance", GetBalance) // 查询余额
				}

				// ==================== 配额管理路由 ====================
				_quota := _billing.Group("/quota")
				{
					_quota.POST("/consume", ConsumeQuota)   // 消费配额
					_quota.GET("/check", CheckQuota)        // 检查配额
					_quota.POST("/reset", ResetQuota)       // 重置配额
					_quota.POST("/rollback", RollbackQuota) // 回滚配额
				}

				// ==================== 账单生成路由 ====================
				_bills := _billing.Group("/bills")
				{
					_bills.POST("/generate", GenerateBill)          // 生成账单
					_bills.GET("/:bill_id", GetBill)                // 查询账单详情
					_bills.GET("", ListBills)                       // 查询账单列表
					_bills.POST("/:bill_id/pay", PayBill)           // 支付账单
					_bills.GET("/:bill_id/invoice", GetInvoice)     // 获取发票
					_bills.GET("/:bill_id/invoice/pdf", DownloadInvoicePDF) // 下载发票PDF
				}
			}
		}
	}
}
