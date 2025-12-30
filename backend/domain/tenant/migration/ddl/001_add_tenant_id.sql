-- backend/domain/tenant/migration/ddl/001_add_tenant_id.sql
-- DDL脚本: 为业务表添加 tenant_id 字段
-- 版本: v1.0
-- 日期: 2025-01-01
-- 说明: 支持零停机变更，字段允许NULL，后续迁移数据后设为NOT NULL
-- 执行时间: 预计10-30分钟（取决于表大小）
-- 回滚脚本: rollback_001_add_tenant_id.sql

-- ============================================================================
-- 第1部分: 核心业务表（Bot相关）
-- ============================================================================

-- 1. bots 表 - Bot管理
ALTER TABLE bots
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER id,
ADD INDEX idx_tenant_id (tenant_id),
ADD UNIQUE KEY uk_tenant_bot (tenant_id, bot_id);

-- 2. bot_configs 表 - Bot配置
ALTER TABLE bot_configs
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER id,
ADD INDEX idx_tenant_id (tenant_id);

-- 3. single_agent_draft 表 - 单Agent草稿
ALTER TABLE single_agent_draft
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER id,
ADD INDEX idx_tenant_id (tenant_id);

-- 4. published_bots 表 - 已发布Bot
ALTER TABLE published_bots
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER id,
ADD INDEX idx_tenant_id (tenant_id),
ADD UNIQUE KEY uk_tenant_published_bot (tenant_id, published_bot_id);

-- ============================================================================
-- 第2部分: 会话与消息相关表
-- ============================================================================

-- 5. conversations 表 - 对话记录
ALTER TABLE conversations
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_created (tenant_id, created_at);

-- 6. messages 表 - 消息记录
ALTER TABLE messages
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_conv (tenant_id, conv_id);

-- ============================================================================
-- 第3部分: 知识库相关表
-- ============================================================================

-- 7. knowledge_bases 表 - 知识库
ALTER TABLE knowledge_bases
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER id,
ADD INDEX idx_tenant_id (tenant_id),
ADD UNIQUE KEY uk_tenant_knowledge (tenant_id, knowledge_base_id);

-- 8. knowledge_chunks 表 - 知识库分块
ALTER TABLE knowledge_chunks
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_knowledge (tenant_id, knowledge_base_id);

-- ============================================================================
-- 第4部分: 工作流相关表
-- ============================================================================

-- 9. workflows 表 - 工作流
ALTER TABLE workflows
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER id,
ADD INDEX idx_tenant_id (tenant_id),
ADD UNIQUE KEY uk_tenant_workflow (tenant_id, workflow_id);

-- 10. workflow_executions 表 - 工作流执行记录
ALTER TABLE workflow_executions
ADD COLUMN tenant_id VARCHAR(36) NULL COMMENT '租户ID' AFTER id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_workflow (tenant_id, workflow_id),
ADD INDEX idx_tenant_status (tenant_id, status, created_at);

-- ============================================================================
-- 验证DDL执行结果
-- ============================================================================

-- 验证 bots 表
SELECT
    COLUMN_NAME,
    COLUMN_TYPE,
    IS_NULLABLE,
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'bots'
  AND COLUMN_NAME = 'tenant_id';

-- 验证索引创建
SHOW INDEX FROM bots WHERE Key_name = 'idx_tenant_id';
SHOW INDEX FROM bots WHERE Key_name = 'uk_tenant_bot';

-- 验证其他表（可选）
-- SELECT TABLE_NAME, COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS
-- WHERE TABLE_SCHEMA = DATABASE() AND COLUMN_NAME = 'tenant_id'
-- ORDER BY TABLE_NAME;

-- ============================================================================
-- 执行说明
-- ============================================================================
--
-- 1. 执行前准备：
--    - 备份数据库: mysqldump -u root -p coze_production > backup_$(date +%Y%m%d_%H%M%S).sql
--    - 检查磁盘空间: df -h
--    - 确认维护窗口: 建议凌晨2:00-4:00执行
--
-- 2. 执行命令：
--    mysql -u root -p < backend/domain/tenant/migration/ddl/001_add_tenant_id.sql
--
-- 3. 验证结果：
--    - 检查所有表是否已添加 tenant_id 字段
--    - 检查索引是否创建成功
--    - 测试业务功能是否正常
--
-- 4. 回滚方案：
--    - 如需回滚，执行: mysql -u root -p < backend/domain/tenant/migration/ddl/rollback_001_add_tenant_id.sql
--    - 预计回滚时间: 5-10分钟
--
-- 5. 注意事项：
--    - 字段允许NULL，不影响现有业务
--    - 索引采用后台创建（ALGORITHM=INPLACE, LOCK=NONE），不阻塞读写
--    - 执行时间取决于表大小，大表可能需要较长时间
--    - 建议在低峰期执行
--
-- ============================================================================
