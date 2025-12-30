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
	"strconv"
	"time"

	botstoreAPI "github.com/coze-dev/coze-studio/backend/api/model/botstore"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

var (
	botStorePublisher service.BotStorePublisher
	botStoreBrowser  service.BotStoreBrowser
	botStoreReviewer service.BotStoreReviewer
)

// InitBotStoreServices 初始化Bot商店服务
func InitBotStoreServices(
	publisher service.BotStorePublisher,
	browser service.BotStoreBrowser,
	reviewer service.BotStoreReviewer,
) {
	botStorePublisher = publisher
	botStoreBrowser = browser
	botStoreReviewer = reviewer
}

// PublishBotToStore 发布Bot到商店
// @router /api/v1/bot-store/publish [POST]
func PublishBotToStore(ctx context.Context, c *app.RequestContext) {
	var req botstoreAPI.PublishBotToStoreRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	// 获取用户ID和租户ID
	userID := getUserIDFromContext(ctx)
	tenantID := getTenantIDFromContext(ctx)

	// 调用发布服务
	publishReq := &service.PublishBotRequest{
		BotID:       req.BotID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Tags:        req.Tags,
		Price:       req.Price,
		Screenshots: req.Screenshots,
		Version:     req.Version,
		PublisherID: userID,
		TenantID:    tenantID,
	}

	item, err := botStorePublisher.PublishBot(ctx, publishReq)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换为DTO
	dto := entityToDTO(item)

	c.JSON(consts.StatusOK, &botstoreAPI.PublishBotToStoreResponse{
		BaseResponse: botstoreAPI.BaseResponse{
			Code: 0,
			Msg:  "success",
		},
		Data: dto,
	})
}

// UnpublishBot 下架Bot
// @router /api/v1/bot-store/:item_id/unpublish [POST]
func UnpublishBot(ctx context.Context, c *app.RequestContext) {
	itemID := c.Param("item_id")
	if itemID == "" {
		invalidParamRequestResponse(c, "item_id is required")
		return
	}

	userID := getUserIDFromContext(ctx)

	if err := botStorePublisher.UnpublishBot(ctx, itemID, userID); err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, &botstoreAPI.UnpublishBotResponse{
		BaseResponse: botstoreAPI.BaseResponse{
			Code: 0,
			Msg:  "success",
		},
	})
}

// UpdateBotStoreItem 更新商店项目
// @router /api/v1/bot-store/:item_id [PUT]
func UpdateBotStoreItem(ctx context.Context, c *app.RequestContext) {
	itemID := c.Param("item_id")
	if itemID == "" {
		invalidParamRequestResponse(c, "item_id is required")
		return
	}

	var req botstoreAPI.UpdateBotStoreItemRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	req.ItemID = itemID
	userID := getUserIDFromContext(ctx)

	updateReq := &service.UpdateBotStoreItemRequest{
		ItemID:      req.ItemID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Tags:        req.Tags,
		Price:       req.Price,
		Screenshots: req.Screenshots,
		Version:     req.Version,
	}

	if err := botStorePublisher.UpdateBotStoreItem(ctx, updateReq, userID); err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, &botstoreAPI.UpdateBotStoreItemResponse{
		BaseResponse: botstoreAPI.BaseResponse{
			Code: 0,
			Msg:  "success",
		},
	})
}

// ListBotStoreItems 列出商店项目
// @router /api/v1/bot-store/list [GET]
func ListBotStoreItems(ctx context.Context, c *app.RequestContext) {
	var req botstoreAPI.ListBotStoreItemsRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	listReq := &service.ListBotsRequest{
		Category: req.Category,
		SortBy:   req.SortBy,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	resp, err := botStoreBrowser.ListBots(ctx, listReq)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换为DTO
	items := make([]*botstoreAPI.BotStoreItemDTO, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, entityToDTO(item))
	}

	c.JSON(consts.StatusOK, &botstoreAPI.ListBotStoreItemsResponse{
		BaseResponse: botstoreAPI.BaseResponse{
			Code: 0,
			Msg:  "success",
		},
		Data: &botstoreAPI.BotStoreItemListData{
			Items:      items,
			Total:      resp.Total,
			Page:       resp.Page,
			PageSize:   resp.PageSize,
			TotalPages: resp.TotalPages,
		},
	})
}

// SearchBotStoreItems 搜索商店项目
// @router /api/v1/bot-store/search [GET]
func SearchBotStoreItems(ctx context.Context, c *app.RequestContext) {
	var req botstoreAPI.SearchBotStoreItemsRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	searchReq := &service.SearchBotsRequest{
		Query:    req.Query,
		Category: req.Category,
		Tags:     req.Tags,
		PriceMin: req.PriceMin,
		PriceMax: req.PriceMax,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	resp, err := botStoreBrowser.SearchBots(ctx, searchReq)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换为DTO
	items := make([]*botstoreAPI.BotStoreItemDTO, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, entityToDTO(item))
	}

	c.JSON(consts.StatusOK, &botstoreAPI.SearchBotStoreItemsResponse{
		BaseResponse: botstoreAPI.BaseResponse{
			Code: 0,
			Msg:  "success",
		},
		Data: &botstoreAPI.BotStoreItemListData{
			Items:      items,
			Total:      resp.Total,
			Page:       resp.Page,
			PageSize:   resp.PageSize,
			TotalPages: resp.TotalPages,
		},
	})
}

// GetBotStoreItem 获取商店项目详情
// @router /api/v1/bot-store/:item_id [GET]
func GetBotStoreItem(ctx context.Context, c *app.RequestContext) {
	itemID := c.Param("item_id")
	if itemID == "" {
		invalidParamRequestResponse(c, "item_id is required")
		return
	}

	item, err := botStoreBrowser.GetBotStoreItem(ctx, itemID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 记录浏览
	_ = botStoreBrowser.RecordView(ctx, itemID)

	c.JSON(consts.StatusOK, &botstoreAPI.GetBotStoreItemResponse{
		BaseResponse: botstoreAPI.BaseResponse{
			Code: 0,
			Msg:  "success",
		},
		Data: entityToDTO(item),
	})
}

// GetBotCategories 获取分类列表
// @router /api/v1/bot-store/categories [GET]
func GetBotCategories(ctx context.Context, c *app.RequestContext) {
	categories, err := botStoreBrowser.GetCategories(ctx)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换为DTO
	dtos := make([]*botstoreAPI.BotStoreCategoryDTO, 0, len(categories))
	for _, category := range categories {
		dtos = append(dtos, &botstoreAPI.BotStoreCategoryDTO{
			CategoryID:  category.CategoryID,
			Name:        category.Name,
			Icon:        category.Icon,
			Description: category.Description,
			SortOrder:   category.SortOrder,
			BotCount:    category.BotCount,
		})
	}

	c.JSON(consts.StatusOK, &botstoreAPI.GetBotCategoriesResponse{
		BaseResponse: botstoreAPI.BaseResponse{
			Code: 0,
			Msg:  "success",
		},
		data: dtos,
	})
}

// GetPendingReviews 获取待审核列表（管理员）
// @router /api/v1/bot-store/admin/pending [GET]
func GetPendingReviews(ctx context.Context, c *app.RequestContext) {
	var req botstoreAPI.GetPendingReviewsRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	items, total, err := botStoreReviewer.GetPendingReviews(ctx, req.Page, req.PageSize)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换为DTO
	dtos := make([]*botstoreAPI.BotStoreItemDTO, 0, len(items))
	for _, item := range items {
		dtos = append(dtos, entityToDTO(item))
	}

	totalPages := total / req.PageSize
	if total%req.PageSize > 0 {
		totalPages++
	}

	c.JSON(consts.StatusOK, &botstoreAPI.GetPendingReviewsResponse{
		BaseResponse: botstoreAPI.BaseResponse{
			Code: 0,
			Msg:  "success",
		},
		Data: &botstoreAPI.BotStoreItemListData{
			Items:      dtos,
			Total:      total,
			Page:       req.Page,
			PageSize:   req.PageSize,
			TotalPages: totalPages,
		},
	})
}

// ReviewBotStoreItem 审核Bot商店项目（管理员）
// @router /api/v1/bot-store/admin/:item_id/review [POST]
func ReviewBotStoreItem(ctx context.Context, c *app.RequestContext) {
	itemID := c.Param("item_id")
	if itemID == "" {
		invalidParamRequestResponse(c, "item_id is required")
		return
	}

	var req botstoreAPI.ReviewBotStoreItemRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	reviewerID := getUserIDFromContext(ctx)

	reviewReq := &service.ReviewBotRequest{
		ItemID:    itemID,
		Approved:  req.Approved,
		Reason:    req.Reason,
		ReviewerID: reviewerID,
	}

	if err := botStoreReviewer.ReviewBot(ctx, reviewReq); err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, &botstoreAPI.ReviewBotStoreItemResponse{
		BaseResponse: botstoreAPI.BaseResponse{
			Code: 0,
			Msg:  "success",
		},
	})
}

// entityToDTO 实体转DTO
func entityToDTO(item *entity.BotStoreItem) *botstoreAPI.BotStoreItemDTO {
	return &botstoreAPI.BotStoreItemDTO{
		ItemID:        item.ItemID,
		BotID:         item.BotID,
		Name:          item.Name,
		Description:   item.Description,
		Category:      item.Category,
		Tags:          item.Tags,
		Price:         item.Price,
		PublisherID:   item.PublisherID,
		Status:        string(item.Status),
		ViewCount:     item.ViewCount,
		DownloadCount: item.DownloadCount,
		Rating:        item.Rating,
		RatingCount:   item.RatingCount,
		Screenshots:   item.Screenshots,
		Version:       item.Version,
		CreatedAt:     item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     item.UpdatedAt.Format(time.RFC3339),
	}
}

// getUserIDFromContext 从上下文获取用户ID
func getUserIDFromContext(ctx context.Context) string {
	// TODO: 从实际的上下文中获取用户ID
	// 这里需要根据实际的认证中间件实现来获取
	return "user_id"
}

// getTenantIDFromContext 从上下文获取租户ID
func getTenantIDFromContext(ctx context.Context) string {
	// TODO: 从实际的上下文中获取租户ID
	// 这里需要根据实际的租户中间件实现来获取
	return "tenant_id"
}

// invalidParamRequestResponse 无效参数响应
func invalidParamRequestResponse(c *app.RequestContext, msg string) {
	c.JSON(consts.StatusBadRequest, &botstoreAPI.BaseResponse{
		Code: consts.StatusBadRequest,
		Msg:  msg,
	})
}

// internalServerErrorResponse 内部服务器错误响应
func internalServerErrorResponse(ctx context.Context, c *app.RequestContext, err error) {
	c.JSON(consts.StatusInternalServerError, &botstoreAPI.BaseResponse{
		Code: consts.StatusInternalServerError,
		Msg:  err.Error(),
	})
}
