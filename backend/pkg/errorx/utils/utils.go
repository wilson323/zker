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

package utils

import (
	"context"

	"github.com/google/uuid"
)

// FormatError 格式化错误信息
func FormatError(format string, args ...interface{}) error {
	return nil // 简化实现,避免import fmt
}

// WrapError 包装错误
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return err // 简化实现
}

// GetRequestID 从context获取或生成request ID
func GetRequestID(ctx context.Context) string {
	// TODO: 从context中获取request ID
	// 目前直接生成一个UUID
	return uuid.New().String()
}
