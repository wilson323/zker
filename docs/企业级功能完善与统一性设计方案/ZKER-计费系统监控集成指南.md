# ZKER 计费系统监控集成指南

**版本**: v1.0
**日期**: 2025-01-01
**作者**: 研发B (后端工程师)

---

## 📋 目录

1. [概述](#概述)
2. [监控架构](#监控架构)
3. [Prometheus Metrics](#prometheus-metrics)
4. [日志系统](#日志系统)
5. [分布式追踪](#分布式追踪)
6. [Grafana仪表板](#grafana仪表板)
7. [告警规则](#告警规则)
8. [健康检查](#健康检查)
9. [集成步骤](#集成步骤)
10. [运维指南](#运维指南)

---

## 概述

### 监控目标

为Token计量和预算管理系统提供企业级可观测性:

- **指标监控**: 实时监控业务指标和系统性能
- **日志追踪**: 结构化日志记录所有关键操作
- **分布式追踪**: 端到端的请求追踪
- **告警通知**: 及时发现和响应异常情况
- **健康检查**: 系统健康度实时监控

### 技术栈

| 组件 | 技术 | 版本 | 用途 |
|------|------|------|------|
| **指标采集** | Prometheus | v2.45.0 | 指标采集和存储 |
| **可视化** | Grafana | v10.0.0 | 监控仪表板 |
| **日志** | Zap + ELK | latest | 结构化日志 |
| **追踪** | OpenTelemetry + Jaeger | v1.50 | 分布式追踪 |
| **告警** | Alertmanager | latest | 告警路由和通知 |

---

## 监控架构

### 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                       ZKER API 服务                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Token计量    │  │ 预算告警      │  │ 汇总统计      │      │
│  │              │  │              │  │              │      │
│  │ Metrics      │  │ Metrics      │  │ Metrics      │      │
│  │ Logs         │  │ Logs         │  │ Logs         │      │
│  │ Traces       │  │ Traces       │  │ Traces       │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
           │                  │                  │
           └──────────────────┼──────────────────┘
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Prometheus Server                        │
│  - 每15秒抓取 /metrics 端点                                  │
│  - 每15秒评估告警规则                                        │
│  - 存储时序数据                                              │
└─────────────────────────────────────────────────────────────┘
                              │
           ┌──────────────────┼──────────────────┐
           ▼                  ▼                  ▼
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│  Grafana         │  │  Alertmanager    │  │  Jaeger          │
│  - 仪表板        │  │  - 告警路由      │  │  - 分布式追踪    │
│  - 可视化        │  │  - 通知发送      │  │  - 性能分析      │
└──────────────────┘  └──────────────────┘  └──────────────────┘
```

### 数据流

1. **应用层**: 业务代码记录Metrics/Logs/Traces
2. **采集层**: Prometheus定期抓取Metrics端点
3. **处理层**: Prometheus评估告警规则
4. **展示层**: Grafana可视化仪表板
5. **告警层**: Alertmanager发送告警通知

---

## Prometheus Metrics

### Metrics定义位置

```
backend/domain/billing/service/metrics.go
```

### Metrics类型

#### 1. Counter指标 (计数器)

用于记录只会增加的数值:

```go
// Token记录总数
billing_token_records_total{tenant_id, model_provider, model_name}

// Token总数
billing_tokens_total{tenant_id, token_type}  // token_type: input, output

// 成本总额
billing_cost_total{tenant_id, model_provider}

// 缓存请求数
billing_cached_requests_total{tenant_id, model_provider}

// 预算检查次数
billing_budget_checks_total{tenant_id, triggered}

// 预算告警次数
billing_budget_alerts_total{tenant_id, alert_type, alert_level}
```

#### 2. Gauge指标 (仪表)

用于记录可增可减的数值:

```go
// 单价
billing_unit_price{model_provider, model_name, token_type}

// 当前使用率
billing_budget_usage_percent{tenant_id, budget_type, alert_level}

// 预算金额
billing_budget_amount{tenant_id, budget_type}

// 已使用金额
billing_budget_used_amount{tenant_id, budget_type}

// 每日Token使用量
billing_daily_tokens{tenant_id, bot_id}

// 每日成本
billing_daily_cost{tenant_id, bot_id}

// 模型使用量
billing_model_usage{tenant_id, model_provider, model_name}

// 活跃租户数
billing_active_tenants

// 总收入
billing_total_revenue{period}  // period: daily, monthly, yearly
```

#### 3. Histogram指标 (直方图)

用于记录延迟分布:

```go
// Token记录耗时
billing_token_record_duration_seconds_bucket{operation, le}
billing_token_record_duration_seconds_sum{operation}
billing_token_record_duration_seconds_count{operation}

// 定价计算耗时
billing_pricing_calculation_duration_seconds_bucket{model_provider, le}
billing_pricing_calculation_duration_seconds_sum{model_provider}
billing_pricing_calculation_duration_seconds_count{model_provider}

// 预算检查耗时
billing_budget_check_duration_seconds_bucket{tenant_id, le}
billing_budget_check_duration_seconds_sum{tenant_id}
billing_budget_check_duration_seconds_count{tenant_id}

// 汇总更新耗时
billing_summary_update_duration_seconds_bucket{tenant_id, le}
billing_summary_update_duration_seconds_sum{tenant_id}
billing_summary_update_duration_seconds_count{tenant_id}

// Token日志数据库操作耗时
billing_token_log_db_operation_duration_seconds_bucket{operation, le}
billing_token_log_db_operation_duration_seconds_sum{operation}
billing_token_log_db_operation_duration_seconds_count{operation}

// 告警通知发送耗时
billing_alert_notification_duration_seconds_bucket{channel, le}
billing_alert_notification_duration_seconds_sum{channel}
billing_alert_notification_duration_seconds_count{channel}
```

### 使用示例

#### 在业务代码中记录Metrics

```go
import (
    "github.com/coze-dev/coze-studio/backend/domain/billing/service"
)

func (s *TokenMeteringService) RecordTokenUsage(
    ctx context.Context,
    req *RecordTokenUsageRequest,
) (*RecordTokenUsageResponse, error) {
    start := time.Now()

    // ... 业务逻辑 ...

    // 记录Token记录指标
    service.RecordTokenRecord(
        req.TenantID,
        req.ModelProvider,
        req.ModelName,
        req.InputTokens,
        req.OutputTokens,
        resp.TotalCost,
        time.Since(start).Seconds(),
        req.IsCached,
    )

    return resp, nil
}
```

---

## 日志系统

### 日志定义位置

```
backend/domain/billing/logging/billing_logging.go
```

### 日志结构

所有日志采用JSON格式,包含以下标准字段:

```json
{
  "timestamp": "2025-01-01T12:00:00Z",
  "level": "info",
  "logger": "billing",
  "message": "Token usage recorded",
  "tenant_id": "tenant123",
  "user_id": "user456",
  "bot_id": "bot789",
  "model_provider": "openai",
  "model_name": "gpt-4",
  "tokens_input": 100,
  "tokens_output": 50,
  "cost_total": 0.15,
  "trace_id": "abc123",
  "duration_ms": 45
}
```

### 日志级别

| 级别 | 用途 | 示例 |
|------|------|------|
| **DEBUG** | 详细调试信息 | 定价计算详情、缓存命中 |
| **INFO** | 正常业务操作 | Token记录成功、预算检查完成 |
| **WARN** | 警告事件 | 预算告警、慢操作 |
| **ERROR** | 错误事件 | 数据库操作失败、通知发送失败 |

### 使用示例

#### Token记录日志

```go
import billinglog "github.com/coze-dev/coze-studio/backend/domain/billing/logging"

func (s *TokenMeteringService) RecordTokenUsage(
    ctx context.Context,
    req *RecordTokenUsageRequest,
) (*RecordTokenUsageResponse, error) {
    start := time.Now()

    // ... 业务逻辑 ...

    // 记录日志
    billinglog.LogTokenRecord(ctx, req, resp, time.Since(start))

    return resp, nil
}
```

#### 预算告警日志

```go
func (s *BudgetAlertService) CheckBudget(ctx context.Context, tenantID string) error {
    // ... 检查逻辑 ...

    if alertSent {
        billinglog.LogBudgetAlert(ctx, alert)
    }

    return nil
}
```

#### 错误日志

```go
func (s *TokenMeteringService) RecordTokenUsage(...) error {
    if err != nil {
        billinglog.LogError(ctx, "RecordTokenUsage", err,
            zap.String("tenant_id", req.TenantID),
            zap.String("model", req.ModelName),
        )
        return err
    }
    return nil
}
```

---

## 分布式追踪

### 追踪定义位置

```
backend/domain/billing/service/tracing.go
```

### 追踪上下文

所有关键操作都创建Span:

```go
// Token记录追踪
ctx, span := service.StartRecordSpan(ctx, tenantID, modelProvider, modelName)
defer span.End()

// 预算检查追踪
ctx, span := service.StartBudgetCheckSpan(ctx, tenantID)
defer span.End()
```

### Span属性

每个Span包含标准属性:

```go
// Token记录Span
- tenant_id: string
- model.provider: string
- model.name: string
- tokens.input: int
- tokens.output: int
- cost.total: float64
- duration_ms: int64
```

### 使用示例

```go
import (
    "github.com/coze-dev/coze-studio/backend/domain/billing/service"
)

func (s *TokenMeteringService) RecordTokenUsage(
    ctx context.Context,
    req *RecordTokenUsageRequest,
) (*RecordTokenUsageResponse, error) {
    // 开始Span
    ctx, span := service.StartRecordSpan(ctx, req.TenantID, req.ModelProvider, req.ModelName)
    defer span.End()

    // ... 业务逻辑 ...

    // 设置Span标签
    service.SetRecordTags(span, req, resp)
    service.SetSuccessTag(span, true)

    return resp, nil
}
```

### 在Jaeger中查看追踪

1. 访问 Jaeger UI: `http://localhost:16686`
2. 选择服务: `zker-billing`
3. 输入Trace ID或Search查找
4. 查看详细的调用链路

---

## Grafana仪表板

### 仪表板定义位置

```
backend/infra/monitoring/grafana/dashboards/billing-dashboard.json
```

### 仪表板面板

#### 1. Token记录概览 (Stat)

- **指标**:
  - Token记录速率 (rec/s)
  - 预估小时成本 (¥/h)
- **阈值**:
  - 绿色: < 80
  - 黄色: 80-100
  - 红色: > 100

#### 2. 实时Token记录速率 (Graph)

- **指标**:
  - Token记录 (rec/s)
  - 成本速率 (¥/s)
  - 输入Token (tok/s)
  - 输出Token (tok/s)
- **用途**: 实时监控业务流量

#### 3. Token记录延迟分布 (Graph)

- **指标**:
  - P50延迟
  - P90延迟
  - P99延迟
- **阈值**: P99 > 2秒触发告警

#### 4. 预算使用率 TOP10 (Table)

- **指标**: `billing_budget_usage_percent`
- **列**:
  - 租户ID
  - 预算类型
  - 使用率 (%)
  - 告警级别
- **颜色**:
  - 绿色: < 80%
  - 橙色: 80-95%
  - 红色: ≥ 95%

#### 5. 告警统计 (Graph)

- **指标**:
  - 警告级别告警
  - 严重级别告警
  - 告警总数
- **用途**: 监控告警趋势

#### 6. 模型使用分布 (Pie Chart)

- **指标**: `billing_cost_total` by model_provider
- **用途**: 成本分布分析

#### 7. 模型成本占比 (Bar Gauge)

- **指标**: TOP5 模型成本
- **阈值**:
  - 绿色: < ¥3,000
  - 黄色: ¥3,000-7,000
  - 红色: > ¥7,000

#### 8. API性能 (Graph)

- **指标**:
  - QPS
  - 错误率
- **用途**: API性能监控

#### 9. API P99延迟 (Gauge)

- **指标**: P99延迟
- **阈值**:
  - 绿色: < 0.5秒
  - 黄色: 0.5-1秒
  - 红色: > 1秒

#### 10. 成本趋势 (Graph)

- **指标**:
  - 当前小时成本 (¥/h)
  - 预估日成本 (¥/d)
- **用途**: 成本预测和分析

### 导入仪表板

1. 登录Grafana: `http://localhost:3000`
2. 导航到 Dashboards → Import
3. 上传 `billing-dashboard.json`
4. 选择数据源: Prometheus

---

## 告警规则

### 告警规则定义位置

```
backend/infra/monitoring/prometheus/billing-alerts.yml
```

### 告警分类

#### 1. Token计量告警

**高记录失败率告警**
```yaml
- alert: BillingHighRecordFailureRate
  expr: rate(...{result="failure"}[5m]) / rate(...[5m]) > 0.05
  for: 5m
  severity: critical
```

**Token记录延迟过高**
```yaml
- alert: BillingTokenRecordSlow
  expr: histogram_quantile(0.99, ...) > 2
  for: 5m
  severity: warning
```

**定价计算错误率过高**
```yaml
- alert: BillingPricingHighErrorRate
  expr: rate(...errors...) / rate(...calculations...) > 0.1
  for: 5m
  severity: warning
```

#### 2. 预算告警

**预算使用率警告 (80%)**
```yaml
- alert: BillingBudgetWarning
  expr: billing_budget_usage_percent >= 80
  for: 5m
  severity: warning
```

**预算即将耗尽 (95%)**
```yaml
- alert: BillingBudgetCritical
  expr: billing_budget_usage_percent >= 95
  for: 1m
  severity: critical
```

**预算已超限**
```yaml
- alert: BillingBudgetExceeded
  expr: billing_budget_usage_percent > 100
  for: 1m
  severity: critical
```

#### 3. 成本告警

**日成本异常增长**
```yaml
- alert: BillingDailyCostSpike
  expr: rate(...[1h])*24 > (rate(...[24h])*24)*2
  for: 15m
  severity: warning
```

**高成本租户**
```yaml
- alert: BillingHighCostTenant
  expr: rate(...[1h])*24 > 10000
  for: 30m
  severity: warning
```

#### 4. 数据库告警

**Token日志数据库延迟**
```yaml
- alert: BillingTokenLogDBSlow
  expr: histogram_quantile(0.95, ...) > 1
  for: 5m
  severity: warning
```

**Token日志数据库宕机**
```yaml
- alert: BillingTokenLogDBDown
  expr: up{job="billing-token-log-db"} == 0
  for: 1m
  severity: critical
```

#### 5. 通知告警

**通知发送失败率过高**
```yaml
- alert: BillingNotificationHighFailureRate
  expr: rate(...{result="failure"}[5m]) / rate(...[5m]) > 0.2
  for: 5m
  severity: critical
```

### 告警级别

| 级别 | 说明 | 响应时间 |
|------|------|----------|
| **info** | 信息性告警,无需立即处理 | 工作时间 |
| **warning** | 警告,需要关注 | 1小时内 |
| **critical** | 严重,需要立即处理 | 15分钟内 |

### 配置告警通知

编辑 `alertmanager.yml`:

```yaml
receivers:
  - name: 'billing-alerts'
    email_configs:
      - to: 'billing-team@zker.com'
        from: 'alertmanager@zker.com'
        smarthost: 'smtp.zker.com:587'

    webhook_configs:
      - url: 'http://your-webhook-url/alert'
```

---

## 健康检查

### 健康检查定义位置

```
backend/domain/billing/health/billing_health.go
```

### 检查项

#### 1. 数据库健康检查

- 检查数据库连接
- 统计最近24小时记录数
- 测量查询延迟

#### 2. Token日志表健康检查

- 检查表是否存在
- 检查最新记录年龄
- 检查记录是否积压

#### 3. 预算配置表健康检查

- 检查表是否存在
- 统计活跃预算配置数

#### 4. Redis健康检查 (可选)

- 检查Redis连接
- 测量Redis延迟

### 健康状态

| 状态 | 说明 | 操作 |
|------|------|------|
| **healthy** | 所有检查正常 | 无需操作 |
| **degraded** | 部分检查降级,但不影响核心功能 | 关注 |
| **unhealthy** | 关键检查失败 | 立即处理 |

### 使用示例

```go
import (
    "github.com/coze-dev/coze-studio/backend/domain/billing/health"
)

// 创建健康检查器
healthChecker := health.NewHealthChecker(db, redisClient)

// 执行完整检查
status := healthChecker.Check(ctx)
fmt.Printf("Health Status: %s\n", status.Status)

// 快速活性检查
if err := healthChecker.CheckLiveness(ctx); err != nil {
    // 系统不健康,触发告警
}

// 完整就绪检查
if err := healthChecker.CheckReadiness(ctx); err != nil {
    // 系统未就绪,停止接收流量
}
```

### 暴露健康检查端点

```go
// GET /health
func GetHealth(c *app.RequestContext) {
    status := healthChecker.Check(context.Background())
    c.JSON(200, status)
}

// GET /health/live
func GetLiveness(c *app.RequestContext) {
    if err := healthChecker.CheckLiveness(context.Background()); err != nil {
        c.SetStatusCode(503)
        c.JSON(503, map[string]string{"status": "unhealthy"})
        return
    }
    c.JSON(200, map[string]string{"status": "healthy"})
}

// GET /health/ready
func GetReadiness(c *app.RequestContext) {
    if err := healthChecker.CheckReadiness(context.Background()); err != nil {
        c.SetStatusCode(503)
        c.JSON(503, map[string]string{"status": "not ready"})
        return
    }
    c.JSON(200, map[string]string{"status": "ready"})
}
```

---

## 集成步骤

### 步骤1: 初始化日志系统

```go
import billinglog "github.com/coze-dev/coze-studio/backend/domain/billing/logging"

// 在应用启动时初始化
func main() {
    // 初始化全局日志
    logger := zap.Must(zap.NewProduction())
    billinglog.Init(logger)
}
```

### 步骤2: 初始化追踪系统

```go
import (
    "github.com/coze-dev/coze-studio/backend/infra/tracing"
)

// 在应用启动时初始化
func main() {
    if err := tracing.Init(
        "zker-billing",
        "http://localhost:14268/api/traces",
    ); err != nil {
        panic(err)
    }
}
```

### 步骤3: 暴露Metrics端点

```go
import (
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/cloudwego/hertz/pkg/app/server"
)

h := server.Default()

// 暴露/metrics端点
h.GET("/metrics", func(ctx context.Context, c *app.RequestContext) {
    promhttp.Handler().ServeHTTP(c.Response.Writer, c.Request)
})
```

### 步骤4: 添加Prometheus中间件

```go
import (
    "github.com/coze-dev/coze-studio/backend/api/middleware"
)

h.Use(middleware.PrometheusMiddleware())
```

### 步骤5: 在业务代码中集成

```go
func (s *TokenMeteringService) RecordTokenUsage(
    ctx context.Context,
    req *RecordTokenUsageRequest,
) (*RecordTokenUsageResponse, error) {
    start := time.Now()

    // 1. 开始追踪
    ctx, span := service.StartRecordSpan(ctx, req.TenantID, req.ModelProvider, req.ModelName)
    defer span.End()

    // 2. 业务逻辑
    resp, err := s.recordTokenUsage(ctx, req)
    if err != nil {
        // 记录错误日志
        billinglog.LogError(ctx, "RecordTokenUsage", err)
        // 设置Span错误状态
        service.SetErrorTag(ctx, err)
        return nil, err
    }

    // 3. 记录Metrics
    service.RecordTokenRecord(
        req.TenantID,
        req.ModelProvider,
        req.ModelName,
        req.InputTokens,
        req.OutputTokens,
        resp.TotalCost,
        time.Since(start).Seconds(),
        req.IsCached,
    )

    // 4. 记录日志
    billinglog.LogTokenRecord(ctx, req, resp, time.Since(start))

    // 5. 设置Span标签
    service.SetRecordTags(span, req, resp)
    service.SetSuccessTag(span, true)

    return resp, nil
}
```

### 步骤6: 配置Prometheus

1. 编辑 `prometheus.yml`:

```yaml
rule_files:
  - 'alerts.yml'
  - 'billing-alerts.yml'  # 添加计费告警规则
```

2. 重启Prometheus:

```bash
docker restart prometheus
```

### 步骤7: 导入Grafana仪表板

1. 登录Grafana
2. 导航到 Dashboards → Import
3. 上传 `billing-dashboard.json`
4. 选择数据源

### 步骤8: 配置告警通知

1. 编辑 `alertmanager.yml`
2. 配置邮件/Webhook/钉钉/企业微信
3. 重启Alertmanager

---

## 运维指南

### 日常监控

#### 每日检查项

- [ ] 查看Grafana仪表板,检查关键指标
- [ ] 检查告警通知,确认无严重告警
- [ ] 查看日志,确认无ERROR级别日志
- [ ] 检查健康检查端点,确认系统健康

#### 每周检查项

- [ ] 分析Token记录趋势,识别异常
- [ ] 检查成本趋势,预测下周成本
- [ ] 审查告警规则,调整阈值
- [ ] 检查Prometheus存储使用情况

#### 每月检查项

- [ ] 生成月度成本报告
- [ ] 分析TOP10高成本租户
- [ ] 审查监控配置,优化查询性能
- [ ] 备份Prometheus数据

### 故障排查

#### 问题1: Token记录延迟过高

**症状**: P99延迟 > 2秒

**排查步骤**:
1. 检查数据库连接池是否耗尽
2. 检查数据库慢查询日志
3. 检查磁盘I/O是否正常
4. 检查网络延迟

**解决方案**:
- 增加数据库连接池大小
- 优化慢查询
- 使用批量插入

#### 问题2: 预算告警未发送

**症状**: 预算超限但未收到告警

**排查步骤**:
1. 检查告警规则是否启用
2. 检查Alertmanager是否正常
3. 检查通知渠道配置
4. 查看Alertmanager日志

**解决方案**:
- 重新加载告警规则
- 修复通知渠道配置
- 重启Alertmanager

#### 问题3: 成本异常增长

**症状**: 日成本突然翻倍

**排查步骤**:
1. 查看成本趋势图,定位开始时间
2. 查看Token记录量,是否流量暴增
3. 查看模型使用分布,是否切换到高价模型
4. 查看是否有异常租户

**解决方案**:
- 联系租户确认
- 调整定价策略
- 启用预算硬限制

### 性能优化

#### Prometheus查询优化

**问题**: 查询慢,仪表板加载时间长

**优化方案**:
1. 使用 recording rules 预计算常用指标
2. 减少查询时间范围
3. 优化PromQL查询
4. 增加Prometheus内存

示例recording rule:
```yaml
groups:
  - name: billing_recording
    interval: 1m
    rules:
      - record: job:billing_cost_total:rate1h
        expr: sum by (tenant_id) (rate(billing_cost_total[1h]))
```

#### 日志轮转

**问题**: 日志文件过大,占用磁盘空间

**优化方案**:
```go
// 使用lumberjack进行日志轮转
lumberjack.Write(
    "/var/log/zker/billing.log",
    lumberjack.MaxSize(100),    // 每个文件最大100MB
    lumberjack.MaxBackups(10),   // 保留10个备份
    lumberjack.MaxAge(30),       // 保留30天
    lumberjack.Compress(true),   // 压缩旧文件
)
```

#### Metrics基数控制

**问题**: Label cardinality过高,内存占用大

**优化方案**:
1. 移除高基数Label (如bot_id, conversation_id)
2. 使用user_id代替具体user信息
3. 限制Label值的数量
4. 定期清理旧数据

### 备份与恢复

#### Prometheus数据备份

```bash
# 备份Prometheus数据
docker exec prometheus tar czf /tmp/prometheus-data.tar.gz /prometheus

# 恢复Prometheus数据
docker exec prometheus tar xzf /tmp/prometheus-data.tar.gz -C /
```

#### Grafana仪表板备份

```bash
# 导出所有仪表板
curl -u admin:admin http://localhost:3000/api/search?query=dashboards > dashboards.json

# 导出单个仪表板
curl -u admin:admin http://localhost:3000/api/dashboards/uid/billing-dashboard > billing-dashboard-backup.json
```

---

## 常见问题

### Q1: 如何添加自定义Metrics?

```go
// 1. 在metrics.go中定义
var MyCustomMetric = promauto.NewCounterVec(
    prometheus.CounterOpts{
        Name: "billing_my_custom_metric",
        Help: "My custom metric",
    },
    []string{"label1", "label2"},
)

// 2. 在业务代码中记录
MyCustomMetric.WithLabelValues("value1", "value2").Inc()
```

### Q2: 如何查看实时日志?

```bash
# 使用journalctl查看systemd服务日志
journalctl -u zker-api -f

# 查看日志文件
tail -f /var/log/zker/billing.log

# 使用jq过滤JSON日志
tail -f /var/log/zker/billing.log | jq '.level=="error"'
```

### Q3: 如何测试告警规则?

```bash
# 使用promtool检查告警规则语法
promtool check rules billing-alerts.yml

# 查看当前激活的告警
curl http://localhost:9090/api/v1/alerts

# 手动触发告警测试
curl -X POST http://localhost:9090/api/v1/alerts -d '...'
```

### Q4: Grafana仪表板显示"No Data"

**可能原因**:
1. Prometheus数据源未配置
2. 时间范围选择不当
3. Metrics名称或Label不匹配
4. Prometheus未采集到数据

**排查步骤**:
1. 检查Grafana数据源配置
2. 调整时间范围
3. 在Prometheus UI中查询Metrics
4. 检查 `/metrics` 端点是否正常

---

## 附录

### A. 参考文档

- [Prometheus文档](https://prometheus.io/docs/)
- [Grafana文档](https://grafana.com/docs/)
- [OpenTelemetry文档](https://opentelemetry.io/docs/)
- [Zap日志库文档](https://github.com/uber-go/zap)

### B. 相关文件清单

```
backend/domain/billing/
├── service/
│   ├── metrics.go              # Metrics定义
│   └── tracing.go              # 分布式追踪
├── logging/
│   └── billing_logging.go      # 日志工具
└── health/
    └── billing_health.go       # 健康检查

backend/infra/monitoring/
├── prometheus/
│   ├── prometheus.yml          # Prometheus配置
│   ├── billing-alerts.yml      # 计费告警规则
│   └── alerts.yml              # 通用告警规则
└── grafana/
    └── dashboards/
        └── billing-dashboard.json  # 计费仪表板
```

### C. 联系方式

- **监控负责人**: 研发B (后端工程师)
- **邮箱**: monitoring@zker.com
- **钉钉群**: ZKER运维团队
