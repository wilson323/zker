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
	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

var (
	botStorePublisher  service.BotStorePublisher
	botStoreBrowser   service.BotStoreBrowser
	botStoreReviewer  service.BotStoreReviewer
	botStoreReviewSvc service.BotStoreReviewService
)

// InitBotStoreServices 初始化Bot商店服务
func InitBotStoreServices(
	publisher service.BotStorePublisher,
	browser service.BotStoreBrowser,
	reviewer service.BotStoreReviewer,
	reviewSvc service.BotStoreReviewService,
) {
	botStorePublisher = publisher
	botStoreBrowser = browser
	botStoreReviewer = reviewer
	botStoreReviewSvc = reviewSvc
}

// PublishBotToStore 发布Bot到商店
// @router /api/v1/bot-store/publish [POST]
func PublishBotToStore(ctx context.Context, c *app.RequestContext) {
	var req botstoreAPI.PublishBotToStoreRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, err.Error(), "参数验证失败", nil))
		return
	}

	// 获取用户ID和租户ID
	userID := strconv.FormatInt(*ctxutil.GetUIDFromCtx(ctx), 10)

	// 获取租户ID (需要从上下文中获取,暂时使用空字符串)
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
		httputil.BuildErrorRespFromEnhanced(c, errno.NewInternalError(ctx, err))
		return
	}

	// 转换为DTO
	dto := entityToDTO(item)

	httputil.BuildSuccessResp(c, dto)
}

// UnpublishBot 下架Bot
// @router /api/v1/bot-store/:item_id/unpublish [POST]
func UnpublishBot(ctx context.Context, c *app.RequestContext) {
	itemID := c.Param("item_id")
	if itemID == "" {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, "item_id is required", "参数验证失败", nil)
		return
	}

	userID := getUserIDFromContext(ctx)

	if err := botStorePublisher.UnpublishBot(ctx, itemID, userID); err != nil {
		httputil.BuildErrorRespFromEnhanced(c, errno.NewInternalError(ctx, err))
		return
	}

	httputil.BuildSuccessResp(c, nil)
}

// UpdateBotStoreItem 更新商店项目
// @router /api/v1/bot-store/:item_id [PUT]
func UpdateBotStoreItem(ctx context.Context, c *app.RequestContext) {
	itemID := c.Param("item_id")
	if itemID == "" {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, "item_id is required", "参数验证失败", nil)
		return
	}

	var req botstoreAPI.UpdateBotStoreItemRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, err.Error(), "参数验证失败", nil))
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
		httputil.BuildErrorRespFromEnhanced(c, errno.NewInternalError(ctx, err))
		return
	}

httputil.BuildSuccessResp(c, nil)
}

// ListBotStoreItems 列出商店项目
// @router /api/v1/bot-store/list [GET]
func ListBotStoreItems(ctx context.Context, c *app.RequestContext) {
	var req botstoreAPI.ListBotStoreItemsRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, err.Error(), "参数验证失败", nil))
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
		httputil.BuildErrorRespFromEnhanced(c, errno.NewInternalError(ctx, err))
		return
	}

	// 转换为DTO
	items := make([]*botstoreAPI.BotStoreItemDTO, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, entityToDTO(item))
	}

httputil.BuildSuccessResp(c, &botstoreAPI.BotStoreItemListData{
		Items:      items,
		Total:      resp.Total,
		Page:       resp.Page,
		PageSize:   resp.PageSize,
		TotalPages: resp.TotalPages,
	})
}

// SearchBotStoreItems 搜索商店项目
// @router /api/v1/bot-store/search [GET]
func SearchBotStoreItems(ctx context.Context, c *app.RequestContext) {
	var req botstoreAPI.SearchBotStoreItemsRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, err.Error(), "参数验证失败", nil))
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
		httputil.BuildErrorRespFromEnhanced(c, errno.NewInternalError(ctx, err))
		return
	}

	// 转换为DTO
	items := make([]*botstoreAPI.BotStoreItemDTO, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, entityToDTO(item))
	}

httputil.BuildSuccessResp(c, &botstoreAPI.BotStoreItemListData{
		Items:      items,
		Total:      resp.Total,
		Page:       resp.Page,
		PageSize:   resp.PageSize,
		TotalPages: resp.TotalPages,
	})
}

// GetBotStoreItem 获取商店项目详情
// @router /api/v1/bot-store/:item_id [GET]
func GetBotStoreItem(ctx context.Context, c *app.RequestContext) {
	itemID := c.Param("item_id")
	if itemID == "" {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, "item_id is required", "参数验证失败", nil)
		return
	}

	item, err := botStoreBrowser.GetBotStoreItem(ctx, itemID)
	if err != nil {
		httputil.BuildErrorRespFromEnhanced(c, errno.NewInternalError(ctx, err))
		return
	}

	// 记录浏览
	_ = botStoreBrowser.RecordView(ctx, itemID)

httputil.BuildSuccessResp(c, entityToDTO(item))
}

// GetBotCategories 获取分类列表
// @router /api/v1/bot-store/categories [GET]
func GetBotCategories(ctx context.Context, c *app.RequestContext) {
	categories, err := botStoreBrowser.GetCategories(ctx)
	if err != nil {
		httputil.BuildErrorRespFromEnhanced(c, errno.NewInternalError(ctx, err))
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

httputil.BuildSuccessResp(c, dtos)
}

// GetPendingReviews 获取待审核列表（管理员）
// @router /api/v1/bot-store/admin/pending [GET]
func GetPendingReviews(ctx context.Context, c *app.RequestContext) {
	var req botstoreAPI.GetPendingReviewsRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, err.Error(), "参数验证失败", nil))
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
		httputil.BuildErrorRespFromEnhanced(c, errno.NewInternalError(ctx, err))
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

httputil.BuildSuccessResp(c, &botstoreAPI.BotStoreItemListData{
		Items:      dtos,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	})
}

// ReviewBotStoreItem 审核Bot商店项目（管理员）
// @router /api/v1/bot-store/admin/:item_id/review [POST]
func ReviewBotStoreItem(ctx context.Context, c *app.RequestContext) {
	itemID := c.Param("item_id")
	if itemID == "" {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, "item_id is required", "参数验证失败", nil)
		return
	}

	var req botstoreAPI.ReviewBotStoreItemRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, err.Error(), "参数验证失败", nil))
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
		httputil.BuildErrorRespFromEnhanced(c, errno.NewInternalError(ctx, err))
		return
	}

httputil.BuildSuccessResp(c, nil)
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
	uid := ctxutil.GetUIDFromCtx(ctx)
	if uid == nil {
		return ""
	}
	return strconv.FormatInt(*uid, 10)
}

// getTenantIDFromContext 从上下文获取租户ID
func getTenantIDFromContext(ctx context.Context) string {
	// TODO: 实现从租户中间件获取租户ID
	// 当前返回空字符串，待实现租户上下文管理器后完善
	return ""
}


// ========== Bot商店评论相关Handler ==========

// CreateReview 创建评论
// @router /api/v1/bot-store/:item_id/reviews [POST]
func CreateReview(ctx context.Context, c *app.RequestContext) {
	itemID := c.Param("item_id")
	if itemID == "" {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, "item_id is required", "参数验证失败", nil)
		return
	}

	var req botstoreAPI.CreateReviewRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, err.Error(), "参数验证失败", nil))
		return
	}

	req.ItemID = itemID

	// 获取用户ID和租户ID
	userID := getUserIDFromContext(ctx)
	tenantID := getTenantIDFromContext(ctx)

	// 调用服务
	createReq := &service.CreateReviewRequest{
		ItemID:  req.ItemID,
		Rating:  req.Rating,
		Comment: req.Comment,
	}

	review, err := botStoreReviewSvc.CreateReview(ctx, createReq, userID, tenantID)
	if err != nil {
		handleReviewError(ctx, c, err)
		return
	}

httputil.BuildSuccessResp(c, reviewToDTO(review))
}

// GetReviews 获取评论列表
// @router /api/v1/bot-store/:item_id/reviews [GET]
func GetReviews(ctx context.Context, c *app.RequestContext) {
	itemID := c.Param("item_id")
	if itemID == "" {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, "item_id is required", "参数验证失败", nil)
		return
	}

	// 设置默认值
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	sortBy := c.Query("sort_by")
	if sortBy == "" {
		sortBy = "latest"
	}

	// 调用服务
	req := &service.ListReviewsRequest{
		ItemID:   itemID,
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
	}

	resp, err := botStoreReviewSvc.ListReviews(ctx, req)
	if err != nil {
		handleReviewError(ctx, c, err)
		return
	}

	// 转换为DTO
	reviews := make([]*botstoreAPI.BotStoreReviewDTO, 0, len(resp.Reviews))
	for _, review := range resp.Reviews {
		reviews = append(reviews, reviewToDTO(review))
	}

httputil.BuildSuccessResp(c, &botstoreAPI.ReviewListData{
		Reviews:       reviews,
		Total:         resp.Total,
		Page:          resp.Page,
		PageSize:      resp.PageSize,
		TotalPages:    resp.TotalPages,
		AverageRating: resp.AverageRating,
		RatingCount:   resp.RatingCount,
	})
}

// UpdateReview 更新评论
// @router /api/v1/bot-store/:item_id/reviews/:review_id [PUT]
func UpdateReview(ctx context.Context, c *app.RequestContext) {
	reviewID := c.Param("review_id")
	if reviewID == "" {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, "review_id is required", "参数验证失败", nil)
		return
	}

	var req botstoreAPI.UpdateReviewRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, err.Error(), "参数验证失败", nil))
		return
	}

	userID := getUserIDFromContext(ctx)

	// 调用服务
	updateReq := &service.UpdateReviewRequest{
		Rating:  req.Rating,
		Comment: req.Comment,
	}

	if err := botStoreReviewSvc.UpdateReview(ctx, reviewID, updateReq, userID); err != nil {
		handleReviewError(ctx, c, err)
		return
	}

httputil.BuildSuccessResp(c, nil)
}

// DeleteReview 删除评论
// @router /api/v1/bot-store/:item_id/reviews/:review_id [DELETE]
func DeleteReview(ctx context.Context, c *app.RequestContext) {
	reviewID := c.Param("review_id")
	if reviewID == "" {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, "review_id is required", "参数验证失败", nil)
		return
	}

	userID := getUserIDFromContext(ctx)

	if err := botStoreReviewSvc.DeleteReview(ctx, reviewID, userID); err != nil {
		handleReviewError(ctx, c, err)
		return
	}

httputil.BuildSuccessResp(c, nil)
}

// GetReviewStatistics 获取评论统计
// @router /api/v1/bot-store/:item_id/reviews/statistics [GET]
func GetReviewStatistics(ctx context.Context, c *app.RequestContext) {
	itemID := c.Param("item_id")
	if itemID == "" {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, "item_id is required", "参数验证失败", nil)
		return
	}

	stats, err := botStoreReviewSvc.GetReviewStatistics(ctx, itemID)
	if err != nil {
		handleReviewError(ctx, c, err)
		return
	}

httputil.BuildSuccessResp(c, &botstoreAPI.ReviewStatisticsDTO{
		AverageRating: stats.AverageRating,
		RatingCount:   stats.RatingCount,
		Rating1Count:  stats.Rating1Count,
		Rating2Count:  stats.Rating2Count,
		Rating3Count:  stats.Rating3Count,
		Rating4Count:  stats.Rating4Count,
		Rating5Count:  stats.Rating5Count,
	})
}

// reviewToDTO 实体转DTO
func reviewToDTO(review *entity.BotStoreReview) *botstoreAPI.BotStoreReviewDTO {
	return &botstoreAPI.BotStoreReviewDTO{
		ReviewID:  review.ReviewID,
		ItemID:    review.ItemID,
		UserID:    review.UserID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt.Format(time.RFC3339),
		UpdatedAt: review.UpdatedAt.Format(time.RFC3339),
	}
}

// handleReviewError 处理评论错误
func handleReviewError(ctx context.Context, c *app.RequestContext, err error) {
	// 根据错误类型返回不同的HTTP状态码
	// 这里简化处理，实际应该根据errno返回对应的状态码
	httputil.InternalError(ctx, c, err)
}
