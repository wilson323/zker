// backend/api/middleware/tenant_isolation_middleware.go
// 租户隔离中间件 - 根据租户隔离策略路由请求到不同的数据库实例
package middleware

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	tenantservice "github.com/coze-dev/coze-studio/backend/domain/tenant/service"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// ================================================================================
// 租户隔离策略中间件
// ================================================================================

// TenantIsolationMiddlewareConfig 租户隔离中间件配置
type TenantIsolationMiddlewareConfig struct {
	IsolationService *tenantservice.IsolationUpgradeService
	DBManager        DatabaseManager
	CacheEnabled     bool
}

// DatabaseManager 数据库管理器接口
type DatabaseManager interface {
	// GetDatabase 根据租户隔离策略获取数据库连接
	GetDatabase(ctx context.Context, tenantID string, strategy tenantservice.IsolationStrategy) (*DatabaseConnection, error)

	// ReleaseDatabase 释放数据库连接
	ReleaseDatabase(ctx context.Context, conn *DatabaseConnection) error
}

// DatabaseConnection 数据库连接
type DatabaseConnection struct {
	DB           interface{} // GORM DB实例或原生DB
	InstanceID   int
	SchemaName   string
	IsolationStrategy tenantservice.IsolationStrategy
}

// TenantIsolationDatabaseMiddleware 租户隔离数据库路由中间件
//
// 功能:
//   1. 检查租户隔离策略
//   2. 根据策略选择数据库实例
//   3. 将数据库连接注入上下文
//   4. 监控隔离策略使用情况
func TenantIsolationDatabaseMiddleware(config TenantIsolationMiddlewareConfig) app.HandlerFunc {
	// 租户策略缓存
	strategyCache := make(map[string]tenantservice.IsolationStrategy)
	var cacheMutex sync.RWMutex

	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 从上下文获取租户信息
		tenant, err := GetTenantFromContext(c)
		if err != nil {
			logs.CtxWarnf(ctx, "[TenantIsolation] 获取租户信息失败: %v", err)
			c.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    "TENANT_ISOLATION_001",
				"message": "获取租户信息失败",
			})
			c.Abort()
			return
		}

		// 2. 获取租户隔离策略
		strategy, err := getIsolationStrategy(ctx, tenant.TenantID, config.IsolationService, strategyCache, &cacheMutex)
		if err != nil {
			logs.CtxErrorf(ctx, "[TenantIsolation] 获取隔离策略失败 [tenant=%s]: %v", tenant.TenantID, err)
			c.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    "TENANT_ISOLATION_002",
				"message": "获取隔离策略失败",
			})
			c.Abort()
			return
		}

		logs.CtxInfof(ctx, "[TenantIsolation] 租户 %s 使用隔离策略: %s", tenant.TenantID, strategy)

		// 3. 根据策略选择数据库实例
		if strategy == tenantservice.StrategyDatabaseLevel && config.DBManager != nil {
			// Database级隔离: 需要切换数据库实例
			dbConn, err := config.DBManager.GetDatabase(ctx, tenant.TenantID, strategy)
			if err != nil {
				logs.CtxErrorf(ctx, "[TenantIsolation] 获取数据库连接失败 [tenant=%s]: %v", tenant.TenantID, err)
				c.JSON(consts.StatusInternalServerError, map[string]interface{}{
					"code":    "TENANT_ISOLATION_003",
					"message": "获取数据库连接失败",
				})
				c.Abort()
				return
			}

			// 将数据库连接注入上下文
			c.Set("db_connection", dbConn)
			defer func() {
				// 释放连接
				if err := config.DBManager.ReleaseDatabase(ctx, dbConn); err != nil {
					logs.CtxWarnf(ctx, "[TenantIsolation] 释放数据库连接失败: %v", err)
				}
			}()

			logs.CtxInfof(ctx, "[TenantIsolation] 已切换到数据库实例 %d", dbConn.InstanceID)
		}

		// 4. 将隔离策略注入上下文
		c.Set("isolation_strategy", string(strategy))

		c.Next(ctx)
	}
}

// getIsolationStrategy 获取租户隔离策略（带缓存）
func getIsolationStrategy(
	ctx context.Context,
	tenantID string,
	isolationService *tenantservice.IsolationUpgradeService,
	cache map[string]tenantservice.IsolationStrategy,
	cacheMutex *sync.RWMutex,
) (tenantservice.IsolationStrategy, error) {
	// 1. 尝试从缓存获取
	if cache != nil {
		cacheMutex.RLock()
		strategy, exists := cache[tenantID]
		cacheMutex.RUnlock()

		if exists {
			return strategy, nil
		}
	}

	// 2. 从服务获取
	status, err := isolationService.GetUpgradeStatus(ctx, tenantID)
	if err != nil {
		return "", fmt.Errorf("查询隔离策略失败: %w", err)
	}

	// 3. 解析策略
	strategy := tenantservice.IsolationStrategy(status.CurrentStrategy)
	if strategy == "" {
		// 默认使用行级隔离
		strategy = tenantservice.StrategyRowLevel
	}

	// 4. 更新缓存
	if cache != nil {
		cacheMutex.Lock()
		cache[tenantID] = strategy
		cacheMutex.Unlock()
	}

	return strategy, nil
}

// ================================================================================
// 租户配额检查中间件
// ================================================================================

// TenantQuotaMiddlewareConfig 租户配额检查中间件配置
type TenantQuotaMiddlewareConfig struct {
	QuotaService interface {
		CheckQuota(ctx context.Context, tenantID string, resourceType string, count int) error
	}
	ResourceType string // 检查的资源类型
	RequiredQuota int   // 需要的配额数量
}

// TenantQuotaMiddleware 租户配额检查中间件
//
// 功能:
//   1. 检查租户配额是否充足
//   2. 如果不足，返回403错误
//   3. 记录配额使用情况
func TenantQuotaMiddleware(config TenantQuotaMiddlewareConfig) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 获取租户信息
		tenant, err := GetTenantFromContext(c)
		if err != nil {
			logs.CtxWarnf(ctx, "[TenantQuota] 获取租户信息失败: %v", err)
			c.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    "TENANT_QUOTA_001",
				"message": "获取租户信息失败",
			})
			c.Abort()
			return
		}

		// 2. 检查配额
		if config.QuotaService != nil {
			err := config.QuotaService.CheckQuota(ctx, tenant.TenantID, config.ResourceType, config.RequiredQuota)
			if err != nil {
				logs.CtxWarnf(ctx, "[TenantQuota] 配额不足 [tenant=%s, resource=%s, required=%d]: %v",
					tenant.TenantID, config.ResourceType, config.RequiredQuota, err)

				c.JSON(consts.StatusForbidden, map[string]interface{}{
					"code":    "TENANT_QUOTA_002",
					"message": "配额不足",
					"details": map[string]interface{}{
						"resource_type": config.ResourceType,
						"required":      config.RequiredQuota,
					},
				})
				c.Abort()
				return
			}
		}

		logs.CtxInfof(ctx, "[TenantQuota] 配额检查通过 [tenant=%s, resource=%s]", tenant.TenantID, config.ResourceType)

		c.Next(ctx)
	}
}

// ================================================================================
// 租户状态检查中间件
// ================================================================================

// TenantStatusMiddleware 租户状态检查中间件
//
// 功能:
//   1. 检查租户状态是否正常
//   2. 如果被禁用或过期，返回403错误
func TenantStatusMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 获取租户信息
		tenant, err := GetTenantFromContext(c)
		if err != nil {
			logs.CtxWarnf(ctx, "[TenantStatus] 获取租户信息失败: %v", err)
			c.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    "TENANT_STATUS_001",
				"message": "获取租户信息失败",
			})
			c.Abort()
			return
		}

		// 2. 检查租户状态
		if tenant.Status != "active" {
			logs.CtxWarnf(ctx, "[TenantStatus] 租户状态异常 [tenant=%s, status=%s]", tenant.TenantID, tenant.Status)

			c.JSON(consts.StatusForbidden, map[string]interface{}{
				"code":    "TENANT_STATUS_002",
				"message": "租户状态异常",
				"details": map[string]interface{}{
					"tenant_id": tenant.TenantID,
					"status":    tenant.Status,
				},
			})
			c.Abort()
			return
		}

		// 3. 检查订阅是否过期
		// TODO: 实现订阅过期检查逻辑

		c.Next(ctx)
	}
}

// ================================================================================
// 辅助函数
// ================================================================================

// GetDatabaseConnection 从上下文获取数据库连接
func GetDatabaseConnection(c *app.RequestContext) (*DatabaseConnection, error) {
	conn, ok := c.Get("db_connection")
	if !ok || conn == nil {
		// 没有独立的数据库连接，使用默认连接
		return nil, nil
	}

	return conn.(*DatabaseConnection), nil
}

// GetIsolationStrategy 从上下文获取隔离策略
func GetIsolationStrategy(c *app.RequestContext) (tenantservice.IsolationStrategy, error) {
	strategy, ok := c.Get("isolation_strategy")
	if !ok || strategy == nil {
		// 默认行级隔离
		return tenantservice.StrategyRowLevel, nil
	}

	return tenantservice.IsolationStrategy(strategy.(string)), nil
}

// MustGetDatabaseConnection 从上下文获取数据库连接（panic if not found）
func MustGetDatabaseConnection(c *app.RequestContext) *DatabaseConnection {
	conn, err := GetDatabaseConnection(c)
	if err != nil {
		panic(err)
	}
	return conn
}

// MustGetIsolationStrategy 从上下文获取隔离策略（panic if not found）
func MustGetIsolationStrategy(c *app.RequestContext) tenantservice.IsolationStrategy {
	strategy, err := GetIsolationStrategy(c)
	if err != nil {
		panic(err)
	}
	return strategy
}
