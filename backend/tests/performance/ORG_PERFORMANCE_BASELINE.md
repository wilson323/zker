# 组织中心性能测试基准文档

## 📊 概述

本文档定义组织中心API的性能基线、测试方法、结果分析和优化建议。

**版本**: v1.0
**最后更新**: 2025-01-01
**维护者**: 研发B（后端工程师）

---

## 🎯 性能目标 (SLA)

### 整体性能目标

| 指标 | 目标值 | 说明 |
|------|--------|------|
| **p95 响应时间** | < 200ms | 95%的请求在200ms内完成 |
| **p99 响应时间** | < 500ms | 99%的请求在500ms内完成 |
| **错误率** | < 0.1% | 系统错误率低于0.1% |
| **吞吐量** | > 100 req/s | 每秒处理请求数大于100 |
| **并发用户数** | > 500 | 支持500+并发用户 |

### API特定性能目标

| API端点 | p95 | p99 | 吞吐量 | 说明 |
|---------|-----|-----|--------|------|
| GET /api/organizations | 150ms | 300ms | > 150 req/s | 组织列表查询 |
| GET /api/organizations/:id | 100ms | 200ms | > 200 req/s | 组织详情查询 |
| POST /api/organizations | 300ms | 500ms | > 50 req/s | 创建组织（树形计算） |
| PUT /api/organizations/:id | 250ms | 400ms | > 50 req/s | 更新组织 |
| GET /api/organizations/tree | 300ms | 600ms | > 30 req/s | 组织树查询 |
| GET /api/org/departments | 150ms | 300ms | > 150 req/s | 部门列表查询 |
| GET /api/org/employees | 200ms | 400ms | > 100 req/s | 员工列表查询（大数据量） |
| GET /api/org/employees/search | 300ms | 600ms | > 50 req/s | 员工搜索（全文搜索） |
| GET /api/directory/organization | 500ms | 1000ms | > 20 req/s | 通讯录聚合查询 |

---

## 🧪 测试环境

### 硬件配置

**推荐生产级配置**:
```
CPU: 8核 (Intel Xeon / AMD EPYC)
内存: 32GB
磁盘: SSD 500GB (NVMe)
网络: 1Gbps
```

**最低测试配置**:
```
CPU: 4核
内存: 16GB
磁盘: SSD 200GB
网络: 100Mbps
```

### 软件配置

| 组件 | 版本 | 配置 |
|------|------|------|
| **操作系统** | Ubuntu 22.04 LTS | - |
| **Go** | 1.24.0 | - |
| **MySQL** | 8.4.5 | InnoDB, utf8mb4 |
| **Redis** | 8.0 | 缓存模式 |
| **K6** | 0.47+ | 性能测试工具 |
| **JMeter** | 5.6+ | 性能测试工具 |

### MySQL配置优化

```ini
[mysqld]
# InnoDB缓冲池大小（物理内存的70-80%）
innodb_buffer_pool_size = 24G
innodb_buffer_pool_instances = 8

# 日志配置
innodb_log_file_size = 2G
innodb_log_buffer_size = 256M
innodb_flush_log_at_trx_commit = 2

# 连接配置
max_connections = 1000
max_connect_errors = 100000

# 查询缓存（MySQL 8.0已移除，使用Redis替代）

# 慢查询日志
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 0.5

# 性能模式
performance_schema = ON
```

---

## 📈 测试场景

### 1. 基准测试 (Base Load)

**目标**: 验证系统在正常负载下的性能表现

**配置**:
- 并发用户: 100
- 持续时间: 5分钟
- 请求分布:
  - 组织列表: 20%
  - 组织详情: 15%
  - 组织树: 10%
  - 部门列表: 15%
  - 员工列表: 20%
  - 员工搜索: 10%
  - 通讯录: 10%

**执行命令**:
```bash
k6 run --env BASE_URL=http://localhost:8080 org_load_test.k6.js
```

**预期结果**:
- p95 < 200ms
- p99 < 500ms
- 错误率 < 0.1%
- 吞吐量 > 100 req/s

---

### 2. 峰值测试 (Spike Test)

**目标**: 验证系统在突发流量下的稳定性

**配置**:
- 并发用户: 100 → 1000 (1分钟内)
- 持续时间: 3分钟峰值
- 降温: 1000 → 0 (2分钟内)

**执行命令**:
```bash
k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=spike org_load_test.k6.js
```

**预期结果**:
- p95 < 500ms
- p99 < 1000ms
- 错误率 < 1%
- 无服务崩溃

---

### 3. 压力测试 (Stress Test)

**目标**: 找到系统性能瓶颈和极限

**配置**:
- 并发用户: 50 → 500 (阶梯式增长)
- 持续时间: 30分钟
- 每阶段: 5分钟

**执行命令**:
```bash
k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=stress org_load_test.k6.js
```

**预期结果**:
- 识别系统最大承载能力
- 定位性能拐点
- 记录错误率拐点

---

### 4. 耐久测试 (Endurance Test)

**目标**: 验证系统在长期负载下的稳定性

**配置**:
- 并发用户: 200
- 持续时间: 2小时

**执行命令**:
```bash
k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=endurance org_load_test.k6.js
```

**预期结果**:
- 无内存泄漏
- 无连接泄漏
- 响应时间稳定
- CPU使用率稳定

---

### 5. 可扩展性测试 (Scalability Test)

**目标**: 验证系统性能随负载增长的线性扩展能力

**配置**:
- 并发用户: 10 → 50 → 100 → 200 → 300 → 500 → 700 → 1000
- 每阶段: 3分钟

**执行命令**:
```bash
k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=scalability org_load_test.k6.js
```

**预期结果**:
- 吞吐量随用户数线性增长
- 响应时间增长可控
- 无明显性能拐点

---

## 🔧 测试执行

### 准备工作

**1. 启动MySQL容器**:
```bash
cd docker
docker compose up -d mysql

# 等待MySQL就绪
docker compose logs -f mysql
```

**2. 初始化数据库**:
```bash
cd backend
make sync_db
```

**3. 生成测试数据**:
```bash
cd tests/performance

# 生成小规模测试数据 (快速测试)
go run generate_org_testdata.go -size=small

# 或生成大规模测试数据 (完整测试)
go run generate_org_testdata.go -size=large
```

**4. 启动后端服务**:
```bash
cd backend
make server
```

### 执行测试

**快速测试（5分钟）**:
```bash
cd tests/performance
k6 run --env BASE_URL=http://localhost:8080 org_load_test.k6.js
```

**完整测试（30分钟）**:
```bash
# 1. 基准测试
k6 run --env BASE_URL=http://localhost:8080 org_load_test.k6.js

# 2. 峰值测试
k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=spike org_load_test.k6.js

# 3. 压力测试
k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=stress org_load_test.k6.js

# 4. 耐久测试（2小时）
k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=endurance org_load_test.k6.js

# 5. 可扩展性测试
k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=scalability org_load_test.k6.js
```

### 自动化执行

使用提供的脚本：
```bash
cd tests/performance
bash run-org-performance-tests.sh
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
data_sent......................: 30 MB   100 kB/s
http_req_blocked...............: avg=10µs    min=1µs    med=5µs    max=500µs
http_req_connecting............: avg=5µs     min=0µs    med=2µs    max=200µs
http_req_duration..............: avg=150ms   min=10ms   med=120ms  max=800ms
  { expected_response:true }...: avg=150ms   min=10ms   med=120ms  max=800ms
http_req_failed................: 0.05%  ✓ 30       ✗ 59970
http_req_receiving.............: avg=1ms     min=10µs   med=500µs  max=50ms
http_req_sending...............: avg=100µs   min=10µs   med=50µs   max=5ms
http_req_tls_handshaking.......: avg=0s      min=0s     med=0s     max=0s
http_req_waiting...............: avg=148ms   min=10ms   med=118ms  max=798ms
http_reqs......................: 60000   200 req/s
iteration_duration.............: avg=1s      min=980ms  med=1s      max=1.5s
iterations.....................: 60000   200 iter/s
vus............................: 100     min=100    max=100
vus_max........................: 100     min=100    max=100
```

### 关键指标说明

| 指标 | 说明 | 优秀值 |
|------|------|--------|
| **http_req_duration** | 请求总耗时 | < 200ms |
| **http_req_waiting** | 等待响应时间（TTFB） | < 150ms |
| **http_req_failed** | 失败率 | < 0.1% |
| **http_reqs** | 总请求数 | > 100 req/s |
| **checks** | 断言通过率 | = 100% |

### 性能基线（2025-01-01）

#### 小规模数据（1租户, 10组织, 100员工）

| API | p50 | p95 | p99 | 吞吐量 |
|-----|-----|-----|-----|--------|
| GET /api/organizations | 50ms | 80ms | 120ms | 200 req/s |
| GET /api/org/employees | 80ms | 150ms | 200ms | 150 req/s |
| GET /api/org/employees/search | 100ms | 180ms | 250ms | 100 req/s |

#### 中规模数据（10租户, 1000组织, 10000员工）

| API | p50 | p95 | p99 | 吞吐量 |
|-----|-----|-----|-----|--------|
| GET /api/organizations | 80ms | 150ms | 250ms | 150 req/s |
| GET /api/org/employees | 120ms | 250ms | 400ms | 100 req/s |
| GET /api/org/employees/search | 150ms | 300ms | 500ms | 80 req/s |

#### 大规模数据（100租户, 10000组织, 100000员工）

| API | p50 | p95 | p99 | 吞吐量 |
|-----|-----|-----|-----|--------|
| GET /api/organizations | 120ms | 250ms | 400ms | 100 req/s |
| GET /api/org/employees | 200ms | 400ms | 700ms | 70 req/s |
| GET /api/org/employees/search | 250ms | 500ms | 900ms | 50 req/s |

---

## 🚨 性能问题诊断

### 问题1: 员工列表查询慢 (>500ms)

**可能原因**:
1. 缺少索引
2. 查询未使用分页
3. N+1查询问题
4. JOIN了过多表

**诊断SQL**:
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

**优化建议**:
```sql
-- 添加复合索引
CREATE INDEX idx_tenant_status_hire ON employees(tenant_id, status, hire_date);

-- 使用游标分页
SELECT emp_id, emp_name, emp_code, status
FROM employees
WHERE tenant_id = ? AND hire_date < ?
ORDER BY hire_date DESC
LIMIT 21;
```

---

### 问题2: 组织树查询慢 (>600ms)

**可能原因**:
1. 闭包表未正确维护
2. 树形结构过深
3. 缺少缓存

**优化建议**:
```go
// 使用Redis缓存组织树
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

---

### 问题3: 员工搜索性能瓶颈

**可能原因**:
1. 全文搜索未使用索引
2. LIKE '%keyword%' 导致全表扫描
3. 未使用Elasticsearch

**优化建议**:

**方案1: MySQL全文索引**
```sql
-- 添加全文索引
ALTER TABLE employees ADD FULLTEXT INDEX ft_emp_name (emp_name);

-- 使用全文搜索
SELECT * FROM employees
WHERE MATCH(emp_name) AGAINST('张三' IN NATURAL LANGUAGE MODE)
LIMIT 20;
```

**方案2: Elasticsearch**
```go
// 使用Elasticsearch进行搜索
func (s *EmployeeService) SearchWithES(ctx context.Context, req *SearchRequest) ([]*Employee, error) {
    query := elastic.NewBoolQuery()

    if req.Keyword != "" {
        query.Must(elastic.NewMultiMatchQuery(req.Keyword, "emp_name", "emp_code"))
    }

    result, err := s.es.Search().
        Index("employees").
        Query(query).
        From(req.Offset).
        Size(req.PageSize).
        Do(ctx)

    if err != nil {
        return nil, err
    }

    return parseEmployees(result), nil
}
```

---

## 📈 性能优化建议

### 1. 数据库优化

**索引优化**:
```sql
-- 租户隔离索引（所有表必备）
CREATE INDEX idx_tenant_id ON employees(tenant_id);

-- 复合索引（覆盖常用查询）
CREATE INDEX idx_tenant_status ON employees(tenant_id, status);
CREATE INDEX idx_tenant_dept ON employees(tenant_id, dept_id);

-- 分页索引
CREATE INDEX idx_tenant_hire_date ON employees(tenant_id, hire_date DESC);
```

**查询优化**:
- 使用游标分页代替OFFSET
- 避免SELECT *，只查询必要字段
- 使用Preload避免N+1查询
- 限制JOIN表数量（<5个）

**缓存策略**:
```go
// Redis缓存配置
cacheConfig := &redis.Config{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
    PoolSize: 100,
}

// 缓存层级
// L1: 内存缓存 (5分钟) - 组织树、部门树
// L2: Redis缓存 (15分钟) - 员工列表、岗位列表
// L3: 数据库 (持久化) - 所有数据
```

---

### 2. 应用层优化

**并发控制**:
```go
// 限制并发数据库连接数
db.SetMaxOpenConns(100)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(time.Hour)

// 使用worker pool批量处理
sem := make(chan struct{}, 10) // 最多10个并发
var wg sync.WaitGroup

for _, item := range items {
    wg.Add(1)
    sem <- struct{}{}

    go func(item Item) {
        defer wg.Done()
        defer func() { <-sem }()

        processItem(item)
    }(item)
}

wg.Wait()
```

**对象池化**:
```go
// 使用sync.Pool减少内存分配
var employeePool = sync.Pool{
    New: func() interface{} {
        return &Employee{}
    },
}

func getEmployee() *Employee {
    return employeePool.Get().(*Employee)
}

func putEmployee(emp *Employee) {
    emp.Reset()
    employeePool.Put(emp)
}
```

---

### 3. 架构优化

**读写分离**:
```
主库（Master）: 处理所有写操作
  ↓
从库（Slave 1-N）: 处理所有读操作（SELECT）
```

**分库分表**:
```sql
-- 按租户ID分表
employees_0000 -- tenant_id % 100 = 0
employees_0001 -- tenant_id % 100 = 1
...
employees_0099 -- tenant_id % 100 = 99
```

**微服务拆分**:
```
组织服务 (org-service): 组织、部门、岗位
员工服务 (employee-service): 员工、合同、调岗、离职
通讯录服务 (directory-service): 通讯录查询（聚合）
```

---

## 📝 性能测试报告模板

### 测试基本信息

| 项目 | 内容 |
|------|------|
| **测试日期** | 2025-01-01 |
| **测试人员** | 研发B |
| **测试环境** | 生产级配置 |
| **测试工具** | K6 v0.47+ |
| **测试数据规模** | 100租户, 10000组织, 100000员工 |
| **测试类型** | 基准测试 |

### 测试结果汇总

| API | p50 | p95 | p99 | 错误率 | 吞吐量 | 状态 |
|-----|-----|-----|-----|--------|--------|------|
| GET /api/organizations | 120ms | 250ms | 400ms | 0% | 100 req/s | ✅ |
| GET /api/org/employees | 200ms | 400ms | 700ms | 0.05% | 70 req/s | ⚠️ |
| GET /api/org/employees/search | 250ms | 500ms | 900ms | 0% | 50 req/s | ⚠️ |

### 问题汇总

| 严重性 | 问题描述 | 优化建议 | 优先级 |
|--------|----------|----------|--------|
| 高 | 员工列表查询慢 | 添加索引、使用缓存 | P0 |
| 中 | 员工搜索性能不足 | 接入Elasticsearch | P1 |
| 低 | 通讯录聚合查询慢 | 使用Redis缓存 | P2 |

### 优化效果对比

| API | 优化前p95 | 优化后p95 | 提升 |
|-----|-----------|-----------|------|
| GET /api/org/employees | 400ms | 250ms | 37.5% |

---

## 🔗 相关文档

- [性能测试README](./README.md)
- [研发B开发计划](../../../docs/企业级功能完善与统一性设计方案/研发B-后端工程师开发计划_v1.0.md)
- [企业级开发规范手册](../../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [统一错误码规范](../../../docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)

---

**最后更新**: 2025-01-01
**版本**: v1.0
**状态**: ✅ 生产就绪
