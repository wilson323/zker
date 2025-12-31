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

package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	hertzconsts "github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"

	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/consts"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// ErrorHandlerMW 企业级错误处理中间件
// 功能:
// 1. 统一错误格式
// 2. 错误日志记录
// 3. 请求ID追踪
// 4. 双语支持(中英文)
// 5. HTTP状态码映射
// 6. 错误详情收集
func ErrorHandlerMW() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 生成请求ID(如果还没有)
		requestIDBytes := c.GetHeader("X-Request-ID")
		requestID := string(requestIDBytes)
		if requestID == "" {
			requestID = uuid.New().String()
			c.Header("X-Request-ID", requestID)
		}

		// 将请求ID存入context
		ctx = context.WithValue(ctx, consts.CtxLogIDKey, requestID)

		// 执行请求
		defer func() {
			// 捕获panic
			if r := recover(); r != nil {
				handlePanic(ctx, c, requestID, r)
			}
		}()

		c.Next(ctx)

		// 处理错误
		handleError(ctx, c, requestID)
	}
}

// handleError 处理错误
func handleError(ctx context.Context, c *app.RequestContext, requestID string) {
	// 获取响应状态码
	status := c.Response.StatusCode()

	// 如果状态码不是错误状态码,直接返回
	if status < http.StatusBadRequest {
		return
	}

	// 获取错误响应体
	responseBody := c.Response.Body()

	// 如果响应体已经被设置,直接返回
	if len(responseBody) > 0 {
		return
	}

	// 检查是否有错误
	if len(c.Errors) == 0 {
		return
	}

	// 获取最后一个错误
	err := c.Errors.Last()

	// 转换为增强错误
	enhancedErr := convertToEnhancedError(ctx, c, err.Err, requestID)

	// 记录错误日志
	logError(ctx, c, enhancedErr)

	// 设置响应
	c.JSON(enhancedErr.GetHTTPStatus(), enhancedErr.ToErrorWithFields())
}

// convertToEnhancedError 转换为增强错误
func convertToEnhancedError(ctx context.Context, c *app.RequestContext, err error, requestID string) *errno.EnhancedError {
	if err == nil {
		return nil
	}

	// 尝试转换为StatusError
	if statusErr, ok := err.(errorx.StatusError); ok {
		enhancedErr := errno.ConvertFromStatusErr(statusErr)
		enhancedErr.WithRequestID(requestID)
		enhancedErr.WithTimestamp(time.Now().Format(time.RFC3339))

		// 添加请求上下文信息
		addRequestContext(ctx, c, enhancedErr)

		return enhancedErr
	}

	// 如果是EnhancedError,直接返回
	if enhancedErr, ok := err.(*errno.EnhancedError); ok {
		enhancedErr.WithRequestID(requestID)
		enhancedErr.WithTimestamp(time.Now().Format(time.RFC3339))
		return enhancedErr
	}

	// 通用错误
	enhancedErr := &errno.EnhancedError{
		Code:       hertzconsts.StatusInternalServerError,
		Message:    "Internal server error",
		MessageZH:  "服务器内部错误",
		HTTPStatus: http.StatusInternalServerError,
		RequestID:  requestID,
		Timestamp:  time.Now().Format(time.RFC3339),
		OriginalErr: err,
		Details: map[string]interface{}{
			"underlying_error": err.Error(),
		},
	}

	// 添加请求上下文信息
	addRequestContext(ctx, c, enhancedErr)

	return enhancedErr
}

// addRequestContext 添加请求上下文信息
func addRequestContext(ctx context.Context, c *app.RequestContext, enhancedErr *errno.EnhancedError) {
	// 添加请求信息
	enhancedErr.WithDetail("http_method", string(c.Request.Method()))
	enhancedErr.WithDetail("http_path", string(c.Request.URI().Path()))
	enhancedErr.WithDetail("client_ip", c.ClientIP())
	enhancedErr.WithDetail("user_agent", string(c.Request.Header.UserAgent()))

	// 从context中提取租户ID和用户ID
	if tenantID := ctx.Value("tenant_id"); tenantID != nil {
		enhancedErr.WithTenantID(tenantID.(string))
	}
	if userID := ctx.Value("user_id"); userID != nil {
		enhancedErr.WithUserID(userID.(string))
	}
	if traceID := ctx.Value("trace_id"); traceID != nil {
		enhancedErr.WithTraceID(traceID.(string))
	}

	// 开发环境添加堆栈跟踪
	if isDevelopment() {
		enhancedErr.WithStackTrace(string(debug.Stack()))
	}
}

// logError 记录错误日志
func logError(ctx context.Context, c *app.RequestContext, enhancedErr *errno.EnhancedError) {
	// 构建日志字段
	logFields := []interface{}{
		"request_id", enhancedErr.RequestID,
		"error_code", enhancedErr.Code,
		"error_message", enhancedErr.Message,
		"error_message_zh", enhancedErr.MessageZH,
		"http_method", string(c.Request.Method()),
		"http_path", string(c.Request.URI().Path()),
		"http_status", enhancedErr.GetHTTPStatus(),
		"client_ip", c.ClientIP(),
	}

	// 添加租户ID和用户ID
	if enhancedErr.TenantID != "" {
		logFields = append(logFields, "tenant_id", enhancedErr.TenantID)
	}
	if enhancedErr.UserID != "" {
		logFields = append(logFields, "user_id", enhancedErr.UserID)
	}

	// 添加追踪ID
	if enhancedErr.TraceID != "" {
		logFields = append(logFields, "trace_id", enhancedErr.TraceID)
	}

	// 添加错误详情
	if len(enhancedErr.Details) > 0 {
		detailsJSON, _ := json.Marshal(enhancedErr.Details)
		logFields = append(logFields, "details", string(detailsJSON))
	}

	// 根据HTTP状态码选择日志级别
	status := enhancedErr.GetHTTPStatus()
	switch {
	case status >= http.StatusInternalServerError:
		// 服务器错误 - Error级别
		if enhancedErr.StackTrace != "" {
			logFields = append(logFields, "stack", enhancedErr.StackTrace)
		}
		logs.CtxErrorf(ctx, "HTTP error: %v", logFields...)
	case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
		// 客户端错误 - Warn级别
		logs.CtxWarnf(ctx, "HTTP client error: %v", logFields...)
	default:
		logs.CtxInfof(ctx, "HTTP request: %v", logFields...)
	}
}

// handlePanic 处理panic
func handlePanic(ctx context.Context, c *app.RequestContext, requestID string, r interface{}) {
	// 记录panic日志
	logs.CtxErrorf(ctx, "Panic recovered: %v\n%s", r, debug.Stack())

	// 创建panic错误响应
	enhancedErr := errno.NewEnhancedError(
		hertzconsts.StatusInternalServerError,
		"Internal server error",
		"服务器内部错误",
	)
	enhancedErr.WithRequestID(requestID)
	enhancedErr.WithTimestamp(time.Now().Format(time.RFC3339))
	enhancedErr.WithDetail("panic", fmt.Sprintf("%v", r))
	enhancedErr.WithDetail("stack_trace", string(debug.Stack()))

	// 添加请求上下文
	addRequestContext(ctx, c, enhancedErr)

	// 记录错误
	logError(ctx, c, enhancedErr)

	// 返回错误响应
	c.JSON(http.StatusInternalServerError, enhancedErr.ToErrorWithFields())
}

// isDevelopment 判断是否为开发环境
func isDevelopment() bool {
	// TODO: 从配置读取环境
	// return config.GetEnv() == "development"
	return true
}

// WrapError 包装错误并添加到Hertz上下文
func WrapError(c *app.RequestContext, err error) {
	if err != nil {
		c.Error(err)
	}
}

// ErrorWithCode 带错误码的错误响应
func ErrorWithCode(c *app.RequestContext, code int32, message, messageZH string) {
	enhancedErr := errno.NewEnhancedError(code, message, messageZH)
	enhancedErr.WithTimestamp(time.Now().Format(time.RFC3339))

	// 从context获取requestID
	if requestID := c.Value("request_id"); requestID != nil {
		if reqID, ok := requestID.(string); ok {
			enhancedErr.WithRequestID(reqID)
		}
	}

	// 添加到错误列表
	c.Error(enhancedErr)
}

// SuccessResponse 成功响应
func SuccessResponse(c *app.RequestContext, data interface{}) {
	response := errno.NewSuccessResponse(data)

	// 从context获取requestID
	if requestID := c.Value("request_id"); requestID != nil {
		if reqID, ok := requestID.(string); ok {
			response.WithRequestID(reqID)
		}
	}

	c.JSON(hertzconsts.StatusOK, response)
}

// CreatedResponse 创建成功响应(201)
func CreatedResponse(c *app.RequestContext, data interface{}) {
	response := errno.NewSuccessResponse(data)
	response.Message = "Resource created successfully"
	response.MessageZH = "资源创建成功"

	// 从context获取requestID
	if requestID := c.Value("request_id"); requestID != nil {
		if reqID, ok := requestID.(string); ok {
			response.WithRequestID(reqID)
		}
	}

	c.JSON(hertzconsts.StatusCreated, response)
}

// NoContentResponse 无内容响应(204)
func NoContentResponse(c *app.RequestContext) {
	c.Status(hertzconsts.StatusNoContent)
}

// BadRequestResponse 错误请求响应(400)
func BadRequestResponse(c *app.RequestContext, message, messageZH string) {
	ErrorWithCode(c, hertzconsts.StatusBadRequest, message, messageZH)
}

// UnauthorizedResponse 未认证响应(401)
func UnauthorizedResponse(c *app.RequestContext, message, messageZH string) {
	ErrorWithCode(c, hertzconsts.StatusUnauthorized, message, messageZH)
}

// ForbiddenResponse 禁止访问响应(403)
func ForbiddenResponse(c *app.RequestContext, message, messageZH string) {
	ErrorWithCode(c, hertzconsts.StatusForbidden, message, messageZH)
}

// NotFoundResponse 未找到响应(404)
func NotFoundResponse(c *app.RequestContext, message, messageZH string) {
	ErrorWithCode(c, hertzconsts.StatusNotFound, message, messageZH)
}

// ConflictResponse 冲突响应(409)
func ConflictResponse(c *app.RequestContext, message, messageZH string) {
	ErrorWithCode(c, hertzconsts.StatusConflict, message, messageZH)
}

// TooManyRequestsResponse 请求过多响应(429)
func TooManyRequestsResponse(c *app.RequestContext, message, messageZH string) {
	ErrorWithCode(c, hertzconsts.StatusTooManyRequests, message, messageZH)
}

// InternalServerErrorResponse 服务器错误响应(500)
func InternalServerErrorResponse(c *app.RequestContext, message, messageZH string) {
	ErrorWithCode(c, hertzconsts.StatusInternalServerError, message, messageZH)
}

// ServiceUnavailableResponse 服务不可用响应(503)
func ServiceUnavailableResponse(c *app.RequestContext, message, messageZH string) {
	ErrorWithCode(c, hertzconsts.StatusServiceUnavailable, message, messageZH)
}
