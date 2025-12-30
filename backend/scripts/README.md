# 租户表DDL快速部署指南

## 📋 概述

本指南提供了租户表DDL脚本的快速部署方法，适用于开发、测试和生产环境。

## 🚀 快速开始

### 前置条件

- MySQL 8.4.5+
- Go 1.24+ (用于运行集成测试)
- Bash 环境（Linux/macOS/Git Bash）

### 方式一：自动化部署（推荐）

```bash
cd backend/scripts

# 设置环境变量
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASSWORD=your_password

# 执行部署
chmod +x test_tenant_ddl.sh
./test_tenant_ddl.sh
```

### 方式二：手动执行

```bash
# 1. 创建表
mysql -u root -p opencoze < docker/atlas/migrations/20251230025000_create_tenant_tables.sql

# 2. 创建触发器
mysql -u root -p opencoze < docker/atlas/migrations/20251230025001_tenant_status_trigger.sql

# 3. 创建视图
mysql -u root -p opencoze < docker/atlas/migrations/20251230025002_tenant_statistics_view.sql

# 4. 初始化数据
mysql -u root -p opencoze < docker/atlas/migrations/20251230025003_init_tenant_data.sql

# 5. 验证
mysql -u root -p opencoze < backend/scripts/verify_tenant_ddl.sql
```

### 方式三：Go集成测试

```bash
cd backend/scripts

# 设置环境变量
export MYSQL_PASSWORD=your_password

# 运行集成测试
go run test_tenant_integration.go \
  -host=localhost \
  -port=3306 \
  -user=root \
  -password=$MYSQL_PASSWORD \
  -database=opencoze
```

## ✅ 验证部署

### SQL验证

```sql
-- 查看所有表
SHOW TABLES LIKE '%tenant%';
SHOW TABLES LIKE '%subscription%';
SHOW TABLES LIKE '%quota%';

-- 查看系统租户
SELECT * FROM tenants WHERE tenant_id = 'system-default';

-- 查看统计视图
SELECT * FROM v_tenant_statistics LIMIT 10;
```

### 完整验证脚本

```bash
mysql -u root -p opencoze < backend/scripts/verify_tenant_ddl.sql > verification_report.txt
cat verification_report.txt
```

## 🔄 回滚操作

```bash
# 执行回滚
mysql -u root -p opencoze < docker/atlas/migrations/rollback/20251230025000_tenant_tables_rollback.sql

# 验证回滚
SHOW TABLES LIKE '%tenant%';
```

## 📁 文件结构

```
backend/scripts/
├── test_tenant_ddl.sh              # 自动化部署脚本
├── test_tenant_integration.go      # Go集成测试
├── verify_tenant_ddl.sql           # SQL验证脚本
└── README.md                       # 本文件

docker/atlas/migrations/
├── 20251230025000_create_tenant_tables.sql         # 创建表
├── 20251230025001_tenant_status_trigger.sql        # 创建触发器
├── 20251230025002_tenant_statistics_view.sql       # 创建视图
├── 20251230025003_init_tenant_data.sql             # 初始化数据
└── rollback/
    └── 20251230025000_tenant_tables_rollback.sql   # 回滚脚本
```

## 🔧 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| MYSQL_HOST | MySQL主机 | localhost |
| MYSQL_PORT | MySQL端口 | 3306 |
| MYSQL_USER | MySQL用户 | root |
| MYSQL_PASSWORD | MySQL密码 | - |
| MYSQL_DATABASE | 数据库名 | opencoze |

## 📊 部署检查清单

- [ ] MySQL服务运行正常
- [ ] 数据库 `opencoze` 已创建
- [ ] 4个核心表已创建（tenants、subscriptions、quotas、quota_usage_log）
- [ ] 2个触发器已创建（trg_tenant_status_update、trg_quota_update_timestamp）
- [ ] 3个视图已创建（v_tenant_statistics、v_quota_usage_rate、v_subscription_expiry_alert）
- [ ] 系统默认租户已初始化
- [ ] 配额数据已初始化
- [ ] 所有验证测试通过

## 🐛 常见问题

### 1. 外键约束错误

```bash
# 检查tenants表是否存在
mysql -u root -p -e "USE opencoze; SHOW TABLES LIKE 'tenants';"

# 如果不存在，重新执行表创建脚本
mysql -u root -p opencoze < docker/atlas/migrations/20251230025000_create_tenant_tables.sql
```

### 2. 触发器已存在

```sql
-- 删除已存在的触发器
DROP TRIGGER IF EXISTS trg_tenant_status_update;
DROP TRIGGER IF EXISTS trg_quota_update_timestamp;

-- 重新创建
SOURCE docker/atlas/migrations/20251230025001_tenant_status_trigger.sql;
```

### 3. 系统租户已存在

```sql
-- 使用ON DUPLICATE KEY UPDATE（脚本已包含）
-- 或先删除再插入
DELETE FROM tenants WHERE tenant_id = 'system-default';
DELETE FROM subscriptions WHERE tenant_id = 'system-default';
DELETE FROM quotas WHERE tenant_id = 'system-default';

-- 重新执行初始化脚本
SOURCE docker/atlas/migrations/20251230025003_init_tenant_data.sql;
```

## 📚 更多文档

- [完整部署文档](../../docs/企业级功能完善与统一性设计方案/ZKER-租户表DDL部署文档_v1.0.md)
- [企业级开发规范手册](../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [数据库设计完整交付清单](../../docs/企业级功能完善与统一性设计方案/数据库设计完整交付清单.md)

## 📞 支持

如有问题，请联系开发团队或查看相关文档。

---

**最后更新**: 2025-01-01
**版本**: v1.0.0
