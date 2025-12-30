# RBAC权限系统完整版集成指南

**文档版本**: v1.0
**创建日期**: 2025-01-01
**适用范围**: 所有需要使用RBAC权限系统的模块

---

## 📋 目录

1. [快速开始](#快速开始)
2. [数据库迁移](#数据库迁移)
3. [应用层集成](#应用层集成)
4. [API层集成](#api层集成)
5. [使用示例](#使用示例)
6. [性能优化](#性能优化)
7. [故障排查](#故障排查)

---

## 快速开始

### 1. 执行数据库迁移

```bash
# 在项目根目录执行
cd backend
mysql -u root -p your_database < domain/permission/internal/dal/permission_migration.sql
```

### 2. 初始化权限组件

在 `backend/application/permission/init.go` 中：

```go
package permission

import (
	"github.com/coze-dev/coze-studio/backend/domain/permission/service"
	"github.com/coze-dev/coze-studio/backend/infra/cache"
	"gorm.io/gorm"
)

var (
	DataPermChecker    service.DataPermissionChecker
	FieldPermChecker   service.FieldPermissionChecker
	PermissionCache    cache.PermissionCache
)

func Init(db *gorm.DB, redisClient *redis.Client, logger *zap.Logger) {
	// 1. 初始化缓存
	PermissionCache = cache.NewPermissionCache(redisClient, logger)

	// 2. 初始化Repository（需要实现）
	userDeptRepo := repository.NewUserDepartmentRepository(db)
	deptRepo := repository.NewDepartmentRepository(db)
	userRoleRepo := repository.NewUserRoleRepository(db)
	dataPermRepo := repository.NewDataPermissionRepository(db)
	fieldPermRepo := repository.NewFieldPermissionRepository(db)

	// 3. 初始化检查器
	DataPermChecker = service.NewDataPermissionChecker(
		db,
		userDeptRepo,
		deptRepo,
		userRoleRepo,
		dataPermRepo,
	)

	FieldPermChecker = service.NewFieldPermissionChecker(
		userRoleRepo,
		fieldPermRepo,
	)
}
```

---

## 数据库迁移

### 迁移脚本说明

迁移脚本位于：
```
backend/domain/permission/internal/dal/permission_migration.sql
```

### 主要变更

1. **data_permissions表扩展**：支持5级数据权限
2. **field_permissions表优化**：支持3级字段权限+脱敏规则
3. **departments表扩展**：支持树形结构
4. **新增审计表**：permission_audit_logs、permission_change_history

### 验证迁移

```sql
-- 验证表结构
DESC data_permissions;
DESC field_permissions;
DESC departments;

-- 验证视图
SELECT * FROM v_user_permissions LIMIT 10;

-- 验证存储过程
SHOW PROCEDURE STATUS WHERE Db = 'your_database';
```

---

## 应用层集成

### 1. 在BotService中使用数据权限

**文件**: `backend/application/bot/bot_service.go`

```go
package bot

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/permission/service"
	"github.com/coze-dev/coze-studio/backend/infra/cache"
)

type BotService struct {
	db              *gorm.DB
	dataPermChecker service.DataPermissionChecker
	permissionCache cache.PermissionCache
}

func NewBotService(
	db *gorm.DB,
	dataPermChecker service.DataPermissionChecker,
	permissionCache cache.PermissionCache,
) *BotService {
	return &BotService{
		db:              db,
		dataPermChecker: dataPermChecker,
		permissionCache: permissionCache,
	}
}

// ListBots 列出Bot（应用数据权限过滤）
func (s *BotService) ListBots(ctx context.Context, userID string, req *ListBotsRequest) (*ListBotsResponse, error) {
	// 1. 获取用户权限级别
	permissionLevel := s.getUserPermissionLevel(ctx, userID, "bots")

	// 2. 生成数据权限过滤条件
	filter, err := s.dataPermChecker.Filter(ctx, userID, permissionLevel, "bots")
	if err != nil {
		return nil, fmt.Errorf("failed to generate permission filter: %w", err)
	}

	// 3. 构建查询
	query := s.db.Table("bots")
	if filter.WhereClause != "" {
		query = query.Where(filter.WhereClause, filter.Args...)
	}

	// 4. 添加其他查询条件
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 5. 分页查询
	var bots []Bot
	var total int64

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&bots).Error; err != nil {
		return nil, err
	}

	return &ListBotsResponse{
		Bots:  bots,
		Total: total,
	}, nil
}

// GetBot 获取Bot详情（应用字段权限）
func (s *BotService) GetBot(ctx context.Context, userID, botID string) (*Bot, error) {
	var bot Bot

	// 1. 查询Bot
	if err := s.db.Where("bot_id = ?", botID).First(&bot).Error; err != nil {
		return nil, err
	}

	// 2. 应用字段权限脱敏
	maskedData, err := s.applyFieldPermissions(ctx, userID, &bot)
	if err != nil {
		return nil, err
	}

	return maskedData, nil
}

// CreateBot 创建Bot（验证字段权限）
func (s *BotService) CreateBot(ctx context.Context, userID string, req *CreateBotRequest) (*Bot, error) {
	// 1. 转换为map
	data := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"api_key":     req.APIKey,
		"webhook_url": req.WebhookURL,
	}

	// 2. 验证字段权限
	fieldPermChecker := s.getFieldPermChecker() // 获取字段权限检查器
	if err := fieldPermChecker.ValidateFieldPermissions(ctx, userID, "bots", data); err != nil {
		return nil, fmt.Errorf("field permission denied: %w", err)
	}

	// 3. 创建Bot
	bot := &Bot{
		BotID:     generateID(),
		Name:      req.Name,
		APIKey:    req.APIKey,
		CreatorID: userID,
	}

	if err := s.db.Create(bot).Error; err != nil {
		return nil, err
	}

	return bot, nil
}

// getUserPermissionLevel 获取用户权限级别
func (s *BotService) getUserPermissionLevel(ctx context.Context, userID, resourceType string) service.DataPermissionLevel {
	// TODO: 从用户角色中获取权限级别
	// 这里是简化示例
	return service.SELF
}

// applyFieldPermissions 应用字段权限
func (s *BotService) applyFieldPermissions(ctx context.Context, userID string, bot *Bot) (*Bot, error) {
	fieldPermChecker := s.getFieldPermChecker()

	// 转换为map
	data := map[string]interface{}{
		"bot_id":      bot.BotID,
		"name":        bot.Name,
		"description": bot.Description,
		"api_key":     bot.APIKey,
		"webhook_url": bot.WebhookURL,
	}

	// 脱敏处理
	maskedData, err := fieldPermChecker.MaskSensitiveFields(ctx, userID, data)
	if err != nil {
		return nil, err
	}

	// 转换回Bot对象
	return &Bot{
		BotID:      maskedData["bot_id"].(string),
		Name:       maskedData["name"].(string),
		APIKey:     maskedData["api_key"].(string),
		WebhookURL: maskedData["webhook_url"].(string),
	}, nil
}
```

---

## API层集成

### 1. 在路由中使用中间件

**文件**: `backend/api/router/bot/bot_router.go`

```go
package bot

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/coze-dev/coze-studio/backend/api/middleware"
	"github.com/coze-dev/coze-studio/backend/domain/permission/service"
	"github.com/coze-dev/coze-studio/backend/handler/bot"
)

// RegisterRoutes 注册Bot路由
func RegisterRoutes(r *server.Server, botHandler *bot.BotHandler) {
	// API分组
	api := r.Group("/api")

	// Bot路由组
	botGroup := api.Group("/bots")

	// 应用数据权限过滤中间件
	botGroup.GET("",
		middleware.DataPermissionFilter(middleware.DataPermissionFilterConfig{
			ResourceType: "bots",
			GetPermissionLevel: func(c *app.RequestContext) service.DataPermissionLevel {
				// 从context中获取用户权限级别
				userID := c.GetHeader("X-User-ID")
				return getUserPermissionLevelFromContext(c, userID)
			},
		}),
		botHandler.ListBots,
	)

	// 单个Bot路由（应用字段权限脱敏）
	botGroup.GET("/:bot_id",
		middleware.FieldPermissionMask("bots"),
		botHandler.GetBot,
	)

	// 创建Bot（验证字段权限）
	botGroup.POST("",
		middleware.ValidateFieldPermissions("bots"),
		botHandler.CreateBot,
	)

	// 更新Bot（验证字段权限）
	botGroup.PUT("/:bot_id",
		middleware.ValidateFieldPermissions("bots"),
		botHandler.UpdateBot,
	)

	// 删除Bot（检查数据权限）
	botGroup.DELETE("/:bot_id",
		middleware.RequirePermission(middleware.PermissionCheckConfig{
			ResourceType: "bots",
			Action:       "delete",
			GetResourceID: func(c *app.RequestContext) string {
				return c.Param("bot_id")
			},
		}),
		botHandler.DeleteBot,
	)
}

func getUserPermissionLevelFromContext(c *app.RequestContext, userID string) service.DataPermissionLevel {
	// 从数据库或缓存中获取用户权限级别
	return service.SELF // 示例
}
```

### 2. 在Handler中使用权限过滤器

**文件**: `backend/api/handler/bot/bot_handler.go`

```go
package bot

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/coze-dev/coze-studio/backend/domain/permission/service"
	berrno "github.com/coze-dev/coze-studio/backend/types/errno"
)

type BotHandler struct {
	botService *bot.BotService
}

func NewBotHandler(botService *bot.BotService) *BotHandler {
	return &BotHandler{
		botService: botService,
	}
}

// ListBots 列出Bot
func (h *BotHandler) ListBots(ctx context.Context, c *app.RequestContext) {
	userID := c.GetHeader("X-User-ID")
	tenantID := c.GetHeader("X-Tenant-ID")

	// 构建请求
	req := &ListBotsRequest{
		Page:     getIntParam(c, "page", 1),
		PageSize: getIntParam(c, "page_size", 20),
		Status:   c.Query("status"),
	}

	// 调用Service（Service内部已应用数据权限过滤）
	resp, err := h.botService.ListBots(ctx, userID, req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, ErrorResponse{
			Code:    berrno.ErrInternalCode.Code,
			Message: err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// GetBot 获取Bot详情
func (h *BotHandler) GetBot(ctx context.Context, c *app.RequestContext) {
	userID := c.GetHeader("X-User-ID")
	botID := c.Param("bot_id")

	// 调用Service（Service内部已应用字段权限脱敏）
	bot, err := h.botService.GetBot(ctx, userID, botID)
	if err != nil {
		c.JSON(consts.StatusNotFound, ErrorResponse{
			Code:    berrno.ErrBotNotFoundCode.Code,
			Message: "Bot not found",
		})
		return
	}

	c.JSON(consts.StatusOK, bot)
}
```

---

## 使用示例

### 示例1：创建角色并分配数据权限

```go
package main

import (
	"context"
	"fmt"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
)

func CreateRoleWithDataPermissions(ctx context.Context, roleRepo repository.RoleRepository, dataPermRepo repository.DataPermissionRepository) error {
	// 1. 创建角色
	role := &entity.Role{
		RoleID:      "role_manager",
		RoleName:    "部门经理",
		RoleCode:    "manager",
		RoleType:    entity.RoleTypeCustom,
		Description: "可以管理部门及子部门的数据",
		TenantID:    "tenant123",
	}

	if err := roleRepo.Create(ctx, role); err != nil {
		return err
	}

	// 2. 分配Bot数据权限（本部门及子部门）
	botPerm := &entity.DataPermission{
		PermissionID:  generateUUID(),
		RoleID:        role.RoleID,
		ResourceType:  entity.ResourceTypeBots,
		Scope:         entity.DataPermissionScopeDepartmentAndSub,
		CustomFilter:  "",
	}

	if err := dataPermRepo.Create(ctx, botPerm); err != nil {
		return err
	}

	// 3. 分配知识库数据权限（仅本人）
	knowledgePerm := &entity.DataPermission{
		PermissionID:  generateUUID(),
		RoleID:        role.RoleID,
		ResourceType:  entity.ResourceTypeKnowledge,
		Scope:         entity.DataPermissionScopeOwn,
		CustomFilter:  "",
	}

	if err := dataPermRepo.Create(ctx, knowledgePerm); err != nil {
		return err
	}

	fmt.Printf("角色 %s 创建成功，数据权限配置完成\n", role.RoleName)
	return nil
}
```

### 示例2：配置字段权限

```go
func ConfigureFieldPermissions(ctx context.Context, fieldPermRepo repository.FieldPermissionRepository) error {
	roleID := "role_manager"

	// 配置Bot字段权限
	fieldPerms := []*entity.FieldPermission{
		{
			PermissionID:    generateUUID(),
			RoleID:          roleID,
			ResourceType:    "bots",
			FieldName:       "api_key",
			PermissionLevel: entity.FieldPermissionLevelReadonly, // 只读
			MaskRule:        "partial",                           // 部分脱敏
		},
		{
			PermissionID:    generateUUID(),
			RoleID:          roleID,
			ResourceType:    "bots",
			FieldName:       "webhook_url",
			PermissionLevel: entity.FieldPermissionLevelHidden, // 隐藏
		},
		{
			PermissionID:    generateUUID(),
			RoleID:          roleID,
			ResourceType:    "bots",
			FieldName:       "name",
			PermissionLevel: entity.FieldPermissionLevelEditable, // 可编辑
		},
		{
			PermissionID:    generateUUID(),
			RoleID:          roleID,
			ResourceType:    "bots",
			FieldName:       "description",
			PermissionLevel: entity.FieldPermissionLevelEditable,
		},
	}

	for _, perm := range fieldPerms {
		if err := fieldPermRepo.Create(ctx, perm); err != nil {
			return err
		}
	}

	fmt.Println("字段权限配置完成")
	return nil
}
```

### 示例3：使用自定义过滤器

```go
func CreateCustomDataPermission(ctx context.Context, dataPermRepo repository.DataPermissionRepository) error {
	// 创建自定义过滤器JSON
	customFilter := map[string]interface{}{
		"conditions": []map[string]interface{}{
			{
				"field":    "status",
				"operator": "=",
				"value":    "published",
			},
			{
				"field":    "created_at",
				"operator": ">",
				"value":    "1704067200", // 2024-01-01的时间戳
			},
		},
		"logic": "AND",
	}

	filterJSON, _ := json.Marshal(customFilter)

	// 创建数据权限
	perm := &entity.DataPermission{
		PermissionID:  generateUUID(),
		RoleID:        "role_custom_viewer",
		ResourceType:  entity.ResourceTypeBots,
		Scope:         entity.DataPermissionScopeCustom,
		CustomFilter:  string(filterJSON),
	}

	return dataPermRepo.Create(ctx, perm)
}
```

---

## 性能优化

### 1. 缓存配置

权限系统使用Redis缓存，配置建议：

```yaml
# application.yaml
redis:
  addr: "localhost:6379"
  password: ""
  db: 0
  pool_size: 100

permission:
  cache:
    user_roles_ttl: 3600        # 用户角色缓存：1小时
    role_permissions_ttl: 1800   # 角色权限缓存：30分钟
    data_filter_ttl: 900        # 数据权限过滤器缓存：15分钟
    field_permissions_ttl: 1800 # 字段权限缓存：30分钟
```

### 2. 性能指标

期望性能指标：

| 指标 | 目标值 | 测量方法 |
|------|--------|----------|
| 权限检查响应时间 | < 50ms | 缓存命中时 |
| 缓存命中率 | > 90% | Prometheus监控 |
| 数据库查询次数 | < 2次/请求 | 使用EXPLAIN分析 |

### 3. 性能测试

```bash
# 运行性能测试
cd backend
go test -bench=. -benchmem ./domain/permission/service/

# 使用pprof分析性能
go test -cpuprofile=cpu.prof -memprofile=mem.prof ./domain/permission/service/
go tool pprof cpu.prof
```

---

## 故障排查

### 问题1：权限检查失败

**现象**：所有请求都返回403 Forbidden

**排查步骤**：

1. 检查用户是否有角色分配
```sql
SELECT * FROM user_roles WHERE user_id = 'user123';
```

2. 检查角色是否有权限配置
```sql
SELECT * FROM data_permissions WHERE role_id IN (
    SELECT role_id FROM user_roles WHERE user_id = 'user123'
);
```

3. 检查权限级别是否正确
```sql
SELECT role_id, resource_type, scope FROM data_permissions WHERE role_id = 'role123';
```

### 问题2：缓存失效导致性能下降

**现象**：权限检查响应时间 > 200ms

**排查步骤**：

1. 检查Redis连接
```bash
redis-cli ping
```

2. 检查缓存命中率
```bash
redis-cli INFO stats | grep keyspace_hits
```

3. 手动清理缓存
```go
// 在应用中添加清理接口
func InvalidateUserCache(ctx context.Context, userID string) error {
    return permissionCache.InvalidateUser(ctx, userID)
}
```

### 问题3：字段脱敏不生效

**现象**：敏感字段仍然可见

**排查步骤**：

1. 检查字段权限配置
```sql
SELECT * FROM field_permissions WHERE role_id = 'role123' AND field_name = 'api_key';
```

2. 确认是否使用了FieldPermissionMask中间件

3. 检查响应是否被正确序列化

---

## 附录

### A. 错误码清单

| 错误码 | 说明 | HTTP状态码 |
|--------|------|------------|
| 102001 | 权限检查失败 | 500 |
| 102002 | 数据权限拒绝 | 403 |
| 102003 | 字段权限拒绝 | 403 |
| 102004 | 角色不存在 | 404 |
| 102005 | 部门不存在 | 404 |

### B. API接口清单

| 接口 | 方法 | 说明 |
|------|------|------|
| /api/roles | POST | 创建角色 |
| /api/roles/:id | GET | 获取角色详情 |
| /api/roles/:id/permissions | PUT | 配置角色权限 |
| /api/users/:id/roles | GET | 获取用户角色列表 |
| /api/users/:id/roles | POST | 分配角色给用户 |
| /api/permissions/check | POST | 检查权限 |

### C. 参考文档

- [企业级开发规范手册](./ZKER-企业级开发规范手册_v1.0.md)
- [统一错误码定义规范](./ZKER-统一错误码定义规范.md)
- [多租户SaaS架构设计](./zker_MultiTenant_SaaS_完整架构设计文档.md)

---

**文档结束**
