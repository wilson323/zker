# ZKER Grafana监控运维手册

**版本**: v1.0.0
**最后更新**: 2025-01-03
**作者**: 研发B (后端工程师)
**适用阶段**: 生产运维、故障排查、性能分析

---

## 目录

- [1. 概述](#1-概述)
- [2. 架构设计](#2-架构设计)
- [3. 快速开始](#3-快速开始)
- [4. Dashboard使用](#4-dashboard使用)
- [5. 告警配置](#5-告警配置)
- [6. 故障排查](#6-故障排查)
- [7. 性能分析](#7-性能分析)
- [8. 运维操作](#8-运维操作)
- [9. 安全加固](#9-安全加固)
- [10. 最佳实践](#10-最佳实践)

---

## 1. 概述

### 1.1 监控体系

ZKER监控体系基于Prometheus + Grafana构建，提供完整的可观测性：

```
监控架构:

┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│   ZKER API  │─────▶│ Prometheus  │─────▶│  Grafana    │
│  :8888      │      │   :9090     │      │   :3000     │
└─────────────┘      └─────────────┘      └─────────────┘
       │                    │                     │
       │                    │                     │
       ▼                    ▼                     ▼
  采集指标             存储指标              可视化告警
└──────────────────────────────────────────────────┘
                      │
                      ▼
              ┌─────────────┐
              │AlertManager │
              │   :9093     │
              └─────────────┘
                      │
                      ▼
                  发送告警
```

### 1.2 监控覆盖范围

| 类别 | 覆盖内容 | Dashboard |
|------|---------|-----------|
| **API性能** | 请求率、延迟、错误率 | API Performance |
| **数据库** | 查询延迟、连接池、慢查询 | Database Monitoring |
| **缓存** | 命中率、操作延迟、大小 | Cache Monitoring |
| **租户** | 配额使用、订阅状态 | Tenant Monitoring |
| **业务** | Bot执行、消息、工作流 | Business Metrics |
| **系统** | CPU、内存、磁盘、网络 | System Monitoring |
| **日志** | 应用日志、审计日志 | Loki Logs |

---

## 2. 架构设计

### 2.1 数据流

```
应用指标采集流程:

ZKER API (Prometheus Middleware)
  │
  ├─ HTTP请求指标 (请求率、延迟、错误率)
  ├─ 数据库指标 (查询延迟、连接数)
  ├─ 缓存指标 (命中率、操作延迟)
  ├─ 业务指标 (Bot执行、消息发送)
  └─ 系统指标 (Goroutine、内存、GC)
  │
  ▼
/metrics端点 (Prometheus格式)
  │
  ▼
Prometheus (每15秒抓取一次)
  │
  ├─ 时序数据库存储 (30天保留)
  ├─ PromQL查询
  └─ 规则评估 (告警触发)
  │
  ▼
Grafana Dashboard (可视化)
  │
  ├─ 实时监控
  ├─ 历史趋势
  └─ 告警通知
```

### 2.2 指标分类

#### Counter (计数器)
- 只增不减
- 用于请求总数、错误总数等
- 示例: `http_requests_total`, `bot_invocation_total`

#### Gauge (仪表)
- 可增可减
- 用于当前值，如连接数、内存使用
- 示例: `db_connections_active`, `active_users_total`

#### Histogram (直方图)
- 分布情况
- 用于延迟、请求大小等
- 示例: `http_request_duration_seconds`, `db_query_duration_seconds`

#### Summary (摘要)
- 类似Histogram，但客户端计算
- 用于不需要聚合的分布
- 示例: (推荐使用Histogram)

---

## 3. 快速开始

### 3.1 启动监控栈

```bash
# 1. 启动所有监控服务
cd D:\code\coze-studio
docker-compose -f docker-compose.monitoring.yml up -d

# 2. 验证服务状态
docker-compose -f docker-compose.monitoring.yml ps

# 3. 查看日志
docker-compose -f docker-compose.monitoring.yml logs -f
```

**服务列表**:
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/admin)
- AlertManager: http://localhost:9093
- Node Exporter: http://localhost:9100/metrics

### 3.2 访问Grafana

1. 打开浏览器: http://localhost:3000
2. 登录 (默认: admin/admin)
3. 首次登录会提示修改密码
4. 导入Dashboard

**导入Dashboard**:
1. 点击 "+" → "Import"
2. 选择"ZKER Performance Overview"
3. 或上传JSON文件: `deploy/monitoring/grafana/dashboards/zker-performance-dashboard.json`

### 3.3 验证指标采集

```bash
# 检查Prometheus目标
curl http://localhost:9090/api/v1/targets | jq

# 查看当前指标
curl http://localhost:9090/api/v1/label/__name__/values | jq

# 查询特定指标
curl 'http://localhost:9090/api/v1/query?query=up' | jq
```

---

## 4. Dashboard使用

### 4.1 核心Dashboard

#### 1. ZKER Performance Overview

**位置**: http://localhost:3000/d/zker-performance

**面板**:
- API Request Rate (请求率)
- API Latency P95 (延迟P95)
- API Error Rate (错误率)
- Cache Hit Rate (缓存命中率)
- Database Query Duration (数据库查询延迟)
- Active Users (活跃用户)

**使用场景**:
- 日常监控
- 性能趋势分析
- 容量规划

#### 2. ZKER API Performance

**位置**: http://localhost:3000/d/zker-api

**面板**:
- 请求率 (按端点)
- 延迟分布 (P50, P95, P99)
- 错误率 (按状态码)
- 请求/响应大小

**使用场景**:
- API性能分析
- 慢端点定位
- 错误诊断

#### 3. ZKER Database Monitoring

**位置**: http://localhost:3000/d/zker-database

**面板**:
- 连接池使用率
- 查询延迟分布
- 慢查询Top10
- 事务成功率

**使用场景**:
- 数据库性能优化
- 慢查询排查
- 连接池调优

#### 4. ZKER Tenant Monitoring

**位置**: http://localhost:3000/d/zker-tenant

**面板**:
- 活跃租户数
- 配额使用率 (Top10)
- 订阅状态分布
- 即将过期订阅

**使用场景**:
- 租户容量管理
- 配额预警
- 订阅管理

### 4.2 查询技巧

#### PromQL基础

**查询当前值**:
```promql
# API请求率
rate(http_requests_total[5m])

# P95延迟
histogram_quantile(0.95, http_request_duration_seconds_bucket)

# 错误率
sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))
```

**聚合查询**:
```promql
# 按端点聚合
sum(rate(http_requests_total[5m])) by (endpoint)

# 按租户聚合
sum(rate(http_requests_total[5m])) by (tenant_id)

# Top10慢查询
topk(10, rate(db_query_duration_seconds_sum[5m]) / rate(db_query_duration_seconds_count[5m]))
```

**时间范围查询**:
```promql
# 比较昨天和今天
rate(http_requests_total[5m]) offset 1d
rate(http_requests_total[5m])

# 增长率
(rate(http_requests_total[5m]) - rate(http_requests_total[5m] offset 1h)) / rate(http_requests_total[5m] offset 1h) * 100
```

### 4.3 变量使用

在Dashboard中使用变量实现动态查询：

**定义变量**:
```yaml
# Dashboard Settings → Variables → New
Name: tenant_id
Type: Query
Query: label_values(coze_http_requests_total, tenant_id)
```

**使用变量**:
```promql
sum(rate(coze_http_requests_total{tenant_id="$tenant_id"}[5m])) by (endpoint)
```

---

## 5. 告警配置

### 5.1 告警规则

**告警规则文件**: `deploy/monitoring/prometheus/rules/zker-alerts.yml`

**告警分类**:

| 类别 | 严重级别 | 通知渠道 |
|------|---------|---------|
| **API错误率过高** | Critical | Email + Slack + PagerDuty |
| **API延迟过高** | Warning | Email + Slack |
| **数据库慢查询** | Warning | Email + Slack |
| **数据库连接池满** | Critical | Email + Slack + PagerDuty |
| **缓存命中率低** | Info | Email |
| **租户配额超限** | Warning | Email |
| **订阅即将过期** | Warning | Email |

### 5.2 告警规则示例

```yaml
# API错误率告警
- alert: HighAPIErrorRate
  expr: |
    sum(rate(coze_http_requests_total{status=~"5.."}[5m]))
    /
    sum(rate(coze_http_requests_total[5m])) > 0.05
  for: 5m
  labels:
    severity: critical
    category: api
  annotations:
    summary: "API错误率过高"
    description: "过去5分钟错误率为 {{ $value | humanizePercentage }}"
    runbook_url: "https://docs.coze-studio.com/runbooks/api-errors"
```

**参数说明**:
- `expr`: PromQL表达式
- `for`: 持续时间 (避免瞬时峰值)
- `severity`: 严重级别
- `annotations`: 告警描述

### 5.3 沉默告警

**临时维护**:
```bash
# 通过API创建沉默规则
curl -X POST http://localhost:9093/api/v1/silences \
  -H 'Content-Type: application/json' \
  -d '{
    "matchers": [
      {"name": "alertname", "value": "HighAPIErrorRate"},
      {"name": "severity", "value": "critical"}
    ],
    "startsAt": "2025-01-03T10:00:00Z",
    "endsAt": "2025-01-03T12:00:00Z",
    "comment": "数据库维护中"
  }'
```

**通过Web UI**:
1. 访问 AlertManager: http://localhost:9093
2. 点击 "Silence"
3. 选择匹配条件
4. 设置持续时间
5. 添加备注

### 5.4 告警通知

#### Email通知

配置`alertmanager.yml`:
```yaml
global:
  smtp_smarthost: 'smtp.gmail.com:587'
  smtp_from: 'alerts@coze-studio.com'
  smtp_auth_username: 'alerts@coze-studio.com'
  smtp_auth_password: '${SMTP_PASSWORD}'

receivers:
  - name: 'critical'
    email_configs:
      - to: 'oncall@coze-studio.com'
        headers:
          Subject: '[CRITICAL] [ZKER] {{ .GroupLabels.alertname }}'
```

#### Slack通知

```yaml
receivers:
  - name: 'critical'
    slack_configs:
      - channel: '#alerts-critical'
        title: '[CRITICAL] {{ .GroupLabels.alertname }}'
        text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'
        color: 'danger'
        send_resolved: true
```

#### PagerDuty通知

```yaml
receivers:
  - name: 'critical'
    pagerduty_configs:
      - service_key: '${PAGERDUTY_SERVICE_KEY}'
        description: '{{ .GroupLabels.alertname }}'
```

---

## 6. 故障排查

### 6.1 常见问题诊断

#### 问题1: API错误率突增

**排查步骤**:
1. 查看Grafana "API Error Rate"面板
2. 确认错误分布 (按端点、状态码)
3. 查看应用日志
4. 检查依赖服务状态

**PromQL查询**:
```promql
# 错误率按端点
sum(rate(coze_http_requests_total{status=~"5.."}[5m])) by (endpoint)

# 错误码分布
sum(rate(error_total[5m])) by (error_code)
```

#### 问题2: API延迟过高

**排查步骤**:
1. 查看P95延迟趋势
2. 定位慢端点
3. 检查数据库慢查询
4. 检查缓存命中率

**PromQL查询**:
```promql
# P95延迟按端点
histogram_quantile(0.95, sum(rate(coze_http_request_duration_seconds_bucket[5m])) by (le, endpoint))

# 慢查询Top10
topk(10, rate(db_query_duration_seconds_sum[5m]) / rate(db_query_duration_seconds_count[5m]))
```

#### 问题3: 数据库连接池耗尽

**排查步骤**:
1. 查看连接池使用率
2. 检查慢查询
3. 检查连接泄露
4. 分析连接等待时间

**PromQL查询**:
```promql
# 连接池使用率
db_connections_active / db_connections_max

# 慢查询
rate(db_query_duration_seconds_bucket{le="0.5"}[5m])
```

#### 问题4: 缓存命中率低

**排查步骤**:
1. 查看缓存命中率
2. 分析缓存未命中原因
3. 检查缓存配置
4. 优化缓存策略

**PromQL查询**:
```promql
# 缓存命中率
sum(rate(cache_hit_total[5m])) / (sum(rate(cache_hit_total[5m])) + sum(rate(cache_miss_total[5m])))

# 按缓存类型
sum(rate(cache_hit_total[5m])) by (cache_type) / (sum(rate(cache_hit_total[5m])) by (cache_type) + sum(rate(cache_miss_total[5m])) by (cache_type))
```

### 6.2 日志查询

使用Loki查询日志：

**查询所有错误日志**:
```logncli
{job="coze-studio"} |= "error"
```

**查询特定租户的日志**:
```logncli
{tenant_id="tenant-123"}
```

**查询慢查询日志**:
```logncli
{job="coze-studio"} |= "slow query" > 1000ms
```

---

## 7. 性能分析

### 7.1 容量规划

**计算当前容量**:
```promql
# 当前QPS
sum(rate(http_requests_total[5m]))

# P95延迟
histogram_quantile(0.95, http_request_duration_seconds_bucket)

# 系统资源使用率
avg(rate(process_cpu_seconds_total{mode!="idle"}[5m])) by (instance)
process_resident_memory_bytes / node_memory_MemTotal_bytes
```

**预测未来容量**:
```promql
# QPS增长率
rate(http_requests_total[7d])

# 预测30天后的QPS
predict_linear(http_requests_total[7d], 30*24*3600)
```

### 7.2 性能基线对比

**对比上周同期**:
```promql
# 今天
rate(http_requests_total[1h])

# 上周同期
rate(http_requests_total[1h] offset 1w)

# 增长率
(rate(http_requests_total[1h]) - rate(http_requests_total[1h] offset 1w)) / rate(http_requests_total[1h] offset 1w) * 100
```

### 7.3 性能优化验证

**优化前后对比**:
```promql
# 优化前 (标记时间点)
http_request_duration_seconds_bucket{deploy_env="prod"} < timestamp_before_optimization

# 优化后
http_request_duration_seconds_bucket{deploy_env="prod"} > timestamp_after_optimization
```

---

## 8. 运维操作

### 8.1 数据备份

**备份Prometheus数据**:
```bash
# 1. 停止Prometheus
docker-compose -f docker-compose.monitoring.yml stop prometheus

# 2. 备份数据目录
docker run --rm \
  -v zker-prometheus-data:/data \
  -v $(pwd)/backup:/backup \
  alpine tar czf /backup/prometheus-$(date +%Y%m%d).tar.gz -C /data .

# 3. 重启Prometheus
docker-compose -f docker-compose.monitoring.yml start prometheus
```

**备份Grafana Dashboard**:
```bash
# 导出所有Dashboard
curl -u admin:admin http://localhost:3000/api/search?query= | \
  jq -r '.[] | .uid' | \
  while read uid; do
    curl -u admin:admin "http://localhost:3000/api/dashboards/uid/$uid" | \
      jq -r '.dashboard | .id = "" | del(.overwrite)' > "backup/dashboard-$uid.json"
  done
```

### 8.2 配置热更新

**更新Prometheus配置**:
```bash
# 1. 修改配置文件
vim deploy/monitoring/prometheus/prometheus.yml

# 2. 验证配置
docker exec zker-prometheus promtool check config /etc/prometheus/prometheus.yml

# 3. 热加载
curl -X POST http://localhost:9090/-/reload
```

**更新告警规则**:
```bash
# 1. 修改规则文件
vim deploy/monitoring/prometheus/rules/zker-alerts.yml

# 2. 验证规则
docker exec zker-prometheus promtool check rules /etc/prometheus/rules/zker-alerts.yml

# 3. 热加载
curl -X POST http://localhost:9090/-/reload
```

### 8.3 监控数据清理

**删除旧数据**:
```bash
# 通过API删除 (不推荐)
curl -X POST http://localhost:9090/api/v1/admin/tsdb/delete_series \
  -d 'match[]={__name__="http_requests_total"}'

# 通过配置保留期 (推荐)
# 修改prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    cluster: 'coze-studio'
    environment: 'production'

# 设置数据保留时间为30天
storage:
  tsdb:
    path: /prometheus
    retention:
      time: 30d
```

### 8.4 服务扩容

**Prometheus扩容** (使用Thanos):

```yaml
# docker-compose.monitoring.yml
services:
  thanos-sidecar:
    image: quay.io/thanos/thanos:v0.31.0
    container_name: zker-thanos-sidecar
    command:
      - sidecar
      - --prometheus.url=http://prometheus:9090
      - --objconfig.config-file=/etc/thanos/objconfig.yml
    volumes:
      - ./deploy/monitoring/thanos/objconfig.yml:/etc/thanos/objconfig.yml:ro
      - prometheus-data:/prometheus
    networks:
      - monitoring
    restart: unless-stopped
```

---

## 9. 安全加固

### 9.1 访问控制

**Grafana匿名访问禁用**:
```ini
# grafana.ini
[auth.anonymous]
enabled = false

[auth.basic]
enabled = true
```

**Prometheus基础认证**:
```yaml
# prometheus.yml
global:
  external_labels:
    cluster: 'coze-studio'

# 使用nginx反向代理
```

### 9.2 数据加密

**TLS配置**:
```nginx
# nginx.conf
server {
    listen 443 ssl;
    server_name monitor.coze-studio.com;

    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;

    location /prometheus/ {
        proxy_pass http://localhost:9090/;
    }

    location /grafana/ {
        proxy_pass http://localhost:3000/;
    }
}
```

### 9.3 敏感信息保护

**使用环境变量**:
```yaml
# alertmanager.yml
global:
  smtp_auth_password: '${SMTP_PASSWORD}'
  slack_api_url: '${SLACK_WEBHOOK_URL}'
```

**加载环境变量**:
```bash
# .env
SMTP_PASSWORD=your-password
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/xxx
PAGERDUTY_SERVICE_KEY=your-key

# docker-compose.yml
services:
  alertmanager:
    env_file:
      - .env
```

---

## 10. 最佳实践

### 10.1 告警策略

**告警分级**:
- **Critical**: 立即处理 (5分钟内)
- **Warning**: 尽快处理 (30分钟内)
- **Info**: 记录观察 (工作时间内)

**告警抑制**:
- 子告警抑制父告警
- 避免告警风暴
- 合理设置持续时间

**告警分组**:
- 按类别分组 (API, Database, Cache)
- 按严重级别分组
- 避免重复通知

### 10.2 Dashboard设计

**最佳实践**:
1. 一屏展示关键指标
2. 合理使用颜色 (绿=正常, 黄=警告, 红=异常)
3. 提供上下文 (对比基线、去年同期)
4. 支持下钻分析

**命名规范**:
- Dashboard: `[类别] [用途]`
- Panel: `[指标名] [聚合方式]`
- Variable: `类别_属性` (如`tenant_id`)

### 10.3 指标设计

**命名规范**:
- 使用小写字母和下划线
- 使用单位后缀 (`_seconds`, `_bytes`, `_total`)
- 包含标签 (`endpoint`, `tenant_id`, `status`)

**标签规范**:
```promql
# 好的标签
http_requests_total{tenant_id="123", endpoint="/api/bots", method="GET", status="200"}

# 不好的标签
http_requests_total{tenant="123", url="/api/bots", method="GET", statusCode="200"}
```

**指标类型选择**:
- **Counter**: 只增不减的量 (请求总数、错误总数)
- **Gauge**: 可增可减的量 (连接数、内存使用)
- **Histogram**: 分布值 (延迟、请求大小)

---

## 附录

### A. 常用PromQL查询

```promql
# QPS
sum(rate(http_requests_total[5m]))

# P95延迟
histogram_quantile(0.95, http_request_duration_seconds_bucket)

# 错误率
sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))

# 缓存命中率
sum(rate(cache_hit_total[5m])) / (sum(rate(cache_hit_total[5m])) + sum(rate(cache_miss_total[5m])))

# 数据库连接使用率
db_connections_active / db_connections_max

# CPU使用率
sum(rate(process_cpu_seconds_total{mode!="idle"}[5m])) by (instance)

# 内存使用率
process_resident_memory_bytes / node_memory_MemTotal_bytes
```

### B. 参考资源

- [Prometheus官方文档](https://prometheus.io/docs/)
- [Grafana官方文档](https://grafana.com/docs/)
- [AlertManager配置指南](https://prometheus.io/docs/alerting/latest/alertmanager/)
- [PromQL查询语言](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [企业级开发规范手册](./ZKER-企业级开发规范手册_v1.0.md)

### C. 联系方式

- **监控运维团队**: ops-team@coze-studio.com
- **问题反馈**: [GitHub Issues](https://github.com/coze-dev/coze-studio/issues)

---

**版本**: v1.0.0 | **作者**: 研发B | **日期**: 2025-01-03
