# ZKER API 设计规范

**版本**: v3.0.0 | **更新**: 2025-01-03 | **状态**: 强制执行

---

## 目录

- [RESTful 设计原则](#restful-设计原则)
- [URL 设计规范](#url-设计规范)
- [HTTP 方法规范](#http-方法规范)
- [请求规范](#请求规范)
- [响应规范](#响应规范)
- [错误处理规范](#错误处理规范)
- [API 版本控制](#api-版本控制)
- [分页规范](#分页规范)
- [排序规范](#排序规范)
- [过滤规范](#过滤规范)
- [检查清单](#检查清单)

---

## RESTful 设计原则

### 核心原则

| 原则 | 说明 | 示例 |
|------|------|------|
| **资源导向** | URL 表示资源，不是动作 | `/tenants` ✅ vs `/getTenants` ❌ |
| **HTTP 方法语义** | 使用正确的 HTTP 方法 | `GET` 查询, `POST` 创建, `PUT` 更新, `DELETE` 删除 |
| **无状态** | 每个请求包含所有必要信息 | 不依赖会话状态 |
| **统一接口** | 使用统一的 URL 结构 | `/api/v1/{resource}/{id}` |

### RESTful 成熟度模型

```mermaid
graph LR
    A[Level 0: Swamp of POX] --> B[Level 1: Resources]
    B --> C[Level 2: HTTP Verbs]
    C --> D[Level 3: Hypermedia Controls]

    style D fill:#4caf50
    style C fill:#8bc34a
    style B fill:#cddc39
    style A fill:#ffeb3b
```

**ZKER 目标**: Level 2（使用正确的 HTTP 方法）

---

## URL 设计规范

### 基本格式

```
{scheme}://{host}/{api}/{version}/{resource}/{id}
```

### 规则

| 规则 | 说明 | 示例 |
|------|------|------|
| **全小写** | URL 全部小写 | `/api/v1/tenants` ✅ |
| **使用连字符** | 多个单词使用连字符分隔 | `/bot-templates` ✅ vs `/botTemplates` ❌ |
| **使用复数** | 资源名使用复数形式 | `/tenants` ✅ vs `/tenant` ❌ |
| **使用名词** | 不使用动词 | `/tenants` ✅ vs `/getTenants` ❌ |
| **层级深度 ≤ 3** | 避免过深的层级 | `/tenants` ✅ vs `/a/b/c/d/tenants` ❌ |

### 示例

```go
// ✅ Good: 标准的 RESTful API
GET    /api/v1/tenants                    # 获取租户列表
POST   /api/v1/tenants                    # 创建租户
GET    /api/v1/tenants/:id                # 获取租户详情
PUT    /api/v1/tenants/:id                # 更新租户
DELETE /api/v1/tenants/:id                # 删除租户
GET    /api/v1/tenants/:id/users          # 获取租户的用户列表
POST   /api/v1/tenants/:id/users          # 为租户创建用户

// ❌ Bad: 不符合 RESTful 规范
GET    /api/v1/tenant                     // 要使用复数
POST   /api/v1/createTenant               // 不要在 URL 中使用动词
GET    /api/v1/getTenantById              // 不要在 URL 中使用动词
GET    /api/v1/Tenants                    // 不要大写
GET    /api/v1/tenants/:tenantId/users    // 参数名使用小写 + 下划线
```

---

## HTTP 方法规范

### 标准方法

| 方法 | 说明 | 幂等性 | 安全性 | 示例 |
|------|------|--------|--------|------|
| **GET** | 获取资源 | ✅ 是 | ✅ 是 | `GET /api/v1/tenants` |
| **POST** | 创建资源 | ❌ 否 | ❌ 否 | `POST /api/v1/tenants` |
| **PUT** | 完整更新资源 | ✅ 是 | ❌ 否 | `PUT /api/v1/tenants/:id` |
| **PATCH** | 部分更新资源 | ❌ 否 | ❌ 否 | `PATCH /api/v1/tenants/:id` |
| **DELETE** | 删除资源 | ✅ 是 | ❌ 否 | `DELETE /api/v1/tenants/:id` |

### 方法使用场景

```go
// ✅ Good: 正确使用 HTTP 方法
GET    /api/v1/tenants              # 获取租户列表
POST   /api/v1/tenants              # 创建租户
GET    /api/v1/tenants/:id          # 获取单个租户
PUT    /api/v1/tenants/:id          # 完整更新租户
PATCH  /api/v1/tenants/:id          # 部分更新租户（如只更新状态）
DELETE /api/v1/tenants/:id          # 删除租户

// ❌ Bad: 错误使用 HTTP 方法
POST   /api/v1/tenants/:id/delete   // 应该使用 DELETE
GET    /api/v1/tenants/create       // 应该使用 POST
POST   /api/v1/tenants/:id/update   // 应该使用 PUT/PATCH
GET    /api/v1/tenants/:id/delete   // 不要使用 GET 执行删除操作
```

### 幂等性

**定义**: 同一个请求执行多次，结果相同。

| 方法 | 幂等性 | 说明 |
|------|--------|------|
| **GET** | ✅ 是 | 多次获取资源，结果相同 |
| **POST** | ❌ 否 | 多次创建，会创建多个资源 |
| **PUT** | ✅ 是 | 多次更新，结果相同 |
| **PATCH** | ❌ 否 | 取决于操作 |
| **DELETE** | ✅ 是 | 多次删除，结果相同（已删除返回 404） |

---

## 请求规范

### 请求头

#### 必需的请求头

```http
# 通用请求头
Content-Type: application/json
Accept: application/json
Authorization: Bearer {token}

# 租户上下文（如果未在 JWT 中包含）
X-Tenant-ID: {tenant_id}

# 请求追踪（用于日志关联）
X-Request-ID: {request_id}

# 语言偏好
Accept-Language: zh-CN
```

#### 请求头示例

```bash
# 使用 curl
curl -X GET http://localhost:8888/api/v1/tenants \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "X-Tenant-ID: tenant_123"

# 使用 fetch
fetch('http://localhost:8888/api/v1/tenants', {
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
    'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
    'X-Tenant-ID': 'tenant_123',
  },
});
```

### 请求体

#### 创建资源 (POST)

```json
// ✅ Good: 清晰的请求体结构
POST /api/v1/tenants

{
  "tenant_name": "示例公司",
  "plan_type": "pro",
  "admin": {
    "username": "admin",
    "email": "admin@example.com",
    "password": "password123"
  }
}
```

#### 更新资源 (PUT/PATCH)

```json
// ✅ Good: PUT 完整更新
PUT /api/v1/tenants/{tenant_id}

{
  "tenant_name": "新公司名称",
  "plan_type": "enterprise",
  "status": "active"
}

// ✅ Good: PATCH 部分更新
PATCH /api/v1/tenants/{tenant_id}

{
  "status": "suspended"
}
```

### 参数命名

#### 规则

| 规则 | 说明 | 示例 |
|------|------|------|
| **全小写** | 参数名全小写 | `page_size` ✅ vs `pageSize` ❌ |
| **使用下划线** | 多个单词使用下划线分隔 | `created_at` ✅ vs `createdAt` ❌ |
| **清晰表达** | 参数名要清晰表达意图 | `tenant_id` ✅ vs `id` ❌ |

---

## 响应规范

### 响应格式

#### 成功响应

```json
// ✅ Good: 标准的成功响应
{
  "code": 0,
  "message": "success",
  "data": {
    "tenant_id": "tenant_123",
    "tenant_name": "示例公司",
    "plan_type": "pro",
    "status": "active",
    "created_at": "2025-01-03T10:00:00Z"
  },
  "request_id": "req_abc123"
}

// ✅ Good: 列表响应
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "tenant_id": "tenant_123",
        "tenant_name": "示例公司",
        "plan_type": "pro",
        "status": "active"
      },
      {
        "tenant_id": "tenant_456",
        "tenant_name": "另一家公司",
        "plan_type": "free",
        "status": "active"
      }
    ],
    "total": 2,
    "page": 1,
    "page_size": 20
  },
  "request_id": "req_abc123"
}
```

#### 错误响应

```json
// ✅ Good: 标准的错误响应
{
  "code": 10001,
  "message": "租户不存在",
  "details": {
    "tenant_id": "tenant_123",
    "reason": "not found"
  },
  "request_id": "req_abc123",
  "timestamp": "2025-01-03T10:00:00Z"
}

// ✅ Good: 参数验证错误
{
  "code": 20001,
  "message": "参数验证失败",
  "details": {
    "errors": [
      {
        "field": "tenant_name",
        "message": "租户名称长度必须在 3-100 之间"
      },
      {
        "field": "email",
        "message": "邮箱格式不正确"
      }
    ]
  },
  "request_id": "req_abc123",
  "timestamp": "2025-01-03T10:00:00Z"
}
```

### HTTP 状态码

#### 成功状态码

| 状态码 | 说明 | 使用场景 |
|--------|------|---------|
| **200 OK** | 请求成功 | GET、PATCH、DELETE |
| **201 Created** | 创建成功 | POST |
| **204 No Content** | 成功但无返回内容 | DELETE |

#### 客户端错误状态码

| 状态码 | 说明 | 使用场景 |
|--------|------|---------|
| **400 Bad Request** | 请求参数错误 | 参数验证失败 |
| **401 Unauthorized** | 未认证 | Token 无效或过期 |
| **403 Forbidden** | 无权限 | 权限不足 |
| **404 Not Found** | 资源不存在 | 查询的资源不存在 |
| **409 Conflict** | 资源冲突 | 资源已存在 |
| **422 Unprocessable Entity** | 无法处理 | 业务逻辑验证失败 |
| **429 Too Many Requests** | 请求过多 | 超出限流 |

#### 服务器错误状态码

| 状态码 | 说明 | 使用场景 |
|--------|------|---------|
| **500 Internal Server Error** | 服务器内部错误 | 未预期的服务器错误 |
| **502 Bad Gateway** | 网关错误 | 上游服务错误 |
| **503 Service Unavailable** | 服务不可用 | 服务维护中 |

---

## 错误处理规范

### 统一错误码

#### 错误码结构

```
{业务模块}{错误类型}{具体错误}
```

| 部分 | 位数 | 说明 | 示例 |
|------|------|------|------|
| **业务模块** | 1 位 | 1=租户, 2=权限, 3=路由, 4=Agent | `1` |
| **错误类型** | 1 位 | 0=系统, 1=业务, 2=参数 | `0` |
| **具体错误** | 3 位 | 具体错误编号 | `001` |

#### 错误码示例

```go
// backend/types/errno/tenant.go
package errno

const (
    // 1xxxx: 租户模块

    // 10xxx: 租户系统错误
    TenantNotFound          = 10001  // 租户不存在
    TenantAlreadyExists     = 10002  // 租户已存在
    TenantNameInvalid       = 10003  // 租户名称无效
    TenantSuspended         = 10004  // 租户已停用

    // 12xxx: 租户参数错误
    TenantNameTooShort      = 12001  // 租户名称太短
    TenantNameTooLong       = 12002  // 租户名称太长
    TenantPlanTypeInvalid   = 12003  // 套餐类型无效

    // 11xxx: 租户业务错误
    TenantQuotaExceeded     = 11001  // 配额超限
    TenantExpired           = 11002  // 租户已过期
    TenantInactive          = 11003  // 租户未激活
)

// backend/types/errno/permission.go
package errno

const (
    // 2xxxx: 权限模块

    // 20xxx: 权限系统错误
    PermissionDenied        = 20001  // 权限不足
    RoleNotFound            = 20002  // 角色不存在
    RoleAlreadyExists       = 20003  // 角色已存在

    // 21xxx: 权限业务错误
    UserNotInRole           = 21001  // 用户不在角色中
    RoleHasUsers            = 21002  // 角色下有用户，无法删除
)
```

### 错误处理示例

```go
// ✅ Good: 使用统一错误码
package handler

import (
    "backend/types/errno"
    "github.com/cloudwego/hertz/pkg/app"
)

type TenantHandler struct {
    tenantApp *application.TenantApplication
}

func (h *TenantHandler) GetTenant(ctx context.Context, c *app.RequestContext) {
    tenantID := c.Param("id")

    tenant, err := h.tenantApp.GetTenantByID(ctx, tenantID)
    if err != nil {
        // 使用统一错误码
        if errors.Is(err, errno.TenantNotFound) {
            c.JSON(404, ErrorResponse(errno.TenantNotFound))
            return
        }
        c.JSON(500, ErrorResponse(errno.InternalError))
        return
    }

    c.JSON(200, SuccessResponse(tenant))
}

// 错误响应构建
func ErrorResponse(errCode int) map[string]interface{} {
    return map[string]interface{}{
        "code":      errCode,
        "message":   errno.GetMessage(errCode),
        "request_id": getRequestID(),
        "timestamp": time.Now().Format(time.RFC3339),
    }
}

func SuccessResponse(data interface{}) map[string]interface{} {
    return map[string]interface{}{
        "code":      0,
        "message":   "success",
        "data":      data,
        "request_id": getRequestID(),
    }
}

// ❌ Bad: 硬编码错误信息
func (h *TenantHandler) GetTenant(ctx context.Context, c *app.RequestContext) {
    tenantID := c.Param("id")

    tenant, err := h.tenantApp.GetTenantByID(ctx, tenantID)
    if err != nil {
        c.JSON(404, map[string]interface{}{
            "message": "tenant not found",  // 硬编码，不利于国际化
        })
        return
    }

    c.JSON(200, tenant)
}
```

---

## API 版本控制

### 版本策略

**URL 版本控制**（推荐）

```
/api/v1/tenants
/api/v2/tenants
```

**优点**:
- 清晰明确
- 易于路由
- 支持多版本共存

**请求头版本控制**（不推荐）

```
GET /api/tenants
Accept: application/vnd.zker.v1+json
```

**缺点**:
- 不够直观
- 难以调试
- 缓存困难

### 版本升级策略

| 场景 | 策略 | 示例 |
|------|------|------|
| **新增字段** | 向后兼容，不需要升级版本 | v1 响应中新增字段 |
| **修改字段** | 向后兼容，保留旧字段，新增字段 | v1 响应中保留 `tenantName`，新增 `tenant_name` |
| **删除字段** | 需要升级版本 | v2 移除 `tenantName` |
| **修改业务逻辑** | 需要升级版本 | v2 改变权限检查逻辑 |

### 版本废弃

```
# 响应头添加废弃警告
Warning: 299 - "Deprecated API, use /api/v2/tenants instead"

# 设置废弃日期
{
  "code": 0,
  "message": "success",
  "data": {...},
  "meta": {
    "deprecated": true,
    "sunset": "2025-06-01",
    "migration_guide": "https://docs.zker.com/api/v2-migration"
  }
}
```

---

## 分页规范

### 请求参数

```
GET /api/v1/tenants?page=1&page_size=20
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| **page** | int | 1 | 页码（从 1 开始） |
| **page_size** | int | 20 | 每页数量（1-100） |

### 响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
  }
}
```

### 游标分页（大数据量场景）

```
GET /api/v1/tenants?limit=20&cursor=MToxNjE0NTYyNDU4MDAwOjE2MTQ1NjI0NjAwMDA=
```

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [...],
    "next_cursor": "MjoxNjE0NTYyNDU4MDAwOjE2MTQ1NjI0NjAwMDA=",
    "has_more": true
  }
}
```

---

## 排序规范

### 请求参数

```
GET /api/v1/tenants?sort_by=created_at&order=desc
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| **sort_by** | string | created_at | 排序字段 |
| **order** | string | desc | 排序方向（asc/desc） |

### 多字段排序

```
GET /api/v1/tenants?sort=status:desc,created_at:desc
```

---

## 过滤规范

### 基本过滤

```
GET /api/v1/tenants?status=active&plan_type=pro
```

### 范围过滤

```
GET /api/v1/tenants?created_at[gte]=2025-01-01&created_at[lte]=2025-01-31
```

### 搜索过滤

```
GET /api/v1/tenants?search=示例
```

### 数组过滤

```
GET /api/v1/tenants?status=active,created_at:desc
```

---

## 检查清单

### API 设计检查清单

- [ ] URL 使用复数形式
- [ ] URL 使用连字符分隔
- [ ] 使用正确的 HTTP 方法
- [ ] 请求参数使用下划线分隔
- [ ] 响应包含 code、message、data
- [ ] 使用统一错误码
- [ ] 错误响应包含详细信息
- [ ] 支持分页、排序、过滤
- [ ] 包含 request_id
- [ ] 支持 CORS

### 安全检查清单

- [ ] 所有 API 都需要认证（公开 API 除外）
- [ ] 敏感操作需要二次验证
- [ ] 参数验证完整
- [ ] SQL 注入防护
- [ ] XSS 防护
- [ ] CSRF 防护
- [ ] 限流保护

---

## 相关文档

| 文档 | 说明 |
|------|------|
| [naming-conventions.md](naming-conventions.md) | 命名规范 |
| [database-design.md](database-design.md) | 数据库设计规范 |
| [error-handling.md](error-handling.md) | 错误处理规范 |
| [backend-dev-guide.md](backend-dev-guide.md) | 后端开发指南 |

---

**🎯 目标**: RESTful、一致、易用的 API 设计！
