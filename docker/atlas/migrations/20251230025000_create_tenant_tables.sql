/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

-- ============================================================================
-- 租户系统核心表创建脚本
-- 版本: v1.0.0
-- 日期: 2025-01-01
-- 作者: ZKER Enterprise Team
-- 目的: 创建完整的租户管理系统表结构
-- ============================================================================

-- 1. 租户表（tenants）
CREATE TABLE IF NOT EXISTS `opencoze`.`tenants` (
    -- 主键
    `tenant_id` VARCHAR(36) PRIMARY KEY COMMENT '租户ID（UUID）',

    -- 基本信息
    `tenant_name` VARCHAR(200) NOT NULL COMMENT '租户名称',
    `tenant_type` ENUM('individual', 'team', 'enterprise') NOT NULL DEFAULT 'team' COMMENT '租户类型',
    `subdomain` VARCHAR(64) UNIQUE COMMENT '子域名',
    `status` ENUM('active', 'suspended', 'deleted') NOT NULL DEFAULT 'active' COMMENT '租户状态',

    -- 订阅信息
    `subscription_tier` ENUM('free', 'pro', 'enterprise') NOT NULL DEFAULT 'free' COMMENT '订阅等级',

    -- 联系信息
    `contact_name` VARCHAR(100) COMMENT '联系人姓名',
    `contact_email` VARCHAR(255) COMMENT '联系人邮箱',
    `contact_phone` VARCHAR(20) COMMENT '联系人电话',

    -- 配置信息
    `logo_url` VARCHAR(500) COMMENT '企业logo URL',
    `settings` JSON COMMENT '租户配置（JSON格式）',

    -- 审计字段
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（Unix时间戳）',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（Unix时间戳）',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间（软删除）',

    -- 索引
    INDEX `idx_tenant_type` (`tenant_type`),
    INDEX `idx_status` (`status`),
    INDEX `idx_subscription_tier` (`subscription_tier`),
    INDEX `idx_created_at` (`created_at`),
    UNIQUE INDEX `uk_tenant_name` (`tenant_name`, `deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户表';

-- 2. 订阅表（subscriptions）
CREATE TABLE IF NOT EXISTS `opencoze`.`subscriptions` (
    -- 主键
    `subscription_id` VARCHAR(36) PRIMARY KEY COMMENT '订阅ID（UUID）',

    -- 关联租户
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',

    -- 订阅信息
    `plan_tier` ENUM('free', 'pro', 'enterprise') NOT NULL COMMENT '订阅等级',
    `billing_cycle` ENUM('monthly', 'quarterly', 'yearly') NOT NULL DEFAULT 'monthly' COMMENT '计费周期',
    `status` ENUM('active', 'past_due', 'canceled', 'expired') NOT NULL DEFAULT 'active' COMMENT '订阅状态',

    -- 订阅周期
    `start_date` DATE NOT NULL COMMENT '订阅开始日期',
    `end_date` DATE COMMENT '订阅结束日期',

    -- 配额配置
    `quota_bots` INT NOT NULL DEFAULT 10 COMMENT 'Bot数量配额',
    `quota_messages_per_month` INT NOT NULL DEFAULT 1000 COMMENT '每月消息配额',
    `quota_storage_gb` INT NOT NULL DEFAULT 10 COMMENT '存储空间配额（GB）',
    `quota_team_members` INT NOT NULL DEFAULT 1 COMMENT '团队成员配额',

    -- 自动续费
    `auto_renew` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否自动续费',

    -- 审计字段
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间',

    -- 外键约束
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,

    -- 索引
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_plan_tier` (`plan_tier`),
    INDEX `idx_status` (`status`),
    INDEX `idx_end_date` (`end_date`),
    UNIQUE INDEX `uk_tenant_plan` (`tenant_id`, `plan_tier`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订阅表';

-- 3. 配额表（quotas）
CREATE TABLE IF NOT EXISTS `opencoze`.`quotas` (
    -- 主键
    `quota_id` VARCHAR(36) PRIMARY KEY COMMENT '配额ID（UUID）',

    -- 关联租户
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',

    -- 资源类型
    `resource_type` ENUM('bots', 'users', 'api_call', 'message', 'token', 'storage', 'knowledge_base', 'workflow') NOT NULL COMMENT '资源类型',

    -- 配额限制
    `max_limit` INT NOT NULL COMMENT '最大限制（-1表示无限制）',
    `used_count` INT NOT NULL DEFAULT 0 COMMENT '已使用数量',

    -- 重置周期
    `reset_cycle` ENUM('daily', 'weekly', 'monthly', 'yearly', 'never') NOT NULL DEFAULT 'monthly' COMMENT '重置周期',
    `last_reset_at` BIGINT DEFAULT 0 COMMENT '上次重置时间',

    -- 审计字段
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间',

    -- 外键约束
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,

    -- 索引
    INDEX `idx_tenant_resource` (`tenant_id`, `resource_type`),
    UNIQUE INDEX `uk_tenant_resource` (`tenant_id`, `resource_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='配额表';

-- 4. 配额使用日志表（quota_usage_log）
CREATE TABLE IF NOT EXISTS `opencoze`.`quota_usage_log` (
    -- 主键
    `log_id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '日志ID',

    -- 关联租户
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',

    -- 资源信息
    `resource_type` ENUM('bots', 'users', 'api_call', 'message', 'token', 'storage', 'knowledge_base', 'workflow') NOT NULL COMMENT '资源类型',
    `resource_id` VARCHAR(36) COMMENT '资源ID',

    -- 操作信息
    `action` ENUM('consume', 'rollback', 'reset') NOT NULL COMMENT '操作类型',
    `amount` INT NOT NULL COMMENT '数量（正数为消费，负数为回滚）',

    -- 配额信息
    `before_count` INT NOT NULL COMMENT '操作前数量',
    `after_count` INT NOT NULL COMMENT '操作后数量',

    -- 元数据
    `request_id` VARCHAR(36) COMMENT '请求ID',
    `user_id` VARCHAR(36) COMMENT '操作用户ID',
    `reason` VARCHAR(500) COMMENT '操作原因',

    -- 时间戳
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',

    -- 外键约束
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,

    -- 索引
    INDEX `idx_tenant_resource` (`tenant_id`, `resource_type`),
    INDEX `idx_created_at` (`created_at`),
    INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='配额使用日志表';

-- 5. 初始化系统默认租户
INSERT INTO `opencoze`.`tenants` (
    `tenant_id`, `tenant_name`, `tenant_type`, `status`,
    `subscription_tier`, `created_at`, `updated_at`
) VALUES
(
    'system-default',
    'System Default',
    'individual',
    'active',
    'free',
    UNIX_TIMESTAMP() * 1000,
    UNIX_TIMESTAMP() * 1000
) ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP() * 1000;
