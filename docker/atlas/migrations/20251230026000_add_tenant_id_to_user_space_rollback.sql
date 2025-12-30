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
-- 回滚脚本: 移除User和Space表的tenant_id字段
-- 版本: v1.0.0
-- 日期: 2025-01-01
-- 说明: 安全回滚，仅移除tenant_id字段，不删除tenants表数据
-- =====================================================

-- =====================================================
-- Step 1: 删除外键约束（如果已创建）
-- =====================================================

-- 1.1 删除user表外键约束
-- ALTER TABLE `opencoze`.`user` DROP FOREIGN KEY `fk_user_tenant`;

-- 1.2 删除space表外键约束
-- ALTER TABLE `opencoze`.`space` DROP FOREIGN KEY `fk_space_tenant`;

-- 1.3 删除space_user表外键约束
-- ALTER TABLE `opencoze`.`space_user` DROP FOREIGN KEY `fk_space_user_tenant`;

-- =====================================================
-- Step 2: 删除复合索引
-- =====================================================

-- 2.1 删除user表复合索引
ALTER TABLE `opencoze`.`user` DROP INDEX `idx_tenant_id_email`;

-- 2.2 删除space表复合索引
ALTER TABLE `opencoze`.`space` DROP INDEX `idx_tenant_id_owner`;

-- 2.3 删除space_user表复合索引
ALTER TABLE `opencoze`.`space_user` DROP INDEX `idx_tenant_id_user`;

-- =====================================================
-- Step 3: 删除tenant_id索引
-- =====================================================

-- 3.1 删除user表索引
ALTER TABLE `opencoze`.`user` DROP INDEX `idx_tenant_id`;

-- 3.2 删除space表索引
ALTER TABLE `opencoze`.`space` DROP INDEX `idx_tenant_id`;

-- 3.3 删除space_user表索引
ALTER TABLE `opencoze`.`space_user` DROP INDEX `idx_tenant_id`;

-- =====================================================
-- Step 4: 删除tenant_id列
-- =====================================================

-- 4.1 删除user表的tenant_id列
ALTER TABLE `opencoze`.`user` DROP COLUMN `tenant_id`;

-- 4.2 删除space表的tenant_id列
ALTER TABLE `opencoze`.`space` DROP COLUMN `tenant_id`;

-- 4.3 删除space_user表的tenant_id列
ALTER TABLE `opencoze`.`space_user` DROP COLUMN `tenant_id`;

-- =====================================================
-- Step 5: 验证回滚结果
-- =====================================================

-- 验证1: 确认user表不再有tenant_id列
SELECT
    COUNT(*) AS column_exists,
    CASE
        WHEN COUNT(*) = 0 THEN 'SUCCESS: tenant_id column removed from user table'
        ELSE 'ERROR: tenant_id column still exists in user table'
    END AS message
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = 'opencoze'
  AND TABLE_NAME = 'user'
  AND COLUMN_NAME = 'tenant_id';
-- 期望结果: 0

-- 验证2: 确认space表不再有tenant_id列
SELECT
    COUNT(*) AS column_exists,
    CASE
        WHEN COUNT(*) = 0 THEN 'SUCCESS: tenant_id column removed from space table'
        ELSE 'ERROR: tenant_id column still exists in space table'
    END AS message
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = 'opencoze'
  AND TABLE_NAME = 'space'
  AND COLUMN_NAME = 'tenant_id';
-- 期望结果: 0

-- 验证3: 确认space_user表不再有tenant_id列
SELECT
    COUNT(*) AS column_exists,
    CASE
        WHEN COUNT(*) = 0 THEN 'SUCCESS: tenant_id column removed from space_user table'
        ELSE 'ERROR: tenant_id column still exists in space_user table'
    END AS message
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = 'opencoze'
  AND TABLE_NAME = 'space_user'
  AND COLUMN_NAME = 'tenant_id';
-- 期望结果: 0

-- 验证4: 确认数据完整性
SELECT
    'Total Users' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`user`
WHERE `deleted_at` IS NULL

UNION ALL

SELECT
    'Total Spaces' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`space`
WHERE `deleted_at` IS NULL

UNION ALL

SELECT
    'Total Space Users' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`space_user`;

-- =====================================================
-- 注意事项
-- =====================================================

-- 注意：
-- 1. 本回滚脚本不会删除tenants表中的数据，仅删除user/space/space_user表的tenant_id字段
-- 2. 如果需要完全清理tenants表，请手动执行：
--    DROP TABLE IF EXISTS `opencoze`.`tenants`;
--    DROP TABLE IF EXISTS `opencoze`.`user_tenants`;
--    DROP TABLE IF EXISTS `opencoze`.`subscriptions`;
--    DROP TABLE IF EXISTS `opencoze`.`quotas`;
-- 3. 建议在回滚前备份数据库

-- =====================================================
-- 回滚完成
-- =====================================================

SELECT '====================================' AS '';
SELECT 'Rollback Completed Successfully!' AS '';
SELECT CONCAT('Timestamp: ', NOW()) AS '';
SELECT '====================================' AS '';
