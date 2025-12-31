# ZKER 数据迁移 FAQ 与故障排查指南

> **版本**: v1.0
> **更新日期**: 2025-01-03
> **维护团队**: 数据库架构组

---

## 📋 目录

1. [常见问题FAQ](#常见问题faq)
2. [故障排查指南](#故障排查指南)
3. [应急处理流程](#应急处理流程)
4. [最佳实践](#最佳实践)
5. [工具参考](#工具参考)

---

## 常见问题FAQ

### Q1: 迁移前如何评估时间？

**问题**: 如何准确估算迁移时间？

**回答**:

```bash
# 使用评估工具
./migration-tool estimate \
  --tables=bots,conversations,knowledge \
  --sample-size=10000

# 输出示例：
# 表 bots:
#   总记录数: 1,000,000
#   样本迁移速率: 5000 条/秒
#   预估时间: 200 秒 (3.3 分钟)
#
# 表 conversations:
#   总记录数: 5,000,000
#   样本迁移速率: 3000 条/秒
#   预估时间: 1667 秒 (27.8 分钟)
#
# 总预估时间: 31.1 分钟
```

**影响因素**:
1. 数据量大小
2. 服务器性能（CPU、内存、磁盘I/O）
3. 网络带宽
4. 批处理大小
5. 并发度

### Q2: 双写模式和直接模式如何选择？

**问题**: 两种迁移模式的区别和使用场景？

**回答**:

| 对比项 | 双写模式 | 直接模式 |
|-------|---------|---------|
| **停机时间** | 零停机 | 需要停机 |
| **实现复杂度** | 高 | 低 |
| **迁移周期** | 长（1-2周） | 短（数小时） |
| **数据安全性** | 高 | 中 |
| **适用规模** | 大规模（≥1000万） | 小规模（<100万） |
| **回滚难度** | 低 | 中 |

**推荐选择**:
- 生产环境、大规模数据：双写模式
- 测试环境、小规模数据：直接模式

### Q3: 如何避免锁表？

**问题**: 迁移期间避免长时间锁表影响业务？

**回答**:

**策略1: 使用pt-online-schema-change**

```bash
pt-online-schema-change \
  --alter "ADD COLUMN tenant_id VARCHAR(36) NULL" \
  --charset=utf8mb4 \
  --critical-load="Threads_running=50" \
  --max-load="Threads_running=100" \
  --chunk-size=1000 \
  --progress=time,30 \
  --execute \
  D=coze_studio,t=bots
```

**策略2: 使用gh-ost**

```bash
gh-ost \
  --max-load=Threads_running=100 \
  --critical-load=Threads_running=1000 \
  --chunk-size=1000 \
  --throttle-control-replicas="..." \
  --database=coze_studio \
  --table=bots \
  --alter="ADD COLUMN tenant_id VARCHAR(36) NULL" \
  --allow-on-master \
  --execute
```

**策略3: 分批小事务**

```sql
-- 不推荐（锁表时间长）
UPDATE bots SET tenant_id = 'xxx';

-- 推荐（分批小事务）
UPDATE bots
SET tenant_id = 'xxx'
WHERE id > 0 AND id <= 10000;

UPDATE bots
SET tenant_id = 'xxx'
WHERE id > 10000 AND id <= 20000;
```

### Q4: 如何处理外键约束？

**问题**: 迁移表时遇到外键约束怎么办？

**回答**:

**方案1: 临时禁用外键检查**

```sql
-- 迁移前
SET FOREIGN_KEY_CHECKS = 0;

-- 执行迁移
ALTER TABLE bots ADD COLUMN tenant_id VARCHAR(36);
UPDATE bots SET tenant_id = 'xxx';

-- 迁移后
SET FOREIGN_KEY_CHECKS = 1;
```

**方案2: 按依赖顺序迁移**

```bash
# 1. 先迁移父表
./migration-tool migrate --tables=tenants,users

# 2. 再迁移子表
./migration-tool migrate --tables=bots,conversations

# 3. 最后迁移关联表
./migration-tool migrate --tables=bot_tags,conversation_tags
```

### Q5: 内存不足怎么办？

**问题**: 迁移大表时遇到OOM（Out of Memory）？

**回答**:

**临时方案**:

```bash
# 1. 减小批处理大小
./migration-tool migrate --batch-size=500  # 从1000降到500

# 2. 减少并发度
./migration-tool migrate --workers=5  # 从10降到5

# 3. 清理缓存
mysql -u root -p -e "FLUSH TABLES;"
```

**永久方案**:

```sql
-- 调整MySQL缓冲池大小
SET GLOBAL innodb_buffer_pool_size = 4294967296; -- 4GB

-- 调整排序缓冲区
SET GLOBAL sort_buffer_size = 67108864; -- 64MB

-- 调整连接缓冲区
SET GLOBAL read_buffer_size = 2097152; -- 2MB
SET GLOBAL read_rnd_buffer_size = 8388608; -- 8MB
```

### Q6: 如何验证数据完整性？

**问题**: 迁移后如何确保数据没有丢失或损坏？

**回答**:

**方法1: 记录数对比**

```sql
-- 迁移前记录数
SELECT COUNT(*) FROM bots;

-- 迁移后记录数
SELECT COUNT(*) FROM bots;

-- 应该完全一致
```

**方法2: 哈希对比**

```bash
# 迁移前计算哈希
./data_validator.go --tables=bots --output=before.json

# 迁移后计算哈希
./data_validator.go --tables=bots --output=after.json

# 对比
diff before.json after.json
```

**方法3: 采样对比**

```sql
-- 随机抽取1000条记录对比
SELECT * FROM bots TABLESAMPLE SYSTEM(0.1) LIMIT 1000;
```

### Q7: 迁移过程中断怎么办？

**问题**: 迁移中断了，如何恢复？

**回答**:

```bash
# 检查点机制会自动保存进度
# 重新运行相同命令即可恢复
./migration-tool migrate --tables=bots

# 输出：
# 检测到检查点: last_id=500000
# 从检查点恢复...
```

**手动清理检查点**:

```sql
-- 如果需要从头开始
DELETE FROM migration_checkpoints
WHERE migration_id = 'migration_xxx'
  AND table_name = 'bots';
```

### Q8: 性能下降多少是正常的？

**问题**: 迁移期间性能会下降多少？

**回答**:

**预期影响**:
- 在线迁移（双写模式）: 性能下降 5-15%
- DDL操作（添加字段）: 性能下降 50-100%（数分钟）
- 数据回填: 性能下降 10-30%

**监控指标**:
```bash
# 实时监控QPS
./monitor --web-port=8080

# 关注：
# - QPS是否低于基线的80%
# - P95延迟是否增加50%以上
# - 错误率是否超过1%
```

### Q9: 如何灰度发布？

**问题**: 如何安全地切换到新表结构？

**回答**:

**灰度策略**:

```bash
# 第1天：10%流量
./migration-tool switch --percent=10 --duration=24h

# 第2天：50%流量
./migration-tool switch --percent=50 --duration=24h

# 第3天：100%流量
./migration-tool switch --percent=100 --duration=24h

# 第4天：清理旧逻辑
./migration-tool cleanup --after-days=7
```

**回滚**:
```bash
# 如果发现问题，立即回滚
./migration-tool rollback-switch --percent=100
```

### Q10: binlog空间不足怎么办？

**问题**: 迁移期间binlog增长过快占满磁盘？

**回答**:

**临时方案**:

```bash
# 清理旧binlog（谨慎！）
mysql -u root -p -e "
  PURGE BINARY LOGS BEFORE DATE_SUB(NOW(), INTERVAL 1 DAY);
"
```

**永久方案**:

```sql
-- 设置binlog过期时间
SET GLOBAL binlog_expire_logs_seconds = 86400; -- 1天

-- 或者按文件数量保留
SET GLOBAL max_binlog_files = 10;
```

---

## 故障排查指南

### 故障1: 迁移卡住不动

**症状**: 进度长时间（>30分钟）无更新

**排查步骤**:

```bash
# 1. 检查MySQL进程列表
mysql -u root -p -e "SHOW PROCESSLIST;"

# 2. 检查锁等待
mysql -u root -p -e "
  SELECT *
  FROM performance_schema.events_waits_current
  WHERE EVENT_NAME LIKE '%lock%';
"

# 3. 检查长事务
mysql -u root -p -e "
  SELECT *
  FROM information_schema.innodb_trx
  WHERE TIME_TO_SEC(TIMEDIFF(NOW(), trx_started)) > 60;
"

# 4. 检查死锁
mysql -u root -p -e "SHOW ENGINE INNODB STATUS\G" | grep -A 20 "LATEST DETECTED DEADLOCK"
```

**解决方案**:

```sql
-- 如果有锁等待，杀死会话
KILL <process_id>;

-- 如果有死锁，回滚事务
ROLLBACK;

-- 如果DDL操作卡住，可能是表太大，考虑使用pt-online-schema-change
```

### 故障2: 数据不一致

**症状**: 验证时发现记录数不匹配

**排查步骤**:

```bash
# 1. 对比迁移前后记录数
mysql -u root -p -e "
  SELECT
    TABLE_NAME,
    TABLE_ROWS
  FROM information_schema.TABLES
  WHERE TABLE_SCHEMA = 'coze_studio'
    AND TABLE_NAME IN ('bots', 'conversations');
"

# 2. 检查是否有NULL值
mysql -u root -p -e "
  SELECT
    TABLE_NAME,
    COUNT(*) as total,
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) as null_count
  FROM information_schema.COLUMNS c
  JOIN information_schema.TABLES t ON c.TABLE_NAME = t.TABLE_NAME
  WHERE c.COLUMN_NAME = 'tenant_id'
    AND t.TABLE_SCHEMA = 'coze_studio'
  GROUP BY TABLE_NAME;
"

# 3. 检查迁移日志
grep -i "error\|fail" /var/log/migration.log
```

**解决方案**:

```sql
-- 修复NULL值
UPDATE bots
SET tenant_id = (
  SELECT u.tenant_id
  FROM users u
  WHERE u.user_id = bots.creator_id
  LIMIT 1
)
WHERE tenant_id IS NULL OR tenant_id = '';

-- 如果没有关联用户，使用默认租户
UPDATE bots
SET tenant_id = 'system_tenant'
WHERE tenant_id IS NULL OR tenant_id = '';
```

### 故障3: 性能严重下降

**症状**: QPS下降 > 50%，响应时间增加 > 200%

**排查步骤**:

```bash
# 1. 检查慢查询
mysql -u root -p -e "
  SELECT *
  FROM mysql.slow_log
  WHERE start_time > DATE_SUB(NOW(), INTERVAL 10 MINUTE)
  ORDER BY query_time DESC
  LIMIT 10;
"

# 2. 查看执行计划
mysql -u root -p -e "
  EXPLAIN SELECT *
  FROM bots
  WHERE tenant_id = 'xxx'
  ORDER BY created_at DESC
  LIMIT 20;
"

# 3. 检查索引使用情况
mysql -u root -p -e "
  SELECT
    TABLE_NAME,
    INDEX_NAME,
    SEQ_IN_INDEX,
    COLUMN_NAME
  FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = 'coze_studio'
    AND TABLE_NAME = 'bots'
    AND INDEX_NAME LIKE '%tenant%';
"
```

**解决方案**:

```sql
-- 如果缺少索引，添加索引
CREATE INDEX idx_tenant_created ON bots(tenant_id, created_at);

-- 如果索引失效，重建索引
ALTER TABLE bots DROP INDEX idx_tenant_id;
CREATE INDEX idx_tenant_id ON bots(tenant_id);

-- 如果统计信息过期，更新统计信息
ANALYZE TABLE bots;
```

### 故障4: 磁盘空间不足

**症状**: 磁盘使用率 > 90%

**排查步骤**:

```bash
# 检查磁盘使用
df -h /var/lib/mysql

# 检查binlog大小
ls -lh /var/lib/mysql/mysql-bin.*

# 检查临时文件
ls -lh /tmp/
```

**解决方案**:

```bash
# 1. 清理binlog
mysql -u root -p -e "
  PURGE BINARY LOGS BEFORE DATE_SUB(NOW(), INTERVAL 1 DAY);
"

# 2. 清理临时文件
rm -f /tmp/migration_*.tmp

# 3. 清理旧备份
rm -f /backup/migration/backup_20250101.sql

# 4. 如果仍然不足，暂停迁移
./migration-tool pause --migration-id=migration_xxx
```

### 故障5: 连接数耗尽

**症状**: "Too many connections" 错误

**排查步骤**:

```bash
# 检查当前连接数
mysql -u root -p -e "SHOW STATUS LIKE 'Threads_connected';"

# 检查最大连接数
mysql -u root -p -e "SHOW VARIABLES LIKE 'max_connections';"

# 检查连接详情
mysql -u root -p -e "SHOW PROCESSLIST;"
```

**解决方案**:

```sql
-- 临时增加最大连接数
SET GLOBAL max_connections = 500;

-- 查看空闲连接
KILL <idle_process_id>;

-- 优化连接池配置（应用侧）
-- spring.datasource.hikari.maximum-pool-size=20
```

---

## 应急处理流程

### 场景1: 迁移失败需要回滚

```bash
# 1. 立即停止迁移
./migration-tool abort --migration-id=migration_xxx

# 2. 停止应用服务
sudo systemctl stop coze-studio

# 3. 恢复数据库备份
gunzip -c /backup/migration/coze_studio_full_YYYYMMDD.sql.gz | \
  mysql -u root -p coze_studio

# 4. 验证数据
mysql -u root -p -e "USE coze_studio; SELECT COUNT(*) FROM bots;"

# 5. 重启应用服务
sudo systemctl start coze-studio

# 6. 验证服务
curl http://localhost:8888/health
```

### 场景2: 性能严重下降影响业务

```bash
# 1. 暂停迁移（保留进度）
./migration-tool pause --migration-id=migration_xxx

# 2. 调整迁移参数
./migration-tool resume \
  --migration-id=migration_xxx \
  --batch-size=500 \
  --workers=3

# 3. 监控性能恢复情况
./monitor --web-port=8080
```

### 场景3: 数据损坏

```bash
# 1. 停止迁移
./migration-tool abort --migration-id=migration_xxx

# 2. 使用binlog恢复
# 查找损坏前的binlog位置
mysqlbinlog --start-datetime="2025-01-03 02:00:00" \
  --stop-datetime="2025-01-03 06:00:00" \
  /var/lib/mysql/mysql-bin.000001 | \
  mysql -u root -p coze_studio

# 3. 验证数据完整性
./data_validator.go --tables=bots --verbose
```

---

## 最佳实践

### 1. 迁移前

- ✅ **强制备份**: 全量备份 + binlog
- ✅ **测试环境验证**: 在测试环境完整演练
- ✅ **性能基线**: 记录迁移前的性能指标
- ✅ **容量评估**: 确保磁盘、内存充足
- ✅ **维护窗口**: 选择业务低峰期
- ✅ **通知相关方**: 提前通知开发、测试、产品

### 2. 迁移中

- ✅ **小步快跑**: 分批次迁移，每次验证
- ✅ **实时监控**: 监控进度、性能、错误
- ✅ **保留日志**: 详细记录所有操作
- ✅ **及时沟通**: 定期同步进度和问题
- ✅ **准备回滚**: 随时可以回滚

### 3. 迁移后

- ✅ **全面验证**: 数据完整性 + 业务功能 + 性能
- ✅ **持续监控**: 至少监控7天
- ✅ **清理资源**: 删除临时文件、检查点
- ✅ **文档记录**: 更新架构文档、操作手册
- ✅ **复盘总结**: 记录经验教训

---

## 工具参考

### 迁移工具

| 工具 | 功能 | 使用场景 |
|------|------|---------|
| `migration-tool` | 主迁移工具 | 执行迁移、回滚 |
| `data_validator` | 数据验证工具 | 验证数据完整性 |
| `performance_benchmark` | 性能测试工具 | 对比迁移前后性能 |
| `monitor` | 监控面板 | 实时监控迁移进度 |

### MySQL工具

| 工具 | 功能 | 文档 |
|------|------|------|
| `mysqldump` | 逻辑备份 | [文档](https://dev.mysql.com/doc/refman/8.0/en/mysqldump.html) |
| `mysqlbinlog` | binlog恢复 | [文档](https://dev.mysql.com/doc/refman/8.0/en/mysqlbinlog.html) |
| `pt-online-schema-change` | 在线DDL | [文档](https://www.percona.com/doc/percona-toolkit/LATEST/pt-online-schema-change.html) |
| `gh-ost` | 在线DDL | [GitHub](https://github.com/github/gh-ost) |

### 监控工具

| 工具 | 功能 | URL |
|------|------|-----|
| Prometheus | 指标采集 | http://localhost:9090 |
| Grafana | 可视化面板 | http://localhost:3000 |
| phpMyAdmin | Web管理 | http://localhost:8080 |

---

**维护团队**: 数据库架构组
**最后更新**: 2025-01-03
**版本**: v1.0
