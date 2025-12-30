# ZKER 组织中心监控配置指南

## 📋 目录

- [监控架构](#监控架构)
- [快速开始](#快速开始)
- [指标说明](#指标说明)
- [Grafana 仪表盘](#grafana-仪表盘)
- [告警配置](#告警配置)
- [故障排查](#故障排查)
- [最佳实践](#最佳实践)

---

## 🏗️ 监控架构

### 架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                        监控架构                             │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐ │
│  │ Go 应用服务   │────▶│ Prometheus   │────▶│   Grafana    │ │
│  │ :8889/metrics│     │   :9090      │     │   :3000      │ │
│  └──────────────┘     └──────────────┘     └──────────────┘ │
│         │                    │                     │        │
│         │                    ▼                     ▼        │
│         │            ┌──────────────┐     ┌──────────────┐ │
│         │            │Alertmanager  │     │  仪表盘展示   │ │
│         │            │   :9093      │     └──────────────┘ │
│         │            └──────────────┘                       │
│         │                    │                              │
│         ▼                    ▼                              │
│  ┌──────────────┐     ┌──────────────┐                     │
│  │ Node Exporter│     │ MySQL Exporter│                    │
│  │   :9100      │     │   :9104      │                     │
│  └──────────────┘     └──────────────┘                     │
│                                                              │
│  ┌──────────────┐     ┌──────────────┐                     │
│  │Redis Exporter│     │ES Exporter   │                     │
│  │   :9121      │     │   :9114      │                     │
│  └──────────────┘     └──────────────┘                     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 监控组件

| 组件 | 版本 | 端口 | 说明 |
|------|------|------|------|
| **Prometheus** | v2.45.0 | 9090 | 指标采集和存储 |
| **Grafana** | 10.0.0 | 3000 | 可视化展示 |
| **Alertmanager** | v0.26.0 | 9093 | 告警管理 |
| **Node Exporter** | v1.7.0 | 9100 | 系统指标 |
| **MySQL Exporter** | v0.15.1 | 9104 | 数据库指标 |
| **Redis Exporter** | v1.55.0 | 9121 | 缓存指标 |
| **Elasticsearch Exporter** | v1.6.0 | 9114 | 搜索指标 |
| **cAdvisor** | v0.47.2 | 8080 | 容器指标 |
| **Loki** | 2.9.2 | 3100 | 日志聚合 |
| **Promtail** | 2.9.2 | 9080 | 日志采集 |

---

## 🚀 快速开始

### 1. 环境准备

确保已安装 Docker 和 Docker Compose：

```bash
docker --version    # Docker 20.10+
docker compose version  # Docker Compose 2.0+
```

### 2. 配置环境变量

创建 `.env` 文件（如果不存在）：

```bash
# 监控服务端口
PROMETHEUS_PORT=9090
GRAFANA_PORT=3000
ALERTMANAGER_PORT=9093

# Exporter 端口
NODE_EXPORTER_PORT=9100
MYSQLD_EXPORTER_PORT=9104
REDIS_EXPORTER_PORT=9121
ES_EXPORTER_PORT=9114
CADVISOR_PORT=8080

# Grafana 配置
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=admin

# MySQL Exporter DSN
MYSQL_EXPORTER_DSN=root:root@(zker-mysql:3306)/

# 网络配置
MONITORING_NETWORK_SUBNET=172.21.0.0/16
```

### 3. 创建必要的配置文件

#### Prometheus 配置 (`volumes/monitoring/prometheus/prometheus.yml`)

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    cluster: 'zker-production'
    environment: 'production'

alerting:
  alertmanagers:
    - static_configs:
        - targets:
            - 'alertmanager:9093'

rule_files:
  - 'alerts.yml'

scrape_configs:
  # Go 应用服务
  - job_name: 'zker-api'
    static_configs:
      - targets: ['host.docker.internal:8889']
        labels:
          service: 'zker-api'
          team: 'backend'
    metrics_path: '/metrics'
    scrape_interval: 10s
    scrape_timeout: 5s

  # Node Exporter
  - job_name: 'node'
    static_configs:
      - targets: ['node_exporter:9100']
        labels:
          service: 'node-exporter'
    scrape_interval: 15s

  # MySQL Exporter
  - job_name: 'mysql'
    static_configs:
      - targets: ['mysqld_exporter:9104']
        labels:
          service: 'mysql'
    scrape_interval: 30s

  # Redis Exporter
  - job_name: 'redis'
    static_configs:
      - targets: ['redis_exporter:9121']
        labels:
          service: 'redis'
    scrape_interval: 30s

  # Elasticsearch Exporter
  - job_name: 'elasticsearch'
    static_configs:
      - targets: ['elasticsearch_exporter:9114']
        labels:
          service: 'elasticsearch'
    scrape_interval: 30s

  # Prometheus 自身
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
        labels:
          service: 'prometheus'
```

#### 告警规则 (`volumes/monitoring/prometheus/alerts.yml`)

```yaml
groups:
  # 组织中心告警规则
  - name: org_center_alerts
    interval: 30s
    rules:
      # 组织操作延迟过高
      - alert: HighOrgOperationLatency
        expr: |
          histogram_quantile(0.95,
            rate(organization_operation_duration_seconds_bucket[5m])
          ) > 1
        for: 5m
        labels:
          severity: warning
          team: backend
        annotations:
          summary: "组织操作延迟过高"
          description: "P95延迟 {{ $value }}s 超过1秒"

      # 员工离职率异常
      - alert: HighEmployeeResignationRate
        expr: |
          rate(employee_resignation_total[7d]) /
          sum(employee_total{emp_status="active"}) > 0.1
        for: 1h
        labels:
          severity: warning
          team: hr
        annotations:
          summary: "员工离职率异常"
          description: "周离职率 {{ $value | humanizePercentage }} 超过10%"

      # 孤立组织节点检测
      - alert: OrphanOrgNodesDetected
        expr: sum(org_orphan_total) > 0
        for: 10m
        labels:
          severity: warning
          team: backend
        annotations:
          summary: "检测到孤立组织节点"
          description: "发现 {{ $value }} 个孤立组织节点"

      # HR流程处理超时
      - alert: HRProcessTimeout
        expr: |
          histogram_quantile(0.95,
            rate(hr_onboarding_duration_seconds_bucket[5m])
          ) > 3600
        for: 10m
        labels:
          severity: warning
          team: hr
        annotations:
          summary: "HR流程处理超时"
          description: "入职流程P95时长 {{ $value }}s 超过1小时"

      # 组织层级深度过深
      - alert: OrgTreeTooDeep
        expr: max(organization_depth_max) > 10
        for: 5m
        labels:
          severity: warning
          team: backend
        annotations:
          summary: "组织层级深度过深"
          description: "最大层级深度 {{ $value }} 超过10层"

      # 岗位填充率过低
      - alert: LowPositionFillRate
        expr: position_fill_rate < 0.7
        for: 1h
        labels:
          severity: info
          team: hr
        annotations:
          summary: "岗位填充率过低"
          description: "Level {{ $labels.position_level }} 岗位填充率 {{ $value | humanizePercentage }} 低于70%"

      # 目录缓存命中率低
      - alert: LowDirectoryCacheHitRate
        expr: directory_cache_hit_rate < 0.7
        for: 15m
        labels:
          severity: info
          team: backend
        annotations:
          summary: "目录缓存命中率低"
          description: "{{ $labels.cache_type }} 缓存命中率 {{ $value | humanizePercentage }} 低于70%"
```

#### Alertmanager 配置 (`volumes/monitoring/alertmanager/alertmanager.yml`)

```yaml
global:
  resolve_timeout: 5m
  smtp_smarthost: 'smtp.example.com:587'
  smtp_from: 'alerts@example.com'
  smtp_auth_username: 'alerts@example.com'
  smtp_auth_password: 'password'

route:
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'default'
  routes:
    - match:
        severity: critical
      receiver: 'critical'
    - match:
        severity: warning
      receiver: 'warning'
    - match:
        severity: info
      receiver: 'info'

receivers:
  - name: 'default'
    email_configs:
      - to: 'team@example.com'
        headers:
          Subject: '[ZKER] {{ .GroupLabels.alertname }}'

  - name: 'critical'
    email_configs:
      - to: 'oncall@example.com'
        headers:
          Subject: '[CRITICAL] {{ .GroupLabels.alertname }}'
    webhook_configs:
      - url: 'http://your-webhook-url/alert'
        send_resolved: true

  - name: 'warning'
    email_configs:
      - to: 'team@example.com'
        headers:
          Subject: '[WARNING] {{ .GroupLabels.alertname }}'

  - name: 'info'
    email_configs:
      - to: 'info@example.com'
        headers:
          Subject: '[INFO] {{ .GroupLabels.alertname }}'
```

### 4. 启动监控服务

```bash
# 进入 docker 目录
cd docker

# 启动监控服务
docker compose -f docker-compose-monitoring.yml up -d

# 查看服务状态
docker compose -f docker-compose-monitoring.yml ps

# 查看日志
docker compose -f docker-compose-monitoring.yml logs -f
```

### 5. 访问监控界面

| 服务 | URL | 默认账号 |
|------|-----|---------|
| **Prometheus** | http://localhost:9090 | - |
| **Grafana** | http://localhost:3000 | admin/admin |
| **Alertmanager** | http://localhost:9093 | - |
| **cAdvisor** | http://localhost:8080 | - |

### 6. 导入 Grafana 仪表盘

1. 登录 Grafana (http://localhost:3000)
2. 导航到 **Dashboards** → **Import**
3. 上传仪表盘 JSON 文件：
   - `backend/infra/monitoring/grafana/dashboards/org-dashboard.json`
4. 选择数据源：**Prometheus**
5. 点击 **Import**

### 7. 验证监控

```bash
# 检查 Prometheus targets
curl http://localhost:9090/api/v1/targets

# 检查指标是否正常采集
curl http://localhost:9090/api/v1/query?query=organization_total

# 检查告警规则
curl http://localhost:9090/api/v1/rules

# 查看 Alertmanager 状态
curl http://localhost:9093/api/v2/status
```

---

## 📊 指标说明

### 组织管理指标 (Organization)

| 指标名称 | 类型 | 标签 | 说明 |
|---------|------|------|------|
| `organization_total` | Gauge | tenant_id, org_type, status | 组织总数 |
| `organization_creation_total` | Counter | tenant_id, org_type, result | 组织创建总数 |
| `organization_deletion_total` | Counter | tenant_id, org_type, result | 组织删除总数 |
| `organization_update_total` | Counter | tenant_id, org_type, update_type, result | 组织更新总数 |
| `organization_move_total` | Counter | tenant_id, org_type, result | 组织移动总数 |
| `organization_depth_max` | Gauge | tenant_id | 组织最大层级深度 |
| `organization_children_count` | Histogram | tenant_id, org_type | 组织子节点数分布 |
| `organization_operation_duration_seconds` | Histogram | tenant_id, operation_type | 组织操作延迟 |

### 部门管理指标 (Department)

| 指标名称 | 类型 | 标签 | 说明 |
|---------|------|------|------|
| `department_total` | Gauge | tenant_id, org_id, status | 部门总数 |
| `department_member_total` | Gauge | tenant_id, dept_id | 部门成员总数 |
| `department_member_count` | Histogram | tenant_id, org_id | 部门成员数分布 |
| `department_depth_max` | Gauge | tenant_id, org_id | 部门最大层级深度 |
| `department_operation_duration_seconds` | Histogram | tenant_id, operation_type | 部门操作延迟 |

### 员工管理指标 (Employee)

| 指标名称 | 类型 | 标签 | 说明 |
|---------|------|------|------|
| `employee_total` | Gauge | tenant_id, org_id, emp_status | 员工总数 |
| `employee_creation_total` | Counter | tenant_id, org_id, result | 员工创建总数 |
| `employee_deletion_total` | Counter | tenant_id, org_id, result | 员工删除总数 |
| `employee_transfer_total` | Counter | tenant_id, transfer_type, result | 员工调转总数 |
| `employee_promotion_total` | Counter | tenant_id, result | 员工晋升总数 |
| `employee_resignation_total` | Counter | tenant_id, org_id, reason_type | 员工离职总数 |
| `employee_active_days` | Histogram | tenant_id, org_id | 员工在职天数 |
| `employee_active_daily` | Gauge | tenant_id, org_id | 每日活跃员工数 |
| `employee_active_weekly` | Gauge | tenant_id, org_id | 每周活跃员工数 |
| `employee_active_monthly` | Gauge | tenant_id, org_id | 每月活跃员工数 |

### 岗位管理指标 (Position)

| 指标名称 | 类型 | 标签 | 说明 |
|---------|------|------|------|
| `position_total` | Gauge | tenant_id, dept_id, position_level, status | 岗位总数 |
| `position_occupied_total` | Gauge | tenant_id, position_level | 岗位占用总数 |
| `position_vacant_total` | Gauge | tenant_id, position_level | 岗位空缺总数 |
| `position_fill_rate` | Gauge | tenant_id, position_level | 岗位填充率 |

### HR流程指标 (HR Process)

| 指标名称 | 类型 | 标签 | 说明 |
|---------|------|------|------|
| `hr_onboarding_duration_seconds` | Histogram | tenant_id, step | 入职办理延迟 |
| `hr_resignation_duration_seconds` | Histogram | tenant_id, step | 离职办理延迟 |
| `hr_transfer_duration_seconds` | Histogram | tenant_id, transfer_type | 调转办理延迟 |
| `hr_approval_duration_seconds` | Histogram | tenant_id, approval_type | HR审批延迟 |
| `hr_process_total` | Counter | tenant_id, process_type, result | HR流程总数 |

### 组织架构指标 (Org Structure)

| 指标名称 | 类型 | 标签 | 说明 |
|---------|------|------|------|
| `org_tree_health` | Gauge | tenant_id, health_dimension | 组织树健康度 |
| `org_orphan_total` | Gauge | tenant_id, node_type | 孤立组织节点数 |
| `org_query_duration_seconds` | Histogram | tenant_id, query_type | 组织查询延迟 |
| `org_traversal_duration_seconds` | Histogram | tenant_id, traversal_type | 组织遍历延迟 |

### 目录服务指标 (Directory)

| 指标名称 | 类型 | 标签 | 说明 |
|---------|------|------|------|
| `directory_query_duration_seconds` | Histogram | tenant_id, query_type | 目录查询延迟 |
| `directory_search_duration_seconds` | Histogram | tenant_id, search_type | 目录搜索延迟 |
| `directory_cache_hit_rate` | Gauge | tenant_id, cache_type | 目录缓存命中率 |
| `directory_index_size_bytes` | Gauge | tenant_id, index_type | 目录索引大小 |

---

## 📈 Grafana 仪表盘

### 组织中心监控大盘

**仪表盘 ID**: `zker-org-dashboard`

**面板分组**:

1. **组织概览**
   - 组织总数
   - 部门总数
   - 在职员工总数
   - 岗位总数
   - 组织类型分布（饼图）
   - 组织状态分布（饼图）
   - 员工状态分布（饼图）

2. **员工管理**
   - 员工入职趋势（时序图）
   - 员工离职趋势（时序图）
   - 员工在职时长分布（热力图）
   - 活跃员工统计（时序图）

3. **HR流程**
   - HR流程处理时长（时序图）
   - HR审批时长（时序图）
   - HR流程成功率（饼图）

4. **组织架构**
   - 组织最大层级深度（Stat）
   - 部门最大层级深度（Stat）
   - 孤立组织节点数（Stat）
   - 组织子节点数分布（热力图）
   - 部门成员数分布（热力图）

5. **岗位管理**
   - 岗位职级分布（饼图）
   - 岗位占用vs空缺（饼图）
   - 岗位填充率趋势（时序图）

6. **性能监控**
   - 组织操作延迟（时序图）
   - 部门操作延迟（时序图）
   - 组织查询延迟（时序图）
   - 目录查询延迟（时序图）
   - 目录缓存命中率（时序图）

### 导入仪表盘

```bash
# 方法1: 通过 Web UI
1. 登录 Grafana
2. Dashboards → Import
3. Upload JSON file
4. 选择 Prometheus 数据源
5. Import

# 方法2: 通过 API
curl -X POST http://localhost:3000/api/dashboards/import \
  -H "Content-Type: application/json" \
  -u admin:admin \
  -d @/path/to/org-dashboard.json
```

---

## 🚨 告警配置

### 告警规则分类

#### P0 - 严重（立即处理）

| 告警名称 | 条件 | 说明 |
|---------|------|------|
| `HighOrgOperationLatency` | P95延迟 > 1s | 组织操作严重延迟 |
| `OrphanOrgNodesDetected` | 孤立节点 > 0 | 检测到孤立组织节点 |

#### P1 - 高（1小时内处理）

| 告警名称 | 条件 | 说明 |
|---------|------|------|
| `HRProcessTimeout` | 入职流程 > 1h | HR流程处理超时 |
| `OrgTreeTooDeep` | 层级深度 > 10 | 组织层级过深 |

#### P2 - 中（当天处理）

| 告警名称 | 条件 | 说明 |
|---------|------|------|
| `HighEmployeeResignationRate` | 离职率 > 10% | 员工离职率异常 |
| `LowPositionFillRate` | 填充率 < 70% | 岗位填充率过低 |

#### P3 - 低（计划处理）

| 告警名称 | 条件 | 说明 |
|---------|------|------|
| `LowDirectoryCacheHitRate` | 命中率 < 70% | 目录缓存命中率低 |

### 配置告警通知

#### Email 通知

在 `alertmanager.yml` 中配置：

```yaml
global:
  smtp_smarthost: 'smtp.example.com:587'
  smtp_from: 'alerts@example.com'

receivers:
  - name: 'critical'
    email_configs:
      - to: 'oncall@example.com'
        headers:
          Subject: '[CRITICAL] {{ .GroupLabels.alertname }}'
```

#### Webhook 通知

```yaml
receivers:
  - name: 'webhook'
    webhook_configs:
      - url: 'http://your-webhook-url/alert'
        send_resolved: true
```

#### 钉钉通知（推荐）

创建钉钉通知模板 `volumes/monitoring/alertmanager/templates/dingtalk.tmpl`：

```yaml
{{ define "dingtalk.title" }}
[{{ .Status | toUpper }}{{ if eq .Status "firing" }}告警{{ else }}恢复{{ end }}]
{{ .GroupLabels.alertname }}
{{ end }}

{{ define "dingtalk.content" }}
**告警名称**: {{ .GroupLabels.alertname }}

**告警级别**: {{ .CommonLabels.severity }}

**告警详情**:
{{ range .Alerts }}
**实例**: {{ .Labels.instance }}
**值**: {{ .Value }}
**描述**: {{ .Annotations.description }}
{{ end }}

**开始时间**: {{ .StartsAt.Format "2006-01-02 15:04:05" }}
{{ end }}
```

配置 Alertmanager 使用钉钉模板：

```yaml
templates:
  - '/etc/alertmanager/templates/dingtalk.tmpl'

receivers:
  - name: 'dingtalk'
    webhook_configs:
      - url: 'https://oapi.dingtalk.com/robot/send?access_token=YOUR_TOKEN'
        send_resolved: true
```

---

## 🔧 故障排查

### 常见问题

#### 1. Prometheus 无法采集指标

**症状**: Prometheus UI 显示 target 为 `DOWN`

**排查步骤**:

```bash
# 检查服务是否运行
docker compose -f docker-compose-monitoring.yml ps

# 检查 metrics 端点
curl http://localhost:8889/metrics

# 检查 Prometheus 配置
docker exec zker-prometheus promtool check config /etc/prometheus/prometheus.yml

# 查看 Prometheus 日志
docker logs zker-prometheus --tail 100
```

**解决方案**:

- 确保 Go 服务正在运行且暴露 metrics 端口
- 检查防火墙设置
- 验证 Prometheus 配置中的 targets 地址正确

#### 2. Grafana 无法连接 Prometheus

**症状**: Grafana 数据源测试失败

**排查步骤**:

```bash
# 测试 Prometheus 连接
curl http://localhost:9090/api/v1/query?query=up

# 检查 Grafana 容器网络
docker exec zker-grafana ping prometheus

# 查看 Grafana 日志
docker logs zker-grafana --tail 100
```

**解决方案**:

- 确保 Prometheus 和 Grafana 在同一网络
- 检查 Grafana 数据源配置中的 URL
- 验证网络连通性

#### 3. 告警未触发

**症状**: 满足告警条件但未收到通知

**排查步骤**:

```bash
# 检查告警规则
curl http://localhost:9090/api/v1/rules

# 查看 Alertmanager 状态
curl http://localhost:9093/api/v2/status

# 检查 Alertmanager 日志
docker logs zker-alertmanager --tail 100
```

**解决方案**:

- 确保告警规则已加载
- 检查告警的 `for` 子句时长
- 验证 Alertmanager 配置
- 检查通知渠道配置

#### 4. 指标数据缺失

**症状**: Grafana 仪表盘显示 "No Data"

**排查步骤**:

```bash
# 检查指标是否存在
curl 'http://localhost:9090/api/v1/query?query=organization_total'

# 查看最近采集的数据
curl 'http://localhost:9090/api/v1/query?query=organization_total&time='$(date +%s)

# 检查 Go 应用是否正确暴露指标
curl http://localhost:8889/metrics | grep organization
```

**解决方案**:

- 确保业务代码正确调用指标记录函数
- 检查 Prometheus middleware 是否正确注册
- 验证指标标签值是否正确

### 性能优化

#### Prometheus 优化

```yaml
# prometheus.yml
global:
  scrape_interval: 30s        # 降低采集频率
  evaluation_interval: 30s    # 降低评估频率

# 存储优化
--storage.tsdb.retention.time=15d   # 减少数据保留时间
--storage.tsdb.retention.size=10GB  # 限制数据大小
```

#### Grafana 优化

- 使用变量（Variables）减少面板数量
- 合理设置刷新间隔（建议 30s-1m）
- 使用查询缓存
- 限制时间范围查询

---

## 📚 最佳实践

### 1. 指标命名规范

```go
// ✅ Good: 清晰的指标名称
organization_total
employee_creation_total
hr_onboarding_duration_seconds

// ❌ Bad: 模糊的指标名称
org_count
emp_metric
hr_time
```

### 2. 标签使用规范

```go
// ✅ Good: 适当的标签基数
OrganizationTotal.WithLabelValues(tenantID, orgType, status).Set(1)

// ❌ Bad: 过多的标签值
OrganizationTotal.WithLabelValues(employeeID, timestamp, ipAddress).Set(1)
```

### 3. 指标类型选择

```go
// ✅ Counter: 只增不减的计数
OrganizationCreationTotal.WithLabelValues(...).Inc()

// ✅ Gauge: 可增可减的当前值
OrganizationTotal.WithLabelValues(...).Set(1)

// ✅ Histogram: 延迟、大小等分布
OrganizationOperationDuration.WithLabelValues(...).Observe(duration)
```

### 4. 采样率控制

```go
// ✅ Good: 记录关键操作
metrics.RecordOrganizationCreation(tenantID, orgType, "success", duration)

// ❌ Bad: 记录所有细节（高基数）
metrics.RecordOrganizationDetail(tenantID, empID, timestamp, detail)
```

### 5. 告警规则设计

```yaml
# ✅ Good: 合理的阈值和持续时间
- alert: HighOrgOperationLatency
  expr: histogram_quantile(0.95, rate(...[5m])) > 1
  for: 5m  # 持续5分钟才触发

# ❌ Bad: 过于敏感
- alert: OrgSlow
  expr: rate(...[1m]) > 0.1
  for: 1m  # 1分钟就触发，误报高
```

### 6. 仪表盘设计

- **使用变量**: 如 `tenant_id`, `org_type`
- **合理分组**: 按业务功能分组面板
- **设置阈值**: 使用颜色和阈值线标识异常
- **添加描述**: 每个面板添加清晰的描述

### 7. 文档维护

- **指标文档**: 为每个指标添加 `Help` 说明
- **告警文档**: 为每个告警添加 runbook 链接
- **变更记录**: 记录告警规则的变更历史

---

## 🔗 相关资源

- [Prometheus 官方文档](https://prometheus.io/docs/)
- [Grafana 官方文档](https://grafana.com/docs/)
- [Alertmanager 官方文档](https://prometheus.io/docs/alerting/latest/alertmanager/)
- [Go client_golang 文档](https://github.com/prometheus/client_golang)
- [ZKER 监控最佳实践](../企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)

---

## 📞 支持

如有问题，请联系：

- **技术支持**: tech-support@example.com
- **监控团队**: monitoring-team@example.com
- **值班电话**: +86-xxx-xxxx-xxxx

---

**最后更新**: 2025-01-01
**维护者**: ZKER 监控团队
**版本**: v1.0.0
