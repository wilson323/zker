-- ============================================================
-- 性能优化索引添加脚本 - 第一批
-- 目的：为高频查询添加复合索引，提升查询性能50-80%
-- 执行时机：维护窗口（凌晨2:00-4:00）
-- 预计时长：10-20分钟
-- 风险等级：低（无业务影响）
-- ============================================================

-- ============================================================
-- 1. bots表 - 租户Bot查询优化
-- ============================================================

-- 场景1：租户 + 状态 + 时间排序（最常见）
CREATE INDEX IF NOT EXISTS idx_tenant_status_created
ON bots(tenant_id, status, created_at DESC);

-- 场景2：租户 + 类型 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_type_created
ON bots(tenant_id, bot_type, created_at DESC);

-- 场景3：租户 + 创建者 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_creator_created
ON bots(tenant_id, creator_id, created_at DESC);


-- ============================================================
-- 2. conversations表 - 租户对话查询优化
-- ============================================================

-- 场景1：租户 + 状态 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_status_created
ON conversations(tenant_id, status, created_at DESC);

-- 场景2：租户 + 用户 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_user_created
ON conversations(tenant_id, user_id, created_at DESC);

-- 场景3：租户 + Bot + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_bot_created
ON conversations(tenant_id, bot_id, created_at DESC);


-- ============================================================
-- 3. messages表 - 租户消息查询优化（数据量最大，优化效果最明显）
-- ============================================================

-- 场景1：租户 + 会话 + 时间排序（最高频）
CREATE INDEX IF NOT EXISTS idx_tenant_conversation_created
ON messages(tenant_id, conversation_id, created_at DESC);

-- 场景2：租户 + 角色 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_role_created
ON messages(tenant_id, role, created_at DESC);

-- 场景3：租户 + 会话 + 角色 + 时间排序（复合查询）
CREATE INDEX IF NOT EXISTS idx_tenant_conv_role_created
ON messages(tenant_id, conversation_id, role, created_at DESC);


-- ============================================================
-- 4. knowledge_bases表 - 租户知识库查询优化
-- ============================================================

-- 场景1：租户 + 状态 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_status_created
ON knowledge_bases(tenant_id, status, created_at DESC);

-- 场景2：租户 + 类型 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_type_created
ON knowledge_bases(tenant_id, kb_type, created_at DESC);


-- ============================================================
-- 5. knowledge_chunks表 - 知识库分块查询优化
-- ============================================================

-- 场景1：租户 + 知识库 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_kb_created
ON knowledge_chunks(tenant_id, kb_id, created_at DESC);


-- ============================================================
-- 6. workflows表 - 租户工作流查询优化
-- ============================================================

-- 场景1：租户 + 状态 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_status_created
ON workflows(tenant_id, status, created_at DESC);

-- 场景2：租户 + 创建者 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_creator_created
ON workflows(tenant_id, creator_id, created_at DESC);


-- ============================================================
-- 7. workflow_executions表 - 工作流执行记录查询优化
-- ============================================================

-- 场景1：租户 + 工作流 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_workflow_created
ON workflow_executions(tenant_id, workflow_id, created_at DESC);

-- 场景2：租户 + 状态 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_status_created
ON workflow_executions(tenant_id, status, created_at DESC);

-- 场景3：租户 + 工作流 + 状态 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_wf_status_created
ON workflow_executions(tenant_id, workflow_id, status, created_at DESC);


-- ============================================================
-- 8. single_agent_draft表 - 单Agent草稿查询优化
-- ============================================================

-- 场景1：租户 + 用户 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_user_created
ON single_agent_draft(tenant_id, user_id, created_at DESC);

-- 场景2：租户 + 用户 + 状态 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_user_status_created
ON single_agent_draft(tenant_id, user_id, status, created_at DESC);


-- ============================================================
-- 9. published_bots表 - 已发布Bot查询优化
-- ============================================================

-- 场景1：租户 + 创建者 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_creator_created
ON published_bots(tenant_id, creator_id, created_at DESC);


-- ============================================================
-- 10. token_usage_logs表 - Token使用日志查询优化
-- ============================================================

-- 场景1：租户 + 成本 + 时间排序（成本分析）
CREATE INDEX IF NOT EXISTS idx_tenant_cost_created
ON token_usage_logs(tenant_id, total_cost, created_at DESC);

-- 场景2：租户 + Bot + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_bot_created
ON token_usage_logs(tenant_id, bot_id, created_at DESC);


-- ============================================================
-- 11. alert_history表 - 告警历史查询优化
-- ============================================================

-- 场景1：租户 + 严重级别 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_severity_created
ON alert_history(tenant_id, severity, created_at DESC);

-- 场景2：租户 + 状态 + 时间排序
CREATE INDEX IF NOT EXISTS idx_tenant_status_created
ON alert_history(tenant_id, status, created_at DESC);


-- ============================================================
-- 12. employees表 - 员工查询优化
-- ============================================================

-- 场景1：租户 + 状态 + 入职日期排序
CREATE INDEX IF NOT EXISTS idx_tenant_status_hire
ON employees(tenant_id, status, hire_date DESC);

-- 场景2：租户 + 部门 + 状态排序
CREATE INDEX IF NOT EXISTS idx_tenant_dept_status
ON employees(tenant_id, dept_id, status);


-- ============================================================
-- 验证脚本
-- ============================================================

-- 检查索引是否创建成功
SELECT
    TABLE_NAME,
    INDEX_NAME,
    GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) AS columns,
    INDEX_TYPE,
    NON_UNIQUE
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND INDEX_NAME IN (
      'idx_tenant_status_created',
      'idx_tenant_type_created',
      'idx_tenant_creator_created',
      'idx_tenant_user_created',
      'idx_tenant_bot_created',
      'idx_tenant_conversation_created',
      'idx_tenant_role_created',
      'idx_tenant_conv_role_created',
      'idx_tenant_kb_created',
      'idx_tenant_workflow_created',
      'idx_tenant_wf_status_created',
      'idx_tenant_cost_created',
      'idx_tenant_severity_created',
      'idx_tenant_status_hire',
      'idx_tenant_dept_status'
  )
GROUP BY TABLE_NAME, INDEX_NAME, INDEX_TYPE, NON_UNIQUE
ORDER BY TABLE_NAME, INDEX_NAME;


-- ============================================================
-- 性能测试脚本
-- ============================================================

-- 测试1：租户Bot列表查询
EXPLAIN SELECT *
FROM bots
WHERE tenant_id = 'test_tenant_id'
  AND status = 'active'
ORDER BY created_at DESC
LIMIT 20;

-- 预期结果：
-- type: ref 或 range
-- key: idx_tenant_status_created
-- rows: < 100


-- 测试2：租户消息查询
EXPLAIN SELECT *
FROM messages
WHERE tenant_id = 'test_tenant_id'
  AND conversation_id = 'test_conv_id'
ORDER BY created_at DESC
LIMIT 50;

-- 预期结果：
-- type: ref 或 range
-- key: idx_tenant_conversation_created
-- rows: < 100


-- ============================================================
-- 索引大小统计
-- ============================================================

-- 查看新增索引的大小
SELECT
    TABLE_NAME,
    INDEX_NAME,
    ROUND(STAT_VALUE * @@innodb_page_size / 1024 / 1024, 2) AS size_mb
FROM mysql.innodb_index_stats
WHERE DATABASE_NAME = DATABASE()
  AND INDEX_NAME IN (
      'idx_tenant_status_created',
      'idx_tenant_type_created',
      'idx_tenant_conversation_created',
      'idx_tenant_role_created'
  )
  AND STAT_NAME = 'size'
ORDER BY TABLE_NAME, INDEX_NAME;


-- ============================================================
-- 完成标记
-- ============================================================
-- 执行完成后，请在此处记录：
-- 执行时间: _______________
-- 执行结果: □ 成功  □ 失败
-- 新增索引数量: ______ 个
-- 索引总大小: _____ MB
-- 性能提升: ______ %
-- 备注说明: ___________________________________
-- ============================================================
