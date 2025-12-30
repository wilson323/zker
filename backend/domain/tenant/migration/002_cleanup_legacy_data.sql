-- ============================================
-- 清理旧数据SQL脚本
-- ⚠️ 警告: 此脚本将删除旧字段，不可恢复！
-- 执行时机: 切换完成并观察1周后
-- 最后更新: 2025-12-30
-- ============================================
--
-- ⚠️ 重要提示:
-- 1. 建议在切换完成并观察1周后执行
-- 2. 执行前必须完整备份数据库
-- 3. 逐个表执行，不要批量执行
-- 4. 每个表执行后验证业务正常
--
-- ============================================

-- ============================================
-- 执行前检查清单
-- ============================================
-- [ ] 备份数据库（完整备份）
-- [ ] 确认切换已完成1周以上
-- [ ] 确认业务运行正常
-- [ ] 确认没有查询使用旧字段
-- [ ] 在测试环境验证
-- [ ] 选择业务低峰期执行（凌晨2:00-4:00）
-- ============================================


-- ============================================
-- ⚠️ 危险操作：删除旧字段
-- ============================================
--
-- 注意：以下为示例SQL，实际旧字段名称需要根据项目情况调整
-- 如果项目中没有这些旧字段，请勿执行
--
-- ============================================

-- 示例1: 删除 bots 表的旧字段（如果存在）
-- ALTER TABLE bots DROP COLUMN space_id;

-- 示例2: 删除 conversations 表的旧字段（如果存在）
-- ALTER TABLE conversations DROP COLUMN space_id;

-- 示例3: 删除 messages 表的旧字段（如果存在）
-- ALTER TABLE messages DROP COLUMN space_id;

-- 示例4: 删除 knowledge_bases 表的旧字段（如果存在）
-- ALTER TABLE knowledge_bases DROP COLUMN space_id;

-- 示例5: 删除 workflows 表的旧字段（如果存在）
-- ALTER TABLE workflows DROP COLUMN space_id;

-- ... 其他表类似处理


-- ============================================
-- 删除旧索引（如果存在）
-- ============================================

-- 示例1: 删除 bots 表的旧索引（如果存在）
-- ALTER TABLE bots DROP INDEX idx_space_id;

-- 示例2: 删除 conversations 表的旧索引（如果存在）
-- ALTER TABLE conversations DROP INDEX idx_space_created;


-- ============================================
-- 执行后验证脚本
-- ============================================

-- 1. 检查所有表是否都只有tenant_id字段
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    DATA_TYPE,
    IS_NULLABLE
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME IN (
      'bots', 'bot_configs', 'conversations', 'messages',
      'knowledge_bases', 'knowledge_chunks', 'workflows',
      'workflow_executions', 'single_agent_draft', 'published_bots'
  )
  AND COLUMN_NAME LIKE '%space%'
ORDER BY TABLE_NAME, COLUMN_NAME;

-- 预期结果: 空（说明旧字段已删除）


-- 2. 检查tenant_id字段是否为NOT NULL
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    IS_NULLABLE,
    COLUMN_TYPE
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'tenant_id'
  AND TABLE_NAME IN (
      'bots', 'bot_configs', 'conversations', 'messages',
      'knowledge_bases', 'knowledge_chunks', 'workflows',
      'workflow_executions', 'single_agent_draft', 'published_bots'
  )
ORDER BY TABLE_NAME;

-- 预期结果: IS_NULLABLE = NO（所有表的tenant_id都为NOT NULL）


-- 3. 验证业务查询是否正常
-- （需要根据实际业务查询调整）

-- 示例1: 查询bot数量
-- SELECT COUNT(*) FROM bots WHERE tenant_id = 'tenant_xxx';

-- 示例2: 查询对话列表
-- SELECT * FROM conversations WHERE tenant_id = 'tenant_xxx' LIMIT 10;


-- ============================================
-- 回滚脚本（仅紧急情况使用！）
-- ============================================
--
-- ⚠️ 警告: 如果删除了旧字段，无法直接回滚！
-- 需要从备份恢复数据库，然后重新添加旧字段
--
/*
-- 如果误删旧字段，需要从备份恢复
-- 步骤:
-- 1. 停止应用
-- 2. 从备份恢复数据库
-- 3. 重新添加旧字段
-- 4. 重启应用
*/


-- ============================================
-- 执行完成标记
-- ============================================
-- 执行完成后，请在此处记录执行时间和结果
-- 执行时间: _______________
-- 执行结果: □ 成功  □ 失败
-- 业务验证: □ 正常  □ 异常
-- 备注说明: ___________________________________
-- ============================================


-- ============================================
-- 清理完成清单
-- ============================================
-- [ ] 旧字段已删除
-- [ ] 旧索引已删除
-- [ ] tenant_id字段为NOT NULL
-- [ ] 业务查询正常
-- [ ] 应用运行正常
-- [ ] 性能监控正常
-- [ ] 日志无异常
-- ============================================
