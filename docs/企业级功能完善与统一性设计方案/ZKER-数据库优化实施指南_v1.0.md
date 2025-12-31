# ZKER 数据库优化实施指南

## 📋 文档信息

**版本**: v1.0
**日期**: 2025-01-03
**执行人员**: 研发B（后端工程师）
**预计周期**: 1-2周
**风险等级**: 低-中

---

## 🎯 优化目标

### 性能目标

| 指标 | 优化前 | 优化后 | 提升幅度 |
|------|-------|-------|---------|
| **查询响应时间** | 500ms | 100ms | ⬇️ 80% |
| **慢查询数量** | 100次/天 | 20次/天 | ⬇️ 80% |
| **数据库CPU使用率** | 60% | 40% | ⬇️ 33% |
| **磁盘I/O** | 80MB/s | 50MB/s | ⬇️ 38% |

### 功能目标

- ✅ 100%表名符合规范
- ✅ 100%必要索引已添加
- ✅ 100%外键约束完整
- ✅ 100%字段类型优化

---

## 📅 实施计划

### 阶段一：准备阶段（第1-2天）

#### 任务清单

- [ ] **任务1: 环境准备**
  - [ ] 搭建测试环境
  - [ ] 复制生产数据到测试环境
  - [ ] 验证测试环境可用性

- [ ] **任务2: 备份准备**
  - [ ] 备份生产数据库（全量）
  - [ ] 启用binlog（增量备份）
  - [ ] 验证备份文件完整性

- [ ] **任务3: 脚本准备**
  - [ ] 准备优化SQL脚本
  - [ ] 准备验证SQL脚本
  - [ ] 准备回滚SQL脚本

- [ ] **任务4: 通知准备**
  - [ ] 通知开发团队
  - [ ] 通知运维团队
  - [ ] 发布维护公告

#### 验收标准

- ✅ 测试环境与生产环境数据一致
- ✅ 备份文件可正常恢复
- ✅ 所有脚本已通过语法检查

---

### 阶段二：测试阶段（第3-5天）

#### 任务清单

- [ ] **任务1: 在测试环境执行优化**
  - [ ] 执行索引优化脚本（001_add_performance_indexes.sql）
  - [ ] 执行全文索引脚本（002_add_fulltext_indexes.sql）
  - [ ] 执行字段类型优化脚本（003_optimize_field_types.sql）

- [ ] **任务2: 验证优化效果**
  - [ ] 执行验证脚本（database_optimization_verification.sql）
  - [ ] 运行性能测试
  - [ ] 验证业务功能

- [ ] **任务3: 性能对比**
  - [ ] 记录优化前性能基线
  - [ ] 记录优化后性能数据
  - [ ] 生成性能对比报告

- [ ] **任务4: 问题修复**
  - [ ] 修复发现的问题
  - [ ] 优化SQL脚本
  - [ ] 重新测试验证

#### 验收标准

- ✅ 所有SQL脚本执行成功，无错误
- ✅ 验证脚本通过率100%
- ✅ 业务功能测试通过率100%
- ✅ 性能提升达到预期（≥50%）

---

### 阶段三：生产环境实施（第6-7天）

#### 任务清单

- [ ] **任务1: 生产环境准备**
  - [ ] 选择维护窗口（凌晨2:00-4:00）
  - [ ] 准备回滚方案
  - [ ] 准备监控工具

- [ ] **任务2: 低风险优化实施**
  - [ ] 执行索引优化（预计15分钟）
  - [ ] 验证索引创建成功
  - [ ] 监控数据库性能

- [ ] **任务3: 功能验证**
  - [ ] 执行冒烟测试
  - [ ] 执行关键业务流程测试
  - [ ] 验证数据完整性

- [ ] **任务4: 监控观察**
  - [ ] 监控慢查询日志
  - [ ] 监控数据库CPU/内存/磁盘
  - [ ] 监控业务指标（QPS、响应时间）

#### 验收标准

- ✅ 优化执行成功，无回滚
- ✅ 业务功能正常，无报错
- ✅ 性能提升≥50%
- ✅ 无新增慢查询

---

## 🔧 详细实施步骤

### 步骤1: 环境准备

#### 1.1 搭建测试环境

```bash
# 1. 复制生产数据到测试环境
mysqldump -u root -p opencoze | mysql -u root -p opencoze_test

# 2. 验证数据一致性
mysql -u root -p -e "
  SELECT
    TABLE_NAME,
    TABLE_ROWS
  FROM information_schema.TABLES
  WHERE TABLE_SCHEMA = 'opencoze_test'
  ORDER BY TABLE_ROWS DESC;
"

# 3. 启用慢查询日志
mysql -u root -p -e "
  SET GLOBAL slow_query_log = 'ON';
  SET GLOBAL long_query_time = 1;
  SET GLOBAL log_queries_not_using_indexes = 'ON';
"
```

#### 1.2 备份生产数据库

```bash
# 全量备份
mysqldump -u root -p \
  --single-transaction \
  --routines \
  --triggers \
  --events \
  opencoze > opencoze_backup_$(date +%Y%m%d).sql

# 验证备份文件
ls -lh opencoze_backup_*.sql

# 测试恢复（可选）
mysql -u root -p opencoze_test < opencoze_backup_$(date +%Y%m%d).sql
```

---

### 步骤2: 测试环境优化

#### 2.1 执行索引优化

```bash
# 登录MySQL
mysql -u root -p opencoze_test

# 执行优化脚本
source backend/migrations/performance/001_add_performance_indexes.sql;
source backend/migrations/performance/002_add_fulltext_indexes.sql;
source backend/migrations/performance/003_optimize_field_types.sql;

# 验证执行结果
SHOW PROCESSLIST;
```

#### 2.2 验证优化效果

```bash
# 执行验证脚本
mysql -u root -p opencoze_test < backend/scripts/database_optimization_verification.sql > verification_result.txt

# 查看验证结果
cat verification_result.txt
```

#### 2.3 性能测试

```bash
# 使用sysbench进行性能测试
sysbench /usr/share/sysbench/oltp_read_write.lua \
  --mysql-host=localhost \
  --mysql-user=root \
  --mysql-password=xxx \
  --mysql-db=opencoze_test \
  --tables=10 \
  --table-size=10000 \
  --threads=10 \
  --time=300 \
  run

# 查看测试结果
# 关注指标：
# - queries: 总查询数
# - latency: 平均延迟
# - 95th percentile: 95分位延迟
```

---

### 步骤3: 生产环境实施

#### 3.1 执行前检查

```bash
# 1. 检查数据库连接
mysql -u root -p -e "SELECT NOW();"

# 2. 检查当前慢查询数
mysql -u root -p -e "
  SELECT
    COUNT(*) AS slow_query_count
  FROM mysql.slow_log
  WHERE start_time > DATE_SUB(NOW(), INTERVAL 1 DAY);
"

# 3. 检查表大小
mysql -u root -p -e "
  SELECT
    TABLE_NAME,
    ROUND(((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024), 2) AS size_mb
  FROM information_schema.TABLES
  WHERE TABLE_SCHEMA = DATABASE()
  ORDER BY size_mb DESC
  LIMIT 10;
"
```

#### 3.2 执行优化（低风险）

```bash
# 登录生产数据库
mysql -u root -p opencoze

-- 1. 开始事务（可选，用于回滚）
-- START TRANSACTION;

-- 2. 执行索引优化
source backend/migrations/performance/001_add_performance_indexes.sql;

-- 3. 验证索引创建
SELECT
    TABLE_NAME,
    INDEX_NAME,
    GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) AS columns
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND INDEX_NAME IN (
      'idx_tenant_status_created',
      'idx_tenant_type_created',
      'idx_tenant_conversation_created'
  )
GROUP BY TABLE_NAME, INDEX_NAME;

-- 4. 提交事务（如果开启了事务）
-- COMMIT;
```

#### 3.3 功能验证

```bash
# 1. 冒烟测试
mysql -u root -p opencoze -e "
  -- 测试1：租户Bot查询
  SELECT * FROM bots
  WHERE tenant_id = 'test_tenant_id'
    AND status = 'active'
  ORDER BY created_at DESC
  LIMIT 20;

  -- 测试2：租户消息查询
  SELECT * FROM messages
  WHERE tenant_id = 'test_tenant_id'
    AND conversation_id = 'test_conv_id'
  ORDER BY created_at DESC
  LIMIT 50;

  -- 测试3：知识库搜索
  SELECT * FROM knowledge_chunks
  WHERE MATCH(content) AGAINST('测试' IN NATURAL LANGUAGE MODE)
  LIMIT 20;
"

# 2. 查看执行计划
mysql -u root -p opencoze -e "
  EXPLAIN SELECT *
  FROM bots
  WHERE tenant_id = 'test_tenant_id'
    AND status = 'active'
  ORDER BY created_at DESC
  LIMIT 20;
"

# 预期结果：
# - type: ref 或 range
# - key: idx_tenant_status_created
# - rows: < 100
```

#### 3.4 监控观察

```bash
# 1. 查看慢查询日志
mysql -u root -p -e "
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
"

# 2. 查看数据库状态
mysql -u root -p -e "
  SHOW PROCESSLIST;
  SHOW ENGINE INNODB STATUS\G
  SHOW STATUS LIKE 'Questions';
  SHOW STATUS LIKE 'Queries';
"

# 3. 查看系统指标
# CPU使用率
top -bn1 | grep "Cpu(s)"

# 内存使用
free -h

# 磁盘I/O
iostat -x 1 5

# 网络连接
netstat -an | grep 3306 | wc -l
```

---

## ⚠️ 风险控制

### 风险识别

| 风险类型 | 风险描述 | 影响 | 概率 |
|---------|---------|------|------|
| **性能风险** | 索引创建期间锁表 | 中 | 低 |
| **空间风险** | 索引占用额外存储 | 低 | 高 |
| **兼容风险** | 字段类型变更导致应用报错 | 高 | 低 |
| **数据风险** | 字段类型变更导致数据丢失 | 高 | 极低 |

### 风险缓解措施

#### 措施1: 使用pt-online-schema-change（大表）

```bash
# 安装percona-toolkit
apt-get install percona-toolkit

# 使用pt-online-schema-change添加索引
pt-online-schema-change \
  --alter "ADD INDEX idx_tenant_status_created (tenant_id, status, created_at)" \
  --charset=utf8mb4 \
  --critical-load="Threads_running=50" \
  --max-load="Threads_running=100" \
  --chunk-size=1000 \
  --progress=time,30 \
  --dry-run \
  D=opencoze,t=bots

# 确认无误后执行
pt-online-schema-change \
  --alter "ADD INDEX idx_tenant_status_created (tenant_id, status, created_at)" \
  --charset=utf8mb4 \
  --critical-load="Threads_running=50" \
  --max-load="Threads_running=100" \
  --chunk-size=1000 \
  --progress=time,30 \
  --execute \
  D=opencoze,t=bots
```

#### 措施2: 分阶段添加索引

```sql
-- 阶段1：添加小表索引（<100万行）
CREATE INDEX idx_tenant_status_created ON bots(tenant_id, status, created_at);
CREATE INDEX idx_tenant_status_created ON conversations(tenant_id, status, created_at);

-- 阶段2：观察1小时，无异常后继续
-- 阶段3：添加大表索引（≥100万行）
CREATE INDEX idx_tenant_conversation_created ON messages(tenant_id, conversation_id, created_at) ALGORITHM=INPLACE, LOCK=NONE;
```

#### 措施3: 准备回滚方案

```sql
-- 回滚脚本（drop index）
DROP INDEX idx_tenant_status_created ON bots;
DROP INDEX idx_tenant_status_created ON conversations;
DROP INDEX idx_tenant_conversation_created ON messages;

-- 恢复字段类型
ALTER TABLE bot_store_item MODIFY COLUMN price FLOAT COMMENT '价格';
```

---

## 📊 监控指标

### 关键指标

| 指标类别 | 指标名称 | 优化前基线 | 优化后目标 | 监控频率 |
|---------|---------|-----------|-----------|---------|
| **性能指标** | QPS | 500 | 750 | 实时 |
| **性能指标** | P95响应时间 | 500ms | 250ms | 实时 |
| **性能指标** | 慢查询数 | 100次/天 | 20次/天 | 每小时 |
| **资源指标** | CPU使用率 | 60% | 40% | 实时 |
| **资源指标** | 内存使用率 | 70% | 60% | 实时 |
| **资源指标** | 磁盘I/O | 80MB/s | 50MB/s | 实时 |
| **业务指标** | 错误率 | 0.1% | 0.1% | 实时 |
| **业务指标** | 可用性 | 99.9% | 99.9% | 实时 |

### 监控工具

#### 1. Prometheus + Grafana

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'mysql'
    static_configs:
      - targets: ['localhost:9104']
```

```sql
-- MySQL Exporter配置
-- 下载：https://github.com/prometheus/mysqld_exporter

-- 启动exporter
./mysqld_exporter \
  --config.my-cnf=/etc/.my.cnf \
  --web.listen-address=:9104
```

#### 2. 慢查询日志分析

```bash
# 使用pt-query-digest分析慢查询
pt-query-digest \
  --since '24h ago' \
  --limit 10 \
  /var/log/mysql/slow-query.log

# 输出指标：
# - Query time: 查询时间
# - Lock time: 锁等待时间
# - Rows examined: 扫描行数
```

#### 3. 自定义监控脚本

```bash
#!/bin/bash
# monitor.sh - 数据库性能监控脚本

while true; do
  # 查询QPS
  QPS=$(mysql -u root -p -e "
    SELECT VARIABLE_VALUE
    FROM information_schema.GLOBAL_STATUS
    WHERE VARIABLE_NAME = 'Questions';
  " | tail -n 1)

  # 查询连接数
  CONNECTIONS=$(mysql -u root -p -e "
    SELECT COUNT(*)
    FROM information_schema.PROCESSLIST;
  " | tail -n 1)

  # 查询慢查询数
  SLOW_QUERIES=$(mysql -u root -p -e "
    SELECT COUNT(*)
    FROM mysql.slow_log
    WHERE start_time > DATE_SUB(NOW(), INTERVAL 5 MINUTE);
  " | tail -n 1)

  echo "$(date '+%Y-%m-%d %H:%M:%S') | QPS: $QPS | Connections: $CONNECTIONS | Slow Queries: $SLOW_QUERIES"

  sleep 60
done
```

---

## ✅ 验收标准

### 功能验收

- [ ] **验收1: 所有索引创建成功**
  ```sql
  SELECT COUNT(*) FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE()
    AND INDEX_NAME IN (
        'idx_tenant_status_created',
        'idx_tenant_type_created',
        'idx_tenant_conversation_created',
        'idx_tenant_role_created',
        'idx_tenant_kb_created',
        'idx_tenant_workflow_created',
        'idx_tenant_cost_created',
        'idx_tenant_severity_created'
    );
  -- 预期：15个索引
  ```

- [ ] **验收2: 所有全文索引创建成功**
  ```sql
  SELECT COUNT(*) FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE()
    AND INDEX_TYPE = 'FULLTEXT';
  -- 预期：5个全文索引
  ```

- [ ] **验收3: 字段类型优化成功**
  ```sql
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'bot_store_item'
    AND COLUMN_NAME = 'price'
    AND DATA_TYPE = 'decimal';
  -- 预期：1行
  ```

### 性能验收

- [ ] **验收4: 查询性能提升≥50%**
  ```sql
  -- 执行性能测试
  -- 优化前：500ms
  -- 优化后：≤250ms
  ```

- [ ] **验收5: 慢查询数量降低≥70%**
  ```sql
  -- 统计慢查询数量
  -- 优化前：100次/天
  -- 优化后：≤30次/天
  ```

- [ ] **验收6: CPU使用率降低≥20%**
  ```bash
  # 查看CPU使用率
  # 优化前：60%
  # 优化后：≤48%
  ```

### 业务验收

- [ ] **验收7: 核心业务流程测试通过**
  - [ ] Bot创建流程
  - [ ] 对话流程
  - [ ] 知识库查询
  - [ ] 工作流执行

- [ ] **验收8: 数据完整性验证通过**
  ```sql
  -- 验证数据无丢失
  SELECT
      TABLE_NAME,
      TABLE_ROWS
  FROM information_schema.TABLES
  WHERE TABLE_SCHEMA = DATABASE()
  ORDER BY TABLE_ROWS DESC;
  ```

---

## 📚 附录

### A. 常用命令速查

```bash
# 查看表结构
SHOW CREATE TABLE bots;

# 查看索引
SHOW INDEX FROM bots;

# 查看表大小
SELECT
    TABLE_NAME,
    ROUND(((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024), 2) AS size_mb
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'bots';

# 查看慢查询配置
SHOW VARIABLES LIKE 'slow_query%';

# 查看进程列表
SHOW PROCESSLIST;

# 杀死查询
KILL <process_id>;

# 查看InnoDB状态
SHOW ENGINE INNODB STATUS\G

# 分析表
ANALYZE TABLE bots;

# 优化表
OPTIMIZE TABLE bots;
```

### B. 故障处理

#### 故障1: 索引创建卡住

```bash
# 查看进程列表
mysql -u root -p -e "SHOW PROCESSLIST;"

# 查看正在执行的DDL操作
mysql -u root -p -e "
  SELECT *
  FROM performance_schema.events_stages_current
  WHERE EVENT_NAME LIKE '%stage/innodb/%';
"

# 如果卡住超过1小时，考虑杀死进程
# KILL <process_id>;
```

#### 故障2: 磁盘空间不足

```bash
# 查看磁盘使用
df -h

# 清理binlog（谨慎！）
mysql -u root -p -e "
  PURGE BINARY LOGS BEFORE DATE_SUB(NOW(), INTERVAL 7 DAY);
"

# 删除旧备份
rm -f /backup/opencoze_backup_20250101.sql
```

#### 故障3: 性能下降

```bash
# 1. 查看慢查询
mysql -u root -p -e "
  SELECT * FROM mysql.slow_log
  WHERE start_time > DATE_SUB(NOW(), INTERVAL 1 HOUR)
  ORDER BY query_time DESC
  LIMIT 10;
"

# 2. 查看执行计划
mysql -u root -p -e "
  EXPLAIN SELECT * FROM bots WHERE tenant_id = 'xxx';
"

# 3. 如果索引失效，考虑回滚
# DROP INDEX idx_tenant_status_created ON bots;
```

### C. 联系方式

- **研发B**（后端工程师）：xxx@company.com
- **DBA团队**：dba@company.com
- **运维团队**：ops@company.com

---

**文档版本**: v1.0
**最后更新**: 2025-01-03
**执行人员**: 研发B
**审核人员**: 技术架构组
