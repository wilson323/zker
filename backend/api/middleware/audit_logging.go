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
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	hertzconsts "github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/coze-dev/coze-studio/backend/domain/audit/entity"
	"github.com/coze-dev/coze-studio/backend/domain/audit/service"
	securitysvc "github.com/coze-dev/coze-studio/backend/domain/security/service"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

// AuditConfig 审计日志配置
type AuditConfig struct {
	// 敏感操作路径(需要记录详细日志)
	SensitivePaths map[string]bool
	// 不需要记录的路径
	ExcludedPaths map[string]bool
	// 需要记录请求体的路径
	RequestBodyPaths map[string]bool
}

// AuditLoggingMiddleware 审计日志中间件
type AuditLoggingMiddleware struct {
	auditSvc   service.AuditLogService
	maskingSvc *securitysvc.DataMaskingService
	config     *AuditConfig
}

// NewAuditLoggingMiddleware 创建审计日志中间件
func NewAuditLoggingMiddleware(
	auditSvc service.AuditLogService,
	maskingSvc *securitysvc.DataMaskingService,
) *AuditLoggingMiddleware {
	config := &AuditConfig{
		SensitivePaths: map[string]bool{
			"/api/user/delete":       true,
			"/api/role/delete":       true,
			"/api/bot/delete":        true,
			"/api/knowledge/delete":  true,
			"/api/workflow/delete":   true,
			"/api/data/export":       true,
			"/api/config/update":     true,
			"/api/user/password":     true,
			"/api/mfa/enable":        true,
			"/api/mfa/disable":       true,
		},
		ExcludedPaths: map[string]bool{
			"/api/health":            true,
			"/api/metrics":           true,
			"/api/ping":              true,
			"/api/passport/login":    true, // 登录由专门的中间件处理
			"/api/static/":           true,
			"/api/favicon.ico":       true,
		},
		RequestBodyPaths: map[string]bool{
			"/api/bot/create":        true,
			"/api/bot/update":        true,
			"/api/user/create":       true,
			"/api/user/update":       true,
			"/api/role/create":       true,
			"/api/role/update":       true,
			"/api/config/update":     true,
		},
	}

	return &AuditLoggingMiddleware{
		auditSvc:   auditSvc,
		maskingSvc: maskingSvc,
		config:     config,
	}
}

// Middleware 返回Hertz中间件函数
func (m *AuditLoggingMiddleware) Middleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		startTime := time.Now()

		// 1. 检查是否需要记录日志
		path := string(ctx.Request.URI().Path())
		if m.shouldExclude(path) {
			ctx.Next(c)
			return
		}

		// 2. 获取用户信息
		userID, username, userEmail, tenantID := m.getUserInfo(ctx)

		// 3. 读取请求体(如果需要)
		var requestData string
		if m.shouldLogRequestBody(path) {
			requestData = m.getRequestBody(ctx)
		}

		// 4. 执行请求
		ctx.Next(c)

		// 5. 异步记录审计日志
		go m.logAudit(c, ctx, startTime, userID, username, userEmail, tenantID, requestData)
	}
}

// shouldExclude 判断是否排除该路径
func (m *AuditLoggingMiddleware) shouldExclude(path string) bool {
	// 检查精确匹配
	if m.config.ExcludedPaths[path] {
		return true
	}

	// 检查前缀匹配
	for excludedPath := range m.config.ExcludedPaths {
		if strings.HasSuffix(excludedPath, "/") && strings.HasPrefix(path, excludedPath) {
			return true
		}
	}

	return false
}

// shouldLogRequestBody 判断是否需要记录请求体
func (m *AuditLoggingMiddleware) shouldLogRequestBody(path string) bool {
	return m.config.RequestBodyPaths[path]
}

// getUserInfo 获取用户信息
func (m *AuditLoggingMiddleware) getUserInfo(ctx *app.RequestContext) (userID, username, userEmail, tenantID string) {
	// 从上下文获取session信息
	sessionData, exists := ctx.Get(consts.SessionDataKeyInCtx)
	if exists {
		if session, ok := sessionData.(*SessionData); ok {
			userID = session.UserID
			username = session.Username
			userEmail = session.UserEmail
			tenantID = session.TenantID
		}
	}

	// 从上下文获取租户ID
	if tenantID == "" {
		if tid, ok := ctx.Get("tenant_id"); ok {
			if tidStr, ok := tid.(string); ok {
				tenantID = tidStr
			}
		}
	}

	return
}

// getRequestBody 获取请求体
func (m *AuditLoggingMiddleware) getRequestBody(ctx *app.RequestContext) string {
	body := ctx.Request.Body()

	// 限制大小(避免记录过大的请求体)
	maxSize := 10240 // 10KB
	if len(body) > maxSize {
		body = body[:maxSize]
	}

	// 脱敏敏感数据
	masked := m.maskingSvc.SanitizeForLog(string(body))
	return masked
}

// logAudit 记录审计日志
func (m *AuditLoggingMiddleware) logAudit(
	c context.Context,
	ctx *app.RequestContext,
	startTime time.Time,
	userID, username, userEmail, tenantID string,
	requestData string,
) {
	// 构建审计日志
	log := &entity.AuditLog{
		LogID:        generateLogID(),
		TenantID:     tenantID,
		UserID:       userID,
		Username:     username,
		UserEmail:    userEmail,

		Action:       m.deriveAction(ctx),
		Resource:     m.deriveResource(ctx),
		ResourceID:   m.getResourceID(ctx),
		ResourceName: m.getResourceName(ctx),

		RequestMethod: string(ctx.Request.Method()),
		RequestPath:   string(ctx.Request.URI().Path()),
		RequestIP:     m.getClientIP(ctx),
		UserAgent:     string(ctx.Request.Header.Get("User-Agent")),

		RequestData:  requestData,
		ResponseData: m.getResponseData(ctx),

		Status:       m.deriveStatus(ctx),
		ErrorCode:    m.getErrorCode(ctx),
		ErrorMsg:     m.getErrorMsg(ctx),

		SessionID:    m.getSessionID(ctx),
		TraceID:      m.getTraceID(ctx),
	}

	// 记录日志
	if err := m.auditSvc.LogOperationAsync(c, log); err != nil {
		logs.Errorf("[AuditLoggingMiddleware] failed to log audit: %v", err)
	}
}

// deriveAction 推断操作类型
func (m *AuditLoggingMiddleware) deriveAction(ctx *app.RequestContext) entity.AuditAction {
	path := string(ctx.Request.URI().Path())
	method := string(ctx.Request.Method())

	// 根据路径和方法推断操作
	switch {
	case strings.Contains(path, "/user/") && strings.Contains(path, "/delete"):
		return entity.AuditActionUserDelete
	case strings.Contains(path, "/user/") && strings.Contains(path, "/password"):
		return entity.AuditActionUserPassword
	case strings.Contains(path, "/user/") && method == "POST":
		return entity.AuditActionUserCreate
	case strings.Contains(path, "/user/") && method == "PUT":
		return entity.AuditActionUserUpdate
	case strings.Contains(path, "/role/") && strings.Contains(path, "/delete"):
		return entity.AuditActionRoleDelete
	case strings.Contains(path, "/role/") && method == "POST":
		return entity.AuditActionRoleCreate
	case strings.Contains(path, "/role/") && method == "PUT":
		return entity.AuditActionRoleUpdate
	case strings.Contains(path, "/bot/") && method == "DELETE":
		return entity.AuditActionDataDelete
	case strings.Contains(path, "/bot/") && method == "POST":
		return entity.AuditActionDataCreate
	case strings.Contains(path, "/bot/") && method == "PUT":
		return entity.AuditActionDataUpdate
	case strings.Contains(path, "/export"):
		return entity.AuditActionDataExport
	case strings.Contains(path, "/import"):
		return entity.AuditActionDataImport
	case strings.Contains(path, "/config/") && method == "PUT":
		return entity.AuditActionConfigUpdate
	case strings.Contains(path, "/config/") && method == "DELETE":
		return entity.AuditActionConfigDelete
	default:
		return entity.AuditActionDataQuery
	}
}

// deriveResource 推断资源类型
func (m *AuditLoggingMiddleware) deriveResource(ctx *app.RequestContext) entity.AuditResource {
	path := string(ctx.Request.URI().Path())

	switch {
	case strings.Contains(path, "/user/"):
		return entity.AuditResourceUser
	case strings.Contains(path, "/role/"):
		return entity.AuditResourceRole
	case strings.Contains(path, "/tenant/"):
		return entity.AuditResourceTenant
	case strings.Contains(path, "/bot/"):
		return entity.AuditResourceBot
	case strings.Contains(path, "/conversation/"):
		return entity.AuditResourceConversation
	case strings.Contains(path, "/knowledge/"):
		return entity.AuditResourceKnowledge
	case strings.Contains(path, "/workflow/"):
		return entity.AuditResourceWorkflow
	case strings.Contains(path, "/database/"):
		return entity.AuditResourceDatabase
	case strings.Contains(path, "/config/"):
		return entity.AuditResourceConfig
	default:
		return entity.AuditResourceSystem
	}
}

// getResourceID 获取资源ID
func (m *AuditLoggingMiddleware) getResourceID(ctx *app.RequestContext) string {
	// 从路径参数中提取资源ID
	path := string(ctx.Request.URI().Path())
	parts := strings.Split(path, "/")

	// 示例: /api/bot/{bot_id} -> 提取bot_id
	for i, part := range parts {
		if part == "bot" && i+1 < len(parts) {
			return parts[i+1]
		}
		if part == "user" && i+1 < len(parts) {
			return parts[i+1]
		}
		if part == "role" && i+1 < len(parts) {
			return parts[i+1]
		}
		// 可以继续添加其他资源类型
	}

	return ""
}

// getResourceName 获取资源名称(从响应体中提取)
func (m *AuditLoggingMiddleware) getResourceName(ctx *app.RequestContext) string {
	// 简化实现，实际应该从响应体中解析
	return ""
}

// getClientIP 获取客户端IP
func (m *AuditLoggingMiddleware) getClientIP(ctx *app.RequestContext) string {
	// 尝试从X-Forwarded-For获取
	if xff := ctx.GetHeader("X-Forwarded-For"); len(xff) > 0 {
		ips := strings.Split(string(xff), ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 尝试从X-Real-IP获取
	if xri := ctx.GetHeader("X-Real-IP"); len(xri) > 0 {
		return string(xri)
	}

	// 从RemoteAddr获取
	if addr := ctx.RemoteAddr(); addr != nil {
		return addr.String()
	}
	return ""
}

// getResponseData 获取响应数据
func (m *AuditLoggingMiddleware) getResponseData(ctx *app.RequestContext) string {
	// 只记录失败的响应
	statusCode := ctx.Response.StatusCode()
	if statusCode >= 400 {
		body := ctx.Response.Body()
		if len(body) > 0 {
			// 限制大小
			maxSize := 2048 // 2KB
			if len(body) > maxSize {
				body = body[:maxSize]
			}
			return m.maskingSvc.SanitizeForLog(string(body))
		}
	}
	return ""
}

// deriveStatus 推断状态
func (m *AuditLoggingMiddleware) deriveStatus(ctx *app.RequestContext) entity.AuditStatus {
	statusCode := ctx.Response.StatusCode()
	if statusCode >= 200 && statusCode < 300 {
		return entity.AuditStatusSuccess
	}
	if statusCode >= 400 && statusCode < 500 {
		return entity.AuditStatusFailed
	}
	return entity.AuditStatusPending
}

// getErrorCode 获取错误码
func (m *AuditLoggingMiddleware) getErrorCode(ctx *app.RequestContext) string {
	statusCode := ctx.Response.StatusCode()
	if statusCode >= 400 {
		// 尝试从响应体中解析错误码
		var resp map[string]interface{}
		if err := json.Unmarshal(ctx.Response.Body(), &resp); err == nil {
			if code, ok := resp["code"]; ok {
				return fmt.Sprintf("%v", code)
			}
		}
		return fmt.Sprintf("HTTP_%d", statusCode)
	}
	return ""
}

// getErrorMsg 获取错误消息
func (m *AuditLoggingMiddleware) getErrorMsg(ctx *app.RequestContext) string {
	statusCode := ctx.Response.StatusCode()
	if statusCode >= 400 {
		var resp map[string]interface{}
		if err := json.Unmarshal(ctx.Response.Body(), &resp); err == nil {
			if msg, ok := resp["message"]; ok {
				return fmt.Sprintf("%v", msg)
			}
			if msg, ok := resp["msg"]; ok {
				return fmt.Sprintf("%v", msg)
			}
		}
		return hertzconsts.StatusMessage(statusCode)
	}
	return ""
}

// getSessionID 获取会话ID
func (m *AuditLoggingMiddleware) getSessionID(ctx *app.RequestContext) string {
	if sessionID := ctx.GetHeader("X-Session-ID"); len(sessionID) > 0 {
		return string(sessionID)
	}
	return ""
}

// getTraceID 获取追踪ID
func (m *AuditLoggingMiddleware) getTraceID(ctx *app.RequestContext) string {
	if traceID := ctx.GetHeader("X-Trace-ID"); len(traceID) > 0 {
		return string(traceID)
	}
	return ""
}

// generateLogID 生成日志ID
func generateLogID() string {
	return fmt.Sprintf("audit_%d", time.Now().UnixNano())
}

// SessionData 会话数据(简化)
type SessionData struct {
	UserID    string
	Username  string
	UserEmail string
	TenantID  string
}

// StatusMessage 获取HTTP状态码消息
func StatusMessage(statusCode int) string {
	switch statusCode {
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return hertzconsts.StatusMessage(statusCode)
	}
}
