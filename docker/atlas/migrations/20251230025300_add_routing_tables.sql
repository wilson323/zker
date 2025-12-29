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

-- Create "routing_rules" table
CREATE TABLE `opencoze`.`routing_rules` (
    `rule_id` VARCHAR(36) PRIMARY KEY COMMENT '规则ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `rule_name` VARCHAR(200) NOT NULL COMMENT '规则名称',
    `rule_type` ENUM('keyword', 'regex', 'intent', 'category') NOT NULL COMMENT '规则类型',
    `priority` INT NOT NULL DEFAULT 0 COMMENT '优先级（数字越大优先级越高）',
    `condition` JSON NOT NULL COMMENT '匹配条件',
    `target_bot_id` VARCHAR(36) NULL COMMENT '目标Bot ID',
    `target_workflow_id` VARCHAR(36) NULL COMMENT '目标工作流 ID',
    `is_active` BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒）',
    `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒）',
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,
    INDEX `idx_tenant_priority` (`tenant_id`, `priority` DESC),
    INDEX `idx_is_active` (`is_active`),
    INDEX `idx_rule_type` (`rule_type`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='路由规则表';

-- Create "routing_logs" table for logging routing decisions
CREATE TABLE `opencoze`.`routing_logs` (
    `log_id` VARCHAR(36) PRIMARY KEY COMMENT '日志ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `user_input` TEXT NOT NULL COMMENT '用户输入',
    `matched_rule_id` VARCHAR(36) NULL COMMENT '匹配的规则ID',
    `matched_bot_id` VARCHAR(36) NULL COMMENT '匹配的Bot ID',
    `matched_workflow_id` VARCHAR(36) NULL COMMENT '匹配的工作流ID',
    `confidence` DOUBLE NULL COMMENT '置信度',
    `match_type` VARCHAR(50) NULL COMMENT '匹配类型',
    `routing_score` DOUBLE NULL COMMENT '路由评分',
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒）',
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_created_at` (`created_at`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='路由日志表';
