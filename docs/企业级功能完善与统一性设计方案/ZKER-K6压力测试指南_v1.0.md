# ZKER K6压力测试指南

**版本**: v1.0.0
**最后更新**: 2025-01-03
**作者**: 研发B (后端工程师)
**适用阶段**: 性能测试、性能优化、CI/CD集成

---

## 目录

- [1. 概述](#1-概述)
- [2. K6简介](#2-k6简介)
- [3. 测试环境准备](#3-测试环境准备)
- [4. 测试场景详解](#4-测试场景详解)
- [5. 测试执行](#5-测试执行)
- [6. 结果分析](#6-结果分析)
- [7. 性能基线](#7-性能基线)
- [8. 故障排查](#8-故障排查)
- [9. CI/CD集成](#9-cicd集成)
- [10. 最佳实践](#10-最佳实践)

---

## 1. 概述

### 1.1 测试目标

ZKER压力测试体系旨在验证系统在高负载下的性能表现，确保：

- **吞吐量**: 系统能够支持预期的并发用户数
- **延迟**: API响应时间在可接受范围内
- **稳定性**: 长时间运行下系统稳定可靠
- **可扩展性**: 系统能够线性扩展

### 1.2 测试范围

```
tests/performance/k6/
├── knowledge_api_test.js        # 知识库API测试
├── bot_api_test.js              # Bot API测试
├── conversation_load_test.js    # 会话负载测试
├── permission_test.js           # 权限检查测试
├── mixed_load_test.js           # 混合负载测试
├── k6.config.json               # 测试配置
└── README.md                    # 使用文档
```

---

## 2. K6简介

### 2.1 为什么选择K6?

**优势**:
- ✅ 开源免费，社区活跃
- ✅ 基于JavaScript，易于编写和维护
- ✅ 支持多种负载模式
- ✅ 丰富的指标和可视化
- ✅ 良好的CI/CD集成

**核心概念**:
- **VU (Virtual User)**: 虚拟用户
- **Iteration**: 迭代，VU执行脚本的一次完整循环
- **Scenario**: 测试场景，定义负载模式
- **Stage**: 阶段，定义测试过程中的负载变化

### 2.2 安装K6

#### Linux/Mac
```bash
# Mac
brew install k6

# Linux (Debian/Ubuntu)
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6
```

#### Windows
```powershell
# 使用Chocolatey
choco install k6

# 或下载二进制文件
# https://k6.io/docs/getting-started/installation/
```

验证安装:
```bash
k6 version
# 输出: k6 v0.45.0
```

---

## 3. 测试环境准备

### 3.1 启动服务

```bash
# 1. 启动监控栈
cd D:\code\coze-studio
docker-compose -f docker-compose.monitoring.yml up -d

# 2. 启动ZKER服务
make server

# 3. 验证服务健康
curl http://localhost:8888/health
```

### 3.2 配置环境变量

创建`.env`文件：

```bash
# ZKER服务配置
BASE_URL=http://localhost:8888

# 认证配置
TEST_USERNAME=test@example.com
TEST_PASSWORD=password123
API_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# 测试Bot配置
TEST_BOT_ID=your-bot-id-here

# K6输出配置
K6_OUTPUT_DIR=tests/performance/results
```

加载环境变量：
```bash
# Linux/Mac
source .env

# Windows
set BASE_URL=http://localhost:8888
set API_TOKEN=your-token
```

### 3.3 准备测试数据

```bash
# 1. 创建测试账号
curl -X POST http://localhost:8888/api/passport/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test@example.com","password":"password123"}'

# 2. 创建测试Bot
curl -X POST http://localhost:8888/api/bots \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Performance Test Bot","description":"用于性能测试"}'

# 3. 创建测试知识库
curl -X POST http://localhost:8888/api/knowledge \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Performance Test Knowledge","description":"用于性能测试"}'
```

---

## 4. 测试场景详解

### 4.1 知识库API测试

**测试文件**: `knowledge_api_test.js`

**测试目标**:
- P95延迟 < 500ms
- 错误率 < 5%
- 吞吐量 > 100 req/s

**测试场景**:
```javascript
// 1. 知识库列表查询
GET /api/knowledge

// 2. 知识库详情查询
GET /api/knowledge/{id}

// 3. 知识库搜索
POST /api/knowledge/{id}/search
Body: { "query": "test", "top_k": 10 }

// 4. 文档列表查询
GET /api/knowledge/{id}/documents
```

**负载模式**:
```javascript
export const options = {
  stages: [
    { duration: '1m', target: 10 },   // 启动: 10用户
    { duration: '3m', target: 50 },   // 正常: 50用户
    { duration: '2m', target: 100 },  // 高负载: 100用户
    { duration: '2m', target: 50 },   // 降负载: 50用户
    { duration: '1m', target: 0 },    // 冷却
  ],
};
```

**运行测试**:
```bash
k6 run tests/performance/k6/knowledge_api_test.js
```

### 4.2 Bot API测试

**测试文件**: `bot_api_test.js`

**测试目标**:
- P95延迟 < 1000ms
- 错误率 < 2%
- 吞吐量 > 50 req/s

**测试场景**:
- Bot列表查询
- Bot创建 (30%概率)
- Bot详情查询
- Bot更新 (20%概率)
- Bot删除

**运行测试**:
```bash
k6 run tests/performance/k6/bot_api_test.js
```

### 4.3 会话负载测试

**测试文件**: `conversation_load_test.js`

**测试目标**:
- P95延迟 < 2000ms (Agent运行)
- 错误率 < 3%
- 并发支持 > 50 会话/秒

**测试场景**:
```javascript
// 1. 创建会话
POST /api/conversation
Body: { "bot_id": "xxx", "title": "Test" }

// 2. 发送消息
POST /api/conversation/{id}/messages
Body: { "content": "Hello", "role": "user" }

// 3. 查询历史
GET /api/conversation/{id}/messages?limit=10

// 4. 删除会话
DELETE /api/conversation/{id}
```

**运行测试**:
```bash
k6 run tests/performance/k6/conversation_load_test.js
```

### 4.4 权限检查测试

**测试文件**: `permission_test.js`

**测试目标**:
- P95延迟 < 100ms (权限检查应该非常快)
- 错误率 < 1%
- 吞吐量 > 500 req/s

**测试场景**:
- 数据权限检查 (Bot访问)
- 字段权限检查 (用户敏感字段)
- 角色权限验证 (管理功能)
- 资源访问控制 (跨租户隔离)

**运行测试**:
```bash
k6 run tests/performance/k6/permission_test.js
```

### 4.5 混合负载测试

**测试文件**: `mixed_load_test.js`

**测试目标**:
- 整体P95延迟 < 800ms
- 错误率 < 3%
- 系统吞吐量 > 200 req/s

**操作分布**:
| 操作 | 占比 | 说明 |
|------|------|------|
| 浏览Bot列表 | 30% | 最常见操作 |
| 创建/编辑Bot | 10% | 写入操作 |
| 查询知识库 | 25% | 读取操作 |
| 发送消息 | 25% | 复杂操作 |
| 检查权限 | 10% | 高频操作 |

**运行测试**:
```bash
k6 run tests/performance/k6/mixed_load_test.js
```

---

## 5. 测试执行

### 5.1 交互式执行

#### Linux/Mac
```bash
chmod +x tests/performance/run_k6_tests.sh
./tests/performance/run_k6_tests.sh
```

#### Windows
```batch
tests\performance\run_k6_tests.bat
```

**交互式菜单**:
```
=========================================
  ZKER Performance Testing Suite (K6)
  企业级性能测试套件
=========================================

[INFO] 选择测试模式:
1) 运行所有测试
2) 运行单个测试
3) 快速测试(仅核心功能)
4) 退出

请输入选择 [1-4]:
```

### 5.2 命令行执行

#### 运行单个测试
```bash
k6 run tests/performance/k6/knowledge_api_test.js
```

#### 指定输出格式
```bash
# JSON输出
k6 run --out json=test_result.json tests/performance/k6/bot_api_test.js

# InfluxDB输出 (用于Grafana可视化)
k6 run --out influxdb=http://localhost:8086/k6 tests/performance/k6/mixed_load_test.js
```

#### 指定VU数量和持续时间
```bash
k6 run --vus 100 --duration 10m tests/performance/k6/conversation_load_test.js
```

#### 分阶段加压
```bash
k6 run \
  --stage 2m:10 \
  --stage 5m:50 \
  --stage 3m:100 \
  --stage 2m:50 \
  tests/performance/k6/permission_test.js
```

### 5.3 K6 Cloud执行

```bash
# 登录K6 Cloud
k6 login cloud --token YOUR_K6_CLOUD_TOKEN

# 运行并上传结果
k6 cloud tests/performance/k6/mixed_load_test.js
```

---

## 6. 结果分析

### 6.1 终端输出

K6会在终端显示实时统计：

```
█Iteration: 49/1000 (4%)
  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓

checks.........................: 98.50% ✓ 2462  ✗ 38
data_received..................: 45 MB  750 kB/s
data_sent......................: 6.0 MB 100 kB/s
http_req_blocked...............: avg=1ms    min=0µs    med=1ms    max=50ms
http_req_connecting............: avg=1ms    min=0µs    med=1ms    max=40ms
http_req_duration..............: avg=250ms  min=10ms   med=200ms  max=1500ms
  { expected_response:true }...: avg=250ms  min=10ms   med=200ms  max=1500ms
http_req_failed................: 1.50%   ✓ 38     ✗ 2462
http_req_receiving.............: avg=10ms   min=10µs   med=5ms    max=100ms
http_req_sending...............: avg=1ms    min=10µs   med=1ms    max=20ms
http_req_tls_handshaking.......: avg=0s     min=0s     med=0s     max=0s
http_req_waiting...............: avg=239ms  min=10ms   med=190ms  max=1490ms
http_reqs......................: 2500    41.65/s
iteration_duration.............: avg=2.5s   min=1s     med=2s     max=5s
iterations.....................: 2500    41.65/s
vus............................: 50     min=10   max=100
vus_max........................: 100    min=100  max=100
```

**关键指标解释**:

| 指标 | 说明 |
|------|------|
| `checks` | 断言通过率 |
| `http_req_duration` | HTTP请求延迟 |
| `http_req_failed` | 请求失败率 |
| `http_reqs` | 总请求数 |
| `vus` | 当前VU数量 |
| `data_received` | 接收数据量 |
| `data_sent` | 发送数据量 |

### 6.2 JSON结果分析

```bash
# 使用jq分析结果
jq '.metrics.http_req_duration.values' tests/performance/results/test_20250103.json

# 提取P95延迟
jq '.metrics.http_req_duration.values."p(95)"' test.json

# 提取错误率
jq '.metrics.http_req_failed.values.rate' test.json
```

### 6.3 Grafana可视化

**步骤**:

1. 打开Grafana: http://localhost:3000
2. 导入Dashboard: "ZKER Performance Overview"
3. 选择时间范围
4. 查看实时指标

**关键Dashboard**:
- API Request Rate (请求率)
- API Latency P95 (延迟)
- API Error Rate (错误率)
- Cache Hit Rate (缓存命中率)
- Database Query Duration (数据库查询延迟)

---

## 7. 性能基线

### 7.1 API性能基线

| API类别 | P50 (ms) | P95 (ms) | P99 (ms) | QPS | 错误率 |
|---------|---------|---------|---------|-----|-------|
| **知识库列表** | < 100 | < 300 | < 500 | > 200 | < 1% |
| **知识库搜索** | < 200 | < 500 | < 1000 | > 100 | < 2% |
| **Bot列表** | < 100 | < 300 | < 500 | > 200 | < 1% |
| **Bot创建** | < 300 | < 800 | < 1500 | > 50 | < 2% |
| **发送消息** | < 500 | < 1500 | < 3000 | > 50 | < 3% |
| **权限检查** | < 20 | < 50 | < 100 | > 500 | < 0.5% |

### 7.2 系统资源基线

| 资源 | 100 QPS | 500 QPS | 1000 QPS | 5000 QPS |
|------|---------|---------|----------|----------|
| **CPU使用率** | < 10% | < 30% | < 50% | < 80% |
| **内存使用** | < 2GB | < 4GB | < 6GB | < 16GB |
| **数据库连接** | < 20 | < 50 | < 100 | < 180 |
| **Redis连接** | < 10 | < 20 | < 30 | < 50 |

### 7.3 更新性能基线

当性能基线需要更新时：

1. 在稳定环境运行测试
2. 收集3次测试结果
3. 取平均值作为新基线
4. 更新本文件
5. 通知团队

---

## 8. 故障排查

### 8.1 常见问题

#### 问题1: 连接被拒绝

**错误**:
```
ERRO[0000] GoError: dial tcp 127.0.0.1:8888: connect: connection refused
```

**解决**:
```bash
# 检查服务是否运行
curl http://localhost:8888/health

# 启动服务
make server
```

#### 问题2: 认证失败

**错误**:
```
✗ login successful
  status: 401
```

**解决**:
```bash
# 检查Token是否有效
export API_TOKEN=$(curl -X POST http://localhost:8888/api/passport/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test@example.com","password":"password123"}' \
  | jq -r '.data.token')
```

#### 问题3: 超时

**错误**:
```
context deadline exceeded
```

**解决**:
```javascript
// 增加超时时间
export const options = {
  thresholds: {
    http_req_duration: ['p(95)<3000'], // 从1000ms增加到3000ms
  },
};
```

### 8.2 性能问题定位

#### 高延迟定位

**步骤**:
1. 查看Grafana Dashboard，确认慢端点
2. 查看Prometheus慢查询日志
3. 检查数据库索引
4. 检查缓存命中率

**命令**:
```bash
# 查看慢查询
curl http://localhost:9090/api/v1/query?query=db_query_duration_seconds_bucket

# 查看缓存命中率
curl http://localhost:9090/api/v1/query?query=cache_hit_total/(cache_hit_total+cache_miss_total)
```

#### 高错误率定位

**步骤**:
1. 查看日志文件
2. 查看错误码分布
3. 检查依赖服务状态

**命令**:
```bash
# 查看错误码分布
curl http://localhost:9090/api/v1/query?query=sum(rate(error_total[5m]))%20by%20(error_code)

# 查看数据库连接
curl http://localhost:9090/api/v1/query?query=db_connections_active
```

---

## 9. CI/CD集成

### 9.1 GitHub Actions配置

创建`.github/workflows/performance.yml`:

```yaml
name: Performance Tests

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]
  schedule:
    - cron: '0 2 * * *'  # 每天凌晨2点

jobs:
  performance:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Install K6
        run: |
          curl https://github.com/grafana/k6/releases/download/v0.45.0/k6-v0.45.0-linux-amd64.tar.gz -L | tar xvz
          sudo mv k6-v0.45.0-linux-amd64/k6 /usr/local/bin/

      - name: Start services
        run: |
          docker-compose -f docker-compose.monitoring.yml up -d
          make server

      - name: Wait for services
        run: |
          for i in {1..30}; do
            if curl -s http://localhost:8888/health > /dev/null; then
              echo "Services ready"
              break
            fi
            echo "Waiting for services... ($i/30)"
            sleep 2
          done

      - name: Run performance tests
        run: |
          export BASE_URL=http://localhost:8888
          export API_TOKEN=${{ secrets.API_TOKEN }}
          ./tests/performance/run_k6_tests.sh
        continue-on-error: true

      - name: Upload results
        uses: actions/upload-artifact@v3
        with:
          name: performance-results
          path: tests/performance/results/*.json

      - name: Publish results to K6 Cloud
        if: github.event_name == 'schedule'
        run: |
          k6 login cloud --token ${{ secrets.K6_CLOUD_TOKEN }}
          k6 cloud tests/performance/k6/mixed_load_test.js
```

### 9.2 性能回归检测

在PR中检测性能回归：

```yaml
- name: Performance regression check
  run: |
    # 获取当前PR的性能数据
    CURRENT_P95=$(jq '.metrics.http_req_duration.values."p(95)"' tests/performance/results/current.json)

    # 获取main分支的性能基线
    BASELINE_P95=$(jq '.metrics.http_req_duration.values."p(95)"' tests/performance/results/baseline.json)

    # 计算性能退化
    REGRESSION=$(echo "scale=2; ($CURRENT_P95 - $BASELINE_P95) / $BASELINE_P95 * 100" | bc)

    # 检查是否超过阈值(10%)
    if (( $(echo "$REGRESSION > 10" | bc -l) )); then
      echo "❌ 性能退化: ${REGRESSION}%"
      exit 1
    else
      echo "✅ 性能正常: 退化 ${REGRESSION}%"
    fi
```

---

## 10. 最佳实践

### 10.1 测试设计原则

1. **真实性**: 模拟真实用户行为
2. **独立性**: 测试之间互不影响
3. **可重复**: 结果可重现
4. **可维护**: 代码清晰易读

### 10.2 脚本编写技巧

#### 1. 使用自定义指标

```javascript
const botCreateCount = new Counter('bot_create_total');
const botOperationLatency = new Trend('bot_operation_latency');

botCreateCount.inc();
botOperationLatency.add(duration);
```

#### 2. 合理的思考时间

```javascript
// 随机等待，模拟真实用户
sleep(Math.random() * 3 + 2); // 2-5秒
```

#### 3. 分组测试

```javascript
group('Bot Operations', () => {
  // 相关的操作放在一起
  createBot();
  updateBot();
  deleteBot();
});
```

#### 4. 错误处理

```javascript
const success = check(res, {
  'status 200 or 201': (r) => r.status === 200 || r.status === 201,
});

if (!success) {
  console.error('Failed to create bot');
  return null;
}
```

### 10.3 性能优化建议

1. **减少不必要的数据传输**
2. **使用连接池**
3. **启用缓存**
4. **优化数据库查询**
5. **异步处理**

### 10.4 持续改进

- 定期审查性能基线
- 识别性能瓶颈
- 优化慢查询
- 监控生产环境

---

## 附录

### A. 参考资源

- [K6官方文档](https://k6.io/docs/)
- [K6示例脚本](https://k6.io/docs/examples/)
- [企业级开发规范手册](./ZKER-企业级开发规范手册_v1.0.md)
- [Grafana监控运维手册](./ZKER-Grafana监控运维手册_v1.0.md)

### B. 联系方式

- **性能测试团队**: perf-team@coze-studio.com
- **问题反馈**: [GitHub Issues](https://github.com/coze-dev/coze-studio/issues)

### C. 更新日志

| 版本 | 日期 | 更新内容 |
|------|------|---------|
| v1.0.0 | 2025-01-03 | 初始版本 |

---

**版本**: v1.0.0 | **作者**: 研发B | **日期**: 2025-01-03
