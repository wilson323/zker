// backend/api/middleware/quota_config.go
package middleware

import "time"

const (
	// 性能配置
	quotaCheckTimeout         = 50 * time.Millisecond
	quotaCacheCleanupInterval = 5 * time.Minute

	// 缓存配置
	maxPathCacheSize = 10000
	pathCacheTTL      = 10 * time.Minute
)

// QuotaCheckConfig 配额检查配置
type QuotaCheckConfig struct {
	// 是否启用缓存
	EnableCache bool
	// 缓存TTL
	CacheTTL time.Duration
	// 最大缓存大小
	MaxCacheSize int
	// 是否启用性能监控
	EnableMetrics bool
	// 慢查询阈值（超过此值记录日志）
	SlowQueryThreshold time.Duration
}

// DefaultConfig 默认配置
var DefaultConfig = QuotaCheckConfig{
	EnableCache:         true,
	CacheTTL:            60 * time.Second,
	MaxCacheSize:        10000,
	EnableMetrics:       true,
	SlowQueryThreshold:  20 * time.Millisecond,
}

// DevelopmentConfig 开发环境配置
// 开发环境关闭缓存，便于调试
var DevelopmentConfig = QuotaCheckConfig{
	EnableCache:         false,
	CacheTTL:            0,
	MaxCacheSize:        0,
	EnableMetrics:       true,
	SlowQueryThreshold:  100 * time.Millisecond,
}

// ProductionConfig 生产环境配置
// 生产环境启用所有优化
var ProductionConfig = QuotaCheckConfig{
	EnableCache:         true,
	CacheTTL:            60 * time.Second,
	MaxCacheSize:        50000,
	EnableMetrics:       true,
	SlowQueryThreshold:  20 * time.Millisecond,
}

// TestingConfig 测试环境配置
var TestingConfig = QuotaCheckConfig{
	EnableCache:         false,
	CacheTTL:            0,
	MaxCacheSize:        0,
	EnableMetrics:       false,
	SlowQueryThreshold:  0,
}

// GetConfig 根据环境获取配置
func GetConfig(env string) QuotaCheckConfig {
	switch env {
	case "production", "prod":
		return ProductionConfig
	case "development", "dev":
		return DevelopmentConfig
	case "testing", "test":
		return TestingConfig
	default:
		return DefaultConfig
	}
}

// Validate 验证配置有效性
func (c *QuotaCheckConfig) Validate() error {
	if c.EnableCache {
		if c.CacheTTL <= 0 {
			c.CacheTTL = 60 * time.Second
		}
		if c.MaxCacheSize <= 0 {
			c.MaxCacheSize = 10000
		}
	}

	if c.SlowQueryThreshold <= 0 {
		c.SlowQueryThreshold = 20 * time.Millisecond
	}

	return nil
}
