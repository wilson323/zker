-- ============================================================
-- 字段类型优化脚本 - 第三批
-- 目的：优化字段类型，提升存储效率和查询性能
-- 执行时机：维护窗口（凌晨2:00-4:00）
-- 预计时长：30-60分钟
-- 风险等级：中（需要数据迁移，需先备份）
-- ============================================================

-- ============================================================
-- 前置检查
-- ============================================================

-- 检查表是否存在
SELECT TABLE_NAME, TABLE_ROWS
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME IN ('bot_store_item', 'tenant_metrics', 'agent_metrics');

-- 备份提醒
-- ⚠️ 重要：执行前请先备份相关表！
-- mysqldump -u root -p opencoze bot_store_item > bot_store_item_backup.sql
-- mysqldump -u root -p opencoze tenant_metrics > tenant_metrics_backup.sql
-- mysqldump -u root -p opencoze agent_metrics > agent_metrics_backup.sql


-- ============================================================
-- 1. bot_store_item表 - 字段类型优化
-- ============================================================

-- 优化1：price字段 FLOAT -> DECIMAL(10,2)
-- 原因：FLOAT存在精度丢失，DECIMAL适合存储金额
ALTER TABLE bot_store_item
MODIFY COLUMN price DECIMAL(10,2) NOT NULL COMMENT '价格（保留2位小数）';

-- 优化2：description字段 TEXT -> VARCHAR(500)
-- 原因：VARCHAR可以创建索引，性能更好
ALTER TABLE bot_store_item
MODIFY COLUMN description VARCHAR(500) COMMENT '描述';

-- 优化3：reject_reason字段 TEXT -> VARCHAR(500)
ALTER TABLE bot_store_item
MODIFY COLUMN reject_reason VARCHAR(500) COMMENT '驳回原因';


-- ============================================================
-- 2. 监控表时间戳统一（TIMESTAMP -> BIGINT）
-- ============================================================

-- ⚠️ 注意：监控表使用BIGINT毫秒时间戳，保持与业务表一致

-- tenant_metrics表已使用BIGINT，无需修改

-- agent_metrics表已使用BIGINT，无需修改

-- 如果有其他表使用TIMESTAMP，按以下方式修改：

-- 步骤1：添加新字段（BIGINT）
-- ALTER TABLE example_table
-- ADD COLUMN created_at_big BIGINT COMMENT '创建时间（毫秒）';

-- 步骤2：数据迁移
-- UPDATE example_table
-- SET created_at_big = UNIX_TIMESTAMP(created_at) * 1000;

-- 步骤3：验证数据
-- SELECT
--     created_at,
--     created_at_big,
--     FROM_UNIXTIME(created_at_big / 1000) AS converted_back
-- FROM example_table
-- LIMIT 10;

-- 步骤4：删除旧字段（谨慎！）
-- ALTER TABLE example_table DROP COLUMN created_at;

-- 步骤5：重命名新字段
-- ALTER TABLE example_table
-- CHANGE created_at_big created_at BIGINT NOT NULL COMMENT '创建时间（毫秒）';


-- ============================================================
-- 3. JSON字段生成列优化（可选）
-- ============================================================

-- bot_store_item.tags字段：添加生成列用于索引
-- 原因：JSON字段无法直接索引，生成列可以

-- 步骤1：添加生成列（提取category）
ALTER TABLE bot_store_item
ADD COLUMN category_extracted VARCHAR(50)
    AS (JSON_UNQUOTE(JSON_EXTRACT(category, '$'))) STORED
    COMMENT '类别提取（用于索引）';

-- 步骤2：添加索引
CREATE INDEX idx_category ON bot_store_item(category_extracted);

-- 步骤3：查询示例
-- 旧查询（无法使用索引）
-- SELECT * FROM bot_store_item WHERE JSON_CONTAINS(category, '"AI"');

-- 新查询（可以使用索引）
-- SELECT * FROM bot_store_item WHERE category_extracted = 'AI';


-- ============================================================
-- 4. rating字段优化（FLOAT -> DECIMAL）
-- ============================================================

ALTER TABLE bot_store_item
MODIFY COLUMN rating DECIMAL(3,2) DEFAULT 0.00 COMMENT '评分（0.00-5.00）';

-- 说明：
-- - DECIMAL(3,2): 总共3位，2位小数，范围0.00-9.99
-- - 适合存储0-5星评分


-- ============================================================
-- 5. view_count和download_count字段优化（INT -> MEDIUMINT）
-- ============================================================

-- 如果数据量可控（< 16,777,215），使用MEDIUMINT节省空间
ALTER TABLE bot_store_item
MODIFY COLUMN view_count MEDIUMINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '浏览次数';

ALTER TABLE bot_store_item
MODIFY COLUMN download_count MEDIUMINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '下载次数';


-- ============================================================
-- 验证脚本
-- ============================================================

-- 1. 验证字段类型修改成功
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    DATA_TYPE,
    COLUMN_TYPE,
    IS_NULLABLE,
    COLUMN_DEFAULT,
    COLUMN_COMMENT
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'bot_store_item'
  AND COLUMN_NAME IN (
      'price',
      'description',
      'reject_reason',
      'rating',
      'view_count',
      'download_count',
      'category_extracted'
  )
ORDER BY ORDINAL_POSITION;

-- 2. 验证数据完整性
SELECT
    COUNT(*) AS total_rows,
    COUNT(CASE WHEN price IS NULL THEN 1 END) AS null_price,
    COUNT(CASE WHEN rating < 0 OR rating > 5 THEN 1 END) AS invalid_rating,
    MIN(price) AS min_price,
    MAX(price) AS max_price
FROM bot_store_item;

-- 预期结果：
-- - null_price = 0（price字段无NULL值）
-- - invalid_rating = 0（rating值在0-5范围内）


-- ============================================================
-- 性能测试
-- ============================================================

-- 测试1：价格范围查询（DECIMAL索引）
EXPLAIN SELECT *
FROM bot_store_item
WHERE category_extracted = 'AI'
  AND price BETWEEN 10 AND 100
ORDER BY rating DESC
LIMIT 20;

-- 预期结果：
-- - type: ref 或 range
-- - key: idx_category
-- - rows: < 100

-- 测试2：评分排序查询
EXPLAIN SELECT *
FROM bot_store_item
WHERE category_extracted = 'AI'
  AND rating >= 4.0
ORDER BY rating DESC, created_at DESC
LIMIT 20;


-- ============================================================
-- 存储空间对比
-- ============================================================

-- 查看优化前后的表大小
SELECT
    TABLE_NAME,
    ROUND(((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024), 2) AS total_size_mb,
    ROUND((DATA_LENGTH / 1024 / 1024), 2) AS data_size_mb,
    ROUND((INDEX_LENGTH / 1024 / 1024), 2) AS index_size_mb,
    TABLE_ROWS
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'bot_store_item';


-- ============================================================
-- 回滚脚本（仅在优化失败时使用）
-- ============================================================

/*
-- ⚠️ 警告：回滚将丢失优化效果

-- 回滚bot_store_item表字段类型
ALTER TABLE bot_store_item
MODIFY COLUMN price FLOAT COMMENT '价格';

ALTER TABLE bot_store_item
MODIFY COLUMN description TEXT COMMENT '描述';

ALTER TABLE bot_store_item
MODIFY COLUMN reject_reason TEXT COMMENT '驳回原因';

ALTER TABLE bot_store_item
MODIFY COLUMN rating FLOAT DEFAULT 0 COMMENT '评分';

ALTER TABLE bot_store_item
MODIFY COLUMN view_count INT NOT NULL DEFAULT 0 COMMENT '浏览次数';

ALTER TABLE bot_store_item
MODIFY COLUMN download_count INT NOT NULL DEFAULT 0 COMMENT '下载次数';

-- 删除生成列和索引
ALTER TABLE bot_store_item DROP COLUMN category_extracted;
DROP INDEX idx_category ON bot_store_item;
*/


-- ============================================================
-- 注意事项
-- ============================================================

-- 1. 数据类型选择原则：
--    - 金额：DECIMAL(M, D)（精确计算）
--    - 整数：根据范围选择 TINYINT、SMALLINT、MEDIUMINT、INT、BIGINT
--    - 字符串：根据长度选择 VARCHAR(N)（N < 65535）或 TEXT
--    - 布尔：BOOLEAN（TINYINT(1)）
--    - 时间：BIGINT（毫秒时间戳）或 TIMESTAMP（自动更新）

-- 2. 性能影响：
--    - DECIMAL vs FLOAT: DECIMAL精度高，性能略低（可忽略）
--    - VARCHAR vs TEXT: VARCHAR可索引，性能更好
--    - 生成列: 占用额外存储空间，但查询性能提升显著

-- 3. 存储空间：
--    - DECIMAL(10,2): 5字节
--    - FLOAT: 4字节（但精度有损失）
--    - VARCHAR(N): 实际长度+1或+2字节
--    - TEXT: 实际长度+2字节（不可索引）

-- 4. 兼容性：
--    - 修改字段类型可能导致应用层代码报错
--    - 需要同步修改ORM模型定义
--    - 需要充分测试业务功能


-- ============================================================
-- 完成标记
-- ============================================================
-- 执行完成后，请在此处记录：
-- 执行时间: _______________
-- 执行结果: □ 成功  □ 失败
-- 优化字段数量: ______ 个
-- 节省存储空间: _____ MB
-- 性能提升: ______ %
-- 备注说明: ___________________________________
-- ============================================================
