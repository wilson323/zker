-- ============================================
-- ZKER 性能优化 - 数据库索引优化脚本
-- 版本: v1.0.0
-- 生成时间: 2025-01-03
-- 预期效果: 查询性能提升 5-10倍
-- ============================================

-- ============================================
-- 第一部分: 知识库模块索引优化 (P0-严重)
-- ============================================

-- 1.1 knowledge 表索引 (3个)
-- ----------------------------

-- 索引1.1: 租户知识库列表查询 (最频繁的查询)
-- 查询: SELECT * FROM knowledge WHERE app_id = ? AND space_id = ? AND status = ? ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_knowledge_app_space_status_created
ON knowledge(app_id, space_id, status, created_at DESC)
COMMENT '租户知识库列表查询 - 支持app_id+space_id+status过滤和created_at排序';

-- 索引1.2: 知识库状态查询
-- 查询: SELECT * FROM knowledge WHERE status = ? ORDER BY updated_at DESC
CREATE INDEX IF NOT EXISTS idx_knowledge_status_updated
ON knowledge(status, updated_at DESC)
COMMENT '知识库状态查询 - 支持status过滤和updated_at排序';

-- 索引1.3: 知识库唯一性约束 (防止同名知识库)
CREATE UNIQUE INDEX IF NOT EXISTS uk_knowledge_space_name
ON knowledge(space_id, name)
WHERE deleted_at IS NULL
COMMENT '知识库名称唯一性约束 - 防止同一空间下重名';

-- 1.2 knowledge_document 表索引 (4个)
-- -------------------------------------

-- 索引2.1: 文档列表查询 (核心查询)
-- 查询: SELECT * FROM knowledge_document WHERE knowledge_id = ? AND status = ? ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_document_knowledge_status_created
ON knowledge_document(knowledge_id, status, created_at DESC)
COMMENT '文档列表查询 - 支持knowledge_id+status过滤和created_at排序';

-- 索引2.2: 文档覆盖索引 (避免回表)
-- 查询: SELECT id, name, status, created_at, size FROM knowledge_document WHERE knowledge_id = ? AND status = ?
CREATE INDEX IF NOT EXISTS idx_document_cover
ON knowledge_document(knowledge_id, status, created_at, id, name, size)
COMMENT '文档列表覆盖索引 - 包含常用字段避免回表查询';

-- 索引2.3: 文档类型查询
-- 查询: SELECT * FROM knowledge_document WHERE document_type = ? AND status = ?
CREATE INDEX IF NOT EXISTS idx_document_type_status
ON knowledge_document(document_type, status)
COMMENT '文档类型查询 - 支持document_type+status过滤';

-- 索引2.4: 文档去重约束
CREATE UNIQUE INDEX IF NOT EXISTS uk_document_knowledge_name
ON knowledge_document(knowledge_id, name)
WHERE deleted_at IS NULL
COMMENT '文档名称唯一性约束 - 防止同一知识库下重名';

-- 1.3 knowledge_document_slice 表索引 (3个)
-- ------------------------------------------

-- 索引3.1: Slice列表查询
-- 查询: SELECT * FROM knowledge_document_slice WHERE document_id = ? ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_slice_document_created
ON knowledge_document_slice(document_id, created_at DESC)
COMMENT 'Slice列表查询 - 支持document_id过滤和created_at排序';

-- 索引3.2: Slice知识库查询
-- 查询: SELECT * FROM knowledge_document_slice WHERE knowledge_id = ? ORDER BY sequence
CREATE INDEX IF NOT EXISTS idx_slice_knowledge_sequence
ON knowledge_document_slice(knowledge_id, sequence)
COMMENT 'Slice知识库查询 - 支持knowledge_id过滤和sequence排序';

-- 索引3.3: Slice统计查询
-- 查询: SELECT COUNT(*) FROM knowledge_document_slice WHERE knowledge_id = ? AND hit = ?
CREATE INDEX IF NOT EXISTS idx_slice_knowledge_hit
ON knowledge_document_slice(knowledge_id, hit)
COMMENT 'Slice统计查询 - 支持knowledge_id+hit统计';

-- ============================================
-- 第二部分: 对话模块索引优化 (P0-严重)
-- ============================================

-- 2.1 message 表索引 (5个)
-- -------------------------

-- 索引4.1: 对话消息列表 (最频繁查询)
-- 查询: SELECT * FROM message WHERE conversation_id = ? AND status = ? ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_message_conversation_status_created
ON message(conversation_id, status, created_at DESC)
COMMENT '对话消息列表 - 支持conversation_id+status过滤和created_at排序';

-- 索引4.2: Run消息查询
-- 查询: SELECT * FROM message WHERE run_id IN (?) AND status = ? ORDER BY created_at ASC
CREATE INDEX IF NOT EXISTS idx_message_run_status_created
ON message(run_id, status, created_at ASC)
COMMENT 'Run消息查询 - 支持run_id+status过滤和created_at升序';

-- 索引4.3: 游标分页优化 (解决OFFSET深度分页问题)
-- 查询: SELECT * FROM message WHERE conversation_id = ? AND created_at < ? ORDER BY created_at DESC, id DESC
CREATE INDEX IF NOT EXISTS idx_message_conversation_created_cursor
ON message(conversation_id, created_at, id)
COMMENT '消息游标分页 - 支持游标分页避免深度OFFSET';

-- 索引4.4: 消息类型查询
-- 查询: SELECT * FROM message WHERE message_type = ? ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_message_type_created
ON message(message_type, created_at DESC)
COMMENT '消息类型查询 - 支持message_type过滤和created_at排序';

-- 索引4.5: Agent消息查询
-- 查询: SELECT * FROM message WHERE agent_id = ? ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_message_agent_created
ON message(agent_id, created_at DESC)
COMMENT 'Agent消息查询 - 支持agent_id过滤和created_at排序';

-- 2.2 conversation 表索引 (3个)
-- -------------------------------

-- 索引5.1: 租户对话列表 (核心查询)
-- 查询: SELECT * FROM conversation WHERE tenant_id = ? AND space_id = ? ORDER BY updated_at DESC
CREATE INDEX IF NOT EXISTS idx_conversation_tenant_space_updated
ON conversation(tenant_id, space_id, updated_at DESC)
COMMENT '租户对话列表 - 支持tenant_id+space_id过滤和updated_at排序';

-- 索引5.2: 对话状态查询
-- 查询: SELECT * FROM conversation WHERE status = ? ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_conversation_status_created
ON conversation(status, created_at DESC)
COMMENT '对话状态查询 - 支持status过滤和created_at排序';

-- 索引5.3: Agent对话查询
-- 查询: SELECT * FROM conversation WHERE agent_id = ? ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_conversation_agent_created
ON conversation(agent_id, created_at DESC)
COMMENT 'Agent对话查询 - 支持agent_id过滤和created_at排序';

-- ============================================
-- 第三部分: 插件模块索引优化 (P0-严重)
-- ============================================

-- 3.1 plugin 表索引 (3个)
-- ------------------------

-- 索引6.1: 空间插件列表 (核心查询)
-- 查询: SELECT * FROM plugin WHERE space_id = ? ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_plugin_space_created
ON plugin(space_id, created_at DESC)
COMMENT '空间插件列表 - 支持space_id过滤和created_at排序';

-- 索引6.2: 插件类型查询
-- 查询: SELECT * FROM plugin WHERE plugin_type = ? AND status = ?
CREATE INDEX IF NOT EXISTS idx_plugin_type_status
ON plugin(plugin_type, status)
COMMENT '插件类型查询 - 支持plugin_type+status过滤';

-- 索引6.3: 插件唯一性约束
CREATE UNIQUE INDEX IF NOT EXISTS uk_plugin_space_developer_name
ON plugin(space_id, developer_id, name)
WHERE deleted_at IS NULL
COMMENT '插件名称唯一性约束 - 防止同一空间+开发者下重名';

-- ============================================
-- 第四部分: 工作流模块索引优化 (P1-高优先级)
-- ============================================

-- 4.1 workflow 相关表索引 (4个)
-- -------------------------------

-- 索引7.1: 工作流列表
-- 查询: SELECT * FROM workflow WHERE space_id = ? AND status = ? ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_workflow_space_status_created
ON workflow(space_id, status, created_at DESC)
COMMENT '工作流列表 - 支持space_id+status过滤和created_at排序';

-- 索引7.2: 工作流执行历史 (核心查询)
-- 查询: SELECT * FROM workflow_execution WHERE workflow_id = ? AND status = ? ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_workflow_execution_workflow_status_created
ON workflow_execution(workflow_id, status, created_at DESC)
COMMENT '工作流执行历史 - 支持workflow_id+status过滤和created_at排序';

-- 索引7.3: 节点执行查询
-- 查询: SELECT * FROM node_execution WHERE run_id = ? ORDER BY created_at ASC
CREATE INDEX IF NOT EXISTS idx_node_execution_run_created
ON node_execution(run_id, created_at ASC)
COMMENT '节点执行查询 - 支持run_id过滤和created_at升序';

-- 索引7.4: 工作流版本查询
-- 查询: SELECT * FROM workflow WHERE workflow_id = ? ORDER BY version DESC
CREATE INDEX IF NOT EXISTS idx_workflow_workflow_id_version
ON workflow(workflow_id, version DESC)
COMMENT '工作流版本查询 - 支持workflow_id过滤和version降序';

-- ============================================
-- 第五部分: 权限模块索引优化 (P0-严重)
-- ============================================

-- 5.1 permission 相关表索引 (3个)
-- -------------------------------

-- 索引8.1: 角色权限查询 (核心查询,高频)
-- 查询: SELECT * FROM data_permission WHERE role_id = ? AND resource_type = ? AND resource_id = ?
CREATE INDEX IF NOT EXISTS idx_data_permission_role_resource
ON data_permission(role_id, resource_type, resource_id)
COMMENT '角色权限查询 - 支持role_id+resource_type+resource_id联合查询';

-- 索引8.2: 用户角色查询
-- 查询: SELECT * FROM user_roles WHERE user_id = ? AND role_id = ?
CREATE INDEX IF NOT EXISTS idx_user_role_user_role
ON user_roles(user_id, role_id)
COMMENT '用户角色查询 - 支持user_id+role_id联合查询';

-- 索引8.3: 字段权限查询
-- 查询: SELECT * FROM field_permission WHERE role_id = ? AND resource_type = ? AND field_name = ?
CREATE INDEX IF NOT EXISTS idx_field_permission_role_resource
ON field_permission(role_id, resource_type, field_name)
COMMENT '字段权限查询 - 支持role_id+resource_type+field_name联合查询';

-- ============================================
-- 第六部分: 租户模块索引优化 (P0-严重)
-- ============================================

-- 6.1 tenant 相关表索引 (2个)
-- ----------------------------

-- 索引9.1: 租户用户列表 (核心查询)
-- 查询: SELECT * FROM users WHERE tenant_id = ? AND status = ? ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_user_tenant_status_created
ON users(tenant_id, status, created_at DESC)
COMMENT '租户用户列表 - 支持tenant_id+status过滤和created_at排序';

-- 索引9.2: 配额查询
-- 查询: SELECT * FROM quota WHERE tenant_id = ? AND quota_type = ?
CREATE INDEX IF NOT EXISTS idx_quota_tenant_type
ON quota(tenant_id, quota_type)
COMMENT '配额查询 - 支持tenant_id+quota_type联合查询';

-- ============================================
-- 第七部分: 其他业务表索引优化 (P1-高优先级)
-- ============================================

-- 7.1 其他业务表索引 (5个)
-- -------------------------

-- 索引10.1: API密钥查询
-- 查询: SELECT * FROM api_key WHERE space_id = ? AND status = ?
CREATE INDEX IF NOT EXISTS idx_api_key_space_status
ON api_key(space_id, status)
COMMENT 'API密钥查询 - 支持space_id+status过滤';

-- 索引10.2: 数据库配置查询
-- 查询: SELECT * FROM database_info WHERE space_id = ? AND agent_id = ?
CREATE INDEX IF NOT EXISTS idx_database_space_agent
ON database_info(space_id, agent_id)
COMMENT '数据库配置查询 - 支持space_id+agent_id联合查询';

-- 索引10.3: 变量实例查询
-- 查询: SELECT * FROM variable_instance WHERE space_id = ? AND variable_key = ?
CREATE INDEX IF NOT EXISTS idx_variable_instance_space_key
ON variable_instance(space_id, variable_key)
COMMENT '变量实例查询 - 支持space_id+variable_key联合查询';

-- 索引10.4: 模板查询
-- 查询: SELECT * FROM template WHERE space_id = ? AND template_type = ?
CREATE INDEX IF NOT EXISTS idx_template_space_type
ON template(space_id, template_type)
COMMENT '模板查询 - 支持space_id+template_type联合查询';

-- 索引10.5: 快捷命令查询
-- 查询: SELECT * FROM shortcut_command WHERE space_id = ? AND agent_id = ?
CREATE INDEX IF NOT EXISTS idx_shortcut_space_agent
ON shortcut_command(space_id, agent_id)
COMMENT '快捷命令查询 - 支持space_id+agent_id联合查询';

-- ============================================
-- 第八部分: 索引维护与验证
-- ============================================

-- 8.1 索引使用情况分析
-- ---------------------
-- 查询所有索引的详细信息
SELECT
    TABLE_NAME,
    INDEX_NAME,
    SEQ_IN_INDEX,
    COLUMN_NAME,
    CARDINALITY,
    INDEX_TYPE
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME IN (
    'knowledge',
    'knowledge_document',
    'knowledge_document_slice',
    'message',
    'conversation',
    'plugin',
    'workflow',
    'workflow_execution',
    'node_execution',
    'data_permission',
    'user_roles',
    'field_permission',
    'users',
    'quota',
    'api_key',
    'database_info',
    'variable_instance',
    'template',
    'shortcut_command'
  )
ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX;

-- 8.2 查看未使用的索引 (运行一段时间后执行)
-- ------------------------------------------
SELECT
    object_schema AS table_schema,
    object_name AS table_name,
    index_name,
    count_star AS usage_count,
    sum_timer_wait / 1000000000000 AS total_latency_sec
FROM performance_schema.table_io_waits_summary_by_index_usage
WHERE index_name IS NOT NULL
  AND count_star = 0
  AND object_schema = DATABASE()
ORDER BY object_schema, object_name;

-- 8.3 索引大小统计
-- -----------------
SELECT
    TABLE_NAME,
    INDEX_NAME,
    ROUND(STAT_VALUE * @@innodb_page_size / 1024 / 1024, 2) AS size_mb
FROM mysql.innodb_index_stats
WHERE database_name = DATABASE()
  AND stat_name = 'size'
  AND stat_description != 'Number of pages in the index'
ORDER BY size_mb DESC
LIMIT 20;

-- 8.4 索引效率分析
-- -----------------
-- 分析索引的选择性(基数)
SELECT
    TABLE_NAME,
    INDEX_NAME,
    GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) AS index_columns,
    CARDINALITY,
    TABLE_ROWS,
    ROUND(CARDINALITY / TABLE_ROWS * 100, 2) AS selectivity_pct
FROM information_schema.STATISTICS s
JOIN information_schema.TABLES t
  ON s.TABLE_SCHEMA = t.TABLE_SCHEMA
  AND s.TABLE_NAME = t.TABLE_NAME
WHERE s.TABLE_SCHEMA = DATABASE()
  AND s.INDEX_NAME != 'PRIMARY'
  AND t.TABLE_ROWS > 0
GROUP BY TABLE_NAME, INDEX_NAME, CARDINALITY, TABLE_ROWS
ORDER BY selectivity_pct DESC;

-- ============================================
-- 第九部分: 性能验证查询
-- ============================================

-- 9.1 验证知识库查询优化效果
-- -----------------------------
EXPLAIN SELECT *
FROM knowledge
WHERE app_id = 1
  AND space_id = 1
  AND status = 1
ORDER BY created_at DESC
LIMIT 20;
-- 预期: type=ref, Extra=Using where; Using index

-- 9.2 验证文档查询优化效果
-- ---------------------------
EXPLAIN SELECT id, name, status, created_at, size
FROM knowledge_document
WHERE knowledge_id = 1
  AND status = 1
ORDER BY created_at DESC
LIMIT 20;
-- 预期: type=ref, Extra=Using where; Using index (覆盖索引)

-- 9.3 验证消息查询优化效果
-- ---------------------------
EXPLAIN SELECT *
FROM message
WHERE conversation_id = 1
  AND status = 1
ORDER BY created_at DESC
LIMIT 50;
-- 预期: type=ref, Extra=Using where

-- 9.4 验证权限查询优化效果
-- ---------------------------
EXPLAIN SELECT *
FROM data_permission
WHERE role_id = 1
  AND resource_type = 'bot'
  AND resource_id = '123';
-- 预期: type=ref, Extra=Using index condition

-- ============================================
-- 执行说明
-- ============================================
--
-- 1. 执行前准备:
--    - 确认MySQL版本 >= 8.0 (支持IF NOT EXISTS)
--    - 备份数据库: mysqldump -u root -p coze_studio > backup.sql
--    - 检查磁盘空间: 至少需要2倍当前数据库大小的空间
--
-- 2. 执行方式:
--    - 完整执行: mysql -u root -p coze_studio < 20250103120000_performance_indexes.sql
--    - 分步执行: 逐段执行,观察每个索引的创建时间
--
-- 3. 执行时间估算:
--    - 小表(<10万行): 每个索引 1-5秒
--    - 中表(10-100万行): 每个索引 10-30秒
--    - 大表(>100万行): 每个索引 30-120秒
--    - 总计: 约30-60分钟(取决于数据量)
--
-- 4. 执行后验证:
--    - 运行第八部分的索引维护查询
--    - 运行第九部分的性能验证查询
--    - 执行应用性能测试,对比优化前后效果
--
-- 5. 回滚方案:
--    - 如需回滚,执行以下命令删除索引:
--      DROP INDEX idx_xxx ON table_name;
--    - 或从备份恢复数据库
--
-- 6. 注意事项:
--    - 索引创建会锁表(在线DDL会最小化锁表时间)
--    - 建议在业务低峰期执行
--    - 监控磁盘IO和CPU使用率
--    - 执行完成后监控索引碎片,必要时执行 ANALYZE TABLE
--
-- ============================================
-- 版本历史
-- ============================================
-- v1.0.0 (2025-01-03): 初始版本,创建28个性能优化索引
-- ============================================
