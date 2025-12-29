-- ============================================================================
-- ZKER 多租户迁移脚本：添加 tenant_id 字段到所有业务表
-- ============================================================================
-- 版本: v1.0.0
-- 创建日期: 2025-01-01
-- 执行时间: 预计 5-10 分钟（取决于数据量）
-- 影响范围: 所有核心业务表
--
-- 功能说明:
-- 1. 为所有业务表添加 tenant_id VARCHAR(64) 字段
-- 2. 创建 idx_tenant_id 索引以提升查询性能
-- 3. 设置默认值为 'default_tenant' 以兼容现有数据
-- 4. 包含数据验证脚本
--
-- 注意事项:
-- ⚠️ 执行前请先备份数据库！
-- ⚠️ 建议在低峰期执行！
-- ⚠️ 大表操作可能需要较长时间！
--
-- 回滚方案:
-- ALTER TABLE {table_name} DROP COLUMN tenant_id;
-- ALTER TABLE {table_name} DROP INDEX idx_tenant_id;
-- ============================================================================

-- ============================================================================
-- 第一部分：核心业务表迁移
-- ============================================================================

-- ----------------------------------------------------------------------------
-- 表1: bots - Bot管理表
-- ----------------------------------------------------------------------------
ALTER TABLE bots
    ADD COLUMN tenant_id VARCHAR(64) NOT NULL DEFAULT 'default_tenant'
        COMMENT '租户ID，用于多租户数据隔离'
        AFTER bot_id;

ALTER TABLE bots
    ADD INDEX idx_tenant_id (tenant_id);

ALTER TABLE bots
    ADD INDEX idx_tenant_status_created (tenant_id, status, created_at);

-- ----------------------------------------------------------------------------
-- 表2: conversations - 对话管理表
-- ----------------------------------------------------------------------------
ALTER TABLE conversations
    ADD COLUMN tenant_id VARCHAR(64) NOT NULL DEFAULT 'default_tenant'
        COMMENT '租户ID，用于多租户数据隔离'
        AFTER conversation_id;

ALTER TABLE conversations
    ADD INDEX idx_tenant_id (tenant_id);

ALTER TABLE conversations
    ADD INDEX idx_tenant_created (tenant_id, created_at);

-- ----------------------------------------------------------------------------
-- 表3: knowledge_bases - 知识库表
-- ----------------------------------------------------------------------------
ALTER TABLE knowledge_bases
    ADD COLUMN tenant_id VARCHAR(64) NOT NULL DEFAULT 'default_tenant'
        COMMENT '租户ID，用于多租户数据隔离'
        AFTER knowledge_id;

ALTER TABLE knowledge_bases
    ADD INDEX idx_tenant_id (tenant_id);

-- ----------------------------------------------------------------------------
-- 表4: workflows - 工作流表
-- ----------------------------------------------------------------------------
ALTER TABLE workflows
    ADD COLUMN tenant_id VARCHAR(64) NOT NULL DEFAULT 'default_tenant'
        COMMENT '租户ID，用于多租户数据隔离'
        AFTER workflow_id;

ALTER TABLE workflows
    ADD INDEX idx_tenant_id (tenant_id);

-- ----------------------------------------------------------------------------
-- 表5: agent_draft - Agent草稿表
-- ----------------------------------------------------------------------------
ALTER TABLE agent_draft
    ADD COLUMN tenant_id VARCHAR(64) NOT NULL DEFAULT 'default_tenant'
        COMMENT '租户ID，用于多租户数据隔离'
        AFTER draft_id;

ALTER TABLE agent_draft
    ADD INDEX idx_tenant_id (tenant_id);

-- ============================================================================
-- 第二部分：关联业务表迁移
-- ============================================================================

-- ----------------------------------------------------------------------------
-- 表6: bot_messages - 消息表
-- ----------------------------------------------------------------------------
ALTER TABLE bot_messages
    ADD COLUMN tenant_id VARCHAR(64) NOT NULL DEFAULT 'default_tenant'
        COMMENT '租户ID，用于多租户数据隔离'
        AFTER message_id;

ALTER TABLE bot_messages
    ADD INDEX idx_tenant_id (tenant_id);

ALTER TABLE bot_messages
    ADD INDEX idx_tenant_conversation (tenant_id, conversation_id);

-- ----------------------------------------------------------------------------
-- 表7: bot_tools - Bot工具表
-- ----------------------------------------------------------------------------
ALTER TABLE bot_tools
    ADD COLUMN tenant_id VARCHAR(64) NOT NULL DEFAULT 'default_tenant'
        COMMENT '租户ID，用于多租户数据隔离'
        AFTER tool_id;

ALTER TABLE bot_tools
    ADD INDEX idx_tenant_id (tenant_id);

-- ----------------------------------------------------------------------------
-- 表8: bot_components - Bot组件表
-- ----------------------------------------------------------------------------
ALTER TABLE bot_components
    ADD COLUMN tenant_id VARCHAR(64) NOT NULL DEFAULT 'default_tenant'
        COMMENT '租户ID，用于多租户数据隔离'
        AFTER component_id;

ALTER TABLE bot_components
    ADD INDEX idx_tenant_id (tenant_id);

-- ============================================================================
-- 第三部分：数据验证
-- ============================================================================

-- 验证1: 检查所有表是否成功添加了tenant_id字段
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    COLUMN_TYPE,
    IS_NULLABLE,
    COLUMN_DEFAULT,
    COLUMN_COMMENT
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'tenant_id'
  AND TABLE_NAME IN (
    'bots', 'conversations', 'knowledge_bases', 'workflows', 'agent_draft',
    'bot_messages', 'bot_tools', 'bot_components'
  )
ORDER BY TABLE_NAME;

-- 验证2: 检查所有表是否成功创建了索引
SELECT
    TABLE_NAME,
    INDEX_NAME,
    COLUMN_NAME,
    SEQ_IN_INDEX
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND INDEX_NAME LIKE 'idx_tenant%'
  AND TABLE_NAME IN (
    'bots', 'conversations', 'knowledge_bases', 'workflows', 'agent_draft',
    'bot_messages', 'bot_tools', 'bot_components'
  )
ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX;

-- 验证3: 检查数据完整性（确保没有NULL值）
SELECT
    'bots' AS table_name,
    COUNT(*) AS total_rows,
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) AS null_count
FROM bots
UNION ALL
SELECT
    'conversations',
    COUNT(*),
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END)
FROM conversations
UNION ALL
SELECT
    'knowledge_bases',
    COUNT(*),
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END)
FROM knowledge_bases
UNION ALL
SELECT
    'workflows',
    COUNT(*),
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END)
FROM workflows
UNION ALL
SELECT
    'agent_draft',
    COUNT(*),
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END)
FROM agent_draft;

-- ============================================================================
-- 第四部分：性能优化（可选，根据实际情况执行）
-- ============================================================================

-- 分析表以更新查询优化器统计信息
ANALYZE TABLE bots;
ANALYZE TABLE conversations;
ANALYZE TABLE knowledge_bases;
ANALYZE TABLE workflows;
ANALYZE TABLE agent_draft;
ANALYZE TABLE bot_messages;
ANALYZE TABLE bot_tools;
ANALYZE TABLE bot_components;

-- ============================================================================
-- 第五部分：回滚脚本（仅用于紧急回滚）
-- ============================================================================
-- ⚠️ 警告：回滚将永久删除tenant_id列及相关数据，请谨慎操作！

/*
-- 回滚 bots 表
ALTER TABLE bots DROP COLUMN tenant_id;
ALTER TABLE bots DROP INDEX idx_tenant_id;
ALTER TABLE bots DROP INDEX idx_tenant_status_created;

-- 回滚 conversations 表
ALTER TABLE conversations DROP COLUMN tenant_id;
ALTER TABLE conversations DROP INDEX idx_tenant_id;
ALTER TABLE conversations DROP INDEX idx_tenant_created;

-- 回滚 knowledge_bases 表
ALTER TABLE knowledge_bases DROP COLUMN tenant_id;
ALTER TABLE knowledge_bases DROP INDEX idx_tenant_id;

-- 回滚 workflows 表
ALTER TABLE workflows DROP COLUMN tenant_id;
ALTER TABLE workflows DROP INDEX idx_tenant_id;

-- 回滚 agent_draft 表
ALTER TABLE agent_draft DROP COLUMN tenant_id;
ALTER TABLE agent_draft DROP INDEX idx_tenant_id;

-- 回滚 bot_messages 表
ALTER TABLE bot_messages DROP COLUMN tenant_id;
ALTER TABLE bot_messages DROP INDEX idx_tenant_id;
ALTER TABLE bot_messages DROP INDEX idx_tenant_conversation;

-- 回滚 bot_tools 表
ALTER TABLE bot_tools DROP COLUMN tenant_id;
ALTER TABLE bot_tools DROP INDEX idx_tenant_id;

-- 回滚 bot_components 表
ALTER TABLE bot_components DROP COLUMN tenant_id;
ALTER TABLE bot_components DROP INDEX idx_tenant_id;
*/

-- ============================================================================
-- 执行完成
-- ============================================================================
-- 执行完成后，请运行以下命令验证迁移结果：
-- 1. 检查应用日志，确认没有报错
-- 2. 运行验证脚本（见第三部分）
-- 3. 测试基本功能是否正常
-- 4. 监控数据库性能指标
-- ============================================================================
