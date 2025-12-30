# 计费系统监控和日志使用指南

**版本**: v1.0.0
**最后更新**: 2025-01-01

---

## 📋 目录

- [监控概览](#监控概览)
- [快速开始](#快速开始)
- [Prometheus指标](#prometheus指标)
- [Grafana Dashboard](#grafana-dashboard)
- [告警规则](#告警规则)
- [ELK日志](#elk日志)
- [Jaeger追踪](#jaeger追踪)
- [使用示例](#使用示例)
- [故障排查](#故障排查)

---

## 🎯 监控概览

### 监控组件

计费系统包含以下监控组件：

1. **Prometheus指标** - 实时指标采集和存储
2. **Grafana Dashboard** - 可视化监控面板
3. **Alertmanager** - 告警规则和通知
4. **ELK Stack** - 日志收集、存储和分析
5. **Jaeger** - 分布式追踪

### 监控覆盖范围

| 模块 | 指标 | 日志 | 追踪 | 告警 |
|------|------|------|------|------|
| Token计量 | ✅ | ✅ | ✅ | ✅ |
| 定价引擎 | ✅ | ✅ | ✅ | ✅ |
| 预算告警 | ✅ | ✅ | ✅ | ✅ |
| 通知服务 | ✅ | ✅ | ✅ | ✅ |
| 汇总统计 | ✅ | ✅ | ✅ | ✅ |

---

## 🚀 快速开始

### 1. 前置条件

确保以下组件已安装并运行：

```bash
# Prometheus（指标采集）
prometheus --config.file=/etc/prometheus/prometheus.yml

# Grafana（可视化）
grafana-server --config=/etc/grafana/grafana.ini

# Elasticsearch + Logstash + Kibana（日志）
systemctl start elasticsearch logstash kibana

# Jaeger（追踪）
jaeger-all-in-one
```

### 2. 导入Grafana Dashboard

**方式一：通过UI导入**

1. 登录 Grafana：`http://localhost:3000`
2. 导航到 **Dashboards** → **Import**
3. 上传 JSON 文件：`deploy/monitoring/dashboards/billing-dashboard.json`

**方式二：通过API导入**

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d @deploy/monitoring/dashboards/billing-dashboard.json \
  http://admin:admin@localhost:3000/api/dashboards/db
```

### 3. 配置告警规则

将告警规则复制到Prometheus配置目录：

```bash
cp deploy/monitoring/alerts/billing-alerts.yml /etc/prometheus/alerts/
promtool check rules /etc/prometheus/alerts/billing-alerts.yml
systemctl reload prometheus
```

### 4. 配置Filebeat日志采集

```bash
cp deploy/logging/filebeat-billing.yml /etc/filebeat/modules.d/
filebeat modules enable billing
systemctl restart filebeat
```

---

## 📊 Prometheus指标

### 指标分类

#### 1. Token计量指标

```promql
# Token记录速率
rate(billing_token_records_total[5m])

# Token记录P95耗时
histogram_quantile(0.95, rate(billing_token_record_duration_seconds_bucket[5m]))

# Token处理总量
sum(billing_tokens_total) by (tenant_id, token_type)

# 成本总额
sum(billing_cost_total) by (tenant_id)
```

#### 2. 定价引擎指标

```promql
# 定价计算速率
rate(billing_pricing_calculations_total[5m])

# 定价计算错误率
rate(billing_pricing_errors_total[5m]) / rate(billing_pricing_calculations_total[5m])

# 单价
billing_unit_price{model_provider="openai", model_name="gpt-4"}
```

#### 3. 预算告警指标

```promql
# 预算使用率
billing_budget_usage_percent{tenant_id="tenant-001"}

# 预算告警数量
sum(increase(billing_budget_alerts_total[1h])) by (alert_level)

# 预算检查速率
rate(billing_budget_checks_total[5m])
```

### 完整指标列表

| 指标名称 | 类型 | 标签 | 说明 |
|---------|------|------|------|
| `billing_token_records_total` | Counter | tenant_id, model_provider, model_name | Token记录总数 |
| `billing_token_record_duration_seconds` | Histogram | operation | Token记录耗时 |
| `billing_tokens_total` | Counter | tenant_id, token_type | Token总数 |
| `billing_cost_total` | Counter | tenant_id, model_provider | 成本总额 |
| `billing_pricing_calculations_total` | Counter | model_provider, model_name, result | 定价计算次数 |
| `billing_pricing_errors_total` | Counter | model_provider, model_name, error_type | 定价错误数 |
| `billing_budget_checks_total` | Counter | tenant_id, triggered | 预算检查次数 |
| `billing_budget_alerts_total` | Counter | tenant_id, alert_type, alert_level | 预算告警次数 |
| `billing_budget_usage_percent` | Gauge | tenant_id, budget_type | 预算使用率 |
| `billing_cached_requests_total` | Counter | tenant_id, model_provider | 缓存请求数 |

---

## 📈 Grafana Dashboard

### Dashboard概览

计费系统Dashboard包含15个面板：

1. **Token记录速率** - 每秒Token记录次数
2. **Token处理总量** - Input/Output Token总数
3. **成本趋势** - 按租户统计成本
4. **预算使用率** - 实时预算使用百分比（Gauge）
5. **预算告警数量** - 告警级别分布
6. **Token记录耗时分布** - 热力图
7. **Token记录P95耗时** - 时间序列图
8. **定价计算速率** - 每秒计算次数
9. **定价计算错误率** - 百分比
10. **缓存命中率** - 统计卡片
11. **活跃租户数** - 实时统计
12. **总收入** - 日收入统计
13. **预算检查速率** - 每秒检查次数
14. **按模型统计成本** - 饼图
15. **按租户统计Token使用** - 表格

### 阈值配置

| 面板 | 警告阈值 | 严重阈值 |
|------|---------|---------|
| Token记录速率 | > 1000次/5秒 | > 5000次/5秒 |
| Token记录P95耗时 | > 500ms | > 1000ms |
| 预算使用率 | > 80% | > 95% |
| 定价计算错误率 | > 5% | > 10% |
| 缓存命中率 | < 30% | < 20% |

### Dashboard查询示例

```promql
# 查询特定租户的Token使用
sum(billing_tokens_total{tenant_id="tenant-001"}) by (token_type)

# 查询最近1小时的成本趋势
sum(rate(billing_cost_total[1h])) by (tenant_id)

# 查询预算告警趋势
sum(increase(billing_budget_alerts_total[1h])) by (alert_level)

# 查询Top 10成本最高的租户
topk(10, sum(billing_cost_total) by (tenant_id))
```

---

## 🚨 告警规则

### 告警分类

#### 1. Token计量告警

| 告警名称 | 级别 | 触发条件 | 持续时间 |
|---------|------|---------|---------|
| `HighTokenRecordRate` | warning | 速率 > 1000次/5秒 | 5分钟 |
| `TokenRecordSlow` | warning | P95耗时 > 500ms | 10分钟 |
| `HighTokenRecordFailureRate` | warning | 失败率 > 5% | 5分钟 |

#### 2. 定价引擎告警

| 告警名称 | 级别 | 触发条件 | 持续时间 |
|---------|------|---------|---------|
| `HighPricingErrorRate` | warning | 错误率 > 5% | 5分钟 |
| `PricingCalculationSlow` | warning | P95耗时 > 10ms | 10分钟 |
| `PricingModelNotFound` | critical | 模型缺失 | 1分钟 |

#### 3. 预算告警

| 告警名称 | 级别 | 触发条件 | 持续时间 |
|---------|------|---------|---------|
| `BudgetNearLimit` | warning | 使用率 > 80% | 1分钟 |
| `BudgetCritical` | critical | 使用率 > 95% | 1分钟 |
| `BudgetExceeded` | emergency | 使用率 ≥ 100% | 0分钟 |

#### 4. 成本告警

| 告警名称 | 级别 | 触发条件 | 持续时间 |
|---------|------|---------|---------|
| `TenantCostSpike` | warning | 成本突增 > 2倍 | 5分钟 |
| `HighModelCost` | warning | 模型成本 > 1000元/小时 | 10分钟 |
| `RevenueDrop` | warning | 收入下降 > 50% | 10分钟 |

### 告警通知

告警通过以下渠道发送：

1. **Email** - 所有告警
2. **Webhook** - 严重及以上级别
3. **企业微信/钉钉** - 警告及以上级别
4. **短信/电话** - 严重级别

### 告警静音规则

```yaml
# 工作时间静音非紧急告警（周一至周五 9:00-18:00）
mute_alerts:
  - severity: warning
    schedule: "0 9 * * 1-5"
    duration: "9h"

# 维护窗口静音所有告警
maintenance_window:
  - start: "2025-01-01T02:00:00Z"
    end: "2025-01-01T04:00:00Z"
    mute_all: true
```

---

## 📝 ELK日志

### 日志索引

计费系统日志索引：

- `coze-studio-billing-token-metering-*` - Token计量日志
- `coze-studio-billing-budget-alert-*` - 预算告警日志
- `coze-studio-billing-pricing-engine-*` - 定价引擎日志
- `coze-studio-billing-notification-*` - 通知服务日志
- `coze-studio-billing-*` - 默认索引

### 日志格式

```json
{
  "@timestamp": "2025-01-01T12:00:00.000Z",
  "service": "billing",
  "module": "token_metering",
  "environment": "production",
  "level": "Info",
  "tenant_id": "tenant-001",
  "message": "[TokenMetering] Record completed: tenant=tenant-001 model=openai/gpt-4 tokens=100/200/300 cost=0.0045 duration=15ms",
  "tokens": {
    "input": 100,
    "output": 200,
    "total": 300
  },
  "cost": {
    "input": 0.0015,
    "output": 0.0030,
    "total": 0.0045
  },
  "trace_id": "abc123def456",
  "span_id": "789xyz"
}
```

### Kibana查询示例

```json
// 查询特定租户的所有日志
{
  "query": {
    "term": { "tenant_id": "tenant-001" }
  }
}

// 查询错误日志
{
  "query": {
    "match": { "level": "Error" }
  }
}

// 查询慢操作日志
{
  "query": {
    "range": {
      "duration_ms": { "gte": 1000 }
    }
  }
}

// 查询特定时间范围的日志
{
  "query": {
    "range": {
      "@timestamp": {
        "gte": "now-1h",
        "lte": "now"
      }
    }
  }
}
```

### 日志保留策略

| 索引模式 | 保留时间 | 说明 |
|---------|---------|------|
| `coze-studio-billing-*` | 30天 | 生产环境 |
| `coze-studio-billing-*` | 7天 | 测试环境 |
| `coze-studio-billing-alert-*` | 90天 | 告警日志保留更久 |

---

## 🔍 Jaeger追踪

### 追踪架构

计费系统采用OpenTelemetry进行分布式追踪：

```
API请求 → TokenMetering.Record
            ↓
         PricingEngine.Calculate
            ↓
         DB.Query (token_usage_logs)
            ↓
         TokenMetering.UpdateSummary
            ↓
         BudgetAlertService.Check
```

### Span命名规范

| Span名称 | 父Span | 说明 |
|---------|--------|------|
| `TokenMetering.Record` | HTTP请求 | Token记录根操作 |
| `PricingEngine.Calculate` | Record | 定价计算 |
| `DB.Query` | Record | 数据库查询 |
| `TokenMetering.UpdateSummary` | Record | 汇总更新 |
| `BudgetAlertService.Check` | Record | 预算检查 |
| `NotificationService.SendAlert` | Check | 告警通知 |

### 查看追踪

1. 访问 Jaeger UI：`http://localhost:16686`
2. 选择服务：`coze-studio-api`
3. 搜索Trace ID（从日志中获取）
4. 查看Span详情和性能数据

### 性能分析

```python
# 查找慢Trace
SELECT
  traceID,
  duration / 1000000 as duration_ms,
  spanName
FROM jaeger_spans
WHERE duration > 1000000000  # 1秒
  AND serviceName = 'coze-studio-api'
  AND operationName LIKE '%TokenMetering%'
ORDER BY duration DESC
LIMIT 100
```

---

## 💡 使用示例

### 示例1：监控Token记录性能

```go
// 在代码中使用监控集成
monitoring := NewMonitoringIntegration(tokenMeteringService)

// 记录Token使用（自动集成监控）
result, err := monitoring.RecordTokenWithMonitoring(ctx, &RecordTokenUsageRequest{
    TenantID:       "tenant-001",
    ModelProvider:  "openai",
    ModelName:      "gpt-4",
    InputTokens:    100,
    OutputTokens:   200,
    TotalTokens:    300,
    RequestType:    "chat",
    IsCached:       false,
})
```

### 示例2：查询Prometheus指标

```bash
# 查询Token记录速率
curl http://localhost:9090/api/v1/query?query=rate(billing_token_records_total[5m])

# 查询特定租户的成本
curl http://localhost:9090/api/v1/query?query=billing_cost_total{tenant_id="tenant-001"}

# 查询预算使用率
curl http://localhost:9090/api/v1/query?query=billing_budget_usage_percent
```

### 示例3：在Kibana中分析日志

```json
POST /coze-studio-billing-token-metering-*/_search
{
  "size": 100,
  "query": {
    "bool": {
      "must": [
        { "term": { "tenant_id": "tenant-001" } },
        { "range": { "@timestamp": { "gte": "now-1h" } } },
        { "term": { "level": "Error" } }
      ]
    }
  },
  "sort": [
    { "@timestamp": { "order": "desc" } }
  ]
}
```

### 示例4：创建自定义告警

```yaml
# 自定义告警：单个Token记录成本过高
- alert: HighSingleTokenCost
  expr: billing_cost_total / billing_tokens_total > 0.1
  for: 5m
  labels:
    severity: warning
    category: billing
  annotations:
    summary: "单个Token成本过高"
    description: "租户 {{ $labels.tenant_id }} 单Token平均成本超过0.1元"
```

---

## 🔧 故障排查

### 问题1：Prometheus指标未出现

**症状**：Grafana Dashboard显示"No Data"

**排查步骤**：

1. 检查Prometheus是否正常运行
```bash
curl http://localhost:9090/-/healthy
```

2. 检查指标端点是否可访问
```bash
curl http://localhost:8080/metrics | grep billing_
```

3. 检查Prometheus配置
```bash
promtool check config /etc/prometheus/prometheus.yml
```

4. 检查服务发现
```bash
curl http://localhost:9090/api/v1/targets
```

### 问题2：告警未触发

**症状**：达到阈值但未收到告警通知

**排查步骤**：

1. 检查Alertmanager状态
```bash
curl http://localhost:9093/api/v1/alerts
```

2. 检查告警规则
```bash
promtool check rules /etc/prometheus/alerts/billing-alerts.yml
```

3. 检查告警路由配置
```bash
curl http://localhost:9093/api/v1/receivers
```

4. 查看Alertmanager日志
```bash
journalctl -u alertmanager -f
```

### 问题3：日志未发送到Elasticsearch

**症状**：Kibana中无计费系统日志

**排查步骤**：

1. 检查Filebeat状态
```bash
filebeat test output
filebeat test config
```

2. 检查Elasticsearch索引
```bash
curl http://localhost:9200/_cat/indices?v | grep billing
```

3. 检查日志文件权限
```bash
ls -la /var/log/coze-studio/billing/
```

4. 查看Filebeat日志
```bash
journalctl -u filebeat -f
```

### 问题4：Jaeger追踪无数据

**症状**：Jaeger UI中无Trace数据

**排查步骤**：

1. 检查Jaeger服务状态
```bash
curl http://localhost:16686/api/services
```

2. 检查追踪初始化
```bash
# 查看应用日志
grep "tracing" /var/log/coze-studio/app.log
```

3. 检查采样率配置
```bash
# 默认采样率为10%，调高以获取更多Trace
```

---

## 📚 最佳实践

### 1. 指标命名规范

- 使用小写字母和下划线
- 指标名包含业务领域前缀：`billing_*`
- 使用标准单位：`_seconds`（时间）、`_bytes`（字节）、`_total`（计数）

### 2. 日志记录规范

- 使用结构化日志（JSON格式）
- 包含必要的上下文信息（tenant_id、trace_id等）
- 使用适当的日志级别：Debug < Info < Warn < Error
- 避免记录敏感信息（API密钥、密码等）

### 3. 告警配置规范

- 设置合理的持续时间和阈值
- 包含详细的告警描述和Runbook链接
- 使用适当的告警级别
- 配置告警静音和抑制规则

### 4. 性能优化

- 使用Prometheus recording rules预计算复杂查询
- 配置合理的日志保留策略
- 优化Elasticsearch索引mapping
- 使用Filebeat批量发送日志

---

## 📖 相关文档

- [Prometheus最佳实践](https://prometheus.io/docs/practices/)
- [Grafana Dashboard指南](./GRAFANA_GUIDE.md)
- [ELK Stack文档](https://www.elastic.co/guide/)
- [Jaeger文档](https://www.jaegertracing.io/docs/)
- [企业级开发规范手册](./ZKER-企业级开发规范手册_v1.0.md)

---

**更新日志**：

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2025-01-01 | v1.0.0 | 初始版本，完整的监控和日志集成指南 | Claude AI |
