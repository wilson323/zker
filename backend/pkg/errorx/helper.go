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
	"fmt"
	"net/http"
)

// GetMessageByCode 根据错误码获取默认错误消息
// 遵循KISS原则:简单映射,不过度设计
func GetMessageByCode(code int32) string {
	// 根据错误码范围生成默认消息
	million := int32(1000000)
	codeSegment := code / million

	switch codeSegment {
	case 2:
		return "Tenant error"
	case 3:
		return "Quota error"
	case 4:
		return "Subscription error"
	case 201:
		return "Bot error"
	case 202:
		return "Conversation error"
	case 203:
		return "Workflow error"
	case 204:
		return "Knowledge error"
	default:
		return fmt.Sprintf("Error %d", code)
	}
}

// GetHTTPStatusByCode 根据错误码获取HTTP状态码
// 使用范围映射避免为每个错误码单独映射(遵循DRY原则)
func GetHTTPStatusByCode(code int32) int {
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

	// 默认 → 500
	return http.StatusInternalServerError
}

// EnhancedError 增强的错误结构接口
// 避免直接引用errno包的EnhancedError,防止循环依赖
type EnhancedError interface {
	Code() int32
	Message() string
	MessageZH() string
	Details() map[string]interface{}
	RequestID() string
	TenantID() string
}
