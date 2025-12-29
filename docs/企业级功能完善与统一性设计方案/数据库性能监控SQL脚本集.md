# 数据库性能监控SQL脚本集

**文档编号**: DE-DD-2025-MONITOR-001
**版本**: v1.0.0
**创建日期**: 2025-12-30
**参考**: 《数据库设计最佳实践与优化指南》

---

## 📋 文档说明

本文档提供 Coze Studio 项目的**数据库性能监控SQL脚本**,用于实时监控数据库性能和诊断问题。

**监控维度**:
- ✅ 连接数监控
- ✅ 慢查询监控
- ✅ 索引使用情况
- ✅ 表空间使用
- ✅ InnoDB缓冲池状态
- ✅ 主从复制状态
- ✅ 锁等待情况

---

## 1. 连接数监控

### 1.1 当前连接数统计

```sql
-- ============================================
-- 监控脚本: 当前连接数统计
-- 执行频率: 每分钟
-- 告警阈值: 使用率>80%
-- ============================================

SELECT
    -- 当前连接数
    COUNT(*) AS current_connections,

    -- 最大连接数配置
    (SELECT @@max_connections) AS max_connections,

    -- 连接使用率
    ROUND(COUNT(*) / (SELECT @@max_connections) * 100, 2) AS usage_percent,

    -- 按状态分组
    SUM(CASE WHEN command = 'Sleep' THEN 1 ELSE 0 END) AS sleep_connections,
    SUM(CASE WHEN command = 'Query' THEN 1 ELSE 0 END) AS query_connections,
    SUM(CASE WHEN command IN ('INSERT', 'UPDATE', 'DELETE') THEN 1 ELSE 0 END) AS write_connections,
    SUM(CASE WHEN command = 'SELECT' THEN 1 ELSE 0 END) AS read_connections
FROM information_schema.PROCESSLIST
WHERE ID != CONNECTION_ID();

-- 告警条件
-- IF usage_percent > 80 THEN
--     -- 发送告警
--     INSERT INTO monitoring_alerts (metric_name, metric_value, alert_level, message)
--     VALUES ('connection_usage', usage_percent, 'WARNING', '数据库连接数使用率超过80%');
-- END IF;
```

### 1.2 连接数趋势(历史)

```sql
-- ============================================
-- 监控脚本: 连接数趋势(需要历史表)
-- 执行频率: 每小时
-- ============================================

-- 创建历史记录表
CREATE TABLE IF NOT EXISTS monitoring_connection_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    recorded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    current_connections INT NOT NULL,
    max_connections INT NOT NULL,
    usage_percent DECIMAL(5, 2) NOT NULL,
    INDEX idx_recorded_at (recorded_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='连接数监控历史';

-- 记录当前连接数
INSERT INTO monitoring_connection_history (current_connections, max_connections, usage_percent)
SELECT
    COUNT(*) AS current_connections,
    (SELECT @@max_connections) AS max_connections,
    ROUND(COUNT(*) / (SELECT @@max_connections) * 100, 2) AS usage_percent
FROM information_schema.PROCESSLIST;

-- 查询最近24小时连接数趋势
SELECT
    DATE_FORMAT(recorded_at, '%Y-%m-%d %H:00') AS hour,
    AVG(current_connections) AS avg_connections,
    MAX(current_connections) AS max_connections,
    AVG(usage_percent) AS avg_usage_percent
FROM monitoring_connection_history
WHERE recorded_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
GROUP BY DATE_FORMAT(recorded_at, '%Y-%m-%d %H:00')
ORDER BY hour DESC;
```

### 1.3 长时间运行的连接

```sql
-- ============================================
-- 监控脚本: 长时间运行的连接
-- 执行频率: 每分钟
-- 告警阈值: 运行时间>300秒
-- ============================================

SELECT
    ID AS connection_id,
    USER AS connection_user,
    HOST AS connection_host,
    DB AS database_name,
    COMMAND AS command,
    TIME AS duration_seconds,
    ROUND(TIME / 60, 2) AS duration_minutes,
    STATE AS state,
    LEFT(INFO, 100) AS query_preview
FROM information_schema.PROCESSLIST
WHERE TIME > 300  -- 超过5分钟
  AND COMMAND != 'Sleep'
ORDER BY TIME DESC;

-- 如果结果不为空,可能需要优化或终止长时间运行的查询
-- KILL QUERY connection_id;  -- 终止查询,保留连接
-- KILL connection_id;         -- 终止连接
```

---

## 2. 慢查询监控

### 2.1 实时慢查询统计

```sql
-- ============================================
-- 监控脚本: 慢查询统计
-- 执行频率: 每分钟
-- 数据源: 慢查询日志(需开启slow_query_log)
-- ============================================

SELECT
    -- 慢查询总数
    (SELECT @@slow_queries) AS total_slow_queries,

    -- 慢查询日志状态
    (SELECT @@slow_query_log) AS slow_query_log_enabled,

    -- 慢查询阈值(秒)
    (SELECT @@long_query_time) AS long_query_time_seconds,

    -- 当前慢查询数量(从performance_schema)
    COUNT(*) AS current_slow_queries
FROM performance_schema.events_statements_summary_by_digest
WHERE SUM_TIMER_WAIT > 500000000000  -- 超过500ms
  AND (SELECT @@slow_query_log) = 'ON';
```

### 2.2 慢查询TOP10

```sql
-- ============================================
-- 监控脚本: 慢查询TOP10
-- 执行频率: 每小时
-- 数据源: performance_schema
-- ============================================

SELECT
    DIGEST_TEXT AS query,
    COUNT_STAR AS exec_count,
    ROUND(AVG_TIMER_WAIT / 1000000000000, 2) AS avg_time_sec,
    ROUND(MAX_TIMER_WAIT / 1000000000000, 2) AS max_time_sec,
    ROUND(SUM_TIMER_WAIT / 1000000000000, 2) AS total_time_sec,
    SUM_ROWS_EXAMINED AS total_rows_examined,
    ROUND(SUM_ROWS_EXAMINED / COUNT_STAR, 2) AS avg_rows_examined,
    SUM_ERRORS AS errors,
    SUM_WARNINGS AS warnings
FROM performance_schema.events_statements_summary_by_digest
WHERE SCHEMA_NAME = 'coze_studio'
  -- 筛选实际执行过的查询
  AND COUNT_STAR > 10
  -- 按平均执行时间排序
ORDER BY AVG_TIMER_WAIT DESC
LIMIT 10;
```

### 2.3 慢查询详细分析

```sql
-- ============================================
-- 分析脚本: 慢查询详细执行计划
-- ============================================

-- 替换查询文本进行分析
EXPLAIN SELECT * FROM bot_messages WHERE conversation_id = ? ORDER BY created_at DESC LIMIT 20;

-- 关键指标:
-- type: 查询类型(ALL=全表扫描, index=索引扫描, range=范围扫描, ref=引用)
-- rows: 预估扫描行数
-- Extra: Using index(覆盖索引), Using filesort(文件排序), Using temporary(临时表)
```

---

## 3. 索引使用监控

### 3.1 索引使用情况统计

```sql
-- ============================================
-- 监控脚本: 索引使用情况
-- 执行频率: 每小时
-- ============================================

SELECT
    object_name AS table_name,
    index_name,
    IF(index_name = 'PRIMARY', '主键索引',
       IF(index_name LIKE 'uk_%', '唯一索引',
       IF(index_name LIKE 'idx_%', '普通索引', '其他'))) AS index_type,
    count_star AS usage_count,
    ROUND(sum_timer_wait / 1000000000000, 2) AS total_time_sec,
    ROUND(avg_timer_wait / 1000000000000, 4) AS avg_time_sec,
    ROUND(count_fetch / NULLIF(count_star, 0), 2) AS fetch_per_query
FROM performance_schema.table_io_waits_summary_by_index_usage
WHERE object_schema = 'coze_studio'
  AND object_name IN ('bots', 'bot_conversations', 'bot_messages', 'users')  -- 关键业务表
ORDER BY object_name, count_star DESC;
```

### 3.2 未使用的索引

```sql
-- ============================================
-- 监控脚本: 查找未使用的索引
-- 执行频率: 每天
-- 用途: 清理冗余索引
-- ============================================

SELECT
    TABLE_NAME,
    INDEX_NAME,
    ROUND(STAT_VALUE * @@innodb_page_size / 1024 / 1024, 2) AS size_mb,
    COLUMN_NAME
FROM information_schema.STATISTICS s
WHERE TABLE_SCHEMA = 'coze_studio'
  AND INDEX_NAME != 'PRIMARY'
  -- 索引未被使用
  AND NOT EXISTS (
    SELECT 1
    FROM performance_schema.table_io_waits_summary_by_index_usage p
    WHERE p.object_schema = s.TABLE_SCHEMA
      AND p.object_name = s.TABLE_NAME
      AND p.index_name = s.INDEX_NAME
      AND p.count_star > 0
  )
ORDER BY size_mb DESC;

-- 注意: 主键索引和唯一约束索引不应删除
-- 删除前先确认是否有业务需要
-- ALTER TABLE bots DROP INDEX idx_xxx;
```

### 3.3 索引选择性分析

```sql
-- ============================================
-- 分析脚本: 索引选择性
-- 执行频率: 每次/新增索引时
-- 用途: 评估索引效果
-- ============================================

SELECT
    TABLE_NAME,
    INDEX_NAME,
    COLUMN_NAME,
    CARDINALITY,
    TABLE_ROWS,
    -- 选择性 = 唯一值数量 / 总行数
    ROUND(CARDINALITY / NULLIF(TABLE_ROWS, 0), 4) AS selectivity,
    CASE
        WHEN ROUND(CARDINALITY / NULLIF(TABLE_ROWS, 0), 4) > 0.9 THEN '优秀'
        WHEN ROUND(CARDINALITY / NULLIF(TABLE_ROWS, 0), 4) > 0.7 THEN '良好'
        WHEN ROUND(CARDINALITY / NULLIF(TABLE_ROWS, 0), 4) > 0.5 THEN '一般'
        ELSE '较差'
    END AS selectivity_rating
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = 'coze_studio'
  AND INDEX_NAME != 'PRIMARY'
  AND TABLE_ROWS > 0
ORDER BY selectivity DESC, TABLE_NAME, INDEX_NAME;

-- 选择性解释:
-- > 0.9: 优秀,强烈建议建索引
-- 0.7-0.9: 良好,建议建索引
-- 0.5-0.7: 一般,考虑建索引
-- 0.1-0.5: 较差,谨慎建索引
-- < 0.1:   很差,不建议建索引(如gender, is_active)
```

---

## 4. 表空间监控

### 4.1 表大小统计

```sql
-- ============================================
-- 监控脚本: 表大小统计
-- 执行频率: 每小时
-- 告警阈值: 单表>10GB
-- ============================================

SELECT
    TABLE_NAME,
    ENGINE AS storage_engine,
    TABLE_ROWS,
    ROUND((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024, 2) AS total_mb,
    ROUND(DATA_LENGTH / 1024 / 1024, 2) AS data_mb,
    ROUND(INDEX_LENGTH / 1024 / 1024, 2) AS index_mb,
    ROUND(INDEX_LENGTH / NULLIF(DATA_LENGTH, 1), 2) AS index_data_ratio,
    ROUND(DATA_FREE / 1024 / 1024, 2) AS data_free_mb,
    AUTO_INCREMENT,
    TABLE_COLLATION,
    CREATE_TIME,
    UPDATE_TIME
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = 'coze_studio'
  AND TABLE_TYPE = 'BASE TABLE'
ORDER BY (DATA_LENGTH + INDEX_LENGTH) DESC;

-- 告警逻辑
-- IF total_mb > 10240 THEN  -- 10GB
--     INSERT INTO monitoring_alerts (metric_name, metric_value, alert_level, message)
--     VALUES (TABLE_NAME, total_mb, 'WARNING', CONCAT('表 ', TABLE_NAME, ' 大小超过10GB'));
-- END IF;
```

### 4.2 表碎片率分析

```sql
-- ============================================
-- 监控脚本: 表碎片率
-- 执行频率: 每天
-- 告警阈值: 碎片率>10%
-- ============================================

SELECT
    TABLE_NAME,
    TABLE_ROWS,
    ROUND(DATA_LENGTH / 1024 / 1024, 2) AS data_mb,
    ROUND(DATA_FREE / 1024 / 1024, 2) AS data_free_mb,
    -- 碎片率 = 碎片空间 / 总空间
    ROUND(DATA_FREE / NULLIF(DATA_LENGTH, 0) * 100, 2) AS fragmentation_percent,
    CASE
        WHEN ROUND(DATA_FREE / NULLIF(DATA_LENGTH, 0) * 100, 2) > 10 THEN '需要优化'
        ELSE '正常'
    END AS status
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = 'coze_studio'
  AND TABLE_TYPE = 'BASE TABLE'
  AND DATA_FREE > 0
ORDER BY fragmentation_percent DESC;

-- 优化建议(碎片率>10%的表)
-- OPTIMIZE TABLE bots;
-- OPTIMIZE TABLE bot_conversations;
-- 注意: OPTIMIZE会锁表,建议在业务低峰期执行
```

### 4.3 数据库总大小

```sql
-- ============================================
-- 监控脚本: 数据库总大小
-- 执行频率: 每小时
-- ============================================

SELECT
    table_schema AS database_name,
    COUNT(*) AS table_count,
    ROUND(SUM(DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024, 2) AS total_mb,
    ROUND(SUM(DATA_LENGTH) / 1024 / 1024, 2) AS data_mb,
    ROUND(SUM(INDEX_LENGTH) / 1024 / 1024, 2) AS index_mb,
    ROUND(SUM(DATA_FREE) / 1024 / 1024, 2) AS data_free_mb,
    SUM(TABLE_ROWS) AS total_rows
FROM information_schema.TABLES
WHERE table_schema = 'coze_studio'
GROUP BY table_schema;
```

---

## 5. InnoDB缓冲池监控

### 5.1 缓冲池状态

```sql
-- ============================================
-- 监控脚本: InnoDB缓冲池状态
-- 执行频率: 每分钟
-- ============================================

SELECT
    -- 缓冲池大小
    (SELECT @@innodb_buffer_pool_size) AS buffer_pool_bytes,
    ROUND((SELECT @@innodb_buffer_pool_size) / 1024 / 1024 / 1024, 2) AS buffer_pool_gb,

    -- 缓冲池实例数
    (SELECT @@innodb_buffer_pool_instances) AS buffer_pool_instances,

    -- 缓冲池块大小
    (SELECT @@innodb_buffer_pool_chunk_size) AS buffer_pool_chunk_bytes,

    -- 读取统计
    (SELECT VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Innodb_buffer_pool_read_requests') AS read_requests,
    (SELECT VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Innodb_buffer_pool_reads') AS read_from_disk,

    -- 缓冲池命中率 = 1 - (磁盘读取 / 总读取)
    ROUND(
        1 - (
            (SELECT VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Innodb_buffer_pool_reads') /
            NULLIF((SELECT VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Innodb_buffer_pool_read_requests'), 0)
        ),
        4
    ) AS hit_rate,

    CASE
        WHEN ROUND(
            1 - (
                (SELECT VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Innodb_buffer_pool_reads') /
                NULLIF((SELECT VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Innodb_buffer_pool_read_requests'), 0)
            ),
            4
        ) > 0.95 THEN '优秀'
        WHEN ROUND(
            1 - (
                (SELECT VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Innodb_buffer_pool_reads') /
                NULLIF((SELECT VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Innodb_buffer_pool_read_requests'), 0)
            ),
            4
        ) > 0.9 THEN '良好'
        ELSE '需优化'
    END AS hit_rate_rating;
```

### 5.2 缓冲池使用详情

```sql
-- ============================================
-- 监控脚本: 缓冲池使用详情
-- 执行频率: 每小时
-- ============================================

SELECT
    pool_id,
    pool_size,
    free_buffers,
    database_pages,
    old_database_pages,
    modified_database_pages,
    pending_reads,
    pending_flush_lru,
    pending_flush_list
FROM information_schema.INNODB_BUFFER_POOL_STATS;
```

---

## 6. 锁等待监控

### 6.1 当前锁等待

```sql
-- ============================================
-- 监控脚本: 当前锁等待
-- 执行频率: 每分钟
-- 告警条件: 等待时间>10秒
-- ============================================

SELECT
    r.trx_id AS waiting_trx_id,
    r.trx_mysql_thread_id AS waiting_thread,
    r.trx_query AS waiting_query,
    r.trx_wait_time AS wait_time_seconds,
    r.trx_state AS waiting_state,
    b.trx_id AS blocking_trx_id,
    b.trx_mysql_thread_id AS blocking_thread,
    b.trx_query AS blocking_query,
    b.trx_state AS blocking_state,
    CONCAT('KILL ', b.trx_mysql_thread_id) AS kill_command
FROM information_schema.innodb_lock_waits w
JOIN information_schema.innodb_trx b ON b.trx_id = w.blocking_trx_id
JOIN information_schema.innodb_trx r ON r.trx_id = w.requesting_trx_id
WHERE r.trx_wait_time > 10  -- 等待超过10秒
ORDER BY r.trx_wait_time DESC;

-- 如果结果不为空,说明有锁等待
-- 解决方案:
-- 1. 优化长时间运行的事务
-- 2. 终止阻塞事务: KILL blocking_thread
-- 3. 优化SQL,减少事务持有锁的时间
```

### 6.2 历史锁等待统计

```sql
-- ============================================
-- 监控脚本: 锁等待历史统计
-- 执行频率: 每小时
-- ============================================

-- 创建锁等待历史表
CREATE TABLE IF NOT EXISTS monitoring_lock_waits_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    recorded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    waiting_trx_id VARCHAR(50),
    waiting_query TEXT,
    wait_time_seconds INT,
    blocking_trx_id VARCHAR(50),
    blocking_query TEXT,
    INDEX idx_recorded_at (recorded_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='锁等待历史';

-- 记录当前锁等待
INSERT INTO monitoring_lock_waits_history (waiting_trx_id, waiting_query, wait_time_seconds, blocking_trx_id, blocking_query)
SELECT
    r.trx_id,
    r.trx_query,
    r.trx_wait_time,
    b.trx_id,
    b.trx_query
FROM information_schema.innodb_lock_waits w
JOIN information_schema.innodb_trx b ON b.trx_id = w.blocking_trx_id
JOIN information_schema.innodb_trx r ON r.trx_id = w.requesting_trx_id
WHERE r.trx_wait_time > 5;  -- 等待超过5秒

-- 查询最近24小时锁等待统计
SELECT
    DATE_FORMAT(recorded_at, '%Y-%m-%d %H:00') AS hour,
    COUNT(*) AS lock_wait_count,
    AVG(wait_time_seconds) AS avg_wait_seconds,
    MAX(wait_time_seconds) AS max_wait_seconds
FROM monitoring_lock_waits_history
WHERE recorded_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
GROUP BY DATE_FORMAT(recorded_at, '%Y-%m-%d %H:00')
ORDER BY hour DESC;
```

---

## 7. 主从复制监控

### 7.1 主从复制状态

```sql
-- ============================================
-- 监控脚本: 主从复制状态
-- 执行频率: 每分钟
-- 告警阈值: 延迟>10秒
-- ============================================

SHOW SLAVE STATUS\G

-- 关键指标:
-- Slave_IO_Running: Yes/No (IO线程是否运行)
-- Slave_SQL_Running: Yes/No (SQL线程是否运行)
-- Seconds_Behind_Master: 延迟秒数(0=无延迟)
-- Last_Error: 最后错误信息
-- Exec_Master_Log_Pos: 执行到的主日志位置

-- 自动化查询版本
SELECT
    Slave_IO_Running AS io_running,
    Slave_SQL_Running AS sql_running,
    Seconds_Behind_Master AS delay_seconds,
    Last_IO_Error AS io_error,
    Last_SQL_Error AS sql_error,
    CASE
        WHEN Slave_IO_Running = 'Yes' AND Slave_SQL_Running = 'Yes' AND Seconds_Behind_Master < 10 THEN '正常'
        WHEN Slave_IO_Running = 'No' OR Slave_SQL_Running = 'No' THEN '异常'
        WHEN Seconds_Behind_Master >= 10 THEN '延迟过高'
        ELSE '未知'
    END AS status
FROM (SHOW SLAVE STATUS) AS slave_status;
```

### 7.2 主从同步趋势

```sql
-- ============================================
-- 监控脚本: 主从同步延迟趋势
-- 执行频率: 每小时
-- ============================================

-- 创建同步延迟历史表
CREATE TABLE IF NOT EXISTS monitoring_replication_lag_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    recorded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    delay_seconds INT NOT NULL,
    io_running VARCHAR(10),
    sql_running VARCHAR(10),
    INDEX idx_recorded_at (recorded_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='主从复制延迟历史';

-- 记录当前延迟
INSERT INTO monitoring_replication_lag_history (delay_seconds, io_running, sql_running)
SELECT
    Seconds_Behind_Master,
    Slave_IO_Running,
    Slave_SQL_Running
FROM (SHOW SLAVE STATUS) AS slave_status;

-- 查询最近24小时延迟趋势
SELECT
    DATE_FORMAT(recorded_at, '%Y-%m-%d %H:00') AS hour,
    AVG(delay_seconds) AS avg_delay,
    MAX(delay_seconds) AS max_delay,
    MIN(delay_seconds) AS min_delay
FROM monitoring_replication_lag_history
WHERE recorded_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
GROUP BY DATE_FORMAT(recorded_at, '%Y-%m-%d %H:00')
ORDER BY hour DESC;
```

---

## 8. 综合性能仪表板

### 8.1 性能总览

```sql
-- ============================================
-- 监控脚本: 性能总览仪表板
-- 执行频率: 每分钟
-- ============================================

SELECT
    -- 连接数
    (SELECT COUNT(*) FROM information_schema.PROCESSLIST) AS current_connections,
    ROUND((SELECT COUNT(*) FROM information_schema.PROCESSLIST) / (SELECT @@max_connections) * 100, 2) AS connection_usage_percent,

    -- QPS (每秒查询数)
    (SELECT VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Questions') /
    (SELECT VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Uptime') AS qps,

    -- 慢查询
    (SELECT @@slow_queries) AS total_slow_queries,

    -- InnoDB缓冲池命中率
    ROUND(
        1 - (
            (SELECT VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Innodb_buffer_pool_reads') /
            NULLIF((SELECT VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Innodb_buffer_pool_read_requests'), 0)
        ),
        4
    ) AS buffer_pool_hit_rate,

    -- 主从延迟
    (SELECT Seconds_Behind_Master FROM (SHOW SLAVE STATUS) AS slave_status) AS replication_delay_seconds,

    -- 当前时间
    NOW() AS current_time;
```

### 8.2 表性能统计

```sql
-- ============================================
-- 监控脚本: 关键表性能统计
-- 执行频率: 每小时
-- ============================================

SELECT
    OBJECT_NAME AS table_name,
    COUNT_READ AS read_count,
    COUNT_WRITE AS write_count,
    ROUND(COUNT_READ + COUNT_WRITE) AS total_count,
    ROUND(SUM_TIMER_WAIT / 1000000000000, 2) AS total_time_sec,
    ROUND(AVG_TIMER_WAIT / 1000000000000, 4) AS avg_time_sec,
    ROUND(SUM_NUMBER_OF_BYTES_READ / 1024 / 1024, 2) AS data_read_mb,
    ROUND(SUM_NUMBER_OF_BYTES_WRITE / 1024 / 1024, 2) AS data_written_mb
FROM performance_schema.table_io_waits_summary_by_table
WHERE OBJECT_SCHEMA = 'coze_studio'
  AND OBJECT_NAME IN ('bots', 'bot_conversations', 'bot_messages', 'users', 'subscriptions')
ORDER BY total_count DESC;
```

---

## 9. 告警规则

### 9.1 告警规则定义

```sql
-- ============================================
-- 告警规则表
-- ============================================

CREATE TABLE IF NOT EXISTS monitoring_alert_rules (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    rule_name VARCHAR(255) NOT NULL UNIQUE,
    metric_name VARCHAR(100) NOT NULL,
    operator ENUM('>', '<', '>=', '<=', '=') NOT NULL,
    threshold_value DECIMAL(20, 4) NOT NULL,
    alert_level ENUM('INFO', 'WARNING', 'CRITICAL') NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    description TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='监控告警规则';

-- 插入告警规则
INSERT INTO monitoring_alert_rules (rule_name, metric_name, operator, threshold_value, alert_level, description) VALUES
('连接数使用率过高', 'connection_usage_percent', '>', 80, 'WARNING', '数据库连接数使用率超过80%'),
('缓冲池命中率过低', 'buffer_pool_hit_rate', '<', 0.9, 'WARNING', 'InnoDB缓冲池命中率低于90%'),
('主从复制延迟过高', 'replication_delay_seconds', '>', 10, 'WARNING', '主从复制延迟超过10秒'),
('表大小过大', 'table_size_mb', '>', 10240, 'WARNING', '表大小超过10GB'),
('锁等待时间过长', 'lock_wait_seconds', '>', 10, 'WARNING', '锁等待时间超过10秒');
```

### 9.2 告警历史记录

```sql
-- ============================================
-- 告警历史表
-- ============================================

CREATE TABLE IF NOT EXISTS monitoring_alerts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    rule_name VARCHAR(255) NOT NULL,
    metric_name VARCHAR(100) NOT NULL,
    metric_value DECIMAL(20, 4) NOT NULL,
    alert_level ENUM('INFO', 'WARNING', 'CRITICAL') NOT NULL,
    message TEXT,
    resolved BOOLEAN DEFAULT FALSE,
    resolved_at DATETIME DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_alert_level (alert_level),
    INDEX idx_resolved (resolved),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='监控告警历史';

-- 查询未解决的告警
SELECT
    rule_name,
    metric_name,
    metric_value,
    alert_level,
    message,
    created_at
FROM monitoring_alerts
WHERE resolved = FALSE
ORDER BY
    CASE alert_level
        WHEN 'CRITICAL' THEN 1
        WHEN 'WARNING' THEN 2
        WHEN 'INFO' THEN 3
    END,
    created_at DESC;
```

---

## 10. 监控仪表板SQL汇总

### 10.1 Grafana/Prometheus查询

```sql
-- ============================================
-- Prometheus Exporter查询示例
-- ============================================

-- 导出器: mysqld_exporter
-- 指标: mysql_global_status_questions

-- 连接数使用率
mysql_global_variables_max_connections - mysql_global_status_threads_connected

-- QPS
rate(mysql_global_status_questions[1m])

-- 慢查询
rate(mysql_global_status_slow_queries[5m])

-- InnoDB缓冲池命中率
(
    mysql_global_status_innodb_buffer_pool_read_requests -
    mysql_global_status_innodb_buffer_pool_reads
) / mysql_global_status_innodb_buffer_pool_read_requests
```

### 10.2 健康检查脚本

```sql
-- ============================================
-- 健康检查脚本(用于监控探针)
-- 执行频率: 每30秒
-- ============================================

SELECT
    'coze_studio' AS service_name,
    'mysql' AS component_name,
    CASE
        -- 检查1: 连接数正常
        WHEN ROUND((SELECT COUNT(*) FROM information_schema.PROCESSLIST) / (SELECT @@max_connections) * 100, 2) > 90 THEN 'UNHEALTHY'
        -- 检查2: 主从复制正常
        WHEN (SELECT Seconds_Behind_Master FROM (SHOW SLAVE STATUS) AS slave_status) > 60 THEN 'UNHEALTHY'
        -- 检查3: 有未解决的CRITICAL告警
        WHEN (SELECT COUNT(*) FROM monitoring_alerts WHERE resolved = FALSE AND alert_level = 'CRITICAL') > 0 THEN 'UNHEALTHY'
        ELSE 'HEALTHY'
    END AS health_status,
    NOW() AS check_time,
    (SELECT COUNT(*) FROM information_schema.PROCESSLIST) AS current_connections,
    (SELECT Seconds_Behind_Master FROM (SHOW SLAVE STATUS) AS slave_status) AS replication_delay_seconds;
```

---

## 11. 定期维护脚本

### 11.1 每日维护任务

```sql
-- ============================================
-- 维护脚本: 每日任务
-- 执行时间: 每天凌晨2点
-- ============================================

-- 任务1: 清理监控历史数据(保留30天)
DELETE FROM monitoring_connection_history
WHERE recorded_at < DATE_SUB(NOW(), INTERVAL 30 DAY);

DELETE FROM monitoring_lock_waits_history
WHERE recorded_at < DATE_SUB(NOW(), INTERVAL 30 DAY);

DELETE FROM monitoring_replication_lag_history
WHERE recorded_at < DATE_SUB(NOW(), INTERVAL 30 DAY);

-- 任务2: 分析表(更新统计信息)
ANALYZE TABLE bots;
ANALYZE TABLE bot_conversations;
ANALYZE TABLE bot_messages;
ANALYZE TABLE users;
ANALYZE TABLE subscriptions;

-- 任务3: 清理已解决的告警(保留90天)
DELETE FROM monitoring_alerts
WHERE resolved = TRUE
  AND resolved_at < DATE_SUB(NOW(), INTERVAL 90 DAY);
```

### 11.2 每周维护任务

```sql
-- ============================================
-- 维护脚本: 每周任务
-- 执行时间: 每周日凌晨3点
-- ============================================

-- 任务1: 优化表(清理碎片)
-- 注意: 会锁表,建议在业务低峰期执行
OPTIMIZE TABLE bots;
OPTIMIZE TABLE bot_conversations;
OPTIMIZE TABLE bot_messages;
OPTIMIZE TABLE bot_store_listings;

-- 任务2: 检查表
CHECK TABLE bots;
CHECK TABLE bot_conversations;
CHECK TABLE bot_messages;

-- 任务3: 汇总监控数据(按天聚合)
-- (根据需要创建汇总表和聚合脚本)
```

---

## 12. 快速诊断脚本

### 12.1 性能快速诊断

```sql
-- ============================================
-- 诊断脚本: 快速性能诊断
-- 用途: 性能问题排查
-- ============================================

-- 1. 当前运行中的查询
SELECT * FROM information_schema.PROCESSLIST WHERE COMMAND != 'Sleep' ORDER BY TIME DESC LIMIT 20;

-- 2. 当前锁等待
SELECT * FROM information_schema.innodb_lock_waits;

-- 3. 最近的慢查询
SELECT * FROM mysql.slow_log ORDER BY query_time DESC LIMIT 20;

-- 4. 表碎片率
SELECT
    TABLE_NAME,
    ROUND(DATA_FREE / 1024 / 1024, 2) AS data_free_mb,
    ROUND(DATA_FREE / NULLIF(DATA_LENGTH, 0) * 100, 2) AS fragmentation_percent
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = 'coze_studio'
  AND DATA_FREE > 0
ORDER BY fragmentation_percent DESC
LIMIT 10;

-- 5. 未使用的索引
-- (见上文"未使用的索引"脚本)
```

### 12.2 瓶颈分析

```sql
-- ============================================
-- 诊断脚本: 系统瓶颈分析
-- ============================================

-- 1. CPU瓶颈
SELECT * FROM information_schema.PROCESSLIST WHERE TIME > 60 ORDER BY TIME DESC;

-- 2. IO瓶颈
SELECT * FROM performance_schema.file_summary_by_instance
ORDER BY SUM_TIMER_WAIT DESC LIMIT 10;

-- 3. 内存瓶颈
SELECT
    (SELECT @@innodb_buffer_pool_size) AS buffer_pool_size_bytes,
    ROUND((SELECT @@innodb_buffer_pool_size) / 1024 / 1024 / 1024, 2) AS buffer_pool_size_gb,
    (SELECT VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Innodb_buffer_pool_pages_dirty') AS dirty_pages;

-- 4. 网络瓶颈
SELECT
    (SELECT @@max_connections) AS max_connections,
    (SELECT COUNT(*) FROM information_schema.PROCESSLIST) AS current_connections,
    (SELECT SUM(@@read_buffer_size + @@sort_buffer_size) / 1024 / 1024 FROM information_schema.PROCESSLIST) AS total_buffer_mb;
```

---

**文档版本历史**:
- v1.0.0 (2025-12-30): 初始版本,包含12大类监控脚本
