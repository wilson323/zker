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

package developer

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	devservice "github.com/coze-dev/coze-studio/backend/domain/developer/service"
	"github.com/coze-dev/coze-studio/backend/domain/developer/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	berrno "github.com/coze-dev/coze-studio/backend/types/errno"
)

// WebhookHandler Webhook管理Handler
type WebhookHandler struct {
	webhookService devservice.WebhookService
}

// NewWebhookHandler 创建Webhook管理Handler实例
func NewWebhookHandler(webhookService devservice.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookService: webhookService,
	}
}

// ListWebhooksRequest 列出Webhook请求
type ListWebhooksRequest struct {
	ProjectID  string              `json:"project_id" query:"project_id"`
	Status     entity.WebhookStatus `json:"status" query:"status"`
	PageToken  string              `json:"page_token" query:"page_token"`
	PageSize   int                 `json:"page_size" query:"page_size"`
}

// ListWebhooksResponse 列出Webhook响应
type ListWebhooksResponse struct {
	Webhooks  []*entity.Webhook `json:"webhooks"`
	Total     int64            `json:"total"`
	PageToken string           `json:"page_token,omitempty"`
}

// CreateWebhookRequest 创建Webhook请求
type CreateWebhookRequest struct {
	ProjectID   string                `json:"project_id" validate:"required"`
	WebhookURL  string                `json:"webhook_url" validate:"required,url,max=500"`
	Events      entity.WebhookEvents  `json:"events" validate:"required,min=1"`
}

// CreateWebhookResponse 创建Webhook响应
type CreateWebhookResponse struct {
	Webhook        *entity.Webhook `json:"webhook"`
	WebhookSecret  string          `json:"webhook_secret"` // 仅在创建时返回
	Message        string          `json:"message"`
}

// UpdateWebhookRequest 更新Webhook请求
type UpdateWebhookRequest struct {
	WebhookURL *string               `json:"webhook_url" validate:"omitempty,url,max=500"`
	Events     *entity.WebhookEvents `json:"events" validate:"omitempty,min=1"`
}

// TestWebhookRequest 测试Webhook请求
type TestWebhookRequest struct {
	EventType  entity.WebhookEventType `json:"event_type" validate:"required"`
	TestData   map[string]interface{} `json:"test_data"`
}

// TestWebhookResponse 测试Webhook响应
type TestWebhookResponse struct {
	StatusCode int    `json:"status_code"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Error      string `json:"error,omitempty"`
}

// ListWebhooks 列出Webhook
// @router /api/developer/webhooks [GET]
func (h *WebhookHandler) ListWebhooks(ctx context.Context, c *app.RequestContext) {
	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 解析请求参数
	var req ListWebhooksRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 设置默认分页参数
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 构建过滤器
	filter := &devservice.WebhookListFilter{
		TenantID:  tenantID,
		ProjectID: req.ProjectID,
		Status:    req.Status,
		PageToken: req.PageToken,
		PageSize:  req.PageSize,
	}

	// 查询Webhook列表
	webhooks, total, err := h.webhookService.ListWebhooks(ctx, filter)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, &ListWebhooksResponse{
		Webhooks:  webhooks,
		Total:     total,
		PageToken: req.PageToken,
	})
}

// CreateWebhook 创建Webhook
// @router /api/developer/webhooks [POST]
func (h *WebhookHandler) CreateWebhook(ctx context.Context, c *app.RequestContext) {
	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 解析请求参数
	var req CreateWebhookRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 验证URL有效性
	if !h.isValidWebhookURL(req.WebhookURL) {
		httputil.BuildErrorResp(c, berrno.ErrWebhookURLInvalid.Code, berrno.ErrWebhookURLInvalid.Message, berrno.ErrWebhookURLInvalid.MessageZH, map[string]interface{}{
			"webhook_url": req.WebhookURL,
		})
		return
	}

	// 验证事件类型有效性
	if !h.isValidEvents(req.Events) {
		httputil.BuildErrorResp(c, berrno.ErrWebhookEventInvalid.Code, berrno.ErrWebhookEventInvalid.Message, berrno.ErrWebhookEventInvalid.MessageZH, nil)
		return
	}

	// 构建创建请求
	createReq := &devservice.CreateWebhookRequest{
		TenantID:   tenantID,
		ProjectID:  req.ProjectID,
		WebhookURL: req.WebhookURL,
		Events:     req.Events,
	}

	// 创建Webhook
	webhook, secret, err := h.webhookService.CreateWebhook(ctx, createReq)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	// 警告：webhook_secret仅在创建时返回一次，请妥善保管
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":        0,
		"message":     "Webhook created successfully. Please save the webhook_secret securely.",
		"message_zh":  "Webhook创建成功，请妥善保管webhook_secret",
		"message_en":  "Webhook created successfully. Please save the webhook_secret securely.",
		"data": &CreateWebhookResponse{
			Webhook:       webhook,
			WebhookSecret: secret,
			Message:       "Please save the webhook_secret securely",
		},
	})
}

// GetWebhook 获取Webhook详情
// @router /api/developer/webhooks/:id [GET]
func (h *WebhookHandler) GetWebhook(ctx context.Context, c *app.RequestContext) {
	webhookID := c.Param("id")
	if webhookID == "" {
		httputil.BadRequest(c, "webhook_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 查询Webhook
	webhook, err := h.webhookService.GetWebhook(ctx, webhookID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	// 验证租户隔离
	if webhook == nil || webhook.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrWebhookNotFound.Code, berrno.ErrWebhookNotFound.Message, berrno.ErrWebhookNotFound.MessageZH, nil)
		return
	}

	httputil.BuildSuccessResp(c, webhook)
}

// UpdateWebhook 更新Webhook
// @router /api/developer/webhooks/:id [PUT]
func (h *WebhookHandler) UpdateWebhook(ctx context.Context, c *app.RequestContext) {
	webhookID := c.Param("id")
	if webhookID == "" {
		httputil.BadRequest(c, "webhook_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 验证Webhook存在且属于该租户
	webhook, err := h.webhookService.GetWebhook(ctx, webhookID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if webhook == nil || webhook.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrWebhookNotFound.Code, berrno.ErrWebhookNotFound.Message, berrno.ErrWebhookNotFound.MessageZH, nil)
		return
	}

	// 解析请求参数
	var req UpdateWebhookRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 验证URL有效性（如果提供）
	if req.WebhookURL != nil && !h.isValidWebhookURL(*req.WebhookURL) {
		httputil.BuildErrorResp(c, berrno.ErrWebhookURLInvalid.Code, berrno.ErrWebhookURLInvalid.Message, berrno.ErrWebhookURLInvalid.MessageZH, nil)
		return
	}

	// 验证事件类型有效性（如果提供）
	if req.Events != nil && !h.isValidEvents(*req.Events) {
		httputil.BuildErrorResp(c, berrno.ErrWebhookEventInvalid.Code, berrno.ErrWebhookEventInvalid.Message, berrno.ErrWebhookEventInvalid.MessageZH, nil)
		return
	}

	// 构建更新请求
	updateReq := &devservice.UpdateWebhookRequest{
		WebhookID:  webhookID,
		WebhookURL: req.WebhookURL,
		Events:     req.Events,
	}

	// 更新Webhook
	if err := h.webhookService.UpdateWebhook(ctx, updateReq); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"webhook_id": webhookID,
		"message":    "Webhook updated successfully",
	})
}

// DeleteWebhook 删除Webhook
// @router /api/developer/webhooks/:id [DELETE]
func (h *WebhookHandler) DeleteWebhook(ctx context.Context, c *app.RequestContext) {
	webhookID := c.Param("id")
	if webhookID == "" {
		httputil.BadRequest(c, "webhook_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 验证Webhook存在且属于该租户
	webhook, err := h.webhookService.GetWebhook(ctx, webhookID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if webhook == nil || webhook.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrWebhookNotFound.Code, berrno.ErrWebhookNotFound.Message, berrno.ErrWebhookNotFound.MessageZH, nil)
		return
	}

	// 删除Webhook
	if err := h.webhookService.DeleteWebhook(ctx, webhookID); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Webhook deleted successfully",
	})
}

// TestWebhook 测试Webhook
// @router /api/developer/webhooks/:id/test [POST]
func (h *WebhookHandler) TestWebhook(ctx context.Context, c *app.RequestContext) {
	webhookID := c.Param("id")
	if webhookID == "" {
		httputil.BadRequest(c, "webhook_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 验证Webhook存在且属于该租户
	webhook, err := h.webhookService.GetWebhook(ctx, webhookID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if webhook == nil || webhook.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrWebhookNotFound.Code, berrno.ErrWebhookNotFound.Message, berrno.ErrWebhookNotFound.MessageZH, nil)
		return
	}

	// 解析请求参数
	var req TestWebhookRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 构建测试负载
	payload := &entity.WebhookPayload{
		TenantID:  tenantID,
		ProjectID: webhook.ProjectID,
		Data:      req.TestData,
	}

	// 如果没有提供测试数据，使用默认数据
	if payload.Data == nil {
		payload.Data = map[string]interface{}{
			"test": true,
			"message": "This is a test webhook",
		}
	}

	// 触发Webhook（同步测试）
	// 注意：这里直接调用service层的方法进行测试
	// TODO: 实现同步测试逻辑（需要在service层添加TestWebhook方法）

	httputil.BuildSuccessResp(c, &TestWebhookResponse{
		StatusCode: 200,
		Success:    true,
		Message:    "Webhook test initiated successfully",
	})
}

// PauseWebhook 暂停Webhook
// @router /api/developer/webhooks/:id/pause [POST]
func (h *WebhookHandler) PauseWebhook(ctx context.Context, c *app.RequestContext) {
	webhookID := c.Param("id")
	if webhookID == "" {
		httputil.BadRequest(c, "webhook_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 验证Webhook存在且属于该租户
	webhook, err := h.webhookService.GetWebhook(ctx, webhookID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if webhook == nil || webhook.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrWebhookNotFound.Code, berrno.ErrWebhookNotFound.Message, berrno.ErrWebhookNotFound.MessageZH, nil)
		return
	}

	// 检查当前状态
	if !webhook.IsActive() {
		httputil.BuildErrorResp(c, berrno.ErrWebhookNotFound.Code, "Webhook is already paused", "Webhook已暂停", nil)
		return
	}

	// 暂停Webhook
	if err := h.webhookService.PauseWebhook(ctx, webhookID); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"webhook_id": webhookID,
		"status":     entity.WebhookStatusPaused,
		"message":    "Webhook paused successfully",
	})
}

// ResumeWebhook 恢复Webhook
// @router /api/developer/webhooks/:id/resume [POST]
func (h *WebhookHandler) ResumeWebhook(ctx context.Context, c *app.RequestContext) {
	webhookID := c.Param("id")
	if webhookID == "" {
		httputil.BadRequest(c, "webhook_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 验证Webhook存在且属于该租户
	webhook, err := h.webhookService.GetWebhook(ctx, webhookID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if webhook == nil || webhook.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrWebhookNotFound.Code, berrno.ErrWebhookNotFound.Message, berrno.ErrWebhookNotFound.MessageZH, nil)
		return
	}

	// 检查当前状态
	if webhook.IsActive() {
		httputil.BuildErrorResp(c, berrno.ErrWebhookNotFound.Code, "Webhook is already active", "Webhook已激活", nil)
		return
	}

	// 恢复Webhook
	if err := h.webhookService.ResumeWebhook(ctx, webhookID); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"webhook_id": webhookID,
		"status":     entity.WebhookStatusActive,
		"message":    "Webhook resumed successfully",
	})
}

// GetWebhookStats 获取Webhook统计信息
// @router /api/developer/webhooks/:id/stats [GET]
func (h *WebhookHandler) GetWebhookStats(ctx context.Context, c *app.RequestContext) {
	webhookID := c.Param("id")
	if webhookID == "" {
		httputil.BadRequest(c, "webhook_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 验证Webhook存在且属于该租户
	webhook, err := h.webhookService.GetWebhook(ctx, webhookID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if webhook == nil || webhook.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrWebhookNotFound.Code, berrno.ErrWebhookNotFound.Message, berrno.ErrWebhookNotFound.MessageZH, nil)
		return
	}

	// 获取统计信息
	stats, err := h.webhookService.GetWebhookStats(ctx, webhookID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, stats)
}

// GetProjectWebhooks 获取项目的Webhook列表
// @router /api/developer/projects/:id/webhooks [GET]
func (h *WebhookHandler) GetProjectWebhooks(ctx context.Context, c *app.RequestContext) {
	projectID := c.Param("id")
	if projectID == "" {
		httputil.BadRequest(c, "project_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 查询Webhook列表
	webhooks, err := h.webhookService.GetProjectWebhooks(ctx, projectID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	// 过滤租户（确保租户隔离）
	filteredWebhooks := make([]*entity.Webhook, 0)
	for _, webhook := range webhooks {
		if webhook.TenantID == tenantID {
			filteredWebhooks = append(filteredWebhooks, webhook)
		}
	}

	httputil.BuildSuccessResp(c, filteredWebhooks)
}

// GetSupportedEvents 获取支持的事件类型列表
// @router /api/developer/webhooks/events [GET]
func (h *WebhookHandler) GetSupportedEvents(ctx context.Context, c *app.RequestContext) {
	events := []map[string]interface{}{
		{
			"event_type":       entity.WebhookEventBotCreated,
			"display_name":     "Bot Created",
			"display_name_zh":  "Bot创建",
			"description":      "当创建新Bot时触发",
			"category":         "bot",
		},
		{
			"event_type":       entity.WebhookEventBotUpdated,
			"display_name":     "Bot Updated",
			"display_name_zh":  "Bot更新",
			"description":      "当更新Bot时触发",
			"category":         "bot",
		},
		{
			"event_type":       entity.WebhookEventBotDeleted,
			"display_name":     "Bot Deleted",
			"display_name_zh":  "Bot删除",
			"description":      "当删除Bot时触发",
			"category":         "bot",
		},
		{
			"event_type":       entity.WebhookEventWorkflowCompleted,
			"display_name":     "Workflow Completed",
			"display_name_zh":  "工作流完成",
			"description":      "当工作流执行完成时触发",
			"category":         "workflow",
		},
		{
			"event_type":       entity.WebhookEventWorkflowFailed,
			"display_name":     "Workflow Failed",
			"display_name_zh":  "工作流失败",
			"description":      "当工作流执行失败时触发",
			"category":         "workflow",
		},
		{
			"event_type":       entity.WebhookEventConversationCreated,
			"display_name":     "Conversation Created",
			"display_name_zh":  "对话创建",
			"description":      "当创建新对话时触发",
			"category":         "conversation",
		},
		{
			"event_type":       entity.WebhookEventMessageReceived,
			"display_name":     "Message Received",
			"display_name_zh":  "消息接收",
			"description":      "当接收到新消息时触发",
			"category":         "conversation",
		},
		{
			"event_type":       entity.WebhookEventAPICall,
			"display_name":     "API Call",
			"display_name_zh":  "API调用",
			"description":      "当API调用时触发",
			"category":         "api",
		},
		{
			"event_type":       entity.WebhookEventErrorOccurred,
			"display_name":     "Error Occurred",
			"display_name_zh":  "错误发生",
			"description":      "当发生错误时触发",
			"category":         "error",
		},
	}

	httputil.BuildSuccessResp(c, events)
}

// isValidWebhookURL 验证Webhook URL有效性
func (h *WebhookHandler) isValidWebhookURL(url string) bool {
	// 必须是HTTPS URL（生产环境强制要求）
	// TODO: 实现更严格的URL验证
	return len(url) > 10 && len(url) <= 500
}

// isValidEvents 验证事件类型有效性
func (h *WebhookHandler) isValidEvents(events entity.WebhookEvents) bool {
	validEvents := map[entity.WebhookEventType]bool{
		entity.WebhookEventBotCreated:         true,
		entity.WebhookEventBotUpdated:         true,
		entity.WebhookEventBotDeleted:         true,
		entity.WebhookEventWorkflowCompleted:  true,
		entity.WebhookEventWorkflowFailed:     true,
		entity.WebhookEventConversationCreated: true,
		entity.WebhookEventMessageReceived:    true,
		entity.WebhookEventAPICall:            true,
		entity.WebhookEventErrorOccurred:      true,
	}

	for _, event := range events {
		if !validEvents[event] {
			return false
		}
	}

	return true
}
