# 预算管理API文档

## 概述

预算管理API提供完整的预算配置、使用监控、告警历史查询等功能，帮助企业客户有效控制AI成本。

**基础路径**: `/api/v1/tenants/:tenant_id/budget`

---

## 目录

- [创建预算配置](#1-创建预算配置)
- [获取预算配置](#2-获取预算配置)
- [更新预算配置](#3-更新预算配置)
- [删除预算配置](#4-删除预算配置)
- [获取预算使用情况](#5-获取预算使用情况)
- [手动触发预算检查](#6-手动触发预算检查)
- [获取告警历史](#7-获取告警历史)
- [数据模型](#数据模型)
- [错误码](#错误码)
- [使用示例](#使用示例)

---

## API端点

### 1. 创建预算配置

创建新的预算配置。

**端点**: `POST /api/v1/tenants/:tenant_id/budget`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tenant_id | string | 是 | 租户ID（路径参数） |
| budget_type | string | 是 | 预算类型：`monthly`, `quarterly`, `yearly` |
| budget_amount | float64 | 是 | 预算金额（CNY） |
| currency | string | 是 | 货币代码：`CNY`, `USD` |
| alert_threshold_1 | int | 是 | 一级告警阈值（%），如：80 |
| alert_threshold_2 | int | 是 | 二级告警阈值（%），如：95 |
| hard_cap_enabled | bool | 是 | 是否启用硬性上限 |
| hard_cap_amount | float64 | 否 | 硬性上限金额 |
| auto_downgrade_enabled | bool | 是 | 是否启用自动降级 |
| downgrade_threshold | int | 是 | 降级阈值（%） |
| original_model | string | 否 | 原始模型名称 |
| fallback_model | string | 否 | 降级模型名称 |
| notification_channels | []string | 是 | 通知渠道：`email`, `webhook`, `sms` |
| notification_recipients | []string | 是 | 通知接收人邮箱列表 |

**请求示例**:

```json
{
  "tenant_id": "tenant-001",
  "budget_type": "monthly",
  "budget_amount": 10000.00,
  "currency": "CNY",
  "alert_threshold_1": 80,
  "alert_threshold_2": 95,
  "hard_cap_enabled": true,
  "hard_cap_amount": 12000.00,
  "auto_downgrade_enabled": false,
  "downgrade_threshold": 90,
  "notification_channels": ["email", "webhook"],
  "notification_recipients": ["admin@example.com"]
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Budget created successfully",
  "data": {
    "tenant_id": "tenant-001",
    "budget_type": "monthly",
    "budget_amount": 10000.00,
    "currency": "CNY",
    "alert_threshold_1": 80,
    "alert_threshold_2": 95,
    "hard_cap_enabled": true,
    "hard_cap_amount": 12000.00,
    "auto_downgrade_enabled": false,
    "downgrade_threshold": 90,
    "notification_channels": ["email", "webhook"],
    "notification_recipients": ["admin@example.com"],
    "created_at": 1704067200,
    "updated_at": 1704067200
  }
}
```

**状态码**:
- `201 Created`: 创建成功
- `400 Bad Request`: 请求参数错误
- `409 Conflict`: 预算配置已存在
- `500 Internal Server Error`: 服务器错误

---

### 2. 获取预算配置

获取指定租户的预算配置。

**端点**: `GET /api/v1/tenants/:tenant_id/budget`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tenant_id | string | 是 | 租户ID（路径参数） |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "tenant_id": "tenant-001",
    "budget_type": "monthly",
    "budget_amount": 10000.00,
    "currency": "CNY",
    "alert_threshold_1": 80,
    "alert_threshold_2": 95,
    "hard_cap_enabled": true,
    "hard_cap_amount": 12000.00,
    "auto_downgrade_enabled": false,
    "downgrade_threshold": 90,
    "notification_channels": ["email", "webhook"],
    "notification_recipients": ["admin@example.com"],
    "created_at": 1704067200,
    "updated_at": 1704067200
  }
}
```

**状态码**:
- `200 OK`: 成功
- `404 Not Found`: 预算配置不存在
- `500 Internal Server Error`: 服务器错误

---

### 3. 更新预算配置

更新现有的预算配置。只更新提供的字段，未提供的字段保持不变。

**端点**: `PUT /api/v1/tenants/:tenant_id/budget`

**请求参数**:

所有字段都是可选的，只更新提供的字段。

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tenant_id | string | 是 | 租户ID（路径参数） |
| budget_amount | float64 | 否 | 新的预算金额 |
| alert_threshold_1 | int | 否 | 一级告警阈值 |
| alert_threshold_2 | int | 否 | 二级告警阈值 |
| hard_cap_enabled | bool | 否 | 是否启用硬性上限 |
| hard_cap_amount | float64 | 否 | 硬性上限金额 |
| auto_downgrade_enabled | bool | 否 | 是否启用自动降级 |
| downgrade_threshold | int | 否 | 降级阈值 |
| original_model | string | 否 | 原始模型名称 |
| fallback_model | string | 否 | 降级模型名称 |
| notification_channels | []string | 否 | 通知渠道 |
| notification_recipients | []string | 否 | 通知接收人 |

**请求示例**:

```json
{
  "budget_amount": 15000.00,
  "alert_threshold_1": 75,
  "alert_threshold_2": 90
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Budget updated successfully",
  "data": {
    "tenant_id": "tenant-001",
    "budget_type": "monthly",
    "budget_amount": 15000.00,
    "currency": "CNY",
    "alert_threshold_1": 75,
    "alert_threshold_2": 90,
    "hard_cap_enabled": true,
    "hard_cap_amount": 12000.00,
    "auto_downgrade_enabled": false,
    "downgrade_threshold": 90,
    "notification_channels": ["email", "webhook"],
    "notification_recipients": ["admin@example.com"],
    "created_at": 1704067200,
    "updated_at": 1704153600
  }
}
```

**状态码**:
- `200 OK`: 更新成功
- `400 Bad Request`: 请求参数错误
- `404 Not Found`: 预算配置不存在
- `500 Internal Server Error`: 服务器错误

---

### 4. 删除预算配置

删除预算配置。

**端点**: `DELETE /api/v1/tenants/:tenant_id/budget`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tenant_id | string | 是 | 租户ID（路径参数） |

**响应示例**:

```json
{
  "code": 0,
  "message": "Budget deleted successfully"
}
```

**状态码**:
- `200 OK`: 删除成功
- `404 Not Found`: 预算配置不存在
- `500 Internal Server Error`: 服务器错误

---

### 5. 获取预算使用情况

获取当前预算周期的使用情况统计。

**端点**: `GET /api/v1/tenants/:tenant_id/budget/usage`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tenant_id | string | 是 | 租户ID（路径参数） |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "budget_amount": 10000.00,
    "used_amount": 5200.50,
    "remaining_amount": 4799.50,
    "usage_percent": 52.0,
    "period_start": 1704067200,
    "period_end": 1706745600,
    "total_tokens": 5200000,
    "total_requests": 500,
    "average_cost": 0.001,
    "alert_count": 1,
    "last_alert_at": 1704500000,
    "last_alert_level": "warning",
    "estimated_daily_usage": 173.35,
    "days_remaining": 15,
    "will_exceed_budget": false,
    "updated_at": 1704500000
  }
}
```

**字段说明**:

| 字段 | 说明 |
|------|------|
| budget_amount | 预算金额 |
| used_amount | 已使用金额 |
| remaining_amount | 剩余金额 |
| usage_percent | 使用率（%） |
| period_start | 当前周期开始时间（Unix时间戳） |
| period_end | 当前周期结束时间（Unix时间戳） |
| total_tokens | 总Token数 |
| total_requests | 总请求数 |
| average_cost | 平均每1K Token成本 |
| alert_count | 当前周期告警数 |
| last_alert_at | 最后一次告警时间 |
| last_alert_level | 最后一次告警级别 |
| estimated_daily_usage | 预计每日使用量 |
| days_remaining | 周期剩余天数 |
| will_exceed_budget | 是否预测会超预算 |
| updated_at | 更新时间 |

**状态码**:
- `200 OK`: 成功
- `404 Not Found`: 预算配置不存在
- `500 Internal Server Error`: 服务器错误

---

### 6. 手动触发预算检查

手动触发一次预算检查，返回当前使用状态和告警信息。

**端点**: `POST /api/v1/tenants/:tenant_id/budget/check`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tenant_id | string | 是 | 租户ID（路径参数） |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "tenant_id": "tenant-001",
    "budget_amount": 10000.00,
    "used_amount": 8500.00,
    "usage_percent": 85.0,
    "alert_triggered": true,
    "alert_level": "critical",
    "message": "Critical: Usage at 85.0% (≥95%)",
    "checked_at": 1704500000
  }
}
```

**字段说明**:

| 字段 | 说明 |
|------|------|
| alert_triggered | 是否触发告警 |
| alert_level | 告警级别：`warning`, `critical`, `emergency` |
| message | 告警消息 |
| checked_at | 检查时间 |

**状态码**:
- `200 OK`: 成功
- `500 Internal Server Error`: 服务器错误

---

### 7. 获取告警历史

获取预算告警历史记录，支持过滤和分页。

**端点**: `GET /api/v1/tenants/:tenant_id/budget/alerts`

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tenant_id | string | 是 | 租户ID（路径参数） |
| alert_type | string | 否 | 告警类型：`threshold_1`, `threshold_2`, `hard_cap`, `downgrade` |
| alert_level | string | 否 | 告警级别：`warning`, `critical`, `emergency` |
| notification_status | string | 否 | 通知状态：`pending`, `sent`, `failed` |
| start_date | string | 否 | 开始日期（RFC3339格式） |
| end_date | string | 否 | 结束日期（RFC3339格式） |
| page_token | string | 否 | 分页token |
| page_size | int | 否 | 每页数量（默认20，最大100） |

**请求示例**:

```
GET /api/v1/tenants/tenant-001/budget/alerts?alert_level=critical&page_size=20
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "tenant_id": "tenant-001",
      "alert_type": "threshold_2",
      "budget_amount": 10000.00,
      "used_amount": 9500.00,
      "usage_percent": 95.0,
      "alert_level": "critical",
      "alert_message": "🚨 **预算告警**\n\n租户ID: tenant-001\n预算金额: ¥10000.00\n已使用: ¥9500.00 (95.0%)\n\n警告: 即将达到预算上限，建议立即充值或调整预算配置。",
      "notification_channels": ["email"],
      "notification_status": "sent",
      "sent_at": 1704500000,
      "error_message": null,
      "created_at": 1704500000
    }
  ],
  "pagination": {
    "next_page_token": "eyJpZCI6MjB9",
    "total_count": 45,
    "has_more": true
  }
}
```

**状态码**:
- `200 OK`: 成功
- `400 Bad Request`: 请求参数错误
- `500 Internal Server Error`: 服务器错误

---

## 数据模型

### BudgetSettingsDTO 预算配置

```typescript
{
  tenant_id: string;              // 租户ID
  budget_type: string;            // 预算类型: monthly, quarterly, yearly
  budget_amount: number;          // 预算金额
  currency: string;               // 货币代码

  // 告警配置
  alert_threshold_1: number;      // 一级告警阈值（%）
  alert_threshold_2: number;      // 二级告警阈值（%）
  hard_cap_enabled: boolean;      // 是否启用硬性上限
  hard_cap_amount?: number;       // 硬性上限金额

  // 降级配置
  auto_downgrade_enabled: boolean; // 是否启用自动降级
  downgrade_threshold: number;     // 降级阈值（%）
  original_model?: string;         // 原始模型
  fallback_model?: string;         // 降级模型

  // 通知配置
  notification_channels: string[];   // 通知渠道
  notification_recipients: string[]; // 通知接收人

  created_at: number;              // 创建时间（Unix时间戳）
  updated_at: number;              // 更新时间（Unix时间戳）
}
```

### BudgetUsageDTO 预算使用情况

```typescript
{
  // 预算信息
  budget_amount: number;         // 预算金额
  used_amount: number;           // 已使用金额
  remaining_amount: number;      // 剩余金额
  usage_percent: number;         // 使用率（%）

  // 当前周期统计
  period_start: number;          // 周期开始时间
  period_end: number;            // 周期结束时间
  total_tokens: number;          // 总Token数
  total_requests: number;        // 总请求数
  average_cost: number;          // 平均每1K Token成本

  // 告警统计
  alert_count: number;           // 当前周期告警数
  last_alert_at?: number;        // 最后告警时间
  last_alert_level?: string;     // 最后告警级别

  // 预测
  estimated_daily_usage: number; // 预计每日使用量
  days_remaining: number;        // 剩余天数
  will_exceed_budget: boolean;   // 是否会超预算

  updated_at: number;            // 更新时间
}
```

### BudgetAlertDTO 告警记录

```typescript
{
  id: number;                    // 告警ID
  tenant_id: string;             // 租户ID
  alert_type: string;            // 告警类型
  budget_amount: number;         // 预算金额
  used_amount: number;           // 已使用金额
  usage_percent: number;         // 使用率（%）
  alert_level: string;           // 告警级别
  alert_message?: string;        // 告警消息
  notification_channels: string[]; // 通知渠道
  notification_status: string;   // 通知状态
  sent_at?: number;              // 发送时间
  error_message?: string;        // 错误信息
  created_at: number;            // 创建时间
}
```

---

## 错误码

| 错误码 | 说明 | 处理建议 |
|--------|------|----------|
| 0 | 成功 | - |
| 400 | 请求参数错误 | 检查请求参数格式和必填字段 |
| 404 | 资源不存在 | 检查租户ID是否正确 |
| 409 | 资源冲突 | 预算配置已存在，请使用更新接口 |
| 500 | 服务器错误 | 联系技术支持 |

---

## 预算类型说明

### Monthly（月度预算）

- **周期**: 每月1日重置
- **适用场景**: 常规使用，月度成本控制
- **示例**: 每月预算¥10,000

### Quarterly（季度预算）

- **周期**: 每季度首月（1/4/7/10月）1日重置
- **适用场景**: 项目制、季度计划
- **示例**: 每季度预算¥30,000

### Yearly（年度预算）

- **周期**: 每年1月1日重置
- **适用场景**: 年度规划、长期合同
- **示例**: 每年预算¥120,000

---

## 告警级别说明

### Warning（一级告警）

- **触发条件**: 使用率 ≥ alert_threshold_1（默认80%）
- **建议措施**: 注意控制使用量，避免超额
- **通知频率**: 每日最多1次

### Critical（二级告警）

- **触发条件**: 使用率 ≥ alert_threshold_2（默认95%）
- **建议措施**: 立即充值或调整预算配置
- **通知频率**: 每日最多1次

### Emergency（紧急告警）

- **触发条件**: 达到硬性上限
- **自动处理**:
  - 如果启用自动降级：切换到降级模型
  - 如果未启用降级：暂停服务
- **通知频率**: 立即发送

---

## 使用示例

### 场景1: 为新租户创建月度预算

```bash
curl -X POST "http://api.example.com/api/v1/tenants/tenant-001/budget" \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "tenant-001",
    "budget_type": "monthly",
    "budget_amount": 10000.00,
    "currency": "CNY",
    "alert_threshold_1": 80,
    "alert_threshold_2": 95,
    "hard_cap_enabled": true,
    "hard_cap_amount": 12000.00,
    "auto_downgrade_enabled": false,
    "downgrade_threshold": 90,
    "notification_channels": ["email"],
    "notification_recipients": ["admin@example.com"]
  }'
```

### 场景2: 查看当前使用情况

```bash
curl -X GET "http://api.example.com/api/v1/tenants/tenant-001/budget/usage"
```

### 场景3: 手动触发预算检查

```bash
curl -X POST "http://api.example.com/api/v1/tenants/tenant-001/budget/check"
```

### 场景4: 查询近30天的所有告警

```bash
curl -X GET "http://api.example.com/api/v1/tenants/tenant-001/budget/alerts?start_date=2024-01-01T00:00:00Z&end_date=2024-01-30T23:59:59Z&page_size=50"
```

### 场景5: 调整预算金额

```bash
curl -X PUT "http://api.example.com/api/v1/tenants/tenant-001/budget" \
  -H "Content-Type: application/json" \
  -d '{
    "budget_amount": 15000.00
  }'
```

---

## 最佳实践

### 1. 合理设置告警阈值

建议的阈值配置：
- **一级告警**: 70-80%（提前预警）
- **二级告警**: 90-95%（紧急处理）
- **硬性上限**: 100-120%（容错空间）

### 2. 启用自动降级

对于对成本敏感的场景，建议启用自动降级：
```json
{
  "auto_downgrade_enabled": true,
  "downgrade_threshold": 90,
  "original_model": "gpt-4",
  "fallback_model": "gpt-3.5-turbo"
}
```

### 3. 多渠道通知

配置多个通知渠道确保告警及时送达：
```json
{
  "notification_channels": ["email", "webhook", "sms"],
  "notification_recipients": ["admin@example.com", "ops@example.com"]
}
```

### 4. 定期检查

建议每小时自动执行预算检查：
```bash
# Crontab: 0 * * * *
curl -X POST "http://api.example.com/api/v1/tenants/tenant-001/budget/check"
```

### 5. 监控使用趋势

使用 `GetBudgetUsage` API 监控使用趋势，及时调整策略：
```bash
curl -X GET "http://api.example.com/api/v1/tenants/tenant-001/budget/usage"
```

---

## 集成说明

### 与BudgetAlertService的集成

`BudgetManagementService` 内部集成了 `BudgetAlertService`，提供以下功能：

1. **自动预算检查**: `BudgetManagementService.CheckBudget()` 调用 `BudgetAlertService.CheckBudget()` 进行预算检查
2. **告警发送**: 检测到超阈值时自动发送告警通知
3. **告警历史**: 所有告警记录存储在 `budget_alerts_history` 表中

### 与TokenMetering的集成

预算使用数据来源于 `TokenMetering` 系统：
1. **实时记录**: 每次API调用记录Token使用量
2. **定期汇总**: 按天/小时汇总使用数据到 `token_usage_summary` 表
3. **成本计算**: 根据模型定价计算成本
4. **预算监控**: 定期检查预算使用率并触发告警

---

## 附录

### 相关文档

- [Token Metering API文档](./TOKEN_METERING_API.md)
- [预算告警服务文档](./BUDGET_ALERT_SERVICE.md)
- [企业级开发规范手册](../../企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)

### 技术支持

如有问题，请联系技术支持团队或查阅相关文档。
