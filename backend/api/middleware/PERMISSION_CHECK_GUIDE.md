# 权限检查中间件使用指南

**文档版本**: v1.0
**创建日期**: 2025-01-01
**负责人**: 后端架构团队

---

## 📋 功能概述

权限检查中间件用于在 HTTP 请求处理前自动检查用户的权限，确保用户只能访问其有权访问的资源。

### 核心功能

- ✅ **数据权限检查**: 支持全部/部门/个人/自定义/无权限 5级数据权限
- ✅ **角色检查**: 支持单一角色和任一角色检查
- ✅ **字段权限过滤**: 根据用户角色过滤敏感字段（editable/readonly/hidden）
- ✅ **手动权限检查**: 提供辅助函数在Handler内部进行权限检查
- ✅ **灵活配置**: 支持多种资源类型和自定义权限提取逻辑

---

## 🚀 快速开始

### 1. 初始化中间件

在应用启动时注入权限检查服务：

```go
// backend/server.go 或 main.go

package main

import (
    "context"
    "github.com/coze-studio/backend/api/middleware"
    "github.com/coze-studio/backend/domain/permission/repository"
    "github.com/coze-studio/backend/domain/permission/service"
    "gorm.io/gorm"
)

func main() {
    db := initDB()

    // 创建仓储
    roleRepo := repository.NewRoleRepository(db)
    dataPermRepo := repository.NewDataPermissionRepository(db)
    fieldPermRepo := repository.NewFieldPermissionRepository(db)
    userRoleRepo := repository.NewUserRoleRepository(db)
    userDeptRepo := repository.NewUserDepartmentRepository(db)
    departmentRepo := repository.NewDepartmentRepository(db)

    // 创建基础权限检查服务
    permissionSvc := service.NewPermissionChecker(
        db,
        roleRepo,
        dataPermRepo,
        fieldPermRepo,
        userRoleRepo,
        userDeptRepo,
        departmentRepo,
    )

    // 创建增强权限检查器（包装基础服务）
    enhancedChecker := middleware.NewEnhancedPermissionChecker(
        permissionSvc,
        userRoleRepo,
    )

    // 初始化权限中间件
    middleware.InitEnhancedPermissionMiddleware(enhancedChecker)

    // ... 启动服务器 ...
}
```

### 2. 在路由中使用

#### 场景 1: 检查Bot访问权限（数据权限）

```go
// backend/api/handler/bot_handler.go

package handler

import (
    "github.com/cloudwego/hertz/pkg/server"
    "github.com/coze-studio/backend/api/middleware"
    "github.com/coze-studio/backend/domain/permission/entity"
)

// RegisterBotRoutes 注册Bot相关路由
func RegisterBotRoutes(r *server.Hertz) {
    // 获取Bot：检查读权限
    r.GET("/api/v1/bots/:bot_id",
        middleware.RequirePermission(middleware.PermissionCheckConfig{
            ResourceType: entity.ResourceTypeBots,
            Action:       "read",
            GetResourceID: middleware.ParseResourceIDFromPath("bot_id"),
        }),
        h.GetBot,
    )

    // 更新Bot：检查写权限
    r.PUT("/api/v1/bots/:bot_id",
        middleware.RequirePermission(middleware.PermissionCheckConfig{
            ResourceType: entity.ResourceTypeBots,
            Action:       "write",
            GetResourceID: middleware.ParseResourceIDFromPath("bot_id"),
        }),
        h.UpdateBot,
    )

    // 删除Bot：检查删除权限
    r.DELETE("/api/v1/bots/:bot_id",
        middleware.RequirePermission(middleware.PermissionCheckConfig{
            ResourceType: entity.ResourceTypeBots,
            Action:       "delete",
            GetResourceID: middleware.ParseResourceIDFromPath("bot_id"),
        }),
        h.DeleteBot,
    )
}
```

#### 场景 2: 检查管理员权限（角色检查）

```go
// 注册管理员路由
func RegisterAdminRoutes(r *server.Hertz) {
    // 创建用户：需要admin角色
    r.POST("/api/v1/admin/users",
        middleware.RequireRole("admin"),
        h.CreateUser,
    )

    // 批量导入用户：需要admin或super_admin任一角色
    r.POST("/api/v1/admin/users/batch",
        middleware.RequireAnyRole("admin", "super_admin"),
        h.BatchImportUsers,
    )
}
```

#### 场景 3: 字段权限过滤（敏感信息隐藏）

```go
// 获取Bot详情（过滤敏感字段）
r.GET("/api/v1/bots/:bot_id",
    middleware.RequirePermission(middleware.PermissionCheckConfig{
        ResourceType: entity.ResourceTypeBots,
        Action:       "read",
        GetResourceID: middleware.ParseResourceIDFromPath("bot_id"),
    }),
    middleware.RequireFieldPermission("bots", []string{"id", "name", "description"}),
    func(ctx context.Context, c *app.RequestContext) {
        // 1. 从数据库获取Bot
        bot := h.getBotFromDB(c.Param("bot_id"))

        // 2. 根据字段权限过滤敏感字段
        filteredBot, err := middleware.FilterFieldsByPermission(ctx, c, &bot)
        if err != nil {
            c.JSON(500, ErrorResponse(err))
            return
        }

        // 3. 返回过滤后的数据
        c.JSON(200, SuccessResponse(filteredBot))
    },
)
```

#### 场景 4: 手动权限检查（复杂业务逻辑）

```go
// backend/api/handler/bot_handler.go

package handler

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/coze-studio/backend/api/middleware"
    "github.com/coze-studio/backend/domain/permission/entity"
)

func CreateBot(ctx context.Context, c *app.RequestContext) {
    var req CreateBotRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, ErrorResponse(err))
        return
    }

    // 方式1：手动数据权限检查
    err := middleware.ManualPermissionCheck(ctx, c, middleware.PermissionCheckConfig{
        ResourceType: entity.ResourceTypeBots,
        Action:       "create",
    })
    if err != nil {
        c.JSON(403, ErrorResponse(err))
        return
    }

    // 方式2：手动角色检查
    err = middleware.ManualRoleCheck(ctx, c, "bot_creator")
    if err != nil {
        c.JSON(403, ErrorResponse(err))
        return
    }

    // 方式3：获取用户所有权限，进行复杂判断
    perms, err := middleware.GetCurrentUserPermissions(ctx, c)
    if err != nil {
        c.JSON(500, ErrorResponse(err))
        return
    }

    // 检查是否有bot.create权限
    if !contains(perms, "bot.create") {
        c.JSON(403, ErrorResponse("permission denied: bot.create required"))
        return
    }

    // ... 继续业务逻辑 ...
}

func contains(perms []string, perm string) bool {
    for _, p := range perms {
        if p == perm {
            return true
        }
    }
    return false
}
```

---

## 🔧 高级用法

### 1. 自定义资源ID提取

如果资源ID不在路径参数或查询参数中，可以自定义提取函数：

```go
// 从请求体中提取资源ID
r.POST("/api/v1/bots/delete",
    middleware.RequirePermission(middleware.PermissionCheckConfig{
        ResourceType: entity.ResourceTypeBots,
        Action:       "delete",
        GetResourceID: func(c *app.RequestContext) string {
            var req struct {
                BotID string `json:"bot_id"`
            }
            c.Bind(&req)
            return req.BotID
        },
    }),
    h.DeleteBot,
)

// 从多个来源提取（优先级：Path > Query > Body）
func extractResourceID(c *app.RequestContext) string {
    // 1. 尝试从路径参数获取
    if botID := c.Param("bot_id"); botID != "" {
        return botID
    }

    // 2. 尝试从查询参数获取
    if botID := c.Query("bot_id"); botID != "" {
        return botID
    }

    // 3. 尝试从请求体获取
    var req struct {
        BotID string `json:"bot_id"`
    }
    if err := c.Bind(&req); err == nil && req.BotID != "" {
        return req.BotID
    }

    return ""
}
```

### 2. 权限检查中间件链

可以组合多个权限检查中间件：

```go
// 需要同时满足：admin角色 + 有Bot的写权限
r.PUT("/api/v1/bots/:bot_id",
    middleware.RequireRole("admin"),                    // 1. 必须是admin
    middleware.RequirePermission(middleware.PermissionCheckConfig{ // 2. 必须有Bot写权限
        ResourceType: entity.ResourceTypeBots,
        Action:       "write",
        GetResourceID: middleware.ParseResourceIDFromPath("bot_id"),
    }),
    middleware.RequireFieldPermission("bots", []string{"id", "name"}), // 3. 字段权限过滤
    h.UpdateBot,
)
```

### 3. 条件权限检查

根据业务逻辑决定是否进行权限检查：

```go
func ConditionalPermissionCheck() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 公开Bot不需要权限检查
        if isPublicBot(c.Param("bot_id")) {
            c.Next(ctx)
            return
        }

        // 非公开Bot需要权限检查
        middleware.RequirePermission(middleware.PermissionCheckConfig{
            ResourceType: entity.ResourceTypeBots,
            Action:       "read",
            GetResourceID: middleware.ParseResourceIDFromPath("bot_id"),
        })(ctx, c)
    }
}

func isPublicBot(botID string) bool {
    // 检查Bot是否公开
    return botService.IsPublic(botID)
}
```

### 4. 权限检查失败的自定义处理

```go
func CustomPermissionCheck() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        err := middleware.ManualPermissionCheck(ctx, c, middleware.PermissionCheckConfig{
            ResourceType: entity.ResourceTypeBots,
            Action:       "read",
        })

        if err != nil {
            // 自定义错误响应
            if isPermissionDeniedError(err) {
                c.JSON(403, map[string]interface{}{
                    "code":    40301,
                    "message": "您没有权限访问此Bot",
                    "hint":    "请联系管理员分配权限",
                })
                c.Abort()
                return
            }

            // 其他错误
            c.JSON(500, ErrorResponse(err))
            c.Abort()
            return
        }

        c.Next(ctx)
    }
}
```

---

## 📊 权限检查流程

### 数据权限检查流程（RequirePermission）

```
用户请求
    ↓
提取 user_id 和 tenant_id
    ↓ (失败)
返回 401 Unauthorized
    ↓ (成功)
检查权限检查器是否初始化
    ↓ (未初始化)
跳过检查，继续处理（降级）
    ↓ (已初始化)
获取资源ID
    ↓
检查数据权限
    ↓ (失败)
返回 403 Forbidden
    ↓ (成功)
继续处理请求
```

### 角色检查流程（RequireRole）

```
用户请求
    ↓
提取 user_id 和 tenant_id
    ↓ (失败)
返回 401 Unauthorized
    ↓ (成功)
检查权限检查器是否初始化
    ↓ (未初始化)
跳过检查，继续处理（降级）
    ↓ (已初始化)
获取用户的所有角色
    ↓
检查是否包含 required_role
    ↓ (不包含)
返回 403 Forbidden
    ↓ (包含)
继续处理请求
```

### 字段权限过滤流程（FilterFieldsByPermission）

```
Handler获取数据
    ↓
获取用户的字段权限配置
    ↓
使用反射遍历结构体字段
    ↓
对于每个字段：
  - 获取字段名（JSON tag）
  - 检查权限级别（editable/readonly/hidden）
  - hidden字段：过滤掉
  - editable/readonly：保留
    ↓
返回过滤后的数据
```

---

## ⚠️ 注意事项

### 1. user_id 和 tenant_id 传递

确保请求中包含认证信息：

```go
// 方式1：通过Header传递（推荐）
req.Header.Set("X-User-ID", "user_123")
req.Header.Set("X-Tenant-ID", "tenant_123")

// 方式2：通过认证中间件设置
func AuthMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        token := c.GetHeader("Authorization")
        claims := validateToken(token)

        // 设置到上下文，供权限中间件使用
        c.Set("user_id", claims.UserID)
        c.Set("tenant_id", claims.TenantID)

        c.Next(ctx)
    }
}
```

### 2. 错误处理

权限检查失败返回以下HTTP状态码：

| 状态码 | 错误码 | 说明 |
|-------|--------|------|
| 401 | 40100 | 未认证（缺少user_id） |
| 400 | 20004 | 缺少tenant_id |
| 403 | 40301 | 数据权限拒绝 |
| 403 | 40300 | 角色权限拒绝 |
| 500 | 50015 | 权限检查失败 |

**403 数据权限拒绝响应示例**：
```json
{
  "code": 40301,
  "message": "data permission denied",
  "data": null,
  "resource_type": "bots",
  "resource_id": "bot_123"
}
```

**403 角色权限拒绝响应示例**：
```json
{
  "code": 40300,
  "message": "permission denied",
  "data": null,
  "required_role": "admin"
}
```

### 3. 性能考虑

- ✅ **缓存角色**: 可以缓存用户的角色信息，减少数据库查询
- ✅ **批量检查**: 批量操作时一次性检查所有权限
- ✅ **降级策略**: 权限检查器未初始化时可以降级（跳过检查）

**缓存角色示例**：
```go
type CachedPermissionChecker struct {
    checker       PermissionCheckerInterface
    roleCache     *lru.Cache  // github.com/hashicorp/golang-lru
    cacheDuration time.Duration
}

func (c *CachedPermissionChecker) UserHasRole(ctx context.Context, tenantID, userID, roleCode string) (bool, error) {
    cacheKey := fmt.Sprintf("roles:%s:%s", tenantID, userID)

    // 1. 尝试从缓存获取
    if cached, ok := c.roleCache.Get(cacheKey); ok {
        roles := cached.([]*entity.Role)
        for _, role := range roles {
            if role.RoleCode == roleCode {
                return true, nil
            }
        }
        return false, nil
    }

    // 2. 从数据库获取
    roles, err := c.checker.(*EnhancedPermissionChecker).userRoleRepo.GetRolesByUser(ctx, userID, tenantID)
    if err != nil {
        return false, err
    }

    // 3. 写入缓存
    c.roleCache.Add(cacheKey, roles)

    // 4. 检查角色
    for _, role := range roles {
        if role.RoleCode == roleCode {
            return true, nil
        }
    }

    return false, nil
}
```

### 4. 字段过滤限制

`FilterFieldsByPermission` 使用反射过滤字段，有以下限制：

- ⚠️ 仅支持结构体，不支持slice/map
- ⚠️ 需要正确的JSON tag
- ⚠️ 嵌套结构体需要递归处理（当前未实现）

**不支持的情况**：
```go
// ❌ 不支持：直接返回map
botMap := map[string]interface{}{
    "id": "bot_123",
    "name": "My Bot",
}
filtered, err := middleware.FilterFieldsByPermission(ctx, c, botMap) // 错误！

// ✅ 正确：使用结构体
type Bot struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
bot := Bot{ID: "bot_123", Name: "My Bot"}
filtered, err := middleware.FilterFieldsByPermission(ctx, c, &bot) // 正确！
```

---

## 🧪 测试

### 单元测试

```bash
# 运行权限检查中间件测试
cd backend
go test ./api/middleware -run TestRequirePermission -v

# 运行所有测试
go test ./api/middleware -v

# 生成覆盖率报告
go test ./api/middleware -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 集成测试

```go
// backend/api/middleware/integration_test.go

func TestPermissionCheckIntegration(t *testing.T) {
    // 1. 启动测试服务器
    testServer := setupTestServer()
    defer testServer.Close()

    // 2. 创建测试用户和角色
    createTestUserWithRole("user_123", "tenant_123", "admin")

    // 3. 发送请求（有权限）
    resp := getBot(testServer, "user_123", "tenant_123", "bot_123")
    assert.Equal(t, 200, resp.StatusCode)

    // 4. 发送请求（无权限）
    createTestUserWithRole("user_456", "tenant_123", "guest")
    resp = getBot(testServer, "user_456", "tenant_123", "bot_123")
    assert.Equal(t, 403, resp.StatusCode)
}
```

---

## 📚 相关文档

- [ZKER-企业级开发规范手册_v1.0.md](../../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [ZKER-统一错误码定义规范.md](../../../docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)
- [权限管理设计文档](../../../docs/企业级功能完善与统一性设计方案/权限管理设计文档.md)

---

## 🐛 常见问题

### Q1: 权限检查失败但用户应该有权限？

**A**: 检查以下几点：

1. 确认user_id和tenant_id正确传递
2. 确认用户已分配到角色
3. 确认角色已配置数据权限
4. 确认资源ID正确提取

```go
// 调试：打印用户的所有角色
roles, _ := userRoleRepo.GetRolesByUser(ctx, userID, tenantID)
logs.CtxInfof(ctx, "User roles: %+v", roles)

// 调试：打印角色的数据权限
for _, role := range roles {
    perms, _ := dataPermRepo.GetByRoleAndResource(ctx, role.RoleID, entity.ResourceTypeBots)
    logs.CtxInfof(ctx, "Role %s permissions: %+v", role.RoleCode, perms)
}
```

### Q2: 如何跳过某些路由的权限检查？

**A**: 不添加权限中间件即可。或者创建白名单中间件：

```go
func SkipPermissionCheckForPublic() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        path := c.Request.URL.Path
        if isPublicPath(path) {
            c.Next(ctx)
            return
        }

        // 执行权限检查
        middleware.RequirePermission(config)(ctx, c)
    }
}

func isPublicPath(path string) bool {
    publicPaths := []string{
        "/api/v1/public",
        "/health",
        "/metrics",
    }
    for _, prefix := range publicPaths {
        if strings.HasPrefix(path, prefix) {
            return true
        }
    }
    return false
}
```

### Q3: 如何实现动态权限（例如只能操作自己创建的资源）？

**A**: 使用数据权限的"OWN"（仅自己数据）级别：

```sql
-- 在数据库中配置角色的数据权限为 OWN
INSERT INTO data_permissions (permission_id, role_id, resource_type, scope)
VALUES ('perm_123', 'role_123', 'bots', 'OWN');
```

权限检查器会自动检查资源的`creator_id`字段是否与当前用户ID相同。

### Q4: 字段过滤不生效？

**A**: 检查以下几点：

1. 确保使用结构体而不是map
2. 确保字段有正确的JSON tag
3. 确保调用了`FilterFieldsByPermission`
4. 确保字段权限已配置

```go
// 检查字段权限配置
fieldPerms, err := enhancedPermissionChecker.checker.GetFieldPermissions(ctx, userID, tenantID, "bots")
logs.CtxInfof(ctx, "Field permissions: %+v", fieldPerms)
```

---

**文档变更历史**:

| 版本 | 日期 | 变更内容 | 作者 |
|-----|------|---------|------|
| v1.0 | 2025-01-01 | 初始版本 | 后端架构团队 |

---

**© 2025 ZKER Project. All rights reserved.**
