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

package errorx

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/coze-dev/coze-studio/backend/pkg/errorx/i18n"
	// 注意: 移除了对types/errno的导入,避免循环依赖
)

// ErrorResponse 统一错误响应格式
type ErrorResponse struct {
	Code      int32                  `json:"code"`                // 错误码
	Message   string                 `json:"message"`             // 错误消息（根据请求语言）
	MessageZH string                 `json:"message_zh,omitempty"` // 中文错误消息（可选）
	Details   map[string]interface{} `json:"details,omitempty"`   // 错误详情
	RequestID string                 `json:"request_id"`          // 请求ID（用于追踪）
	Timestamp  int64                  `json:"timestamp"`           // 时间戳（毫秒）
	TenantID   string                 `json:"tenant_id,omitempty"` // 租户ID（可选）
	UserID     string                 `json:"user_id,omitempty"`   // 用户ID（可选）
	TraceID    string                 `json:"trace_id,omitempty"`  // 追踪ID（可选）
}

// ResponseBuilder 响应构造器
type ResponseBuilder struct {
	ctx          context.Context
	errorCode    int32
	errorMsg     string
	lang         i18n.Language
	details      map[string]interface{}
	requestID    string
	tenantID     string
	userID       string
	traceID      string
	timestamp    int64
	httpStatus   int
}

// NewResponseBuilder 创建响应构造器
func NewResponseBuilder(ctx context.Context, errorCode int32) *ResponseBuilder {
	return &ResponseBuilder{
		ctx:        ctx,
		errorCode:  errorCode,
		lang:       i18n.LanguageDefault, // 默认英文
		details:    make(map[string]interface{}),
		timestamp:  time.Now().UnixMilli(),
		httpStatus: http.StatusOK, // 默认200，实际会根据错误码映射
	}
}

// WithLanguage 设置响应语言
func (b *ResponseBuilder) WithLanguage(lang i18n.Language) *ResponseBuilder {
	b.lang = lang
	return b
}

// WithMessage 设置错误消息（会覆盖i18n翻译）
func (b *ResponseBuilder) WithMessage(msg string) *ResponseBuilder {
	b.errorMsg = msg
	return b
}

// WithDetail 添加错误详情
func (b *ResponseBuilder) WithDetail(key string, value interface{}) *ResponseBuilder {
	b.details[key] = value
	return b
}

// WithDetails 添加多个错误详情
func (b *ResponseBuilder) WithDetails(details map[string]interface{}) *ResponseBuilder {
	for k, v := range details {
		b.details[k] = v
	}
	return b
}

// WithRequestID 设置请求ID
func (b *ResponseBuilder) WithRequestID(requestID string) *ResponseBuilder {
	b.requestID = requestID
	return b
}

// WithTenantID 设置租户ID
func (b *ResponseBuilder) WithTenantID(tenantID string) *ResponseBuilder {
	b.tenantID = tenantID
	return b
}

// WithUserID 设置用户ID
func (b *ResponseBuilder) WithUserID(userID string) *ResponseBuilder {
	b.userID = userID
	return b
}

// WithTraceID 设置追踪ID
func (b *ResponseBuilder) WithTraceID(traceID string) *ResponseBuilder {
	b.traceID = traceID
	return b
}

// WithHTTPStatus 设置HTTP状态码（如果不设置会根据错误码自动映射）
func (b *ResponseBuilder) WithHTTPStatus(status int) *ResponseBuilder {
	b.httpStatus = status
	return b
}

// Build 构建错误响应
func (b *ResponseBuilder) Build() *ErrorResponse {
	// 生成request_id（如果未设置）
	if b.requestID == "" {
		b.requestID = uuid.New().String()
	}

	// 获取错误消息
	var message string
	if b.errorMsg != "" {
		// 使用自定义消息
		message = b.errorMsg
	} else {
		// 使用i18n翻译，从i18n包获取默认消息
		if defaultMsg, ok := i18n.Translate(b.errorCode, i18n.LanguageEN); ok {
			message = i18n.TranslateWithFallback(b.errorCode, b.lang, defaultMsg)
		} else {
			// 如果i18n中也没有，使用通用消息
			message = i18n.TranslateWithFallback(b.errorCode, b.lang, "Internal error")
		}
	}

	// 获取中文消息（用于message_zh字段）
	messageZH := i18n.TranslateWithFallback(b.errorCode, i18n.LanguageZH, message)

	// 如果未设置HTTP状态码，根据错误码映射
	if b.httpStatus == http.StatusOK {
		b.httpStatus = getDefaultHTTPStatus(b.errorCode)
	}

	return &ErrorResponse{
		Code:      b.errorCode,
		Message:   message,
		MessageZH: messageZH,
		Details:   b.details,
		RequestID: b.requestID,
		Timestamp: b.timestamp,
		TenantID:  b.tenantID,
		UserID:    b.userID,
		TraceID:   b.traceID,
	}
}

// HTTPStatus 获取HTTP状态码
func (b *ResponseBuilder) HTTPStatus() int {
	return b.httpStatus
}

// 便捷函数

// BuildErrorResponse 快速构造错误响应（重命名避免与ErrorResponse类型冲突）
func BuildErrorResponse(ctx context.Context, errorCode int32) *ErrorResponse {
	return NewResponseBuilder(ctx, errorCode).Build()
}

// ErrorResponseWithDetails 构造带详情的错误响应
func ErrorResponseWithDetails(ctx context.Context, errorCode int32, details map[string]interface{}) *ErrorResponse {
	return NewResponseBuilder(ctx, errorCode).WithDetails(details).Build()
}

// ErrorResponseWithLanguage 构造指定语言的错误响应
func ErrorResponseWithLanguage(ctx context.Context, errorCode int32, lang i18n.Language) *ErrorResponse {
	return NewResponseBuilder(ctx, errorCode).WithLanguage(lang).Build()
}

// ErrorResponseFromError 从标准error构造响应
func ErrorResponseFromError(ctx context.Context, err error) *ErrorResponse {
	if err == nil {
		return nil
	}

	// 尝试转换为StatusError
	if statusErr, ok := err.(StatusError); ok {
		return NewResponseBuilder(ctx, statusErr.Code()).
			WithMessage(statusErr.Msg()).
			Build()
	}

	// 如果不是StatusError，返回通用错误
	return NewResponseBuilder(ctx, 500000000).
		WithMessage(err.Error()).
		WithHTTPStatus(http.StatusInternalServerError).
		Build()
}

// ErrorResponseFromMap 从map构造响应（用于从其他层传递错误信息）
func ErrorResponseFromMap(ctx context.Context, errMap map[string]interface{}) *ErrorResponse {
	builder := NewResponseBuilder(ctx, 0)

	// 提取code
	if code, ok := errMap["code"].(int32); ok {
		builder.errorCode = code
	}

	// 提取message
	if msg, ok := errMap["message"].(string); ok {
		builder.WithMessage(msg)
	}

	// 提取details
	if details, ok := errMap["details"].(map[string]interface{}); ok {
		builder.WithDetails(details)
	}

	// 提取request_id
	if requestID, ok := errMap["request_id"].(string); ok {
		builder.WithRequestID(requestID)
	}

	// 提取tenant_id
	if tenantID, ok := errMap["tenant_id"].(string); ok {
		builder.WithTenantID(tenantID)
	}

	// 提取user_id
	if userID, ok := errMap["user_id"].(string); ok {
		builder.WithUserID(userID)
	}

	// 提取trace_id
	if traceID, ok := errMap["trace_id"].(string); ok {
		builder.WithTraceID(traceID)
	}

	return builder.Build()
}

// SuccessResponse 成功响应格式（用于对比）
type SuccessResponse struct {
	Code      int32                  `json:"code"`                // 业务状态码（成功时为0）
	Message   string                 `json:"message"`             // 成功消息
	Data      interface{}            `json:"data"`                // 响应数据
	RequestID string                 `json:"request_id"`          // 请求ID
	Timestamp  int64                  `json:"timestamp"`           // 时间戳
	Metadata  map[string]interface{} `json:"metadata,omitempty"`  // 元数据
}

// NewSuccessResponse 构造成功响应
func NewSuccessResponse(ctx context.Context, data interface{}, message ...string) *SuccessResponse {
	msg := "Success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	return &SuccessResponse{
		Code:      0,
		Message:   msg,
		Data:      data,
		RequestID: GetRequestID(ctx),
		Timestamp: time.Now().UnixMilli(),
		Metadata:  make(map[string]interface{}),
	}
}

// GetRequestID 从context中提取request ID
// 如果context中没有request ID，则生成一个新的
func GetRequestID(ctx context.Context) string {
	// 尝试从context中获取request_id
	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		return requestID
	}

	// 尝试从context中获取x-request-id header
	if requestID, ok := ctx.Value("x-request-id").(string); ok && requestID != "" {
		return requestID
	}

	// 如果都没有，生成一个新的request ID
	return uuid.New().String()
}

// getDefaultHTTPStatus 根据错误码获取默认HTTP状态码
// 使用范围映射避免为每个错误码单独映射
func getDefaultHTTPStatus(code int32) int {
	// 提取错误码的最后3位
	errorType := code % 1000

	// 001-009: NotFound → 404
	if errorType >= 1 && errorType <= 9 {
		return http.StatusNotFound
	}

	// 010-019: AlreadyExists → 409
	if errorType >= 10 && errorType <= 19 {
		return http.StatusConflict
	}

	// 020-039: InvalidParam → 400
	if errorType >= 20 && errorType <= 39 {
		return http.StatusBadRequest
	}

	// 040-049: PermissionDenied → 403
	if errorType >= 40 && errorType <= 49 {
		return http.StatusForbidden
	}

	// 050-089: Failed → 500
	if errorType >= 50 && errorType <= 89 {
		return http.StatusInternalServerError
	}

	// 090-099: Other → 400
	if errorType >= 90 && errorType <= 99 {
		return http.StatusBadRequest
	}

	// 默认 → 500
	return http.StatusInternalServerError
}
