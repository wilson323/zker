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
-- 回滚脚本: 删除租户系统表
-- 版本: v1.0.0
-- 日期: 2025-01-01
-- 目的: 完整回滚租户系统相关对象
-- ============================================================================

-- 1. 删除视图
DROP VIEW IF EXISTS `opencoze`.`v_subscription_expiry_alert`;
DROP VIEW IF EXISTS `opencoze`.`v_quota_usage_rate`;
DROP VIEW IF EXISTS `opencoze`.`v_tenant_statistics`;

-- 2. 删除触发器
DROP TRIGGER IF EXISTS `opencoze`.`trg_quota_update_timestamp`;
DROP TRIGGER IF EXISTS `opencoze`.`trg_tenant_status_update`;

-- 3. 删除表（按依赖关系逆序）
-- 注意：需要先禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `opencoze`.`quota_usage_log`;
DROP TABLE IF EXISTS `opencoze`.`quotas`;
DROP TABLE IF EXISTS `opencoze`.`subscriptions`;
DROP TABLE IF EXISTS `opencoze`.`tenants`;

SET FOREIGN_KEY_CHECKS = 1;

-- 4. 清理初始化数据
-- 注意：由于表已删除，数据已自动清理

-- 5. 记录回滚信息（如果存在migration_history表）
-- UPDATE `migration_history`
-- SET `status` = 'rolled_back',
--     `rolled_back_at` = NOW()
-- WHERE `migration_name` = '20251230025000_create_tenant_tables';

SELECT 'Tenant tables rollback completed successfully' as status;
