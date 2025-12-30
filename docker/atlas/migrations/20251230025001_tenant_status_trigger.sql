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
-- 租户状态更新触发器
-- 版本: v1.0.0
-- 日期: 2025-01-01
-- 目的: 当租户状态变更时自动处理相关订阅
-- ============================================================================

DELIMITER $$

-- 删除已存在的触发器（如果存在）
DROP TRIGGER IF EXISTS `trg_tenant_status_update`$$

-- 创建触发器
CREATE TRIGGER `trg_tenant_status_update`
BEFORE UPDATE ON `opencoze`.`tenants`
FOR EACH ROW
BEGIN
    -- 当租户状态从active变为suspended时
    IF NEW.status != OLD.status AND OLD.status = 'active' AND NEW.status = 'suspended' THEN
        -- 暂停所有订阅
        UPDATE `opencoze`.`subscriptions`
        SET `status` = 'past_due',
            `updated_at` = UNIX_TIMESTAMP() * 1000
        WHERE `tenant_id` = NEW.tenant_id
          AND `status` = 'active';

        -- 记录日志
        INSERT INTO `opencoze`.`quota_usage_log`
        (`tenant_id`, `resource_type`, `action`, `amount`, `before_count`, `after_count`, `reason`, `created_at`)
        VALUES
        (NEW.tenant_id, 'bots', 'consume', 0, 0, 0,
         CONCAT('Tenant suspended from ', OLD.status, ' to ', NEW.status),
         UNIX_TIMESTAMP() * 1000);
    END IF;

    -- 当租户状态从suspended变为active时
    IF NEW.status != OLD.status AND OLD.status = 'suspended' AND NEW.status = 'active' THEN
        -- 恢复订阅
        UPDATE `opencoze`.`subscriptions`
        SET `status` = 'active',
            `updated_at` = UNIX_TIMESTAMP() * 1000
        WHERE `tenant_id` = NEW.tenant_id
          AND `status` = 'past_due'
          AND `end_date` >= CURDATE();
    END IF;
END$$

DELIMITER ;

-- ============================================================================
-- 配额更新触发器
-- 目的: 自动更新配额使用时间戳
-- ============================================================================

DELIMITER $$

-- 删除已存在的触发器（如果存在）
DROP TRIGGER IF EXISTS `trg_quota_update_timestamp`$$

-- 创建触发器
CREATE TRIGGER `trg_quota_update_timestamp`
BEFORE UPDATE ON `opencoze`.`quotas`
FOR EACH ROW
BEGIN
    -- 当used_count发生变化时，自动更新updated_at
    IF NEW.used_count != OLD.used_count THEN
        SET NEW.updated_at = UNIX_TIMESTAMP() * 1000;
    END IF;
END$$

DELIMITER ;
