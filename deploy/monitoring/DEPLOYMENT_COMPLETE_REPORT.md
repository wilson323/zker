# ZKER 企业级监控系统部署完成报告

**版本**: v2.0.0
**更新**: 2025-01-03
**状态**: ✅ 部署完成
**负责人**: DevOps团队

---

## 📋 目录

- [1. 部署概述](#1-部署概述)
- [2. 配置文件清单](#2-配置文件清单)
- [3. 核心功能](#3-核心功能)
- [4. 部署步骤](#4-部署步骤)
- [5. 验证测试](#5-验证测试)
- [6. 告警规则](#6-告警规则)
- [7. Grafana仪表板](#7-grafana仪表板)
- [8. 性能指标](#8-性能指标)
- [9. 故障排查](#9-故障排查)
- [10. 最佳实践](#10-最佳实践)

---

## 1. 部署概述

### 1.1 监控架构

```
┌─────────────────────────────────────────────────────────────┐
│                      ZKER 监控系统架构                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐         ┌──────────────┐                │
│  │   ZKER API   │────────▶│  Prometheus  │                │
│  │  :8888/metrics│         │    :9090     │                │
│  └──────────────┘         └──────┬───────┘                │
│                                    │                        │
│  ┌──────────────┐                 │                        │
│  │  MySQL       │─────────────────┤                        │
│  │  Exporter    │                 │                        │
│  │  :9104       │                 ▼                        │
│  └──────────────┘         ┌──────────────┐                │
│                           │AlertManager  │                │
│  ┌──────────────┐         │   :9093      │                │
│  │  Redis       │────────▶└──────────────┘                │
│  │  Exporter    │                                        │
│  │  :9121       │         ┌──────────────┐                │
│  └──────────────┘         │   Grafana    │                │
│                           │   :3000      │                │
│  ┌──────────────┐         └──────────────┘                │
│  │  Node        │                                        │
│  │  Exporter    │─────────────────────────────────┐      │
│  │  :9100       │                                  │      │
│  └──────────────┘                                  │      │
│                                                     ▼      │
│                                          ┌──────────────┐  │
│                                          │   告警通知    │  │
│                                          │ (邮件/钉钉)   │  │
│                                          └──────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 核心组件

| 组件 | 版本 | 端口 | 用途 |
|------|------|------|------|
| **Prometheus** | v2.45.0 | 9090 | 指标采集和存储 |
| **Grafana** | v10.0.0 | 3000 | 可视化仪表板 |
| **AlertManager** | v0.26.0 | 9093 | 告警路由和通知 |
| **Node Exporter** | v1.7.0 | 9100 | 系统指标采集 |
| **MySQL Exporter** | v0.15.1 | 9104 | MySQL指标采集 |
| **Redis Exporter** | v1.55.0 | 9121 | Redis指标采集 |
| **cAdvisor** | v0.47.2 | 8080 | 容器指标采集 |

### 1.3 监控覆盖范围

- ✅ **HTTP性能**: QPS、P50/P95/P99延迟、错误率
- ✅ **业务指标**: Bot调用、工作流执行、配额使用
- ✅ **数据库**: 连接池、查询延迟、慢查询
- ✅ **缓存**: 命中率、操作延迟、使用量
- ✅ **系统资源**: CPU、内存、磁盘、网络
- ✅ **租户业务**: 配额、订阅、计费

---

## 2. 配置文件清单

### 2.1 Prometheus配置

| 文件 | 路径 | 说明 |
|------|------|------|
| **基础配置** | `deploy/monitoring/prometheus.yml` | 原始Prometheus配置 |
| **增强配置** | `deploy/monitoring/prometheus-enhanced.yml` | 支持高QPS和业务指标 ⭐ |
| **告警规则** | `deploy/monitoring/alerts.yml` | 企业级告警规则 |
| **增强告警** | `deploy/monitoring/alerts-enhanced.yml` | 高QPS/P99延迟告警 ⭐ |

### 2.2 Grafana配置

| 文件 | 路径 | 说明 |
|------|------|------|
| **业务仪表板** | `deploy/monitoring/dashboards/zker-business-metrics.json` | 核心业务指标 ⭐ |
| **API性能** | `deploy/monitoring/dashboards/api-performance.json` | API性能监控 |
| **系统资源** | `deploy/monitoring/dashboards/system-resources.json` | 系统资源监控 |
| **路由性能** | `deploy/monitoring/dashboards/routing-performance.json` | 智能路由监控 |
| **租户业务** | `deploy/monitoring/dashboards/tenant-business.json` | 租户业务监控 |
| **计费仪表板** | `deploy/monitoring/dashboards/billing-dashboard.json` | 计费指标监控 |

### 2.3 Docker Compose

| 文件 | 路径 | 说明 |
|------|------|------|
| **监控服务** | `docker/docker-compose-monitoring.yml` | 完整监控栈 |

### 2.4 脚本工具

| 脚本 | 路径 | 用途 |
|------|------|------|
| **验证脚本** | `deploy/monitoring/scripts/verify-monitoring.sh` | 健康检查和验证 ⭐ |
| **指标生成器** | `deploy/monitoring/scripts/generate-metrics.sh` | 模拟业务指标 ⭐ |

---

## 3. 核心功能

### 3.1 高QPS监控

✅ **支持10K+ QPS监控**
- 抓取间隔：5秒（高频）
- 采样限制：10000条/次
- 远程写入：支持Thanos长期存储

```yaml
scrape_configs:
  - job_name: 'zker-api'
    scrape_interval: 5s
    sample_limit: 10000
```

### 3.2 P99延迟监控

✅ **P99延迟告警**
- API P99延迟 > 1秒：警告
- API P99延迟 > 2秒：严重
- 数据库P99延迟 > 500ms：警告
- Bot执行P99延迟 > 30秒：警告

```promql
# API P99延迟
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))
```

### 3.3 高错误率监控

✅ **5xx错误率告警**
- 错误率 > 1%：警告
- 错误率 > 5%：严重
- 4xx错误率 > 10%：警告

```promql
# 5xx错误率
(rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])) * 100
```

### 3.4 数据库连接池监控

✅ **连接池使用率告警**
- MySQL连接池 > 80%：警告
- MySQL连接池 > 90%：严重
- Redis连接池 > 80%：警告

```promql
# MySQL连接池使用率
(mysql_global_status_threads_connected / mysql_global_variables_max_connections) * 100
```

---

## 4. 部署步骤

### 4.1 前置条件

- ✅ Docker & Docker Compose已安装
- ✅ 端口未被占用（9090, 3000, 9093, 9100, 9104, 9121）
- ✅ 后端服务已启动（端口8888）

### 4.2 快速部署

```bash
# 1. 进入docker目录
cd docker

# 2. 复制增强配置（可选）
cp ../deploy/monitoring/prometheus-enhanced.yml \
   ./volumes/monitoring/prometheus/prometheus.yml

cp ../deploy/monitoring/alerts-enhanced.yml \
   ./volumes/monitoring/prometheus/alerts.yml

# 3. 启动监控服务
docker compose -f docker-compose-monitoring.yml up -d

# 4. 查看服务状态
docker compose -f docker-compose-monitoring.yml ps

# 5. 查看日志
docker compose -f docker-compose-monitoring.yml logs -f
```

### 4.3 访问地址

| 服务 | URL | 默认账号 |
|------|-----|----------|
| **Prometheus** | http://localhost:9090 | - |
| **Grafana** | http://localhost:3000 | admin/admin |
| **AlertManager** | http://localhost:9093 | - |

### 4.4 导入Grafana仪表板

```bash
# 方法1: 通过UI导入
1. 登录Grafana
2. 导航到 Dashboards -> Import
3. 上传 deploy/monitoring/dashboards/*.json

# 方法2: 通过API导入
curl -X POST http://localhost:3000/api/dashboards/import \
  -u admin:admin \
  -H "Content-Type: application/json" \
  -d @deploy/monitoring/dashboards/zker-business-metrics.json
```

---

## 5. 验证测试

### 5.1 自动化验证

```bash
# 赋予执行权限
chmod +x deploy/monitoring/scripts/verify-monitoring.sh

# 运行验证脚本
cd deploy/monitoring/scripts
./verify-monitoring.sh

# 带性能测试的验证
RUN_PERFORMANCE_TEST=true ./verify-monitoring.sh
```

**预期输出**:
```
[INFO] === 基础健康检查 ===
[✓] Prometheus 健康检查通过
[✓] Grafana 健康检查通过
[✓] AlertManager 健康检查通过

[INFO] === Prometheus 检查 ===
[✓] 指标 http_requests_total 存在 (150 条数据)
[✓] 指标 http_request_duration_seconds 存在 (80 条数据)
[✓] 指标 quota_usage 存在 (20 条数据)

========================================
  验证结果汇总
========================================
通过: 15
警告: 2
失败: 0
```

### 5.2 手动验证

#### 5.2.1 检查Prometheus Targets

访问: http://localhost:9090/targets

**期望状态**: 所有Targets为 `UP`

```
✓ zker-api                 (1/1 up)
✓ mysql                    (1/1 up)
✓ redis                    (1/1 up)
✓ node                     (1/1 up)
✓ cadvisor                 (1/1 up)
```

#### 5.2.2 查询指标

在Prometheus UI中执行以下查询:

```promql
# 1. API QPS
rate(http_requests_total[5m])

# 2. API P99延迟
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))

# 3. 5xx错误率
(rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])) * 100

# 4. MySQL连接池使用率
(mysql_global_status_threads_connected / mysql_global_variables_max_connections) * 100
```

#### 5.2.3 生成测试指标

```bash
# 赋予执行权限
chmod +x deploy/monitoring/scripts/generate-metrics.sh

# 生成API请求指标
cd deploy/monitoring/scripts
./generate-metrics.sh api

# 生成所有业务指标
./generate-metrics.sh all
```

---

## 6. 告警规则

### 6.1 告警分类

| 分类 | 规则数 | 严重级别 |
|------|--------|---------|
| **高QPS告警** | 3 | Warning/Critical |
| **P99延迟告警** | 4 | Warning/Critical |
| **高错误率告警** | 4 | Warning/Critical |
| **数据库连接池** | 4 | Warning/Critical |
| **业务指标** | 4 | Warning/Critical |
| **系统资源** | 3 | Warning |
| **缓存性能** | 1 | Warning |
| **租户业务** | 2 | Warning |
| **路由性能** | 2 | Warning |

### 6.2 关键告警规则

#### 6.2.1 API QPS过高

```yaml
- alert: APIQPSCritical
  expr: rate(http_requests_total[5m]) > 10000
  for: 2m
  labels:
    severity: critical
  annotations:
    summary: "API QPS严重过高"
    description: "QPS为 {{ $value }}/s，需要扩容"
```

**阈值**:
- Warning: QPS > 5000
- Critical: QPS > 10000

#### 6.2.2 API P99延迟过高

```yaml
- alert: APIP99LatencyCritical
  expr: histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m])) > 2
  for: 3m
  labels:
    severity: critical
  annotations:
    summary: "API P99延迟严重"
    description: "P99延迟为 {{ $value }}s"
```

**阈值**:
- Warning: P99 > 1秒
- Critical: P99 > 2秒

#### 6.2.3 5xx错误率过高

```yaml
- alert: API5xxErrorRateCritical
  expr: (rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])) * 100 > 5
  for: 2m
  labels:
    severity: critical
  annotations:
    summary: "API 5xx错误率严重"
    description: "5xx错误率为 {{ $value }}%"
```

**阈值**:
- Warning: 错误率 > 1%
- Critical: 错误率 > 5%

#### 6.2.4 MySQL连接池告警

```yaml
- alert: MySQLConnectionPoolCritical
  expr: (mysql_global_status_threads_connected / mysql_global_variables_max_connections) * 100 > 90
  for: 2m
  labels:
    severity: critical
  annotations:
    summary: "MySQL连接池严重不足"
    description: "连接使用率为 {{ $value }}%"
```

**阈值**:
- Warning: 使用率 > 80%
- Critical: 使用率 > 90%

---

## 7. Grafana仪表板

### 7.1 仪表板清单

| 仪表板 | UID | 说明 |
|--------|-----|------|
| **ZKER业务指标** | `zker-business-metrics` | 核心业务指标监控 ⭐ |
| **API性能** | `api-performance` | API性能分析 |
| **系统资源** | `system-resources` | 服务器资源监控 |
| **路由性能** | `routing-performance` | 智能路由监控 |
| **租户业务** | `tenant-business` | 租户业务分析 |
| **计费仪表板** | `billing-dashboard` | 计费指标分析 |

### 7.2 核心仪表板: ZKER业务指标

**访问路径**: Grafana -> Dashboards -> ZKER业务指标

**面板列表**:
1. **API QPS (请求/秒)**
   - 阈值: 5000 (警告), 10000 (严重)
   - 显示: 各端点实时QPS

2. **API P99延迟**
   - 阈值: 1秒 (警告), 2秒 (严重)
   - 显示: 各端点P99延迟

3. **API 5xx错误率 (%)**
   - 阈值: 1% (警告), 5% (严重)
   - 显示: 各端点5xx错误率

4. **Bot调用QPS**
   - 显示: 各Bot调用频率

5. **MySQL连接池使用率**
   - 阈值: 80% (警告), 90% (严重)
   - 显示: 当前连接池使用情况

6. **租户配额使用率分布**
   - 显示: 各租户配额使用占比

### 7.3 仪表板变量

| 变量名 | 类型 | 说明 |
|--------|------|------|
| `$datasource` | Datasource | Prometheus数据源 |
| `$endpoint` | Query | API端点选择 |
| `$tenant_id` | Query | 租户ID选择 |
| `$bot_id` | Query | Bot ID选择 |

---

## 8. 性能指标

### 8.1 Prometheus性能

| 指标 | 数值 | 说明 |
|------|------|------|
| **采集间隔** | 5秒 | 高频采集 |
| **数据保留** | 30天 | 本地存储 |
| **采样限制** | 10000/次 | 防止OOM |
| **内存使用** | ~2GB | 正常负载 |
| **磁盘使用** | ~10GB/天 | 取决于指标量 |

### 8.2 高QPS支持

| 场景 | QPS | CPU使用率 | 内存使用率 |
|------|-----|----------|-----------|
| **正常** | 1000 | 20% | 30% |
| **高负载** | 5000 | 60% | 50% |
| **峰值** | 10000+ | 80%+ | 70%+ |

**优化建议**:
- QPS > 5000: 建议使用远程存储（Thanos/VictoriaMetrics）
- QPS > 10000: 建议增加Prometheus副本

### 8.3 告警响应时间

| 告警级别 | 响应时间 | 通知延迟 |
|----------|----------|---------|
| **Critical** | < 1分钟 | < 30秒 |
| **Warning** | < 5分钟 | < 1分钟 |

---

## 9. 故障排查

### 9.1 常见问题

#### 9.1.1 Prometheus无法采集指标

**症状**: Targets显示为 `DOWN`

**排查步骤**:
```bash
# 1. 检查服务是否运行
curl http://localhost:8888/metrics

# 2. 检查Prometheus配置
docker logs zker-prometheus

# 3. 检查网络连通性
docker exec zker-prometheus ping coze-studio
```

**解决方案**:
- 确保后端服务 `/metrics` 端点可访问
- 检查Docker网络配置
- 验证Prometheus配置文件语法

#### 9.1.2 Grafana无法查询数据

**症状**: 仪表板显示 "No Data"

**排查步骤**:
```bash
# 1. 检查Prometheus数据源
curl http://localhost:9090/api/v1/query?query=up

# 2. 检查Grafana数据源配置
# 导航到 Configuration -> Data Sources -> Prometheus

# 3. 测试查询
curl http://localhost:9090/api/v1/query?query=http_requests_total
```

**解决方案**:
- 确认Grafana数据源URL正确
- 检查Prometheus是否有数据
- 验证时间范围设置

#### 9.1.3 告警未触发

**症状**: 达到阈值但未收到告警

**排查步骤**:
```bash
# 1. 检查告警规则
curl http://localhost:9090/api/v1/rules

# 2. 检查AlertManager状态
curl http://localhost:9093/api/v2/alerts

# 3. 查看Prometheus日志
docker logs zker-prometheus | grep -i alert
```

**解决方案**:
- 确认告警规则已加载
- 检查 `for` 持续时间设置
- 验证AlertManager通知配置

### 9.2 日志查看

```bash
# 查看Prometheus日志
docker logs -f zker-prometheus

# 查看Grafana日志
docker logs -f zker-grafana

# 查看AlertManager日志
docker logs -f zker-alertmanager

# 查看所有监控服务日志
docker compose -f docker-compose-monitoring.yml logs -f
```

### 9.3 性能优化

#### 9.3.1 减少基数（Cardinality）

```yaml
# 方案1: 限制标签值
metrics.RelabelConfig{
  SourceLabels: []string{"__address__"},
  TargetLabel:  "instance",
  Replacement:  "zker-api-primary",
}

# 方案2: 删除不需要的标签
metrics.RelabelConfig{
  Regex:         "label_to_drop",
  Action:        "labeldrop",
}
```

#### 9.3.2 调整存储配置

```yaml
# prometheus.yml
global:
  scrape_interval: 10s  # 降低抓取频率

# 启动参数
--storage.tsdb.retention.time=15d  # 缩短保留时间
--storage.tsdb.retention.size=50GB  # 限制磁盘使用
```

---

## 10. 最佳实践

### 10.1 告警配置

✅ **DO**:
- 设置合理的 `for` 持续时间（避免抖动）
- 使用分级告警（Warning/Critical）
- 添加详细的 `annotations` 和 `runbook_url`
- 定期审查和更新告警规则

❌ **DON'T**:
- 不要设置过低的阈值（避免告警疲劳）
- 不要忽略 `for` 参数（避免瞬时峰值触发）
- 不要使用过于复杂的PromQL（影响性能）

### 10.2 仪表板设计

✅ **DO**:
- 使用一致的变量命名
- 添加阈值线（可视化告警点）
- 使用颜色区分状态（绿/黄/红）
- 提供下钻能力（点击跳转）

❌ **DON'T**:
- 不要在一个面板中显示过多数据
- 不要使用过短的刷新间隔（影响性能）
- 不要忽略单位设置

### 10.3 指标采集

✅ **DO**:
- 使用Histogram测量延迟（支持分位数）
- 使用Counter测量请求量（支持rate计算）
- 使用Gauge测量瞬时值（连接数、队列长度）
- 限制标签基数（< 10000个唯一值）

❌ **DON'T**:
- 不要在标签中包含高基数值（用户ID、请求ID）
- 不要创建过多指标（< 10000个指标）
- 不要忽略 `le` 标签（Histogram必须）

### 10.4 性能优化

✅ **DO**:
- 使用 `recording rules` 预计算复杂查询
- 配置 `remote_write` 实现长期存储
- 使用 `snapshot` 定期备份数据
- 监控Prometheus自身的性能

❌ **DON'T**:
- 不要在生产环境使用 `--enable-feature=ntpd`
- 不要忽视磁盘IO（使用SSD）
- 不要让Prometheus内存使用超过80%

---

## 11. 附录

### 11.1 常用PromQL查询

```promql
# 1. API QPS
sum(rate(http_requests_total{endpoint=~"/api/.*"}[5m])) by (endpoint)

# 2. API P50/P95/P99延迟
histogram_quantile(0.50, rate(http_request_duration_seconds_bucket[5m]))
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))

# 3. 错误率
(rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])) * 100

# 4. 数据库连接池
(mysql_global_status_threads_connected / mysql_global_variables_max_connections) * 100

# 5. Redis内存使用率
(redis_memory_used_bytes / redis_memory_max_bytes) * 100

# 6. 缓存命中率
(rate(cache_hit_total[5m]) / (rate(cache_hit_total[5m]) + rate(cache_miss_total[5m]))) * 100

# 7. Bot调用失败率
(rate(bot_invocation_total{status="failed"}[5m]) / rate(bot_invocation_total[5m])) * 100

# 8. CPU使用率
100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance))

# 9. 内存使用率
(1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) * 100

# 10. 磁盘使用率
(1 - node_filesystem_avail_bytes / node_filesystem_size_bytes) * 100
```

### 11.2 有用的链接

- [Prometheus官方文档](https://prometheus.io/docs/)
- [Grafana官方文档](https://grafana.com/docs/)
- [PromQL入门指南](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [最佳实践白皮书](https://prometheus.io/docs/practices/naming/)

### 11.3 监控指标参考

| 指标名称 | 类型 | 标签 | 说明 |
|----------|------|------|------|
| `http_requests_total` | Counter | method, endpoint, status | HTTP请求总数 |
| `http_request_duration_seconds` | Histogram | method, endpoint | HTTP请求延迟 |
| `quota_usage` | Gauge | tenant_id, resource_type | 配额使用量 |
| `bot_invocation_total` | Counter | tenant_id, bot_id, status | Bot调用总数 |
| `db_query_duration_seconds` | Histogram | database, operation, table | 数据库查询延迟 |
| `cache_hit_total` | Counter | cache_type, key_prefix | 缓存命中数 |
| `cache_miss_total` | Counter | cache_type, key_prefix | 缓存未命中数 |

---

## 12. 部署确认清单

部署完成后，请确认以下项目：

- [ ] Prometheus UI可访问 (http://localhost:9090)
- [ ] Grafana UI可访问 (http://localhost:3000)
- [ ] AlertManager UI可访问 (http://localhost:9093)
- [ ] 所有Targets状态为 `UP`
- [ ] 核心指标已采集（http_requests_total, bot_invocation_total等）
- [ ] Grafana仪表板已导入并显示数据
- [ ] 告警规则已加载（Prometheus -> Alerts）
- [ ] 验证脚本全部通过
- [ ] 性能测试通过（如需要）

---

## 13. 联系方式

如有问题，请联系：

- **DevOps团队**: devops@zker.com
- **后端团队**: backend@zker.com
- **文档维护**: docs@zker.com

---

**文档版本**: v2.0.0
**最后更新**: 2025-01-03
**维护者**: DevOps团队
