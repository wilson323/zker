# ZKER 生产环境压力测试执行指南

**版本**: v1.0.0
**最后更新**: 2025-01-01
**适用环境**: 生产环境、预发布环境

---

## 📋 目录

- [测试概述](#测试概述)
- [执行前准备](#执行前准备)
- [测试环境要求](#测试环境要求)
- [测试场景配置](#测试场景配置)
- [执行步骤](#执行步骤)
- [结果分析](#结果分析)
- [优化建议](#优化建议)
- [回滚方案](#回滚方案)
- [应急预案](#应急预案)

---

## 🎯 测试概述

### 测试目的

**生产环境压力测试** 的目的是：

1. **验证系统容量**：确认系统是否能支撑预期峰值流量
2. **发现性能瓶颈**：识别系统在高负载下的弱点和限制
3. **验证弹性能力**：测试自动扩缩容、限流、降级等机制
4. **验证监控告警**：确认监控和告警系统在压力下是否正常工作
5. **评估稳定性**：验证系统在持续高压下的稳定性

### 测试类型

| 测试类型 | 目标 | 并发用户 | 持续时间 | 执行时机 |
|---------|------|---------|---------|---------|
| **负载测试** | 验证正常负载下的性能 | 100-200 | 10-15分钟 | 上线前、重大变更后 |
| **压力测试** | 找到系统性能极限 | 200-500 | 20-30分钟 | 上线前、容量评估 |
| **浸泡测试** | 验证长时间运行的稳定性 | 100-150 | 2-8小时 | 上线前、重大变更后 |
| **峰值测试** | 验证应对突发流量的能力 | 1000-2000 | 5-10分钟 | 上线前、营销活动前 |

### 风险评估

**⚠️ 高风险操作**：生产环境压力测试可能影响真实用户

**风险等级**: 🔴 **高风险**

**主要风险**：
- 服务响应变慢，影响用户体验
- 数据库连接池耗尽，导致服务不可用
- Redis缓存击穿，导致数据库压力过大
- 第三方API调用超限，产生额外费用
- 资源耗尽导致服务器宕机

**缓解措施**：
- ✅ **选择低峰时段**（凌晨2-6点）
- ✅ **使用隔离的测试环境**（预发布环境）
- ✅ **准备快速回滚方案**（5分钟内可回滚）
- ✅ **配置合理的限流规则**（保护关键服务）
- ✅ **实时监控关键指标**（随时可以中断）
- ✅ **通知所有相关人员**（开发、运维、测试）

---

## 📦 执行前准备

### 1. 环境检查清单

**生产环境检查**：
- [ ] 确认当前时间段是低峰期（凌晨2-6点）
- [ ] 确认有足够的监控资源（磁盘空间、日志收集）
- [ ] 确认所有告警规则已启用且通知渠道畅通
- [ ] 确认有运维人员在线待命
- [ ] 确认已通知所有相关人员（测试开始/结束通知）

**预发布环境检查**（推荐）：
- [ ] 环境数据已从生产环境同步（脱敏后）
- [ ] 所有服务已启动并运行正常
- [ ] 监控和告警系统已配置
- [ ] 数据库和缓存性能与生产环境一致

**测试工具准备**：
- [ ] K6已安装并配置（`k6 version`）
- [ ] 测试脚本已准备并review通过
- [ ] 测试数据已准备（测试账号、测试租户）
- [ ] 测试结果存储位置已确定

### 2. 监控指标准备

**关键监控指标**：

| 指标类型 | 指标名称 | 告警阈值 | 监控工具 |
|---------|---------|---------|---------|
| **系统资源** | CPU使用率 | > 80% 持续5分钟 | Prometheus |
| **系统资源** | 内存使用率 | > 85% 持续5分钟 | Prometheus |
| **API性能** | P95响应时间 | > 2秒 | Prometheus |
| **API性能** | 错误率 | > 5% | Prometheus |
| **数据库** | 连接池使用率 | > 80% | Prometheus |
| **数据库** | 慢查询数量 | > 10/分钟 | Prometheus |
| **缓存** | 缓存命中率 | < 90% | Prometheus |
| **业务** | Bot调用失败率 | > 10% | Grafana |

**监控大盘准备**：
- [ ] Grafana监控大盘已打开并全屏显示
- [ ] 关键指标面板已调整到合理的时间范围（实时）
- [ ] 告警通知渠道已测试（企业微信/钉钉/短信）

### 3. 回滚准备

**回滚方案检查**：
- [ ] 已记录当前Git版本号（`git rev-parse HEAD`）
- [ ] 已备份当前数据库（`mysqldump`）
- [ ] 已备份当前配置文件（`/etc/coze-studio/`）
- [ ] 回滚脚本已准备并测试（`scripts/rollback.sh`）
- [ ] 回滚执行人已确定（运维负责人）
- [ ] 回滚决策标准已明确（何时触发回滚）

**回滚触发条件**：
- ❌ 错误率持续 > 10% 超过5分钟
- ❌ P95响应时间持续 > 10秒 超过5分钟
- ❌ 系统资源使用率持续 > 95% 超过5分钟
- ❌ 任何服务完全不可用（502/503错误）
- ❌ 数据库连接池耗尽

---

## 🏗️ 测试环境要求

### 生产环境要求

**最低配置**（适用于小型部署）：
```
CPU: 8核心
内存: 16GB
磁盘: 500GB SSD
网络: 1 Gbps
数据库: MySQL 8.4.5, 200连接
缓存: Redis 8.0, 4GB内存
```

**推荐配置**（适用于中大型部署）：
```
CPU: 16核心
内存: 32GB
磁盘: 1TB SSD
网络: 10 Gbps
数据库: MySQL 8.4.5, 500连接
缓存: Redis 8.0, 16GB内存
负载均衡: Nginx/HAProxy
```

### 预发布环境要求

**与生产环境一致**：
- 硬件配置一致（或按比例缩小）
- 软件版本一致（Go、MySQL、Redis等）
- 数据量一致（使用脱敏数据）
- 配置参数一致（连接池、超时等）
- 监控告警一致

**差异点**：
- 使用测试域名（如 `test-api.coze-studio.com`）
- 数据已脱敏（用户名、邮箱等）
- 不发送真实的外部通知（邮件、短信）

---

## 🧪 测试场景配置

### 场景1：负载测试（Load Test）

**目标**：验证系统在正常负载下的性能

**配置**：
```javascript
// tests/performance/scenarios/load_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '2m', target: 50 },   // 预热：2分钟爬坡到50用户
    { duration: '5m', target: 100 },  // 正常负载：100用户持续5分钟
    { duration: '2m', target: 0 },    // 降温：2分钟降到0
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],   // P95响应时间 < 500ms
    http_req_failed: ['rate<0.01'],     // 错误率 < 1%
  },
};

const BASE_URL = __ENV.API_URL || 'https://api.coze-studio.com';

export default function () {
  // 1. 登录
  const loginResp = http.post(`${BASE_URL}/api/passport/login`, JSON.stringify({
    email: 'test@example.com',
    password: 'test123',
  }), {
    headers: { 'Content-Type': 'application/json' },
  });

  check(loginResp, {
    'login successful': (r) => r.status === 200,
  });

  const token = loginResp.json('data.token');

  // 2. 获取Bot列表
  const listResp = http.get(`${BASE_URL}/api/bots`, {
    headers: {
      'Authorization': `Bearer ${token}`,
      'X-Tenant-ID': 'test-tenant-001',
    },
  });

  check(listResp, {
    'list bots successful': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });

  sleep(1);
}
```

**性能基线**：
- P95响应时间: < 500ms
- P99响应时间: < 1000ms
- 错误率: < 1%
- CPU使用率: < 60%
- 内存使用率: < 70%

**执行命令**：
```bash
# 设置环境变量
export API_URL="https://api.coze-studio.com"

# 执行测试
k6 run tests/performance/scenarios/load_test.js

# 查看实时结果
k6 run --out json=load_test_results.json tests/performance/scenarios/load_test.js
```

### 场景2：压力测试（Stress Test）

**目标**：找到系统性能瓶颈和极限

**配置**：
```javascript
// tests/performance/scenarios/stress_test.js
export const options = {
  stages: [
    { duration: '2m', target: 100 },   // 预热
    { duration: '5m', target: 200 },   // 逐步加压到200用户
    { duration: '5m', target: 300 },   // 继续加压到300用户
    { duration: '5m', target: 400 },   // 继续加压到400用户（寻找瓶颈）
    { duration: '2m', target: 0 },     // 快速降压
  ],
  thresholds: {
    http_req_duration: ['p(95)<3000'],  // P95 < 3秒（宽松）
    http_req_failed: ['rate<0.10'],     // 错误率 < 10%（宽松）
  },
};
```

**性能基线**（压力测试允许降级）：
- P95响应时间: < 3秒
- 错误率: < 10%
- CPU使用率: < 90%
- 内存使用率: < 95%

**观察指标**：
- 系统何时开始出现错误（错误率开始上升）
- 响应时间何时急剧上升（指数增长点）
- CPU/内存使用率峰值
- 是否有服务崩溃或重启
- 数据库连接池是否耗尽

**执行命令**：
```bash
k6 run tests/performance/scenarios/stress_test.js
```

### 场景3：浸泡测试（Soak Test）

**目标**：验证系统长时间运行的稳定性

**配置**：
```javascript
// tests/performance/scenarios/soak_test.js
export const options = {
  stages: [
    { duration: '10m', target: 50 },   // 预热
    { duration: '4h', target: 100 },   // 持续运行4小时
    { duration: '10m', target: 0 },    // 降温
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'],  // P95 < 1秒
    http_req_failed: ['rate<0.01'],     // 错误率 < 1%
  },
};
```

**观察指标**：
- **内存泄漏**：内存使用是否持续上升
- **连接泄漏**：数据库连接数是否持续上升
- **性能退化**：响应时间是否逐渐变慢
- **GC压力**：Go的GC频率是否增加
- **磁盘空间**：日志文件是否占满磁盘

**执行命令**：
```bash
# 后台运行，并将输出重定向到文件
nohup k6 run tests/performance/scenarios/soak_test.js > soak_test.log 2>&1 &

# 监控进程
ps aux | grep k6

# 实时查看日志
tail -f soak_test.log
```

### 场景4：峰值测试（Spike Test）

**目标**：验证系统应对突发流量的能力

**配置**：
```javascript
// tests/performance/scenarios/spike_test.js
export const options = {
  stages: [
    { duration: '1m', target: 100 },   // 正常负载
    { duration: '30s', target: 1000 },  // 突然增加到1000用户（峰值）
    { duration: '2m', target: 1000 },   // 保持峰值2分钟
    { duration: '1m', target: 100 },    // 快速降压
  ],
  thresholds: {
    http_req_duration: ['p(95)<5000'],  // P95 < 5秒（非常宽松）
    http_req_failed: ['rate<0.20'],     // 错误率 < 20%（非常宽松）
  },
};
```

**观察指标**：
- 系统是否崩溃（502/503错误）
- 自动扩缩容是否生效（K8s HPA）
- 限流机制是否触发（返回429）
- 服务降级是否生效（关闭非核心功能）
- 是否有请求堆积（队列溢出）

**执行命令**：
```bash
k6 run tests/performance/scenarios/spike_test.js
```

---

## 🚀 执行步骤

### 步骤1：执行前检查（T-30分钟）

**T-30分钟**：环境准备
```bash
# 1. 检查当前Git版本
git rev-parse HEAD > /tmp/before_test_git_hash.txt
cat /tmp/before_test_git_hash.txt

# 2. 备份数据库
mysqldump -u root -p coze_studio > /tmp/before_test_db_backup_$(date +%Y%m%d_%H%M%S).sql

# 3. 备份配置文件
tar -czf /tmp/before_test_config_$(date +%Y%m%d_%H%M%S).tar.gz /etc/coze-studio/

# 4. 检查服务状态
systemctl status coze-studio-api
systemctl status coze-studio-worker
```

**T-20分钟**：监控准备
```bash
# 1. 打开Grafana监控大盘
# 浏览器访问: http://grafana.coze-studio.com

# 2. 检查告警通知渠道
# 发送测试告警，确认企业微信/钉钉正常

# 3. 设置测试期间的告警阈值（临时放宽）
# 编辑 alerts.yml，降低阈值，避免频繁告警
```

**T-10分钟**：最后确认
```bash
# 1. 确认所有服务正常运行
curl http://api.coze-studio.com/api/health/status

# 2. 检查数据库连接数
mysql -u root -p -e "SHOW PROCESSLIST" | wc -l

# 3. 检查Redis连接数
redis-cli INFO clients

# 4. 通知所有相关人员
echo "压力测试将在10分钟后开始，请密切关注监控大盘" | mail -s "压力测试通知" ops-team@coze-studio.com
```

### 步骤2：执行测试（T+0）

**T+0分钟**：开始负载测试
```bash
# 记录开始时间
echo "Test started at $(date)" > /tmp/test_start_time.txt

# 执行负载测试
k6 run tests/performance/scenarios/load_test.js | tee load_test_output.log
```

**T+10分钟**：检查负载测试结果
```bash
# 查看测试摘要
grep "http_req_duration" load_test_output.log
grep "http_req_failed" load_test_output.log

# 检查是否有错误
if grep -q "status=500" load_test_output.log; then
  echo "⚠️ 发现500错误，请检查日志"
  tail -100 /var/log/coze-studio/api.log
fi
```

**T+15分钟**：开始压力测试（如果负载测试通过）
```bash
k6 run tests/performance/scenarios/stress_test.js | tee stress_test_output.log
```

**T+40分钟**：检查压力测试结果
```bash
# 查看系统资源使用情况
top -b -n 1 | head -20

# 查看数据库慢查询
mysql -u root -p -e "SELECT * FROM information_schema.PROCESSLIST WHERE TIME > 5"

# 检查Redis性能
redis-cli INFO stats
```

**T+45分钟**：决定是否继续
- ✅ 如果系统稳定，继续执行峰值测试
- ❌ 如果发现问题，停止测试，分析原因

### 步骤3：峰值测试（T+50分钟）

```bash
k6 run tests/performance/scenarios/spike_test.js | tee spike_test_output.log
```

**观察重点**：
- 系统是否崩溃
- 限流是否生效
- 自动扩容是否触发

### 步骤4：执行后清理（T+60分钟）

```bash
# 1. 停止所有K6进程
pkill -9 k6

# 2. 记录结束时间
echo "Test ended at $(date)" > /tmp/test_end_time.txt

# 3. 收集测试结果
mkdir -p /tmp/performance_test_results_$(date +%Y%m%d_%H%M%S)
cp *.log /tmp/performance_test_results_$(date +%Y%m%d_%H%M%S)/
cp *_output.log /tmp/performance_test_results_$(date +%Y%m%d_%H%M%S)/

# 4. 检查系统状态
curl http://api.coze-studio.com/api/health/status

# 5. 恢复告警阈值
git checkout alerts.yml
systemctl reload prometheus

# 6. 通知测试完成
echo "压力测试已完成，结果已保存到 /tmp/performance_test_results_*" | mail -s "压力测试完成" ops-team@coze-studio.com
```

---

## 📊 结果分析

### 1. 性能指标分析

**响应时间分析**：
```
指标                    基线值    实测值    状态
-----------------------------------------------
P50 响应时间            100ms     120ms     ✅ 通过
P95 响应时间            500ms     800ms     ⚠️ 警告
P99 响应时间            1000ms    2500ms    ❌ 失败
```

**错误率分析**：
```
错误类型                基线值    实测值    状态
-----------------------------------------------
4xx 错误率              0.1%      0.5%      ⚠️ 警告
5xx 错误率              0.05%     0.2%      ❌ 失败
总体错误率              0.15%     0.7%      ❌ 失败
```

**资源使用分析**：
```
资源                    基线值    实测值    峰值     状态
-----------------------------------------------
CPU 使用率              60%       75%       95%     ⚠️ 警告
内存使用率              70%       85%       98%     ❌ 失败
数据库连接数            100/200   180/200   200/200 ❌ 失败
Redis 命中率            95%       88%       75%     ⚠️ 警告
```

### 2. 瓶颈识别

**常见瓶颈及解决方案**：

| 瓶颈类型 | 症状 | 根本原因 | 解决方案 |
|---------|------|---------|---------|
| **CPU瓶颈** | CPU使用率 > 90% | 算法复杂度高、无缓存 | 优化算法、增加缓存 |
| **内存瓶颈** | 内存使用率 > 95% | 内存泄漏、缓存过大 | 修复内存泄漏、限制缓存大小 |
| **数据库瓶颈** | 慢查询增多 | 缺少索引、N+1查询 | 添加索引、优化SQL |
| **网络瓶颈** | 响应时间高、吞吐量低 | 带宽不足、延迟高 | 升级带宽、使用CDN |
| **锁竞争** | 吞吐量上不去 | 数据库锁、互斥锁 | 优化锁粒度、使用乐观锁 |

### 3. 趋势分析

**性能退化趋势**：
```
时间点    P95响应时间   错误率   CPU使用率
---------------------------------------
0分钟     300ms        0.1%     50%
10分钟    450ms        0.3%     65%
20分钟    800ms        1.2%     85%
30分钟    1500ms       3.5%     95%  ❌ 性能退化
```

**分析**：
- 前10分钟性能正常
- 10-20分钟性能开始下降（可能是连接池接近上限）
- 20-30分钟性能急剧下降（连接池耗尽，请求开始堆积）

---

## 🔧 优化建议

### 1. 数据库优化

**问题**：慢查询增多，P95响应时间 > 1秒

**解决方案**：
```sql
-- 1. 添加复合索引
CREATE INDEX idx_tenant_status_created
ON bots(tenant_id, status, created_at);

-- 2. 优化N+1查询
-- 优化前：N+1查询
SELECT * FROM bots WHERE tenant_id = 'xxx';
-- 循环N次查询
SELECT * FROM bot_configs WHERE bot_id = '...';

-- 优化后：使用JOIN
SELECT b.*, bc.*
FROM bots b
LEFT JOIN bot_configs bc ON b.bot_id = bc.bot_id
WHERE b.tenant_id = 'xxx';

-- 3. 增加连接池大小
SET GLOBAL max_connections = 500;
```

**效果预期**：
- 查询时间减少 50-80%
- P95响应时间降低到 < 500ms

### 2. 缓存优化

**问题**：Redis命中率 < 90%，数据库压力大

**解决方案**：
```go
// 1. 增加缓存时间
// 优化前：TTL = 1分钟
redis.Set("tenant:"+tenantID, data, 1*time.Minute)

// 优化后：TTL = 5分钟
redis.Set("tenant:"+tenantID, data, 5*time.Minute)

// 2. 使用缓存预热
func WarmupCache(tenantID string) {
    // 预加载常用数据到缓存
    tenant := tenantService.GetTenant(tenantID)
    redis.Set("tenant:"+tenantID, tenant, 5*time.Minute)

    quotas := quotaService.GetQuotas(tenantID)
    redis.Set("quotas:"+tenantID, quotas, 5*time.Minute)
}

// 3. 使用缓存穿透保护
func GetTenantWithCache(tenantID string) (*Tenant, error) {
    // 尝试从缓存获取
    cached, err := redis.Get("tenant:" + tenantID)
    if err == nil {
        return cached, nil
    }

    // 使用布隆过滤器检查是否存在
    if !bloomFilter MightContain(tenantID) {
        return nil, ErrTenantNotFound
    }

    // 从数据库查询
    tenant, err := db.GetTenant(tenantID)
    if err != nil {
        return nil, err
    }

    // 写入缓存
    redis.Set("tenant:"+tenantID, tenant, 5*time.Minute)
    return tenant, nil
}
```

**效果预期**：
- 缓存命中率提升到 > 95%
- 数据库查询减少 80%

### 3. 代码优化

**问题**：CPU使用率高，响应慢

**解决方案**：
```go
// 1. 减少序列化开销
// 优化前：多次序列化
for _, item := range items {
    json.Marshal(item)  // 每次都序列化
}

// 优化后：批量序列化
jsonData, _ := json.Marshal(items)  // 只序列化一次

// 2. 使用连接池
// 优化前：每次都创建新连接
func callExternalAPI(url string) {
    client := &http.Client{Timeout: 10 * time.Second}
    client.Get(url)
}

// 优化后：复用连接池
var httpClient = &http.Client{
    Timeout: 10 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

func callExternalAPI(url string) {
    httpClient.Get(url)
}

// 3. 使用goroutine池
// 优化前：无限制创建goroutine
for _, id := range botIDs {
    go processBot(id)  // 可能创建数万个goroutine
}

// 优化后：限制并发数
sem := make(chan struct{}, 100)  // 最多100个并发
var wg sync.WaitGroup

for _, id := range botIDs {
    wg.Add(1)
    sem <- struct{}{}
    go func(botID string) {
        defer wg.Done()
        defer func() { <-sem }()
        processBot(botID)
    }(id)
}
wg.Wait()
```

**效果预期**：
- CPU使用率降低 30-50%
- 响应时间降低 40-60%

---

## 🔄 回滚方案

### 回滚决策标准

**立即回滚**（触发任何一项）：
- ❌ 错误率持续 > 10% 超过5分钟
- ❌ P95响应时间持续 > 10秒 超过5分钟
- ❌ 系统资源使用率持续 > 95% 超过5分钟
- ❌ 任何服务完全不可用（连续502/503错误）
- ❌ 数据库连接池耗尽
- ❌ 监控系统本身不可用

**考虑回滚**（触发多项）：
- ⚠️ 错误率持续 > 5% 超过10分钟
- ⚠️ P95响应时间持续 > 5秒 超过10分钟
- ⚠️ CPU使用率持续 > 85% 超过15分钟
- ⚠️ 用户投诉明显增加

### 回滚执行步骤

**步骤1**：停止测试（T+0）
```bash
# 立即停止所有K6进程
pkill -9 k9

# 确认已停止
ps aux | grep k6
```

**步骤2**：恢复服务（T+1分钟）
```bash
# 如果服务已崩溃，重启服务
systemctl restart coze-studio-api
systemctl restart coze-studio-worker

# 检查服务状态
systemctl status coze-studio-api
```

**步骤3**：恢复数据库（如果需要，T+5分钟）
```bash
# 停止应用服务
systemctl stop coze-studio-api

# 恢复数据库备份
mysql -u root -p coze_studio < /tmp/before_test_db_backup_XXXXX.sql

# 启动应用服务
systemctl start coze-studio-api
```

**步骤4**：恢复配置（如果需要，T+3分钟）
```bash
# 恢复配置文件
cd /etc
tar -xzf /tmp/before_test_config_XXXXX.tar.gz

# 重启相关服务
systemctl reload nginx
systemctl reload prometheus
```

**步骤5**：验证恢复（T+10分钟）
```bash
# 检查服务健康状态
curl http://api.coze-studio.com/api/health/status

# 检查关键API
curl http://api.coze-studio.com/api/bots

# 检查数据库
mysql -u root -p -e "SELECT COUNT(*) FROM bots"

# 检查Redis
redis-cli PING
```

**步骤6**：通知相关人员（T+15分钟）
```bash
echo "压力测试已回滚，服务已恢复" | mail -s "压力测试回滚通知" ops-team@coze-studio.com
```

### 一键回滚脚本

```bash
#!/bin/bash
# scripts/rollback_performance_test.sh

set -e

echo "🔄 开始回滚压力测试环境..."

# 1. 停止测试
echo "⏹️ 停止测试进程..."
pkill -9 k6 || true

# 2. 恢复Git版本
echo "📦 恢复Git版本..."
BEFORE_HASH=$(cat /tmp/before_test_git_hash.txt)
git reset --hard $BEFORE_HASH

# 3. 重启服务
echo "🔄 重启服务..."
systemctl restart coze-studio-api
systemctl restart coze-studio-worker

# 4. 等待服务启动
echo "⏳ 等待服务启动..."
sleep 30

# 5. 验证服务
echo "✅ 验证服务状态..."
HEALTH_CHECK=$(curl -s http://api.coze-studio.com/api/health/status)
if echo "$HEALTH_CHECK" | grep -q "ok"; then
    echo "✅ 服务恢复成功"
else
    echo "❌ 服务恢复失败，需要人工介入"
    exit 1
fi

# 6. 发送通知
echo "📧 发送回滚通知..."
echo "压力测试已回滚，服务已恢复正常" | mail -s "压力测试回滚完成" ops-team@coze-studio.com

echo "✅ 回滚完成"
```

---

## 🚨 应急预案

### 场景1：服务完全不可用

**症状**：所有API返回502/503，无法访问

**应急处理**：
```bash
# 1. 立即停止测试
pkill -9 k9

# 2. 检查服务状态
systemctl status coze-studio-api

# 3. 如果服务崩溃，重启服务
systemctl restart coze-studio-api

# 4. 检查错误日志
tail -100 /var/log/coze-studio/api.log | grep ERROR

# 5. 如果是数据库问题，检查数据库
mysql -u root -p -e "SHOW PROCESSLIST"

# 6. 必要时切换到备用数据库
# 修改配置文件，指向备用数据库
# systemctl restart coze-studio-api
```

### 场景2：数据库连接池耗尽

**症状**：API报错"Too many connections"

**应急处理**：
```bash
# 1. 立即停止测试
pkill -9 k9

# 2. 查看当前连接数
mysql -u root -p -e "SHOW PROCESSLIST" | wc -l

# 3. 杀掉空闲连接
mysql -u root -p -e "KILL IDLE CONNECTIONS"

# 4. 临时增加最大连接数
mysql -u root -p -e "SET GLOBAL max_connections = 1000"

# 5. 重启应用服务
systemctl restart coze-studio-api

# 6. 监控连接数
watch -n 5 'mysql -u root -p -e "SHOW PROCESSLIST" | wc -l'
```

### 场景3：磁盘空间不足

**症状**：日志文件占满磁盘，服务无法写入

**应急处理**：
```bash
# 1. 立即停止测试
pkill -9 k9

# 2. 检查磁盘使用
df -h

# 3. 清理旧日志
find /var/log/coze-studio/ -name "*.log" -mtime +7 -delete

# 4. 清理K6输出文件
rm -rf /tmp/k6-*.json

# 5. 如果仍然不足，清理Docker镜像
docker system prune -a

# 6. 监控磁盘使用
watch -n 5 df -h
```

### 场景4：Redis缓存失效

**症状**：缓存命中率 < 50%，数据库压力巨大

**应急处理**：
```bash
# 1. 立即停止测试
pkill -9 k9

# 2. 检查Redis状态
redis-cli PING
redis-cli INFO stats

# 3. 如果Redis挂了，重启Redis
systemctl restart redis

# 4. 预热缓存（重新加载常用数据）
curl -X POST http://api.coze-studio.com/api/admin/cache/warmup \
  -H "Authorization: Bearer <admin-token>"

# 5. 监控缓存命中率
redis-cli INFO stats | grep keyspace_hits
```

---

## 📚 相关文档

- [性能基线文档](./ZKER-性能基线文档.md)
- [监控告警规则说明](../../deploy/monitoring/README.md)
- [故障排查手册](./ZKER-故障排查手册_v1.0.md)
- [灰度发布策略](./ZKER-灰度发布策略_v1.0.md)

---

**更新日志**：

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2025-01-01 | v1.0.0 | 初始版本，完整的生产环境压力测试执行指南 | Claude AI |
