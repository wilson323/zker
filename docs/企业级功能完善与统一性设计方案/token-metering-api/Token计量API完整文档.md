# Token计量API完整文档

## 📋 目录

- [概述](#概述)
- [快速开始](#快速开始)
- [API端点](#api端点)
- [数据模型](#数据模型)
- [错误码](#错误码)
- [性能指标](#性能指标)
- [使用示例](#使用示例)
- [最佳实践](#最佳实践)

---

## 概述

### 功能说明

Token计量API提供了完整的AI模型使用追踪和成本计算功能，支持：

- ✅ **Token使用记录**: 记录每次AI模型调用的详细Token消耗和成本
- ✅ **批量记录**: 高性能批量记录API，支持最多1000条/批次
- ✅ **使用统计**: 多维度统计数据（总量、按模型、按日期）
- ✅ **成本分析**: 13个主流AI模型的精确成本计算
- ✅ **预算告警**: 自动预算检查和告警触发
- ✅ **性能优化**: 汇总数据异步更新，不阻塞主流程

### 技术特性

- **定价引擎**: 支持13个AI模型的价格计算（OpenAI、Anthropic、通义千问、百度文心、智谱ChatGLM）
- **数据库优化**: 批量插入使用`CreateInBatches`（100条/批）
- **异步处理**: 汇总计算和预算告警异步执行
- **并发安全**: 使用GORM事务保证数据一致性

### 集成组件

```
TokenMeteringService
├── PricingEngine (定价引擎)
├── TokenUsageLogRepository (日志仓储)
├── TokenUsageSummaryRepository (汇总仓储)
└── BudgetAlertService (预算告警)
```

---

## 快速开始

### 基础示例

```bash
# 1. 记录单次Token使用
curl -X POST "http://localhost:8080/api/v1/tenants/tenant-123/tokens/usage" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1001,
    "bot_id": "bot-456",
    "model_provider": "openai",
    "model_name": "gpt-4",
    "input_tokens": 1000,
    "output_tokens": 500,
    "total_tokens": 1500,
    "request_type": "chat",
    "response_time_ms": 1200,
    "is_cached": false
  }'

# 2. 获取使用统计
curl -X GET "http://localhost:8080/api/v1/tenants/tenant-123/tokens/stats?start_date=2024-01-01&end_date=2024-01-07"

# 3. 获取每日趋势
curl -X GET "http://localhost:8080/api/v1/tenants/tenant-123/tokens/daily?start_date=2024-01-01&end_date=2024-01-07"
```

---

## API端点

### 1. 记录Token使用

记录单次AI模型调用的Token使用情况和成本。

**端点**: `POST /api/v1/tenants/:tenant_id/tokens/usage`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tenant_id | string | 是 | 租户ID（URL路径参数） |
| user_id | uint64 | 否 | 用户ID |
| bot_id | string | 否 | Bot ID |
| conversation_id | string | 否 | 会话ID |
| message_id | string | 否 | 消息ID |
| model_provider | string | 是 | 模型提供商（openai, anthropic, qwen, baidu, zhipu） |
| model_name | string | 是 | 模型名称（gpt-4, claude-3-opus等） |
| model_version | string | 否 | 模型版本 |
| input_tokens | int | 是 | 输入Token数 |
| output_tokens | int | 是 | 输出Token数 |
| total_tokens | int | 是 | 总Token数（必须等于input+output） |
| request_type | string | 是 | 请求类型（chat, completion, embedding, rerank） |
| response_time_ms | int | 否 | 响应时间（毫秒） |
| latency_ms | int | 否 | 首字延迟（毫秒） |
| is_cached | bool | 否 | 是否缓存命中（默认false） |
| metadata | string | 否 | 额外元数据（JSON字符串） |

**请求示例**:

```json
{
  "user_id": 1001,
  "bot_id": "bot-456",
  "conversation_id": "conv-789",
  "message_id": "msg-001",
  "model_provider": "openai",
  "model_name": "gpt-4",
  "input_tokens": 1000,
  "output_tokens": 500,
  "total_tokens": 1500,
  "request_type": "chat",
  "response_time_ms": 1200,
  "latency_ms": 300,
  "is_cached": false
}
```

**响应示例**:

```json
{
  "log_id": 1,
  "input_cost": 0.03,
  "output_cost": 0.03,
  "total_cost": 0.06,
  "created_at": 1704067200
}
```

**错误响应**:

```json
{
  "code": "INVALID_PARAMETER",
  "message": "Invalid request parameter",
  "details": "total_tokens must equal input_tokens + output_tokens"
}
```

---

### 2. 批量记录Token使用

批量记录多次AI模型调用的Token使用情况，性能优化。

**端点**: `POST /api/v1/tenants/:tenant_id/tokens/usage/batch`

**限制**: 最多1000条/批次

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tenant_id | string | 是 | 租户ID（URL路径参数） |
| records | array | 是 | 记录数组（1-1000条） |

**请求示例**:

```json
{
  "records": [
    {
      "model_provider": "openai",
      "model_name": "gpt-4",
      "input_tokens": 1000,
      "output_tokens": 500,
      "total_tokens": 1500,
      "request_type": "chat"
    },
    {
      "model_provider": "anthropic",
      "model_name": "claude-3-opus",
      "input_tokens": 2000,
      "output_tokens": 1000,
      "total_tokens": 3000,
      "request_type": "chat"
    }
  ]
}
```

**响应示例**:

```json
{
  "success_count": 2,
  "failed_count": 0,
  "log_ids": [1, 2],
  "total_cost": 0.54
}
```

**性能指标**:
- 批量记录1000条: < 1秒
- 批量插入批次大小: 100条/批

---

### 3. 获取使用统计

获取租户的Token使用统计汇总数据。

**端点**: `GET /api/v1/tenants/:tenant_id/tokens/stats`

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tenant_id | string | 是 | 租户ID（URL路径参数） |
| bot_id | string | 否 | Bot ID（过滤特定Bot） |
| model_name | string | 否 | 模型名称（过滤特定模型） |
| start_date | string | 否 | 开始日期（格式: 2006-01-02，默认: 7天前） |
| end_date | string | 否 | 结束日期（格式: 2006-01-02，默认: 今天） |
| granularity | string | 否 | 粒度（daily/hourly，默认: daily） |

**响应示例**:

```json
{
  "tenant_id": "tenant-123",
  "total_input_tokens": 100000,
  "total_output_tokens": 50000,
  "total_tokens": 150000,
  "total_cost": 9.5,
  "total_requests": 150,
  "cached_requests": 30,
  "avg_response_time_ms": 850.5,
  "start_date": "2024-01-01",
  "end_date": "2024-01-07"
}
```

---

### 4. 获取每日使用趋势

获取租户每日的Token使用趋势数据。

**端点**: `GET /api/v1/tenants/:tenant_id/tokens/daily`

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tenant_id | string | 是 | 租户ID（URL路径参数） |
| bot_id | string | 否 | Bot ID（过滤特定Bot） |
| start_date | string | 否 | 开始日期（默认: 30天前） |
| end_date | string | 否 | 结束日期（默认: 今天） |

**响应示例**:

```json
{
  "daily_stats": [
    {
      "date": "2024-01-01",
      "total_tokens": 15000,
      "total_cost": 0.95,
      "total_requests": 15,
      "cached_requests": 3,
      "avg_response_time_ms": 820.0
    },
    {
      "date": "2024-01-02",
      "total_tokens": 22000,
      "total_cost": 1.38,
      "total_requests": 22,
      "cached_requests": 5,
      "avg_response_time_ms": 900.5
    }
  ],
  "total": {
    "total_tokens": 37000,
    "total_cost": 2.33,
    "total_requests": 37,
    "top_model": "2024-01-02",
    "top_model_cost": 1.38
  }
}
```

---

### 5. 获取模型使用统计

按模型维度统计Token使用情况和成本分布。

**端点**: `GET /api/v1/tenants/:tenant_id/tokens/models`

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tenant_id | string | 是 | 租户ID（URL路径参数） |
| start_date | string | 否 | 开始日期（默认: 7天前） |
| end_date | string | 否 | 结束日期（默认: 今天） |

**响应示例**:

```json
{
  "model_stats": [
    {
      "model_provider": "openai",
      "model_name": "gpt-4",
      "total_tokens": 150000,
      "total_cost": 9.0,
      "total_requests": 100,
      "avg_cost_per_1k_token": 0.06,
      "cost_percentage": 75.0
    },
    {
      "model_provider": "anthropic",
      "model_name": "claude-3-opus",
      "total_tokens": 50000,
      "total_cost": 3.0,
      "total_requests": 50,
      "avg_cost_per_1k_token": 0.06,
      "cost_percentage": 25.0
    }
  ],
  "total": {
    "total_tokens": 200000,
    "total_cost": 12.0,
    "total_requests": 150,
    "top_model": "gpt-4",
    "top_model_cost": 9.0
  }
}
```

---

## 数据模型

### RecordTokenUsageRequest

```typescript
{
  tenant_id: string;              // 必填
  user_id?: number;               // 可选
  bot_id?: string;                // 可选
  conversation_id?: string;       // 可选
  message_id?: string;            // 可选
  model_provider: string;         // 必填: openai, anthropic, qwen, baidu, zhipu
  model_name: string;             // 必填: gpt-4, claude-3-opus, qwen-max等
  model_version?: string;         // 可选
  input_tokens: number;           // 必填, >= 0
  output_tokens: number;          // 必填, >= 0
  total_tokens: number;           // 必填, = input + output
  request_type: string;           // 必填: chat, completion, embedding, rerank
  response_time_ms?: number;      // 可选
  latency_ms?: number;            // 可选
  is_cached: boolean;             // 默认false
  metadata?: string;              // 可选, JSON字符串
}
```

### RecordTokenUsageResponse

```typescript
{
  log_id: number;       // 日志ID
  input_cost: number;   // 输入成本（CNY）
  output_cost: number;  // 输出成本（CNY）
  total_cost: number;   // 总成本（CNY）
  created_at: number;   // 创建时间（Unix时间戳）
}
```

---

## 错误码

| 错误码 | HTTP状态 | 说明 |
|--------|----------|------|
| INVALID_PARAMETER | 400 | 请求参数无效 |
| INTERNAL_ERROR | 500 | 服务器内部错误 |

### 常见错误场景

**1. Token数不匹配**

```json
{
  "code": "INVALID_PARAMETER",
  "message": "Invalid request parameter",
  "details": "total_tokens must equal input_tokens + output_tokens"
}
```

**2. 无效的请求类型**

```json
{
  "code": "INVALID_PARAMETER",
  "message": "Invalid request parameter",
  "details": "invalid request_type: invalid_type"
}
```

**3. 批量超过限制**

```json
{
  "code": "INVALID_PARAMETER",
  "message": "Invalid request parameter",
  "details": "batch size exceeds limit of 1000"
}
```

---

## 性能指标

### 单次记录

- **响应时间**: < 100ms
- **数据库事务**: 1次
- **异步任务**: 2个（汇总更新、预算检查）

### 批量记录

- **1000条记录**: < 1秒
- **批处理大小**: 100条/批
- **数据库事务**: 1次（批量插入）

### 统计查询

- **使用统计**: < 200ms（7天数据）
- **每日趋势**: < 300ms（30天数据）
- **模型统计**: < 500ms（7天数据，按模型聚合）

### 数据库优化

- **索引**:
  - `idx_tenant_created`: (tenant_id, created_at)
  - `idx_bot_created`: (bot_id, created_at)
  - `idx_model`: (model_provider, model_name)
  - `idx_cost`: (total_cost)

---

## 使用示例

### 前端集成（React + TypeScript）

```typescript
import { useMutation, useQuery } from '@tanstack/react-query';

// 1. 记录Token使用
const recordTokenUsage = async (data: RecordTokenUsageRequest) => {
  const response = await fetch(
    `/api/v1/tenants/${tenantId}/tokens/usage`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    }
  );
  return response.json();
};

// 在组件中使用
function ChatComponent() {
  const recordMutation = useMutation({
    mutationFn: recordTokenUsage,
    onSuccess: (data) => {
      console.log('Token usage recorded:', data.total_cost);
    },
  });

  const handleSendMessage = async () => {
    // ... 调用AI模型
    const response = await callAIModel(message);

    // 记录Token使用
    recordMutation.mutate({
      bot_id: botId,
      model_provider: 'openai',
      model_name: 'gpt-4',
      input_tokens: response.usage.prompt_tokens,
      output_tokens: response.usage.completion_tokens,
      total_tokens: response.usage.total_tokens,
      request_type: 'chat',
      response_time_ms: response.response_time,
    });
  };

  return <div>...</div>;
}

// 2. 获取使用统计
function UsageStats() {
  const { data, isLoading } = useQuery({
    queryKey: ['usageStats', tenantId],
    queryFn: async () => {
      const response = await fetch(
        `/api/v1/tenants/${tenantId}/tokens/stats?start_date=2024-01-01&end_date=2024-01-07`
      );
      return response.json();
    },
  });

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      <h3>总成本: ¥{data.total_cost.toFixed(2)}</h3>
      <p>总Token: {data.total_tokens.toLocaleString()}</p>
      <p>总请求数: {data.total_requests}</p>
    </div>
  );
}
```

### 后端集成（Go）

```go
// 1. 在AI模型调用后记录Token使用
func (s *ChatService) SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error) {
    // 调用AI模型
    aiResp, err := s.aiClient.Chat(ctx, req.Message)
    if err != nil {
        return nil, err
    }

    // 记录Token使用
    _, err = s.tokenMeteringSvc.RecordTokenUsage(ctx, &service.RecordTokenUsageRequest{
        TenantID:       req.TenantID,
        UserID:         &req.UserID,
        BotID:          &req.BotID,
        ConversationID: &req.ConversationID,
        MessageID:      &resp.MessageID,
        ModelProvider:  "openai",
        ModelName:      "gpt-4",
        InputTokens:    aiResp.Usage.PromptTokens,
        OutputTokens:   aiResp.Usage.CompletionTokens,
        TotalTokens:    aiResp.Usage.TotalTokens,
        RequestType:    entity.RequestTypeChat,
        ResponseTimeMs: &aiResp.ResponseTimeMs,
    })

    if err != nil {
        // 记录失败不影响主流程
        logrus.WithError(err).Warn("[Chat] Failed to record token usage")
    }

    return resp, nil
}
```

---

## 最佳实践

### 1. 批量记录优化

对于高频场景（如批量处理），使用批量API：

```typescript
// 不推荐：循环调用单次记录
for (const item of items) {
  await recordTokenUsage(item); // N次网络请求
}

// 推荐：批量记录
const records = items.map(item => ({
  model_provider: "openai",
  model_name: "gpt-4",
  input_tokens: item.input,
  output_tokens: item.output,
  total_tokens: item.input + item.output,
  request_type: "chat",
}));
await batchRecordTokenUsage({ records }); // 1次网络请求
```

### 2. 错误处理

Token记录失败不应影响主流程：

```go
// 记录失败仅记录日志，不返回错误
if err := s.tokenMeteringSvc.RecordTokenUsage(ctx, req); err != nil {
    logrus.WithError(err).Warn("[Chat] Failed to record token usage")
    // 继续处理，不影响用户
}
```

### 3. 异步汇总

汇总数据异步更新，确保响应速度：

```go
// 主流程：保存日志（同步）
if err := s.usageLogRepo.Create(ctx, log); err != nil {
    return err
}

// 异步：更新汇总（不阻塞）
go s.updateSummaryAsync(context.Background(), req.TenantID, req.BotID, log)
```

### 4. 预算告警

每次记录后自动检查预算，触发告警：

```go
// 异步检查预算
go func() {
    if err := s.budgetAlertSvc.CheckBudget(context.Background(), req.TenantID); err != nil {
        s.logger.WithError(err).Warn("[TokenMetering] Failed to check budget")
    }
}()
```

### 5. 数据一致性

使用数据库事务保证原子性：

```go
err := s.db.Transaction(func(tx *gorm.DB) error {
    // 1. 保存日志
    if err := s.usageLogRepo.Create(ctx, log); err != nil {
        return err
    }

    // 2. 其他相关操作
    // ...

    return nil
})
```

---

## 支持的AI模型和价格

### OpenAI

| 模型 | 输入价格（¥/1K tokens） | 输出价格（¥/1K tokens） |
|------|------------------------|------------------------|
| gpt-4 | ¥0.03 | ¥0.06 |
| gpt-4-32k | ¥0.06 | ¥0.12 |
| gpt-3.5-turbo | ¥0.003 | ¥0.006 |
| gpt-3.5-turbo-16k | ¥0.004 | ¥0.008 |

### Anthropic

| 模型 | 输入价格（¥/1K tokens） | 输出价格（¥/1K tokens） |
|------|------------------------|------------------------|
| claude-3-opus | ¥0.09 | ¥0.27 |
| claude-3-sonnet | ¥0.015 | ¥0.045 |
| claude-3-haiku | ¥0.0025 | ¥0.0125 |

### 通义千问

| 模型 | 输入价格（¥/1K tokens） | 输出价格（¥/1K tokens） |
|------|------------------------|------------------------|
| qwen-max | ¥0.02 | ¥0.06 |
| qwen-plus | ¥0.004 | ¥0.012 |
| qwen-turbo | ¥0.001 | ¥0.002 |

### 百度文心

| 模型 | 输入价格（¥/1K tokens） | 输出价格（¥/1K tokens） |
|------|------------------------|------------------------|
| ernie-bot-4 | ¥0.012 | ¥0.012 |
| ernie-bot-turbo | ¥0.008 | ¥0.008 |

### 智谱ChatGLM

| 模型 | 输入价格（¥/1K tokens） | 输出价格（¥/1K tokens） |
|------|------------------------|------------------------|
| chatglm-turbo | ¥0.005 | ¥0.005 |
| chatglm-pro | ¥0.01 | ¥0.01 |

---

## 总结

Token计量API提供了完整的AI模型使用追踪和成本计算功能，核心特性：

✅ **精确计费**: 支持13个主流AI模型的价格计算
✅ **高性能**: 批量记录1000条 < 1秒
✅ **异步优化**: 汇总和告警异步执行，不阻塞主流程
✅ **多维度统计**: 总量、按模型、按日期等多维度数据
✅ **预算告警**: 自动预算检查和告警触发
✅ **企业级**: 事务保证、错误处理、日志完善

立即开始使用，轻松管理AI成本！
