-- backend/domain/tenant/migration/ddl/rollback_001_add_tenant_id.sql
-- DDL回滚脚本: 删除业务表的 tenant_id 字段
-- 版本: v1.0
-- 日期: 2025-01-01
-- 说明: 快速回滚到迁移前状态
-- 执行时间: 预计5-10分钟
-- 前置脚本: 001_add_tenant_id.sql
--
-- ⚠️  警告: 执行此脚本将删除 tenant_id 字段和相关索引
--      请确保已备份重要数据

-- ============================================================================
-- 第1部分: 删除索引（必须先删除索引，才能删除字段）
-- ============================================================================

-- 1. bots 表
ALTER TABLE bots DROP INDEX uk_tenant_bot;
ALTER TABLE bots DROP INDEX idx_tenant_id;

-- 2. bot_configs 表
ALTER TABLE bot_configs DROP INDEX idx_tenant_id;

-- 3. single_agent_draft 表
ALTER TABLE single_agent_draft DROP INDEX idx_tenant_id;

-- 4. published_bots 表
ALTER TABLE published_bots DROP INDEX uk_tenant_published_bot;
ALTER TABLE published_bots DROP INDEX idx_tenant_id;

-- 5. conversations 表
ALTER TABLE conversations DROP INDEX idx_tenant_created;
ALTER TABLE conversations DROP INDEX idx_tenant_id;

-- 6. messages 表
ALTER TABLE messages DROP INDEX idx_tenant_conv;
ALTER TABLE messages DROP INDEX idx_tenant_id;

-- 7. knowledge_bases 表
ALTER TABLE knowledge_bases DROP INDEX uk_tenant_knowledge;
ALTER TABLE knowledge_bases DROP INDEX idx_tenant_id;

-- 8. knowledge_chunks 表
ALTER TABLE knowledge_chunks DROP INDEX idx_tenant_knowledge;
ALTER TABLE knowledge_chunks DROP INDEX idx_tenant_id;

-- 9. workflows 表
ALTER TABLE workflows DROP INDEX uk_tenant_workflow;
ALTER TABLE workflows DROP INDEX idx_tenant_id;

-- 10. workflow_executions 表
ALTER TABLE workflow_executions DROP INDEX idx_tenant_status;
ALTER TABLE workflow_executions DROP INDEX idx_tenant_workflow;
ALTER TABLE workflow_executions DROP INDEX idx_tenant_id;

-- ============================================================================
-- 第2部分: 删除 tenant_id 字段
-- ============================================================================

-- 1. bots 表
ALTER TABLE bots DROP COLUMN tenant_id;

-- 2. bot_configs 表
ALTER TABLE bot_configs DROP COLUMN tenant_id;

-- 3. single_agent_draft 表
ALTER TABLE single_agent_draft DROP COLUMN tenant_id;

-- 4. published_bots 表
ALTER TABLE published_bots DROP COLUMN tenant_id;

-- 5. conversations 表
ALTER TABLE conversations DROP COLUMN tenant_id;

-- 6. messages 表
ALTER TABLE messages DROP COLUMN tenant_id;

-- 7. knowledge_bases 表
ALTER TABLE knowledge_bases DROP COLUMN tenant_id;

-- 8. knowledge_chunks 表
ALTER TABLE knowledge_chunks DROP COLUMN tenant_id;

-- 9. workflows 表
ALTER TABLE workflows DROP COLUMN tenant_id;

-- 10. workflow_executions 表
ALTER TABLE workflow_executions DROP COLUMN tenant_id;

-- ============================================================================
-- 验证回滚结果
-- ============================================================================

-- 验证字段已删除（应返回空结果）
SELECT TABLE_NAME, COLUMN_NAME
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'tenant_id'
  AND TABLE_NAME IN (
    'bots', 'bot_configs', 'single_agent_draft', 'published_bots',
    'conversations', 'messages',
    'knowledge_bases', 'knowledge_chunks',
    'workflows', 'workflow_executions'
  );

-- 如果上述查询返回空结果，说明回滚成功

-- ============================================================================
-- 执行说明
-- ============================================================================
--
-- 1. 回滚场景：
--    - DDL执行后发现严重问题
--    - 数据迁移无法继续
--    - 需要快速恢复到原始状态
--
-- 2. 执行命令：
--    mysql -u root -p < backend/domain/tenant/migration/ddl/rollback_001_add_tenant_id.sql
--
-- 3. 回滚影响：
--    - 删除所有 tenant_id 字段
--    - 删除相关索引
--    - 数据无影响（因为字段允许NULL）
--
-- 4. 回滚后操作：
--    - 验证业务功能正常
--    - 检查应用日志
--    - 通知相关团队
--
-- 5. 预计执行时间: 5-10分钟
--
-- ============================================================================
