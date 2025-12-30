-- backend/domain/tenant/migration/ddl/002_alter_tenant_id_not_null.sql
-- 清理脚本1: 将tenant_id设为NOT NULL
-- 执行时机: 灰度发布完成，100%流量切换后

-- ========== 警告：此操作不可逆 ========== --
-- 执行前确保：
-- 1. 所有数据已迁移完成
-- 2. 100%流量已切换到新架构
-- 3. 数据一致性验证通过
-- 4. 已完整备份数据库

SET @OLD_SQL_MODE=@@SQL_MODE, @@SQL_MODE='TRADITIONAL';

-- 1. bots表
ALTER TABLE `bots`
    MODIFY COLUMN `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    ALGORITHM=INPLACE, LOCK=NONE;

-- 2. bot_configs表
ALTER TABLE `bot_configs`
    MODIFY COLUMN `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    ALGORITHM=INPLACE, LOCK=NONE;

-- 3. conversations表
ALTER TABLE `conversations`
    MODIFY COLUMN `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    ALGORITHM=INPLACE, LOCK=NONE;

-- 4. messages表
ALTER TABLE `messages`
    MODIFY COLUMN `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    ALGORITHM=INPLACE, LOCK=NONE;

-- 5. knowledge_bases表
ALTER TABLE `knowledge_bases`
    MODIFY COLUMN `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    ALGORITHM=INPLACE, LOCK=NONE;

-- 6. knowledge_chunks表
ALTER TABLE `knowledge_chunks`
    MODIFY COLUMN `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    ALGORITHM=INPLACE, LOCK=NONE;

-- 7. workflows表
ALTER TABLE `workflows`
    MODIFY COLUMN `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    ALGORITHM=INPLACE, LOCK=NONE;

-- 8. workflow_executions表
ALTER TABLE `workflow_executions`
    MODIFY COLUMN `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    ALGORITHM=INPLACE, LOCK=NONE;

-- 9. single_agent_draft表
ALTER TABLE `single_agent_draft`
    MODIFY COLUMN `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    ALGORITHM=INPLACE, LOCK=NONE;

-- 10. published_bots表
ALTER TABLE `published_bots`
    MODIFY COLUMN `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    ALGORITHM=INPLACE, LOCK=NONE;

SET @@SQL_MODE=@OLD_SQL_MODE;

-- 验证
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    IS_NULLABLE
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'tenant_id'
  AND TABLE_NAME IN (
    'bots', 'bot_configs', 'conversations', 'messages',
    'knowledge_bases', 'knowledge_chunks', 'workflows',
    'workflow_executions', 'single_agent_draft', 'published_bots'
  );

-- 预期结果: 所有表的 IS_NULLABLE 应为 NO
