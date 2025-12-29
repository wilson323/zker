# ZKER 端到端性能测试指南

**版本**: v1.0.0
**最后更新**: 2025-01-01
**工具**: K6 (https://k6.io/)

---

## 📋 目录

- [测试概述](#测试概述)
- [快速开始](#快速开始)
- [测试场景说明](#测试场景说明)
- [性能基线](#性能基线)
- [结果分析](#结果分析)
- [故障排查](#故障排查)
- [持续集成](#持续集成)

---

## 🎯 测试概述

### 测试目标

本测试套件旨在验证 ZKER 系统在以下方面的性能表现：

1. **API 响应时间**: P95 < 500ms, P99 < 1000ms
2. **系统稳定性**: 错误率 < 1%
3. **并发能力**: 支持 100 并发用户
4. **极限性能**: 1000 并发用户下可降级服务

### 测试覆盖

| 测试类别 | 测试场景 | 并发用户 | 持续时间 |
|---------|---------|---------|---------|
| **负载测试** | 基础功能测试 | 10 → 100 | 5 分钟 |
| **压力测试** | 极限压力测试 | 100 → 300 | 16 分钟 |
| **浸泡测试** | 长时间稳定性 | 50 | 1 小时 |
| **峰值测试** | 突发流量测试 | 100 → 1000 | 4 分钟 |

### 测试 API

1. **租户管理 API**:
   - GET `/api/tenants` - 获取租户列表
   - GET `/api/tenants/:id` - 获取租户详情

2. **配额管理 API**:
   - GET `/api/tenants/:id/quotas` - 获取配额状态
   - POST `/api/tenants/:id/quotas/check` - 检查配额

3. **Bot 服务 API**:
   - GET `/api/bots` - 列出 Bots
   - POST `/api/bots` - 创建 Bot

4. **权限检查 API**:
   - GET `/api/users/:id/roles` - 获取用户角色
   - POST `/api/permissions/check-data` - 检查数据权限

---

## 🚀 快速开始

### 1. 安装 K6

```bash
# macOS
brew install k6

# Linux
sudo apt-get install k6

# Windows
choco install k6

# Docker
docker pull grafana/k6
```

### 2. 配置环境变量

```bash
# 设置目标环境
export BASE_URL="http://localhost:8888"
export TENANT_ID="test_tenant_001"
export USER_ID="test_user_001"
```

### 3. 运行测试

```bash
# 基础负载测试
k6 run performance_test.js

# 指定输出格式
k6 run --out json=output.json performance_test.js

# 生成 HTML 报告
k6 run --out json=output.json performance_test.js
k6-to-html output.json > report.html
```

### 4. Docker 运行

```bash
docker run --rm \
  -e BASE_URL="http://host.docker.internal:8888" \
  -e TENANT_ID="test_tenant_001" \
  -e USER_ID="test_user_001" \
  -v $(pwd):/scripts \
  grafana/k6 run /scripts/performance_test.js
```

---

## 📊 测试场景说明

### 1. 基础负载测试 (Load Test)

**目标**: 验证系统在正常负载下的性能表现

**配置**:
```javascript
stages: [
  { duration: '1m', target: 10 },   // 启动阶段
  { duration: '2m', target: 50 },   // 加压阶段
  { duration: '1m', target: 100 },  // 峰值阶段
  { duration: '1m', target: 0 },    // 降压阶段
]
```

**性能基线**:
- P95 响应时间: < 500ms
- P99 响应时间: < 1000ms
- 错误率: < 1%

### 2. 压力测试 (Stress Test)

**目标**: 找到系统性能瓶颈和极限

**配置**:
```javascript
stages: [
  { duration: '2m', target: 100 },  // 2分钟爬坡到100用户
  { duration: '5m', target: 100 },  // 持续5分钟100用户
  { duration: '2m', target: 200 },  // 2分钟爬坡到200用户
  { duration: '5m', target: 200 },  // 持续5分钟200用户
  { duration: '2m', target: 300 },  // 2分钟爬坡到300用户
  { duration: '5m', target: 300 },  // 持续5分钟300用户
  { duration: '2m', target: 0 },    // 降压到0
]
```

**性能基线** (压力测试阈值更宽松):
- P95 响应时间: < 3000ms
- 错误率: < 10%

**观察指标**:
- 系统何时开始出现错误
- 响应时间何时急剧上升
- CPU/内存使用率峰值

### 3. 浸泡测试 (Soak Test)

**目标**: 验证系统长时间运行的稳定性

**配置**:
```javascript
stages: [
  { duration: '5m', target: 10 },   // 预热5分钟
  { duration: '1h', target: 50 },   // 持续1小时50用户
  { duration: '5m', target: 10 },   // 降温5分钟
  { duration: '5m', target: 0 },    // 降压到0
]
```

**性能基线** (稳定性测试要求更严格):
- P95 响应时间: < 1000ms
- 错误率: < 1%

**观察指标**:
- 内存泄漏 (内存使用持续上升)
- 连接泄漏 (数据库连接数持续上升)
- 性能退化 (响应时间逐渐变慢)

### 4. 峰值测试 (Spike Test)

**目标**: 验证系统应对突发流量的能力

**配置**:
```javascript
stages: [
  { duration: '1m', target: 100 },   // 1分钟爬坡到100用户
  { duration: '1m', target: 1000 },  // 突增到1000用户
  { duration: '1m', target: 1000 },  // 保持1000用户
  { duration: '1m', target: 0 },     // 降压到0
]
```

**性能基线** (峰值测试允许降级):
- P95 响应时间: < 5000ms
- 错误率: < 20%

**观察指标**:
- 系统是否崩溃
- 自动扩缩容是否生效
- 限流机制是否触发

---

## 🎯 性能基线

### API 性能基线

| API 端点 | 方法 | P95 (ms) | P99 (ms) | 说明 |
|---------|------|---------|---------|------|
| `/api/tenants` | GET | 300 | 500 | 租户列表查询 |
| `/api/tenants/:id` | GET | 200 | 300 | 租户详情查询 |
| `/api/tenants/:id/quotas` | GET | 500 | 1000 | 配额状态查询 |
| `/api/tenants/:id/quotas/check` | POST | 500 | 1000 | 配额检查（含验证） |
| `/api/bots` | GET | 1000 | 2000 | Bot 列表（可能含分页） |
| `/api/bots` | POST | 2000 | 3000 | Bot 创建（含权限检查） |
| `/api/users/:id/roles` | GET | 500 | 1000 | 用户角色查询 |
| `/api/permissions/check-data` | POST | 200 | 500 | 权限检查（应快速） |

### 业务指标基线

| 指标类别 | 指标名称 | 基线值 | 说明 |
|---------|---------|-------|------|
| **并发能力** | 最大并发用户 | 100 | 正常运行 |
| **并发能力** | 极限并发用户 | 1000 | 可降级服务 |
| **错误率** | 正常负载错误率 | < 1% | 100 用户以下 |
| **错误率** | 峰值负载错误率 | < 20% | 1000 用户 |
| **吞吐量** | 请求/秒 (RPS) | > 100 | 100 用户并发 |
| **资源使用** | CPU 使用率 | < 80% | 100 用户并发 |
| **资源使用** | 内存使用率 | < 85% | 100 用户并发 |
| **资源使用** | 数据库连接数 | < 80% | 连接池利用率 |

---

## 📈 结果分析

### 1. 终端输出解读

```
scenarios: (100.00%) 1 scenario, 100 max VUs

     ✓ Tenant Management APIs - List Tenants
     ✓ Tenant Management APIs - Response Time < 500ms
     ✓ Quota Management APIs - Get Quotas
     ✓ Bot Service APIs - List Bots
     ✓ Permission Check APIs - Get User Roles

     checks.........................: 99.50% ✓ 19999/ 20098
     data_received..................: 15 MB  250 kB/s
     data_sent......................: 2.1 MB 35 kB/s
     http_req_blocked...............: avg=1.2ms    min=1µs      med=4µs      max=500ms    p(95)=10µs     p(99)=50µs
     http_req_connecting............: avg=800µs    min=0s       med=0s       max=400ms    p(95)=0s       p(99)=2ms
     http_req_duration..............: avg=350ms    min=10ms     med=200ms    max=5s       p(95)=800ms   p(99)=2s
       { expected_response:true }...: avg=340ms    min=10ms     med=190ms    max=4s       p(95)=700ms   p(99)=1.8s
     http_req_failed................: 0.50%  ✓ 100   ✗ 19998
     http_req_receiving.............: avg=15ms     min=10µs     med=100µs    max=2s       p(95)=50ms    p(99)=200ms
     http_req_sending...............: avg=5ms      min=5µs      med=50µs     max=500ms    p(95)=20ms    p(99)=100ms
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s       p(95)=0s       p(99)=0s
     http_req_waiting...............: avg=330ms    min=10ms     med=190ms    max=4.8s     p(95)=750ms   p(99)=1.9s
     http_reqs......................: 20098  334.931439/s
     iteration_duration.............: avg=5.5s     min=4.9s     med=5.3s     max=10s      p(95)=8s       p(99)=9.5s
     iterations.....................: 3333   55.549776/s
     vus............................: 10     min=10   max=100
     vus_max........................: 100    min=100  max=100
```

**关键指标说明**:

- `checks`: 检查通过率，应 ≥ 99%
- `http_req_duration`: HTTP 请求总时长
  - `p(95)`: 95% 请求的响应时间，应 < 500ms（正常负载）
  - `p(99)`: 99% 请求的响应时间，应 < 1000ms
- `http_req_failed`: 请求失败率，应 < 1%
- `vus`: 虚拟用户数
- `iterations`: 迭代次数（场景执行次数）

### 2. 生成 HTML 报告

```bash
# 1. 运行测试并输出 JSON
k6 run --out json=output.json performance_test.js

# 2. 安装 k6-to-html
npm install -g k6-to-html

# 3. 生成 HTML 报告
k6-to-html output.json > report.html

# 4. 打开报告
open report.html  # macOS
xdg-open report.html  # Linux
start report.html  # Windows
```

### 3. 性能瓶颈分析

**常见性能瓶颈及解决方案**:

| 症状 | 可能原因 | 解决方案 |
|------|---------|---------|
| P95 响应时间 > 2秒 | 数据库慢查询 | 优化索引、重构查询 |
| 错误率突然上升 | 连接池耗尽 | 增加连接池大小 |
| 内存使用持续上升 | 内存泄漏 | 排查 goroutine 泄漏 |
| CPU 使用率 > 90% | 计算密集操作 | 优化算法、异步处理 |
| 某个 API 特别慢 | N+1 查询 | 使用 Preload、批量查询 |

---

## 🔧 故障排查

### 问题 1: 无法连接到服务器

**错误**: `Error: connect ECONNREFUSED 127.0.0.1:8888`

**解决方案**:
```bash
# 1. 检查服务是否运行
curl http://localhost:8888/health

# 2. 检查防火墙
sudo ufw status

# 3. 检查 BASE_URL 环境变量
echo $BASE_URL
```

### 问题 2: 大量 401/403 错误

**错误**: `status: 401` 或 `status: 403`

**解决方案**:
```bash
# 1. 检查 TENANT_ID 和 USER_ID
export TENANT_ID="valid_tenant_id"
export USER_ID="valid_user_id"

# 2. 确认租户和用户存在
curl -H "X-Tenant-ID: $TENANT_ID" http://localhost:8888/api/tenants/$TENANT_ID
```

### 问题 3: 内存不足

**错误**: `JavaScript heap out of memory`

**解决方案**:
```bash
# 增加 Node.js 内存限制
export NODE_OPTIONS="--max-old-space-size=4096"
k6 run performance_test.js

# 或减少并发用户数
# 修改 performance_test.js 中的 vus 和 stages
```

### 问题 4: 数据库连接池耗尽

**错误**: `Error: too many connections`

**解决方案**:
```bash
# 1. 检查数据库连接数
mysql> SHOW STATUS LIKE 'Threads_connected';

# 2. 增加连接池大小
# 修改 backend/conf/config.yaml
database:
  max_open_conns: 200
  max_idle_conns: 50

# 3. 减少测试并发用户数
```

---

## 🔄 持续集成

### GitHub Actions 集成

创建 `.github/workflows/performance-test.yml`:

```yaml
name: Performance Test

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]
  schedule:
    # 每天凌晨 2 点运行
    - cron: '0 2 * * *'

jobs:
  performance-test:
    runs-on: ubuntu-latest

    services:
      mysql:
        image: mysql:8.4
        env:
          MYSQL_ROOT_PASSWORD: password
          MYSQL_DATABASE: coze_studio_test
        ports:
          - 3306:3306
        options: --health-cmd="mysqladmin ping" --health-interval=10s --health-timeout=5s --health-retries=3

      redis:
        image: redis:8.0
        ports:
          - 6379:6379
        options: --health-cmd="redis-cli ping" --health-interval=10s --health-timeout=5s --health-retries=3

    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24.0'

      - name: Install K6
        run: |
          sudo gpg -k
          sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
          echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
          sudo apt-get update
          sudo apt-get install k6

      - name: Install dependencies
        run: |
          cd backend
          go mod download

      - name: Build backend
        run: |
          cd backend
          go build -o coze-studio ./cmd/coze-studio

      - name: Start backend server
        run: |
          cd backend
          ./coze-studio --config ../conf/config.yaml &
          sleep 10

      - name: Run performance test
        env:
          BASE_URL: http://localhost:8888
          TENANT_ID: test_tenant_001
          USER_ID: test_user_001
        run: |
          k6 run --out json=output.json tests/performance/performance_test.js

      - name: Generate HTML report
        run: |
          npm install -g k6-to-html
          k6-to-html output.json > report.html

      - name: Upload report
        uses: actions/upload-artifact@v3
        with:
          name: performance-report
          path: report.html

      - name: Comment PR with results
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v6
        with:
          script: |
            const fs = require('fs');
            const report = fs.readFileSync('report.html', 'utf8');
            // 解析报告数据并评论到 PR
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: '## 性能测试报告\n\n测试已完成，详细报告请查看 Artifacts。'
            });
```

### 性能阈值检查

在 K6 脚本中配置阈值，测试失败时 CI 会报错：

```javascript
export const options = {
  thresholds: {
    // P95 响应时间必须 < 2秒
    http_req_duration: ['p(95)<2000'],

    // 错误率必须 < 5%
    http_req_failed: ['rate<0.05'],

    // 租户 API 成功率必须 > 95%
    'tenant_api_success_rate': ['rate>0.95'],
  },
};
```

---

## 📚 相关文档

- [企业级开发规范手册](../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [性能基线文档](../../docs/企业级功能完善与统一性设计方案/性能基线文档.md) (待创建)
- [故障排查手册](../../docs/企业级功能完善与统一性设计方案/ZKER-故障排查手册_v1.0.md)
- [K6 官方文档](https://k6.io/docs/)

---

## 📞 支持

如有问题，请联系：
- 性能测试团队：perf-team@coze-studio.com
- 后端团队：backend-team@coze-studio.com

---

**更新日志**：

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2025-01-01 | v1.0.0 | 初始版本，完整的K6性能测试套件 | Claude AI |
