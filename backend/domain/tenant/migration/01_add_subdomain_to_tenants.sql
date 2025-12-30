-- ================================================================================
-- Coze Studio - 数据库迁移脚本
-- ================================================================================
-- 迁移名称: 添加子域名字段到 tenants 表
-- 迁移版本: v1.0.1
-- 创建时间: 2025-01-01
-- 作者: ZKER Backend Team
-- 描述: 支持多租户子域名路由，添加 subdomain 字段到 tenants 表
-- ================================================================================

-- 检查迁移表是否存在（用于幂等性）
CREATE TABLE IF NOT EXISTS `migrations` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `migration_name` VARCHAR(255) NOT NULL COMMENT '迁移名称',
    `executed_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '执行时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_migration_name` (`migration_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据库迁移记录表';

-- ================================================================================
-- Step 1: 添加 subdomain 列（如果不存在）
-- ================================================================================

SET @dbname = DATABASE();
SET @tablename = 'tenants';
SET @columnname = 'subdomain';
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = @dbname
    AND TABLE_NAME = @tablename
    AND COLUMN_NAME = @columnname
  ) > 0,
  'SELECT 1',
  CONCAT('ALTER TABLE ', @tablename, ' ADD COLUMN ', @columnname, ' VARCHAR(64) NOT NULL UNIQUE COMMENT ''租户子域名（用于子域名路由，如 tenant-a.saas.com）''')
));

PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

-- ================================================================================
-- Step 2: 为现有租户生成默认子域名（如果 subdomain 为空）
-- ================================================================================

-- 使用 tenant_id 的前缀作为默认子域名（如果子域名不存在）
-- 规则: tenant_id 转小写，移除连字符
UPDATE tenants
SET subdomain = LOWER(SUBSTRING_INDEX(REPLACE(tenant_id, '-', ''), '_', 64))
WHERE subdomain IS NULL OR subdomain = '';

-- ================================================================================
-- Step 3: 创建索引（如果不存在）
-- ================================================================================

-- 创建唯一索引（如果不存在）
CREATE INDEX IF NOT EXISTS `uk_tenant_subdomain` ON tenants(subdomain);

-- ================================================================================
-- Step 4: 记录迁移
-- ================================================================================

INSERT INTO migrations (migration_name, executed_at)
VALUES ('add_subdomain_to_tenants', NOW())
ON DUPLICATE KEY UPDATE executed_at = NOW();

-- ================================================================================
-- 验证迁移
-- ================================================================================

SELECT
    'Migration completed successfully!' AS message,
    COUNT(*) AS total_tenants,
    SUM(CASE WHEN subdomain IS NOT NULL THEN 1 ELSE 0 END) AS tenants_with_subdomain
FROM tenants;

-- ================================================================================
-- 回滚脚本（如需回滚，请执行以下命令）
-- ================================================================================

-- -- 删除 subdomain 列
-- ALTER TABLE tenants DROP COLUMN subdomain;
--
-- -- 删除迁移记录
-- DELETE FROM migrations WHERE migration_name = 'add_subdomain_to_tenants';
--
-- ================================================================================
