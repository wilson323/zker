-- ============================================================================
-- ZKER 多租户迁移：model_instance 表添加 tenant_id 字段
-- ============================================================================
-- 版本: v1.0.0
-- 创建日期: 2025-12-30
-- 执行时间: 预计 2-5 分钟（取决于数据量）
-- 影响范围: model_instance 表
--
-- 功能说明:
-- 1. 为 model_instance 表添加 tenant_id VARCHAR(64) 字段
-- 2. 创建 idx_tenant_id 索引以提升查询性能
-- 3. 设置默认值为 'default' 以兼容现有数据
--
-- 注意事项:
-- ⚠️ 执行前请先备份数据库！
-- ⚠️ 建议在低峰期执行！
--
-- 回滚方案:
-- ALTER TABLE model_instance DROP COLUMN tenant_id;
-- ALTER TABLE model_instance DROP INDEX idx_model_instance_tenant_id;
-- ============================================================================

-- ============================================================================
-- 步骤1: 为 model_instance 表添加 tenant_id 字段
-- ============================================================================
ALTER TABLE `model_instance`
ADD COLUMN `tenant_id` VARCHAR(64) NOT NULL DEFAULT 'default'
    COMMENT '租户ID（多租户隔离）'
    AFTER `id`;

-- 验证字段添加
SHOW COLUMNS FROM `model_instance` LIKE 'tenant_id';

-- ============================================================================
-- 步骤2: 创建索引
-- ============================================================================

-- 创建单列索引
ALTER TABLE `model_instance`
ADD INDEX `idx_model_instance_tenant_id` (`tenant_id`);

-- 创建复合索引（用于租户隔离查询）
ALTER TABLE `model_instance`
ADD INDEX `idx_model_instance_tenant_id_deleted` (`tenant_id`, `deleted_at`);

-- 验证索引创建
SHOW INDEX FROM `model_instance` WHERE Key_name LIKE 'idx_model_instance_tenant%';

-- ============================================================================
-- 步骤3: 数据回填（使用默认租户ID）
-- ============================================================================
UPDATE `model_instance`
SET `tenant_id` = 'default'
WHERE `tenant_id` = 'default';

-- 验证数据完整性
SELECT
    COUNT(*) AS total_records,
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) AS null_tenant_count,
    SUM(CASE WHEN tenant_id = 'default' THEN 1 ELSE 0 END) AS default_tenant_count,
    COUNT(DISTINCT tenant_id) AS unique_tenant_count
FROM `model_instance`;

-- ============================================================================
-- 步骤4: 验证脚本
-- ============================================================================

-- 验证1: 检查字段是否成功添加
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    COLUMN_TYPE,
    IS_NULLABLE,
    COLUMN_DEFAULT,
    COLUMN_COMMENT
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'model_instance'
  AND COLUMN_NAME = 'tenant_id';

-- 验证2: 检查索引是否成功创建
SELECT
    TABLE_NAME,
    INDEX_NAME,
    COLUMN_NAME,
    SEQ_IN_INDEX,
    INDEX_TYPE
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'model_instance'
  AND INDEX_NAME LIKE '%tenant%'
ORDER BY INDEX_NAME, SEQ_IN_INDEX;

-- 验证3: 检查数据完整性
SELECT
    COUNT(*) AS total_records,
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) AS null_count,
    SUM(CASE WHEN tenant_id IS NOT NULL THEN 1 ELSE 0 END) AS has_tenant_count,
    ROUND(SUM(CASE WHEN tenant_id IS NOT NULL THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2) AS completion_rate
FROM `model_instance`;

-- ============================================================================
-- 步骤5: 性能优化
-- ============================================================================
ANALYZE TABLE `model_instance`;

-- ============================================================================
-- 回滚脚本（⚠️ 仅在迁移失败时使用）
-- ============================================================================
/*
ALTER TABLE `model_instance` DROP COLUMN tenant_id;
ALTER TABLE `model_instance` DROP INDEX idx_model_instance_tenant_id;
ALTER TABLE `model_instance` DROP INDEX idx_model_instance_tenant_id_deleted;
*/

-- ============================================================================
-- 执行完成标记
-- ============================================================================
-- 执行完成后，请在此处记录执行时间和结果
-- 执行时间: _______________
-- 执行结果: □ 成功  □ 失败
-- 备注说明: ___________________________________
-- 验证通过: □ 是  □ 否
-- ============================================================================
