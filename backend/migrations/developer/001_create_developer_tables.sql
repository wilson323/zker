-- Developer Platform Database Migration
-- Copyright 2025 coze-dev Authors

-- 创建开发者表
CREATE TABLE IF NOT EXISTS `developers` (
    `developer_id` VARCHAR(36) PRIMARY KEY COMMENT '开发者ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `developer_name` VARCHAR(100) NOT NULL COMMENT '开发者名称',
    `developer_email` VARCHAR(100) NOT NULL COMMENT '开发者邮箱',
    `status` ENUM('active', 'suspended', 'deleted') DEFAULT 'active' NOT NULL COMMENT '状态',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间（毫秒时间戳）',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_deleted_at` (`deleted_at`),
    UNIQUE KEY `uk_tenant_user` (`tenant_id`, `user_id`, `deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='开发者表';

-- 创建项目表
CREATE TABLE IF NOT EXISTS `developer_projects` (
    `project_id` VARCHAR(36) PRIMARY KEY COMMENT '项目ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `developer_id` VARCHAR(36) NOT NULL COMMENT '开发者ID',
    `project_name` VARCHAR(100) NOT NULL COMMENT '项目名称',
    `project_type` ENUM('bot', 'workflow', 'integration') NOT NULL COMMENT '项目类型',
    `description` TEXT COMMENT '项目描述',
    `status` ENUM('development', 'production', 'archived') DEFAULT 'development' NOT NULL COMMENT '状态',
    `config` JSON COMMENT '项目配置',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间（毫秒时间戳）',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_developer_id` (`developer_id`),
    INDEX `idx_type` (`project_type`),
    INDEX `idx_status` (`status`),
    INDEX `idx_deleted_at` (`deleted_at`),
    UNIQUE KEY `uk_tenant_name` (`tenant_id`, `project_name`, `deleted_at`),

    FOREIGN KEY (`developer_id`) REFERENCES `developers`(`developer_id`) ON DELETE CASCADE,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenants`(`tenant_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='开发者项目表';

-- 创建API密钥表
CREATE TABLE IF NOT EXISTS `developer_api_keys` (
    `key_id` VARCHAR(36) PRIMARY KEY COMMENT '密钥ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `project_id` VARCHAR(36) NOT NULL COMMENT '项目ID',
    `key_name` VARCHAR(100) NOT NULL COMMENT '密钥名称',
    `key_secret` VARCHAR(255) NOT NULL COMMENT '密钥（加密存储）',
    `key_prefix` ENUM('sk', 'pk', 'tk') NOT NULL COMMENT '密钥前缀',
    `key_masked` VARCHAR(50) NOT NULL COMMENT '脱敏密钥（用于显示）',
    `scopes` JSON NOT NULL COMMENT '权限范围',
    `expires_at` BIGINT DEFAULT NULL COMMENT '过期时间（毫秒时间戳）',
    `last_used_at` BIGINT DEFAULT NULL COMMENT '最后使用时间（毫秒时间戳）',
    `status` ENUM('active', 'revoked', 'expired') DEFAULT 'active' NOT NULL COMMENT '状态',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间（毫秒时间戳）',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_project_id` (`project_id`),
    INDEX `idx_expires_at` (`expires_at`),
    INDEX `idx_last_used` (`last_used_at`),
    INDEX `idx_status` (`status`),
    INDEX `idx_deleted_at` (`deleted_at`),
    UNIQUE KEY `uk_key_secret` (`key_secret`),

    FOREIGN KEY (`project_id`) REFERENCES `developer_projects`(`project_id`) ON DELETE CASCADE,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenants`(`tenant_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='开发者API密钥表';

-- 创建Webhook表
CREATE TABLE IF NOT EXISTS `developer_webhooks` (
    `webhook_id` VARCHAR(36) PRIMARY KEY COMMENT 'Webhook ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `project_id` VARCHAR(36) NOT NULL COMMENT '项目ID',
    `webhook_url` VARCHAR(500) NOT NULL COMMENT 'Webhook URL',
    `webhook_secret` VARCHAR(100) NOT NULL COMMENT 'Webhook密钥（用于签名验证）',
    `events` JSON NOT NULL COMMENT '事件类型',
    `status` ENUM('active', 'paused') DEFAULT 'active' NOT NULL COMMENT '状态',
    `last_triggered_at` BIGINT DEFAULT NULL COMMENT '最后触发时间（毫秒时间戳）',
    `success_count` BIGINT DEFAULT 0 COMMENT '成功次数',
    `failure_count` BIGINT DEFAULT 0 COMMENT '失败次数',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间（毫秒时间戳）',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_project_id` (`project_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_deleted_at` (`deleted_at`),

    FOREIGN KEY (`project_id`) REFERENCES `developer_projects`(`project_id`) ON DELETE CASCADE,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenants`(`tenant_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='开发者Webhook表';

-- 创建Webhook日志表
CREATE TABLE IF NOT EXISTS `developer_webhook_logs` (
    `log_id` VARCHAR(36) PRIMARY KEY COMMENT '日志ID',
    `webhook_id` VARCHAR(36) NOT NULL COMMENT 'Webhook ID',
    `event_type` VARCHAR(50) NOT NULL COMMENT '事件类型',
    `status_code` INT COMMENT 'HTTP状态码',
    `response` TEXT COMMENT '响应内容',
    `duration_ms` BIGINT COMMENT '耗时（毫秒）',
    `success` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否成功',
    `retry_count` INT DEFAULT 0 COMMENT '重试次数',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',

    INDEX `idx_webhook_id` (`webhook_id`),
    INDEX `idx_event_type` (`event_type`),
    INDEX `idx_created_at` (`created_at`),

    FOREIGN KEY (`webhook_id`) REFERENCES `developer_webhooks`(`webhook_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Webhook日志表';

-- 创建SDK表
CREATE TABLE IF NOT EXISTS `developer_sdks` (
    `sdk_id` VARCHAR(36) PRIMARY KEY COMMENT 'SDK ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `project_id` VARCHAR(36) NOT NULL COMMENT '项目ID',
    `sdk_name` VARCHAR(100) NOT NULL COMMENT 'SDK名称',
    `language` ENUM('python', 'javascript', 'go', 'java') NOT NULL COMMENT '编程语言',
    `version` VARCHAR(20) NOT NULL COMMENT '版本号',
    `description` TEXT COMMENT '描述',
    `code` LONGTEXT COMMENT '生成的代码',
    `package_url` VARCHAR(500) COMMENT '包下载地址',
    `readme` TEXT COMMENT '使用文档',
    `status` ENUM('draft', 'published', 'archived') DEFAULT 'draft' NOT NULL COMMENT '状态',
    `download_count` BIGINT DEFAULT 0 COMMENT '下载次数',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间（毫秒时间戳）',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_project_id` (`project_id`),
    INDEX `idx_language` (`language`),
    INDEX `idx_status` (`status`),
    INDEX `idx_deleted_at` (`deleted_at`),

    FOREIGN KEY (`project_id`) REFERENCES `developer_projects`(`project_id`) ON DELETE CASCADE,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenants`(`tenant_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='开发者SDK表';
