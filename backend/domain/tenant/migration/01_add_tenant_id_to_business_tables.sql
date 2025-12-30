-- ================================================================================
-- ZKER 多租户迁移 - tenant_id 字段添加脚本
-- ================================================================================
-- 版本: v1.0
-- 创建日期: 2025-01-01
-- 说明: 为业务表添加 tenant_id 字段,实现多租户数据隔离
--       执行前请务必备份数据库!
-- ================================================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ================================================================================
-- 第一部分: 应用管理相关表 (3张表)
-- ================================================================================

-- 1.1 应用草稿表 (app_draft)
ALTER TABLE `app_draft`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `space_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 1.2 应用发布记录表 (app_release_record)
ALTER TABLE `app_release_record`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `space_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 1.3 应用对话模板草稿 (app_conversation_template_draft)
ALTER TABLE `app_conversation_template_draft`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `space_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- ================================================================================
-- 第二部分: 对话和消息相关表 (2张表)
-- ================================================================================

-- 2.1 对话表 (conversation)
ALTER TABLE `conversation`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `connector_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 2.2 消息表 (message)
ALTER TABLE `message`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `conversation_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- ================================================================================
-- 第三部分: 知识库相关表 (3张表)
-- ================================================================================

-- 3.1 知识库表 (knowledge)
ALTER TABLE `knowledge`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `space_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 3.2 知识库文档表 (knowledge_document)
ALTER TABLE `knowledge_document`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `knowledge_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 3.3 知识库文档分块表 (knowledge_document_slice)
ALTER TABLE `knowledge_document_slice`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `knowledge_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- ================================================================================
-- 第四部分: 文件和资源相关表 (1张表)
-- ================================================================================

-- 4.1 文件表 (files)
ALTER TABLE `files`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `coze_account_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- ================================================================================
-- 第五部分: Agent 工具相关表 (3张表)
-- ================================================================================

-- 5.1 Agent 工具草稿表 (agent_tool_draft)
ALTER TABLE `agent_tool_draft`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `agent_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 5.2 Agent 工具版本表 (agent_tool_version)
ALTER TABLE `agent_tool_version`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `agent_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 5.3 Agent 数据库关联表 (agent_to_database)
ALTER TABLE `agent_to_database`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `agent_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- ================================================================================
-- 第六部分: 工作流相关表 (2张表)
-- ================================================================================

-- 6.1 对话流角色配置表 (chat_flow_role_config)
ALTER TABLE `chat_flow_role_config`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `workflow_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 6.2 连接器工作流版本表 (connector_workflow_version)
ALTER TABLE `connector_workflow_version`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `app_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- ================================================================================
-- 第七部分: 数据库相关表 (2张表)
-- ================================================================================

-- 7.1 在线数据库信息表 (online_database_info)
ALTER TABLE `online_database_info`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `app_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 7.2 草稿数据库信息表 (draft_database_info)
ALTER TABLE `draft_database_info`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `space_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- ================================================================================
-- 第八部分: 其他业务表 (5张表)
-- ================================================================================

-- 8.1 API 密钥表 (api_key)
ALTER TABLE `api_key`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `user_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 8.2 应用连接器发布关联表 (app_connector_release_ref)
ALTER TABLE `app_connector_release_ref`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `record_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 8.3 应用动态对话草稿 (app_dynamic_conversation_draft)
ALTER TABLE `app_dynamic_conversation_draft`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `app_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 8.4 应用动态对话在线 (app_dynamic_conversation_online)
ALTER TABLE `app_dynamic_conversation_online`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `app_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 8.5 应用静态对话草稿 (app_static_conversation_draft)
ALTER TABLE `app_static_conversation_draft`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `template_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 8.6 应用静态对话在线 (app_static_conversation_online)
ALTER TABLE `app_static_conversation_online`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `template_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 8.7 应用对话模板在线 (app_conversation_template_online)
ALTER TABLE `app_conversation_template_online`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `app_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 8.8 数据复制任务表 (data_copy_task)
ALTER TABLE `data_copy_task`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `target_app_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 8.9 节点执行表 (node_execution)
ALTER TABLE `node_execution`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `execute_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 8.10 知识库文档审查表 (knowledge_document_review)
ALTER TABLE `knowledge_document_review`
    ADD COLUMN `tenant_id` VARCHAR(64) NULL COMMENT '租户ID' AFTER `knowledge_id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

SET FOREIGN_KEY_CHECKS = 1;

-- ================================================================================
-- 执行说明
-- ================================================================================
-- 1. 执行前备份: mysqldump -u root -p database_name > backup_before_tenant_id.sql
-- 2. 执行本脚本: mysql -u root -p database_name < 01_add_tenant_id_to_business_tables.sql
-- 3. 验证表结构: SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS
--                 WHERE TABLE_SCHEMA = 'database_name' AND COLUMN_NAME = 'tenant_id';
-- 4. 继续执行数据回填脚本 (02_backfill_tenant_id_data.sql)
-- ================================================================================

-- 迁移完成标记
-- 请记录执行完成时间: ____________________
-- 执行人: ____________________
-- 审核人: ____________________
