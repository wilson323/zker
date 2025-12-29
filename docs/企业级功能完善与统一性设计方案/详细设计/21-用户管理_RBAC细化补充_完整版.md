# 21-用户管理_RBAC细化补充完整版

**模块**: 21-用户管理扩展
**扩展内容**: RBAC权限细化、数据权限、字段权限、临时授权
**版本**: v2.0
**日期**: 2025-01-03
**优先级**: P0
**工期**: 3周

---

## 目录

- [1. 功能概述](#1-功能概述)
- [2. 数据库设计](#2-数据库设计)
- [3. 数据权限控制](#3-数据权限控制)
- [4. 字段权限控制](#4-字段权限控制)
- [5. 临时授权机制](#5-临时授权机制)
- [6. 权限验证引擎](#6-权限验证引擎)
- [7. API设计](#7-api设计)
- [8. 前端组件](#8-前端组件)
- [9. 实施计划](#9-实施计划)

---

## 1. 功能概述

### 1.1 核心目标

**RBAC权限细化** 提供企业级的细粒度权限控制能力：

- ✅ **5级数据权限** - 全部/本部门及下级/本部门/仅本人/自定义
- ✅ **3级字段权限** - 可见/脱敏/隐藏
- ✅ **临时授权** - 支持临时权限授予与回收
- ✅ **权限继承** - 基于组织架构的权限继承
- ✅ **权限审计** - 完整的权限变更日志

### 1.2 应用场景

**场景1：数据权限隔离**
```
部门经理查看用户列表：
- 可见：本部门及下级部门的用户（50人）
- 不可见：其他部门用户
- 操作：可编辑、可删除本部门用户
```

**场景2：敏感字段脱敏**
```
普通员工查看客户信息：
- 客户姓名：张三（可见）
- 手机号：138****1234（脱敏）
- 身份证号：3301**********1234（脱敏）
- 银行卡号：****（隐藏）
```

**场景3：临时权限授予**
```
项目负责人临时授予访客权限：
- 授予人：管理员
- 接受人：外部顾问
- 权限：Bot编辑权限
- 有效期：3天
- 自动过期后权限回收
```

---

## 2. 数据库设计

### 2.1 数据权限表

```sql
CREATE TABLE data_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '权限ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    resource_type VARCHAR(64) NOT NULL COMMENT '资源类型 (user/bot/knowledge_base)',

    -- 权限范围
    permission_scope ENUM('all', 'department_and_below', 'department', 'self', 'custom') NOT NULL COMMENT '权限范围',
    custom_org_ids JSON COMMENT '自定义组织ID列表 [1,2,3]',
    custom_user_ids JSON COMMENT '自定义用户ID列表 [100,200,300]',

    -- 权限操作
    permissions JSON COMMENT '权限操作 ["read","write","delete"]',

    -- 元数据
    created_by BIGINT COMMENT '创建人ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_tenant_role (tenant_id, role_id),
    INDEX idx_resource_type (resource_type),
    UNIQUE KEY uk_role_resource (role_id, resource_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='数据权限表';
```

### 2.2 字段权限表

```sql
CREATE TABLE field_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    resource_type VARCHAR(64) NOT NULL COMMENT '资源类型',

    -- 字段配置
    field_name VARCHAR(64) NOT NULL COMMENT '字段名称',
    field_display_name VARCHAR(100) COMMENT '字段显示名称',
    permission ENUM('visible', 'masked', 'hidden') NOT NULL COMMENT '权限级别',

    -- 脱敏规则
    mask_rule VARCHAR(128) COMMENT '脱敏规则 (如:手机号:138****1234)',
    mask_pattern VARCHAR(128) COMMENT '正则表达式模式',
    mask_replacement VARCHAR(32) DEFAULT '****' COMMENT '替换字符',

    -- 元数据
    created_by BIGINT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_role (tenant_id, role_id),
    INDEX idx_resource_type (resource_type),
    UNIQUE KEY uk_role_resource_field (role_id, resource_type, field_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='字段权限表';
```

### 2.3 临时授权表

```sql
CREATE TABLE temporary_grants (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    grant_code VARCHAR(64) NOT NULL UNIQUE COMMENT '授权码',

    -- 授权关系
    grantor_id BIGINT NOT NULL COMMENT '授予人ID',
    grantee_id BIGINT NOT NULL COMMENT '接受人ID',

    -- 权限配置
    permission_type ENUM('role', 'resource', 'operation') NOT NULL COMMENT '权限类型',
    permission_config JSON NOT NULL COMMENT '权限配置',
    reason VARCHAR(500) COMMENT '授权原因',

    -- 时间
    start_time DATETIME NOT NULL COMMENT '开始时间',
    end_time DATETIME NOT NULL COMMENT '结束时间',

    -- 状态
    is_revoked BOOLEAN DEFAULT FALSE COMMENT '是否已撤销',
    revoked_at DATETIME COMMENT '撤销时间',
    revoked_by BIGINT COMMENT '撤销人ID',
    revoke_reason VARCHAR(500) COMMENT '撤销原因',

    -- 审计
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_grantee (tenant_id, grantee_id),
    INDEX idx_grant_code (grant_code),
    INDEX idx_time_range (start_time, end_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='临时授权表';
```

### 2.4 权限审计日志表

```sql
CREATE TABLE permission_audit_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',

    -- 操作信息
    operator_id BIGINT NOT NULL COMMENT '操作人ID',
    operator_type ENUM('user', 'system') NOT NULL COMMENT '操作人类型',
    operation_type ENUM('grant', 'revoke', 'modify', 'check') NOT NULL COMMENT '操作类型',

    -- 权限信息
    target_type ENUM('role', 'user', 'resource') NOT NULL COMMENT '目标类型',
    target_id VARCHAR(64) NOT NULL COMMENT '目标ID',
    permission_detail JSON COMMENT '权限详情',

    -- 变更信息
    old_value JSON COMMENT '变更前值',
    new_value JSON COMMENT '变更后值',

    -- 结果
    operation_result ENUM('success', 'failed') NOT NULL,
    error_message TEXT COMMENT '错误信息',

    -- 元数据
    ip_address VARCHAR(50) COMMENT 'IP地址',
    user_agent VARCHAR(500) COMMENT 'User-Agent',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_operator (tenant_id, operator_id),
    INDEX idx_operation_type (operation_type),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='权限审计日志表';
```

---

## 3. 数据权限控制

### 3.1 数据权限服务

```go
package rbac

import (
    "context"
    "github.com/casbin/casbin/v2"
)

// DataPermissionService 数据权限服务
type DataPermissionService struct {
    dataPermRepo  repository.DataPermissionRepository
    orgRepo       repository.OrganizationRepository
    enforcer      *casbin.Enforcer
}

// FilterDataByPermission 根据权限过滤数据
func (s *DataPermissionService) FilterDataByPermission(
    ctx context.Context,
    userID int64,
    resourceType string,
    dataIDs []interface{},
) ([]interface{}, error) {
    // 1. 获取用户的角色
    roles, _ := s.getUserRoles(ctx, userID)

    // 2. 获取每个角色的数据权限
    var allowedIDs []interface{}
    hasAllPermission := false

    for _, role := range roles {
        perm, _ := s.dataPermRepo.GetByRoleAndResource(ctx, role.ID, resourceType)

        if perm == nil {
            continue
        }

        switch perm.PermissionScope {
        case "all":
            // 全部数据权限
            hasAllPermission = true
            break

        case "department_and_below":
            // 本部门及下级部门
            orgIDs := s.getDepartmentAndBelow(ctx, userID)
            allowedIDs = append(allowedIDs, s.filterByOrgs(ctx, dataIDs, orgIDs)...)

        case "department":
            // 仅本部门
            orgID := s.getUserDepartment(ctx, userID)
            orgIDs := []int64{orgID}
            allowedIDs = append(allowedIDs, s.filterByOrgs(ctx, dataIDs, orgIDs)...)

        case "self":
            // 仅本人数据
            allowedIDs = append(allowedIDs, s.filterByOwner(ctx, dataIDs, userID)...)

        case "custom":
            // 自定义范围
            if perm.CustomOrgIDs != nil {
                allowedIDs = append(allowedIDs, s.filterByOrgs(ctx, dataIDs, perm.CustomOrgIDs)...)
            }
            if perm.CustomUserIDs != nil {
                allowedIDs = append(allowedIDs, s.filterByUsers(ctx, dataIDs, perm.CustomUserIDs)...)
            }
        }
    }

    // 如果有全部数据权限，直接返回所有数据
    if hasAllPermission {
        return dataIDs, nil
    }

    // 去重
    return uniqueIDs(allowedIDs), nil
}

// CheckDataPermission 检查数据权限
func (s *DataPermissionService) CheckDataPermission(
    ctx context.Context,
    userID int64,
    resourceType string,
    resourceID interface{},
    permission string, // read/write/delete
) (bool, error) {
    // 1. 获取用户的角色
    roles, _ := s.getUserRoles(ctx, userID)

    // 2. 检查每个角色
    for _, role := range roles {
        perm, _ := s.dataPermRepo.GetByRoleAndResource(ctx, role.ID, resourceType)

        if perm == nil {
            continue
        }

        // 检查操作权限
        if !s.hasOperationPermission(perm.Permissions, permission) {
            continue
        }

        // 检查数据范围
        hasAccess, _ := s.checkDataScope(ctx, perm, userID, resourceID)
        if hasAccess {
            return true, nil
        }
    }

    return false, nil
}

// checkDataScope 检查数据范围
func (s *DataPermissionService) checkDataScope(
    ctx context.Context,
    perm *entity.DataPermission,
    userID int64,
    resourceID interface{},
) (bool, error) {
    // 获取资源的组织ID和所有者ID
    orgID, ownerID, _ := s.getResourceInfo(ctx, resourceID)

    switch perm.PermissionScope {
    case "all":
        return true, nil

    case "department_and_below":
        userOrgID := s.getUserDepartment(ctx, userID)
        // 检查资源组织是否是用户部门或下级部门
        return s.isDepartmentOrBelow(ctx, userOrgID, orgID), nil

    case "department":
        userOrgID := s.getUserDepartment(ctx, userID)
        return userOrgID == orgID, nil

    case "self":
        return ownerID == userID, nil

    case "custom":
        // 检查自定义组织列表
        if perm.CustomOrgIDs != nil {
            for _, customOrgID := range perm.CustomOrgIDs {
                if customOrgID == orgID {
                    return true, nil
                }
            }
        }
        // 检查自定义用户列表
        if perm.CustomUserIDs != nil {
            for _, customUserID := range perm.CustomUserIDs {
                if customUserID == ownerID {
                    return true, nil
                }
            }
        }
        return false, nil
    }

    return false, nil
}

// getDepartmentAndBelow 获取部门及下级部门ID列表
func (s *DataPermissionService) getDepartmentAndBelow(
    ctx context.Context,
    userID int64,
) []int64 {
    // 1. 获取用户所在部门
    userDeptID := s.getUserDepartment(ctx, userID)

    // 2. 递归获取所有下级部门
    deptIDs := []int64{userDeptID}
    childDepts := s.getChildDepartments(ctx, userDeptID)
    deptIDs = append(deptIDs, childDepts...)

    return deptIDs
}
```

### 3.2 组织服务集成

```go
package rbac

// OrganizationService 组织服务
type OrganizationService struct {
    orgRepo repository.OrganizationRepository
    cache   *redis.Client
}

// isDepartmentOrBelow 检查是否是部门或下级部门
func (s *OrganizationService) isDepartmentOrBelow(
    ctx context.Context,
    parentDeptID, targetDeptID int64,
) bool {
    // 检查是否是同一部门
    if parentDeptID == targetDeptID {
        return true
    }

    // 递归检查下级部门
    childDepts := s.getChildDepartments(ctx, parentDeptID)
    for _, childID := range childDepts {
        if childID == targetDeptID {
            return true
        }
        // 递归检查
        if s.isDepartmentOrBelow(ctx, childID, targetDeptID) {
            return true
        }
    }

    return false
}

// getChildDepartments 获取子部门ID列表
func (s *OrganizationService) getChildDepartments(
    ctx context.Context,
    parentDeptID int64,
) []int64 {
    // 使用缓存加速查询
    cacheKey := fmt.Sprintf("org:children:%d", parentDeptID)
    cached, _ := s.cache.Get(ctx, cacheKey).Result()

    if cached != "" {
        var deptIDs []int64
        json.Unmarshal([]byte(cached), &deptIDs)
        return deptIDs
    }

    // 查询数据库
    depts, _ := s.orgRepo.GetByParentID(ctx, parentDeptID)
    var deptIDs []int64
    for _, dept := range depts {
        deptIDs = append(deptIDs, dept.ID)
    }

    // 写入缓存（1小时）
    data, _ := json.Marshal(deptIDs)
    s.cache.Set(ctx, cacheKey, data, 1*time.Hour)

    return deptIDs
}
```

---

## 4. 字段权限控制

### 4.1 字段权限服务

```go
package rbac

import (
    "context"
    "regexp"
)

// FieldPermissionService 字段权限服务
type FieldPermissionService struct {
    fieldPermRepo repository.FieldPermissionRepository
    cache         *redis.Client
}

// GetVisibleFields 获取可见字段列表
func (s *FieldPermissionService) GetVisibleFields(
    ctx context.Context,
    userID int64,
    resourceType string,
) (map[string]string, error) {
    // 1. 获取用户的角色
    roles, _ := s.getUserRoles(ctx, userID)

    // 2. 获取所有字段权限配置
    allPerms := make(map[string]*FieldPermConfig)

    for _, role := range roles {
        perms, _ := s.fieldPermRepo.GetByRoleAndResource(ctx, role.ID, resourceType)

        for _, perm := range perms {
            if existing, ok := allPerms[perm.FieldName]; !ok || perm.Permission < existing.Permission {
                allPerms[perm.FieldName] = &FieldPermConfig{
                    FieldName:    perm.FieldName,
                    DisplayName:  perm.FieldDisplayName,
                    Permission:   perm.Permission,
                    MaskRule:     perm.MaskRule,
                    MaskPattern:  perm.MaskPattern,
                }
            }
        }
    }

    // 3. 只返回可见和脱敏字段
    visibleFields := make(map[string]string)
    for fieldName, config := range allPerms {
        if config.Permission != "hidden" {
            visibleFields[fieldName] = config.DisplayName
        }
    }

    return visibleFields, nil
}

// ApplyFieldPermissions 应用字段权限（过滤敏感字段）
func (s *FieldPermissionService) ApplyFieldPermissions(
    ctx context.Context,
    userID int64,
    resourceType string,
    data map[string]interface{},
) (map[string]interface{}, error) {
    // 1. 获取字段权限配置
    roles, _ := s.getUserRoles(ctx, userID)
    fieldPerms := make(map[string]*FieldPermConfig)

    for _, role := range roles {
        perms, _ := s.fieldPermRepo.GetByRoleAndResource(ctx, role.ID, resourceType)

        for _, perm := range perms {
            if existing, ok := fieldPerms[perm.FieldName]; !ok || perm.Permission < existing.Permission {
                fieldPerms[perm.FieldName] = &FieldPermConfig{
                    FieldName:    perm.FieldName,
                    Permission:   perm.Permission,
                    MaskRule:     perm.MaskRule,
                    MaskPattern:  perm.MaskPattern,
                    Replacement:  perm.MaskReplacement,
                }
            }
        }
    }

    // 2. 应用权限
    result := make(map[string]interface{})

    for fieldName, value := range data {
        perm, ok := fieldPerms[fieldName]

        if !ok {
            // 没有配置权限，默认可见
            result[fieldName] = value
            continue
        }

        switch perm.Permission {
        case "visible":
            // 可见
            result[fieldName] = value

        case "masked":
            // 脱敏
            result[fieldName] = s.maskValue(value, perm)

        case "hidden":
            // 隐藏，不返回
            continue
        }
    }

    return result, nil
}

// maskValue 脱敏处理
func (s *FieldPermissionService) maskValue(
    value interface{},
    perm *FieldPermConfig,
) interface{} {
    if value == nil {
        return nil
    }

    strValue := fmt.Sprintf("%v", value)

    // 优先使用正则表达式脱敏
    if perm.MaskPattern != "" {
        pattern := regexp.MustCompile(perm.MaskPattern)
        return pattern.ReplaceAllString(strValue, perm.Replacement)
    }

    // 使用预定义规则脱敏
    if perm.MaskRule != "" {
        switch perm.MaskRule {
        case "phone":
            // 手机号：138****1234
            if len(strValue) == 11 {
                return strValue[:3] + "****" + strValue[7:]
            }
        case "id_card":
            // 身份证：3301**********1234
            if len(strValue) == 18 {
                return strValue[:4] + "**********" + strValue[14:]
            }
        case "email":
            // 邮箱：u***@example.com
            parts := strings.Split(strValue, "@")
            if len(parts) == 2 {
                return string(parts[0][0]) + "***@" + parts[1]
            }
        case "bank_card":
            // 银行卡：****
            return "****"
        }
    }

    // 默认脱敏
    return perm.Replacement
}
```

---

## 5. 临时授权机制

### 5.1 临时授权服务

```go
package rbac

import (
    "context"
    "time"
    "github.com/google/uuid"
)

// TemporaryGrantService 临时授权服务
type TemporaryGrantService struct {
    grantRepo  repository.TemporaryGrantRepository
    auditRepo  repository.AuditLogRepository
    enforcer   *casbin.Enforcer
    scheduler  *cron.Cron
}

// CreateGrant 创建临时授权
func (s *TemporaryGrantService) CreateGrant(
    ctx context.Context,
    req *CreateGrantRequest,
) (*CreateGrantResponse, error) {
    // 1. 生成授权码
    grantCode := "TG-" + uuid.New().String()

    // 2. 构建授权记录
    grant := &entity.TemporaryGrant{
        TenantID:        req.TenantID,
        GrantCode:       grantCode,
        GrantorID:       req.GrantorID,
        GranteeID:       req.GranteeID,
        PermissionType:  req.PermissionType,
        PermissionConfig: req.PermissionConfig,
        Reason:          req.Reason,
        StartTime:       req.StartTime,
        EndTime:         req.EndTime,
    }

    // 3. 持久化授权
    if err := s.grantRepo.Create(ctx, grant); err != nil {
        return nil, err
    }

    // 4. 记录审计日志
    s.auditRepo.Create(ctx, &entity.AuditLog{
        TenantID:       req.TenantID,
        OperatorID:     req.GrantorID,
        OperationType:  "grant",
        TargetType:     "user",
        TargetID:       fmt.Sprintf("%d", req.GranteeID),
        PermissionDetail: req.PermissionConfig,
        OperationResult: "success",
    })

    // 5. 添加到自动回收任务
    s.scheduleRevoke(grant)

    return &CreateGrantResponse{
        GrantCode: grantCode,
        Message:   "临时授权创建成功",
    }, nil
}

// CheckGrant 检查临时授权
func (s *TemporaryGrantService) CheckGrant(
    ctx context.Context,
    granteeID int64,
    resourceType string,
    resourceID string,
    action string,
) (bool, error) {
    now := time.Now()

    // 1. 获取有效的临时授权
    grants, _ := s.grantRepo.GetValidGrants(ctx, granteeID, now)

    // 2. 检查每个授权
    for _, grant := range grants {
        hasPermission := s.checkGrantPermission(grant, resourceType, resourceID, action)
        if hasPermission {
            return true, nil
        }
    }

    return false, nil
}

// RevokeGrant 撤销临时授权
func (s *TemporaryGrantService) RevokeGrant(
    ctx context.Context,
    grantCode string,
    revokedBy int64,
    reason string,
) error {
    // 1. 获取授权记录
    grant, err := s.grantRepo.GetByCode(ctx, grantCode)
    if err != nil {
        return err
    }

    // 2. 撤销授权
    grant.IsRevoked = true
    grant.RevokedAt = time.Now()
    grant.RevokedBy = revokedBy
    grant.RevokeReason = reason

    if err := s.grantRepo.Update(ctx, grant); err != nil {
        return err
    }

    // 3. 记录审计日志
    s.auditRepo.Create(ctx, &entity.AuditLog{
        TenantID:       grant.TenantID,
        OperatorID:     revokedBy,
        OperationType:  "revoke",
        TargetType:     "user",
        TargetID:       fmt.Sprintf("%d", grant.GranteeID),
        PermissionDetail: grant.PermissionConfig,
        OperationResult: "success",
    })

    return nil
}

// scheduleRevoke 调度自动回收任务
func (s *TemporaryGrantService) scheduleRevoke(
    grant *entity.TemporaryGrant,
) {
    // 计算回收时间（提前1分钟）
    revokeTime := grant.EndTime.Add(-1 * time.Minute)

    // 添加定时任务
    s.scheduler.AddFunc(revokeTime.Format("2006-01-02 15:04:05"), func() {
        s.RevokeGrant(context.Background(), grant.GrantCode, 0, "授权到期自动回收")
    })
}
```

---

## 6. 权限验证引擎

### 6.1 Casbin集成

```go
package rbac

import (
    "context"
    "github.com/casbin/casbin/v2"
    "github.com/casbin/casbin/v2/model"
    xormadapter "github.com/casbin/xorm-adapter/v2"
)

// PermissionEngine 权限引擎
type PermissionEngine struct {
    enforcer *casbin.Enforcer
}

// NewPermissionEngine 创建权限引擎
func NewPermissionEngine(dsn string) (*PermissionEngine, error) {
    // 1. 初始化Xorm适配器
    adapter, _ := xormadapter.NewAdapter("mysql", dsn)

    // 2. 加载RBAC模型
    modelText := `
    [request_definition]
    r = sub, obj, act

    [policy_definition]
    p = sub, obj, act

    [role_definition]
    g = _, _

    [policy_effect]
    e = some(where (p.eft == allow))

    [matchers]
    m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
    `

    m, _ := model.NewModelFromString(modelText)

    // 3. 创建Enforcer
    enforcer, _ := casbin.NewEnforcer(m, adapter)

    return &PermissionEngine{
        enforcer: enforcer,
    }, nil
}

// CheckPermission 检查权限
func (e *PermissionEngine) CheckPermission(
    ctx context.Context,
    userID int64,
    resource string,
    action string,
) (bool, error) {
    // 1. 检查RBAC权限
    allowed, _ := e.enforcer.Enforce(
        fmt.Sprintf("user_%d", userID),
        resource,
        action,
    )

    if !allowed {
        // 2. 检查临时授权
        return e.checkTemporaryGrant(ctx, userID, resource, action)
    }

    return true, nil
}

// AddRoleForUser 为用户添加角色
func (e *PermissionEngine) AddRoleForUser(
    ctx context.Context,
    userID int64,
    role string,
) error {
    _, err := e.enforcer.AddRoleForUser(
        fmt.Sprintf("user_%d", userID),
        role,
    )
    return err
}

// AddPermissionForRole 为角色添加权限
func (e *PermissionEngine) AddPermissionForRole(
    ctx context.Context,
    role, resource, action string,
) error {
    _, err := e.enforcer.AddPolicy(
        role,
        resource,
        action,
    )
    return err
}
```

---

## 7. API设计

### 7.1 数据权限API

| 方法 | 路径 | 功能 | 请求示例 | 响应示例 |
|------|------|------|----------|----------|
| GET | /api/v1/rbac/data-permissions | 查询数据权限 | `?role_id=1&resource_type=user` | `[{permission_scope:"department",...}]` |
| POST | /api/v1/rbac/data-permissions | 创建数据权限 | `{"role_id":1,"resource_type":"user","permission_scope":"department"}` | `{"id":123}` |
| PUT | /api/v1/rbac/data-permissions/:id | 更新数据权限 | `{"permission_scope":"all"}` | `{"success":true}` |
| DELETE | /api/v1/rbac/data-permissions/:id | 删除数据权限 | - | `{"success":true}` |

### 7.2 字段权限API

| 方法 | 路径 | 功能 | 请求示例 | 响应示例 |
|------|------|------|----------|----------|
| GET | /api/v1/rbac/field-permissions | 查询字段权限 | `?role_id=1&resource_type=user` | `[{field_name:"phone",permission:"masked"}]` |
| POST | /api/v1/rbac/field-permissions | 创建字段权限 | `{"role_id":1,"field_name":"phone","permission":"masked"}` | `{"id":456}` |
| PUT | /api/v1/rbac/field-permissions/:id | 更新字段权限 | `{"permission":"visible"}` | `{"success":true}` |

### 7.3 临时授权API

| 方法 | 路径 | 功能 | 请求示例 | 响应示例 |
|------|------|------|----------|----------|
| POST | /api/v1/rbac/temporary-grants | 创建临时授权 | `{"grantee_id":100,"permission_type":"role","start_time":"...","end_time":"..."}` | `{"grant_code":"TG-xxx"}` |
| GET | /api/v1/rbac/temporary-grants | 查询临时授权 | `?grantee_id=100` | `[{grant_code:"TG-xxx",...}]` |
| POST | /api/v1/rbac/temporary-grants/:code/revoke | 撤销临时授权 | `{"reason":"提前结束"}` | `{"success":true}` |

---

## 8. 前端组件

### 8.1 数据权限配置组件

```typescript
// components/rbac/DataPermissionConfig.tsx
import React, { useState, useEffect } from 'react';
import { Form, Select, TreeSelect, Button, Card } from '@douyinfe/semi-ui';

export const DataPermissionConfig: React.FC<Props> = ({ roleId, resourceType }) => {
  const [form] = Form.useForm();
  const [orgTree, setOrgTree] = useState([]);

  useEffect(() => {
    fetchOrgTree();
  }, []);

  const fetchOrgTree = async () => {
    const response = await fetch('/api/v1/organizations/tree');
    const data = await response.json();
    setOrgTree(data);
  };

  const handleSubmit = async (values) => {
    await fetch('/api/v1/rbac/data-permissions', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        role_id: roleId,
        resource_type: resourceType,
        ...values,
      }),
    });
    Toast.success('保存成功');
  };

  return (
    <Form
      form={form}
      onSubmit={handleSubmit}
      initValues={{
        permission_scope: 'department',
      }}
    >
      <Form.Select
        field="permission_scope"
        label="权限范围"
        style={{ width: 300 }}
      >
        <Select.Option value="all">全部数据</Select.Option>
        <Select.Option value="department_and_below">本部门及下级</Select.Option>
        <Select.Option value="department">本部门</Select.Option>
        <Select.Option value="self">仅本人</Select.Option>
        <Select.Option value="custom">自定义</Select.Option>
      </Form.Select>

      <Form.Dependency>
        {({ permission_scope }) =>
          permission_scope === 'custom' && (
            <Form.TreeSelect
              field="custom_org_ids"
              label="选择组织"
              treeData={orgTree}
              multiple
              placeholder="请选择允许访问的组织"
              style={{ width: 400 }}
            />
          )
        }
      </Form.Dependency>

      <Button type="primary" htmlType="submit">
        保存配置
      </Button>
    </Form>
  );
};
```

---

## 9. 实施计划

### 9.1 开发阶段划分

| 阶段 | 任务 | 工期 | 交付物 |
|------|------|------|--------|
| **第1周** | 数据库设计与创建 | 2天 | 4张表 |
| | 数据权限服务开发 | 3天 | 权限过滤+检查 |
| **第2周** | 字段权限服务 | 2天 | 脱敏引擎 |
| | 临时授权服务 | 2天 | 授权+回收 |
| | Casbin集成 | 1天 | 权限引擎 |
| **第3周** | 前端组件开发 | 3天 | 3个组件 |
| | API集成测试 | 2天 | 测试报告 |

### 9.2 技术风险与缓解

| 风险 | 影响 | 概率 | 缓解措施 |
|------|------|------|----------|
| 性能损耗（权限检查） | 高 | 中 | Redis缓存+权限预计算 |
| 组织架构变化 | 中 | 中 | 事件驱动更新缓存 |
| Casbin配置复杂 | 中 | 低 | 提供配置模板+可视化配置 |

### 9.3 成功指标

- ✅ 权限检查延迟: **< 50ms** (99%请求)
- ✅ 脱敏准确率: **100%**
- ✅ 临时授权准时回收: **100%**
- ✅ 权限审计完整性: **100%**

---

**文档结束**
