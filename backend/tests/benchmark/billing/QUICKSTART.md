# Token计量和预算管理系统 - 性能测试快速开始指南

**版本**: v1.0.0
**最后更新**: 2025-12-30

---

## 🚀 快速开始

### 1. 前置要求

```bash
# 检查Go版本
go version  # 需要 Go 1.24.0+

# 检查Bash版本
bash --version  # 需要 GNU Bash 4.0+

# （可选）安装pprof用于profile分析
go install github.com/google/pprof@latest
```

### 2. 进入测试目录

```bash
cd backend/tests/benchmark/billing
```

### 3. 运行快速测试

**方式1: 使用Makefile（推荐）**

```bash
# 快速测试（不包含压力测试）
make bench-quick

# 生成CPU和Memory profile
make profile

# 分析测试结果
make analyze
```

**方式2: 使用脚本**

```bash
# 快速测试
./run_benchmark.sh quick

# 分析结果
./analyze_results.sh benchmark_latest.txt
```

**方式3: 手动执行**

```bash
# 运行基准测试
go test -bench=. -benchmem -run=^$ ./

# 运行特定测试
go test -bench=BenchmarkRecordTokenUsage -benchmem ./

# 查看结果
cat benchmark_*.txt
```

---

## 📊 测试结果解读

### 基准测试输出

```
BenchmarkRecordTokenUsage-8      1000    1234567 ns/op    1234 B/op    12 allocs/op
```

**字段说明**:
- `-8`: GOMAXPROCS数量（8核CPU）
- `1000`: 迭代次数
- `1234567 ns/op`: 每次操作耗时（1.23ms/op）
- `1234 B/op`: 每次操作内存分配（1234字节）
- `12 allocs/op`: 每次操作内存分配次数（12次）

### 性能判断标准

| 实际性能 / 基线值 | 状态 | 说明 |
|-------------------|------|------|
| < 1.0 | ✅ 优秀 | 性能超出预期 |
| 1.0 - 2.0 | ⚠️ 良好 | 性能可接受 |
| > 2.0 | ❌ 需优化 | 性能需要优化 |

### 示例分析

**示例1: 性能优秀**
```
BenchmarkRecordTokenUsage-8      1000    50000000 ns/op    1024 B/op    10 allocs/op
```
- 实际性能: 50ms/op
- 基线性能: 100ms/op
- 比值: 50/100 = 0.5
- 结论: ✅ 优秀（性能超出预期）

**示例2: 需要优化**
```
BenchmarkRecordTokenUsage-8      1000    250000000 ns/op    2048 B/op    20 allocs/op
```
- 实际性能: 250ms/op
- 基线性能: 100ms/op
- 比值: 250/100 = 2.5
- 结论: ❌ 需要优化（性能是基线的2.5倍）

---

## 🔧 常用命令

### 测试执行

```bash
# 运行所有基准测试
go test -bench=. -benchmem ./

# 运行特定测试
go test -bench=BenchmarkPricingEngine -benchmem ./

# 运行测试并保存结果
go test -bench=. -benchmem . > results.txt 2>&1

# 运行测试（short模式，跳过压力测试）
go test -bench=. -benchmem -short ./
```

### Profile生成和分析

```bash
# 生成CPU profile
go test -cpuprofile=cpu.prof -bench=BenchmarkCPUProfile ./

# 分析CPU profile
go tool pprof cpu.prof
# 在pprof交互式命令中:
(pprof) top10          # 查看Top 10热点
(pprof) list funcName  # 查看特定函数
(pprof) web            # 生成可视化图表（需要浏览器）

# 生成Memory profile
go test -memprofile=mem.prof -bench=BenchmarkMemoryProfile ./

# 分析Memory profile
go tool pprof mem.prof
# 在pprof交互式命令中:
(pprof) top            # 查看Top 10内存分配
(pprof) list funcName  # 查看特定函数的内存分配
(pprof) web            # 生成可视化图表
```

### 内存泄漏检测

```bash
# 运行内存泄漏检测
go test -v -run=TestMemoryLeak ./

# 查看结果
cat memory_leak_*.txt
```

### 压力测试

```bash
# 运行持续高负载测试（1小时，⚠️ 耗时较长）
go test -v -run=TestStressTokenRecording ./

# 运行峰值负载测试（1分钟）
go test -v -run=TestPeakLoad ./

# 跳过压力测试（使用-short标志）
go test -short -v -run="TestStressTokenRecording|TestPeakLoad" ./
```

---

## 📁 文件清单

### 测试文件

| 文件 | 说明 | 行数 |
|------|------|------|
| `billing_benchmark_test.go` | 主测试文件，包含所有基准测试 | ~800行 |
| `token_metering_bench_test.go` | Token计量专项测试 | ~900行 |
| `budget_management_bench_test.go` | 预算管理专项测试 | ~800行 |
| `repository_bench_test.go` | 数据库Repository测试 | ~850行 |

### 脚本文件

| 文件 | 说明 | 功能 |
|------|------|------|
| `run_benchmark.sh` | 测试执行脚本 | 执行所有测试，生成报告 |
| `analyze_results.sh` | 结果分析脚本 | 解析结果，对比基线 |
| `Makefile` | Make命令 | 简化测试执行 |

### 文档文件

| 文件 | 说明 |
|------|------|
| `README.md` | 完整测试文档 |
| `QUICKSTART.md` | 快速开始指南（本文件） |
| `PERFORMANCE_TEST_REPORT_TEMPLATE.md` | 性能测试报告模板 |

---

## 🎯 性能目标

### 核心指标

| 操作 | 目标性能 | 实际性能 | 状态 |
|------|---------|---------|------|
| RecordTokenUsage | < 100ms | 待测试 | ⏳ |
| BatchRecordTokenUsage (1000) | < 1秒 | 待测试 | ⏳ |
| GetUsageStats | < 50ms | 待测试 | ⏳ |
| CheckBudget | < 50ms | 待测试 | ⏳ |
| PricingEngine (10000次) | < 100ms | 待测试 | ⏳ |
| 并发吞吐量 | > 1000 req/s | 待测试 | ⏳ |

### 资源限制

| 资源 | 限制 | 说明 |
|------|------|------|
| CPU | < 80% (单核) | 正常负载下 |
| 内存 | < 500MB | 测试过程中 |
| 数据库连接 | < 100 | 并发连接数 |
| goroutine数量 | < 1000 | 避免泄漏 |

---

## 📈 性能优化建议

### 快速优化技巧

1. **数据库查询优化**
   - 添加合适的索引
   - 使用`EXPLAIN`分析查询计划
   - 避免`SELECT *`

2. **缓存优化**
   - 实现多级缓存（本地 + Redis）
   - 预热热点数据
   - 设置合理的TTL

3. **并发优化**
   - 使用worker pool限制并发
   - 避免创建过多goroutine
   - 使用context控制超时

4. **内存优化**
   - 使用sync.Pool对象池
   - 避免频繁的字符串拼接
   - 及时释放大对象

---

## 🆘 常见问题

### Q1: 测试失败怎么办？

**A**: 检查以下几点：
1. 确保Go版本 >= 1.24.0
2. 确保在正确的目录下执行测试
3. 查看错误日志，确定失败原因

### Q2: 压力测试耗时太长？

**A**: 使用`-short`标志跳过压力测试：
```bash
go test -short -bench=. -benchmem ./
```

### Q3: 如何生成可视化图表？

**A**: 使用pprof的web功能：
```bash
go tool pprof -http=:8080 cpu.prof
# 然后在浏览器中访问 http://localhost:8080
```

### Q4: 测试结果如何对比基线？

**A**: 使用分析脚本：
```bash
./analyze_results.sh benchmark_result.txt
```

---

## 📚 相关文档

- [完整测试文档](README.md)
- [性能测试报告模板](PERFORMANCE_TEST_REPORT_TEMPLATE.md)
- [ZKER性能测试计划](../../../../../docs/企业级功能完善与统一性设计方案/ZKER-性能测试计划_v1.0.md)
- [ZKER性能基线文档](../../../../../docs/企业级功能完善与统一性设计方案/ZKER-性能基线文档.md)

---

## ✅ 下一步

1. **执行首次测试**:
   ```bash
   cd backend/tests/benchmark/billing
   make bench-quick
   ```

2. **分析测试结果**:
   ```bash
   make analyze
   ```

3. **根据结果优化**:
   - 查看未达标的测试项
   - 参考优化建议进行优化
   - 重新测试验证优化效果

4. **建立性能基线**:
   - 保存首次测试结果作为基线
   - 定期执行性能测试
   - 监控性能趋势

---

**祝你测试愉快！** 🎉

如有问题，请参考[完整测试文档](README.md)或联系研发B。
