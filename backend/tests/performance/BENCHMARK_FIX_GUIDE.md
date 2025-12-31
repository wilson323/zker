# Go Benchmark 测试修复指南

**日期**: 2025-01-03
**状态**: ⚠️ 需要修复

---

## 问题描述

### 文件位置
`backend/tests/performance/benchmark_handler_test.go`

### 错误1: 导入路径错误
```go
// ❌ 错误
tenanthandler "github.com/coze-dev/coze-studio/backend/api/handler/coze/tenant"

// ✅ 应该是（如果该包存在）
// 但实际项目中没有tenant子包
```

### 错误2: 假设的API不存在
```go
// ❌ 函数不存在
tenanthandler.RegisterTenantRoutes(h)
```

### 实际的Tenant API位置
实际项目中tenant相关的handler在：
- `backend/api/handler/coze/tenant_service.go`
- `backend/api/handler/coze/tenant_management_service.go`

这些文件包含：
- `CreateTenant`
- `GetTenant`
- `UpdateTenant`
- `DeleteTenant`
- `ListTenants`

---

## 修复方案

### 方案1: 删除不可用的测试（推荐）

由于这些benchmark测试假设的API不存在，最简单的方案是删除或注释掉它们：

```bash
# 重命名文件，标记为不可用
cd backend/tests/performance
mv benchmark_handler_test.go benchmark_handler_test.go.disabled
```

### 方案2: 重写测试使用实际API

需要重写测试以使用实际存在的handler函数：

```go
// benchmark_handler_test.go
package performance_test

import (
    "context"
    "testing"
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
)

// 直接测试handler函数，不使用路由注册
func BenchmarkTenantHandler_CreateTenant(b *testing.B) {
    // 使用实际的handler函数
    // tenant.CreateTenant(ctx, c)
}

func BenchmarkTenantHandler_GetTenant(b *testing.B) {
    // tenant.GetTenant(ctx, c)
}
```

### 方案3: 创建Mock测试（推荐用于CI/CD）

创建单元测试级别的benchmark，不依赖完整的HTTP服务器：

```go
// benchmark_tenant_service_test.go
package performance_test

import (
    "context"
    "testing"
    "github.com/coze-dev/coze-studio/backend/domain/tenant/service"
)

func BenchmarkTenantService_Create(b *testing.B) {
    // Mock repository
    // 测试service层性能
}
```

---

## 推荐行动

### 立即行动（5分钟）
1. **禁用不可用的测试**
   ```bash
   cd backend/tests/performance
   mv benchmark_handler_test.go benchmark_handler_test.go.disabled
   ```

2. **验证其他benchmark测试**
   ```bash
   cd backend
   go test -bench=. -benchmem ./tests/performance/ -run=^$
   ```

### 后续行动（可选）
1. 重写handler benchmark测试
2. 使用实际API端点
3. 集成到CI/CD

---

## 影响评估

| 受影响文件 | 状态 | 影响 |
|-----------|------|------|
| benchmark_handler_test.go | ⚠️ 不可用 | 无法运行handler层benchmark |
| benchmark_repository_test.go | ✅ 可用 | Repository层benchmark正常 |
| benchmark_service_test.go | ✅ 可用 | Service层benchmark正常 |

**结论**: 仅handler层benchmark受影响，repository和service层benchmark应该可以正常运行。

---

## 测试状态

- K6性能测试: ✅ 100% 可用
- Go Repository Benchmark: ✅ 应该可用
- Go Service Benchmark: ✅ 应该可用
- Go Handler Benchmark: ⚠️ 需要修复或禁用

**建议**: 暂时禁用handler benchmark，优先运行K6测试验证QPS目标。
