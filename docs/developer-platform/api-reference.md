# 开发者平台 API参考文档

## 基础信息

**Base URL**: `https://api.coze.com/api/developer`

**认证方式**: Bearer Token (API Key)

**请求格式**: JSON

**响应格式**: JSON

## 通用响应格式

### 成功响应

```json
{
  "code": "SUCCESS",
  "message": "操作成功",
  "message_zh": "操作成功",
  "message_en": "Operation successful",
  "data": { ... }
}
```

### 错误响应

```json
{
  "code": "DEV401002",
  "message": "无效的API密钥",
  "message_zh": "无效的API密钥",
  "message_en": "Invalid API key",
  "details": { ... }
}
```

## 错误码列表

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| DEV201001 | 404 | 项目不存在 |
| DEV201002 | 409 | 项目名称已存在 |
| DEV401002 | 401 | 无效的API密钥 |
| DEV401003 | 401 | API密钥已过期 |
| DEV401006 | 403 | API密钥权限不足 |
| DEV402001 | 404 | Webhook不存在 |
| DEV403002 | 400 | 不支持的SDK语言 |

## API端点

### 1. 项目管理

#### 1.1 创建项目

**POST** `/projects`

**请求头**：
```
Authorization: Bearer sk_abc123...
Content-Type: application/json
```

**请求体**：
```json
{
  "project_name": "My AI Assistant",
  "project_type": "bot",
  "description": "An intelligent assistant",
  "config": {
    "enable_webhook": true,
    "enable_rate_limit": true,
    "rate_limit_per_min": 100,
    "enable_analytics": true
  }
}
```

**响应**：
```json
{
  "code": "SUCCESS",
  "data": {
    "project_id": "proj_abc123...",
    "tenant_id": "tenant_123",
    "developer_id": "dev_456",
    "project_name": "My AI Assistant",
    "project_type": "bot",
    "description": "An intelligent assistant",
    "status": "development",
    "config": { ... },
    "created_at": 1234567890000,
    "updated_at": 1234567890000
  }
}
```

#### 1.2 查询项目列表

**GET** `/projects`

**查询参数**：
- `project_type` (可选): 项目类型（bot/workflow/integration）
- `status` (可选): 项目状态（development/production/archived）
- `page_token` (可选): 分页令牌
- `page_size` (可选): 每页数量（默认20）

**响应**：
```json
{
  "code": "SUCCESS",
  "data": {
    "projects": [ ... ],
    "total": 100,
    "page_token": "next_page_token"
  }
}
```

#### 1.3 获取项目详情

**GET** `/projects/:project_id`

**响应**：
```json
{
  "code": "SUCCESS",
  "data": {
    "project_id": "proj_abc123...",
    "project_name": "My AI Assistant",
    "project_type": "bot",
    "status": "development",
    ...
  }
}
```

#### 1.4 更新项目

**PUT** `/projects/:project_id`

**请求体**：
```json
{
  "project_name": "Updated Name",
  "description": "Updated description",
  "config": { ... }
}
```

#### 1.5 发布项目

**POST** `/projects/:project_id/publish`

**响应**：
```json
{
  "code": "SUCCESS",
  "data": {
    "project_id": "proj_abc123...",
    "status": "production"
  }
}
```

#### 1.6 归档项目

**POST** `/projects/:project_id/archive`

#### 1.7 删除项目

**DELETE** `/projects/:project_id`

### 2. API密钥管理

#### 2.1 创建API密钥

**POST** `/api-keys`

**请求体**：
```json
{
  "project_id": "proj_abc123...",
  "key_name": "Production Key",
  "key_prefix": "sk",
  "scopes": ["read", "write"],
  "expires_in": 365
}
```

**响应**：
```json
{
  "code": "SUCCESS",
  "data": {
    "key_id": "key_abc123...",
    "key_name": "Production Key",
    "key_prefix": "sk",
    "key_masked": "sk_****1234",
    "key_secret": "sk_abc123...", // ⚠️ 只在创建时返回一次
    "scopes": ["read", "write"],
    "status": "active",
    "created_at": 1234567890000
  }
}
```

#### 2.2 查询API密钥列表

**GET** `/api-keys`

**查询参数**：
- `project_id` (可选): 项目ID
- `key_prefix` (可选): 密钥前缀（sk/pk/tk）
- `status` (可选): 密钥状态（active/revoked/expired）

**响应**：
```json
{
  "code": "SUCCESS",
  "data": {
    "api_keys": [
      {
        "key_id": "key_abc123...",
        "key_name": "Production Key",
        "key_masked": "sk_****1234",
        "scopes": ["read", "write"],
        "status": "active",
        "created_at": 1234567890000
      }
    ],
    "total": 10
  }
}
```

#### 2.3 获取API密钥详情

**GET** `/api-keys/:key_id`

#### 2.4 撤销API密钥

**POST** `/api-keys/:key_id/revoke`

#### 2.5 重新生成API密钥

**POST** `/api-keys/:key_id/regenerate`

**响应**：
```json
{
  "code": "SUCCESS",
  "data": {
    "key_id": "key_abc123...",
    "key_secret": "sk_new123...", // ⚠️ 新密钥，只返回一次
    "key_masked": "sk_****5678"
  }
}
```

### 3. Webhook管理

#### 3.1 创建Webhook

**POST** `/webhooks`

**请求体**：
```json
{
  "project_id": "proj_abc123...",
  "webhook_url": "https://your-app.com/webhooks/coze",
  "events": ["bot.created", "message.received"]
}
```

**响应**：
```json
{
  "code": "SUCCESS",
  "data": {
    "webhook_id": "webhook_abc123...",
    "webhook_url": "https://your-app.com/webhooks/coze",
    "events": ["bot.created", "message.received"],
    "status": "active",
    "webhook_secret": "wh_secret_...", // ⚠️ 只在创建时返回一次
    "created_at": 1234567890000
  }
}
```

#### 3.2 查询Webhook列表

**GET** `/webhooks`

**查询参数**：
- `project_id` (可选): 项目ID
- `status` (可选): Webhook状态（active/paused）

#### 3.3 获取Webhook详情

**GET** `/webhooks/:webhook_id`

#### 3.4 更新Webhook

**PUT** `/webhooks/:webhook_id`

**请求体**：
```json
{
  "webhook_url": "https://updated-url.com/webhooks",
  "events": ["bot.created", "bot.updated"]
}
```

#### 3.5 暂停Webhook

**POST** `/webhooks/:webhook_id/pause`

#### 3.6 恢复Webhook

**POST** `/webhooks/:webhook_id/resume`

#### 3.7 删除Webhook

**DELETE** `/webhooks/:webhook_id`

#### 3.8 获取Webhook统计

**GET** `/webhooks/:webhook_id/stats`

**响应**：
```json
{
  "code": "SUCCESS",
  "data": {
    "webhook_id": "webhook_abc123...",
    "total_calls": 1000,
    "success_calls": 950,
    "failure_calls": 50,
    "success_rate": 95.0,
    "avg_duration": 123.45,
    "last_triggered_at": 1234567890000
  }
}
```

### 4. SDK管理

#### 4.1 生成SDK

**POST** `/sdks/generate`

**请求体**：
```json
{
  "project_id": "proj_abc123...",
  "sdk_name": "my-assistant-sdk",
  "language": "python",
  "version": "1.0.0",
  "description": "Python SDK for My AI Assistant"
}
```

**响应**：
```json
{
  "code": "SUCCESS",
  "data": {
    "sdk_id": "sdk_abc123...",
    "sdk_name": "my-assistant-sdk",
    "language": "python",
    "version": "1.0.0",
    "status": "draft",
    "code": "import requests...",
    "readme": "# My Assistant SDK...",
    "created_at": 1234567890000
  }
}
```

#### 4.2 查询SDK列表

**GET** `/sdks`

**查询参数**：
- `project_id` (可选): 项目ID
- `language` (可选): 编程语言（python/javascript/go/java）
- `status` (可选): SDK状态（draft/published/archived）

#### 4.3 获取SDK详情

**GET** `/sdks/:sdk_id`

#### 4.4 发布SDK

**POST** `/sdks/:sdk_id/publish`

#### 4.5 下载SDK

**GET** `/sdks/:sdk_id/download`

**响应**：二进制文件（ZIP格式）

#### 4.6 删除SDK

**DELETE** `/sdks/:sdk_id`

## Webhook事件

### 事件类型

| 事件类型 | 说明 | 触发时机 |
|---------|------|---------|
| `bot.created` | Bot创建 | 新Bot创建成功后 |
| `bot.updated` | Bot更新 | Bot配置更新后 |
| `bot.deleted` | Bot删除 | Bot被删除后 |
| `workflow.completed` | 工作流完成 | 工作流执行完成 |
| `workflow.failed` | 工作流失败 | 工作流执行失败 |
| `conversation.created` | 对话创建 | 新对话创建后 |
| `message.received` | 消息接收 | 收到用户消息后 |
| `api.call` | API调用 | API被调用时 |
| `error.occurred` | 错误发生 | 系统错误发生时 |

### 事件格式

所有Webhook POST请求包含以下Headers：

```
Content-Type: application/json
X-Coze-Signature: sha256=abc123...
X-Coze-Event-ID: evt_abc123...
X-Coze-Event-Type: bot.created
X-Coze-Timestamp: 1234567890000
```

**请求体示例**：

```json
{
  "event_id": "evt_abc123...",
  "event_type": "bot.created",
  "timestamp": 1234567890000,
  "tenant_id": "tenant_123",
  "project_id": "proj_456",
  "data": {
    "bot_id": "bot_789",
    "name": "My Bot",
    "description": "A helpful assistant",
    "project_type": "bot"
  }
}
```

### 签名验证

Webhook使用HMAC-SHA256签名验证请求的真实性。

**验证步骤**：

1. 从请求头获取签名：`X-Coze-Signature`
2. 读取请求体
3. 使用Webhook Secret计算HMAC-SHA256
4. 比对签名

**Python验证示例**：

```python
import hmac
import hashlib

def verify_signature(payload: bytes, signature: str, secret: str) -> bool:
    expected = hmac.new(
        secret.encode(),
        payload,
        hashlib.sha256
    ).hexdigest()
    return hmac.compare_digest(f"sha256={expected}", signature)

# 使用
signature = request.headers.get("X-Coze-Signature")
payload = await request.body()

if not verify_signature(payload, signature, WEBHOOK_SECRET):
    return {"error": "Invalid signature"}, 401
```

## 速率限制

所有API端点都有速率限制，以防止滥用。

| 计划类型 | 限制 |
|---------|------|
| 免费版 | 100次/分钟 |
| 专业版 | 1,000次/分钟 |
| 企业版 | 10,000次/分钟 |

**响应头**：

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1234567890
```

**超限响应**：

```json
{
  "code": "COMMON429001",
  "message": "请求过于频繁，请稍后再试",
  "message_en": "Rate limit exceeded",
  "retry_after": 60
}
```

## SDK客户端库

### Python

```python
pip install coze-studio-sdk
```

```python
from coze_studio import CozeClient

client = CozeClient(api_key="sk_abc123...")
```

### JavaScript/TypeScript

```bash
npm install @coze-studio/sdk
```

```typescript
import { CozeClient } from '@coze-studio/sdk';

const client = new CozeClient({ apiKey: 'sk_abc123...' });
```

### Go

```bash
go get github.com/coze-studio/sdk-go
```

```go
import "github.com/coze-studio/sdk-go"

client := coze.NewClient("sk_abc123...")
```

### Java

```xml
<dependency>
    <groupId>com.coze-studio</groupId>
    <artifactId>sdk-java</artifactId>
    <version>1.0.0</version>
</dependency>
```

```java
import com.coze.studio.sdk.CozeClient;

CozeClient client = new CozeClient("sk_abc123...");
```

---

**API版本**：v1.0
**最后更新**：2025-01-01
