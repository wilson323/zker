# ZKER 监控系统快速启动指南

**版本**: v2.0.0
**预计耗时**: 10分钟

---

## 🚀 快速开始（3步部署）

### 第1步: 启动监控服务

```bash
# 进入docker目录
cd docker

# 启动监控栈
docker compose -f docker-compose-monitoring.yml up -d

# 等待服务启动（约30秒）
docker compose -f docker-compose-monitoring.yml ps
```

**期望输出**:
```
NAME                    STATUS         PORTS
zker-prometheus         Up             0.0.0.0:9090->9090/tcp
zker-grafana            Up             0.0.0.0:3000->3000/tcp
zker-alertmanager       Up             0.0.0.0:9093->9093/tcp
zker-mysqld-exporter    Up             0.0.0.0:9104->9104/tcp
zker-redis-exporter     Up             0.0.0.0:9121->9121/tcp
zker-node-exporter      Up             0.0.0.0:9100->9100/tcp
zker-cadvisor           Up             0.0.0.0:8080->8080/tcp
```

---

### 第2步: 访问监控UI

#### Prometheus (指标采集)
- **URL**: http://localhost:9090
- **功能**: 查询指标、查看Targets、告警规则

**快速检查**:
1. 访问 http://localhost:9090/targets
2. 确认所有Targets状态为 `UP`
3. 访问 http://localhost:9090/graph
4. 输入查询: `rate(http_requests_total[5m])`

#### Grafana (可视化)
- **URL**: http://localhost:3000
- **默认账号**: admin / admin
- **功能**: 仪表板、告警通知

**快速导入仪表板**:
1. 登录后，导航到 Dashboards -> Import
2. 上传 `deploy/monitoring/dashboards/zker-business-metrics.json`
3. 点击 "Import"

#### AlertManager (告警管理)
- **URL**: http://localhost:9093
- **功能**: 查看活跃告警、静默规则

---

### 第3步: 验证监控数据

```bash
# 赋予执行权限
chmod +x deploy/monitoring/scripts/verify-monitoring.sh

# 运行验证
cd deploy/monitoring/scripts
./verify-monitoring.sh
```

**期望输出**:
```
[✓] Prometheus 健康检查通过
[✓] Grafana 健康检查通过
[✓] AlertManager 健康检查通过
[✓] 指标 http_requests_total 存在 (150 条数据)
[✓] 指标 http_request_duration_seconds 存在 (80 条数据)

通过: 15
警告: 2
失败: 0
```

---

## 📊 核心指标查询

### 在Prometheus中执行以下查询

#### 1. API性能指标

```promql
# QPS (每秒请求数)
sum(rate(http_requests_total[5m])) by (endpoint)

# P99延迟 (最慢1%的请求)
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))

# 5xx错误率
(rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])) * 100
```

#### 2. 数据库指标

```promql
# MySQL连接池使用率
(mysql_global_status_threads_connected / mysql_global_variables_max_connections) * 100

# 数据库查询P99延迟
histogram_quantile(0.99, rate(db_query_duration_seconds_bucket[5m]))
```

#### 3. 缓存指标

```promql
# Redis内存使用率
(redis_memory_used_bytes / redis_memory_max_bytes) * 100

# 缓存命中率
(rate(cache_hit_total[5m]) / (rate(cache_hit_total[5m]) + rate(cache_miss_total[5m]))) * 100
```

#### 4. Bot业务指标

```promql
# Bot调用QPS
sum(rate(bot_invocation_total[5m])) by (bot_id)

# Bot失败率
(rate(bot_invocation_total{status="failed"}[5m]) / rate(bot_invocation_total[5m])) * 100
```

---

## 🚨 告警规则验证

### 查看已加载的告警规则

访问: http://localhost:9090/alerts

**关键告警**:
- ✅ APIQPSWarning: QPS > 5000
- ✅ APIQPSCritical: QPS > 10000
- ✅ APIP99LatencyHigh: P99 > 1秒
- ✅ API5xxErrorRateCritical: 错误率 > 5%
- ✅ MySQLConnectionPoolCritical: 连接池 > 90%

### 查看AlertManager中的告警

访问: http://localhost:9093/#/alerts

---

## 🧪 生成测试数据（可选）

### 生成API请求指标

```bash
# 赋予执行权限
chmod +x deploy/monitoring/scripts/generate-metrics.sh

# 生成API请求（100 QPS，持续60秒）
cd deploy/monitoring/scripts
./generate-metrics.sh api
```

### 生成所有业务指标

```bash
# 生成所有类型的指标（Bot、DB、缓存、配额）
./generate-metrics.sh all
```

---

## 🎯 验证清单

部署完成后，请确认以下项目：

### 基础服务
- [ ] Prometheus UI可访问 (http://localhost:9090)
- [ ] Grafana UI可访问 (http://localhost:3000)
- [ ] AlertManager UI可访问 (http://localhost:9093)

### Targets健康状态
- [ ] zker-api: UP
- [ ] mysql: UP
- [ ] redis: UP
- [ ] node: UP
- [ ] cadvisor: UP

### 指标采集
- [ ] http_requests_total 有数据
- [ ] http_request_duration_seconds 有数据
- [ ] bot_invocation_total 有数据
- [ ] db_query_duration_seconds 有数据
- [ ] cache_hit_total 有数据

### 仪表板
- [ ] ZKER业务指标仪表板已导入
- [ ] 仪表板显示实时数据
- [ ] 仪表板刷新正常（默认5秒）

### 告警规则
- [ ] 告警规则已加载（Prometheus -> Alerts）
- [ ] AlertManager可接收告警

---

## 🔧 常见问题

### Q1: Targets显示DOWN

**原因**: 后端服务未启动或网络不通

**解决方案**:
```bash
# 检查后端服务
curl http://localhost:8888/metrics

# 检查网络连通性
docker exec zker-prometheus ping coze-studio

# 重启Prometheus
docker restart zker-prometheus
```

### Q2: Grafana仪表板显示"No Data"

**原因**: Prometheus数据源配置错误或无数据

**解决方案**:
```bash
# 1. 检查Prometheus是否有数据
curl http://localhost:9090/api/v1/query?query=up

# 2. 在Grafana中检查数据源配置
# Configuration -> Data Sources -> Prometheus
# URL: http://prometheus:9090

# 3. 测试查询
# 在Grafana Explore中执行: rate(http_requests_total[5m])
```

### Q3: 告警未触发

**原因**: 未达到阈值或持续时间不足

**解决方案**:
```bash
# 1. 检查告警规则
curl http://localhost:9090/api/v1/rules

# 2. 查看当前告警状态
curl http://localhost:9093/api/v2/alerts

# 3. 生成测试数据触发告警
./deploy/monitoring/scripts/generate-metrics.sh all
```

---

## 📈 性能基准

| 指标 | 目标值 | 说明 |
|------|--------|------|
| **Prometheus采集间隔** | 5秒 | 高频采集 |
| **Grafana刷新间隔** | 5秒 | 实时监控 |
| **告警响应时间** | < 1分钟 | Critical告警 |
| **数据保留时间** | 30天 | 本地存储 |
| **支持QPS** | 10K+ | 高性能采集 |

---

## 📞 获取帮助

- **文档**: `deploy/monitoring/DEPLOYMENT_COMPLETE_REPORT.md`
- **验证脚本**: `deploy/monitoring/scripts/verify-monitoring.sh`
- **DevOps团队**: devops@zker.com

---

**版本**: v2.0.0
**更新**: 2025-01-03
