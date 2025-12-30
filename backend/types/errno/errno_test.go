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

// TestCacheErrorCodes 测试缓存错误码
func TestCacheErrorCodes(t *testing.T) {
	tests := []struct {
		name           string
		errCode        *BaseErrorCode
		expectedCode   string
		expectedZH     string
		expectedEN     string
		expectedStatus int
	}{
		{
			name:           "CACHE500001 - 缓存连接失败",
			errCode:        CACHE500001,
			expectedCode:   "CACHE500001",
			expectedZH:     "缓存连接失败",
			expectedEN:     "Cache connection failed",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "CACHE500002 - 缓存操作失败",
			errCode:        CACHE500002,
			expectedCode:   "CACHE500002",
			expectedZH:     "缓存操作失败",
			expectedEN:     "Cache operation failed",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "CACHE500003 - 缓存雪崩",
			errCode:        CACHE500003,
			expectedCode:   "CACHE500003",
			expectedZH:     "缓存雪崩",
			expectedEN:     "Cache avalanche",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "CACHE500004 - 缓存穿透",
			errCode:        CACHE500004,
			expectedCode:   "CACHE500004",
			expectedZH:     "缓存穿透",
			expectedEN:     "Cache penetration",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "CACHE500005 - 缓存击穿",
			errCode:        CACHE500005,
			expectedCode:   "CACHE500005",
			expectedZH:     "缓存击穿",
			expectedEN:     "Cache breakdown",
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
		// 租户错误
		{errCode: 2001001, expectCode: http.StatusNotFound},       // 租户不存在
		{errCode: 2001002, expectCode: http.StatusConflict},        // 租户已存在
		{errCode: 2002001, expectCode: http.StatusBadRequest},      // 参数无效
		{errCode: 2003001, expectCode: http.StatusForbidden},       // 租户已暂停
		{errCode: 2005001, expectCode: http.StatusInternalServerError}, // 迁移失败

		// 注意: 数据库(DB500001)和缓存(CACHE500001)错误使用字符串常量
		// 它们是 *BaseErrorCode 类型，不通过 GetHTTPStatusForError(int32) 处理
		// 配额错误使用 300xxxx 范围 (如 ErrQuotaExceededCode = 300000003)
		// 因此不存在数字 500001 的冲突
		// 以下测试仅演示 int32 类型错误码的范围映射
		{errCode: 404001, expectCode: http.StatusNotFound}, // 缓存键不存在
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
