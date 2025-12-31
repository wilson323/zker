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
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/domain/user/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	typesconsts "github.com/coze-dev/coze-studio/backend/types/consts"
	"github.com/coze-dev/coze-studio/backend/infra/tracing"
)

// SessionUserTenantIsolationMiddleware Session/User表租户隔离中间件
// 用于 Session 和 User 表迁移后的租户隔离
func SessionUserTenantIsolationMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 1. 检查是否需要跳过租户隔离
		if shouldSkipTenantIsolation(ctx) {
			ctx.Next(c)
			return
		}

		// 2. 提取tenant_id
		tenantID, err := extractTenantIDWithSessionSupport(c, ctx)
		if err != nil {
			logs.CtxErrorf(c, "[SessionUserTenantIsolation] extract tenant_id failed: %v", err)
			httputil.Unauthorized(ctx, "missing or invalid tenant_id")
			return
		}

		// 3. 验证租户
		if err := validateTenant(c, tenantID); err != nil {
			logs.CtxErrorf(c, "[SessionUserTenantIsolation] validate tenant failed: %v", err)
			httputil.Unauthorized(ctx, err.Error())
			return
		}

		// 4. 将tenant_id注入context
		ctxcache.Store(c, typesconsts.TenantIDKeyInCtx, tenantID)
		ctx.Set(TenantIDKey, tenantID)

		// 5. 记录追踪和日志
		tracing.AddSpanAttributes(c,
			tracing.AttrTenantID.String(tenantID),
		)
		logs.CtxInfof(c, "[SessionUserTenantIsolation] tenant_id=%s (from session/user)", tenantID)

		ctx.Next(c)
	}
}

// extractTenantIDWithSessionSupport 增强的tenant_id提取逻辑
// 支持 Session 和 User 表的 tenant_id
func extractTenantIDWithSessionSupport(c context.Context, ctx *app.RequestContext) (string, error) {
	// 1. 从HTTP Header获取（优先级最高）
	if tenantID := string(ctx.GetHeader(TenantIDHeader)); tenantID != "" {
		return tenantID, nil
	}

	// 2. 从Session获取（支持新的Session.TenantID）
	if session, ok := ctxcache.Get[*entity.Session](c, typesconsts.SessionDataKeyInCtx); ok {
		if session.HasTenantID() {
			logs.CtxInfof(c, "[SessionUserTenantIsolation] tenant_id from session: %s", session.TenantID)
			return session.TenantID, nil
		}
		// 兼容旧Session（没有tenant_id），尝试从users表获取
		if session.UserID > 0 {
			tenantID, err := getTenantIDByUserID(c, session.UserID)
			if err == nil && tenantID != "" {
				logs.CtxInfof(c, "[SessionUserTenantIsolation] tenant_id from user table: %s", tenantID)
				return tenantID, nil
			}
		}
	}

	// 3. 从JWT Token获取（TODO: JWT改造后启用）
	// if claims := getJWTClaims(ctx); claims != nil && claims.TenantID != "" {
	// 	return claims.TenantID, nil
	// }

	// 4. 兼容模式：使用默认租户ID（数据迁移完成后移除）
	logs.CtxWarnf(c, "[SessionUserTenantIsolation] no tenant_id found, using default: %s", DefaultTenantID)
	return DefaultTenantID, nil
}

// getTenantIDByUserID 通过用户ID获取租户ID
// 注意：这是临时实现，实际应该调用user service
func getTenantIDByUserID(c context.Context, userID int64) (string, error) {
	// TODO: 实现从user service获取用户租户ID的逻辑
	// 目前暂时返回空字符串，让上层使用默认租户ID
	logs.CtxWarnf(c, "[SessionUserTenantIsolation] getTenantIDByUserID not implemented, user_id=%d", userID)
	return "", fmt.Errorf("getTenantIDByUserID not implemented")
}

// RequireTenantIDForSessionUser Session/User表强制租户ID中间件
// 用于确保 Session 和 User 表操作必须有有效的 tenant_id
func RequireTenantIDForSessionUser() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		tenantID := GetTenantIDFromContext(c)

		// 检查是否为默认租户
		if tenantID == "" || tenantID == DefaultTenantID {
			logs.CtxWarnf(c, "[RequireTenantIDForSessionUser] missing valid tenant_id, got: %s", tenantID)

			// 检查是否为 Session/User 相关操作
			path := string(ctx.GetRequest().URI().Path())
			if isSessionUserOperation(path) {
				httputil.BadRequest(ctx,
					fmt.Sprintf("valid tenant_id is required for %s operations", getOperationType(path)),
				)
				return
			}
		}

		ctx.Next(c)
	}
}

// isSessionUserOperation 检查是否为 Session/User 操作
func isSessionUserOperation(path string) bool {
	// 检查是否为用户相关操作
	userOps := []string{
		"/api/passport/",
		"/api/user/",
		"/api/users/",
		"/api/session/",
		"/api/sessions/",
	}

	for _, op := range userOps {
		if len(path) >= len(op) && path[:len(op)] == op {
			return true
		}
	}

	return false
}

// getOperationType 获取操作类型
func getOperationType(path string) string {
	if len(path) >= 12 && path[:12] == "/api/passport" {
		return "passport"
	}
	if len(path) >= 9 && path[:9] == "/api/user" {
		return "user"
	}
	if len(path) >= 11 && path[:11] == "/api/session" {
		return "session"
	}
	return "unknown"
}

// GetTenantIDFromContextV2 增强的tenant_id获取函数
// 支持 Session 和 User 表迁移后的场景
func GetTenantIDFromContextV2(c context.Context) string {
	// 优先从Hertz context获取
	if tenantID, ok := c.Value(TenantIDKey).(string); ok {
		return tenantID
	}

	// 其次从ctxcache获取
	if tenantID, ok := ctxcache.Get[string](c, typesconsts.TenantIDKeyInCtx); ok {
		return tenantID
	}

	// 尝试从Session获取
	if session, ok := ctxcache.Get[*entity.Session](c, typesconsts.SessionDataKeyInCtx); ok {
		if session.HasTenantID() {
			return session.TenantID
		}
	}

	// 兜底：返回默认租户ID
	logs.CtxWarnf(c, "[GetTenantIDFromContextV2] no tenant_id in context, using default")
	return DefaultTenantID
}

// MustGetTenantIDFromContextV2 强制获取tenant_id
// Session/User 表迁移后必须使用此函数
func MustGetTenantIDFromContextV2(c context.Context) string {
	tenantID := GetTenantIDFromContextV2(c)
	if tenantID == "" || tenantID == DefaultTenantID {
		// 检查是否为 Session/User 操作
		// 如果是，则严格报错
		logs.CtxErrorf(c, "[MustGetTenantIDFromContextV2] invalid tenant_id: %s", tenantID)
		panic("tenant_id is required for Session/User operations")
	}
	return tenantID
}

// SessionUserTenantGetter Session/User表租户ID获取器
// 用于DAL层获取租户ID
type SessionUserTenantGetter struct {
	db *gorm.DB
}

// NewSessionUserTenantGetter 创建获取器
func NewSessionUserTenantGetter(db *gorm.DB) *SessionUserTenantGetter {
	return &SessionUserTenantGetter{db: db}
}

// GetTenantIDByUserID 根据用户ID获取租户ID
func (g *SessionUserTenantGetter) GetTenantIDByUserID(ctx context.Context, userID int64) (string, error) {
	var tenantID string
	err := g.db.WithContext(ctx).Table("user").
		Select("tenant_id").
		Where("id = ?", userID).
		Pluck("tenant_id", &tenantID).Error

	if err != nil {
		return "", err
	}

	if tenantID == "" {
		return "", fmt.Errorf("user %d has no tenant_id", userID)
	}

	return tenantID, nil
}

// GetTenantIDBySessionID 根据会话获取租户ID
func (g *SessionUserTenantGetter) GetTenantIDBySessionID(ctx context.Context, sessionID string) (string, error) {
	// 如果session表中直接有tenant_id
	var tenantID string
	err := g.db.WithContext(ctx).Table("session").
		Select("tenant_id").
		Where("id = ?", sessionID).
		Pluck("tenant_id", &tenantID).Error

	if err != nil {
		return "", err
	}

	if tenantID != "" {
		return tenantID, nil
	}

	// 兼容：通过user_id获取
	var userID int64
	err = g.db.WithContext(ctx).Table("session").
		Select("user_id").
		Where("id = ?", sessionID).
		Pluck("user_id", &userID).Error

	if err != nil {
		return "", err
	}

	return g.GetTenantIDByUserID(ctx, userID)
}

// GetUserTenantRelations 获取用户的所有租户关联
func (g *SessionUserTenantGetter) GetUserTenantRelations(
	ctx context.Context,
	userID int64,
) ([]entity.UserTenantRelation, error) {
	var relations []entity.UserTenantRelation

	err := g.db.WithContext(ctx).Raw(`
		SELECT
			user_id,
			tenant_id,
			role,
			is_default
		FROM user_tenant
		WHERE user_id = ?
		  AND deleted_at IS NULL
		ORDER BY is_default DESC, created_at ASC
	`, userID).Scan(&relations).Error

	if err != nil {
		return nil, err
	}

	return relations, nil
}

// GetUserDefaultTenant 获取用户的默认租户
func (g *SessionUserTenantGetter) GetUserDefaultTenant(ctx context.Context, userID int64) (string, error) {
	var tenantID string
	err := g.db.WithContext(ctx).Raw(`
		SELECT tenant_id
		FROM user_tenant
		WHERE user_id = ?
		  AND is_default = 1
		  AND deleted_at IS NULL
		LIMIT 1
	`, userID).Scan(&tenantID).Error

	if err != nil {
		return "", err
	}

	if tenantID == "" {
		return "", fmt.Errorf("user %d has no default tenant", userID)
	}

	return tenantID, nil
}

// SetUserDefaultTenant 设置用户的默认租户
func (g *SessionUserTenantGetter) SetUserDefaultTenant(
	ctx context.Context,
	userID int64,
	tenantID string,
) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 取消当前默认租户
		if err := tx.Exec(`
			UPDATE user_tenant
			SET is_default = 0
			WHERE user_id = ?
			  AND deleted_at IS NULL
		`, userID).Error; err != nil {
			return err
		}

		// 2. 设置新的默认租户
		if err := tx.Exec(`
			UPDATE user_tenant
			SET is_default = 1
			WHERE user_id = ?
			  AND tenant_id = ?
			  AND deleted_at IS NULL
		`, userID, tenantID).Error; err != nil {
			return err
		}

		return nil
	})
}
