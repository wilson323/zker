-- ================================================================================
-- ZKER 多租户系统 - 系统角色初始化脚本
-- ================================================================================
-- 版本: v1.0
-- 创建日期: 2025-01-01
-- 说明: 初始化5个预置系统角色及其权限配置
--       必须在创建 roles、data_permissions、field_permissions 表后执行
-- ================================================================================

-- ================================================================================
-- 第一部分：系统角色定义
-- ================================================================================

-- 1. 超级管理员 (Super Admin)
-- 拥有所有权限，包括租户管理、角色管理、所有资源的管理权限
INSERT INTO `roles` (
    `role_id`,
    `tenant_id`,
    `role_name`,
    `role_code`,
    `role_type`,
    `description`
) VALUES (
    'role_super_admin',
    'system',           -- system 租户（全局系统角色）
    '超级管理员',
    'super_admin',
    'system',
    '拥有所有权限，包括租户管理、角色管理、所有资源的完全访问权限'
);

-- 2. 管理员 (Admin)
-- 拥有租户管理权限，但不能管理角色
INSERT INTO `roles` (
    `role_id`,
    `tenant_id`,
    `role_name`,
    `role_code`,
    `role_type`,
    `description`
) VALUES (
    'role_admin',
    'system',
    '管理员',
    'admin',
    'system',
    '拥有租户管理权限，可以管理Bot、对话、知识库、工作流等资源，但不能管理角色和权限'
);

-- 3. 开发者 (Developer)
-- 可以创建和编辑Bot、工作流等资源，但无管理权限
INSERT INTO `roles` (
    `role_id`,
    `tenant_id`,
    `role_name`,
    `role_code`,
    `role_type`,
    `description`
) VALUES (
    'role_developer',
    'system',
    '开发者',
    'developer',
    'system',
    '可以创建和编辑Bot、对话、知识库、工作流等资源，但无管理权限'
);

-- 4. 访客 (Guest)
-- 只有只读权限
INSERT INTO `roles` (
    `role_id`,
    `tenant_id`,
    `role_name`,
    `role_code`,
    `role_type`,
    `description`
) VALUES (
    'role_guest',
    'system',
    '访客',
    'guest',
    'system',
    '只有只读权限，可以查看Bot、对话、知识库等资源，但不能创建或编辑'
);

-- 5. 租户所有者 (Tenant Owner)
-- 拥有租户内的所有权限，包括角色管理
INSERT INTO `roles` (
    `role_id`,
    `tenant_id`,
    `role_name`,
    `role_code`,
    `role_type`,
    `description`
) VALUES (
    'role_tenant_owner',
    'system',
    '租户所有者',
    'tenant_owner',
    'system',
    '拥有租户内的所有权限，包括角色管理、配额管理、订阅管理等'
);

-- ================================================================================
-- 第二部分：数据权限配置
-- ================================================================================

-- 超级管理员：所有资源的全部数据权限
INSERT INTO `data_permissions` (`permission_id`, `role_id`, `resource_type`, `scope`)
VALUES
    ('dp_super_admin_bots', 'role_super_admin', 'bots', 'ALL'),
    ('dp_super_admin_conversations', 'role_super_admin', 'conversations', 'ALL'),
    ('dp_super_admin_knowledge', 'role_super_admin', 'knowledge', 'ALL'),
    ('dp_super_admin_workflows', 'role_super_admin', 'workflows', 'ALL'),
    ('dp_super_admin_plugins', 'role_super_admin', 'plugins', 'ALL');

-- 管理员：所有资源的全部数据权限
INSERT INTO `data_permissions` (`permission_id`, `role_id`, `resource_type`, `scope`)
VALUES
    ('dp_admin_bots', 'role_admin', 'bots', 'ALL'),
    ('dp_admin_conversations', 'role_admin', 'conversations', 'ALL'),
    ('dp_admin_knowledge', 'role_admin', 'knowledge', 'ALL'),
    ('dp_admin_workflows', 'role_admin', 'workflows', 'ALL'),
    ('dp_admin_plugins', 'role_admin', 'plugins', 'ALL');

-- 开发者：所有资源的仅自己数据权限
INSERT INTO `data_permissions` (`permission_id`, `role_id`, `resource_type`, `scope`)
VALUES
    ('dp_developer_bots', 'role_developer', 'bots', 'OWN'),
    ('dp_developer_conversations', 'role_developer', 'conversations', 'OWN'),
    ('dp_developer_knowledge', 'role_developer', 'knowledge', 'OWN'),
    ('dp_developer_workflows', 'role_developer', 'workflows', 'OWN'),
    ('dp_developer_plugins', 'role_developer', 'plugins', 'OWN');

-- 访客：所有资源的仅自己数据权限（只读）
INSERT INTO `data_permissions` (`permission_id`, `role_id`, `resource_type`, `scope`)
VALUES
    ('dp_guest_bots', 'role_guest', 'bots', 'OWN'),
    ('dp_guest_conversations', 'role_guest', 'conversations', 'OWN'),
    ('dp_guest_knowledge', 'role_guest', 'knowledge', 'OWN'),
    ('dp_guest_workflows', 'role_guest', 'workflows', 'OWN');

-- 租户所有者：所有资源的全部数据权限
INSERT INTO `data_permissions` (`permission_id`, `role_id`, `resource_type`, `scope`)
VALUES
    ('dp_tenant_owner_bots', 'role_tenant_owner', 'bots', 'ALL'),
    ('dp_tenant_owner_conversations', 'role_tenant_owner', 'conversations', 'ALL'),
    ('dp_tenant_owner_knowledge', 'role_tenant_owner', 'knowledge', 'ALL'),
    ('dp_tenant_owner_workflows', 'role_tenant_owner', 'workflows', 'ALL'),
    ('dp_tenant_owner_plugins', 'role_tenant_owner', 'plugins', 'ALL');

-- ================================================================================
-- 第三部分：字段权限配置
-- ================================================================================

-- 超级管理员：所有字段可编辑
INSERT INTO `field_permissions` (`permission_id`, `role_id`, `resource_type`, `field_name`, `permission_level`)
VALUES
    -- Bots 表字段
    ('fp_super_admin_bots_id', 'role_super_admin', 'bots', 'bot_id', 'editable'),
    ('fp_super_admin_bots_name', 'role_super_admin', 'bots', 'name', 'editable'),
    ('fp_super_admin_bots_description', 'role_super_admin', 'bots', 'description', 'editable'),
    ('fp_super_admin_bots_config', 'role_super_admin', 'bots', 'config', 'editable'),
    ('fp_super_admin_bots_api_token', 'role_super_admin', 'bots', 'api_token', 'editable'),
    ('fp_super_admin_bots_creator_id', 'role_super_admin', 'bots', 'creator_id', 'readonly'),
    ('fp_super_admin_bots_created_at', 'role_super_admin', 'bots', 'created_at', 'readonly'),

    -- Conversations 表字段
    ('fp_super_admin_conversations_id', 'role_super_admin', 'conversations', 'conversation_id', 'editable'),
    ('fp_super_admin_conversations_name', 'role_super_admin', 'conversations', 'name', 'editable'),
    ('fp_super_admin_conversations_messages', 'role_super_admin', 'conversations', 'messages', 'editable'),

    -- Knowledge 表字段
    ('fp_super_admin_knowledge_id', 'role_super_admin', 'knowledge', 'knowledge_id', 'editable'),
    ('fp_super_admin_knowledge_name', 'role_super_admin', 'knowledge', 'name', 'editable'),
    ('fp_super_admin_knowledge_documents', 'role_super_admin', 'knowledge', 'documents', 'editable');

-- 管理员：敏感字段只读，其他可编辑
INSERT INTO `field_permissions` (`permission_id`, `role_id`, `resource_type`, `field_name`, `permission_level`)
VALUES
    -- Bots 表字段
    ('fp_admin_bots_id', 'role_admin', 'bots', 'bot_id', 'readonly'),
    ('fp_admin_bots_name', 'role_admin', 'bots', 'name', 'editable'),
    ('fp_admin_bots_description', 'role_admin', 'bots', 'description', 'editable'),
    ('fp_admin_bots_config', 'role_admin', 'bots', 'config', 'editable'),
    ('fp_admin_bots_api_token', 'role_admin', 'bots', 'api_token', 'hidden'),  -- 隐藏API Token
    ('fp_admin_bots_creator_id', 'role_admin', 'bots', 'creator_id', 'readonly'),
    ('fp_admin_bots_created_at', 'role_admin', 'bots', 'created_at', 'readonly');

-- 开发者：自己创建的资源可编辑，其他只读
INSERT INTO `field_permissions` (`permission_id`, `role_id`, `resource_type`, `field_name`, `permission_level`)
VALUES
    -- Bots 表字段
    ('fp_developer_bots_id', 'role_developer', 'bots', 'bot_id', 'readonly'),
    ('fp_developer_bots_name', 'role_developer', 'bots', 'name', 'editable'),
    ('fp_developer_bots_description', 'role_developer', 'bots', 'description', 'editable'),
    ('fp_developer_bots_config', 'role_developer', 'bots', 'config', 'editable'),
    ('fp_developer_bots_api_token', 'role_developer', 'bots', 'api_token', 'hidden'),
    ('fp_developer_bots_creator_id', 'role_developer', 'bots', 'creator_id', 'readonly');

-- 访客：所有字段只读
INSERT INTO `field_permissions` (`permission_id`, `role_id`, `resource_type`, `field_name`, `permission_level`)
VALUES
    -- Bots 表字段
    ('fp_guest_bots_id', 'role_guest', 'bots', 'bot_id', 'readonly'),
    ('fp_guest_bots_name', 'role_guest', 'bots', 'name', 'readonly'),
    ('fp_guest_bots_description', 'role_guest', 'bots', 'description', 'readonly'),
    ('fp_guest_bots_config', 'role_guest', 'bots', 'config', 'hidden'),      -- 隐藏配置
    ('fp_guest_bots_api_token', 'role_guest', 'bots', 'api_token', 'hidden'),  -- 隐藏API Token
    ('fp_guest_bots_creator_id', 'role_guest', 'bots', 'creator_id', 'readonly'),
    ('fp_guest_bots_created_at', 'role_guest', 'bots', 'created_at', 'readonly');

-- 租户所有者：所有字段可编辑
INSERT INTO `field_permissions` (`permission_id`, `role_id`, `resource_type`, `field_name`, `permission_level`)
VALUES
    -- Bots 表字段
    ('fp_tenant_owner_bots_id', 'role_tenant_owner', 'bots', 'bot_id', 'editable'),
    ('fp_tenant_owner_bots_name', 'role_tenant_owner', 'bots', 'name', 'editable'),
    ('fp_tenant_owner_bots_description', 'role_tenant_owner', 'bots', 'description', 'editable'),
    ('fp_tenant_owner_bots_config', 'role_tenant_owner', 'bots', 'config', 'editable'),
    ('fp_tenant_owner_bots_api_token', 'role_tenant_owner', 'bots', 'api_token', 'editable'),
    ('fp_tenant_owner_bots_creator_id', 'role_tenant_owner', 'bots', 'creator_id', 'readonly'),
    ('fp_tenant_owner_bots_created_at', 'role_tenant_owner', 'bots', 'created_at', 'readonly');

-- ================================================================================
-- 第四部分：为每个租户自动创建角色副本
-- ================================================================================

-- 创建存储过程：为租户初始化角色
DELIMITER $$

DROP PROCEDURE IF EXISTS `InitTenantRoles`$$

CREATE PROCEDURE `InitTenantRoles`(
    IN p_tenant_id VARCHAR(36)
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;

    START TRANSACTION;

    -- 1. 为租户创建5个系统角色副本
    INSERT INTO `roles` (`role_id`, `tenant_id`, `role_name`, `role_code`, `role_type`, `description`, `created_at`, `updated_at`)
    SELECT
        CONCAT(p_tenant_id, '_', SUBSTR(role_id, 6)) as role_id,  -- 租户ID_角色ID
        p_tenant_id as tenant_id,
        role_name,
        role_code,
        'system' as role_type,
        description,
        UNIX_TIMESTAMP(NOW()) * 1000 as created_at,
        UNIX_TIMESTAMP(NOW()) * 1000 as updated_at
    FROM `roles`
    WHERE tenant_id = 'system';

    -- 2. 为租户角色复制数据权限配置
    INSERT INTO `data_permissions` (
        `permission_id`,
        `role_id`,
        `resource_type`,
        `scope`,
        `created_at`,
        `updated_at`
    )
    SELECT
        CONCAT(p_tenant_id, '_', SUBSTR(dp.permission_id, 6)) as permission_id,
        (SELECT role_id FROM `roles` WHERE tenant_id = p_tenant_id AND role_code = sys_role.role_code LIMIT 1) as role_id,
        dp.resource_type,
        dp.scope,
        UNIX_TIMESTAMP(NOW()) * 1000 as created_at,
        UNIX_TIMESTAMP(NOW()) * 1000 as updated_at
    FROM `data_permissions` dp
    JOIN `roles` sys_role ON dp.role_id = sys_role.role_id
    WHERE sys_role.tenant_id = 'system';

    -- 3. 为租户角色复制字段权限配置
    INSERT INTO `field_permissions` (
        `permission_id`,
        `role_id`,
        `resource_type`,
        `field_name`,
        `permission_level`,
        `created_at`,
        `updated_at`
    )
    SELECT
        CONCAT(p_tenant_id, '_', SUBSTR(fp.permission_id, 6)) as permission_id,
        (SELECT role_id FROM `roles` WHERE tenant_id = p_tenant_id AND role_code = sys_role.role_code LIMIT 1) as role_id,
        fp.resource_type,
        fp.field_name,
        fp.permission_level,
        UNIX_TIMESTAMP(NOW()) * 1000 as created_at,
        UNIX_TIMESTAMP(NOW()) * 1000 as updated_at
    FROM `field_permissions` fp
    JOIN `roles` sys_role ON fp.role_id = sys_role.role_id
    WHERE sys_role.tenant_id = 'system';

    COMMIT;

END$$

DELIMITER ;

-- ================================================================================
-- 第五部分：验证脚本
-- ================================================================================

-- 验证系统角色是否创建成功
SELECT
    '系统角色初始化验证' as check_name,
    role_code as '角色代码',
    role_name as '角色名称',
    role_type as '角色类型',
    description as '描述'
FROM `roles`
WHERE tenant_id = 'system'
ORDER BY role_code;

-- 验证每个角色的数据权限数量
SELECT
    r.role_code as '角色代码',
    r.role_name as '角色名称',
    COUNT(DISTINCT dp.resource_type) as '数据权限数量',
    GROUP_CONCAT(DISTINCT dp.resource_type, ':', dp.scope) as '权限详情'
FROM `roles` r
LEFT JOIN `data_permissions` dp ON r.role_id = dp.role_id
WHERE r.tenant_id = 'system'
GROUP BY r.role_id, r.role_code, r.role_name
ORDER BY r.role_code;

-- 验证每个角色的字段权限数量
SELECT
    r.role_code as '角色代码',
    r.role_name as '角色名称',
    COUNT(DISTINCT CONCAT(fp.resource_type, '.', fp.field_name)) as '字段权限数量',
    SUM(CASE WHEN fp.permission_level = 'hidden' THEN 1 ELSE 0 END) as '隐藏字段数',
    SUM(CASE WHEN fp.permission_level = 'readonly' THEN 1 ELSE 0 END) as '只读字段数',
    SUM(CASE WHEN fp.permission_level = 'editable' THEN 1 ELSE 0 END) as '可编辑字段数'
FROM `roles` r
LEFT JOIN `field_permissions` fp ON r.role_id = fp.role_id
WHERE r.tenant_id = 'system'
GROUP BY r.role_id, r.role_code, r.role_name
ORDER BY r.role_code;

-- ================================================================================
-- 使用说明
-- ================================================================================
-- 1. 为新租户初始化角色（示例：为租户 tenant_123 初始化角色）
-- CALL InitTenantRoles('tenant_123');
--
-- 2. 为现有所有租户初始化角色
-- INSERT INTO tenant_role_init (tenant_id)
-- SELECT DISTINCT tenant_id FROM tenants
-- WHERE tenant_id NOT IN (SELECT tenant_id FROM roles WHERE tenant_id != 'system');
--
-- 3. 验证租户角色是否创建成功
-- SELECT tenant_id, COUNT(*) as role_count
-- FROM roles
-- WHERE tenant_id != 'system'
-- GROUP BY tenant_id;
-- ================================================================================

-- 初始化完成标记
-- 执行时间: ____________________
-- 执行人: ____________________
-- 验证结果: [] 全部通过  [ ] 存在问题需要修复
