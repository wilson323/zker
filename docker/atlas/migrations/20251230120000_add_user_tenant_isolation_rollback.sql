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
-- 回滚脚本: 用户租户隔离 - 完整回滚
-- 版本: v1.0
-- 日期: 2025-01-01
-- 作者: 研发B（后端工程师）
-- 说明: 回滚用户租户隔离相关的所有变更
-- ⚠️ 警告: 此操作将删除数据，请谨慎执行！
-- =====================================================

-- =====================================================
-- Step 1: 删除user_tenants关联表
-- =====================================================

-- 删除外键约束
ALTER TABLE `opencoze`.`user_tenants` DROP FOREIGN KEY `user_tenants_ibfk_1`;
ALTER TABLE `opencoze`.`user_tenants` DROP FOREIGN KEY `user_tenants_ibfk_2`;

-- 删除表
DROP TABLE IF EXISTS `opencoze`.`user_tenants`;

-- 验证表是否删除成功
SELECT COUNT(*) AS table_exists
FROM INFORMATION_SCHEMA.TABLES
WHERE TABLE_SCHEMA = 'opencoze'
  AND TABLE_NAME = 'user_tenants';
-- 期望结果: 0

-- =====================================================
-- Step 2: 删除users表的tenant_id列
-- =====================================================

-- 删除索引
DROP INDEX IF EXISTS `idx_tenant_id` ON `opencoze`.`users`;
DROP INDEX IF EXISTS `idx_tenant_id_email` ON `opencoze`.`users`;

-- 删除列
ALTER TABLE `opencoze`.`users`
DROP COLUMN IF EXISTS `tenant_id`;

-- 验证列是否删除成功
SELECT COUNT(*) AS column_exists
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = 'opencoze'
  AND TABLE_NAME = 'users'
  AND COLUMN_NAME = 'tenant_id';
-- 期望结果: 0

-- =====================================================
-- Step 3: 删除为用户创建的个人租户
-- =====================================================

-- 删除个人租户的配额记录
DELETE FROM `opencoze`.`quotas`
WHERE `tenant_id` IN (
    SELECT `tenant_id`
    FROM `opencoze`.`tenants`
    WHERE `tenant_id` LIKE 'individual-user-%'
);

-- 删除个人租户的订阅记录
DELETE FROM `opencoze`.`subscriptions`
WHERE `tenant_id` IN (
    SELECT `tenant_id`
    FROM `opencoze`.`tenants`
    WHERE `tenant_id` LIKE 'individual-user-%'
);

-- 删除个人租户
DELETE FROM `opencoze`.`tenants`
WHERE `tenant_id` LIKE 'individual-user-%';

-- 验证个人租户是否删除成功
SELECT COUNT(*) AS individual_tenants_remaining
FROM `opencoze`.`tenants`
WHERE `tenant_id` LIKE 'individual-user-%';
-- 期望结果: 0

-- =====================================================
-- Step 4: 验证回滚结果
-- =====================================================

-- 验证1: users表不应该有tenant_id列
SELECT
    COLUMN_NAME,
    'ERROR: tenant_id column still exists in users table' AS message
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = 'opencoze'
  AND TABLE_NAME = 'users'
  AND COLUMN_NAME = 'tenant_id';
-- 期望结果: 0 rows

-- 验证2: user_tenants表不应该存在
SELECT
    TABLE_NAME,
    'ERROR: user_tenants table still exists' AS message
FROM INFORMATION_SCHEMA.TABLES
WHERE TABLE_SCHEMA = 'opencoze'
  AND TABLE_NAME = 'user_tenants';
-- 期望结果: 0 rows

-- 验证3: 个人租户应该被删除
SELECT COUNT(*) AS individual_tenants_count
FROM `opencoze`.`tenants`
WHERE `tenant_id` LIKE 'individual-user-%';
-- 期望结果: 0

-- =====================================================
-- 回滚完成
-- =====================================================

-- 记录回滚完成时间
SELECT '====================================' AS '';
SELECT 'Rollback Completed Successfully!' AS '';
SELECT CONCAT('Timestamp: ', NOW()) AS '';
SELECT '====================================' AS '';
