# API接口设计规范

**文档编号**: DE-API-2025-SPEC-001
**版本**: v1.0.0
**创建日期**: 2025-12-30
**适用范围**: Coze Studio 全平台API接口设计

---

## 📋 目录

1. [设计原则](#1-设计原则)
2. [RESTful规范](#2-restful规范)
3. [URL设计规范](#3-url设计规范)
4. [HTTP方法使用](#4-http方法使用)
5. [状态码规范](#5-状态码规范)
6. [请求格式规范](#6-请求格式规范)
7. [响应格式规范](#7-响应格式规范)
8. [错误处理规范](#8-错误处理规范)
9. [认证与授权](#9-认证与授权)
10. [API版本控制](#10-api版本控制)
11. [分页规范](#11-分页规范)
12. [排序与过滤](#12-排序与过滤)
13. [限流与节流](#13-限流与节流)
14. [幂等性设计](#14-幂等性设计)
15. [API文档规范](#15-api文档规范)

---

## 1. 设计原则

### 1.1 核心原则

- **资源导向**: URL表示资源,使用名词而非动词
- **统一接口**: 遵循RESTful架构风格
- **无状态**: 每个请求包含完整信息,服务器不保存客户端状态
- **分层系统**: 客户端无法判断是否连接到中间服务器
- **按需代码**: 可选,通过下载代码扩展功能

### 1.2 API设计理念

```yaml
设计理念:
  - 简洁性: API简单易用,学习成本低
  - 一致性: 所有接口遵循统一的设计模式
  - 可预测性: 接口行为可从URL和HTTP方法推断
  - 可扩展性: 支持未来功能扩展
  - 向后兼容: 新版本不破坏旧版本客户端
```

---

## 2. RESTful规范

### 2.1 资源命名

**规则:**
- ✅ 使用名词复数形式
- ✅ 使用小写字母
- ✅ 使用连字符(-)分隔多词
- ✅ 资源层级不超过3层

```bash
# ✅ 正确示例
GET    /api/v1/bots                    # 获取Bot列表
GET    /api/v1/bots/{id}               # 获取指定Bot
POST   /api/v1/bots                    # 创建Bot
PUT    /api/v1/bots/{id}               # 更新Bot
DELETE /api/v1/bots/{id}               # 删除Bot
GET    /api/v1/bots/{id}/conversations # 获取Bot的对话列表
POST   /api/v1/bots/{id}/publish       # 发布Bot到渠道

# ❌ 错误示例
GET /api/v1/getBots                   # 使用了动词
GET /api/v1/Bot                       # 使用了单数
GET /api/v1/bot_conversations         # 使用了下划线
GET /api/v1/bots/{id}/conversations/{conversationId}/messages/{msgId}/attachments # 层级过深
```

### 2.2 URL层级设计

```bash
# 标准层级: /api/{version}/{resource}/{id}/{sub-resource}/{sub-id}
GET    /api/v1/bots/{bot_id}/conversations/{conversation_id}/messages

# 特殊操作使用子资源而非动词
POST   /api/v1/bots/{bot_id}/actions/publish   # ✅ 发布Bot
POST   /api/v1/subscriptions/{id}/cancel       # ✅ 取消订阅
```

---

## 3. URL设计规范

### 3.1 URL结构

```bash
# 标准URL格式
{scheme}://{host}/{api}/{version}/{resource}/{id}?{query_params}

# 实际示例
https://api.coze.com/api/v1/bots/123?include=published_conversations
```

### 3.2 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| API路径前缀 | `/api` | `/api/v1/bots` |
| 版本号 | `/v{major}` | `/v1`, `/v2` |
| 资源名 | 小写+连字符+复数 | `bot-templates`, `speech-records` |
| 资源ID | `{resource}_id` | `{bot_id}`, `{conversation_id}` |
| 查询参数 | 小写+下划线 | `page_size`, `sort_by` |

### 3.3 特殊场景URL设计

```bash
# 1. 搜索操作
GET /api/v1/bots/search?q=chatbot&category=productivity

# 2. 统计操作
GET /api/v1/bots/{id}/analytics

# 3. 批量操作
POST /api/v1/bots/batch/delete
Body: {"bot_ids": [1, 2, 3]}

# 4. 异步操作
POST /api/v1/bots/{id}/train
Response: {"job_id": "xxx", "status": "processing"}
GET /api/v1/jobs/{job_id}

# 5. 导入导出
POST /api/v1/bots/import
GET /api/v1/bots/export?format=json
```

---

## 4. HTTP方法使用

### 4.1 方法映射

| HTTP方法 | 操作 | 幂等性 | 安全性 | 请求体 | 响应体 |
|----------|------|--------|--------|--------|--------|
| GET | 查询 | ✅ | ✅ | ❌ | ✅ |
| POST | 创建 | ❌ | ❌ | ✅ | ✅ |
| PUT | 完整更新 | ✅ | ❌ | ✅ | ✅ |
| PATCH | 部分更新 | ❌ | ❌ | ✅ | ✅ |
| DELETE | 删除 | ✅ | ❌ | ❌ | ✅ |
| HEAD | 获取元数据 | ✅ | ✅ | ❌ | 仅Header |
| OPTIONS | 获取支持方法 | ✅ | ✅ | ❌ | ✅ |

### 4.2 使用示例

```bash
# GET - 查询资源
GET /api/v1/bots
GET /api/v1/bots/{id}
GET /api/v1/bots?status=published&category=assistant

# POST - 创建资源
POST /api/v1/bots
Content-Type: application/json
{
  "name": "Customer Service Bot",
  "description": "AI assistant for customer support",
  "avatar": "https://cdn.example.com/avatar.png"
}

# PUT - 完整更新(需要提供所有必填字段)
PUT /api/v1/bots/{id}
Content-Type: application/json
{
  "name": "Updated Bot Name",
  "description": "Updated description",
  "avatar": "https://cdn.example.com/new-avatar.png",
  "config": {...}  # 所有必填字段
}

# PATCH - 部分更新(仅更新提供的字段)
PATCH /api/v1/bots/{id}
Content-Type: application/json
{
  "name": "Updated Bot Name"  # 仅更新name字段
}

# DELETE - 删除资源
DELETE /api/v1/bots/{id}

# 批量删除
DELETE /api/v1/bots?bot_ids=1,2,3
```

---

## 5. 状态码规范

### 5.1 状态码分类

| 分类 | 状态码 | 含义 |
|------|--------|------|
| 成功 | 2xx | 请求成功 |
| 重定向 | 3xx | 需要进一步操作 |
| 客户端错误 | 4xx | 客户端请求错误 |
| 服务器错误 | 5xx | 服务器处理错误 |

### 5.2 常用状态码详解

#### 成功响应 (2xx)

| 状态码 | 名称 | 使用场景 | 响应体 |
|--------|------|----------|--------|
| 200 | OK | GET/PUT/PATCH成功 | ✅ 返回资源数据 |
| 201 | Created | POST创建成功 | ✅ 返回新创建资源 |
| 202 | Accepted | 异步操作已接受 | ✅ 返回任务ID |
| 204 | No Content | DELETE/PATCH成功 | ❌ 无响应体 |

```json
// 200 OK 示例
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 123,
    "name": "My Bot",
    "status": "published"
  }
}

// 201 Created 示例
{
  "code": 0,
  "message": "Bot created successfully",
  "data": {
    "id": 456,
    "created_at": "2025-12-30T10:00:00Z"
  },
  "meta": {
    "location": "/api/v1/bots/456"
  }
}

// 202 Accepted 示例
{
  "code": 0,
  "message": "Bot training started",
  "data": {
    "job_id": "train-20251230-001",
    "status": "processing",
    "estimated_time": 120
  }
}
```

#### 客户端错误 (4xx)

| 状态码 | 名称 | 使用场景 | 业务错误码 |
|--------|------|----------|------------|
| 400 | Bad Request | 请求参数错误 | 40001-49999 |
| 401 | Unauthorized | 未认证/Token过期 | 40101-40199 |
| 403 | Forbidden | 无权限访问 | 40301-40399 |
| 404 | Not Found | 资源不存在 | 40401-40499 |
| 409 | Conflict | 资源冲突(如重复创建) | 40901-40999 |
| 422 | Unprocessable Entity | 业务逻辑验证失败 | 42201-42299 |
| 429 | Too Many Requests | 超过限流 | 42901-42999 |

```json
// 400 Bad Request 示例
{
  "code": 40001,
  "message": "Invalid request parameters",
  "errors": [
    {
      "field": "bot_name",
      "message": "Bot name must be between 2 and 50 characters"
    },
    {
      "field": "avatar",
      "message": "Avatar URL is invalid"
    }
  ],
  "request_id": "req-abc123",
  "timestamp": "2025-12-30T10:00:00Z"
}

// 401 Unauthorized 示例
{
  "code": 40101,
  "message": "Authentication token expired",
  "data": {
    "refresh_token_required": true
  },
  "request_id": "req-def456",
  "timestamp": "2025-12-30T10:00:00Z"
}

// 403 Forbidden 示例
{
  "code": 40303,
  "message": "You don't have permission to access this bot",
  "data": {
    "required_permission": "bot:write",
    "user_permissions": ["bot:read"]
  },
  "request_id": "req-ghi789",
  "timestamp": "2025-12-30T10:00:00Z"
}

// 404 Not Found 示例
{
  "code": 40401,
  "message": "Bot not found",
  "data": {
    "resource_type": "bot",
    "resource_id": 999
  },
  "request_id": "req-jkl012",
  "timestamp": "2025-12-30T10:00:00Z"
}

// 422 Unprocessable Entity 示例
{
  "code": 42201,
  "message": "Business rule validation failed",
  "errors": [
    {
      "field": "status",
      "message": "Cannot publish bot without completing required configuration"
    }
  ],
  "request_id": "req-mno345",
  "timestamp": "2025-12-30T10:00:00Z"
}

// 429 Too Many Requests 示例
{
  "code": 42901,
  "message": "Rate limit exceeded",
  "data": {
    "limit": 100,
    "remaining": 0,
    "reset_at": "2025-12-30T10:01:00Z",
    "retry_after": 45
  },
  "request_id": "req-pqr678",
  "timestamp": "2025-12-30T10:00:00Z"
}
```

#### 服务器错误 (5xx)

| 状态码 | 名称 | 使用场景 |
|--------|------|----------|
| 500 | Internal Server Error | 服务器内部错误 |
| 502 | Bad Gateway | 上游服务错误 |
| 503 | Service Unavailable | 服务不可用(维护中) |
| 504 | Gateway Timeout | 上游服务超时 |

```json
// 500 Internal Server Error 示例
{
  "code": 50001,
  "message": "Internal server error",
  "data": {
    "error_id": "err-internal-xyz789",
    "support_contact": "support@coze.com"
  },
  "request_id": "req-stu901",
  "timestamp": "2025-12-30T10:00:00Z"
}
```

### 5.3 业务错误码规范

```yaml
错误码格式: "{HTTP_STATUS}{SEQUENCE}"

分类:
  # 通用错误 (40001-40099)
  40001: "请求参数缺失"
  40002: "请求参数格式错误"
  40003: "请求参数值无效"
  40004: "请求体过大"
  40005: "Content-Type不支持"

  # 认证错误 (40101-40199)
  40101: "Token缺失"
  40102: "Token无效"
  40103: "Token过期"
  40104: "Refresh Token无效"
  40105: "Refresh Token过期"

  # 权限错误 (40301-40399)
  40301: "无权限访问资源"
  40302: "角色权限不足"
  40303: "跨租户访问被拒绝"
  40304: "资源已被锁定"

  # 资源错误 (40401-40499)
  40401: "资源不存在"
  40402: "API端点不存在"
  40403: "关联资源不存在"

  # 冲突错误 (40901-40999)
  40901: "资源已存在"
  40902: "资源版本冲突"
  40903: "资源状态冲突"

  # 业务验证错误 (42201-42299)
  42201: "业务规则验证失败"
  42202: "操作不允许"
  42203: "配额已用尽"
  42204: "订阅已过期"

  # 限流错误 (42901-42999)
  42901: "超过API调用频率限制"
  42902: "超过并发请求数限制"
  42903: "超过配额限制"

  # 服务器错误 (50001-50099)
  50001: "内部服务器错误"
  50002: "数据库错误"
  50003: "第三方服务错误"
  50004: "消息队列错误"
```

---

## 6. 请求格式规范

### 6.1 请求头

```http
# 必须的请求头
Content-Type: application/json          # 请求体格式
Accept: application/json                # 期望的响应格式
Authorization: Bearer {access_token}    # 认证Token
X-Request-ID: {uuid}                    # 请求追踪ID(可选)
X-Client-Version: 1.0.0                 # 客户端版本(可选)

# 多租户场景
X-Tenant-ID: {tenant_id}                # 租户ID(从Token解析,可选)

# 国际化
Accept-Language: zh-CN                  # 语言偏好
```

### 6.2 请求体格式

#### 标准请求体结构

```json
{
  "data": {
    // 业务数据字段
    "name": "My Bot",
    "description": "Customer service assistant"
  }
}
```

#### 批量操作请求体

```json
{
  "data": {
    "bot_ids": [1, 2, 3],
    "operation": "delete"
  }
}
```

#### 部分更新(PATCH)请求体

```json
{
  "data": {
    "name": "Updated Name",
    "status": "published"  // 仅更新提供的字段
  }
}
```

### 6.3 请求体字段规范

| 字段类型 | 数据类型 | 命名规范 | 示例 |
|----------|----------|----------|------|
| 资源ID | Long | `{resource}_id` | `bot_id: 123` |
| 名称 | String | `name` | `"Customer Service Bot"` |
| 描述 | String | `description` | `"AI assistant for support"` |
| 状态 | Enum | `status` | `"active"` |
| 时间 | ISO8601 | `{action}_at` | `"2025-12-30T10:00:00Z"` |
| 布尔值 | Boolean | `is_{feature}` | `is_public: true` |
| 数量 | Integer | `{resource}_count` | `message_count: 100` |
| 列表 | Array | `{resource}s` | `tags: ["ai", "chatbot"]` |
| 配置 | Object | `config` | `{"model": "gpt-4"}` |

---

## 7. 响应格式规范

### 7.1 成功响应结构

```json
{
  "code": 0,
  "message": "success",
  "data": {
    // 业务数据
  },
  "meta": {
    // 元数据(分页、统计等)
  }
}
```

### 7.2 单资源响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 123,
    "name": "Customer Service Bot",
    "description": "AI-powered customer assistant",
    "avatar": "https://cdn.coze.com/bot/123/avatar.png",
    "status": "published",
    "created_at": "2025-12-30T10:00:00Z",
    "updated_at": "2025-12-30T11:00:00Z"
  }
}
```

### 7.3 资源列表响应(带分页)

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "Bot 1",
      "status": "published"
    },
    {
      "id": 2,
      "name": "Bot 2",
      "status": "draft"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_count": 100,
      "total_pages": 5,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

### 7.4 批量操作响应

```json
{
  "code": 0,
  "message": "Batch operation completed",
  "data": {
    "successful_count": 8,
    "failed_count": 2,
    "results": [
      {
        "id": 1,
        "status": "success"
      },
      {
        "id": 2,
        "status": "failed",
        "error": {
          "code": 42201,
          "message": "Bot is published, cannot delete"
        }
      }
    ]
  }
}
```

### 7.5 异步操作响应

```json
{
  "code": 0,
  "message": "Job started",
  "data": {
    "job_id": "train-20251230-001",
    "status": "processing",
    "progress": 0,
    "estimated_time": 120,
    "created_at": "2025-12-30T10:00:00Z"
  },
  "meta": {
    "job_status_url": "/api/v1/jobs/train-20251230-001",
    "websocket_url": "wss://api.coze.com/ws/jobs/train-20251230-001"
  }
}
```

---

## 8. 错误处理规范

### 8.1 错误响应结构

```json
{
  "code": 40001,
  "message": "Invalid request parameters",
  "errors": [
    {
      "field": "bot_name",
      "message": "Bot name must be between 2 and 50 characters",
      "rejected_value": "",
      "constraint": {
        "min_length": 2,
        "max_length": 50
      }
    }
  ],
  "data": {
    // 额外的错误上下文数据
  },
  "request_id": "req-abc123",
  "timestamp": "2025-12-30T10:00:00Z"
}
```

### 8.2 错误处理最佳实践

```yaml
开发规范:
  1. 统一错误处理:
     - 所有错误通过统一的错误处理器返回
     - 使用预定义的业务错误码
     - 记录错误日志(包含request_id)

  2. 错误信息设计:
     - 对用户友好,避免技术术语
     - 提供明确的错误原因
     - 给出修复建议

  3. 敏感信息保护:
     - 不暴露内部实现细节
     - 不暴露数据库结构
     - 不暴露敏感的系统信息

  4. 错误监控:
     - 所有5xx错误记录并告警
     - 4xx错误统计分析
     - 异常错误模式告警
```

### 8.3 Go后端错误处理示例

```go
// 定义业务错误类型
type BusinessError struct {
    Code       int                    `json:"code"`
    Message    string                 `json:"message"`
    Errors     []FieldError           `json:"errors,omitempty"`
    Data       map[string]interface{} `json:"data,omitempty"`
    RequestID  string                 `json:"request_id"`
    Timestamp  time.Time              `json:"timestamp"`
}

type FieldError struct {
    Field         string                 `json:"field"`
    Message       string                 `json:"message"`
    RejectedValue interface{}            `json:"rejected_value,omitempty"`
    Constraint    map[string]interface{} `json:"constraint,omitempty"`
}

// 预定义错误码
const (
    ErrCodeInvalidParams    = 40001
    ErrCodeUnauthorized     = 40101
    ErrCodeTokenExpired     = 40103
    ErrCodeForbidden        = 40301
    ErrCodeNotFound         = 40401
    ErrCodeResourceConflict = 40901
    ErrCodeValidationFailed = 42201
    ErrCodeRateLimitExceed  = 42901
)

// 创建错误响应
func NewBusinessError(code int, message string) *BusinessError {
    return &BusinessError{
        Code:      code,
        Message:   message,
        RequestID: getRequestID(),
        Timestamp: time.Now(),
    }
}

func (e *BusinessError) AddFieldError(field string, message string) *BusinessError {
    e.Errors = append(e.Errors, FieldError{
        Field:   field,
        Message: message,
    })
    return e
}

// API处理器示例
func (h *BotHandler) CreateBot(ctx context.Context, req *CreateBotRequest) (*BotResponse, error) {
    // 参数验证
    if err := validateBotName(req.Name); err != nil {
        return nil, NewBusinessError(ErrCodeInvalidParams, "Invalid bot name").
            AddFieldError("name", err.Error())
    }

    // 业务逻辑
    bot, err := h.botService.Create(ctx, req)
    if err != nil {
        if errors.Is(err, ErrDuplicateBot) {
            return nil, NewBusinessError(ErrCodeResourceConflict, "Bot name already exists")
        }
        return nil, err
    }

    return &BotResponse{Data: bot}, nil
}
```

---

## 9. 认证与授权

### 9.1 认证机制

#### JWT Token认证

```yaml
认证流程:
  1. 客户端使用用户名密码登录
  2. 服务端验证成功后返回JWT Token
  3. 客户端后续请求在Header中携带Token
  4. 服务端验证Token有效性

Token结构:
  access_token:
    expire: 2小时
    contains: user_id, permissions, tenant_id

  refresh_token:
    expire: 30天
    purpose: 刷新access_token
```

#### 请求示例

```http
# 登录获取Token
POST /api/v1/auth/login
Content-Type: application/json

{
  "data": {
    "email": "user@example.com",
    "password": "hashed_password"
  }
}

# 响应
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 7200,
    "token_type": "Bearer",
    "user": {
      "id": 123,
      "name": "John Doe",
      "email": "user@example.com"
    }
  }
}

# 使用Token访问API
GET /api/v1/bots
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

#### Token刷新

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "data": {
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}

# 响应
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "new_access_token",
    "expires_in": 7200
  }
}
```

### 9.2 授权机制

#### 基于角色的访问控制(RBAC)

```yaml
权限模型:
  roles:
    - admin:        # 管理员,全部权限
    - developer:    # 开发者,Bot读写权限
    - user:         # 普通用户,Bot使用权限

  permissions:
    bot:read:       # 查看Bot
    bot:write:      # 创建/编辑Bot
    bot:delete:     # 删除Bot
    bot:publish:    # 发布Bot
    billing:read:   # 查看账单
    billing:write:  # 管理订阅
```

#### 权限检查流程

```go
// 权限检查中间件
func AuthMiddleware(requiredPermission string) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 1. 从Token解析用户信息
        user := getUserFromToken(c)

        // 2. 检查租户隔离
        tenantID := c.GetHeader("X-Tenant-ID")
        if user.TenantID != tenantID {
            c.JSON(403, gin.H{"code": 40303, "message": "Cross-tenant access denied"})
            c.Abort()
            return
        }

        // 3. 检查权限
        if !user.HasPermission(requiredPermission) {
            c.JSON(403, gin.H{
                "code": 40301,
                "message": "Permission denied",
                "data": {
                    "required_permission": requiredPermission,
                    "user_permissions": user.Permissions,
                },
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// 使用示例
botGroup := r.Group("/api/v1/bots")
botGroup.GET("", AuthMiddleware("bot:read"), h.ListBots)
botGroup.POST("", AuthMiddleware("bot:write"), h.CreateBot)
botGroup.DELETE("/:id", AuthMiddleware("bot:delete"), h.DeleteBot)
```

### 9.3 API密钥认证

#### 适用场景
- 服务间调用
- 第三方集成
- Webhook验证

#### API Key格式

```http
# 在Header中传递
X-API-Key: sk_live_xxxxxxxxxxxx

# 或在Query参数传递(不推荐,仅用于Webhook)
?api_key=sk_live_xxxxxxxxxxxx
```

---

## 10. API版本控制

### 10.1 版本控制策略

```yaml
版本规则:
  - URL路径版本控制: /api/v1/, /api/v2/
  - 主版本号递增: v1 -> v2
  - 向后兼容: v1客户端继续可用
  - 废弃通知: 提前6个月通知v1下线
```

### 10.2 版本演进示例

```bash
# v1 API
POST /api/v1/bots
{
  "data": {
    "name": "My Bot",
    "config": {...}  # 单一配置对象
  }
}

# v2 API (增强功能,保持兼容)
POST /api/v2/bots
{
  "data": {
    "name": "My Bot",
    "config": {...},
    "capabilities": {    # 新增字段
      "multimodal": true,
      "voice": true
    }
  }
}
```

### 10.3 版本废弃策略

```http
# 废弃通知响应头
HTTP/1.1 200 OK
X-API-Deprecated: true
X-API-Sunset: 2026-06-30
X-API-Alternative: /api/v2/bots
```

---

## 11. 分页规范

### 11.1 分页参数

```bash
# 标准分页
GET /api/v1/bots?page=1&page_size=20

# 响应
{
  "data": [...],
  "meta": {
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_count": 100,
      "total_pages": 5,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

### 11.2 游标分页(Cursor Pagination)

```bash
# 适用于大数据集
GET /api/v1/messages?cursor=eyJpZCI6MTIzfQ&limit=50

# 响应
{
  "data": [...],
  "meta": {
    "pagination": {
      "next_cursor": "eyJpZCI6MTMzfQ",
      "has_next": true,
      "limit": 50
    }
  }
}
```

### 11.3 分页参数规范

| 参数 | 类型 | 默认值 | 最大值 | 说明 |
|------|------|--------|--------|------|
| page | Integer | 1 | - | 页码(从1开始) |
| page_size | Integer | 20 | 100 | 每页数量 |
| cursor | String | - | - | 游标(游标分页) |
| limit | Integer | 20 | 100 | 数量(游标分页) |

---

## 12. 排序与过滤

### 12.1 排序参数

```bash
# 单字段排序
GET /api/v1/bots?sort_by=created_at&order=desc

# 多字段排序
GET /api/v1/bots?sort=created_at:desc,name:asc

# 支持的排序字段
sort_by: created_at, updated_at, name, message_count
order: asc, desc
```

### 12.2 过滤参数

```bash
# 等值过滤
GET /api/v1/bots?status=published

# 多值过滤
GET /api/v1/bots?status=published,draft

# 范围过滤
GET /api/v1/bots?created_atgte=2025-01-01&created_atlte=2025-12-31

# 搜索过滤
GET /api/v1/bots?q=customer+service

# 复杂过滤
GET /api/v1/bots?status=published&category=assistant&created_atgte=2025-01-01
```

### 12.3 字段投影

```bash
# 仅返回指定字段
GET /api/v1/bots?fields=id,name,status

# 排除指定字段
GET /api/v1/bots?exclude=config,logs
```

---

## 13. 限流与节流

### 13.1 限流策略

```yaml
限流规则:
  认证用户:
    default: 100 requests/minute
    write_ops: 20 requests/minute

  未认证用户:
    default: 20 requests/minute

  API Key:
    tier_1: 1000 requests/minute
    tier_2: 10000 requests/minute

限流算法:
  - 令牌桶算法(Token Bucket)
  - 滑动窗口算法(Sliding Window)
```

### 13.2 限流响应

```http
HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1703925600
Retry-After: 45

{
  "code": 42901,
  "message": "Rate limit exceeded",
  "data": {
    "limit": 100,
    "remaining": 0,
    "reset_at": "2025-12-30T10:01:00Z",
    "retry_after": 45
  }
}
```

---

## 14. 幂等性设计

### 14.1 幂等性要求

```yaml
天然幂等:
  - GET: 查询操作天然幂等
  - PUT: 完整更新幂等(多次结果相同)
  - DELETE: 删除幂等(删除已删除资源返回404)

需要实现幂等:
  - POST: 创建操作(需要幂等键)
  - PATCH: 部分更新(需要条件更新)
```

### 14.2 幂等键实现

```http
# 客户端生成幂等键
POST /api/v1/bots
X-Idempotency-Key: uuid-v4-unique-key
Content-Type: application/json

{
  "data": {
    "name": "My Bot"
  }
}

# 服务端处理
# 1. 检查幂等键是否已处理
# 2. 如果已处理,返回之前的结果
# 3. 如果未处理,处理请求并缓存结果(24小时)
```

---

## 15. API文档规范

### 15.1 OpenAPI规范

```yaml
# 使用OpenAPI 3.0规范编写API文档
openapi: 3.0.0
info:
  title: Coze Studio API
  version: 1.0.0
  description: 一站式AI Agent开发平台API

servers:
  - url: https://api.coze.com/api/v1
    description: Production
  - url: https://api-staging.coze.com/api/v1
    description: Staging

paths:
  /bots:
    get:
      summary: List bots
      operationId: listBots
      tags:
        - Bots
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            default: 1
        - name: page_size
          in: query
          schema:
            type: integer
            default: 20
            maximum: 100
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/BotListResponse'
        '401':
          $ref: '#/components/responses/Unauthorized'

    post:
      summary: Create bot
      operationId: createBot
      tags:
        - Bots
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateBotRequest'
      responses:
        '201':
          description: Bot created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/BotResponse'

components:
  schemas:
    Bot:
      type: object
      required:
        - id
        - name
        - status
      properties:
        id:
          type: integer
          format: int64
        name:
          type: string
          minLength: 2
          maxLength: 50
        description:
          type: string
          maxLength: 500
        status:
          type: string
          enum: [draft, published, archived]
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

  responses:
    Unauthorized:
      description: Authentication required
      content:
        application/json:
          schema:
            type: object
            properties:
              code:
                type: integer
                example: 40101
              message:
                type: string
                example: "Authentication token required"
```

### 15.2 文档生成工具

```bash
# 使用Swag生成Go文档
# 安装
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g cmd/api/main.go -o docs

# 访问Swagger UI
# http://localhost:8888/swagger/index.html
```

### 15.3 API注释规范

```go
// CreateBot 创建Bot
// @Summary      创建Bot
// @Description  创建一个新的AI Bot
// @Tags         bots
// @Accept       json
// @Produce      json
// @Param        request body CreateBotRequest true "创建Bot请求"
// @Success      201  {object}  BotResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      422  {object}  ErrorResponse
// @Router       /bots [post]
func (h *BotHandler) CreateBot(ctx context.Context, c *app.RequestContext) {
    // ...
}
```

---

## 附录

### A. 完整API示例

#### 场景: Bot全生命周期管理

```bash
# 1. 创建Bot
POST /api/v1/bots
Authorization: Bearer {token}
Content-Type: application/json

{
  "data": {
    "name": "Customer Service Bot",
    "description": "AI-powered customer assistant",
    "avatar": "https://cdn.example.com/avatar.png",
    "config": {
      "model": "gpt-4",
      "temperature": 0.7,
      "max_tokens": 2000
    }
  }
}

# 响应: 201 Created
{
  "code": 0,
  "message": "Bot created successfully",
  "data": {
    "id": 123,
    "name": "Customer Service Bot",
    "status": "draft",
    "created_at": "2025-12-30T10:00:00Z"
  }
}

# 2. 更新Bot
PATCH /api/v1/bots/123
Authorization: Bearer {token}
Content-Type: application/json

{
  "data": {
    "description": "Updated description"
  }
}

# 响应: 200 OK
{
  "code": 0,
  "message": "Bot updated successfully",
  "data": {
    "id": 123,
    "name": "Customer Service Bot",
    "description": "Updated description",
    "updated_at": "2025-12-30T10:05:00Z"
  }
}

# 3. 发布Bot
POST /api/v1/bots/123/publish
Authorization: Bearer {token}

# 响应: 202 Accepted
{
  "code": 0,
  "message": "Bot publishing started",
  "data": {
    "job_id": "publish-20251230-001",
    "status": "processing",
    "estimated_time": 30
  }
}

# 4. 查询发布状态
GET /api/v1/jobs/publish-20251230-001
Authorization: Bearer {token}

# 响应: 200 OK
{
  "code": 0,
  "message": "success",
  "data": {
    "job_id": "publish-20251230-001",
    "status": "completed",
    "progress": 100,
    "result": {
      "bot_id": 123,
      "published_at": "2025-12-30T10:07:00Z"
    }
  }
}

# 5. 获取Bot列表
GET /api/v1/bots?page=1&page_size=20&status=published&sort_by=created_at&order=desc
Authorization: Bearer {token}

# 响应: 200 OK
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 123,
      "name": "Customer Service Bot",
      "status": "published",
      "created_at": "2025-12-30T10:00:00Z"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_count": 1,
      "total_pages": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}

# 6. 删除Bot
DELETE /api/v1/bots/123
Authorization: Bearer {token}

# 响应: 204 No Content
```

### B. Go后端实现模板

```go
package handler

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
)

// BotHandler Bot处理器
type BotHandler struct {
    botService BotService
}

// CreateBotRequest 创建Bot请求
type CreateBotRequest struct {
    Name        string                 `json:"name" validate:"required,min=2,max=50"`
    Description string                 `json:"description" validate:"max=500"`
    Avatar      string                 `json:"avatar" validate:"url"`
    Config      map[string]interface{} `json:"config"`
}

// BotResponse Bot响应
type BotResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

// Meta 元数据
type Meta struct {
    Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination 分页信息
type Pagination struct {
    Page       int  `json:"page"`
    PageSize   int  `json:"page_size"`
    TotalCount int  `json:"total_count"`
    TotalPages int  `json:"total_pages"`
    HasNext    bool `json:"has_next"`
    HasPrev    bool `json:"has_prev"`
}

// CreateBot 创建Bot
// @Summary      创建Bot
// @Description  创建一个新的AI Bot
// @Tags         bots
// @Accept       json
// @Produce      json
// @Param        request body CreateBotRequest true "创建Bot请求"
// @Success      201  {object}  BotResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      422  {object}  ErrorResponse
// @Router       /bots [post]
func (h *BotHandler) CreateBot(ctx context.Context, c *app.RequestContext) {
    var req CreateBotRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, NewBusinessError(40001, "Invalid request parameters"))
        return
    }

    bot, err := h.botService.Create(ctx, &req)
    if err != nil {
        handleBusinessError(c, err)
        return
    }

    c.JSON(201, BotResponse{
        Code:    0,
        Message: "Bot created successfully",
        Data:    bot,
        Meta: &Meta{
            Pagination: nil,
        },
    })
}

// ListBots 获取Bot列表
func (h *BotHandler) ListBots(ctx context.Context, c *app.RequestContext) {
    // 解析分页参数
    page := c.DefaultQuery("page", "1")
    pageSize := c.DefaultQuery("page_size", "20")

    // 解析过滤参数
    status := c.Query("status")
    sortBy := c.DefaultQuery("sort_by", "created_at")
    order := c.DefaultQuery("order", "desc")

    // 调用服务层
    bots, pagination, err := h.botService.List(ctx, &ListBotQuery{
        Page:     parseInt(page),
        PageSize: parseInt(pageSize),
        Status:   status,
        SortBy:   sortBy,
        Order:    order,
    })

    if err != nil {
        handleBusinessError(c, err)
        return
    }

    c.JSON(200, BotResponse{
        Code:    0,
        Message: "success",
        Data:    bots,
        Meta: &Meta{
            Pagination: pagination,
        },
    })
}

func handleBusinessError(c *app.RequestContext, err error) {
    if businessErr, ok := err.(*BusinessError); ok {
        c.JSON(businessErr.Code, businessErr)
    } else {
        c.JSON(500, NewBusinessError(50001, "Internal server error"))
    }
}
```

### C. 前端TypeScript类型定义

```typescript
// types/api.ts

// 通用响应结构
interface ApiResponse<T = any> {
  code: number;
  message: string;
  data?: T;
  meta?: {
    pagination?: Pagination;
    [key: string]: any;
  };
}

// 错误响应
interface ApiError {
  code: number;
  message: string;
  errors?: FieldError[];
  data?: Record<string, any>;
  request_id: string;
  timestamp: string;
}

interface FieldError {
  field: string;
  message: string;
  rejected_value?: any;
  constraint?: Record<string, any>;
}

// 分页参数
interface PaginationParams {
  page?: number;
  page_size?: number;
  sort_by?: string;
  order?: 'asc' | 'desc';
}

// 分页信息
interface Pagination {
  page: number;
  page_size: number;
  total_count: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

// Bot相关类型
interface Bot {
  id: number;
  name: string;
  description?: string;
  avatar?: string;
  status: 'draft' | 'published' | 'archived';
  config: Record<string, any>;
  created_at: string;
  updated_at: string;
}

interface CreateBotRequest {
  name: string;
  description?: string;
  avatar?: string;
  config?: Record<string, any>;
}

interface UpdateBotRequest {
  name?: string;
  description?: string;
  avatar?: string;
  config?: Record<string, any>;
}

interface ListBotsParams extends PaginationParams {
  status?: 'draft' | 'published' | 'archived';
  q?: string; // 搜索关键词
}

// API客户端
import axios, { AxiosInstance } from 'axios';

class CozeApiClient {
  private client: AxiosInstance;

  constructor(baseURL: string, token: string) {
    this.client = axios.create({
      baseURL,
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    });

    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          // Token过期,刷新Token
          return this.refreshToken().then(() => {
            return this.client.request(error.config);
          });
        }
        return Promise.reject(error);
      }
    );
  }

  async createBot(request: CreateBotRequest): Promise<Bot> {
    const response = await this.client.post<ApiResponse<Bot>>('/bots', {
      data: request,
    });
    return response.data.data!;
  }

  async listBots(params?: ListBotsParams): Promise<{ bots: Bot[]; pagination: Pagination }> {
    const response = await this.client.get<ApiResponse<Bot[]>>('/bots', { params });
    return {
      bots: response.data.data!,
      pagination: response.data.meta!.pagination!,
    };
  }

  async getBot(id: number): Promise<Bot> {
    const response = await this.client.get<ApiResponse<Bot>>(`/bots/${id}`);
    return response.data.data!;
  }

  async updateBot(id: number, request: UpdateBotRequest): Promise<Bot> {
    const response = await this.client.patch<ApiResponse<Bot>>(`/bots/${id}`, {
      data: request,
    });
    return response.data.data!;
  }

  async deleteBot(id: number): Promise<void> {
    await this.client.delete(`/bots/${id}`);
  }

  private async refreshToken(): Promise<void> {
    // Token刷新逻辑
  }
}

export const apiClient = new CozeApiClient(
  process.env.API_BASE_URL || 'https://api.coze.com/api/v1',
  localStorage.getItem('access_token') || ''
);
```

---

**文档版本历史**:
- v1.0.0 (2025-12-30): 初始版本,定义完整的API接口设计规范
