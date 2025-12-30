-- ============================================================
-- 计费系统和权限系统数据库迁移脚本
-- 版本: v1.0
-- 日期: 2025-01-03
-- 说明: 创建10张核心表（权限系统5张 + 计费系统5张）
-- ============================================================

-- ============================================================
-- 第一部分：权限系统表（5张表）
-- 基于设计文档: 权限系统使用指南.md
-- ============================================================

-- 1. 角色表 (roles)
CREATE TABLE IF NOT EXISTS `roles` (
    `role_id` VARCHAR(36) PRIMARY KEY COMMENT '角色ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `role_name` VARCHAR(50) NOT NULL COMMENT '角色名称',
    `role_code` VARCHAR(50) NOT NULL COMMENT '角色编码',
    `description` VARCHAR(500) COMMENT '角色描述',

    -- 角色类型：system(系统预置) | custom(自定义)
    `role_type` VARCHAR(20) NOT NULL COMMENT '角色类型',

    -- 父角色ID（支持角色继承）
    `parent_role_id` VARCHAR(36) COMMENT '父角色ID',

    -- 是否启用
    `is_enabled` BOOLEAN NOT NULL DEFAULT TRUE COMMENT '是否启用',

    -- 是否系统预置角色
    `is_system` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否系统预置',

    -- 时间戳
    `created_by` VARCHAR(36) COMMENT '创建者ID',
    `updated_by` VARCHAR(36) COMMENT '更新者ID',
    `created_at` BIGINT NOT NULL COMMENT '创建时间（毫秒时间戳）',
    `updated_at` BIGINT NOT NULL COMMENT '更新时间（毫秒时间戳）',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间（软删除）',

    -- 索引
    INDEX `idx_tenant_role` (`tenant_id`, `role_code`),
    INDEX `idx_deleted_at` (`deleted_at`),
    UNIQUE KEY `uk_tenant_code` (`tenant_id`, `role_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 2. 权限表 (permissions)
CREATE TABLE IF NOT EXISTS `permissions` (
    `permission_id` VARCHAR(36) PRIMARY KEY COMMENT '权限ID',
    `permission_code` VARCHAR(100) NOT NULL UNIQUE COMMENT '权限编码',
    `name` VARCHAR(100) NOT NULL COMMENT '权限名称',
    `description` VARCHAR(500) COMMENT '权限描述',

    -- 权限类型：operation(操作) | data(数据) | field(字段)
    `type` VARCHAR(20) NOT NULL COMMENT '权限类型',

    -- 资源类型（如：bot, conversation, workflow等）
    `resource_type` VARCHAR(50) NOT NULL COMMENT '资源类型',

    -- 操作类型：create, read, update, delete, export, import, approve, reject
    `operation` VARCHAR(20) NOT NULL COMMENT '操作类型',

    -- 字段权限相关
    `field_name` VARCHAR(50) COMMENT '字段名称',
    `is_sensitive` BOOLEAN DEFAULT FALSE COMMENT '是否敏感字段',

    -- 系统预置权限
    `is_system` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否系统预置',

    -- 时间戳
    `created_at` BIGINT NOT NULL COMMENT '创建时间（毫秒时间戳）',
    `updated_at` BIGINT NOT NULL COMMENT '更新时间（毫秒时间戳）',

    -- 索引
    INDEX `idx_type_resource` (`type`, `resource_type`),
    INDEX `idx_operation` (`operation`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 3. 用户角色关联表 (user_roles)
CREATE TABLE IF NOT EXISTS `user_roles` (
    `user_role_id` VARCHAR(36) PRIMARY KEY COMMENT '用户角色关联ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `role_id` VARCHAR(36) NOT NULL COMMENT '角色ID',

    -- 授予者信息
    `granted_by` VARCHAR(36) COMMENT '授予者ID',
    `granted_at` BIGINT NOT NULL COMMENT '授予时间（毫秒时间戳）',

    -- 过期时间（可选）
    `expires_at` BIGINT DEFAULT NULL COMMENT '过期时间（毫秒时间戳）',

    -- 是否启用
    `is_enabled` BOOLEAN DEFAULT TRUE COMMENT '是否启用',

    -- 时间戳
    `created_at` BIGINT NOT NULL COMMENT '创建时间（毫秒时间戳）',

    -- 索引
    INDEX `idx_tenant_user` (`tenant_id`, `user_id`),
    INDEX `idx_role` (`role_id`),
    INDEX `idx_expires_at` (`expires_at`),
    UNIQUE KEY `uk_tenant_user_role` (`tenant_id`, `user_id`, `role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 4. 数据权限表 (data_permissions)
CREATE TABLE IF NOT EXISTS `data_permissions` (
    `permission_id` VARCHAR(36) PRIMARY KEY COMMENT '权限ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `role_id` VARCHAR(36) NOT NULL COMMENT '角色ID',
    `permission_name` VARCHAR(100) NOT NULL COMMENT '权限名称',
    `permission_code` VARCHAR(100) NOT NULL UNIQUE COMMENT '权限编码',
    `description` VARCHAR(500) COMMENT '权限描述',

    -- 数据权限范围类型：all, department, team, own, custom
    `scope` VARCHAR(20) NOT NULL COMMENT '数据权限范围',

    -- 自定义过滤条件（JSON格式）
    `custom_filter` TEXT COMMENT '自定义过滤条件',

    -- 关联的资源类型
    `resource_type` VARCHAR(50) NOT NULL COMMENT '资源类型',

    -- 是否启用
    `is_enabled` BOOLEAN DEFAULT TRUE COMMENT '是否启用',

    -- 时间戳
    `created_by` VARCHAR(36) COMMENT '创建者ID',
    `updated_by` VARCHAR(36) COMMENT '更新者ID',
    `created_at` BIGINT NOT NULL COMMENT '创建时间（毫秒时间戳）',
    `updated_at` BIGINT NOT NULL COMMENT '更新时间（毫秒时间戳）',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间（软删除）',

    -- 索引
    INDEX `idx_tenant` (`tenant_id`),
    INDEX `idx_role` (`role_id`),
    INDEX `idx_resource_type` (`resource_type`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据权限表';

-- 5. 字段权限表 (field_permissions)
CREATE TABLE IF NOT EXISTS `field_permissions` (
    `permission_id` VARCHAR(36) PRIMARY KEY COMMENT '权限ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `role_id` VARCHAR(36) NOT NULL COMMENT '角色ID',
    `permission_name` VARCHAR(100) NOT NULL COMMENT '权限名称',
    `permission_code` VARCHAR(100) NOT NULL UNIQUE COMMENT '权限编码',
    `description` VARCHAR(500) COMMENT '权限描述',

    -- 资源类型
    `resource_type` VARCHAR(50) NOT NULL COMMENT '资源类型',

    -- 字段名称
    `field_name` VARCHAR(50) NOT NULL COMMENT '字段名称',

    -- 权限级别：editable(可编辑), readonly(只读), hidden(隐藏)
    `permission_level` VARCHAR(20) NOT NULL COMMENT '权限级别',

    -- 是否启用
    `is_enabled` BOOLEAN DEFAULT TRUE COMMENT '是否启用',

    -- 时间戳
    `created_by` VARCHAR(36) COMMENT '创建者ID',
    `updated_by` VARCHAR(36) COMMENT '更新者ID',
    `created_at` BIGINT NOT NULL COMMENT '创建时间（毫秒时间戳）',
    `updated_at` BIGINT NOT NULL COMMENT '更新时间（毫秒时间戳）',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间（软删除）',

    -- 索引
    INDEX `idx_tenant` (`tenant_id`),
    INDEX `idx_role` (`role_id`),
    INDEX `idx_resource_type` (`resource_type`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='字段权限表';

-- ============================================================
-- 第二部分：计费系统表（5张表）
-- 基于设计文档: 23-租户计费系统_TokenMetering补充_完整版.md
-- ============================================================

-- 6. Token使用明细日志表 (token_usage_logs)
CREATE TABLE IF NOT EXISTS `token_usage_logs` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',

    -- 租户与用户信息
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `user_id` BIGINT COMMENT '用户ID',

    -- 关联对象信息
    `bot_id` VARCHAR(64) COMMENT 'Bot ID',
    `conversation_id` VARCHAR(64) COMMENT '会话ID',
    `message_id` VARCHAR(64) COMMENT '消息ID',

    -- Token统计
    `input_tokens` INT NOT NULL COMMENT '输入Token数',
    `output_tokens` INT NOT NULL COMMENT '输出Token数',
    `total_tokens` INT NOT NULL COMMENT '总Token数',

    -- 模型信息
    `model_provider` VARCHAR(50) NOT NULL COMMENT '模型提供商(openai/anthropic/通义千问)',
    `model_name` VARCHAR(50) NOT NULL COMMENT '模型名称(gpt-4/claude-3-opus/qwen-max)',
    `model_version` VARCHAR(20) COMMENT '模型版本',

    -- 成本计算
    `unit_price` DECIMAL(10,6) NOT NULL COMMENT '每1K Token单价(CNY)',
    `input_cost` DECIMAL(10,6) NOT NULL COMMENT '输入成本',
    `output_cost` DECIMAL(10,6) NOT NULL COMMENT '输出成本',
    `total_cost` DECIMAL(10,6) NOT NULL COMMENT '总成本',

    -- 性能指标
    `response_time_ms` INT COMMENT '响应时间(毫秒)',
    `latency_ms` INT COMMENT '首字延迟(毫秒)',
    `is_cached` BOOLEAN DEFAULT FALSE COMMENT '是否缓存命中',

    -- 元数据
    `request_type` VARCHAR(20) NOT NULL COMMENT '请求类型(chat/completion/embedding/rerank)',
    `metadata` JSON COMMENT '额外元数据',

    -- 时间戳
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    -- 索引
    INDEX `idx_tenant_created` (`tenant_id`, `created_at`),
    INDEX `idx_bot_created` (`bot_id`, `created_at`),
    INDEX `idx_model` (`model_provider`, `model_name`),
    INDEX `idx_cost` (`total_cost`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Token使用明细日志表';

-- 7. Token使用汇总表 (token_usage_summary)
CREATE TABLE IF NOT EXISTS `token_usage_summary` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '汇总ID',

    -- 分组维度
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `bot_id` VARCHAR(64) COMMENT 'Bot ID (NULL表示整体汇总)',
    `summary_date` DATE NOT NULL COMMENT '汇总日期',
    `summary_hour` TINYINT COMMENT '小时 (0-23, NULL表示日汇总)',

    -- Token汇总
    `total_input_tokens` BIGINT NOT NULL DEFAULT 0 COMMENT '总输入Token数',
    `total_output_tokens` BIGINT NOT NULL DEFAULT 0 COMMENT '总输出Token数',
    `total_tokens` BIGINT NOT NULL DEFAULT 0 COMMENT '总Token数',

    -- 成本汇总
    `total_cost` DECIMAL(12,6) NOT NULL DEFAULT 0.000000 COMMENT '总成本',

    -- 请求统计
    `total_requests` INT NOT NULL DEFAULT 0 COMMENT '总请求数',
    `cached_requests` INT NOT NULL DEFAULT 0 COMMENT '缓存命中请求数',
    `avg_response_time` DECIMAL(8,2) COMMENT '平均响应时间(ms)',

    -- 模型分布
    `model_distribution` JSON COMMENT '模型使用分布 {model_name: {tokens, cost}}',

    -- 时间戳
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    -- 唯一索引
    UNIQUE KEY `uk_tenant_bot_date_hour` (`tenant_id`, `bot_id`, `summary_date`, `summary_hour`),

    -- 索引
    INDEX `idx_tenant_date` (`tenant_id`, `summary_date`),
    INDEX `idx_bot_date` (`bot_id`, `summary_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Token使用汇总表(按天/小时)';

-- 8. 预算配置表 (budget_settings)
CREATE TABLE IF NOT EXISTS `budget_settings` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '配置ID',
    `tenant_id` VARCHAR(64) NOT NULL UNIQUE COMMENT '租户ID',

    -- 预算配置
    `budget_type` VARCHAR(20) NOT NULL DEFAULT 'monthly' COMMENT '预算周期(monthly/quarterly/yearly)',
    `budget_amount` DECIMAL(12,2) NOT NULL COMMENT '预算金额(CNY)',
    `currency` VARCHAR(3) DEFAULT 'CNY' COMMENT '货币',

    -- 告警阈值
    `alert_threshold_1` INT DEFAULT 80 COMMENT '一级告警阈值(%)',
    `alert_threshold_2` INT DEFAULT 95 COMMENT '二级告警阈值(%)',
    `hard_cap_enabled` BOOLEAN DEFAULT FALSE COMMENT '是否启用硬性上限',
    `hard_cap_amount` DECIMAL(12,2) COMMENT '硬性上限金额',

    -- 降级策略
    `auto_downgrade_enabled` BOOLEAN DEFAULT FALSE COMMENT '是否启用自动降级',
    `downgrade_threshold` INT DEFAULT 90 COMMENT '降级阈值(%)',
    `original_model` VARCHAR(100) COMMENT '原模型',
    `fallback_model` VARCHAR(100) COMMENT '降级模型',

    -- 通知配置
    `notification_channels` JSON COMMENT '通知渠道 ["email","sms","webhook"]',
    `notification_recipients` JSON COMMENT '通知接收人列表',

    -- 时间戳
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    -- 索引
    INDEX `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预算配置表';

-- 9. 预算告警历史表 (budget_alerts_history)
CREATE TABLE IF NOT EXISTS `budget_alerts_history` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '告警ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',

    -- 告警类型
    `alert_type` VARCHAR(30) NOT NULL COMMENT '告警类型(threshold_1/threshold_2/hard_cap/downgrade)',

    -- 预算状态
    `budget_amount` DECIMAL(12,2) NOT NULL COMMENT '预算金额',
    `used_amount` DECIMAL(12,2) NOT NULL COMMENT '已使用金额',
    `usage_percent` DECIMAL(5,2) NOT NULL COMMENT '使用率(%)',

    -- 告警信息
    `alert_level` VARCHAR(20) NOT NULL COMMENT '告警级别(warning/critical/emergency)',
    `alert_message` TEXT COMMENT '告警消息',

    -- 发送状态
    `notification_channels` JSON COMMENT '通知渠道',
    `notification_status` VARCHAR(20) DEFAULT 'pending' COMMENT '通知状态(pending/sent/failed)',
    `sent_at` DATETIME COMMENT '发送时间',
    `error_message` TEXT COMMENT '错误信息',

    -- 时间戳
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    -- 索引
    INDEX `idx_tenant_created` (`tenant_id`, `created_at`),
    INDEX `idx_alert_type` (`alert_type`),
    INDEX `idx_status` (`notification_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预算告警历史表';

-- 10. 成本优化建议表 (cost_optimization_suggestions)
CREATE TABLE IF NOT EXISTS `cost_optimization_suggestions` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '建议ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',

    -- 建议类型
    `suggestion_type` VARCHAR(50) NOT NULL COMMENT '建议类型(model_downgrade/enable_cache/batch_request/prompt_optimization)',

    -- 分析数据
    `analysis_period_start` DATE NOT NULL COMMENT '分析开始日期',
    `analysis_period_end` DATE NOT NULL COMMENT '分析结束日期',

    -- 建议内容
    `target_bot_id` VARCHAR(64) COMMENT '目标Bot ID',
    `current_model` VARCHAR(100) COMMENT '当前模型',
    `suggested_model` VARCHAR(100) COMMENT '建议模型',
    `reason` TEXT COMMENT '优化原因',

    -- 预期效果
    `estimated_monthly_saving` DECIMAL(12,2) COMMENT '预计月节省金额',
    `estimated_saving_percent` DECIMAL(5,2) COMMENT '预计节省比例(%)',

    -- 状态
    `status` VARCHAR(20) DEFAULT 'pending' COMMENT '状态(pending/approved/rejected/applied)',
    `applied_at` DATETIME COMMENT '应用时间',
    `applied_by` BIGINT COMMENT '应用操作人ID',

    -- 时间戳
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    -- 索引
    INDEX `idx_tenant_status` (`tenant_id`, `status`),
    INDEX `idx_type` (`suggestion_type`),
    INDEX `idx_saving` (`estimated_monthly_saving`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='成本优化建议表';

-- ============================================================
-- 迁移完成
-- ============================================================
