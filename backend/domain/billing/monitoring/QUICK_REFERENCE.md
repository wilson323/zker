# 计费系统监控快速参考

**版本**: v1.0
**最后更新**: 2025-01-01

---

## 🎯 核心Metrics速查

### Token计量指标

```go
// 记录Token使用
service.RecordTokenRecord(
    tenantID,          // 租户ID
    modelProvider,     // 模型提供商
    modelName,         // 模型名称
    inputTokens,       // 输入Token数
    outputTokens,      // 输出Token数
    totalCost,         // 总成本
    duration.Seconds(),// 耗时(秒)
    isCached,          // 是否缓存
)

// 记录定价计算
service.RecordPricingCalculation(
    provider,          // 提供商
    model,             // 模型
    "success",         // 结果: success/failure
    duration.Seconds(),// 耗时
    unitPrice,         // 单价
    tokenType,         // Token类型
)

// 记录定价错误
service.RecordPricingError(
    provider,          // 提供商
    model,             // 模型
    "model_not_found", // 错误类型
)
```

### 预算告警指标

```go
// 记录预算检查
service.RecordBudgetCheck(
    tenantID,          // 租户ID
    true,              // 是否触发告警
    duration.Seconds(),// 耗时
)

// 记录预算告警
service.RecordBudgetAlert(
    tenantID,          // 租户ID
    "threshold1",      // 告警类型
    "warning",         // 告警级别
)

// 更新预算使用率
service.UpdateBudgetUsage(
    tenantID,          // 租户ID
    "monthly",         // 预算类型
    "warning",         // 告警级别
    85.5,              // 使用率(%)
    10000.0,           // 预算金额
    8550.0,            // 已使用金额
)
```

### 数据库操作指标

```go
// 记录Token日志数据库操作
service.RecordTokenLogDBOperation(
    "batch_create",    // 操作: create/batch_create/get_by_date_range
    "success",         // 结果: success/failure
    duration.Seconds(),// 耗时
)
```

---

## 📝 日志记录速查

### Token计量日志

```go
import billinglog "github.com/coze-dev/coze-studio/backend/domain/billing/logging"

// 记录Token记录
billinglog.LogTokenRecord(ctx, req, resp, duration)

// 记录批量Token记录
billinglog.LogBatchTokenRecord(ctx, tenantID, recordCount, totalCost, duration)

// 记录定价计算
billinglog.LogPricingCalculation(ctx, provider, model, inputTokens, outputTokens,
    inputCost, outputCost, totalCost, duration, err)

// 记录使用统计查询
billinglog.LogUsageStatsRetrieval(ctx, tenantID, startDate, endDate, recordCount, duration)
```

### 预算告警日志

```go
// 记录预算检查
billinglog.LogBudgetCheck(ctx, tenantID, budgetType, budgetAmount,
    usedAmount, usagePercent, triggered, duration)

// 记录预算告警
billinglog.LogBudgetAlert(ctx, alert)

// 记录批量预算检查
billinglog.LogBatchBudgetCheck(ctx, tenantCount, alertCount, duration)
```

### 汇总更新日志

```go
// 记录汇总更新
billinglog.LogSummaryUpdate(ctx, tenantID, summaryDate, dailyTokens,
    dailyCost, duration, err)
```

### 通知发送日志

```go
// 记录告警通知
billinglog.LogAlertNotification(ctx, tenantID, channel, alertID, duration, err)
```

### 通用日志

```go
// 记录错误
billinglog.LogError(ctx, "operation_name", err)

// 记录带字段的错误
billinglog.LogErrorWithFields(ctx, "operation_name", err, map[string]interface{}{
    "tenant_id": tenantID,
    "model": model,
})

// 记录慢操作
billinglog.LogSlowOperation(ctx, "operation_name", duration, threshold)

// 记录审计事件
billinglog.LogAuditEvent(ctx, eventType, actor, action, target, details)
```

---

## 🔍 分布式追踪速查

### 开始Span

```go
import "github.com/coze-dev/coze-studio/backend/domain/billing/service"

// Token记录Span
ctx, span := service.StartRecordSpan(ctx, tenantID, modelProvider, modelName)
defer span.End()

// 批量Token记录Span
ctx, span := service.StartBatchRecordSpan(ctx, tenantID, recordCount)
defer span.End()

// 定价计算Span
ctx, span := service.StartPricingCalculationSpan(ctx, provider, model,
    inputTokens, outputTokens)
defer span.End()

// 预算检查Span
ctx, span := service.StartBudgetCheckSpan(ctx, tenantID)
defer span.End()

// 批量预算检查Span
ctx, span := service.StartBatchBudgetCheckSpan(ctx, tenantCount)
defer span.End()

// 告警发送Span
ctx, span := service.StartAlertSendSpan(ctx, tenantID, channel)
defer span.End()

// 汇总更新Span
ctx, span := service.StartSummaryUpdateSpan(ctx, tenantID, summaryDate)
defer span.End()

// 使用统计查询Span
ctx, span := service.StartUsageStatsRetrievalSpan(ctx, tenantID, startDate, endDate)
defer span.End()
```

### 设置Span属性

```go
// 设置Token记录标签
service.SetRecordTags(span, req, resp)

// 设置定价标签
service.SetPricingTags(span, inputCost, outputCost, totalCost, unitPrice)

// 设置预算检查标签
service.SetBudgetCheckTags(span, budgetAmount, usedAmount, usagePercent, triggered)

// 设置告警标签
service.SetAlertTags(span, alertType, alertLevel, usagePercent)

// 设置成功状态
service.SetSuccessTag(span, true)

// 添加事件
service.AddEvent(ctx, "event_name", attrs...)
```

### 获取Trace信息

```go
// 获取Trace ID
traceID := service.GetTraceID(ctx)

// 获取Span ID
spanID := service.GetSpanID(ctx)
```

---

## 🚨 告警规则速查

### 关键告警

| 告警名称 | 级别 | 触发条件 | 响应时间 |
|---------|------|----------|----------|
| BillingHighRecordFailureRate | critical | 失败率 > 5% | 15分钟 |
| BillingTokenRecordSlow | warning | P99 > 2秒 | 1小时 |
| BillingBudgetCritical | critical | 使用率 ≥ 95% | 15分钟 |
| BillingBudgetExceeded | critical | 使用率 > 100% | 立即 |
| BillingDailyCostSpike | warning | 日成本翻倍 | 1小时 |
| BillingNotificationHighFailureRate | critical | 失败率 > 20% | 15分钟 |
| BillingTokenLogDBDown | critical | 数据库宕机 | 立即 |

### 告警查询

```bash
# 查看当前激活的告警
curl http://localhost:9090/api/v1/alerts | jq '.data.alerts[] | select(.labels.component=="billing")'

# 查看告警规则评估结果
curl http://localhost:9090/api/v1/rules | jq '.data.groups[] | select(.name=="billing_budget_alerts")'
```

---

## 🏥 健康检查速查

### 健康检查端点

```bash
# 完整健康检查
curl http://localhost:8080/health

# 活性检查 (Kubernetes liveness)
curl http://localhost:8080/health/live

# 就绪检查 (Kubernetes readiness)
curl http://localhost:8080/health/ready
```

### 健康状态

```json
{
  "status": "healthy",
  "timestamp": "2025-01-01T12:00:00Z",
  "version": "1.0.0",
  "checks": {
    "database": {
      "status": "healthy",
      "duration": "5ms",
      "details": {
        "recent_records": 1000,
        "latency_ms": 5
      }
    },
    "token_log_table": {
      "status": "healthy",
      "duration": "3ms",
      "details": {
        "latest_record_age_ms": 60000,
        "latest_record_id": "12345"
      }
    },
    "budget_config_table": {
      "status": "healthy",
      "duration": "2ms",
      "details": {
        "active_budgets": 50
      }
    }
  },
  "uptime": "10ms"
}
```

---

## 📊 Grafana查询速查

### PromQL查询示例

```promql
# Token记录速率 (每秒)
sum(rate(billing_token_records_total[5m]))

# 成本速率 (每小时)
sum(rate(billing_cost_total[1h])) * 3600

# P99延迟
histogram_quantile(0.99, rate(billing_token_record_duration_seconds_bucket[5m]))

# 预算使用率 TOP10
topk(10, billing_budget_usage_percent)

# 模型成本分布
sum by (model_provider) (billing_cost_total)

# 错误率
rate(billing_token_log_db_operations_total{result="failure"}[5m]) /
rate(billing_token_log_db_operations_total[5m])

# 日成本预测
rate(billing_cost_total[1h]) * 24

# 活跃租户数
billing_active_tenants

# 缓存命中率
rate(cache_hit_total{cache_type="pricing"}[5m]) /
(rate(cache_hit_total{cache_type="pricing"}[5m]) +
 rate(cache_miss_total{cache_type="pricing"}[5m]))
```

---

## 🛠️ 常用运维命令

### Prometheus

```bash
# 检查Prometheus状态
curl http://localhost:9090/-/healthy

# 重新加载配置
curl -X POST http://localhost:9090/-/reload

# 查询Metrics
curl 'http://localhost:9090/api/v1/query?query=billing_token_records_total'

# 查看配置
curl http://localhost:9090/api/v1/status/config

# 检查告警规则
promtool check rules billing-alerts.yml
```

### Grafana

```bash
# 导出仪表板
curl -u admin:admin http://localhost:3000/api/dashboards/uid/billing-dashboard

# 导入仪表板
curl -u admin:admin -X POST http://localhost:3000/api/dashboards/db \
  -H "Content-Type: application/json" \
  -d @billing-dashboard.json

# 查看数据源
curl -u admin:admin http://localhost:3000/api/datasources
```

### 日志查询

```bash
# 查看ERROR日志
tail -f /var/log/zker/billing.log | jq '.level=="error"'

# 按租户ID查询
tail -f /var/log/zker/billing.log | jq 'select(.tenant_id=="tenant123")'

# 按时间范围查询
grep "2025-01-01T12:" /var/log/zker/billing.log | jq

# 统计ERROR数量
jq -r 'select(.level=="error") | .message' /var/log/zker/billing.log | wc -l
```

---

## 🚀 快速诊断

### 问题: Token记录延迟高

```bash
# 1. 查看P99延迟
curl 'http://localhost:9090/api/v1/query?query=histogram_quantile(0.99,rate(billing_token_record_duration_seconds_bucket[5m]))' | jq

# 2. 查看数据库延迟
curl 'http://localhost:9090/api/v1/query?query=histogram_quantile(0.95,rate(billing_token_log_db_operation_duration_seconds_bucket[5m]))' | jq

# 3. 查看慢查询日志
tail -f /var/log/zker/billing.log | jq 'select(.duration_ms>1000)'

# 4. 检查数据库连接
curl http://localhost:8080/health | jq '.checks.database'
```

### 问题: 预算告警未发送

```bash
# 1. 查看告警规则状态
curl http://localhost:9090/api/v1/rules | jq '.data.groups[] | select(.name=="billing_budget_alerts")'

# 2. 查看当前告警
curl http://localhost:9090/api/v1/alerts | jq '.data.alerts[] | select(.labels.alert_type=="budget")'

# 3. 查看告警通知日志
tail -f /var/log/zker/billing.log | jq 'select(.message=="Alert notification sent")'

# 4. 检查Alertmanager
curl http://localhost:9093/api/v1/status
```

### 问题: 成本异常增长

```bash
# 1. 查看成本趋势
curl 'http://localhost:9090/api/v1/query?query=rate(billing_cost_total[1h])*24' | jq

# 2. 查看Token记录量
curl 'http://localhost:9090/api/v1/query?query=sum(rate(billing_token_records_total[1h]))' | jq

# 3. 查看模型使用分布
curl 'http://localhost:9090/api/v1/query?query=sum by(model_provider,model_name)(billing_cost_total)' | jq

# 4. 查看TOP租户
curl 'http://localhost:9090/api/v1/query?query=topk(10,sum by(tenant_id)(billing_cost_total))' | jq
```

---

## 📞 联系方式

- **文档**: [ZKER-计费系统监控集成指南](./ZKER-计费系统监控集成指南.md)
- **负责人**: 研发B (后端工程师)
- **邮箱**: monitoring@zker.com
- **钉钉**: ZKER运维团队
