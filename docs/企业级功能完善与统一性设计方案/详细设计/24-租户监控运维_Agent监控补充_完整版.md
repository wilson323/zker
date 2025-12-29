# 24-租户监控运维_Agent监控补充完整版

**模块**: 24-租户监控运维扩展
**扩展内容**: AI Agent应用级监控与告警
**版本**: v2.0
**日期**: 2025-01-03
**优先级**: P0
**工期**: 2周

---

## 目录

- [1. 功能概述](#1-功能概述)
- [2. 数据库设计](#2-数据库设计)
- [3. Agent监控指标体系](#3-agent监控指标体系)
- [4. 指标收集服务](#4-指标收集服务)
- [5. 告警引擎](#5-告警引擎)
- [6. 监控仪表盘](#6-监控仪表盘)
- [7. API设计](#7-api设计)
- [8. 前端组件](#8-前端组件)
- [9. 实施计划](#9-实施计划)

---

## 1. 功能概述

### 1.1 核心目标

**Agent应用级监控** 提供AI Agent特有的细粒度监控能力：

- ✅ **性能监控** - QPS、响应时间（P50/P95/P99）、并发数
- ✅ **质量监控** - 错误率、用户满意度、Token消耗、成功率
- ✅ **资源监控** - 模型调用次数、成本追踪、缓存命中率
- ✅ **智能告警** - 多级阈值告警、趋势异常检测
- ✅ **可视化大盘** - 租户级、Bot级、会话级多维监控

### 1.2 应用场景

**场景1：实时性能监控**
```
运维人员查看Agent监控大盘：
- 当前QPS: 150 req/s
- 响应时间: P50=800ms, P95=1.5s, P99=2.3s
- 错误率: 2.3%
- Token消耗: 15K tokens/min
- 用户满意度: 4.3/5.0
```

**场景2：异常告警**
```
触发告警：
- 指标: 响应时间P99 > 3秒
- 阈值: 持续5分钟
- 告警: 发送钉钉通知 + 邮件
- 处理: 运维人员立即介入排查
```

**场景3：成本追踪**
```
财务人员查看成本监控：
- 本月Token消耗: 10M tokens
- 本月成本: ¥3,000
- Top 3消耗Bot: 客服助手(40%), 文档助手(30%), 数据分析(20%)
- 预算使用率: 60%
```

---

## 2. 数据库设计

### 2.1 Agent性能指标表

```sql
CREATE TABLE agent_performance_metrics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '指标ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    bot_id VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    conversation_id VARCHAR(64) COMMENT '会话ID',

    -- 性能指标
    metric_type ENUM('qps', 'response_time', 'concurrent', 'error_rate', 'token_usage') NOT NULL,
    metric_value DECIMAL(20,6) NOT NULL COMMENT '指标值',
    metric_unit VARCHAR(20) COMMENT '指标单位 (req/s, ms, count, %, tokens)',

    -- 分位数
    percentile VARCHAR(10) COMMENT '分位数 (P50, P95, P99)',

    -- 时间
    timestamp DATETIME NOT NULL COMMENT '时间戳',
    time_window ENUM('1m', '5m', '15m', '1h', '1d') NOT NULL COMMENT '时间窗口',

    -- 标签
    tags JSON COMMENT '额外标签',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_bot_time (tenant_id, bot_id, timestamp),
    INDEX idx_metric_type (metric_type),
    INDEX idx_timestamp (timestamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Agent性能指标表';
```

### 2.2 Agent质量指标表

```sql
CREATE TABLE agent_quality_metrics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    bot_id VARCHAR(64) NOT NULL,

    -- 质量指标
    total_conversations INT DEFAULT 0 COMMENT '总会话数',
    successful_conversations INT DEFAULT 0 COMMENT '成功会话数',
    success_rate DECIMAL(5,2) COMMENT '成功率 (%)',
    avg_rating DECIMAL(3,2) COMMENT '平均评分 (1-5)',
    total_evaluations INT DEFAULT 0 COMMENT '总评价数',

    -- Token统计
    total_input_tokens BIGINT DEFAULT 0 COMMENT '输入Token总数',
    total_output_tokens BIGINT DEFAULT 0 COMMENT '输出Token总数',
    total_tokens BIGINT DEFAULT 0 COMMENT '总Token数',
    avg_tokens_per_conversation INT COMMENT '平均每会话Token数',

    -- 成本统计
    total_cost DECIMAL(12,6) COMMENT '总成本(CNY)',
    avg_cost_per_conversation DECIMAL(10,6) COMMENT '平均每会话成本',

    -- 模型分布
    model_distribution JSON COMMENT '模型使用分布',

    -- 时间
    summary_date DATE NOT NULL COMMENT '汇总日期',
    summary_hour TINYINT COMMENT '小时 (NULL表示日汇总)',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_tenant_bot_date_hour (tenant_id, bot_id, summary_date, summary_hour),
    INDEX idx_tenant_date (tenant_id, summary_date),
    INDEX idx_bot_date (bot_id, summary_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Agent质量指标汇总表';
```

### 2.3 Agent告警规则表

```sql
CREATE TABLE agent_alert_rules (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    bot_id VARCHAR(64) COMMENT 'Bot ID (NULL表示全局规则)',

    -- 规则配置
    rule_name VARCHAR(100) NOT NULL COMMENT '规则名称',
    metric_type ENUM('qps', 'response_time', 'error_rate', 'token_usage', 'rating') NOT NULL,
    condition_type ENUM('gt', 'lt', 'gte', 'lte', 'eq', 'not_eq') NOT NULL,
    threshold_value DECIMAL(20,6) NOT NULL COMMENT '阈值',

    -- 高级配置
    duration_seconds INT DEFAULT 300 COMMENT '持续时间(秒)',
    percentile VARCHAR(10) COMMENT '分位数 (P50/P95/P99)',

    -- 告警级别
    severity ENUM('info', 'warning', 'critical', 'emergency') NOT NULL,

    -- 通知配置
    notification_channels JSON NOT NULL COMMENT '通知渠道',
    notification_cooldown_seconds INT DEFAULT 3600 COMMENT '通知冷却时间(秒)',

    -- 状态
    is_active BOOLEAN DEFAULT TRUE,
    created_by BIGINT COMMENT '创建人ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_tenant_bot (tenant_id, bot_id),
    INDEX idx_metric_type (metric_type),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Agent告警规则表';
```

### 2.4 Agent告警历史表

```sql
CREATE TABLE agent_alert_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    rule_id BIGINT NOT NULL COMMENT '规则ID',
    tenant_id VARCHAR(64) NOT NULL,
    bot_id VARCHAR(64) COMMENT 'Bot ID',

    -- 告警信息
    alert_type VARCHAR(50) NOT NULL COMMENT '告警类型',
    metric_value DECIMAL(20,6) NOT NULL COMMENT '实际值',
    threshold_value DECIMAL(20,6) NOT NULL COMMENT '阈值',
    severity ENUM('info', 'warning', 'critical', 'emergency') NOT NULL,

    -- 状态
    status ENUM('fired', 'acknowledged', 'resolved') DEFAULT 'fired',
    fired_at DATETIME NOT NULL COMMENT '触发时间',
    acknowledged_at DATETIME COMMENT '确认时间',
    acknowledged_by BIGINT COMMENT '确认人ID',
    resolved_at DATETIME COMMENT '解决时间',

    -- 通知
    notification_sent BOOLEAN DEFAULT FALSE,
    notification_channels JSON COMMENT '通知渠道',

    -- 描述
    alert_message TEXT COMMENT '告警消息',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_fired_at (fired_at),
    INDEX idx_severity (severity)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Agent告警历史表';
```

---

## 3. Agent监控指标体系

### 3.1 性能指标

| 指标名称 | 说明 | 单位 | 采集频率 | 正常阈值 |
|---------|------|------|----------|----------|
| **QPS** | 每秒请求数 | req/s | 1分钟 | > 10 |
| **响应时间P50** | 中位数响应时间 | ms | 1分钟 | < 1000 |
| **响应时间P95** | 95分位响应时间 | ms | 1分钟 | < 2000 |
| **响应时间P99** | 99分位响应时间 | ms | 1分钟 | < 3000 |
| **并发数** | 同时处理的请求数 | count | 实时 | < 100 |
| **队列长度** | 等待处理的请求数 | count | 实时 | < 50 |

### 3.2 质量指标

| 指标名称 | 说明 | 单位 | 采集频率 | 正常阈值 |
|---------|------|------|----------|----------|
| **错误率** | 错误请求占比 | % | 1分钟 | < 5% |
| **成功率** | 成功会话占比 | % | 5分钟 | > 90% |
| **用户满意度** | 平均评分 | score/5 | 1小时 | > 4.0 |
| **Token消耗** | 每分钟Token数 | tokens/min | 1分钟 | < 50K |
| **成本** | 每小时成本 | CNY/hour | 1小时 | < ¥500 |

### 3.3 资源指标

| 指标名称 | 说明 | 单位 | 采集频率 | 正常阈值 |
|---------|------|------|----------|----------|
| **模型调用次数** | LLM API调用次数 | count/min | 1分钟 | - |
| **缓存命中率** | 语义缓存命中占比 | % | 5分钟 | > 30% |
| **知识库检索次数** | 向量检索次数 | count/min | 1分钟 | - |

---

## 4. 指标收集服务

### 4.1 指标收集器

```go
package monitoring

import (
    "context"
    "time"
    "github.com/cloudwego/hertz/pkg/common/hlog"
)

// AgentMetricsCollector Agent指标收集器
type AgentMetricsCollector struct {
    perfMetricRepo  repository.PerformanceMetricRepository
    qualityMetricRepo repository.QualityMetricRepository
}

// CollectQPS 收集QPS指标
func (c *AgentMetricsCollector) CollectQPS(
    ctx context.Context,
    tenantID, botID string,
) error {
    // 1. 统计最近1分钟的请求数
    endTime := time.Now()
    startTime := endTime.Add(-1 * time.Minute)

    count, err := c.countRequests(ctx, tenantID, botID, startTime, endTime)
    if err != nil {
        return err
    }

    qps := float64(count) / 60.0 // 每秒请求数

    // 2. 写入指标
    metric := &entity.PerformanceMetric{
        TenantID:     tenantID,
        BotID:        botID,
        MetricType:   "qps",
        MetricValue:  qps,
        MetricUnit:   "req/s",
        Timestamp:    endTime,
        TimeWindow:   "1m",
    }

    return c.perfMetricRepo.Create(ctx, metric)
}

// CollectResponseTime 收集响应时间指标
func (c *AgentMetricsCollector) CollectResponseTime(
    ctx context.Context,
    tenantID, botID string,
) error {
    // 1. 获取最近1分钟的所有响应时间
    endTime := time.Now()
    startTime := endTime.Add(-1 * time.Minute)

    responseTimes, err := c.getResponseTimes(ctx, tenantID, botID, startTime, endTime)
    if err != nil {
        return err
    }

    if len(responseTimes) == 0 {
        return nil
    }

    // 2. 计算分位数
    p50 := calculatePercentile(responseTimes, 50)
    p95 := calculatePercentile(responseTimes, 95)
    p99 := calculatePercentile(responseTimes, 99)

    // 3. 批量写入指标
    metrics := []*entity.PerformanceMetric{
        {
            TenantID:     tenantID,
            BotID:        botID,
            MetricType:   "response_time",
            MetricValue:  p50,
            MetricUnit:   "ms",
            Percentile:   "P50",
            Timestamp:    endTime,
            TimeWindow:   "1m",
        },
        {
            TenantID:     tenantID,
            BotID:        botID,
            MetricType:   "response_time",
            MetricValue:  p95,
            MetricUnit:   "ms",
            Percentile:   "P95",
            Timestamp:    endTime,
            TimeWindow:   "1m",
        },
        {
            TenantID:     tenantID,
            BotID:        botID,
            MetricType:   "response_time",
            MetricValue:  p99,
            MetricUnit:   "ms",
            Percentile:   "P99",
            Timestamp:    endTime,
            TimeWindow:   "1m",
        },
    }

    for _, metric := range metrics {
        c.perfMetricRepo.Create(ctx, metric)
    }

    return nil
}

// CollectErrorRate 收集错误率指标
func (c *AgentMetricsCollector) CollectErrorRate(
    ctx context.Context,
    tenantID, botID string,
) error {
    // 1. 统计最近1分钟的总请求数和错误数
    endTime := time.Now()
    startTime := endTime.Add(-1 * time.Minute)

    totalRequests, _ := c.countRequests(ctx, tenantID, botID, startTime, endTime)
    errorRequests, _ := c.countErrors(ctx, tenantID, botID, startTime, endTime)

    if totalRequests == 0 {
        return nil
    }

    // 2. 计算错误率
    errorRate := float64(errorRequests) / float64(totalRequests) * 100

    // 3. 写入指标
    metric := &entity.PerformanceMetric{
        TenantID:     tenantID,
        BotID:        botID,
        MetricType:   "error_rate",
        MetricValue:  errorRate,
        MetricUnit:   "%",
        Timestamp:    endTime,
        TimeWindow:   "1m",
    }

    return c.perfMetricRepo.Create(ctx, metric)
}

// calculatePercentile 计算分位数
func calculatePercentile(values []int, percentile int) float64 {
    if len(values) == 0 {
        return 0
    }

    // 排序
    sorted := make([]int, len(values))
    copy(sorted, values)
    sort.Ints(sorted)

    // 计算索引
    index := (len(sorted) * percentile) / 100
    if index >= len(sorted) {
        index = len(sorted) - 1
    }

    return float64(sorted[index])
}
```

### 4.2 定时收集任务

```go
package monitoring

import (
    "context"
    "time"
)

// MetricsScheduler 指标收集调度器
type MetricsScheduler struct {
    collector *AgentMetricsCollector
    botRepo   repository.BotRepository
}

// StartCollectionTasks 启动收集任务
func (s *MetricsScheduler) StartCollectionTasks(
    ctx context.Context,
) {
    // 每1分钟收集性能指标
    go s.schedulePerfMetrics(ctx, 1*time.Minute)

    // 每5分钟汇总质量指标
    go s.scheduleQualityMetrics(ctx, 5*time.Minute)

    // 每1小时汇总成本指标
    go s.scheduleCostMetrics(ctx, 1*time.Hour)
}

// schedulePerfMetrics 调度性能指标收集
func (s *MetricsScheduler) schedulePerfMetrics(
    ctx context.Context,
    interval time.Duration,
) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // 获取所有活跃的Bot
            bots, _ := s.botRepo.GetActiveBots(ctx)

            for _, bot := range bots {
                // 并发收集指标
                go func(bot *entity.Bot) {
                    // QPS
                    s.collector.CollectQPS(ctx, bot.TenantID, bot.ID)

                    // 响应时间
                    s.collector.CollectResponseTime(ctx, bot.TenantID, bot.ID)

                    // 错误率
                    s.collector.CollectErrorRate(ctx, bot.TenantID, bot.ID)

                    // Token消耗
                    s.collector.CollectTokenUsage(ctx, bot.TenantID, bot.ID)
                }(bot)
            }

        case <-ctx.Done():
            return
        }
    }
}
```

---

## 5. 告警引擎

### 5.1 规则评估引擎

```go
package monitoring

import (
    "context"
    "time"
)

// AlertEngine 告警引擎
type AlertEngine struct {
    ruleRepo     repository.AlertRuleRepository
    alertRepo    repository.AlertHistoryRepository
    notifier     *NotificationService
}

// EvaluateRules 评估告警规则
func (e *AlertEngine) EvaluateRules(
    ctx context.Context,
    tenantID, botID string,
) error {
    // 1. 获取活跃的告警规则
    rules, err := e.ruleRepo.GetActiveRulesByBot(ctx, tenantID, botID)
    if err != nil {
        return err
    }

    // 2. 评估每个规则
    for _, rule := range rules {
        go e.evaluateRule(ctx, rule)
    }

    return nil
}

// evaluateRule 评估单个规则
func (e *AlertEngine) evaluateRule(
    ctx context.Context,
    rule *entity.AlertRule,
) {
    // 1. 获取当前指标值
    currentValue, err := e.getCurrentMetricValue(
        ctx,
        rule.TenantID,
        rule.BotID,
        rule.MetricType,
        rule.Percentile,
    )
    if err != nil {
        hlog.Errorf("Failed to get metric value: %v", err)
        return
    }

    // 2. 判断是否触发告警
    triggered := e.checkCondition(
        currentValue,
        rule.ThresholdValue,
        rule.ConditionType,
    )

    if !triggered {
        return
    }

    // 3. 检查是否在持续时间内保持触发状态
    if rule.DurationSeconds > 0 {
        sustained := e.checkSustainedCondition(
            ctx,
            rule,
            rule.DurationSeconds,
        )
        if !sustained {
            return
        }
    }

    // 4. 检查是否在冷却期内
    inCooldown, _ := e.checkCooldown(ctx, rule)
    if inCooldown {
        return
    }

    // 5. 触发告警
    e.fireAlert(ctx, rule, currentValue)
}

// checkCondition 检查条件是否满足
func (e *AlertEngine) checkCondition(
    currentValue, thresholdValue float64,
    conditionType string,
) bool {
    switch conditionType {
    case "gt":
        return currentValue > thresholdValue
    case "lt":
        return currentValue < thresholdValue
    case "gte":
        return currentValue >= thresholdValue
    case "lte":
        return currentValue <= thresholdValue
    case "eq":
        return currentValue == thresholdValue
    case "not_eq":
        return currentValue != thresholdValue
    default:
        return false
    }
}

// checkSustainedCondition 检查条件是否持续满足
func (e *AlertEngine) checkSustainedCondition(
    ctx context.Context,
    rule *entity.AlertRule,
    durationSeconds int,
) bool {
    // 查询过去durationSeconds时间内的所有指标值
    endTime := time.Now()
    startTime := endTime.Add(-time.Duration(durationSeconds) * time.Second)

    values, _ := e.getMetricHistory(
        ctx,
        rule.TenantID,
        rule.BotID,
        rule.MetricType,
        startTime,
        endTime,
    )

    // 检查是否所有值都满足条件
    for _, value := range values {
        if !e.checkCondition(value, rule.ThresholdValue, rule.ConditionType) {
            return false
        }
    }

    return true
}

// fireAlert 触发告警
func (e *AlertEngine) fireAlert(
    ctx context.Context,
    rule *entity.AlertRule,
    currentValue float64,
) {
    // 1. 构建告警消息
    alertMessage := e.buildAlertMessage(rule, currentValue)

    // 2. 创建告警历史记录
    alert := &entity.AlertHistory{
        RuleID:         rule.ID,
        TenantID:       rule.TenantID,
        BotID:          rule.BotID,
        AlertType:      rule.MetricType,
        MetricValue:    currentValue,
        ThresholdValue: rule.ThresholdValue,
        Severity:       rule.Severity,
        Status:         "fired",
        FiredAt:        time.Now(),
        AlertMessage:   alertMessage,
        NotificationChannels: rule.NotificationChannels,
    }

    e.alertRepo.Create(ctx, alert)

    // 3. 发送通知
    e.notifier.SendAlert(ctx, alert)

    hlog.Warnf("Alert fired: rule=%d, bot=%s, value=%.2f, threshold=%.2f",
        rule.ID, rule.BotID, currentValue, rule.ThresholdValue)
}

// buildAlertMessage 构建告警消息
func (e *AlertEngine) buildAlertMessage(
    rule *entity.AlertRule,
    currentValue float64,
) string {
    var emoji, unit string

    switch rule.Severity {
    case "info":
        emoji = "ℹ️"
    case "warning":
        emoji = "⚠️"
    case "critical":
        emoji = "🚨"
    case "emergency":
        emoji = "🛑"
    }

    switch rule.MetricType {
    case "qps":
        unit = "req/s"
    case "response_time":
        unit = "ms"
    case "error_rate":
        unit = "%"
    case "token_usage":
        unit = "tokens/min"
    }

    return fmt.Sprintf(
        "%s **%s告警**\n\n"+
        "Bot: %s\n"+
        "指标: %s\n"+
        "当前值: %.2f %s\n"+
        "阈值: %.2f %s\n\n"+
        "请立即处理！",
        emoji,
        rule.Severity,
        rule.BotID,
        rule.MetricType,
        currentValue,
        unit,
        rule.ThresholdValue,
        unit,
    )
}
```

### 5.2 智能异常检测

```go
package monitoring

import (
    "context"
    "math"
    "time"
)

// AnomalyDetector 异常检测器
type AnomalyDetector struct {
    metricRepo repository.MetricRepository
}

// DetectAnomaly 检测异常（基于标准差）
func (d *AnomalyDetector) DetectAnomaly(
    ctx context.Context,
    tenantID, botID, metricType string,
) (*AnomalyResult, error) {
    // 1. 获取最近7天的历史数据
    endTime := time.Now()
    startTime := endTime.Add(-7 * 24 * time.Hour)

    values, _ := d.metricRepo.GetHistory(
        ctx,
        tenantID,
        botID,
        metricType,
        startTime,
        endTime,
    )

    if len(values) < 10 {
        return nil, errors.New("insufficient data")
    }

    // 2. 计算均值和标准差
    mean, stdDev := d.calculateMeanAndStdDev(values)

    // 3. 获取当前值
    currentValue, _ := d.getCurrentValue(ctx, tenantID, botID, metricType)

    // 4. 计算Z分数
    zScore := math.Abs((currentValue - mean) / stdDev)

    // 5. 判断是否异常（Z分数 > 3）
    isAnomaly := zScore > 3

    return &AnomalyResult{
        MetricType:     metricType,
        CurrentValue:   currentValue,
        Mean:           mean,
        StdDev:         stdDev,
        ZScore:         zScore,
        IsAnomaly:      isAnomaly,
        AnomalyLevel:   d.getAnomalyLevel(zScore),
        DetectedAt:     time.Now(),
    }, nil
}

// calculateMeanAndStdDev 计算均值和标准差
func (d *AnomalyDetector) calculateMeanAndStdDev(
    values []float64,
) (mean, stdDev float64) {
    n := len(values)

    // 计算均值
    var sum float64
    for _, v := range values {
        sum += v
    }
    mean = sum / float64(n)

    // 计算标准差
    var variance float64
    for _, v := range values {
        diff := v - mean
        variance += diff * diff
    }
    variance /= float64(n)
    stdDev = math.Sqrt(variance)

    return
}

// getAnomalyLevel 获取异常等级
func (d *AnomalyDetector) getAnomalyLevel(zScore float64) string {
    if zScore > 5 {
        return "extreme" // 极端异常
    } else if zScore > 4 {
        return "severe" // 严重异常
    } else if zScore > 3 {
        return "moderate" // 中度异常
    }
    return "normal"
}
```

---

## 6. 监控仪表盘

### 6.1 仪表盘配置

```json
{
  "dashboard_name": "Agent监控大盘",
  "panels": [
    {
      "panel_id": 1,
      "panel_title": "QPS趋势",
      "panel_type": "line_chart",
      "metrics": ["qps"],
      "group_by": ["bot_id"],
      "time_range": "last_1h",
      "refresh_interval": "1m"
    },
    {
      "panel_id": 2,
      "panel_title": "响应时间分布",
      "panel_type": "heatmap",
      "metrics": ["response_time_p50", "response_time_p95", "response_time_p99"],
      "time_range": "last_1h",
      "refresh_interval": "1m"
    },
    {
      "panel_id": 3,
      "panel_title": "错误率Top 10",
      "panel_type": "bar_chart",
      "metrics": ["error_rate"],
      "group_by": ["bot_id"],
      "order_by": "error_rate_desc",
      "limit": 10,
      "time_range": "last_24h",
      "refresh_interval": "5m"
    },
    {
      "panel_id": 4,
      "panel_title": "用户满意度",
      "panel_type": "gauge",
      "metrics": ["avg_rating"],
      "time_range": "last_7d",
      "refresh_interval": "1h"
    },
    {
      "panel_id": 5,
      "panel_title": "成本趋势",
      "panel_type": "line_chart",
      "metrics": ["total_cost"],
      "time_range": "last_30d",
      "refresh_interval": "1h"
    }
  ]
}
```

---

## 7. API设计

### 7.1 监控指标API

| 方法 | 路径 | 功能 | 请求示例 | 响应示例 |
|------|------|------|----------|----------|
| GET | /api/v1/monitoring/metrics | 查询指标 | `?bot_id=bot123&metric_type=qps&start_time=...&end_time=...` | 见下方完整响应 |
| GET | /api/v1/monitoring/dashboard | 获取仪表盘数据 | `?dashboard_id=1` | `{panels:[{panel_id:1,data:...}]}` |
| GET | /api/v1/monitoring/alerts | 查询告警 | `?status=fired&limit=20` | `[{rule_id:1,severity:"critical",...}]` |

**指标查询响应**:
```json
{
  "bot_id": "bot123",
  "metric_type": "qps",
  "data_points": [
    {
      "timestamp": "2025-01-03T10:00:00Z",
      "value": 150.5
    },
    {
      "timestamp": "2025-01-03T10:01:00Z",
      "value": 162.3
    }
  ],
  "summary": {
    "min": 120.0,
    "max": 200.0,
    "avg": 158.5,
    "p95": 185.0,
    "p99": 195.0
  }
}
```

### 7.2 告警管理API

| 方法 | 路径 | 功能 | 请求示例 | 响应示例 |
|------|------|------|----------|----------|
| GET | /api/v1/monitoring/alert-rules | 查询告警规则 | `?bot_id=bot123` | `[{rule_id:1,rule_name:"QPS告警",...}]` |
| POST | /api/v1/monitoring/alert-rules | 创建告警规则 | `{"metric_type":"qps","condition_type":"lt","threshold_value":10,...}` | `{"rule_id":123}` |
| PUT | /api/v1/monitoring/alert-rules/:id | 更新告警规则 | `{"is_active":false}` | `{"success":true}` |
| POST | /api/v1/monitoring/alerts/:id/acknowledge | 确认告警 | - | `{"success":true}` |

---

## 8. 前端组件

### 8.1 监控仪表盘组件

```typescript
// components/monitoring/MonitoringDashboard.tsx
import React, { useEffect, useState } from 'react';
import { Card, Row, Col, Metric } from '@douyinfe/semi-ui';
import { Line, Heatmap, Gauge } from '@douyinfe/semi-charts';

interface MetricData {
  timestamp: string;
  value: number;
}

export const MonitoringDashboard: React.FC<Props> = ({ botId, timeRange }) => {
  const [metrics, setMetrics] = useState<{
    qps: MetricData[];
    responseTime: MetricData[];
    errorRate: MetricData[];
    avgRating: number;
    totalCost: MetricData[];
  } | null>(null);

  useEffect(() => {
    fetchMetrics();
    const interval = setInterval(fetchMetrics, 60000); // 每分钟刷新
    return () => clearInterval(interval);
  }, [botId, timeRange]);

  const fetchMetrics = async () => {
    // 并发获取多个指标
    const [qps, responseTime, errorRate, rating, cost] = await Promise.all([
      fetch(`/api/v1/monitoring/metrics?bot_id=${botId}&metric_type=qps&time_range=${timeRange}`).then(r => r.json()),
      fetch(`/api/v1/monitoring/metrics?bot_id=${botId}&metric_type=response_time&time_range=${timeRange}`).then(r => r.json()),
      fetch(`/api/v1/monitoring/metrics?bot_id=${botId}&metric_type=error_rate&time_range=${timeRange}`).then(r => r.json()),
      fetch(`/api/v1/monitoring/metrics?bot_id=${botId}&metric_type=rating&time_range=${timeRange}`).then(r => r.json()),
      fetch(`/api/v1/monitoring/metrics?bot_id=${botId}&metric_type=cost&time_range=${timeRange}`).then(r => r.json()),
    ]);

    setMetrics({ qps, responseTime, errorRate, avgRating: rating.avg_rating, totalCost: cost });
  };

  if (!metrics) return <Spin />;

  return (
    <div className="monitoring-dashboard">
      {/* 关键指标卡片 */}
      <Row gutter={[16, 16]}>
        <Col span={6}>
          <Card title="当前QPS">
            <Metric
              value={metrics.qps.data_points[metrics.qps.data_points.length - 1].value}
              fmt="value"
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card title="响应时间P95">
            <Metric
              value={metrics.responseTime.summary.p95}
              suffix="ms"
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card title="错误率">
            <Metric
              value={metrics.errorRate.data_points[metrics.errorRate.data_points.length - 1].value}
              suffix="%"
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card title="用户满意度">
            <Metric
              value={metrics.avgRating}
              suffix="/5.0"
            />
          </Card>
        </Col>
      </Row>

      {/* QPS趋势图 */}
      <Card title="QPS趋势" style={{ marginTop: 16 }}>
        <Line
          data={metrics.qps.data_points.map(p => ({
            time: p.timestamp,
            value: p.value,
          }))}
          xField="time"
          yField="value"
          yAxis={{ label: { formatter: (v) => `${v} req/s` } }}
        />
      </Card>

      {/* 响应时间热力图 */}
      <Card title="响应时间分布" style={{ marginTop: 16 }}>
        <Heatmap
          data={buildHeatmapData(metrics.responseTime)}
          xField="time"
          yField="percentile"
          valueField="value"
        />
      </Card>

      {/* 成本趋势 */}
      <Card title="成本趋势" style={{ marginTop: 16 }}>
        <Line
          data={metrics.totalCost.data_points.map(p => ({
            time: p.timestamp,
            value: p.value,
          }))}
          xField="time"
          yField="value"
          yAxis={{ label: { formatter: (v) => `¥${v}` } }}
        />
      </Card>
    </div>
  );
};
```

---

## 9. 实施计划

### 9.1 开发阶段划分

| 阶段 | 任务 | 工期 | 交付物 |
|------|------|------|--------|
| **第1周** | 数据库设计与创建 | 2天 | 4张表 |
| | 指标收集服务开发 | 3天 | 收集器+调度器 |
| **第2周** | 告警引擎开发 | 2天 | 规则评估+异常检测 |
| | 前端仪表盘开发 | 3天 | 监控大盘组件 |

### 9.2 技术风险与缓解

| 风险 | 影响 | 概率 | 缓解措施 |
|------|------|------|----------|
| 指标数据量爆炸 | 高 | 高 | 使用InfluxDB时序数据库+数据降采样 |
| 告警风暴 | 中 | 中 | 告警聚合+冷却期 |
| 性能损耗 | 中 | 低 | 异步采集+批量写入 |

### 9.3 成功指标

- ✅ 指标采集延迟: **< 5秒**
- ✅ 告警响应时间: **< 30秒**
- ✅ 仪表盘刷新: **实时** (1分钟)
- ✅ 数据存储成本: **< ¥500/月** (100万指标/天)

---

**文档结束**
