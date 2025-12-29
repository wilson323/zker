# 24-MultiTenant_SaaS核心_租户监控运维 详细设计说明书

**文档编号**: DE-DD-2025-024
**模块名称**: 租户监控运维 (TenantMonitoring)
**版本**: v1.0.0
**作者**: ZKER Enterprise Team
**创建日期**: 2025-01-03

---

## 1. 模块概述

**租户监控运维** 是 MultiTenant SaaS 的核心运维模块，通过**多维度监控 + 自动告警**，确保多租户系统的稳定性、性能和安全性。

**核心设计理念**：
- ✅ **多租户隔离监控**：独立监控每个租户的资源使用
- ✅ **实时告警**：支持多种告警规则和通知渠道
- ✅ **可视化**：租户级监控大盘

**实现策略**：✅ 50% 编码（监控框架） + 50% 配置（告警规则）

---

## 2. 核心功能

### 2.1 监控维度

**F1 - 性能监控**
- F1.1 API响应时间（P50/P95/P99）
- F1.2 错误率
- F1.3 并发数
- F1.4 队列长度

**F2 - 资源监控**
- F2.1 CPU使用率
- F2.2 内存使用率
- F2.3 磁盘IO
- F2.4 网络带宽

**F3 - 业务监控**
- F3.1 活跃用户数
- F3.2 API调用量
- F3.3 Token消耗量
- F3.4 Bot使用量

### 2.2 告警管理

**F4 - 告警规则**
- F4.1 阈值告警
- F4.2 趋势告警
- F4.3 异常检测

**F5 - 通知渠道**
- F5.1 邮件通知
- F5.2 短信通知
- F5.3 Webhook
- F5.4 企业微信/钉钉

---

## 3. 数据库设计

### 3.1 核心表结构

#### 3.1.1 监控指标表 (monitoring_metrics)

```sql
CREATE TABLE monitoring_metrics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '指标ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    metric_name VARCHAR(50) NOT NULL COMMENT '指标名称',
    metric_value DECIMAL(20,6) NOT NULL COMMENT '指标值',
    metric_tags JSON COMMENT '指标标签',
    timestamp BIGINT NOT NULL COMMENT '时间戳',

    INDEX idx_tenant_metric_time (tenant_id, metric_name, timestamp),
    INDEX idx_timestamp (timestamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='监控指标表';
```

#### 3.1.2 告警规则表 (alert_rules)

```sql
CREATE TABLE alert_rules (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '规则ID',
    tenant_id VARCHAR(64) COMMENT '租户ID (NULL表示平台级)',
    name VARCHAR(100) NOT NULL COMMENT '规则名称',
    description VARCHAR(500) COMMENT '规则描述',

    -- 规则配置
    metric_name VARCHAR(50) NOT NULL COMMENT '监控指标',
    condition ENUM('gt', 'lt', 'eq', 'gte', 'lte') NOT NULL COMMENT '条件',
    threshold DECIMAL(20,6) NOT NULL COMMENT '阈值',
    duration INT DEFAULT 60 COMMENT '持续时间(秒)',

    -- 通知配置
    notification_channels JSON NOT NULL COMMENT '通知渠道',
    severity ENUM('info', 'warning', 'critical') NOT NULL COMMENT '严重程度',

    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_metric_name (metric_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='告警规则表';
```

**通知渠道 JSON 示例**：

```json
{
  "email": ["admin@example.com"],
  "sms": ["+86138xxxxxxxx"],
  "webhook": ["https://hooks.example.com/alert"],
  "wechat": {
    "corpId": "ww123456",
    "agentId": 123456,
    "toUser": ["admin"]
  }
}
```

#### 3.1.3 告警历史表 (alert_history)

```sql
CREATE TABLE alert_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '告警ID',
    rule_id BIGINT NOT NULL COMMENT '规则ID',
    tenant_id VARCHAR(64) COMMENT '租户ID',
    metric_name VARCHAR(50) NOT NULL COMMENT '指标名称',
    actual_value DECIMAL(20,6) NOT NULL COMMENT '实际值',
    threshold DECIMAL(20,6) NOT NULL COMMENT '阈值',
    severity ENUM('info', 'warning', 'critical') NOT NULL COMMENT '严重程度',
    status ENUM('fired', 'resolved', 'acknowledged') DEFAULT 'fired',
    fired_at DATETIME NOT NULL COMMENT '触发时间',
    resolved_at DATETIME COMMENT '解决时间',
    acknowledged_by BIGINT COMMENT '确认人ID',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status),
    INDEX idx_fired_at (fired_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='告警历史表';
```

---

## 4. 监控框架设计

### 4.1 指标收集器

```go
package monitoring

// MetricsCollector 指标收集器
type MetricsCollector struct {
    influxdbClient *influxdb.Client
    logger         *zap.Logger
}

// CollectMetric 收集指标
func (c *MetricsCollector) CollectMetric(
    ctx context.Context,
    tenantID, metricName string,
    value float64,
    tags map[string]string,
) error {
    metric := &Metric{
        TenantID:    tenantID,
        MetricName:  metricName,
        MetricValue: value,
        MetricTags:  tags,
        Timestamp:   time.Now().Unix(),
    }

    // 写入 InfluxDB
    point := influxdb.NewPoint(
        "tenant_metrics",
        tags,
        map[string]interface{}{"value": value},
        time.Now(),
    )

    return c.influxdbClient.Write(point)
}

// CollectTenantMetrics 收集租户所有指标
func (c *MetricsCollector) CollectTenantMetrics(
    ctx context.Context,
    tenantID string,
) error {
    // 1. API调用量
    apiCalls := c.getAPICalls(ctx, tenantID)
    c.CollectMetric(ctx, tenantID, "api_calls", apiCalls, nil)

    // 2. 活跃用户数
    activeUsers := c.getActiveUsers(ctx, tenantID)
    c.CollectMetric(ctx, tenantID, "active_users", activeUsers, nil)

    // 3. Token消耗
    tokens := c.getTokenUsage(ctx, tenantID)
    c.CollectMetric(ctx, tenantID, "token_usage", tokens, nil)

    return nil
}
```

### 4.2 告警引擎

```go
package monitoring

// AlertEngine 告警引擎
type AlertEngine struct {
    ruleRepo    repository.AlertRuleRepository
    alertRepo   repository.AlertHistoryRepository
    notifier    *NotificationService
    evaluator   *RuleEvaluator
}

// EvaluateRules 评估告警规则
func (e *AlertEngine) EvaluateRules(
    ctx context.Context,
    tenantID, metricName string,
    value float64,
) error {
    // 1. 获取相关规则
    rules, err := e.ruleRepo.GetActiveRulesByMetric(ctx, tenantID, metricName)
    if err != nil {
        return err
    }

    // 2. 评估每条规则
    for _, rule := range rules {
        triggered := e.evaluator.Evaluate(rule, value)
        if triggered {
            // 触发告警
            go e.fireAlert(ctx, rule, value)
        }
    }

    return nil
}

// fireAlert 触发告警
func (e *AlertEngine) fireAlert(
    ctx context.Context,
    rule *AlertRule,
    value float64,
) error {
    // 1. 创建告警记录
    alert := &Alert{
        RuleID:      rule.ID,
        TenantID:    rule.TenantID,
        MetricName:  rule.MetricName,
        ActualValue: value,
        Threshold:   rule.Threshold,
        Severity:    rule.Severity,
        Status:      "fired",
        FiredAt:     time.Now(),
    }

    if err := e.alertRepo.Create(ctx, alert); err != nil {
        return err
    }

    // 2. 发送通知
    return e.notifier.SendNotification(ctx, alert, rule.NotificationChannels)
}
```

---

## 5. 配置系统设计

### 5.1 预置告警规则

```sql
-- 平台级告警规则
INSERT INTO alert_rules (tenant_id, name, metric_name, condition, threshold, duration, severity, notification_channels) VALUES
-- 性能告警
(NULL, 'API响应时间过长', 'api_response_time_p99', 'gt', 5000, 300, 'critical',
'{"email": ["ops@zker.com"],"webhook": ["https://hooks.example.com/critical"]}'),

(NULL, '错误率过高', 'api_error_rate', 'gt', 0.05, 60, 'warning',
'{"email": ["ops@zker.com"]}'),

-- 资源告警
(NULL, '磁盘使用率过高', 'disk_usage_percent', 'gt', 80, 300, 'warning',
'{"email": ["ops@zker.com"]}'),

(NULL, '内存使用率过高', 'memory_usage_percent', 'gt', 85, 300, 'warning',
'{"email": ["ops@zker.com"]});
```

### 5.2 企业自定义告警规则

```sql
-- 企业级自定义告警
INSERT INTO alert_rules (tenant_id, name, metric_name, condition, threshold, notification_channels) VALUES
('tenant-abc-001', 'Bot调用次数异常', 'bot_calls_per_minute', 'gt', 1000,
'{"wechat": {"corpId": "ww12345","agentId": 123,"toUser": ["admin"]}}');
```

---

## 6. 监控大盘

### 6.1 租户监控大盘（Grafana）

**Grafana Dashboard JSON 配置**：

```json
{
  "dashboard": {
    "title": "租户监控大盘",
    "panels": [
      {
        "id": 1,
        "title": "API调用量(Top 10)",
        "type": "graph",
        "targets": [
          {
            "query": "SELECT sum(\"value\") FROM \"tenant_metrics\" WHERE \"metric_name\" = 'api_calls' GROUP BY \"tenant_id\" ORDER BY time DESC LIMIT 10"
          }
        ]
      },
      {
        "id": 2,
        "title": "错误率",
        "type": "graph",
        "targets": [
          {
            "query": "SELECT mean(\"value\") FROM \"tenant_metrics\" WHERE \"metric_name\" = 'api_error_rate' GROUP BY time(5m)"
          }
        ]
      },
      {
        "id": 3,
        "title": "活跃租户数",
        "type": "stat",
        "targets": [
          {
            "query": "SELECT count(DISTINCT(\"tenant_id\")) FROM \"tenant_metrics\" WHERE time > now() - 1h"
          }
        ]
      }
    ]
  }
}
```

---

## 7. API 设计

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | /api/v1/monitoring/metrics | 查询监控指标 |
| GET | /api/v1/monitoring/alerts | 列出告警 |
| GET | /api/v1/monitoring/alert-rules | 列出告警规则 |
| POST | /api/v1/monitoring/alert-rules | 创建告警规则 |
| PUT | /api/v1/monitoring/alert-rules/:id | 更新告警规则 |
| POST | /api/v1/monitoring/alerts/:id/acknowledge | 确认告警 |

---

## 8. 总结

### 8.1 实施策略总结

| 实施项 | 实施方式 | 工作量 |
|-------|---------|--------|
| **监控框架** | 💻 独立编码 | 50% |
| **告警规则** | 📊 配置驱动 | 50% |

**总计**：50% 编码 + 50% 配置

### 8.2 核心优势

- ✅ **多租户隔离**：独立监控每个租户的资源使用
- ✅ **实时告警**：支持多种通知渠道
- ✅ **可视化**：Grafana 集成

---

**文档结束**
