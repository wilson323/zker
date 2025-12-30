-- ================================================================================
-- 001_add_tenant_id_to_core_tables.sql
-- 为核心业务表添加 tenant_id 字段
-- ================================================================================
-- 用途: 从单租户迁移到多租户架构
-- 策略: 添加字段 -> 数据回填 -> 删除旧字段(可选)
-- ================================================================================

-- 设置安全模式
SET FOREIGN_KEY_CHECKS = 0;
SET SQL_MODE = '';

-- ================================================================================
-- 1. 核心表（高优先级）
-- ================================================================================

-- 1.1 用户表
ALTER TABLE `users`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 1.2 Bot相关表
ALTER TABLE `bot_draft`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`),
ADD INDEX `idx_tenant_updated` (`tenant_id`, `updated_at`);

ALTER TABLE `bot_published`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 1.3 对话相关表
ALTER TABLE `conversation`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_bot` (`tenant_id`, `bot_id`),
ADD INDEX `idx_tenant_created` (`tenant_id`, `created_at`);

ALTER TABLE `message`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`),
ADD INDEX `idx_tenant_conv` (`tenant_id`, `conversation_id`);

-- ================================================================================
-- 2. 知识库表
-- ================================================================================

ALTER TABLE `knowledge`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`);

ALTER TABLE `knowledge_document`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`);

ALTER TABLE `knowledge_document_slice`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`);

-- ================================================================================
-- 3. 工作流表
-- ================================================================================

ALTER TABLE `workflow`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`);

-- ================================================================================
-- 4. 文件表
-- ================================================================================

ALTER TABLE `files`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`);

-- ================================================================================
-- 5. 插件相关表
-- ================================================================================

ALTER TABLE `plugin_draft`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`);

ALTER TABLE `plugin_published`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`);

-- ================================================================================
-- 6. 数据库表
-- ================================================================================

ALTER TABLE `online_database_info`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`);

ALTER TABLE `draft_database_info`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`);

-- ================================================================================
-- 7. API认证表
-- ================================================================================

ALTER TABLE `api_key`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`);

-- ================================================================================
-- 8. 其他业务表
-- ================================================================================

ALTER TABLE `agent_to_database`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`);

ALTER TABLE `data_copy_task`
ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT 'system_tenant' COMMENT '租户ID' AFTER `id`,
ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 恢复安全模式
SET FOREIGN_KEY_CHECKS = 1;

-- ================================================================================
-- 验证脚本
-- ================================================================================

-- 检查所有表是否都已添加tenant_id字段
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    COLUMN_DEFAULT
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'tenant_id'
  AND TABLE_NAME IN (
    'users',
    'bot_draft',
    'bot_published',
    'conversation',
    'message',
    'knowledge',
    'knowledge_document',
    'knowledge_document_slice',
    'workflow',
    'files',
    'plugin_draft',
    'plugin_published',
    'online_database_info',
    'draft_database_info',
    'api_key',
    'agent_to_database',
    'data_copy_task'
  )
ORDER BY TABLE_NAME;

-- 检查索引是否已创建
SELECT
    TABLE_NAME,
    INDEX_NAME,
    COLUMN_NAME
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'tenant_id'
  AND TABLE_NAME IN (
    'users',
    'bot_draft',
    'bot_published',
    'conversation',
    'message',
    'knowledge',
    'knowledge_document',
    'knowledge_document_slice',
    'workflow',
    'files',
    'plugin_draft',
    'plugin_published',
    'online_database_info',
    'draft_database_info',
    'api_key',
    'agent_to_database',
    'data_copy_task'
  )
ORDER BY TABLE_NAME, INDEX_NAME;
