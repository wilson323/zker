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
	"fmt"
	"strings"

	codepkg "github.com/coze-dev/coze-studio/backend/pkg/errorx/code"
)

// ErrMsgByCode 根据错误码获取错误消息
// 优先返回中文消息，如果没有则返回英文消息
func ErrMsgByCode(code int32) string {
	// 1. 尝试从已注册的CodeDefinition读取
	if def := codepkg.GetCodeDefinition(code); def != nil {
		// 优先返回中文消息（如果有）
		if def.MessageZH != "" {
			return def.MessageZH
		}
		// 否则返回英文消息
		if def.Message != "" {
			return def.Message
		}
	}

	// 2. 使用范围映射生成消息
	return generateChineseMessageByRange(code)
}

// GetHTTPStatusForError 根据错误码获取建议的HTTP状态码
// 用于API响应时确定HTTP状态码
func GetHTTPStatusForError(code int32) int {
	// 1. 如果有精确映射，使用精确映射
	if status, ok := HTTPStatusMapping[code]; ok {
		return status
	}

	// 2. 按错误码段进行范围映射
	return mapHTTPStatusByErrorCode(code)
}

// FormatError 格式化错误消息，替换占位符
// 使用示例：FormatError(ErrTenantNotFoundCode, map[string]interface{}{"tenant_id": "123"})
// 返回："租户不存在: tenant_id=123"
func FormatError(code int32, params map[string]interface{}) string {
	msgTemplate := ErrMsgByCode(code)

	// 如果没有参数，直接返回模板
	if len(params) == 0 {
		return msgTemplate
	}

	// 替换占位符 {key} 为实际值
	result := msgTemplate
	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}

	return result
}

// FormatErrorWithCode 格式化错误消息，包含错误码
// 使用示例：FormatErrorWithCode(ErrTenantNotFoundCode, map[string]interface{}{"tenant_id": "123"})
// 返回："[200000001] 租户不存在: tenant_id=123"
func FormatErrorWithCode(code int32, params map[string]interface{}) string {
	msg := FormatError(code, params)
	return fmt.Sprintf("[%d] %s", code, msg)
}

// IsNotFoundError 判断是否为"不存在"类型的错误
// 错误码范围：xxx xxx 001-009
func IsNotFoundError(code int32) bool {
	errorType := code % 1000
	return errorType >= 1 && errorType <= 9
}

// IsAlreadyExistsError 判断是否为"已存在"类型的错误
// 错误码范围：xxx xxx 010-019
func IsAlreadyExistsError(code int32) bool {
	errorType := code % 1000
	return errorType >= 10 && errorType <= 19
}

// IsInvalidParamError 判断是否为"参数无效"类型的错误
// 错误码范围：xxx xxx 020-039
func IsInvalidParamError(code int32) bool {
	errorType := code % 1000
	return errorType >= 20 && errorType <= 39
}

// IsPermissionDeniedError 判断是否为"权限拒绝"类型的错误
// 错误码范围：xxx xxx 040-049
func IsPermissionDeniedError(code int32) bool {
	errorType := code % 1000
	return errorType >= 40 && errorType <= 49
}

// IsOperationFailedError 判断是否为"操作失败"类型的错误
// 错误码范围：xxx xxx 050-089
func IsOperationFailedError(code int32) bool {
	errorType := code % 1000
	return errorType >= 50 && errorType <= 89
}

// IsRetryableError 判断错误是否可重试
// 一般操作失败（050-089）可以重试，参数错误（020-039）和权限错误（040-049）不可重试
func IsRetryableError(code int32) bool {
	// 参数错误不可重试
	if IsInvalidParamError(code) {
		return false
	}

	// 权限错误不可重试
	if IsPermissionDeniedError(code) {
		return false
	}

	// 已存在错误不可重试
	if IsAlreadyExistsError(code) {
		return false
	}

	// 不存在错误通常不可重试（除非是竞态条件）
	if IsNotFoundError(code) {
		return false
	}

	// 操作失败可以重试
	if IsOperationFailedError(code) {
		return true
	}

	// 其他错误默认不可重试
	return false
}

// GetModuleFromErrorCode 从错误码提取模块名称
// 例如：200000001 -> "租户", 201001001 -> "Bot"
func GetModuleFromErrorCode(code int32) string {
	return getModuleName(code)
}

// GetErrorTypeFromErrorCode 从错误码提取错误类型
// 例如：200000001 -> "NotFound", 200000010 -> "AlreadyExists"
func GetErrorTypeFromErrorCode(code int32) string {
	errorType := code % 1000

	switch {
	case errorType >= 1 && errorType <= 9:
		return "NotFound"
	case errorType >= 10 && errorType <= 19:
		return "AlreadyExists"
	case errorType >= 20 && errorType <= 39:
		return "InvalidParam"
	case errorType >= 40 && errorType <= 49:
		return "PermissionDenied"
	case errorType >= 50 && errorType <= 89:
		return "OperationFailed"
	case errorType >= 90 && errorType <= 99:
		return "Other"
	default:
		return "Unknown"
	}
}

// ErrorCategory 错误分类
type ErrorCategory string

const (
	ErrorCategoryNotFound        ErrorCategory = "NOT_FOUND"         // 资源不存在
	ErrorCategoryAlreadyExists   ErrorCategory = "ALREADY_EXISTS"    // 资源已存在
	ErrorCategoryInvalidParam    ErrorCategory = "INVALID_PARAM"     // 参数无效
	ErrorCategoryPermission      ErrorCategory = "PERMISSION"        // 权限不足
	ErrorCategoryOperationFailed ErrorCategory = "OPERATION_FAILED"  // 操作失败
	ErrorCategoryOther           ErrorCategory = "OTHER"             // 其他错误
	ErrorCategoryUnknown         ErrorCategory = "UNKNOWN"           // 未知错误
)

// GetErrorCategory 获取错误的分类
func GetErrorCategory(code int32) ErrorCategory {
	if IsNotFoundError(code) {
		return ErrorCategoryNotFound
	}
	if IsAlreadyExistsError(code) {
		return ErrorCategoryAlreadyExists
	}
	if IsInvalidParamError(code) {
		return ErrorCategoryInvalidParam
	}
	if IsPermissionDeniedError(code) {
		return ErrorCategoryPermission
	}
	if IsOperationFailedError(code) {
		return ErrorCategoryOperationFailed
	}
	errorType := code % 1000
	if errorType >= 90 && errorType <= 99 {
		return ErrorCategoryOther
	}
	return ErrorCategoryUnknown
}

// ShouldLogError 判断错误是否应该记录日志
// NOT_FOUND 和 ALREADY_EXISTS 通常是预期的业务错误，不需要记录错误日志
// PERMISSION 和 OPERATION_FAILED 需要记录
func ShouldLogError(code int32) bool {
	category := GetErrorCategory(code)
	return category == ErrorCategoryPermission ||
	       category == ErrorCategoryOperationFailed ||
	       category == ErrorCategoryUnknown
}

// ShouldAlert 判断错误是否应该触发告警
// 只有影响稳定性的错误才应该触发告警
func ShouldAlert(code int32) bool {
	// 从 CodeDefinition 读取 affect_stability
	if def := codepkg.GetCodeDefinition(code); def != nil {
		return def.IsAffectStability
	}

	// 默认：操作失败和未知错误需要告警
	category := GetErrorCategory(code)
	return category == ErrorCategoryOperationFailed ||
	       category == ErrorCategoryUnknown
}

// ClientSafeError 判断错误是否对客户端安全
// 某些内部错误（如数据库连接失败）不应该直接暴露给客户端
func IsClientSafeError(code int32) bool {
	// 参数错误、权限错误、不存在错误是安全的
	if IsInvalidParamError(code) || IsPermissionDeniedError(code) || IsNotFoundError(code) {
		return true
	}

	// 已存在错误是安全的
	if IsAlreadyExistsError(code) {
		return true
	}

	// 操作失败需要判断具体类型
	if IsOperationFailedError(code) {
		// 某些操作失败是安全的（如配额超限）
		// 这里可以根据具体错误码判断
		return true
	}

	// 其他错误默认不安全，不应该暴露详细信息给客户端
	return false
}

// GetSafeErrorMessage 获取对客户端安全的错误消息
// 如果错误不安全，返回通用消息
func GetSafeErrorMessage(code int32) string {
	if IsClientSafeError(code) {
		return ErrMsgByCode(code)
	}

	// 不安全的错误返回通用消息
	return "操作失败，请稍后重试"
}
