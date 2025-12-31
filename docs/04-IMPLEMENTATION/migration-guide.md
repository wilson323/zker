# ZKER 数据迁移操作手册

> **版本**: v1.0
> **更新日期**: 2025-01-03
> **适用范围**: 企业级多租户数据迁移
> **执行人员**: 数据库管理员、后端工程师

---

## 📋 目录

1. [迁移概述](#迁移概述)
2. [迁移前准备](#迁移前准备)
3. [迁移策略](#迁移策略)
4. [迁移执行](#迁移执行)
5. [验证与测试](#验证与测试)
6. [回滚方案](#回滚方案)
7. [监控与告警](#监控与告警)
8. [故障排查](#故障排查)
9. [FAQ](#faq)

---

## 迁移概述

### 迁移目标

为所有业务表添加 `tenant_id` 字段，实现企业级多租户数据隔离。

### 迁移范围

| 表名 | 记录数（预估） | 迁移时间（预估） |
|------|--------------|---------------|
| bots | 100万 | 30分钟 |
| conversations | 500万 | 2小时 |
| knowledge | 50万 | 15分钟 |
| knowledge_document | 200万 | 1小时 |
| knowledge_document_slice | 1000万 | 4小时 |
| workflows | 30万 | 10分钟 |
| messages | 2000万 | 8小时 |
| **总计** | **~3680万** | **~16小时** |

### 质量标准

- ✅ 数据丢失率：0%
- ✅ 数据准确率：100%
- ✅ 迁移成功率：≥99.9%
- ✅ 业务停机时间：≤4小时

---

## 迁移前准备

### 环境检查

#### 1. 数据库版本检查

```bash
mysql -u root -p -e "SELECT VERSION();"
```

**要求**: MySQL 8.0+

#### 2. 磁盘空间检查

```bash
# 检查可用空间（需要至少2倍数据库大小）
df -h /var/lib/mysql
```

**要求**: 可用空间 ≥ 200GB

#### 3. 内存检查

```bash
free -h
```

**要求**: 可用内存 ≥ 8GB

#### 4. 网络连接检查

```bash
ping -c 3 db_hostname
telnet db_hostname 3306
```

### 备份数据库

#### 全量备份

```bash
#!/bin/bash
# backup_database.sh

BACKUP_DIR="/backup/migration"
BACKUP_FILE="$BACKUP_DIR/coze_studio_full_$(date +%Y%m%d_%H%M%S).sql"

mkdir -p "$BACKUP_DIR"

mysqldump -u root -p \
  --single-transaction \
  --routines \
  --triggers \
  --events \
  --master-data=2 \
  --flush-logs \
  coze_studio | gzip > "$BACKUP_FILE.gz"

echo "备份完成: $BACKUP_FILE.gz"
ls -lh "$BACKUP_FILE.gz"
```

#### 验证备份

```bash
# 测试恢复（到测试环境）
gunzip -c "$BACKUP_FILE.gz" | mysql -u root -p test_db
```

### 创建迁移表

```sql
-- 创建迁移记录表
CREATE TABLE IF NOT EXISTS migration_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    migration_id VARCHAR(100) NOT NULL UNIQUE,
    config JSON NOT NULL,
    status VARCHAR(50) NOT NULL,
    progress_percent DECIMAL(5,2) DEFAULT 0,
    total_records BIGINT DEFAULT 0,
    migrated_records BIGINT DEFAULT 0,
    failed_records BIGINT DEFAULT 0,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP NULL,
    error TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_status (status),
    INDEX idx_started_at (started_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据迁移记录表';

-- 创建迁移检查点表
CREATE TABLE IF NOT EXISTS migration_checkpoints (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    migration_id VARCHAR(100) NOT NULL,
    table_name VARCHAR(100) NOT NULL,
    last_id BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_migration_table (migration_id, table_name),
    INDEX idx_migration_id (migration_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='迁移检查点表';
```

---

## 迁移策略

### 策略选择

#### 1. 停机迁移（推荐用于中小规模）

**适用场景**:
- 数据量 < 1000万
- 可接受停机时间 ≤ 4小时
- 迁移窗口：凌晨2:00-6:00

**优点**:
- 实现简单
- 数据一致性保证
- 无需双写逻辑

**缺点**:
- 业务中断
- 风险集中

#### 2. 在线迁移（推荐用于大规模）

**适用场景**:
- 数据量 ≥ 1000万
- 要求零停机
- 24/7业务

**优点**:
- 零停机
- 风险分散
- 可灰度

**缺点**:
- 实现复杂
- 需要双写支持
- 周期较长

### 迁移模式

#### 双写模式（推荐）

```
阶段1: 添加字段 (10分钟)
  ↓
阶段2: 启用双写 (5分钟)
  ↓
阶段3: 数据回填 (主要耗时)
  ↓
阶段4: 数据验证 (30分钟)
  ↓
阶段5: 切换读写 (5分钟)
  ↓
阶段6: 停止双写 (5分钟)
```

---

## 迁移执行

### 方式一：使用Go工具（推荐）

#### 1. 编译迁移工具

```bash
cd backend/scripts/migration
go build -o migration-tool *.go
```

#### 2. 执行迁移

```bash
# 查看帮助
./migration-tool --help

# 执行完整迁移（双写模式）
./migration-tool migrate \
  --tables=bots,conversations,knowledge \
  --mode=double_write \
  --batch-size=1000 \
  --workers=10 \
  --default-tenant=default_tenant

# 仅执行数据验证
./migration-tool validate \
  --tables=bots,conversations

# 查看迁移进度
./migration-tool progress --migration-id=migration_xxx

# 回滚迁移
./migration-tool rollback --migration-id=migration_xxx
```

#### 3. 启动监控面板

```bash
# 启动Web监控面板
./monitor \
  --host=localhost \
  --port=3306 \
  --user=root \
  --password=xxx \
  --database=coze_studio \
  --web-port=8080

# 访问监控面板
open http://localhost:8080
```

### 方式二：使用Shell脚本

```bash
#!/bin/bash
# execute_migration.sh

# 配置
export DB_HOST="localhost"
export DB_PORT="3306"
export DB_USER="root"
export DB_NAME="coze_studio"

# 1. 备份数据库
bash scripts/migrate_tenant_id.sh backup

# 2. 执行迁移
bash scripts/migrate_tenant_id.sh migrate

# 3. 验证结果
bash scripts/migrate_tenant_id.sh validate

# 4. 如有问题，回滚
# bash scripts/migrate_tenant_id.sh rollback
```

### 执行步骤（双写模式）

#### 步骤1：添加tenant_id字段

```sql
-- 示例：bots表
ALTER TABLE bots
ADD COLUMN tenant_id VARCHAR(36) NULL
    COMMENT '租户ID'
    AFTER creator_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_creator (tenant_id, creator_id);
```

#### 步骤2：数据回填

```bash
# 使用批处理器并行回填
cd backend/domain/tenant/migration
go run batch_processor.go \
  --table=bots \
  --batch-size=1000 \
  --concurrency=10 \
  --default-tenant=default_tenant
```

#### 步骤3：数据验证

```bash
# 使用验证工具
cd backend/scripts/migration
go run data_validator.go \
  --tables=bots,conversations,knowledge \
  --output=validation_report.json
```

#### 步骤4：切换读写

```sql
-- 设置字段为NOT NULL（确保所有记录都有值）
ALTER TABLE bots
MODIFY COLUMN tenant_id VARCHAR(36) NOT NULL;
```

---

## 验证与测试

### 数据完整性验证

#### 1. 记录数验证

```sql
-- 验证迁移前后记录数一致
SELECT
    TABLE_NAME,
    TABLE_ROWS
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = 'coze_studio'
  AND TABLE_NAME IN ('bots', 'conversations', 'knowledge')
ORDER BY TABLE_ROWS DESC;
```

#### 2. NULL值检查

```sql
-- 检查是否有NULL的tenant_id
SELECT
    TABLE_NAME,
    COUNT(*) as total_rows,
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) as null_rows
FROM information_schema.COLUMNS c
JOIN information_schema.TABLES t ON c.TABLE_NAME = t.TABLE_NAME
WHERE t.TABLE_SCHEMA = 'coze_studio'
  AND c.COLUMN_NAME = 'tenant_id'
  AND c.TABLE_SCHEMA = 'coze_studio'
GROUP BY TABLE_NAME;
```

**预期**: null_rows = 0

#### 3. 哈希对比

```bash
# 迁移前计算哈希
./data_validator.go --tables=bots --output=before_hash.json

# 迁移后计算哈希
./data_validator.go --tables=bots --output=after_hash.json

# 对比哈希值（应该一致）
diff before_hash.json after_hash.json
```

### 业务功能验证

#### 1. 核心业务流程测试

- [ ] 创建Bot
- [ ] 创建对话
- [ ] 上传知识库文档
- [ ] 创建工作流
- [ ] 查询列表（带租户过滤）

#### 2. API测试

```bash
# 测试Bot创建API
curl -X POST http://localhost:8888/api/bot \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant_123" \
  -d '{"name": "test_bot"}'

# 测试对话查询API（应该只返回当前租户的对话）
curl -X GET http://localhost:8888/api/conversations \
  -H "X-Tenant-ID: tenant_123"
```

### 性能验证

#### 1. 基准测试

```bash
# 迁移前性能基线
./performance_benchmark.go \
  --duration=5m \
  --threads=10 \
  --output=before_performance.json

# 迁移后性能测试
./performance_benchmark.go \
  --duration=5m \
  --threads=10 \
  --output=after_performance.json

# 对比性能
# 关注指标：QPS、P95延迟、错误率
```

#### 2. 慢查询检查

```sql
-- 查看慢查询日志
SHOW VARIABLES LIKE 'slow_query%';

-- 分析慢查询
SELECT
    query_time,
    lock_time,
    rows_sent,
    rows_examined,
    sql_text
FROM mysql.slow_log
WHERE start_time > DATE_SUB(NOW(), INTERVAL 1 HOUR)
ORDER BY query_time DESC
LIMIT 10;
```

---

## 回滚方案

### 回滚触发条件

- 数据丢失率 > 0%
- 数据准确率 < 99.9%
- 迁移成功率 < 99%
- 性能下降 > 50%
- 关键业务功能异常

### 回滚步骤

#### 方式一：使用Go工具

```bash
# 自动回滚
./migration-tool rollback --migration-id=migration_xxx
```

#### 方式二：手动回滚

```sql
-- 1. 停止应用服务
sudo systemctl stop coze-studio

-- 2. 恢复数据库备份
gunzip -c /backup/migration/coze_studio_full_YYYYMMDD_HHMMSS.sql.gz | \
  mysql -u root -p coze_studio

-- 3. 验证数据
mysql -u root -p -e "USE coze_studio; SELECT COUNT(*) FROM bots;"

-- 4. 重启应用服务
sudo systemctl start coze_studio

-- 5. 验证服务
curl http://localhost:8888/health
```

---

## 监控与告警

### 启动监控面板

```bash
# 方式1：使用Go监控工具
./monitor --web-port=8080

# 方式2：使用Grafana + Prometheus
cd deploy/monitoring
docker-compose up -d
```

### 关键监控指标

| 指标 | 告警阈值 | 说明 |
|------|---------|------|
| 迁移进度 | - | 实时显示百分比 |
| QPS | < 100 | 迁移速率过低 |
| 错误率 | > 1% | 失败率过高 |
| CPU使用率 | > 80% | 系统负载过高 |
| 内存使用率 | > 85% | 内存不足 |
| 磁盘I/O | > 80MB/s | I/O瓶颈 |
| 慢查询数 | > 10/min | 性能问题 |

### 日志查看

```bash
# 查看迁移日志
tail -f /var/log/migration.log

# 查看MySQL慢查询日志
tail -f /var/log/mysql/slow-query.log

# 查看应用日志
tail -f /var/log/coze-studio/app.log
```

---

## 故障排查

### 问题1：迁移卡住不动

**症状**: 进度长时间无更新

**排查**:
```bash
# 1. 查看MySQL进程列表
mysql -u root -p -e "SHOW PROCESSLIST;"

# 2. 查看锁等待
mysql -u root -p -e "
  SELECT *
  FROM performance_schema.events_waits_current
  WHERE EVENT_NAME LIKE '%lock%';
"

# 3. 查看长事务
mysql -u root -p -e "
  SELECT *
  FROM information_schema.innodb_trx
  WHERE TIME_TO_SEC(TIMEDIFF(NOW(), trx_started)) > 60;
"
```

**解决方案**:
```sql
-- 如果有长时间锁表，考虑杀死会话
KILL <process_id>;
```

### 问题2：内存不足

**症状**: OOM错误

**排查**:
```bash
# 查看MySQL内存使用
mysql -u root -p -e "SHOW STATUS LIKE 'Innodb_buffer_pool_%';"
```

**解决方案**:
```sql
-- 减小batch_size
SET GLOBAL innodb_buffer_pool_size = 4294967296; -- 4GB
```

### 问题3：数据不一致

**症状**: 验证失败

**排查**:
```bash
# 重新验证
./data_validator.go --tables=bots --verbose

# 查看不一致详情
cat validation_report.json | jq '.results.bots.issues'
```

**解决方案**:
```sql
-- 修复NULL值
UPDATE bots
SET tenant_id = 'default_tenant'
WHERE tenant_id IS NULL OR tenant_id = '';
```

---

## FAQ

### Q1: 迁移需要多长时间？

**A**: 取决于数据量和硬件配置。参考：
- 小规模（<100万）：1-2小时
- 中规模（100-1000万）：4-8小时
- 大规模（>1000万）：8-24小时

### Q2: 迁移期间会影响业务吗？

**A**: 取决于迁移策略：
- 停机迁移：需要停机4小时
- 在线迁移：零停机，但性能可能下降10-20%

### Q3: 如何确保数据不丢失？

**A**:
1. 迁移前强制全量备份
2. 启用binlog增量备份
3. 双写模式确保数据同步
4. 严格的数据验证

### Q4: 迁移失败怎么办？

**A**:
1. 立即停止迁移
2. 分析失败原因（查看日志）
3. 如果是数据问题，修复后继续
4. 如果是系统问题，回滚后重试

### Q5: 可以暂停和恢复迁移吗？

**A**: 可以。使用检查点机制：
```bash
# 暂停（Ctrl+C）
# 恢复（使用相同命令，会从检查点继续）
./migration-tool migrate --tables=bots
```

### Q6: 迁移后性能会下降吗？

**A**:
- 索引优化后性能可能提升10-30%
- 如果性能下降，检查索引是否正确创建
- 考虑使用查询缓存

---

## 附录

### A. 快速检查清单

迁移前：
- [ ] 数据库备份完成
- [ ] 磁盘空间充足（≥200GB）
- [ ] 内存充足（≥8GB）
- [ ] 维护窗口已确认
- [ ] 迁移脚本已测试
- [ ] 监控工具已就绪

迁移中：
- [ ] 进度正常更新
- [ ] 错误率 < 1%
- [ ] 系统资源正常
- [ ] 无长锁等待

迁移后：
- [ ] 数据验证通过
- [ ] 业务功能正常
- [ ] 性能测试通过
- [ ] 监控告警正常

### B. 联系方式

- **数据库团队**: dba@company.com
- **后端团队**: backend@company.com
- **运维团队**: ops@company.com
- **紧急联系**: +86-xxx-xxxx-xxxx

### C. 相关文档

- [多租户架构设计](../03-DESIGN/multi-tenant/README.md)
- [数据库优化指南](./database-optimization-guide.md)
- [系统监控手册](../05-OPERATIONS/monitoring.md)

---

**文档维护**: 技术架构组
**最后更新**: 2025-01-03
**版本**: v1.0
