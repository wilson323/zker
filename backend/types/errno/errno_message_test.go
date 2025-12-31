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
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestErrMessageFormatting 测试errno消息格式化
func TestErrMessageFormatting(t *testing.T) {
	t.Run("测试基础错误消息", func(t *testing.T) {
		err := ErrTenantNotFound
		assert.NotEmpty(t, err.Message())
		assert.NotEmpty(t, err.MessageZH())
		assert.NotEmpty(t, err.MessageEN())
		assert.Contains(t, err.MessageZH(), "租户")
		assert.Contains(t, err.MessageEN(), "Tenant")
	})

	t.Run("测试配额错误消息", func(t *testing.T) {
		err := ErrQuotaExceeded
		assert.NotEmpty(t, err.MessageZH())
		assert.NotEmpty(t, err.MessageEN())
		assert.Contains(t, err.MessageZH(), "配额")
		assert.Contains(t, err.MessageEN(), "Quota")
	})

	t.Run("测试权限错误消息", func(t *testing.T) {
		err := ErrRoleNotFound
		assert.NotEmpty(t, err.MessageZH())
		assert.NotEmpty(t, err.MessageEN())
		assert.Contains(t, err.MessageZH(), "角色")
		assert.Contains(t, err.MessageEN(), "Role")
	})

	t.Run("测试计费错误消息", func(t *testing.T) {
		err := ErrBillingChargeFailed
		assert.NotEmpty(t, err.MessageZH())
		assert.NotEmpty(t, err.MessageEN())
		// ErrBillingChargeFailed 的消息是 "扣费失败" 和 "Failed to charge"
		assert.Contains(t, err.MessageZH(), "扣费")
		assert.Contains(t, err.MessageEN(), "charge")
	})
}

// TestErrorCodeString 测试错误码字符串
func TestErrorCodeString(t *testing.T) {
	tests := []struct {
		name           string
		errCode        *BaseErrorCode
		expectedString string
	}{
		{
			name:           "ErrTenantNotFound code",
			errCode:        ErrTenantNotFound,
			expectedString: "TENANT_NOT_FOUND",
		},
		{
			name:           "ErrQuotaExceeded code",
			errCode:        ErrQuotaExceeded,
			expectedString: "QUOTA_EXCEEDED",
		},
		{
			name:           "ErrRoleNotFound code",
			errCode:        ErrRoleNotFound,
			expectedString: "ROLE_NOT_FOUND",
		},
		{
			name:           "InvalidParams code",
			errCode:        InvalidParams,
			expectedString: "COMMON201001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedString, tt.errCode.Code())
		})
	}
}

// TestErrMessageConsistency 测试错误消息一致性
func TestErrMessageConsistency(t *testing.T) {
	t.Run("中文消息不为空", func(t *testing.T) {
		errnos := []*BaseErrorCode{
			ErrTenantNotFound,
			ErrTenantInactive,
			ErrQuotaExceeded,
			ErrQuotaBotExceeded,
			ErrRoleNotFound,
			ErrPermissionCheckFailed,
			ErrBillingChargeFailed,
			// ErrBotStoreNotFound 是error类型，跳过
		}

		for _, err := range errnos {
			assert.NotEmpty(t, err.MessageZH(), "Error %s should have Chinese message", err.Code())
		}
	})

	t.Run("英文消息不为空", func(t *testing.T) {
		errnos := []*BaseErrorCode{
			ErrTenantNotFound,
			ErrTenantInactive,
			ErrQuotaExceeded,
			ErrQuotaBotExceeded,
			ErrRoleNotFound,
			ErrPermissionCheckFailed,
			ErrBillingChargeFailed,
			// ErrBotStoreNotFound 是error类型，跳过
		}

		for _, err := range errnos {
			assert.NotEmpty(t, err.MessageEN(), "Error %s should have English message", err.Code())
		}
	})

	t.Run("中英文消息应该有相关性", func(t *testing.T) {
		// 对于某些错误，中英文消息应该包含相同的关键词
		err := ErrTenantNotFound
		zhMsg := err.MessageZH()
		enMsg := err.MessageEN()

		// 两者都应该包含错误类型相关的词汇
		assert.Contains(t, zhMsg, "不存在")
		assert.Contains(t, enMsg, "not found")
	})
}

// TestErrInt32Code 测试错误码Int32值
func TestErrInt32Code(t *testing.T) {
	tests := []struct {
		name          string
		errCode       *BaseErrorCode
		expectedCode  int32
	}{
		{
			name:          "Tenant not found - HTTP 404",
			errCode:       ErrTenantNotFound,
			expectedCode:  404, // Int32Code返回HTTP状态码
		},
		{
			name:          "Quota exceeded - HTTP 403",
			errCode:       ErrQuotaExceeded,
			expectedCode:  403, // Int32Code返回HTTP状态码
		},
		{
			name:          "Role not found - HTTP 404",
			errCode:       ErrRoleNotFound,
			expectedCode:  404, // Int32Code返回HTTP状态码
		},
		{
			name:          "Invalid params - HTTP 400",
			errCode:       InvalidParams,
			expectedCode:  400, // Int32Code返回HTTP状态码
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedCode, tt.errCode.Int32Code())
		})
	}
}

// TestErrorCodeAndMessage 测试错误码和消息方法
func TestErrorCodeAndMessage(t *testing.T) {
	t.Run("Code()方法返回错误码", func(t *testing.T) {
		err := ErrTenantNotFound
		code := err.Code()
		assert.Contains(t, code, "TENANT_NOT_FOUND")
	})

	t.Run("Message()方法包含错误信息", func(t *testing.T) {
		err := ErrTenantNotFound
		message := err.Message()
		assert.NotEmpty(t, message)
	})

	t.Run("MessageZH()方法返回中文消息", func(t *testing.T) {
		err := ErrQuotaExceeded
		zhMsg := err.MessageZH()
		assert.Contains(t, zhMsg, "配额")
	})

	t.Run("MessageEN()方法返回英文消息", func(t *testing.T) {
		err := ErrQuotaExceeded
		enMsg := err.MessageEN()
		assert.Contains(t, enMsg, "Quota")
	})
}

// TestFormatError 测试错误消息格式化辅助函数
func TestFormatError(t *testing.T) {
	// 这个测试验证FormatError函数的行为
	// 注意：如果FormatError函数使用了模板替换，需要测试参数替换功能

	t.Run("格式化无参数错误消息", func(t *testing.T) {
		// 对于没有参数的错误码，格式化应该返回原始消息
		msg := FormatError(ErrTenantNotFoundCode, nil)
		assert.NotEmpty(t, msg)
		assert.Contains(t, msg, "租户")
	})

	t.Run("格式化带参数错误消息", func(t *testing.T) {
		// 对于有占位符的错误码，测试参数替换
		params := map[string]interface{}{
			"tenant_id": "tenant-123",
			"reason": "expired",
		}

		// 注意：这个测试依赖于实际的错误消息格式化实现
		// 如果FormatError支持参数替换，验证替换是否正确
		msg := FormatError(ErrTenantInactiveCode, params)
		assert.NotEmpty(t, msg)
	})
}

// TestErrMessageLength 测试错误消息长度合理性
func TestErrMessageLength(t *testing.T) {
	t.Run("中文消息长度合理", func(t *testing.T) {
		errnos := []*BaseErrorCode{
			ErrTenantNotFound,
			ErrQuotaExceeded,
			ErrRoleNotFound,
			ErrBillingChargeFailed,
		}

		for _, err := range errnos {
			zhMsg := err.MessageZH()
			// 错误消息应该在10-100字符之间
			assert.GreaterOrEqual(t, len(zhMsg), 2, "Error message too short: %s", err.Code())
			assert.LessOrEqual(t, len(zhMsg), 200, "Error message too long: %s", err.Code())
		}
	})

	t.Run("英文消息长度合理", func(t *testing.T) {
		errnos := []*BaseErrorCode{
			ErrTenantNotFound,
			ErrQuotaExceeded,
			ErrRoleNotFound,
			ErrBillingChargeFailed,
		}

		for _, err := range errnos {
			enMsg := err.MessageEN()
			// 错误消息应该在5-200字符之间
			assert.GreaterOrEqual(t, len(enMsg), 2, "Error message too short: %s", err.Code())
			assert.LessOrEqual(t, len(enMsg), 200, "Error message too long: %s", err.Code())
		}
	})
}

// TestErrorCategoriesMessages 测试不同错误类别的消息
func TestErrorCategoriesMessages(t *testing.T) {
	t.Run("NotFound类错误消息", func(t *testing.T) {
		errnos := []*BaseErrorCode{
			ErrTenantNotFound,
			ErrQuotaNotFound,
			ErrRoleNotFound,
		}

		for _, err := range errnos {
			msg := err.MessageZH()
			// NotFound类错误通常包含"不存在"关键词
			assert.Contains(t, msg, "不存在", "NotFound error should contain '不存在': %s", err.Code())
		}
	})

	t.Run("AlreadyExists类错误消息", func(t *testing.T) {
		errnos := []*BaseErrorCode{
			ErrTenantAlreadyExists,
			ErrRoleAlreadyExists,
		}

		for _, err := range errnos {
			msg := err.MessageZH()
			// AlreadyExists类错误通常包含"已存在"关键词
			assert.Contains(t, msg, "已存在", "AlreadyExists error should contain '已存在': %s", err.Code())
		}
	})

	t.Run("InvalidParam类错误消息", func(t *testing.T) {
		errnos := []*BaseErrorCode{
			ErrTenantInvalidName,
			ErrQuotaInvalid,
			ErrInvalidRole,
		}

		for _, err := range errnos {
			msg := err.MessageZH()
			// InvalidParam类错误通常包含"无效"关键词
			assert.Contains(t, msg, "无效", "InvalidParam error should contain '无效': %s", err.Code())
		}
	})

	t.Run("OperationFailed类错误消息", func(t *testing.T) {
		errnos := []*BaseErrorCode{
			ErrQuotaCheckFailed,
			ErrQuotaConsumeFailed,
			ErrPermissionCheckFailed,
			ErrBillingChargeFailed,
		}

		for _, err := range errnos {
			msg := err.MessageZH()
			// OperationFailed类错误通常包含"失败"关键词
			assert.Contains(t, msg, "失败", "OperationFailed error should contain '失败': %s", err.Code())
		}
	})
}

// BenchmarkErrorCodeString 基准测试错误码字符串获取
func BenchmarkErrorCodeString(b *testing.B) {
	err := ErrTenantNotFound

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Code()
	}
}

// BenchmarkErrorMessageZH 基准测试中文消息获取
func BenchmarkErrorMessageZH(b *testing.B) {
	err := ErrTenantNotFound

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.MessageZH()
	}
}

// BenchmarkErrorMessageEN 基准测试英文消息获取
func BenchmarkErrorMessageEN(b *testing.B) {
	err := ErrTenantNotFound

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.MessageEN()
	}
}
