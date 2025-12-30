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
-- 迁移脚本: 用户租户隔离 - 完整迁移
-- 版本: v1.0
-- 日期: 2025-01-01
-- 作者: 研发B（后端工程师）
-- 说明: 为现有系统添加租户隔离能力
-- =====================================================

-- =====================================================
-- Step 1: user表添加tenant_id列
-- =====================================================

-- 添加tenant_id列（允许NULL，兼容现有数据）
ALTER TABLE `opencoze`.`user`
ADD COLUMN `tenant_id` VARCHAR(36) NULL COMMENT '租户ID（UUID）'
AFTER `session_key`;

-- 创建索引（优化查询性能）
CREATE INDEX `idx_tenant_id` ON `opencoze`.`user` (`tenant_id`);

-- 创建复合索引（优化租户+用户查询）
CREATE INDEX `idx_tenant_id_email` ON `opencoze`.`user` (`tenant_id`, `email`);

-- 验证列是否添加成功
SELECT
    COLUMN_NAME,
    DATA_TYPE,
    IS_NULLABLE,
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = 'opencoze'
  AND TABLE_NAME = 'user'
  AND COLUMN_NAME = 'tenant_id';

-- =====================================================
-- Step 2: 创建user_tenants关联表
-- =====================================================

CREATE TABLE `opencoze`.`user_tenants` (
    -- 主键
    `user_tenant_id` VARCHAR(36) PRIMARY KEY COMMENT '用户租户关联ID（UUID）',

    -- 外键
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',

    -- 用户在租户中的角色
    `role` ENUM('owner', 'admin', 'member') DEFAULT 'member' COMMENT '角色',
    `is_default` BOOLEAN DEFAULT FALSE COMMENT '是否为默认租户',

    -- 时间戳
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒）',
    `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒）',
    `deleted_at` BIGINT UNSIGNED NULL COMMENT '删除时间（毫秒）',

    -- 外键约束
    FOREIGN KEY (`user_id`) REFERENCES `opencoze`.`user`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,

    -- 唯一约束：一个用户在一个租户中只能有一条记录
    UNIQUE INDEX `uk_user_tenant` (`user_id`, `tenant_id`),

    -- 索引
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_is_default` (`is_default`, `deleted_at`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户租户关联表';

-- 验证表是否创建成功
SELECT
    TABLE_NAME,
    TABLE_COMMENT
FROM INFORMATION_SCHEMA.TABLES
WHERE TABLE_SCHEMA = 'opencoze'
  AND TABLE_NAME = 'user_tenants';

-- =====================================================
-- Step 3: 为现有user创建个人租户并关联
-- 策略: 为每个用户创建一个individual类型的租户
-- =====================================================

-- 3.1 为每个用户创建个人租户
INSERT INTO `opencoze`.`tenants` (
    `tenant_id`,
    `tenant_name`,
    `tenant_type`,
    `status`,
    `subscription_tier`,
    `created_at`,
    `updated_at`
)
SELECT
    CONCAT('individual-user-', CAST(`id` AS CHAR)) AS `tenant_id`,
    CONCAT(`name`, '\'s Personal Space') AS `tenant_name`,
    'individual' AS `tenant_type`,
    'active' AS `status`,
    'free' AS `subscription_tier`,
    UNIX_TIMESTAMP(NOW()) * 1000 AS `created_at`,
    UNIX_TIMESTAMP(NOW()) * 1000 AS `updated_at`
FROM `opencoze`.`user`
WHERE `deleted_at` IS NULL;

-- 3.2 更新users表的tenant_id列
UPDATE `opencoze`.`user` u
SET u.`tenant_id` = CONCAT('individual-user-', CAST(u.`id` AS CHAR))
WHERE u.`deleted_at` IS NULL
  AND u.`tenant_id` IS NULL;

-- 3.3 为每个租户创建配额记录（默认配额）
INSERT INTO `opencoze`.`quotas` (
    `quota_id`,
    `tenant_id`,
    `resource_type`,
    `max_limit`,
    `used_count`,
    `reset_cycle`,
    `last_reset_at`,
    `created_at`,
    `updated_at`
)
SELECT
    UUID() AS `quota_id`,
    t.`tenant_id`,
    qt.resource_type,
    qt.max_limit,
    0 AS used_count,
    'monthly' AS reset_cycle,
    UNIX_TIMESTAMP(NOW()) * 1000 AS last_reset_at,
    UNIX_TIMESTAMP(NOW()) * 1000 AS created_at,
    UNIX_TIMESTAMP(NOW()) * 1000 AS updated_at
FROM `opencoze`.`tenants` t
CROSS JOIN (
    SELECT 'bots' AS resource_type, 100 AS max_limit
    UNION ALL SELECT 'messages', 10000
    UNION ALL SELECT 'storage', 1024
) qt
WHERE t.`tenant_type` = 'individual'
  AND t.`deleted_at` IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM `opencoze`.`quotas` q
    WHERE q.`tenant_id` = t.`tenant_id`
    AND q.`resource_type` = qt.resource_type
);

-- 3.4 在user_tenants表中创建关联记录
INSERT INTO `opencoze`.`user_tenants` (
    `user_tenant_id`,
    `user_id`,
    `tenant_id`,
    `role`,
    `is_default`,
    `created_at`,
    `updated_at`
)
SELECT
    UUID() AS `user_tenant_id`,
    u.`id` AS `user_id`,
    u.`tenant_id` AS `tenant_id`,
    'owner' AS `role`,
    TRUE AS `is_default`,
    UNIX_TIMESTAMP(NOW()) * 1000 AS `created_at`,
    UNIX_TIMESTAMP(NOW()) * 1000 AS `updated_at`
FROM `opencoze`.`user` u
WHERE u.`deleted_at` IS NULL
  AND u.`tenant_id` IS NOT NULL
  AND NOT EXISTS (
    SELECT 1 FROM `opencoze`.`user_tenants` ut
    WHERE ut.`user_id` = u.`id`
    AND ut.`tenant_id` = u.`tenant_id`
  );

-- =====================================================
-- Step 4: 数据验证
-- =====================================================

-- 验证1: 所有user都应该有tenant_id
SELECT
    COUNT(*) AS users_without_tenant,
    'ERROR: User without tenant_id' AS message
FROM `opencoze`.`user`
WHERE `tenant_id` IS NULL
  AND `deleted_at` IS NULL;
-- 期望结果: 0

-- 验证2: 所有user在user_tenants表中都应该有对应记录
SELECT
    COUNT(*) AS orphan_users,
    'ERROR: User without user_tenant relation' AS message
FROM `opencoze`.`user` u
LEFT JOIN `opencoze`.`user_tenants` ut ON u.`id` = ut.`user_id` AND ut.`deleted_at` IS NULL
WHERE u.`deleted_at` IS NULL
  AND ut.`user_tenant_id` IS NULL;
-- 期望结果: 0

-- 验证3: user_tenants表中的所有tenant_id都应该在tenants表中存在
SELECT
    COUNT(*) AS invalid_tenants,
    'ERROR: User_tenants with non-existent tenants' AS message
FROM `opencoze`.`user_tenants` ut
LEFT JOIN `opencoze`.`tenants` t ON ut.`tenant_id` = t.`tenant_id`
WHERE ut.`deleted_at` IS NULL
  AND t.`tenant_id` IS NULL;
-- 期望结果: 0

-- 验证4: 每个用户至少有一个默认租户
SELECT
    u.`id` AS user_id,
    u.`email`,
    'WARNING: User without default tenant' AS message
FROM `opencoze`.`user` u
LEFT JOIN `opencoze`.`user_tenants` ut ON u.`id` = ut.`user_id` AND ut.`is_default` = TRUE AND ut.`deleted_at` IS NULL
WHERE u.`deleted_at` IS NULL
  AND ut.`user_tenant_id` IS NULL
LIMIT 10;
-- 期望结果: 0 rows

-- 验证5: 统计数据
SELECT
    'Total User' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`user`
WHERE `deleted_at` IS NULL

UNION ALL

SELECT
    'User with Tenant ID' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`user`
WHERE `tenant_id` IS NOT NULL
  AND `deleted_at` IS NULL

UNION ALL

SELECT
    'Total Tenants' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`tenants`
WHERE `deleted_at` IS NULL

UNION ALL

SELECT
    'Individual Tenants' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`tenants`
WHERE `tenant_type` = 'individual'
  AND `deleted_at` IS NULL

UNION ALL

SELECT
    'User-Tenant Relations' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`user_tenants`
WHERE `deleted_at` IS NULL;

-- =====================================================
-- 迁移完成
-- =====================================================

-- 记录迁移完成时间
SELECT '====================================' AS '';
SELECT 'Migration Completed Successfully!' AS '';
SELECT CONCAT('Timestamp: ', NOW()) AS '';
SELECT '====================================' AS '';
