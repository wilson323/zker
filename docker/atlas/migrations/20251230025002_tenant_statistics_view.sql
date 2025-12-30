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
-- 租户统计视图
-- 版本: v1.0.0
-- 日期: 2025-01-01
-- 目的: 提供租户使用情况统计
-- ============================================================================

-- 删除已存在的视图（如果存在）
DROP VIEW IF EXISTS `opencoze`.`v_tenant_statistics`;

-- 创建租户统计视图
CREATE VIEW `opencoze`.`v_tenant_statistics` AS
SELECT
    t.tenant_id,
    t.tenant_name,
    t.tenant_type,
    t.status,
    t.subscription_tier,

    -- 订阅信息
    s.plan_tier,
    s.start_date,
    s.end_date,
    s.auto_renew,

    -- Bot数量统计
    (SELECT COUNT(*) FROM `opencoze`.`bots` b WHERE b.tenant_id = t.tenant_id AND b.deleted_at IS NULL) as bot_count,
    (SELECT COUNT(*) FROM `opencoze`.`bots` b WHERE b.tenant_id = t.tenant_id AND b.status = 'published' AND b.deleted_at IS NULL) as published_bot_count,

    -- 用户统计
    (SELECT COUNT(DISTINCT user_id) FROM `opencoze`.`space_users` su WHERE su.space_id = t.tenant_id) as user_count,

    -- 配额使用情况
    (SELECT used_count FROM `opencoze`.`quotas` q WHERE q.tenant_id = t.tenant_id AND q.resource_type = 'bots') as quota_bots_used,
    (SELECT max_limit FROM `opencoze`.`quotas` q WHERE q.tenant_id = t.tenant_id AND q.resource_type = 'bots') as quota_bots_limit,

    (SELECT used_count FROM `opencoze`.`quotas` q WHERE q.tenant_id = t.tenant_id AND q.resource_type = 'messages') as quota_messages_used,
    (SELECT max_limit FROM `opencoze`.`quotas` q WHERE q.tenant_id = t.tenant_id AND q.resource_type = 'messages') as quota_messages_limit,

    -- 时间统计
    t.created_at as created_at,
    t.updated_at as last_updated_at

FROM `opencoze`.`tenants` t
LEFT JOIN `opencoze`.`subscriptions` s ON s.tenant_id = t.tenant_id AND s.status = 'active'
WHERE t.deleted_at IS NULL;

-- ============================================================================
-- 配额使用率视图
-- ============================================================================

-- 删除已存在的视图（如果存在）
DROP VIEW IF EXISTS `opencoze`.`v_quota_usage_rate`;

-- 创建配额使用率视图
CREATE VIEW `opencoze`.`v_quota_usage_rate` AS
SELECT
    q.tenant_id,
    t.tenant_name,
    q.resource_type,
    q.max_limit,
    q.used_count,
    CASE
        WHEN q.max_limit = -1 THEN 0  -- 无限制
        WHEN q.max_limit = 0 THEN 100  -- 避免除零错误
        ELSE ROUND((q.used_count * 100.0 / q.max_limit), 2)
    END as usage_rate,
    q.reset_cycle,
    q.last_reset_at,
    q.updated_at
FROM `opencoze`.`quotas` q
JOIN `opencoze`.`tenants` t ON t.tenant_id = q.tenant_id
WHERE t.deleted_at IS NULL;

-- ============================================================================
-- 订阅到期提醒视图
-- ============================================================================

-- 删除已存在的视图（如果存在）
DROP VIEW IF EXISTS `opencoze`.`v_subscription_expiry_alert`;

-- 创建订阅到期提醒视图
CREATE VIEW `opencoze`.`v_subscription_expiry_alert` AS
SELECT
    s.subscription_id,
    s.tenant_id,
    t.tenant_name,
    s.plan_tier,
    s.end_date,
    DATEDIFF(s.end_date, CURDATE()) as days_until_expiry,
    s.auto_renew,
    s.status,
    CASE
        WHEN DATEDIFF(s.end_date, CURDATE()) <= 7 THEN 'critical'
        WHEN DATEDIFF(s.end_date, CURDATE()) <= 30 THEN 'warning'
        ELSE 'normal'
    END as alert_level
FROM `opencoze`.`subscriptions` s
JOIN `opencoze`.`tenants` t ON t.tenant_id = s.tenant_id
WHERE s.status = 'active'
  AND s.end_date IS NOT NULL
  AND t.deleted_at IS NULL;
