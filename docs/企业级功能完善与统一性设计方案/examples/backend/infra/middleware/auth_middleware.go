// Package middleware 提供HTTP中间件
//
// 本文件展示认证中间件的完整实现，包括：
// - JWT令牌验证
// - 用户上下文注入
// - 租户识别（header/subdomain/path）
// - 请求追踪
package middleware

import (
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang-jwt/jwt/v5"
)

// AuthConfig 认证配置
type AuthConfig struct {
	// JWTSecret JWT密钥
	JWTSecret string
	// TokenExpiration Token过期时间（默认24小时）
	TokenExpiration time.Duration
	// RefreshTokenExpiration 刷新Token过期时间（默认7天）
	RefreshTokenExpiration time.Duration
	// TenantIdHeader 租户ID请求头名称
	TenantIdHeader string
	// EnableSubdomainTenant 是否启用子域名识别租户
	EnableSubdomainTenant bool
	// EnablePathTenant 是否启用路径识别租户
	EnablePathTenant bool
}

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	config *AuthConfig
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(config *AuthConfig) *AuthMiddleware {
	// 设置默认值
	if config.TokenExpiration == 0 {
		config.TokenExpiration = 24 * time.Hour
	}
	if config.RefreshTokenExpiration == 0 {
		config.RefreshTokenExpiration = 7 * 24 * time.Hour
	}
	if config.TenantIdHeader == "" {
		config.TenantIdHeader = "X-Tenant-ID"
	}

	return &AuthMiddleware{
		config: config,
	}
}

// Middleware 返回Hertz中间件函数
func (m *AuthMiddleware) Middleware() app.HandlerFunc {
	return func(ctx *app.RequestContext) {
		// 步骤1: 从请求中提取Token
		token := m.extractToken(ctx)
		if token == "" {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    consts.StatusUnauthorized,
				"message": "missing authorization token",
			})
			ctx.Abort()
			return
		}

		// 步骤2: 验证Token
		claims, err := m.validateToken(token)
		if err != nil {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    consts.StatusUnauthorized,
				"message": "invalid or expired token",
				"error":   err.Error(),
			})
			ctx.Abort()
			return
		}

		// 步骤3: 识别租户ID
		tenantID := m.identifyTenant(ctx)
		if tenantID == "" {
			ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
				"code":    consts.StatusBadRequest,
				"message": "unable to identify tenant",
			})
			ctx.Abort()
			return
		}

		// 步骤4: 构建用户上下文
		userContext := &UserContext{
			UserID:       claims.UserID,
			Username:     claims.Username,
			Email:        claims.Email,
			TenantID:     tenantID,
			Roles:        claims.Roles,
			Permissions:  claims.Permissions,
			TokenID:      claims.TokenID,
			LoginTime:    claims.LoginTime,
		}

		// 步骤5: 将用户上下文注入到请求上下文中
		ctx.Set("user_context", userContext)

		// 步骤6: 设置请求追踪ID（如果没有的话）
		requestID := ctx.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
			ctx.SetHeader("X-Request-ID", requestID)
		}
		ctx.Set("request_id", requestID)

		// 继续处理请求
		ctx.Next(ctx)
	}
}

// extractToken 从请求中提取Token
//
// 支持三种方式：
// 1. Authorization header (Bearer token)
// 2. Query parameter (token)
// 3. Cookie (access_token)
func (m *AuthMiddleware) extractToken(ctx *app.RequestContext) string {
	// 方式1: Authorization header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// 方式2: Query parameter
	token := ctx.Query("token")
	if token != "" {
		return token
	}

	// 方式3: Cookie
	cookie, err := ctx.Cookie("access_token")
	if err == nil && cookie != "" {
		return cookie
	}

	return ""
}

// validateToken 验证JWT Token
func (m *AuthMiddleware) validateToken(tokenString string) (*JWTClaims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(m.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// 提取claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	// 检查Token是否过期
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return nil, jwt.ErrTokenExpired
	}

	return claims, nil
}

// identifyTenant 识别租户ID
//
// 支持三种方式（按优先级）：
// 1. 请求头（X-Tenant-ID）
// 2. 子域名（tenant-id.example.com）
// 3. 路径参数（/api/v1/tenants/:tenant-id/...）
func (m *AuthMiddleware) identifyTenant(ctx *app.RequestContext) string {
	// 方式1: 请求头
	tenantID := ctx.GetHeader(m.config.TenantIdHeader)
	if tenantID != "" {
		return tenantID
	}

	// 方式2: 子域名
	if m.config.EnableSubdomainTenant {
		host := ctx.Host()
		parts := strings.Split(host, ".")
		if len(parts) >= 2 {
			// 第一个部分可能是租户ID
			potentialTenantID := parts[0]
			// 验证是否是有效的租户ID格式（例如：排除www等公共子域名）
			if potentialTenantID != "www" && potentialTenantID != "api" {
				return potentialTenantID
			}
		}
	}

	// 方式3: 路径参数
	if m.config.EnablePathTenant {
		// 从路径中提取租户ID（例如：/api/v1/tenants/tenant123/bots/...）
		path := string(ctx.Path())
		parts := strings.Split(path, "/")
		for i, part := range parts {
			if part == "tenants" && i+1 < len(parts) {
				return parts[i+1]
			}
		}
	}

	return ""
}

// ============ 辅助类型和函数 ============

// JWTClaims JWT声明
type JWTClaims struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	TokenID     string   `json:"token_id"`
	LoginTime   int64    `json:"login_time"`
	jwt.RegisteredClaims
}

// UserContext 用户上下文
type UserContext struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	TenantID    string   `json:"tenant_id"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	TokenID     string   `json:"token_id"`
	LoginTime   int64    `json:"login_time"`
}

// GetUserContext 从上下文获取用户信息
func GetUserContext(ctx *app.RequestContext) *UserContext {
	value, exists := ctx.Get("user_context")
	if !exists {
		return nil
	}
	return value.(*UserContext)
}

// generateRequestID 生成请求追踪ID
func generateRequestID() string {
	// 简化实现：使用时间戳
	// 实际项目中可以使用UUID或其他唯一ID生成算法
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// ============ Token生成工具（用于登录） ============

// GenerateToken 生成访问Token
func (m *AuthMiddleware) GenerateToken(
	userID, username, email string,
	roles, permissions []string,
) (string, error) {
	now := time.Now()
	claims := &JWTClaims{
		UserID:      userID,
		Username:    username,
		Email:       email,
		Roles:       roles,
		Permissions: permissions,
		TokenID:     generateTokenID(),
		LoginTime:   now.Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.TokenExpiration)),
			Issuer:    "zker-enterprise",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.JWTSecret))
}

// GenerateRefreshToken 生成刷新Token
func (m *AuthMiddleware) GenerateRefreshToken(
	userID string,
) (string, error) {
	now := time.Now()
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.config.RefreshTokenExpiration)),
		Issuer:    "zker-enterprise",
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.JWTSecret))
}

// generateTokenID 生成Token ID
func generateTokenID() string {
	return fmt.Sprintf("tok_%d", time.Now().UnixNano())
}
