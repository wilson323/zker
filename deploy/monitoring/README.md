# ZKER 监控告警规则说明

**版本**: v1.0.0
**最后更新**: 2025-01-01

## 📋 目录

- [告警规则概述](#告警规则概述)
- [告警规则分类](#告警规则分类)
- [使用指南](#使用指南)
- [告警响应流程](#告警响应流程)
- [告警调优建议](#告警调优建议)

---

## 🎯 告警规则概述

本配置文件定义了ZKER系统的完整监控告警规则，涵盖：

- **12个告警组** (Alert Groups)
- **40+条告警规则** (Alert Rules)
- **多维度监控** (系统、API、业务、成本)

### 告警级别

| 级别 | 说明 | 响应时间 | 示例 |
|------|------|----------|------|
| **Critical** | 严重故障，影响核心业务 | 立即响应（<5分钟） | 磁盘使用率>95%、API错误率>15% |
| **Warning** | 性能降级或潜在风险 | 尽快处理（<30分钟） | CPU>80%、配额使用率>80% |

---

## 📊 告警规则分类

### 1. 系统资源告警 (system_resources)

监控CPU、内存、磁盘等基础资源。

| 告警名称 | 触发条件 | 级别 | 处理建议 |
|---------|---------|------|---------|
| HighCPUUsage | CPU使用率 > 80% (持续5分钟) | Warning | 检查高进程、考虑扩容 |
| HighMemoryUsage | 内存使用率 > 85% (持续5分钟) | Warning | 检查内存泄漏、优化缓存 |
| HighDiskUsage | 磁盘使用率 > 85% (持续5分钟) | Warning | 清理日志、扩容磁盘 |
| DiskSpaceCritical | 磁盘使用率 > 95% (持续2分钟) | Critical | 立即扩容、清理空间 |

### 2. API性能告警 (api_performance)

监控API响应时间和错误率。

| 告警名称 | 触发条件 | 级别 | 处理建议 |
|---------|---------|------|---------|
| HighAPILatency | P95响应时间 > 2秒 | Warning | 优化慢查询、增加缓存 |
| HighAPIErrorRate | 5xx错误率 > 5% | Warning | 检查应用日志、修复Bug |
| APIErrorRateCritical | 5xx错误率 > 15% | Critical | 立即回滚、排查故障 |

### 3. 数据库告警 (database_performance)

监控数据库连接、查询性能。

| 告警名称 | 触发条件 | 级别 | 处理建议 |
|---------|---------|------|---------|
| HighDBConnections | 连接使用率 > 80% | Warning | 优化连接池、扩容数据库 |
| HighDBQueryLatency | P95查询延迟 > 1秒 | Warning | 优化索引、重构查询 |
| SlowDBQuery | P99查询延迟 > 5秒 | Warning | 分析慢查询日志、优化SQL |

### 4. 缓存告警 (cache_performance)

监控缓存命中率和性能。

| 告警名称 | 触发条件 | 级别 | 处理建议 |
|---------|---------|------|---------|
| LowCacheHitRate | 缓存命中率 < 70% (持续10分钟) | Warning | 优化缓存策略、增加缓存 |
| CacheHitRateCritical | 缓存命中率 < 50% (持续5分钟) | Critical | 检查缓存配置、重启缓存服务 |

### 5. 配额告警 (quota_alerts)

监控租户配额使用情况。

| 告警名称 | 触发条件 | 级别 | 处理建议 |
|---------|---------|------|---------|
| HighQuotaUsage | 配额使用率 > 80% | Warning | 通知用户、升级订阅 |
| QuotaNearLimit | 配额使用率 > 90% | Critical | 限制操作、通知升级 |
| HighQuotaDenialRate | 配额拒绝率 > 5% | Warning | 分析拒绝原因、优化流程 |

### 6. Bot服务告警 (bot_service_health)

监控Bot调用性能和稳定性。

| 告警名称 | 触发条件 | 级别 | 处理建议 |
|---------|---------|------|---------|
| HighBotErrorRate | Bot失败率 > 10% | Warning | 检查Bot配置、模型状态 |
| HighBotLatency | P95响应时间 > 30秒 | Warning | 优化Bot逻辑、增加超时 |
| HighTokenUsage | Token使用 > 100万/小时 | Warning | 分析使用量、优化Prompt |

### 7. 权限系统告警 (permission_system)

监控权限检查性能和拒绝率。

| 告警名称 | 触发条件 | 级别 | 处理建议 |
|---------|---------|------|---------|
| HighPermissionDenialRate | 权限拒绝率 > 20% | Warning | 分析拒绝原因、调整权限 |
| HighPermissionCheckLatency | P95检查延迟 > 10ms | Warning | 优化权限缓存、简化规则 |

### 8. 存储告警 (storage_alerts)

监控存储使用情况。

| 告警名称 | 触发条件 | 级别 | 处理建议 |
|---------|---------|------|---------|
| HighStorageUsage | 存储使用率 > 80% | Warning | 清理无用文件、升级套餐 |
| StorageUsageCritical | 存储使用率 > 95% | Critical | 立即清理、暂停写入 |

### 9. 模型服务告警 (model_service_alerts)

监控模型调用和成本。

| 告警名称 | 触发条件 | 级别 | 处理建议 |
|---------|---------|------|---------|
| HighModelErrorRate | 模型失败率 > 15% | Warning | 检查模型状态、切换模型 |
| HighModelLatency | P95响应时间 > 60秒 | Warning | 优化Prompt、检查网络 |
| HighModelCost | 每小时成本 > $10 | Warning | 分析使用量、优化调用 |

### 10. 业务健康度告警 (business_health)

监控业务整体健康度。

| 告警名称 | 触发条件 | 级别 | 处理建议 |
|---------|---------|------|---------|
| LowHealthScore | 健康度评分 < 60 | Warning | 查看详细指标、制定改进计划 |
| SLAComplianceLow | SLA合规率 < 95% | Warning | 分析违规原因、优化服务 |
| SLAViolationCritical | SLA合规率 < 90% | Critical | 紧急修复、通知管理层 |

### 11. 订阅告警 (subscription_alerts)

监控订阅状态。

| 告警名称 | 触发条件 | 级别 | 处理建议 |
|---------|---------|------|---------|
| SubscriptionExpiringSoon | 订阅7天内过期 | Warning | 提前通知续费 |
| TooManySuspendedTenants | 暂停租户比例 > 10% | Warning | 分析原因、优化服务 |

### 12. 事故告警 (incident_alerts)

监控事故发生情况。

| 告警名称 | 触发条件 | 级别 | 处理建议 |
|---------|---------|------|---------|
| CriticalIncident | 发生严重事故 | Critical | 立即响应、启动应急预案 |
| FrequentIncidents | 每小时事故 > 5次 | Warning | 分析根本原因、改进系统 |

---

## 🚀 使用指南

### 部署告警规则

1. **复制告警规则文件到Prometheus服务器**：
   ```bash
   cp deploy/monitoring/alerts.yml /etc/prometheus/alerts/
   ```

2. **配置Prometheus加载告警规则**：
   ```yaml
   # /etc/prometheus/prometheus.yml
   global:
     evaluation_interval: 30s
     external_labels:
       cluster: 'coze-studio'
       env: 'production'

   rule_files:
     - '/etc/prometheus/alerts/*.yml'

   alerting:
     alertmanagers:
       - static_configs:
           - targets:
               - 'localhost:9093'
   ```

3. **重启Prometheus**：
   ```bash
   systemctl restart prometheus
   ```

4. **验证规则加载**：
   ```bash
   # 查看已加载的规则
   curl http://localhost:9090/api/v1/rules

   # 查看当前活跃的告警
   curl http://localhost:9090/api/v1/alerts
   ```

### 配置AlertManager路由

创建AlertManager配置文件 `alertmanager.yml`：

```yaml
global:
  resolve_timeout: 5m

# 路由配置
route:
  group_by: ['alertname', 'cluster', 'service']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'default'

  子路由：按严重级别路由
  routes:
    - match:
        severity: critical
      receiver: 'critical_alerts'
      continue: true

    - match:
        severity: warning
      receiver: 'warning_alerts'

    - match:
        category: quota
      receiver: 'quota_team'

    - match:
        category: bot
      receiver: 'bot_team'

# 接收器配置
receivers:
  - name: 'default'
    webhook_configs:
      - url: 'http://localhost:8000/webhook/alert'

  - name: 'critical_alerts'
    pagerduty_configs:
      - service_key: 'YOUR_PAGERDUTY_KEY'
    email_configs:
      - to: 'oncall@example.com'
        headers:
          Subject: '[CRITICAL] {{ .GroupLabels.alertname }}'

  - name: 'warning_alerts'
    email_configs:
      - to: 'alerts@example.com'
        headers:
          Subject: '[WARNING] {{ .GroupLabels.alertname }}'

  - name: 'quota_team'
    webhook_configs:
      - url: 'http://internal-webhook/quota-alerts'

  - name: 'bot_team'
    webhook_configs:
      - url: 'http://internal-webhook/bot-alerts'

# 抑制规则
inhibit_rules:
  - source_match:
      severity: 'critical'
    target_match:
      severity: 'warning'
    equal: ['alertname', 'instance']
```

### 告警测试

测试告警是否正常工作：

```bash
# 1. 触发一个测试告警
curl -XPOST http://localhost:9090/api/v1/alerts \
  -d '[{
    "labels": {
      "alertname": "HighCPUUsage",
      "severity": "warning",
      "instance": "test-server"
    },
    "annotations": {
      "summary": "测试告警",
      "description": "这是一个测试告警"
    }
  }]'

# 2. 查看告警状态
curl http://localhost:9090/api/v1/alerts

# 3. 查看AlertManager状态
curl http://localhost:9093/api/v1/status
```

---

## 🔄 告警响应流程

### Critical告警响应流程

1. **立即通知** (0-5分钟)
   - 发送短信/电话给值班人员
   - 创建事故工单
   - 通知管理层

2. **初步响应** (5-15分钟)
   - 确认告警严重性
   - 评估影响范围
   - 启动应急预案

3. **故障处理** (15-60分钟)
   - 执行Runbook中的处理步骤
   - 更新事故状态
   - 协调相关团队

4. **恢复验证** (60-120分钟)
   - 确认问题解决
   - 验证服务恢复
   - 解除告警

5. **事后分析** (1-3天)
   - 编写事故报告
   - 制定改进措施
   - 更新Runbook

### Warning告警响应流程

1. **通知** (0-30分钟)
   - 发送邮件/企业微信
   - 记录到告警系统

2. **分析处理** (30分钟-4小时)
   - 分析告警原因
   - 执行预防性措施
   - 记录处理结果

3. **持续监控**
   - 观察告警趋势
   - 必要时升级为Critical

---

## 🔧 告警调优建议

### 调优原则

1. **避免告警风暴**
   - 设置合理的持续时间 (for)
   - 使用分组和抑制规则
   - 设置合适的阈值

2. **提高告警质量**
   - 定期审查误报率
   - 调整阈值和持续时间
   - 优化告警描述和Runbook

3. **分层告警**
   - Warning: 早期预警，预防性处理
   - Critical: 严重故障，立即响应

### 调优流程

1. **收集反馈** (每周)
   - 统计告警数量和准确率
   - 收集值班人员反馈
   - 记录误报和漏报

2. **分析优化** (每月)
   - 分析告警趋势
   - 调整不合理的阈值
   - 合并重复的告警
   - 拆分过于宽泛的告警

3. **验证改进** (每季度)
   - 测试新规则有效性
   - 更新Runbook
   - 培训团队

### 示例调优场景

#### 场景1：CPU告警过于频繁

**问题**：HighCPUUsage告警每天触发数百次

**分析**：
- 大部分是瞬时峰值（持续1-2分钟）
- 实际影响不大，造成告警疲劳

**优化**：
```yaml
# 修改前
- alert: HighCPUUsage
  expr: 100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance) > 80
  for: 5m  # 持续5分钟

# 修改后
- alert: HighCPUUsage
  expr: 100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance) > 85  # 阈值提高到85%
  for: 10m  # 持续时间延长到10分钟
```

#### 场景2：数据库慢查询告警漏报

**问题**：实际有慢查询但未触发告警

**分析**：
- P99阈值5秒太宽松
- 某些表的查询确实很慢但未超过5秒

**优化**：
```yaml
# 修改前
- alert: SlowDBQuery
  expr: histogram_quantile(0.99, rate(db_query_duration_seconds_bucket[5m])) > 5
  for: 5m

# 修改后
- alert: SlowDBQuery
  expr: histogram_quantile(0.95, rate(db_query_duration_seconds_bucket[5m])) > 2  # P95超过2秒
  for: 10m  # 持续10分钟
```

---

## 📚 相关文档

- [Runbook目录](https://docs.coze-studio.com/runbooks)
- [监控最佳实践](https://docs.coze-studio.com/monitoring/best-practices)
- [告警管理规范](https://docs.coze-studio.com/operations/alert-guidelines)

---

## 📞 联系方式

如有问题，请联系：
- 监控团队：monitoring-team@coze-studio.com
- 运维团队：ops-team@coze-studio.com
- 紧急联系：+86-xxx-xxxx-xxxx

---

**更新日志**：

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2025-01-01 | v1.0.0 | 初始版本，12个告警组，40+条规则 | Claude AI |
