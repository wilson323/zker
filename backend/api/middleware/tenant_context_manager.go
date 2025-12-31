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

// Package middleware 提供租户上下文管理中间件
package middleware

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	tenantentity "github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	tenantrepo "github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

const (
	// TenantContextKey 租户上下文 Key
	TenantContextKey = "tenant_context"

	// SubdomainRegex 子域名匹配正则表达式
	// 示例匹配: tenant-a.saas.coze.com, company-b.app.example.com
	SubdomainRegex = `^([a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)\.`

	// PathTenantRegex 路径租户匹配正则表达式
	// 示例匹配: /tenant-a/bots, /company-b/api/v1/...
	PathTenantRegex = `^/([a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)/`
)

var (
	// 编译后的正则表达式
	subdomainPattern = regexp.MustCompile(SubdomainRegex)
	pathPattern      = regexp.MustCompile(PathTenantRegex)

	// 租户缓存（提高性能）
	tenantCache      = make(map[string]*tenantentity.Tenant)
	tenantCacheMutex sync.RWMutex
	tenantCacheTTL   = 300 // 秒
)

// TenantContextManager 租户上下文管理器接口
type TenantContextManager interface {
	// IdentifyFromSubdomain 从子域名识别租户
	// 示例: tenant-a.saas.coze.com -> Tenant{tenant_id: "tenant-a"}
	IdentifyFromSubdomain(ctx context.Context, host string) (*tenantentity.Tenant, error)

	// IdentifyFromPath 从路径识别租户
	// 示例: /tenant-a/bots -> Tenant{tenant_id: "tenant-a"}
	IdentifyFromPath(ctx context.Context, path string) (*tenantentity.Tenant, error)

	// IdentifyFromHeader 从 Header 识别租户（兼容模式）
	// 示例: X-Tenant-ID: tenant-a
	IdentifyFromHeader(ctx context.Context, headers map[string]string) (*tenantentity.Tenant, error)

	// GetTenantContext 从上下文获取租户
	GetTenantContext(ctx context.Context) (*tenantentity.Tenant, error)

	// SetTenantContext 设置租户到上下文
	SetTenantContext(ctx context.Context, tenant *tenantentity.Tenant) context.Context
}

// tenantContextManagerImpl 租户上下文管理器实现
type tenantContextManagerImpl struct {
	tenantRepo   tenantrepo.TenantRepository
	subdomainMap map[string]string // 子域名 -> 租户ID 映射（可选：用于加速）
	baseDomains  []string           // 基础域名列表（如 saas.coze.com）
}

// NewTenantContextManager 创建租户上下文管理器
func NewTenantContextManager(
	tenantRepo tenantrepo.TenantRepository,
	baseDomains []string,
) TenantContextManager {
	return &tenantContextManagerImpl{
		tenantRepo:   tenantRepo,
		subdomainMap: make(map[string]string),
		baseDomains:  baseDomains,
	}
}

// ================================================================================
// IdentifyFromSubdomain 从子域名识别租户
// ================================================================================
func (m *tenantContextManagerImpl) IdentifyFromSubdomain(
	ctx context.Context,
	host string,
) (*tenantentity.Tenant, error) {
	// 1. 提取子域名
	subdomain := m.extractSubdomain(host)
	if subdomain == "" {
		return nil, fmt.Errorf("无效的子域名: %s", host)
	}

	// 2. 从缓存获取
	if tenant := m.getFromCache(subdomain); tenant != nil {
		logs.CtxInfof(ctx, "[TenantContext] 子域名命中缓存: %s -> %s", subdomain, tenant.TenantID)
		return tenant, nil
	}

	// 3. 从数据库查询
	tenant, err := m.tenantRepo.GetBySubdomain(ctx, subdomain)
	if err != nil {
		return nil, fmt.Errorf("查询租户失败 (subdomain=%s): %w", subdomain, err)
	}

	if tenant == nil {
		return nil, fmt.Errorf("租户不存在 (subdomain=%s)", subdomain)
	}

	// 4. 缓存租户信息
	m.putToCache(subdomain, tenant)

	logs.CtxInfof(ctx, "[TenantContext] 子域名识别成功: %s -> %s", subdomain, tenant.TenantID)
	return tenant, nil
}

// ================================================================================
// IdentifyFromPath 从路径识别租户
// ================================================================================
func (m *tenantContextManagerImpl) IdentifyFromPath(
	ctx context.Context,
	path string,
) (*tenantentity.Tenant, error) {
	// 1. 提取路径前缀的租户标识
	matches := pathPattern.FindStringSubmatch(path)
	if len(matches) < 2 {
		return nil, fmt.Errorf("无效的路径格式: %s", path)
	}

	tenantIdentifier := matches[1]

	// 2. 从缓存获取
	if tenant := m.getFromCache(tenantIdentifier); tenant != nil {
		logs.CtxInfof(ctx, "[TenantContext] 路径租户命中缓存: %s -> %s", tenantIdentifier, tenant.TenantID)
		return tenant, nil
	}

	// 3. 从数据库查询
	tenant, err := m.tenantRepo.GetByID(ctx, tenantIdentifier)
	if err != nil {
		return nil, fmt.Errorf("查询租户失败 (path=%s): %w", tenantIdentifier, err)
	}

	if tenant == nil {
		return nil, fmt.Errorf("租户不存在 (path=%s)", tenantIdentifier)
	}

	// 4. 缓存租户信息
	m.putToCache(tenantIdentifier, tenant)

	logs.CtxInfof(ctx, "[TenantContext] 路径租户识别成功: %s -> %s", tenantIdentifier, tenant.TenantID)
	return tenant, nil
}

// ================================================================================
// IdentifyFromHeader 从 Header 识别租户（兼容模式）
// ================================================================================
func (m *tenantContextManagerImpl) IdentifyFromHeader(
	ctx context.Context,
	headers map[string]string,
) (*tenantentity.Tenant, error) {
	// 1. 优先使用 X-Tenant-ID Header
	tenantID, ok := headers["X-Tenant-ID"]
	if !ok {
		// 兼容旧版本：使用 X-Tenant-Id
		tenantID = headers["X-Tenant-Id"]
	}

	if tenantID == "" {
		return nil, fmt.Errorf("缺少租户标识 Header")
	}

	// 2. 从缓存获取
	if tenant := m.getFromCache(tenantID); tenant != nil {
		logs.CtxInfof(ctx, "[TenantContext] Header 租户命中缓存: %s", tenantID)
		return tenant, nil
	}

	// 3. 从数据库查询
	tenant, err := m.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("查询租户失败 (header=%s): %w", tenantID, err)
	}

	if tenant == nil {
		return nil, fmt.Errorf("租户不存在 (header=%s)", tenantID)
	}

	// 4. 缓存租户信息
	m.putToCache(tenantID, tenant)

	logs.CtxInfof(ctx, "[TenantContext] Header 租户识别成功: %s", tenantID)
	return tenant, nil
}

// ================================================================================
// GetTenantContext 从上下文获取租户
// ================================================================================
func (m *tenantContextManagerImpl) GetTenantContext(
	ctx context.Context,
) (*tenantentity.Tenant, error) {
	tenant, ok := ctx.Value(TenantContextKey).(*tenantentity.Tenant)
	if !ok || tenant == nil {
		return nil, fmt.Errorf("上下文中缺少租户信息")
	}
	return tenant, nil
}

// ================================================================================
// SetTenantContext 设置租户到上下文
// ================================================================================
func (m *tenantContextManagerImpl) SetTenantContext(
	ctx context.Context,
	tenant *tenantentity.Tenant,
) context.Context {
	return context.WithValue(ctx, TenantContextKey, tenant)
}

// ================================================================================
// extractSubdomain 提取子域名
// ================================================================================
func (m *tenantContextManagerImpl) extractSubdomain(host string) string {
	// 1. 移除端口号
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// 2. 检查是否是 IP 地址（跳过）
	if isIPAddress(host) {
		return ""
	}

	// 3. 提取子域名
	matches := subdomainPattern.FindStringSubmatch(host)
	if len(matches) < 2 {
		return ""
	}

	subdomain := matches[1]

	// 4. 检查是否是基础域名
	for _, baseDomain := range m.baseDomains {
		if host == baseDomain || strings.HasSuffix(host, "."+baseDomain) {
			// 如果是基础域名，返回空（无子域名）
			return ""
		}
	}

	return subdomain
}

// ================================================================================
// isIPAddress 判断是否是 IP 地址
// ================================================================================
func isIPAddress(host string) bool {
	// 简单的 IPv4 检查
	parts := strings.Split(host, ".")
	if len(parts) != 4 {
		return false
	}
	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}
		for _, c := range part {
			if c < '0' || c > '9' {
				return false
			}
		}
	}
	return true
}

// ================================================================================
// getFromCache 从缓存获取租户
// ================================================================================
func (m *tenantContextManagerImpl) getFromCache(key string) *tenantentity.Tenant {
	tenantCacheMutex.RLock()
	defer tenantCacheMutex.RUnlock()
	return tenantCache[key]
}

// ================================================================================
// putToCache 将租户放入缓存
// ================================================================================
func (m *tenantContextManagerImpl) putToCache(key string, tenant *tenantentity.Tenant) {
	tenantCacheMutex.Lock()
	defer tenantCacheMutex.Unlock()
	tenantCache[key] = tenant
}

// ================================================================================
// ClearCache 清空租户缓存（用于测试或租户更新）
// ================================================================================
func ClearTenantCache() {
	tenantCacheMutex.Lock()
	defer tenantCacheMutex.Unlock()
	tenantCache = make(map[string]*tenantentity.Tenant)
}

// ================================================================================
// TenantIdentificationMiddleware 租户识别中间件
// ================================================================================

// TenantIdentificationConfig 租户识别配置
type TenantIdentificationConfig struct {
	TenantManager   TenantContextManager
	IdentificationMode string // "subdomain", "path", "header", "auto"
	BaseDomains     []string // 基础域名列表
	FallbackToHeader bool    // 是否回退到 Header 识别
}

// TenantIdentificationMiddleware 租户识别中间件
func TenantIdentificationMiddleware(config TenantIdentificationConfig) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var tenant *tenantentity.Tenant
		var err error

		// 1. 根据配置的识别模式进行租户识别
		switch config.IdentificationMode {
		case "subdomain":
			// 子域名模式：tenant-a.saas.coze.com
			host := string(c.Host())
			tenant, err = config.TenantManager.IdentifyFromSubdomain(ctx, host)

		case "path":
			// 路径模式：/tenant-a/bots
			path := string(c.Path())
			tenant, err = config.TenantManager.IdentifyFromPath(ctx, path)

		case "header":
			// Header 模式：X-Tenant-ID
			headers := make(map[string]string)
			headers["X-Tenant-ID"] = string(c.GetHeader("X-Tenant-ID"))
			headers["X-Tenant-Id"] = string(c.GetHeader("X-Tenant-Id"))
			tenant, err = config.TenantManager.IdentifyFromHeader(ctx, headers)

		case "auto":
			// 自动模式：依次尝试子域名、路径、Header
			// 1. 尝试子域名
			host := string(c.Host())
			tenant, err = config.TenantManager.IdentifyFromSubdomain(ctx, host)
			if err != nil || tenant == nil {
				// 2. 尝试路径
				path := string(c.Path())
				tenant, err = config.TenantManager.IdentifyFromPath(ctx, path)
			}
			if (err != nil || tenant == nil) && config.FallbackToHeader {
				// 3. 回退到 Header
				headers := make(map[string]string)
				headers["X-Tenant-ID"] = string(c.GetHeader("X-Tenant-ID"))
				headers["X-Tenant-Id"] = string(c.GetHeader("X-Tenant-Id"))
				tenant, err = config.TenantManager.IdentifyFromHeader(ctx, headers)
			}

		default:
			c.JSON(consts.StatusBadRequest, map[string]interface{}{
				"code":    "TENANT_001",
				"message": "无效的租户识别模式",
			})
			c.Abort()
			return
		}

		// 2. 处理识别失败
		if err != nil {
			logs.CtxWarnf(ctx, "[TenantMiddleware] 租户识别失败: %v", err)
			c.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    "TENANT_002",
				"message": "租户识别失败",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		if tenant == nil {
			logs.CtxWarnf(ctx, "[TenantMiddleware] 租户不存在")
			c.JSON(consts.StatusNotFound, map[string]interface{}{
				"code":    "TENANT_003",
				"message": "租户不存在",
			})
			c.Abort()
			return
		}

		// 3. 将租户信息注入上下文
		ctx = config.TenantManager.SetTenantContext(ctx, tenant)

		// 4. 将租户ID设置到 Header（供下游使用）
		c.Response.Header.Set("X-Tenant-ID", tenant.TenantID)

		// 5. 将租户上下文传递给 Hertz 的 Context
		c.Set("tenant", tenant)

		logs.CtxInfof(ctx, "[TenantMiddleware] 租户识别成功: %s (mode=%s)", tenant.TenantID, config.IdentificationMode)

		// 6. 继续执行后续中间件
		c.Next(ctx)
	}
}

// ================================================================================
// 辅助函数：从 Hertz Context 获取租户
// ================================================================================

// GetTenantFromContext 从 Hertz RequestContext 获取租户
func GetTenantFromContext(c *app.RequestContext) (*tenantentity.Tenant, error) {
	tenant, ok := c.Get("tenant")
	if !ok || tenant == nil {
		return nil, fmt.Errorf("上下文中缺少租户信息")
	}
	return tenant.(*tenantentity.Tenant), nil
}

// MustGetTenantFromContext 从 Hertz RequestContext 获取租户（panic if not found）
func MustGetTenantFromContext(c *app.RequestContext) *tenantentity.Tenant {
	tenant, err := GetTenantFromContext(c)
	if err != nil {
		panic(err)
	}
	return tenant
}
