# ZKER Grafana 监控大盘使用指南

**版本**: v1.0.0
**最后更新**: 2025-01-01

---

## 📋 目录

- [监控大盘概述](#监控大盘概述)
- [快速开始](#快速开始)
- [监控面板说明](#监控面板说明)
- [告警集成](#告警集成)
- [最佳实践](#最佳实践)

---

## 🎯 监控大盘概述

### 监控大盘功能

**ZKER 企业级监控大盘** 提供全面的系统监控，涵盖：

1. **系统资源监控**
   - CPU 使用率（阈值：80% 警告，95% 严重）
   - 内存使用率（阈值：85% 警告）
   - 磁盘使用率（阈值：85% 警告，95% 严重）

2. **API 性能监控**
   - P95 响应时间（阈值：500ms 警告，2000ms 严重）
   - API 错误率（阈值：1% 警告，5% 严重）
   - API 请求量（RPS）

3. **业务指标监控**
   - Bot 调用速率
   - Token 使用速率
   - 配额使用率（阈值：80% 警告，90% 严重）

4. **数据库性能监控**
   - 连接池使用率（阈值：80% 警告）
   - 查询延迟（P95阈值：1秒警告，5秒严重）

### 监控维度

| 维度 | 说明 | 示例 |
|------|------|------|
| **按实例** | 不同服务器实例的指标 | instance-1, instance-2 |
| **按端点** | API 端点维度 | /api/bots, /api/tenants |
| **按租户** | 租户维度（多租户） | tenant_001, tenant_002 |
| **按 Bot** | Bot 维度 | bot_001, bot_002 |
| **按状态** | HTTP 状态码维度 | 200, 404, 500 |

---

## 🚀 快速开始

### 1. 前置条件

确保以下组件已安装并运行：

```bash
# Prometheus（指标采集）
prometheus --config.file=/etc/prometheus/prometheus.yml

# Grafana（可视化）
grafana-server --config=/etc/grafana/grafana.ini
```

### 2. 导入监控大盘

**方式一：通过 Grafana UI 导入**

1. 登录 Grafana：`http://localhost:3000`
2. 导航到 **Dashboards** → **Import**
3. 上传 JSON 文件：`deploy/monitoring/grafana-dashboard.json`
4. 点击 **Import** 导入

**方式二：通过 API 导入**

```bash
# 使用 Grafana API 导入
curl -X POST \
  -H "Content-Type: application/json" \
  -d @deploy/monitoring/grafana-dashboard.json \
  http://admin:admin@localhost:3000/api/dashboards/db
```

**方式三：使用 Provisioning（自动化）**

将 JSON 文件复制到 Grafana provisioning 目录：

```bash
cp deploy/monitoring/grafana-dashboard.json \
   /var/lib/grafana/dashboards/zker-enterprise-dashboard.json

# 重启 Grafana
systemctl restart grafana-server
```

### 3. 配置 Prometheus 数据源

1. 在 Grafana 中，导航到 **Configuration** → **Data Sources**
2. 添加新的 Prometheus 数据源：
   - **Name**: `Prometheus`
   - **URL**: `http://localhost:9090`
   - **Access**: `Server (default)`
3. 点击 **Save & Test**

### 4. 验证监控大盘

导入后，应该能看到：

- ✅ 4 个监控面板组（系统资源、API 性能、业务指标、数据库）
- ✅ 11 个监控图表
- ✅ 实时数据更新（每30秒）
- ✅ 阈值告警显示（彩色区域）

---

## 📊 监控面板说明

### 面板 1：系统资源监控

#### 1.1 CPU 使用率

**查询语句**:
```promql
100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance)
```

**指标说明**:
- 显示每个服务器实例的 CPU 使用率
- 时间窗口：过去5分钟
- 阈值设置：
  - 🟢 绿色：< 80%
  - 🔴 红色：≥ 80%

#### 1.2 内存使用率

**查询语句**:
```promql
(1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) * 100
```

**阈值设置**:
- 🟢 绿色：< 85%
- 🔴 红色：≥ 85%

#### 1.3 磁盘使用率

**查询语句**:
```promql
(1 - node_filesystem_avail_bytes{fstype!~"tmpfs|fuse.*"} / node_filesystem_size_bytes) * 100
```

**阈值设置**:
- 🟢 绿色：< 85%
- 🟡 黄色：85% - 95%
- 🔴 红色：≥ 95%

---

### 面板 2：API 性能监控

#### 2.1 API P95 响应时间

**查询语句**:
```promql
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```

**阈值设置**:
- 🟢 绿色：< 500ms
- 🟡 黄色：500ms - 2秒
- 🔴 红色：≥ 2秒

#### 2.2 API 错误率

**查询语句**:
```promql
(rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])) * 100
```

**阈值设置**:
- 🟢 绿色：< 1%
- 🟡 黄色：1% - 5%
- 🔴 红色：≥ 5%

#### 2.3 API 请求量（RPS）

**查询语句**:
```promql
rate(http_requests_total[5m])
```

**说明**:
- 显示每个 API 端点的请求速率
- 单位：请求/秒（RPS）
- 使用堆叠面积图

---

### 面板 3：业务指标监控

#### 3.1 Bot 调用速率

**查询语句**:
```promql
rate(bot_invocation_total[5m])
```

**说明**:
- 显示每个 Bot 的调用速率
- 用于监控 Bot 使用情况和趋势

#### 3.2 Token 使用速率

**查询语句**:
```promql
rate(bot_token_usage_total[5m])
```

**说明**:
- 显示每个 Bot 的 Token 使用速率
- 用于监控 Token 成本和趋势

#### 3.3 Bot 配额使用率

**查询语句**:
```promql
(quota_usage{resource_type="bots"} / quota_limit{resource_type="bots"}) * 100
```

**阈值设置**:
- 🟢 绿色：< 80%
- 🟡 黄色：80% - 90%
- 🔴 红色：≥ 90%

---

### 面板 4：数据库性能监控

#### 4.1 数据库连接池使用率

**查询语句**:
```promql
db_connections_active / db_connections_active_max * 100
```

**阈值设置**:
- 🟢 绿色：< 80%
- 🔴 红色：≥ 80%

#### 4.2 数据库查询 P95 延迟

**查询语句**:
```promql
histogram_quantile(0.95, rate(db_query_duration_seconds_bucket[5m]))
```

**阈值设置**:
- 🟢 绿色：< 1秒
- 🟡 黄色：1秒 - 5秒
- 🔴 红色：≥ 5秒

---

## 🚨 告警集成

### 告警规则配置

监控大盘已与 Prometheus 告警规则集成，详见 `alerts.yml`。

**告警通知**：
- 警告级别：发送到企业微信/钉钉/邮件
- 严重级别：发送短信 + 电话 + PagerDuty

### Grafana 内置告警

也可以在 Grafana 中创建告警规则：

1. 点击图表标题 → **Edit**
2. 切换到 **Alert** 标签页
3. 配置告警条件：
   ```yaml
   # 示例：CPU 使用率告警
   WHEN avg() OF query(A, 5m, now)
   IS ABOVE 80
   ```
4. 配置通知渠道（Webhook/Email）

---

## 🎨 自定义监控面板

### 添加新图表

1. 点击面板右上角 **+ Add panel**
2. 选择可视化类型（Time series / Gauge / Stat 等）
3. 输入 PromQL 查询语句
4. 配置样式和阈值
5. 保存面板

### 推荐查询语句

**按租户统计 Bot 数量**:
```promql
count(bot_invocation_total) by (tenant_id)
```

**按租户统计 Token 使用**:
```promql
sum(rate(bot_token_usage_total[1h])) by (tenant_id)
```

**工作流执行成功率**:
```promql
sum(rate(workflow_execution_total{status="success"}[5m])) by (tenant_id) /
sum(rate(workflow_execution_total[5m])) by (tenant_id) * 100
```

---

## 📈 最佳实践

### 1. 监控大盘维护

**定期检查清单**：
- [ ] 每周检查告警阈值是否合理
- [ ] 每月审查监控面板覆盖度
- [ ] 每季度优化查询性能
- [ ] 更新监控大盘文档

### 2. 性能优化

**查询优化技巧**:
- 使用 `rate()` 计算速率，避免 Counter 重置问题
- 使用 `by()` 分组，避免时间序列爆炸
- 使用时间窗口 `[5m]` 平衡准确性和性能
- 避免过于复杂的正则表达式

**示例 - 优化前**:
```promql
# 查询所有可能的 HTTP 状态码
http_requests_total{status=~".*"}
```

**示例 - 优化后**:
```promql
# 仅查询 5xx 错误
http_requests_total{status=~"5.."}
```

### 3. 告警降噪

**避免告警风暴**:
- 设置合理的持续时间（`for` 子句）
- 使用告警抑制规则
- 分组告警通知
- 在低峰期调整告警阈值

### 4. 监控大盘权限管理

**访问控制**：
- 生产环境大盘：仅运维和开发团队可访问
- 测试环境大盘：所有团队成员可访问
- 使用 Grafana 的 Role-Based Access Control (RBAC)

---

## 🔗 相关文档

- [Prometheus 告警规则说明](./README.md)
- [性能测试指南](../../tests/performance/README.md)
- [业务指标定义](../../backend/infra/monitoring/metrics/business_metrics.go)

---

## 📞 支持

如有问题，请联系：
- 监控团队：monitoring-team@coze-studio.com
- 运维团队：ops-team@coze-studio.com

---

**更新日志**：

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2025-01-01 | v1.0.0 | 初始版本，11个监控面板，4个面板组 | Claude AI |
