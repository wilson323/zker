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
-- 租户系统初始化数据
-- 版本: v1.0.0
-- 日期: 2025-01-01
-- 目的: 初始化默认订阅配额配置
-- ============================================================================

-- 1. 为system默认租户创建订阅
INSERT INTO `opencoze`.`subscriptions` (
    `subscription_id`, `tenant_id`, `plan_tier`, `billing_cycle`,
    `status`, `start_date`, `end_date`,
    `quota_bots`, `quota_messages_per_month`, `quota_storage_gb`, `quota_team_members`,
    `auto_renew`, `created_at`, `updated_at`
) VALUES (
    'sub-system-default',
    'system-default',
    'free',
    'monthly',
    'active',
    CURDATE(),
    DATE_ADD(CURDATE(), INTERVAL 100 YEAR),
    10,
    1000,
    10,
    1,
    FALSE,
    UNIX_TIMESTAMP() * 1000,
    UNIX_TIMESTAMP() * 1000
)
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP() * 1000;

-- 2. 为system默认租户初始化配额
INSERT INTO `opencoze`.`quotas` (
    `quota_id`, `tenant_id`, `resource_type`, `max_limit`, `used_count`,
    `reset_cycle`, `last_reset_at`, `created_at`, `updated_at`
) VALUES
('quota-system-bots', 'system-default', 'bots', 10, 0, 'monthly', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
('quota-system-users', 'system-default', 'users', 1, 0, 'monthly', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
('quota-system-messages', 'system-default', 'message', 1000, 0, 'monthly', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
('quota-system-token', 'system-default', 'token', 100000, 0, 'monthly', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
('quota-system-storage', 'system-default', 'storage', 10, 0, 'monthly', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
('quota-system-knowledge_base', 'system-default', 'knowledge_base', 5, 0, 'monthly', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
('quota-system-workflow', 'system-default', 'workflow', 10, 0, 'monthly', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000)
ON DUPLICATE KEY UPDATE `updated_at` = UNIX_TIMESTAMP() * 1000;

-- 3. 记录初始化日志
INSERT INTO `opencoze`.`quota_usage_log` (
    `tenant_id`, `resource_type`, `action`, `amount`, `before_count`, `after_count`,
    `reason`, `created_at`
) VALUES
('system-default', 'bots', 'reset', 0, 0, 10, 'System initialization', UNIX_TIMESTAMP() * 1000),
('system-default', 'users', 'reset', 0, 0, 1, 'System initialization', UNIX_TIMESTAMP() * 1000),
('system-default', 'message', 'reset', 0, 0, 1000, 'System initialization', UNIX_TIMESTAMP() * 1000),
('system-default', 'token', 'reset', 0, 0, 100000, 'System initialization', UNIX_TIMESTAMP() * 1000),
('system-default', 'storage', 'reset', 0, 0, 10, 'System initialization', UNIX_TIMESTAMP() * 1000),
('system-default', 'knowledge_base', 'reset', 0, 0, 5, 'System initialization', UNIX_TIMESTAMP() * 1000),
('system-default', 'workflow', 'reset', 0, 0, 10, 'System initialization', UNIX_TIMESTAMP() * 1000);
