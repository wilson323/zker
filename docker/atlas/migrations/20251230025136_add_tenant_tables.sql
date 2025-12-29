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

-- Create "tenants" table
CREATE TABLE `opencoze`.`tenants` (
    `tenant_id` VARCHAR(36) PRIMARY KEY COMMENT '租户ID',
    `tenant_name` VARCHAR(200) NOT NULL COMMENT '租户名称',
    `tenant_type` ENUM('individual', 'team', 'enterprise') NOT NULL COMMENT '租户类型',
    `status` ENUM('active', 'suspended', 'deleted') DEFAULT 'active' COMMENT '状态',
    `subscription_tier` ENUM('free', 'pro', 'enterprise') DEFAULT 'free' COMMENT '订阅等级',
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒）',
    `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒）',
    `deleted_at` BIGINT UNSIGNED NULL COMMENT '删除时间（毫秒）',
    UNIQUE INDEX `uk_tenant_name` (`tenant_name`, `deleted_at`),
    INDEX `idx_status_type` (`status`, `tenant_type`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户表';

-- Create "subscriptions" table
CREATE TABLE `opencoze`.`subscriptions` (
    `subscription_id` VARCHAR(36) PRIMARY KEY COMMENT '订阅ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `plan_tier` ENUM('free', 'pro', 'enterprise') NOT NULL COMMENT '套餐等级',
    `billing_cycle` ENUM('monthly', 'yearly') NOT NULL COMMENT '计费周期',
    `start_date` DATE NOT NULL COMMENT '开始日期',
    `end_date` DATE NULL COMMENT '结束日期',
    `auto_renew` BOOLEAN DEFAULT TRUE COMMENT '自动续费',
    `status` ENUM('active', 'expired', 'cancelled') DEFAULT 'active' COMMENT '状态',
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒）',
    `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒）',
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_status_end_date` (`status`, `end_date`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订阅表';

-- Create "quotas" table
CREATE TABLE `opencoze`.`quotas` (
    `quota_id` VARCHAR(36) PRIMARY KEY COMMENT '配额ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `resource_type` ENUM('bots', 'messages', 'storage', 'team_members') NOT NULL COMMENT '资源类型',
    `max_limit` INT NOT NULL COMMENT '最大限制',
    `used_count` INT DEFAULT 0 COMMENT '已使用数量',
    `reset_cycle` ENUM('daily', 'monthly', 'yearly', 'never') DEFAULT 'monthly' COMMENT '重置周期',
    `last_reset_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '上次重置时间（毫秒）',
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒）',
    `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒）',
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,
    UNIQUE INDEX `uk_tenant_resource` (`tenant_id`, `resource_type`),
    INDEX `idx_tenant_id` (`tenant_id`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='配额表';
