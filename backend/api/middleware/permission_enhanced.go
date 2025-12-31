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
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/coze-dev/coze-studio/backend/domain/permission/service"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	berrno "github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/infra/cache"
)

var (
	// 增强权限检查器（在应用启动时注入）
	dataPermChecker    service.DataPermissionChecker
	fieldPermChecker   service.FieldPermissionChecker
	permissionCache    cache.PermissionCache
)

// InitDataPermissionMiddleware 初始化数据权限中间件依赖
func InitDataPermissionMiddleware(
	dpc service.DataPermissionChecker,
	fpc service.FieldPermissionChecker,
	pc cache.PermissionCache,
) {
	dataPermChecker = dpc
	fieldPermChecker = fpc
	permissionCache = pc
}

// DataPermissionFilterConfig 数据权限过滤配置
type DataPermissionFilterConfig struct {
	ResourceType     string // 资源类型（bots, conversations, knowledge等）
	GetPermissionLevel func(c *app.RequestContext) service.DataPermissionLevel // 获取权限级别函数
	SkipFilter       func(c *app.RequestContext) bool // 跳过过滤的条件
}

// DataPermissionFilter 数据权限过滤中间件
// 自动在请求上下文中注入数据权限过滤条件
//
// 使用示例：
//
//	r.GET("/api/bots", middleware.DataPermissionFilter(middleware.DataPermissionFilterConfig{
//	    ResourceType: "bots",
//	    GetPermissionLevel: func(c *app.RequestContext) service.DataPermissionLevel {
//	        // 从用户角色中获取权限级别
//	        return getUserPermissionLevel(c)
//	    },
//	}), handler.ListBots)
func DataPermissionFilter(config DataPermissionFilterConfig) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 获取user_id
		userIDBytes := c.GetHeader("X-User-ID")
		if len(userIDBytes) == 0 {
			c.JSON(consts.StatusUnauthorized, responseWithError(berrno.ErrUnauthorizedCode))
			c.Abort()
			return
		}
		userID := string(userIDBytes)

		// 2. 检查是否跳过过滤
		if config.SkipFilter != nil && config.SkipFilter(c) {
			c.Next(ctx)
			return
		}

		// 3. 获取权限级别
		var permissionLevel service.DataPermissionLevel
		if config.GetPermissionLevel != nil {
			permissionLevel = config.GetPermissionLevel(c)
		} else {
			// 默认从用户角色获取权限级别
			level, err := getDataPermissionLevelFromRole(ctx, userID, config.ResourceType)
			if err != nil {
				logs.CtxErrorf(ctx, "[DataPermissionFilter] get permission level failed: %v", err)
				c.JSON(consts.StatusInternalServerError, responseWithError(berrno.ErrPermissionCheckFailedCode))
				c.Abort()
				return
			}
			permissionLevel = level
		}

		// 4. 尝试从缓存获取过滤器
		var filter *service.DataPermissionFilter
		if permissionCache != nil {
			cachedFilter, err := permissionCache.GetDataPermissionFilter(ctx, userID, config.ResourceType)
			if err == nil && cachedFilter != nil {
				filter = &service.DataPermissionFilter{
					WhereClause: cachedFilter.WhereClause,
					Args:        cachedFilter.Args,
				}
				logs.CtxDebugf(ctx, "[DataPermissionFilter] cache hit for filter: user=%s, resource=%s",
					userID, config.ResourceType)
			}
		}

		// 5. 缓存未命中，生成新的过滤器
		if filter == nil && dataPermChecker != nil {
			var err error
			filter, err = dataPermChecker.Filter(ctx, userID, permissionLevel, config.ResourceType)
			if err != nil {
				logs.CtxErrorf(ctx, "[DataPermissionFilter] generate filter failed: %v", err)
				c.JSON(consts.StatusInternalServerError, responseWithError(berrno.ErrPermissionCheckFailedCode))
				c.Abort()
				return
			}

			// 6. 缓存过滤器
			if permissionCache != nil {
				filterInfo := &cache.DataFilterInfo{
					WhereClause: filter.WhereClause,
					Args:        filter.Args,
					CachedAt:    time.Now(),
				}
				_ = permissionCache.SetDataPermissionFilter(ctx, userID, config.ResourceType, filterInfo)
			}
		}

		// 7. 将过滤器存入上下文，供Handler使用
		if filter != nil {
			c.Set("data_permission_filter", filter)
			logs.CtxInfof(ctx, "[DataPermissionFilter] applied filter: user=%s, resource=%s, level=%d, where=%s",
				userID, config.ResourceType, permissionLevel, filter.WhereClause)
		}

		c.Next(ctx)
	}
}

// FieldPermissionMask 字段权限脱敏中间件
// 自动对响应数据进行字段权限过滤和脱敏
//
// 使用示例：
//
//	r.GET("/api/bots/:bot_id",
//	    middleware.FieldPermissionMask("bots"),
//	    handler.GetBot)
func FieldPermissionMask(resourceType string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 获取user_id
		userIDBytes := c.GetHeader("X-User-ID")
		if len(userIDBytes) == 0 {
			c.Next(ctx)
			return
		}
		userID := string(userIDBytes)

		// 2. 执行Handler
		c.Next(ctx)

		// 3. 检查响应状态码
		if c.Response.StatusCode() != consts.StatusOK {
			return
		}

		// 4. 解析响应数据
		var responseData map[string]interface{}
		if err := json.Unmarshal(c.Response.Body(), &responseData); err != nil {
			// 解析失败，可能是列表响应
			var responseList []map[string]interface{}
			if err2 := json.Unmarshal(c.Response.Body(), &responseList); err2 != nil {
				logs.CtxDebugf(ctx, "[FieldPermissionMask] failed to parse response: %v", err)
				return
			}

			// 处理列表响应
			maskedList := make([]map[string]interface{}, 0)
			for _, item := range responseList {
				maskedItem, err := maskFields(ctx, userID, resourceType, item)
				if err != nil {
					logs.CtxErrorf(ctx, "[FieldPermissionMask] mask fields failed: %v", err)
					return
				}
				maskedList = append(maskedList, maskedItem)
			}

			// 重新序列化响应
			newBody, _ := json.Marshal(maskedList)
			c.Response.SetBody(newBody)
			c.Response.Header.SetContentType("application/json; charset=utf-8")
			return
		}

		// 处理单个对象响应
		maskedData, err := maskFields(ctx, userID, resourceType, responseData)
		if err != nil {
			logs.CtxErrorf(ctx, "[FieldPermissionMask] mask fields failed: %v", err)
			return
		}

		// 重新序列化响应
		newBody, _ := json.Marshal(maskedData)
		c.Response.SetBody(newBody)
		c.Response.Header.SetContentType("application/json; charset=utf-8")
	}
}

// ValidateFieldPermissions 字段权限验证中间件
// 验证写入操作的字段权限
//
// 使用示例：
//
//	r.POST("/api/bots", middleware.ValidateFieldPermissions("bots"), handler.CreateBot)
func ValidateFieldPermissions(resourceType string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 获取user_id
		userIDBytes := c.GetHeader("X-User-ID")
		if len(userIDBytes) == 0 {
			c.JSON(consts.StatusUnauthorized, responseWithError(berrno.ErrUnauthorizedCode))
			c.Abort()
			return
		}
		userID := string(userIDBytes)

		// 2. 解析请求数据
		var requestData map[string]interface{}
		if err := json.Unmarshal(c.Request.Body(), &requestData); err != nil {
			logs.CtxDebugf(ctx, "[ValidateFieldPermissions] failed to parse request: %v", err)
			c.Next(ctx)
			return
		}

		// 3. 验证字段权限
		if fieldPermChecker != nil {
			if err := fieldPermChecker.ValidateFieldPermissions(ctx, userID, resourceType, requestData); err != nil {
				logs.CtxWarnf(ctx, "[ValidateFieldPermissions] validation failed: %v", err)

				if permErr, ok := err.(*service.FieldPermissionDeniedError); ok {
					c.JSON(consts.StatusForbidden, responseWithError(berrno.ErrFieldPermissionDeniedCode,
						map[string]interface{}{
							"field":  permErr.FieldName,
							"reason": permErr.Reason,
						}))
					c.Abort()
					return
				}

				c.JSON(consts.StatusForbidden, responseWithError(berrno.ErrFieldPermissionDeniedCode))
				c.Abort()
				return
			}
		}

		c.Next(ctx)
	}
}

// RequireDataPermissionLevel 要求数据权限级别的中间件
//
// 使用示例：
//
//	r.GET("/api/bots/all", middleware.RequireDataPermissionLevel(service.ALL), handler.ListAllBots)
func RequireDataPermissionLevel(requiredLevel service.DataPermissionLevel) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userIDBytes := c.GetHeader("X-User-ID")
		if len(userIDBytes) == 0 {
			c.JSON(consts.StatusUnauthorized, responseWithError(berrno.ErrUnauthorizedCode))
			c.Abort()
			return
		}
		_ = string(userIDBytes) // userID available for future use

		// 获取用户权限级别（简化示例）
		userLevel := service.SELF // 默认权限级别

		// 检查权限级别是否满足要求
		if userLevel < requiredLevel {
			c.JSON(consts.StatusForbidden, responseWithError(berrno.ErrDataPermissionDeniedCode,
				map[string]interface{}{
					"required_level": requiredLevel.String(),
					"user_level":     userLevel.String(),
				}))
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// maskFields 脱敏处理字段
func maskFields(ctx context.Context, userID, resourceType string, data map[string]interface{}) (map[string]interface{}, error) {
	if fieldPermChecker == nil {
		return data, nil
	}

	maskedData, err := fieldPermChecker.MaskSensitiveFields(ctx, userID, data)
	if err != nil {
		return nil, err
	}

	logs.CtxDebugf(ctx, "[FieldPermissionMask] masked fields for user=%s, resource=%s",
		userID, resourceType)

	return maskedData, nil
}

// getDataPermissionLevelFromRole 从用户角色获取数据权限级别
func getDataPermissionLevelFromRole(ctx context.Context, userID, resourceType string) (service.DataPermissionLevel, error) {
	// TODO: 实现从角色获取权限级别的逻辑
	// 这里是一个简化示例
	return service.SELF, nil
}

// ==========================================
// 辅助函数
// ==========================================

// responseWithError 返回错误响应
func responseWithError(errnoCode int32, kv ...interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"code":    errnoCode,
		"message": berrno.ErrMsgByCode(errnoCode),
	}

	if len(kv) > 0 {
		for _, item := range kv {
			if kvMap, ok := item.(map[string]interface{}); ok {
				for k, v := range kvMap {
					resp[k] = v
				}
			}
		}
	}

	return resp
}

// ==========================================
// 性能监控中间件
// ==========================================

// PermissionMetrics 权限性能指标
type PermissionMetrics struct {
	CheckCount      int64         // 检查次数
	CacheHitCount   int64         // 缓存命中次数
	CacheMissCount  int64         // 缓存未命中次数
	AvgCheckTime    time.Duration // 平均检查时间
	PermissionDeniedCount int64   // 权限拒绝次数
}

// PermissionMetricsMiddleware 权限性能监控中间件
func PermissionMetricsMiddleware() app.HandlerFunc {
	metrics := &PermissionMetrics{}

	return func(ctx context.Context, c *app.RequestContext) {
		startTime := time.Now()

		// 执行请求
		c.Next(ctx)

		// 记录指标
		duration := time.Since(startTime)
		metrics.CheckCount++
		metrics.AvgCheckTime = (metrics.AvgCheckTime*time.Duration(metrics.CheckCount-1) + duration) / time.Duration(metrics.CheckCount)

		// 检查是否有权限相关的响应头
		if c.Response.Header.Get("X-Permission-Cache") == "HIT" {
			metrics.CacheHitCount++
		} else if c.Response.Header.Get("X-Permission-Cache") == "MISS" {
			metrics.CacheMissCount++
		}

		if c.Response.StatusCode() == consts.StatusForbidden {
			metrics.PermissionDeniedCount++
		}

		// 添加性能响应头
		c.Response.Header.Set("X-Permission-Duration", duration.String())

		// 定期记录日志（每100次请求）
		if metrics.CheckCount%100 == 0 {
			logs.CtxInfof(ctx, "[PermissionMetrics] total=%d, cache_hit=%d, cache_miss=%d, denied=%d, avg_time=%v",
				metrics.CheckCount,
				metrics.CacheHitCount,
				metrics.CacheMissCount,
				metrics.PermissionDeniedCount,
				metrics.AvgCheckTime)
		}
	}
}
