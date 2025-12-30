-- ================================================================
-- ZKER 系统角色初始化脚本
-- ================================================================
-- 功能：初始化4个系统预置角色及其权限
-- 版本：v1.0
-- 创建日期：2025-01-01
-- ================================================================
-- 说明：
--   1. 本脚本初始化4个系统角色：TenantAdmin, Developer, Operator, Viewer
--   2. 系统角色的 tenant_id 为 'SYSTEM'，作为所有租户的模板
--   3. 每个角色都配置了数据权限和字段权限
--   4. 新租户创建时，应复制这些系统角色并替换 tenant_id
-- ================================================================

-- ================================================================
-- 第1步：插入4个系统角色
-- ================================================================

INSERT INTO `roles` (
    `role_id`,
    `tenant_id`,
    `role_name`,
    `role_code`,
    `role_type`,
    `description`,
    `created_at`,
    `updated_at`
) VALUES
-- 1. TenantAdmin (租户管理员)
(
    'role-tenant-admin',
    'SYSTEM',
    'Tenant Administrator',
    'TENANT_ADMIN',
    'system',
    'Full tenant administration with all data permissions and all fields editable',
    UNIX_TIMESTAMP(NOW()) * 1000,
    UNIX_TIMESTAMP(NOW()) * 1000
),
-- 2. Developer (开发者)
(
    'role-developer',
    'SYSTEM',
    'Developer',
    'DEVELOPER',
    'system',
    'Developer with department-level data permissions and development fields editable',
    UNIX_TIMESTAMP(NOW()) * 1000,
    UNIX_TIMESTAMP(NOW()) * 1000
),
-- 3. Operator (运营人员)
(
    'role-operator',
    'SYSTEM',
    'Operator',
    'OPERATOR',
    'system',
    'Operator with own data permissions only and readonly access to most fields',
    UNIX_TIMESTAMP(NOW()) * 1000,
    UNIX_TIMESTAMP(NOW()) * 1000
),
-- 4. Viewer (查看者)
(
    'role-viewer',
    'SYSTEM',
    'Viewer',
    'VIEWER',
    'system',
    'Viewer with own data permissions only and readonly access to all fields',
    UNIX_TIMESTAMP(NOW()) * 1000,
    UNIX_TIMESTAMP(NOW()) * 1000
);

-- ================================================================
-- 第2步：配置数据权限 (Data Permissions)
-- ================================================================

-- TenantAdmin: ALL for all resources
INSERT INTO `data_permissions` (`permission_id`, `role_id`, `resource_type`, `scope`, `created_at`, `updated_at`)
SELECT
    UUID() AS permission_id,
    'role-tenant-admin' AS role_id,
    resource_type,
    'ALL' AS scope,
    UNIX_TIMESTAMP(NOW()) * 1000 AS created_at,
    UNIX_TIMESTAMP(NOW()) * 1000 AS updated_at
FROM (
    SELECT 'bots' AS resource_type UNION
    SELECT 'conversations' UNION
    SELECT 'knowledge' UNION
    SELECT 'workflows' UNION
    SELECT 'plugins'
) AS resources;

-- Developer: DEPARTMENT for bots/workflows, OWN for others
INSERT INTO `data_permissions` (`permission_id`, `role_id`, `resource_type`, `scope`, `created_at`, `updated_at`)
VALUES
(UUID(), 'role-developer', 'bots', 'DEPARTMENT', UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
(UUID(), 'role-developer', 'workflows', 'DEPARTMENT', UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
(UUID(), 'role-developer', 'conversations', 'OWN', UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
(UUID(), 'role-developer', 'knowledge', 'OWN', UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
(UUID(), 'role-developer', 'plugins', 'OWN', UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000);

-- Operator: OWN for all resources
INSERT INTO `data_permissions` (`permission_id`, `role_id`, `resource_type`, `scope`, `created_at`, `updated_at`)
SELECT
    UUID() AS permission_id,
    'role-operator' AS role_id,
    resource_type,
    'OWN' AS scope,
    UNIX_TIMESTAMP(NOW()) * 1000 AS created_at,
    UNIX_TIMESTAMP(NOW()) * 1000 AS updated_at
FROM (
    SELECT 'bots' AS resource_type UNION
    SELECT 'conversations' UNION
    SELECT 'knowledge' UNION
    SELECT 'workflows' UNION
    SELECT 'plugins'
) AS resources;

-- Viewer: OWN for all resources (readonly enforced via field permissions)
INSERT INTO `data_permissions` (`permission_id`, `role_id`, `resource_type`, `scope`, `created_at`, `updated_at`)
SELECT
    UUID() AS permission_id,
    'role-viewer' AS role_id,
    resource_type,
    'OWN' AS scope,
    UNIX_TIMESTAMP(NOW()) * 1000 AS created_at,
    UNIX_TIMESTAMP(NOW()) * 1000 AS updated_at
FROM (
    SELECT 'bots' AS resource_type UNION
    SELECT 'conversations' UNION
    SELECT 'knowledge' UNION
    SELECT 'workflows' UNION
    SELECT 'plugins'
) AS resources;

-- ================================================================
-- 第3步：配置字段权限 (Field Permissions)
-- ================================================================

-- 敏感字段定义
-- 这些字段对所有非管理员角色限制访问：
-- - api_key: API密钥
-- - secret_key: 密钥
-- - webhook_url: Webhook URL
-- - auth_token: 认证令牌
-- - private_key: 私钥

-- TenantAdmin: editable for all fields (no entries needed, default is editable)

-- Developer: readonly for sensitive fields
INSERT INTO `field_permissions` (`permission_id`, `role_id`, `resource_type`, `field_name`, `permission_level`, `created_at`, `updated_at`)
SELECT
    UUID() AS permission_id,
    'role-developer' AS role_id,
    resource_type,
    field_name,
    'readonly' AS permission_level,
    UNIX_TIMESTAMP(NOW()) * 1000 AS created_at,
    UNIX_TIMESTAMP(NOW()) * 1000 AS updated_at
FROM (
    SELECT 'bots' AS resource_type, 'api_key' AS field_name UNION
    SELECT 'bots', 'secret_key' UNION
    SELECT 'bots', 'webhook_url' UNION
    SELECT 'workflows', 'api_key' UNION
    SELECT 'workflows', 'auth_token' UNION
    SELECT 'plugins', 'api_key' UNION
    SELECT 'plugins', 'secret_key'
) AS sensitive_fields;

-- Operator: hidden for sensitive fields, readonly for others
INSERT INTO `field_permissions` (`permission_id`, `role_id`, `resource_type`, `field_name`, `permission_level`, `created_at`, `updated_at`)
SELECT
    UUID() AS permission_id,
    'role-operator' AS role_id,
    resource_type,
    field_name,
    CASE
        WHEN field_name IN ('api_key', 'secret_key', 'webhook_url', 'auth_token', 'private_key') THEN 'hidden'
        ELSE 'readonly'
    END AS permission_level,
    UNIX_TIMESTAMP(NOW()) * 1000 AS created_at,
    UNIX_TIMESTAMP(NOW()) * 1000 AS updated_at
FROM (
    -- 常见字段（所有只读）
    SELECT 'bots' AS resource_type, 'name' AS field_name UNION
    SELECT 'bots', 'description' UNION
    SELECT 'bots', 'status' UNION
    SELECT 'bots', 'created_at' UNION
    SELECT 'bots', 'updated_at' UNION
    -- 敏感字段（隐藏）
    SELECT 'bots', 'api_key' UNION
    SELECT 'bots', 'secret_key' UNION
    SELECT 'bots', 'webhook_url' UNION
    SELECT 'conversations', 'message' UNION
    SELECT 'conversations', 'created_at' UNION
    SELECT 'knowledge', 'name' UNION
    SELECT 'knowledge', 'document_count' UNION
    SELECT 'workflows', 'name' UNION
    SELECT 'workflows', 'status' UNION
    SELECT 'workflows', 'created_at' UNION
    SELECT 'workflows', 'api_key' UNION
    SELECT 'workflows', 'auth_token' UNION
    SELECT 'plugins', 'name' UNION
    SELECT 'plugins', 'version' UNION
    SELECT 'plugins', 'api_key' UNION
    SELECT 'plugins', 'secret_key'
) AS operator_fields;

-- Viewer: readonly for all fields (enforce readonly on key fields)
INSERT INTO `field_permissions` (`permission_id`, `role_id`, `resource_type`, `field_name`, `permission_level`, `created_at`, `updated_at`)
SELECT
    UUID() AS permission_id,
    'role-viewer' AS role_id,
    resource_type,
    field_name,
    'readonly' AS permission_level,
    UNIX_TIMESTAMP(NOW()) * 1000 AS created_at,
    UNIX_TIMESTAMP(NOW()) * 1000 AS updated_at
FROM (
    SELECT 'bots' AS resource_type, 'name' AS field_name UNION
    SELECT 'bots', 'description' UNION
    SELECT 'bots', 'status' UNION
    SELECT 'bots', 'created_at' UNION
    SELECT 'bots', 'updated_at' UNION
    SELECT 'conversations', 'message' UNION
    SELECT 'conversations', 'created_at' UNION
    SELECT 'knowledge', 'name' UNION
    SELECT 'knowledge', 'document_count' UNION
    SELECT 'knowledge', 'created_at' UNION
    SELECT 'workflows', 'name' UNION
    SELECT 'workflows', 'status' UNION
    SELECT 'workflows', 'created_at' UNION
    SELECT 'plugins', 'name' UNION
    SELECT 'plugins', 'version' UNION
    SELECT 'plugins', 'created_at'
) AS viewer_fields;

-- Viewer: hidden for sensitive fields
INSERT INTO `field_permissions` (`permission_id`, `role_id`, `resource_type`, `field_name`, `permission_level`, `created_at`, `updated_at`)
SELECT
    UUID() AS permission_id,
    'role-viewer' AS role_id,
    resource_type,
    field_name,
    'hidden' AS permission_level,
    UNIX_TIMESTAMP(NOW()) * 1000 AS created_at,
    UNIX_TIMESTAMP(NOW()) * 1000 AS updated_at
FROM (
    SELECT 'bots' AS resource_type, 'api_key' AS field_name UNION
    SELECT 'bots', 'secret_key' UNION
    SELECT 'bots', 'webhook_url' UNION
    SELECT 'workflows', 'api_key' UNION
    SELECT 'workflows', 'auth_token' UNION
    SELECT 'plugins', 'api_key' UNION
    SELECT 'plugins', 'secret_key'
) AS viewer_sensitive_fields;

-- ================================================================
-- 验证脚本执行结果
-- ================================================================

-- 检查系统角色数量
-- 应返回4
SELECT COUNT(*) AS system_role_count FROM `roles`
WHERE `tenant_id` = 'SYSTEM' AND `role_type` = 'system';

-- 检查数据权限配置
-- 应返回20（4个角色 × 5个资源类型）
SELECT COUNT(*) AS data_permission_count FROM `data_permissions` dp
JOIN `roles` r ON dp.role_id = r.role_id
WHERE r.tenant_id = 'SYSTEM';

-- 检查字段权限配置
-- 应返回非零值（Developer: 7, Operator: 20, Viewer: 22）
SELECT
    r.role_code,
    COUNT(*) AS field_permission_count
FROM `field_permissions` fp
JOIN `roles` r ON fp.role_id = r.role_id
WHERE r.tenant_id = 'SYSTEM'
GROUP BY r.role_code
ORDER BY r.role_code;

-- ================================================================
-- 使用说明
-- ================================================================
-- 1. 本脚本仅在数据库初始化时执行一次
-- 2. 新租户创建时，应复制这些系统角色并替换 tenant_id
-- 3. 系统角色不应被修改或删除
-- 4. 用户自定义角色应基于 system 角色，但 role_type = 'custom'
-- ================================================================
