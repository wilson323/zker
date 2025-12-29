# 研发 A - 后端架构师 8 周开发计划 v1.0

> **角色定位**: 后端架构师 - 负责多租户架构、RBAC 权限系统、智能路由引擎等核心架构功能
>
> **工作目标**: 构建 ZKER 企业级核心架构，确保多租户隔离、细粒度权限控制、智能路由能力
>
> **核心原则**: 架构优先、接口清晰、高内聚低耦合、可扩展性

---

## 📋 个人职责概述

### 核心负责模块

```
backend/
├── domain/tenant/              ✅ 专属负责 - 租户领域
│   ├── entity/
│   ├── repository/
│   └── service/
├── domain/permission/          ✅ 专属负责 - 权限领域
│   ├── entity/
│   ├── repository/
│   └── service/
├── domain/routing/             ✅ 专属负责 - 智能路由领域
│   ├── entity/
│   ├── repository/
│   └── service/
├── application/tenant/         ✅ 专属负责
├── application/permission/     ✅ 专属负责
├── application/routing/        ✅ 专属负责
├── api/v1/tenant/             ✅ 专属负责
├── api/v1/permission/         ✅ 专属负责
└── api/v1/routing/            ✅ 专属负责
```

### 协作接口

| 协作对象 | 协作内容 | 接口定义位置 | 依赖关系 |
|---------|---------|------------|---------|
| **研发 B** | 错误码定义、性能测试、监控埋点 | `types/errno/errno.go` | 研发 A 定义 → 研发 B 使用 |
| **研发 C** | 前端 API 契约、数据格式 | `idl/tenant.thrift`、`idl/permission.thrift` | 研发 A 定义 IDL → 研发 C 生成代码 |
| **研发 D** | 数据库迁移脚本、Docker 配置 | `migrations/`、`docker/docker-compose.yml` | 研发 A 设计 Schema → 研发 D 编写迁移 |

---

## 🎯 8 周详细开发计划

### Week 1-2: 多租户核心架构 + RBAC 基础

#### Week 1: 租户领域核心实现

**目标**: 完成租户领域的实体定义、仓储接口、核心服务

##### Day 1-2: 实体和仓储设计

**数据库表设计**:
```sql
-- 1. 租户表
CREATE TABLE tenants (
    tenant_id VARCHAR(36) PRIMARY KEY COMMENT '租户ID',
    tenant_name VARCHAR(200) NOT NULL COMMENT '租户名称',
    tenant_type ENUM('individual', 'team', 'enterprise') NOT NULL COMMENT '租户类型',
    status ENUM('active', 'suspended', 'deleted') DEFAULT 'active' COMMENT '状态',
    subscription_tier ENUM('free', 'pro', 'enterprise') DEFAULT 'free' COMMENT '订阅等级',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE KEY uk_tenant_name (tenant_name, deleted_at),
    INDEX idx_status_type (status, tenant_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户表';

-- 2. 订阅表
CREATE TABLE subscriptions (
    subscription_id VARCHAR(36) PRIMARY KEY COMMENT '订阅ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    plan_tier ENUM('free', 'pro', 'enterprise') NOT NULL COMMENT '套餐等级',
    billing_cycle ENUM('monthly', 'yearly') NOT NULL COMMENT '计费周期',
    start_date DATE NOT NULL COMMENT '开始日期',
    end_date DATE COMMENT '结束日期',
    auto_renew BOOLEAN DEFAULT TRUE COMMENT '自动续费',
    status ENUM('active', 'expired', 'cancelled') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status_end_date (status, end_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订阅表';

-- 3. 配额表
CREATE TABLE quotas (
    quota_id VARCHAR(36) PRIMARY KEY COMMENT '配额ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    resource_type ENUM('bots', 'messages', 'storage', 'team_members') NOT NULL COMMENT '资源类型',
    max_limit INT NOT NULL COMMENT '最大限制',
    used_count INT DEFAULT 0 COMMENT '已使用数量',
    reset_cycle ENUM('daily', 'monthly', 'yearly', 'never') DEFAULT 'monthly' COMMENT '重置周期',
    last_reset_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '上次重置时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    UNIQUE KEY uk_tenant_resource (tenant_id, resource_type),
    INDEX idx_tenant_id (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='配额表';
```

**实体层代码** (`domain/tenant/entity/`):
```go
// entity/tenant.go
package tenant

import "time"

// Tenant 租户实体
type Tenant struct {
    TenantID         string    `json:"tenant_id" gorm:"primaryKey;type:varchar(36)"`
    TenantName       string    `json:"tenant_name" gorm:"type:varchar(200);not null"`
    TenantType       string    `json:"tenant_type" gorm:"type:enum('individual','team','enterprise');not null"`
    Status           string    `json:"status" gorm:"type:enum('active','suspended','deleted');default:'active'"`
    SubscriptionTier string    `json:"subscription_tier" gorm:"type:enum('free','pro','enterprise');default:'free'"`
    CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt        *time.Time `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (Tenant) TableName() string {
    return "tenants"
}

// Subscription 订阅实体
type Subscription struct {
    SubscriptionID string    `json:"subscription_id" gorm:"primaryKey;type:varchar(36)"`
    TenantID       string    `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
    PlanTier       string    `json:"plan_tier" gorm:"type:enum('free','pro','enterprise');not null"`
    BillingCycle   string    `json:"billing_cycle" gorm:"type:enum('monthly','yearly');not null"`
    StartDate      time.Time `json:"start_date" gorm:"type:date;not null"`
    EndDate        *time.Time `json:"end_date" gorm:"type:date"`
    AutoRenew      bool      `json:"auto_renew" gorm:"default:true"`
    Status         string    `json:"status" gorm:"type:enum('active','expired','cancelled');default:'active'"`
    CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`

    // 关联
    Tenant Tenant `json:"tenant" gorm:"foreignKey:TenantID"`
}

// Quota 配额实体
type Quota struct {
    QuotaID      string    `json:"quota_id" gorm:"primaryKey;type:varchar(36)"`
    TenantID     string    `json:"tenant_id" gorm:"type:varchar(36);not null;uniqueIndex:uk_tenant_resource"`
    ResourceType string    `json:"resource_type" gorm:"type:enum('bots','messages','storage','team_members');not null;uniqueIndex:uk_tenant_resource"`
    MaxLimit     int       `json:"max_limit" gorm:"not null"`
    UsedCount    int       `json:"used_count" gorm:"default:0"`
    ResetCycle   string    `json:"reset_cycle" gorm:"type:enum('daily','monthly','yearly','never');default:'monthly'"`
    LastResetAt  time.Time `json:"last_reset_at" gorm:"autoCreateTime"`
    CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
```

**仓储接口** (`domain/tenant/repository/`):
```go
// repository/tenant_repository.go
package repository

import (
    "context"
    "github.com/coze-studio/domain/tenant/entity"
)

// TenantRepository 租户仓储接口
type TenantRepository interface {
    // Create 创建租户
    Create(ctx context.Context, tenant *entity.Tenant) error

    // GetByID 根据ID获取租户
    GetByID(ctx context.Context, tenantID string) (*entity.Tenant, error)

    // GetByName 根据名称获取租户
    GetByName(ctx context.Context, name string) (*entity.Tenant, error)

    // Update 更新租户
    Update(ctx context.Context, tenant *entity.Tenant) error

    // Delete 软删除租户
    Delete(ctx context.Context, tenantID string) error

    // List 分页查询租户列表
    List(ctx context.Context, filter *TenantFilter) ([]*entity.Tenant, int64, error)

    // UpdateStatus 更新租户状态
    UpdateStatus(ctx context.Context, tenantID, status string) error
}

// TenantFilter 租户查询过滤器
type TenantFilter struct {
    Status           string
    TenantType       string
    SubscriptionTier string
    PageToken        string
    PageSize         int
}

// SubscriptionRepository 订阅仓储接口
type SubscriptionRepository interface {
    Create(ctx context.Context, sub *entity.Subscription) error
    GetByTenantID(ctx context.Context, tenantID string) (*entity.Subscription, error)
    Update(ctx context.Context, sub *entity.Subscription) error
    List(ctx context.Context, filter *SubscriptionFilter) ([]*entity.Subscription, int64, error)
}

// SubscriptionFilter 订阅查询过滤器
type SubscriptionFilter struct {
    TenantID string
    Status   string
    PlanTier string
    PageToken string
    PageSize int
}

// QuotaRepository 配额仓储接口
type QuotaRepository interface {
    Create(ctx context.Context, quota *entity.Quota) error
    GetByTenantAndResource(ctx context.Context, tenantID, resourceType string) (*entity.Quota, error)
    UpdateUsedCount(ctx context.Context, quotaID string, delta int) error
    ResetUsage(ctx context.Context, quotaID string) error
    List(ctx context.Context, tenantID string) ([]*entity.Quota, error)
}
```

##### Day 3-5: 领域服务实现

**配额检查服务** (`domain/tenant/service/quota_service.go`):
```go
// service/quota_service.go
package service

import (
    "context"
    "fmt"
    "github.com/coze-studio/domain/tenant/entity"
    "github.com/coze-studio/domain/tenant/repository"
    "github.com/coze-studio/types/errno"
)

type QuotaService struct {
    quotaRepo repository.QuotaRepository
}

func NewQuotaService(quotaRepo repository.QuotaRepository) *QuotaService {
    return &QuotaService{quotaRepo: quotaRepo}
}

// CheckQuota 检查配额是否充足
func (s *QuotaService) CheckQuota(ctx context.Context, tenantID, resourceType string, requiredCount int) error {
    quota, err := s.quotaRepo.GetByTenantAndResource(ctx, tenantID, resourceType)
    if err != nil {
        return fmt.Errorf("failed to get quota: %w", err)
    }

    if quota.UsedCount+requiredCount > quota.MaxLimit {
        return &errno.Errno{
            Code:    "QUOTA_EXCEEDED",
            Message: fmt.Sprintf("Quota exceeded for %s", resourceType),
            Details: map[string]interface{}{
                "tenant_id":      tenantID,
                "resource_type":  resourceType,
                "used":          quota.UsedCount,
                "max_limit":     quota.MaxLimit,
                "required":      requiredCount,
            },
        }
    }

    return nil
}

// ConsumeQuota 消费配额
func (s *QuotaService) ConsumeQuota(ctx context.Context, tenantID, resourceType string, count int) error {
    if err := s.CheckQuota(ctx, tenantID, resourceType, count); err != nil {
        return err
    }

    quota, err := s.quotaRepo.GetByTenantAndResource(ctx, tenantID, resourceType)
    if err != nil {
        return err
    }

    return s.quotaRepo.UpdateUsedCount(ctx, quota.QuotaID, count)
}

// RollbackQuota 回滚配额
func (s *QuotaService) RollbackQuota(ctx context.Context, tenantID, resourceType string, count int) error {
    quota, err := s.quotaRepo.GetByTenantAndResource(ctx, tenantID, resourceType)
    if err != nil {
        return err
    }

    return s.quotaRepo.UpdateUsedCount(ctx, quota.QuotaID, -count)
}
```

**订阅管理服务** (`domain/tenant/service/subscription_service.go`):
```go
// service/subscription_service.go
package service

import (
    "context"
    "time"
    "github.com/coze-studio/domain/tenant/entity"
    "github.com/coze-studio/domain/tenant/repository"
)

type SubscriptionService struct {
    subRepo    repository.SubscriptionRepository
    quotaRepo  repository.QuotaRepository
}

func NewSubscriptionService(subRepo repository.SubscriptionRepository, quotaRepo repository.QuotaRepository) *SubscriptionService {
    return &SubscriptionService{
        subRepo:   subRepo,
        quotaRepo: quotaRepo,
    }
}

// CreateSubscription 创建订阅
func (s *SubscriptionService) CreateSubscription(ctx context.Context, tenantID, planTier, billingCycle string) (*entity.Subscription, error) {
    sub := &entity.Subscription{
        SubscriptionID: generateID(),
        TenantID:       tenantID,
        PlanTier:       planTier,
        BillingCycle:   billingCycle,
        StartDate:      time.Now(),
        Status:         "active",
        AutoRenew:      true,
    }

    // 计算结束日期
    if billingCycle == "monthly" {
        endDate := time.Now().AddDate(0, 1, 0)
        sub.EndDate = &endDate
    } else if billingCycle == "yearly" {
        endDate := time.Now().AddDate(1, 0, 0)
        sub.EndDate = &endDate
    }

    // 创建配额
    quotas := s.getQuotasForTier(planTier)
    for _, quota := range quotas {
        quota.TenantID = tenantID
        if err := s.quotaRepo.Create(ctx, quota); err != nil {
            return nil, err
        }
    }

    if err := s.subRepo.Create(ctx, sub); err != nil {
        return nil, err
    }

    return sub, nil
}

// getQuotasForTier 根据订阅等级获取配额
func (s *SubscriptionService) getQuotasForTier(planTier string) []*entity.Quota {
    quotas := make([]*entity.Quota, 0)

    switch planTier {
    case "free":
        quotas = append(quotas, &entity.Quota{
            QuotaID:      generateID(),
            ResourceType: "bots",
            MaxLimit:     10,
            ResetCycle:   "never",
        })
        quotas = append(quotas, &entity.Quota{
            QuotaID:      generateID(),
            ResourceType: "messages",
            MaxLimit:     1000,
            ResetCycle:   "monthly",
        })
    case "pro":
        quotas = append(quotas, &entity.Quota{
            QuotaID:      generateID(),
            ResourceType: "bots",
            MaxLimit:     100,
            ResetCycle:   "never",
        })
        quotas = append(quotas, &entity.Quota{
            QuotaID:      generateID(),
            ResourceType: "messages",
            MaxLimit:     100000,
            ResetCycle:   "monthly",
        })
    case "enterprise":
        quotas = append(quotas, &entity.Quota{
            QuotaID:      generateID(),
            ResourceType: "bots",
            MaxLimit:     -1, // 无限制
            ResetCycle:   "never",
        })
        quotas = append(quotas, &entity.Quota{
            QuotaID:      generateID(),
            ResourceType: "messages",
            MaxLimit:     -1, // 无限制
            ResetCycle:   "monthly",
        })
    }

    return quotas
}
```

#### Week 2: RBAC 权限系统基础

**目标**: 完成角色、权限、数据权限的基础实现

##### Day 1-2: 权限实体和仓储

**数据库表设计**:
```sql
-- 4. 角色表
CREATE TABLE roles (
    role_id VARCHAR(36) PRIMARY KEY COMMENT '角色ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    role_name VARCHAR(100) NOT NULL COMMENT '角色名称',
    role_code VARCHAR(50) NOT NULL COMMENT '角色编码',
    role_type ENUM('system', 'custom') NOT NULL COMMENT '角色类型',
    parent_role_id VARCHAR(36) COMMENT '父角色ID',
    description TEXT COMMENT '角色描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    FOREIGN KEY (parent_role_id) REFERENCES roles(role_id) ON DELETE SET NULL,
    UNIQUE KEY uk_tenant_code (tenant_id, role_code, deleted_at),
    INDEX idx_tenant_id (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 5. 数据权限表
CREATE TABLE data_permissions (
    permission_id VARCHAR(36) PRIMARY KEY COMMENT '权限ID',
    role_id VARCHAR(36) NOT NULL COMMENT '角色ID',
    resource_type ENUM('bots', 'conversations', 'knowledge', 'workflows', 'plugins') NOT NULL COMMENT '资源类型',
    scope ENUM('ALL', 'DEPARTMENT', 'OWN', 'CUSTOM', 'NONE') NOT NULL COMMENT '权限范围',
    custom_filter JSON COMMENT '自定义过滤条件',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    INDEX idx_role_id (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据权限表';

-- 6. 字段权限表
CREATE TABLE field_permissions (
    permission_id VARCHAR(36) PRIMARY KEY COMMENT '权限ID',
    role_id VARCHAR(36) NOT NULL COMMENT '角色ID',
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型',
    field_name VARCHAR(100) NOT NULL COMMENT '字段名称',
    permission_level ENUM('hidden', 'readonly', 'editable') NOT NULL COMMENT '权限级别',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    UNIQUE KEY uk_role_resource_field (role_id, resource_type, field_name),
    INDEX idx_role_id (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='字段权限表';

-- 7. 用户角色关联表
CREATE TABLE user_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(36) NOT NULL COMMENT '用户ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    role_id VARCHAR(36) NOT NULL COMMENT '角色ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_tenant_role (user_id, tenant_id, role_id),
    INDEX idx_user_id (user_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_role_id (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';
```

**实体层代码** (`domain/permission/entity/`):
```go
// entity/role.go
package permission

import "time"

type Role struct {
    RoleID        string     `json:"role_id" gorm:"primaryKey;type:varchar(36)"`
    TenantID      string     `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
    RoleName      string     `json:"role_name" gorm:"type:varchar(100);not null"`
    RoleCode      string     `json:"role_code" gorm:"type:varchar(50);not null"`
    RoleType      string     `json:"role_type" gorm:"type:enum('system','custom');not null"`
    ParentRoleID  *string    `json:"parent_role_id" gorm:"type:varchar(36)"`
    Description   string     `json:"description" gorm:"type:text"`
    CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt     *time.Time `json:"deleted_at" gorm:"index"`

    // 关联
    ParentRole   *Role              `json:"parent_role,omitempty" gorm:"foreignKey:ParentRoleID"`
    DataPerms    []DataPermission   `json:"data_permissions,omitempty" gorm:"foreignKey:RoleID"`
    FieldPerms   []FieldPermission  `json:"field_permissions,omitempty" gorm:"foreignKey:RoleID"`
}

type DataPermission struct {
    PermissionID  string          `json:"permission_id" gorm:"primaryKey;type:varchar(36)"`
    RoleID        string          `json:"role_id" gorm:"type:varchar(36);not null;index:idx_role_id"`
    ResourceType  string          `json:"resource_type" gorm:"type:enum('bots','conversations','knowledge','workflows','plugins');not null"`
    Scope         string          `json:"scope" gorm:"type:enum('ALL','DEPARTMENT','OWN','CUSTOM','NONE');not null"`
    CustomFilter  string          `json:"custom_filter" gorm:"type:json"`
    CreatedAt     time.Time       `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt     time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

type FieldPermission struct {
    PermissionID    string    `json:"permission_id" gorm:"primaryKey;type:varchar(36)"`
    RoleID          string    `json:"role_id" gorm:"type:varchar(36);not null;index:idx_role_id"`
    ResourceType    string    `json:"resource_type" gorm:"type:varchar(50);not null"`
    FieldName       string    `json:"field_name" gorm:"type:varchar(100);not null"`
    PermissionLevel string    `json:"permission_level" gorm:"type:enum('hidden','readonly','editable');not null"`
    CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type UserRole struct {
    ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
    UserID    string    `json:"user_id" gorm:"type:varchar(36);not null;index:idx_user_id"`
    TenantID  string    `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
    RoleID    string    `json:"role_id" gorm:"type:varchar(36);not null;index:idx_role_id"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

    // 关联
    Role Role `json:"role" gorm:"foreignKey:RoleID"`
}
```

##### Day 3-5: 权限服务实现

**权限检查服务** (`domain/permission/service/permission_checker.go`):
```go
// service/permission_checker.go
package service

import (
    "context"
    "fmt"
    "github.com/coze-studio/domain/permission/entity"
    "github.com/coze-studio/domain/permission/repository"
    "github.com/coze-studio/types/errno"
)

type PermissionChecker struct {
    roleRepo       repository.RoleRepository
    dataPermRepo   repository.DataPermissionRepository
    fieldPermRepo  repository.FieldPermissionRepository
    userRoleRepo   repository.UserRoleRepository
}

func NewPermissionChecker(
    roleRepo repository.RoleRepository,
    dataPermRepo repository.DataPermissionRepository,
    fieldPermRepo repository.FieldPermissionRepository,
    userRoleRepo repository.UserRoleRepository,
) *PermissionChecker {
    return &PermissionChecker{
        roleRepo:      roleRepo,
        dataPermRepo:  dataPermRepo,
        fieldPermRepo: fieldPermRepo,
        userRoleRepo:  userRoleRepo,
    }
}

// CheckDataPermission 检查数据权限
func (p *PermissionChecker) CheckDataPermission(ctx context.Context, userID, tenantID, resourceType, resourceID string) error {
    // 1. 获取用户的所有角色
    roles, err := p.userRoleRepo.GetRolesByUser(ctx, userID, tenantID)
    if err != nil {
        return err
    }

    // 2. 检查数据权限
    for _, role := range roles {
        perm, err := p.dataPermRepo.GetByRoleAndResource(ctx, role.RoleID, resourceType)
        if err != nil {
            continue
        }

        switch perm.Scope {
        case "ALL":
            return nil // 有全部权限
        case "NONE":
            return &errno.Errno{
                Code:    "PERMISSION_DENIED",
                Message: "No permission for this resource",
            }
        case "OWN":
            // 检查是否是资源的创建者
            if !p.isOwner(ctx, userID, resourceType, resourceID) {
                return &errno.Errno{
                    Code:    "PERMISSION_DENIED",
                    Message: "Can only access own resources",
                }
            }
        case "DEPARTMENT":
            // 检查是否同部门
            if !p.isSameDepartment(ctx, userID, resourceType, resourceID) {
                return &errno.Errno{
                    Code:    "PERMISSION_DENIED",
                    Message: "Can only access department resources",
                }
            }
        case "CUSTOM":
            // 根据 custom_filter 进行过滤
            if !p.matchCustomFilter(ctx, userID, perm.CustomFilter, resourceType, resourceID) {
                return &errno.Errno{
                    Code:    "PERMISSION_DENIED",
                    Message: "Does not match custom filter",
                }
            }
        }
    }

    return nil
}

// GetFieldPermissions 获取字段权限
func (p *PermissionChecker) GetFieldPermissions(ctx context.Context, userID, tenantID, resourceType string) (map[string]string, error) {
    roles, err := p.userRoleRepo.GetRolesByUser(ctx, userID, tenantID)
    if err != nil {
        return nil, err
    }

    fieldPerms := make(map[string]string)
    for _, role := range roles {
        perms, _ := p.fieldPermRepo.GetByRoleAndResource(ctx, role.RoleID, resourceType)
        for _, perm := range perms {
            // 如果已有权限，优先级更高：editable > readonly > hidden
            if existing, ok := fieldPerms[perm.FieldName]; ok {
                if existing == "editable" || perm.PermissionLevel == "hidden" {
                    continue
                }
            }
            fieldPerms[perm.FieldName] = perm.PermissionLevel
        }
    }

    return fieldPerms, nil
}

// isOwner 检查是否是资源创建者
func (p *PermissionChecker) isOwner(ctx context.Context, userID, resourceType, resourceID string) bool {
    // 根据资源类型查询资源表，检查 creator_id
    // 实现省略...
    return true
}

// isSameDepartment 检查是否同部门
func (p *PermissionChecker) isSameDepartment(ctx context.Context, userID, resourceType, resourceID string) bool {
    // 实现省略...
    return true
}

// matchCustomFilter 匹配自定义过滤条件
func (p *PermissionChecker) matchCustomFilter(ctx context.Context, userID string, customFilter, resourceType, resourceID string) bool {
    // 解析 customFilter JSON，进行过滤
    // 实现省略...
    return true
}
```

---

### Week 3-4: 智能路由引擎

#### Week 3: 意图识别引擎

**目标**: 实现混合意图匹配器（规则+相似度+模型）

##### Day 1-3: 规则匹配器

**数据库表设计**:
```sql
-- 8. 路由规则表
CREATE TABLE routing_rules (
    rule_id VARCHAR(36) PRIMARY KEY COMMENT '规则ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    rule_name VARCHAR(200) NOT NULL COMMENT '规则名称',
    rule_type ENUM('keyword', 'regex', 'intent', 'category') NOT NULL COMMENT '规则类型',
    priority INT NOT NULL DEFAULT 0 COMMENT '优先级（数字越大优先级越高）',
    condition JSON NOT NULL COMMENT '匹配条件',
    target_bot_id VARCHAR(36) COMMENT '目标Bot ID',
    target_workflow_id VARCHAR(36) COMMENT '目标工作流 ID',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    INDEX idx_tenant_priority (tenant_id, priority DESC),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='路由规则表';
```

**规则匹配器实现** (`domain/routing/service/rule_matcher.go`):
```go
// service/rule_matcher.go
package service

import (
    "context"
    "encoding/json"
    "regexp"
    "sort"
    "github.com/coze-studio/domain/routing/entity"
    "github.com/coze-studio/domain/routing/repository"
)

type RuleBasedMatcher struct {
    ruleRepo repository.RoutingRuleRepository
}

func NewRuleBasedMatcher(ruleRepo repository.RoutingRuleRepository) *RuleBasedMatcher {
    return &RuleBasedMatcher{ruleRepo: ruleRepo}
}

type MatchInput struct {
    UserInput  string                 `json:"user_input"`
    TenantID   string                 `json:"tenant_id"`
    Context    map[string]interface{} `json:"context"`
}

type MatchOutput struct {
    BotID       string  `json:"bot_id,omitempty"`
    WorkflowID  string  `json:"workflow_id,omitempty"`
    Confidence  float64 `json:"confidence"`
    MatchType   string  `json:"match_type"`
    RuleID      string  `json:"rule_id"`
}

// Match 规则匹配
func (m *RuleBasedMatcher) Match(ctx context.Context, input *MatchInput) ([]*MatchOutput, error) {
    // 1. 获取所有启用的规则，按优先级排序
    rules, err := m.ruleRepo.GetActiveRulesByTenant(ctx, input.TenantID)
    if err != nil {
        return nil, err
    }

    results := make([]*MatchOutput, 0)

    // 2. 遍历规则进行匹配
    for _, rule := range rules {
        matched, confidence := m.matchRule(ctx, rule, input)
        if matched {
            results = append(results, &MatchOutput{
                BotID:      rule.TargetBotID,
                WorkflowID: rule.TargetWorkflowID,
                Confidence: confidence,
                MatchType:  rule.RuleType,
                RuleID:     rule.RuleID,
            })
        }
    }

    // 3. 按优先级和置信度排序
    sort.Slice(results, func(i, j int) bool {
        if results[i].Confidence != results[j].Confidence {
            return results[i].Confidence > results[j].Confidence
        }
        return results[i].RuleID > results[j].RuleID
    })

    // 4. 返回 Top-3
    if len(results) > 3 {
        results = results[:3]
    }

    return results, nil
}

// matchRule 匹配单个规则
func (m *RuleBasedMatcher) matchRule(ctx context.Context, rule *entity.RoutingRule, input *MatchInput) (bool, float64) {
    var condition map[string]interface{}
    if err := json.Unmarshal([]byte(rule.Condition), &condition); err != nil {
        return false, 0
    }

    switch rule.RuleType {
    case "keyword":
        return m.matchKeyword(condition, input.UserInput)
    case "regex":
        return m.matchRegex(condition, input.UserInput)
    case "intent":
        return m.matchIntent(condition, input.Context)
    case "category":
        return m.matchCategory(condition, input.Context)
    }

    return false, 0
}

// matchKeyword 关键词匹配
func (m *RuleBasedMatcher) matchKeyword(condition map[string]interface{}, userInput string) (bool, float64) {
    keywords, ok := condition["keywords"].([]interface{})
    if !ok {
        return false, 0
    }

    matchCount := 0
    for _, kw := range keywords {
        keyword, ok := kw.(string)
        if !ok {
            continue
        }
        if contains(userInput, keyword) {
            matchCount++
        }
    }

    if matchCount == 0 {
        return false, 0
    }

    // 置信度 = 匹配的关键词数量 / 总关键词数量
    confidence := float64(matchCount) / float64(len(keywords))
    return true, confidence
}

// matchRegex 正则匹配
func (m *RuleBasedMatcher) matchRegex(condition map[string]interface{}, userInput string) (bool, float64) {
    pattern, ok := condition["pattern"].(string)
    if !ok {
        return false, 0
    }

    matched, err := regexp.MatchString(pattern, userInput)
    if err != nil || !matched {
        return false, 0
    }

    return true, 1.0
}

// matchIntent 意图匹配
func (m *RuleBasedMatcher) matchIntent(condition map[string]interface{}, context map[string]interface{}) (bool, float64) {
    requiredIntent, ok := condition["intent"].(string)
    if !ok {
        return false, 0
    }

    currentIntent, ok := context["intent"].(string)
    if !ok {
        return false, 0
    }

    if currentIntent == requiredIntent {
        return true, 1.0
    }

    return false, 0
}

// matchCategory 分类匹配
func (m *RuleBasedMatcher) matchCategory(condition map[string]interface{}, context map[string]interface{}) (bool, float64) {
    requiredCategory, ok := condition["category"].(string)
    if !ok {
        return false, 0
    }

    currentCategory, ok := context["category"].(string)
    if !ok {
        return false, 0
    }

    if currentCategory == requiredCategory {
        return true, 1.0
    }

    return false, 0
}

func contains(s, substr string) bool {
    return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
        (len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr))))
}
```

##### Day 4-5: 相似度匹配器

**相似度匹配器实现** (`domain/routing/service/similarity_matcher.go`):
```go
// service/similarity_matcher.go
package service

import (
    "context"
    "math"
    "sort"
    "github.com/coze-studio/domain/bot/repository"
    "github.com/coze-studio/domain/workflow/repository"
)

type SimilarityMatcher struct {
    botRepo     repository.BotRepository
    workflowRepo repository.WorkflowRepository
    embeddingClient EmbeddingClient // 嵌入向量客户端
}

type EmbeddingClient interface {
    GetEmbedding(ctx context.Context, text string) ([]float32, error)
}

type CandidateBot struct {
    BotID      string
    Name       string
    Desc       string
    Embedding  []float32
}

func NewSimilarityMatcher(
    botRepo repository.BotRepository,
    workflowRepo repository.WorkflowRepository,
    embeddingClient EmbeddingClient,
) *SimilarityMatcher {
    return &SimilarityMatcher{
        botRepo:         botRepo,
        workflowRepo:    workflowRepo,
        embeddingClient: embeddingClient,
    }
}

// Match 相似度匹配
func (m *SimilarityMatcher) Match(ctx context.Context, input *MatchInput) ([]*MatchOutput, error) {
    // 1. 获取用户输入的嵌入向量
    userEmbedding, err := m.embeddingClient.GetEmbedding(ctx, input.UserInput)
    if err != nil {
        return nil, err
    }

    // 2. 获取租户下所有启用的 Bot
    bots, err := m.botRepo.ListByTenant(ctx, input.TenantID, &repository.BotFilter{
        Status: "active",
    })
    if err != nil {
        return nil, err
    }

    results := make([]*MatchOutput, 0)

    // 3. 计算相似度
    for _, bot := range bots {
        if bot.Embedding == nil {
            continue
        }

        similarity := cosineSimilarity(userEmbedding, bot.Embedding)

        if similarity > 0.7 { // 阈值
            results = append(results, &MatchOutput{
                BotID:      bot.BotID,
                Confidence: similarity,
                MatchType:  "similarity",
            })
        }
    }

    // 4. 按相似度排序
    sort.Slice(results, func(i, j int) bool {
        return results[i].Confidence > results[j].Confidence
    })

    // 5. 返回 Top-3
    if len(results) > 3 {
        results = results[:3]
    }

    return results, nil
}

// cosineSimilarity 余弦相似度计算
func cosineSimilarity(a, b []float32) float64 {
    if len(a) != len(b) {
        return 0
    }

    var dotProduct float32
    var normA float32
    var normB float32

    for i := 0; i < len(a); i++ {
        dotProduct += a[i] * b[i]
        normA += a[i] * a[i]
        normB += b[i] * b[i]
    }

    if normA == 0 || normB == 0 {
        return 0
    }

    return float64(dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB)))))
}
```

#### Week 4: 路由决策引擎

**目标**: 实现基于评分的路由决策引擎

##### Day 1-3: 混合意图匹配器

**混合匹配器实现** (`domain/routing/service/hybrid_matcher.go`):
```go
// service/hybrid_matcher.go
package service

import (
    "context"
    "sort"
)

type HybridIntentMatcher struct {
    ruleMatcher      *RuleBasedMatcher
    similarityMatcher *SimilarityMatcher
    modelMatcher     *ModelBasedMatcher // 可选：基于模型的意图识别
}

type ModelBasedMatcher struct {
    // 模型客户端配置
}

func NewHybridIntentMatcher(
    ruleMatcher *RuleBasedMatcher,
    similarityMatcher *SimilarityMatcher,
    modelMatcher *ModelBasedMatcher,
) *HybridIntentMatcher {
    return &HybridIntentMatcher{
        ruleMatcher:       ruleMatcher,
        similarityMatcher: similarityMatcher,
        modelMatcher:      modelMatcher,
    }
}

// Match 混合意图匹配
func (m *HybridIntentMatcher) Match(ctx context.Context, input *MatchInput) ([]*MatchOutput, error) {
    results := make([]*MatchOutput, 0)

    // 1. 规则匹配（权重 40%）
    ruleResults, err := m.ruleMatcher.Match(ctx, input)
    if err == nil {
        for _, r := range ruleResults {
            r.Confidence *= 0.4
            results = append(results, r)
        }
    }

    // 2. 相似度匹配（权重 40%）
    simResults, err := m.similarityMatcher.Match(ctx, input)
    if err == nil {
        for _, r := range simResults {
            r.Confidence *= 0.4
            results = append(results, r)
        }
    }

    // 3. 模型匹配（权重 20%，可选）
    if m.modelMatcher != nil {
        modelResults, err := m.modelMatcher.Match(ctx, input)
        if err == nil {
            for _, r := range modelResults {
                r.Confidence *= 0.2
                results = append(results, r)
            }
        }
    }

    // 4. 聚合结果（相同 Bot 的置信度累加）
    aggregated := m.aggregateResults(results)

    // 5. 按置信度排序
    sort.Slice(aggregated, func(i, j int) bool {
        return aggregated[i].Confidence > aggregated[j].Confidence
    })

    // 6. 返回 Top-3
    if len(aggregated) > 3 {
        aggregated = aggregated[:3]
    }

    return aggregated, nil
}

// aggregateResults 聚合结果
func (m *HybridIntentMatcher) aggregateResults(results []*MatchOutput) []*MatchOutput {
    aggMap := make(map[string]*MatchOutput)

    for _, r := range results {
        key := r.BotID
        if key == "" {
            key = r.WorkflowID
        }

        if existing, ok := aggMap[key]; ok {
            existing.Confidence += r.Confidence
        } else {
            aggMap[key] = r
        }
    }

    aggregated := make([]*MatchOutput, 0, len(aggMap))
    for _, v := range aggMap {
        aggregated = append(aggregated, v)
    }

    return aggregated
}
```

##### Day 4-5: 路由决策引擎

**路由决策引擎实现** (`domain/routing/service/routing_engine.go`):
```go
// service/routing_engine.go
package service

import (
    "context"
    "sort"
    "time"
)

type ScoreBasedRouter struct {
    intentMatcher    IntentMatcher
    serviceRegistry  ServiceRegistry
    loadMonitor      LoadMonitor
}

type IntentMatcher interface {
    Match(ctx context.Context, input *MatchInput) ([]*MatchOutput, error)
}

type ServiceRegistry interface {
    GetBotHealth(ctx context.Context, botID string) (*ServiceHealth, error)
    GetWorkflowHealth(ctx context.Context, workflowID string) (*ServiceHealth, error)
}

type LoadMonitor interface {
    GetServiceLoad(ctx context.Context, serviceID string) (*ServiceLoad, error)
}

type ServiceHealth struct {
    ServiceID    string
    IsHealthy    bool
    SuccessRate  float64
    AvgLatency   time.Duration
}

type ServiceLoad struct {
    ServiceID    string
    CurrentLoad  int // 当前请求数
    MaxCapacity  int // 最大容量
}

type RoutingDecision struct {
    BotID       string
    WorkflowID  string
    Confidence  float64
    Score       float64
    Reasons     []string
}

func NewScoreBasedRouter(
    intentMatcher IntentMatcher,
    serviceRegistry ServiceRegistry,
    loadMonitor LoadMonitor,
) *ScoreBasedRouter {
    return &ScoreBasedRouter{
        intentMatcher:   intentMatcher,
        serviceRegistry: serviceRegistry,
        loadMonitor:     loadMonitor,
    }
}

// Route 路由决策
func (r *ScoreBasedRouter) Route(ctx context.Context, input *MatchInput) (*RoutingDecision, error) {
    // 1. 意图匹配，获取候选列表
    candidates, err := r.intentMatcher.Match(ctx, input)
    if err != nil || len(candidates) == 0 {
        return nil, err
    }

    // 2. 为每个候选计算综合得分
    scoredCandidates := make([]*RoutingDecision, 0)

    for _, candidate := range candidates {
        score, reasons := r.calculateScore(ctx, candidate)

        scoredCandidates = append(scoredCandidates, &RoutingDecision{
            BotID:      candidate.BotID,
            WorkflowID: candidate.WorkflowID,
            Confidence: candidate.Confidence,
            Score:      score,
            Reasons:    reasons,
        })
    }

    // 3. 按得分排序
    sort.Slice(scoredCandidates, func(i, j int) bool {
        return scoredCandidates[i].Score > scoredCandidates[j].Score
    })

    // 4. 返回得分最高的
    return scoredCandidates[0], nil
}

// calculateScore 计算综合得分
func (r *ScoreBasedRouter) calculateScore(ctx context.Context, candidate *MatchOutput) (float64, []string) {
    reasons := make([]string, 0)
    score := 0.0

    // 1. 意图匹配得分（权重 30%）
    intentScore := candidate.Confidence * 0.3
    score += intentScore
    reasons = append(reasons, fmt.Sprintf("Intent match: %.2f", intentScore))

    // 2. 服务健康度得分（权重 20%）
    var serviceID string
    if candidate.BotID != "" {
        serviceID = candidate.BotID
    } else {
        serviceID = candidate.WorkflowID
    }

    var health *ServiceHealth
    if candidate.BotID != "" {
        health, _ = r.serviceRegistry.GetBotHealth(ctx, candidate.BotID)
    } else {
        health, _ = r.serviceRegistry.GetWorkflowHealth(ctx, candidate.WorkflowID)
    }

    healthScore := 0.0
    if health != nil && health.IsHealthy {
        healthScore = health.SuccessRate * 0.2
    }
    score += healthScore
    reasons = append(reasons, fmt.Sprintf("Health: %.2f", healthScore))

    // 3. 负载得分（权重 20%）
    load, _ := r.loadMonitor.GetServiceLoad(ctx, serviceID)
    loadScore := 0.0
    if load != nil {
        loadRatio := float64(load.CurrentLoad) / float64(load.MaxCapacity)
        loadScore = (1 - loadRatio) * 0.2
    }
    score += loadScore
    reasons = append(reasons, fmt.Sprintf("Load: %.2f", loadScore))

    // 4. 区域亲和性得分（权重 10%）
    regionScore := 0.1 // 简化实现
    score += regionScore
    reasons = append(reasons, fmt.Sprintf("Region: %.2f", regionScore))

    // 5. 成本得分（权重 10%）
    costScore := 0.1 // 简化实现
    score += costScore
    reasons = append(reasons, fmt.Sprintf("Cost: %.2f", costScore))

    return score, reasons
}
```

---

### Week 5-6: RBAC 高级特性 + 配额管理

#### Week 5: 数据权限增强

**目标**: 实现部门级数据权限、自定义过滤条件

##### Day 1-3: 部门数据权限

**数据库表设计**:
```sql
-- 9. 部门表
CREATE TABLE departments (
    department_id VARCHAR(36) PRIMARY KEY COMMENT '部门ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    department_name VARCHAR(200) NOT NULL COMMENT '部门名称',
    parent_department_id VARCHAR(36) COMMENT '父部门ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    FOREIGN KEY (parent_department_id) REFERENCES departments(department_id) ON DELETE SET NULL,
    UNIQUE KEY uk_tenant_name (tenant_id, department_name, deleted_at),
    INDEX idx_tenant_id (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='部门表';

-- 10. 用户部门关联表
CREATE TABLE user_departments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(36) NOT NULL COMMENT '用户ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    department_id VARCHAR(36) NOT NULL COMMENT '部门ID',
    is_leader BOOLEAN DEFAULT FALSE COMMENT '是否是部门领导',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    FOREIGN KEY (department_id) REFERENCES departments(department_id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_tenant_dept (user_id, tenant_id, department_id),
    INDEX idx_user_id (user_id),
    INDEX idx_department_id (department_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户部门关联表';
```

**部门权限检查增强** (`domain/permission/service/department_permission_checker.go`):
```go
// service/department_permission_checker.go
package service

import (
    "context"
    "github.com/coze-studio/domain/permission/repository"
)

type DepartmentPermissionChecker struct {
    userDeptRepo    repository.UserDepartmentRepository
    departmentRepo  repository.DepartmentRepository
}

func NewDepartmentPermissionChecker(
    userDeptRepo repository.UserDepartmentRepository,
    departmentRepo repository.DepartmentRepository,
) *DepartmentPermissionChecker {
    return &DepartmentPermissionChecker{
        userDeptRepo:   userDeptRepo,
        departmentRepo: departmentRepo,
    }
}

// GetAccessibleDepartmentIDs 获取用户可访问的部门ID列表
func (d *DepartmentPermissionChecker) GetAccessibleDepartmentIDs(ctx context.Context, userID, tenantID string) ([]string, error) {
    // 1. 获取用户所属部门
    userDepts, err := d.userDeptRepo.GetByUser(ctx, userID, tenantID)
    if err != nil {
        return nil, err
    }

    deptIDs := make([]string, 0)
    for _, ud := range userDepts {
        deptIDs = append(deptIDs, ud.DepartmentID)

        // 2. 如果是部门领导，可以访问子部门
        if ud.IsLeader {
            childDepts, _ := d.departmentRepo.GetDescendants(ctx, ud.DepartmentID)
            for _, child := range childDepts {
                deptIDs = append(deptIDs, child.DepartmentID)
            }
        }
    }

    return deptIDs, nil
}

// FilterResourcesByDepartment 按部门过滤资源
func (d *DepartmentPermissionChecker) FilterResourcesByDepartment(
    ctx context.Context,
    userID, tenantID string,
    resources []interface{},
) ([]interface{}, error) {
    accessibleDeptIDs, err := d.GetAccessibleDepartmentIDs(ctx, userID, tenantID)
    if err != nil {
        return nil, err
    }

    filtered := make([]interface{}, 0)

    for _, resource := range resources {
        // 根据资源类型检查 department_id
        // 如果资源的 department_id 在 accessibleDeptIDs 中，则保留
        // 实现省略...
    }

    return filtered, nil
}
```

##### Day 4-5: 自定义过滤条件

**自定义过滤器引擎** (`domain/permission/service/custom_filter_engine.go`):
```go
// service/custom_filter_engine.go
package service

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/coze-studio/domain/permission/entity"
)

type CustomFilterEngine struct {
    // 可集成表达式引擎，如 govaluate
}

type FilterCondition struct {
    Field    string      `json:"field"`
    Operator string      `json:"operator"` // eq, ne, gt, lt, in, contains
    Value    interface{} `json:"value"`
    Logic    string      `json:"logic"`    // AND, OR
}

func NewCustomFilterEngine() *CustomFilterEngine {
    return &CustomFilterEngine{}
}

// Match 检查资源是否匹配自定义过滤条件
func (e *CustomFilterEngine) Match(ctx context.Context, customFilterJSON string, resource map[string]interface{}) (bool, error) {
    var conditions []FilterCondition
    if err := json.Unmarshal([]byte(customFilterJSON), &conditions); err != nil {
        return false, err
    }

    return e.evaluateConditions(resource, conditions), nil
}

// evaluateConditions 评估条件
func (e *CustomFilterEngine) evaluateConditions(resource map[string]interface{}, conditions []FilterCondition) bool {
    if len(conditions) == 0 {
        return true
    }

    result := true
    currentLogic := "AND"

    for _, condition := range conditions {
        matched := e.evaluateCondition(resource, condition)

        if condition.Logic != "" {
            currentLogic = condition.Logic
        }

        if currentLogic == "AND" {
            result = result && matched
        } else if currentLogic == "OR" {
            result = result || matched
        }
    }

    return result
}

// evaluateCondition 评估单个条件
func (e *CustomFilterEngine) evaluateCondition(resource map[string]interface{}, condition FilterCondition) bool {
    fieldValue, exists := resource[condition.Field]
    if !exists {
        return false
    }

    switch condition.Operator {
    case "eq":
        return fmt.Sprintf("%v", fieldValue) == fmt.Sprintf("%v", condition.Value)
    case "ne":
        return fmt.Sprintf("%v", fieldValue) != fmt.Sprintf("%v", condition.Value)
    case "gt":
        return compareNumbers(fieldValue, condition.Value) > 0
    case "lt":
        return compareNumbers(fieldValue, condition.Value) < 0
    case "in":
        return isInArray(fieldValue, condition.Value)
    case "contains":
        return containsString(fieldValue, condition.Value)
    default:
        return false
    }
}

// 辅助函数省略...
```

#### Week 6: 配额管理和计费

**目标**: 实现配额实时监控、超额处理、计费统计

##### Day 1-3: 配额监控服务

**配额监控服务** (`domain/tenant/service/quota_monitor.go`):
```go
// service/quota_monitor.go
package service

import (
    "context"
    "time"
    "github.com/coze-studio/domain/tenant/repository"
)

type QuotaMonitor struct {
    quotaRepo      repository.QuotaRepository
    alertService   AlertService
}

type AlertService interface {
    SendAlert(ctx context.Context, alert *QuotaAlert) error
}

type QuotaAlert struct {
    TenantID      string
    ResourceType  string
    AlertType     string // "warning", "critical", "exceeded"
    Usage         int
    MaxLimit      int
    UsagePercent  float64
}

func NewQuotaMonitor(quotaRepo repository.QuotaRepository, alertService AlertService) *QuotaMonitor {
    return &QuotaMonitor{
        quotaRepo:    quotaRepo,
        alertService: alertService,
    }
}

// MonitorQuotas 监控所有租户的配额
func (m *QuotaMonitor) MonitorQuotas(ctx context.Context) error {
    // 1. 获取所有配额
    quotas, err := m.quotaRepo.GetAll(ctx)
    if err != nil {
        return err
    }

    // 2. 检查每个配额
    for _, quota := range quotas {
        m.checkQuota(ctx, quota)
    }

    return nil
}

// checkQuota 检查单个配额
func (m *QuotaMonitor) checkQuota(ctx context.Context, quota *entity.Quota) {
    if quota.MaxLimit <= 0 {
        return // 无限制配额
    }

    usagePercent := float64(quota.UsedCount) / float64(quota.MaxLimit) * 100

    // 1. 警告：使用率超过 80%
    if usagePercent >= 80 && usagePercent < 100 {
        m.alertService.SendAlert(ctx, &QuotaAlert{
            TenantID:     quota.TenantID,
            ResourceType: quota.ResourceType,
            AlertType:    "warning",
            Usage:        quota.UsedCount,
            MaxLimit:     quota.MaxLimit,
            UsagePercent: usagePercent,
        })
    }

    // 2. 严重：使用率超过 95%
    if usagePercent >= 95 && usagePercent < 100 {
        m.alertService.SendAlert(ctx, &QuotaAlert{
            TenantID:     quota.TenantID,
            ResourceType: quota.ResourceType,
            AlertType:    "critical",
            Usage:        quota.UsedCount,
            MaxLimit:     quota.MaxLimit,
            UsagePercent: usagePercent,
        })
    }

    // 3. 超额：使用率超过 100%
    if usagePercent >= 100 {
        m.alertService.SendAlert(ctx, &QuotaAlert{
            TenantID:     quota.TenantID,
            ResourceType: quota.ResourceType,
            AlertType:    "exceeded",
            Usage:        quota.UsedCount,
            MaxLimit:     quota.MaxLimit,
            UsagePercent: usagePercent,
        })
    }
}

// ResetPeriodicQuotas 定期重置配额
func (m *QuotaMonitor) ResetPeriodicQuotas(ctx context.Context) error {
    quotas, err := m.quotaRepo.GetAll(ctx)
    if err != nil {
        return err
    }

    now := time.Now()

    for _, quota := range quotas {
        shouldReset := false

        switch quota.ResetCycle {
        case "daily":
            // 如果上次重置不是今天
            if quota.LastResetAt.Day() != now.Day() {
                shouldReset = true
            }
        case "monthly":
            // 如果上次重置不是本月
            if quota.LastResetAt.Month() != now.Month() || quota.LastResetAt.Year() != now.Year() {
                shouldReset = true
            }
        case "yearly":
            // 如果上次重置不是本年
            if quota.LastResetAt.Year() != now.Year() {
                shouldReset = true
            }
        }

        if shouldReset {
            m.quotaRepo.ResetUsage(ctx, quota.QuotaID)
        }
    }

    return nil
}
```

##### Day 4-5: 计费统计服务

**计费统计服务** (`domain/tenant/service/billing_service.go`):
```go
// service/billing_service.go
package service

import (
    "context"
    "time"
    "github.com/coze-studio/domain/tenant/entity"
    "github.com/coze-studio/domain/tenant/repository"
)

type BillingService struct {
    quotaRepo       repository.QuotaRepository
    usageLogRepo    repository.UsageLogRepository
    invoiceRepo     repository.InvoiceRepository
}

type UsageLog struct {
    LogID       string    `json:"log_id"`
    TenantID    string    `json:"tenant_id"`
    ResourceType string   `json:"resource_type"`
    Action      string    `json:"action"` // "create", "update", "delete", "query"
    Quantity    int       `json:"quantity"`
    Timestamp   time.Time `json:"timestamp"`
}

type Invoice struct {
    InvoiceID       string    `json:"invoice_id"`
    TenantID        string    `json:"tenant_id"`
    BillingCycle    string    `json:"billing_cycle"`
    StartDate       time.Time `json:"start_date"`
    EndDate         time.Time `json:"end_date"`
    TotalUsage      int       `json:"total_usage"`
    TotalAmount     float64   `json:"total_amount"`
    Status          string    `json:"status"` // "pending", "paid", "overdue"
    CreatedAt       time.Time `json:"created_at"`
}

func NewBillingService(
    quotaRepo repository.QuotaRepository,
    usageLogRepo repository.UsageLogRepository,
    invoiceRepo repository.InvoiceRepository,
) *BillingService {
    return &BillingService{
        quotaRepo:    quotaRepo,
        usageLogRepo: usageLogRepo,
        invoiceRepo:  invoiceRepo,
    }
}

// RecordUsage 记录使用量
func (s *BillingService) RecordUsage(ctx context.Context, log *UsageLog) error {
    return s.usageLogRepo.Create(ctx, log)
}

// GenerateInvoice 生成账单
func (s *BillingService) GenerateInvoice(ctx context.Context, tenantID string, startDate, endDate time.Time) (*Invoice, error) {
    // 1. 获取该时间段内的使用记录
    logs, err := s.usageLogRepo.GetByTenantAndDateRange(ctx, tenantID, startDate, endDate)
    if err != nil {
        return nil, err
    }

    // 2. 统计总使用量
    totalUsage := 0
    for _, log := range logs {
        totalUsage += log.Quantity
    }

    // 3. 计算费用
    totalAmount := s.calculateAmount(ctx, tenantID, logs)

    // 4. 创建账单
    invoice := &Invoice{
        InvoiceID:    generateID(),
        TenantID:     tenantID,
        BillingCycle: "monthly",
        StartDate:    startDate,
        EndDate:      endDate,
        TotalUsage:   totalUsage,
        TotalAmount:  totalAmount,
        Status:       "pending",
        CreatedAt:    time.Now(),
    }

    if err := s.invoiceRepo.Create(ctx, invoice); err != nil {
        return nil, err
    }

    return invoice, nil
}

// calculateAmount 计算费用
func (s *BillingService) calculateAmount(ctx context.Context, tenantID string, logs []*UsageLog) float64 {
    // 根据订阅等级和使用量计算费用
    // 实现省略...
    return 0.0
}
```

---

### Week 7-8: 集成测试 + 性能优化

#### Week 7: 集成测试

**目标**: 编写完整的集成测试用例

##### Day 1-3: 单元测试

**配额服务单元测试** (`domain/tenant/service/quota_service_test.go`):
```go
// service/quota_service_test.go
package service

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock QuotaRepository
type MockQuotaRepository struct {
    mock.Mock
}

func (m *MockQuotaRepository) GetByTenantAndResource(ctx context.Context, tenantID, resourceType string) (*entity.Quota, error) {
    args := m.Called(ctx, tenantID, resourceType)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Quota), args.Error(1)
}

func (m *MockQuotaRepository) UpdateUsedCount(ctx context.Context, quotaID string, delta int) error {
    args := m.Called(ctx, quotaID, delta)
    return args.Error(0)
}

func TestQuotaService_CheckQuota(t *testing.T) {
    // Arrange
    mockRepo := new(MockQuotaRepository)
    service := NewQuotaService(mockRepo)

    quota := &entity.Quota{
        QuotaID:      "quota_123",
        TenantID:     "tenant_123",
        ResourceType: "bots",
        MaxLimit:     10,
        UsedCount:    5,
    }

    mockRepo.On("GetByTenantAndResource", mock.Anything, "tenant_123", "bots").Return(quota, nil)

    // Act
    err := service.CheckQuota(context.Background(), "tenant_123", "bots", 3)

    // Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}

func TestQuotaService_CheckQuota_Exceeded(t *testing.T) {
    // Arrange
    mockRepo := new(MockQuotaRepository)
    service := NewQuotaService(mockRepo)

    quota := &entity.Quota{
        QuotaID:      "quota_123",
        TenantID:     "tenant_123",
        ResourceType: "bots",
        MaxLimit:     10,
        UsedCount:    9,
    }

    mockRepo.On("GetByTenantAndResource", mock.Anything, "tenant_123", "bots").Return(quota, nil)

    // Act
    err := service.CheckQuota(context.Background(), "tenant_123", "bots", 2)

    // Assert
    assert.Error(t, err)
    assert.Equal(t, "QUOTA_EXCEEDED", err.(*errno.Errno).Code)
    mockRepo.AssertExpectations(t)
}

func TestQuotaService_ConsumeQuota(t *testing.T) {
    // Arrange
    mockRepo := new(MockQuotaRepository)
    service := NewQuotaService(mockRepo)

    quota := &entity.Quota{
        QuotaID:      "quota_123",
        TenantID:     "tenant_123",
        ResourceType: "bots",
        MaxLimit:     10,
        UsedCount:    5,
    }

    mockRepo.On("GetByTenantAndResource", mock.Anything, "tenant_123", "bots").Return(quota, nil)
    mockRepo.On("UpdateUsedCount", mock.Anything, "quota_123", 3).Return(nil)

    // Act
    err := service.ConsumeQuota(context.Background(), "tenant_123", "bots", 3)

    // Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

##### Day 4-5: 集成测试

**集成测试** (`application/tenant/integration_test.go`):
```go
// application/tenant/integration_test.go
package tenant

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/coze-studio/infra/db"
    "github.com/coze-studio/infra/repository/tenant"
    "github.com/coze-studio/domain/tenant/service"
    "github.com/coze-studio/application/tenant/dto"
)

func TestTenantService_CreateTenant(t *testing.T) {
    // Setup
    db := setupTestDB()
    defer cleanupTestDB(db)

    tenantRepo := tenant_repository.NewTenantRepository(db)
    quotaRepo := tenant_repository.NewQuotaRepository(db)
    quotaService := service.NewQuotaService(quotaRepo)
    tenantService := NewTenantService(tenantRepo, quotaService)

    // Act
    req := &dto.CreateTenantRequest{
        TenantName: "Test Tenant",
        TenantType: "enterprise",
    }

    result, err := tenantService.CreateTenant(context.Background(), req)

    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, result.TenantID)
    assert.Equal(t, "Test Tenant", result.TenantName)
    assert.Equal(t, "enterprise", result.TenantType)

    // Verify database
    saved, err := tenantRepo.GetByID(context.Background(), result.TenantID)
    assert.NoError(t, err)
    assert.Equal(t, result.TenantID, saved.TenantID)
}

func TestTenantService_QuotaEnforcement(t *testing.T) {
    // Setup
    db := setupTestDB()
    defer cleanupTestDB(db)

    // Create tenant with quota
    tenantID := createTestTenant(db, 10) // Max 10 bots

    tenantRepo := tenant_repository.NewTenantRepository(db)
    quotaRepo := tenant_repository.NewQuotaRepository(db)
    quotaService := service.NewQuotaService(quotaRepo)
    botService := NewBotService(tenantRepo, quotaService)

    // Act: Create 10 bots (should succeed)
    for i := 0; i < 10; i++ {
        err := botService.CreateBot(context.Background(), tenantID)
        assert.NoError(t, err)
    }

    // Act: Create 11th bot (should fail)
    err := botService.CreateBot(context.Background(), tenantID)
    assert.Error(t, err)
    assert.Equal(t, "QUOTA_EXCEEDED", err.(*errno.Errno).Code)
}
```

#### Week 8: 性能优化 + 文档

**目标**: 优化性能瓶颈，编写技术文档

##### Day 1-3: 性能优化

**优化点**:

1. **数据库查询优化**
```go
// Before: N+1 查询
func (s *TenantService) GetList(ctx context.Context) ([]*entity.Tenant, error) {
    tenants, _ := s.tenantRepo.List(ctx, nil)

    for _, tenant := range tenants {
        subscription, _ := s.subRepo.GetByTenantID(ctx, tenant.TenantID)
        tenant.Subscription = subscription
    }

    return tenants, nil
}

// After: 使用 Preload 避免循环查询
func (s *TenantService) GetList(ctx context.Context) ([]*entity.Tenant, error) {
    var tenants []*entity.Tenant
    err := s.db.
        Preload("Subscription").
        Preload("Quotas").
        Find(&tenants).Error

    return tenants, err
}
```

2. **并发处理优化**
```go
// Before: 串行处理
func (s *QuotaMonitor) MonitorQuotas(ctx context.Context) error {
    quotas, _ := s.quotaRepo.GetAll(ctx)

    for _, quota := range quotas {
        s.checkQuota(ctx, quota)
    }

    return nil
}

// After: 并发处理（限制并发数）
func (s *QuotaMonitor) MonitorQuotas(ctx context.Context) error {
    quotas, _ := s.quotaRepo.GetAll(ctx)

    sem := make(chan struct{}, 10) // 最多 10 个并发
    var wg sync.WaitGroup

    for _, quota := range quotas {
        wg.Add(1)
        sem <- struct{}{} // 获取信号量
        go func(q *entity.Quota) {
            defer wg.Done()
            defer func() { <-sem }() // 释放信号量
            s.checkQuota(ctx, q)
        }(quota)
    }

    wg.Wait()
    return nil
}
```

3. **缓存优化**
```go
// 添加 Redis 缓存
func (s *TenantService) GetByID(ctx context.Context, tenantID string) (*entity.Tenant, error) {
    // 1. 先查缓存
    cacheKey := fmt.Sprintf("tenant:%s", tenantID)
    cached, err := s.redis.Get(ctx, cacheKey)
    if err == nil {
        var tenant entity.Tenant
        if err := json.Unmarshal([]byte(cached), &tenant); err == nil {
            return &tenant, nil
        }
    }

    // 2. 缓存未命中，查数据库
    tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存（5分钟过期）
    data, _ := json.Marshal(tenant)
    s.redis.Set(ctx, cacheKey, data, 5*time.Minute)

    return tenant, nil
}
```

##### Day 4-5: 技术文档

**编写文档**:
1. 租户系统 API 文档
2. 权限系统使用指南
3. 智能路由配置指南
4. 配额管理运维手册

---

## 📦 专属开发规范

### 作为架构师的特殊要求

#### 1. 接口设计优先

✅ **必须**：
- 先定义接口（Interface），再实现
- 接口定义在 `domain/` 层
- 接口方法命名清晰、参数合理

❌ **禁止**：
- 直接实现而不定义接口
- 在 `infra/` 层定义接口
- 接口方法过于复杂

#### 2. 依赖倒置原则

✅ **必须**：
- 高层模块不依赖低层模块
- 抽象不依赖具体实现
- 使用依赖注入

```go
// ✅ Good: 依赖抽象
type BotService struct {
    repo BotRepository // 接口
}

// ❌ Bad: 依赖具体实现
type BotService struct {
    repo *MySQLBotRepository // 具体实现
}
```

#### 3. 领域逻辑纯粹

✅ **必须**：
- 领域层不依赖基础设施
- 业务逻辑在领域服务中
- 实体包含业务规则

❌ **禁止**：
- 领域层直接调用数据库
- 领域层直接调用 HTTP 客户端
- 业务逻辑泄露到 API 层

#### 4. 数据库设计规范

✅ **必须**：
- 所有表都有 `tenant_id`（除系统表）
- 使用 `deleted_at` 软删除
- 外键约束完整
- 索引合理设计

❌ **禁止**：
- 硬删除（`DELETE FROM`）
- 缺少外键约束
- 缺少必要索引

---

## 🤝 协作接口定义

### 与研发 B 协作

**依赖研发 B**:
- 统一错误码定义（`types/errno/`）
- 性能测试脚本
- 监控埋点规范

**向研发 B 提供**:
- OpenAPI 规范文件（`idl/*.thrift`）
- 数据库 Schema（`migrations/`）
- API 契约测试用例

### 与研发 C 协作

**依赖研发 C**:
- 前端 Mock 数据
- UI/UX 反馈

**向研发 C 提供**:
- IDL 文件（`idl/tenant.thrift`、`idl/permission.thrift`）
- API 接口文档
- 数据格式示例

### 与研发 D 协作

**依赖研发 D**:
- Docker 环境配置
- K8s 部署配置
- CI/CD 流水线

**向研发 D 提供**:
- 数据库迁移脚本（`migrations/`）
- 环境变量清单
- 健康检查接口

---

## 📅 里程碑和交付物

### Week 2 交付物

- [ ] 租户实体和仓储接口
- [ ] 配额检查服务
- [ ] 订阅管理服务
- [ ] 单元测试覆盖率 ≥ 80%

### Week 4 交付物

- [ ] RBAC 基础实体和服务
- [ ] 权限检查服务
- [ ] 规则匹配器
- [ ] 相似度匹配器
- [ ] 集成测试用例

### Week 6 交付物

- [ ] 混合意图匹配器
- [ ] 路由决策引擎
- [ ] 部门权限检查
- [ ] 自定义过滤引擎
- [ ] 配额监控服务
- [ ] 计费统计服务

### Week 8 交付物

- [ ] 完整集成测试
- [ ] 性能优化报告
- [ ] API 文档
- [ ] 运维手册

---

## ✅ 质量检查清单

### 代码提交前

- [ ] 代码符合 `ZKER-企业级开发规范手册_v1.0.md`
- [ ] 单元测试通过（`go test ./...`）
- [ ] 代码覆盖率 ≥ 80%（`go test -cover`）
- [ ] golangci-lint 检查通过（`golangci-lint run`）
- [ ] API 文档更新（IDL 文件）

### 集成前

- [ ] 与其他模块接口测试通过
- [ ] 跨模块集成测试通过
- [ ] 性能测试达标（响应时间 < 500ms）
- [ ] 与研发 B/C/D 的协作接口验证

### 发布前

- [ ] 所有测试用例通过
- [ ] 数据库迁移脚本准备就绪
- [ ] API 文档完整
- [ ] 运维手册完整
- [ ] 回滚方案确认

---

## 📊 关键指标

### 代码质量

- 单元测试覆盖率：≥ 80%
- 代码审查通过率：100%
- golangci-lint 警告数：0
- 代码重复率：< 5%

### 性能指标

- API P95 响应时间：< 500ms
- 配额检查延迟：< 50ms
- 权限检查延迟：< 100ms
- 路由决策延迟：< 200ms

### 功能完整性

- Week 2: 多租户核心 100%
- Week 4: RBAC 基础 100%
- Week 6: 智能路由 100%
- Week 8: 高级特性 100%

---

## 📚 参考资料

- [ZKER-企业级开发规范手册_v1.0.md](../ZKER-企业级开发规范手册_v1.0.md)
- [ZKER-全局一致性检查清单_v1.0.md](../ZKER-全局一致性检查清单_v1.0.md)
- [ZKER-实现差距分析与研发计划_v1.0.md](../ZKER-实现差距分析与研发计划_v1.0.md)
- [ZKER-统一错误码定义规范.md](../ZKER-统一错误码定义规范.md)

---

**文档状态**: ✅ 已完成
**最后更新**: 2025-01-01
**责任人**: 研发 A - 后端架构师
**评审人**: 技术负责人
