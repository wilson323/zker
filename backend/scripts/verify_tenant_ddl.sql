-- ============================================================================
-- 租户表DDL验证脚本
-- 版本: v1.0.0
-- 日期: 2025-01-01
-- 目的: 验证租户系统表结构完整性
-- ============================================================================

-- 1. 验证表结构
SELECT '========== Table Structure Verification ==========' as '';
SELECT 'Table: tenants' as check_item;
SHOW CREATE TABLE `opencoze`.`tenants`;

SELECT 'Table: subscriptions' as check_item;
SHOW CREATE TABLE `opencoze`.`subscriptions`;

SELECT 'Table: quotas' as check_item;
SHOW CREATE TABLE `opencoze`.`quotas`;

SELECT 'Table: quota_usage_log' as check_item;
SHOW CREATE TABLE `opencoze`.`quota_usage_log`;

-- 2. 验证索引
SELECT '========== Index Verification ==========' as '';
SELECT 'Indexes: tenants' as check_item;
SHOW INDEX FROM `opencoze`.`tenants`;

SELECT 'Indexes: subscriptions' as check_item;
SHOW INDEX FROM `opencoze`.`subscriptions`;

SELECT 'Indexes: quotas' as check_item;
SHOW INDEX FROM `opencoze`.`quotas`;

SELECT 'Indexes: quota_usage_log' as check_item;
SHOW INDEX FROM `opencoze`.`quota_usage_log`;

-- 3. 验证外键
SELECT '========== Foreign Key Verification ==========' as '';
SELECT
    TABLE_NAME,
    CONSTRAINT_NAME,
    REFERENCED_TABLE_NAME,
    REFERENCED_COLUMN_NAME
FROM
    information_schema.KEY_COLUMN_USAGE
WHERE
    TABLE_SCHEMA = 'opencoze'
    AND REFERENCED_TABLE_NAME IS NOT NULL
    AND TABLE_NAME IN ('tenants', 'subscriptions', 'quotas', 'quota_usage_log')
ORDER BY
    TABLE_NAME, CONSTRAINT_NAME;

-- 4. 验证触发器
SELECT '========== Trigger Verification ==========' as '';
SHOW TRIGGERS LIKE 'trg_%';

-- 5. 验证视图
SELECT '========== View Verification ==========' as '';
SELECT 'View: v_tenant_statistics' as check_item;
SHOW CREATE VIEW `opencoze`.`v_tenant_statistics`;

SELECT 'View: v_quota_usage_rate' as check_item;
SHOW CREATE VIEW `opencoze`.`v_quota_usage_rate`;

SELECT 'View: v_subscription_expiry_alert' as check_item;
SHOW CREATE VIEW `opencoze`.`v_subscription_expiry_alert`;

-- 6. 数据完整性检查
SELECT '========== Data Integrity Check ==========' as '';
SELECT 'Data check: system tenant' as check_item;
SELECT * FROM `opencoze`.`tenants` WHERE `tenant_id` = 'system-default';

SELECT 'Data check: system subscription' as check_item;
SELECT * FROM `opencoze`.`subscriptions` WHERE `tenant_id` = 'system-default';

SELECT 'Data check: system quotas' as check_item;
SELECT * FROM `opencoze`.`quotas` WHERE `tenant_id` = 'system-default';

-- 7. 统计信息
SELECT '========== Statistics Summary ==========' as '';
SELECT
    (SELECT COUNT(*) FROM `opencoze`.`tenants`) as tenants_count,
    (SELECT COUNT(*) FROM `opencoze`.`subscriptions`) as subscriptions_count,
    (SELECT COUNT(*) FROM `opencoze`.`quotas`) as quotas_count,
    (SELECT COUNT(*) FROM `opencoze`.`quota_usage_log`) as usage_logs_count;

-- 8. 字段类型验证
SELECT '========== Column Type Verification ==========' as '';
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    DATA_TYPE,
    IS_NULLABLE,
    COLUMN_DEFAULT,
    COLUMN_COMMENT
FROM
    information_schema.COLUMNS
WHERE
    TABLE_SCHEMA = 'opencoze'
    AND TABLE_NAME IN ('tenants', 'subscriptions', 'quotas', 'quota_usage_log')
ORDER BY
    TABLE_NAME, ORDINAL_POSITION;

-- 9. 枚举类型验证
SELECT '========== ENUM Type Verification ==========' as '';
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    COLUMN_TYPE
FROM
    information_schema.COLUMNS
WHERE
    TABLE_SCHEMA = 'opencoze'
    AND TABLE_NAME IN ('tenants', 'subscriptions', 'quotas', 'quota_usage_log')
    AND DATA_TYPE = 'enum'
ORDER BY
    TABLE_NAME, COLUMN_NAME;

-- 10. 性能检查（索引使用情况）
SELECT '========== Performance Check ==========' as '';
SELECT
    TABLE_NAME,
    INDEX_NAME,
    SEQ_IN_INDEX,
    COLUMN_NAME,
    CARDINALITY
FROM
    information_schema.STATISTICS
WHERE
    TABLE_SCHEMA = 'opencoze'
    AND TABLE_NAME IN ('tenants', 'subscriptions', 'quotas', 'quota_usage_log')
ORDER BY
    TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX;
