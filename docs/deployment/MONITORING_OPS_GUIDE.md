# 监控运维手册

**版本**: v1.0.0
**最后更新**: 2025-01-02
**维护者**: DevOps团队

---

## 📊 监控架构概览

### 监控组件

- **Prometheus**: 指标采集和存储
- **Grafana**: 可视化大盘
- **Alertmanager**: 告警路由和通知
- **Node Exporter**: 主机指标采集
- **cAdvisor**: 容器指标采集

### 数据流

```
应用 → Metrics Exporter → Prometheus → Alertmanager → Slack/PagerDuty
                             ↓
                          Grafana
```

---

## 🎛️ Grafana使用指南

### Dashboard列表

1. **API Performance Dashboard** (`uid: api-performance`)
   - 用途: 监控API性能指标
   - 关键指标: QPS、P95延迟、错误率
   - 刷新频率: 30秒

2. **Tenant Business Dashboard** (`uid: tenant-business`)
   - 用途: 监控租户业务指标
   - 关键指标: 租户数、订阅分布、配额使用率
   - 刷新频率: 1分钟

3. **Routing Performance Dashboard** (`uid: routing-performance`)
   - 用途: 监控智能路由性能
   - 关键指标: 路由延迟、意图置信度、Bot健康状态
   - 刷新频率: 30秒

4. **System Resources Dashboard** (`uid: system-resources`)
   - 用途: 监控系统资源使用
   - 关键指标: CPU、内存、磁盘、网络
   - 刷新频率: 30秒

### 导入Dashboard

```bash
# 方式1: 使用Grafana API
curl -X POST http://localhost:3000/api/dashboards/import \
  -H "Content-Type: application/json" \
  -u admin:admin \
  -d @deploy/monitoring/dashboards/api-performance.json

# 方式2: 通过Web界面
1. 登录Grafana (http://grafana.coze-studio.com)
2. 点击 "+" → "Import"
3. 上传JSON文件或粘贴内容
4. 选择Prometheus数据源
5. 点击"Import"
```

---

## 🔔 告警规则说明

### 告警级别

**Critical（严重）**:
- API错误率 > 15%
- 磁盘空间 < 5%
- Bot不健康
- 租户配额超限

**Warning（警告）**:
- API错误率 > 5%
- CPU使用率 > 80%
- 内存使用率 > 85%
- 数据库连接数 > 80%

### 添加自定义告警

```yaml
# 在deploy/monitoring/alerts.yml中添加
groups:
  - name: custom_alerts
    interval: 1m
    rules:
      - alert: MyCustomAlert
        expr: my_metric > threshold
        for: 5m
        labels:
          severity: warning
          team: backend
        annotations:
          summary: "自定义告警"
          description: "告警详情: {{ $value }}"
```

重新加载告警规则：
```bash
kubectl apply -f deploy/monitoring/alerts.yml
# Prometheus会自动重载配置
```

---

## 📈 常用PromQL查询

### API性能

```promql
# QPS
sum(rate(http_requests_total{job="coze-studio"}[5m]))

# P95延迟
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# 错误率
sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))
```

### 租户业务

```promql
# 活跃租户数
count(tenant_info{status="active"})

# 配额使用率
quota_usage_percent{tenant_id="xxx"}

# Top 10配额消费者
topk(10, quota_usage_percent)
```

### 系统资源

```promql
# CPU使用率
100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance))

# 内存使用率
(1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) * 100

# 磁盘使用率
(1 - node_filesystem_avail_bytes / node_filesystem_size_bytes) * 100
```

---

## 🚨 故障处理流程

### 1. 收到告警

**立即检查**:
1. 登录Grafana查看相关Dashboard
2. 确认告警详情和影响范围
3. 判断严重程度（P1/P2/P3/P4）

### 2. 初步诊断

```bash
# 检查Pod状态
kubectl get pods -n coze-studio-prod | grep -v Running

# 查看相关日志
kubectl logs -n coze-studio-prod deployment/backend --tail=100

# 检查资源使用
kubectl top pods -n coze-studio-prod
kubectl top nodes
```

### 3. 根据告警类型处理

#### API错误率过高

```bash
# 1. 检查错误日志
kubectl logs -n coze-studio-prod deployment/backend | grep ERROR

# 2. 检查数据库连接
kubectl exec -it -n coze-studio-prod deployment/backend -- \
  mysql -h mysql -u root -p -e "SHOW PROCESSLIST;"

# 3. 如必要，回滚部署
./scripts/rollback.sh coze-studio-prod backend
```

#### CPU/内存过高

```bash
# 1. 查看资源使用情况
kubectl top pods -n coze-studio-prod

# 2. 检查HPA状态
kubectl describe hpa backend -n coze-studio-prod

# 3. 手动扩容（如需要）
kubectl scale deployment backend --replicas=10 -n coze-studio-prod
```

#### 磁盘空间不足

```bash
# 1. 检查磁盘使用
kubectl exec -it -n coze-studio-prod deployment/backend -- df -h

# 2. 清理日志文件
kubectl exec -it -n coze-studio-prod deployment/backend -- \
  sh -c "find /var/log -name '*.log' -mtime +7 -delete"

# 3. 扩容PVC（如需要）
kubectl patch pvc data-pvc -n coze-studio-prod -p '{"spec":{"resources":{"requests":{"storage":"200Gi"}}}}'
```

---

## 🔄 定期维护任务

### 每日

- 检查Grafana Dashboard（早晚各一次）
- 处理Critical告警
- 检查磁盘空间使用率

### 每周

- 审查告警规则有效性
- 清理旧备份文件
- 性能趋势分析

### 每月

- 审查监控指标覆盖度
- 优化告警阈值
- 容量规划评估

---

## 📞 联系方式

- **DevOps值班**: devops@coze-studio.com
- **紧急联系**: +86-xxx-xxxx-xxxx
- **Runbook**: https://docs.coze-studio.com/runbooks

---

**保持监控，主动运维！** 📊
