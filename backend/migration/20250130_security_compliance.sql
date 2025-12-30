-- ============================================================================
-- 企业级安全合规系统 - MySQL Migration Script
-- 版本: v1.0.0
-- 日期: 2025-01-30
-- 说明: 创建审计日志、MFA多因素认证、数据脱敏相关表
-- 合规: 等保三级、GDPR、SOC2
-- ============================================================================

-- ============================================================================
-- 1. 审计日志表 (audit_logs)
-- 用途: 记录所有敏感操作，满足等保三级和GDPR审计要求
-- ============================================================================
CREATE TABLE IF NOT EXISTS `audit_logs` (
    -- 主键
    `log_id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '日志ID',

    -- 租户和用户信息
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `username` VARCHAR(100) NOT NULL COMMENT '用户名',
    `user_email` VARCHAR(255) NOT NULL COMMENT '用户邮箱',

    -- 操作信息
    `action` VARCHAR(50) NOT NULL COMMENT '操作类型(user.login, user.delete, role.create等)',
    `resource` VARCHAR(50) NOT NULL COMMENT '资源类型(user, role, bot, knowledge等)',
    `resource_id` VARCHAR(36) DEFAULT NULL COMMENT '资源ID',
    `resource_name` VARCHAR(255) DEFAULT NULL COMMENT '资源名称',

    -- 请求信息
    `request_method` VARCHAR(10) DEFAULT NULL COMMENT 'HTTP方法(GET, POST, PUT, DELETE等)',
    `request_path` VARCHAR(500) DEFAULT NULL COMMENT '请求路径',
    `request_ip` VARCHAR(45) DEFAULT NULL COMMENT '客户端IP(IPv4/IPv6)',
    `user_agent` VARCHAR(500) DEFAULT NULL COMMENT 'User-Agent',

    -- 状态和结果
    `status` ENUM('success', 'failed', 'pending') NOT NULL COMMENT '操作状态',
    `error_code` VARCHAR(50) DEFAULT NULL COMMENT '错误码',
    `error_msg` TEXT DEFAULT NULL COMMENT '错误消息',

    -- 详细信息(已脱敏)
    `request_data` LONGTEXT DEFAULT NULL COMMENT '请求数据(敏感数据已脱敏)',
    `response_data` LONGTEXT DEFAULT NULL COMMENT '响应数据(敏感数据已脱敏)',
    `changes` LONGTEXT DEFAULT NULL COMMENT '变更内容(JSON格式)',
    `metadata` JSON DEFAULT NULL COMMENT '额外元数据',

    -- 合规字段
    `session_id` VARCHAR(36) DEFAULT NULL COMMENT '会话ID',
    `trace_id` VARCHAR(36) DEFAULT NULL COMMENT '追踪ID(用于分布式追踪)',

    -- 时间戳
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `archived_at` TIMESTAMP NULL DEFAULT NULL COMMENT '归档时间',

    -- 安全字段
    `signature` VARCHAR(64) DEFAULT NULL COMMENT '数字签名(SHA-256,用于防篡改)',
    `is_archived` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已归档',

    -- 索引
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_action` (`action`),
    INDEX `idx_resource` (`resource`),
    INDEX `idx_resource_id` (`resource_id`),
    INDEX `idx_request_path` (`request_path`),
    INDEX `idx_request_ip` (`request_ip`),
    INDEX `idx_status` (`status`),
    INDEX `idx_session_id` (`session_id`),
    INDEX `idx_trace_id` (`trace_id`),
    INDEX `idx_created_at` (`created_at`),
    INDEX `idx_is_archived` (`is_archived`),
    INDEX `idx_signature` (`signature`),

    -- 组合索引(常用查询)
    INDEX `idx_tenant_created` (`tenant_id`, `created_at`),
    INDEX `idx_user_action` (`user_id`, `action`),
    INDEX `idx_resource_status` (`resource`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='审计日志表-满足等保三级和GDPR审计要求';

-- ============================================================================
-- 2. MFA配置表 (mfa_configs)
-- 用途: 存储多因素认证配置
-- ============================================================================
CREATE TABLE IF NOT EXISTS `mfa_configs` (
    -- 主键
    `config_id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '配置ID',

    -- 租户和用户信息
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',

    -- MFA配置
    `type` ENUM('totp', 'sms', 'email', 'hardware') NOT NULL COMMENT 'MFA类型',
    `status` ENUM('enabled', 'disabled', 'locked') NOT NULL DEFAULT 'disabled' COMMENT 'MFA状态',
    `secret` VARCHAR(255) NOT NULL COMMENT 'TOTP密钥(加密存储)',
    `backup_codes` TEXT DEFAULT NULL COMMENT '恢复码(加密存储,JSON数组)',

    -- TOTP配置
    `issuer` VARCHAR(100) DEFAULT NULL COMMENT '发行者名称',
    `account` VARCHAR(100) DEFAULT NULL COMMENT '账户名称',
    `algorithm` VARCHAR(20) NOT NULL DEFAULT 'SHA1' COMMENT '算法(SHA1, SHA256, SHA512)',
    `digits` INT NOT NULL DEFAULT 6 COMMENT '位数(6或8)',
    `period` INT NOT NULL DEFAULT 30 COMMENT '时间步长(秒)',

    -- 备用验证方式
    `backup_type` ENUM('totp', 'sms', 'email', 'hardware') DEFAULT NULL COMMENT '备用MFA类型',
    `backup_value` VARCHAR(255) DEFAULT NULL COMMENT '备用方式值(如手机号、邮箱)',

    -- 统计信息
    `used_count` INT NOT NULL DEFAULT 0 COMMENT '使用次数',
    `last_used_at` TIMESTAMP NULL DEFAULT NULL COMMENT '最后使用时间',
    `failed_attempts` INT NOT NULL DEFAULT 0 COMMENT '失败次数',
    `locked_until` TIMESTAMP NULL DEFAULT NULL COMMENT '锁定到期时间',

    -- 时间戳
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间(软删除)',

    -- 索引
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_deleted_at` (`deleted_at`),

    -- 唯一约束(每个用户只能有一个启用的MFA配置)
    UNIQUE KEY `uk_tenant_user` (`tenant_id`, `user_id`, `deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='MFA配置表-多因素认证';

-- ============================================================================
-- 3. 恢复码表 (recovery_codes)
-- 用途: 存储MFA恢复码(当用户无法访问 Authenticator App 时使用)
-- ============================================================================
CREATE TABLE IF NOT EXISTS `recovery_codes` (
    -- 主键
    `code_id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '恢复码ID',

    -- 关联信息
    `config_id` VARCHAR(36) NOT NULL COMMENT 'MFA配置ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',

    -- 恢复码(哈希存储)
    `code_hash` VARCHAR(255) NOT NULL COMMENT '恢复码哈希值(SHA-256)',

    -- 使用状态
    `used` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已使用',
    `used_at` TIMESTAMP NULL DEFAULT NULL COMMENT '使用时间',

    -- 时间戳
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `expires_at` TIMESTAMP NULL DEFAULT NULL COMMENT '过期时间',

    -- 索引
    INDEX `idx_config_id` (`config_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_used` (`used`),
    INDEX `idx_expires_at` (`expires_at`),

    -- 唯一约束(每个恢复码只能使用一次)
    UNIQUE KEY `uk_code_hash` (`code_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='MFA恢复码表-用于紧急情况下的身份验证';

-- ============================================================================
-- 4. MFA验证记录表 (mfa_verifications)
-- 用途: 记录所有MFA验证尝试，用于安全分析
-- ============================================================================
CREATE TABLE IF NOT EXISTS `mfa_verifications` (
    -- 主键
    `verify_id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '验证ID',

    -- 租户和用户信息
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',

    -- 验证信息
    `type` ENUM('totp', 'sms', 'email', 'hardware') NOT NULL COMMENT 'MFA类型',
    `success` TINYINT(1) NOT NULL COMMENT '是否成功',
    `failure_reason` VARCHAR(255) DEFAULT NULL COMMENT '失败原因',

    -- 请求信息
    `request_ip` VARCHAR(45) DEFAULT NULL COMMENT '客户端IP',
    `user_agent` VARCHAR(500) DEFAULT NULL COMMENT 'User-Agent',

    -- 时间戳
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '验证时间',

    -- 索引
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_created_at` (`created_at`),
    INDEX `idx_success` (`success`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='MFA验证记录表-用于安全分析';

-- ============================================================================
-- 5. MFA登录会话表 (mfa_login_sessions)
-- 用途: 支持多步骤MFA登录流程
-- ============================================================================
CREATE TABLE IF NOT EXISTS `mfa_login_sessions` (
    -- 主键
    `session_id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '会话ID',

    -- 租户和用户信息
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',

    -- 第一步认证(密码/邮箱)
    `first_factor_passed` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '第一步是否通过',
    `first_factor_at` TIMESTAMP NULL DEFAULT NULL COMMENT '第一步通过时间',

    -- 第二步认证(MFA)
    `second_factor_passed` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '第二步是否通过',
    `second_factor_at` TIMESTAMP NULL DEFAULT NULL COMMENT '第二步通过时间',

    -- 会话状态
    `completed` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否完成',
    `expires_at` TIMESTAMP NOT NULL COMMENT '过期时间',

    -- 时间戳
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    -- 索引
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_completed` (`completed`),
    INDEX `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='MFA登录会话表-支持多步骤认证';

-- ============================================================================
-- 6. 数据脱敏配置表 (data_masking_rules) - 可选
-- 用途: 存储自定义数据脱敏规则
-- ============================================================================
CREATE TABLE IF NOT EXISTS `data_masking_rules` (
    -- 主键
    `rule_id` VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '规则ID',

    -- 租户信息
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',

    -- 规则配置
    `rule_name` VARCHAR(100) NOT NULL COMMENT '规则名称',
    `field_pattern` VARCHAR(255) NOT NULL COMMENT '字段匹配模式(正则表达式)',
    `masking_level` ENUM('none', 'low', 'medium', 'high') NOT NULL DEFAULT 'medium' COMMENT '脱敏级别',
    `masking_method` VARCHAR(50) NOT NULL COMMENT '脱敏方法(email, phone, idcard等)',

    -- 是否启用
    `enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',

    -- 时间戳
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    -- 索引
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_enabled` (`enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='数据脱敏规则表-自定义脱敏配置';

-- ============================================================================
-- 初始化数据: 插入默认脱敏规则
-- ============================================================================
INSERT INTO `data_masking_rules` (`rule_id`, `tenant_id`, `rule_name`, `field_pattern`, `masking_level`, `masking_method`) VALUES
('rule_email_default', 'system', '邮箱脱敏', '(?i)email|mail', 'medium', 'email'),
('rule_phone_default', 'system', '手机号脱敏', '(?i)phone|mobile|tel', 'medium', 'phone'),
('rule_idcard_default', 'system', '身份证脱敏', '(?i)idcard|id_card', 'high', 'idcard'),
('rule_bankcard_default', 'system', '银行卡脱敏', '(?i)bank|card', 'high', 'bankcard'),
('rule_password_default', 'system', '密码脱敏', '(?i)password|passwd|secret', 'high', 'password')
ON DUPLICATE KEY UPDATE `updated_at` = CURRENT_TIMESTAMP;

-- ============================================================================
-- 性能优化: 创建分区表(可选,用于大数据量场景)
-- ============================================================================

-- 审计日志按月分区(可选)
-- ALTER TABLE `audit_logs` PARTITION BY RANGE (TO_DAYS(created_at)) (
--     PARTITION p202501 VALUES LESS THAN (TO_DAYS('2025-02-01')),
--     PARTITION p202502 VALUES LESS THAN (TO_DAYS('2025-03-01')),
--     PARTITION p202503 VALUES LESS THAN (TO_DAYS('2025-04-01')),
--     -- 继续添加更多分区...
--     PARTITION p_max VALUES LESS THAN MAXVALUE
-- );

-- ============================================================================
-- 安全加固建议
-- ============================================================================

-- 1. 启用MySQL审计插件
-- INSTALL PLUGIN audit_log SONAME 'audit_log.so';
-- SET GLOBAL audit_log_policy = 'ALL';

-- 2. 启用二进制日志(用于数据恢复)
-- SET GLOBAL log_bin = ON;
-- SET GLOBAL binlog_format = 'ROW';

-- 3. 设置自动清理过期审计日志(例如保留180天)
-- CREATE EVENT evt_cleanup_audit_logs
-- ON SCHEDULE EVERY 1 DAY
-- DO DELETE FROM audit_logs WHERE created_at < DATE_SUB(NOW(), INTERVAL 180 DAY);

-- 4. 定期归档审计日志(例如每月归档一次)
-- CREATE EVENT evt_archive_audit_logs
-- ON SCHEDULE EVERY 1 MONTH
-- DO CALL archive_audit_logs(DATE_SUB(NOW(), INTERVAL 90 DAY));

-- ============================================================================
-- 验证脚本
-- ============================================================================

-- 验证表是否创建成功
SELECT
    TABLE_NAME,
    TABLE_COMMENT,
    CREATE_TIME
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
AND TABLE_NAME IN ('audit_logs', 'mfa_configs', 'recovery_codes', 'mfa_verifications', 'mfa_login_sessions', 'data_masking_rules')
ORDER BY TABLE_NAME;

-- 验证索引是否创建成功
SELECT
    TABLE_NAME,
    INDEX_NAME,
    COLUMN_NAME,
    SEQ_IN_INDEX
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
AND TABLE_NAME IN ('audit_logs', 'mfa_configs', 'recovery_codes', 'mfa_verifications', 'mfa_login_sessions')
ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX;

-- ============================================================================
-- 回滚脚本(谨慎使用!)
-- ============================================================================
-- DROP TABLE IF EXISTS `audit_logs`;
-- DROP TABLE IF EXISTS `mfa_configs`;
-- DROP TABLE IF EXISTS `recovery_codes`;
-- DROP TABLE IF EXISTS `mfa_verifications`;
-- DROP TABLE IF EXISTS `mfa_login_sessions`;
-- DROP TABLE IF EXISTS `data_masking_rules`;
