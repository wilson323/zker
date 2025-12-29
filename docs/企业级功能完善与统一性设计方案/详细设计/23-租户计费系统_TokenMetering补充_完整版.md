# 23-租户计费系统_Token Metering补充完整版

**模块**: 23-租户计费系统扩展
**扩展内容**: Token Metering、成本防护、成本优化
**版本**: v2.0
**日期**: 2025-01-03
**优先级**: P0
**工期**: 2周

---

## 目录

- [1. 功能概述](#1-功能概述)
- [2. 数据库设计](#2-数据库设计)
- [3. Token计量引擎](#3-token计量引擎)
- [4. 预算告警系统](#4-预算告警系统)
- [5. 自动降级策略](#5-自动降级策略)
- [6. 成本优化引擎](#6-成本优化引擎)
- [7. API设计](#7-api设计)
- [8. 前端组件](#8-前端组件)
- [9. 实施计划](#9-实施计划)

---

## 1. 功能概述

### 1.1 核心目标

**Token Metering** 提供精细化的AI成本管理能力：

- ✅ **实时计量** - Token级别使用追踪
- ✅ **预算防护** - 多级告警 + 硬性上限
- ✅ **自动降级** - 超预算自动切换模型
- ✅ **成本洞察** - 多维度成本分析
- ✅ **智能优化** - AI驱动的成本优化建议

### 1.2 应用场景

**场景1：预算预警**
```
租户A月度预算 ¥10,000
- 当使用达到 ¥8,000 (80%) → 发送邮件告警
- 当使用达到 ¥9,500 (95%) → 发送短信 + 站内信告警
- 当使用达到 ¥10,000 (100%) → 触发硬性上限，停止服务或降级
```

**场景2：自动降级**
```
租户B配置降级策略：
- 原模型: GPT-4 (¥0.12/1K tokens)
- 降级模型: GPT-3.5-Turbo (¥0.01/1K tokens)
- 阈值: 80%预算使用率
- 效果: 节省 92% 成本
```

**场景3：成本优化建议**
```
AI分析发现：
- Bot A 使用 GPT-4 处理简单问答 → 建议切换到 GPT-3.5
- Bot B 大量重复请求 → 建议启用缓存
- 预计节省: 35%
```

---

## 2. 数据库设计

### 2.1 Token使用明细表

```sql
CREATE TABLE token_usage_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    user_id BIGINT COMMENT '用户ID',
    bot_id VARCHAR(64) COMMENT 'Bot ID',
    conversation_id VARCHAR(64) COMMENT '会话ID',
    message_id VARCHAR(64) COMMENT '消息ID',

    -- Token统计
    input_tokens INT NOT NULL COMMENT '输入Token数',
    output_tokens INT NOT NULL COMMENT '输出Token数',
    total_tokens INT NOT NULL COMMENT '总Token数',

    -- 模型信息
    model_provider VARCHAR(50) NOT NULL COMMENT '模型提供商(openai/anthropic/通义千问)',
    model_name VARCHAR(50) NOT NULL COMMENT '模型名称(gpt-4/claude-3-opus/qwen-max)',
    model_version VARCHAR(20) COMMENT '模型版本',

    -- 成本计算
    unit_price DECIMAL(10,6) NOT NULL COMMENT '每1K Token单价(CNY)',
    input_cost DECIMAL(10,6) NOT NULL COMMENT '输入成本',
    output_cost DECIMAL(10,6) NOT NULL COMMENT '输出成本',
    total_cost DECIMAL(10,6) NOT NULL COMMENT '总成本',

    -- 性能指标
    response_time_ms INT COMMENT '响应时间(毫秒)',
    latency_ms INT COMMENT '首字延迟(毫秒)',
    is_cached BOOLEAN DEFAULT FALSE COMMENT '是否缓存命中',

    -- 元数据
    request_type ENUM('chat', 'completion', 'embedding', 'rerank') NOT NULL,
    metadata JSON COMMENT '额外元数据',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_created (tenant_id, created_at),
    INDEX idx_bot_created (bot_id, created_at),
    INDEX idx_model (model_provider, model_name),
    INDEX idx_cost (total_cost)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Token使用明细日志表';
```

### 2.2 Token使用汇总表

```sql
CREATE TABLE token_usage_summary (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    bot_id VARCHAR(64) COMMENT 'Bot ID (NULL表示整体汇总)',
    summary_date DATE NOT NULL COMMENT '汇总日期',
    summary_hour TINYINT COMMENT '小时 (0-23, NULL表示日汇总)',

    -- Token汇总
    total_input_tokens BIGINT NOT NULL DEFAULT 0,
    total_output_tokens BIGINT NOT NULL DEFAULT 0,
    total_tokens BIGINT NOT NULL DEFAULT 0,

    -- 成本汇总
    total_cost DECIMAL(12,6) NOT NULL DEFAULT 0.000000,

    -- 请求统计
    total_requests INT NOT NULL DEFAULT 0,
    cached_requests INT NOT NULL DEFAULT 0,
    avg_response_time DECIMAL(8,2) COMMENT '平均响应时间(ms)',

    -- 模型分布
    model_distribution JSON COMMENT '模型使用分布 {model_name: {tokens, cost}}',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_tenant_bot_date_hour (tenant_id, bot_id, summary_date, summary_hour),
    INDEX idx_tenant_date (tenant_id, summary_date),
    INDEX idx_bot_date (bot_id, summary_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Token使用汇总表(按天/小时)';
```

### 2.3 预算配置表

```sql
CREATE TABLE budget_settings (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL UNIQUE COMMENT '租户ID',

    -- 预算配置
    budget_type ENUM('monthly', 'quarterly', 'yearly') NOT NULL DEFAULT 'monthly',
    budget_amount DECIMAL(12,2) NOT NULL COMMENT '预算金额(CNY)',
    currency VARCHAR(3) DEFAULT 'CNY',

    -- 告警阈值
    alert_threshold_1 INT DEFAULT 80 COMMENT '一级告警阈值(%)',
    alert_threshold_2 INT DEFAULT 95 COMMENT '二级告警阈值(%)',
    hard_cap_enabled BOOLEAN DEFAULT FALSE COMMENT '是否启用硬性上限',
    hard_cap_amount DECIMAL(12,2) COMMENT '硬性上限金额',

    -- 降级策略
    auto_downgrade_enabled BOOLEAN DEFAULT FALSE COMMENT '是否启用自动降级',
    downgrade_threshold INT DEFAULT 90 COMMENT '降级阈值(%)',
    original_model VARCHAR(100) COMMENT '原模型',
    fallback_model VARCHAR(100) COMMENT '降级模型',

    -- 通知配置
    notification_channels JSON COMMENT '通知渠道 ["email","sms","webhook"]',
    notification_recipients JSON COMMENT '通知接收人列表',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_tenant_id (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='预算配置表';
```

### 2.4 预算告警历史表

```sql
CREATE TABLE budget_alerts_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    alert_type ENUM('threshold_1', 'threshold_2', 'hard_cap', 'downgrade') NOT NULL,

    -- 预算状态
    budget_amount DECIMAL(12,2) NOT NULL COMMENT '预算金额',
    used_amount DECIMAL(12,2) NOT NULL COMMENT '已使用金额',
    usage_percent DECIMAL(5,2) NOT NULL COMMENT '使用率(%)',

    -- 告警信息
    alert_level ENUM('warning', 'critical', 'emergency') NOT NULL,
    alert_message TEXT COMMENT '告警消息',

    -- 发送状态
    notification_channels JSON COMMENT '通知渠道',
    notification_status ENUM('pending', 'sent', 'failed') DEFAULT 'pending',
    sent_at DATETIME COMMENT '发送时间',
    error_message TEXT COMMENT '错误信息',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_created (tenant_id, created_at),
    INDEX idx_alert_type (alert_type),
    INDEX idx_status (notification_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='预算告警历史表';
```

### 2.5 成本优化建议表

```sql
CREATE TABLE cost_optimization_suggestions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    suggestion_type ENUM('model_downgrade', 'enable_cache', 'batch_request', 'prompt_optimization') NOT NULL,

    -- 分析数据
    analysis_period_start DATE NOT NULL COMMENT '分析开始日期',
    analysis_period_end DATE NOT NULL COMMENT '分析结束日期',

    -- 建议内容
    target_bot_id VARCHAR(64) COMMENT '目标Bot ID',
    current_model VARCHAR(100) COMMENT '当前模型',
    suggested_model VARCHAR(100) COMMENT '建议模型',
    reason TEXT COMMENT '优化原因',

    -- 预期效果
    estimated_monthly_saving DECIMAL(12,2) COMMENT '预计月节省金额',
    estimated_saving_percent DECIMAL(5,2) COMMENT '预计节省比例(%)',

    -- 状态
    status ENUM('pending', 'approved', 'rejected', 'applied') DEFAULT 'pending',
    applied_at DATETIME COMMENT '应用时间',
    applied_by BIGINT COMMENT '应用操作人ID',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_type (suggestion_type),
    INDEX idx_saving (estimated_monthly_saving)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='成本优化建议表';
```

---

## 3. Token计量引擎

### 3.1 核心服务实现

```go
package metering

import (
    "context"
    "time"
    "github.com/cloudwego/hertz/pkg/common/hlog"
)

// TokenMeteringService Token计量服务
type TokenMeteringService struct {
    logRepo       repository.TokenUsageLogRepository
    summaryRepo   repository.TokenUsageSummaryRepository
    budgetRepo    repository.BudgetSettingsRepository
    alertSvc      *BudgetAlertService
    pricingEngine *PricingEngine
}

// RecordTokenUsage 记录Token使用
func (s *TokenMeteringService) RecordTokenUsage(
    ctx context.Context,
    req *RecordTokenUsageRequest,
) (*RecordTokenUsageResponse, error) {
    // 1. 计算成本
    cost := s.pricingEngine.CalculateCost(
        req.ModelProvider,
        req.ModelName,
        req.InputTokens,
        req.OutputTokens,
    )

    // 2. 构建日志记录
    log := &entity.TokenUsageLog{
        TenantID:       req.TenantID,
        UserID:         req.UserID,
        BotID:          req.BotID,
        ConversationID: req.ConversationID,
        MessageID:      req.MessageID,
        InputTokens:    req.InputTokens,
        OutputTokens:   req.OutputTokens,
        TotalTokens:    req.InputTokens + req.OutputTokens,
        ModelProvider:  req.ModelProvider,
        ModelName:      req.ModelName,
        ModelVersion:   req.ModelVersion,
        UnitPrice:      cost.UnitPrice,
        InputCost:      cost.InputCost,
        OutputCost:     cost.OutputCost,
        TotalCost:      cost.TotalCost,
        ResponseTimeMs: req.ResponseTimeMs,
        LatencyMs:      req.LatencyMs,
        IsCached:       req.IsCached,
        RequestType:    req.RequestType,
        Metadata:       req.Metadata,
    }

    // 3. 持久化日志
    if err := s.logRepo.Create(ctx, log); err != nil {
        hlog.Errorf("Failed to record token usage: %v", err)
        return nil, err
    }

    // 4. 异步更新汇总表
    go s.updateSummary(context.Background(), req.TenantID, req.BotID, log)

    // 5. 异步检查预算
    go s.alertSvc.CheckBudget(context.Background(), req.TenantID)

    return &RecordTokenUsageResponse{
        LogID:     log.ID,
        TotalCost: cost.TotalCost,
    }, nil
}

// updateSummary 更新汇总表
func (s *TokenMeteringService) updateSummary(
    ctx context.Context,
    tenantID, botID string,
    log *entity.TokenUsageLog,
) {
    now := time.Now()
    summaryDate := now.Format("2006-01-02")
    summaryHour := now.Hour()

    // 更新小时级汇总
    s.summaryRepo.Upsert(ctx, &entity.TokenUsageSummary{
        TenantID:          tenantID,
        BotID:             botID,
        SummaryDate:       summaryDate,
        SummaryHour:       summaryHour,
        TotalInputTokens:  log.InputTokens,
        TotalOutputTokens: log.OutputTokens,
        TotalTokens:       log.TotalTokens,
        TotalCost:         log.TotalCost,
        TotalRequests:     1,
        CachedRequests:    boolToInt(log.IsCached),
    })
}

// GetUsageStats 获取使用统计
func (s *TokenMeteringService) GetUsageStats(
    ctx context.Context,
    req *GetUsageStatsRequest,
) (*UsageStatsResponse, error) {
    // 1. 查询汇总数据
    summary, err := s.summaryRepo.GetByDateRange(
        ctx,
        req.TenantID,
        req.BotID,
        req.StartDate,
        req.EndDate,
        req.Granularity, // daily/hourly
    )
    if err != nil {
        return nil, err
    }

    // 2. 计算统计指标
    var totalTokens int64
    var totalCost decimal.Decimal
    var totalRequests int
    var cachedRequests int

    for _, s := range summary {
        totalTokens += s.TotalTokens
        totalCost = totalCost.Add(s.TotalCost)
        totalRequests += int(s.TotalRequests)
        cachedRequests += int(s.CachedRequests)
    }

    cacheHitRate := float64(0)
    if totalRequests > 0 {
        cacheHitRate = float64(cachedRequests) / float64(totalRequests) * 100
    }

    // 3. 构建响应
    return &UsageStatsResponse{
        Period: Period{
            StartDate: req.StartDate,
            EndDate:   req.EndDate,
        },
        TotalTokens:      totalTokens,
        TotalCost:        totalCost,
        TotalRequests:    totalRequests,
        CacheHitRate:     cacheHitRate,
        AvgTokensPerReq:  float64(totalTokens) / float64(totalRequests),
        DailyBreakdown:   s.buildDailyBreakdown(summary),
        ModelDistribution: s.getModelDistribution(ctx, req.TenantID, req.StartDate, req.EndDate),
    }, nil
}

// getModelDistribution 获取模型分布
func (s *TokenMeteringService) getModelDistribution(
    ctx context.Context,
    tenantID string,
    startDate, endDate string,
) map[string]ModelStats {
    logs, _ := s.logRepo.GetByDateRange(ctx, tenantID, startDate, endDate)

    modelStats := make(map[string]ModelStats)
    for _, log := range logs {
        key := log.ModelProvider + "/" + log.ModelName
        if stat, ok := modelStats[key]; ok {
            stat.TotalTokens += log.TotalTokens
            stat.TotalCost = stat.TotalCost.Add(log.TotalCost)
            stat.RequestCount++
        } else {
            modelStats[key] = ModelStats{
                ModelName:    log.ModelName,
                TotalTokens:  log.TotalTokens,
                TotalCost:    log.TotalCost,
                RequestCount: 1,
            }
        }
    }

    return modelStats
}
```

### 3.2 定价引擎

```go
package metering

import (
    "decimal"
)

// PricingEngine 定价引擎
type PricingEngine struct {
    // 模型价格配置 (每1K tokens价格, CNY)
    prices map[string]ModelPrice
}

type ModelPrice struct {
    InputPrice  decimal.Decimal // 输入价格
    OutputPrice decimal.Decimal // 输出价格
}

// NewPricingEngine 创建定价引擎
func NewPricingEngine() *PricingEngine {
    return &PricingEngine{
        prices: map[string]ModelPrice{
            "openai/gpt-4": {
                InputPrice:  decimal.NewFromFloat(0.03),  // ¥0.03/1K tokens
                OutputPrice: decimal.NewFromFloat(0.06),  // ¥0.06/1K tokens
            },
            "openai/gpt-3.5-turbo": {
                InputPrice:  decimal.NewFromFloat(0.003), // ¥0.003/1K tokens
                OutputPrice: decimal.NewFromFloat(0.006), // ¥0.006/1K tokens
            },
            "anthropic/claude-3-opus": {
                InputPrice:  decimal.NewFromFloat(0.09),
                OutputPrice: decimal.NewFromFloat(0.27),
            },
            "anthropic/claude-3-sonnet": {
                InputPrice:  decimal.NewFromFloat(0.015),
                OutputPrice: decimal.NewFromFloat(0.045),
            },
            "qwen/qwen-max": {
                InputPrice:  decimal.NewFromFloat(0.02),
                OutputPrice: decimal.NewFromFloat(0.06),
            },
            "qwen/qwen-plus": {
                InputPrice:  decimal.NewFromFloat(0.004),
                OutputPrice: decimal.NewFromFloat(0.012),
            },
        },
    }
}

// CalculateCost 计算成本
func (e *PricingEngine) CalculateCost(
    provider, model string,
    inputTokens, outputTokens int,
) *CostBreakdown {
    key := provider + "/" + model
    price, ok := e.prices[key]
    if !ok {
        // 默认价格
        price = ModelPrice{
            InputPrice:  decimal.NewFromFloat(0.01),
            OutputPrice: decimal.NewFromFloat(0.02),
        }
    }

    // 计算成本 (Token数 / 1000 * 单价)
    inputCost := decimal.NewFromInt(int64(inputTokens)).
        Div(decimal.NewFromInt(1000)).
        Mul(price.InputPrice)

    outputCost := decimal.NewFromInt(int64(outputTokens)).
        Div(decimal.NewFromInt(1000)).
        Mul(price.OutputPrice)

    totalCost := inputCost.Add(outputCost)

    // 计算平均单价
    totalTokens := inputTokens + outputTokens
    avgPrice := decimal.NewFromInt(int64(totalTokens)).
        Div(decimal.NewFromInt(1000)).
        Mul(decimal.NewFromFloat(1))

    unitPrice := totalCost.Div(avgPrice)

    return &CostBreakdown{
        UnitPrice:  unitPrice,
        InputCost:  inputCost,
        OutputCost: outputCost,
        TotalCost:  totalCost,
    }
}
```

---

## 4. 预算告警系统

### 4.1 预算检查服务

```go
package metering

import (
    "context"
    "time"
)

// BudgetAlertService 预算告警服务
type BudgetAlertService struct {
    budgetRepo   repository.BudgetSettingsRepository
    summaryRepo  repository.TokenUsageSummaryRepository
    alertRepo    repository.BudgetAlertRepository
    notifier     *NotificationService
}

// CheckBudget 检查预算状态
func (s *BudgetAlertService) CheckBudget(
    ctx context.Context,
    tenantID string,
) error {
    // 1. 获取预算配置
    budget, err := s.budgetRepo.GetByTenantID(ctx, tenantID)
    if err != nil || !budget.HardCapEnabled {
        return nil
    }

    // 2. 计算当前周期使用量
    now := time.Now()
    var startDate time.Time
    switch budget.BudgetType {
    case "monthly":
        startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
    case "quarterly":
        quarter := (now.Month() - 1) / 3
        startDate = time.Date(now.Year(), quarter*3+1, 1, 0, 0, 0, 0, now.Location())
    case "yearly":
        startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
    }

    usedAmount, _ := s.summaryRepo.GetTotalCost(ctx, tenantID, startDate, now)

    // 3. 计算使用率
    usagePercent := usedAmount.Div(budget.BudgetAmount).Mul(decimal.NewFromInt(100))

    // 4. 检查告警阈值
    if usagePercent.Cmp(decimal.NewFromInt(int64(budget.AlertThreshold1))) >= 0 {
        // 一级告警 (80%)
        s.sendAlertIfNeeded(ctx, tenantID, budget, usedAmount, usagePercent, "threshold_1")
    }

    if usagePercent.Cmp(decimal.NewFromInt(int64(budget.AlertThreshold2))) >= 0 {
        // 二级告警 (95%)
        s.sendAlertIfNeeded(ctx, tenantID, budget, usedAmount, usagePercent, "threshold_2")
    }

    // 5. 检查硬性上限
    if budget.HardCapEnabled && usedAmount.Cmp(budget.HardCapAmount) >= 0 {
        s.handleHardCap(ctx, tenantID, budget, usedAmount, usagePercent)
    }

    return nil
}

// sendAlertIfNeeded 发送告警 (防重复)
func (s *BudgetAlertService) sendAlertIfNeeded(
    ctx context.Context,
    tenantID string,
    budget *entity.BudgetSettings,
    usedAmount, usagePercent decimal.Decimal,
    alertType string,
) {
    // 检查是否已发送过相同类型的告警
    today := time.Now().Format("2006-01-02")
    alreadySent, _ := s.alertRepo.CheckAlertSentToday(ctx, tenantID, alertType, today)
    if alreadySent {
        return
    }

    // 构建告警消息
    alertLevel := "warning"
    if alertType == "threshold_2" {
        alertLevel = "critical"
    }

    message := s.buildAlertMessage(budget, usedAmount, usagePercent, alertLevel)

    // 记录告警历史
    alert := &entity.BudgetAlert{
        TenantID:       tenantID,
        AlertType:      alertType,
        BudgetAmount:   budget.BudgetAmount,
        UsedAmount:     usedAmount,
        UsagePercent:   usagePercent,
        AlertLevel:     alertLevel,
        AlertMessage:   message,
        NotificationChannels: budget.NotificationChannels,
    }

    s.alertRepo.Create(ctx, alert)

    // 发送通知
    s.notifier.SendAlert(ctx, alert)
}

// handleHardCap 处理硬性上限
func (s *BudgetAlertService) handleHardCap(
    ctx context.Context,
    tenantID string,
    budget *entity.BudgetSettings,
    usedAmount, usagePercent decimal.Decimal,
) {
    // 发送紧急告警
    alert := &entity.BudgetAlert{
        TenantID:       tenantID,
        AlertType:      "hard_cap",
        BudgetAmount:   budget.BudgetAmount,
        UsedAmount:     usedAmount,
        UsagePercent:   usagePercent,
        AlertLevel:     "emergency",
        AlertMessage:   "⚠️ 硬性上限触发！服务已暂停。请联系管理员充值或调整预算。",
        NotificationChannels: budget.NotificationChannels,
    }

    s.alertRepo.Create(ctx, alert)
    s.notifier.SendAlert(ctx, alert)

    // 执行暂停/降级策略
    if budget.AutoDowngradeEnabled {
        s.executeDowngrade(ctx, tenantID, budget)
    } else {
        // 暂停服务
        s.pauseService(ctx, tenantID)
    }
}

// pauseService 暂停服务
func (s *BudgetAlertService) pauseService(ctx context.Context, tenantID string) {
    // 更新租户状态为暂停
    // 实际实现需要调用租户服务
    hlog.Warnf("Tenant %s service paused due to hard cap", tenantID)
}
```

### 4.2 通知服务

```go
package metering

import (
    "context"
    "fmt"
)

// NotificationService 通知服务
type NotificationService struct {
    emailSvc   *EmailService
    smsSvc     *SMSService
    webhookSvc *WebhookService
}

// SendAlert 发送告警通知
func (s *NotificationService) SendAlert(
    ctx context.Context,
    alert *entity.BudgetAlert,
) error {
    channels := alert.NotificationChannels
    var lastErr error

    // 并发发送多个渠道
    for _, channel := range channels {
        switch channel {
        case "email":
            if err := s.emailSvc.SendAlertEmail(ctx, alert); err != nil {
                lastErr = err
            }
        case "sms":
            if err := s.smsSvc.SendAlertSMS(ctx, alert); err != nil {
                lastErr = err
            }
        case "webhook":
            if err := s.webhookSvc.SendWebhook(ctx, alert); err != nil {
                lastErr = err
            }
        }
    }

    // 更新发送状态
    alert.NotificationStatus = "sent"
    alert.SentAt = time.Now()
    if lastErr != nil {
        alert.NotificationStatus = "failed"
        alert.ErrorMessage = lastErr.Error()
    }

    return lastErr
}

// buildAlertMessage 构建告警消息
func (s *BudgetAlertService) buildAlertMessage(
    budget *entity.BudgetSettings,
    usedAmount, usagePercent decimal.Decimal,
    alertLevel string,
) string {
    emoji := "⚠️"
    if alertLevel == "critical" {
        emoji = "🚨"
    } else if alertLevel == "emergency" {
        emoji = "🛑"
    }

    message := fmt.Sprintf(
        "%s **预算告警**\n\n"+
        "租户ID: %s\n"+
        "预算金额: ¥%.2f\n"+
        "已使用: ¥%.2f (%.1f%%)\n\n",
        emoji,
        budget.TenantID,
        budget.BudgetAmount,
        usedAmount,
        usagePercent,
    )

    if alertLevel == "warning" {
        message += "建议: 请注意控制使用量，避免超额产生额外费用。"
    } else if alertLevel == "critical" {
        message += "警告: 即将达到预算上限，建议立即充值或调整预算配置。"
    }

    return message
}
```

---

## 5. 自动降级策略

### 5.1 降级执行服务

```go
package metering

import (
    "context"
    "github.com/cloudwego/hertz/pkg/common/hlog"
)

// DowngradeService 降级服务
type DowngradeService struct {
    budgetRepo     repository.BudgetSettingsRepository
    tenantSvc      tenant.TenantService
    modelRouter    routing.ModelRouter
}

// executeDowngrade 执行降级
func (s *BudgetAlertService) executeDowngrade(
    ctx context.Context,
    tenantID string,
    budget *entity.BudgetSettings,
) {
    // 1. 记录降级事件
    alert := &entity.BudgetAlert{
        TenantID:     tenantID,
        AlertType:    "downgrade",
        AlertLevel:   "warning",
        AlertMessage: fmt.Sprintf(
            "自动降级已触发: %s → %s",
            budget.OriginalModel,
            budget.FallbackModel,
        ),
    }
    s.alertRepo.Create(ctx, alert)
    s.notifier.SendAlert(ctx, alert)

    // 2. 更新模型路由配置
    // 原模型: GPT-4 → 降级模型: GPT-3.5-Turbo
    routingRule := &routing.ModelRoutingRule{
        TenantID:     tenantID,
        OriginalModel: budget.OriginalModel,
        FallbackModel: budget.FallbackModel,
        Reason:       "budget_overage",
        Enabled:      true,
    }

    if err := s.modelRouter.UpdateRoutingRule(ctx, routingRule); err != nil {
        hlog.Errorf("Failed to update routing rule: %v", err)
        return
    }

    hlog.Infof("Tenant %s model downgraded: %s → %s",
        tenantID, budget.OriginalModel, budget.FallbackModel)
}

// UpgradeModel 升级模型 (手动恢复)
func (s *DowngradeService) UpgradeModel(
    ctx context.Context,
    tenantID string,
) error {
    // 1. 获取预算配置
    budget, err := s.budgetRepo.GetByTenantID(ctx, tenantID)
    if err != nil {
        return err
    }

    // 2. 移除降级规则
    if err := s.modelRouter.RemoveRoutingRule(ctx, tenantID); err != nil {
        return err
    }

    // 3. 记录升级操作
    hlog.Infof("Tenant %s model upgraded back to: %s", tenantID, budget.OriginalModel)

    return nil
}
```

### 5.2 模型路由集成

```go
package routing

import (
    "context"
    "sync"
)

// ModelRouter 模型路由器
type ModelRouter struct {
    mu            sync.RWMutex
    routingRules  map[string]*ModelRoutingRule // tenantID -> rule
}

type ModelRoutingRule struct {
    TenantID       string
    OriginalModel  string
    FallbackModel  string
    Reason         string
    Enabled        bool
    CreatedAt      time.Time
}

// SelectModel 选择模型 (考虑降级规则)
func (r *ModelRouter) SelectModel(
    ctx context.Context,
    tenantID string,
    requestedModel string,
) string {
    r.mu.RLock()
    defer r.mu.RUnlock()

    // 检查是否有降级规则
    if rule, ok := r.routingRules[tenantID]; ok && rule.Enabled {
        if rule.OriginalModel == requestedModel {
            hlog.Infof("Routing %s for tenant %s due to %s",
                rule.FallbackModel, tenantID, rule.Reason)
            return rule.FallbackModel
        }
    }

    return requestedModel
}

// UpdateRoutingRule 更新路由规则
func (r *ModelRouter) UpdateRoutingRule(
    ctx context.Context,
    rule *ModelRoutingRule,
) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    rule.CreatedAt = time.Now()
    r.routingRules[rule.TenantID] = rule

    return nil
}

// RemoveRoutingRule 移除路由规则
func (r *ModelRouter) RemoveRoutingRule(
    ctx context.Context,
    tenantID string,
) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    delete(r.routingRules, tenantID)

    return nil
}
```

---

## 6. 成本优化引擎

### 6.1 优化分析服务

```go
package optimization

import (
    "context"
    "decimal"
    "time"
)

// CostOptimizationEngine 成本优化引擎
type CostOptimizationEngine struct {
    usageRepo      repository.TokenUsageLogRepository
    suggestionRepo repository.OptimizationSuggestionRepository
    analyzer       *UsageAnalyzer
}

// AnalyzeAndGenerateSuggestions 分析并生成优化建议
func (e *CostOptimizationEngine) AnalyzeAndGenerateSuggestions(
    ctx context.Context,
    tenantID string,
) ([]*OptimizationSuggestion, error) {
    // 1. 获取最近30天的使用数据
    endDate := time.Now()
    startDate := endDate.AddDate(0, 0, -30)

    logs, err := e.usageRepo.GetByDateRange(ctx, tenantID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
    if err != nil {
        return nil, err
    }

    var suggestions []*OptimizationSuggestion

    // 2. 模型降级分析
    if modelSuggestion := e.analyzer.AnalyzeModelDowngrade(ctx, tenantID, logs); modelSuggestion != nil {
        suggestions = append(suggestions, modelSuggestion)
    }

    // 3. 缓存优化分析
    if cacheSuggestion := e.analyzer.AnalyzeCacheOptimization(ctx, tenantID, logs); cacheSuggestion != nil {
        suggestions = append(suggestions, cacheSuggestion)
    }

    // 4. Prompt优化分析
    if promptSuggestion := e.analyzer.AnalyzePromptOptimization(ctx, tenantID, logs); promptSuggestion != nil {
        suggestions = append(suggestions, promptSuggestion)
    }

    // 5. 批量请求优化
    if batchSuggestion := e.analyzer.AnalyzeBatchRequests(ctx, tenantID, logs); batchSuggestion != nil {
        suggestions = append(suggestions, batchSuggestion)
    }

    // 6. 保存建议到数据库
    for _, suggestion := range suggestions {
        e.suggestionRepo.Create(ctx, suggestion)
    }

    return suggestions, nil
}
```

### 6.2 使用分析器

```go
package optimization

import (
    "context"
    "decimal"
)

// UsageAnalyzer 使用分析器
type UsageAnalyzer struct {
    pricingEngine *metering.PricingEngine
}

// AnalyzeModelDowngrade 分析模型降级机会
func (a *UsageAnalyzer) AnalyzeModelDowngrade(
    ctx context.Context,
    tenantID string,
    logs []*entity.TokenUsageLog,
) *OptimizationSuggestion {
    // 1. 按Bot分组分析
    botStats := make(map[string]*BotModelStats)

    for _, log := range logs {
        if _, ok := botStats[log.BotID]; !ok {
            botStats[log.BotID] = &BotModelStats{
                BotID:        log.BotID,
                ModelName:    log.ModelName,
                TotalTokens:  0,
                TotalCost:    decimal.Zero,
                RequestCount: 0,
            }
        }

        botStats[log.BotID].TotalTokens += log.TotalTokens
        botStats[log.BotID].TotalCost = botStats[log.BotID].TotalCost.Add(log.TotalCost)
        botStats[log.BotID].RequestCount++
    }

    // 2. 寻找降级机会 (使用GPT-4但场景简单的Bot)
    for _, stats := range botStats {
        // 判断条件: 使用GPT-4 + 平均Token数较低 (< 1000 tokens/request)
        if stats.ModelName == "gpt-4" &&
           float64(stats.TotalTokens)/float64(stats.RequestCount) < 1000 {

            // 计算降级到GPT-3.5的节省
            currentCost := stats.TotalCost
            estimatedCost := decimal.NewFromInt(stats.TotalTokens / 1000).
                Mul(decimal.NewFromFloat(0.009)) // GPT-3.5平均价格

            monthlySaving := currentCost.Sub(estimatedCost)
            savingPercent := monthlySaving.Div(currentCost).Mul(decimal.NewFromInt(100))

            // 只有节省比例 > 20% 才建议
            if savingPercent.Cmp(decimal.NewFromInt(20)) > 0 {
                return &OptimizationSuggestion{
                    TenantID:               tenantID,
                    SuggestionType:         "model_downgrade",
                    TargetBotID:            stats.BotID,
                    CurrentModel:           "gpt-4",
                    SuggestedModel:         "gpt-3.5-turbo",
                    Reason:                 "该Bot平均Token使用量较低，可降级到GPT-3.5而不影响质量",
                    EstimatedMonthlySaving: monthlySaving,
                    EstimatedSavingPercent: savingPercent,
                    Status:                 "pending",
                }
            }
        }
    }

    return nil
}

// AnalyzeCacheOptimization 分析缓存优化机会
func (a *UsageAnalyzer) AnalyzeCacheOptimization(
    ctx context.Context,
    tenantID string,
    logs []*entity.TokenUsageLog,
) *OptimizationSuggestion {
    // 1. 统计缓存命中率
    var totalRequests int
    var cachedRequests int

    for _, log := range logs {
        totalRequests++
        if log.IsCached {
            cachedRequests++
        }
    }

    cacheHitRate := float64(cachedRequests) / float64(totalRequests) * 100

    // 2. 如果缓存命中率 < 30%，建议优化
    if cacheHitRate < 30.0 {
        // 计算潜在节省
        uncachedRequests := totalRequests - cachedRequests
        avgTokensPerReq := 1000.0 // 假设平均1000 tokens
        potentialCachedTokens := float64(uncachedRequests) * avgTokensPerReq * 0.5 // 假设50%可缓存

        avgPrice := 0.03 // ¥0.03/1K tokens
        potentialSaving := potentialCachedTokens / 1000 * avgPrice

        return &OptimizationSuggestion{
            TenantID:               tenantID,
            SuggestionType:         "enable_cache",
            Reason:                 fmt.Sprintf(
                "当前缓存命中率 %.1f%% 较低，建议启用语义缓存。预计可节省 ¥%.2f/月",
                cacheHitRate, potentialSaving),
            EstimatedMonthlySaving: decimal.NewFromFloat(potentialSaving),
            EstimatedSavingPercent: decimal.NewFromFloat(15.0), // 假设节省15%
            Status:                 "pending",
        }
    }

    return nil
}

// AnalyzePromptOptimization 分析Prompt优化机会
func (a *UsageAnalyzer) AnalyzePromptOptimization(
    ctx context.Context,
    tenantID string,
    logs []*entity.TokenUsageLog,
) *OptimizationSuggestion {
    // 分析输入/输出Token比例
    var totalInputTokens, totalOutputTokens int

    for _, log := range logs {
        totalInputTokens += log.InputTokens
        totalOutputTokens += log.OutputTokens
    }

    ioRatio := float64(totalInputTokens) / float64(totalOutputTokens)

    // 如果输入Token远大于输出Token (比例 > 5:1)，说明Prompt可能过于冗长
    if ioRatio > 5.0 {
        totalTokens := totalInputTokens + totalOutputTokens
        potentialReduction := float64(totalInputTokens) * 0.3 // 假设可减少30%

        avgPrice := 0.03
        potentialSaving := potentialReduction / 1000 * avgPrice

        return &OptimizationSuggestion{
            TenantID:               tenantID,
            SuggestionType:         "prompt_optimization",
            Reason:                 fmt.Sprintf(
                "输入/输出Token比例 %.1f:1 过高，建议精简Prompt。预计可减少30%%输入Token，节省 ¥%.2f/月",
                ioRatio, potentialSaving),
            EstimatedMonthlySaving: decimal.NewFromFloat(potentialSaving),
            EstimatedSavingPercent: decimal.NewFromFloat(10.0),
            Status:                 "pending",
        }
    }

    return nil
}
```

---

## 7. API设计

### 7.1 Token计量API

| 方法 | 路径 | 功能 | 请求示例 | 响应示例 |
|------|------|------|----------|----------|
| POST | /api/v1/metering/record | 记录Token使用 | `{"tenant_id":"t1","bot_id":"b1","input_tokens":100,"output_tokens":50,...}` | `{"log_id":123,"total_cost":"0.009"}` |
| GET | /api/v1/metering/usage | 查询使用统计 | `?start_date=2025-01-01&end_date=2025-01-31` | 见下方完整响应 |

**使用统计响应**:
```json
{
  "period": {
    "start_date": "2025-01-01",
    "end_date": "2025-01-31"
  },
  "total_tokens": 15000000,
  "total_cost": "450.00",
  "total_requests": 15000,
  "cache_hit_rate": 25.5,
  "avg_tokens_per_req": 1000.0,
  "daily_breakdown": [
    {
      "date": "2025-01-01",
      "tokens": 500000,
      "cost": "15.00",
      "requests": 500
    }
  ],
  "model_distribution": {
    "openai/gpt-4": {
      "total_tokens": 10000000,
      "total_cost": "400.00",
      "request_count": 8000
    },
    "openai/gpt-3.5-turbo": {
      "total_tokens": 5000000,
      "total_cost": "50.00",
      "request_count": 7000
    }
  }
}
```

### 7.2 预算管理API

| 方法 | 路径 | 功能 | 请求示例 | 响应示例 |
|------|------|------|----------|----------|
| GET | /api/v1/billing/budget | 获取预算配置 | - | `{budget_amount:10000,alert_threshold_1:80,...}` |
| PUT | /api/v1/billing/budget | 更新预算配置 | `{"budget_amount":20000,...}` | `{"success":true}` |
| GET | /api/v1/billing/budget/alerts | 查询告警历史 | `?limit=20` | `[{alert_type:"threshold_1",...}]` |
| GET | /api/v1/billing/budget/status | 获取预算状态 | - | `{used_amount:8000,usage_percent:80,...}` |

### 7.3 成本优化API

| 方法 | 路径 | 功能 | 请求示例 | 响应示例 |
|------|------|------|----------|----------|
| GET | /api/v1/billing/optimization/suggestions | 获取优化建议 | - | `[{suggestion_type:"model_downgrade",...}]` |
| POST | /api/v1/billing/optimization/apply | 应用优化建议 | `{"suggestion_id":123}` | `{"success":true}` |
| POST | /api/v1/billing/optimization/reject | 拒绝优化建议 | `{"suggestion_id":123,"reason":"不需要"}` | `{"success":true}` |

---

## 8. 前端组件

### 8.1 Token使用仪表盘

```typescript
// components/metering/TokenUsageDashboard.tsx
import React, { useEffect, useState } from 'react';
import { Card, Metric, AreaChart, Tablerow } from '@douyinfe/semi-ui';

interface TokenUsageStats {
  totalTokens: number;
  totalCost: string;
  totalRequests: number;
  cacheHitRate: number;
  dailyBreakdown: Array<{
    date: string;
    tokens: number;
    cost: string;
    requests: number;
  }>;
  modelDistribution: Record<string, {
    totalTokens: number;
    totalCost: string;
    requestCount: number;
  }>;
}

export const TokenUsageDashboard: React.FC<Props> = ({ tenantId, dateRange }) => {
  const [stats, setStats] = useState<TokenUsageStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchUsageStats();
  }, [tenantId, dateRange]);

  const fetchUsageStats = async () => {
    setLoading(true);
    const response = await fetch(
      `/api/v1/metering/usage?tenant_id=${tenantId}&start_date=${dateRange.start}&end_date=${dateRange.end}`
    );
    const data = await response.json();
    setStats(data);
    setLoading(false);
  };

  if (loading) return <Spin />;

  return (
    <div className="token-usage-dashboard">
      {/* 关键指标卡片 */}
      <div className="metrics-grid">
        <Card title="总Token使用量">
          <Metric value={stats.totalTokens} fmt="value" />
        </Card>
        <Card title="总成本">
          <Metric value={stats.totalCost} prefix="¥" />
        </Card>
        <Card title="总请求数">
          <Metric value={stats.totalRequests} />
        </Card>
        <Card title="缓存命中率">
          <Metric value={stats.cacheHitRate} suffix="%" />
        </Card>
      </div>

      {/* 每日使用趋势图 */}
      <Card title="Token使用趋势">
        <AreaChart
          data={stats.dailyBreakdown}
          xField="date"
          yField="tokens"
          yAxis={{ label: { formatter: (v) => `${(v/1000).toFixed(0)}K` }}}
        />
      </Card>

      {/* 模型分布饼图 */}
      <Card title="模型使用分布">
        <PieChart
          data={Object.entries(stats.modelDistribution).map(([model, data]) => ({
            name: model,
            value: data.totalCost,
          }))}
        />
      </Card>
    </div>
  );
};
```

### 8.2 预算配置组件

```typescript
// components/billing/BudgetSettings.tsx
import React, { useState } from 'react';
import { Form, Input, Select, Switch, Slider, Button } from '@douyinfe/semi-ui';

interface BudgetConfig {
  budgetType: 'monthly' | 'quarterly' | 'yearly';
  budgetAmount: number;
  alertThreshold1: number; // 80
  alertThreshold2: number; // 95
  hardCapEnabled: boolean;
  hardCapAmount: number | null;
  autoDowngradeEnabled: boolean;
  downgradeThreshold: number;
  originalModel: string;
  fallbackModel: string;
}

export const BudgetSettings: React.FC = () => {
  const [form] = Form.useForm();
  const [saving, setSaving] = useState(false);

  const handleSubmit = async (values: BudgetConfig) => {
    setSaving(true);
    try {
      await fetch('/api/v1/billing/budget', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(values),
      });
      Toast.success('预算配置已保存');
    } catch (error) {
      Toast.error('保存失败');
    } finally {
      setSaving(false);
    }
  };

  return (
    <Form
      form={form}
      onSubmit={handleSubmit}
      initValues={{
        budgetType: 'monthly',
        budgetAmount: 10000,
        alertThreshold1: 80,
        alertThreshold2: 95,
        hardCapEnabled: false,
        autoDowngradeEnabled: false,
        downgradeThreshold: 90,
      }}
    >
      <Form.Select
        field="budgetType"
        label="预算周期"
        style={{ width: 200 }}
      >
        <Select.Option value="monthly">月度</Select.Option>
        <Select.Option value="quarterly">季度</Select.Option>
        <Select.Option value="yearly">年度</Select.Option>
      </Form.Select>

      <Form.InputNumber
        field="budgetAmount"
        label="预算金额(CNY)"
        prefix="¥"
        min={0}
        step={1000}
        style={{ width: 200 }}
      />

      <div className="threshold-group">
        <label>一级告警阈值</label>
        <Form.Slider
          field="alertThreshold1"
          min={50}
          max={100}
          step={5}
          marks={{ 60: '60%', 80: '80%', 90: '90%' }}
          style={{ width: 300 }}
        />
        <span className="threshold-value">{form.getFieldValue('alertThreshold1')}%</span>
      </div>

      <div className="threshold-group">
        <label>二级告警阈值</label>
        <Form.Slider
          field="alertThreshold2"
          min={50}
          max={100}
          step={5}
          marks={{ 80: '80%', 95: '95%', 99: '99%' }}
          style={{ width: 300 }}
        />
        <span className="threshold-value">{form.getFieldValue('alertThreshold2')}%</span>
      </div>

      <Form.Switch
        field="hardCapEnabled"
        label="启用硬性上限"
        extraText="达到上限后自动暂停服务或降级"
      />

      <Form.Switch
        field="autoDowngradeEnabled"
        label="启用自动降级"
        extraText="达到阈值后自动切换到便宜模型"
      />

      {form.getFieldValue('autoDowngradeEnabled') && (
        <>
          <Form.Select
            field="originalModel"
            label="原模型"
            style={{ width: 250 }}
          >
            <Select.Option value="openai/gpt-4">GPT-4</Select.Option>
            <Select.Option value="anthropic/claude-3-opus">Claude-3 Opus</Select.Option>
            <Select.Option value="qwen/qwen-max">通义千问 Max</Select.Option>
          </Form.Select>

          <Form.Select
            field="fallbackModel"
            label="降级模型"
            style={{ width: 250 }}
          >
            <Select.Option value="openai/gpt-3.5-turbo">GPT-3.5 Turbo</Select.Option>
            <Select.Option value="anthropic/claude-3-sonnet">Claude-3 Sonnet</Select.Option>
            <Select.Option value="qwen/qwen-plus">通义千问 Plus</Select.Option>
          </Form.Select>
        </>
      )}

      <Button type="primary" htmlType="submit" loading={saving}>
        保存配置
      </Button>
    </Form>
  );
};
```

### 8.3 成本优化建议组件

```typescript
// components/billing/CostOptimizationSuggestions.tsx
import React, { useEffect, useState } from 'react';
import { List, Card, Button, Tag, Empty } from '@douyinfe/semi-ui';

interface OptimizationSuggestion {
  id: number;
  suggestion_type: 'model_downgrade' | 'enable_cache' | 'prompt_optimization';
  target_bot_id?: string;
  current_model?: string;
  suggested_model?: string;
  reason: string;
  estimated_monthly_saving: string;
  estimated_saving_percent: number;
  status: 'pending' | 'approved' | 'rejected' | 'applied';
}

export const CostOptimizationSuggestions: React.FC = () => {
  const [suggestions, setSuggestions] = useState<OptimizationSuggestion[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchSuggestions();
  }, []);

  const fetchSuggestions = async () => {
    const response = await fetch('/api/v1/billing/optimization/suggestions');
    const data = await response.json();
    setSuggestions(data);
    setLoading(false);
  };

  const handleApply = async (suggestionId: number) => {
    try {
      await fetch('/api/v1/billing/optimization/apply', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ suggestion_id: suggestionId }),
      });
      Toast.success('优化建议已应用');
      fetchSuggestions();
    } catch (error) {
      Toast.error('应用失败');
    }
  };

  const handleReject = async (suggestionId: number) => {
    try {
      await fetch('/api/v1/billing/optimization/reject', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ suggestion_id: suggestionId }),
      });
      Toast.success('已拒绝建议');
      fetchSuggestions();
    } catch (error) {
      Toast.error('操作失败');
    }
  };

  if (loading) return <Spin />;

  if (suggestions.length === 0) {
    return <Empty title="暂无优化建议" description="系统将定期分析并生成成本优化建议" />;
  }

  return (
    <div className="optimization-suggestions">
      <h3>💡 成本优化建议</h3>
      <List
        dataSource={suggestions}
        renderItem={item => (
          <List.Item>
            <Card
              title={
                <div>
                  <Tag color="green">¥{item.estimated_monthly_saving}/月</Tag>
                  <Tag color="blue">{item.estimated_saving_percent}%</Tag>
                  {item.target_bot_id && <Tag>Bot: {item.target_bot_id}</Tag>}
                </div>
              }
              headerExtraContent={
                <div>
                  {item.status === 'pending' && (
                    <>
                      <Button
                        type="primary"
                        size="small"
                        onClick={() => handleApply(item.id)}
                      >
                        应用
                      </Button>
                      <Button
                        size="small"
                        onClick={() => handleReject(item.id)}
                      >
                        拒绝
                      </Button>
                    </>
                  )}
                  {item.status === 'applied' && <Tag color="green">已应用</Tag>}
                  {item.status === 'rejected' && <Tag color="red">已拒绝</Tag>}
                </div>
              }
            >
              <p>{item.reason}</p>
              {item.current_model && item.suggested_model && (
                <p>
                  模型切换: <code>{item.current_model}</code> → <code>{item.suggested_model}</code>
                </p>
              )}
            </Card>
          </List.Item>
        )}
      />
    </div>
  );
};
```

---

## 9. 实施计划

### 9.1 开发阶段划分

| 阶段 | 任务 | 工期 | 交付物 |
|------|------|------|--------|
| **第1周** | 数据库设计与创建 | 2天 | 5张表 |
| | Token计量引擎开发 | 3天 | Go服务 |
| **第2周** | 预算告警系统 | 2天 | 告警服务 |
| | 自动降级策略 | 2天 | 降级服务 |
| | 成本优化引擎 | 1天 | 分析器 |
| **第3周** | 前端组件开发 | 3天 | 3个组件 |
| | API集成测试 | 2天 | 测试报告 |
| **第4周** | 性能优化 | 2天 | 优化报告 |
| | 文档与培训 | 3天 | 用户手册 |

### 9.2 技术风险与缓解

| 风险 | 影响 | 概率 | 缓解措施 |
|------|------|------|----------|
| Token计量性能瓶颈 | 高 | 中 | 使用异步写入 + 批量处理 |
| 成本计算准确性 | 高 | 低 | 定期与云服务商账单对账 |
| 降级误判 | 中 | 中 | 增加人工审核环节 |
| 优化建议质量 | 中 | 中 | 持续训练分析模型 |

### 9.3 成功指标

- ✅ Token计量准确率: **99.9%**
- ✅ 预算告警及时性: **< 1分钟**
- ✅ 成本优化建议采纳率: **> 30%**
- ✅ 平均成本节省: **15-25%**

---

**文档结束**
