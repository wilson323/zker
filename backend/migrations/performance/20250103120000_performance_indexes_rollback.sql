-- ============================================
-- ZKER 性能优化 - 索引回滚脚本
-- 版本: v1.0.0
-- 生成时间: 2025-01-03
-- 用途: 回滚性能索引优化
-- ============================================

-- 警告: 此脚本将删除所有性能优化索引,请谨慎执行!
-- 执行前请确认: 1) 已备份数据库 2) 确实需要回滚 3) 业务低峰期

-- ============================================
-- 第一部分: 知识库模块索引回滚
-- ============================================

-- 1.1 knowledge 表索引
DROP INDEX IF EXISTS idx_knowledge_app_space_status_created ON knowledge;
DROP INDEX IF EXISTS idx_knowledge_status_updated ON knowledge;
DROP INDEX IF EXISTS uk_knowledge_space_name ON knowledge;

-- 1.2 knowledge_document 表索引
DROP INDEX IF EXISTS idx_document_knowledge_status_created ON knowledge_document;
DROP INDEX IF EXISTS idx_document_cover ON knowledge_document;
DROP INDEX IF EXISTS idx_document_type_status ON knowledge_document;
DROP INDEX IF EXISTS uk_document_knowledge_name ON knowledge_document;

-- 1.3 knowledge_document_slice 表索引
DROP INDEX IF EXISTS idx_slice_document_created ON knowledge_document_slice;
DROP INDEX IF EXISTS idx_slice_knowledge_sequence ON knowledge_document_slice;
DROP INDEX IF EXISTS idx_slice_knowledge_hit ON knowledge_document_slice;

-- ============================================
-- 第二部分: 对话模块索引回滚
-- ============================================

-- 2.1 message 表索引
DROP INDEX IF EXISTS idx_message_conversation_status_created ON message;
DROP INDEX IF EXISTS idx_message_run_status_created ON message;
DROP INDEX IF EXISTS idx_message_conversation_created_cursor ON message;
DROP INDEX IF EXISTS idx_message_type_created ON message;
DROP INDEX IF EXISTS idx_message_agent_created ON message;

-- 2.2 conversation 表索引
DROP INDEX IF EXISTS idx_conversation_tenant_space_updated ON conversation;
DROP INDEX IF EXISTS idx_conversation_status_created ON conversation;
DROP INDEX IF EXISTS idx_conversation_agent_created ON conversation;

-- ============================================
-- 第三部分: 插件模块索引回滚
-- ============================================

-- 3.1 plugin 表索引
DROP INDEX IF EXISTS idx_plugin_space_created ON plugin;
DROP INDEX IF EXISTS idx_plugin_type_status ON plugin;
DROP INDEX IF EXISTS uk_plugin_space_developer_name ON plugin;

-- ============================================
-- 第四部分: 工作流模块索引回滚
-- ============================================

-- 4.1 workflow 相关表索引
DROP INDEX IF EXISTS idx_workflow_space_status_created ON workflow;
DROP INDEX IF EXISTS idx_workflow_execution_workflow_status_created ON workflow_execution;
DROP INDEX IF EXISTS idx_node_execution_run_created ON node_execution;
DROP INDEX IF EXISTS idx_workflow_workflow_id_version ON workflow;

-- ============================================
-- 第五部分: 权限模块索引回滚
-- ============================================

-- 5.1 permission 相关表索引
DROP INDEX IF EXISTS idx_data_permission_role_resource ON data_permission;
DROP INDEX IF EXISTS idx_user_role_user_role ON user_roles;
DROP INDEX IF EXISTS idx_field_permission_role_resource ON field_permission;

-- ============================================
-- 第六部分: 租户模块索引回滚
-- ============================================

-- 6.1 tenant 相关表索引
DROP INDEX IF EXISTS idx_user_tenant_status_created ON users;
DROP INDEX IF EXISTS idx_quota_tenant_type ON quota;

-- ============================================
-- 第七部分: 其他业务表索引回滚
-- ============================================

-- 7.1 其他业务表索引
DROP INDEX IF EXISTS idx_api_key_space_status ON api_key;
DROP INDEX IF EXISTS idx_database_space_agent ON database_info;
DROP INDEX IF EXISTS idx_variable_instance_space_key ON variable_instance;
DROP INDEX IF EXISTS idx_template_space_type ON template;
DROP INDEX IF EXISTS idx_shortcut_space_agent ON shortcut_command;

-- ============================================
-- 执行说明
-- ============================================
--
-- 1. 执行前准备:
--    - 确认确实需要回滚(索引删除后性能会下降)
--    - 备份数据库结构: mysqldump -u root -p --no-data coze_studio > schema_backup.sql
--
-- 2. 执行方式:
--    mysql -u root -p coze_studio < 20250103120000_performance_indexes_rollback.sql
--
-- 3. 执行后验证:
--    - 检查索引是否删除成功
--    - 运行性能测试确认影响
--    - 如有需要,重新执行优化脚本
--
-- ============================================
-- 版本历史
-- ============================================
-- v1.0.0 (2025-01-03): 初始版本,支持回滚28个性能优化索引
-- ============================================
