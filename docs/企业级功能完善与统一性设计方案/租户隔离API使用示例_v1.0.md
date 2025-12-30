# 租户隔离机制 API 使用示例

> **版本**: v1.0.0
> **最后更新**: 2025-12-30

---

## 目录

1. [快速开始](#1-快速开始)
2. [数据库层使用](#2-数据库层使用)
3. [缓存层使用](#3-缓存层使用)
4. [搜索层使用](#4-搜索层使用)
5. [存储层使用](#5-存储层使用)
6. [限流中间件使用](#6-限流中间件使用)
7. [完整示例](#7-完整示例)

---

## 1. 快速开始

### 1.1 安装GORM插件

```go
package main

import (
    "github.com/coze-dev/coze-studio/backend/api/middleware"
    "gorm.io/gorm"
)

func main() {
    // 初始化数据库连接
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }

    // 安装租户隔离插件
    db.Use(&middleware.TenantScopePlugin{
        SkipCallbackTables: map[string]bool{
            "tenants":       true,
            "subscriptions": true,
            "quotas":       true,
        },
        TenantIDGetter: middleware.GetTenantIDFromContext,
    })

    // ... 启动应用
}
```

### 1.2 注册中间件

```go
package main

import (
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/coze-dev/coze-studio/backend/api/middleware"
)

func main() {
    h := server.Default()

    // 注册租户隔离中间件
    h.Use(middleware.TenantIsolationMiddleware())

    // 注册限流中间件
    limiter := middleware.NewTenantRateLimiter(redisClient, &middleware.RateLimitConfig{
        Rate:  100,
        Burst: 200,
    })
    h.Use(limiter.Middleware())

    // ... 注册路由
    h.Spin()
}
```

---

## 2. 数据库层使用

### 2.1 基本CRUD操作

```go
package service

import (
    "context"
    "github.com/coze-dev/coze-studio/backend/api/middleware"
    "gorm.io/gorm"
)

// Create 创建记录（自动注入tenant_id）
func (s *BotService) Create(ctx context.Context, bot *Bot) error {
    return s.db.WithContext(ctx).Create(bot).Error
    // tenant_id会被自动注入，无需手动设置
}

// Get 查询单条记录（自动添加WHERE tenant_id = ?）
func (s *BotService) Get(ctx context.Context, botID string) (*Bot, error) {
    var bot Bot
    err := s.db.WithContext(ctx).
        Where("bot_id = ?", botID).
        First(&bot).Error
    // 自动添加: WHERE bot_id = ? AND tenant_id = ?
    return &bot, err
}

// List 查询列表（自动添加WHERE tenant_id = ?）
func (s *BotService) List(ctx context.Context, offset, limit int) ([]*Bot, error) {
    var bots []*Bot
    err := s.db.WithContext(ctx).
        Offset(offset).
        Limit(limit).
        Find(&bots).Error
    // 自动添加: WHERE tenant_id = ? LIMIT ? OFFSET ?
    return bots, err
}

// Update 更新记录（自动添加WHERE tenant_id = ?，阻止修改tenant_id）
func (s *BotService) Update(ctx context.Context, botID string, updates map[string]interface{}) error {
    return s.db.WithContext(ctx).
        Model(&Bot{}).
        Where("bot_id = ?", botID).
        Updates(updates).Error
    // 自动添加: WHERE bot_id = ? AND tenant_id = ?
    // 如果updates包含tenant_id，会被阻止
}

// Delete 删除记录（自动添加WHERE tenant_id = ?）
func (s *BotService) Delete(ctx context.Context, botID string) error {
    return s.db.WithContext(ctx).
        Where("bot_id = ?", botID).
        Delete(&Bot{}).Error
    // 自动添加: WHERE bot_id = ? AND tenant_id = ?
}
```

### 2.2 批量操作

```go
// BatchCreate 批量创建
func (s *BotService) BatchCreate(ctx context.Context, bots []*Bot) error {
    return s.db.WithContext(ctx).Create(&bots).Error
    // 每条记录都会自动注入当前租户的tenant_id
}

// BatchUpdate 批量更新
func (s *BotService) BatchUpdate(ctx context.Context, botIDs []string, updates map[string]interface{}) error {
    return s.db.WithContext(ctx).
        Where("bot_id IN ?", botIDs).
        Updates(updates).Error
    // 自动添加: WHERE bot_id IN (?, ?) AND tenant_id = ?
}

// BatchDelete 批量删除
func (s *BotService) BatchDelete(ctx context.Context, botIDs []string) error {
    return s.db.WithContext(ctx).
        Where("bot_id IN ?", botIDs).
        Delete(&Bot{}).Error
    // 自动添加: WHERE bot_id IN (?, ?) AND tenant_id = ?
}
```

### 2.3 复杂查询

```go
// Search 搜索（带租户过滤）
func (s *BotService) Search(ctx context.Context, keyword string, status string) ([]*Bot, error) {
    var bots []*Bot

    query := s.db.WithContext(ctx).
        Where("name LIKE ?", "%"+keyword+"%")

    if status != "" {
        query = query.Where("status = ?", status)
    }

    // 自动添加: WHERE (name LIKE ?) AND (status = ?) AND tenant_id = ?
    err := query.Find(&bots).Error
    return bots, err
}

// Count 统计（带租户过滤）
func (s *BotService) Count(ctx context.Context) (int64, error) {
    var count int64
    err := s.db.WithContext(ctx).
        Model(&Bot{}).
        Count(&count).Error
    // 自动添加: SELECT COUNT(*) FROM bots WHERE tenant_id = ?
    return count, err
}

// JoinQuery 关联查询（带租户过滤）
func (s *BotService) GetWithCreator(ctx context.Context, botID string) (*Bot, error) {
    var bot Bot
    err := s.db.WithContext(ctx).
        Preload("Creator").
        Where("bot_id = ?", botID).
        First(&bot).Error
    // 自动添加: WHERE bot_id = ? AND tenant_id = ?
    return &bot, err
}
```

### 2.4 事务操作

```go
// CreateWithConfig 创建Bot并配置（事务）
func (s *BotService) CreateWithConfig(ctx context.Context, bot *Bot, config *BotConfig) error {
    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 创建Bot（自动注入tenant_id）
        if err := tx.Create(bot).Error; err != nil {
            return err
        }

        // 创建配置（自动注入tenant_id）
        config.BotID = bot.BotID
        if err := tx.Create(config).Error; err != nil {
            return err
        }

        return nil
    })
}
```

---

## 3. 缓存层使用

### 3.1 基本操作

```go
package service

import (
    "context"
    "time"
    "github.com/coze-dev/coze-studio/backend/infra/cache"
)

type BotService struct {
    db    *gorm.DB
    cache *cache.TenantIsolatedCache
}

// GetWithCache 先查缓存，未命中再查数据库
func (s *BotService) GetWithCache(ctx context.Context, tenantID, botID string) (*Bot, error) {
    // 1. 尝试从缓存获取
    cacheKey := fmt.Sprintf("bot:%s", botID)
    var bot Bot

    err := s.cache.Get(ctx, tenantID, cacheKey, &bot)
    if err == nil {
        return &bot, nil // 缓存命中
    }

    // 2. 缓存未命中，查询数据库
    err = s.db.WithContext(ctx).
        Where("bot_id = ?", botID).
        First(&bot).Error
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    s.cache.Set(ctx, tenantID, cacheKey, &bot, 5*time.Minute)

    return &bot, nil
}

// InvalidateCache 使缓存失效
func (s *BotService) InvalidateCache(ctx context.Context, tenantID, botID string) error {
    cacheKey := fmt.Sprintf("bot:%s", botID)
    return s.cache.Delete(ctx, tenantID, cacheKey)
}
```

### 3.2 批量操作

```go
// MGetBatch 批量获取
func (s *BotService) MGetBatch(ctx context.Context, tenantID string, botIDs []string) (map[string]*Bot, error) {
    // 1. 批量从缓存获取
    cacheKeys := make([]string, len(botIDs))
    for i, id := range botIDs {
        cacheKeys[i] = fmt.Sprintf("bot:%s", id)
    }

    cachedData, err := s.cache.MGet(ctx, tenantID, cacheKeys)
    if err != nil {
        return nil, err
    }

    // 2. 找出未命中的ID
    missedIDs := make([]string, 0)
    result := make(map[string]*Bot)

    for _, id := range botIDs {
        key := fmt.Sprintf("bot:%s", id)
        if bot, ok := cachedData[key]; ok {
            result[id] = bot.(*Bot)
        } else {
            missedIDs = append(missedIDs, id)
        }
    }

    // 3. 批量查询未命中的数据
    if len(missedIDs) > 0 {
        var bots []Bot
        err := s.db.WithContext(ctx).
            Where("bot_id IN ?", missedIDs).
            Find(&bots).Error
        if err != nil {
            return nil, err
        }

        // 4. 写入缓存
        cacheData := make(map[string]interface{})
        for _, bot := range bots {
            result[bot.BotID] = &bot
            cacheKey := fmt.Sprintf("bot:%s", bot.BotID)
            cacheData[cacheKey] = &bot
        }

        s.cache.MSet(ctx, tenantID, cacheData, 5*time.Minute)
    }

    return result, nil
}
```

### 3.3 Hash操作

```go
// CacheUserSession 缓存用户会话（Hash结构）
func (s *UserService) CacheUserSession(ctx context.Context, tenantID, sessionID string, session *Session) error {
    key := fmt.Sprintf("session:%s", sessionID)

    data := map[string]interface{}{
        "user_id":    session.UserID,
        "username":   session.Username,
        "login_time": session.LoginTime,
        "expires_at": session.ExpiresAt,
    }

    return s.cache.HSet(ctx, tenantID, key, data)
}

// GetUserSession 获取用户会话
func (s *UserService) GetUserSession(ctx context.Context, tenantID, sessionID string) (*Session, error) {
    key := fmt.Sprintf("session:%s", sessionID)

    var session Session
    err := s.cache.HGet(ctx, tenantID, key, "user_id", &session.UserID)
    if err != nil {
        return nil, err
    }

    // 获取其他字段...
    return &session, nil
}
```

---

## 4. 搜索层使用

### 4.1 基本搜索

```go
package service

import (
    "context"
    "github.com/coze-dev/coze-studio/backend/infra/search"
    "gopkg.in/olivere/elastic.v8"
)

type BotService struct {
    db     *gorm.DB
    search *search.TenantIsolatedSearch
}

// IndexBot 索引Bot到Elasticsearch
func (s *BotService) IndexBot(ctx context.Context, tenantID string, bot *Bot) error {
    return s.search.IndexDocument(ctx, tenantID, "bots", bot.BotID, bot)
}

// SearchBots 搜索Bot
func (s *BotService) SearchBots(ctx context.Context, tenantID, keyword string) ([]*Bot, error) {
    // 构建查询
    query := elastic.NewBoolQuery().
        Must(elastic.NewMatchQuery("name", keyword))

    // 执行搜索（自动添加租户前缀到索引名）
    response, err := s.search.Search(
        ctx,
        tenantID,
        "bots", // 索引名会被转换为: tenant-{tenant_id}-bots
        query,
        search.WithMaxResults(20),
    )

    if err != nil {
        return nil, err
    }

    // 转换结果
    bots := make([]*Bot, len(response.Results))
    for i, result := range response.Results {
        bot := &Bot{}
        data, _ := json.Marshal(result.Source)
        json.Unmarshal(data, bot)
        bots[i] = bot
    }

    return bots, nil
}
```

### 4.2 批量索引

```go
// BulkIndexBots 批量索引Bot
func (s *BotService) BulkIndexBots(ctx context.Context, tenantID string, bots []*Bot) error {
    docs := make([]search.BulkDocument, len(bots))
    for i, bot := range bots {
        docs[i] = search.BulkDocument{
            DocID:    bot.BotID,
            Document: bot,
        }
    }

    return s.search.BulkIndex(ctx, tenantID, "bots", docs)
}
```

---

## 5. 存储层使用

### 5.1 文件上传

```go
package service

import (
    "context"
    "github.com/coze-dev/coze-studio/backend/infra/storage"
)

type FileService struct {
    storage *storage.TenantIsolatedStorage
}

// UploadAvatar 上传头像
func (s *FileService) UploadAvatar(ctx context.Context, tenantID, userID string, file io.Reader, size int64) (*storage.UploadResult, error) {
    fileID := generateUUID()

    metadata := map[string]string{
        "user_id": userID,
        "type":    "avatar",
    }

    return s.storage.Upload(
        ctx,
        tenantID,
        "avatar",       // 分类
        fileID,         // 文件ID
        file,           // 文件内容
        size,           // 文件大小
        "image/png",    // MIME类型
        metadata,       // 元数据
    )
}

// GetDownloadURL 获取下载URL
func (s *FileService) GetDownloadURL(ctx context.Context, tenantID, fileID string) (string, error) {
    key := fmt.Sprintf("avatar/%s", fileID) // 需要完整路径
    return s.storage.GetPresignedDownloadURL(ctx, tenantID, key, 1*time.Hour)
}
```

### 5.2 前端直传

```go
// GetUploadURL 获取上传URL（前端直接上传到对象存储）
func (s *FileService) GetUploadURL(ctx context.Context, tenantID, category, fileID string) (string, error) {
    return s.storage.GetPresignedUploadURL(
        ctx,
        tenantID,
        category,
        fileID,
        1*time.Hour, // URL有效期1小时
    )
}
```

---

## 6. 限流中间件使用

### 6.1 基本使用

```go
package main

import (
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/coze-dev/coze-studio/backend/api/middleware"
)

func main() {
    h := server.Default()

    // 创建限流器
    limiter := middleware.NewTenantRateLimiter(redisClient, &middleware.RateLimitConfig{
        Rate:                  100,              // 每秒100个请求
        Burst:                 200,              // 突发容量200
        EnableConcurrencyLimit: true,
        MaxConcurrency:        100,              // 最大并发100
    })

    // 注册中间件
    h.Use(limiter.Middleware())

    // ... 注册路由
}
```

### 6.2 动态配置

```go
// UpgradeTenantRateLimit 升级租户限流配置
func (s *TenantService) UpgradeRateLimit(ctx context.Context, tenantID string, rate float64) error {
    config := &middleware.RateLimitConfig{
        Rate:                  rate,
        Burst:                 int(rate * 2),
        EnableConcurrencyLimit: true,
        MaxConcurrency:        int(rate),
    }

    s.rateLimiter.SetTenantConfig(tenantID, config)

    return nil
}

// DowngradeTenantRateLimit 降级租户限流配置
func (s *TenantService) DowngradeRateLimit(ctx context.Context, tenantID string) error {
    s.rateLimiter.RemoveTenantConfig(tenantID)
    return nil
}
```

---

## 7. 完整示例

### 7.1 Bot管理服务

```go
package service

import (
    "context"
    "fmt"
    "time"

    "github.com/coze-dev/coze-studio/backend/api/middleware"
    "github.com/coze-dev/coze-studio/backend/infra/cache"
    "github.com/coze-dev/coze-studio/backend/infra/search"
    "github.com/coze-dev/coze-studio/backend/infra/storage"
    "gorm.io/gorm"
)

type BotService struct {
    db      *gorm.DB
    cache   *cache.TenantIsolatedCache
    search  *search.TenantIsolatedSearch
    storage *storage.TenantIsolatedStorage
}

// NewBotService 创建Bot服务
func NewBotService(
    db *gorm.DB,
    cache *cache.TenantIsolatedCache,
    search *search.TenantIsolatedSearch,
    storage *storage.TenantIsolatedStorage,
) *BotService {
    return &BotService{
        db:      db,
        cache:   cache,
        search:  search,
        storage: storage,
    }
}

// CreateBot 创建Bot
func (s *BotService) CreateBot(ctx context.Context, req *CreateBotRequest) (*Bot, error) {
    tenantID := middleware.GetTenantIDFromContext(ctx)

    // 1. 创建Bot记录
    bot := &Bot{
        BotID:   generateUUID(),
        Name:    req.Name,
        Status:  "active",
        // tenant_id会被GORM插件自动注入
    }

    if err := s.db.WithContext(ctx).Create(bot).Error; err != nil {
        return nil, err
    }

    // 2. 缓存Bot信息
    cacheKey := fmt.Sprintf("bot:%s", bot.BotID)
    s.cache.Set(ctx, tenantID, cacheKey, bot, 5*time.Minute)

    // 3. 索引到Elasticsearch
    s.search.IndexDocument(ctx, tenantID, "bots", bot.BotID, bot)

    return bot, nil
}

// GetBot 获取Bot
func (s *BotService) GetBot(ctx context.Context, botID string) (*Bot, error) {
    tenantID := middleware.GetTenantIDFromContext(ctx)

    // 1. 尝试从缓存获取
    cacheKey := fmt.Sprintf("bot:%s", botID)
    var bot Bot

    err := s.cache.Get(ctx, tenantID, cacheKey, &bot)
    if err == nil {
        return &bot, nil
    }

    // 2. 从数据库获取
    err = s.db.WithContext(ctx).
        Where("bot_id = ?", botID).
        First(&bot).Error
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    s.cache.Set(ctx, tenantID, cacheKey, &bot, 5*time.Minute)

    return &bot, nil
}

// UpdateBot 更新Bot
func (s *BotService) UpdateBot(ctx context.Context, botID string, req *UpdateBotRequest) (*Bot, error) {
    tenantID := middleware.GetTenantIDFromContext(ctx)

    // 1. 更新数据库
    updates := map[string]interface{}{
        "name":        req.Name,
        "description": req.Description,
        "updated_at":  time.Now().UnixMilli(),
    }

    err := s.db.WithContext(ctx).
        Model(&Bot{}).
        Where("bot_id = ?", botID).
        Updates(updates).Error
    if err != nil {
        return nil, err
    }

    // 2. 使缓存失效
    cacheKey := fmt.Sprintf("bot:%s", botID)
    s.cache.Delete(ctx, tenantID, cacheKey)

    // 3. 更新索引
    bot, err := s.GetBot(ctx, botID)
    if err != nil {
        return nil, err
    }

    s.search.IndexDocument(ctx, tenantID, "bots", botID, bot)

    return bot, nil
}

// DeleteBot 删除Bot
func (s *BotService) DeleteBot(ctx context.Context, botID string) error {
    tenantID := middleware.GetTenantIDFromContext(ctx)

    // 1. 删除数据库记录
    err := s.db.WithContext(ctx).
        Where("bot_id = ?", botID).
        Delete(&Bot{}).Error
    if err != nil {
        return err
    }

    // 2. 删除缓存
    cacheKey := fmt.Sprintf("bot:%s", botID)
    s.cache.Delete(ctx, tenantID, cacheKey)

    // 3. 删除索引
    s.search.DeleteDocument(ctx, tenantID, "bots", botID)

    return nil
}

// SearchBots 搜索Bot
func (s *BotService) SearchBots(ctx context.Context, keyword string, page, pageSize int) ([]*Bot, int64, error) {
    tenantID := middleware.GetTenantIDFromContext(ctx)

    // 1. 从Elasticsearch搜索
    query := elastic.NewBoolQuery().
        Must(elastic.NewMatchQuery("name", keyword))

    response, err := s.search.Search(
        ctx,
        tenantID,
        "bots",
        query,
        search.WithMaxResults(pageSize),
        search.WithFrom((page-1)*pageSize),
    )

    if err != nil {
        return nil, 0, err
    }

    // 2. 转换结果
    bots := make([]*Bot, len(response.Results))
    for i, result := range response.Results {
        bot := &Bot{}
        data, _ := json.Marshal(result.Source)
        json.Unmarshal(data, bot)
        bots[i] = bot
    }

    return bots, response.Total, nil
}

// UploadBotIcon 上传Bot图标
func (s *BotService) UploadBotIcon(ctx context.Context, botID string, file io.Reader, size int64) (*storage.UploadResult, error) {
    tenantID := middleware.GetTenantIDFromContext(ctx)

    return s.storage.Upload(
        ctx,
        tenantID,
        "bot_icon", // 分类
        botID,     // 文件ID
        file,
        size,
        "image/png",
        map[string]string{
            "bot_id": botID,
            "type":   "icon",
        },
    )
}

// GetBotIconURL 获取Bot图标URL
func (s *BotService) GetBotIconURL(ctx context.Context, botID string) (string, error) {
    tenantID := middleware.GetTenantIDFromContext(ctx)

    // 构建文件路径
    now := time.Now()
    key := fmt.Sprintf("bot_icon/%d/%02d/%02d/%s",
        now.Year(), now.Month(), now.Day(), botID)

    return s.storage.GetPresignedDownloadURL(ctx, tenantID, key, 1*time.Hour)
}
```

---

## 总结

本文档提供了租户隔离机制的完整使用示例，涵盖：

✅ **数据库层**: 自动租户过滤的CRUD操作
✅ **缓存层**: Redis缓存隔离和批量操作
✅ **搜索层**: Elasticsearch索引和搜索
✅ **存储层**: MinIO文件上传和下载
✅ **限流**: 租户级限流和动态配置
✅ **完整示例**: Bot管理服务的完整实现

**关键要点**:
1. 所有操作都会自动添加租户过滤
2. 租户ID会自动注入，无需手动设置
3. 多层隔离确保数据安全
4. 完整的错误处理和日志记录

**下一步**:
- 运行测试: `go test ./... -v`
- 查看文档: `docs/企业级功能完善与统一性设计方案/租户隔离机制实现总结_v1.0.md`
- 部署脚本: `bash scripts/deploy/tenant_isolation_setup.sh`
