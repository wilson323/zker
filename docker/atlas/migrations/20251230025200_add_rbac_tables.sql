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

-- Create "roles" table
CREATE TABLE `opencoze`.`roles` (
    `role_id` VARCHAR(36) PRIMARY KEY COMMENT '角色ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `role_name` VARCHAR(100) NOT NULL COMMENT '角色名称',
    `role_code` VARCHAR(50) NOT NULL COMMENT '角色编码',
    `role_type` ENUM('system', 'custom') NOT NULL COMMENT '角色类型',
    `parent_role_id` VARCHAR(36) NULL COMMENT '父角色ID',
    `description` TEXT NULL COMMENT '角色描述',
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒）',
    `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒）',
    `deleted_at` BIGINT UNSIGNED NULL COMMENT '删除时间（毫秒）',
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,
    FOREIGN KEY (`parent_role_id`) REFERENCES `opencoze`.`roles`(`role_id`) ON DELETE SET NULL,
    UNIQUE INDEX `uk_tenant_code` (`tenant_id`, `role_code`, `deleted_at`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_role_type` (`role_type`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- Create "data_permissions" table
CREATE TABLE `opencoze`.`data_permissions` (
    `permission_id` VARCHAR(36) PRIMARY KEY COMMENT '权限ID',
    `role_id` VARCHAR(36) NOT NULL COMMENT '角色ID',
    `resource_type` ENUM('bots', 'conversations', 'knowledge', 'workflows', 'plugins') NOT NULL COMMENT '资源类型',
    `scope` ENUM('ALL', 'DEPARTMENT', 'OWN', 'CUSTOM', 'NONE') NOT NULL COMMENT '权限范围',
    `custom_filter` JSON NULL COMMENT '自定义过滤条件',
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒）',
    `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒）',
    FOREIGN KEY (`role_id`) REFERENCES `opencoze`.`roles`(`role_id`) ON DELETE CASCADE,
    INDEX `idx_role_id` (`role_id`),
    INDEX `idx_resource_type` (`resource_type`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据权限表';

-- Create "field_permissions" table
CREATE TABLE `opencoze`.`field_permissions` (
    `permission_id` VARCHAR(36) PRIMARY KEY COMMENT '权限ID',
    `role_id` VARCHAR(36) NOT NULL COMMENT '角色ID',
    `resource_type` VARCHAR(50) NOT NULL COMMENT '资源类型',
    `field_name` VARCHAR(100) NOT NULL COMMENT '字段名称',
    `permission_level` ENUM('hidden', 'readonly', 'editable') NOT NULL COMMENT '权限级别',
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒）',
    `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒）',
    FOREIGN KEY (`role_id`) REFERENCES `opencoze`.`roles`(`role_id`) ON DELETE CASCADE,
    UNIQUE INDEX `uk_role_resource_field` (`role_id`, `resource_type`, `field_name`),
    INDEX `idx_role_id` (`role_id`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='字段权限表';

-- Create "user_roles" table
CREATE TABLE `opencoze`.`user_roles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `role_id` VARCHAR(36) NOT NULL COMMENT '角色ID',
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒）',
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,
    FOREIGN KEY (`role_id`) REFERENCES `opencoze`.`roles`(`role_id`) ON DELETE CASCADE,
    UNIQUE INDEX `uk_user_tenant_role` (`user_id`, `tenant_id`, `role_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_role_id` (`role_id`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- Create "departments" table
CREATE TABLE `opencoze`.`departments` (
    `department_id` VARCHAR(36) PRIMARY KEY COMMENT '部门ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `department_name` VARCHAR(200) NOT NULL COMMENT '部门名称',
    `parent_department_id` VARCHAR(36) NULL COMMENT '父部门ID',
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒）',
    `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒）',
    `deleted_at` BIGINT UNSIGNED NULL COMMENT '删除时间（毫秒）',
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,
    FOREIGN KEY (`parent_department_id`) REFERENCES `opencoze`.`departments`(`department_id`) ON DELETE SET NULL,
    UNIQUE INDEX `uk_tenant_name` (`tenant_id`, `department_name`, `deleted_at`),
    INDEX `idx_tenant_id` (`tenant_id`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';

-- Create "user_departures" table
CREATE TABLE `opencoze`.`user_departures` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `department_id` VARCHAR(36) NOT NULL COMMENT '部门ID',
    `is_leader` BOOLEAN DEFAULT FALSE COMMENT '是否是部门领导',
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒）',
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,
    FOREIGN KEY (`department_id`) REFERENCES `opencoze`.`departments`(`department_id`) ON DELETE CASCADE,
    UNIQUE INDEX `uk_user_tenant_dept` (`user_id`, `tenant_id`, `department_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_department_id` (`department_id`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户部门关联表';
