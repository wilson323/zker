# ZKER 数据库设计规范

**版本**: v3.0.0 | **更新**: 2025-01-03 | **状态**: 强制执行

---

## 目录

- [命名规范](#命名规范)
- [数据类型规范](#数据类型规范)
- [表设计规范](#表设计规范)
- [索引设计规范](#索引设计规范)
- [约束设计规范](#约束设计规范)
- [多租户隔离规范](#多租户隔离规范)
- [迁移管理](#迁移管理)
- [检查清单](#检查清单)

---

## 命名规范

### 表命名

#### 规则

1. **全小写**
2. **使用下划线分隔**
3. **使用复数形式**
4. **以业务含义命名**

#### 示例

```sql
-- ✅ Good
CREATE TABLE tenants (...);
CREATE TABLE users (...);
CREATE TABLE roles (...);
CREATE TABLE permissions (...);
CREATE TABLE user_roles (...);

-- ❌ Bad
CREATE TABLE Tenant (...);              -- 不要大写开头
CREATE TABLE tenant (...);               -- 要使用复数
CREATE TABLE TenantInfo (...);           -- 不要使用驼峰
CREATE TABLE tenant_info (...);          -- 不要使用 Info 后缀（除非必要）
```

### 字段命名

#### 规则

1. **全小写**
2. **使用下划线分隔**
3. **使用单数形式**
4. **以表名缩写作为前缀（多表时有歧义时）**

#### 示例

```sql
-- ✅ Good
CREATE TABLE tenants (
    tenant_id VARCHAR(36) PRIMARY KEY,
    tenant_name VARCHAR(100) NOT NULL,
    plan_type VARCHAR(20) NOT NULL DEFAULT 'free',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    quota_max_tokens BIGINT NOT NULL DEFAULT 1000000,
    quota_used_tokens BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ❌ Bad
CREATE TABLE tenants (
    TenantID VARCHAR(36) PRIMARY KEY,              -- 不要使用驼峰
    tenant_name VARCHAR(100) NOT NULL,             -- 命名不一致
    PlanType VARCHAR(20) NOT NULL,                 -- 不要使用驼峰
    max_tokens BIGINT NOT NULL DEFAULT 1000000,    -- 多表时要有前缀
    CreateTime TIMESTAMP NOT NULL                  -- 不要使用驼峰
);
```

### 标准字段

#### 主键

```sql
-- 使用 UUID 作为主键
tenant_id VARCHAR(36) PRIMARY KEY,
user_id VARCHAR(36) PRIMARY KEY,
role_id VARCHAR(36) PRIMARY KEY,
```

#### 时间戳

```sql
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
deleted_at TIMESTAMP NULL COMMENT '删除时间（软删除）',
```

#### 状态字段

```sql
status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态：active=激活, inactive=未激活, suspended=停用',
```

---

## 数据类型规范

### 字符串类型

| 类型 | 最大长度 | 使用场景 | 示例 |
|------|---------|---------|------|
| **CHAR** | 255 | 固定长度字符串 | `CHAR(36)` UUID |
| **VARCHAR** | 65535 | 变长字符串 | `VARCHAR(100)` 名称 |
| **TEXT** | 65535 | 短文本 | `TEXT` 描述 |
| **MEDIUMTEXT** | 16M | 中等文本 | `MEDIUMTEXT` 内容 |
| **LONGTEXT** | 4G | 长文本 | `LONGTEXT` 日志 |

```sql
-- ✅ Good
tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
tenant_name VARCHAR(100) NOT NULL COMMENT '租户名称',
description TEXT COMMENT '描述',
content LONGTEXT COMMENT '内容',

-- ❌ Bad
tenant_name VARCHAR(255) COMMENT '租户名称',  -- 长度太大
description VARCHAR(255) COMMENT '描述',      -- 应该使用 TEXT
```

### 数值类型

| 类型 | 范围 | 使用场景 | 示例 |
|------|------|---------|------|
| **TINYINT** | -128~127 | 布尔值、小整数 | `TINYINT` 状态 |
| **SMALLINT** | -32768~32767 | 小整数 | `SMALLINT` 优先级 |
| **INT** | ±21亿 | 整数 | `INT` 计数 |
| **BIGINT** | ±922亿亿 | 大整数 | `BIGINT` Token |
| **DECIMAL** | 精确小数 | 金额 | `DECIMAL(10,2)` 价格 |

```sql
-- ✅ Good
status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0=禁用, 1=启用',
priority SMALLINT NOT NULL DEFAULT 0 COMMENT '优先级',
count INT NOT NULL DEFAULT 0 COMMENT '计数',
tokens BIGINT NOT NULL DEFAULT 0 COMMENT 'Token数量',
price DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '价格',

-- ❌ Bad
status INT COMMENT '状态',                    -- 应该使用 TINYINT
tokens INT COMMENT 'Token数量',               -- 应该使用 BIGINT
price FLOAT COMMENT '价格',                   -- 金额应该使用 DECIMAL
```

### 时间类型

| 类型 | 格式 | 使用场景 | 示例 |
|------|------|---------|------|
| **DATE** | YYYY-MM-DD | 日期 | `DATE` 生日 |
| **TIME** | HH:MM:SS | 时间 | `TIME` 营业时间 |
| **DATETIME** | YYYY-MM-DD HH:MM:SS | 日期时间 | `DATETIME` 创建时间 |
| **TIMESTAMP** | YYYY-MM-DD HH:MM:SS | 时间戳 | `TIMESTAMP` 更新时间 |

```sql
-- ✅ Good
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
deleted_at TIMESTAMP NULL COMMENT '删除时间',
birth_date DATE COMMENT '生日',

-- ❌ Bad
create_time VARCHAR(20) COMMENT '创建时间',   -- 不要使用字符串存储时间
created_at INT COMMENT '创建时间',             -- 不要使用时间戳
```

### JSON 类型

```sql
-- ✅ Good: 使用 JSON 类型存储结构化数据
metadata JSON COMMENT '元数据',
permissions JSON COMMENT '权限列表',
quota JSON COMMENT '配额信息',

-- ❌ Bad: 不要使用 TEXT 存储 JSON
metadata TEXT COMMENT '元数据',
```

---

## 表设计规范

### 主键设计

#### 规则

1. **使用 VARCHAR(36) 存储 UUID**
2. **主键名为 `资源名_id`**
3. **每张表必须有主键**

#### 示例

```sql
-- ✅ Good
CREATE TABLE tenants (
    tenant_id VARCHAR(36) PRIMARY KEY COMMENT '租户ID',
    tenant_name VARCHAR(100) NOT NULL COMMENT '租户名称',
    ...
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户表';

-- ❌ Bad
CREATE TABLE tenants (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,  -- 不要使用自增ID（多租户场景）
    tenant_name VARCHAR(100) NOT NULL,
    ...
);

CREATE TABLE tenants (
    tenant_id VARCHAR(36),                  -- 主键要明确声明
    tenant_name VARCHAR(100) NOT NULL,
    ...
);
```

### 外键设计

#### 规则

1. **外键名为 `fk_当前表_关联表_字段`**
2. **外键字段名与关联表主键名一致**
3. **级联删除谨慎使用**

#### 示例

```sql
-- ✅ Good
CREATE TABLE user_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    user_id VARCHAR(36) NOT NULL COMMENT '用户ID',
    role_id VARCHAR(36) NOT NULL COMMENT '角色ID',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    FOREIGN KEY fk_user_roles_users_user_id (user_id) REFERENCES users(user_id),
    FOREIGN KEY fk_user_roles_roles_role_id (role_id) REFERENCES roles(role_id),
    INDEX idx_user_id (user_id),
    INDEX idx_role_id (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- ❌ Bad
CREATE TABLE user_roles (
    ...
    user_id VARCHAR(36) NOT NULL,
    role_id VARCHAR(36) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),              -- 外键名不规范
    FOREIGN KEY (role_id) REFERENCES roles(id)               -- 外键名不规范
);
```

### 软删除设计

#### 规则

1. **使用 `deleted_at` 字段**
2. **类型为 `TIMESTAMP NULL`**
3. **创建唯一索引时包含 `deleted_at`**

#### 示例

```sql
-- ✅ Good
CREATE TABLE tenants (
    tenant_id VARCHAR(36) PRIMARY KEY,
    tenant_name VARCHAR(100) NOT NULL,
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    UNIQUE KEY uk_tenant_name_deleted (tenant_name, deleted_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 查询时过滤软删除数据
SELECT * FROM tenants WHERE deleted_at IS NULL;

-- 软删除
UPDATE tenants SET deleted_at = NOW() WHERE tenant_id = ?;

-- ❌ Bad
CREATE TABLE tenants (
    ...
    is_deleted TINYINT NOT NULL DEFAULT 0,      -- 不要使用 is_deleted
    deleted_at TIMESTAMP NULL,                   -- 不要同时使用两个字段
    UNIQUE KEY uk_tenant_name (tenant_name)      -- 唯一索引要包含 deleted_at
);
```

---

## 索引设计规范

### 索引命名

#### 规则

1. **主键索引：`PRIMARY`**
2. **唯一索引：`uk_表名_字段名`**
3. **普通索引：`idx_表名_字段名`**
4. **全文索引：`ft_表名_字段名`**

#### 示例

```sql
-- ✅ Good
CREATE TABLE tenants (
    tenant_id VARCHAR(36) PRIMARY KEY COMMENT '租户ID',
    tenant_name VARCHAR(100) NOT NULL COMMENT '租户名称',
    email VARCHAR(100) NOT NULL COMMENT '邮箱',
    status VARCHAR(20) NOT NULL COMMENT '状态',
    created_at TIMESTAMP NOT NULL COMMENT '创建时间',
    UNIQUE KEY uk_tenants_tenant_name (tenant_name),
    UNIQUE KEY uk_tenants_email (email),
    INDEX idx_tenants_status (status),
    INDEX idx_tenants_created_at (created_at),
    INDEX idx_tenants_status_created (status, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ❌ Bad
CREATE TABLE tenants (
    ...
    UNIQUE KEY tenant_name (tenant_name),         -- 要有 uk_ 前缀
    INDEX status (status),                        -- 要有 idx_ 前缀
    INDEX created (created_at),                   -- 要有 idx_ 前缀
    KEY idx_status (status)                       -- 不要使用 KEY
);
```

### 索引设计原则

| 原则 | 说明 | 示例 |
|------|------|------|
| **选择性高的列** | 优先创建索引 | `tenant_id`, `email` |
| **频繁查询的列** | WHERE、ORDER BY、JOIN | `status`, `created_at` |
| **联合索引** | 遵循最左前缀原则 | `(status, created_at)` |
| **避免过多索引** | 影响 INSERT/UPDATE 性能 | 单表索引 ≤ 5 个 |

```sql
-- ✅ Good: 选择性高的列
CREATE INDEX idx_tenants_tenant_id ON tenants(tenant_id);
CREATE INDEX idx_tenants_email ON tenants(email);

-- ✅ Good: 联合索引
CREATE INDEX idx_tenants_status_created ON tenants(status, created_at);
-- 查询 1: 使用索引
SELECT * FROM tenants WHERE status = ? AND created_at > ?;
-- 查询 2: 使用索引（最左前缀）
SELECT * FROM tenants WHERE status = ?;
-- 查询 3: 不使用索引（跳过了 status）
SELECT * FROM tenants WHERE created_at > ?;

-- ❌ Bad: 选择性低的列
CREATE INDEX idx_tenants_status ON tenants(status);      -- status 只有 3 个值，选择性低
CREATE INDEX idx_tenants_gender ON users(gender);        -- gender 只有 2 个值，选择性低
```

---

## 约束设计规范

### NOT NULL 约束

#### 规则

1. **主键必须 NOT NULL**
2. **外键必须 NOT NULL（除非允许为空）**
3. **业务关键字段必须 NOT NULL**

```sql
-- ✅ Good
CREATE TABLE tenants (
    tenant_id VARCHAR(36) PRIMARY KEY NOT NULL COMMENT '租户ID',
    tenant_name VARCHAR(100) NOT NULL COMMENT '租户名称',
    plan_type VARCHAR(20) NOT NULL DEFAULT 'free' COMMENT '套餐类型',
    description TEXT COMMENT '描述',                      -- 允许为空
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### DEFAULT 约束

#### 规则

1. **状态字段必须有默认值**
2. **时间戳字段使用 CURRENT_TIMESTAMP**
3. **数值类型字段使用 0 作为默认值**

```sql
-- ✅ Good
status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态',
plan_type VARCHAR(20) NOT NULL DEFAULT 'free' COMMENT '套餐类型',
count INT NOT NULL DEFAULT 0 COMMENT '计数',
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
```

### COMMENT 注释

#### 规则

1. **所有字段都必须有注释**
2. **注释要清晰表达字段含义**
3. **枚举值要列出所有选项**

```sql
-- ✅ Good
CREATE TABLE tenants (
    tenant_id VARCHAR(36) PRIMARY KEY COMMENT '租户ID',
    tenant_name VARCHAR(100) NOT NULL COMMENT '租户名称',
    plan_type VARCHAR(20) NOT NULL DEFAULT 'free' COMMENT '套餐类型：free=免费版, pro=专业版, enterprise=企业版',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态：active=激活, inactive=未激活, suspended=停用',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户表';

-- ❌ Bad
CREATE TABLE tenants (
    tenant_id VARCHAR(36) PRIMARY KEY,                    -- 缺少注释
    tenant_name VARCHAR(100) NOT NULL,                    -- 缺少注释
    plan_type VARCHAR(20) NOT NULL DEFAULT 'free',        -- 缺少注释
    status VARCHAR(20) NOT NULL DEFAULT 'active'          -- 缺少注释
);
```

---

## 多租户隔离规范

### 租户隔离设计

#### 规则

1. **所有业务表必须包含 `tenant_id` 字段**
2. **创建索引时包含 `tenant_id`**
3. **所有查询必须过滤 `tenant_id`**

```sql
-- ✅ Good
CREATE TABLE bots (
    bot_id VARCHAR(36) PRIMARY KEY COMMENT 'Bot ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    bot_name VARCHAR(100) NOT NULL COMMENT 'Bot名称',
    ...
    INDEX idx_bots_tenant_id (tenant_id),
    INDEX idx_bots_tenant_bot (tenant_id, bot_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Bot表';

-- 查询时必须包含 tenant_id
SELECT * FROM bots WHERE tenant_id = ? AND bot_id = ?;
SELECT * FROM bots WHERE tenant_id = ? ORDER BY created_at DESC;

-- ❌ Bad
CREATE TABLE bots (
    bot_id VARCHAR(36) PRIMARY KEY,
    bot_name VARCHAR(100) NOT NULL,
    ...                                        -- 缺少 tenant_id
);

-- 查询时缺少 tenant_id
SELECT * FROM bots WHERE bot_id = ?;           -- 跨租户数据泄漏风险
```

### 中间件自动注入

```go
// ✅ Good: 使用中间件自动注入 tenant_id
func TenantIsolation() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        session := GetSession(c)
        if session == nil {
            c.JSON(401, ErrorResponse(errno.Unauthorized))
            c.Abort()
            return
        }

        // 注入 tenant_id 到上下文
        context.SetTenantID(ctx, session.TenantID)

        c.Next(ctx)
    }
}

// 查询时自动过滤
func (r *BotRepository) FindByID(ctx context.Context, botID string) (*Bot, error) {
    var bot Bot
    tenantID := context.GetTenantID(ctx)

    err := r.db.WithContext(ctx).
        Where("bot_id = ? AND tenant_id = ? AND deleted_at IS NULL", botID, tenantID).
        First(&bot).Error

    return &bot, err
}
```

---

## 迁移管理

### 迁移脚本规范

#### 命名规范

```
{version}_{description}.sql
```

#### 示例

```
20250103_01_create_tenants_table.sql
20250103_02_create_users_table.sql
20250103_03_create_roles_table.sql
```

### 迁移脚本示例

```sql
-- 20250103_01_create_tenants_table.sql
-- 创建租户表

CREATE TABLE IF NOT EXISTS tenants (
    tenant_id VARCHAR(36) PRIMARY KEY COMMENT '租户ID',
    tenant_name VARCHAR(100) NOT NULL COMMENT '租户名称',
    plan_type VARCHAR(20) NOT NULL DEFAULT 'free' COMMENT '套餐类型：free=免费版, pro=专业版, enterprise=企业版',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态：active=激活, inactive=未激活, suspended=停用',
    quota JSON NOT NULL COMMENT '配额信息',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',
    UNIQUE KEY uk_tenants_tenant_name_deleted (tenant_name, deleted_at),
    INDEX idx_tenants_status (status),
    INDEX idx_tenants_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户表';

-- 20250103_02_alter_tenants_add_max_bots.sql
-- 为租户表添加最大Bot数量字段

ALTER TABLE tenants ADD COLUMN max_bots INT NOT NULL DEFAULT 10 COMMENT '最大Bot数量' AFTER quota;
```

### 回滚脚本

```sql
-- 20250103_02_alter_tenants_add_max_bots.rollback.sql
-- 回滚：删除最大Bot数量字段

ALTER TABLE tenants DROP COLUMN max_bots;
```

---

## 检查清单

### 表设计检查清单

- [ ] 表名全小写，使用下划线分隔，复数形式
- [ ] 字段名全小写，使用下划线分隔，单数形式
- [ ] 包含标准字段：`created_at`, `updated_at`, `deleted_at`
- [ ] 所有字段都有 COMMENT 注释
- [ ] 使用 VARCHAR(36) 存储主键
- [ ] 业务表包含 `tenant_id` 字段

### 索引检查清单

- [ ] 主键索引：`PRIMARY`
- [ ] 唯一索引：`uk_表名_字段名`
- [ ] 普通索引：`idx_表名_字段名`
- [ ] 联合索引遵循最左前缀原则
- [ ] 单表索引 ≤ 5 个

### 多租户检查清单

- [ ] 所有业务表包含 `tenant_id` 字段
- [ ] 创建索引时包含 `tenant_id`
- [ ] 所有查询过滤 `tenant_id`
- [ ] 使用中间件自动注入 `tenant_id`

---

## 相关文档

| 文档 | 说明 |
|------|------|
| [naming-conventions.md](naming-conventions.md) | 命名规范 |
| [api-design.md](api-design.md) | API 设计规范 |
| [error-handling.md](error-handling.md) | 错误处理规范 |
| [backend-dev-guide.md](backend-dev-guide.md) | 后端开发指南 |

---

**🎯 目标**: 规范的数据库设计，保证数据一致性和查询性能！
