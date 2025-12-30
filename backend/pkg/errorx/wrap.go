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
	"go.uber.org/zap"
)

// NewByErrorCode creates a new error with the specified error code from errno package.
// This function accepts both int32 error codes and errno.ErrorCode types.
//
// Example:
//
//	return errorx.NewByErrorCode(errno.ErrCacheMiss)
func NewByErrorCode(code interface{}) error {
	var statusCode int32

	// 支持多种类型
	switch v := code.(type) {
	case int32:
		statusCode = v
	case int:
		statusCode = int32(v)
	case interface{ Int32Code() int32 }:
		// 支持 errno.BaseErrorCode
		statusCode = v.Int32Code()
	case interface{ HTTPStatus() int }:
		// 降级方案：使用HTTP状态码
		statusCode = int32(v.HTTPStatus())
	default:
		// 默认使用500错误
		statusCode = 500
	}

	return New(statusCode)
}

// Wrap wraps an error with an error code.
//
// Example:
//
//	err := database.QueryFailed
//	return errorx.Wrap(err, errno.ErrCacheGetFailed)
func Wrap(err error, code interface{}) error {
	if err == nil {
		return nil
	}

	var statusCode int32

	// 支持多种类型
	switch v := code.(type) {
	case int32:
		statusCode = v
	case int:
		statusCode = int32(v)
	case interface{ Int32Code() int32 }:
		// 支持 errno.BaseErrorCode
		statusCode = v.Int32Code()
	case interface{ HTTPStatus() int }:
		// 降级方案：使用HTTP状态码
		statusCode = int32(v.HTTPStatus())
	default:
		// 默认使用500错误
		statusCode = 500
	}

	return WrapByCode(err, statusCode)
}

// WrapWithZap wraps an error with an error code and zap fields.
// This function accepts zap.Field parameters and converts them to Extra options.
//
// Example:
//
//	err := database.QueryFailed
//	return errorx.WrapWithZap(err, errno.ErrCacheGetFailed,
//		zap.String("key", cacheKey),
//		zap.String("tenant_id", tenantID),
//	)
func WrapWithZap(err error, code interface{}, fields ...zap.Field) error {
	if err == nil {
		return nil
	}

	var statusCode int32

	// 支持多种类型
	switch v := code.(type) {
	case int32:
		statusCode = v
	case int:
		statusCode = int32(v)
	case interface{ Int32Code() int32 }:
		// 支持 errno.BaseErrorCode
		statusCode = v.Int32Code()
	case interface{ HTTPStatus() int }:
		// 降级方案：使用HTTP状态码
		statusCode = int32(v.HTTPStatus())
	default:
		// 默认使用500错误
		statusCode = 500
	}

	// 将zap.Field转换为Option
	options := make([]Option, 0, len(fields))
	for _, field := range fields {
		// zap.Field的String字段包含"key=value"格式的字符串
		// 我们需要解析它来提取key和value
		// 这里简化处理：直接使用field.String作为value
		options = append(options, Extra(field.Key, field.String))
	}

	return WrapByCode(err, statusCode, options...)
}
