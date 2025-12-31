# 监控系统快速部署指南

**版本**: v1.0
**日期**: 2025-01-03
**状态**: ✅ 准备就绪

---

## 📋 部署前检查

### ✅ 环境检查
- ✅ Docker: v29.1.2
- ✅ Docker Compose: v2.40.3
- ✅ 配置文件: 已就绪
- ✅ 仪表板: 已创建 (6个)

---

## 🚀 快速部署（5分钟）

### 步骤1: 复制增强型配置

```bash
# 复制增强型Prometheus配置
cp D:/code/coze-studio/deploy/monitoring/prometheus-enhanced.yml \
   D:/code/coze-studio/docker/volumes/monitoring/prometheus/prometheus.yml

# 复制增强型告警规则
cp D:/code/coze-studio/deploy/monitoring/alerts-enhanced.yml \
   D:/code/coze-studio/docker/volumes/monitoring/prometheus/alerts.yml
```

### 步骤2: 创建数据目录

```bash
cd D:/code/coze-studio/docker
mkdir -p data/monitoring/prometheus
mkdir -p data/monitoring/grafana
mkdir -p data/monitoring/alertmanager
```

### 步骤3: 启动监控服务

```bash
cd D:/code/coze-studio/docker
docker-compose -f docker-compose-monitoring.yml up -d
```

### 步骤4: 验证服务状态

```bash
# 检查容器状态
docker-compose -f docker-compose-monitoring.yml ps

# 应该看到以下服务运行中：
# - zker-prometheus (port 9090)
# - zker-grafana (port 3000)
# - zker-alertmanager (port 9093)
# - zker-node-exporter (port 9100)
```

### 步骤5: 访问监控界面

```bash
# Prometheus
open http://localhost:9090

# Grafana (默认账号: admin/admin)
open http://localhost:3000

# AlertManager
open http://localhost:9093
```

---

## 📊 Grafana仪表板配置

### 导入仪表板

1. 登录Grafana (admin/admin)
2. 点击左侧菜单 "Dashboard" → "Import"
3. 上传以下JSON文件：

| 仪表板 | 文件路径 | 说明 |
|--------|----------|------|
| API性能 | `deploy/monitoring/dashboards/api-performance.json` | QPS、延迟、错误率 |
| 租户业务 | `deploy/monitoring/dashboards/tenant-business.json` | 租户指标 |
| 路由性能 | `deploy/monitoring/dashboards/routing-performance.json` | 路由监控 |
| 系统资源 | `deploy/monitoring/dashboards/system-resources.json` | CPU、内存、磁盘 |
| 计费 | `deploy/monitoring/dashboards/billing-dashboard.json` | 计费指标 |
| 综合指标 | `deploy/monitoring/dashboards/zker-business-metrics.json` | 综合业务 |

### 配置Prometheus数据源

1. 在Grafana中: Configuration → Data Sources
2. 添加Prometheus
3. URL: `http://prometheus:9090`
4. 点击 "Save & Test"

---

## ⚠️ 常见问题

### 问题1: 端口被占用

```bash
# 检查端口占用
netstat -ano | findstr :9090
netstat -ano | findstr :3000
netstat -ano | findstr :9093

# 修改端口（在docker-compose-monitoring.yml中）
PROMETHEUS_PORT=9091
GRAFANA_PORT=3001
ALERTMANAGER_PORT=9094
```

### 问题2: 权限错误

```bash
# Windows: 以管理员身份运行PowerShell
# Linux/Mac: 使用sudo
sudo docker-compose -f docker-compose-monitoring.yml up -d
```

### 问题3: 配置文件错误

```bash
# 验证YAML语法
docker run --rm -v \
  D:/code/coze-studio/docker/volumes/monitoring/prometheus:/etc/prometheus \
  prom/prometheus:v2.45.0 \
  promtool check config /etc/prometheus/prometheus.yml
```

---

## 📈 验证监控数据

### 1. 检查Prometheus Targets

```bash
# 访问
open http://localhost:9090/targets

# 应该看到所有targets为 "UP" 状态
```

### 2. 查询指标

```bash
# 在Prometheus UI中运行以下查询：

# API请求总数
rate(http_requests_total[5m])

# P99延迟
histogram_quantile(0.99, api_request_duration_ms)

# 错误率
rate(http_errors_total[5m]) / rate(http_requests_total[5m])
```

### 3. 验证告警规则

```bash
# 访问
open http://localhost:9090/alerts

# 应该看到9类告警规则已加载
```

---

## 🛑 停止监控服务

```bash
cd D:/code/coze-studio/docker
docker-compose -f docker-compose-monitoring.yml down

# 清理数据（可选）
docker-compose -f docker-compose-monitoring.yml down -v
```

---

## 📝 配置文件说明

### Prometheus配置 (`prometheus.yml`)

**关键配置**:
- 抓取间隔: 10秒（高频）
- 数据保留: 30天
- 告警管理器: 集成
- 远程写入: 可选

**服务发现**:
- Coze Studio API: `coze-studio:8888`
- MySQL Exporter: `mysql-exporter:9104`
- Redis Exporter: `redis-exporter:9121`
- Node Exporter: `node-exporter:9100`

### 告警规则 (`alerts.yml`)

**9类告警**:
1. API QPS告警 (>10000)
2. P99延迟告警 (>200ms)
3. 错误率告警 (>1%)
4. MySQL连接池告警 (>180)
5. Redis连接告警
6. 磁盘空间告警 (>80%)
7. 内存使用告警 (>85%)
8. CPU使用告警 (>80%)
9. 租户配额告警

---

## 🎯 下一步

部署完成后，应该看到：

✅ **Prometheus**: http://localhost:9090
- Targets状态: 所有UP
- 指标采集: 正常

✅ **Grafana**: http://localhost:3000
- 仪表板: 已导入6个
- 数据源: Prometheus已连接
- 数据展示: 实时更新

✅ **AlertManager**: http://localhost:9093
- 告警规则: 已加载9类
- 通知配置: 可配置Webhook/Email

---

## 📞 故障排查

### 查看日志

```bash
# 查看所有服务日志
docker-compose -f docker-compose-monitoring.yml logs -f

# 查看特定服务日志
docker logs zker-prometheus -f
docker logs zker-grafana -f
docker logs zker-alertmanager -f
```

### 重启服务

```bash
# 重启所有服务
docker-compose -f docker-compose-monitoring.yml restart

# 重启特定服务
docker restart zker-prometheus
```

---

**部署时间**: 约5-10分钟
**维护需求**: 低（自动化重启）
**资源消耗**: 低（<2GB内存，<10GB磁盘）

**最后更新**: 2025-01-03
