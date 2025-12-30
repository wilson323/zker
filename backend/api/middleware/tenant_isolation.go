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

package middleware

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/service"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	pkgerrorx "github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/consts"
	berrno "github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/infra/tracing"
	"github.com/coze-dev/coze-studio/backend/infra/cache"
)

const (
	// TenantIDKey context中存储tenant_id的key
	TenantIDKey = "tenant_id"

	// TenantIDHeader HTTP Header中传递tenant_id的key（用于API调用）
	TenantIDHeader = "X-Tenant-ID"

	// DefaultTenantID 默认租户ID（兼容模式：当无法获取tenant_id时使用）
	// TODO: 数据迁移完成后移除兼容模式，强制要求租户ID
	DefaultTenantID = "default"

	// EnableStrictTenantIsolation 启用严格租户隔离模式
	// true: 强制要求tenant_id，不存在则拒绝请求
	// false: 兼容模式，使用默认租户ID
	// TODO: 数据迁移完成后设置为 true
	EnableStrictTenantIsolation = false

	// CacheKeyPrefix 租户缓存key前缀
	CacheKeyPrefix = "tenant:info:"

	// CacheTTL 租户缓存过期时间
	CacheTTL = 5 * time.Minute
)

var (
	// 全局依赖（在应用启动时注入）
	tenantService *service.TenantService
	redisCache    cache.Cmdable

	// 不需要租户隔离的路径
	noTenantIsolationPath = map[string]bool{
		"/api/passport/": true,
		"/api/health/":   true,
		"/api/metrics/":  true,
		"/api/public/":   true,
	}
)

// InitTenantMiddleware 初始化租户中间件依赖
//
// **使用场景**：应用启动时调用，注入tenantService和redis cache
//
// **示例**：
//   middleware.InitTenantMiddleware(tenantService, redisClient)
func InitTenantMiddleware(ts *service.TenantService, rc cache.Cmdable) {
	tenantService = ts
	redisCache = rc
}

// TenantIsolationMiddleware 租户隔离中间件
//
// **核心功能**：
// 1. 从请求中提取tenant_id（优先级：HTTP Header > Session > JWT > Default）
// 2. 验证租户存在且激活（防止跨租户访问）✅ 已启用
// 3. 将tenant_id注入context（供后续业务逻辑使用）
// 4. 记录追踪和日志（便于故障排查）
//
// **TODO - 数据迁移后需要完成**：
// [ ] 1. Session添加TenantID字段（修改session.go）
// [ ] 2. users表添加tenant_id字段（数据库迁移）
// [ ] 3. 创建user_tenant关联表（支持用户多租户）
// [ ] 4. JWT token包含tenant_id（token生成逻辑）
// [ ] 5. 移除DefaultTenantID兼容逻辑（强制租户隔离）
func TenantIsolationMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 1. 检查是否需要跳过租户隔离
		if shouldSkipTenantIsolation(ctx) {
			ctx.Next(c)
			return
		}

		// 2. 提取tenant_id（多来源兼容）
		tenantID, err := extractTenantID(c, ctx)
		if err != nil {
			logs.CtxErrorf(c, "[TenantIsolation] extract tenant_id failed: %v", err)
			httputil.Unauthorized(ctx, "missing or invalid tenant_id")
			return
		}

		// 3. 验证租户（✅ 已启用）
		if err := validateTenant(c, tenantID); err != nil {
			logs.CtxErrorf(c, "[TenantIsolation] validate tenant failed: %v", err)
			httputil.Error(ctx, err.Error(), consts.StatusForbidden)
			return
		}

		// 4. 将tenant_id注入context
		ctxcache.Store(c, consts.TenantIDKeyInCtx, tenantID)
		ctx.Set(TenantIDKey, tenantID)

		// 5. 记录追踪和日志
		tracing.AddSpanAttributes(c,
			tracing.AttrTenantID.String(tenantID),
		)
		logs.CtxInfof(c, "[TenantIsolation] tenant_id=%s", tenantID)

		ctx.Next(c)
	}
}

// extractTenantID 从请求中提取tenant_id
//
// **优先级**：
// 1. HTTP Header: X-Tenant-ID（用于API调用）
// 2. Session: session.TenantID（✅ 已启用）
// 3. JWT Token: claims.tenant_id（TODO: JWT改造后启用）
// 4. 严格模式：返回错误（EnableStrictTenantIsolation = true）
// 5. 兼容模式：使用默认租户ID（EnableStrictTenantIsolation = false，数据迁移后移除）
func extractTenantID(c context.Context, ctx *app.RequestContext) (string, error) {
	// 1. 从HTTP Header获取（优先级最高）
	if tenantID := string(ctx.GetHeader(TenantIDHeader)); tenantID != "" {
		return tenantID, nil
	}

	// 2. 从Session获取（✅ 已启用）
	if session, ok := ctxcache.Get[*entity.Session](c, consts.SessionDataKeyInCtx); ok {
		if session.HasTenantID() {
			return session.TenantID, nil
		}
		// 兼容旧Session（没有tenant_id）
		// 继续尝试其他来源
	}

	// 3. 从JWT Token获取（TODO: JWT改造后启用）
	// if claims := getJWTClaims(ctx); claims != nil && claims.TenantID != "" {
	// 	return claims.TenantID, nil
	// }

	// 4. 严格模式检查
	if EnableStrictTenantIsolation {
		logs.CtxErrorf(c, "[TenantIsolation] strict mode enabled but no tenant_id found")
		return "", pkgerrorx.New(berrno.ErrMissingTenantID)
	}

	// 5. 兼容模式：使用默认租户ID（数据迁移后移除）
	logs.CtxWarnf(c, "[TenantIsolation] no tenant_id found, using default: %s", DefaultTenantID)
	return DefaultTenantID, nil
}

// validateTenant 验证租户存在且激活
//
// **验证逻辑**：
// 1. 检查租户是否存在
// 2. 检查租户状态是否为active
// 3. 检查租户是否被删除
//
// **性能优化**：
// - 使用Redis缓存租户信息（TTL=5分钟）
// - 缓存key: tenant:info:{tenant_id}
func validateTenant(c context.Context, tenantID string) error {
	// 跳过默认租户的验证（兼容模式）
	if tenantID == DefaultTenantID {
		return nil
	}

	// 检查是否已初始化依赖
	if tenantService == nil {
		logs.CtxWarnf(c, "[ValidateTenant] tenantService not initialized, skipping validation")
		return nil
	}

	// 1. 尝试从Redis缓存获取
	if redisCache != nil {
		cachedTenant, err := getTenantFromCache(c, tenantID)
		if err == nil && cachedTenant != nil {
			if !cachedTenant.IsActive() {
				return pkgerrorx.New(berrno.ErrTenantSuspended).WithZap(
					zap.String("tenant_id", tenantID),
					zap.String("status", string(cachedTenant.Status)),
				)
			}
			return nil
		}
		// 缓存未命中，继续查询数据库
	}

	// 2. 从数据库查询
	tenant, err := tenantService.GetTenant(c, tenantID)
	if err != nil {
		// 转换为统一错误码
		return pkgerrorx.Wrap(err, berrno.ErrTenantNotFound).WithZap(
			zap.String("tenant_id", tenantID),
		)
	}
	if tenant == nil {
		return pkgerrorx.New(berrno.ErrTenantNotFound).WithZap(
			zap.String("tenant_id", tenantID),
		)
	}

	// 3. 验证状态
	if !tenant.IsActive() {
		return pkgerrorx.New(berrno.ErrTenantSuspended).WithZap(
			zap.String("tenant_id", tenantID),
			zap.String("status", string(tenant.Status)),
		)
	}

	// 4. 检查是否已删除
	if tenant.IsDeleted() {
		return pkgerrorx.New(berrno.ErrTenantDeleted).WithZap(
			zap.String("tenant_id", tenantID),
		)
	}

	// 5. 写入缓存
	if redisCache != nil {
		if err := setTenantToCache(c, tenant, CacheTTL); err != nil {
			logs.CtxWarnf(c, "[ValidateTenant] failed to cache tenant: %v", err)
		}
	}

	return nil
}

// getTenantFromCache 从Redis缓存获取租户信息
func getTenantFromCache(c context.Context, tenantID string) (*entity.Tenant, error) {
	cacheKey := CacheKeyPrefix + tenantID

	data, err := redisCache.Get(c, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	var tenant entity.Tenant
	if err := json.Unmarshal([]byte(data), &tenant); err != nil {
		logs.CtxErrorf(c, "[GetTenantFromCache] failed to unmarshal tenant: %v", err)
		return nil, err
	}

	return &tenant, nil
}

// setTenantToCache 将租户信息写入Redis缓存
func setTenantToCache(c context.Context, tenant *entity.Tenant, ttl time.Duration) error {
	cacheKey := CacheKeyPrefix + tenant.TenantID

	data, err := json.Marshal(tenant)
	if err != nil {
		return err
	}

	return redisCache.Set(c, cacheKey, data, ttl).Err()
}

// shouldSkipTenantIsolation 检查是否需要跳过租户隔离
func shouldSkipTenantIsolation(ctx *app.RequestContext) bool {
	path := string(ctx.GetRequest().URI().Path)

	// 检查精确匹配
	for skipPath := range noTenantIsolationPath {
		if len(path) >= len(skipPath) && path[:len(skipPath)] == skipPath {
			return true
		}
	}

	return false
}

// GetTenantIDFromContext 从context中获取tenant_id
//
// **使用场景**：
// - 业务逻辑需要获取当前租户ID
// - 数据库查询需要添加WHERE tenant_id = ?
// - 日志记录需要标记租户
func GetTenantIDFromContext(c context.Context) string {
	// 优先从Hertz context获取
	if tenantID, ok := c.Value(TenantIDKey).(string); ok {
		return tenantID
	}

	// 其次从ctxcache获取
	if tenantID, ok := ctxcache.Get[string](c, consts.TenantIDKeyInCtx); ok {
		return tenantID
	}

	// 兜底：返回默认租户ID
	logs.CtxWarnf(c, "[GetTenantIDFromContext] no tenant_id in context, using default")
	return DefaultTenantID
}

// MustGetTenantIDFromContext 从context中获取tenant_id，如果没有则panic
//
// **使用场景**：
// - 租户强相关的业务逻辑
// - 必须明确租户的操作
func MustGetTenantIDFromContext(c context.Context) string {
	tenantID := GetTenantIDFromContext(c)
	if tenantID == "" || tenantID == DefaultTenantID {
		logs.CtxErrorf(c, "[MustGetTenantIDFromContext] invalid tenant_id: %s", tenantID)
		panic("tenant_id is required")
	}
	return tenantID
}

// RequireTenantID 强制要求tenant_id的中间件
//
// **使用场景**：
// - 租户管理API
// - 数据隔离要求严格的接口
func RequireTenantID() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		tenantID := GetTenantIDFromContext(c)
		if tenantID == "" || tenantID == DefaultTenantID {
			logs.CtxWarnf(c, "[RequireTenantID] missing valid tenant_id")
			httputil.Unauthorized(ctx, "valid tenant_id is required")
			return
		}
		ctx.Next(c)
	}
}

// TenantIDToInt64 将tenant_id string转换为int64
//
// **注意**：这是一个临时兼容函数
// 当前系统很多地方使用int64类型的user_id
// 数据迁移后，tenant_id统一为string类型（UUID）
func TenantIDToInt64(tenantID string) (int64, error) {
	// 如果是UUID格式，返回错误
	if len(tenantID) == 36 {
		return 0, pkgerrorx.New(berrno.ErrInvalidTenantID,
			zap.String("tenant_id", tenantID),
			zap.String("reason", "UUID format cannot convert to int64"))
	}

	// 尝试转换为int64
	id, err := strconv.ParseInt(tenantID, 10, 64)
	if err != nil {
		return 0, pkgerrorx.Wrap(err, berrno.ErrInvalidTenantID,
			zap.String("tenant_id", tenantID))
	}
	return id, nil
}
