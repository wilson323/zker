// backend/api/middleware/quota_resource_type_test.go
package middleware

import (
	"sync"
	"testing"

	"github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/stretchr/testify/assert"
)

// TestGetResourceTypeFromPath 测试资源类型识别（优化版本）
func TestGetResourceTypeFromPath(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedType   string
		expectedCode   errno.ErrorCode
		description    string
	}{
		{
			name:           "Bot创建路径",
			path:           "/api/v1/bots",
			expectedType:   "bots",
			expectedCode:   errno.QUOTA402001,
			description:    "Bot创建路径应返回bots类型",
		},
		{
			name:           "Bot更新路径",
			path:           "/api/v1/bots/123",
			expectedType:   "bots",
			expectedCode:   errno.QUOTA402001,
			description:    "Bot更新路径应返回bots类型",
		},
		{
			name:           "对话创建路径",
			path:           "/api/v1/conversations",
			expectedType:   "messages",
			expectedCode:   errno.QUOTA402002,
			description:    "对话创建路径应返回messages类型",
		},
		{
			name:           "知识库路径",
			path:           "/api/v1/knowledge",
			expectedType:   "storage",
			expectedCode:   errno.QUOTA402003,
			description:    "知识库路径应返回storage类型",
		},
		{
			name:           "工作流路径",
			path:           "/api/v1/workflows",
			expectedType:   "workflows",
			expectedCode:   errno.QUOTA402004,
			description:    "工作流路径应返回workflows类型",
		},
		{
			name:           "健康检查路径（无配额）",
			path:           "/api/health",
			expectedType:   "",
			expectedCode:   nil,
			description:    "健康检查路径应返回空",
		},
		{
			name:           "登录路径（无配额）",
			path:           "/api/v1/auth/login",
			expectedType:   "",
			expectedCode:   nil,
			description:    "登录路径应返回空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceType, errorCode := GetResourceTypeFromPath(tt.path)
			assert.Equal(t, tt.expectedType, resourceType, tt.description)
			assert.Equal(t, tt.expectedCode, errorCode, tt.description)
		})
	}
}

// TestGetResourceTypeFromPath_Caching 测试路径缓存功能
func TestGetResourceTypeFromPath_Caching(t *testing.T) {
	path := "/api/v1/bots"

	// 第一次调用：缓存未命中
	resourceType1, errorCode1 := GetResourceTypeFromPath(path)
	assert.Equal(t, "bots", resourceType1)
	assert.Equal(t, errno.QUOTA402001, errorCode1)

	// 第二次调用：缓存命中
	resourceType2, errorCode2 := GetResourceTypeFromPath(path)
	assert.Equal(t, "bots", resourceType2)
	assert.Equal(t, errno.QUOTA402001, errorCode2)

	// 清空缓存
	ClearPathCache()

	// 第三次调用：缓存已清空，重新计算
	resourceType3, errorCode3 := GetResourceTypeFromPath(path)
	assert.Equal(t, "bots", resourceType3)
	assert.Equal(t, errno.QUOTA402001, errorCode3)
}

// TestGetResourceTypeFromPath_Concurrency 测试并发安全性
func TestGetResourceTypeFromPath_Concurrency(t *testing.T) {
	const goroutines = 100
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				GetResourceTypeFromPath("/api/v1/bots")
				GetResourceTypeFromPath("/api/v1/conversations")
				GetResourceTypeFromPath("/api/health")
			}
		}()
	}

	wg.Wait()

	// 如果没有 panic 或死锁，测试通过
	assert.True(t, true)
}

// TestRegisterResourceType 测试动态注册资源类型
func TestRegisterResourceType(t *testing.T) {
	// 注册新的资源类型
	newConfig := ResourceTypeConfig{
		PathPrefix:   "/api/v1/plugins",
		ResourceType: "plugins",
		ErrorCode:    errno.QUOTA402001, // 复用现有错误码
	}

	RegisterResourceType(newConfig)

	// 验证新注册的类型可识别
	resourceType, errorCode := GetResourceTypeFromPath("/api/v1/plugins")
	assert.Equal(t, "plugins", resourceType)
	assert.Equal(t, errno.QUOTA402001, errorCode)

	// 清理：清空缓存，恢复原始状态
	ClearPathCache()
}

// BenchmarkGetResourceTypeFromPath 性能基准测试
func BenchmarkGetResourceTypeFromPath(b *testing.B) {
	path := "/api/v1/bots"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetResourceTypeFromPath(path)
	}
}

// BenchmarkGetResourceTypeFromPath_Parallel 并发性能基准测试
func BenchmarkGetResourceTypeFromPath_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		path := "/api/v1/bots"
		for pb.Next() {
			GetResourceTypeFromPath(path)
		}
	})
}
