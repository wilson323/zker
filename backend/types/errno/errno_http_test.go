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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHTTPStatusMappingDetailed 测试HTTP状态码映射详细版
func TestHTTPStatusMappingDetailed(t *testing.T) {
	t.Run("NotFound错误映射到404", func(t *testing.T) {
		tests := []struct {
			name  string
			err   *BaseErrorCode
			code  int32
		}{
			{"Tenant not found", ErrTenantNotFound, 2001001},
			{"Quota not found", ErrQuotaNotFound, 300000001},
			{"Role not found", ErrRoleNotFound, 108010001},
			{"Billing balance not found", ErrBillingBalanceNotFound, 310000001},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// 验证错误码尾段是001-009 (NotFound)
				assert.Equal(t, http.StatusNotFound, tt.err.HTTPStatus())
			})
		}
	})

	t.Run("AlreadyExists错误映射到409", func(t *testing.T) {
		tests := []struct {
			name string
			err  *BaseErrorCode
		}{
			{"Tenant already exists", ErrTenantAlreadyExists},
			{"Role already exists", ErrRoleAlreadyExists},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, http.StatusConflict, tt.err.HTTPStatus())
			})
		}
	})

	t.Run("InvalidParam错误映射到400", func(t *testing.T) {
		tests := []struct {
			name string
			err  *BaseErrorCode
		}{
			{"Tenant invalid name", ErrTenantInvalidName},
			{"Tenant invalid config", ErrTenantInvalidConfig},
			{"Quota invalid", ErrQuotaInvalid},
			{"Quota invalid resource", ErrQuotaInvalidResource},
			{"Invalid role", ErrInvalidRole},
			{"Permission invalid param", ErrPermissionInvalidParam},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, http.StatusBadRequest, tt.err.HTTPStatus())
			})
		}
	})

	t.Run("PermissionDenied错误映射到403", func(t *testing.T) {
		tests := []struct {
			name string
			err  *BaseErrorCode
		}{
			{"Tenant inactive", ErrTenantInactive},
			{"Tenant quota exceeded", ErrTenantQuotaExceeded},
			{"Quota exceeded", ErrQuotaExceeded},
			{"Quota bot exceeded", ErrQuotaBotExceeded},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, http.StatusForbidden, tt.err.HTTPStatus())
			})
		}
	})

	t.Run("OperationFailed错误映射到500", func(t *testing.T) {
		tests := []struct {
			name string
			err  *BaseErrorCode
		}{
			{"Tenant migration failed", ErrTenantMigrationFailed},
			{"Quota check failed", ErrQuotaCheckFailed},
			{"Quota consume failed", ErrQuotaConsumeFailed},
			{"Quota reset failed", ErrQuotaResetFailed},
			{"Permission check failed", ErrPermissionCheckFailed},
			{"Billing charge failed", ErrBillingChargeFailed},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, http.StatusInternalServerError, tt.err.HTTPStatus())
			})
		}
	})

	t.Run("特殊HTTP状态码", func(t *testing.T) {
		tests := []struct {
			name           string
			err            *BaseErrorCode
			expectedStatus int
		}{
			{
				name:           "Quota storage exceeded - 507 Insufficient Storage",
				err:            ErrQuotaStorageExceeded,
				expectedStatus: 507, // StatusInsufficientStorage
			},
			{
				name:           "Quota API call exceeded - 429 Too Many Requests",
				err:            ErrQuotaAPICallExceeded,
				expectedStatus: 429, // StatusTooManyRequests
			},
			{
				name:           "Billing insufficient balance - 402 Payment Required",
				err:            ErrBillingInsufficientBalance,
				expectedStatus: 402, // StatusPaymentRequired
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expectedStatus, tt.err.HTTPStatus())
			})
		}
	})
}

// TestGetHTTPStatusForError 测试GetHTTPStatusForError函数
func TestGetHTTPStatusForError(t *testing.T) {
	t.Run("基于错误码尾段的映射", func(t *testing.T) {
		tests := []struct {
			name       string
			errCode    int32
			expectCode int
		}{
			// NotFound: 001-009 → 404
			{"2001001 (001)", 2001001, 404},
			{"206000001 (001)", 206000001, 404},
			{"300000001 (001)", 300000001, 404},

			// AlreadyExists: 010-019 → 409
			{"2001010 (010)", 2001010, 409},
			{"108010010 (010)", 108010010, 409},

			// InvalidParam: 020-039 → 400
			{"2001020 (020)", 2001020, 400},
			{"2001021 (021)", 2001021, 400},

			// PermissionDenied: 040-049 → 403
			{"2001040 (040)", 2001040, 403},
			{"300000041 (041)", 300000041, 403},

			// OperationFailed: 050-089 → 500
			{"2001050 (050)", 2001050, 500},
			{"300000004 (004)", 300000004, 500},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				status := GetHTTPStatusForError(tt.errCode)
				assert.Equal(t, tt.expectCode, status)
			})
		}
	})

	t.Run("未知错误码默认500", func(t *testing.T) {
		// 使用一个不在任何预定义范围内的错误码
		unknownCode := int32(999999999)
		status := GetHTTPStatusForError(unknownCode)
		assert.Equal(t, 500, status)
	})
}

// TestHTTPStatusCategories 测试HTTP状态码分类
func TestHTTPStatusCategories(t *testing.T) {
	t.Run("4xx客户端错误", func(t *testing.T) {
		// 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, ErrTenantInvalidName.HTTPStatus())
		assert.Equal(t, http.StatusBadRequest, ErrQuotaInvalid.HTTPStatus())
		assert.Equal(t, http.StatusBadRequest, ErrInvalidRole.HTTPStatus())

		// 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, Unauthorized.HTTPStatus())

		// 403 Forbidden
		assert.Equal(t, http.StatusForbidden, ErrTenantInactive.HTTPStatus())
		assert.Equal(t, http.StatusForbidden, ErrQuotaExceeded.HTTPStatus())
		assert.Equal(t, http.StatusForbidden, Forbidden.HTTPStatus())

		// 404 Not Found
		assert.Equal(t, http.StatusNotFound, ErrTenantNotFound.HTTPStatus())
		assert.Equal(t, http.StatusNotFound, ErrQuotaNotFound.HTTPStatus())
		assert.Equal(t, http.StatusNotFound, NotFound.HTTPStatus())

		// 409 Conflict
		assert.Equal(t, http.StatusConflict, ErrTenantAlreadyExists.HTTPStatus())
		assert.Equal(t, http.StatusConflict, ErrRoleAlreadyExists.HTTPStatus())
	})

	t.Run("5xx服务器错误", func(t *testing.T) {
		// 500 Internal Server Error
		assert.Equal(t, http.StatusInternalServerError, ErrTenantMigrationFailed.HTTPStatus())
		assert.Equal(t, http.StatusInternalServerError, ErrQuotaCheckFailed.HTTPStatus())
		assert.Equal(t, http.StatusInternalServerError, InternalError.HTTPStatus())
		assert.Equal(t, http.StatusInternalServerError, DB500001.HTTPStatus())
	})

	t.Run("其他状态码", func(t *testing.T) {
		// 429 Too Many Requests
		assert.Equal(t, 429, ErrQuotaAPICallExceeded.HTTPStatus())

		// 507 Insufficient Storage
		assert.Equal(t, 507, ErrQuotaStorageExceeded.HTTPStatus())

		// 402 Payment Required
		assert.Equal(t, 402, ErrBillingInsufficientBalance.HTTPStatus())
	})
}

// TestHTTPStatusCodeCoverage 测试所有错误码都有HTTP状态码
func TestHTTPStatusCodeCoverage(t *testing.T) {
	t.Run("验证所有便捷变量都有HTTP状态码", func(t *testing.T) {
		errnos := []*BaseErrorCode{
			// Tenant
			ErrTenantNotFound,
			ErrTenantAlreadyExists,
			ErrTenantInactive,
			ErrTenantInvalidName,
			ErrTenantInvalidConfig,
			ErrTenantMigrationFailed,

			// Quota
			ErrQuotaNotFound,
			ErrQuotaInvalid,
			ErrQuotaExceeded,
			ErrQuotaBotExceeded,
			ErrQuotaAPICallExceeded,
			ErrQuotaStorageExceeded,
			ErrQuotaCheckFailed,

			// Permission
			ErrRoleNotFound,
			ErrRoleAlreadyExists,
			ErrInvalidRole,
			ErrPermissionCheckFailed,

			// Billing
			ErrBillingChargeFailed,
			ErrBillingBalanceNotFound,
			ErrBillingInsufficientBalance,

			// Common
			Success,
			InvalidParams,
			Unauthorized,
			Forbidden,
			NotFound,
			InternalError,
		}

		for _, err := range errnos {
			status := err.HTTPStatus()
			assert.Greater(t, status, 0, "Error %s should have valid HTTP status code", err.Code())
			assert.LessOrEqual(t, status, 599, "Error %s HTTP status code too high", err.Code())
		}
	})
}

// TestHTTPStatusMappingConsistency 测试HTTP状态码映射一致性
func TestHTTPStatusMappingConsistency(t *testing.T) {
	t.Run("相同错误类别映射到相同状态码", func(t *testing.T) {
		// 所有NotFound类错误都应该映射到404
		notFoundErrs := []int32{2001001, 206000001, 300000001, 108010001}
		for _, code := range notFoundErrs {
			status := GetHTTPStatusForError(code)
			assert.Equal(t, http.StatusNotFound, status,
				"NotFound error %d should map to 404", code)
		}

		// 所有AlreadyExists类错误都应该映射到409
		alreadyExistsErrs := []int32{2001010, 108010010}
		for _, code := range alreadyExistsErrs {
			status := GetHTTPStatusForError(code)
			assert.Equal(t, http.StatusConflict, status,
				"AlreadyExists error %d should map to 409", code)
		}
	})

	t.Run("精确映射优先于范围映射", func(t *testing.T) {
		// 某些错误码有精确的HTTP状态码映射
		// 这些精确映射应该优先于基于尾段的范围映射

		// 例如：ErrQuotaAPICallExceeded 精确映射到 429
		// 而不是基于尾段 001 映射到 404
		assert.Equal(t, 429, ErrQuotaAPICallExceeded.HTTPStatus())
	})
}

// BenchmarkGetHTTPStatusForError 基准测试HTTP状态码获取
func BenchmarkGetHTTPStatusForError(b *testing.B) {
	errCode := int32(2001001)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetHTTPStatusForError(errCode)
	}
}

// BenchmarkHTTPStatus 基准测试HTTPStatus方法
func BenchmarkHTTPStatus(b *testing.B) {
	err := ErrTenantNotFound

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.HTTPStatus()
	}
}

// TestHTTPStatusCodeRanges 测试HTTP状态码范围有效性
func TestHTTPStatusCodeRanges(t *testing.T) {
	t.Run("所有状态码在有效范围内", func(t *testing.T) {
		errnos := []*BaseErrorCode{
			Success,
			InvalidParams,
			Unauthorized,
			Forbidden,
			NotFound,
			InternalError,
			ErrTenantNotFound,
			ErrTenantAlreadyExists,
			ErrQuotaExceeded,
			ErrQuotaAPICallExceeded,
			ErrQuotaStorageExceeded,
			ErrBillingInsufficientBalance,
		}

		for _, err := range errnos {
			status := err.HTTPStatus()
			// HTTP状态码应该在100-599范围内
			assert.GreaterOrEqual(t, status, 100,
				"HTTP status %d for error %s is too low", status, err.Code())
			assert.LessOrEqual(t, status, 599,
				"HTTP status %d for error %s is too high", status, err.Code())
		}
	})

	t.Run("信息响应(1xx)", func(t *testing.T) {
		// 通常不使用1xx状态码
		// 这个测试确保我们没有错误码映射到1xx
		errnos := []*BaseErrorCode{
			ErrTenantNotFound,
			ErrQuotaExceeded,
			ErrRoleNotFound,
		}

		for _, err := range errnos {
			status := err.HTTPStatus()
			assert.NotEqual(t, status/100, 1,
				"Error %s should not map to 1xx status", err.Code())
		}
	})
}
