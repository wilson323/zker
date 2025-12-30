# Bot商店集成指南

## 快速集成步骤

### 步骤1：数据库迁移

```bash
cd docker
docker compose up -d mysql
atlas migrate apply
```

验证表是否创建成功：
```bash
mysql -u root -p123456 -e "USE opencoze; SHOW TABLES LIKE 'bot_store%';"
```

### 步骤2：初始化服务

在 `backend/application/botstore/init.go` 中添加：

```go
package botstore

import (
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/botstore/internal/dal"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/service"
)

var (
	PublisherSvc service.BotStorePublisher
	BrowserSvc  service.BotStoreBrowser
	ReviewerSvc  service.BotStoreReviewer
)

// Init 初始化Bot商店服务
func Init(db *gorm.DB) {
	storeRepo := dal.NewBotStoreRepository(db)
	categoryRepo := dal.NewBotStoreCategoryRepository(db)

	PublisherSvc = service.NewBotStorePublisher(storeRepo, categoryRepo)
	BrowserSvc = service.NewBotStoreBrowser(storeRepo, categoryRepo)
	ReviewerSvc = service.NewBotStoreReviewer(storeRepo)
}
```

### 步骤3：注册到主应用

在主应用初始化中调用：

```go
// backend/application/application.go
func InitApplication(db *gorm.DB) {
	// ... 其他初始化

	// 初始化Bot商店
	botstore.Init(db)
}
```

### 步骤4：注册路由

在 `backend/api/router/coze/api.go` 中添加：

```go
func RegisterRouter(r *server.Hertz) {
	// ... 其他路由

	// 注册Bot商店路由
	coze.InitBotStoreServices(
		botstore.PublisherSvc,
		botstore.BrowserSvc,
		botstore.ReviewerSvc,
	)
	coze.RegisterBotStoreRoutes(r)
}
```

### 步骤5：实现认证中间件

在 `backend/api/handler/coze/bot_store_service.go` 中实现：

```go
// 从JWT token中获取用户ID
func getUserIDFromContext(ctx context.Context) string {
	// 从context中获取用户信息
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	return ""
}

// 从context中获取租户ID
func getTenantIDFromContext(ctx context.Context) string {
	// 从context中获取租户信息
	if tenantID, ok := ctx.Value("tenant_id").(string); ok {
		return tenantID
	}
	return ""
}
```

### 步骤6：添加中间件保护

更新路由注册，添加认证中间件：

```go
func RegisterBotStoreRoutes(r *server.Hertz) {
	// Bot商店API组
	botStoreGroup := r.Group("/api/v1/bot-store")
	{
		// 需要认证的路由
		authGroup := botStoreGroup.Group("")
		authGroup.Use(middleware.AuthRequired())
		{
			authGroup.POST("/publish", PublishBotToStore)
			authGroup.POST("/:item_id/unpublish", UnpublishBot)
			authGroup.PUT("/:item_id", UpdateBotStoreItem)
		}

		// 公开路由
		botStoreGroup.GET("/list", ListBotStoreItems)
		botStoreGroup.GET("/search", SearchBotStoreItems)
		botStoreGroup.GET("/categories", GetBotCategories)
		botStoreGroup.GET("/:item_id", GetBotStoreItem)
	}

	// Bot商店管理API组（需要管理员权限）
	botStoreAdminGroup := r.Group("/api/v1/bot-store/admin")
	botStoreAdminGroup.Use(middleware.AuthRequired(), middleware.RequireAdmin())
	{
		botStoreAdminGroup.GET("/pending", GetPendingReviews)
		botStoreAdminGroup.POST("/:item_id/review", ReviewBotStoreItem)
	}
}
```

### 步骤7：测试集成

#### 7.1 测试发布Bot

```bash
curl -X POST http://localhost:8080/api/v1/bot-store/publish \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "bot_id": "bot123",
    "name": "智能助手",
    "description": "一个强大的AI助手",
    "category": "cat_productivity",
    "tags": ["效率", "AI"],
    "price": 0.0,
    "screenshots": ["https://example.com/screenshot.jpg"],
    "version": "1.0.0"
  }'
```

预期响应：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "item_id": "uuid-here",
    "bot_id": "bot123",
    "name": "智能助手",
    "status": "pending",
    "created_at": "2025-12-30T10:00:00Z"
  }
}
```

#### 7.2 测试浏览Bot

```bash
curl "http://localhost:8080/api/v1/bot-store/list?page=1&page_size=20&sort_by=popular"
```

预期响应：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
  }
}
```

#### 7.3 测试审核Bot

```bash
curl -X POST http://localhost:8080/api/v1/bot-store/admin/ITEM_ID/review \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -d '{
    "approved": true,
    "reason": ""
  }'
```

## 常见问题

### Q1: 如何获取分类列表？

```bash
curl http://localhost:8080/api/v1/bot-store/categories
```

### Q2: 如何搜索Bot？

```bash
curl "http://localhost:8080/api/v1/bot-store/search?query=效率&category=cat_productivity&page=1"
```

### Q3: 如何下架Bot？

```bash
curl -X POST http://localhost:8080/api/v1/bot-store/ITEM_ID/unpublish \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Q4: 状态码说明

- 200：成功
- 400：请求参数错误
- 403：权限不足
- 404：资源不存在
- 500：服务器错误

### Q5: Bot状态说明

- draft：草稿
- pending：待审核
- published：已发布
- rejected：已拒绝
- offline：已下架

## 部署检查清单

- [ ] 数据库迁移已执行
- [ ] 服务已初始化
- [ ] 路由已注册
- [ ] 认证中间件已配置
- [ ] 日志记录已添加
- [ ] 监控指标已配置
- [ ] API文档已更新
- [ ] 单元测试已通过
- [ ] 集成测试已通过
- [ ] 性能测试已通过

## 监控指标

建议添加以下监控指标：

1. **发布次数**：bot_store_publish_total
2. **浏览次数**：bot_store_view_total
3. **下载次数**：bot_store_download_total
4. **审核次数**：bot_store_review_total
5. **搜索次数**：bot_store_search_total
6. **API响应时间**：bot_store_api_duration_seconds
7. **API错误率**：bot_store_api_errors_total

## 日志记录

建议记录以下日志：

1. **发布日志**：bot_id、user_id、status
2. **审核日志**：item_id、reviewer_id、approved、reason
3. **浏览日志**：item_id、user_id、timestamp
4. **下载日志**：item_id、user_id、timestamp
5. **错误日志**：error、request_id、context

## 性能优化建议

1. **缓存分类列表**：Redis缓存1小时
2. **缓存热门Bot**：Redis缓存10分钟
3. **数据库索引**：确保所有索引已创建
4. **分页限制**：最大page_size=100
5. **CDN加速**：图片使用CDN
6. **搜索优化**：考虑使用Elasticsearch

## 安全建议

1. **认证**：所有写操作需要认证
2. **授权**：检查用户权限
3. **限流**：API限流保护
4. **验证**：严格的输入验证
5. **SQL注入**：使用参数化查询
6. **XSS防护**：输出转义

## 后续支持

如有问题，请联系开发团队或查阅文档：
- 架构文档：`backend/domain/botstore/README.md`
- 实现总结：`docs/botstore-implementation-summary.md`
- API文档：Swagger UI

---

**最后更新**：2025-12-30
**版本**：v1.0.0
