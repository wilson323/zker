-- =====================================================
-- ZKER 迁移审计表 (migration_audit)
-- 版本: v1.0.0
-- 说明: 记录所有数据迁移的执行状态和结果
-- =====================================================

CREATE TABLE IF NOT EXISTS `migration_audit` (
    -- 主键
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',

    -- 迁移标识
    `migration_id` VARCHAR(64) NOT NULL COMMENT '迁移批次ID（格式：MIGRATE_YYYYMMDD_HHMMSS）',
    `migration_name` VARCHAR(200) NOT NULL COMMENT '迁移名称',
    `migration_type` ENUM('schema', 'data', 'backfill', 'cutover', 'verification') NOT NULL COMMENT '迁移类型',

    -- 状态信息
    `status` ENUM('pending', 'running', 'completed', 'failed', 'rolled_back') NOT NULL DEFAULT 'pending' COMMENT '迁移状态',
    `start_time` BIGINT NOT NULL COMMENT '开始时间（Unix毫秒时间戳）',
    `end_time` BIGINT DEFAULT NULL COMMENT '结束时间（Unix毫秒时间戳）',
    `duration_seconds` BIGINT DEFAULT NULL COMMENT '耗时（秒）',

    -- 处理统计
    `total_records` BIGINT DEFAULT NULL COMMENT '总记录数',
    `processed_records` BIGINT DEFAULT NULL COMMENT '已处理记录数',
    `failed_records` BIGINT DEFAULT NULL COMMENT '失败记录数',

    -- 错误信息
    `error_message` TEXT DEFAULT NULL COMMENT '错误信息',

    -- 执行信息
    `executed_by` VARCHAR(100) NOT NULL COMMENT '执行人',
    `metadata` JSON DEFAULT NULL COMMENT '元数据（JSON格式）',

    -- 审计字段
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（Unix毫秒时间戳）',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（Unix毫秒时间戳）',

    -- 索引
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_migration_id` (`migration_id`),
    KEY `idx_status` (`status`),
    KEY `idx_type` (`migration_type`),
    KEY `idx_start_time` (`start_time`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_status_start_time` (`status`, `start_time`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='迁移审计表：记录所有数据迁移的执行状态和结果';

-- =====================================================
-- 初始化数据示例（可选）
-- =====================================================

-- 示例：插入一条迁移记录
-- INSERT INTO `migration_audit` (
--     `migration_id`,
--     `migration_name`,
--     `migration_type`,
--     `status`,
--     `start_time`,
--     `total_records`,
--     `executed_by`,
--     `created_at`,
--     `updated_at`
-- ) VALUES (
--     'MIGRATE_20250130_120000',
--     'tenant_id字段添加和回填',
--     'schema',
--     'pending',
--     UNIX_TIMESTAMP(NOW()) * 1000,
--     1000000,
--     'system',
--     UNIX_TIMESTAMP(NOW()) * 1000,
--     UNIX_TIMESTAMP(NOW()) * 1000
-- );
