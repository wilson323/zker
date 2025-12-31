// backend/types/errno/errno_test.go
package errno

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTenantErrorCodes 测试租户错误码
func TestTenantErrorCodes(t *testing.T) {
	tests := []struct {
		name           string
		errCode        *BaseErrorCode
		expectedCode   string
		expectedZH     string
		expectedEN     string
		expectedStatus int
	}{
		{
			name:           "ErrTenantNotFound",
			errCode:        ErrTenantNotFound,
			expectedCode:   "TENANT_NOT_FOUND",
			expectedZH:     "租户不存在",
			expectedEN:     "Tenant not found",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "ErrTenantAlreadyExists",
			errCode:        ErrTenantAlreadyExists,
			expectedCode:   "TENANT_ALREADY_EXISTS",
			expectedZH:     "租户已存在",
			expectedEN:     "Tenant already exists",
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "ErrTenantInactive",
			errCode:        ErrTenantInactive,
			expectedCode:   "TENANT_INACTIVE",
			expectedZH:     "租户未激活",
			expectedEN:     "Tenant is inactive",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "ErrTenantQuotaExceeded",
			errCode:        ErrTenantQuotaExceeded,
			expectedCode:   "TENANT_QUOTA_EXCEEDED",
			expectedZH:     "租户配额已用尽",
			expectedEN:     "Tenant quota exceeded",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "ErrTenantInvalidName",
			errCode:        ErrTenantInvalidName,
			expectedCode:   "TENANT_INVALID_NAME",
			expectedZH:     "租户名称无效",
			expectedEN:     "Invalid tenant name",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "ErrTenantInvalidConfig",
			errCode:        ErrTenantInvalidConfig,
			expectedCode:   "TENANT_INVALID_CONFIG",
			expectedZH:     "租户配置无效",
			expectedEN:     "Invalid tenant config",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "ErrTenantMigrationFailed",
			errCode:        ErrTenantMigrationFailed,
			expectedCode:   "TENANT_MIGRATION_FAILED",
			expectedZH:     "租户迁移失败",
			expectedEN:     "Tenant migration failed",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "ErrTenantIsolationNotSupported",
			errCode:        ErrTenantIsolationNotSupported,
			expectedCode:   "TENANT_ISOLATION_NOT_SUPPORTED",
			expectedZH:     "不支持租户隔离",
			expectedEN:     "Tenant isolation not supported",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "ErrTenantSubscriptionExpired",
			errCode:        ErrTenantSubscriptionExpired,
			expectedCode:   "TENANT_SUBSCRIPTION_EXPIRED",
			expectedZH:     "租户订阅已过期",
			expectedEN:     "Tenant subscription expired",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedCode, tt.errCode.Code())
			assert.Equal(t, tt.expectedZH, tt.errCode.MessageZH())
			assert.Equal(t, tt.expectedEN, tt.errCode.MessageEN())
			assert.Equal(t, tt.expectedStatus, tt.errCode.HTTPStatus())
		})
	}
}

// TestDatabaseErrorCodes 测试数据库错误码
func TestDatabaseErrorCodes(t *testing.T) {
	tests := []struct {
		name           string
		errCode        *BaseErrorCode
		expectedCode   string
		expectedZH     string
		expectedEN     string
		expectedStatus int
	}{
		{
			name:           "DB500001 - 数据库连接失败",
			errCode:        DB500001,
			expectedCode:   "DB500001",
			expectedZH:     "数据库连接失败",
			expectedEN:     "Database connection failed",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "DB500002 - SQL执行失败",
			errCode:        DB500002,
			expectedCode:   "DB500002",
			expectedZH:     "SQL执行失败",
			expectedEN:     "SQL execution failed",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "DB500003 - 事务失败",
			errCode:        DB500003,
			expectedCode:   "DB500003",
			expectedZH:     "事务失败",
			expectedEN:     "Transaction failed",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "DB500004 - 死锁错误",
			errCode:        DB500004,
			expectedCode:   "DB500004",
			expectedZH:     "死锁错误",
			expectedEN:     "Database deadlock",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "DB500005 - 连接池耗尽",
			errCode:        DB500005,
			expectedCode:   "DB500005",
			expectedZH:     "连接池耗尽",
			expectedEN:     "Connection pool exhausted",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedCode, tt.errCode.Code())
			assert.Equal(t, tt.expectedZH, tt.errCode.MessageZH())
			assert.Equal(t, tt.expectedEN, tt.errCode.MessageEN())
			assert.Equal(t, tt.expectedStatus, tt.errCode.HTTPStatus())
		})
	}
}

// TestEnhancedError 测试EnhancedError功能
func TestEnhancedError(t *testing.T) {
	t.Run("测试基础创建", func(t *testing.T) {
		err := NewEnhancedError(2001001, "Tenant not found", "租户不存在")
		assert.Equal(t, int32(2001001), err.Code)
		assert.Equal(t, "Tenant not found", err.Message)
		assert.Equal(t, "租户不存在", err.MessageZH)
	})

	t.Run("测试WithRequestID", func(t *testing.T) {
		err := NewEnhancedError(2001001, "Tenant not found", "租户不存在")
		err.WithRequestID("req-12345")
		assert.Equal(t, "req-12345", err.RequestID)
	})

	t.Run("测试WithDetails", func(t *testing.T) {
		err := NewEnhancedError(2001001, "Tenant not found", "租户不存在")
		details := map[string]interface{}{
			"field":   "tenant_id",
			"user_id": "12345",
		}
		err.WithDetails(details)
		assert.Equal(t, "tenant_id", err.Details["field"])
		assert.Equal(t, "12345", err.Details["user_id"])
	})

	t.Run("测试WithTenantID", func(t *testing.T) {
		err := NewEnhancedError(2001001, "Tenant not found", "租户不存在")
		err.WithTenantID("tenant-12345")
		assert.Equal(t, "tenant-12345", err.TenantID)
	})

	t.Run("测试Error()方法", func(t *testing.T) {
		err := NewEnhancedError(2001001, "Tenant not found", "租户不存在")
		errStr := err.Error()
		assert.Contains(t, errStr, "2001001")
		assert.Contains(t, errStr, "Tenant not found")
		assert.Contains(t, errStr, "租户不存在")
	})
}

// TestHTTPStatusMapping 测试HTTP状态码映射逻辑
func TestHTTPStatusMapping(t *testing.T) {
	tests := []struct {
		name       string
		errCode    int32
		expectCode int
	}{
		// 租户错误 - 验证精确映射优先
		{errCode: 2001001, expectCode: http.StatusNotFound},       // 租户不存在，精确映射为404
		{errCode: 2001002, expectCode: http.StatusConflict},        // 租户已存在，精确映射为409
		{errCode: 2002001, expectCode: http.StatusBadRequest},      // 参数无效，精确映射为400
		{errCode: 2003001, expectCode: http.StatusForbidden},       // 租户已暂停，精确映射为403
		{errCode: 2005001, expectCode: http.StatusNotFound},       // 2005001 % 1000 = 1 → 404 (无精确映射，走范围映射)

		// 范围映射测试
		{errCode: 404001, expectCode: http.StatusNotFound}, // 404001 % 1000 = 1 → 404 (001-009范围)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := GetHTTPStatusForError(tt.errCode)
			assert.Equal(t, tt.expectCode, status)
		})
	}
}

// TestCommonErrorCodes 测试通用错误码
func TestCommonErrorCodes(t *testing.T) {
	tests := []struct {
		name           string
		errCode        *BaseErrorCode
		expectedCode   string
		expectedZH     string
		expectedStatus int
	}{
		{
			name:           "Success",
			errCode:        Success,
			expectedCode:   "SUCCESS",
			expectedZH:     "操作成功",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "InvalidParams",
			errCode:        InvalidParams,
			expectedCode:   "COMMON201001",
			expectedZH:     "参数错误",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unauthorized",
			errCode:        Unauthorized,
			expectedCode:   "COMMON401001",
			expectedZH:     "未认证",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Forbidden",
			errCode:        Forbidden,
			expectedCode:   "COMMON403001",
			expectedZH:     "权限不足",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "NotFound",
			errCode:        NotFound,
			expectedCode:   "COMMON404001",
			expectedZH:     "资源不存在",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "InternalError",
			errCode:        InternalError,
			expectedCode:   "COMMON500001",
			expectedZH:     "服务器内部错误",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedCode, tt.errCode.Code())
			assert.Equal(t, tt.expectedZH, tt.errCode.MessageZH())
			assert.Equal(t, tt.expectedStatus, tt.errCode.HTTPStatus())
		})
	}
}

// TestQuotaErrorCodes 测试配额错误码便捷变量
func TestQuotaErrorCodes(t *testing.T) {
	tests := []struct {
		name           string
		errCode        *BaseErrorCode
		expectedCode   string
		expectedStatus int
	}{
		{
			name:           "ErrQuotaNotFound",
			errCode:        ErrQuotaNotFound,
			expectedCode:   "QUOTA_NOT_FOUND",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "ErrQuotaInvalid",
			errCode:        ErrQuotaInvalid,
			expectedCode:   "QUOTA_INVALID",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "ErrQuotaExceeded",
			errCode:        ErrQuotaExceeded,
			expectedCode:   "QUOTA_EXCEEDED",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "ErrQuotaAPICallExceeded",
			errCode:        ErrQuotaAPICallExceeded,
			expectedCode:   "QUOTA_API_CALL_EXCEEDED",
			expectedStatus: http.StatusTooManyRequests,
		},
		{
			name:           "ErrQuotaBotExceeded",
			errCode:        ErrQuotaBotExceeded,
			expectedCode:   "QUOTA_BOT_EXCEEDED",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "ErrQuotaStorageExceeded",
			errCode:        ErrQuotaStorageExceeded,
			expectedCode:   "QUOTA_STORAGE_EXCEEDED",
			expectedStatus: 507, // HTTP StatusInsufficientStorage
		},
		{
			name:           "ErrQuotaInvalidResource",
			errCode:        ErrQuotaInvalidResource,
			expectedCode:   "QUOTA_INVALID_RESOURCE",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "ErrQuotaInvalidLimit",
			errCode:        ErrQuotaInvalidLimit,
			expectedCode:   "QUOTA_INVALID_LIMIT",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "ErrQuotaCheckFailed",
			errCode:        ErrQuotaCheckFailed,
			expectedCode:   "QUOTA_CHECK_FAILED",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "ErrQuotaConsumeFailed",
			errCode:        ErrQuotaConsumeFailed,
			expectedCode:   "QUOTA_CONSUME_FAILED",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "ErrQuotaResetFailed",
			errCode:        ErrQuotaResetFailed,
			expectedCode:   "QUOTA_RESET_FAILED",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "ErrQuotaKnowledgeExceeded",
			errCode:        ErrQuotaKnowledgeExceeded,
			expectedCode:   "QUOTA_KNOWLEDGE_EXCEEDED",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "ErrQuotaWorkflowExceeded",
			errCode:        ErrQuotaWorkflowExceeded,
			expectedCode:   "QUOTA_WORKFLOW_EXCEEDED",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedCode, tt.errCode.Code())
			assert.NotEmpty(t, tt.errCode.MessageZH())
			assert.Equal(t, tt.expectedStatus, tt.errCode.HTTPStatus())
		})
	}
}

// TestPermissionErrorCodes 测试权限错误码便捷变量
func TestPermissionErrorCodes(t *testing.T) {
	tests := []struct {
		name           string
		errCode        *BaseErrorCode
		expectedCode   string
		expectedStatus int
	}{
		{
			name:           "ErrRoleNotFound",
			errCode:        ErrRoleNotFound,
			expectedCode:   "ROLE_NOT_FOUND",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "ErrRoleAlreadyExists",
			errCode:        ErrRoleAlreadyExists,
			expectedCode:   "ROLE_ALREADY_EXISTS",
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "ErrInvalidRole",
			errCode:        ErrInvalidRole,
			expectedCode:   "INVALID_ROLE",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "ErrRoleNameInvalid",
			errCode:        ErrRoleNameInvalid,
			expectedCode:   "ROLE_NAME_INVALID",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "ErrPermissionInvalidParam",
			errCode:        ErrPermissionInvalidParam,
			expectedCode:   "PERMISSION_INVALID_PARAM",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "ErrPermissionCheckFailed",
			errCode:        ErrPermissionCheckFailed,
			expectedCode:   "PERMISSION_CHECK_FAILED",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "ErrResourceTypeNotSupported",
			errCode:        ErrResourceTypeNotSupported,
			expectedCode:   "RESOURCE_TYPE_NOT_SUPPORTED",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedCode, tt.errCode.Code())
			assert.NotEmpty(t, tt.errCode.MessageZH())
			assert.Equal(t, tt.expectedStatus, tt.errCode.HTTPStatus())
		})
	}
}

// TestBillingErrorCodes 测试计费错误码便捷变量
func TestBillingErrorCodes(t *testing.T) {
	tests := []struct {
		name           string
		errCode        *BaseErrorCode
		expectedCode   string
		expectedStatus int
	}{
		{
			name:           "ErrBillingChargeFailed",
			errCode:        ErrBillingChargeFailed,
			expectedCode:   "BILLING_CHARGE_FAILED",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "ErrBillingBalanceNotFound",
			errCode:        ErrBillingBalanceNotFound,
			expectedCode:   "BILLING_BALANCE_NOT_FOUND",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "ErrBillingInsufficientBalance",
			errCode:        ErrBillingInsufficientBalance,
			expectedCode:   "BILLING_INSUFFICIENT_BALANCE",
			expectedStatus: 402, // HTTP StatusPaymentRequired
		},
		{
			name:           "ErrBillingRecordFailed",
			errCode:        ErrBillingRecordFailed,
			expectedCode:   "BILLING_RECORD_FAILED",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "ErrBillingCalculateFailed",
			errCode:        ErrBillingCalculateFailed,
			expectedCode:   "BILLING_CALCULATE_FAILED",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedCode, tt.errCode.Code())
			assert.NotEmpty(t, tt.errCode.MessageZH())
			assert.Equal(t, tt.expectedStatus, tt.errCode.HTTPStatus())
		})
	}
}

// TestBotStoreErrorCodes 测试Bot商店错误码便捷变量
func TestBotStoreErrorCodes(t *testing.T) {
	// BotStore errno变量是error类型，只测试它们不为nil
	t.Run("验证BotStore errno变量存在", func(t *testing.T) {
		assert.NotNil(t, ErrBotStoreInvalidParam)
		assert.NotNil(t, ErrBotStoreNotFound)
		assert.NotNil(t, ErrCreateBotStoreItemFailed)
		assert.NotNil(t, ErrUpdateBotStoreItemFailed)
		assert.NotNil(t, ErrDeleteBotStoreItemFailed)
		assert.NotNil(t, ErrBotStoreItemNotFound)
	})
}

// TestErrnoHelperFunctions 测试errno辅助函数
func TestErrnoHelperFunctions(t *testing.T) {
	t.Run("IsNotFoundError - 测试NotFound类型错误", func(t *testing.T) {
		tests := []struct {
			name     string
			errCode  int32
			expected bool
		}{
			{"ErrTenantNotFound (001结尾)", 2001001, true},
			{"ErrBotStoreNotFound (001结尾)", 206000001, true},
			{"ErrTenantAlreadyExists (010结尾)", 2001010, false},
			{"ErrTenantInvalidParam (020结尾)", 2001020, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := IsNotFoundError(tt.errCode)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("IsAlreadyExistsError - 测试AlreadyExists类型错误", func(t *testing.T) {
		tests := []struct {
			name     string
			errCode  int32
			expected bool
		}{
			{"ErrTenantAlreadyExists (010结尾)", 2001010, true},
			{"ErrTenantNotFound (001结尾)", 2001001, false},
			{"ErrTenantInvalidParam (020结尾)", 2001020, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := IsAlreadyExistsError(tt.errCode)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("IsInvalidParamError - 测试InvalidParam类型错误", func(t *testing.T) {
		tests := []struct {
			name     string
			errCode  int32
			expected bool
		}{
			{"ErrTenantInvalidName (020结尾)", 2001020, true},
			{"ErrTenantInvalidConfig (021结尾)", 2001021, true},
			{"ErrTenantNotFound (001结尾)", 2001001, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := IsInvalidParamError(tt.errCode)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("IsPermissionDeniedError - 测试PermissionDenied类型错误", func(t *testing.T) {
		tests := []struct {
			name     string
			errCode  int32
			expected bool
		}{
			{"ErrTenantInactive (040结尾)", 2001040, true},
			{"ErrQuotaExceeded (041结尾)", 300000041, true},
			{"ErrTenantNotFound (001结尾)", 2001001, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := IsPermissionDeniedError(tt.errCode)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("IsOperationFailedError - 测试OperationFailed类型错误", func(t *testing.T) {
		tests := []struct {
			name     string
			errCode  int32
			expected bool
		}{
			{"ErrTenantCreateFailed (050结尾)", 2001050, true},
			{"ErrTenantUpdateFailed (051结尾)", 2001051, true},
			{"ErrTenantNotFound (001结尾)", 2001001, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := IsOperationFailedError(tt.errCode)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("GetErrorCategory - 测试错误分类", func(t *testing.T) {
		tests := []struct {
			name           string
			errCode        int32
			expectedCategory ErrorCategory
		}{
			{"NotFound", 2001001, ErrorCategoryNotFound},
			{"AlreadyExists", 2001010, ErrorCategoryAlreadyExists},
			{"InvalidParam", 2001020, ErrorCategoryInvalidParam},
			{"PermissionDenied", 2001040, ErrorCategoryPermission},
			{"OperationFailed", 2001050, ErrorCategoryOperationFailed},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				category := GetErrorCategory(tt.errCode)
				assert.Equal(t, tt.expectedCategory, category)
			})
		}
	})
}

// TestErrnoCoverage 测试errno便捷变量覆盖率
func TestErrnoCoverage(t *testing.T) {
	t.Run("验证所有便捷变量都正确定义", func(t *testing.T) {
		// 测试quota模块便捷变量不为nil
		assert.NotNil(t, ErrQuotaNotFound)
		assert.NotNil(t, ErrQuotaInvalid)
		assert.NotNil(t, ErrQuotaExceeded)
		assert.NotNil(t, ErrQuotaAPICallExceeded)
		assert.NotNil(t, ErrQuotaBotExceeded)
		assert.NotNil(t, ErrQuotaStorageExceeded)
		assert.NotNil(t, ErrQuotaKnowledgeExceeded)
		assert.NotNil(t, ErrQuotaWorkflowExceeded)

		// 测试permission模块便捷变量不为nil
		assert.NotNil(t, ErrRoleNotFound)
		assert.NotNil(t, ErrRoleAlreadyExists)
		assert.NotNil(t, ErrInvalidRole)
		assert.NotNil(t, ErrRoleNameInvalid)
		assert.NotNil(t, ErrPermissionInvalidParam)
		assert.NotNil(t, ErrPermissionCheckFailed)
		assert.NotNil(t, ErrResourceTypeNotSupported)

		// 测试billing模块便捷变量不为nil
		assert.NotNil(t, ErrBillingChargeFailed)
		assert.NotNil(t, ErrBillingBalanceNotFound)
		assert.NotNil(t, ErrBillingInsufficientBalance)
		assert.NotNil(t, ErrBillingRecordFailed)
		assert.NotNil(t, ErrBillingCalculateFailed)

		// 测试botstore模块便捷变量不为nil
		assert.NotNil(t, ErrBotStoreInvalidParam)
		assert.NotNil(t, ErrBotStoreNotFound)
		assert.NotNil(t, ErrCreateBotStoreItemFailed)
		assert.NotNil(t, ErrUpdateBotStoreItemFailed)
		assert.NotNil(t, ErrDeleteBotStoreItemFailed)
		assert.NotNil(t, ErrBotStoreItemNotFound)
	})
}
