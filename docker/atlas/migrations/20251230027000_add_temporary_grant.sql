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
-- 临时授权功能：添加临时授权支持
-- 版本: v1.0.0
-- 日期: 2025-01-01
-- 目的: 支持临时权限授予和自动过期
-- ============================================================================

-- 1. 为user_roles表添加expires_at字段
ALTER TABLE `opencoze`.`user_roles`
ADD COLUMN `expires_at` BIGINT UNSIGNED NULL COMMENT '过期时间（毫秒时间戳，NULL表示永久）' AFTER `created_at`,
ADD INDEX `idx_expires_at` (`expires_at`);

-- 2. 创建临时授权码表（temporary_grants）
CREATE TABLE IF NOT EXISTS `opencoze`.`temporary_grants` (
    -- 主键
    `grant_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '授权ID',

    -- 授权码信息
    `grant_code` VARCHAR(100) UNIQUE NOT NULL COMMENT '授权码（一次性使用）',

    -- 授权信息
    `grantee_id` VARCHAR(36) NOT NULL COMMENT '被授权人用户ID',
    `grantor_id` VARCHAR(36) NOT NULL COMMENT '授权人用户ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',

    -- 权限类型
    `permission_type` ENUM('role', 'data_permission', 'field_permission') NOT NULL COMMENT '权限类型',
    `permission_data` JSON NOT NULL COMMENT '权限数据（JSON格式）',

    -- 有效期
    `expires_at` BIGINT UNSIGNED NOT NULL COMMENT '过期时间（毫秒时间戳）',

    -- 状态
    `is_used` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否已使用',
    `is_revoked` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否已撤销',
    `used_at` BIGINT UNSIGNED NULL COMMENT '使用时间（毫秒时间戳）',
    `revoked_at` BIGINT UNSIGNED NULL COMMENT '撤销时间（毫秒时间戳）',

    -- 元数据
    `reason` VARCHAR(500) NULL COMMENT '授权原因',
    `request_id` VARCHAR(36) NULL COMMENT '关联请求ID',

    -- 审计字段
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
    `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',

    -- 主键定义
    PRIMARY KEY (`grant_id`),

    -- 外键约束
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,
    FOREIGN KEY (`grantee_id`) REFERENCES `opencoze`.`user`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`grantor_id`) REFERENCES `opencoze`.`user`(`id`) ON DELETE CASCADE,

    -- 唯一索引
    UNIQUE INDEX `uk_grant_code` (`grant_code`),

    -- 普通索引
    INDEX `idx_grantee_id` (`grantee_id`),
    INDEX `idx_grantor_id` (`grantor_id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_expires_at` (`expires_at`),
    INDEX `idx_is_used` (`is_used`),
    INDEX `idx_is_revoked` (`is_revoked`),
    INDEX `idx_status_expires` (`is_used`, `is_revoked`, `expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='临时授权表';

-- 3. 创建临时授权历史表（temporary_grant_history）
CREATE TABLE IF NOT EXISTS `opencoze`.`temporary_grant_history` (
    -- 主键
    `history_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '历史ID',

    -- 关联授权
    `grant_id` BIGINT UNSIGNED NOT NULL COMMENT '临时授权ID',
    `grant_code` VARCHAR(100) NOT NULL COMMENT '授权码',

    -- 操作信息
    `action_type` ENUM('created', 'used', 'expired', 'revoked') NOT NULL COMMENT '操作类型',
    `operator_id` VARCHAR(36) NOT NULL COMMENT '操作人ID',
    `operator_name` VARCHAR(100) NULL COMMENT '操作人姓名',

    -- 快照数据
    `snapshot_data` JSON NULL COMMENT '操作时的权限数据快照',

    -- 元数据
    `reason` VARCHAR(500) NULL COMMENT '操作原因',
    `request_id` VARCHAR(36) NULL COMMENT '关联请求ID',

    -- 时间戳
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',

    -- 主键定义
    PRIMARY KEY (`history_id`),

    -- 外键约束
    FOREIGN KEY (`grant_id`) REFERENCES `opencoze`.`temporary_grants`(`grant_id`) ON DELETE CASCADE,

    -- 索引
    INDEX `idx_grant_id` (`grant_id`),
    INDEX `idx_grant_code` (`grant_code`),
    INDEX `idx_action_type` (`action_type`),
    INDEX `idx_created_at` (`created_at`),
    INDEX `idx_operator_id` (`operator_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='临时授权历史表';

-- 4. 记录迁移（可选，如果项目有migration_history表）
-- INSERT INTO `opencoze`.`migration_history`
-- (`migration_name`, `status`, `checksum`, `created_at`)
-- VALUES
-- ('20251230027000_add_temporary_grant', 'success', MD5('20251230027000'), UNIX_TIMESTAMP(NOW(3)) * 1000);
