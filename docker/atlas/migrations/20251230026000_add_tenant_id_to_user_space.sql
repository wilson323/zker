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
-- P0紧急修复: User和Space表添加tenant_id字段
-- 版本: v1.0.0
-- 日期: 2025-01-01
-- 目的: 完全支持多租户隔离
-- =====================================================

-- =====================================================
-- Step 1: 为user表添加tenant_id字段
-- =====================================================

-- 1.1 添加tenant_id列（允许NULL，兼容现有数据）
ALTER TABLE `opencoze`.`user`
ADD COLUMN `tenant_id` VARCHAR(36) NULL COMMENT '租户ID（UUID）'
AFTER `session_key`;

-- 1.2 创建索引（优化查询性能）
ALTER TABLE `opencoze`.`user`
ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 1.3 创建复合索引（优化租户+邮箱查询）
ALTER TABLE `opencoze`.`user`
ADD INDEX `idx_tenant_id_email` (`tenant_id`, `email`);

-- 1.4 验证列是否添加成功
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
-- Step 2: 为space表添加tenant_id字段
-- =====================================================

-- 2.1 添加tenant_id列
ALTER TABLE `opencoze`.`space`
ADD COLUMN `tenant_id` VARCHAR(36) NULL COMMENT '租户ID（UUID）'
AFTER `owner_id`;

-- 2.2 创建索引
ALTER TABLE `opencoze`.`space`
ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 2.3 创建复合索引（优化租户+所有者查询）
ALTER TABLE `opencoze`.`space`
ADD INDEX `idx_tenant_id_owner` (`tenant_id`, `owner_id`);

-- 2.4 验证列是否添加成功
SELECT
    COLUMN_NAME,
    DATA_TYPE,
    IS_NULLABLE,
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = 'opencoze'
  AND TABLE_NAME = 'space'
  AND COLUMN_NAME = 'tenant_id';

-- =====================================================
-- Step 3: 为space_user表添加tenant_id字段
-- =====================================================

-- 3.1 添加tenant_id列
ALTER TABLE `opencoze`.`space_user`
ADD COLUMN `tenant_id` VARCHAR(36) NULL COMMENT '租户ID（UUID）'
AFTER `user_id`;

-- 3.2 创建索引
ALTER TABLE `opencoze`.`space_user`
ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 3.3 创建复合索引（优化租户+用户查询）
ALTER TABLE `opencoze`.`space_user`
ADD INDEX `idx_tenant_id_user` (`tenant_id`, `user_id`);

-- 3.4 验证列是否添加成功
SELECT
    COLUMN_NAME,
    DATA_TYPE,
    IS_NULLABLE,
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = 'opencoze'
  AND TABLE_NAME = 'space_user'
  AND COLUMN_NAME = 'tenant_id';

-- =====================================================
-- Step 4: 数据迁移 - 为现有记录创建个人租户
-- 策略: 为每个用户/空间创建个人租户并关联
-- =====================================================

-- 4.1 为每个用户创建个人租户（如果不存在）
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
WHERE `deleted_at` IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM `opencoze`.`tenants` t
    WHERE t.`tenant_id` = CONCAT('individual-user-', CAST(`user`.`id` AS CHAR))
);

-- 4.2 更新user表的tenant_id列
UPDATE `opencoze`.`user` u
SET u.`tenant_id` = CONCAT('individual-user-', CAST(u.`id` AS CHAR))
WHERE u.`deleted_at` IS NULL
  AND (u.`tenant_id` IS NULL OR u.`tenant_id` = '');

-- 4.3 为每个空间创建个人租户（如果不存在）
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
    CONCAT('individual-space-', CAST(`id` AS CHAR)) AS `tenant_id`,
    CONCAT(`name`, ' Space') AS `tenant_name`,
    'individual' AS `tenant_type`,
    'active' AS `status`,
    'free' AS `subscription_tier`,
    UNIX_TIMESTAMP(NOW()) * 1000 AS `created_at`,
    UNIX_TIMESTAMP(NOW()) * 1000 AS `updated_at`
FROM `opencoze`.`space`
WHERE `deleted_at` IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM `opencoze`.`tenants` t
    WHERE t.`tenant_id` = CONCAT('individual-space-', CAST(`space`.`id` AS CHAR))
);

-- 4.4 更新space表的tenant_id列
UPDATE `opencoze`.`space` s
SET s.`tenant_id` = CONCAT('individual-space-', CAST(s.`id` AS CHAR))
WHERE s.`deleted_at` IS NULL
  AND (s.`tenant_id` IS NULL OR s.`tenant_id` = '');

-- 4.5 更新space_user表的tenant_id列（继承自space）
UPDATE `opencoze`.`space_user` su
INNER JOIN `opencoze`.`space` s ON su.`space_id` = s.`id`
SET su.`tenant_id` = s.`tenant_id`
WHERE su.`tenant_id` IS NULL OR su.`tenant_id` = '';

-- =====================================================
-- Step 5: 数据验证
-- =====================================================

-- 验证1: user表所有记录都应该有tenant_id
SELECT
    COUNT(*) AS users_without_tenant,
    'ERROR: Users without tenant_id' AS message
FROM `opencoze`.`user`
WHERE (`tenant_id` IS NULL OR `tenant_id` = '')
  AND `deleted_at` IS NULL;
-- 期望结果: 0

-- 验证2: space表所有记录都应该有tenant_id
SELECT
    COUNT(*) AS spaces_without_tenant,
    'ERROR: Spaces without tenant_id' AS message
FROM `opencoze`.`space`
WHERE (`tenant_id` IS NULL OR `tenant_id` = '')
  AND `deleted_at` IS NULL;
-- 期望结果: 0

-- 验证3: space_user表所有记录都应该有tenant_id
SELECT
    COUNT(*) AS space_users_without_tenant,
    'ERROR: Space_users without tenant_id' AS message
FROM `opencoze`.`space_user`
WHERE `tenant_id` IS NULL OR `tenant_id` = '';
-- 期望结果: 0

-- 验证4: space_user的tenant_id应该与space的tenant_id一致
SELECT
    COUNT(*) AS mismatch_tenant_count,
    'ERROR: space_user tenant_id does not match space tenant_id' AS message
FROM `opencoze`.`space_user` su
INNER JOIN `opencoze`.`space` s ON su.`space_id` = s.`id`
WHERE su.`tenant_id` != s.`tenant_id`;
-- 期望结果: 0

-- 验证5: 统计数据
SELECT
    'Total Users' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`user`
WHERE `deleted_at` IS NULL

UNION ALL

SELECT
    'Users with Tenant ID' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`user`
WHERE `tenant_id` IS NOT NULL
  AND `tenant_id` != ''
  AND `deleted_at` IS NULL

UNION ALL

SELECT
    'Total Spaces' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`space`
WHERE `deleted_at` IS NULL

UNION ALL

SELECT
    'Spaces with Tenant ID' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`space`
WHERE `tenant_id` IS NOT NULL
  AND `tenant_id` != ''
  AND `deleted_at` IS NULL

UNION ALL

SELECT
    'Total Space Users' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`space_user`

UNION ALL

SELECT
    'Space Users with Tenant ID' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`space_user`
WHERE `tenant_id` IS NOT NULL
  AND `tenant_id` != '';

-- =====================================================
-- Step 6: 添加外键约束（可选，如果性能允许）
-- =====================================================

-- 注意：外键约束会影响性能，建议在数据验证通过后再添加
-- 如果不需要外键约束，可以跳过此步骤

-- 6.1 为user表添加外键约束
-- ALTER TABLE `opencoze`.`user`
-- ADD CONSTRAINT `fk_user_tenant`
-- FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`)
-- ON DELETE CASCADE;

-- 6.2 为space表添加外键约束
-- ALTER TABLE `opencoze`.`space`
-- ADD CONSTRAINT `fk_space_tenant`
-- FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`)
-- ON DELETE CASCADE;

-- 6.3 为space_user表添加外键约束
-- ALTER TABLE `opencoze`.`space_user`
-- ADD CONSTRAINT `fk_space_user_tenant`
-- FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`)
-- ON DELETE CASCADE;

-- =====================================================
-- 迁移完成
-- =====================================================

-- 记录迁移完成时间
SELECT '====================================' AS '';
SELECT 'Migration Completed Successfully!' AS '';
SELECT CONCAT('Timestamp: ', NOW()) AS '';
SELECT '====================================' AS '';
