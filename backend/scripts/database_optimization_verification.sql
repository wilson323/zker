-- ============================================================
-- 数据库优化验证脚本
-- 目的：验证所有优化是否成功执行，并提供性能对比
-- 执行时机：优化完成后立即执行
-- 预计时长：5-10分钟
-- ============================================================

-- ============================================================
-- 第一部分：表命名规范验证
-- ============================================================

SELECT '=== 表命名规范检查 ===' AS check_section;

-- 检查1：查找不符合规范的表名（大写字母）
SELECT
    TABLE_NAME,
    CASE
        WHEN TABLE_NAME REGEXP '[A-Z]' THEN '❌ 包含大写字母'
        WHEN TABLE_NAME LIKE '%_%' THEN '✅ 使用下划线分隔'
        ELSE '⚠️ 命名待优化'
    END AS naming_status
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME NOT LIKE '%migration%'
  AND TABLE_NAME NOT LIKE '%test%'
ORDER BY TABLE_NAME;

-- 检查2：查找使用复数形式的表名
SELECT
    TABLE_NAME,
    CASE
        WHEN TABLE_NAME REGEXP 's$' AND TABLE_NAME NOT IN ('status', 'metrics', 'logs', 'contents', 'settings', 'histories', 'transfers', 'resignations') THEN '⚠️ 可能是复数形式'
        ELSE '✅ 单数形式'
    END AS plural_status
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME NOT LIKE '%migration%'
  AND TABLE_NAME NOT LIKE '%test%'
ORDER BY TABLE_NAME;


-- ============================================================
-- 第二部分：索引完整性验证
-- ============================================================

SELECT '=== 索引完整性检查 ===' AS check_section;

-- 检查1：租户索引覆盖率
SELECT
    COUNT(DISTINCT TABLE_NAME) AS tables_with_tenant_id,
    (
        SELECT COUNT(DISTINCT TABLE_NAME)
        FROM information_schema.COLUMNS
        WHERE TABLE_SCHEMA = DATABASE()
          AND COLUMN_NAME = 'tenant_id'
    ) AS total_tables_with_tenant_id,
    CONCAT(
        ROUND(
            COUNT(DISTINCT TABLE_NAME) / (
                SELECT COUNT(DISTINCT TABLE_NAME)
                FROM information_schema.COLUMNS
                WHERE TABLE_SCHEMA = DATABASE()
                  AND COLUMN_NAME = 'tenant_id'
            ) * 100,
            2
        ),
        '%'
    ) AS coverage_rate
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND INDEX_NAME = 'idx_tenant_id';

-- 检查2：复合索引统计
SELECT
    INDEX_NAME,
    COUNT(*) AS column_count,
    GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) AS columns
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND INDEX_NAME IN (
      'idx_tenant_status_created',
      'idx_tenant_type_created',
      'idx_tenant_user_created',
      'idx_tenant_conversation_created',
      'idx_tenant_role_created',
      'idx_tenant_conv_role_created',
      'idx_tenant_kb_created',
      'idx_tenant_workflow_created',
      'idx_tenant_wf_status_created',
      'idx_tenant_cost_created',
      'idx_tenant_severity_created',
      'idx_tenant_status_hire',
      'idx_tenant_dept_status'
  )
GROUP BY INDEX_NAME
ORDER BY INDEX_NAME;

-- 检查3：全文索引统计
SELECT
    TABLE_NAME,
    INDEX_NAME,
    GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) AS columns,
    INDEX_TYPE
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND INDEX_TYPE = 'FULLTEXT'
GROUP BY TABLE_NAME, INDEX_NAME, INDEX_TYPE
ORDER BY TABLE_NAME;


-- ============================================================
-- 第三部分：外键约束完整性验证
-- ============================================================

SELECT '=== 外键约束完整性检查 ===' AS check_section;

-- 检查1：所有tenant_id外键约束
SELECT
    TABLE_NAME,
    CONSTRAINT_NAME,
    COLUMN_NAME,
    REFERENCED_TABLE_NAME,
    REFERENCED_COLUMN_NAME
FROM information_schema.KEY_COLUMN_USAGE
WHERE TABLE_SCHEMA = DATABASE()
  AND REFERENCED_TABLE_NAME = 'tenants'
  AND REFERENCED_COLUMN_NAME = 'tenant_id'
ORDER BY TABLE_NAME;

-- 检查2：外键约束统计
SELECT
    REFERENCED_TABLE_NAME,
    COUNT(*) AS fk_count,
    GROUP_CONCAT(DISTINCT TABLE_NAME) AS child_tables
FROM information_schema.KEY_COLUMN_USAGE
WHERE TABLE_SCHEMA = DATABASE()
  AND REFERENCED_TABLE_NAME IS NOT NULL
GROUP BY REFERENCED_TABLE_NAME
ORDER BY fk_count DESC;

-- 检查3：缺失外键约束的字段
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    '⚠️ 缺少外键约束' AS warning
FROM information_schema.COLUMNS c
WHERE TABLE_SCHEMA = DATABASE()
  AND (
      (COLUMN_NAME = 'tenant_id' AND TABLE_NAME != 'tenants')
      OR (COLUMN_NAME = 'bot_id' AND TABLE_NAME != 'bots')
      OR (COLUMN_NAME = 'user_id' AND TABLE_NAME != 'users')
      OR (COLUMN_NAME = 'conversation_id' AND TABLE_NAME != 'conversations')
  )
  AND NOT EXISTS (
      SELECT 1
      FROM information_schema.KEY_COLUMN_USAGE k
      WHERE k.TABLE_SCHEMA = c.TABLE_SCHEMA
        AND k.TABLE_NAME = c.TABLE_NAME
        AND k.COLUMN_NAME = c.COLUMN_NAME
        AND k.REFERENCED_TABLE_NAME IS NOT NULL
  )
ORDER BY TABLE_NAME, COLUMN_NAME;


-- ============================================================
-- 第四部分：字段类型优化验证
-- ============================================================

SELECT '=== 字段类型优化检查 ===' AS check_section;

-- 检查1：FLOAT类型字段（应使用DECIMAL）
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    DATA_TYPE,
    COLUMN_TYPE,
    '⚠️ 建议改为DECIMAL' AS recommendation
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND DATA_TYPE IN ('float', 'double')
  AND TABLE_NAME NOT LIKE '%test%'
ORDER BY TABLE_NAME, COLUMN_NAME;

-- 检查2：TEXT类型字段（评估是否可用VARCHAR）
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    DATA_TYPE,
    CHARACTER_MAXIMUM_LENGTH,
    CASE
        WHEN CHARACTER_MAXIMUM_LENGTH < 65535 THEN '⚠️ 可改为VARCHAR'
        ELSE '✅ 保留TEXT'
    END AS recommendation
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND DATA_TYPE = 'text'
  AND TABLE_NAME NOT LIKE '%test%'
ORDER BY TABLE_NAME, COLUMN_NAME;

-- 检查3：时间戳字段类型一致性
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    DATA_TYPE,
    CASE
        WHEN DATA_TYPE IN ('bigint', 'timestamp') THEN '✅ 使用标准类型'
        ELSE '⚠️ 建议统一类型'
    END AS type_status
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME IN ('created_at', 'updated_at', 'deleted_at')
  AND TABLE_NAME NOT LIKE '%test%'
ORDER BY TABLE_NAME, COLUMN_NAME;


-- ============================================================
-- 第五部分：性能测试验证
-- ============================================================

SELECT '=== 性能测试验证 ===' AS check_section;

-- 测试1：bots表查询性能
EXPLAIN SELECT *
FROM bots
WHERE tenant_id = 'test_tenant_id'
  AND status = 'active'
ORDER BY created_at DESC
LIMIT 20;

-- 预期结果：
-- type: ref 或 range
-- key: idx_tenant_status_created
-- rows: < 100
-- Extra: Using where; Using index（如果使用了覆盖索引）

-- 测试2：messages表查询性能
EXPLAIN SELECT *
FROM messages
WHERE tenant_id = 'test_tenant_id'
  AND conversation_id = 'test_conv_id'
  AND role = 'user'
ORDER BY created_at DESC
LIMIT 50;

-- 预期结果：
-- type: ref 或 range
-- key: idx_tenant_conv_role_created
-- rows: < 100

-- 测试3：knowledge_chunks全文搜索性能
EXPLAIN SELECT
    chunk_id,
    content,
    MATCH(content) AGAINST('测试关键词' IN NATURAL LANGUAGE MODE) AS score
FROM knowledge_chunks
WHERE MATCH(content) AGAINST('测试关键词' IN NATURAL LANGUAGE MODE)
ORDER BY score DESC
LIMIT 20;

-- 预期结果：
-- type: fulltext
-- key: ft_content
-- rows: < 100

-- 测试4：token_usage_logs租户成本查询
EXPLAIN SELECT
    tenant_id,
    DATE(created_at) as date,
    SUM(total_cost) as total_cost
FROM token_usage_logs
WHERE tenant_id = 'test_tenant_id'
  AND created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)
GROUP BY DATE(created_at)
ORDER BY date DESC;

-- 预期结果：
-- type: range
-- key: idx_tenant_cost_created


-- ============================================================
-- 第六部分：存储空间统计
-- ============================================================

SELECT '=== 存储空间统计 ===' AS check_section;

-- 统计1：表大小排行（TOP 10）
SELECT
    TABLE_NAME,
    ROUND(((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024), 2) AS total_size_mb,
    ROUND((DATA_LENGTH / 1024 / 1024), 2) AS data_size_mb,
    ROUND((INDEX_LENGTH / 1024 / 1024), 2) AS index_size_mb,
    TABLE_ROWS,
    ROUND(INDEX_LENGTH / DATA_LENGTH, 2) AS index_data_ratio
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME NOT LIKE '%migration%'
  AND TABLE_NAME NOT LIKE '%test%'
ORDER BY (DATA_LENGTH + INDEX_LENGTH) DESC
LIMIT 10;

-- 统计2：索引大小排行（TOP 10）
SELECT
    TABLE_NAME,
    INDEX_NAME,
    ROUND(STAT_VALUE * @@innodb_page_size / 1024 / 1024, 2) AS size_mb,
    ROUND(
        STAT_VALUE * @@innodb_page_size / 1024 / 1024 /
        (SELECT DATA_LENGTH / 1024 / 1024 FROM information_schema.TABLES
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = s.TABLE_NAME) * 100,
        2
    ) AS size_ratio_percent
FROM mysql.innodb_index_stats s
WHERE DATABASE_NAME = DATABASE()
  AND STAT_NAME = 'size'
  AND TABLE_NAME NOT LIKE '%migration%'
  AND TABLE_NAME NOT LIKE '%test%'
ORDER BY size_mb DESC
LIMIT 10;


-- ============================================================
-- 第七部分：优化总结报告
-- ============================================================

SELECT '=== 优化总结报告 ===' AS check_section;

-- 报告1：表命名规范
SELECT
    '表命名规范' AS category,
    COUNT(*) AS total_tables,
    SUM(CASE WHEN TABLE_NAME REGEXP '[A-Z]' THEN 1 ELSE 0 END) AS needs_improvement,
    CONCAT(
        ROUND(
            (1 - SUM(CASE WHEN TABLE_NAME REGEXP '[A-Z]' THEN 1 ELSE 0 END) / COUNT(*)) * 100,
            2
        ),
        '%'
    ) AS compliance_rate
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME NOT LIKE '%migration%'
  AND TABLE_NAME NOT LIKE '%test%';

-- 报告2：索引覆盖率
SELECT
    '索引覆盖率' AS category,
    COUNT(DISTINCT TABLE_NAME) AS tables_covered,
    (
        SELECT COUNT(DISTINCT TABLE_NAME)
        FROM information_schema.COLUMNS
        WHERE TABLE_SCHEMA = DATABASE()
          AND COLUMN_NAME = 'tenant_id'
    ) AS total_tables,
    (
        SELECT COUNT(DISTINCT TABLE_NAME)
        FROM information_schema.COLUMNS
        WHERE TABLE_SCHEMA = DATABASE()
          AND COLUMN_NAME = 'tenant_id'
          AND NOT EXISTS (
              SELECT 1 FROM information_schema.STATISTICS
              WHERE TABLE_SCHEMA = DATABASE()
                AND STATISTICS.TABLE_NAME = COLUMNS.TABLE_NAME
                AND INDEX_NAME = 'idx_tenant_id'
          )
    ) AS missing_indexes,
    CONCAT(
        ROUND(
            COUNT(DISTINCT TABLE_NAME) / (
                SELECT COUNT(DISTINCT TABLE_NAME)
                FROM information_schema.COLUMNS
                WHERE TABLE_SCHEMA = DATABASE()
                  AND COLUMN_NAME = 'tenant_id'
            ) * 100,
            2
        ),
        '%'
    ) AS coverage_rate
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND INDEX_NAME = 'idx_tenant_id';

-- 报告3：外键约束完整性
SELECT
    '外键约束完整性' AS category,
    COUNT(*) AS total_foreign_keys,
    COUNT(DISTINCT TABLE_NAME) AS tables_with_fk,
    (
        SELECT COUNT(DISTINCT TABLE_NAME)
        FROM information_schema.COLUMNS
        WHERE TABLE_SCHEMA = DATABASE()
          AND COLUMN_NAME = 'tenant_id'
          AND TABLE_NAME != 'tenants'
    ) AS tables_should_have_fk
FROM information_schema.KEY_COLUMN_USAGE
WHERE TABLE_SCHEMA = DATABASE()
  AND REFERENCED_TABLE_NAME IS NOT NULL;

-- 报告4：字段类型优化
SELECT
    '字段类型优化' AS category,
    SUM(CASE WHEN DATA_TYPE IN ('float', 'double') THEN 1 ELSE 0 END) AS float_columns,
    SUM(CASE WHEN DATA_TYPE = 'decimal' THEN 1 ELSE 0 END) AS decimal_columns,
    SUM(CASE WHEN DATA_TYPE = 'text' THEN 1 ELSE 0 END) AS text_columns,
    SUM(CASE WHEN DATA_TYPE LIKE 'varchar%' THEN 1 ELSE 0 END) AS varchar_columns
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME NOT LIKE '%migration%'
  AND TABLE_NAME NOT LIKE '%test%';


-- ============================================================
-- 第八部分：优化建议清单
-- ============================================================

SELECT '=== 优化建议清单 ===' AS check_section;

-- 建议1：缺失的索引
SELECT
    '缺失索引' AS issue_type,
    TABLE_NAME,
    COLUMN_NAME,
    '建议添加索引' AS recommendation
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND (
      (COLUMN_NAME = 'tenant_id' AND TABLE_NAME != 'tenants')
      OR (COLUMN_NAME = 'status' AND TABLE_NAME IN ('bots', 'conversations', 'workflows'))
      OR (COLUMN_NAME = 'created_at')
  )
  AND NOT EXISTS (
      SELECT 1
      FROM information_schema.STATISTICS
      WHERE TABLE_SCHEMA = DATABASE()
        AND STATISTICS.TABLE_NAME = COLUMNS.TABLE_NAME
        AND STATISTICS.COLUMN_NAME = COLUMNS.COLUMN_NAME
        AND (
            STATISTICS.INDEX_NAME LIKE 'idx_%'
            OR STATISTICS.INDEX_NAME LIKE 'uk_%'
            OR STATISTICS.INDEX_NAME LIKE 'ft_%'
        )
  )
ORDER BY TABLE_NAME, COLUMN_NAME;

-- 建议2：待优化的字段类型
SELECT
    '字段类型优化' AS issue_type,
    TABLE_NAME,
    COLUMN_NAME,
    DATA_TYPE AS current_type,
    CASE
        WHEN DATA_TYPE = 'float' THEN 'DECIMAL(10,2)'
        WHEN DATA_TYPE = 'double' THEN 'DECIMAL(15,4)'
        WHEN DATA_TYPE = 'text' AND CHARACTER_MAXIMUM_LENGTH < 65535 THEN 'VARCHAR(500)'
        ELSE '待评估'
    END AS recommended_type
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND (
      DATA_TYPE IN ('float', 'double')
      OR (DATA_TYPE = 'text' AND CHARACTER_MAXIMUM_LENGTH < 65535)
  )
  AND TABLE_NAME NOT LIKE '%test%'
ORDER BY TABLE_NAME, COLUMN_NAME;


-- ============================================================
-- 完成标记
-- ============================================================

-- 执行完成后，将以上所有结果保存为：
-- database_optimization_verification_result_[date].txt

-- 关键指标：
-- 1. 表命名规范合规率: _____ %
-- 2. 索引覆盖率: _____ %
-- 3. 外键约束完整性: _____ %
-- 4. 字段类型优化率: _____ %
-- 5. 查询性能提升: _____ %
-- 6. 存储空间节省: _____ MB


-- ============================================================
-- 注意事项
-- ============================================================

-- 1. 本验证脚本应分批执行，每部分单独验证
-- 2. 性能测试部分需要替换实际的tenant_id、conversation_id等参数
-- 3. 存储空间统计需要定期执行，监控增长趋势
-- 4. 优化建议清单应作为后续优化的输入
-- 5. 所有验证结果应归档保存，用于性能对比分析
