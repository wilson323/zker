# 组织中心性能测试套件

## 📖 概述

完整的组织中心性能测试解决方案，使用K6和JMeter实现多场景性能测试。

**版本**: v1.0
**最后更新**: 2025-01-01
**维护者**: 研发B（后端工程师）

---

## 🎯 测试覆盖

### API端点覆盖

| API端点 | 测试场景 | 性能目标 |
|---------|----------|----------|
| GET /api/organizations | 组织列表查询 | p95 < 200ms |
| GET /api/organizations/:id | 组织详情查询 | p95 < 100ms |
| POST /api/organizations | 创建组织 | p95 < 300ms |
| PUT /api/organizations/:id | 更新组织 | p95 < 250ms |
| GET /api/organizations/tree | 组织树查询 | p95 < 300ms |
| GET /api/org/departments | 部门列表查询 | p95 < 200ms |
| GET /api/org/employees | 员工列表查询 | p95 < 200ms |
| GET /api/org/employees/search | 员工搜索 | p95 < 300ms |
| GET /api/directory/organization | 通讯录查询 | p95 < 500ms |

### 测试场景

| 场景 | 描述 | 并发用户 | 持续时间 |
|------|------|----------|----------|
| **基准测试** | 正常负载下性能 | 100 | 5分钟 |
| **峰值测试** | 突发流量冲击 | 1000 | 3分钟 |
| **压力测试** | 找性能瓶颈 | 500 | 30分钟 |
| **耐久测试** | 长期稳定性 | 200 | 2小时 |
| **可扩展性测试** | 线性扩展能力 | 10-1000 | 阶梯式 |

---

## 📁 文件结构

```
backend/tests/performance/
├── org_load_test.k6.js           # K6性能测试脚本 (~600行)
├── org_jmeter_test.jmx           # JMeter测试计划 (~800行)
├── generate_org_testdata.go       # 测试数据生成器 (~400行)
├── run-org-performance-tests.sh   # 测试执行脚本
├── plot_results.py                # Python可视化工具
├── ORG_PERFORMANCE_BASELINE.md    # 性能基准文档
├── ORG_PERFORMANCE_README.md      # 本文件
└── results/                       # 测试结果目录
    ├── *.json                     # K6 JSON结果
    ├── *.jtl                      # JMeter结果
    └── *.png                      # 性能图表
```

---

## 🚀 快速开始

### 1. 环境准备

**安装依赖**:
```bash
# 安装K6
curl https://github.com/grafana/k6/releases/download/v0.47.0/k6-v0.47.0-linux-amd64.tar.gz -L | tar xvz
sudo mv k6-0.47.0-linux-amd64/k6 /usr/local/bin/

# 安装Go (用于生成测试数据)
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz

# 安装Python (用于生成图表)
sudo apt install python3-pip
pip3 install matplotlib
```

**启动MySQL**:
```bash
cd docker
docker compose up -d mysql
```

**启动后端服务**:
```bash
cd backend
make server
```

### 2. 生成测试数据

```bash
cd backend/tests/performance

# 生成小规模数据 (1租户, 10组织, 100员工) - 快速测试
go run generate_org_testdata.go -size=small

# 生成中等规模数据 (10租户, 100组织, 10000员工) - 标准测试
go run generate_org_testdata.go -size=medium

# 生成大规模数据 (100租户, 10000组织, 100000员工) - 压力测试
go run generate_org_testdata.go -size=large

# 自定义规模
go run generate_org_testdata.go \
  -tenants=5 \
  -orgs=50 \
  -depts=500 \
  -positions=50 \
  -emps=5000
```

### 3. 执行性能测试

**方法1: 使用K6** (推荐)
```bash
cd backend/tests/performance

# 基准测试 (100用户, 5分钟)
k6 run --env BASE_URL=http://localhost:8080 org_load_test.k6.js

# 峰值测试 (1000用户, 短时间冲击)
k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=spike org_load_test.k6.js

# 压力测试 (500用户, 30分钟)
k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=stress org_load_test.k6.js

# 耐久测试 (200用户, 2小时)
k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=endurance org_load_test.k6.js

# 可扩展性测试 (阶梯式负载)
k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=scalability org_load_test.k6.js
```

**方法2: 使用自动化脚本** (推荐)
```bash
cd backend/tests/performance

# 完整测试（包括数据生成）
bash run-org-performance-tests.sh --full

# 仅执行测试（跳过数据生成）
bash run-org-performance-tests.sh --skip-data

# 仅执行基准测试
bash run-org-performance-tests.sh --test-type=base

# 自定义URL
bash run-org-performance-tests.sh --base-url=http://192.168.1.100:8080
```

**方法3: 使用JMeter**
```bash
# 启动JMeter GUI
jmeter -t org_jmeter_test.jmx

# 命令行模式（无GUI）
jmeter -n -t org_jmeter_test.jmx -l results.jtl -e -o report/

# 使用自定义参数
jmeter -n -t org_jmeter_test.jmx \
  -JBASE_URL=http://localhost:8080 \
  -JTENANT_ID=test-tenant-001 \
  -l results.jtl
```

---

## 📊 结果分析

### K6输出解读

```
✓ status is 200
✓ response time < 200ms
✓ has organizations array

checks.........................: 100.00% ✓ 60000      ✗ 0
data_received..................: 150 MB  500 kB/s
http_req_duration..............: avg=150ms   min=10ms   med=120ms  max=800ms
  { expected_response:true }...: avg=150ms   min=10ms   med=120ms  max=800ms
http_req_failed................: 0.05%  ✓ 30       ✗ 59970
http_reqs......................: 60000   200 req/s
```

**关键指标**:
- **http_req_duration**: 请求总耗时 (目标: p95 < 200ms)
- **http_req_failed**: 失败率 (目标: < 0.1%)
- **http_reqs**: 吞吐量 (目标: > 100 req/s)
- **checks**: 断言通过率 (目标: = 100%)

### 生成性能图表

```bash
cd backend/tests/performance

# 生成单个测试结果图表
python3 plot_results.py results/base_load_test.json

# 生成所有测试结果图表
python3 plot_results.py results/

# 生成HTML报告
python3 plot_results.py results/ --html

# 对比两个测试结果
python3 plot_results.py results/before.json results/after.json --compare
```

### 查看详细报告

**K6 HTML报告**:
```bash
# 生成HTML报告
k6 run --env BASE_URL=http://localhost:8080 \
  --out json=results/test.json \
  org_load_test.k6.js

# 使用k6-reporter生成HTML
k6-reporter html --input results/test.json --output results/report.html
```

**JMeter HTML报告**:
```bash
jmeter -n -t org_jmeter_test.jmx -l results.jtl -e -o report/
# 报告将生成在 report/ 目录
```

---

## 🎯 性能基线

### 小规模数据 (1租户, 10组织, 100员工)

| API | p50 | p95 | p99 | 吞吐量 |
|-----|-----|-----|-----|--------|
| GET /api/organizations | 50ms | 80ms | 120ms | 200 req/s |
| GET /api/org/employees | 80ms | 150ms | 200ms | 150 req/s |
| GET /api/org/employees/search | 100ms | 180ms | 250ms | 100 req/s |

### 中规模数据 (10租户, 1000组织, 10000员工)

| API | p50 | p95 | p99 | 吞吐量 |
|-----|-----|-----|-----|--------|
| GET /api/organizations | 80ms | 150ms | 250ms | 150 req/s |
| GET /api/org/employees | 120ms | 250ms | 400ms | 100 req/s |
| GET /api/org/employees/search | 150ms | 300ms | 500ms | 80 req/s |

**详细基线**: 参考 [ORG_PERFORMANCE_BASELINE.md](./ORG_PERFORMANCE_BASELINE.md)

---

## 🔍 性能问题诊断

### 常见性能问题

#### 1. 员工列表查询慢 (>500ms)

**诊断**:
```sql
-- 检查慢查询
SELECT * FROM mysql.slow_log
WHERE sql_text LIKE '%employees%'
ORDER BY query_time DESC
LIMIT 10;

-- 检查索引使用情况
EXPLAIN SELECT * FROM employees
WHERE tenant_id = 'xxx' AND status = 'active'
ORDER BY hire_date DESC
LIMIT 20;
```

**优化**:
```sql
-- 添加复合索引
CREATE INDEX idx_tenant_status_hire ON employees(tenant_id, status, hire_date);

-- 使用游标分页
SELECT emp_id, emp_name, emp_code
FROM employees
WHERE tenant_id = ? AND hire_date < ?
ORDER BY hire_date DESC
LIMIT 21;
```

#### 2. 组织树查询慢 (>600ms)

**优化**: 使用Redis缓存
```go
func (s *OrganizationService) GetTreeWithCache(ctx context.Context, tenantID string) (*TreeVO, error) {
    cacheKey := fmt.Sprintf("org:tree:%s", tenantID)

    // 尝试从缓存获取
    var tree TreeVO
    if err := s.cache.Get(ctx, cacheKey, &tree); err == nil {
        return &tree, nil
    }

    // 查询数据库
    tree, err := s.getTreeFromDB(ctx, tenantID)
    if err != nil {
        return nil, err
    }

    // 缓存结果（5分钟）
    s.cache.Set(ctx, cacheKey, tree, 5*time.Minute)

    return &tree, nil
}
```

#### 3. 员工搜索性能瓶颈

**优化**: 使用MySQL全文索引或Elasticsearch
```sql
-- MySQL全文索引
ALTER TABLE employees ADD FULLTEXT INDEX ft_emp_name (emp_name);
SELECT * FROM employees
WHERE MATCH(emp_name) AGAINST('张三' IN NATURAL LANGUAGE MODE)
LIMIT 20;
```

---

## 📈 持续集成

### GitHub Actions自动化测试

测试工作流: `.github/workflows/org-performance-test.yml`

**触发条件**:
- 每天凌晨2点定时执行
- PR创建或更新
- 手动触发

**执行流程**:
1. 启动MySQL容器
2. 生成测试数据
3. 启动后端服务
4. 执行性能测试
5. 生成测试报告
6. 检查性能回归

**查看结果**: GitHub Actions -> Summary

---

## 🛠️ 自定义测试

### 修改测试配置

编辑 `org_load_test.k6.js`:
```javascript
// 修改并发用户数
export const options = {
    stages: [
        { duration: '1m', target: 100 },  // 修改这里
        { duration: '5m', target: 100 },
    ],
};

// 修改性能阈值
thresholds: {
    http_req_duration: ['p(95)<200', 'p(99)<500'],  // 修改这里
    http_req_failed: ['rate<0.001'],
}
```

### 添加新的测试场景

```javascript
// 在 org_load_test.k6.js 中添加新函数
function scenario10_customTest(data) {
    const resp = http.get(
        `${BASE_URL}/api/custom-endpoint`,
        {
            headers: { 'X-Tenant-ID': data.tenant_id },
            tags: { name: 'CustomTest' },
        }
    );

    const success = check(resp, {
        'status is 200': (r) => r.status === 200,
        'response time < 100ms': (r) => r.timings.duration < 100,
    });

    errorRate.add(!success);
    responseTime.add(resp.timings.duration);
}
```

---

## 📚 相关文档

- [性能基准文档](./ORG_PERFORMANCE_BASELINE.md) - 详细性能基线和优化建议
- [企业级开发规范手册](../../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [研发B开发计划](../../../docs/企业级功能完善与统一性设计方案/研发B-后端工程师开发计划_v1.0.md)

---

## 🤝 贡献指南

### 添加新测试用例

1. 在 `org_load_test.k6.js` 中添加新场景函数
2. 更新 `ORG_PERFORMANCE_BASELINE.md` 中的性能目标
3. 更新本文档中的API覆盖表
4. 运行测试验证

### 报告性能问题

1. 记录测试环境和配置
2. 导出测试结果 (JSON/HTML)
3. 提供慢查询日志
4. 提交Issue到项目仓库

---

## 📞 支持

**问题反馈**: GitHub Issues
**技术支持**: 研发B（后端工程师）

---

**最后更新**: 2025-01-01
**版本**: v1.0
**状态**: ✅ 生产就绪
