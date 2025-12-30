# Session/User 表 tenant_id 迁移报告

> **执行日期**: 2025-12-30
> **优先级**: P0（企业级多租户架构完整性）
> **状态**: ✅ 代码准备完成，待执行

---

## 📋 执行摘要

本迁移方案为 Session 和 User 表添加 `tenant_id` 字段，并创建 `user_tenant` 关联表，实现企业级多租户架构的完整数据隔离。

### 迁移范围

| 表名 | 变更类型 | 说明 |
|------|---------|------|
| `session` | 添加 `tenant_id` 字段 | 支持会话级别的租户隔离 |
| `user` | 添加 `tenant_id` 字段 | 支持用户级别的租户隔离 |
| `user_tenant` | 新建关联表 | 支持用户多租户场景 |

---

## 🗄️ 数据库变更

### 1. session 表变更

```sql
ALTER TABLE session
ADD COLUMN tenant_id VARCHAR(64) NULL
    COMMENT '租户ID（多租户隔离）'
    AFTER user_id,
ADD INDEX idx_session_tenant (tenant_id),
ADD INDEX idx_session_user_tenant (user_id, tenant_id);
```

### 2. user 表变更

```sql
ALTER TABLE `user`
ADD COLUMN tenant_id VARCHAR(64) NULL
    COMMENT '租户ID（多租户隔离）'
    AFTER id,
ADD INDEX idx_user_tenant (tenant_id),
ADD INDEX idx_user_email_tenant (email, tenant_id);
```

### 3. user_tenant 关联表（新建）

```sql
CREATE TABLE IF NOT EXISTS user_tenant (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    role VARCHAR(64) DEFAULT 'member' COMMENT '角色: owner, admin, member',
    is_default TINYINT(1) DEFAULT 0 COMMENT '是否为默认租户',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME DEFAULT NULL COMMENT '删除时间（软删除）',

    UNIQUE KEY uk_user_tenant (user_id, tenant_id, deleted_at),
    KEY idx_user_tenant_user (user_id),
    KEY idx_user_tenant_tenant (tenant_id),
    KEY idx_user_tenant_default (user_id, is_default, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='用户-租户关联表（支持用户多租户）';
```

---

## 💻 代码变更

### 1. 实体层（Entity）

| 文件 | 变更说明 |
|------|---------|
| `backend/domain/user/entity/user.go` | ✅ 已有 `TenantID` 字段 |
| `backend/domain/user/entity/session.go` | ✅ 已有 `TenantID` 字段 |
| `backend/domain/tenant/entity/user_tenant.go` | ✅ 新建 `UserTenant` 实体 |

### 2. 数据模型（GORM Model）

| 文件 | 变更说明 |
|------|---------|
| `backend/domain/user/internal/dal/model/user.gen.go` | ✅ 已添加 `TenantID` 字段 |
| `backend/domain/user/internal/dal/model/session.gen.go` | ✅ 新建 `Session` 模型 |
| `backend/domain/user/internal/dal/model/user_tenant.gen.go` | ✅ 新建 `UserTenant` 模型 |

### 3. 迁移服务（Migration Service）

| 文件 | 功能 |
|------|------|
| `backend/domain/tenant/migration/session_user_migrator.go` | ✅ 迁移执行器 |
| `backend/domain/tenant/migration/session_user_dual_write.go` | ✅ 双写适配器 |
| `backend/domain/tenant/migration/session_user_validator.go` | ✅ 迁移验证器 |

### 4. 中间件（Middleware）

| 文件 | 功能 |
|------|------|
| `backend/api/middleware/tenant_isolation.go` | ✅ 租户隔离中间件（已有） |
| `backend/api/middleware/tenant_isolation_v2.go` | ✅ Session/User表增强版中间件 |

---

## 🔄 迁移流程（Saga 模式）

### 阶段 1：准备阶段 ✅

1. ✅ 创建迁移 SQL 脚本
2. ✅ 更新数据库模型
3. ✅ 实现迁移服务
4. ✅ 实现双写支持
5. ✅ 实现验证逻辑

### 阶段 2：迁移阶段（待执行）

```bash
# 1. 执行 DDL 变更
mysql -u root -p < backend/domain/tenant/migration/003_session_user_tenant_id.sql

# 2. 运行迁移验证
go run backend/cmd/migrate/main.go --validate

# 3. 启动双写监控（可选）
go run backend/cmd/migrate/main.go --monitor --duration=24h
```

### 阶段 3：切换阶段（待执行）

1. 灰度切换（10% -> 50% -> 100%）
2. 验证数据一致性
3. 回滚准备

### 阶段 4：清理阶段（待执行）

1. 移除双写逻辑
2. 移除兼容代码
3. 清理旧数据

---

## ✅ 验证清单

### 数据库验证

- [ ] session 表 `tenant_id` 字段存在
- [ ] user 表 `tenant_id` 字段存在
- [ ] user_tenant 表已创建
- [ ] 所有索引已创建
- [ ] 外键约束已建立

### 数据验证

- [ ] session 表所有记录都有 `tenant_id`
- [ ] user 表所有记录都有 `tenant_id`
- [ ] user_tenant 关联数据完整
- [ ] 无孤儿数据

### 代码验证

- [ ] GORM 模型正确映射
- [ ] 迁移服务可正常执行
- [ ] 双写逻辑正确工作
- [ ] 验证逻辑正确工作

---

## 📊 迁移进度

| 阶段 | 状态 | 进度 |
|------|------|------|
| 阶段1：准备阶段 | ✅ 完成 | 100% |
| 阶段2：迁移阶段 | ⏸️ 待执行 | 0% |
| 阶段3：切换阶段 | ⏸️ 待开始 | 0% |
| 阶段4：清理阶段 | ⏸️ 待开始 | 0% |

**总体进度**: 25% （代码准备完成）

---

## 🛡️ 风险控制

### 双写保护 ✅

- ✅ 实现了 `SessionUserDualWriteAdapter`
- ✅ 支持数据一致性验证
- ✅ 支持实时监控

### 回滚方案 ✅

```sql
-- 回滚脚本（谨慎使用！）
DROP TABLE IF EXISTS user_tenant;
ALTER TABLE `user` DROP COLUMN tenant_id;
ALTER TABLE session DROP COLUMN tenant_id;
```

### 监控告警 ✅

- ✅ 实现了 `ProgressMonitor`
- ✅ 实现了 `CheckpointManager`
- ✅ 支持断点续传

### 灰度发布 ✅

- ✅ 支持分阶段切换
- ✅ 支持实时验证
- ✅ 支持快速回滚

---

## 📝 使用指南

### 执行迁移

```go
package main

import (
    "context"
    "log"
    "github.com/coze-dev/coze-studio/backend/domain/tenant/migration"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    // 1. 连接数据库
    dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 2. 创建迁移器
    migrator := migration.NewSessionUserMigrator(db, 1000, "default")

    // 3. 执行迁移
    ctx := context.Background()
    if err := migrator.Migrate(ctx); err != nil {
        log.Fatal(err)
    }

    log.Println("迁移完成！")
}
```

### 验证迁移

```go
// 创建验证器
validator := migration.NewSessionUserValidator(db)

// 执行验证
report, err := validator.ValidateMigration(ctx)
if err != nil {
    log.Fatal(err)
}

// 打印报告
validator.PrintReport(report)
```

### 监控双写

```go
// 创建双写适配器
adapter := migration.NewSessionUserDualWriteAdapter(db)

// 启动监控（每5分钟检查一次，持续24小时）
err := adapter.MonitorDualWrite(ctx, 5*time.Minute, 24*time.Hour)
if err != nil {
    log.Fatal(err)
}
```

---

## 🔗 相关文档

- [02-SPECS/multi-tenant-architecture.md](../02-SPECS/multi-tenant-architecture.md) - 多租户架构设计
- [02-SPECS/data-migration-guide.md](../02-SPECS/data-migration-guide.md) - 数据迁移指南
- [03-DESIGN/multi-tenant/migration-strategy.md](../03-DESIGN/multi-tenant/migration-strategy.md) - 迁移策略设计

---

## 📞 联系方式

如有问题，请联系：

- **研发A（后端架构师）**: 负责 tenant/permission/routing 模块
- **Issue**: [GitHub Issues](https://github.com/coze-dev/coze-studio/issues)

---

**报告生成时间**: 2025-12-30
**版本**: v1.0.0
