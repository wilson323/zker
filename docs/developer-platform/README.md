# 开发者平台（Developer Platform）

> 完整的企业级开发者平台实现，提供项目管理、API密钥管理、SDK生成和Webhook集成功能。

## 📋 目录

- [概述](#概述)
- [功能特性](#功能特性)
- [实现内容](#实现内容)
- [快速开始](#快速开始)
- [文档](#文档)
- [技术架构](#技术架构)
- [API文档](#api文档)
- [贡献指南](#贡献指南)

## 概述

开发者平台是为企业用户提供的完整API开发能力，包括：

- ✅ **项目管理**：创建和管理AI项目（Bot、Workflow、Integration）
- ✅ **API密钥管理**：安全的密钥生成、权限控制和过期管理
- ✅ **SDK生成**：自动生成多语言SDK（Python、JavaScript、Go、Java）
- ✅ **Webhook集成**：事件订阅、异步回调和签名验证

## 功能特性

### 项目管理
- 创建和管理AI项目（Bot、Workflow、Integration）
- 项目生命周期管理（开发、生产、归档）
- 项目配置管理（速率限制、缓存、分析等）
- 租户隔离和数据权限控制

### API密钥管理
- 安全的API密钥生成（AES-256加密存储）
- 密钥权限控制（read/write/admin）
- 密钥脱敏显示（sk_****1234）
- 密钥过期管理
- 密钥使用统计

### SDK生成
- 多语言SDK自动生成：
  - Python SDK（使用requests库）
  - JavaScript/TypeScript SDK（使用axios）
  - Go SDK（使用net/http）
  - Java SDK（使用OkHttp）
- SDK版本管理
- 一键下载

### Webhook集成
- 9种事件类型（bot.created、message.received等）
- 异步触发（不阻塞主流程）
- HMAC-SHA256签名验证
- Webhook统计（成功率、平均耗时）

## 实现内容

### 代码统计

| 层级 | 文件数 | 代码行数 | 说明 |
|------|--------|----------|------|
| Entity | 5 | ~750 | 5个实体定义 |
| Repository | 6 | ~1,200 | 接口+实现 |
| Service | 4 | ~2,000 | 4个核心服务 |
| Database | 1 | ~250 | 迁移SQL |
| Errno | 1 | ~50 | 错误码定义 |
| 文档 | 4 | ~2,000 | 完整文档 |
| **总计** | **21** | **~6,250** | **完整实现** |

### 文件清单

**实体层**：
- `backend/domain/developer/entity/developer.go`
- `backend/domain/developer/entity/project.go`
- `backend/domain/developer/entity/api_key.go`
- `backend/domain/developer/entity/webhook.go`
- `backend/domain/developer/entity/sdk.go`

**Repository层**：
- `backend/domain/developer/repository/developer_repository.go`
- `backend/domain/developer/repository/developer_repository_impl.go`
- `backend/domain/developer/repository/project_repository_impl.go`
- `backend/domain/developer/repository/api_key_repository_impl.go`
- `backend/domain/developer/repository/webhook_repository_impl.go`
- `backend/domain/developer/repository/sdk_repository_impl.go`

**Service层**：
- `backend/domain/developer/service/project_management_service.go`
- `backend/domain/developer/service/api_key_management_service.go`
- `backend/domain/developer/service/sdk_generator_service.go`
- `backend/domain/developer/service/webhook_service.go`

**数据库**：
- `backend/migrations/developer/001_create_developer_tables.sql`

**错误码**：
- `backend/types/errno/developer.go`

**文档**：
- `docs/developer-platform/implementation-summary.md`
- `docs/developer-platform/quick-start.md`
- `docs/developer-platform/api-reference.md`
- `docs/developer-platform/README.md`

## 快速开始

### 1. 数据库初始化

```bash
mysql -u root -p coze_studio < backend/migrations/developer/001_create_developer_tables.sql
```

### 2. 创建第一个项目

```go
project, err := projectService.CreateProject(ctx, &service.CreateProjectRequest{
    TenantID:    "tenant_123",
    DeveloperID: "dev_456",
    ProjectName: "Customer Support Bot",
    ProjectType: entity.ProjectTypeBot,
    Description: "AI-powered customer support assistant",
})
```

### 3. 创建API密钥

```go
apiKey, keySecret, err := apiKeyService.CreateAPIKey(ctx, &service.CreateAPIKeyRequest{
    TenantID:  "tenant_123",
    ProjectID: project.ProjectID,
    KeyName:   "Production Key",
    KeyPrefix: entity.APIKeyPrefixSecret,
    Scopes:    entity.APIKeyScopes{entity.APIKeyScopeRead, entity.APIKeyScopeWrite},
    ExpiresIn: &days365,
})

// ⚠️ 重要：keySecret只在创建时返回一次，请妥善保存！
fmt.Println("API Key:", keySecret)
```

### 4. 生成Python SDK

```go
sdk, err := sdkService.GenerateSDK(ctx, &service.GenerateSDKRequest{
    TenantID:    "tenant_123",
    ProjectID:   project.ProjectID,
    SDKName:     "customer-support-sdk",
    Language:    entity.SDKLanguagePython,
    Version:     "1.0.0",
    Description: "Python SDK for Customer Support Bot",
})
```

### 5. 创建Webhook

```go
webhook, secret, err := webhookService.CreateWebhook(ctx, &service.CreateWebhookRequest{
    TenantID:   "tenant_123",
    ProjectID:  project.ProjectID,
    WebhookURL: "https://your-app.com/webhooks/coze",
    Events:     entity.WebhookEvents{entity.WebhookEventBotCreated},
})

// ⚠️ 重要：secret用于验证Webhook签名
fmt.Println("Webhook Secret:", secret)
```

更多使用示例请查看[快速开始指南](./quick-start.md)。

## 文档

- [实现总结](./implementation-summary.md) - 完整的实现总结和代码统计
- [快速开始](./quick-start.md) - 快速上手指南和使用示例
- [API参考](./api-reference.md) - 完整的API文档

## 技术架构

### 分层架构

```
┌─────────────────────────────────────────┐
│         API Handler Layer              │  ← HTTP接口（待实现）
├─────────────────────────────────────────┤
│         Service Layer                  │  ← 业务逻辑层
│  - ProjectManagementService            │
│  - APIKeyManagementService             │
│  - SDKGeneratorService                 │
│  - WebhookService                      │
├─────────────────────────────────────────┤
│         Repository Layer               │  ← 数据访问层
│  - DeveloperRepository                 │
│  - ProjectRepository                   │
│  - APIKeyRepository                    │
│  - WebhookRepository                   │
│  - SDKRepository                       │
├─────────────────────────────────────────┤
│         Entity Layer                   │  ← 实体定义层
│  - Developer                           │
│  - Project                             │
│  - APIKey                              │
│  - Webhook                             │
│  - SDK                                 │
└─────────────────────────────────────────┘
         ↓
┌─────────────────────────────────────────┐
│         Database Layer                 │  ← MySQL 8.4.5
│  - developers                          │
│  - developer_projects                  │
│  - developer_api_keys                  │
│  - developer_webhooks                  │
│  - developer_webhook_logs              │
│  - developer_sdks                      │
└─────────────────────────────────────────┘
```

### 核心设计原则

- ✅ **SOLID原则**：单一职责、开闭原则、里氏替换、接口隔离、依赖倒置
- ✅ **DDD分层**：领域驱动设计，清晰的分层架构
- ✅ **租户隔离**：所有表带tenant_id，实现多租户隔离
- ✅ **软删除**：使用deleted_at，避免数据丢失
- ✅ **安全优先**：密钥加密、签名验证、权限控制
- ✅ **可扩展性**：插件化设计、多语言SDK支持

## API文档

### Base URL
```
https://api.coze.com/api/developer
```

### 主要端点

**项目管理**：
- `POST /projects` - 创建项目
- `GET /projects` - 查询项目列表
- `GET /projects/:id` - 获取项目详情
- `PUT /projects/:id` - 更新项目
- `POST /projects/:id/publish` - 发布项目
- `DELETE /projects/:id` - 删除项目

**API密钥**：
- `POST /api-keys` - 创建API密钥
- `GET /api-keys` - 查询API密钥列表
- `POST /api-keys/:id/revoke` - 撤销API密钥
- `POST /api-keys/:id/regenerate` - 重新生成API密钥

**Webhook**：
- `POST /webhooks` - 创建Webhook
- `GET /webhooks` - 查询Webhook列表
- `POST /webhooks/:id/pause` - 暂停Webhook
- `POST /webhooks/:id/resume` - 恢复Webhook

**SDK**：
- `POST /sdks/generate` - 生成SDK
- `GET /sdks` - 查询SDK列表
- `POST /sdks/:id/publish` - 发布SDK
- `GET /sdks/:id/download` - 下载SDK

完整的API文档请查看[API参考](./api-reference.md)。

## 安全特性

### API密钥安全
- ✅ AES-256加密存储
- ✅ 密钥脱敏显示（sk_****1234）
- ✅ 支持多种前缀（sk/pk/tk）
- ✅ 权限范围控制（read/write/admin）
- ✅ 过期时间自动检查

### Webhook安全
- ✅ HMAC-SHA256签名验证
- ✅ 事件ID追踪
- ✅ 超时控制（30秒）
- ✅ 失败日志记录

### 租户隔离
- ✅ 所有表带tenant_id
- ✅ Repository层自动过滤
- ✅ Service层权限检查

## 待实现功能

虽然核心功能已完成，但以下功能建议进一步完善：

### 1. API Handler层（~1,200行）
需要创建HTTP处理器：
- `ProjectManagementHandler`
- `APIKeyManagementHandler`
- `SDKGeneratorHandler`
- `WebhookHandler`

### 2. 单元测试（~1,000行）
需要为每个Service编写单元测试：
- Repository层测试（Mock GORM）
- Service层测试（Mock Repository）
- 覆盖率要求：≥80%

### 3. 集成测试
- API端到端测试
- Webhook触发测试
- SDK生成测试

### 4. 前端界面
- 项目管理页面
- API密钥管理页面
- SDK生成和下载页面
- Webhook配置页面

## 贡献指南

欢迎贡献代码！请遵循以下规范：

### 代码规范
- 遵循[企业级开发规范手册](../../企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- 函数长度 < 50行
- 参数数量 < 5个
- 使用统一的错误码

### Git提交规范
```bash
feat(developer): add project management API
fix(developer): resolve api key validation bug
docs(developer): update api documentation
test(developer): add unit tests for sdk service
```

### Pull Request清单
- [ ] 代码符合开发规范
- [ ] 单元测试通过，覆盖率≥80%
- [ ] 文档已更新
- [ ] API文档已更新
- [ ] 没有新的安全漏洞

## 许可证

Apache License 2.0

## 联系方式

- 项目地址：https://github.com/coze-dev/coze-studio
- 问题反馈：https://github.com/coze-dev/coze-studio/issues

---

**实现日期**：2025-01-01
**版本**：v1.0
**代码行数**：~6,250行（含文档）
**文件数量**：21个文件
