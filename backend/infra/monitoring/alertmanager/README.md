# Prometheus Alertmanager 运维手册

## 概述

本手册涵盖 ZKER 监告警系统的配置、管理和故障排查。

---

## 📋 目录

1. [快速开始](#快速开始)
2. [告警分类](#告警分类)
3. [告警处理流程](#告警处理流程)
4. [常见告警处理](#常见告警处理)
5. [配置修改](#配置修改)
6. [测试验证](#测试验证)
7. [故障排查](#故障排查)

---

## 快速开始

### 1. 启动 Alertmanager

```bash
# 使用 Docker 启动
docker run -d \
  --name alertmanager \
  -p 9093:9093 \
  -v $(pwd)/alertmanager.yml:/etc/alertmanager/alertmanager.yml \
  -v $(pwd)/templates:/etc/alertmanager/templates \
  prom/alertmanager:latest

# 使用 Docker Compose
cd backend/infra/monitoring
docker compose up -d alertmanager
```

### 2. 访问 Alertmanager UI

```
http://localhost:9093
```

### 3. 查看告警状态

```bash
# 查看所有活跃告警
curl http://localhost:9093/api/v2/alerts | jq

# 查看告警分组
curl http://localhost:9093/api/v2/alerts/groups | jq

# 查看接收器状态
curl http://localhost:9093/api/v2/receivers | jq
```

---

## 告警分类

### 1. API 告警（api-team）

| 告警名称 | 级别 | 触发条件 | 通知渠道 |
|---------|------|---------|---------|
| HighErrorRate | critical | 错误率 > 1% 持续5分钟 | Email + Slack + 钉钉 |
| HighLatency | warning | P95延迟 > 500ms 持续10分钟 | Email + Slack + 钉钉 |
| TrafficAnomaly | warning | 流量异常（增长 > 200%） | Email + Slack + 钉钉 |

### 2. 业务告警（sales-team）

| 告警名称 | 级别 | 触发条件 | 通知渠道 |
|---------|------|---------|---------|
| QuotaNearLimit | warning | 配额使用率 > 80% 持续10分钟 | Email + 钉钉 |
| QuotaExceeded | critical | 配额已耗尽 | Email + 钉钉 |
| SubscriptionExpiring | warning | 订阅7天内到期 | Email + 钉钉 |
| TenantSuspended | critical | 租户被暂停 | Email + 钉钉 |

### 3. 数据库告警（dba-team）

| 告警名称 | 级别 | 触发条件 | 通知渠道 |
|---------|------|---------|---------|
| DBConnectionPoolHigh | warning | 连接池使用率 > 80% 持续5分钟 | Email + Slack |
| DBQuerySlow | warning | 查询延迟 > 1s 持续10分钟 | Email + Slack |
| DBConnectionHigh | critical | 连接数 > 400 | Email + Slack |

### 4. 缓存告警（backend-team）

| 告警名称 | 级别 | 触发条件 | 通知渠道 |
|---------|------|---------|---------|
| CacheHitRateLow | warning | 命中率 < 70% 持续15分钟 | Email + Slack |
| CacheSizeHigh | warning | 缓存大小 > 5GB | Email + Slack |

### 5. 系统告警（ops-team）

| 告警名称 | 级别 | 触发条件 | 通知渠道 |
|---------|------|---------|---------|
| HighGoroutines | warning | Goroutine数 > 5000 持续10分钟 | Email + Slack + 钉钉 |
| HighMemoryUsage | critical | 堆内存 > 4GB 持续5分钟 | Email + Slack + 钉钉 |
| HighGCFrequency | warning | GC频率 > 1次/秒 持续10分钟 | Email + Slack + 钉钉 |

---

## 告警处理流程

### 1. 接收告警

当收到告警通知时，按以下步骤处理：

1. **查看告警详情**
   - 点击通知中的链接访问 Grafana
   - 或访问 Alertmanager UI: http://localhost:9093

2. **评估严重性**
   - **critical**: 立即处理（15分钟内响应）
   - **warning**: 计划处理（1小时内响应）

3. **确认影响范围**
   - 检查受影响的租户/服务
   - 评估业务影响

### 2. 诊断问题

1. **收集信息**
   ```bash
   # 查看相关日志
   tail -f /var/log/zker/api.log | grep ERROR

   # 查看指标趋势
   # 在 Grafana 中查看相关面板
   ```

2. **确认根因**
   - 查看相关日志和监控指标
   - 检查最近的变更（部署、配置等）
   - 查看是否有相关历史问题

3. **记录诊断结果**
   - 在工单系统中记录问题
   - 更新告警状态（正在处理）

### 3. 解决问题

1. **实施修复**
   - 根据诊断结果执行修复方案
   - 参考运维手册中的"常见告警处理"

2. **验证修复**
   - 确认告警已恢复
   - 验证服务正常
   - 检查相关指标

3. **更新状态**
   - 标记告警为已解决
   - 更新工单状态

### 4. 事后分析

1. **编写事故报告**
   - 描述问题影响
   - 分析根本原因
   - 提出改进建议

2. **实施改进**
   - 更新监控规则
   - 改进代码/配置
   - 更新运维文档

---

## 常见告警处理

### HighErrorRate（API错误率过高）

**症状**: API错误率 > 1% 持续5分钟

**诊断步骤**:
1. 查看Grafana "Error Rate" 面板，确认哪个端点错误率高
2. 查看API日志，搜索错误信息
3. 检查最近是否有部署

**可能原因**:
- 代码bug导致错误
- 数据库连接问题
- 外部服务不可用
- 配置错误

**解决方案**:
```bash
# 1. 如果是代码bug，立即回滚
make rollback

# 2. 如果是数据库问题，检查连接
docker exec -it mysql mysql -uroot -p
SHOW PROCESSLIST;
KILL <process_id>;

# 3. 如果是外部服务问题，检查服务状态
curl -v https://external-service.com/health
```

**预防措施**:
- 加强代码审查
- 增加自动化测试覆盖
- 设置数据库连接池监控

---

### HighLatency（API延迟过高）

**症状**: P95延迟 > 500ms 持续10分钟

**诊断步骤**:
1. 查看Grafana "P95 Latency" 面板
2. 查看Jaeger追踪，定位慢查询
3. 检查数据库慢查询日志

**可能原因**:
- 慢SQL查询
- N+1查询问题
- 缓存未命中
- 外部服务调用慢

**解决方案**:
```bash
# 1. 优化慢查询
# 添加索引
ALTER TABLE bots ADD INDEX idx_tenant_created (tenant_id, created_at);

# 2. 启用查询缓存
# 检查缓存配置
redis-cli INFO stats
# hits: misses 应 > 3:1

# 3. 优化N+1查询
# 使用 GORM Preload
db.Preload("Config").Preload("Creator").Find(&bots)
```

**预防措施**:
- 定期Review慢查询日志
- 使用数据库索引优化
- 实施查询缓存策略

---

### QuotaNearLimit（配额接近限制）

**症状**: 租户配额使用率 > 80%

**诊断步骤**:
1. 查看受影响的租户
2. 检查配额使用趋势
3. 确认是否有异常消耗

**解决方案**:
```bash
# 1. 通知租户配额即将耗尽
# 发送邮件提醒

# 2. 如果是合理增长，建议租户升级套餐
# 创建工单给销售团队

# 3. 如果是异常消耗，检查是否有滥用
# 查看API调用日志
grep "tenant_id:XXX" /var/log/zker/api.log | tail -100
```

**预防措施**:
- 实施配额预警（80%时通知）
- 提供配额使用监控面板
- 优化资源使用效率

---

### QuotaExceeded（配额已耗尽）

**症状**: 租户配额已完全耗尽

**紧急处理**:
```bash
# 1. 确认影响范围
# 查看哪些请求被拒绝

# 2. 紧急情况可临时提高配额
UPDATE quotas SET limit = limit * 1.2 WHERE tenant_id = 'XXX';

# 3. 通知相关团队
# - 销售团队：联系租户
# - 技术团队：监控服务状态
```

**长期解决**:
- 引导租户升级套餐
- 优化资源使用效率
- 实施资源配额弹性扩展

---

### DBConnectionPoolHigh（数据库连接池使用率高）

**症状**: 连接池使用率 > 80% 持续5分钟

**诊断步骤**:
```bash
# 1. 检查当前连接数
docker exec -it mysql mysql -uroot -p -e "SHOW PROCESSLIST;"

# 2. 查看慢查询
docker exec -it mysql mysql -uroot -p -e "SHOW FULL PROCESSLIST;" | grep "Query time"

# 3. 检查连接池配置
cat backend/conf/config.toml | grep -A 5 mysql
```

**可能原因**:
- 慢查询导致连接堆积
- 连接未正确释放
- 并发请求过多

**解决方案**:
```sql
-- 1. 杀掉长时间运行的查询
KILL <process_id>;

-- 2. 增加连接池大小
-- 修改 backend/conf/config.toml
[mysql]
max_open_conns = 200  # 从100增加到200

-- 3. 优化慢查询
-- 添加索引或优化SQL
```

**预防措施**:
- 设置慢查询阈值（>1s）
- 定期Review慢查询日志
- 实施查询超时机制

---

### CacheHitRateLow（缓存命中率低）

**症状**: 缓存命中率 < 70% 持续15分钟

**诊断步骤**:
```bash
# 1. 检查Redis缓存统计
redis-cli INFO stats
# hits: 1000
# misses: 500
# 命中率 = 1000 / (1000 + 500) = 66.7%

# 2. 检查缓存大小
redis-cli INFO memory
# used_memory: 6GB
```

**可能原因**:
- 缓存容量不足，导致淘汰
- 缓存策略不当
- 冷数据过多

**解决方案**:
```bash
# 1. 增加Redis内存
# 修改 docker-compose.yml
services:
  redis:
    command: redis-server --maxmemory 8gb --maxmemory-policy allkeys-lru

# 2. 优化缓存策略
# - 使用 LRU 淘汰策略
# - 增加缓存预热
# - 延长热点数据TTL

# 3. 分析缓存访问模式
redis-cli --bigkeys
# 找出大key，优化存储结构
```

**预防措施**:
- 监控缓存命中率趋势
- 定期分析缓存访问模式
- 实施缓存预热策略

---

### HighMemoryUsage（内存使用过高）

**症状**: 堆内存 > 4GB 持续5分钟

**诊断步骤**:
```bash
# 1. 获取heap profile
curl http://localhost:8080/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# 2. 查看内存分配TOP10
go tool pprof -top heap.prof

# 3. 查看Goroutine数量
curl http://localhost:8080/debug/pprof/goroutine?debug=1
```

**可能原因**:
- 内存泄漏
- Goroutine泄漏
- 大对象未释放

**解决方案**:
```go
// 1. 检查Goroutine泄漏
// 添加监控
go func() {
    for {
        time.Sleep(1 * time.Minute)
        log.Printf("Goroutine count: %d", runtime.NumGoroutine())
    }
}()

// 2. 使用sync.Pool重用对象
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

// 3. 及时释放大对象
buffer.Reset()
bufferPool.Put(buffer)
```

**预防措施**:
- 定期进行内存profile
- 使用pprof监控内存
- 实施内存泄漏检测

---

## 配置修改

### 修改告警阈值

```bash
# 1. 编辑告警规则
vim infra/monitoring/prometheus/alerts.yml

# 2. 示例：修改API错误率阈值
# - alert: HighErrorRate
#   expr: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.02  # 从0.01改为0.02

# 3. 重新加载配置
curl -X POST http://localhost:9090/-/reload

# 4. 验证配置
curl http://localhost:9090/api/v1/status/config
```

### 修改通知渠道

```bash
# 1. 编辑Alertmanager配置
vim infra/monitoring/alertmanager/alertmanager.yml

# 2. 示例：添加新的Slack频道
# receivers:
#   - name: 'api-team'
#     slack_configs:
#       - api_url: '${SLACK_WEBHOOK_URL}'
#         channel: '#api-alerts-new'

# 3. 重新加载配置
curl -X POST http://localhost:9093/-/reload

# 4. 验证配置
curl http://localhost:9093/api/v1/status/config
```

### 添加新的告警规则

```yaml
# 1. 在 prometheus/alerts.yml 中添加
- alert: CustomAlert
  expr: custom_metric > threshold
  for: 5m
  labels:
    severity: warning
    team: backend-team
  annotations:
    summary: "Custom alert triggered"
    description: "Custom metric {{ $labels.instance }} is above threshold"

# 2. 重新加载Prometheus
curl -X POST http://localhost:9090/-/reload

# 3. 验证规则
curl http://localhost:9090/api/v1/rules | jq
```

---

## 测试验证

### 测试告警规则

```bash
# 1. 触发测试告警
# 使用 Prometheus 表达式浏览器
# 访问: http://localhost:9090/graph
# 输入: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])

# 2. 查看当前告警
curl http://localhost:9090/api/v1/alerts | jq

# 3. 查看规则评估结果
curl http://localhost:9090/api/v1/rules | jq '.data.groups[].rules[] | select(.name=="HighErrorRate")'
```

### 测试通知渠道

```bash
# 1. 手动触发测试告警
cat <<EOF | curl -X POST http://localhost:9093/api/v1/alerts -H 'Content-Type: application/json' -d @-
[
  {
    "labels": {
      "alertname": "TestAlert",
      "severity": "warning",
      "instance": "localhost"
    },
    "annotations": {
      "description": "This is a test alert"
    }
  }
]
EOF

# 2. 检查Alertmanager日志
docker logs alertmanager | tail -50

# 3. 验证通知是否发送
# - 检查Email
# - 检查Slack频道
# - 检查钉钉群
```

---

## 故障排查

### 告警未触发

**症状**: 预期会触发的告警未触发

**检查步骤**:
```bash
# 1. 确认Prometheus正在抓取指标
curl http://localhost:9090/api/v1/targets | jq '.data.activeTargets[] | select(.labels.job=="zker-api")'

# 2. 检查告警规则是否加载
curl http://localhost:9090/api/v1/rules | jq '.data.groups[].rules[] | select(.name=="HighErrorRate")'

# 3. 检查规则评估时间
# 访问 http://localhost:9090/graph
# 查看规则评估结果

# 4. 验证阈值设置
# 确认告警表达式的阈值是否合理
```

### 告警未发送通知

**症状**: 告警已触发但未收到通知

**检查步骤**:
```bash
# 1. 检查Alertmanager状态
curl http://localhost:9093/api/v1/status | jq

# 2. 检查告警是否到达Alertmanager
curl http://localhost:9093/api/v2/alerts | jq

# 3. 检查通知配置
cat infra/monitoring/alertmanager/alertmanager.yml | grep -A 10 "receivers:"

# 4. 检查Alertmanager日志
docker logs alertmanager | grep -i error

# 5. 验证环境变量
env | grep DINGTALK_WEBHOOK_URL
env | grep SLACK_WEBHOOK_URL
```

### 通知发送失败

**症状**: 通知发送报错

**常见原因**:
- Webhook URL错误或失效
- 钉钉/Slack token过期
- 邮件服务器配置错误
- 网络问题

**解决方案**:
```bash
# 1. 验证Webhook URL
curl -X POST ${DINGTALK_WEBHOOK_URL} \
  -H 'Content-Type: application/json' \
  -d '{"msgtype":"text","text":{"content":"test"}}'

# 2. 重新配置Alertmanager
vim infra/monitoring/alertmanager/alertmanager.yml
# 更新正确的URL和token

# 3. 重新加载配置
docker restart alertmanager

# 4. 检查防火墙规则
sudo iptables -L -n | grep 9093
```

---

## 最佳实践

### 1. 告警分级

- **Critical**: 立即处理（15分钟内）
- **Warning**: 计划处理（1小时内）
- **Info**: 信息提示（可忽略）

### 2. 告警静默

```bash
# 在Alertmanager UI中创建静默规则
# 访问: http://localhost:9093
# 点击: Silence -> New Silence

# 或使用API
cat <<EOF | curl -X POST http://localhost:9093/api/v2/silences -H 'Content-Type: application/json' -d @-
{
  "matchers": [
    {
      "name": "alertname",
      "value": "HighLatency",
      "isRegex": false
    }
  ],
  "startsAt": "2025-01-01T00:00:00Z",
  "endsAt": "2025-01-01T01:00:00Z",
  "createdBy": "ops-team",
  "comment": "Scheduled maintenance"
}
EOF
```

### 3. 告警聚合

- 利用 `group_by` 减少重复通知
- 合理设置 `group_wait` 和 `group_interval`
- 使用抑制规则避免告警风暴

### 4. 告警优化

- 定期Review告警规则
- 删除无效或误报的告警
- 调整阈值减少噪音
- 添加必要的上下文信息

---

## 相关文档

- [Prometheus告警规则配置](../prometheus/alerts.yml)
- [Grafana监控大盘](../grafana/dashboards/api-dashboard.json)
- [ZKER-故障排查手册](../../../../../docs/企业级功能完善与统一性设计方案/ZKER-故障排查手册_v1.0.md)

---

**最后更新**: 2025-01-01
**维护人**: 研发B - 后端工程师
