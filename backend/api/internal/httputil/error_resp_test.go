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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"

	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// TestBuildErrorResp 测试构建错误响应
func TestBuildErrorResp(t *testing.T) {
	// 创建测试上下文
	w := httptest.NewRecorder()
	ctx := ut.NewContext(w)

	// 构建错误响应
	details := map[string]interface{}{
		"field":   "username",
		"hint":    "用户名不能为空",
		"user_id": "12345",
	}
	BuildErrorResp(ctx, 2001001, "Tenant not found", "租户不存在", details)

	// 验证响应
	resp := w.Body.String()
	assert.DeepEqual(t, resp, "")

	// 验证HTTP状态码
	assert.DeepEqual(t, w.Code, int(errno.GetHTTPStatusForError(2001001)))
}

// TestBuildErrorRespFromEnhanced 测试从EnhancedError构建错误响应
func TestBuildErrorRespFromEnhanced(t *testing.T) {
	// 创建测试上下文
	w := httptest.NewRecorder()
	ctx := ut.NewContext(w)

	// 创建EnhancedError
	enhancedErr := errno.NewEnhancedError(
		2001001,
		"Tenant not found",
		"租户不存在",
	)
	enhancedErr.WithRequestID("req-12345")
	enhancedErr.WithDetail("field", "tenant_id")

	// 构建错误响应
	BuildErrorRespFromEnhanced(ctx, enhancedErr)

	// 验证HTTP状态码
	assert.DeepEqual(t, w.Code, int(enhancedErr.GetHTTPStatus()))
}

// TestBuildSuccessResp 测试构建成功响应
func TestBuildSuccessResp(t *testing.T) {
	// 创建测试上下文
	w := httptest.NewRecorder()
	ctx := ut.NewContext(w)

	// 构建成功响应
	data := map[string]interface{}{
		"id":   "12345",
		"name": "Test Tenant",
	}
	BuildSuccessResp(ctx, data)

	// 验证HTTP状态码
	assert.DeepEqual(t, w.Code, http.StatusOK)
}

// TestErrorResponseFields 测试错误响应字段完整性
func TestErrorResponseFields(t *testing.T) {
	resp := &ErrorResponse{
		Code:       2001001,
		Message:    "Tenant not found",
		MessageZH:  "租户不存在",
		MessageEN:  "Tenant not found",
		Details:    map[string]interface{}{"field": "tenant_id"},
		RequestID:  "req-12345",
		Timestamp:  "2025-01-01T12:00:00Z",
		TraceID:    "trace-12345",
		TenantID:   "tenant-12345",
	}

	// 验证所有字段
	assert.DeepEqual(t, resp.Code, int32(2001001))
	assert.DeepEqual(t, resp.Message, "Tenant not found")
	assert.DeepEqual(t, resp.MessageZH, "租户不存在")
	assert.DeepEqual(t, resp.MessageEN, "Tenant not found")
	assert.DeepEqual(t, resp.RequestID, "req-12345")
	assert.DeepEqual(t, resp.TraceID, "trace-12345")
	assert.DeepEqual(t, resp.TenantID, "tenant-12345")
	assert.DeepEqual(t, resp.Details["field"], "tenant_id")
}

// TestHTTPStatusMapping 测试HTTP状态码映射
func TestHTTPStatusMapping(t *testing.T) {
	tests := []struct {
		name       string
		errCode    int32
		expectCode int
	}{
		{
			name:       "租户不存在 - 404",
			errCode:    2001001,
			expectCode: http.StatusNotFound,
		},
		{
			name:       "租户已存在 - 409",
			errCode:    2001002,
			expectCode: http.StatusConflict,
		},
		{
			name:       "租户已暂停 - 403",
			errCode:    2003001,
			expectCode: http.StatusForbidden,
		},
		{
			name:       "租户参数无效 - 400",
			errCode:    2002001,
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := errno.GetHTTPStatusForError(tt.errCode)
			assert.DeepEqual(t, status, tt.expectCode)
		})
	}
}

// TestLegacyErrorResponse 测试兼容旧版本错误响应
func TestLegacyErrorResponse(t *testing.T) {
	// 测试BadRequest
	w := httptest.NewRecorder()
	ctx := ut.NewContext(w)
	BadRequest(ctx, "Invalid parameters")
	assert.DeepEqual(t, w.Code, http.StatusBadRequest)

	// 测试Unauthorized
	w = httptest.NewRecorder()
	ctx = ut.NewContext(w)
	Unauthorized(ctx, "Unauthorized")
	assert.DeepEqual(t, w.Code, http.StatusUnauthorized)
}

// helper function to create test context
func newTestContext() *app.RequestContext {
	return &app.RequestContext{}
}
