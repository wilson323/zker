# 数据库迁移应用指南

> 最后更新: 2025-12-31
> 状态: 待应用迁移

---

## 📋 待应用迁移列表

### 1. Tenant联系字段扩展

**文件**: `docker/atlas/migrations/20251231121812_update_tenant_contact_fields.sql`

**变更**: 扩展`tenants.contact_phone`字段长度 VARCHAR(20) → VARCHAR(32)

**原因**: 匹配实体定义 `backend/domain/tenant/entity/tenant.go`

---

## 🚀 应用迁移方法

### 方法1: 使用Atlas CLI（推荐）

```bash
# 1. 安装Atlas（如果未安装）
curl -sSf https://atlasgo.sh | sh

# 2. 检查迁移状态
atlas migrate diff \
  --dir "file://docker/atlas/migrations" \
  --url "mysql://root:root_password_2025@localhost:3306/opencoze"

# 3. 应用迁移
atlas migrate apply \
  --dir "file://docker/atlas/migrations" \
  --url "mysql://root:root_password_2025@localhost:3306/opencoze" \
  --baseline 20251231121812
```

### 方法2: 使用Docker执行SQL

```bash
# 连接到MySQL容器
docker exec -it coze-mysql mysql -uroot -p

# 在MySQL命令行中执行
USE opencoze;
SOURCE /docker-entrypoint-initdb.d/20251231121812_update_tenant_contact_fields.sql;

# 或者直接执行ALTER语句
ALTER TABLE tenants MODIFY COLUMN contact_phone VARCHAR(32) COMMENT '联系人电话';
```

### 方法3: 使用MySQL客户端

```bash
mysql -h localhost -P 3306 -u root -p opencoze < docker/atlas/migrations/20251231121812_update_tenant_contact_fields.sql
```

---

## ✅ 验证迁移

```bash
# 检查表结构
docker exec coze-mysql mysql -uroot -p -e "
USE opencoze;
DESCRIBE tenants;
SHOW CREATE TABLE tenants\G
"

# 验证字段长度
docker exec coze-mysql mysql -uroot -p -e "
USE opencoze;
SELECT COLUMN_NAME, COLUMN_TYPE, CHARACTER_MAXIMUM_LENGTH
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_NAME = 'tenants' AND COLUMN_NAME = 'contact_phone';
"
```

**期望输出**:
```
+---------------+-------------+------------------------+
| COLUMN_NAME   | COLUMN_TYPE | CHARACTER_MAXIMUM_LENGTH |
+---------------+-------------+------------------------+
| contact_phone | varchar(32) |                     32 |
+---------------+-------------+------------------------+
```

---

## 🔄 回滚迁移

如果需要回滚：

```bash
# 方法1: Atlas回滚
atlas migrate apply \
  --dir "file://docker/atlas/migrations" \
  --url "mysql://root:root_password_2025@localhost:3306/opencoze" \
  --down 20251231121812

# 方法2: 手动回滚SQL
ALTER TABLE tenants MODIFY COLUMN contact_phone VARCHAR(20) COMMENT '联系人电话';
```

---

## 📝 注意事项

1. **备份数据**: 应用迁移前先备份数据库
   ```bash
   docker exec coze-mysql mysqldump -uroot -p opencoze > backup_before_migration.sql
   ```

2. **检查依赖**: 确认没有代码依赖旧字段长度

3. **测试环境**: 先在测试环境验证迁移

4. **生产环境**: 建议在低峰期执行

---

## 🔍 相关文件

- Entity定义: `backend/domain/tenant/entity/tenant.go`
- 迁移文件: `docker/atlas/migrations/20251231121812_update_tenant_contact_fields.sql`
- Atlas配置: `docker/atlas/atlas.hcl`

---

**维护者**: Backend Team
**文档版本**: v1.0
