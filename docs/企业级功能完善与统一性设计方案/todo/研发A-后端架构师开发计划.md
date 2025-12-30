# 研发A - 后端架构师开发计划

**负责人**: 研发A（后端架构师）
**开发周期**: 5周（4人并行）
**总工作量**: 11人天
**核心职责**: RBAC权限系统、权限检查中间件、服务健康检查集成
**最后更新**: 2025-12-30

---

## 📋 工作量总览

| 周次 | 主要任务 | 工作量 | 优先级 |
|------|---------|--------|--------|
| **Week 1-2** | 权限检查中间件实现 | 4人天 | **P0** |
| **Week 2** | 系统预置角色初始化 | 2人天 | **P1** |
| **Week 3** | 协助tenant_id迁移 | 2人天 | **P0** |
| **Week 4** | 服务健康检查集成 | 3人天 | **P1** |
| **Week 5** | 代码审查和Bug修复 | - | - |

**📌 注意**: 你的工作**完全不涉及**前端代码和数据库迁移（除了协助），专注在后端权限和路由领域。

---

## 🎯 Week 1-2: 权限检查中间件实现 (P0 - 阻塞性任务)

### 任务1.1: RequireDataPermission中间件 (2人天)

**文件**: `backend/api/middleware/permission_check_enhanced.go`

**实施步骤**:

#### Day 1: 核心逻辑实现

```go
// RequireDataPermission 数据权限检查中间件
// 支持5级数据权限: ALL, DEPARTMENT_SUB, DEPARTMENT, OWN, CUSTOM, NONE
func RequireDataPermission(resourceType, action string) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 1. 提取user_id和tenant_id
        userID := ctx.GetString("user_id")
        tenantID := ctx.GetString("tenant_id")

        if userID == "" || tenantID == "" {
            c.JSON(consts.StatusUnauthorized, map[string]interface{}{
                "code": "PERM401",
                "message": "缺少用户或租户信息",
                "message_zh": "缺少用户或租户信息",
                "message_en": "Missing user or tenant information",
            })
            c.Abort()
            return
        }

        // 2. 提取resource_id
        resourceID := c.Param("id")
        if resourceID == "" {
            // 尝试从请求体获取
            var body map[string]interface{}
            if err := json.Unmarshal(c.GetRequest().Body(), &body); err == nil {
                if id, ok := body["id"].(string); ok {
                    resourceID = id
                }
            }
        }

        // 3. 调用权限检查服务
        permissionSvc := permissionapp.PermissionAppSVC
        req := &permission.CheckDataPermissionRequest{
            UserID:       userID,
            TenantID:     tenantID,
            ResourceType: resourceType,
            ResourceID:   resourceID,
            Action:       action,
        }

        resp, err := permissionSvc.CheckDataPermission(ctx, req)
        if err != nil || !resp.Data.Allowed {
            c.JSON(consts.StatusForbidden, map[string]interface{}{
                "code": "PERM403",
                "message": "权限不足",
                "message_zh": "权限不足",
                "message_en": "Permission denied",
                "reason": "data_permission_denied",
            })
            c.Abort()
            return
        }

        c.Next(ctx)
    }
}
```

#### Day 2: 单元测试 + 集成测试

```go
// backend/api/middleware/permission_check_enhanced_test.go
package middleware_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestRequireDataPermission_Allowed(t *testing.T) {
    // 测试ALL权限场景
}

func TestRequireDataPermission_OwnOnly(t *testing.T) {
    // 测试OWN权限场景
}

func TestRequireDataPermission_DepartmentOnly(t *testing.T) {
    // 测试DEPARTMENT权限场景
}

func TestRequireDataPermission_Denied(t *testing.T) {
    // 测试拒绝场景
}
```

**⚠️ 注意事项**:
1. **不要修改现有的 `permission_check.go`**，创建新文件 `permission_check_enhanced.go`
2. **严格遵循5级权限定义**，不要自定义权限级别
3. **所有权限检查失败必须返回统一的PERM403错误码**
4. **context传递必须包含request_id**，用于日志追踪

**📖 开发规范**:
- 函数不超过50行
- 错误处理使用 `%w` 包装
- 所有日志必须包含 `request_id`、`user_id`、`tenant_id`
- 单元测试覆盖率 ≥ 80%

**🔗 设计文档链接**:
- [21-用户管理_RBAC细化补充.md](../21-用户管理_RBAC细化补充.md)
- [ZKER-企业级开发规范手册_v1.0.md](../ZKER-企业级开发规范手册_v1.0.md) - 后端开发规范章节

---

### 任务1.2: RequireFieldPermission中间件 (1人天)

**文件**: `backend/api/middleware/permission_check_enhanced.go` (继续在同一文件)

**实施步骤**:

```go
// RequireFieldPermission 字段权限检查中间件
// 支持3级字段权限: hidden, readonly, editable
func RequireFieldPermission(resourceType string) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 1. 提取user_id和tenant_id
        userID := ctx.GetString("user_id")
        tenantID := ctx.GetString("tenant_id")

        // 2. 调用字段权限检查服务
        permissionSvc := permissionapp.PermissionAppSVC
        fields, err := permissionSvc.GetFieldPermissions(ctx, userID, tenantID, resourceType)

        if err != nil {
            // 记录错误但继续请求，使用默认权限
            logs.CtxErrorf(ctx, "[FieldPermission] failed to get field permissions: %v", err)
            fields = getDefaultFieldPermissions()
        }

        // 3. 将字段权限注入context，供后续Handler使用
        ctx.Set("field_permissions", fields)

        c.Next(ctx)
    }
}

// 辅助函数：获取默认字段权限
func getDefaultFieldPermissions() map[string]string {
    return map[string]string{
        "default": "readonly",
    }
}
```

**⚠️ 注意事项**:
1. **字段权限检查失败不应阻塞请求**，使用默认权限
2. **字段权限信息存放在context中**，由Handler层实际过滤
3. **不要在这里删除字段**，只提供权限信息

**📖 开发规范**:
- 失败降级而非阻塞请求
- Context key使用常量定义

**🔗 设计文档链接**:
- [21-用户管理_RBAC细化补充.md](../21-用户管理_RBAC细化补充.md)

---

### 任务1.3: RequireRole中间件 (1人天)

**文件**: `backend/api/middleware/permission_check_enhanced.go` (继续在同一文件)

**实施步骤**:

```go
// RequireRole 角色检查中间件
// 支持多个角色，拥有任一角色即可通过
func RequireRole(roles ...string) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 1. 提取user_id和tenant_id
        userID := ctx.GetString("user_id")
        tenantID := ctx.GetString("tenant_id")

        // 2. 获取用户所有角色
        permissionSvc := permissionapp.PermissionAppSVC
        userRoles, err := permissionSvc.GetUserRoles(ctx, userID, tenantID)

        if err != nil {
            c.JSON(consts.StatusInternalServerError, map[string]interface{}{
                "code": "PERM500",
                "message": "获取用户角色失败",
                "message_zh": "获取用户角色失败",
                "message_en": "Failed to get user roles",
            })
            c.Abort()
            return
        }

        // 3. 检查是否拥有所需角色
        hasRole := false
        for _, requiredRole := range roles {
            for _, userRole := range userRoles.Data.Roles {
                if userRole.RoleCode == requiredRole {
                    hasRole = true
                    break
                }
            }
            if hasRole {
                break
            }
        }

        if !hasRole {
            c.JSON(consts.StatusForbidden, map[string]interface{}{
                "code": "PERM403",
                "message": "需要以下角色之一：" + strings.Join(roles, ","),
                "message_zh": "需要以下角色之一：" + strings.Join(roles, ","),
                "message_en": "One of the following roles is required: " + strings.Join(roles, ","),
                "required_roles": roles,
            })
            c.Abort()
            return
        }

        c.Next(ctx)
    }
}
```

**⚠️ 注意事项**:
1. **支持多个角色**（OR逻辑），拥有任一角色即可
2. **角色不存在不是错误**，只是没有权限
3. **错误消息返回required_roles**，便于前端提示

**📖 开发规范**:
- 错误消息清晰具体
- 包含required_roles信息

**🔗 设计文档链接**:
- [API接口文档_用户管理RBAC.md](../API接口文档_用户管理RBAC.md)

---

## 🎯 Week 2: 系统预置角色初始化 (P1)

### 任务2.1: 创建系统预置角色SQL脚本 (0.5人天)

**文件**: `backend/domain/permission/migration/001_init_system_roles.sql`

**实施步骤**:

```sql
-- ============================================
-- 系统预置角色初始化脚本
-- 执行时机: 新租户创建时
-- ============================================

-- 1. 创建系统预置角色
INSERT INTO roles (role_id, tenant_id, role_name, role_code, role_type, description, created_at) VALUES
('role_tenant_owner_sys', 'system', '租户所有者', 'tenant_owner', 'system', '拥有租户内所有资源的完整权限', NOW()),
('role_tenant_admin_sys', 'system', '租户管理员', 'tenant_admin', 'system', '拥有租户内部门及以下资源的完整权限', NOW()),
('role_tenant_member_sys', 'system', '普通成员', 'tenant_member', 'system', '仅能访问自己创建的资源', NOW()),
('role_tenant_viewer_sys', 'system', '只读成员', 'tenant_viewer', 'system', '仅能查看自己创建的资源', NOW());

-- 2. 数据权限默认配置
-- tenant_owner: 所有资源类型都是ALL权限
INSERT INTO data_permissions (permission_id, role_id, resource_type, scope, created_at) VALUES
('perm_tenant_owner_bots', 'role_tenant_owner_sys', 'bots', 'ALL', NOW()),
('perm_tenant_owner_conversations', 'role_tenant_owner_sys', 'conversations', 'ALL', NOW()),
('perm_tenant_owner_knowledge', 'role_tenant_owner_sys', 'knowledge', 'ALL', NOW()),
('perm_tenant_owner_workflows', 'role_tenant_owner_sys', 'workflows', 'ALL', NOW());

-- tenant_admin: 部门及以下权限
INSERT INTO data_permissions (permission_id, role_id, resource_type, scope, created_at) VALUES
('perm_tenant_admin_bots', 'role_tenant_admin_sys', 'bots', 'DEPARTMENT_SUB', NOW()),
('perm_tenant_admin_conversations', 'role_tenant_admin_sys', 'conversations', 'DEPARTMENT_SUB', NOW()),
('perm_tenant_admin_knowledge', 'role_tenant_admin_sys', 'knowledge', 'DEPARTMENT_SUB', NOW()),
('perm_tenant_admin_workflows', 'role_tenant_admin_sys', 'workflows', 'DEPARTMENT_SUB', NOW());

-- tenant_member: 仅自己的资源
INSERT INTO data_permissions (permission_id, role_id, resource_type, scope, created_at) VALUES
('perm_tenant_member_bots', 'role_tenant_member_sys', 'bots', 'OWN', NOW()),
('perm_tenant_member_conversations', 'role_tenant_member_sys', 'conversations', 'OWN', NOW()),
('perm_tenant_member_knowledge', 'role_tenant_member_sys', 'knowledge', 'OWN', NOW()),
('perm_tenant_member_workflows', 'role_tenant_member_sys', 'workflows', 'OWN', NOW());

-- tenant_viewer: 仅自己的资源（只读）
INSERT INTO data_permissions (permission_id, role_id, resource_type, scope, created_at) VALUES
('perm_tenant_viewer_bots', 'role_tenant_viewer_sys', 'bots', 'OWN', NOW()),
('perm_tenant_viewer_conversations', 'role_tenant_viewer_sys', 'conversations', 'OWN', NOW()),
('perm_tenant_viewer_knowledge', 'role_tenant_viewer_sys', 'knowledge', 'OWN', NOW()),
('perm_tenant_viewer_workflows', 'role_tenant_viewer_sys', 'workflows', 'OWN', NOW());
```

**⚠️ 注意事项**:
1. **使用 `system` 作为tenant_id**，表示系统级角色
2. **role_type = 'system'**，这些角色不能被删除
3. **每个角色配置4种资源类型**的权限：bots, conversations, knowledge, workflows
4. **不要修改已存在的role_id**，保持唯一性

**📖 开发规范**:
- SQL注释清晰
- 使用幂等性设计（可重复执行）
- 外键约束完整

**🔗 设计文档链接**:
- [数据库设计完整交付清单.md](../数据库设计完整交付清单.md)

---

### 任务2.2: Go初始化函数实现 (1人天)

**文件**: `backend/domain/permission/service/role_init.go`

**实施步骤**:

```go
// backend/domain/permission/service/role_init.go
package service

import (
    "context"
    "github.com/coze-studio/backend/domain/permission/entity"
    "github.com/coze-studio/backend/domain/permission/repository"
)

// InitializeSystemRoles 初始化系统预置角色
// 调用时机: 新租户创建时
func (s *RoleService) InitializeSystemRoles(ctx context.Context, tenantID string) error {
    // 1. 检查是否已初始化
    existingRoles, err := s.repo.GetByTenant(ctx, tenantID)
    if err == nil && len(existingRoles) > 0 {
        logs.CtxInfof(ctx, "[RoleInit] roles already initialized for tenant %s", tenantID)
        return nil // 已初始化，幂等性
    }

    // 2. 创建系统预置角色
    roles := []*entity.Role{
        {
            RoleID:     generateRoleID("tenant_owner", tenantID),
            TenantID:   tenantID,
            RoleName:   "租户所有者",
            RoleCode:   "tenant_owner",
            RoleType:   entity.RoleTypeSystem,
            Description: "拥有租户内所有资源的完整权限",
        },
        {
            RoleID:     generateRoleID("tenant_admin", tenantID),
            TenantID:   tenantID,
            RoleName:   "租户管理员",
            RoleCode:   "tenant_admin",
            RoleType:   entity.RoleTypeSystem,
            Description: "拥有租户内部门及以下资源的完整权限",
        },
        {
            RoleID:     generateRoleID("tenant_member", tenantID),
            TenantID:   tenantID,
            RoleName:   "普通成员",
            RoleCode:   "tenant_member",
            RoleType:   entity.RoleTypeSystem,
            Description: "仅能访问自己创建的资源",
        },
        {
            RoleID:     generateRoleID("tenant_viewer", tenantID),
            TenantID:   tenantID,
            RoleName:   "只读成员",
            RoleCode:   "tenant_viewer",
            RoleType:   entity.RoleTypeSystem,
            Description: "仅能查看自己创建的资源",
        },
    }

    // 3. 批量创建角色
    for _, role := range roles {
        if err := s.repo.Create(ctx, role); err != nil {
            logs.CtxErrorf(ctx, "[RoleInit] failed to create role %s: %v", role.RoleCode, err)
            return err
        }

        // 4. 为每个角色初始化数据权限
        if err := s.initializeDataPermissions(ctx, role); err != nil {
            logs.CtxErrorf(ctx, "[RoleInit] failed to init data permissions for role %s: %v", role.RoleCode, err)
            return err
        }
    }

    logs.CtxInfof(ctx, "[RoleInit] successfully initialized %d system roles for tenant %s", len(roles), tenantID)
    return nil
}

// initializeDataPermissions 初始化角色的数据权限
func (s *RoleService) initializeDataPermissions(ctx context.Context, role *entity.Role) error {
    // 根据角色类型配置数据权限
    var scope entity.DataPermissionScope
    switch role.RoleCode {
    case "tenant_owner":
        scope = entity.DataPermissionScopeAll
    case "tenant_admin":
        scope = entity.DataPermissionScopeDepartmentSub
    case "tenant_member", "tenant_viewer":
        scope = entity.DataPermissionScopeOwn
    default:
        scope = entity.DataPermissionScopeNone
    }

    // 为每种资源类型创建数据权限
    resourceTypes := []entity.ResourceType{
        entity.ResourceTypeBots,
        entity.ResourceTypeConversations,
        entity.ResourceTypeKnowledge,
        entity.ResourceTypeWorkflows,
    }

    for _, rt := range resourceTypes {
        perm := &entity.DataPermission{
            PermissionID: generatePermissionID(role.RoleID, string(rt)),
            RoleID:       role.RoleID,
            ResourceType: rt,
            Scope:        scope,
        }
        if err := s.dataPermRepo.Create(ctx, perm); err != nil {
            return err
        }
    }

    return nil
}

// generateRoleID 生成角色ID
func generateRoleID(roleCode, tenantID string) string {
    return fmt.Sprintf("role_%s_%s", roleCode, tenantID[:8])
}

// generatePermissionID 生成权限ID
func generatePermissionID(roleID, resourceType string) string {
    return fmt.Sprintf("perm_%s_%s", roleID, resourceType)
}
```

**⚠️ 注意事项**:
1. **幂等性设计**：如果角色已存在，不重复创建
2. **事务性**：角色创建失败要回滚，保证数据一致性
3. **日志完整**：每个步骤都要记录日志，便于排查问题
4. **不要硬编码role_id**，使用函数生成

**📖 开发规范**:
- 函数命名清晰（Initialize*）
- 错误处理完整
- 日志包含上下文信息

**🔗 设计文档链接**:
- [21-用户管理_RBAC细化补充.md](../21-用户管理_RBAC细化补充.md)
- [ZKER-企业级开发规范手册_v1.0.md](../ZKER-企业级开发规范手册_v1.0.md) - 领域层设计章节

---

### 任务2.3: 集成到租户创建流程 (0.5人天)

**修改文件**: `backend/application/tenant/tenant_service.go`

**实施步骤**:

```go
// 在CreateTenant函数中添加角色初始化
func (s *TenantAppSVC) CreateTenant(ctx context.Context, req *tenant.CreateTenantRequest) (*tenant.CreateTenantResponse, error) {
    // ... 现有租户创建逻辑 ...

    // 新增：初始化系统预置角色
    if err := s.roleSvc.InitializeSystemRoles(ctx, tenant.TenantID); err != nil {
        logs.CtxErrorf(ctx, "[TenantSvc] failed to initialize system roles for tenant %s: %v", tenant.TenantID, err)
        // 回滚租户创建
        s.tenantRepo.Delete(ctx, tenant.TenantID)
        return nil, fmt.Errorf("初始化系统角色失败: %w", err)
    }

    logs.CtxInfof(ctx, "[TenantSvc] tenant %s created with system roles initialized", tenant.TenantID)

    // ... 返回响应 ...
}
```

**⚠️ 注意事项**:
1. **角色初始化失败必须回滚整个租户创建**
2. **不要使用异步初始化**，必须同步完成
3. **错误日志要包含tenant_id**

**📖 开发规范**:
- 事务完整性
- 失败回滚机制

**🔗 设计文档链接**:
- [21-MultiTenant_SaaS核心_租户识别与管理.md](../21-MultiTenant_SaaS核心_租户识别与管理.md)

---

## 🎯 Week 3: 协助tenant_id迁移 (P0)

### 任务3.1: 代码审查和架构指导 (2人天)

**职责范围**: 协助研发B进行数据迁移，但不直接写迁移代码

**具体任务**:

#### Day 1: 迁移架构审查

1. **审查迁移脚本设计**
   - 检查分批处理逻辑
   - 检查幂等性实现
   - 检查回滚方案
   - 确认数据一致性验证方法

2. **评估性能影响**
   - 迁移期间对线上服务的影响
   - 锁表时间评估
   - 索引创建策略

3. **风险评估**
   - 识别潜在的数据丢失风险
   - 识别业务中断风险
   - 制定缓解措施

**⚠️ 注意事项**:
1. **不要直接修改迁移脚本**，只提供审查意见
2. **重点关注数据一致性和回滚方案**
3. **确保迁移期间服务可用性**

**🔗 设计文档链接**:
- [ZKER-数据迁移方案_v1.0.md](../ZKER-数据迁移方案_v1.0.md)
- [Session和User表tenant_id字段迁移方案_v1.0.md](../Session和User表tenant_id字段迁移方案_v1.0.md)

#### Day 2: 双写验证指导

1. **审查双写逻辑**
   - 检查双写时机正确性
   - 检查双写失败处理
   - 确认数据对比方法

2. **监控指标设计**
   - 双写成功率监控
   - 数据不一致告警
   - 性能指标监控

**⚠️ 注意事项**:
1. **双写期间要有详细的监控**
2. **数据不一致要有告警**
3. **准备快速回滚方案**

**🔗 设计文档链接**:
- [ZKER-数据迁移方案_v1.0.md](../ZKER-数据迁移方案_v1.0.md) - 双写验证章节

---

## 🎯 Week 4: 服务健康检查集成 (P1)

### 任务4.1: 集成健康检查到路由决策 (3人天)

**文件**: `backend/domain/routing/service/health_integration.go`

**实施步骤**:

#### Day 1: 健康检查集成核心逻辑

```go
// backend/domain/routing/service/health_integration.go
package service

import (
    "context"
    "time"
)

// IntegrateHealthCheck 集成健康检查到路由决策
func (e *RoutingEngine) IntegrateHealthCheck(ctx context.Context, tenantID string) error {
    // 1. 获取所有活跃的Bot
    bots, err := e.botRepo.GetActiveBots(ctx, tenantID)
    if err != nil {
        return fmt.Errorf("获取活跃Bot失败: %w", err)
    }

    // 2. 批量检查Bot健康状态
    for _, bot := range bots {
        healthStatus := e.healthSVC.CheckHealth(ctx, bot.BotID)

        // 3. 更新Bot健康状态指标
        metrics.UpdateBotHealth(
            tenantID,
            bot.BotID,
            healthStatus.IsHealthy,
            healthStatus.SuccessRate,
            healthStatus.CurrentLoad,
            healthStatus.MaxCapacity,
        )

        logs.CtxInfof(ctx, "[RoutingEngine] bot %s health updated: healthy=%v, success_rate=%.2f, load=%d/%d",
            bot.BotID, healthStatus.IsHealthy, healthStatus.SuccessRate,
            healthStatus.CurrentLoad, healthStatus.MaxCapacity)
    }

    return nil
}

// GetHealthyBots 获取健康的Bot列表
// 过滤条件: is_healthy=true AND success_rate>=0.8
func (e *RoutingEngine) GetHealthyBots(ctx context.Context, tenantID string) ([]*Bot, error) {
    // 1. 获取所有Bot
    bots, err := e.botRepo.GetByTenant(ctx, tenantID)
    if err != nil {
        return nil, err
    }

    // 2. 过滤健康的Bot
    healthyBots := make([]*Bot, 0)
    for _, bot := range bots {
        healthStatus := e.healthSVC.CheckHealth(ctx, bot.BotID)
        if healthStatus.IsHealthy && healthStatus.SuccessRate >= 0.8 {
            healthyBots = append(healthyBots, bot)
        } else {
            logs.CtxWarnf(ctx, "[RoutingEngine] bot %s is unhealthy: healthy=%v, success_rate=%.2f",
                bot.BotID, healthStatus.IsHealthy, healthStatus.SuccessRate)
        }
    }

    logs.CtxInfof(ctx, "[RoutingEngine] got %d healthy bots out of %d total bots",
        len(healthyBots), len(bots))

    return healthyBots, nil
}

// UpdateBotMetrics 定期更新Bot指标
// 应该由定时任务调用，每分钟执行一次
func (e *RoutingEngine) UpdateBotMetrics(ctx context.Context) error {
    // 1. 获取所有租户
    tenants, err := e.tenantRepo.ListActive(ctx)
    if err != nil {
        return fmt.Errorf("获取活跃租户列表失败: %w", err)
    }

    // 2. 遍历每个租户
    for _, tenant := range tenants {
        // 3. 更新该租户下所有Bot的指标
        if err := e.IntegrateHealthCheck(ctx, tenant.TenantID); err != nil {
            // 记录日志但继续处理下一个租户
            logs.CtxErrorf(ctx, "[RoutingEngine] failed to update bot metrics for tenant %s: %v",
                tenant.TenantID, err)
        }
    }

    return nil
}
```

**⚠️ 注意事项**:
1. **健康检查失败不应阻塞路由**，记录日志并继续
2. **每分钟执行一次**，不要过于频繁
3. **健康阈值：success_rate >= 0.8**，可配置
4. **不要同步调用健康检查**，避免阻塞路由决策

**📖 开发规范**:
- 错误处理优雅降级
- 日志详细
- 性能优化（批量处理）

**🔗 设计文档链接**:
- [16-五大AI引擎核心_智能路由引擎.md](../16-五大AI引擎核心_智能路由引擎.md)
- [24-租户监控运维_Agent监控补充.md](../24-租户监控运维_Agent监控补充.md)

#### Day 2: 路由决策集成健康检查

**修改文件**: `backend/domain/routing/service/routing_engine.go`

```go
// 在ExecuteRouting函数中集成健康检查
func (e *RoutingEngine) ExecuteRouting(ctx context.Context, req *RoutingRequest) (*RoutingDecision, error) {
    // ... 现有意图匹配逻辑 ...

    // 新增：只从健康的Bot列表中选择
    healthyBots, err := e.GetHealthyBots(ctx, req.TenantID)
    if err != nil {
        logs.CtxErrorf(ctx, "[RoutingEngine] failed to get healthy bots: %v", err)
        // 降级：使用所有Bot
        healthyBots, _ = e.botRepo.GetByTenant(ctx, req.TenantID)
    }

    if len(healthyBots) == 0 {
        return nil, fmt.Errorf("没有可用的健康Bot")
    }

    // ... 评分路由逻辑，仅从healthyBots中选择 ...

    return decision, nil
}
```

**⚠️ 注意事项**:
1. **健康检查失败要降级**，使用所有Bot
2. **没有健康Bot要返回错误**，不要返回空结果
3. **日志记录降级原因**

**🔗 设计文档链接**:
- [16-五大AI引擎核心_智能路由引擎.md](../16-五大AI引擎核心_智能路由引擎.md)

#### Day 3: 定时任务配置 + 单元测试

**文件**: `backend/cmd/server.go`

```go
// 启动Bot健康检查定时任务
func startBotHealthCheckTask(routingEngine *routingservice.RoutingEngine) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    logs.Info("[HealthCheck] bot health check task started")

    for range ticker.C {
        ctx := context.Background()
        if err := routingEngine.UpdateBotMetrics(ctx); err != nil {
            logs.Errorf("[HealthCheck] failed to update bot metrics: %v", err)
        }
    }
}

func main() {
    // ... 现有初始化代码 ...

    // 启动Bot健康检查定时任务
    go startBotHealthCheckTask(routingEngine)

    // ... 启动HTTP服务 ...
}
```

**单元测试**:

```go
// backend/domain/routing/service/health_integration_test.go
func TestGetHealthyBots(t *testing.T) {
    // 测试健康Bot过滤逻辑
}

func TestIntegrateHealthCheck(t *testing.T) {
    // 测试健康检查集成
}

func TestUpdateBotMetrics(t *testing.T) {
    // 测试指标更新
}
```

**⚠️ 注意事项**:
1. **定时任务在独立goroutine中运行**
2. **ticker要记得停止**
3. **单元测试覆盖率 ≥ 80%**

**📖 开发规范**:
- 并发安全
- 资源清理
- 完整的单元测试

**🔗 设计文档链接**:
- [ZKER-企业级开发规范手册_v1.0.md](../ZKER-企业级开发规范手册_v1.0.md) - 并发编程章节

---

## 🎯 Week 5: 代码审查和Bug修复

### 任务5.1: 交叉代码审查 (2天)

**审查范围**:

1. **研发B的代码**:
   - tenant_id迁移脚本
   - 配额检查中间件
   - 错误码定义

2. **研发C的前端代码**:
   - API调用层（不审查UI部分）
   - 类型定义

**审查重点**:
- [ ] DDD分层架构是否正确
- [ ] 错误处理是否完整
- [ ] 日志是否包含必要上下文
- [ ] 并发安全性
- [ ] 命名规范
- [ ] 函数复杂度（不超过50行）

**🔗 设计文档链接**:
- [ZKER-全局一致性检查清单_v1.0.md](../ZKER-全局一致性检查清单_v1.0.md)

### 任务5.2: Bug修复和优化 (3天)

**修复自己代码的Bug**
**优化性能瓶颈**
**补充单元测试**

---

## 📊 关键文档链接汇总

### 必读文档（开发前必读）

1. **[ZKER-企业级开发规范手册_v1.0.md](../ZKER-企业级开发规范手册_v1.0.md)** ⭐⭐⭐
   - 后端开发规范章节
   - DDD分层架构章节
   - 并发编程章节

2. **[ZKER-全局一致性检查清单_v1.0.md](../ZKER-全局一致性检查清单_v1.0.md)** ⭐⭐⭐
   - L1个人检查清单
   - L2模块检查清单

3. **[21-用户管理_RBAC细化补充.md](../21-用户管理_RBAC细化补充.md)** ⭐⭐⭐
   - 5级数据权限定义
   - 3级字段权限定义
   - 13张新增表结构

4. **[ZKER-数据迁移方案_v1.0.md](../ZKER-数据迁移方案_v1.0.md)** ⭐⭐
   - 数据迁移流程
   - 双写验证方案

### 参考文档

5. **[API接口文档_用户管理RBAC.md](../API接口文档_用户管理RBAC.md)** - RBAC API设计
6. **[16-五大AI引擎核心_智能路由引擎.md](../16-五大AI引擎核心_智能路由引擎.md)** - 智能路由设计
7. **[24-租户监控运维_Agent监控补充.md](../24-租户监控运维_Agent监控补充.md)** - Agent监控指标

---

## ⚠️ 核心注意事项

### 1. 代码隔离原则

**你的职责范围**:
- ✅ RBAC权限系统（后端）
- ✅ 权限检查中间件
- ✅ 系统预置角色初始化
- ✅ 服务健康检查集成
- ✅ 代码审查

**不要触碰**:
- ❌ 前端代码（研发C负责）
- ❌ 数据库迁移脚本（研发B负责，你只审查）
- ❌ CI/CD配置（研发D负责）
- ❌ 监控大盘配置（研发D负责）

### 2. 开发规范红线

**必须遵守**:
1. 所有函数不超过50行
2. 所有错误使用统一错误码
3. 所有日志包含request_id、user_id、tenant_id
4. 单元测试覆盖率 ≥ 80%
5. 代码必须通过golangci-lint检查

**禁止行为**:
1. 使用panic（必须用error返回）
2. 硬编码配置（必须从配置文件读取）
3. 直接修改生产代码（先在测试环境验证）
4. 提交未测试的代码

### 3. Git提交规范

**Commit Message格式**:
```
feat(permission): add RequireDataPermission middleware

- Implement 5-level data permission checking
- Add unit tests with 80%+ coverage
- Support resource types: bots, conversations, knowledge, workflows

Refs: #123
```

**提交前检查清单**:
- [ ] 代码通过所有单元测试
- [ ] golangci-lint无警告
- [ ] 代码符合开发规范
- [ ] 自我Code Review通过
- [ ] 更新相关文档

---

## 📅 每日工作检查清单

### 开发前
- [ ] 阅读相关设计文档
- [ ] 理解API契约和数据模型
- [ ] 确认依赖关系
- [ ] 拉取最新代码

### 开发中
- [ ] 遵循SOLID、KISS、DRY、YAGNI原则
- [ ] 每个函数 < 50行
- [ ] 错误处理完整
- [ ] 日志包含上下文
- [ ] 边开发边写单元测试

### 提交前
- [ ] 所有测试通过（go test ./...）
- [ ] 代码覆盖率 ≥ 80%
- [ ] golangci-lint run通过
- [ ] 自我Code Review
- [ ] 更新文档

---

## 🎯 成功标准

### Week 1-2结束时
- [ ] 3个权限检查中间件全部实现
- [ ] 单元测试覆盖率 ≥ 80%
- [ ] 中间件集成到路由
- [ ] 代码审查通过

### Week 2结束时
- [ ] 系统预置角色SQL脚本完成
- [ ] Go初始化函数完成
- [ ] 集成到租户创建流程
- [ ] 新租户自动创建角色

### Week 4结束时
- [ ] 健康检查集成到路由决策
- [ ] 定时任务正常运行
- [ ] 路由引擎优先选择健康Bot
- [ ] Prometheus指标正确更新

### Week 5结束时
- [ ] 所有代码审查完成
- [ ] 所有Bug修复完成
- [ ] 所有单元测试通过
- [ ] 文档完整交付

---

**🎉 祝你开发顺利！遇到问题随时查阅设计文档或与团队讨论。**
