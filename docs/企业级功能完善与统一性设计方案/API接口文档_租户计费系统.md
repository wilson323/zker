# API接口文档：租户计费系统模块

**模块名称**: 租户计费系统 (TenantBilling)
**设计文档**: 23-MultiTenant_SaaS核心_租户计费系统.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 订阅计划API](#2-订阅计划api)
- [3. 订阅管理API](#3-订阅管理api)
- [4. 账单管理API](#4-账单管理api)
- [5. 使用计量API](#5-使用计量api)
- [6. 支付管理API](#6-支付管理api)
- [7. 数据模型](#7-数据模型)
- [8. 错误码定义](#8-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

租户计费系统是MultiTenant SaaS的核心盈利模块，提供灵活的计费模式和自动化账单管理：

- ✅ **多计费模式** - 订阅制、按量计费、混合计费
- ✅ **订阅管理** - 月付/年付、套餐升级/降级、自动续费
- ✅ **账单管理** - 自动出账、账单审核、在线支付
- ✅ **使用计量** - API调用、Token消耗、存储空间计量
- ✅ **发票管理** - 发票开具、发票抬头管理
- ✅ **预算告警** - 预算设置、超额告警、自动降级

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5
- 缓存: Redis 8.0
- 支付: 支付宝、微信支付、银行转账

### 1.3 计费模式

| 计费模式 | 说明 | 适用场景 |
|---------|------|---------|
| **订阅制** | 固定周期费用 | SaaS标准订阅 |
| **按量计费** | 按实际使用量计费 | API调用、Token消耗 |
| **混合计费** | 基础订阅 + 超额按量 | 企业弹性扩展 |

---

## 2. 订阅计划API

### 2.1 获取订阅计划列表

**接口地址**: `GET /api/v1/billing/plans`

**功能说明**: 获取所有可用的订阅计划

**查询参数**:
- type: 计划类型过滤 (free/professional/enterprise)
- include_inactive: 是否包含已停用计划 (默认false)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "plan-free",
        "name": "免费版",
        "type": "free",
        "description": "适合个人用户和小团队试用",
        "prices": {
          "monthly": 0,
          "yearly": 0,
          "currency": "CNY"
        },
        "quotas": {
          "maxUsers": 5,
          "maxBots": 1,
          "maxKnowledgeBases": 1,
          "maxDocuments": 10,
          "maxStorageGB": 1,
          "maxAPICallsPerMonth": 1000,
          "maxTokensPerMonth": 100000,
          "maxConcurrentUsers": 2
        },
        "features": [
          "basic_bot",
          "knowledge_base"
        ],
        "isActive": true
      },
      {
        "id": "plan-professional",
        "name": "专业版",
        "type": "professional",
        "description": "适合成长型企业和团队",
        "prices": {
          "monthly": 299,
          "yearly": 2990,
          "currency": "CNY",
          "yearlyDiscount": 0.17
        },
        "quotas": {
          "maxUsers": 50,
          "maxBots": 10,
          "maxKnowledgeBases": 5,
          "maxDocuments": 500,
          "maxStorageGB": 10,
          "maxAPICallsPerMonth": 100000,
          "maxTokensPerMonth": 10000000,
          "maxConcurrentUsers": 10
        },
        "features": [
          "basic_bot",
          "knowledge_base",
          "api_access",
          "custom_domain",
          "priority_support"
        ],
        "isActive": true
      },
      {
        "id": "plan-enterprise",
        "name": "企业版",
        "type": "enterprise",
        "description": "适合大型企业和定制需求",
        "prices": {
          "monthly": 999,
          "yearly": 9990,
          "currency": "CNY",
          "yearlyDiscount": 0.17
        },
        "quotas": {
          "maxUsers": 500,
          "maxBots": 100,
          "maxKnowledgeBases": 20,
          "maxDocuments": 10000,
          "maxStorageGB": 100,
          "maxAPICallsPerMonth": 1000000,
          "maxTokensPerMonth": 100000000,
          "maxConcurrentUsers": 50
        },
        "features": [
          "basic_bot",
          "knowledge_base",
          "api_access",
          "custom_domain",
          "priority_support",
          "sla_guarantee",
          "dedicated_manager"
        ],
        "isActive": true
      }
    ],
    "total": 3
  }
}
```

### 2.2 获取计划详情

**接口地址**: `GET /api/v1/billing/plans/{plan_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "plan-professional",
    "name": "专业版",
    "type": "professional",
    "description": "适合成长型企业和团队的完整解决方案",
    "prices": {
      "monthly": 299,
      "yearly": 2990,
      "currency": "CNY",
      "yearlyDiscount": 0.17
    },
    "quotas": {
      "maxUsers": 50,
      "maxBots": 10,
      "maxKnowledgeBases": 5,
      "maxDocuments": 500,
      "maxStorageGB": 10,
      "maxAPICallsPerMonth": 100000,
      "maxTokensPerMonth": 10000000,
      "maxConcurrentUsers": 10
    },
    "features": [
      {
        "key": "basic_bot",
        "name": "基础Bot功能",
        "description": "创建和使用智能Bot"
      },
      {
        "key": "api_access",
        "name": "API访问",
        "description": "开放API接口访问"
      }
    ],
    "isActive": true,
    "createdAt": "2025-01-01T00:00:00Z"
  }
}
```

---

## 3. 订阅管理API

### 3.1 创建订阅

**接口地址**: `POST /api/v1/billing/subscribe`

**功能说明**: 创建新的订阅

**请求参数**:
```json
{
  "plan_id": "plan-professional",
  "cycle": "yearly",
  "auto_renew": true
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| plan_id | string | 是 | 计划ID |
| cycle | string | 是 | 计费周期 (monthly/yearly) |
| auto_renew | boolean | 否 | 是否自动续费，默认true |

**响应示例**:
```json
{
  "code": 0,
  "message": "订阅创建成功",
  "data": {
    "id": "sub-001",
    "tenantId": "tenant-123",
    "planId": "plan-professional",
    "planName": "专业版",
    "cycle": "yearly",
    "startDate": "2025-01-03",
    "endDate": "2026-01-03",
    "autoRenew": true,
    "status": "active",
    "amount": 2990,
    "currency": "CNY",
    "nextBillingDate": "2026-01-03",
    "createdAt": "2025-01-03T10:00:00Z"
  }
}
```

### 3.2 获取当前订阅

**接口地址**: `GET /api/v1/billing/subscription`

**功能说明**: 获取当前租户的订阅信息

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "sub-001",
    "tenantId": "tenant-123",
    "plan": {
      "id": "plan-professional",
      "name": "专业版",
      "type": "professional"
    },
    "cycle": "yearly",
    "startDate": "2025-01-03",
    "endDate": "2026-01-03",
    "autoRenew": true,
    "status": "active",
    "trialEndDate": null,
    "discountRate": 0,
    "usage": {
      "users": 25,
      "bots": 5,
      "knowledgeBases": 2,
      "documents": 150,
      "storageGB": 3.5,
      "apiCalls": 45000,
      "tokens": 3500000
    },
    "quotas": {
      "maxUsers": 50,
      "maxBots": 10,
      "maxKnowledgeBases": 5,
      "maxDocuments": 500,
      "maxStorageGB": 10,
      "maxAPICallsPerMonth": 100000,
      "maxTokensPerMonth": 10000000
    },
    "createdAt": "2025-01-03T10:00:00Z",
    "updatedAt": "2025-01-03T10:00:00Z"
  }
}
```

### 3.3 升级/降级订阅

**接口地址**: `POST /api/v1/billing/subscription/change-plan`

**功能说明**: 更改订阅计划（升级或降级）

**请求参数**:
```json
{
  "target_plan_id": "plan-enterprise",
  "effective": "immediate"
}
```

**参数说明**:
- target_plan_id: 目标计划ID
- effective: 生效时间 (immediate/next_cycle)

**响应示例**:
```json
{
  "code": 0,
  "message": "订阅升级成功",
  "data": {
    "subscriptionId": "sub-001",
    "oldPlan": "plan-professional",
    "newPlan": "plan-enterprise",
    "effectiveDate": "2025-01-03",
    "proratedAmount": 580.50,
    "newAmount": 9990
  }
}
```

### 3.4 取消订阅

**接口地址**: `POST /api/v1/billing/subscription/cancel`

**功能说明**: 取消订阅（周期结束后不再续费）

**请求参数**:
```json
{
  "reason": "不再需要",
  "feedback": "功能不符合需求"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "订阅取消成功",
  "data": {
    "subscriptionId": "sub-001",
    "status": "cancelled",
    "endDate": "2026-01-03",
    "message": "您的订阅将在到期后停止，期间可正常使用所有功能"
  }
}
```

---

## 4. 账单管理API

### 4.1 获取账单列表

**接口地址**: `GET /api/v1/billing/invoices`

**功能说明**: 获取当前租户的账单列表

**查询参数**:
- page: 页码
- page_size: 每页数量
- status: 状态过滤 (draft/sent/paid/overdue/cancelled)
- start_date: 开始日期
- end_date: 结束日期

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 12,
    "items": [
      {
        "id": "inv-001",
        "invoiceNumber": "INV2025010001",
        "issueDate": "2025-01-03",
        "dueDate": "2025-01-17",
        "subtotal": 2990,
        "tax": 179.4,
        "discount": 0,
        "total": 3169.4,
        "paidAmount": 0,
        "status": "sent",
        "invoiceStatus": "not_issued",
        "createdAt": "2025-01-03T10:00:00Z"
      },
      {
        "id": "inv-002",
        "invoiceNumber": "INV2024120001",
        "issueDate": "2024-12-03",
        "dueDate": "2024-12-17",
        "subtotal": 2990,
        "tax": 179.4,
        "discount": 0,
        "total": 3169.4,
        "paidAmount": 3169.4,
        "status": "paid",
        "invoiceStatus": "issued",
        "paidAt": "2024-12-10T10:00:00Z"
      }
    ]
  }
}
```

### 4.2 获取账单详情

**接口地址**: `GET /api/v1/billing/invoices/{invoice_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "inv-001",
    "tenantId": "tenant-123",
    "subscriptionId": "sub-001",
    "invoiceNumber": "INV2025010001",
    "issueDate": "2025-01-03",
    "dueDate": "2025-01-17",
    "subtotal": 2990,
    "tax": 179.4,
    "discount": 0,
    "total": 3169.4,
    "paidAmount": 0,
    "status": "sent",
    "invoiceType": "company",
    "invoiceTitle": "XX科技有限公司",
    "taxNumber": "91110000XXXXXXXX",
    "invoiceStatus": "not_issued",
    "items": [
      {
        "description": "专业版年费订阅",
        "quantity": 1,
        "unitPrice": 2990,
        "amount": 2990
      }
    ],
    "createdAt": "2025-01-03T10:00:00Z"
  }
}
```

### 4.3 更新发票信息

**接口地址**: `PUT /api/v1/billing/invoices/{invoice_id}/invoice-info`

**功能说明**: 更新发票抬头信息

**请求参数**:
```json
{
  "invoice_type": "company",
  "invoice_title": "XX科技有限公司",
  "tax_number": "91110000XXXXXXXX"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "发票信息更新成功"
}
```

### 4.4 开具发票

**接口地址**: `POST /api/v1/billing/invoices/{invoice_id}/issue`

**功能说明**: 申请开具发票

**响应示例**:
```json
{
  "code": 0,
  "message": "发票申请已提交",
  "data": {
    "invoiceId": "inv-001",
    "status": "sending",
    "estimatedDays": 7
  }
}
```

---

## 5. 使用计量API

### 5.1 获取使用量统计

**接口地址**: `GET /api/v1/billing/usage`

**功能说明**: 获取当前计费周期的使用量统计

**查询参数**:
- start_date: 开始日期（可选）
- end_date: 结束日期（可选）

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": {
      "startDate": "2025-01-01",
      "endDate": "2025-01-31"
    },
    "quotas": {
      "maxUsers": 50,
      "maxBots": 10,
      "maxKnowledgeBases": 5,
      "maxDocuments": 500,
      "maxStorageGB": 10,
      "maxAPICallsPerMonth": 100000,
      "maxTokensPerMonth": 10000000
    },
    "usage": {
      "users": 25,
      "bots": 5,
      "knowledgeBases": 2,
      "documents": 150,
      "storageGB": 3.5,
      "apiCalls": 45000,
      "tokens": 3500000
    },
    "usagePercentage": {
      "users": 50,
      "bots": 50,
      "knowledgeBases": 40,
      "documents": 30,
      "storageGB": 35,
      "apiCalls": 45,
      "tokens": 35
    },
    "overages": {
      "apiCalls": 0,
      "tokens": 0,
      "estimatedAmount": 0
    }
  }
}
```

### 5.2 获取使用量明细

**接口地址**: `GET /api/v1/billing/usage/details`

**功能说明**: 获取每日使用量明细

**查询参数**:
- metric_type: 指标类型 (api_calls/tokens/storage/users)
- start_date: 开始日期
- end_date: 结束日期

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "metricType": "api_calls",
    "items": [
      {
        "date": "2025-01-01",
        "value": 1500
      },
      {
        "date": "2025-01-02",
        "value": 1800
      }
    ],
    "total": 3300
  }
}
```

### 5.3 设置预算告警

**接口地址**: `POST /api/v1/billing/budget`

**功能说明**: 设置月度预算告警阈值

**请求参数**:
```json
{
  "monthly_budget": 5000,
  "alert_threshold": 80,
  "alert_channels": ["email", "webhook"],
  "auto_downgrade": false
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "预算告警设置成功",
  "data": {
    "tenantId": "tenant-123",
    "monthlyBudget": 5000,
    "alertThreshold": 80,
    "alertChannels": ["email", "webhook"],
    "autoDowngrade": false
  }
}
```

---

## 6. 支付管理API

### 6.1 创建支付订单

**接口地址**: `POST /api/v1/billing/invoices/{invoice_id}/pay`

**功能说明**: 创建支付订单

**请求参数**:
```json
{
  "payment_method": "alipay"
}
```

**支付方式**:
- alipay: 支付宝
- wechat: 微信支付
- bank_transfer: 银行转账

**响应示例**:
```json
{
  "code": 0,
  "message": "支付订单创建成功",
  "data": {
    "orderId": "order-001",
    "invoiceId": "inv-001",
    "amount": 3169.4,
    "currency": "CNY",
    "paymentMethod": "alipay",
    "paymentUrl": "https://openapi.alipay.com/gateway.do?...",
    "qrCode": "https://qr.example.com/...",
    "expiresAt": "2025-01-03T11:00:00Z",
    "createdAt": "2025-01-03T10:00:00Z"
  }
}
```

### 6.2 查询支付状态

**接口地址**: `GET /api/v1/billing/payments/{order_id}/status`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "orderId": "order-001",
    "invoiceId": "inv-001",
    "amount": 3169.4,
    "status": "paid",
    "paidAt": "2025-01-03T10:15:00Z",
    "transactionId": "2025010322001XXXXXXXX"
  }
}
```

---

## 7. 数据模型

### 7.1 SubscriptionPlan（订阅计划）
```typescript
interface SubscriptionPlan {
  id: string;
  name: string;
  type: 'free' | 'professional' | 'enterprise';
  description?: string;
  prices: {
    monthly: number;
    yearly: number;
    currency: string;
    yearlyDiscount?: number;
  };
  quotas: {
    maxUsers: number;
    maxBots: number;
    maxKnowledgeBases: number;
    maxDocuments: number;
    maxStorageGB: number;
    maxAPICallsPerMonth: number;
    maxTokensPerMonth: number;
    maxConcurrentUsers: number;
  };
  features: string[];
  isActive: boolean;
  createdAt: string;
}
```

### 7.2 Subscription（订阅）
```typescript
interface Subscription {
  id: string;
  tenantId: string;
  planId: string;
  cycle: 'monthly' | 'yearly';
  startDate: string;
  endDate: string;
  autoRenew: boolean;
  status: 'active' | 'suspended' | 'cancelled' | 'expired';
  trialEndDate?: string;
  discountRate: number;
  createdAt: string;
  updatedAt: string;
}
```

### 7.3 Invoice（账单）
```typescript
interface Invoice {
  id: string;
  tenantId: string;
  subscriptionId?: string;
  invoiceNumber: string;
  issueDate: string;
  dueDate: string;
  subtotal: number;
  tax: number;
  discount: number;
  total: number;
  paidAmount: number;
  status: 'draft' | 'sent' | 'paid' | 'overdue' | 'cancelled';
  invoiceType?: 'personal' | 'company';
  invoiceTitle?: string;
  taxNumber?: string;
  invoiceStatus: 'not_issued' | 'issued' | 'sending';
  createdAt: string;
  updatedAt: string;
}
```

### 7.4 UsageMetric（使用计量）
```typescript
interface UsageMetric {
  id: number;
  tenantId: string;
  metricType: 'api_calls' | 'tokens' | 'storage' | 'users';
  metricDate: string;
  metricValue: number;
  createdAt: string;
}
```

---

## 8. 错误码定义

| 错误码 | HTTP状态码 | 说明 | 处理建议 |
|--------|-----------|------|----------|
| 50001 | 404 | 订阅计划不存在 | 检查计划ID |
| 50002 | 400 | 已有活跃订阅 | 请先取消现有订阅 |
| 50003 | 400 | 订阅已过期 | 请续费 |
| 50004 | 403 | 功能超出配额 | 升级订阅计划 |
| 50005 | 400 | 降级到免费版失败 | 需要先减少资源使用量 |
| 50101 | 404 | 账单不存在 | 检查账单ID |
| 50102 | 400 | 账单已支付 | 无需重复支付 |
| 50103 | 400 | 账单已过期 | 重新生成账单 |
| 50201 | 400 | 支付订单已过期 | 重新创建支付订单 |
| 50202 | 400 | 支付金额不匹配 | 联系客服 |
| 50203 | 400 | 支付方式不支持 | 选择其他支付方式 |
| 50301 | 400 | 超出预算阈值 | 考虑升级计划或减少使用量 |
| 50302 | 400 | 自动降级已触发 | 系统已自动降级服务 |

---

**文档结束**
