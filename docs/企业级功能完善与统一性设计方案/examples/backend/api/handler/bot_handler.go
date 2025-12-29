// Package handler 提供HTTP处理器
//
// API层职责：
// - 处理HTTP请求和响应
// - 参数解析和验证
// - 调用应用服务
// - 错误处理和响应格式化
// - 请求日志和追踪
package handler

import (
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/coze-studio/application/dto"
	"github.com/coze-studio/application/service"
	"github.com/coze-studio/crossdomain/auth"
)

// BotHandler Bot HTTP处理器
type BotHandler struct {
	botService *service.BotApplicationService
}

// NewBotHandler 创建Bot处理器
func NewBotHandler(
	botService *service.BotApplicationService,
) *BotHandler {
	return &BotHandler{
		botService: botService,
	}
}

// ============ 路由注册 ============

// RegisterRoutes 注册路由
//
// 本方法展示RESTful API路由设计：
// - POST   /api/v1/bots          创建Bot
// - GET    /api/v1/bots/:id      获取Bot详情
// - PUT    /api/v1/bots/:id      更新Bot
// - DELETE /api/v1/bots/:id      删除Bot
// - POST   /api/v1/bots/:id/publish  发布Bot
// - GET    /api/v1/bots          查询Bot列表
func (h *BotHandler) RegisterRoutes(r *app.Server) {
	botGroup := r.Group("/api/v1/bots")
	{
		// 需要认证的路由
		botGroup.POST("", h.CreateBot)
		botGroup.GET("", h.ListBots)
		botGroup.GET("/:id", h.GetBot)
		botGroup.PUT("/:id", h.UpdateBot)
		botGroup.DELETE("/:id", h.DeleteBot)
		botGroup.POST("/:id/publish", h.PublishBot)
	}
}

// ============ 请求处理器 ============

// CreateBot 创建Bot
//
// @Summary      创建Bot
// @Description  创建一个新的AI对话Bot
// @Tags         bots
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateBotRequest true "创建请求"
// @Success      200 {object} Response{data=dto.BotResponse}
// @Failure      400 {object} Response
// @Failure      401 {object} Response
// @Failure      403 {object} Response
// @Failure      500 {object} Response
// @Router       /api/v1/bots [post]
func (h *BotHandler) CreateBot(ctx *app.RequestContext) {
	// 步骤1: 解析请求体
	var req dto.CreateBotRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		h.respondError(ctx, consts.StatusBadRequest, "invalid request body", err)
		return
	}

	// 步骤2: 从上下文获取用户信息（由认证中间件设置）
	currentUser := auth.GetUserContext(ctx)
	if currentUser == nil {
		h.respondError(ctx, consts.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// 步骤3: 调用应用服务
	bot, err := h.botService.CreateBot(ctx, &req, currentUser)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	// 步骤4: 返回成功响应
	h.respondSuccess(ctx, bot, "Bot created successfully")
}

// GetBot 获取Bot详情
//
// @Summary      获取Bot详情
// @Description  根据ID获取Bot的详细信息
// @Tags         bots
// @Accept       json
// @Produce      json
// @Param        id path string true "Bot ID"
// @Success      200 {object} Response{data=dto.BotResponse}
// @Failure      400 {object} Response
// @Failure      401 {object} Response
// @Failure      404 {object} Response
// @Router       /api/v1/bots/{id} [get]
func (h *BotHandler) GetBot(ctx *app.RequestContext) {
	// 步骤1: 解析路径参数
	botID := ctx.Param("id")
	if botID == "" {
		h.respondError(ctx, consts.StatusBadRequest, "bot ID is required", nil)
		return
	}

	// 步骤2: 获取用户信息
	currentUser := auth.GetUserContext(ctx)
	if currentUser == nil {
		h.respondError(ctx, consts.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// 步骤3: 调用应用服务（复用ListBots或创建新方法）
	// 这里假设有一个GetBot方法
	// bot, err := h.botService.GetBot(ctx, botID, currentUser)

	// 步骤4: 返回响应
	// h.respondSuccess(ctx, bot, "Bot retrieved successfully")
}

// UpdateBot 更新Bot
//
// @Summary      更新Bot
// @Description  更新Bot的配置信息
// @Tags         bots
// @Accept       json
// @Produce      json
// @Param        id path string true "Bot ID"
// @Param        request body dto.UpdateBotRequest true "更新请求"
// @Success      200 {object} Response{data=dto.BotResponse}
// @Failure      400 {object} Response
// @Failure      401 {object} Response
// @Failure      403 {object} Response
// @Failure      404 {object} Response
// @Router       /api/v1/bots/{id} [put]
func (h *BotHandler) UpdateBot(ctx *app.RequestContext) {
	// 步骤1: 解析路径参数
	botID := ctx.Param("id")
	if botID == "" {
		h.respondError(ctx, consts.StatusBadRequest, "bot ID is required", nil)
		return
	}

	// 步骤2: 解析请求体
	var req dto.UpdateBotRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		h.respondError(ctx, consts.StatusBadRequest, "invalid request body", err)
		return
	}

	// 步骤3: 获取用户信息
	currentUser := auth.GetUserContext(ctx)
	if currentUser == nil {
		h.respondError(ctx, consts.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// 步骤4: 调用应用服务
	bot, err := h.botService.UpdateBot(ctx, botID, &req, currentUser)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	// 步骤5: 返回成功响应
	h.respondSuccess(ctx, bot, "Bot updated successfully")
}

// DeleteBot 删除Bot
//
// @Summary      删除Bot
// @Description  删除指定的Bot（软删除）
// @Tags         bots
// @Accept       json
// @Produce      json
// @Param        id path string true "Bot ID"
// @Success      200 {object} Response
// @Failure      400 {object} Response
// @Failure      401 {object} Response
// @Failure      403 {object} Response
// @Failure      404 {object} Response
// @Router       /api/v1/bots/{id} [delete]
func (h *BotHandler) DeleteBot(ctx *app.RequestContext) {
	// 步骤1: 解析路径参数
	botID := ctx.Param("id")
	if botID == "" {
		h.respondError(ctx, consts.StatusBadRequest, "bot ID is required", nil)
		return
	}

	// 步骤2: 获取用户信息
	currentUser := auth.GetUserContext(ctx)
	if currentUser == nil {
		h.respondError(ctx, consts.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// 步骤3: 调用应用服务
	if err := h.botService.DeleteBot(ctx, botID, currentUser); err != nil {
		h.handleError(ctx, err)
		return
	}

	// 步骤4: 返回成功响应
	h.respondSuccess(ctx, nil, "Bot deleted successfully")
}

// PublishBot 发布Bot
//
// @Summary      发布Bot
// @Description  发布Bot到指定渠道
// @Tags         bots
// @Accept       json
// @Produce      json
// @Param        id path string true "Bot ID"
// @Param        request body dto.PublishBotRequest true "发布请求"
// @Success      200 {object} Response{data=dto.BotResponse}
// @Failure      400 {object} Response
// @Failure      401 {object} Response
// @Failure      403 {object} Response
// @Failure      404 {object} Response
// @Router       /api/v1/bots/{id}/publish [post]
func (h *BotHandler) PublishBot(ctx *app.RequestContext) {
	// 步骤1: 解析路径参数
	botID := ctx.Param("id")
	if botID == "" {
		h.respondError(ctx, consts.StatusBadRequest, "bot ID is required", nil)
		return
	}

	// 步骤2: 解析请求体
	var req dto.PublishBotRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		h.respondError(ctx, consts.StatusBadRequest, "invalid request body", err)
		return
	}

	// 步骤3: 获取用户信息
	currentUser := auth.GetUserContext(ctx)
	if currentUser == nil {
		h.respondError(ctx, consts.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// 步骤4: 调用应用服务
	bot, err := h.botService.PublishBot(ctx, botID, &req, currentUser)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	// 步骤5: 返回成功响应
	h.respondSuccess(ctx, bot, "Bot published successfully")
}

// ListBots 查询Bot列表
//
// @Summary      查询Bot列表
// @Description  分页查询租户下的Bot列表，支持过滤和排序
// @Tags         bots
// @Accept       json
// @Produce      json
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页大小" default(20)
// @Param        sort_by query string false "排序字段" Enums(id, name, created_at, updated_at)
// @Param        sort_order query string false "排序方向" Enums(asc, desc)
// @Param        name query string false "名称模糊搜索"
// @Param        status query string false "状态过滤" Enums(draft, published, archived)
// @Param        published query bool false "是否发布过滤"
// @Param        tag query string false "标签过滤"
// @Success      200 {object} Response{data=dto.BotListResponse}
// @Failure      400 {object} Response
// @Failure      401 {object} Response
// @Router       /api/v1/bots [get]
func (h *BotHandler) ListBots(ctx *app.RequestContext) {
	// 步骤1: 解析查询参数
	req := &dto.ListBotsRequest{}

	// 分页参数
	if page := ctx.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			req.Page = p
		}
	}

	if pageSize := ctx.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			req.PageSize = ps
		}
	}

	// 排序参数
	req.SortBy = ctx.Query("sort_by")
	req.SortOrder = ctx.Query("sort_order")

	// 过滤参数
	req.Name = ctx.Query("name")
	req.Status = stringPtr(ctx.Query("status"))
	if published := ctx.Query("published"); published != "" {
		if p, err := strconv.ParseBool(published); err == nil {
			req.Published = &p
		}
	}
	req.Tag = ctx.Query("tag")

	// 步骤2: 获取用户信息
	currentUser := auth.GetUserContext(ctx)
	if currentUser == nil {
		h.respondError(ctx, consts.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// 步骤3: 调用应用服务
	result, err := h.botService.ListBots(ctx, req, currentUser)
	if err != nil {
		h.handleError(ctx, err)
		return
	}

	// 步骤4: 返回成功响应
	h.respondSuccess(ctx, result, "Bots retrieved successfully")
}

// ============ 响应辅助方法 ============

// Response 标准响应结构
type Response struct {
	// Code 响应码（0表示成功）
	Code int `json:"code"`
	// Message 响应消息
	Message string `json:"message"`
	// Data 响应数据
	Data interface{} `json:"data,omitempty"`
	// RequestID 请求追踪ID
	RequestID string `json:"request_id"`
}

// respondSuccess 返回成功响应
func (h *BotHandler) respondSuccess(ctx *app.RequestContext, data interface{}, message string) {
	response := Response{
		Code:      0,
		Message:   message,
		Data:      data,
		RequestID: getRequestID(ctx),
	}
	ctx.JSON(consts.StatusOK, response)
}

// respondError 返回错误响应
func (h *BotHandler) respondError(
	ctx *app.RequestContext,
	httpStatus int,
	message string,
	err error,
) {
	response := Response{
		Code:      httpStatus,
		Message:   message,
		RequestID: getRequestID(ctx),
	}

	// 如果有详细错误，添加到响应中
	if err != nil {
		response.Data = map[string]interface{}{
			"error": err.Error(),
		}
	}

	ctx.JSON(httpStatus, response)
}

// handleError 处理应用层错误
func (h *BotHandler) handleError(ctx *app.RequestContext, err error) {
	// 这里需要根据错误类型返回相应的HTTP状态码
	// 简化实现：所有错误都返回500
	h.respondError(ctx, consts.StatusInternalServerError, "internal server error", err)
}

// getRequestID 从上下文获取请求ID
func getRequestID(ctx *app.RequestContext) string {
	// 假设请求ID由中间件设置在请求头中
	requestID := ctx.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = "unknown"
	}
	return requestID
}

// stringPtr 辅助函数：返回字符串指针
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
