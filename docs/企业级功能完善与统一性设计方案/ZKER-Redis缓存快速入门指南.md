# ZKER Redis缓存层 - 快速入门指南

**📅 更新时间**: 2025-12-31
**🎯 目标读者**: 后端工程师、DevOps工程师
**⏱️ 预计阅读时间**: 15分钟

---

## 🚀 快速开始

### 1. 环境准备

```bash
# 1. 确保Redis已启动
docker-compose up -d redis

# 2. 验证Redis连接
redis-cli ping
# 输出: PONG

# 3. 安装Go依赖（如果需要）
cd backend
go mod tidy
```

### 2. 使用权限缓存

```go
package main

import (
    "context"
    "fmt"
    "github.com/coze-dev/coze-studio/backend/application/permission"
    "github.com/coze-dev/coze-studio/backend/infra/cache"
    "gorm.io/gorm"
)

func main() {
    // 1. 初始化数据库连接
    db := initDB()

    // 2. 初始化Redis客户端
    redisClient := cache.NewRedisCache("localhost:6379", "", 0, 5*time.Minute)

    // 3. 初始化带缓存的权限检查器
    permChecker := permission.InitPermissionService(db, redisClient)

    // 4. 使用权限检查（自动缓存）
    ctx := context.Background()
    allowed, err := permChecker.CheckDataPermission(
        ctx,
        "tenant123",    // 租户ID
        "user456",      // 用户ID
        "bots",         // 资源类型
        "read",         // 操作类型
        "bot789",       // 资源ID
    )

    if err != nil {
        fmt.Printf("权限检查失败: %v\n", err)
        return
    }

    fmt.Printf("是否有权限: %v\n", allowed)

    // 5. 查看缓存统计
    stats := permChecker.GetCacheStats()
    hitRate := float64(stats.HitCount) / float64(stats.HitCount+stats.MissCount) * 100
    fmt.Printf("缓存命中率: %.2f%%\n", hitRate)
}
```

### 3. 使用配额缓存

```go
package main

import (
    "context"
    "fmt"
    "github.com/coze-dev/coze-studio/backend/application/tenant"
    "github.com/coze-dev/coze-studio/backend/infra/cache"
)

func main() {
    // 1. 初始化配额Repository
    quotaRepo := initQuotaRepo()

    // 2. 初始化Redis客户端
    redisClient := cache.NewRedisCache("localhost:6379", "", 0, 5*time.Minute)

    // 3. 初始化带缓存的配额服务
    quotaService := tenant.InitQuotaService(quotaRepo, redisClient)

    // 4. 检查配额（自动缓存）
    ctx := context.Background()
    err := quotaService.CheckQuota(ctx, "tenant123", "bots", 10)

    if err != nil {
        fmt.Printf("配额检查失败: %v\n", err)
        return
    }

    fmt.Println("配额充足")

    // 5. 消费配额（自动失效缓存）
    previousUsage, err := quotaService.ConsumeQuota(ctx, "tenant123", "bots", 5)
    if err != nil {
        fmt.Printf("消费配额失败: %v\n", err)
        return
    }

    fmt.Printf("之前的使用量: %d\n", previousUsage)

    // 6. 查看缓存统计
    stats := quotaService.GetCacheStats()
    hitRate := float64(stats.HitCount) / float64(stats.HitCount+stats.MissCount) * 100
    fmt.Printf("缓存命中率: %.2f%%\n", hitRate)
}
```

### 4. 使用配置缓存

```go
package main

import (
    "context"
    "fmt"
    "github.com/coze-dev/coze-studio/backend/infra/cache"
)

func main() {
    // 1. 初始化Redis客户端
    redisClient := cache.NewRedisCache("localhost:6379", "", 0, 10*time.Minute)

    // 2. 初始化配置缓存
    configCache := cache.NewConfigCache(redisClient)

    // 3. 获取配置（自动缓存）
    ctx := context.Background()
    value, err := configCache.GetConfig(ctx, "model.timeout", func() (string, error) {
        // 缓存未命中时从数据库查询
        return configRepo.Get(ctx, "model.timeout")
    })

    if err != nil {
        fmt.Printf("获取配置失败: %v\n", err)
        return
    }

    fmt.Printf("配置值: %s\n", value)

    // 4. 更新配置（自动失效缓存）
    err = configCache.SetConfig(ctx, "model.timeout", "30", func(value string) error {
        return configRepo.Set(ctx, "model.timeout", value)
    })

    if err != nil {
        fmt.Printf("更新配置失败: %v\n", err)
        return
    }

    fmt.Println("配置已更新")
}
```

---

## 📊 监控缓存性能

### 1. 查看缓存统计

```go
// 权限缓存统计
permStats := permChecker.GetCacheStats()
fmt.Printf("权限缓存:\n")
fmt.Printf("  命中次数: %d\n", permStats.HitCount)
fmt.Printf("  未命中次数: %d\n", permStats.MissCount)
fmt.Printf("  设置次数: %d\n", permStats.SetCount)
fmt.Printf("  删除次数: %d\n", permStats.DeleteCount)
fmt.Printf("  命中率: %.2f%%\n", float64(permStats.HitCount)/float64(permStats.HitCount+permStats.MissCount)*100)

// 配额缓存统计
quotaStats := quotaService.GetCacheStats()
fmt.Printf("配额缓存:\n")
fmt.Printf("  命中次数: %d\n", quotaStats.HitCount)
fmt.Printf("  未命中次数: %d\n", quotaStats.MissCount)
fmt.Printf("  消费次数: %d\n", quotaStats.ConsumeCount)
fmt.Printf("  错误次数: %d\n", quotaStats.ErrorCount)
fmt.Printf("  命中率: %.2f%%\n", float64(quotaStats.HitCount)/float64(quotaStats.HitCount+quotaStats.MissCount)*100)
```

### 2. Prometheus监控

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

func init() {
    // 注册监控指标
    prometheus.MustRegister(cacheHitRate)
    prometheus.MustRegister(cacheQueries)
    prometheus.MustRegister(cacheResponseTime)

    // 启动监控端点
    http.Handle("/metrics", promhttp.Handler())
    go http.ListenAndServe(":9090", nil)
}

// 使用示例
func recordCacheMetrics(cacheType string, duration time.Duration, hit bool) {
    cacheResponseTime.WithLabelValues(cacheType, "get").Observe(float64(duration.Milliseconds()))

    result := "miss"
    if hit {
        result = "hit"
    }
    cacheQueries.WithLabelValues(cacheType, result).Inc()
}
```

### 3. Grafana监控大盘

访问 `http://localhost:3000` 导入以下JSON：

```json
{
  "dashboard": {
    "title": "ZKER 缓存性能监控",
    "panels": [
      {
        "title": "缓存命中率",
        "targets": [
          {
            "expr": "cache_hit_rate{cache_type=\"permission\"}"
          },
          {
            "expr": "cache_hit_rate{cache_type=\"quota\"}"
          },
          {
            "expr": "cache_hit_rate{cache_type=\"config\"}"
          }
        ]
      },
      {
        "title": "缓存响应时间 (P95)",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, cache_response_duration_milliseconds)"
          }
        ]
      },
      {
        "title": "缓存QPS",
        "targets": [
          {
            "expr": "rate(cache_queries_total[1m])"
          }
        ]
      }
    ]
  }
}
```

---

## 🔧 缓存管理

### 1. 手动失效缓存

```go
// 失效用户的所有权限缓存
err := permChecker.InvalidateUserPermissions(ctx, "user456")
if err != nil {
    fmt.Printf("失效缓存失败: %v\n", err)
}

// 失效特定数据权限缓存
err = permChecker.InvalidateDataPermission(ctx, "tenant123", "user456", "bots", "read", "bot789")
if err != nil {
    fmt.Printf("失效缓存失败: %v\n", err)
}

// 失效配额缓存
err = quotaService.InvalidateQuota(ctx, "tenant123", "bots")
if err != nil {
    fmt.Printf("失效缓存失败: %v\n", err)
}

// 失效所有配置缓存
err = configCache.InvalidateAllConfigs(ctx)
if err != nil {
    fmt.Printf("失效缓存失败: %v\n", err)
}
```

### 2. 批量预热缓存

```go
// 预热常用权限缓存
func warmUpPermissionCache(ctx context.Context, permChecker *permission.CachedPermissionChecker) error {
    // 获取所有活跃用户
    activeUsers, err := userRepo.GetActiveUsers(ctx)
    if err != nil {
        return err
    }

    // 预热这些用户的常用权限
    for _, user := range activeUsers {
        for _, resourceType := range []string{"bots", "conversations", "knowledge"} {
            // 触发缓存加载（实际会检查权限，这里仅作示例）
            _, _ = permChecker.CheckDataPermission(
                ctx,
                user.TenantID,
                user.UserID,
                resourceType,
                "read",
                "*",
            )
        }
    }

    return nil
}

// 应用启动时预热
func main() {
    permChecker := permission.InitPermissionService(db, redisClient)

    // 预热缓存
    ctx := context.Background()
    if err := warmUpPermissionCache(ctx, permChecker); err != nil {
        log.Errorf("预热缓存失败: %v", err)
    }

    // 启动应用
    startServer()
}
```

---

## 🐛 故障排查

### 问题1：缓存命中率低

**症状**: 缓存命中率 < 50%

**排查步骤**:
```bash
# 1. 检查缓存键是否正确
redis-cli keys "zker:perm:*"

# 2. 检查缓存TTL设置
redis-cli TTL "zker:perm:tenant1:user1:bots:read:bot123"

# 3. 查看缓存统计
curl http://localhost:8080/api/admin/cache/stats
```

**解决方案**:
- 增加本地缓存大小
- 延长Redis缓存TTL
- 检查缓存键是否包含必要参数

### 问题2：Redis连接失败

**症状**: 日志中出现 `Redis connection refused`

**排查步骤**:
```bash
# 1. 检查Redis是否运行
docker ps | grep redis

# 2. 检查Redis端口
netstat -tuln | grep 6379

# 3. 测试Redis连接
redis-cli ping
```

**解决方案**:
```go
// 添加降级策略
func (c *PermissionCacheImpl) CheckDataPermission(...) (bool, error) {
    // L1: 尝试从本地缓存读取
    if val, ok := c.localCache.Get(cacheKey); ok {
        return val.(bool), nil
    }

    // L2: 尝试从Redis缓存读取（失败时跳过）
    val, err := c.cache.Get(ctx, buildKey("perm", cacheKey)).Result()
    if err != nil {
        log.Warnf("Redis不可用，跳过Redis缓存: %v", err)
    } else {
        c.localCache.Set(cacheKey, val == "1", time.Minute)
        return val == "1", nil
    }

    // L3: 执行数据库查询
    result, err := checkFn()
    if err != nil {
        return false, err
    }

    // 仅写回本地缓存
    c.localCache.Set(cacheKey, result, time.Minute)

    return result, nil
}
```

### 问题3：内存占用过高

**症状**: Redis内存使用率 > 80%

**排查步骤**:
```bash
# 1. 查看Redis内存使用
redis-cli INFO memory

# 2. 查看key数量
redis-cli DBSIZE

# 3. 查看大key
redis-cli --bigkeys
```

**解决方案**:
```bash
# 1. 调整Redis最大内存
redis-cli CONFIG SET maxmemory 2gb

# 2. 设置淘汰策略
redis-cli CONFIG SET maxmemory-policy allkeys-lru

# 3. 手动清理部分缓存
redis-cli --scan --pattern "zker:perm:tenant*" | xargs redis-cli DEL
```

---

## 📚 更多资源

### 相关文档
- [Redis缓存层实现报告](./ZKER-Redis缓存层实现报告_v1.0.md)
- [性能分析深度报告](./ZKER-性能分析深度报告_v1.0.md)
- [企业级开发规范手册](./ZKER-企业级开发规范手册_v1.0.md)

### API文档
- [权限缓存API文档](../../backend/infra/cache/permission_cache_impl.go)
- [配额缓存API文档](../../backend/infra/cache/quota_cache.go)
- [配置缓存API文档](../../backend/infra/cache/config_cache.go)

### 测试用例
- [权限缓存测试](../../backend/infra/cache/permission_cache_test.go)
- [配额缓存测试](../../backend/infra/cache/quota_cache_test.go)
- [性能测试](../../tests/performance/cache_performance_test.go)

---

**🎉 祝您使用愉快！**

如有问题，请查阅相关文档或联系技术支持团队。
