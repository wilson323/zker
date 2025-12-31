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

package httputil

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/pkg/conv"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// SuccessResponse 企业级统一成功响应格式
type SuccessResponse struct {
	Code       int32       `json:"code"`                 // 业务状态码 (0表示成功)
	Message    string      `json:"message"`              // 英文成功消息
	MessageZH  string      `json:"message_zh,omitempty"` // 中文成功消息
	MessageEN  string      `json:"message_en,omitempty"` // 英文成功消息(冗余)
	Data       interface{} `json:"data,omitempty"`       // 响应数据
	Timestamp  string      `json:"timestamp"`            // 响应时间戳
	RequestID  string      `json:"request_id,omitempty"` // 请求追踪ID
	TraceID    string      `json:"trace_id,omitempty"`   // 分布式追踪ID
	TenantID   string      `json:"tenant_id,omitempty"`  // 租户ID
	ServerTime int64       `json:"server_time"`          // 服务器时间戳(Unix)
}

// PageData 分页数据结构
type PageData struct {
	Items      interface{} `json:"items"`       // 数据列表
	Total      int64       `json:"total"`       // 总记录数
	Page       int         `json:"page"`        // 当前页码
	PageSize   int         `json:"page_size"`   // 每页大小
	TotalPages int         `json:"total_pages"` // 总页数
}

// SuccessResp 构建成功响应(无数据)
// 用途: DELETE、PUT等操作的成功响应
// 遵循KISS原则:简单直接,参数最少化
func SuccessResp(c *app.RequestContext) {
	SuccessRespWithData(c, nil)
}

// SuccessRespWithData 构建成功响应(带数据)
// 用途: GET、POST等需要返回数据的操作
func SuccessRespWithData(c *app.RequestContext, data interface{}) {
	buildSuccessResponse(c, data, http.StatusOK)
}

// SuccessRespCreated 构建创建成功响应(201)
// 用途: POST创建资源的成功响应
func SuccessRespCreated(c *app.RequestContext, data interface{}) {
	buildSuccessResponse(c, data, http.StatusCreated)
}

// SuccessRespAccepted 构建已接受响应(202)
// 用途: 异步操作的接受响应
func SuccessRespAccepted(c *app.RequestContext, data interface{}) {
	buildSuccessResponse(c, data, http.StatusAccepted)
}

// SuccessRespNoContent 构建无内容响应(204)
// 用途: DELETE等无返回内容的操作
func SuccessRespNoContent(c *app.RequestContext) {
	c.Status(http.StatusNoContent)
}

// SuccessRespWithPage 构建分页成功响应
// 用途: 列表查询的分页响应
func SuccessRespWithPage(c *app.RequestContext, items interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	pageData := &PageData{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	SuccessRespWithData(c, pageData)
}

// buildSuccessResponse 内部方法:构建成功响应
// 遵循DRY原则:避免重复代码
func buildSuccessResponse(c *app.RequestContext, data interface{}, httpStatus int) {
	response := &SuccessResponse{
		Code:       0,
		Message:    "SUCCESS",
		MessageZH:  "操作成功",
		MessageEN:  "Operation successful",
		Data:       data,
		Timestamp:  time.Now().Format(time.RFC3339),
		ServerTime: time.Now().Unix(),
	}

	// 从context获取追踪信息
	if reqID := c.Query("request_id"); reqID != "" {
		response.RequestID = reqID
	}
	if traceID := conv.BytesToString(c.GetHeader("X-Trace-ID")); traceID != "" {
		response.TraceID = traceID
	}
	if tenantID := conv.BytesToString(c.GetHeader("X-Tenant-ID")); tenantID != "" {
		response.TenantID = tenantID
	}

	c.JSON(httpStatus, response)
}

// ErrorRespEnhanced 构建增强错误响应(使用errno.EnhancedError)
// 用途: 统一的错误响应处理,支持多语言和详细信息
func ErrorRespEnhanced(c *app.RequestContext, enhancedErr *errno.EnhancedError) {
	response := &ErrorResponse{
		Code:       enhancedErr.Code,
		Message:    enhancedErr.Message,
		MessageZH:  enhancedErr.MessageZH,
		MessageEN:  enhancedErr.Message,
		Details:    enhancedErr.Details,
		RequestID:  enhancedErr.RequestID,
		Timestamp:  enhancedErr.Timestamp,
		TraceID:    enhancedErr.TraceID,
		TenantID:   enhancedErr.TenantID,
	}

	httpStatus := enhancedErr.GetHTTPStatus()
	logs.CtxWarnf(context.Background(), "[API] Error response: code=%d, msg=%s, status=%d",
		response.Code, response.Message, httpStatus)

	c.AbortWithStatusJSON(httpStatus, response)
}

// ErrorRespCode 构建错误响应(使用错误码)
// 用途: 快速构建标准错误响应
func ErrorRespCode(c *app.RequestContext, errCode int32, message, messageZH string, details map[string]interface{}) {
	BuildErrorResp(c, errCode, message, messageZH, details)
}

// ErrorResp 构建简单错误响应
// 用途: 快速构建常见错误响应
func ErrorResp(c *app.RequestContext, errCode errno.ErrorCode, details map[string]interface{}) {
	enhancedErr := errno.NewEnhancedError(
		errCode.Code(),
		errCode.Message(),
		errCode.MessageZH(),
		details,
	)
	ErrorRespEnhanced(c, enhancedErr)
}

// BadRequestResp 构建400错误响应
func BadRequestResp(c *app.RequestContext, message string) {
	ErrorRespCode(c, errno.InvalidParamsCode, message, "参数错误", nil)
}

// UnauthorizedResp 构建401错误响应
func UnauthorizedResp(c *app.RequestContext, message string) {
	ErrorRespCode(c, errno.UnauthorizedCode, message, "未认证", nil)
}

// ForbiddenResp 构建403错误响应
func ForbiddenResp(c *app.RequestContext, message string) {
	ErrorRespCode(c, errno.ForbiddenCode, message, "权限不足", nil)
}

// NotFoundResp 构建404错误响应
func NotFoundResp(c *app.RequestContext, message string) {
	ErrorRespCode(c, errno.NotFoundCode, message, "资源不存在", nil)
}

// InternalErrorResp 构建500错误响应
func InternalErrorResp(ctx context.Context, c *app.RequestContext, err error) {
	logs.CtxErrorf(ctx, "[API] Internal error: %v", err)
	ErrorRespCode(c, errno.InternalErrorCode, "Internal server error", "服务器内部错误", nil)
}

// ValidationErrorResp 构建参数验证错误响应
// 用途: 统一的参数验证失败响应
func ValidationErrorResp(c *app.RequestContext, field, message string) {
	details := map[string]interface{}{
		"field":   field,
		"message": message,
	}
	ErrorRespCode(c, errno.InvalidParamsCode, "Validation failed", "参数验证失败", details)
}

// RateLimitResp 构建限流错误响应
// 用途: API限流时的错误响应
func RateLimitResp(c *app.RequestContext, retryAfter int) {
	details := map[string]interface{}{
		"retry_after": retryAfter,
	}
	ErrorRespCode(c, errno.ErrRateLimitExceededCode, "Rate limit exceeded", "请求过于频繁", details)
	c.Header("X-Retry-After", string(rune(retryAfter)))
	c.Header("X-RateLimit-Remaining", "0")
}

// ServiceUnavailableResp 构建服务不可用响应(503)
// 用途: 熔断、降级时的响应
func ServiceUnavailableResp(c *app.RequestContext, message string) {
	ErrorRespCode(c, 503000001, message, "服务暂时不可用", nil)
	c.Header("Retry-After", "60")
}
