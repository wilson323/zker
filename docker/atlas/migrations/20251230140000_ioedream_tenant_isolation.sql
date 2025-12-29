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

-- =====================================================
-- 租户隔离迁移 - ioedream数据库适配版本
-- 版本: v1.0
-- 日期: 2025-01-01
-- 说明: 为ioedream数据库添加完整的租户隔离功能
-- =====================================================

-- =====================================================
-- Step 1: 创建租户相关表
-- =====================================================

-- 1.1 创建tenants表
CREATE TABLE IF NOT EXISTS `tenants` (
    `tenant_id` VARCHAR(36) PRIMARY KEY COMMENT '租户ID（UUID）',
    `tenant_name` VARCHAR(200) NOT NULL COMMENT '租户名称',
    `tenant_type` ENUM('individual', 'team', 'enterprise') NOT NULL COMMENT '租户类型',
    `status` ENUM('active', 'suspended', 'deleted') DEFAULT 'active' COMMENT '状态',
    `subscription_tier` ENUM('free', 'pro', 'enterprise') DEFAULT 'free' COMMENT '订阅等级',
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX `idx_status_type` (`status`, `tenant_type`),
    INDEX `idx_subscription_tier` (`subscription_tier`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户表';

-- 1.2 创建subscriptions表（订阅表）
CREATE TABLE IF NOT EXISTS `subscriptions` (
    `subscription_id` VARCHAR(36) PRIMARY KEY COMMENT '订阅ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `plan_tier` ENUM('free', 'pro', 'enterprise') NOT NULL COMMENT '套餐等级',
    `billing_cycle` ENUM('monthly', 'yearly') NOT NULL COMMENT '计费周期',
    `start_date` DATE NOT NULL COMMENT '开始日期',
    `end_date` DATE NULL COMMENT '结束日期',
    `auto_renew` BOOLEAN DEFAULT TRUE COMMENT '自动续费',
    `status` ENUM('active', 'expired', 'cancelled') DEFAULT 'active' COMMENT '状态',
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    FOREIGN KEY (`tenant_id`) REFERENCES `tenants`(`tenant_id`) ON DELETE CASCADE,
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_status_end_date` (`status`, `end_date`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订阅表';

-- 1.3 创建quotas表（配额表）
CREATE TABLE IF NOT EXISTS `quotas` (
    `quota_id` VARCHAR(36) PRIMARY KEY COMMENT '配额ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `resource_type` ENUM('bots', 'messages', 'storage', 'team_members', 'visitors') NOT NULL COMMENT '资源类型',
    `max_limit` INT NOT NULL COMMENT '最大限制',
    `used_count` INT DEFAULT 0 COMMENT '已使用数量',
    `reset_cycle` ENUM('daily', 'monthly', 'yearly', 'never') DEFAULT 'monthly' COMMENT '重置周期',
    `last_reset_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '上次重置时间',
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    FOREIGN KEY (`tenant_id`) REFERENCES `tenants`(`tenant_id`) ON DELETE CASCADE,
    UNIQUE INDEX `uk_tenant_resource` (`tenant_id`, `resource_type`),
    INDEX `idx_tenant_id` (`tenant_id`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='配额表';

-- 1.4 创建user_tenants关联表
CREATE TABLE IF NOT EXISTS `user_tenants` (
    `user_tenant_id` VARCHAR(36) PRIMARY KEY COMMENT '用户租户关联ID（UUID）',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `role` ENUM('owner', 'admin', 'member') DEFAULT 'member' COMMENT '角色',
    `is_default` BOOLEAN DEFAULT FALSE COMMENT '是否为默认租户',
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    FOREIGN KEY (`user_id`) REFERENCES `t_common_user`(`user_id`) ON DELETE CASCADE,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenants`(`tenant_id`) ON DELETE CASCADE,
    UNIQUE INDEX `uk_user_tenant` (`user_id`, `tenant_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_is_default` (`is_default`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户租户关联表';

-- =====================================================
-- Step 2: 为t_common_user表添加tenant_id列
-- =====================================================

-- 2.1 添加tenant_id列
ALTER TABLE `t_common_user`
ADD COLUMN `tenant_id` VARCHAR(36) NULL COMMENT '租户ID（UUID）' AFTER `user_id`;

-- 2.2 创建索引
ALTER TABLE `t_common_user`
ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 2.3 创建复合索引（租户+邮箱）
ALTER TABLE `t_common_user`
ADD INDEX `idx_tenant_id_email` (`tenant_id`, `email`);

-- 验证列是否添加成功
SELECT
    COLUMN_NAME,
    DATA_TYPE,
    IS_NULLABLE,
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 't_common_user'
  AND COLUMN_NAME = 'tenant_id';

-- =====================================================
-- Step 3: 为现有用户创建个人租户
-- =====================================================

-- 3.1 为每个用户创建个人租户
INSERT INTO `tenants` (
    `tenant_id`,
    `tenant_name`,
    `tenant_type`,
    `status`,
    `subscription_tier`
)
SELECT
    CONCAT('individual-user-', CAST(`user_id` AS CHAR)) AS `tenant_id`,
    CONCAT(
        COALESCE(`nickname`, `username`),
        '\'s Personal Space'
    ) AS `tenant_name`,
    'individual' AS `tenant_type`,
    'active' AS `status`,
    'free' AS `subscription_tier`
FROM `t_common_user`
WHERE `deleted_flag` = 0;

-- 3.2 为每个租户创建默认配额
INSERT INTO `quotas` (
    `quota_id`,
    `tenant_id`,
    `resource_type`,
    `max_limit`,
    `used_count`,
    `reset_cycle`
)
SELECT
    UUID() AS `quota_id`,
    t.`tenant_id`,
    qt.resource_type,
    qt.max_limit,
    0 AS used_count,
    'monthly' AS reset_cycle
FROM `tenants` t
CROSS JOIN (
    SELECT 'visitors' AS resource_type, 1000 AS max_limit
    UNION ALL SELECT 'messages', 10000
    UNION ALL SELECT 'storage', 10240
) qt
WHERE t.`tenant_type` = 'individual'
  AND NOT EXISTS (
    SELECT 1 FROM `quotas` q
    WHERE q.`tenant_id` = t.`tenant_id`
    AND q.`resource_type` = qt.resource_type
);

-- 3.3 为每个租户创建免费订阅
INSERT INTO `subscriptions` (
    `subscription_id`,
    `tenant_id`,
    `plan_tier`,
    `billing_cycle`,
    `start_date`,
    `auto_renew`,
    `status`
)
SELECT
    UUID() AS `subscription_id`,
    t.`tenant_id`,
    'free' AS plan_tier,
    'monthly' AS billing_cycle,
    CURDATE() AS start_date,
    TRUE AS auto_renew,
    'active' AS status
FROM `tenants` t
WHERE t.`tenant_type` = 'individual'
  AND NOT EXISTS (
    SELECT 1 FROM `subscriptions` s
    WHERE s.`tenant_id` = t.`tenant_id`
);

-- 3.4 更新t_common_user表的tenant_id列
UPDATE `t_common_user` u
SET u.`tenant_id` = CONCAT('individual-user-', CAST(u.`user_id` AS CHAR))
WHERE u.`deleted_flag` = 0
  AND u.`tenant_id` IS NULL;

-- 3.5 在user_tenants表中创建关联记录
INSERT INTO `user_tenants` (
    `user_tenant_id`,
    `user_id`,
    `tenant_id`,
    `role`,
    `is_default`
)
SELECT
    UUID() AS `user_tenant_id`,
    u.`user_id` AS `user_id`,
    u.`tenant_id` AS `tenant_id`,
    'owner' AS `role`,
    TRUE AS `is_default`
FROM `t_common_user` u
WHERE u.`deleted_flag` = 0
  AND u.`tenant_id` IS NOT NULL
  AND NOT EXISTS (
    SELECT 1 FROM `user_tenants` ut
    WHERE ut.`user_id` = u.`user_id`
    AND ut.`tenant_id` = u.`tenant_id`
);

-- =====================================================
-- Step 4: 数据验证
-- =====================================================

-- 验证1: tenants表是否创建成功
SELECT COUNT(*) AS tenants_count
FROM INFORMATION_SCHEMA.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'tenants';
-- 期望: 1

-- 验证2: 所有t_common_user都应该有tenant_id
SELECT
    COUNT(*) AS users_without_tenant,
    'ERROR: Users without tenant_id' AS message
FROM `t_common_user`
WHERE `tenant_id` IS NULL
  AND `deleted_flag` = 0;
-- 期望: 0

-- 验证3: 所有t_common_user在user_tenants表中都应该有对应记录
SELECT
    COUNT(*) AS orphan_users,
    'ERROR: Users without user_tenant relation' AS message
FROM `t_common_user` u
LEFT JOIN `user_tenants` ut ON u.`user_id` = ut.`user_id`
WHERE u.`deleted_flag` = 0
  AND ut.`user_tenant_id` IS NULL;
-- 期望: 0

-- 验证4: user_tenants表中的所有tenant_id都应该在tenants表中存在
SELECT
    COUNT(*) AS invalid_tenants,
    'ERROR: User_tenants with non-existent tenants' AS message
FROM `user_tenants` ut
LEFT JOIN `tenants` t ON ut.`tenant_id` = t.`tenant_id`
WHERE t.`tenant_id` IS NULL;
-- 期望: 0

-- 验证5: 统计数据
SELECT
    'Total Users' AS metric,
    COUNT(*) AS value
FROM `t_common_user`
WHERE `deleted_flag` = 0

UNION ALL

SELECT
    'Users with Tenant ID' AS metric,
    COUNT(*) AS value
FROM `t_common_user`
WHERE `tenant_id` IS NOT NULL
  AND `deleted_flag` = 0

UNION ALL

SELECT
    'Total Tenants' AS metric,
    COUNT(*) AS value
FROM `tenants`

UNION ALL

SELECT
    'Individual Tenants' AS metric,
    COUNT(*) AS value
FROM `tenants`
WHERE `tenant_type` = 'individual'

UNION ALL

SELECT
    'User-Tenant Relations' AS metric,
    COUNT(*) AS value
FROM `user_tenants`;

-- =====================================================
-- 迁移完成
-- =====================================================

SELECT '====================================' AS '';
SELECT 'Migration Completed Successfully!' AS '';
SELECT CONCAT('Timestamp: ', NOW()) AS '';
SELECT '====================================' AS '';
