# API 设计规范文档

> **版本**: v1.0.0
> **更新日期**: 2025-12-30
> **适用范围**: Coze Studio 企业级功能全部API接口

## 目录

- [1. 概述](#1-概述)
- [2. RESTful API 设计原则](#2-restful-api-设计原则)
- [3. 统一响应结构](#3-统一响应结构)
- [4. 认证与授权](#4-认证与授权)
- [5. API 版本控制](#5-api-版本控制)
- [6. 分页标准](#6-分页标准)
- [7. 排序与过滤](#7-排序与过滤)
- [8. 速率限制](#8-速率限制)
- [9. 错误处理](#9-错误处理)
- [10. API 文档规范](#10-api-文档规范)
- [11. 接口设计示例](#11-接口设计示例)
- [12. OpenAPI 规范](#12-openapi-规范)
- [13. 检查清单](#13-检查清单)

---

## 1. 概述

本文档定义了 Coze Studio 企业级功能的统一 API 设计规范,确保所有模块接口的一致性、可维护性和可扩展性。

### 1.1 设计原则

- **统一性**: 所有接口遵循统一的命名、结构和错误处理规范
- **RESTful**: 遵循 REST 架构风格,合理使用 HTTP 方法和状态码
- **版本化**: 所有 API 必须版本化,确保向后兼容性
- **安全性**: 统一的认证授权机制,数据加密传输
- **可观测性**: 统一的日志、追踪和监控标准
- **多租户**: 支持多租户隔离,tenant_id 贯穿所有请求

### 1.2 技术栈

- **后端框架**: CloudWeGo Hertz (Go)
- **API 文档**: OpenAPI 3.0 Specification
- **认证方式**: JWT Bearer Token + API Key
- **数据格式**: JSON (Content-Type: application/json)
- **字符编码**: UTF-8

---

## 2. RESTful API 设计原则

### 2.1 URI 命名规范

#### 基本规则

```
✅ 正确示例:
GET    /api/v1/bots                    # 获取 Bot 列表
GET    /api/v1/bots/{bot_id}           # 获取特定 Bot
POST   /api/v1/bots                    # 创建 Bot
PUT    /api/v1/bots/{bot_id}           # 完整更新 Bot
PATCH  /api/v1/bots/{bot_id}           # 部分更新 Bot
DELETE /api/v1/bots/{bot_id}           # 删除 Bot
POST   /api/v1/bots/{bot_id}/publish   # 发布 Bot (子资源)

❌ 错误示例:
GET    /api/v1/getBots                 # 不要在 URI 中包含动词
GET    /api/v1/bot                     # 资源名应使用复数
POST   /api/v1/bot/create              # 动作应通过 HTTP 方法表达
GET    /api/v1/bots/{bot_id}/delete    # 删除应使用 DELETE 方法
```

#### 命名约定

| 类型 | 规范 | 示例 |
|------|------|------|
| 资源名 | 小写,复数形式,单词间用 `-` 分隔 | `/api/v1/knowledge-bases` |
| 路径参数 | 小写,单词间用 `_` 分隔 | `/{bot_id}`, `/{knowledge_base_id}` |
| 查询参数 | 小写,单词间用 `_` 分隔 | `?page_size=20&sort_by=created_at` |
| 特殊标识 | 使用 UUID 或雪花 ID | `bot_id: "123e4567-e89b-12d3-a456-426614174000"` |

#### 层级结构

```
# 标准资源层级
/api/v1/{resource}                    # 资源集合
/api/v1/{resource}/{resource_id}      # 特定资源
/api/v1/{resource}/{resource_id}/{sub_resource}  # 子资源集合

# 多层级示例
/api/v1/bots/{bot_id}/versions/{version_id}/deployments

# 避免过深层级 (超过 3 层考虑使用查询参数)
❌ /api/v1/tenants/{tenant_id}/users/{user_id}/bots/{bot_id}/versions/{version_id}
✅ /api/v1/bots/{version_id}?tenant_id={tenant_id}&user_id={user_id}
```

### 2.2 HTTP 方法使用规范

| 方法 | 用途 | 幂等性 | 安全性 | 请求体 | 示例 |
|------|------|--------|--------|--------|------|
| GET | 查询资源 | ✅ 是 | ✅ 是 | ❌ 否 | `GET /api/v1/bots` |
| POST | 创建资源 | ❌ 否 | ❌ 否 | ✅ 是 | `POST /api/v1/bots` |
| PUT | 完整更新资源 | ✅ 是 | ❌ 否 | ✅ 是 | `PUT /api/v1/bots/{bot_id}` |
| PATCH | 部分更新资源 | ❌ 否 | ❌ 否 | ✅ 是 | `PATCH /api/v1/bots/{bot_id}` |
| DELETE | 删除资源 | ✅ 是 | ❌ 否 | ❌ 否 | `DELETE /api/v1/bots/{bot_id}` |

#### 方法使用决策树

```
需要操作数据?
├─ 是,创建新资源 → POST
├─ 是,完整更新现有资源 → PUT
├─ 是,部分更新现有资源 → PATCH
├─ 是,删除资源 → DELETE
└─ 否,仅查询 → GET
```

#### 特殊场景处理

```go
// 批量操作使用 POST
POST /api/v1/bots/batch/delete
Body: { "bot_ids": ["id1", "id2", "id3"] }

// 复杂查询使用 POST (当查询参数过长时)
POST /api/v1/bots/search
Body: {
  "filters": { "status": "published" },
  "sort": { "field": "created_at", "order": "desc" }
}

// 异步操作返回 202 Accepted
POST /api/v1/bots/{bot_id}/train
→ 202 Accepted
→ Response: { "task_id": "task_123", "status": "processing" }
```

### 2.3 HTTP 状态码规范

#### 成功响应 (2xx)

| 状态码 | 含义 | 使用场景 | 示例 |
|--------|------|----------|------|
| 200 OK | 请求成功 | GET/PUT/PATCH 成功 | 获取 Bot 详情成功 |
| 201 Created | 创建成功 | POST 创建资源成功 | 创建 Bot 成功 |
| 202 Accepted | 已接受 | 异步任务已提交 | Bot 训练任务已提交 |
| 204 No Content | 成功无返回 | DELETE 成功 | 删除 Bot 成功 |

#### 客户端错误 (4xx)

| 状态码 | 含义 | 使用场景 | 错误代码 |
|--------|------|----------|----------|
| 400 Bad Request | 请求参数错误 | 参数校验失败 | `INVALID_PARAMETER` |
| 401 Unauthorized | 未认证 | 缺少或无效的 Token | `UNAUTHORIZED` |
| 403 Forbidden | 无权限 | Token 有效但权限不足 | `FORBIDDEN` |
| 404 Not Found | 资源不存在 | 资源 ID 不存在 | `RESOURCE_NOT_FOUND` |
| 409 Conflict | 资源冲突 | 资源已存在 | `RESOURCE_ALREADY_EXISTS` |
| 422 Unprocessable Entity | 语义错误 | 业务逻辑校验失败 | `VALIDATION_ERROR` |
| 429 Too Many Requests | 超出限流 | 超出速率限制 | `RATE_LIMIT_EXCEEDED` |

#### 服务器错误 (5xx)

| 状态码 | 含义 | 使用场景 | 错误代码 |
|--------|------|----------|----------|
| 500 Internal Server Error | 服务器内部错误 | 未预期的服务器错误 | `INTERNAL_ERROR` |
| 502 Bad Gateway | 网关错误 | 上游服务错误 | `BAD_GATEWAY` |
| 503 Service Unavailable | 服务不可用 | 服务维护或过载 | `SERVICE_UNAVAILABLE` |
| 504 Gateway Timeout | 网关超时 | 上游服务超时 | `GATEWAY_TIMEOUT` |

---

## 3. 统一响应结构

### 3.1 标准响应格式

```typescript
// 成功响应
interface SuccessResponse<T> {
  code: number;           // 业务状态码,0 表示成功
  message: string;        // 响应消息
  data: T;               // 响应数据
  timestamp: number;     // 服务器时间戳 (毫秒)
  trace_id: string;      // 请求追踪 ID
}

// 示例
{
  "code": 0,
  "message": "success",
  "data": {
    "bot_id": "bot_123",
    "name": "客服助手",
    "status": "published"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_abc123"
}
```

```typescript
// 错误响应
interface ErrorResponse {
  code: number;           // 错误业务码 (非 0)
  message: string;        // 错误消息
  error: {
    code: string;        // 错误代码 (如: INVALID_PARAMETER)
    details: any;        // 错误详情
    stack_trace?: string; // 堆栈跟踪 (仅开发环境)
    request_id: string;  // 请求 ID
  };
  timestamp: number;
  trace_id: string;
}

// 示例
{
  "code": 40001,
  "message": "参数校验失败",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": {
      "field": "bot_name",
      "reason": "Bot 名称长度必须在 2-50 个字符之间"
    },
    "request_id": "req_456"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_abc123"
}
```

### 3.2 分页响应格式

```typescript
interface PaginatedResponse<T> {
  code: number;
  message: string;
  data: {
    items: T[];          // 数据列表
    pagination: {
      page: number;      // 当前页码 (从 1 开始)
      page_size: number; // 每页大小
      total: number;     // 总记录数
      total_pages: number; // 总页数
      has_next: boolean; // 是否有下一页
      has_prev: boolean; // 是否有上一页
    };
  };
  timestamp: number;
  trace_id: string;
}

// 示例
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      { "bot_id": "bot_1", "name": "Bot 1" },
      { "bot_id": "bot_2", "name": "Bot 2" }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5,
      "has_next": true,
      "has_prev": false
    }
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_123"
}
```

### 3.3 批量操作响应格式

```typescript
interface BatchResponse<T> {
  code: number;
  message: string;
  data: {
    succeeded: T[];      // 成功的项目
    failed: Array<{
      item: T;           // 失败的项目
      error: {
        code: string;
        message: string;
      };
    }>;
    summary: {
      total: number;     // 总数
      succeeded_count: number; // 成功数
      failed_count: number;    // 失败数
    };
  };
  timestamp: number;
  trace_id: string;
}
```

### 3.4 Go 后端实现示例

```go
// pkg/response/response.go
package response

import (
    "context"
    "time"
    "github.com/cloudwego/hertz/pkg/app"
)

type Response struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    Timestamp int64       `json:"timestamp"`
    TraceID   string      `json:"trace_id"`
}

type ErrorDetail struct {
    Code       string `json:"code"`
    Details    interface{} `json:"details,omitempty"`
    RequestID  string `json:"request_id"`
}

type ErrorResponse struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Error     ErrorDetail `json:"error"`
    Timestamp int64       `json:"timestamp"`
    TraceID   string      `json:"trace_id"`
}

// 成功响应
func Success(c *app.RequestContext, data interface{}) {
    traceID := c.GetString("trace_id")
    c.JSON(200, &Response{
        Code:      0,
        Message:   "success",
        Data:      data,
        Timestamp: time.Now().UnixMilli(),
        TraceID:   traceID,
    })
}

// 分页响应
func Paginated(c *app.RequestContext, items interface{}, pagination *PaginationInfo) {
    Success(c, &struct {
        Items      interface{}   `json:"items"`
        Pagination *PaginationInfo `json:"pagination"`
    }{
        Items:      items,
        Pagination: pagination,
    })
}

type PaginationInfo struct {
    Page       int  `json:"page"`
    PageSize   int  `json:"page_size"`
    Total      int  `json:"total"`
    TotalPages int  `json:"total_pages"`
    HasNext    bool `json:"has_next"`
    HasPrev    bool `json:"has_prev"`
}

// 错误响应
func Error(c *app.RequestContext, httpStatus int, errCode int, errCodeStr string, message string, details interface{}) {
    traceID := c.GetString("trace_id")
    c.JSON(httpStatus, &ErrorResponse{
        Code:    errCode,
        Message: message,
        Error: ErrorDetail{
            Code:      errCodeStr,
            Details:   details,
            RequestID: traceID,
        },
        Timestamp: time.Now().UnixMilli(),
        TraceID:   traceID,
    })
}
```

---

## 4. 认证与授权

### 4.1 认证机制

#### JWT Bearer Token (用户认证)

```http
GET /api/v1/bots
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Token 结构**:
```typescript
interface JWTPayload {
  sub: string;      // 用户 ID (subject)
  iss: string;      // 签发者 (issuer)
  aud: string;      // 受众 (audience)
  exp: number;      // 过期时间 (expiration)
  iat: number;      // 签发时间 (issued at)
  jti: string;      // Token ID (JWT ID)
  tenant_id: string; // 租户 ID
  roles: string[];  // 用户角色
}
```

#### API Key (应用认证)

```http
GET /api/v1/bots
X-API-Key: sk_live_1234567890abcdef
X-API-Secret: [HMAC-SHA256签名]
```

**API Key 格式**:
```
sk_live_[32位随机字符]  # 生产环境
sk_test_[32位随机字符]  # 测试环境
```

### 4.2 Token 刷新机制

```go
// 刷新 Token
POST /api/v1/auth/refresh
Request:
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

Response:
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600
  }
}
```

### 4.3 授权模型

#### RBAC (基于角色的访问控制)

```typescript
// 角色定义
enum Role {
  ADMIN = 'admin',           // 系统管理员
  TENANT_ADMIN = 'tenant_admin', // 租户管理员
  BOT_DEVELOPER = 'bot_developer', // Bot 开发者
  BOT_VIEWER = 'bot_viewer',     // Bot 查看者
  BOT_OPERATOR = 'bot_operator'   // Bot 运营者
}

// 权限定义
interface Permission {
  resource: string;    // 资源类型 (bot, knowledge_base, etc.)
  action: string;      // 操作 (create, read, update, delete, publish)
  condition?: object;  // 条件限制 (如: 仅限自己创建的资源)
}

// 权限检查中间件示例
func RequirePermission(resource string, action string) app.HandlerFunc {
  return func(c context.Context, ctx *app.RequestContext) {
      user := GetCurrentUser(ctx)
      if !user.HasPermission(resource, action) {
          response.Error(ctx, 403, 40301, "FORBIDDEN", "权限不足", nil)
          ctx.Abort()
          return
      }
      ctx.Next(c)
  }
}
```

#### 多租户隔离

```go
// 所有请求必须包含 tenant_id
// 方式 1: 从 JWT Token 中提取 (推荐)
func TenantFromToken() app.HandlerFunc {
  return func(c context.Context, ctx *app.RequestContext) {
      claims := GetJWTClaims(ctx)
      tenantID := claims.TenantID
      ctx.Set("tenant_id", tenantID)
      ctx.Next(c)
  }
}

// 方式 2: 从 Header 中提取 (跨租户调用)
ctx.GetHeader("X-Tenant-ID") // string

// 方式 3: 从查询参数中提取 (仅限管理接口)
ctx.Query("tenant_id") // string

// 数据库查询自动添加租户过滤
db.Where("tenant_id = ?", tenantID).Find(&bots)
```

### 4.4 安全最佳实践

```go
// 1. Token 过期时间配置
const (
    AccessTokenExpiration  = 1 * time.Hour  // 访问令牌 1 小时
    RefreshTokenExpiration = 30 * 24 * time.Hour // 刷新令牌 30 天
)

// 2. HTTPS 强制
func RequireHTTPS() app.HandlerFunc {
  return func(c context.Context, ctx *app.RequestContext) {
      if !ctx.IsTLS() {
          response.Error(ctx, 400, 40001, "HTTPS_REQUIRED", "必须使用 HTTPS", nil)
          ctx.Abort()
          return
      }
      ctx.Next(c)
  }
}

// 3. 敏感操作二次验证
// 如: 删除 Bot、修改订阅计划等
type SensitiveOperationRequest struct {
  ConfirmPassword string `json:"confirm_password" binding:"required"`
  TOTPCode       string `json:"totp_code" binding:"required"`
}

// 4. API 限流
func RateLimitByUser() app.HandlerFunc {
  return func(c context.Context, ctx *app.RequestContext) {
      userID := ctx.GetString("user_id")
      key := fmt.Sprintf("rate_limit:user:%s", userID)
      count, _ := redis.Incr(ctx, key)
      if count == 1 {
          redis.Expire(ctx, key, 1*time.Minute)
      }
      if count > 100 { // 每分钟 100 次
          response.Error(ctx, 429, 42901, "RATE_LIMIT_EXCEEDED", "超出速率限制", nil)
          ctx.Abort()
          return
      }
      ctx.Next(c)
  }
}
```

---

## 5. API 版本控制

### 5.1 版本策略

#### URL 版本控制 (推荐)

```
/api/v1/bots    # 当前稳定版本
/api/v2/bots    # 新版本 (不兼容的变更)
```

**版本规则**:
- **主版本号**: 不兼容的 API 变更 (v1 → v2)
- **次版本号**: 向后兼容的功能新增 (不在 URL 中体现)
- **修订号**: 向后兼容的问题修正 (不在 URL 中体现)

#### 版本生命周期

```
v1 (Deprecated) → v2 (Current) → v3 (Beta)
     ↓                    ↓
  6个月后退役          当前推荐版本
```

### 5.2 版本兼容性

#### 向后兼容的变更 (不改变版本号)

```typescript
// ✅ 允许的变更
1. 新增可选字段
{ "name": "Bot", "description": "...", "new_field": "optional" }

2. 新增端点
GET /api/v1/bots/new-endpoint

3. 扩展枚举值
enum Status { PUBLISHED, DRAFT, ARCHIVED, NEW_STATUS }

4. 新增查询参数
GET /api/v1/bots?new_filter=value
```

```typescript
// ❌ 不兼容的变更 (需要升级版本号)
1. 删除或重命名字段
{ "bot_name": "..." } → { "name": "..." }

2. 修改字段类型
{ "count": 123 } → { "count": "123" }

3. 修改必填字段
{ "name": string } → { "name": string, "required_field": string }

4. 修改响应结构
{ "data": {...} } → { "result": {...} }
```

### 5.3 版本弃用流程

```http
# 响应头中包含弃用警告
GET /api/v1/bots
Warning: 299 - "Deprecated API. Use /api/v2/bots instead. Will be removed on 2026-06-30."
Sunset: Sun, 30 Jun 2026 00:00:00 GMT
Link: </api/v2/bots>; rel="successor-version"

# 响应体中包含弃用信息
{
  "code": 0,
  "message": "success",
  "data": {...},
  "meta": {
    "deprecation": {
      "deprecated_since": "2025-01-01",
      "sunset_date": "2026-06-30",
      "migration_guide": "https://docs.coze.studio/api/v2-migration"
    }
  }
}
```

---

## 6. 分页标准

### 6.1 基于页码的分页 (适用于小数据集)

```http
GET /api/v1/bots?page=1&page_size=20

Response:
{
  "data": {
    "items": [...],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

**参数限制**:
- `page`: 默认 1,最小值 1
- `page_size`: 默认 20,最小值 1,最大值 100

### 6.2 基于游标的分页 (适用于大数据集)

```http
GET /api/v1/bots?limit=20&cursor=MzAK

Response:
{
  "data": {
    "items": [...],
    "pagination": {
      "next_cursor": "NDQK",
      "prev_cursor": "MjAK",
      "has_next": true,
      "has_prev": false
    }
  }
}
```

**游标编码**: Base64 编码的 `(last_id, last_updated_at)` 元组

```go
// 游标生成和解析
func EncodeCursor(id string, updatedAt time.Time) string {
  data := fmt.Sprintf("%s|%d", id, updatedAt.UnixNano())
  return base64.StdEncoding.EncodeToString([]byte(data))
}

func DecodeCursor(cursor string) (id string, updatedAt time.Time, err error) {
  data, err := base64.StdEncoding.DecodeString(cursor)
  if err != nil {
    return
  }
  parts := strings.Split(string(data), "|")
  // ... 解析逻辑
}
```

### 6.3 分页实现示例

```go
// pkg/pagination/paginator.go
type PageRequest struct {
  Page      int    `query:"page" default:"1" validate:"min=1"`
  PageSize  int    `query:"page_size" default:"20" validate:"min=1,max=100"`
  SortBy    string `query:"sort_by" default:"created_at"`
  SortOrder string `query:"sort_order" default:"desc" validate:"oneof=asc desc"`
}

type CursorRequest struct {
  Limit  int    `query:"limit" default:"20" validate:"min=1,max=100"`
  Cursor string `query:"cursor"`
}

func Paginate(db *gorm.DB, req *PageRequest, model interface{}) (*PaginationInfo, error) {
  var total int64
  db.Model(model).Count(&total)

  offset := (req.Page - 1) * req.PageSize
  db = db.Offset(offset).Limit(req.PageSize)

  // 排序
  if req.SortBy != "" {
    order := req.SortBy
    if req.SortOrder == "desc" {
      order += " DESC"
    }
    db = db.Order(order)
  }

  totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))
  return &PaginationInfo{
    Page:       req.Page,
    PageSize:   req.PageSize,
    Total:      int(total),
    TotalPages: totalPages,
    HasNext:    req.Page < totalPages,
    HasPrev:    req.Page > 1,
  }, nil
}
```

---

## 7. 排序与过滤

### 7.1 排序参数

```http
# 单字段排序
GET /api/v1/bots?sort_by=created_at&sort_order=desc

# 多字段排序
GET /api/v1/bots?sort=created_at:desc,name:asc
```

**可排序字段白名单**:
```go
var sortableFields = map[string]bool{
  "created_at": true,
  "updated_at": true,
  "name":       true,
  "status":     true,
}

func ValidateSortBy(field string) bool {
  return sortableFields[field]
}
```

### 7.2 过滤参数

```http
# 等值过滤
GET /api/v1/bots?status=published&visibility=public

# 范围过滤
GET /api/v1/bots?created_at[gte]=2024-01-01&created_at[lte]=2024-12-31

# 模糊搜索
GET /api/v1/bots?name[contains]=客服

# 数组包含
GET /api/v1/bots?tags[in]=ai,assistant

# 多条件组合 (AND)
GET /api/v1/bots?status=published&tags[in]=ai&created_at[gte]=2024-01-01
```

### 7.3 过滤实现示例

```go
// pkg/filter/filter.go
type Filter struct {
  Field    string
  Operator string // eq, ne, gt, gte, lt, lte, contains, in
  Value    interface{}
}

func ApplyFilter(db *gorm.DB, filters map[string][]string) *gorm.DB {
  for field, conditions := range filters {
    for _, condition := range conditions {
      // 解析操作符和值
      // created_at[gte]=2024-01-01
      re := regexp.MustCompile(`(\w+)\[(\w+)\]=(.+)`)
      matches := re.FindStringSubmatch(condition)
      if len(matches) == 4 {
        fieldName := matches[1]
        operator := matches[2]
        value := matches[3]

        switch operator {
        case "gte":
          db = db.Where(fieldName+" >= ?", value)
        case "lte":
          db = db.Where(fieldName+" <= ?", value)
        case "contains":
          db = db.Where(fieldName+" LIKE ?", "%"+value+"%")
        case "in":
          values := strings.Split(value, ",")
          db = db.Where(fieldName+" IN ?", values)
        case "eq":
          db = db.Where(fieldName+" = ?", value)
        }
      }
    }
  }
  return db
}
```

---

## 8. 速率限制

### 8.1 限流策略

#### 按用户限流

```go
// 限制: 每用户每分钟 100 次
func RateLimitByUser() app.HandlerFunc {
  return func(c context.Context, ctx *app.RequestContext) {
      userID := ctx.GetString("user_id")
      key := fmt.Sprintf("rate_limit:user:%s", userID)

      count, _ := redis.Incr(ctx, key)
      if count == 1 {
          redis.Expire(ctx, key, 1*time.Minute)
      }

      if count > 100 {
          ctx.SetHeader("X-RateLimit-Limit", "100")
          ctx.SetHeader("X-RateLimit-Remaining", "0")
          ctx.SetHeader("X-RateLimit-Reset", time.Now().Add(1*time.Minute).Unix())
          response.Error(ctx, 429, 42901, "RATE_LIMIT_EXCEEDED", "超出速率限制", nil)
          ctx.Abort()
          return
      }

      ctx.SetHeader("X-RateLimit-Limit", "100")
      ctx.SetHeader("X-RateLimit-Remaining", fmt.Sprintf("%d", 100-count))
      ctx.Next(c)
  }
}
```

#### 按租户限流

```go
// 限制: 每租户每秒 50 次
func RateLimitByTenant() app.HandlerFunc {
  return func(c context.Context, ctx *app.RequestContext) {
      tenantID := ctx.GetString("tenant_id")
      key := fmt.Sprintf("rate_limit:tenant:%s", tenantID)

      // 使用令牌桶算法
      allowed, _ := tokenBucket.Allow(tenantID, 50, time.Second)
      if !allowed {
          response.Error(ctx, 429, 42901, "RATE_LIMIT_EXCEEDED", "超出租户速率限制", nil)
          ctx.Abort()
          return
      }
      ctx.Next(c)
  }
}
```

#### 按接口限流

```go
// 昂贵操作更严格的限制
var rateLimits = map[string]int{
  "/api/v1/bots":           100, // 普通接口
  "/api/v1/bots/search":    20,  // 搜索接口
  "/api/v1/bots/export":    5,   // 导出接口
  "/api/v1/bots/train":     2,   // 训练接口
}
```

### 8.2 限流响应头

```http
HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1735575600
Retry-After: 60

{
  "code": 42901,
  "message": "超出速率限制",
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "details": {
      "limit": 100,
      "window": "1m",
      "retry_after": 60
    }
  }
}
```

---

## 9. 错误处理

### 9.1 错误代码规范

#### 错误代码结构

```
[HTTP 状态码][业务模块码][具体错误码]

示例: 4040101
  404   - HTTP 状态码
  01    - 业务模块 (Bot 管理)
  01    - 具体错误 (Bot 不存在)
```

#### 业务模块代码

```
01 - Bot 管理
02 - 知识库管理
03 - 用户与权限
04 - 订阅与付费
05 - 多渠道发布
06 - 多模态功能
07 - Bot 商店
99 - 系统级错误
```

### 9.2 常见错误代码

```go
// pkg/errors/codes.go
const (
  // 认证授权错误 (401xx)
  ErrCodeUnauthorized       = 40101
  ErrCodeTokenExpired       = 40102
  ErrCodeInvalidAPIKey      = 40103
  ErrCodeInvalidSignature   = 40104

  // 权限错误 (403xx)
  ErrCodeForbidden          = 40301
  ErrCodeInsufficientQuota  = 40302
  ErrCodeTenantInactive     = 40303

  // 资源错误 (404xx)
  ErrCodeNotFound           = 40401
  ErrCodeBotNotFound        = 4040101
  ErrCodeKnowledgeBaseNotFound = 4040201

  // 参数错误 (400xx)
  ErrCodeInvalidParameter   = 40001
  ErrCodeMissingParameter   = 40002
  ErrCodeValidationFailed   = 40003

  // 资源冲突 (409xx)
  ErrCodeConflict           = 40901
  ErrCodeDuplicateResource  = 40902
  ErrCodeResourceLocked     = 40903

  // 速率限制 (429xx)
  ErrCodeRateLimitExceeded  = 42901

  // 服务器错误 (500xx)
  ErrCodeInternalError      = 50001
  ErrCodeDatabaseError      = 50002
  ErrCodeUpstreamError      = 50003
)
```

### 9.3 错误响应示例

```go
// 参数校验错误
{
  "code": 40001,
  "message": "参数校验失败",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": {
      "field_errors": [
        {
          "field": "bot_name",
          "message": "Bot 名称长度必须在 2-50 个字符之间"
        },
        {
          "field": "tags",
          "message": "标签数量不能超过 10 个"
        }
      ]
    },
    "request_id": "req_123"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_abc"
}

// 业务逻辑错误
{
  "code": 40302,
  "message": "订阅配额不足",
  "error": {
    "code": "INSUFFICIENT_QUOTA",
    "details": {
      "current_plan": "free",
      "current_quota": 3,
      "required": 5,
      "upgrade_url": "https://coze.studio/pricing"
    },
    "request_id": "req_456"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_def"
}
```

### 9.4 Go 错误处理最佳实践

```go
// pkg/errors/errors.go
package errors

import "github.com/cloudwego/hertz/pkg/app"

type AppError struct {
  HTTPStatus int    // HTTP 状态码
  Code       int    // 业务错误码
  CodeStr    string // 错误代码字符串
  Message    string // 用户友好的错误消息
  Details    interface{} // 错误详情
}

func (e *AppError) Error() string {
  return e.Message
}

// 预定义错误
var (
  ErrBotNotFound = &AppError{
    HTTPStatus: 404,
    Code:       ErrCodeBotNotFound,
    CodeStr:    "BOT_NOT_FOUND",
    Message:    "Bot 不存在",
  }

  ErrInsufficientQuota = &AppError{
    HTTPStatus: 403,
    Code:       ErrCodeInsufficientQuota,
    CodeStr:    "INSUFFICIENT_QUOTA",
    Message:    "订阅配额不足",
  }
)

// 错误处理中间件
func HandleError() app.HandlerFunc {
  return func(c context.Context, ctx *app.RequestContext) {
    ctx.Next(c)

    // 检查是否有错误
    if len(ctx.Errors) > 0 {
      err := ctx.Errors.Last()

      // 类型断言,获取自定义错误
      if appErr, ok := err.Err.(*AppError); ok {
        response.Error(
          ctx,
          appErr.HTTPStatus,
          appErr.Code,
          appErr.CodeStr,
          appErr.Message,
          appErr.Details,
        )
        return
      }

      // 默认服务器错误
      response.Error(
        ctx,
        500,
        ErrCodeInternalError,
        "INTERNAL_ERROR",
        "服务器内部错误",
        nil,
      )
    }
  }
}
```

---

## 10. API 文档规范

### 10.1 OpenAPI 规范

```yaml
# openapi.yaml
openapi: 3.0.3
info:
  title: Coze Studio API
  version: 1.0.0
  description: |
    Coze Studio 企业级功能 API

    **认证方式**:
    - JWT Bearer Token (用户认证)
    - API Key (应用认证)

    **联系方式**:
    - 技术支持: support@coze.studio
    - 文档: https://docs.coze.studio

  contact:
    name: API Support
    email: support@coze.studio
    url: https://docs.coze.studio
  license:
    name: MIT
    url: https://opensource.org/licenses/MIT

servers:
  - url: https://api.coze.studio/api/v1
    description: 生产环境
  - url: https://staging-api.coze.studio/api/v1
    description: 预发布环境
  - url: http://localhost:8080/api/v1
    description: 开发环境

security:
  - BearerAuth: []
  - APIKeyAuth: []

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: JWT Bearer Token 认证
    APIKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
      description: API Key 认证

  schemas:
    Bot:
      type: object
      required:
        - bot_id
        - name
        - status
      properties:
        bot_id:
          type: string
          format: uuid
          description: Bot ID
        name:
          type: string
          minLength: 2
          maxLength: 50
          description: Bot 名称
        description:
          type: string
          maxLength: 500
          description: Bot 描述
        status:
          type: string
          enum: [draft, published, archived]
          description: Bot 状态
        created_at:
          type: integer
          format: int64
          description: 创建时间戳 (毫秒)
        updated_at:
          type: integer
          format: int64
          description: 更新时间戳 (毫秒)

    SuccessResponse:
      type: object
      properties:
        code:
          type: integer
          description: 业务状态码 (0 表示成功)
        message:
          type: string
          description: 响应消息
        data:
          type: object
          description: 响应数据
        timestamp:
          type: integer
          format: int64
          description: 服务器时间戳 (毫秒)
        trace_id:
          type: string
          description: 请求追踪 ID

    ErrorResponse:
      type: object
      properties:
        code:
          type: integer
          description: 错误业务码
        message:
          type: string
          description: 错误消息
        error:
          type: object
          properties:
            code:
              type: string
              description: 错误代码
            details:
              type: object
              description: 错误详情
            request_id:
              type: string
              description: 请求 ID
        timestamp:
          type: integer
          format: int64
        trace_id:
          type: string

paths:
  /bots:
    get:
      summary: 获取 Bot 列表
      description: |
        获取当前用户有权访问的 Bot 列表,支持分页、排序和过滤。

        **权限要求**: `bot.read`

        **速率限制**: 100 次/分钟/用户
      operationId: listBots
      tags:
        - Bots
      parameters:
        - name: page
          in: query
          description: 页码 (从 1 开始)
          required: false
          schema:
            type: integer
            minimum: 1
            default: 1
        - name: page_size
          in: query
          description: 每页大小
          required: false
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 20
        - name: status
          in: query
          description: 按 Bot 状态过滤
          required: false
          schema:
            type: string
            enum: [draft, published, archived]
        - name: sort_by
          in: query
          description: 排序字段
          required: false
          schema:
            type: string
            enum: [created_at, updated_at, name]
            default: created_at
        - name: sort_order
          in: query
          description: 排序顺序
          required: false
          schema:
            type: string
            enum: [asc, desc]
            default: desc
      responses:
        '200':
          description: 成功
          content:
            application/json:
              schema:
                allOf:
                  - $ref: '#/components/schemas/SuccessResponse'
                  - type: object
                    properties:
                      data:
                        type: object
                        properties:
                          items:
                            type: array
                            items:
                              $ref: '#/components/schemas/Bot'
                          pagination:
                            $ref: '#/components/schemas/Pagination'
        '400':
          description: 请求参数错误
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '401':
          description: 未认证
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '403':
          description: 无权限
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '429':
          description: 超出速率限制
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

    post:
      summary: 创建 Bot
      description: |
        创建新的 Bot。

        **权限要求**: `bot.create`

        **速率限制**: 20 次/分钟/用户
      operationId: createBot
      tags:
        - Bots
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - name
              properties:
                name:
                  type: string
                  minLength: 2
                  maxLength: 50
                  description: Bot 名称
                  example: "客服助手"
                description:
                  type: string
                  maxLength: 500
                  description: Bot 描述
                  example: "智能客服 Bot,用于处理常见问题"
                avatar_url:
                  type: string
                  format: uri
                  description: Bot 头像 URL
                tags:
                  type: array
                  maxItems: 10
                  items:
                    type: string
                  description: Bot 标签
                  example: ["客服", "AI"]
      responses:
        '201':
          description: 创建成功
          content:
            application/json:
              schema:
                allOf:
                  - $ref: '#/components/schemas/SuccessResponse'
                  - type: object
                    properties:
                      data:
                        $ref: '#/components/schemas/Bot'
        '400':
          description: 请求参数错误
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '401':
          description: 未认证
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '403':
          description: 无权限
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /bots/{bot_id}:
    get:
      summary: 获取 Bot 详情
      description: |
        获取指定 Bot 的详细信息。

        **权限要求**: `bot.read`
      operationId: getBot
      tags:
        - Bots
      parameters:
        - name: bot_id
          in: path
          required: true
          description: Bot ID
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: 成功
          content:
            application/json:
              schema:
                allOf:
                  - $ref: '#/components/schemas/SuccessResponse'
                  - type: object
                    properties:
                      data:
                        $ref: '#/components/schemas/Bot'
        '404':
          description: Bot 不存在
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

    put:
      summary: 完整更新 Bot
      description: |
        完整更新 Bot 的所有字段。

        **权限要求**: `bot.update`
      operationId: updateBot
      tags:
        - Bots
      parameters:
        - name: bot_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - name
                - status
              properties:
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
      responses:
        '200':
          description: 更新成功
          content:
            application/json:
              schema:
                allOf:
                  - $ref: '#/components/schemas/SuccessResponse'
                  - type: object
                    properties:
                      data:
                        $ref: '#/components/schemas/Bot'
        '400':
          description: 请求参数错误
        '404':
          description: Bot 不存在

    patch:
      summary: 部分更新 Bot
      description: |
        部分更新 Bot 的指定字段。

        **权限要求**: `bot.update`
      operationId: patchBot
      tags:
        - Bots
      parameters:
        - name: bot_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                  minLength: 2
                  maxLength: 50
                description:
                  type: string
                  maxLength: 500
      responses:
        '200':
          description: 更新成功
        '400':
          description: 请求参数错误
        '404':
          description: Bot 不存在

    delete:
      summary: 删除 Bot
      description: |
        删除指定 Bot。

        **权限要求**: `bot.delete`

        **注意**: 删除操作不可恢复,请谨慎操作。
      operationId: deleteBot
      tags:
        - Bots
      parameters:
        - name: bot_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '204':
          description: 删除成功
        '404':
          description: Bot 不存在
```

### 10.2 接口注释规范

```go
// api/handler/bot_handler.go
package handler

// CreateBot 创建 Bot
//
//	@Summary		创建 Bot
//	@Description	创建新的 Bot,需要 bot.create 权限
//	@Tags			bots
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		request.CreateBotRequest	true	"创建 Bot 请求"
//	@Success		201		{object}	response.SuccessResponse{data=domain.Bot}
//	@Failure		400		{object}	response.ErrorResponse	"请求参数错误"
//	@Failure		401		{object}	response.ErrorResponse	"未认证"
//	@Failure		403		{object}	response.ErrorResponse	"无权限"
//	@Failure		500		{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/bots [post]
func (h *BotHandler) CreateBot(c context.Context, ctx *app.RequestContext) {
	// ...
}
```

---

## 11. 接口设计示例

### 11.1 Bot 管理接口

#### 创建 Bot

```http
POST /api/v1/bots
Authorization: Bearer eyJhbGci...
Content-Type: application/json

{
  "name": "智能客服助手",
  "description": "处理常见客户问题的 AI Bot",
  "avatar_url": "https://cdn.example.com/avatar.png",
  "tags": ["客服", "AI", "自动化"],
  "config": {
    "model": "gpt-4",
    "temperature": 0.7,
    "max_tokens": 2000
  }
}

Response 201 Created:
{
  "code": 0,
  "message": "success",
  "data": {
    "bot_id": "bot_1234567890",
    "name": "智能客服助手",
    "description": "处理常见客户问题的 AI Bot",
    "avatar_url": "https://cdn.example.com/avatar.png",
    "tags": ["客服", "AI", "自动化"],
    "status": "draft",
    "config": {
      "model": "gpt-4",
      "temperature": 0.7,
      "max_tokens": 2000
    },
    "created_by": "user_123",
    "created_at": 1735574400000,
    "updated_at": 1735574400000
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_abc123"
}
```

#### 获取 Bot 列表

```http
GET /api/v1/bots?page=1&page_size=20&status=published&sort_by=created_at&sort_order=desc
Authorization: Bearer eyJhbGci...

Response 200 OK:
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "bot_id": "bot_1",
        "name": "Bot 1",
        "status": "published",
        "created_at": 1735574400000
      },
      {
        "bot_id": "bot_2",
        "name": "Bot 2",
        "status": "published",
        "created_at": 1735574300000
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5,
      "has_next": true,
      "has_prev": false
    }
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_def456"
}
```

#### 发布 Bot

```http
POST /api/v1/bots/bot_123/publish
Authorization: Bearer eyJhbGci...

{
  "channels": ["web", "wechat"],
  "version_note": "V1.0.0 初始版本"
}

Response 202 Accepted:
{
  "code": 0,
  "message": "发布任务已提交",
  "data": {
    "task_id": "task_789",
    "status": "processing",
    "estimated_time": 30
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_ghi789"
}
```

### 11.2 订阅管理接口

#### 创建订阅

```http
POST /api/v1/subscriptions
Authorization: Bearer eyJhbGci...

{
  "plan_id": "plan_pro",
  "billing_cycle": "monthly",
  "payment_method_id": "pm_123456"
}

Response 201 Created:
{
  "code": 0,
  "message": "success",
  "data": {
    "subscription_id": "sub_123456",
    "plan_id": "plan_pro",
    "status": "active",
    "current_period_start": 1735574400000,
    "current_period_end": 1738166400000,
    "cancel_at_period_end": false
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_jkl012"
}
```

#### 升级/降级订阅

```http
POST /api/v1/subscriptions/sub_123456/change-plan
Authorization: Bearer eyJhbGci...

{
  "new_plan_id": "plan_enterprise",
  "effective_date": "immediate"
}

Response 200 OK:
{
  "code": 0,
  "message": "订阅计划已更新",
  "data": {
    "subscription_id": "sub_123456",
    "previous_plan": "plan_pro",
    "new_plan": "plan_enterprise",
    "status": "active",
    "effective_at": 1735574400000
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_mno345"
}
```

### 11.3 多渠道发布接口

#### 发布到指定渠道

```http
POST /api/v1/bots/bot_123/publish
Authorization: Bearer eyJhbGci...

{
  "channels": [
    {
      "type": "web",
      "config": {
        "theme": "dark",
        "position": "bottom-right"
      }
    },
    {
      "type": "wechat",
      "config": {
        "app_id": "wx123456",
        "auto_reply": true
      }
    }
  ]
}

Response 200 OK:
{
  "code": 0,
  "message": "发布成功",
  "data": {
    "deployment_id": "deploy_789",
    "channels": [
      {
        "type": "web",
        "status": "deployed",
        "url": "https://bot.coze.studio/bot_123"
      },
      {
        "type": "wechat",
        "status": "deployed",
        "app_id": "wx123456"
      }
    ],
    "deployed_at": 1735574400000
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_pqr678"
}
```

---

## 12. OpenAPI 规范

### 12.1 完整的 OpenAPI YAML 模板

参见上述第 10.1 节中的完整 OpenAPI 规范示例。

### 12.2 自动生成 API 文档

```bash
# 使用 swag 工具从代码注释生成 Swagger 文档
go install github.com/swaggo/swag/cmd/swag@latest

# 在项目根目录执行
swag init -g cmd/server.go -o docs

# 访问生成的文档
# http://localhost:8080/swagger/index.html
```

### 12.3 API 版本文档管理

```
docs/
├── api/
│   ├── openapi-v1.yaml     # v1 版本 API 规范
│   ├── openapi-v2.yaml     # v2 版本 API 规范
│   └── migration-v1-to-v2.md # v1 到 v2 迁移指南
└── examples/
    ├── bot-management/
    │   ├── create-bot.http
    │   ├── list-bots.http
    │   └── publish-bot.http
    └── subscription/
        ├── create-subscription.http
        └── upgrade-plan.http
```

---

## 13. 检查清单

### 13.1 API 设计检查清单

#### RESTful 设计
- [ ] URI 使用名词复数形式
- [ ] 正确使用 HTTP 方法 (GET/POST/PUT/PATCH/DELETE)
- [ ] 资源层级不超过 3 层
- [ ] URI 中的参数使用 kebab-case
- [ ] 查询参数使用 snake_case

#### 响应结构
- [ ] 使用统一的响应结构
- [ ] 成功响应包含 code=0
- [ ] 错误响应包含详细的错误信息
- [ ] 分页响应包含完整的分页信息
- [ ] 所有响应包含 timestamp 和 trace_id

#### 状态码
- [ ] 正确使用 HTTP 状态码
- [ ] 200 OK: GET/PUT/PATCH 成功
- [ ] 201 Created: POST 创建成功
- [ ] 202 Accepted: 异步任务已提交
- [ ] 204 No Content: DELETE 成功
- [ ] 400: 参数错误
- [ ] 401: 未认证
- [ ] 403: 无权限
- [ ] 404: 资源不存在
- [ ] 409: 资源冲突
- [ ] 429: 速率限制
- [ ] 500: 服务器错误

#### 认证授权
- [ ] 所有 API (除登录外) 需要认证
- [ ] Token 包含必要的用户信息 (sub, tenant_id, roles)
- [ ] 敏感操作需要权限检查
- [ ] 跨租户数据访问需要额外授权
- [ ] 所有数据库查询自动添加 tenant_id 过滤

#### 分页排序过滤
- [ ] 列表接口支持分页
- [ ] page_size 有合理的最大值限制
- [ ] 支持排序参数 (sort_by, sort_order)
- [ ] 排序字段有白名单限制
- [ ] 支持常见的过滤条件
- [ ] 过滤字段有白名单限制

#### 错误处理
- [ ] 使用预定义的错误代码
- [ ] 错误消息对用户友好
- [ ] 开发环境包含堆栈跟踪
- [ ] 生产环境隐藏敏感信息
- [ ] 所有错误包含 trace_id

#### 速率限制
- [ ] 按用户限流
- [ ] 按租户限流
- [ ] 昂贵操作有更严格的限制
- [ ] 限流响应包含 Retry-After 头

#### 文档
- [ ] 有完整的 OpenAPI 规范
- [ ] 代码注释包含 Swagger 注释
- [ ] 示例请求和响应完整
- [ ] 有版本变更日志
- [ ] 有迁移指南 (跨版本)

#### 多租户
- [ ] 所有请求包含 tenant_id
- [ ] 数据库查询自动添加租户过滤
- [ ] 租户间数据完全隔离
- [ ] 跨租户操作有审计日志

### 13.2 代码实现检查清单

#### Go 后端
- [ ] 遵循 DDD 分层架构 (api/application/domain/infra)
- [ ] 使用统一响应结构 (`response.Success`, `response.Error`)
- [ ] 使用预定义错误类型 (`errors.ErrBotNotFound`)
- [ ] 所有 handler 方法添加 Swagger 注释
- [ ] 请求参数校验 (`binding:"required"`)
- [ ] 数据库事务处理正确
- [ ] 敏感信息不记录到日志

#### 安全
- [ ] SQL 注入防护 (使用参数化查询)
- [ ] XSS 防护 (输入过滤和输出转义)
- [ ] CSRF 防护 (CSRF Token)
- [ ] 敏感数据加密存储
- [ ] HTTPS 强制
- [ ] 敏感操作二次验证

---

## 附录

### A. HTTP 状态码速查表

| 状态码 | 名称 | 使用场景 |
|--------|------|----------|
| 200 | OK | 请求成功 |
| 201 | Created | 创建成功 |
| 202 | Accepted | 异步任务已接受 |
| 204 | No Content | 删除成功 |
| 400 | Bad Request | 参数错误 |
| 401 | Unauthorized | 未认证 |
| 403 | Forbidden | 无权限 |
| 404 | Not Found | 资源不存在 |
| 409 | Conflict | 资源冲突 |
| 422 | Unprocessable Entity | 业务逻辑错误 |
| 429 | Too Many Requests | 速率限制 |
| 500 | Internal Server Error | 服务器错误 |
| 502 | Bad Gateway | 网关错误 |
| 503 | Service Unavailable | 服务不可用 |

### B. 错误代码速查表

| 错误代码 | 错误字符串 | HTTP 状态码 | 说明 |
|----------|-----------|-------------|------|
| 40101 | UNAUTHORIZED | 401 | 未认证 |
| 40102 | TOKEN_EXPIRED | 401 | Token 过期 |
| 40301 | FORBIDDEN | 403 | 无权限 |
| 40302 | INSUFFICIENT_QUOTA | 403 | 配额不足 |
| 40401 | RESOURCE_NOT_FOUND | 404 | 资源不存在 |
| 4040101 | BOT_NOT_FOUND | 404 | Bot 不存在 |
| 40001 | INVALID_PARAMETER | 400 | 参数错误 |
| 40002 | VALIDATION_ERROR | 400 | 校验失败 |
| 40901 | RESOURCE_ALREADY_EXISTS | 409 | 资源已存在 |
| 42901 | RATE_LIMIT_EXCEEDED | 429 | 速率限制 |
| 50001 | INTERNAL_ERROR | 500 | 内部错误 |

### C. 参考资源

- [RESTful API 设计指南](https://restfulapi.net/)
- [OpenAPI 3.0 规范](https://swagger.io/specification/)
- [HTTP 状态码完整列表](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Status)
- [CloudWeGo Hertz 文档](https://www.cloudwego.io/docs/hertz/)
- [JWT 最佳实践](https://datatracker.ietf.org/doc/html/rfc8725)

---

**文档维护**: 本文档应随 API 演进持续更新,所有新增或修改的接口都必须同步更新本文档。
