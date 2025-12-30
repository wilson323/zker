// backend/domain/tenant/migration/routing_middleware.go
// API网关集成中间件 - 灰度发布路由
package migration

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// RoutingMiddleware 路由中间件
type RoutingMiddleware struct {
	grayController *GrayReleaseController
	logger         logs.CtxLogger
	enableMetrics  bool
}

// NewRoutingMiddleware 创建路由中间件
func NewRoutingMiddleware(grayController *GrayReleaseController) *RoutingMiddleware {
	return &RoutingMiddleware{
		grayController: grayController,
		logger:         logs.DefaultLogger(),
		enableMetrics:  true,
	}
}

// Middleware 返回Hertz中间件函数
func (m *RoutingMiddleware) Middleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		startTime := time.Now()

		// 1. 从请求中提取tenant_id
		tenantID := m.extractTenantID(c)
		if tenantID == "" {
			m.logger.CtxWarnf(ctx, "[Routing] 请求缺少tenant_id")
			c.Next(ctx)
			return
		}

		// 2. 判断是否路由到新架构
		shouldRouteToNew := m.grayController.ShouldRouteToNew(ctx, tenantID)

		// 3. 设置路由标识到context
		c.Set("route_to_new", shouldRouteToNew)
		c.Set("tenant_id", tenantID)
		c.Set("routing_stage", m.grayController.GetCurrentStage())

		// 4. 记录路由决策
		m.logger.CtxDebugf(ctx, "[Routing] tenant_id=%s, route_to_new=%v, stage=%d",
			tenantID, shouldRouteToNew, m.grayController.GetCurrentStage())

		// 5. 执行请求
		c.Next(ctx)

		// 6. 记录指标
		if m.enableMetrics {
			latency := time.Since(startTime).Seconds()
			m.grayController.RecordMetrics(shouldRouteToNew, latency)
		}
	}
}

// extractTenantID 从请求中提取tenant_id
func (m *RoutingMiddleware) extractTenantID(c *app.RequestContext) string {
	// 优先级:
	// 1. Header
	// 2. Query Parameter
	// 3. JWT Token (如果使用JWT)
	// 4. Cookie

	// 1. 从Header提取
	if tenantIDBytes := c.GetHeader("X-Tenant-ID"); len(tenantIDBytes) > 0 {
		return string(tenantIDBytes)
	}

	// 2. 从Query参数提取
	if tenantIDBytes := c.Query("tenant_id"); len(tenantIDBytes) > 0 {
		return string(tenantIDBytes)
	}

	// 3. 从Cookie提取
	if cookieValueBytes := c.Cookie("tenant_id"); len(cookieValueBytes) > 0 {
		return string(cookieValueBytes)
	}

	// 4. 从路径参数提取（例如: /api/v1/tenants/:tenant_id/bots）
	// 这里简化处理，实际应该从路径参数中提取

	return ""
}

// IsRouteToNew 检查是否路由到新架构
func IsRouteToNew(ctx context.Context, c *app.RequestContext) bool {
	if v, exists := c.Get("route_to_new"); exists {
		if routeToNew, ok := v.(bool); ok {
			return routeToNew
		}
	}
	return false
}

// GetTenantID 获取tenant_id
func GetTenantID(ctx context.Context, c *app.RequestContext) string {
	if v, exists := c.Get("tenant_id"); exists {
		if tenantID, ok := v.(string); ok {
			return tenantID
		}
	}
	return ""
}

// GetRoutingStage 获取当前路由阶段
func GetRoutingStage(ctx context.Context, c *app.RequestContext) int {
	if v, exists := c.Get("routing_stage"); exists {
		if stage, ok := v.(int); ok {
			return stage
		}
	}
	return 0
}

// HashBasedRouting 基于哈希的路由决策
func (m *RoutingMiddleware) HashBasedRouting(tenantID string, percentage int) bool {
	// 使用tenant_id的哈希值确保一致性路由
	hash := sha256.Sum256([]byte(tenantID))
	hashStr := hex.EncodeToString(hash[:])

	// 取哈希值的前几位计算百分比
	// 简化实现：取第一个字符的数值
	hashValue := int(hashStr[0])
	return hashValue%100 < percentage
}

// UserBasedRouting 基于用户ID的路由决策
func (m *RoutingMiddleware) UserBasedRouting(userID string, percentage int) bool {
	// 使用user_id的哈希值确保同一用户始终路由到同一架构
	hash := sha256.Sum256([]byte(userID))
	hashValue := int(hash[0]) % 100
	return hashValue < percentage
}

// ============================================
// 灰度发布管理API
// ============================================

// GrayReleaseAPI 灰度发布管理API
type GrayReleaseAPI struct {
	grayController *GrayReleaseController
	whitelist      *TenantWhitelist
	logger         logs.CtxLogger
}

// NewGrayReleaseAPI 创建灰度发布API
func NewGrayReleaseAPI(grayController *GrayReleaseController, whitelist *TenantWhitelist) *GrayReleaseAPI {
	return &GrayReleaseAPI{
		grayController: grayController,
		whitelist:      whitelist,
		logger:         logs.DefaultLogger(),
	}
}

// StartGrayRelease 启动灰度发布
// POST /api/v1/admin/gray-release/start
func (api *GrayReleaseAPI) StartGrayRelease(ctx context.Context, c *app.RequestContext) {
	if err := api.grayController.StartGrayRelease(ctx); err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    "GRAY_RELEASE_START_FAILED",
			"message": "启动灰度发布失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    "SUCCESS",
		"message": "灰度发布已启动",
		"stage":   api.grayController.GetCurrentStage(),
	})
}

// PromoteStage 晋级到下一阶段
// POST /api/v1/admin/gray-release/promote
func (api *GrayReleaseAPI) PromoteStage(ctx context.Context, c *app.RequestContext) {
	currentStage := api.grayController.GetCurrentStage()
	if currentStage == 0 {
		c.JSON(400, map[string]interface{}{
			"code":    "GRAY_RELEASE_NOT_STARTED",
			"message": "灰度发布未启动",
		})
		return
	}

	if err := api.grayController.PromoteStage(ctx, currentStage); err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    "GRAY_RELEASE_PROMOTE_FAILED",
			"message": "阶段晋升失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    "SUCCESS",
		"message": fmt.Sprintf("已晋升到阶段 %d", api.grayController.GetCurrentStage()),
		"stage":   api.grayController.GetCurrentStage(),
	})
}

// Rollback 回滚
// POST /api/v1/admin/gray-release/rollback
func (api *GrayReleaseAPI) Rollback(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Reason string `json:"reason"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    "INVALID_PARAMS",
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	if err := api.grayController.Rollback(ctx, req.Reason); err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    "GRAY_RELEASE_ROLLBACK_FAILED",
			"message": "回滚失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    "SUCCESS",
		"message": "已回滚",
		"stage":   api.grayController.GetCurrentStage(),
	})
}

// GetStatus 获取灰度发布状态
// GET /api/v1/admin/gray-release/status
func (api *GrayReleaseAPI) GetStatus(ctx context.Context, c *app.RequestContext) {
	stages := api.grayController.GetStages()
	metrics := api.grayController.GetMetrics()
	currentStage := api.grayController.GetCurrentStage()

	c.JSON(200, map[string]interface{}{
		"code":          "SUCCESS",
		"current_stage": currentStage,
		"stages":        stages,
		"metrics":       metrics,
	})
}

// AddWhitelist 添加租户到白名单
// POST /api/v1/admin/gray-release/whitelist
func (api *GrayReleaseAPI) AddWhitelist(ctx context.Context, c *app.RequestContext) {
	var req struct {
		TenantID string `json:"tenant_id" binding:"required"`
		Reason   string `json:"reason"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    "INVALID_PARAMS",
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	if err := api.whitelist.Add(ctx, req.TenantID, req.Reason, "admin"); err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    "WHITELIST_ADD_FAILED",
			"message": "添加白名单失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":      "SUCCESS",
		"message":   "租户已添加到白名单",
		"tenant_id": req.TenantID,
	})
}

// GetWhitelist 获取白名单
// GET /api/v1/admin/gray-release/whitelist
func (api *GrayReleaseAPI) GetWhitelist(ctx context.Context, c *app.RequestContext) {
	entries, err := api.whitelist.GetAll(ctx)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    "WHITELIST_GET_FAILED",
			"message": "获取白名单失败",
			"error":   err.Error(),
		})
		return
	}

	count, _ := api.whitelist.GetCount(ctx)

	c.JSON(200, map[string]interface{}{
		"code":    "SUCCESS",
		"count":   count,
		"entries": entries,
	})
}
