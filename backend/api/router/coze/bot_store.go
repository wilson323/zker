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

// RegisterBotStoreRoutes 注册Bot商店路由
func RegisterBotStoreRoutes(r *server.Hertz) {
	// Bot商店API组
	botStoreGroup := r.Group("/api/v1/bot-store")
	{
		// 发布相关（需要认证）
		botStoreGroup.POST("/publish", PublishBotToStore)
		botStoreGroup.POST("/:item_id/unpublish", UnpublishBot)
		botStoreGroup.PUT("/:item_id", UpdateBotStoreItem)

		// 浏览相关（公开访问）
		botStoreGroup.GET("/list", ListBotStoreItems)
		botStoreGroup.GET("/search", SearchBotStoreItems)
		botStoreGroup.GET("/categories", GetBotCategories)
		botStoreGroup.GET("/:item_id", GetBotStoreItem)

		// 评论相关（需要认证）
		botStoreGroup.POST("/:item_id/reviews", CreateReview)               // 创建评论
		botStoreGroup.GET("/:item_id/reviews", GetReviews)                  // 获取评论列表
		botStoreGroup.GET("/:item_id/reviews/statistics", GetReviewStatistics) // 获取评论统计
		botStoreGroup.PUT("/:item_id/reviews/:review_id", UpdateReview)     // 更新评论
		botStoreGroup.DELETE("/:item_id/reviews/:review_id", DeleteReview)  // 删除评论
	}

	// Bot商店管理API组（需要管理员权限）
	botStoreAdminGroup := r.Group("/api/v1/bot-store/admin")
	{
		botStoreAdminGroup.GET("/pending", GetPendingReviews)
		botStoreAdminGroup.POST("/:item_id/review", ReviewBotStoreItem)
	}
}
