# ZKER 智能路由与权限管理系统 - 性能测试计划

**文档版本**: v1.0
**创建日期**: 2025-01-01
**最后更新**: 2025-01-01
**负责人**: 性能测试团队
**审批人**: 技术架构委员会

---

## 📋 文档概述

### 测试目标

本性能测试计划旨在全面验证 ZKER 系统在高并发、大数据量场景下的性能表现，确保系统满足以下核心指标：

**核心性能指标（KPI）**:
- **响应时间**: P95 < 2s, P99 < 5s
- **系统可用性**: ≥ 99.9% (月度)
- **并发处理能力**: 10,000 QPS (Query Per Second)
- **吞吐量**: 50,000 TPS (Transaction Per Second)
- **资源利用率**: CPU < 70%, 内存 < 80%, 磁盘 I/O < 60%

### 测试范围

**测试模块**:
1. **API 性能测试**: REST API 接口性能
2. **数据库性能测试**: SQL 查询、事务处理
3. **缓存性能测试**: Redis 读写性能
4. **消息队列性能测试**: NSQ 吞吐量
5. **向量检索性能测试**: Milvus 相似度搜索
6. **全文检索性能测试**: Elasticsearch 查询
7. **前端性能测试**: 页面加载、交互响应
8. **压力测试**: 系统极限承载能力
9. **稳定性测试**: 长时间运行稳定性

**不测试范围**:
- 第三方服务商 API（如 OpenAI、Azure OpenAI 等）
- 非核心功能模块（如管理后台报表）

---

## 🎯 性能指标详解

### 1. 响应时间指标

| API 类别 | P50 (中位数) | P95 | P99 | P99.9 |
|---------|-------------|-----|-----|-------|
| 用户认证 (登录/注册) | < 200ms | < 500ms | < 1s | < 2s |
| Bot 管理 (列表/详情) | < 100ms | < 300ms | < 500ms | < 1s |
| 对话交互 (单轮) | < 500ms | < 2s | < 5s | < 10s |
| 知识库检索 | < 300ms | < 1s | < 2s | < 5s |
| 工作流执行 | < 1s | < 3s | < 5s | < 10s |
| 权限校验 | < 50ms | < 100ms | < 200ms | < 500ms |
| 路由决策 | < 100ms | < 200ms | < 300ms | < 500ms |

**测量方式**:
```
响应时间 = 应用层处理时间 + 数据库查询时间 + 外部 API 调用时间
不包括: 网络传输时间 (RTT)、客户端渲染时间
```

### 2. 并发处理能力

**QPS 指标** (Query Per Second):

| 场景 | 目标 QPS | 说明 |
|-----|---------|------|
| 静态资源访问 | 20,000 | 图片、CSS、JS 等 |
| API 请求总 QPS | 10,000 | 包含所有 API 端点 |
| 登录/注册 | 2,000 | 高并发认证场景 |
| 对话交互 | 5,000 | 核心业务场景 |
| 知识库检索 | 3,000 | 向量 + 关键词检索 |

**TPS 指标** (Transaction Per Second):

| 事务类型 | 目标 TPS |
|---------|---------|
| 用户注册事务 | 1,000 |
| 对话创建事务 | 2,000 |
| 消息发送事务 | 5,000 |
| 知识库写入事务 | 500 |

### 3. 系统可用性指标

| 可用性等级 | 目标值 | 年度停机时间 | 月度停机时间 |
|-----------|-------|-------------|-------------|
| 生产环境 | 99.9% | < 8.76 小时 | < 43.2 分钟 |
| 数据库 | 99.95% | < 4.38 小时 | < 21.6 分钟 |
| 缓存服务 | 99.99% | < 52.56 分钟 | < 4.32 分钟 |

**可用性计算**:
```
可用性 = (总时间 - 故障时间) / 总时间 × 100%
月度停机时间预算 = 30 天 × 24 小时 × 60 分钟 × (1 - 99.9%) = 43.2 分钟
```

### 4. 资源利用率指标

| 资源类型 | 正常运行 | 高负载告警 | 极限阈值 |
|---------|---------|-----------|---------|
| CPU | < 60% | 70%-80% | > 90% (持续 5 分钟) |
| 内存 | < 70% | 80%-85% | > 90% (持续 5 分钟) |
| 磁盘 I/O | < 50% | 60%-70% | > 80% (持续 10 分钟) |
| 网络带宽 | < 60% | 70%-80% | > 90% (持续 10 分钟) |
| 数据库连接数 | < 70% | 80%-90% | > 95% |

### 5. 数据库性能指标

**查询性能**:
| 查询类型 | 平均响应时间 | P95 | 每秒查询数 (QPS) |
|---------|-------------|-----|----------------|
| 单行查询 (主键) | < 1ms | < 5ms | 50,000 |
| 索引范围查询 | < 10ms | < 50ms | 10,000 |
| 复杂 JOIN 查询 | < 100ms | < 300ms | 1,000 |
| 全文检索查询 | < 200ms | < 500ms | 500 |

**连接池配置**:
```
最大连接数 = 200
活跃连接数 = 150
空闲连接数 = 50
连接等待超时 = 30s
连接最大生命周期 = 1 小时
```

### 6. 缓存性能指标

| 缓存操作 | 目标 QPS | 命中率 | 平均延迟 |
|---------|---------|-------|---------|
| Redis GET | 100,000 | > 90% | < 1ms |
| Redis SET | 50,000 | N/A | < 2ms |
| Redis MGET | 20,000 | > 95% | < 5ms (批量) |
| 本地缓存 | 200,000 | > 80% | < 0.1ms |

---

## 🧪 测试类型与策略

### 1. 负载测试 (Load Testing)

**目标**: 验证系统在预期负载下的性能表现

**测试场景**:

**场景 1: 渐进式负载增长**
```
起始负载: 100 并发用户
增长策略: 每分钟增加 100 并发
目标负载: 5,000 并发用户
持续时间: 每个负载级别持续 10 分钟
```

**场景 2: 稳定负载持续测试**
```
负载级别: 2,000 并发用户
持续时间: 2 小时
监控指标: 响应时间稳定性、内存泄漏、CPU 波动
```

**场景 3: 日常业务场景模拟**
```
用户分布:
  - 30% 浏览 Bot 列表
  - 25% 查看对话历史
  - 20% 发送对话消息
  - 15% 搜索知识库
  - 10% 其他操作

请求比例:
  - GET 请求: 60%
  - POST 请求: 30%
  - PUT/PATCH: 8%
  - DELETE: 2%
```

### 2. 压力测试 (Stress Testing)

**目标**: 找出系统性能瓶颈和极限承载能力

**测试策略**:

**策略 1: 并发用户极限测试**
```
起始并发: 1,000 用户
增长步长: 500 用户
每阶段持续时间: 5 分钟
终止条件:
  - 错误率 > 5%
  - 响应时间 P95 > 10s
  - 系统崩溃或无响应
```

**策略 2: 数据量极限测试**
```
测试数据规模:
  - 用户数: 100万
  - Bot 数: 50万
  - 对话数: 1000万
  - 消息数: 1亿
  - 知识库文档: 500万

验证指标:
  - 分页查询性能
  - 索引效率
  - 存储空间增长
```

**策略 3: 混合压力测试**
```
同时施加以下压力:
  - 高并发用户: 5000 并发
  - 大数据量查询: 每秒 100 次
  - 高频写入: 每秒 500 次
  - 复杂计算: 工作流执行每秒 50 次

持续时间: 30 分钟
监控: 系统资源、错误率、响应时间
```

### 3. 峰值测试 (Spike Testing)

**目标**: 验证系统应对突发流量冲击的能力

**测试场景**:

**场景 1: 瞬间流量激增**
```
基线负载: 500 并发用户
峰值负载: 5,000 并发用户（10 倍）
激增方式: 10 秒内从基线达到峰值
峰值持续时间: 5 分钟
恢复阶段: 逐步降至基线

验证: 系统是否崩溃、响应时间是否可接受、恢复能力
```

**场景 2: 营销活动模拟**
```
活动开始前: 1,000 用户在线
活动第 1 分钟: 10,000 用户涌入
活动第 5 分钟: 达到峰值 20,000 用户
活动持续: 30 分钟
活动结束后: 逐步回落

验证: 限流机制、降级策略、用户体验
```

### 4. 耐久测试 (Endurance Testing)

**目标**: 验证系统长时间运行的稳定性

**测试配置**:
```
测试时长: 7×24 小时（1 周）
负载水平: 正常负载的 80%（约 2,000 并发）
监控频率: 每 1 分钟采集一次数据
检查项:
  - 内存泄漏检测
  - 连接池状态
  - 磁盘空间增长
  - 日志文件大小
  - 缓存命中率
  - 响应时间漂移
```

**异常检测**:
```
告警规则:
  - 内存使用率持续增长 > 10%/小时
  - 响应时间 P95 递增 > 20%
  - 错误率突增 > 1%
  - 磁盘空间剩余 < 20%
```

### 5. 容量测试 (Capacity Testing)

**目标**: 确定系统在不同资源配置下的承载能力

**测试矩阵**:

| 配置 | CPU | 内存 | 磁盘 | 数据库 | 最大并发 | 最大 QPS |
|-----|-----|------|------|-------|---------|---------|
| 最小配置 | 4 Core | 8GB | SSD 100GB | 1 节点 | 500 | 1,000 |
| 标准配置 | 8 Core | 16GB | SSD 200GB | 1 主 1 从 | 2,000 | 5,000 |
| 推荐配置 | 16 Core | 32GB | SSD 500GB | 1 主 2 从 | 5,000 | 10,000 |
| 高性能配置 | 32 Core | 64GB | SSD 1TB | 集群模式 | 10,000 | 20,000 |

**测试方法**:
```
1. 从最小配置开始，逐步增加负载
2. 记录每种配置下的最大稳定负载
3. 绘制配置-性能曲线
4. 建立容量规划模型
```

### 6. 并发测试 (Concurrency Testing)

**目标**: 验证系统在并发场景下的数据一致性和死锁风险

**测试场景**:

**场景 1: 同一资源并发修改**
```
测试对象: Bot 配置更新
并发数: 100 个用户同时更新同一个 Bot
预期结果:
  - 所有请求串行化处理
  - 最终状态一致（基于乐观锁或悲观锁）
  - 无数据丢失
```

**场景 2: 计数器并发递增**
```
测试对象: 对话消息计数
操作: 1000 个并发请求 +1
初始值: 0
预期结果: 最终值 = 1000
验证: 数据库行锁、Redis 原子操作
```

**场景 3: 死锁检测**
```
测试场景: 复杂事务交叉更新
操作:
  - 事务 A: 更新 Bot 表 → 更新 Conversation 表
  - 事务 B: 更新 Conversation 表 → 更新 Bot 表

并发执行: 100 对交叉事务
预期结果:
  - 死锁检测机制生效
  - 自动重试或返回友好错误
  - 无系统卡死
```

### 7. 配置测试 (Configuration Testing)

**目标**: 找到系统最优配置参数

**测试维度**:

**数据库连接池配置**:
```yaml
测试参数:
  最大连接数: [50, 100, 200, 300, 500]
  空闲连接数: [10, 20, 50, 100]
  连接超时: [5s, 10s, 30s, 60s]

评估指标: 吞吐量、响应时间、连接等待时间
```

**线程池配置**:
```yaml
测试参数:
  核心线程数: [CPU 核心数 × 1, × 2, × 4]
  最大线程数: [核心线程数 × 2, × 4, × 8]
  队列大小: [100, 500, 1000, 5000]

评估指标: CPU 利用率、任务处理延迟、拒绝策略触发次数
```

**缓存配置**:
```yaml
测试参数:
  Redis 连接池大小: [10, 50, 100, 200]
  本地缓存大小: [1000, 5000, 10000]
  缓存过期时间: [5min, 15min, 30min, 1h]

评估指标: 缓存命中率、内存占用、数据一致性
```

**JVM/Go GC 配置** (如果适用):
```yaml
堆内存: [4GB, 8GB, 16GB, 32GB]
GC 算法: [G1, ZGC, Parallel]
GC 日志: 启用详细 GC 日志分析

评估指标: GC 频率、GC 停顿时间、吞吐量影响
```

---

## 🛠️ 测试工具与框架

### 1. 负载测试工具

**推荐工具**: **Apache JMeter** (主力工具) + **K6** (辅助验证)

**JMeter 优势**:
- 开源免费、功能强大
- 支持多种协议 (HTTP/HTTPS, JDBC, Redis, etc.)
- 丰富的插件生态系统
- 可视化测试报告

**K6 优势**:
- 现代化的脚本语言 (JavaScript)
- 更好的性能和资源占用
- 适合 CI/CD 集成

**JMeter 测试脚本结构**:
```xml
<hashTree>
  <TestPlan guiclass="TestPlanGui">
    <elementProp name="TestPlan.user_defined_variables">
      <collectionProp name="Arguments.arguments">
        <elementProp name="BASE_URL" elementType="Argument">
          <stringProp name="Argument.name">BASE_URL</stringProp>
          <stringProp name="Argument.value">http://localhost:8001</stringProp>
        </elementProp>
        <elementProp name="TOKEN" elementType="Argument">
          <stringProp name="Argument.name">TOKEN</stringProp>
          <stringProp name="Argument.value">${__P(token,)}</stringProp>
        </elementProp>
      </collectionProp>
    </elementProp>
  </TestPlan>

  <!-- 线程组 -->
  <hashTree>
    <ThreadGroup guiclass="ThreadGroupGui">
      <stringProp name="ThreadGroup.num_threads">${__P(users,1000)}</stringProp>
      <stringProp name="ThreadGroup.ramp_time">60</stringProp>
      <stringProp name="ThreadGroup.duration">600</stringProp>
    </ThreadGroup>

    <!-- HTTP 请求默认值 -->
    <ConfigTestElement guiclass="HttpDefaultsGui">
      <stringProp name="HTTPSampler.domain">${BASE_URL}</stringProp>
      <stringProp name="HTTPSampler.port">8001</stringProp>
      <stringProp name="HTTPSampler.protocol">http</stringProp>
    </ConfigTestElement>

    <!-- HTTP 请求 - 获取 Bot 列表 -->
    <HTTPSamplerProxy guiclass="HttpTestSampleGui">
      <stringProp name="HTTPSampler.path">/api/v1/bots</stringProp>
      <stringProp name="HTTPSampler.method">GET</stringProp>
      <elementProp name="HTTPsampler.Arguments">
        <collectionProp name="Arguments.arguments">
          <elementProp name="page" elementType="HTTPArgument">
            <stringProp name="Argument.value">1</stringProp>
          </elementProp>
          <elementProp name="page_size" elementType="HTTPArgument">
            <stringProp name="Argument.value">20</stringProp>
          </elementProp>
        </collectionProp>
      </elementProp>
      <stringProp name="HTTPSampler.header_manager">
        <collectionProp name="HeaderManager.headers">
          <elementProp name="Authorization" elementType="Header">
            <stringProp name="Header.name">Authorization</stringProp>
            <stringProp name="Header.value">Bearer ${TOKEN}</stringProp>
          </elementProp>
        </collectionProp>
      </stringProp>
    </HTTPSamplerProxy>

    <!-- 断言 -->
    <ResultCollector guiclass="ViewResultsFullVisualizer"/>
    <ResultCollector guiclass="SummaryReport"/>
    <ResultCollector guiclass="GraphVisualizer"/>
  </hashTree>
</hashTree>
```

**K6 测试脚本示例**:
```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '2m', target: 100 },   // 2 分钟内爬坡到 100 用户
    { duration: '5m', target: 100 },   // 维持 100 用户 5 分钟
    { duration: '2m', target: 200 },   // 爬坡到 200 用户
    { duration: '5m', target: 200 },   // 维持 200 用户 5 分钟
    { duration: '2m', target: 0 },     // 爬坡到 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 95% 请求 < 2s
    http_req_failed: ['rate<0.01'],    // 错误率 < 1%
  },
};

const BASE_URL = 'http://localhost:8001';
let AUTH_TOKEN = '';

export function setup() {
  // 登录获取 Token
  const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`, JSON.stringify({
    username: 'test_user',
    password: 'test_password',
  }), {
    headers: { 'Content-Type': 'application/json' },
  });

  const token = loginRes.json('data.token');
  return { token };
}

export default function(data) {
  const headers = {
    'Authorization': `Bearer ${data.token}`,
    'Content-Type': 'application/json',
  };

  // 测试场景 1: 获取 Bot 列表
  const listBots = http.get(`${BASE_URL}/api/v1/bots?page=1&page_size=20`, { headers });
  check(listBots, {
    'status is 200': (r) => r.status === 200,
    'has bots': (r) => r.json('data.bots.length') > 0,
  });

  sleep(1);

  // 测试场景 2: 创建对话
  const createConv = http.post(`${BASE_URL}/api/v1/conversations`, JSON.stringify({
    bot_id: 'test_bot_id',
  }), { headers });
  check(createConv, {
    'conversation created': (r) => r.status === 201,
  });

  sleep(2);
}

export function teardown(data) {
  // 清理测试数据（可选）
  console.log('Test completed');
}
```

### 2. 数据库性能测试工具

**工具选择**: **sysbench** (通用基准测试) + 自定义 SQL 脚本

**sysbench 安装**:
```bash
# Ubuntu/Debian
apt-get install sysbench

# CentOS/RHEL
yum install sysbench

# macOS
brew install sysbench
```

**sysbench OLTP 测试**:
```bash
# 准备测试数据（10 张表，每表 100 万行）
sysbench /usr/share/sysbench/oltp_read_write.lua \
  --mysql-host=localhost \
  --mysql-port=3306 \
  --mysql-user=root \
  --mysql-password=password \
  --mysql-db=zker_production \
  --tables=10 \
  --table-size=1000000 \
  prepare

# 运行测试（64 线程，持续 10 分钟）
sysbench /usr/share/sysbench/oltp_read_write.lua \
  --mysql-host=localhost \
  --mysql-port=3306 \
  --mysql-user=root \
  --mysql-password=password \
  --mysql-db=zker_production \
  --tables=10 \
  --table-size=1000000 \
  --threads=64 \
  --time=600 \
  --report-interval=10 \
  run

# 清理测试数据
sysbench /usr/share/sysbench/oltp_read_write.lua \
  --mysql-host=localhost \
  --mysql-port=3306 \
  --mysql-user=root \
  --mysql-password=password \
  --mysql-db=zker_production \
  --tables=10 \
  --table-size=1000000 \
  cleanup
```

**自定义数据库性能测试脚本**:
```sql
-- 测试 1: 复杂 JOIN 查询性能
EXPLAIN ANALYZE
SELECT
  b.bot_id,
  b.name,
  COUNT(c.conv_id) as conversation_count,
  COUNT(m.msg_id) as message_count,
  AVG(
    TIMESTAMPDIFF(SECOND, c.created_at, m.created_at)
  ) as avg_response_time
FROM bots b
LEFT JOIN conversations c ON b.bot_id = c.bot_id
LEFT JOIN messages m ON c.conv_id = m.conv_id
WHERE b.tenant_id = 'test_tenant'
  AND b.status = 'published'
  AND c.created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
GROUP BY b.bot_id, b.name
ORDER BY conversation_count DESC
LIMIT 20;

-- 测试 2: 索引效率验证
SHOW INDEX FROM bots;
SHOW INDEX FROM conversations;
SHOW INDEX FROM messages;

-- 测试 3: 分页查询性能
SELECT * FROM conversations
WHERE tenant_id = 'test_tenant'
ORDER BY created_at DESC
LIMIT 20 OFFSET 10000;

-- 测试 4: 全文检索性能
SELECT * FROM knowledge_documents
WHERE MATCH(title, content) AGAINST('人工智能 机器学习' IN NATURAL LANGUAGE MODE)
  AND tenant_id = 'test_tenant'
LIMIT 20;
```

### 3. 缓存性能测试工具

**工具**: **redis-benchmark** (Redis 官方工具)

**测试命令**:
```bash
# 基础性能测试
redis-benchmark -h localhost -p 6379 -c 50 -n 100000

# SET 操作测试
redis-benchmark -h localhost -p 6379 -t set -c 50 -n 100000 -q

# GET 操作测试
redis-benchmark -h localhost -p 6379 -t get -c 50 -n 100000 -q

# 混合操作测试（SET 占 20%，GET 占 80%）
redis-benchmark -h localhost -p 6379 -t set,get -c 50 -n 100000 -q -r 100000

# Pipeline 批量操作测试
redis-benchmark -h localhost -p 6379 -P 16 -t set,get -c 50 -n 100000 -q

# 自定义脚本测试
redis-benchmark -h localhost -p 6379 -n 100000 -q script load "return redis.call('set', KEYS[1], ARGV[1])"
```

**自定义缓存测试脚本** (Python):
```python
import redis
import time
import random
import string

# 连接 Redis
r = redis.Redis(host='localhost', port=6379, db=0)

# 生成随机数据
def random_key(length=10):
    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))

def random_value(length=100):
    return ''.join(random.choices(string.ascii_letters + string.digits + ' ', k=length))

# 测试 1: SET 性能
def test_set_performance(num_ops=10000):
    start_time = time.time()
    for i in range(num_ops):
        key = f"test:{random_key()}"
        value = random_value()
        r.set(key, value)
    end_time = time.time()
    elapsed = end_time - start_time
    qps = num_ops / elapsed
    print(f"SET 操作: {num_ops} 次耗时 {elapsed:.2f}s, QPS = {qps:.2f}")

# 测试 2: GET 性能
def test_get_performance(num_ops=10000):
    # 预先插入数据
    keys = []
    for i in range(num_ops):
        key = f"test:{random_key()}"
        value = random_value()
        r.set(key, value)
        keys.append(key)

    start_time = time.time()
    for key in keys:
        r.get(key)
    end_time = time.time()
    elapsed = end_time - start_time
    qps = num_ops / elapsed
    print(f"GET 操作: {num_ops} 次耗时 {elapsed:.2f}s, QPS = {qps:.2f}")

# 测试 3: MGET 批量获取性能
def test_mget_performance(num_ops=1000, batch_size=100):
    # 预先插入数据
    keys = []
    for i in range(num_ops):
        key = f"test:{random_key()}"
        value = random_value()
        r.set(key, value)
        keys.append(key)

    start_time = time.time()
    for i in range(0, num_ops, batch_size):
        batch_keys = keys[i:i+batch_size]
        r.mget(batch_keys)
    end_time = time.time()
    elapsed = end_time - start_time
    qps = num_ops / elapsed
    print(f"MGET 操作: {num_ops} 次耗时 {elapsed:.2f}s, QPS = {qps:.2f}, batch_size = {batch_size}")

if __name__ == '__main__':
    test_set_performance(10000)
    test_get_performance(10000)
    test_mget_performance(10000, 100)
```

### 4. 消息队列性能测试工具

**NSQ 性能测试** (使用官方工具):
```bash
# 安装 nsq_stat
go get github.com/nsqio/nsq/nsq_stat

# 查看 NSQ 状态
nsq_stat --lookupd-http-address=localhost:4161

# 使用 nsq_pubsub 进行压测
# 创建生产者
nsq_pubsub --topic=test_topic --channel=test_channel --mode=pub

# 创建消费者
nsq_pubsub --topic=test_topic --channel=test_channel --mode=sub
```

**自定义 NSQ 性能测试脚本** (Go):
```go
package main

import (
  "fmt"
  "log"
  "sync/atomic"
  "time"

  "github.com/nsqio/go-nsq"
)

var (
  publishedCount int64
  consumedCount int64
)

func producer() {
  config := nsq.NewConfig()
  producer, err := nsq.NewProducer("localhost:4150", config)
  if err != nil {
    log.Fatal(err)
  }
  defer producer.Stop()

  startTime := time.Now()
  duration := 60 * time.Second

  for time.Since(startTime) < duration {
    message := []byte(fmt.Sprintf("message_%d", atomic.AddInt64(&publishedCount, 1)))
    err := producer.Publish("test_topic", message)
    if err != nil {
      log.Printf("Publish error: %v", err)
    }
  }

  fmt.Printf("Published: %d messages\n", atomic.LoadInt64(&publishedCount))
}

func consumer() {
  config := nsq.NewConfig()
  consumer, err := nsq.NewConsumer("test_topic", "test_channel", config)
  if err != nil {
    log.Fatal(err)
  }

  consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
    atomic.AddInt64(&consumedCount, 1)
    return nil
  }))

  err = consumer.ConnectToNSQD("localhost:4150")
  if err != nil {
    log.Fatal(err)
  }

  time.Sleep(65 * time.Second)
  consumer.Stop()

  fmt.Printf("Consumed: %d messages\n", atomic.LoadInt64(&consumedCount))
}

func main() {
  go producer()
  time.Sleep(1 * time.Second)
  go consumer()

  time.Sleep(70 * time.Second)
}
```

### 5. 向量数据库性能测试

**Milvus 性能测试** (使用 Python SDK):
```python
from pymilvus import connections, Collection, FieldSchema, CollectionSchema, DataType, utility
import time
import random
import numpy as np

# 连接 Milvus
connections.connect(host='localhost', port='19530')

# 创建测试集合
def create_test_collection():
  if utility.has_collection("test_collection"):
    utility.drop_collection("test_collection")

  fields = [
    FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
    FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=1536)
  ]
  schema = CollectionSchema(fields, "Test collection")
  collection = Collection("test_collection", schema)

  # 创建索引
  index_params = {
    "index_type": "IVF_FLAT",
    "metric_type": "L2",
    "params": {"nlist": 128}
  }
  collection.create_index(field_name="embedding", index_params=index_params)

  return collection

# 插入性能测试
def test_insert_performance(collection, num_entities=100000):
  entities = [
    [i for i in range(num_entities)],  # IDs
    [[random.random() for _ in range(1536)] for _ in range(num_entities)]  # embeddings
  ]

  start_time = time.time()
  collection.insert(entities)
  collection.flush()
  end_time = time.time()

  elapsed = end_time - start_time
  qps = num_entities / elapsed
  print(f"插入 {num_entities} 条数据耗时 {elapsed:.2f}s, QPS = {qps:.2f}")

# 查询性能测试
def test_search_performance(collection, num_searches=1000, top_k=10):
  collection.load()

  search_time_total = 0
  for _ in range(num_searches):
    query_vectors = [[random.random() for _ in range(1536)]]

    start_time = time.time()
    results = collection.search(
      query_vectors,
      anns_field="embedding",
      param={"metric_type": "L2", "params": {"nprobe": 10}},
      limit=top_k
    )
    end_time = time.time()

    search_time_total += (end_time - start_time)

  avg_search_time = search_time_total / num_searches
  qps = num_searches / search_time_total
  print(f"平均查询延迟: {avg_search_time*1000:.2f}ms, QPS = {qps:.2f}")

if __name__ == '__main__':
  collection = create_test_collection()
  test_insert_performance(collection, 100000)
  test_search_performance(collection, 1000, 10)
```

### 6. 前端性能测试工具

**工具选择**:
- **Lighthouse**: 页面性能综合评分
- **WebPageTest**: 网络瀑布图分析
- **Chrome DevTools**: 实时性能分析

**Lighthouse 测试命令**:
```bash
# 安装 Lighthouse
npm install -g lighthouse

# 运行测试
lighthouse http://localhost:8888 --output html --output-path ./report.html

# 性能指标
lighthouse http://localhost:8888 --only-categories=performance

# 自定义配置
lighthouse http://localhost:8888 \
  --throttling-method=devtools \
  --throttling.rttMs=40 \
  --throttling.throughputKbps=10240 \
  --emulated-form-factor=desktop \
  --output html
```

**性能指标解读**:
```
核心 Web 指标 (Core Web Vitals):
  - LCP (Largest Contentful Paint): < 2.5s (良好)
  - FID (First Input Delay): < 100ms (良好)
  - CLS (Cumulative Layout Shift): < 0.1 (良好)

其他重要指标:
  - FCP (First Contentful Paint): < 1.8s
  - TTI (Time to Interactive): < 3.8s
  - SI (Speed Index): < 3.4s
```

### 7. 监控与分析工具

**系统监控**: **Prometheus** + **Grafana**

**Prometheus 配置示例**:
```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'zker-api'
    static_configs:
      - targets: ['localhost:8001']
    metrics_path: '/metrics'

  - job_name: 'mysql'
    static_configs:
      - targets: ['localhost:9104']

  - job_name: 'redis'
    static_configs:
      - targets: ['localhost:9121']

  - job_name: 'node'
    static_configs:
      - targets: ['localhost:9100']
```

**Grafana Dashboard JSON** (关键指标):
```json
{
  "dashboard": {
    "title": "ZKER 性能监控",
    "panels": [
      {
        "title": "QPS (每秒请求数)",
        "targets": [
          {
            "expr": "rate(http_requests_total[1m])"
          }
        ]
      },
      {
        "title": "P95 响应时间",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[1m]))"
          }
        ]
      },
      {
        "title": "错误率",
        "targets": [
          {
            "expr": "rate(http_requests_total{status=~\"5..\"}[1m]) / rate(http_requests_total[1m])"
          }
        ]
      },
      {
        "title": "CPU 使用率",
        "targets": [
          {
            "expr": "100 - (avg by (instance) (irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100)"
          }
        ]
      },
      {
        "title": "内存使用率",
        "targets": [
          {
            "expr": "(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100"
          }
        ]
      }
    ]
  }
}
```

**日志分析**: **ELK Stack** (Elasticsearch + Logstash + Kibana)

**APM 工具**: **Jaeger** (分布式追踪)

---

## 🏗️ 测试环境搭建

### 1. 测试环境配置

**环境隔离原则**:
- 开发环境、测试环境、预发布环境、生产环境完全隔离
- 测试环境配置尽量模拟生产环境
- 数据量级应接近生产实际情况（至少 30%）

**推荐配置**:

| 组件 | 测试环境 | 预发布环境 | 生产环境 |
|-----|---------|-----------|---------|
| 应用服务器 | 2 台 (8C/16G) | 4 台 (16C/32G) | 6+ 台 (16C/32G) |
| 负载均衡器 | 1 台 (Nginx) | 2 台 (HAProxy) | 2+ 台 (HAProxy) |
| 数据库 | 1 主 1 从 (8C/32G) | 1 主 2 从 (16C/64G) | 1 主 3 从 (32C/128G) |
| Redis | 1 节点 (4C/8G) | 哨兵模式 (3 节点) | 集群模式 (6 节点) |
| Milvus | 单机版 (8C/32G) | 分布式 (3 节点) | 分布式 (6+ 节点) |
| Elasticsearch | 单节点 (8C/16G) | 3 节点集群 | 6+ 节点集群 |

### 2. 测试数据准备

**数据量规划**:
```sql
-- 用户数据（100 万用户）
INSERT INTO users (user_id, tenant_id, username, email, status)
SELECT
  CONCAT('user_', n) as user_id,
  CONCAT('tenant_', MOD(n, 1000)) as tenant_id,  -- 1000 个租户
  CONCAT('user_', n) as username,
  CONCAT('user_', n, '@example.com') as email,
  'active' as status
FROM (
  SELECT a.N + b.N * 10 + c.N * 100 + d.N * 1000 + e.N * 10000 + f.N * 100000 as n
  FROM
    (SELECT 0 AS N UNION SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9) a,
    (SELECT 0 AS N UNION SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9) b,
    (SELECT 0 AS N UNION SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9) c,
    (SELECT 0 AS N UNION SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9) d,
    (SELECT 0 AS N UNION SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9) e,
    (SELECT 0 AS N UNION SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9) f
) numbers
WHERE n < 1000000;

-- Bot 数据（50 万 Bot）
INSERT INTO bots (bot_id, tenant_id, created_by, name, description, status)
SELECT
  CONCAT('bot_', n) as bot_id,
  CONCAT('tenant_', MOD(n, 1000)) as tenant_id,
  CONCAT('user_', MOD(n, 1000000)) as created_by,
  CONCAT('Bot ', n) as name,
  CONCAT('Description for bot ', n) as description,
  IF(MOD(n, 10) < 8, 'published', 'draft') as status
FROM numbers
WHERE n < 500000;

-- 对话数据（1000 万对话）
INSERT INTO conversations (conv_id, tenant_id, bot_id, created_by, title, status)
SELECT
  CONCAT('conv_', n) as conv_id,
  CONCAT('tenant_', MOD(n, 1000)) as tenant_id,
  CONCAT('bot_', MOD(n, 500000)) as bot_id,
  CONCAT('user_', MOD(n, 1000000)) as created_by,
  CONCAT('Conversation ', n) as title,
  'active' as status
FROM numbers
WHERE n < 10000000;

-- 消息数据（1 亿条消息）
-- 分批插入（每批 100 万条）
-- 注意：实际插入时需要分批执行，避免单次事务过大
```

**数据准备脚本** (Python):
```python
import pymysql
import random
import string
from datetime import datetime, timedelta

def generate_random_id(prefix, length=10):
  return f"{prefix}_{''.join(random.choices(string.ascii_letters + string.digits, k=length))}"

def generate_test_data(num_users, num_bots, num_convs, num_messages):
  conn = pymysql.connect(
    host='localhost',
    user='root',
    password='password',
    database='zker_test'
  )
  cursor = conn.cursor()

  # 插入用户
  print(f"插入 {num_users} 个用户...")
  for i in range(num_users):
    user_id = generate_random_id('user')
    tenant_id = f"tenant_{random.randint(1, 1000)}"
    username = f"user_{i}"
    email = f"user_{i}@example.com"
    cursor.execute(
      "INSERT INTO users (user_id, tenant_id, username, email, status) VALUES (%s, %s, %s, %s, 'active')",
      (user_id, tenant_id, username, email)
    )
    if i % 10000 == 0:
      conn.commit()
      print(f"已插入 {i} 个用户")
  conn.commit()

  # 插入 Bots
  print(f"插入 {num_bots} 个 Bots...")
  for i in range(num_bots):
    bot_id = generate_random_id('bot')
    tenant_id = f"tenant_{random.randint(1, 1000)}"
    created_by = generate_random_id('user')
    cursor.execute(
      "INSERT INTO bots (bot_id, tenant_id, created_by, name, status) VALUES (%s, %s, %s, %s, 'published')",
      (bot_id, tenant_id, created_by, f"Bot {i}")
    )
    if i % 5000 == 0:
      conn.commit()
      print(f"已插入 {i} 个 Bots")
  conn.commit()

  # ... 其他数据插入

  cursor.close()
  conn.close()

if __name__ == '__main__':
  generate_test_data(
    num_users=1000000,
    num_bots=500000,
    num_convs=10000000,
    num_messages=100000000
  )
```

### 3. 测试环境部署

**Docker Compose 部署**:
```yaml
version: '3.8'

services:
  # 应用服务
  api-server:
    image: zker/api-server:latest
    ports:
      - "8001:8001"
    environment:
      - ENVIRONMENT=test
      - DATABASE_HOST=mysql
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis
    deploy:
      replicas: 2

  # 负载均衡器
  nginx:
    image: nginx:latest
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - api-server

  # 数据库
  mysql:
    image: mysql:8.4.5
    ports:
      - "3306:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=password
      - MYSQL_DATABASE=zker_test
    volumes:
      - mysql_data:/var/lib/mysql
    command: --max-connections=500 --innodb-buffer-pool-size=2G

  # 缓存
  redis:
    image: redis:8.0
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --maxmemory 2gb --maxmemory-policy allkeys-lru

  # 消息队列
  nsqlookupd:
    image: nsqio/nsq:v1.3.1
    command: /nsqlookupd
    ports:
      - "4160:4160"
      - "4161:4161"

  nsqd:
    image: nsqio/nsq:v1.3.1
    command: /nsqd --lookupd-tcp-address=nsqlookupd:4160
    ports:
      - "4150:4150"
      - "4151:4151"
    depends_on:
      - nsqlookupd

  # 监控
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin

volumes:
  mysql_data:
  redis_data:
```

**部署命令**:
```bash
# 启动所有服务
docker compose up -d

# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f api-server

# 停止服务
docker compose down

# 清理数据（危险操作）
docker compose down -v
```

---

## 📊 测试执行流程

### 阶段 1: 基准测试建立 (Baseline Establishment)

**时间安排**: 2 天

**目标**: 建立系统在正常负载下的性能基准

**执行步骤**:

1. **环境验证** (半天)
   ```bash
   # 检查所有服务状态
   docker compose ps

   # 验证数据库连接
   mysql -h127.0.0.1 -P3306 -uroot -p -e "SELECT 1"

   # 验证缓存连接
   redis-cli ping

   # 验证 API 健康检查
   curl http://localhost:8001/health
   ```

2. **功能冒烟测试** (半天)
   ```bash
   # 运行最小化测试套件
   rush test --smoke

   # 手动验证核心功能
   # - 用户登录
   # - 创建 Bot
   # - 发送对话
   # - 搜索知识库
   ```

3. **基准性能测试** (1 天)
   ```bash
   # JMeter 脚本：baseline_test.jmx
   # 负载级别：500 并发用户
   # 持续时间：30 分钟

   jmeter -n -t baseline_test.jmx -l baseline_results.jtl -e -o baseline_report

   # 记录基准指标
   # - 平均响应时间: xxx ms
   # - P95 响应时间: xxx ms
   # - QPS: xxx
   # - 错误率: xxx%
   # - CPU 使用率: xxx%
   # - 内存使用率: xxx%
   ```

**输出物**:
- [ ] 基准性能测试报告 (`baseline_report.html`)
- [ ] 基准指标清单 (`baseline_metrics.csv`)
- [ ] 环境配置记录 (`environment_config.yaml`)

### 阶段 2: 负载测试 (Load Testing)

**时间安排**: 3 天

**目标**: 验证系统在不同负载级别下的性能表现

**测试矩阵**:

| 测试场景 | 并发用户 | 持续时间 | 预期响应时间 P95 |
|---------|---------|---------|----------------|
| 低负载 | 500 | 1 小时 | < 500ms |
| 中负载 | 1,000 | 1 小时 | < 1s |
| 高负载 | 2,000 | 1 小时 | < 2s |
| 接近峰值 | 3,000 | 30 分钟 | < 3s |
| 峰值负载 | 5,000 | 15 分钟 | < 5s |

**执行命令**:
```bash
# 场景 1: 低负载测试
jmeter -n -t load_test_low.jmx \
  -Jusers=500 \
  -Jduration=3600 \
  -l load_low_results.jtl \
  -e -o load_low_report

# 场景 2: 中负载测试
jmeter -n -t load_test_medium.jmx \
  -Jusers=1000 \
  -Jduration=3600 \
  -l load_medium_results.jtl \
  -e -o load_medium_report

# 场景 3: 高负载测试
jmeter -n -t load_test_high.jmx \
  -Jusers=2000 \
  -Jduration=3600 \
  -l load_high_results.jtl \
  -e -o load_high_report

# 场景 4: 接近峰值测试
jmeter -n -t load_test_near_peak.jmx \
  -Jusers=3000 \
  -Jduration=1800 \
  -l load_near_peak_results.jtl \
  -e -o load_near_peak_report

# 场景 5: 峰值负载测试
jmeter -n -t load_test_peak.jmx \
  -Jusers=5000 \
  -Jduration=900 \
  -l load_peak_results.jtl \
  -e -o load_peak_report
```

**监控重点**:
- 响应时间增长曲线
- 错误率变化
- 资源利用率趋势
- 数据库慢查询日志
- 应用日志中的错误和警告

**输出物**:
- [ ] 各负载级别测试报告 (5 个 HTML 报告)
- [ ] 性能趋势图 (`performance_trend.png`)
- [ ] 瓶颈分析报告 (`bottleneck_analysis.md`)

### 阶段 3: 压力测试 (Stress Testing)

**时间安排**: 2 天

**目标**: 找出系统性能极限和故障点

**测试策略**:

**策略 1: 渐进式加压**
```
起始: 1000 并发
增量: 每分钟增加 200 并发
上限: 直到错误率 > 10% 或系统崩溃
```

**JMeter 配置**:
```xml
<ThreadGroup guiclass="SteppingThreadGroup">
  <stringProp name="ThreadGroup.num_threads">10000</stringProp>
  <stringProp name="ThreadGroup.ramp_time">600</stringProp>
  <stringProp name="SteppingThreadGroup.startUsers">1000</stringProp>
  <stringProp name="SteppingThreadGroup.addUsers">200</stringProp>
  <stringProp name="SteppingThreadGroup.addUsersPeriod">60</stringProp>
</ThreadGroup>
```

**执行命令**:
```bash
jmeter -n -t stress_test_ramp_up.jmx \
  -l stress_ramp_up_results.jtl \
  -e -o stress_ramp_up_report
```

**策略 2: 瞬间高压**
```
场景: 瞬间达到 10000 并发
持续时间: 5 分钟
目的: 测试系统抗冲击能力
```

**执行命令**:
```bash
jmeter -n -t stress_test_spike.jmx \
  -Jusers=10000 \
  -Jramp_time=10 \
  -Jduration=300 \
  -l stress_spike_results.jtl \
  -e -o stress_spike_report
```

**故障点识别**:
```
监控以下指标，记录首次出现异常的时刻:
1. HTTP 500 错误出现
2. 响应时间 P95 > 10s
3. 数据库连接池耗尽
4. Redis 连接超时
5. CPU 使用率 > 95%
6. 内存使用率 > 95%
7. 磁盘 I/O 等待时间 > 1s
```

**输出物**:
- [ ] 压力测试报告 (`stress_test_report.md`)
- [ ] 系统极限承载能力 (`system_limit.md`)
- [ ] 故障点清单 (`failure_points.md`)

### 阶段 4: 峰值测试 (Spike Testing)

**时间安排**: 1 天

**目标**: 验证系统应对突发流量的能力

**测试场景**:

**场景 1: 10 倍瞬时流量**
```
基线: 500 并发
峰值: 5,000 并发（10 倍）
爬坡时间: 10 秒
峰值持续时间: 5 分钟
```

**K6 脚本**:
```javascript
export const options = {
  stages: [
    { duration: '2m', target: 500 },    // 基线
    { duration: '10s', target: 5000 },  // 爬坡到峰值
    { duration: '5m', target: 5000 },   // 维持峰值
    { duration: '2m', target: 500 },    // 回到基线
  ],
};

export default function() {
  // 测试逻辑...
}
```

**场景 2: 营销活动模拟**
```
活动前: 1,000 用户
活动开始: 10 秒内涌入 20,000 用户
活动持续: 30 分钟
活动结束: 逐步回落
```

**验证点**:
- [ ] 限流机制是否生效
- [ ] 降级策略是否触发
- [ ] 用户体验是否可接受
- [ ] 系统恢复能力

**输出物**:
- [ ] 峰值测试报告 (`spike_test_report.md`)
- [ ] 限流和降级日志 (`throttle_logs.txt`)
- [ ] 用户体验评估 (`user_experience.md`)

### 阶段 5: 耐久测试 (Endurance Testing)

**时间安排**: 7 天（可并行）

**目标**: 验证系统长时间运行的稳定性

**测试配置**:
```
负载: 2,000 并发用户（80% 额定负载）
持续时间: 7×24 小时
监控频率: 每 1 分钟
数据收集: 所有指标
```

**JMeter 配置**:
```xml
<ThreadGroup guiclass="ThreadGroupGui">
  <stringProp name="ThreadGroup.num_threads">2000</stringProp>
  <stringProp name="ThreadGroup.ramp_time">300</stringProp>
  <stringProp name="ThreadGroup.duration">604800</stringProp>  <!-- 7 天 -->
  <stringProp name="ThreadGroup.scheduler">true</stringProp>
  <stringProp name="ThreadGroup.duration">604800</stringProp>
</ThreadGroup>
```

**执行命令**:
```bash
# 后台运行测试
nohup jmeter -n -t endurance_test.jmx \
  -l endurance_results.jtl \
  -e -o endurance_report > endurance_test.log 2>&1 &

# 监控测试进度
tail -f endurance_test.log

# 定期检查系统状态
watch -n 60 'ps aux | grep jmeter'
```

**每日检查项**:
```bash
# 每天执行一次健康检查

# 1. 资源使用率
top -b -n 1 | head -20

# 2. 内存泄漏检测
ps aux | grep jmeter | awk '{print $6}' | awk '{sum+=$1} END {print sum}'

# 3. 磁盘空间
df -h

# 4. 数据库连接数
mysql -e "SHOW PROCESSLIST" | wc -l

# 5. 错误日志统计
tail -10000 /var/log/zker/application.log | grep "ERROR" | wc -l

# 6. 响应时间趋势
grep "elapsed" endurance_results.jtl | awk -F',' '{print $2}' | awk '{sum+=$1; count++} END {print sum/count}'
```

**异常检测**:
```
告警触发条件:
- 内存使用率 24 小时内持续增长 > 20%
- 响应时间 P95 递增 > 30%
- 错误率突增 > 2%
- 磁盘空间剩余 < 15%
- CPU 长时间 > 90% (> 10 分钟)
```

**输出物**:
- [ ] 7 天耐久测试报告 (`endurance_test_report.md`)
- [ ] 资源使用趋势图 (`resource_trend_7days.png`)
- [ ] 内存泄漏分析 (`memory_leak_analysis.md`)
- [ ] 稳定性评估 (`stability_assessment.md`)

### 阶段 6: 容量测试 (Capacity Testing)

**时间安排**: 2 天

**目标**: 建立配置-性能映射模型

**测试矩阵**:

| 测试用例 | CPU | 内存 | 数据库 | 最大并发 | 最大 QPS |
|---------|-----|------|-------|---------|---------|
| TC-1 | 4C | 8GB | 1 节点 | 500 | 1,000 |
| TC-2 | 8C | 16GB | 1 主 1 从 | 2,000 | 5,000 |
| TC-3 | 16C | 32GB | 1 主 2 从 | 5,000 | 10,000 |
| TC-4 | 32C | 64GB | 集群 | 10,000 | 20,000 |

**测试方法**:
```
对每种配置:
1. 部署对应规格的测试环境
2. 运行负载测试，从低负载逐步增加
3. 记录最大稳定负载（错误率 < 1%，P95 < 2s）
4. 绘制性能曲线
```

**容量规划模型**:
```
建立公式: 最大 QPS = f(CPU, 内存, 数据库节点数)

示例（基于测试数据）:
QPS = 125 * CPU核心数 * sqrt(内存GB) * (0.8 + 0.2 * 数据库节点数)

验证:
TC-2: QPS = 125 * 8 * sqrt(16) * (0.8 + 0.2 * 2) = 125 * 8 * 4 * 1.2 = 4,800 ≈ 5,000 ✓
```

**输出物**:
- [ ] 容量测试报告 (`capacity_test_report.md`)
- [ ] 配置-性能对照表 (`config_performance_map.csv`)
- [ ] 容量规划模型 (`capacity_planning_model.xlsx`)

### 阶段 7: 专项性能测试

**时间安排**: 3 天

**测试范围**:

**1. 数据库性能** (1 天)
```bash
# sysbench OLTP 测试
sysbench oltp_read_write \
  --mysql-host=localhost \
  --mysql-port=3306 \
  --mysql-db=zker_test \
  --tables=10 \
  --table-size=1000000 \
  --threads=64 \
  --time=600 \
  run

# 慢查询分析
mysql> SELECT * FROM mysql.slow_log ORDER BY query_time DESC LIMIT 100;

# 索引效率验证
EXPLAIN SELECT ...;
```

**2. 缓存性能** (半天)
```bash
# Redis 性能测试
redis-benchmark -h localhost -p 6379 -c 50 -n 100000 -t set,get -q

# 缓存命中率分析
redis-cli info stats | grep keyspace_hits
```

**3. 消息队列性能** (半天)
```bash
# NSQ 吞吐量测试
# 运行自定义 Go 测试程序
go run nsq_perf_test.go

# 监控消息堆积
nsq_stat --lookupd-http-address=localhost:4161
```

**4. 向量检索性能** (半天)
```python
# 运行 Milvus 性能测试脚本
python milvus_perf_test.py

# 测试不同索引类型
# - IVF_FLAT
# - IVF_SQ8
# - HNSW
```

**输出物**:
- [ ] 数据库性能测试报告 (`database_perf_report.md`)
- [ ] 缓存性能测试报告 (`cache_perf_report.md`)
- [ ] 消息队列性能测试报告 (`mq_perf_report.md`)
- [ ] 向量检索性能测试报告 (`vector_search_perf_report.md`)

---

## 🐛 性能问题诊断与优化

### 常见性能瓶颈与解决方案

### 1. 数据库性能问题

**问题 1: 慢查询**
```
症状: 某些 API 响应时间长，数据库日志出现慢查询

诊断步骤:
1. 开启慢查询日志
   SET GLOBAL slow_query_log = 'ON';
   SET GLOBAL long_query_time = 0.5;

2. 分析慢查询日志
   mysqldumpslow -s t -t 10 /var/log/mysql/slow-query.log

3. 使用 EXPLAIN 分析查询计划
   EXPLAIN SELECT * FROM conversations WHERE tenant_id = 'xxx';

常见原因:
- 缺少索引
- 索引未生效
- 查询字段过多（SELECT *）
- 子查询未优化
- JOIN 顺序不当

解决方案:
- 添加合适的索引
- 使用覆盖索引
- 避免 SELECT *，只查询需要的字段
- 将子查询改为 JOIN
- 优化 JOIN 顺序
```

**问题 2: 连接池耗尽**
```
症状: API 返回 "Connection pool exhausted" 错误

诊断:
1. 检查当前连接数
   SHOW PROCESSLIST;

2. 检查连接池配置
   max_connections
   max_idle_connections

3. 检查是否有连接泄漏

解决方案:
1. 增加连接池大小
   max_connections = 500

2. 优化连接生命周期
   set_max_lifetime = 30 分钟

3. 检查代码是否正确释放连接
   defer db.Close()

4. 使用连接池监控
   定期检查活跃连接数
```

**问题 3: 锁等待和死锁**
```
症状: 事务超时、响应时间长

诊断:
1. 查看锁等待情况
   SELECT * FROM information_schema.INNODB_LOCKS;
   SELECT * FROM information_schema.INNODB_LOCK_WAITS;

2. 查看死锁日志
   SHOW ENGINE INNODB STATUS;

解决方案:
1. 减少事务持有锁的时间
   - 快速提交事务
   - 避免事务中执行长时间操作

2. 优化锁的粒度
   - 尽量使用行锁，避免表锁
   - 使用乐观锁替代悲观锁

3. 统一访问顺序
   - 不同事务按相同顺序访问表

4. 添加死锁重试机制
   retry: 3 次
```

### 2. 应用服务性能问题

**问题 1: CPU 密集型操作**
```
症状: CPU 使用率持续 > 80%

诊断:
1. 使用 pprof 进行性能分析
   go tool pprof http://localhost:8001/debug/pprof/profile?seconds=30

2. 查看 CPU 火焰图
   go tool pprof -http=:8080 cpu.prof

常见原因:
- 复杂计算（如相似度计算）
- 正则表达式匹配
- JSON 序列化/反序列化
- 加密/解密操作

解决方案:
1. 算法优化
   - 使用更高效的算法
   - 减少不必要的计算

2. 异步处理
   - 将耗时操作放入消息队列异步处理

3. 缓存计算结果
   - 使用 Redis 缓存计算结果

4. 并行计算
   - 使用 goroutine 并行处理
```

**问题 2: 内存泄漏**
```
症状: 内存使用率持续增长，最终 OOM

诊断:
1. 使用 pprof 进行内存分析
   go tool pprof http://localhost:8001/debug/pprof/heap

2. 对比不同时间点的内存快照
   go tool pprof -base heap_old.prof heap_new.prof

常见原因:
- 未关闭的资源（文件、连接）
- goroutine 泄漏
- 缓存无限增长
- 循环引用

解决方案:
1. 及时关闭资源
   defer file.Close()
   defer conn.Close()

2. 控制 goroutine 数量
   - 使用 worker pool 模式
   - 添加超时控制

3. 限制缓存大小
   - 使用 LRU 缓存
   - 设置过期时间

4. 定期触发 GC
   runtime.GC()
```

**问题 3: 协程泄漏**
```
症状: goroutine 数量持续增长

诊断:
1. 查看 goroutine 数量
   http://localhost:8001/debug/pprof/goroutine

2. 分析 goroutine 堆栈
   go tool pprof goroutine.prof

常见原因:
- 未处理的 channel 阻塞
- 未设置超时的 context
- for-select 循环未退出

解决方案:
1. 添加超时控制
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()

2. 使用 select 防止阻塞
   select {
   case <-ch:
   case <-ctx.Done():
     return
   }

3. 及时退出 goroutine
   - 使用 done channel
   - 使用 context 取消
```

### 3. 缓存性能问题

**问题 1: 缓存穿透**
```
症状: 大量请求查询不存在的数据，直接打到数据库

解决方案:
1. 布隆过滤器
   - 将所有合法 key 存入布隆过滤器
   - 查询前先检查 key 是否可能存在

2. 缓存空对象
   - 将查询为空的结果也缓存（TTL 较短）
   - SET key "NULL" EX 60

3. 请求参数校验
   - 在应用层先校验参数合法性
```

**问题 2: 缓存雪崩**
```
症状: 大量缓存同时失效，数据库瞬间压力激增

解决方案:
1. 缓存过期时间加随机值
   TTL = base_ttl + random(0, 300)

2. 多级缓存
   - 本地缓存（L1）+ Redis（L2）
   - 分级降级

3. 缓存预热
   - 系统启动时加载热点数据
   - 定期刷新缓存

4. 互斥锁重建缓存
   - 使用分布式锁（Redis SETNX）
   - 只允许一个请求重建缓存
```

**问题 3: 缓存击穿**
```
症状: 热点 key 过期瞬间，大量请求打到数据库

解决方案:
1. 热点数据永不过期
   - 逻辑过期：在 value 中记录过期时间
   - 后台异步更新

2. 互斥锁
   - SETNX lock_key 1 EX 10
   - 获取锁的请求负责重建缓存
   - 其他请求等待或返回旧缓存

3. 限流降级
   - 对热点 key 进行限流
   - 降级返回默认值
```

### 4. 网络性能问题

**问题 1: 高延迟**
```
症状: 响应时间长，但应用处理快

诊断:
1. 使用 traceroute 追踪路由
   traceroute api-server

2. 使用 tcpdump 抓包分析
   tcpdump -i any host api-server -w capture.pcap

3. 检查 DNS 查询
   dig api-server.example.com

解决方案:
1. 使用 CDN 加速静态资源
2. 启用 HTTP/2 多路复用
3. 使用连接池
4. 优化 TCP 参数
   net.ipv4.tcp_fastopen = 3
   net.core.somaxconn = 65535
```

**问题 2: 带宽瓶颈**
```
症状: 传输大文件时速度慢

解决方案:
1. 启用压缩
   - Gzip / Brotli 压缩响应内容

2. 分片传输
   - 大文件分片上传/下载
   - 使用断点续传

3. 使用更快的协议
   - gRPC 替代 REST API
   - WebSocket 替代轮询
```

---

## 📈 性能优化建议

### 数据库优化

**1. 索引优化**
```sql
-- 1. 添加合适的索引
CREATE INDEX idx_tenant_created ON conversations(tenant_id, created_at);
CREATE INDEX idx_bot_status ON bots(bot_id, status);

-- 2. 使用覆盖索引
CREATE INDEX idx_conv_cover ON conversations(conv_id, tenant_id, bot_id, status, created_at);

-- 3. 删除冗余索引
-- 检查未使用的索引
SELECT * FROM sys.schema_unused_indexes WHERE object_schema = 'zker_production';

-- 4. 定期分析和优化表
ANALYZE TABLE bots;
OPTIMIZE TABLE conversations;
```

**2. 查询优化**
```sql
-- 避免 SELECT *
-- ❌ 不好的写法
SELECT * FROM conversations WHERE tenant_id = 'xxx';

-- ✅ 好的写法
SELECT conv_id, bot_id, title, created_at
FROM conversations
WHERE tenant_id = 'xxx'
ORDER BY created_at DESC
LIMIT 20;

-- 使用 LIMIT 避免全表扫描
SELECT * FROM messages WHERE conv_id = 'xxx' LIMIT 1000;

-- 使用 UNION ALL 替代 UNION（避免去重开销）
SELECT user_id FROM bots WHERE status = 'published'
UNION ALL
SELECT user_id FROM conversations WHERE status = 'active';
```

**3. 架构优化**
```
1. 读写分离
   - 主库处理写操作
   - 从库处理读操作
   - 使用代理（如 ProxySQL）

2. 分库分表
   - 按租户 ID 分库
   - 按时间分表（如按月）

3. 归档历史数据
   - 将 6 个月前的数据迁移到归档库
   - 减少主库数据量
```

### 应用服务优化

**1. 代码优化**
```go
// ❌ 不好的写法：N+1 查询
func GetBotsWithConversations(ctx context.Context, botIDs []string) ([]*Bot, error) {
  bots := make([]*Bot, 0, len(botIDs))
  for _, botID := range botIDs {
    bot, err := botRepo.FindByID(ctx, botID)
    if err != nil {
      return nil, err
    }
    bots = append(bots, bot)
  }
  return bots, nil
}

// ✅ 好的写法：批量查询
func GetBotsWithConversations(ctx context.Context, botIDs []string) ([]*Bot, error) {
  return botRepo.FindByIDs(ctx, botIDs)  // 一次查询所有 bot
}

// ❌ 不好的写法：在循环中处理
func ProcessConversations(convIDs []string) error {
  for _, convID := range convIDs {
    result := heavyCalculation(convID)
    saveResult(convID, result)
  }
  return nil
}

// ✅ 好的写法：并行处理
func ProcessConversations(convIDs []string) error {
  var wg sync.WaitGroup
  sem := make(chan struct{}, 10)  // 限制并发数为 10

  for _, convID := range convIDs {
    wg.Add(1)
    go func(id string) {
      defer wg.Done()
      sem <- struct{}{}        // 获取信号量
      defer func() { <-sem }() // 释放信号量

      result := heavyCalculation(id)
      saveResult(id, result)
    }(convID)
  }

  wg.Wait()
  return nil
}
```

**2. 连接池优化**
```go
// 数据库连接池配置
db.SetMaxOpenConns(200)       // 最大连接数
db.SetMaxIdleConns(50)        // 最大空闲连接数
db.SetConnMaxLifetime(1 * time.Hour)  // 连接最大生命周期
db.SetConnMaxIdleTime(10 * time.Minute) // 空闲连接最大存活时间

// Redis 连接池配置
redisClient := redis.NewClient(&redis.Options{
  Addr:         "localhost:6379",
  PoolSize:     50,
  MinIdleConns: 10,
  MaxRetries:   3,
  DialTimeout:  5 * time.Second,
  ReadTimeout:  3 * time.Second,
  WriteTimeout: 3 * time.Second,
  PoolTimeout:  4 * time.Second,
})
```

**3. 内存优化**
```go
// 使用对象池减少内存分配
var botPool = sync.Pool{
  New: func() interface{} {
    return &Bot{}
  },
}

func getBotFromPool() *Bot {
  return botPool.Get().(*Bot)
}

func putBotToPool(bot *Bot) {
  bot.Reset()  // 重置对象
  botPool.Put(bot)
}

// 使用 bytes.Buffer 避免 string 拼接
// ❌ 不好的写法
var s string
for i := 0; i < 1000; i++ {
  s += strconv.Itoa(i)  // 每次都创建新的 string
}

// ✅ 好的写法
var b bytes.Buffer
for i := 0; i < 1000; i++ {
  b.WriteString(strconv.Itoa(i))
}
s := b.String()
```

### 缓存优化

**1. 多级缓存**
```go
// L1: 本地缓存（内存）
// L2: Redis 缓存
// L3: 数据库

func GetBot(ctx context.Context, botID string) (*Bot, error) {
  // L1: 本地缓存
  if bot, ok := localCache.Get(botID); ok {
    return bot, nil
  }

  // L2: Redis 缓存
  val, err := redisClient.Get(ctx, "bot:"+botID).Result()
  if err == nil {
    var bot Bot
    json.Unmarshal([]byte(val), &bot)
    localCache.Set(botID, &bot, 5*time.Minute)  // 写入 L1
    return &bot, nil
  }

  // L3: 数据库
  bot, err := botRepo.FindByID(ctx, botID)
  if err != nil {
    return nil, err
  }

  // 回写缓存
  data, _ := json.Marshal(bot)
  redisClient.Set(ctx, "bot:"+botID, data, 30*time.Minute)
  localCache.Set(botID, bot, 5*time.Minute)

  return bot, nil
}
```

**2. 缓存预热**
```go
// 系统启动时加载热点数据
func WarmUpCache(ctx context.Context) error {
  // 1. 加载热门 Bots
  popularBots, _ := botRepo.FindPopular(ctx, 1000)
  for _, bot := range popularBots {
    data, _ := json.Marshal(bot)
    redisClient.Set(ctx, "bot:"+bot.BotID, data, 1*time.Hour)
  }

  // 2. 加载活跃用户
  activeUsers, _ := userRepo.FindActive(ctx, 10000)
  for _, user := range activeUsers {
    data, _ := json.Marshal(user)
    redisClient.Set(ctx, "user:"+user.UserID, data, 30*time.Minute)
  }

  return nil
}
```

**3. 缓存更新策略**
```go
// Cache Aside 模式
func UpdateBot(ctx context.Context, bot *Bot) error {
  // 1. 先更新数据库
  err := botRepo.Update(ctx, bot)
  if err != nil {
    return err
  }

  // 2. 再删除缓存（而非更新缓存）
  redisClient.Del(ctx, "bot:"+bot.BotID)
  localCache.Del(bot.BotID)

  return nil
}

// 为什么删除而非更新？
// - 避免并发更新导致缓存不一致
// - 下次查询时再加载，实现懒加载
```

---

## 📝 性能测试报告模板

### 报告结构

```markdown
# ZKER 性能测试报告

## 1. 测试概述

### 1.1 测试目标
### 1.2 测试范围
### 1.3 测试环境
### 1.4 测试工具

## 2. 测试结果汇总

### 2.1 关键指标
| 指标 | 目标值 | 实际值 | 达标情况 |
|-----|-------|-------|---------|
| P95 响应时间 | < 2s | 1.5s | ✅ |
| QPS | 10,000 | 12,000 | ✅ |
| 错误率 | < 0.1% | 0.05% | ✅ |
| CPU 使用率 | < 70% | 65% | ✅ |
| 内存使用率 | < 80% | 75% | ✅ |

### 2.2 性能趋势图
[插入图表]

## 3. 详细测试结果

### 3.1 负载测试
[各负载级别下的详细数据]

### 3.2 压力测试
[极限负载测试结果]

### 3.3 峰值测试
[突发流量测试结果]

### 3.4 耐久测试
[7 天稳定性测试结果]

### 3.5 专项测试
[数据库、缓存、MQ 等专项测试结果]

## 4. 性能瓶颈分析

### 4.1 发现的问题
1. **问题 1**: Bot 列表查询慢
   - 影响: P95 响应时间 2.3s
   - 原因: 缺少索引
   - 优先级: P0

2. **问题 2**: 数据库连接池耗尽
   - 影响: 高并发下错误率 1.5%
   - 原因: 连接池配置过小
   - 优先级: P0

### 4.2 优化建议
[详细的优化建议和实施计划]

## 5. 容量规划

### 5.1 配置-性能对照表
[不同配置下的性能数据]

### 5.2 容量规划模型
[容量规划公式和计算]

## 6. 结论

### 6.1 总体评估
[系统是否满足上线要求]

### 6.2 风险评估
[已知风险和缓解措施]

### 6.3 后续计划
[下一步优化计划]
```

---

## 🎯 性能测试验收标准

### 验收清单

**功能验收**:
- [ ] 所有核心功能在高负载下正常运行
- [ ] 无数据丢失、数据不一致
- [ ] 无内存泄漏、连接泄漏
- [ ] 错误处理机制正常

**性能验收**:
- [ ] P95 响应时间 < 2s
- [ ] P99 响应时间 < 5s
- [ ] QPS ≥ 10,000
- [ ] 错误率 < 0.1%
- [ ] CPU 使用率 < 70%
- [ ] 内存使用率 < 80%

**稳定性验收**:
- [ ] 7×24 小时稳定运行
- [ ] 无服务崩溃
- [ ] 无性能严重退化（> 20%）
- [ ] 资源使用稳定

**可扩展性验收**:
- [ ] 水平扩展能力验证
- [ ] 负载均衡有效性
- [ ] 数据库读写分离
- [ ] 缓存集群稳定性

### 上线决策

**满足以下条件可上线**:
- ✅ 所有 P0 性能指标达标
- ✅ 无 P0 级性能问题
- ✅ 压力测试下系统稳定
- ✅ 监控和告警完善
- ✅ 降级和容灾方案就绪

**不满足以下条件暂缓上线**:
- ❌ 存在 P0 性能问题未解决
- ❌ 压力测试下系统崩溃
- ❌ 监控盲区
- ❌ 无降级方案

---

## 📚 附录

### A. 性能测试常用指标定义

| 指标 | 定义 | 计算方式 |
|-----|------|---------|
| QPS | Queries Per Second，每秒请求数 | 总请求数 / 总时间 |
| TPS | Transactions Per Second，每秒事务数 | 完成事务数 / 总时间 |
| RT | Response Time，响应时间 | 请求完成时间 - 请求发送时间 |
| P50 | 中位数响应时间 | 50% 请求的响应时间 |
| P95 | 95 分位响应时间 | 95% 请求的响应时间 |
| P99 | 99 分位响应时间 | 99% 请求的响应时间 |
| 并发用户数 | 同时在线的用户数 | N |
| 吞吐量 | 系统单位时间处理的数据量 | 数据量 / 时间 |

### B. JMeter 常用断言

```xml
<!-- 响应代码断言 -->
<ResponseAssertion guiclass="ResponseCodeGui">
  <collectionProp name="Asserion.test_strings">
    <stringProp name="49586">200</stringProp>
  </collectionProp>
</ResponseAssertion>

<!-- 响应时间断言 -->
<DurationAssertion guiclass="DurationAssertionGui">
  <stringProp name="DurationAssertion.duration">2000</stringProp>
</DurationAssertion>

<!-- JSON 响应断言 -->
<JSONPathAssertion guiclass="JSONPathAssertionGui">
  <stringProp name="JSON_PATH">$.code</stringProp>
  <stringProp name="EXPECTED_VALUE">SUCCESS</stringProp>
</JSONPathAssertion>
```

### C. 常用性能监控命令

```bash
# 系统资源监控
top -b -n 1 | head -20
vmstat 1 10
iostat -x 1 10
free -h
df -h

# 网络监控
netstat -an | grep ESTABLISHED | wc -l
ss -s
iftop

# 数据库监控
mysql -e "SHOW PROCESSLIST"
mysql -e "SHOW ENGINE INNODB STATUS"
mysql -e "SHOW STATUS LIKE 'Threads_%'"

# Redis 监控
redis-cli info
redis-cli --stat
redis-cli info stats

# 应用监控
curl http://localhost:8001/health
curl http://localhost:8001/metrics
```

### D. 性能测试相关资源

**工具下载**:
- Apache JMeter: https://jmeter.apache.org/download_jmeter.cgi
- K6: https://k6.io/
- sysbench: https://github.com/akopytov/sysbench
- Prometheus: https://prometheus.io/download/
- Grafana: https://grafana.com/grafana/download

**参考文档**:
- JMeter User Manual: https://jmeter.apache.org/usermanual/index.html
- Go Profiling: https://go.dev/doc/diagnostics
- MySQL Performance: https://dev.mysql.com/doc/refman/8.4/en/optimization.html
- Redis Performance: https://redis.io/topics/benchmarks

---

**文档变更历史**:

| 版本 | 日期 | 变更内容 | 作者 |
|-----|------|---------|------|
| v1.0 | 2025-01-01 | 初始版本 | 性能测试团队 |

**审批记录**:

| 角色 | 姓名 | 审批意见 | 日期 |
|-----|------|---------|------|
| 技术架构委员会 | [待填写] | [待审批] | [待审批] |
| 测试负责人 | [待填写] | [待审批] | [待审批] |
| 项目经理 | [待填写] | [待审批] | [待审批] |

---

**© 2025 ZKER Project. All rights reserved.**
