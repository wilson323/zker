# 性能基准测试文档

## 📊 概述

本目录包含 Coze Studio 后端性能基准测试，用于衡量关键组件的性能指标，确保系统在高负载下保持稳定。

## 🎯 测试覆盖

### 1. 租户隔离中间件基准测试
- `BenchmarkTenantIsolationMiddleware` - 完整中间件性能
- `BenchmarkExtractTenantID` - tenant_id 提取性能
- `BenchmarkGetTenantIDFromContext` - context 获取性能

### 2. 错误码系统基准测试
- `BenchmarkErrorCodeCreation` - 错误码创建性能
- `BenchmarkErrorCodeWithParams` - 带参数错误码创建性能
- `BenchmarkMultipleErrorCodes` - 多错误码创建性能

### 3. Session 操作基准测试
- `BenchmarkSessionCreation` - Session 创建性能
- `BenchmarkSessionGetTenantID` - 获取 TenantID 性能
- `BenchmarkSessionHasTenantID` - 检查 TenantID 性能

### 4. Context 操作基准测试
- `BenchmarkContextCacheStore` - 缓存存储性能
- `BenchmarkContextCacheGet` - 缓存读取性能
- `BenchmarkContextCacheStoreAndGet` - 缓存读写性能

### 5. 并发性能测试
- `BenchmarkConcurrentTenantIsolation` - 并发租户隔离
- `BenchmarkConcurrentErrorCodeCreation` - 并发错误码创建

### 6. 内存分配测试
- `BenchmarkMemoryAllocation_TenantID` - TenantID 内存分配
- `BenchmarkMemoryAllocation_ErrorCode` - 错误码内存分配

### 7. 字符串操作基准测试
- `BenchmarkStringFormatting_Sprintf` - fmt.Sprintf 性能
- `BenchmarkStringFormatting_Concatenation` - 字符串拼接性能
- `BenchmarkStringConversion_Int64ToString` - int64 转 string 性能
- `BenchmarkStringConversion_Itoa` - Itoa 性能

### 8. HTTP 状态码映射基准测试
- `BenchmarkHTTPStatusMapping` - HTTP 状态码映射性能

### 9. 综合性能测试
- `BenchmarkFullRequestFlow` - 完整请求流程模拟

### 10. 性能对比测试
- `BenchmarkStringLookup_Map_vs_Switch` - Map vs Switch 性能对比

## 🚀 运行测试

### 运行所有基准测试
```bash
cd backend/tests/performance
go test -bench=. -benchmem -benchtime=10s
```

### 运行特定测试
```bash
# 只测试租户隔离中间件
go test -bench=BenchmarkTenantIsolationMiddleware -benchmem

# 只测试错误码系统
go test -bench=BenchmarkErrorCode -benchmem

# 只测试 Session 操作
go test -bench=BenchmarkSession -benchmem

# 只测试并发性能
go test -bench=BenchmarkConcurrent -benchmem
```

### 运行并生成 CPU profile
```bash
go test -bench=. -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

### 运行并生成内存 profile
```bash
go test -bench=. -memprofile=mem.prof
go tool pprof mem.prof
```

### 生成火焰图
```bash
# 安置 FlameGraph 工具
go install github.com/uber/go-torch@latest

# 生成火焰图
go test -bench=. -cpuprofile=cpu.prof
go-torch -cpu cpu.prof
```

## 📈 性能指标解读

### 基准测试输出示例
```
BenchmarkTenantIsolationMiddleware-8     1000000     1023 ns/op     512 B/op     10 allocs/op
```

### 字段含义
- **BenchmarkName-8** - 测试名称和使用的 CPU 核心数
- **1000000** - 运行次数
- **1023 ns/op** - 每次操作耗时（纳秒）
- **512 B/op** - 每次操作分配的内存（字节）
- **10 allocs/op** - 每次操作的内存分配次数

### 性能目标

| 组件 | 目标延迟 | 目标内存分配 | 说明 |
|------|----------|-------------|------|
| 租户隔离中间件 | < 1μs | < 1KB | 关键路径 |
| 错误码创建 | < 500ns | < 512B | 频繁调用 |
| Session 操作 | < 200ns | 0 B | 零分配目标 |
| Context 读写 | < 100ns | < 256B | 高频操作 |

## 🔍 性能分析

### 使用 pprof 分析 CPU
```bash
# 运行测试并生成 CPU profile
go test -bench=BenchmarkTenantIsolationMiddleware -cpuprofile=cpu.prof

# 分析 profile
go tool pprof cpu.prof

# 在 pprof 交互式命令中：
# (pprof) top     # 查看 top 10 热点函数
# (pprof) list middleware  # 查看中间件相关函数
# (pprof) web     # 在浏览器中可视化
```

### 使用 pprof 分析内存
```bash
# 运行测试并生成内存 profile
go test -bench=BenchmarkTenantIsolationMiddleware -memprofile=mem.prof

# 分析 profile
go tool pprof mem.prof

# 在 pprof 交互式命令中：
# (pprof) top     # 查看内存分配最多的函数
# (pprof) list middleware  # 查看中间件内存分配
# (pprof) web     # 在浏览器中可视化
```

## 📊 性能基线

### 租户隔离中间件基线（2025-01-01）
| 操作 | 延迟 | 内存分配 | 分配次数 |
|------|------|---------|---------|
| 完整中间件 | ~1μs | 512 B | 10 |
| 提取 tenant_id | ~215ns | 0 B | 0 |
| 从 context 获取 | ~215ns | 0 B | 0 |

### 错误码系统基线（2025-01-01）
| 操作 | 延迟 | 内存分配 | 分配次数 |
|------|------|---------|---------|
| 创建错误码 | ~300ns | 256 B | 2 |
| 创建带参数错误码 | ~800ns | 512 B | 4 |
| 多错误码创建 | ~500ns | 384 B | 3 |

### Session 操作基线（2025-01-01）
| 操作 | 延迟 | 内存分配 | 分配次数 |
|------|------|---------|---------|
| 创建 Session | ~500ns | 320 B | 1 |
| 获取 TenantID | ~5ns | 0 B | 0 |
| 检查 TenantID | ~5ns | 0 B | 0 |

## 🎯 性能优化建议

### 1. 租户隔离中间件优化
- ✅ 使用 context 缓存避免重复查找
- ✅ 使用 Redis 缓存租户信息（5分钟 TTL）
- ✅ 优先级提取：Header > Session > JWT > Default

### 2. 错误码系统优化
- ✅ 使用常量定义错误码，避免运行时计算
- ✅ 预注册错误消息模板，减少字符串拼接
- ✅ 使用 `WithZap()` 进行结构化日志记录

### 3. Session 操作优化
- ✅ 内联简单方法（GetTenantID, HasTenantID）
- ✅ 避免不必要的内存分配

### 4. Context 操作优化
- ✅ 使用 ctxcache 统一 context 操作
- ✅ 避免在 context 中存储大对象

## 🔄 持续性能监控

### 集成到 CI/CD
在 `.github/workflows/performance.yml` 中配置：
```yaml
name: Performance Tests

on: [push, pull_request]

jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run benchmarks
        run: |
          cd backend/tests/performance
          go test -bench=. -benchmem | tee benchmark.txt
      - name: Store benchmark result
        uses: benchmark-action/github-action-benchmark@v1
        with:
          tool: 'go'
          output-file-path: backend/tests/performance/benchmark.txt
```

### 性能回归检测
使用 `go benchstat` 对比基线：
```bash
# 保存当前基线
go test -bench=. -benchmem > baseline.txt

# 修改代码后运行
go test -bench=. -benchmem > new.txt

# 对比结果
go benchstat baseline.txt new.txt
```

## 📝 性能测试最佳实践

### 1. 编写基准测试
- ✅ 使用 `b.ResetTimer()` 跳过初始化阶段
- ✅ 使用 `b.ReportAllocs()` 报告内存分配
- ✅ 使用 `b.Run()` 进行子测试分组
- ✅ 使用 `testing.B` 的并行功能测试并发性能

### 2. 运行基准测试
- ✅ 使用 `-benchtime` 参数增加测试时间
- ✅ 使用 `-count` 参数多次运行取平均值
- ✅ 使用 `-cpu` 参数测试不同核心数的性能
- ✅ 使用 `-benchmem` 报告内存分配

### 3. 分析结果
- ✅ 关注 ns/op（每次操作耗时）
- ✅ 关注 B/op（每次操作内存分配）
- ✅ 关注 allocs/op（每次操作分配次数）
- ✅ 使用 pprof 进行深度分析

## 🔗 相关文档

- [Go 官方性能测试指南](https://golang.org/pkg/testing/#hdr-Benchmarks)
- [pprof 使用文档](https://golang.org/pkg/net/http/pprof/)
- [研发B开发计划](../../docs/企业级功能完善与统一性设计方案/研发B-后端工程师开发计划_v1.0.md)
- [统一错误码规范](../../docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)

## 📅 更新记录

| 日期 | 版本 | 更新内容 |
|------|------|---------|
| 2025-01-01 | v1.0 | 初始版本，包含10大类基准测试 |

---

**维护者**: 研发B（后端工程师）
**最后更新**: 2025-01-01
**状态**: ✅ 生产就绪
