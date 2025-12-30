# 开发者平台快速开始指南

## 概述

开发者平台（Developer Platform）为企业用户提供完整的API开发能力，包括项目管理、API密钥管理、SDK生成和Webhook集成。

## 核心功能

### 1. 项目管理
- 创建和管理AI项目（Bot、Workflow、Integration）
- 项目生命周期管理（开发、生产、归档）
- 项目配置管理

### 2. API密钥管理
- 安全的API密钥生成和存储
- 密钥权限控制（read/write/admin）
- 密钥过期管理
- 密钥使用统计

### 3. SDK生成
- 多语言SDK自动生成（Python、JavaScript、Go、Java）
- 版本管理
- 一键下载

### 4. Webhook集成
- 事件订阅
- 异步回调
- 签名验证
- 统计监控

## 数据库初始化

### 1. 执行迁移SQL

```bash
# 连接到MySQL数据库
mysql -u root -p coze_studio

# 执行迁移文件
source backend/migrations/developer/001_create_developer_tables.sql
```

### 2. 验证表创建

```sql
-- 查看创建的表
SHOW TABLES LIKE 'developer%';

-- 预期输出：
-- developer_projects
-- developer_api_keys
-- developer_sdks
-- developer_webhooks
-- developer_webhook_logs
-- developers
```

## 服务初始化

### 1. 初始化Repository

```go
package developer

import (
    "gorm.io/gorm"
    "github.com/coze-dev/coze-studio/backend/domain/developer/repository"
)

// InitRepositories 初始化Repository层
func InitRepositories(db *gorm.DB) {
    developerRepo := repository.NewDeveloperRepository(db)
    projectRepo := repository.NewProjectRepository(db)
    apiKeyRepo := repository.NewAPIKeyRepository(db)
    webhookRepo := repository.NewWebhookRepository(db)
    webhookLogRepo := repository.NewWebhookLogRepository(db)
    sdkRepo := repository.NewSDKRepository(db)

    // 注册到依赖注入容器
    // ...
}
```

### 2. 初始化Service

```go
package developer

import (
    "github.com/coze-dev/coze-studio/backend/domain/developer/service"
)

// InitServices 初始化Service层
func InitServices(
    projectRepo repository.ProjectRepository,
    apiKeyRepo repository.APIKeyRepository,
    webhookRepo repository.WebhookRepository,
    logRepo repository.WebhookLogRepository,
    sdkRepo repository.SDKRepository,
    secretKey string,
) {
    // 项目管理服务
    projectService := service.NewProjectManagementService(
        projectRepo,
        apiKeyRepo,
        webhookRepo,
        sdkRepo,
    )

    // API密钥管理服务
    apiKeyService := service.NewAPIKeyManagementService(
        apiKeyRepo,
        projectRepo,
        secretKey, // 从配置文件读取
    )

    // SDK生成服务
    sdkService := service.NewSDKGeneratorService(
        sdkRepo,
        projectRepo,
    )

    // Webhook服务
    webhookService := service.NewWebhookService(
        webhookRepo,
        logRepo,
    )

    // 注册到依赖注入容器
    // ...
}
```

## 使用示例

### 场景1：创建第一个项目

```go
// 1. 创建项目
project, err := projectService.CreateProject(ctx, &service.CreateProjectRequest{
    TenantID:    "tenant_123",
    DeveloperID: "dev_456",
    ProjectName: "Customer Support Bot",
    ProjectType: entity.ProjectTypeBot,
    Description: "AI-powered customer support assistant",
    Config: &entity.ProjectConfig{
        EnableWebhook:    true,
        EnableRateLimit:  true,
        RateLimitPerMin:  100,
        EnableAnalytics:  true,
    },
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Project created: %s\n", project.ProjectID)
```

### 场景2：创建API密钥

```go
// 2. 创建生产环境API密钥
expiresIn := 365 // 365天

apiKey, keySecret, err := apiKeyService.CreateAPIKey(ctx, &service.CreateAPIKeyRequest{
    TenantID:  "tenant_123",
    ProjectID: project.ProjectID,
    KeyName:   "Production Key",
    KeyPrefix: entity.APIKeyPrefixSecret,
    Scopes: entity.APIKeyScopes{
        entity.APIKeyScopeRead,
        entity.APIKeyScopeWrite,
    },
    ExpiresIn: &expiresIn,
})

if err != nil {
    log.Fatal(err)
}

// ⚠️ 重要：keySecret只在创建时返回一次，请妥善保存！
fmt.Printf("API Key Secret: %s\n", keySecret)
fmt.Printf("API Key Masked: %s\n", apiKey.KeyMasked) // sk_****1234
```

### 场景3：验证API密钥

```go
// 3. 验证API密钥
apiKey, err := apiKeyService.ValidateAPIKey(ctx, keySecret)
if err != nil {
    // 处理错误：密钥无效、已过期或已撤销
    log.Printf("API key validation failed: %v", err)
    return
}

// 检查权限
if !apiKey.HasScope(entity.APIKeyScopeWrite) {
    log.Printf("Permission denied: missing write scope")
    return
}

// 更新最后使用时间
_ = apiKeyService.UpdateLastUsed(ctx, apiKey.KeyID)

log.Printf("API key validated successfully: %s", apiKey.KeyName)
```

### 场景4：生成Python SDK

```go
// 4. 为项目生成Python SDK
sdk, err := sdkService.GenerateSDK(ctx, &service.GenerateSDKRequest{
    TenantID:    "tenant_123",
    ProjectID:   project.ProjectID,
    SDKName:     "customer-support-sdk",
    Language:    entity.SDKLanguagePython,
    Version:     "1.0.0",
    Description: "Python SDK for Customer Support Bot",
})

if err != nil {
    log.Fatal(err)
}

// 5. 发布SDK
err = sdkService.PublishSDK(ctx, sdk.SDKID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("SDK generated and published: %s\n", sdk.SDKID)
```

### 场景5：创建Webhook

```go
// 6. 创建Webhook监听Bot创建事件
webhook, secret, err := webhookService.CreateWebhook(ctx, &service.CreateWebhookRequest{
    TenantID:   "tenant_123",
    ProjectID:  project.ProjectID,
    WebhookURL: "https://your-app.com/webhooks/coze",
    Events: entity.WebhookEvents{
        entity.WebhookEventBotCreated,
        entity.WebhookEventMessageReceived,
    },
})

if err != nil {
    log.Fatal(err)
}

// ⚠️ 重要：secret用于验证Webhook签名，请妥善保存！
fmt.Printf("Webhook created: %s\n", webhook.WebhookID)
fmt.Printf("Webhook Secret: %s\n", secret)
```

### 场景6：验证Webhook签名

```go
// 7. 在你的Webhook接收端验证签名
func handleWebhook(w http.ResponseWriter, r *http.Request) {
    // 读取请求体
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read body", http.StatusBadRequest)
        return
    }

    // 获取签名
    signature := r.Header.Get("X-Coze-Signature")

    // 验证签名
    if !webhookService.VerifyWebhookSignature(r.Context(), body, signature, webhookSecret) {
        http.Error(w, "Invalid signature", http.StatusUnauthorized)
        return
    }

    // 处理Webhook事件
    var payload entity.WebhookPayload
    if err := json.Unmarshal(body, &payload); err != nil {
        http.Error(w, "Invalid payload", http.StatusBadRequest)
        return
    }

    log.Printf("Received webhook event: %s", payload.EventType)

    // 返回成功
    w.WriteHeader(http.StatusOK)
}
```

## SDK使用示例

### Python SDK

```python
from coze_studio import CozeClient

# 初始化客户端
client = CozeClient(api_key="sk_abc123...")

# 创建Bot
bot = client.create_bot(
    name="My Bot",
    description="A helpful assistant"
)
print(f"Bot created: {bot['bot_id']}")

# 发送消息
response = client.send_message(
    bot_id=bot['bot_id'],
    message="Hello, World!"
)
print(response)
```

### JavaScript/TypeScript SDK

```typescript
import { CozeClient } from '@coze-studio/customer-support-sdk';

// 初始化客户端
const client = new CozeClient({
  apiKey: 'sk_abc123...',
});

// 创建Bot
const bot = await client.createBot({
  name: 'My Bot',
  description: 'A helpful assistant',
});
console.log(`Bot created: ${bot.bot_id}`);

// 发送消息
const response = await client.sendMessage({
  botId: bot.bot_id,
  message: 'Hello, World!',
});
console.log(response);
```

### Go SDK

```go
package main

import (
    "fmt"
    "github.com/coze-studio/customer-support-sdk"
)

func main() {
    // 初始化客户端
    client := coze.NewClient("sk_abc123...")

    // 创建Bot
    resp, err := client.CreateBot(&coze.CreateBotRequest{
        Name: "My Bot",
        Description: "A helpful assistant",
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Bot created: %s\n", resp.BotID)
}
```

### Java SDK

```java
import com.coze.studio.sdk.CozeClient;
import com.coze.studio.sdk.CreateBotRequest;

public class Main {
    public static void main(String[] args) throws IOException {
        // 初始化客户端
        CozeClient client = new CozeClient("sk_abc123...");

        // 创建Bot
        CreateBotRequest request = new CreateBotRequest("My Bot");
        request.setDescription("A helpful assistant");

        CreateBotResponse response = client.createBot(request);
        System.out.println("Bot created: " + response.getBotId());
    }
}
```

## Webhook事件格式

所有Webhook POST请求都包含以下Headers：

```
Content-Type: application/json
X-Coze-Signature: sha256=...
X-Coze-Event-ID: evt_abc123...
X-Coze-Event-Type: bot.created
X-Coze-Timestamp: 1234567890000
```

请求体格式：

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
    "description": "A helpful assistant"
  }
}
```

## 配置说明

### 环境变量

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=coze_studio

# API密钥加密密钥（必须保密！）
API_KEY_ENCRYPTION_KEY=your-32-character-secret-key

# Webhook超时配置（秒）
WEBHOOK_TIMEOUT=30
```

### 配置文件

```yaml
# config/developer.yaml
developer:
  # API密钥配置
  api_key:
    # 密钥前缀
    prefixes:
      secret: "sk"    # Secret Key
      public: "pk"    # Public Key
      test: "tk"      # Test Key

    # 默认过期时间（天）
    default_expires_in: 365

    # 密钥数量限制
    max_keys_per_project: 10

  # SDK配置
  sdk:
    # 支持的语言
    languages:
      - python
      - javascript
      - go
      - java

    # 版本号格式
    version_format: "semver"

  # Webhook配置
  webhook:
    # 超时时间（秒）
    timeout: 30

    # 重试次数
    max_retries: 3

    # 重试间隔（秒）
    retry_interval: 60

    # 日志保留天数
    log_retention_days: 30
```

## 常见问题

### Q1: API密钥丢失怎么办？

API密钥只在创建时完整显示一次，丢失后需要重新生成：

```go
newKey, newSecret, err := apiKeyService.RegenerateAPIKey(ctx, oldKeyID)
// 旧密钥将立即失效
```

### Q2: 如何限制API密钥权限？

创建密钥时指定`scopes`：

```go
apiKey, _, err := apiKeyService.CreateAPIKey(ctx, &service.CreateAPIKeyRequest{
    Scopes: entity.APIKeyScopes{
        entity.APIKeyScopeRead, // 只读权限
    },
    // ...
})
```

### Q3: Webhook触发失败怎么办？

1. 检查Webhook URL是否可访问
2. 查看Webhook日志获取详细错误信息
3. 验证Webhook Secret是否正确

```go
stats, err := webhookService.GetWebhookStats(ctx, webhookID)
if stats.FailureRate > 50 {
    // 失败率过高，需要检查
}
```

### Q4: 如何更新SDK版本？

```go
// 生成新版本SDK
newSDK, err := sdkService.GenerateSDK(ctx, &service.GenerateSDKRequest{
    ProjectID: projectID,
    Language:  entity.SDKLanguagePython,
    Version:   "2.0.0", // 新版本
    // ...
})
```

## 监控和运维

### 查看项目统计

```go
stats, err := projectService.GetProjectStats(ctx, projectID)
fmt.Printf("Total API Keys: %d\n", stats.TotalAPIKeys)
fmt.Printf("Active Webhooks: %d\n", stats.ActiveWebhooks)
```

### 查看Webhook统计

```go
stats, err := webhookService.GetWebhookStats(ctx, webhookID)
fmt.Printf("Success Rate: %.2f%%\n", stats.SuccessRate)
fmt.Printf("Avg Duration: %.2fms\n", stats.AvgDuration)
```

### 清理过期数据

```go
// 删除30天前的Webhook日志
deletedCount, err := webhookLogRepo.DeleteOldLogs(ctx, 30)
fmt.Printf("Deleted %d old logs\n", deletedCount)
```

## 下一步

1. 查看完整API文档：`docs/developer-platform/api.md`
2. 查看SDK使用指南：`docs/developer-platform/sdk-guide.md`
3. 查看Webhook集成指南：`docs/developer-platform/webhook-guide.md`

---

**文档版本**：v1.0
**最后更新**：2025-01-01
