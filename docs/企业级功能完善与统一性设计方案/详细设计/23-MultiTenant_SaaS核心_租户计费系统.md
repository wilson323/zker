# 23-MultiTenant_SaaS核心_租户计费系统 详细设计说明书

**文档编号**: DE-DD-2025-023
**模块名称**: 租户计费系统 (TenantBilling)
**版本**: v1.0.0
**作者**: ZKER Enterprise Team
**创建日期**: 2025-01-03

---

## 1. 模块概述

**租户计费系统** 是 MultiTenant SaaS 的核心盈利模块，通过**灵活计费模式 + 自动化账单管理**，实现订阅制、按量计费、混合计费等多种计费模式。

**核心设计理念**：
- ✅ **多计费模式**：订阅制、按量计费、阶梯计费
- ✅ **自动化**：自动计量、自动出账、自动扣款
- ✅ **企业级**：支持发票管理、账期管理、欠费处理

**实现策略**：✅ 40% 编码（计费引擎） + 60% 配置（计费规则）

---

## 2. 核心功能

### 2.1 计费模式

**F1 - 订阅制计费**
- F1.1 按月订阅
- F1.2 按年订阅（享受折扣）
- F1.3 套餐管理（免费版、专业版、企业版）

**F2 - 按量计费**
- F2.1 API调用计费
- F2.2 Token消耗计费
- F2.3 存储空间计费

**F3 - 混合计费**
- F3.1 基础订阅 + 超额按量
- F3.2 阶梯定价

### 2.2 账单管理

**F4 - 账单生命周期**
- F4.1 自动生成账单
- F4.2 账单审核
- F4.3 账单通知
- F4.4 在线支付
- F4.5 发票开具

---

## 3. 数据库设计

### 3.1 核心表结构

#### 3.1.1 订阅计划表 (subscription_plans)

```sql
CREATE TABLE subscription_plans (
    id VARCHAR(64) PRIMARY KEY COMMENT '计划ID',
    name VARCHAR(50) NOT NULL COMMENT '计划名称',
    description VARCHAR(500) COMMENT '计划描述',
    type ENUM('free', 'professional', 'enterprise') NOT NULL,

    -- 定价
    monthly_price DECIMAL(10,2) NOT NULL COMMENT '月付价格',
    yearly_price DECIMAL(10,2) NOT NULL COMMENT '年付价格',
    currency VARCHAR(3) DEFAULT 'CNY',

    -- 配额
    quotas JSON NOT NULL COMMENT '配额定义 {"maxUsers":5,"maxBots":10,...}',

    -- 功能
    features JSON COMMENT '功能列表 ["custom_domain","api_access",...]',

    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uk_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='订阅计划表';
```

**配额定义 JSON 示例**：

```json
{
  "maxUsers": 5,
  "maxBots": 10,
  "maxKnowledgeBases": 3,
  "maxDocuments": 100,
  "maxStorageGB": 10,
  "maxAPICallsPerMonth": 10000,
  "maxTokensPerMonth": 1000000,
  "maxConcurrentUsers": 2
}
```

#### 3.1.2 订阅表 (subscriptions)

```sql
CREATE TABLE subscriptions (
    id VARCHAR(64) PRIMARY KEY COMMENT '订阅ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    plan_id VARCHAR(64) NOT NULL COMMENT '计划ID',

    -- 订阅周期
    cycle ENUM('monthly', 'yearly') NOT NULL COMMENT '计费周期',
    start_date DATE NOT NULL COMMENT '开始日期',
    end_date DATE NOT NULL COMMENT '结束日期',
    auto_renew BOOLEAN DEFAULT TRUE COMMENT '是否自动续费',

    -- 状态
    status ENUM('active', 'suspended', 'cancelled', 'expired') DEFAULT 'active',

    -- 优惠
    discount_rate DECIMAL(5,4) DEFAULT 0.0000 COMMENT '折扣率',
    trial_end_date DATE COMMENT '试用期结束日期',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status),
    INDEX idx_end_date (end_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='订阅表';
```

#### 3.1.3 账单表 (invoices)

```sql
CREATE TABLE invoices (
    id VARCHAR(64) PRIMARY KEY COMMENT '账单ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    subscription_id VARCHAR(64) COMMENT '订阅ID',

    -- 账单信息
    invoice_number VARCHAR(50) NOT NULL UNIQUE COMMENT '账单编号',
    issue_date DATE NOT NULL COMMENT '开票日期',
    due_date DATE NOT NULL COMMENT '应付日期',

    -- 金额
    subtotal DECIMAL(12,2) NOT NULL COMMENT '小计',
    tax DECIMAL(12,2) DEFAULT 0.00 COMMENT '税额',
    discount DECIMAL(12,2) DEFAULT 0.00 COMMENT '折扣',
    total DECIMAL(12,2) NOT NULL COMMENT '总计',
    paid_amount DECIMAL(12,2) DEFAULT 0.00 COMMENT '已付金额',

    -- 状态
    status ENUM('draft', 'sent', 'paid', 'overdue', 'cancelled') DEFAULT 'draft',

    -- 发票
    invoice_type ENUM('personal', 'company') COMMENT '发票类型',
    invoice_title VARCHAR(200) COMMENT '发票抬头',
    tax_number VARCHAR(50) COMMENT '税号',
    invoice_status ENUM('not_issued', 'issued', 'sending') DEFAULT 'not_issued',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status),
    INDEX idx_due_date (due_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='账单表';
```

#### 3.1.4 使用计量表 (usage_metrics)

```sql
CREATE TABLE usage_metrics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '计量ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    metric_type ENUM('api_calls', 'tokens', 'storage', 'users') NOT NULL COMMENT '指标类型',
    metric_date DATE NOT NULL COMMENT '计量日期',
    metric_value BIGINT NOT NULL COMMENT '指标值',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uk_tenant_type_date (tenant_id, metric_type, metric_date),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_metric_date (metric_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='使用计量表';
```

---

## 4. 计费引擎设计

### 4.1 计费规则引擎

```go
package billing

// BillingEngine 计费引擎
type BillingEngine struct {
    planRepo      repository.PlanRepository
    usageRepo     repository.UsageRepository
    invoiceSvc    *InvoiceService
}

// CalculateBill 计算账单
func (e *BillingEngine) CalculateBill(
    ctx context.Context,
    tenantID string,
    startDate, endDate time.Time,
) (*Bill, error) {
    // 1. 获取订阅信息
    subscription, err := e.getActiveSubscription(ctx, tenantID)
    if err != nil {
        return nil, err
    }

    // 2. 获取计划配置
    plan, err := e.planRepo.GetByID(ctx, subscription.PlanID)
    if err != nil {
        return nil, err
    }

    // 3. 计算基础费用
    baseAmount := e.calculateBaseFee(subscription, plan)

    // 4. 计算使用量
    usage, err := e.usageRepo.GetUsage(ctx, tenantID, startDate, endDate)
    if err != nil {
        return nil, err
    }

    // 5. 计算超额费用
    overageAmount := e.calculateOverage(plan.Quotas, usage)

    // 6. 应用折扣
    discount := baseAmount * subscription.DiscountRate

    // 7. 生成账单
    bill := &Bill{
        TenantID:    tenantID,
        StartDate:   startDate,
        EndDate:     endDate,
        BaseAmount:  baseAmount,
        OverageAmount: overageAmount,
        Discount:    discount,
        Tax:         (baseAmount + overageAmount - discount) * 0.06, // 6%增值税
        Total:       baseAmount + overageAmount - discount + tax,
    }

    return bill, nil
}

// calculateOverage 计算超额费用
func (e *BillingEngine) calculateOverage(
    quotas map[string]interface{},
    usage map[string]int64,
) decimal.Decimal {
    var total decimal.Decimal

    // API调用超额
    maxAPICalls := int64(quotas["maxAPICallsPerMonth"].(float64))
    if apiCalls := usage["api_calls"]; apiCalls > maxAPICalls {
        overage := apiCalls - maxAPICalls
        total = total.Add(decimal.NewFromInt(overage).Mul(decimal.NewFromFloat(0.01))) // ¥0.01/次
    }

    // Token超额
    maxTokens := int64(quotas["maxTokensPerMonth"].(float64))
    if tokens := usage["tokens"]; tokens > maxTokens {
        overage := tokens - maxTokens
        total = total.Add(decimal.NewFromInt(overage/1000).Mul(decimal.NewFromFloat(0.1))) // ¥0.1/千tokens
    }

    return total
}
```

---

## 5. 配置系统设计

### 5.1 预置订阅计划

```sql
INSERT INTO subscription_plans (id, name, type, monthly_price, yearly_price, quotas, features) VALUES
-- 免费版
('plan-free', '免费版', 'free', 0, 0,
'{"maxUsers":5,"maxBots":1,"maxKnowledgeBases":1,"maxDocuments":10,"maxStorageGB":1,"maxAPICallsPerMonth":1000,"maxTokensPerMonth":100000}',
'["basic_bot","knowledge_base"]'),

-- 专业版
('plan-professional', '专业版', 'professional', 299, 2990,
'{"maxUsers":50,"maxBots":10,"maxKnowledgeBases":5,"maxDocuments":500,"maxStorageGB":10,"maxAPICallsPerMonth":100000,"maxTokensPerMonth":10000000}',
'["basic_bot","knowledge_base","api_access","custom_domain","priority_support"]'),

-- 企业版
('plan-enterprise', '企业版', 'enterprise', 999, 9990,
'{"maxUsers":500,"maxBots":100,"maxKnowledgeBases":20,"maxDocuments":10000,"maxStorageGB":100,"maxAPICallsPerMonth":1000000,"maxTokensPerMonth":100000000}',
'["basic_bot","knowledge_base","api_access","custom_domain","priority_support","sla_guarantee","dedicated_manager"]');
```

---

## 6. API 设计

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | /api/v1/billing/plans | 列出订阅计划 |
| POST | /api/v1/billing/subscribe | 订阅计划 |
| GET | /api/v1/billing/subscription | 查看当前订阅 |
| GET | /api/v1/billing/invoices | 列出账单 |
| GET | /api/v1/billing/invoices/:id | 查看账单详情 |
| POST | /api/v1/billing/invoices/:id/pay | 支付账单 |
| GET | /api/v1/billing/usage | 查看使用量 |

---

## 7. 总结

### 7.1 实施策略总结

| 实施项 | 实施方式 | 工作量 |
|-------|---------|--------|
| **计费引擎** | 💻 独立编码 | 40% |
| **计费规则** | 📊 配置驱动 | 60% |

**总计**：40% 编码 + 60% 配置

### 7.2 核心优势

- ✅ **多计费模式**：支持订阅、按量、混合计费
- ✅ **自动化**：自动计量、自动出账、自动扣款
- ✅ **灵活配置**：计费规则数据库配置

---

**文档结束**
