# 多租户开发助手

**版本**: v3.0.0 | **更新**: 2025-01-03

协助多租户架构相关的开发工作，确保租户隔离和配额管理符合企业级 SaaS 标准。

---

## 🎯 使用场景

当你需要开发多租户相关功能时，使用此技能：

### 核心场景

- **创建租户相关实体** - Tenant、Subscription、Quota、QuotaUsage
- **实现租户隔离** - 数据隔离、查询过滤、权限检查
- **配额管理** - 配额检查、配额使用、配额告警
- **订阅管理** - 订阅计划、订阅状态、订阅升级
- **数据迁移** - tenant_id 字段迁移、数据回填

### 典型任务

```
"为 conversations 表添加租户隔离"
"创建配额检查服务"
"实现订阅升级逻辑"
"生成 tenant_id 迁移脚本"
"检查代码是否遗漏租户隔离"
```

---

## 🔧 代码模板

### 1. 租户实体（Tenant Entity）

\`\`\`go
// domain/tenant/entity/tenant.go
package tenant

import (
    "time"
    "gorm.io/gorm"
)

// Tenant 租户实体
type Tenant struct {
    TenantID   string    `gorm:"primaryKey;type:varchar(36)" json:"tenant_id"`
    TenantName string    `gorm:"type:varchar(100);not null" json:"tenant_name"`
    PlanType   string    `gorm:"type:varchar(20);not null;default:'FREE'" json:"plan_type"`
    Status     string    `gorm:"type:varchar(20);not null;default:'ACTIVE'" json:"status"`
    MaxBots    int32     `gorm:"not null;default:10" json:"max_bots"`
    UsedBots   int32     `gorm:"not null;default:0" json:"used_bots"`
    CreatedAt  time.Time      `json:"created_at"`
    UpdatedAt  time.Time      `json:"updated_at"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (t *Tenant) IsActive() bool {
    return t.Status == "ACTIVE" && t.DeletedAt.Time.IsZero()
}
\`\`\`

### 2. 配额检查服务

\`\`\`go
// application/quota/checker.go
func (s *QuotaChecker) CheckCreateBot(ctx context.Context, tenantID string) error {
    tenant, err := s.tenantRepo.GetByTenantID(ctx, tenantID)
    if err != nil {
        return errorx.Wrapf(err, errno.TenantNotFound)
    }
    
    if !tenant.IsActive() {
        return errorx.New(errno.TenantInactive)
    }
    
    if tenant.UsedBots >= tenant.MaxBots {
        return errorx.New(errno.QuotaExceeded,
            "bot count exceeds quota: %d/%d", tenant.UsedBots, tenant.MaxBots)
    }
    
    return nil
}
\`\`\`

---

## 📋 检查清单

### 创建租户相关功能

- [ ] 所有实体都有 TenantID 字段
- [ ] 所有查询都包含 tenant_id 过滤
- [ ] 创建资源前检查配额
- [ ] 使用统一错误码

### 数据迁移

- [ ] 备份原表数据
- [ ] 添加 tenant_id 字段
- [ ] 回填历史数据
- [ ] 添加索引

---

## 📖 相关文档

- [03-DESIGN/multi-tenant/](../../docs/03-DESIGN/multi-tenant/)
- [研发A-后端架构师开发计划](../../docs/企业级功能完善与统一性设计方案/研发A-后端架构师开发计划_v1.0.md)

**🎯 目标**: 确保多租户架构完整、安全、高效！
