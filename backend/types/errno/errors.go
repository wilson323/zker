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

package errno

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
)

// HTTPStatusMapping 错误码到HTTP状态码的映射
var HTTPStatusMapping = map[int32]int{
	// 租户错误 (2xx)
	ErrTenantNotFoundCode:        http.StatusNotFound,
	ErrTenantAlreadyExistsCode:   http.StatusConflict,
	ErrTenantSuspendedCode:       http.StatusForbidden,
	ErrTenantDeletedCode:         http.StatusGone,
	ErrTenantInvalidParamCode:    http.StatusBadRequest,
	ErrInvalidTenantIDCode:       http.StatusBadRequest,

	// 配额错误 (3xx)
	ErrQuotaExceededCode:          http.StatusForbidden, // 403 Forbidden (quota exceeded)
	ErrQuotaInvalidCode:           http.StatusBadRequest,
	ErrQuotaCheckFailedCode:       http.StatusInternalServerError,
	ErrQuotaBotExceededCode:       http.StatusForbidden,
	ErrQuotaKnowledgeExceededCode:  http.StatusForbidden,
	ErrQuotaWorkflowExceededCode:   http.StatusForbidden,
	ErrQuotaAPICallExceededCode:    http.StatusTooManyRequests, // 429 Too Many Requests
	ErrQuotaStorageExceededCode:    http.StatusForbidden,
	ErrQuotaConcurrentExceededCode: http.StatusTooManyRequests, // 429 Too Many Requests

	// 订阅错误 (4xx)
	ErrSubscriptionNotFoundCode:      http.StatusNotFound,
	ErrSubscriptionExpiredCode:     http.StatusForbidden,
	ErrSubscriptionInactiveCode:    http.StatusForbidden,
	ErrSubscriptionSuspendedCode:   http.StatusForbidden,
	ErrSubscriptionPaymentFailedCode: http.StatusPaymentRequired, // 402 Payment Required

	// 用户错误 (7xx)
	ErrUserAuthenticationFailed:    http.StatusUnauthorized,
	ErrUserEmailAlreadyExistCode:   http.StatusConflict,
	ErrUserInfoInvalidateCode:      http.StatusUnauthorized,
	ErrUserResourceNotFound:        http.StatusNotFound,
	ErrUserInvalidParamCode:        http.StatusBadRequest,
	ErrUserPermissionCode:          http.StatusForbidden,

	// 权限错误 (108)
	ErrPermissionInvalidParamCode:   http.StatusBadRequest,
}

// GetHTTPStatusForError 根据错误码获取HTTP状态码
// 使用范围映射避免为每个错误码单独映射，遵循DRY原则
func GetHTTPStatusForError(code int32) int {
	// 如果有精确映射，使用精确映射
	if status, ok := HTTPStatusMapping[code]; ok {
		return status
	}

	// 按错误码段进行范围映射（遵循KISS原则）
	million := int32(1000000)
	codeSegment := code / million

	switch codeSegment {
	case 201: // Bot模块错误 (201 xxx xxx)
		return mapHTTPStatusByErrorCode(code)
	case 202: // Conversation模块错误 (202 xxx xxx)
		return mapHTTPStatusByErrorCode(code)
	case 203: // Workflow模块错误 (203 xxx xxx)
		return mapHTTPStatusByErrorCode(code)
	default:
		return DefaultHTTPStatus
	}
}

// mapHTTPStatusByErrorCode 根据错误码的最后3位映射HTTP状态码
// 遵循KISS原则，使用统一的映射规则
func mapHTTPStatusByErrorCode(code int32) int {
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

	// 020-039: InvalidParam/InvalidConfig → 400
	if errorType >= 20 && errorType <= 39 {
		return http.StatusBadRequest
	}

	// 040-049: PermissionDenied → 403
	if errorType >= 40 && errorType <= 49 {
		return http.StatusForbidden
	}

	// 050-089: Failed (各种失败) → 500
	if errorType >= 50 && errorType <= 89 {
		return http.StatusInternalServerError
	}

	// 090-099: 其他 → 400
	if errorType >= 90 && errorType <= 99 {
		return http.StatusBadRequest
	}

	// 默认 → 500
	return http.StatusInternalServerError
}

// DefaultHTTPStatus 默认HTTP状态码
const DefaultHTTPStatus = http.StatusInternalServerError

// EnhancedError 增强的错误结构,支持企业级特性
type EnhancedError struct {
	// 基础错误信息
	Code       int32                  `json:"code"`                 // 错误码
	Message    string                 `json:"message"`              // 英文错误消息
	MessageZH  string                 `json:"message_zh,omitempty"` // 中文错误消息
	HTTPStatus int                    `json:"-"`                    // HTTP状态码(不序列化到JSON)

	// 企业级增强字段
	RequestID   string                 `json:"request_id,omitempty"`   // 请求ID
	Timestamp   string                 `json:"timestamp,omitempty"`    // 时间戳
	Details     map[string]interface{} `json:"details,omitempty"`      // 错误详情
	StackTrace  string                 `json:"stack_trace,omitempty"`  // 堆栈跟踪(仅开发环境)
	TenantID    string                 `json:"tenant_id,omitempty"`    // 租户ID
	UserID      string                 `json:"user_id,omitempty"`      // 用户ID
	TraceID     string                 `json:"trace_id,omitempty"`     // 追踪ID

	// 内部原始错误
	OriginalErr error `json:"-"` // 原始错误(不序列化到JSON)
}

// Error 实现error接口
func (e *EnhancedError) Error() string {
	if e.MessageZH != "" {
		return fmt.Sprintf("[%d] %s (中文: %s)", e.Code, e.Message, e.MessageZH)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// GetHTTPStatus 获取HTTP状态码
func (e *EnhancedError) GetHTTPStatus() int {
	if e.HTTPStatus != 0 {
		return e.HTTPStatus
	}
	// 优先使用精确映射，否则使用范围映射（支持新增模块）
	return GetHTTPStatusForError(e.Code)
}

// WithRequestID 添加请求ID
func (e *EnhancedError) WithRequestID(requestID string) *EnhancedError {
	e.RequestID = requestID
	return e
}

// WithTimestamp 添加时间戳
func (e *EnhancedError) WithTimestamp(timestamp string) *EnhancedError {
	e.Timestamp = timestamp
	return e
}

// WithDetails 添加错误详情
func (e *EnhancedError) WithDetails(details map[string]interface{}) *EnhancedError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

// WithDetail 添加单个错误详情
func (e *EnhancedError) WithDetail(key string, value interface{}) *EnhancedError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithTenantID 添加租户ID
func (e *EnhancedError) WithTenantID(tenantID string) *EnhancedError {
	e.TenantID = tenantID
	return e
}

// WithUserID 添加用户ID
func (e *EnhancedError) WithUserID(userID string) *EnhancedError {
	e.UserID = userID
	return e
}

// WithTraceID 添加追踪ID
func (e *EnhancedError) WithTraceID(traceID string) *EnhancedError {
	e.TraceID = traceID
	return e
}

// WithStackTrace 添加堆栈跟踪
func (e *EnhancedError) WithStackTrace(stackTrace string) *EnhancedError {
	e.StackTrace = stackTrace
	return e
}

// ToJSON 转换为JSON
func (e *EnhancedError) ToJSON() []byte {
	data, _ := json.Marshal(e)
	return data
}

// ToMap 转换为map
func (e *EnhancedError) ToMap() map[string]interface{} {
	data, _ := json.Marshal(e)
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result
}

// NewEnhancedError 创建增强错误
func NewEnhancedError(code int32, message, messageZH string) *EnhancedError {
	return &EnhancedError{
		Code:       code,
		Message:    message,
		MessageZH:  messageZH,
		HTTPStatus: HTTPStatusMapping[code],
		Timestamp:  time.Now().Format(time.RFC3339),
	}
}

// WrapError 包装错误为增强错误
func WrapError(err error, code int32, message, messageZH string) *EnhancedError {
	if err == nil {
		return nil
	}

	enhancedErr := &EnhancedError{
		Code:       code,
		Message:    message,
		MessageZH:  messageZH,
		HTTPStatus: HTTPStatusMapping[code],
		Timestamp:  time.Now().Format(time.RFC3339),
		OriginalErr: err,
		Details: map[string]interface{}{
			"underlying_error": err.Error(),
		},
	}

	// 如果原错误是StatusError,提取额外信息
	if statusErr, ok := err.(errorx.StatusError); ok {
		enhancedErr.Details["extra"] = statusErr.Extra()
		enhancedErr.Details["affect_stability"] = statusErr.IsAffectStability()
	}

	return enhancedErr
}

// ConvertFromStatusError 从StatusError转换为EnhancedError
func ConvertFromStatusErr(err error) *EnhancedError {
	if err == nil {
		return nil
	}

	statusErr, ok := err.(errorx.StatusError)
	if !ok {
		return &EnhancedError{
			Code:      DefaultHTTPStatus,
			Message:   err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	// 获取错误码和消息
	code := statusErr.Code()
	msg := statusErr.Msg()

	// 创建增强错误
	enhancedErr := &EnhancedError{
		Code:       code,
		Message:    msg,
		MessageZH:  getChineseMessage(code), // 尝试获取中文消息
		HTTPStatus: HTTPStatusMapping[code],
		Timestamp:  time.Now().Format(time.RFC3339),
		OriginalErr: err,
		Details: map[string]interface{}{
			"affect_stability": statusErr.IsAffectStability(),
			"extra":            statusErr.Extra(),
		},
	}

	return enhancedErr
}

// getChineseMessage 根据错误码获取中文消息
// 优先从CodeDefinition读取，如果没有则使用范围映射生成（遵循DRY原则）
func getChineseMessage(code int32) string {
	// 1. 优先从已注册的CodeDefinition读取
	if def := code.GetCodeDefinition(code); def != nil && def.MessageZH != "" {
		return def.MessageZH
	}

	// 2. 使用范围映射生成中文消息（避免为每个错误码单独映射）
	return generateChineseMessageByRange(code)
}

// generateChineseMessageByRange 根据错误码范围生成中文消息
// 遵循KISS原则，使用统一的映射规则，避免为每个错误码单独添加case语句
func generateChineseMessageByRange(code int32) string {
	// 提取错误码的最后3位
	errorType := code % 1000

	// 根据错误码段和类型生成中文消息
	module := getModuleName(code)

	// 001-009: NotFound（不存在）
	if errorType >= 1 && errorType <= 9 {
		return module + "不存在"
	}

	// 010-019: AlreadyExists（已存在）
	if errorType >= 10 && errorType <= 19 {
		return module + "已存在"
	}

	// 020-039: InvalidParam/InvalidConfig（参数无效）
	if errorType >= 20 && errorType <= 39 {
		return module + "参数无效"
	}

	// 040-049: PermissionDenied（权限不足）
	if errorType >= 40 && errorType <= 49 {
		return module + "权限不足"
	}

	// 050-089: Failed（操作失败）
	if errorType >= 50 && errorType <= 89 {
		return module + "操作失败"
	}

	// 090-099: Other（其他错误）
	if errorType >= 90 && errorType <= 99 {
		return module + "错误"
	}

	// 默认：内部错误
	return "内部错误"
}

// getModuleName 根据错误码获取模块名称
func getModuleName(code int32) string {
	million := int32(1000000)
	codeSegment := code / million

	switch codeSegment {
	case 2:
		return "租户" // 2xx: 租户错误
	case 3:
		return "配额" // 3xx: 配额错误
	case 4:
		return "订阅" // 4xx: 订阅错误
	case 201:
		return "Bot" // 201xxx: Bot模块错误
	case 202:
		return "对话" // 202xxx: Conversation模块错误
	case 203:
		return "工作流" // 203xxx: Workflow模块错误
	case 7:
		return "用户" // 7xx: 用户错误
	case 108:
		return "权限" // 108xxx: 权限错误
	default:
		return "系统"
	}
}

// ErrorWithFields 带字段的错误响应结构
type ErrorWithFields struct {
	Code      int32                  `json:"code"`                // 错误码
	Message   string                 `json:"message"`             // 英文消息
	MessageZH string                 `json:"message_zh,omitempty"` // 中文消息
	Data      map[string]interface{} `json:"data,omitempty"`      // 附加数据
	RequestID string                 `json:"request_id,omitempty"` // 请求ID
	Timestamp string                 `json:"timestamp,omitempty"`  // 时间戳
}

// ToErrorWithFields 转换为ErrorWithFields格式
func (e *EnhancedError) ToErrorWithFields() *ErrorWithFields {
	return &ErrorWithFields{
		Code:      e.Code,
		Message:   e.Message,
		MessageZH: e.MessageZH,
		Data:      e.Details,
		RequestID: e.RequestID,
		Timestamp: e.Timestamp,
	}
}

// SuccessResponse 成功响应结构
type SuccessResponse struct {
	Code      int32                  `json:"code"`                 // 状态码(0表示成功)
	Message   string                 `json:"message"`              // 成功消息
	MessageZH string                 `json:"message_zh,omitempty"` // 中文成功消息
	Data      interface{}            `json:"data,omitempty"`       // 响应数据
	RequestID string                 `json:"request_id,omitempty"` // 请求ID
	Timestamp string                 `json:"timestamp,omitempty"`  // 时间戳
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Code:      0,
		Message:   "SUCCESS",
		MessageZH: "操作成功",
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// WithRequestID 添加请求ID
func (r *SuccessResponse) WithRequestID(requestID string) *SuccessResponse {
	r.RequestID = requestID
	return r
}
