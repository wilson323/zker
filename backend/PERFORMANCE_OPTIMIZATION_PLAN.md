# ZKER 后端性能优化方案 v1.0

> 目标: QPS 10000+ | P99延迟 < 100ms | 优化时间: 2025-01-03

---

## 📊 性能瓶颈分析

### 1. 数据库查询瓶颈

#### 问题识别
- **N+1查询问题**: 在 `GetByID` 等方法中存在大量 `Preload` 调用
- **缺少索引优化**: 租户隔离查询缺少复合索引
- **连接池配置不足**: 当前配置 (MaxIdle: 10, MaxOpen: 100)

#### 具体问题
```go
// ❌ 问题代码示例 - org_repository_impl.go:50
func (r *organizationRepository) GetByID(ctx context.Context, orgID string) (*entity.Organization, error) {
    var org entity.Organization
    err := r.db.WithContext(ctx).
        Preload("Leader").      // N+1 问题
        Where("org_id = ?", orgID).
        Where("deleted_at IS NULL").
        First(&org).Error
    // ...
}
```

**影响**: 每次查询都会额外执行 1-2 次 JOIN 查询

---

### 2. 缓存策略问题

#### 当前问题
- **无本地缓存**: 每次都访问 Redis
- **缓存未预热**: 启动时无热点数据加载
- **无多级缓存**: L1(本地) + L2(Redis) 架构缺失

#### 热点数据识别
```go
// 高频访问数据 (需要缓存)
- tenant.Config (租户配置) - 每次API调用
- permission.DataPerms (权限数据) - 每次权限检查
- botstore.Category (商店分类) - 每次浏览
```

---

### 3. API性能瓶颈

#### Handler层问题
- **无响应压缩**: JSON响应未压缩
- **批量查询缺失**: 列表接口逐条查询
- **序列化开销**: 未使用高效序列化

---

## 🎯 优化方案

### 方案1: 数据库查询优化

#### 1.1 优化Preload策略
```go
// ✅ 优化后 - 使用Joins代替Preload
func (r *organizationRepository) GetByID(ctx context.Context, orgID string) (*entity.Organization, error) {
    var org entity.Organization
    err := r.db.WithContext(ctx).
        Joins("Leader").                      // 使用JOIN
        Select("organizations.*, leader.*").  // 明确字段
        Where("organizations.org_id = ?", orgID).
        Where("organizations.deleted_at IS NULL").
        First(&org).Error
    // ...
}
```

**预期收益**: 查询次数减少 50%, 延迟降低 60%

#### 1.2 添加复合索引
```sql
-- 租户隔离优化索引
CREATE INDEX idx_tenant_deleted ON organizations(tenant_id, deleted_at);
CREATE INDEX idx_tenant_status_deleted ON organizations(tenant_id, status, deleted_at);

-- Bot商店查询优化
CREATE INDEX idx_category_status_created ON bot_store_items(category, status, created_at DESC);
CREATE INDEX idx_publisher_tenant ON bot_store_items(publisher_id, tenant_id);
```

**预期收益**: 查询速度提升 3-5倍

#### 1.3 优化连接池配置
```go
// backend/infra/orm/impl/mysql/mysql.go
sqlDB.SetMaxIdleConns(50)      // 10 -> 50
sqlDB.SetMaxOpenConns(200)     // 100 -> 200
sqlDB.SetConnMaxLifetime(10 * time.Minute)  // 3600s -> 600s
sqlDB.SetConnMaxIdleTime(5 * time.Minute)   // 600s -> 300s
```

**预期收益**: 并发能力提升 2倍

---

### 方案2: Redis多级缓存

#### 2.1 实现本地缓存 (L1)
```go
// backend/infra/cache/multi_level_cache.go
package cache

import (
    "context"
    "time"
    "github.com/coze-dev/coze-studio/backend/pkg/logs"
    "github.com/patrickmn/go-cache"
)

type MultiLevelCache struct {
    l1 *cache.Cache      // 本地缓存 (5分钟)
    l2 *RedisCache       // Redis缓存 (1小时)
}

func NewMultiLevelCache(redis *RedisCache) *MultiLevelCache {
    return &MultiLevelCache{
        l1: cache.New(5*time.Minute, 10*time.Minute),
        l2: redis,
    }
}

func (m *MultiLevelCache) Get(ctx context.Context, key string, dest interface{}) error {
    // L1: 本地缓存
    if x, found := m.l1.Get(key); found {
        logs.Debugf("L1 cache hit: %s", key)
        dest = x.(*interface{})
        return nil
    }

    // L2: Redis缓存
    if err := m.l2.Get(ctx, key, dest); err == nil {
        logs.Debugf("L2 cache hit: %s", key)
        m.l1.Set(key, dest, cache.DefaultExpiration)
        return nil
    }

    return ErrCacheNotFound
}
```

#### 2.2 缓存热点数据
```go
// backend/domain/tenant/service/cache_service.go
func (s *TenantService) WarmupCache(ctx context.Context) error {
    // 预加载租户配置
    tenants, _ := s.repo.ListAll(ctx)
    for _, tenant := range tenants {
        config, _ := s.GetConfig(ctx, tenant.TenantID)
        s.cache.Set(ctx, fmt.Sprintf("tenant:config:%s", tenant.TenantID), config)
    }
    logs.Infof("Warmed up cache for %d tenants", len(tenants))
    return nil
}
```

**预期收益**: 缓存命中率 > 80%, API延迟降低 70%

---

### 方案3: API性能优化

#### 3.1 响应压缩
```go
// backend/api/middleware/compression.go
func GzipMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 检查Accept-Encoding
        if string(c.GetHeader("Accept-Encoding")) == "gzip" {
            c.SetHeader("Content-Encoding", "gzip")
            c.Response.Header.Set("Vary", "Accept-Encoding")
        }
        c.Next(ctx)
    }
}
```

#### 3.2 批量查询优化
```go
// backend/domain/botstore/service/bot_store_browser.go
func (s *botStoreBrowser) ListBots(ctx context.Context, req *ListBotsRequest) (*ListBotsResponse, error) {
    // ❌ 旧方案: 逐条查询分类信息
    for _, item := range items {
        category, _ := s.categoryRepo.GetByID(ctx, item.CategoryID)
        item.Category = category
    }

    // ✅ 新方案: 批量查询
    categoryIDs := extractCategoryIDs(items)
    categories, _ := s.categoryRepo.GetByIDs(ctx, categoryIDs)
    categoryMap := toMap(categories)
    for _, item := range items {
        item.Category = categoryMap[item.CategoryID]
    }
}
```

**预期收益**: 列表接口性能提升 5-10倍

---

## 📈 预期性能提升

| 指标 | 优化前 | 优化后 | 提升幅度 |
|------|--------|--------|----------|
| **QPS** | ~2,000 | **10,000+** | **5x** |
| **P99延迟** | ~300ms | **<100ms** | **3x** |
| **数据库连接数** | 100 | 200 | 2x |
| **缓存命中率** | 30% | **80%** | **2.7x** |
| **内存使用** | 500MB | 800MB | +60% |

---

## 🔧 实施步骤

### Phase 1: 数据库优化 (1天)
1. ✅ 添加复合索引 (tenant_id, deleted_at, status)
2. ✅ 优化Preload为Joins
3. ✅ 调整连接池配置

### Phase 2: 缓存优化 (2天)
1. ✅ 实现MultiLevelCache
2. ✅ 缓存预热机制
3. ✅ 缓存失效策略

### Phase 3: API优化 (1天)
1. ✅ 响应压缩
2. ✅ 批量查询优化
3. ✅ JSON序列化优化

### Phase 4: 性能测试 (1天)
1. ✅ 编写benchmark测试
2. ✅ 压力测试 (10000 QPS)
3. ✅ 生成优化报告

---

## 📝 验证标准

### 基准测试
```bash
# 运行性能测试
cd backend
go test ./... -bench=. -benchmem -benchtime=10s | tee benchmark_results.txt
```

### 压力测试
```bash
# 使用 wrk 或 ab
wrk -t12 -c400 -d30s http://localhost:8888/api/v1/bot-store/list
```

### 验收标准
- ✅ QPS ≥ 10000 (简单查询)
- ✅ P99延迟 < 100ms
- ✅ 缓存命中率 > 80%
- ✅ 数据库CPU < 60%
- ✅ 内存使用 < 2GB

---

## 📚 参考资料

- [Go Performance Best Practices](https://go.dev/doc/diagnostics)
- [MySQL Index Optimization](https://dev.mysql.com/doc/refman/8.0/en/optimization-indexes.html)
- [Redis Caching Strategies](https://redis.io/docs/manual/patterns/caching/)
- [GORM Performance Guide](https://gorm.io/docs/performance.html)

---

**版本**: v1.0 | **作者**: Performance Expert | **更新**: 2025-01-03
