# 租户表DDL部署文档

**版本**: v1.0.0
**日期**: 2025-01-01
**作者**: ZKER Enterprise Team
**状态**: ✅ 已完成

---

## 📋 文档概述

本文档提供了完整的租户表DDL脚本部署指南，包括表结构、触发器、视图、验证脚本和回滚方案。

### 核心特性

✅ **完整的表结构** - 4个核心表（tenants、subscriptions、quotas、quota_usage_log）
✅ **自动化触发器** - 租户状态变更时自动处理订阅
✅ **统计视图** - 3个视图提供全面的使用统计
✅ **验证脚本** - 10项完整检查确保数据完整性
✅ **回滚方案** - 一键回滚所有变更
✅ **初始化数据** - 系统默认租户和配额配置

---

## 📁 文件清单

### DDL脚本文件

| 文件名 | 路径 | 说明 |
|--------|------|------|
| 创建租户表 | `docker/atlas/migrations/20251230025000_create_tenant_tables.sql` | 核心表结构DDL |
| 创建触发器 | `docker/atlas/migrations/20251230025001_tenant_status_trigger.sql` | 状态更新触发器 |
| 创建视图 | `docker/atlas/migrations/20251230025002_tenant_statistics_view.sql` | 统计视图定义 |
| 初始化数据 | `docker/atlas/migrations/20251230025003_init_tenant_data.sql` | 默认数据初始化 |

### 辅助脚本文件

| 文件名 | 路径 | 说明 |
|--------|------|------|
| 验证脚本 | `backend/scripts/verify_tenant_ddl.sql` | 完整性验证脚本 |
| 部署脚本 | `backend/scripts/test_tenant_ddl.sh` | 自动化部署脚本 |
| 回滚脚本 | `docker/atlas/migrations/rollback/20251230025000_tenant_tables_rollback.sql` | 回滚DDL |

---

## 🗄️ 表结构说明

### 1. tenants（租户表）

**用途**: 存储租户基本信息

**核心字段**:
- `tenant_id` (VARCHAR(36)) - 租户ID（主键）
- `tenant_name` (VARCHAR(200)) - 租户名称
- `tenant_type` (ENUM) - 租户类型（individual/team/enterprise）
- `status` (ENUM) - 状态（active/suspended/deleted）
- `subscription_tier` (ENUM) - 订阅等级（free/pro/enterprise）
- `contact_email` (VARCHAR(255)) - 联系人邮箱
- `settings` (JSON) - 租户配置（JSON格式）

**索引**:
- PRIMARY KEY: `tenant_id`
- UNIQUE INDEX: `uk_tenant_name` (`tenant_name`, `deleted_at`)
- INDEX: `idx_tenant_type`, `idx_status`, `idx_subscription_tier`, `idx_created_at`

### 2. subscriptions（订阅表）

**用途**: 管理租户订阅信息

**核心字段**:
- `subscription_id` (VARCHAR(36)) - 订阅ID（主键）
- `tenant_id` (VARCHAR(36)) - 租户ID（外键）
- `plan_tier` (ENUM) - 套餐等级
- `billing_cycle` (ENUM) - 计费周期（monthly/quarterly/yearly）
- `start_date` (DATE) - 订阅开始日期
- `end_date` (DATE) - 订阅结束日期
- `quota_bots` (INT) - Bot数量配额
- `quota_messages_per_month` (INT) - 每月消息配额
- `auto_renew` (BOOLEAN) - 自动续费

**索引**:
- PRIMARY KEY: `subscription_id`
- FOREIGN KEY: `tenant_id` → `tenants(tenant_id)`
- UNIQUE INDEX: `uk_tenant_plan` (`tenant_id`, `plan_tier`)
- INDEX: `idx_tenant_id`, `idx_plan_tier`, `idx_status`, `idx_end_date`

### 3. quotas（配额表）

**用途**: 管理租户资源配额

**核心字段**:
- `quota_id` (VARCHAR(36)) - 配额ID（主键）
- `tenant_id` (VARCHAR(36)) - 租户ID（外键）
- `resource_type` (ENUM) - 资源类型（bots/users/api_call/message/token/storage/knowledge_base/workflow）
- `max_limit` (INT) - 最大限制（-1表示无限制）
- `used_count` (INT) - 已使用数量
- `reset_cycle` (ENUM) - 重置周期（daily/weekly/monthly/yearly/never）
- `last_reset_at` (BIGINT) - 上次重置时间

**索引**:
- PRIMARY KEY: `quota_id`
- FOREIGN KEY: `tenant_id` → `tenants(tenant_id)`
- UNIQUE INDEX: `uk_tenant_resource` (`tenant_id`, `resource_type`)
- INDEX: `idx_tenant_resource`

### 4. quota_usage_log（配额使用日志表）

**用途**: 记录配额使用历史

**核心字段**:
- `log_id` (BIGINT) - 日志ID（主键，自增）
- `tenant_id` (VARCHAR(36)) - 租户ID（外键）
- `resource_type` (ENUM) - 资源类型
- `resource_id` (VARCHAR(36)) - 资源ID
- `action` (ENUM) - 操作类型（consume/rollback/reset）
- `amount` (INT) - 数量（正数为消费，负数为回滚）
- `before_count` (INT) - 操作前数量
- `after_count` (INT) - 操作后数量
- `request_id` (VARCHAR(36)) - 请求ID
- `user_id` (VARCHAR(36)) - 操作用户ID
- `reason` (VARCHAR(500)) - 操作原因

**索引**:
- PRIMARY KEY: `log_id` (AUTO_INCREMENT)
- FOREIGN KEY: `tenant_id` → `tenants(tenant_id)`
- INDEX: `idx_tenant_resource`, `idx_created_at`, `idx_user_id`

---

## 🎯 触发器说明

### 1. trg_tenant_status_update

**触发时机**: BEFORE UPDATE ON `tenants`

**功能**:
1. **租户暂停**: 当状态从 `active` → `suspended` 时
   - 自动暂停所有订阅（status → `past_due`）
   - 记录配额使用日志

2. **租户恢复**: 当状态从 `suspended` → `active` 时
   - 恢复未过期的订阅（status → `active`）
   - 只恢复 `end_date >= CURDATE()` 的订阅

### 2. trg_quota_update_timestamp

**触发时机**: BEFORE UPDATE ON `quotas`

**功能**:
- 当 `used_count` 发生变化时，自动更新 `updated_at` 时间戳

---

## 📊 视图说明

### 1. v_tenant_statistics（租户统计视图）

**用途**: 提供租户全面的使用情况统计

**查询字段**:
- 租户基本信息（tenant_id, tenant_name, tenant_type, status）
- 订阅信息（plan_tier, start_date, end_date, auto_renew）
- Bot统计（bot_count, published_bot_count）
- 用户统计（user_count）
- 配额使用（quota_bots_used/limit, quota_messages_used/limit）

**使用示例**:
```sql
SELECT * FROM v_tenant_statistics WHERE tenant_id = 'system-default';
```

### 2. v_quota_usage_rate（配额使用率视图）

**用途**: 计算各资源配额的使用率百分比

**查询字段**:
- 租户信息（tenant_id, tenant_name）
- 配额信息（resource_type, max_limit, used_count）
- 使用率（usage_rate，保留2位小数）
- 重置信息（reset_cycle, last_reset_at）

**使用示例**:
```sql
SELECT * FROM v_quota_usage_rate WHERE usage_rate > 80 ORDER BY usage_rate DESC;
```

### 3. v_subscription_expiry_alert（订阅到期提醒视图）

**用途**: 监控订阅到期情况

**查询字段**:
- 订阅信息（subscription_id, tenant_id, plan_tier）
- 到期信息（end_date, days_until_expiry）
- 告警级别（alert_level: critical/warning/normal）

**告警规则**:
- **critical**: 距离到期 ≤ 7天
- **warning**: 距离到期 ≤ 30天
- **normal**: 其他

**使用示例**:
```sql
SELECT * FROM v_subscription_expiry_alert WHERE alert_level = 'critical';
```

---

## 🚀 部署步骤

### 方式一：自动化部署（推荐）

#### Linux/macOS

```bash
cd backend/scripts

# 设置环境变量
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASSWORD=your_password

# 执行部署脚本
chmod +x test_tenant_ddl.sh
./test_tenant_ddl.sh
```

#### Windows (Git Bash)

```bash
cd backend/scripts

# 设置环境变量
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASSWORD=your_password

# 执行部署脚本
bash test_tenant_ddl.sh
```

### 方式二：手动部署

```bash
# 1. 创建租户表
mysql -h localhost -u root -p opencoze < docker/atlas/migrations/20251230025000_create_tenant_tables.sql

# 2. 创建触发器
mysql -h localhost -u root -p opencoze < docker/atlas/migrations/20251230025001_tenant_status_trigger.sql

# 3. 创建视图
mysql -h localhost -u root -p opencoze < docker/atlas/migrations/20251230025002_tenant_statistics_view.sql

# 4. 初始化数据
mysql -h localhost -u root -p opencoze < docker/atlas/migrations/20251230025003_init_tenant_data.sql

# 5. 验证部署
mysql -h localhost -u root -p opencoze < backend/scripts/verify_tenant_ddl.sql
```

---

## ✅ 验证部署

### 1. 基础验证

```sql
-- 查看所有租户相关表
SHOW TABLES LIKE '%tenant%';
SHOW TABLES LIKE '%subscription%';
SHOW TABLES LIKE '%quota%';
```

**预期输出**:
```
+------------------------------+
| Tables_in_opencoze (%quota%) |
+------------------------------+
| quota_usage_log              |
| quotas                       |
+------------------------------+
+--------------------------------------+
| Tables_in_opencoze (%subscription%) |
+--------------------------------------+
| subscriptions                        |
+--------------------------------------+
+----------------------------------+
| Tables_in_opencoze (%tenant%)   |
+----------------------------------+
| tenants                          |
+----------------------------------+
```

### 2. 数据验证

```sql
-- 查看系统默认租户
SELECT * FROM tenants WHERE tenant_id = 'system-default';

-- 查看系统订阅
SELECT * FROM subscriptions WHERE tenant_id = 'system-default';

-- 查看系统配额
SELECT * FROM quotas WHERE tenant_id = 'system-default';
```

### 3. 触发器验证

```sql
SHOW TRIGGERS LIKE 'trg_%';
```

**预期输出**:
```
+---------------------------+--------+--------+-----------+------------------+----------------------+----------------------+--------------------+--------------+--------+--------+-----------+----------+----------------+---------------+-------------------------+
| Trigger                    | Event  | Table   | Statement | Timing           | Created              | sql_mode             | definer             | character_set_client | collation_connection | database_collation | engine | trigger_name |
+---------------------------+--------+--------+-----------+------------------+----------------------+----------------------+--------------------+--------------+--------+--------+-----------+----------+----------------+---------------+
| trg_tenant_status_update   | UPDATE | tenants | BEGIN ... | BEFORE           | 2025-01-01 00:00:00  | ONLY_FULL_GROUP_BY... | root@localhost      | utf8mb4              | utf8mb4_unicode_ci   | utf8mb4_unicode_ci  | InnoDB  | ...          |
| trg_quota_update_timestamp | UPDATE | quotas  | BEGIN ... | BEFORE           | 2025-01-01 00:00:00  | ONLY_FULL_GROUP_BY... | root@localhost      | utf8mb4              | utf8mb4_unicode_ci   | utf8mb4_unicode_ci  | InnoDB  | ...          |
+---------------------------+--------+--------+-----------+------------------+----------------------+----------------------+--------------------+--------------+--------+--------+-----------+----------+----------------+---------------+
```

### 4. 视图验证

```sql
-- 查看所有视图
SHOW FULL TABLES WHERE TABLE_TYPE LIKE 'VIEW';

-- 测试租户统计视图
SELECT * FROM v_tenant_statistics LIMIT 10;

-- 测试配额使用率视图
SELECT * FROM v_quota_usage_rate WHERE tenant_id = 'system-default';

-- 测试订阅到期提醒视图
SELECT * FROM v_subscription_expiry_alert WHERE alert_level = 'critical';
```

### 5. 运行完整验证脚本

```bash
mysql -h localhost -u root -p opencoze < backend/scripts/verify_tenant_ddl.sql > verification_report.txt
```

检查输出文件 `verification_report.txt` 确保所有检查项都通过。

---

## 🔄 回滚操作

### 警告 ⚠️

**回滚操作将删除以下内容**:
- 所有租户数据
- 所有订阅数据
- 所有配额数据
- 所有使用日志
- 所有视图和触发器

**请确保已备份重要数据！**

### 回滚步骤

```bash
# 执行回滚脚本
mysql -h localhost -u root -p opencoze < docker/atlas/migrations/rollback/20251230025000_tenant_tables_rollback.sql

# 验证回滚结果
SHOW TABLES LIKE '%tenant%';
SHOW TABLES LIKE '%subscription%';
SHOW TABLES LIKE '%quota%';
```

**预期结果**: 不应该有任何表存在

---

## 📝 使用示例

### 创建新租户

```sql
INSERT INTO tenants (
    tenant_id, tenant_name, tenant_type, status,
    subscription_tier, contact_email, created_at, updated_at
) VALUES (
    UUID(),
    'My Company',
    'enterprise',
    'active',
    'pro',
    'admin@mycompany.com',
    UNIX_TIMESTAMP() * 1000,
    UNIX_TIMESTAMP() * 1000
);

-- 获取刚创建的租户ID
SELECT LAST_INSERT_ID() as tenant_id;
```

### 创建订阅

```sql
INSERT INTO subscriptions (
    subscription_id, tenant_id, plan_tier, billing_cycle,
    status, start_date, end_date,
    quota_bots, quota_messages_per_month, quota_storage_gb, quota_team_members,
    auto_renew, created_at, updated_at
) VALUES (
    UUID(),
    'tenant-uuid-here',  -- 替换为实际的租户ID
    'pro',
    'monthly',
    'active',
    CURDATE(),
    DATE_ADD(CURDATE(), INTERVAL 1 MONTH),
    50,
    50000,
    100,
    10,
    TRUE,
    UNIX_TIMESTAMP() * 1000,
    UNIX_TIMESTAMP() * 1000
);
```

### 初始化配额

```sql
INSERT INTO quotas (
    quota_id, tenant_id, resource_type, max_limit, used_count,
    reset_cycle, last_reset_at, created_at, updated_at
) VALUES
(UUID(), 'tenant-uuid-here', 'bots', 50, 0, 'monthly', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
(UUID(), 'tenant-uuid-here', 'users', 10, 0, 'monthly', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
(UUID(), 'tenant-uuid-here', 'message', 50000, 0, 'monthly', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
(UUID(), 'tenant-uuid-here', 'token', 1000000, 0, 'monthly', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
(UUID(), 'tenant-uuid-here', 'storage', 100, 0, 'monthly', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
(UUID(), 'tenant-uuid-here', 'knowledge_base', 20, 0, 'monthly', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
(UUID(), 'tenant-uuid-here', 'workflow', 50, 0, 'monthly', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000);
```

### 使用配额（记录日志）

```sql
-- 消费配额
INSERT INTO quota_usage_log (
    tenant_id, resource_type, resource_id, action, amount,
    before_count, after_count, user_id, reason, created_at
) VALUES (
    'tenant-uuid-here',
    'bots',
    'bot-uuid-here',
    'consume',
    1,
    0,
    1,
    'user-uuid-here',
    'Create new bot',
    UNIX_TIMESTAMP() * 1000
);

-- 更新配额使用数量
UPDATE quotas
SET used_count = used_count + 1,
    updated_at = UNIX_TIMESTAMP() * 1000
WHERE tenant_id = 'tenant-uuid-here'
  AND resource_type = 'bots';
```

### 查询租户统计

```sql
-- 查看租户统计
SELECT * FROM v_tenant_statistics WHERE tenant_id = 'tenant-uuid-here';

-- 查看配额使用率
SELECT * FROM v_quota_usage_rate
WHERE tenant_id = 'tenant-uuid-here'
  AND usage_rate > 80
ORDER BY usage_rate DESC;

-- 查看到期提醒
SELECT * FROM v_subscription_expiry_alert
WHERE tenant_id = 'tenant-uuid-here'
  AND alert_level IN ('critical', 'warning');
```

---

## 🔧 故障排查

### 问题1: 外键约束错误

**错误信息**:
```
ERROR 1822 (HY000): Cannot add foreign key constraint
```

**原因**: `tenants` 表不存在或被删除

**解决方案**:
```sql
-- 检查tenants表是否存在
SHOW TABLES LIKE 'tenants';

-- 如果不存在，重新创建
SOURCE docker/atlas/migrations/20251230025000_create_tenant_tables.sql;
```

### 问题2: 触发器创建失败

**错误信息**:
```
ERROR 1359 (HY000): Trigger already exists
```

**原因**: 触发器已存在

**解决方案**:
```sql
-- 删除已存在的触发器
DROP TRIGGER IF EXISTS trg_tenant_status_update;
DROP TRIGGER IF EXISTS trg_quota_update_timestamp;

-- 重新创建
SOURCE docker/atlas/migrations/20251230025001_tenant_status_trigger.sql;
```

### 问题3: 视图查询失败

**错误信息**:
```
ERROR 1356 (HY000): View's select contains invalid reference
```

**原因**: 基础表结构变更

**解决方案**:
```sql
-- 删除视图
DROP VIEW IF EXISTS v_tenant_statistics;
DROP VIEW IF EXISTS v_quota_usage_rate;
DROP VIEW IF EXISTS v_subscription_expiry_alert;

-- 重新创建
SOURCE docker/atlas/migrations/20251230025002_tenant_statistics_view.sql;
```

### 问题4: 系统租户已存在

**错误信息**:
```
ERROR 1062 (23000): Duplicate entry 'system-default' for key 'PRIMARY'
```

**原因**: 系统租户已经存在

**解决方案**:
```sql
-- 使用ON DUPLICATE KEY UPDATE（脚本已包含）
-- 或者先删除再插入
DELETE FROM tenants WHERE tenant_id = 'system-default';
DELETE FROM subscriptions WHERE tenant_id = 'system-default';
DELETE FROM quotas WHERE tenant_id = 'system-default';

-- 重新执行初始化脚本
SOURCE docker/atlas/migrations/20251230025003_init_tenant_data.sql;
```

---

## 📚 附录

### A. 字段命名规范

| 类型 | 格式 | 示例 |
|------|------|------|
| 主键 | `{table}_id` | `tenant_id`, `subscription_id` |
| 外键 | `{referenced_table}_id` | `tenant_id` |
| 布尔字段 | `is_{property}` | `is_active`, `is_public` |
| 时间戳 | `{action}_at` | `created_at`, `updated_at`, `deleted_at` |
| 枚举类型 | 使用ENUM | `status`, `tenant_type` |

### B. 索引设计原则

1. **所有外键列**必须有索引
2. **常用查询条件**需要索引（status, tenant_type）
3. **唯一约束**使用UNIQUE INDEX
4. **复合索引**遵循最左前缀原则

### C. 数据类型选择

| 数据类型 | 说明 | 示例 |
|----------|------|------|
| VARCHAR(36) | UUID存储 | `tenant_id` |
| ENUM | 固定选项 | `status` ENUM('active', 'suspended', 'deleted') |
| BIGINT | 时间戳（毫秒） | `created_at` |
| JSON | 配置信息 | `settings` |
| DATE | 日期 | `start_date` |
| BOOLEAN | 布尔值 | `auto_renew` |

### D. 性能优化建议

1. **索引优化**:
   - 为高频查询字段添加索引
   - 定期分析索引使用情况
   - 删除未使用的索引

2. **查询优化**:
   - 避免SELECT *，只查询需要的字段
   - 使用视图简化复杂查询
   - 使用EXPLAIN分析查询计划

3. **数据归档**:
   - 定期归档历史日志数据（quota_usage_log）
   - 考虑使用分区表（按时间分区）

---

## 📞 联系方式

如有问题，请联系：
- **开发团队**: ZKER Enterprise Team
- **文档版本**: v1.0.0
- **最后更新**: 2025-01-01

---

**文档结束** ✅
