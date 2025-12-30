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

-- Create "digital_employee_profiles" table
CREATE TABLE `opencoze`.`digital_employee_profiles` (
    `employee_id` VARCHAR(36) PRIMARY KEY COMMENT '员工ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `name` VARCHAR(100) NOT NULL COMMENT '员工名称',
    `avatar` VARCHAR(255) NULL COMMENT '头像URL',
    `role` ENUM('customer_service', 'sales', 'tech_support', 'consultant', 'trainer') NOT NULL COMMENT '员工角色',
    `skills` JSON NOT NULL COMMENT '技能标签数组',
    `personality` VARCHAR(100) NULL COMMENT '性格特征',
    `knowledge_base_id` VARCHAR(36) NULL COMMENT '关联知识库ID',
    `bot_id` VARCHAR(36) NULL COMMENT '关联Bot ID',
    `status` ENUM('active', 'inactive', 'deleted') NOT NULL DEFAULT 'active' COMMENT '员工状态',
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒）',
    `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒）',
    `deleted_at` BIGINT UNSIGNED NULL COMMENT '删除时间（毫秒）',
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_role` (`role`),
    INDEX `idx_status` (`status`),
    INDEX `idx_bot_id` (`bot_id`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数字员工画像表';

-- Create "digital_employee_task_assignments" table
CREATE TABLE `opencoze`.`digital_employee_task_assignments` (
    `assignment_id` VARCHAR(36) PRIMARY KEY COMMENT '分配ID',
    `task_id` VARCHAR(36) NOT NULL COMMENT '任务ID',
    `employee_id` VARCHAR(36) NOT NULL COMMENT '员工ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `task_type` ENUM('customer_service', 'sales_follow', 'tech_support', 'consultation', 'training') NOT NULL COMMENT '任务类型',
    `priority` ENUM('high', 'medium', 'low') NOT NULL COMMENT '任务优先级',
    `status` ENUM('assigned', 'in_progress', 'completed', 'failed', 'cancelled') NOT NULL DEFAULT 'assigned' COMMENT '任务状态',
    `assigned_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '分配时间（毫秒）',
    `completed_at` BIGINT UNSIGNED NULL COMMENT '完成时间（毫秒）',
    `result` TEXT NULL COMMENT '任务结果',
    `failure_reason` VARCHAR(255) NULL COMMENT '失败原因',
    FOREIGN KEY (`employee_id`) REFERENCES `opencoze`.`digital_employee_profiles`(`employee_id`) ON DELETE CASCADE,
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,
    INDEX `idx_task_id` (`task_id`),
    INDEX `idx_employee_id` (`employee_id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_task_type` (`task_type`),
    INDEX `idx_status` (`status`),
    INDEX `idx_priority` (`priority`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数字员工任务分配表';

-- Create "digital_employee_performance" table
CREATE TABLE `opencoze`.`digital_employee_performance` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `employee_id` VARCHAR(36) NOT NULL COMMENT '员工ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `total_tasks` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '总任务数',
    `completed_tasks` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '完成任务数',
    `failed_tasks` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '失败任务数',
    `completion_rate` DECIMAL(5,2) NOT NULL DEFAULT 0.00 COMMENT '完成率（%）',
    `avg_response_time` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '平均响应时间（秒）',
    `customer_rating` DECIMAL(3,2) NOT NULL DEFAULT 0.00 COMMENT '客户评分（1-5分）',
    `period` ENUM('daily', 'weekly', 'monthly') NOT NULL COMMENT '统计周期',
    `date` DATE NOT NULL COMMENT '统计日期',
    `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒）',
    PRIMARY KEY (`id`),
    FOREIGN KEY (`employee_id`) REFERENCES `opencoze`.`digital_employee_profiles`(`employee_id`) ON DELETE CASCADE,
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,
    UNIQUE KEY `uk_employee_period_date` (`employee_id`, `period`, `date`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_date` (`date`),
    INDEX `idx_period` (`period`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数字员工绩效统计表';
