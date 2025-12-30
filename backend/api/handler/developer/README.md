# 开发者平台 Handler 层实现文档

> 企业级开发者平台 RESTful API Handler 层
>
> 版本: v1.0.0 | 更新: 2025-01-03

---

## 目录结构

```
backend/api/handler/developer/
├── project_handler.go     # 项目管理 Handler
├── apikey_handler.go      # API密钥管理 Handler
├── sdk_handler.go         # SDK生成 Handler
├── webhook_handler.go     # Webhook管理 Handler
└── README.md             # 本文档

backend/api/router/developer/
└── developer_router.go    # 路由注册
```

---

## 核心功能

### 1. 项目管理 (ProjectManagementHandler)

| API | 方法 | 路径 | 说明 |
|-----|------|------|------|
| ListProjects | GET | `/api/developer/projects` | 列出项目 |
| CreateProject | POST | `/api/developer/projects` | 创建项目 |
| GetProject | GET | `/api/developer/projects/:id` | 获取项目详情 |
| UpdateProject | PUT | `/api/developer/projects/:id` | 更新项目 |
| DeleteProject | DELETE | `/api/developer/projects/:id` | 删除项目 |
| PublishProject | POST | `/api/developer/projects/:id/publish` | 发布项目 |
| ArchiveProject | POST | `/api/developer/projects/:id/archive` | 归档项目 |
| GetProjectStats | GET | `/api/developer/projects/:id/stats` | 获取统计信息 |
| GetDeveloperProjects | GET | `/api/developer/developers/:id/projects` | 获取开发者项目 |
| GetTenantProjects | GET | `/api/developer/tenant/projects` | 获取租户项目 |

**请求示例**：
```bash
# 创建项目
curl -X POST http://localhost:8888/api/developer/projects \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant_123" \
  -d '{
    "project_name": "My AI Project",
    "project_type": "bot",
    "description": "An AI assistant project",
    "developer_id": "dev_123"
  }'
```

**响应示例**：
```json
{
  "code": 0,
  "message": "SUCCESS",
  "message_zh": "操作成功",
  "data": {
    "project_id": "proj_abc123",
    "project_name": "My AI Project",
    "project_type": "bot",
    "status": "development"
  }
}
```

---

### 2. API密钥管理 (APIKeyManagementHandler)

| API | 方法 | 路径 | 说明 |
|-----|------|------|------|
| ListAPIKeys | GET | `/api/developer/api-keys` | 列出API密钥 |
| CreateAPIKey | POST | `/api/developer/api-keys` | 创建API密钥 |
| GetAPIKey | GET | `/api/developer/api-keys/:id` | 获取密钥详情 |
| UpdateAPIKey | PUT | `/api/developer/api-keys/:id` | 更新密钥 |
| DeleteAPIKey | DELETE | `/api/developer/api-keys/:id` | 删除密钥 |
| RevokeAPIKey | POST | `/api/developer/api-keys/:id/revoke` | 撤销密钥 |
| RegenerateAPIKey | POST | `/api/developer/api-keys/:id/regenerate` | 重新生成密钥 |
| GetExpiringKeys | GET | `/api/developer/api-keys/expiring` | 获取即将过期的密钥 |

**请求示例**：
```bash
# 创建API密钥
curl -X POST http://localhost:8888/api/developer/api-keys \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant_123" \
  -d '{
    "project_id": "proj_abc123",
    "key_name": "Production Key",
    "key_prefix": "sk",
    "scopes": ["read", "write"],
    "expires_in": 365
  }'
```

**响应示例**（仅在创建时返回完整密钥）：
```json
{
  "code": 0,
  "message": "API key created successfully. Please save the key_secret securely.",
  "message_zh": "API密钥创建成功，请妥善保管key_secret",
  "data": {
    "api_key": {...},
    "key_secret": "sk_xxxxxxxxxxxxx",  // 仅在创建时返回
    "message": "Please save the key_secret securely"
  }
}
```

---

### 3. SDK生成 (SDKGeneratorHandler)

| API | 方法 | 路径 | 说明 |
|-----|------|------|------|
| ListSDKs | GET | `/api/developer/sdks` | 列出SDK |
| GenerateSDK | POST | `/api/developer/sdks/generate` | 生成SDK |
| GetSDK | GET | `/api/developer/sdks/:id` | 获取SDK详情 |
| PublishSDK | POST | `/api/developer/sdks/:id/publish` | 发布SDK |
| DeleteSDK | DELETE | `/api/developer/sdks/:id` | 删除SDK |
| DownloadSDK | GET | `/api/developer/sdks/:id/download` | 下载SDK |
| GetSDKCode | GET | `/api/developer/sdks/:id/code` | 获取SDK代码预览 |
| GetSupportedLanguages | GET | `/api/developer/sdks/languages` | 获取支持的语言 |

**请求示例**：
```bash
# 生成Python SDK
curl -X POST http://localhost:8888/api/developer/sdks/generate \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant_123" \
  -d '{
    "project_id": "proj_abc123",
    "sdk_name": "MyBotSDK",
    "language": "python",
    "version": "1.0.0",
    "description": "Python SDK for MyBot"
  }'
```

---

### 4. Webhook管理 (WebhookHandler)

| API | 方法 | 路径 | 说明 |
|-----|------|------|------|
| ListWebhooks | GET | `/api/developer/webhooks` | 列出Webhook |
| CreateWebhook | POST | `/api/developer/webhooks` | 创建Webhook |
| GetWebhook | GET | `/api/developer/webhooks/:id` | 获取Webhook详情 |
| UpdateWebhook | PUT | `/api/developer/webhooks/:id` | 更新Webhook |
| DeleteWebhook | DELETE | `/api/developer/webhooks/:id` | 删除Webhook |
| TestWebhook | POST | `/api/developer/webhooks/:id/test` | 测试Webhook |
| PauseWebhook | POST | `/api/developer/webhooks/:id/pause` | 暂停Webhook |
| ResumeWebhook | POST | `/api/developer/webhooks/:id/resume` | 恢复Webhook |
| GetWebhookStats | GET | `/api/developer/webhooks/:id/stats` | 获取统计信息 |
| GetSupportedEvents | GET | `/api/developer/webhooks/events` | 获取支持的事件 |

**请求示例**：
```bash
# 创建Webhook
curl -X POST http://localhost:8888/api/developer/webhooks \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant_123" \
  -d '{
    "project_id": "proj_abc123",
    "webhook_url": "https://example.com/webhook",
    "events": ["bot.created", "bot.updated", "workflow.completed"]
  }'
```

---

## 集成方式

### 1. 依赖注入

在应用启动时初始化 Handler：

```go
package main

import (
    devrouter "github.com/coze-dev/coze-studio/backend/api/router/developer"
    devservice "github.com/coze-dev/coze-studio/backend/domain/developer/service"
)

func main() {
    // 初始化服务
    projectSvc := devservice.NewProjectManagementService(...)
    apiKeySvc := devservice.NewAPIKeyManagementService(...)
    sdkSvc := devservice.NewSDKGeneratorService(...)
    webhookSvc := devservice.NewWebhookService(...)

    // 初始化Handler
    devrouter.InitHandlers(projectSvc, apiKeySvc, sdkSvc, webhookSvc)

    // 注册路由
    devrouter.RegisterRoutes(r)
}
```

### 2. 中间件配置

所有开发者平台 API 自动应用以下中间件：
- **TenantIsolationMiddleware**: 租户隔离
- **PermissionCheckMiddleware**: 权限检查

### 3. 请求头

所有请求需要携带租户ID：
```
X-Tenant-ID: tenant_123
```

---

## 错误码

| 错误码 | 说明 | HTTP状态码 |
|-------|------|-----------|
| DEV201001 | 项目不存在 | 404 |
| DEV201002 | 项目名称已存在 | 409 |
| DEV201003 | 无效的项目类型 | 400 |
| DEV201004 | 无效的项目状态 | 400 |
| DEV201005 | 项目数量超限 | 403 |
| DEV201006 | 项目正在生产环境 | 403 |
| DEV401001 | API密钥不存在 | 404 |
| DEV401002 | 无效的API密钥 | 401 |
| DEV401003 | API密钥已过期 | 401 |
| DEV401004 | API密钥已撤销 | 401 |
| DEV401005 | API密钥数量超限 | 403 |
| DEV402001 | Webhook不存在 | 404 |
| DEV402002 | 无效的Webhook URL | 400 |
| DEV402003 | Webhook签名验证失败 | 401 |
| DEV402004 | 无效的Webhook事件 | 400 |
| DEV402005 | Webhook数量超限 | 403 |
| DEV403001 | SDK不存在 | 404 |
| DEV403002 | 不支持的SDK语言 | 400 |
| DEV403003 | 无效的SDK版本 | 400 |
| DEV403004 | SDK生成失败 | 500 |

---

## 开发规范

### 代码风格

1. **遵循 SOLID 原则**
   - 单一职责：每个 Handler 只处理一类资源
   - 接口隔离：使用清晰的服务接口
   - 依赖注入：通过构造函数注入依赖

2. **统一错误处理**
   ```go
   // 使用统一错误码
   httputil.BuildErrorResp(c, berrno.ErrProjectNotFound.Code, ...)

   // 或使用内部错误处理
   httputil.InternalError(ctx, c, err)
   ```

3. **租户隔离**
   ```go
   // 获取租户ID
   tenantID := ctxcache.GetTenantIDFromCtx(ctx)

   // 验证租户隔离
   if resource.TenantID != tenantID {
       httputil.BuildErrorResp(c, berrno.ErrProjectNotFound.Code, ...)
       return
   }
   ```

4. **统一响应格式**
   ```go
   httputil.BuildSuccessResp(c, data)
   ```

---

## 测试

### 单元测试

```bash
go test ./backend/api/handler/developer/...
```

### 集成测试

```bash
# 启动服务
make server

# 测试 API
curl http://localhost:8888/api/developer/projects
```

---

## 更新日志

### v1.0.0 (2025-01-03)

- [x] 实现项目管理 Handler (10个API)
- [x] 实现API密钥管理 Handler (8个API)
- [x] 实现SDK生成 Handler (8个API)
- [x] 实现Webhook管理 Handler (10个API)
- [x] 创建路由注册文件
- [x] 添加企业级错误码支持
- [x] 添加租户隔离检查
- [x] 添加权限检查

---

## 相关文档

- [错误码定义](../../types/errno/developer.go)
- [Domain层服务接口](../../domain/developer/service/)
- [实体定义](../../domain/developer/entity/)
- [开发规范](../../../docs/02-SPECS/backend-dev-guide.md)
