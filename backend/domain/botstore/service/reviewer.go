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

package service

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
)

// ReviewBotRequest 审核请求
type ReviewBotRequest struct {
	ItemID     string `json:"item_id"`
	Approved   bool   `json:"approved"`
	Reason     string `json:"reason"`
	ReviewerID string `json:"reviewer_id"`
}

// BotStoreReviewer Bot商店审核器接口
type BotStoreReviewer interface {
	// ReviewBot 审核Bot
	ReviewBot(ctx context.Context, req *ReviewBotRequest) error

	// GetPendingReviews 获取待审核列表
	GetPendingReviews(ctx context.Context, page, pageSize int) ([]*entity.BotStoreItem, int, error)

	// GetReviewStatistics 获取审核统计
	GetReviewStatistics(ctx context.Context) (*ReviewStatistics, error)
}

// ReviewStatistics 审核统计
type ReviewStatistics struct {
	TotalPending   int `json:"total_pending"`
	TotalPublished int `json:"total_published"`
	TotalRejected  int `json:"total_rejected"`
	TotalDraft     int `json:"total_draft"`
}
