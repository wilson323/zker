-- ============================================
-- tenant_id 迁移 DDL 变更
-- 执行时机: 维护窗口（建议凌晨2:00-4:00）
-- 执行时长: 预计10-20分钟
-- 最后更新: 2025-12-30
-- ============================================
--
-- ⚠️ 重要提示:
-- 1. 所有字段先设为 NULL，迁移完成后再改为 NOT NULL
-- 2. 外键约束使用 ON DELETE CASCADE，级联删除
-- 3. 索引命名规范: idx_tenant_id, idx_tenant_created
-- 4. 每张表的 ALTER 语句独立执行，不要在一个事务中
-- 5. 使用 LOW_PRIORITY 降低锁影响（MySQL 8.0.24+）
--
-- ============================================

-- ============================================
-- 表1: bots - Bot表
-- ============================================
ALTER TABLE bots
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER bot_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_created (tenant_id, created_at),
ADD CONSTRAINT fk_bots_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(tenant_id)
    ON DELETE CASCADE;

-- 验证
SHOW COLUMNS FROM bots LIKE 'tenant_id';
SHOW INDEX FROM bots WHERE Key_name = 'idx_tenant_id';


-- ============================================
-- 表2: bot_configs - Bot配置表
-- ============================================
ALTER TABLE bot_configs
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER config_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD CONSTRAINT fk_bot_configs_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(tenant_id)
    ON DELETE CASCADE;

-- 验证
SHOW COLUMNS FROM bot_configs LIKE 'tenant_id';


-- ============================================
-- 表3: conversations - 对话表
-- ============================================
ALTER TABLE conversations
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER conversation_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_created (tenant_id, created_at),
ADD INDEX idx_tenant_status (tenant_id, status),
ADD CONSTRAINT fk_conversations_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(tenant_id)
    ON DELETE CASCADE;

-- 验证
SHOW COLUMNS FROM conversations LIKE 'tenant_id';


-- ============================================
-- 表4: messages - 消息表（数据量最大，~1000万条）
-- ============================================
ALTER TABLE messages
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER message_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_conversation (tenant_id, conversation_id),
ADD INDEX idx_tenant_created (tenant_id, created_at),
ADD CONSTRAINT fk_messages_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(tenant_id)
    ON DELETE CASCADE;

-- 验证
SHOW COLUMNS FROM messages LIKE 'tenant_id';

-- ⚠️ 注意: messages表数据量大，建议单独执行，使用LOW_PRIORITY
-- ALTER TABLE messages ... LOCK=NONE;


-- ============================================
-- 表5: knowledge_bases - 知识库表
-- ============================================
ALTER TABLE knowledge_bases
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER kb_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_created (tenant_id, created_at),
ADD CONSTRAINT fk_knowledge_bases_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(tenant_id)
    ON DELETE CASCADE;

-- 验证
SHOW COLUMNS FROM knowledge_bases LIKE 'tenant_id';


-- ============================================
-- 表6: knowledge_chunks - 知识库分块表
-- ============================================
ALTER TABLE knowledge_chunks
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER chunk_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_kb (tenant_id, kb_id),
ADD CONSTRAINT fk_knowledge_chunks_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(tenant_id)
    ON DELETE CASCADE;

-- 验证
SHOW COLUMNS FROM knowledge_chunks LIKE 'tenant_id';


-- ============================================
-- 表7: workflows - 工作流表
-- ============================================
ALTER TABLE workflows
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER workflow_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_created (tenant_id, created_at),
ADD INDEX idx_tenant_status (tenant_id, status),
ADD CONSTRAINT fk_workflows_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(tenant_id)
    ON DELETE CASCADE;

-- 验证
SHOW COLUMNS FROM workflows LIKE 'tenant_id';


-- ============================================
-- 表8: workflow_executions - 工作流执行记录表
-- ============================================
ALTER TABLE workflow_executions
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER execution_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_workflow (tenant_id, workflow_id),
ADD INDEX idx_tenant_created (tenant_id, created_at),
ADD CONSTRAINT fk_workflow_executions_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(tenant_id)
    ON DELETE CASCADE;

-- 验证
SHOW COLUMNS FROM workflow_executions LIKE 'tenant_id';


-- ============================================
-- 表9: single_agent_draft - 单Agent草稿表
-- ============================================
ALTER TABLE single_agent_draft
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER draft_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_user (tenant_id, user_id),
ADD CONSTRAINT fk_single_agent_draft_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(tenant_id)
    ON DELETE CASCADE;

-- 验证
SHOW COLUMNS FROM single_agent_draft LIKE 'tenant_id';


-- ============================================
-- 表10: published_bots - 已发布Bot表
-- ============================================
ALTER TABLE published_bots
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER published_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_created (tenant_id, created_at),
ADD CONSTRAINT fk_published_bots_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(tenant_id)
    ON DELETE CASCADE;

-- 验证
SHOW COLUMNS FROM published_bots LIKE 'tenant_id';


-- ============================================
-- 执行后验证脚本
-- ============================================

-- 1. 检查所有表是否都有tenant_id字段
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    COLUMN_TYPE,
    IS_NULLABLE,
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'tenant_id'
  AND TABLE_NAME IN (
      'bots', 'bot_configs', 'conversations', 'messages',
      'knowledge_bases', 'knowledge_chunks', 'workflows',
      'workflow_executions', 'single_agent_draft', 'published_bots'
  )
ORDER BY TABLE_NAME;

-- 2. 检查所有索引是否创建
SELECT
    TABLE_NAME,
    INDEX_NAME,
    COLUMN_NAME,
    SEQ_IN_INDEX
FROM INFORMATION_SCHEMA.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND INDEX_NAME LIKE 'idx_tenant%'
  AND TABLE_NAME IN (
      'bots', 'bot_configs', 'conversations', 'messages',
      'knowledge_bases', 'knowledge_chunks', 'workflows',
      'workflow_executions', 'single_agent_draft', 'published_bots'
  )
ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX;

-- 3. 检查外键约束是否创建
SELECT
    CONSTRAINT_NAME,
    TABLE_NAME,
    COLUMN_NAME,
    REFERENCED_TABLE_NAME,
    REFERENCED_COLUMN_NAME
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
WHERE TABLE_SCHEMA = DATABASE()
  AND REFERENCED_TABLE_NAME = 'tenants'
  AND REFERENCED_COLUMN_NAME = 'tenant_id'
  AND TABLE_NAME IN (
      'bots', 'bot_configs', 'conversations', 'messages',
      'knowledge_bases', 'knowledge_chunks', 'workflows',
      'workflow_executions', 'single_agent_draft', 'published_bots'
  )
ORDER BY TABLE_NAME;


-- ============================================
-- 回滚脚本（谨慎使用！）
-- ============================================

/*
-- ⚠️ 仅在迁移失败时执行回滚

-- 回滚 bots 表
ALTER TABLE bots DROP COLUMN tenant_id;

-- 回滚 bot_configs 表
ALTER TABLE bot_configs DROP COLUMN tenant_id;

-- 回滚 conversations 表
ALTER TABLE conversations DROP COLUMN tenant_id;

-- 回滚 messages 表
ALTER TABLE messages DROP COLUMN tenant_id;

-- 回滚 knowledge_bases 表
ALTER TABLE knowledge_bases DROP COLUMN tenant_id;

-- 回滚 knowledge_chunks 表
ALTER TABLE knowledge_chunks DROP COLUMN tenant_id;

-- 回滚 workflows 表
ALTER TABLE workflows DROP COLUMN tenant_id;

-- 回滚 workflow_executions 表
ALTER TABLE workflow_executions DROP COLUMN tenant_id;

-- 回滚 single_agent_draft 表
ALTER TABLE single_agent_draft DROP COLUMN tenant_id;

-- 回滚 published_bots 表
ALTER TABLE published_bots DROP COLUMN tenant_id;
*/


-- ============================================
-- 执行完成标记
-- ============================================
-- 执行完成后，请在此处记录执行时间和结果
-- 执行时间: _______________
-- 执行结果: □ 成功  □ 失败
-- 备注说明: ___________________________________
-- ============================================
