-- ============================================================
-- 系统初始化数据脚本
-- 版本: v1.0
-- 日期: 2025-01-03
-- 说明: 初始化系统预置角色、权限和模型定价
-- ============================================================

-- ============================================================
-- 第一部分：初始化系统预置角色
-- ============================================================

-- 注意：这些角色需要在每个租户创建时初始化
-- 下面的SQL仅为示例，实际应用中应该在创建租户时动态插入

-- 1. 租户所有者 (tenant_owner) - 拥有租户内所有资源的完整权限
INSERT INTO `roles` (
    `role_id`, `tenant_id`, `role_name`, `role_code`, `description`,
    `role_type`, `is_system`, `is_enabled`,
    `created_at`, `updated_at`
) VALUES (
    'role_tenant_owner_TEMPLATE', 'TENANT_ID_TEMPLATE', '租户所有者', 'tenant_owner',
    '拥有租户内所有资源的完整权限',
    'system', TRUE, TRUE,
    UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000
) ON DUPLICATE KEY UPDATE `role_name` = VALUES(`role_name`);

-- 2. 租户管理员 (tenant_admin) - 拥有租户内部门及以下资源的完整权限
INSERT INTO `roles` (
    `role_id`, `tenant_id`, `role_name`, `role_code`, `description`,
    `role_type`, `is_system`, `is_enabled`,
    `created_at`, `updated_at`
) VALUES (
    'role_tenant_admin_TEMPLATE', 'TENANT_ID_TEMPLATE', '租户管理员', 'tenant_admin',
    '拥有租户内部门及以下资源的完整权限',
    'system', TRUE, TRUE,
    UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000
) ON DUPLICATE KEY UPDATE `role_name` = VALUES(`role_name`);

-- 3. 普通成员 (tenant_member) - 仅能访问自己创建的资源
INSERT INTO `roles` (
    `role_id`, `tenant_id`, `role_name`, `role_code`, `description`,
    `role_type`, `is_system`, `is_enabled`,
    `created_at`, `updated_at`
) VALUES (
    'role_tenant_member_TEMPLATE', 'TENANT_ID_TEMPLATE', '普通成员', 'tenant_member',
    '仅能访问自己创建的资源',
    'system', TRUE, TRUE,
    UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000
) ON DUPLICATE KEY UPDATE `role_name` = VALUES(`role_name`);

-- 4. 只读成员 (tenant_viewer) - 仅能查看自己创建的资源
INSERT INTO `roles` (
    `role_id`, `tenant_id`, `role_name`, `role_code`, `description`,
    `role_type`, `is_system`, `is_enabled`,
    `created_at`, `updated_at`
) VALUES (
    'role_tenant_viewer_TEMPLATE', 'TENANT_ID_TEMPLATE', '只读成员', 'tenant_viewer',
    '仅能查看自己创建的资源',
    'system', TRUE, TRUE,
    UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000
) ON DUPLICATE KEY UPDATE `role_name` = VALUES(`role_name`);

-- 5. 运营人员 (tenant_operator) - 运营人员角色
INSERT INTO `roles` (
    `role_id`, `tenant_id`, `role_name`, `role_code`, `description`,
    `role_type`, `is_system`, `is_enabled`,
    `created_at`, `updated_at`
) VALUES (
    'role_tenant_operator_TEMPLATE', 'TENANT_ID_TEMPLATE', '运营人员', 'tenant_operator',
    '运营人员角色，可查看和操作租户数据',
    'system', TRUE, TRUE,
    UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000
) ON DUPLICATE KEY UPDATE `role_name` = VALUES(`role_name`);

-- ============================================================
-- 第二部分：初始化系统预置权限
-- ============================================================

-- Bot相关权限
INSERT INTO `permissions` (`permission_id`, `permission_code`, `name`, `description`, `type`, `resource_type`, `operation`, `is_system`, `created_at`, `updated_at`) VALUES
('perm_bot_create', 'bot:create', '创建Bot', '创建新Bot', 'operation', 'bot', 'create', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_bot_read', 'bot:read', '查看Bot', '查看Bot详情', 'operation', 'bot', 'read', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_bot_update', 'bot:update', '更新Bot', '更新Bot配置', 'operation', 'bot', 'update', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_bot_delete', 'bot:delete', '删除Bot', '删除Bot', 'operation', 'bot', 'delete', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_bot_export', 'bot:export', '导出Bot', '导出Bot数据', 'operation', 'bot', 'export', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- 租户管理权限
INSERT INTO `permissions` (`permission_id`, `permission_code`, `name`, `description`, `type`, `resource_type`, `operation`, `is_system`, `created_at`, `updated_at`) VALUES
('perm_tenant_create', 'tenant:create', '创建租户', '创建新租户', 'operation', 'tenant', 'create', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_tenant_read', 'tenant:read', '查看租户', '查看租户详情', 'operation', 'tenant', 'read', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_tenant_update', 'tenant:update', '更新租户', '更新租户信息', 'operation', 'tenant', 'update', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_tenant_delete', 'tenant:delete', '删除租户', '删除租户', 'operation', 'tenant', 'delete', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- 用户管理权限
INSERT INTO `permissions` (`permission_id`, `permission_code`, `name`, `description`, `type`, `resource_type`, `operation`, `is_system`, `created_at`, `updated_at`) VALUES
('perm_user_create', 'user:create', '创建用户', '创建新用户', 'operation', 'user', 'create', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_user_read', 'user:read', '查看用户', '查看用户详情', 'operation', 'user', 'read', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_user_update', 'user:update', '更新用户', '更新用户信息', 'operation', 'user', 'update', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_user_delete', 'user:delete', '删除用户', '删除用户', 'operation', 'user', 'delete', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_user_role', 'user:role', '分配角色', '为用户分配角色', 'operation', 'user', 'approve', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- 对话管理权限
INSERT INTO `permissions` (`permission_id`, `permission_code`, `name`, `description`, `type`, `resource_type`, `operation`, `is_system`, `created_at`, `updated_at`) VALUES
('perm_conversation_create', 'conversation:create', '创建对话', '创建新对话', 'operation', 'conversation', 'create', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_conversation_read', 'conversation:read', '查看对话', '查看对话详情', 'operation', 'conversation', 'read', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_conversation_update', 'conversation:update', '更新对话', '更新对话信息', 'operation', 'conversation', 'update', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_conversation_delete', 'conversation:delete', '删除对话', '删除对话', 'operation', 'conversation', 'delete', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- 工作流管理权限
INSERT INTO `permissions` (`permission_id`, `permission_code`, `name`, `description`, `type`, `resource_type`, `operation`, `is_system`, `created_at`, `updated_at`) VALUES
('perm_workflow_create', 'workflow:create', '创建工作流', '创建新工作流', 'operation', 'workflow', 'create', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_workflow_read', 'workflow:read', '查看工作流', '查看工作流详情', 'operation', 'workflow', 'read', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_workflow_update', 'workflow:update', '更新工作流', '更新工作流配置', 'operation', 'workflow', 'update', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_workflow_delete', 'workflow:delete', '删除工作流', '删除工作流', 'operation', 'workflow', 'delete', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_workflow_publish', 'workflow:publish', '发布工作流', '发布工作流', 'operation', 'workflow', 'approve', TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- 敏感字段权限
INSERT INTO `permissions` (`permission_id`, `permission_code`, `name`, `description`, `type`, `resource_type`, `operation`, `field_name`, `is_sensitive`, `is_system`, `created_at`, `updated_at`) VALUES
('perm_field_phone', 'field:phone', '手机号字段', '访问用户手机号', 'field', 'user', 'read', 'phone', TRUE, TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_field_email', 'field:email', '邮箱字段', '访问用户邮箱', 'field', 'user', 'read', 'email', TRUE, TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_field_idcard', 'field:idcard', '身份证字段', '访问用户身份证', 'field', 'user', 'read', 'idcard', TRUE, TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_field_salary', 'field:salary', '薪资字段', '访问用户薪资', 'field', 'user', 'read', 'salary', TRUE, TRUE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_field_address', 'field:address', '地址字段', '访问用户地址', 'field', 'user', 'read', 'address', TRUE, FALSE, UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- ============================================================
-- 第三部分：初始化模型定价配置
-- ============================================================

-- 注意：这些定价配置存储在 PricingEngine 代码中，此处仅为文档说明
-- 实际应用中，可以通过配置表来动态调整定价

-- OpenAI 模型定价
-- openai/gpt-4: 输入 ¥0.03/1K tokens, 输出 ¥0.06/1K tokens
-- openai/gpt-4-32k: 输入 ¥0.06/1K tokens, 输出 ¥0.12/1K tokens
-- openai/gpt-3.5-turbo: 输入 ¥0.003/1K tokens, 输出 ¥0.006/1K tokens
-- openai/gpt-3.5-turbo-16k: 输入 ¥0.004/1K tokens, 输出 ¥0.008/1K tokens

-- Anthropic 模型定价
-- anthropic/claude-3-opus: 输入 ¥0.09/1K tokens, 输出 ¥0.27/1K tokens
-- anthropic/claude-3-sonnet: 输入 ¥0.015/1K tokens, 输出 ¥0.045/1K tokens
-- anthropic/claude-3-haiku: 输入 ¥0.0025/1K tokens, 输出 ¥0.0125/1K tokens

-- 通义千问 模型定价
-- qwen/qwen-max: 输入 ¥0.02/1K tokens, 输出 ¥0.06/1K tokens
-- qwen/qwen-plus: 输入 ¥0.004/1K tokens, 输出 ¥0.012/1K tokens
-- qwen/qwen-turbo: 输入 ¥0.001/1K tokens, 输出 ¥0.002/1K tokens

-- 百度文心 模型定价
-- baidu/ernie-bot-4: 输入 ¥0.012/1K tokens, 输出 ¥0.012/1K tokens
-- baidu/ernie-bot-turbo: 输入 ¥0.008/1K tokens, 输出 ¥0.008/1K tokens

-- 智谱 ChatGLM 模型定价
-- zhipu/chatglm-turbo: 输入 ¥0.005/1K tokens, 输出 ¥0.005/1K tokens
-- zhipu/chatglm-pro: 输入 ¥0.01/1K tokens, 输出 ¥0.01/1K tokens

-- ============================================================
-- 第四部分：示例数据权限配置
-- ============================================================

-- 为 tenant_owner 角色配置数据权限（所有资源全部访问）
-- 注意：实际应用中应该在创建租户时动态插入

INSERT INTO `data_permissions` (
    `permission_id`, `tenant_id`, `role_id`,
    `permission_name`, `permission_code`,
    `scope`, `resource_type`,
    `created_at`, `updated_at`
) VALUES
('perm_data_bot_all_TEMPLATE', 'TENANT_ID_TEMPLATE', 'role_tenant_owner_TEMPLATE',
 'Bot全部数据权限', 'data:bot:all',
 'all', 'bot',
 UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_data_conv_all_TEMPLATE', 'TENANT_ID_TEMPLATE', 'role_tenant_owner_TEMPLATE',
 '对话全部数据权限', 'data:conversation:all',
 'all', 'conversation',
 UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000)
ON DUPLICATE KEY UPDATE `permission_name` = VALUES(`permission_name`);

-- 为 tenant_member 角色配置数据权限（仅自己创建的资源）
INSERT INTO `data_permissions` (
    `permission_id`, `tenant_id`, `role_id`,
    `permission_name`, `permission_code`,
    `scope`, `resource_type`,
    `created_at`, `updated_at`
) VALUES
('perm_data_bot_own_TEMPLATE', 'TENANT_ID_TEMPLATE', 'role_tenant_member_TEMPLATE',
 'Bot自己创建的数据权限', 'data:bot:own',
 'own', 'bot',
 UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_data_conv_own_TEMPLATE', 'TENANT_ID_TEMPLATE', 'role_tenant_member_TEMPLATE',
 '对话自己创建的数据权限', 'data:conversation:own',
 'own', 'conversation',
 UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000)
ON DUPLICATE KEY UPDATE `permission_name` = VALUES(`permission_name`);

-- ============================================================
-- 第五部分：示例字段权限配置
-- ============================================================

-- 为 tenant_viewer 角色配置敏感字段隐藏权限
INSERT INTO `field_permissions` (
    `permission_id`, `tenant_id`, `role_id`,
    `permission_name`, `permission_code`,
    `resource_type`, `field_name`, `permission_level`,
    `created_at`, `updated_at`
) VALUES
('perm_field_api_key_hidden_TEMPLATE', 'TENANT_ID_TEMPLATE', 'role_tenant_viewer_TEMPLATE',
 'API Key隐藏权限', 'field:api_key:hidden',
 'bots', 'api_key', 'hidden',
 UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('perm_field_secret_key_hidden_TEMPLATE', 'TENANT_ID_TEMPLATE', 'role_tenant_viewer_TEMPLATE',
 'Secret Key隐藏权限', 'field:secret_key:hidden',
 'bots', 'secret_key', 'hidden',
 UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000)
ON DUPLICATE KEY UPDATE `permission_name` = VALUES(`permission_name`);

-- ============================================================
-- 初始化完成
-- ============================================================

-- 说明：
-- 1. 上述SQL中的 'TENANT_ID_TEMPLATE' 应该替换为实际的租户ID
-- 2. 在创建新租户时，应该动态执行这些初始化SQL
-- 3. 可以通过存储过程或应用代码来实现自动初始化
-- ============================================================
