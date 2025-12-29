# ZKER 数据库迁移脚本

本目录包含 Coze Studio/ZKER 项目的数据库迁移脚本。

## 📋 可用迁移脚本

| 脚本文件 | 功能描述 | 预计执行时间 | 风险等级 |
|---------|---------|-------------|---------|
| `add_tenant_id_to_tables.sql` | 为所有业务表添加 tenant_id 字段以支持多租户 | 5-10 分钟 | 🟡 中等 |

## 🚀 执行迁移步骤

### 前置准备

1. **备份数据库**（必须！）
   ```bash
   # 使用 mysqldump 备份
   mysqldump -u root -p coze_studio > backup_before_migration_$(date +%Y%m%d_%H%M%S).sql

   # 或使用 Docker 容器备份
   docker exec mysql-container mysqldump -u root -p coze_studio > backup.sql
   ```

2. **检查磁盘空间**
   ```bash
   # 确保至少有 2倍数据库大小的可用空间
   df -h
   ```

3. **停止应用写入**（可选但推荐）
   - 将应用设置为只读模式
   - 或停止应用服务

### 执行迁移

#### 方式一：直接使用 MySQL 客户端

```bash
# 连接到数据库
mysql -u root -p coze_studio

# 执行迁移脚本
source /path/to/migrations/add_tenant_id_to_tables.sql

# 或使用命令行
mysql -u root -p coze_studio < /path/to/migrations/add_tenant_id_to_tables.sql
```

#### 方式二：使用 Docker 容器

```bash
# 复制脚本到容器
docker cp /path/to/migrations/add_tenant_id_to_tables.sql mysql-container:/tmp/

# 执行迁移
docker exec -i mysql-container mysql -u root -p coze_studio < /tmp/add_tenant_id_to_tables.sql
```

#### 方式三：使用 Go 迁移工具（推荐）

```bash
cd backend/migrations
go run migrate.go execute --script=add_tenant_id_to_tables.sql
```

### 验证迁移

迁移脚本包含内置验证查询，执行完成后会自动运行验证。你也可以手动运行：

```sql
-- 1. 检查字段是否添加成功
SELECT TABLE_NAME, COLUMN_NAME, COLUMN_TYPE
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'tenant_id';

-- 2. 检查索引是否创建成功
SELECT TABLE_NAME, INDEX_NAME
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND INDEX_NAME LIKE 'idx_tenant%';

-- 3. 检查数据完整性
SELECT TABLE_NAME, COUNT(*) AS total_rows,
       SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) AS null_count
FROM (
    SELECT 'bots' AS TABLE_NAME, * FROM bots
    UNION ALL
    SELECT 'conversations', * FROM conversations
    -- ... 其他表
) AS all_tables
GROUP BY TABLE_NAME;
```

### 应用层代码更新

迁移完成后，需要更新应用代码以使用 tenant_id：

1. **更新查询逻辑**
   ```go
   // ❌ 旧代码：没有 tenant_id 过滤
   db.Where("bot_id = ?", botID).First(&bot)

   // ✅ 新代码：添加 tenant_id 过滤
   tenantID := middleware.GetTenantIDFromContext(ctx)
   db.Where("tenant_id = ? AND bot_id = ?", tenantID, botID).First(&bot)
   ```

2. **更新插入逻辑**
   ```go
   // ❌ 旧代码：没有设置 tenant_id
   bot := &Bot{BotID: botID, Name: name}

   // ✅ 新代码：自动设置 tenant_id
   tenantID := middleware.GetTenantIDFromContext(ctx)
   bot := &Bot{TenantID: tenantID, BotID: botID, Name: name}
   ```

3. **启用租户隔离中间件验证**
   - 确保应用启动时调用了 `middleware.InitTenantMiddleware()`
   - 参见 `backend/api/middleware/tenant_isolation.go`

## 🔄 回滚方案

如果迁移出现问题，可以使用回滚脚本：

```sql
-- 方式一：手动回滚（见迁移脚本第五部分）
-- 执行 add_tenant_id_to_tables.sql 中的回滚部分

-- 方式二：从备份恢复
mysql -u root -p coze_studio < backup_before_migration_YYYYMMDD_HHMMSS.sql
```

## ⚠️ 注意事项

### 执行前

- ✅ **必须备份数据库**
- ✅ 在测试环境先执行一遍
- ✅ 确认有足够的磁盘空间
- ✅ 通知相关人员迁移计划
- ✅ 准备回滚方案

### 执行中

- ⏱️ 大表迁移可能需要较长时间，请耐心等待
- 📊 监控数据库性能指标（CPU、内存、磁盘IO）
- 📝 记录执行日志

### 执行后

- ✅ 运行验证脚本确认迁移成功
- ✅ 测试应用基本功能
- ✅ 监控错误日志
- ✅ 保留备份至少 7 天

## 📊 影响范围

### 数据库变更

| 表名 | 变更类型 | 预计影响 |
|-----|---------|---------|
| bots | 添加列+索引 | 中等 |
| conversations | 添加列+索引 | 中等 |
| knowledge_bases | 添加列+索引 | 小 |
| workflows | 添加列+索引 | 小 |
| agent_draft | 添加列+索引 | 小 |
| bot_messages | 添加列+索引 | 大（可能较慢） |
| bot_tools | 添加列+索引 | 小 |
| bot_components | 添加列+索引 | 小 |

### 应用层变更

- ✅ 租户隔离中间件已启用验证
- ✅ 需要更新查询逻辑（添加 WHERE tenant_id = ?）
- ✅ 需要更新插入逻辑（设置 tenant_id 字段）

## 🆘 故障排查

### 问题1：迁移脚本执行超时

**症状**：脚本执行时间过长（> 30分钟）

**解决方案**：
```sql
-- 分批执行，先执行小表
ALTER TABLE agent_draft ADD COLUMN tenant_id ...;
ALTER TABLE bot_tools ADD COLUMN tenant_id ...;

-- 再执行大表
ALTER TABLE bot_messages ADD COLUMN tenant_id ...;
```

### 问题2：索引创建失败

**症状**：`ERROR 1061 (42000): Duplicate key name`

**解决方案**：
```sql
-- 检查索引是否已存在
SHOW INDEX FROM bots WHERE Key_name = 'idx_tenant_id';

-- 如果存在，先删除
ALTER TABLE bots DROP INDEX idx_tenant_id;

-- 重新创建
ALTER TABLE bots ADD INDEX idx_tenant_id (tenant_id);
```

### 问题3：磁盘空间不足

**症状**：`ERROR 1114 (HY000): The table is full`

**解决方案**：
```bash
# 1. 清理临时文件
# 2. 扩容磁盘
# 3. 或使用在线DDL（ALGORITHM=INPLACE）
```

## 📞 联系方式

如有问题，请联系：
- 数据库团队：db-team@coze-studio.com
- DevOps团队：devops@coze-studio.com
- 技术支持：support@coze-studio.com

## 📝 更新日志

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2025-01-01 | v1.0.0 | 初始版本，添加 tenant_id 迁移脚本 | Claude AI |
