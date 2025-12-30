-- =====================================================
-- 配置中心数据库表初始化
-- 功能：多租户配置管理、版本控制、变更历史
-- 作者：企业级功能完善项目组
-- 日期：2025-01-01
-- =====================================================

-- =====================================================
-- 1. 配置项表 (config_items)
-- =====================================================
CREATE TABLE IF NOT EXISTS `config_items` (
    -- 主键
    `config_id` VARCHAR(36) PRIMARY KEY COMMENT '配置ID',

    -- 租户关联
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',

    -- 配置信息
    `config_key` VARCHAR(255) NOT NULL COMMENT '配置键',
    `config_value` TEXT NOT NULL COMMENT '配置值',
    `config_type` VARCHAR(50) NOT NULL DEFAULT 'string' COMMENT '配置类型：string, int, float, bool, json',

    -- 描述信息
    `description` TEXT COMMENT '配置描述',
    `is_encrypted` BOOLEAN DEFAULT FALSE COMMENT '是否加密存储',

    -- 审计字段
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `updated_by` VARCHAR(255) COMMENT '更新人',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间（软删除）',

    -- 索引
    INDEX `idx_tenant_key` (`tenant_id`, `config_key`),
    UNIQUE KEY `uk_tenant_key` (`tenant_id`, `config_key`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='配置项表 - 多租户配置管理';

-- =====================================================
-- 2. 配置变更历史表 (config_history)
-- =====================================================
CREATE TABLE IF NOT EXISTS `config_history` (
    -- 主键
    `history_id` VARCHAR(36) PRIMARY KEY COMMENT '历史记录ID',

    -- 配置关联
    `config_id` VARCHAR(36) NOT NULL COMMENT '配置ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `config_key` VARCHAR(255) NOT NULL COMMENT '配置键',

    -- 变更信息
    `old_value` TEXT COMMENT '旧值',
    `new_value` TEXT COMMENT '新值',
    `change_reason` TEXT COMMENT '变更原因',
    `changed_by` VARCHAR(255) NOT NULL COMMENT '变更人',
    `changed_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '变更时间',

    -- 版本控制
    `version_number` INT NOT NULL COMMENT '版本号',

    -- 软删除
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间（软删除）',

    -- 索引
    INDEX `idx_config_id` (`config_id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_changed_at` (`changed_at`),
    INDEX `idx_version_number` (`version_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='配置变更历史表 - 版本控制和审计';

-- =====================================================
-- 3. 初始化示例配置数据
-- =====================================================
-- 注意：这里是示例数据，实际使用时应该通过API创建配置

-- 系统级默认配置（租户ID为 'system'）
INSERT INTO `config_items` (`config_id`, `tenant_id`, `config_key`, `config_value`, `config_type`, `description`, `is_encrypted`) VALUES
('config_sys_001', 'system', 'feature.flag.smart_routing', 'true', 'bool', '智能路由功能开关', FALSE),
('config_sys_002', 'system', 'feature.flag.multi_tenant', 'true', 'bool', '多租户功能开关', FALSE),
('config_sys_003', 'system', 'routing.health_check_interval', '30', 'int', '健康检查间隔（秒）', FALSE),
('config_sys_004', 'system', 'routing.circuit_breaker_threshold', '5', 'int', '熔断器失败阈值', FALSE),
('config_sys_005', 'system', 'monitoring.metrics_retention_days', '30', 'int', '监控指标保留天数', FALSE)
ON DUPLICATE KEY UPDATE `config_value` = VALUES(`config_value`);

-- =====================================================
-- 4. 权限建议
-- =====================================================
-- 建议创建只读用户用于监控和审计：
-- CREATE USER 'config_reader'@'%' IDENTIFIED BY 'password';
-- GRANT SELECT ON config_items, config_history TO 'config_reader'@'%';

-- =====================================================
-- 5. 数据验证查询
-- =====================================================
-- 检查表是否创建成功
SELECT
    TABLE_NAME,
    TABLE_COMMENT,
    CREATE_TIME
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME IN ('config_items', 'config_history')
ORDER BY TABLE_NAME;

-- 检查示例数据
SELECT
    tenant_id,
    config_key,
    config_value,
    config_type,
    description,
    created_at
FROM config_items
WHERE deleted_at IS NULL
ORDER BY tenant_id, config_key;
