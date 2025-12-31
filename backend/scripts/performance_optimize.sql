-- =====================================================
-- ZKER 数据库性能优化索引脚本
-- 目标: 支持QPS 10000+, P99延迟 < 100ms
-- 创建时间: 2025-01-03
-- =====================================================

-- =====================================================
-- 1. 组织表优化 (organizations)
-- =====================================================

-- 租户隔离查询优化索引
CREATE INDEX IF NOT EXISTS idx_org_tenant_deleted
ON organizations(tenant_id, deleted_at)
COMMENT '租户隔离查询优化 - tenant_id + deleted_at';

-- 组织状态筛选优化
CREATE INDEX IF NOT EXISTS idx_org_tenant_status_deleted
ON organizations(tenant_id, status, deleted_at)
COMMENT '组织状态筛选 - tenant_id + status + deleted_at';

-- 组织编码查询优化
CREATE INDEX IF NOT EXISTS idx_org_tenant_code_deleted
ON organizations(tenant_id, org_code, deleted_at)
COMMENT '组织编码查询 - tenant_id + org_code + deleted_at';

-- 层级查询优化
CREATE INDEX IF NOT EXISTS idx_org_parent_level
ON organizations(parent_id, level, deleted_at)
COMMENT '组织层级查询 - parent_id + level + deleted_at';

-- =====================================================
-- 2. 部门表优化 (departments)
-- =====================================================

CREATE INDEX IF NOT EXISTS idx_dept_tenant_deleted
ON departments(tenant_id, deleted_at)
COMMENT '租户隔离查询 - tenant_id + deleted_at';

CREATE INDEX IF NOT EXISTS idx_dept_org_tenant
ON departments(org_id, tenant_id, deleted_at)
COMMENT '组织部门查询 - org_id + tenant_id + deleted_at';

CREATE INDEX IF NOT EXISTS idx_dept_parent_level
ON departments(parent_department_id, level, deleted_at)
COMMENT '部门层级查询 - parent_id + level';

-- =====================================================
-- 3. 权限表优化 (roles, data_permissions, field_permissions)
-- =====================================================

-- 角色表索引
CREATE INDEX IF NOT EXISTS idx_role_tenant_deleted
ON roles(tenant_id, deleted_at)
COMMENT '租户角色查询 - tenant_id + deleted_at';

CREATE INDEX IF NOT EXISTS idx_role_tenant_type_deleted
ON roles(tenant_id, role_type, deleted_at)
COMMENT '角色类型查询 - tenant_id + role_type + deleted_at';

CREATE INDEX IF NOT EXISTS idx_role_code_tenant
ON roles(role_code, tenant_id, deleted_at)
COMMENT '角色编码查询 - role_code + tenant_id + deleted_at';

-- 数据权限索引
CREATE INDEX IF NOT EXISTS idx_dataperm_role_resource
ON data_permissions(role_id, resource_type)
COMMENT '角色数据权限 - role_id + resource_type';

-- 字段权限索引
CREATE INDEX IF NOT EXISTS idx_fieldperm_role_resource
ON field_permissions(role_id, resource_type)
COMMENT '角色字段权限 - role_id + resource_type';

-- 用户角色关联索引
CREATE INDEX IF NOT EXISTS idx_userrole_user_tenant
ON user_roles(user_id, tenant_id)
COMMENT '用户角色查询 - user_id + tenant_id';

CREATE INDEX IF NOT EXISTS idx_userrole_role_tenant
ON user_roles(role_id, tenant_id)
COMMENT '角色用户查询 - role_id + tenant_id';

-- =====================================================
-- 4. Bot商店表优化 (bot_store_items, bot_store_categories)
-- =====================================================

-- Bot商品索引
CREATE INDEX IF NOT EXISTS idx_botstore_category_status_created
ON bot_store_items(category, status, created_at DESC)
COMMENT '商品列表查询 - category + status + created_at';

CREATE INDEX IF NOT EXISTS idx_botstore_publisher_tenant
ON bot_store_items(publisher_id, tenant_id)
COMMENT '发布者商品查询 - publisher_id + tenant_id';

CREATE INDEX IF NOT EXISTS idx_botstore_status_rating
ON bot_store_items(status, rating DESC)
COMMENT '商品评分排序 - status + rating';

CREATE INDEX IF NOT EXISTS idx_botstore_status_downloads
ON bot_store_items(status, download_count DESC)
COMMENT '商品下载量排序 - status + download_count';

-- 全文搜索索引 (MySQL 5.7.6+)
-- CREATE FULLTEXT INDEX IF NOT EXISTS idx_botstore_fulltext
-- ON bot_store_items(name, description, tags);

-- =====================================================
-- 5. Bot评论表优化 (bot_store_reviews)
-- =====================================================

CREATE INDEX IF NOT EXISTS idx_review_item_created
ON bot_store_reviews(item_id, created_at DESC)
COMMENT '商品评论列表 - item_id + created_at';

CREATE INDEX IF NOT EXISTS idx_review_user_tenant
ON bot_store_reviews(user_id, tenant_id)
COMMENT '用户评论查询 - user_id + tenant_id';

CREATE INDEX IF NOT EXISTS idx_review_item_user
ON bot_store_reviews(item_id, user_id)
COMMENT '用户商品评论 - item_id + user_id';

-- =====================================================
-- 6. 审计日志表优化 (audit_logs)
-- =====================================================

CREATE INDEX IF NOT EXISTS idx_audit_tenant_created
ON audit_logs(tenant_id, created_at DESC)
COMMENT '租户审计日志 - tenant_id + created_at';

CREATE INDEX IF NOT EXISTS idx_audit_tenant_action_created
ON audit_logs(tenant_id, action_type, created_at DESC)
COMMENT '租户操作日志 - tenant_id + action_type + created_at';

CREATE INDEX IF NOT EXISTS idx_audit_user_created
ON audit_logs(user_id, created_at DESC)
COMMENT '用户操作日志 - user_id + created_at';

-- 分区表建议 (大数据量场景)
-- ALTER TABLE audit_logs
-- PARTITION BY RANGE (UNIX_TIMESTAMP(created_at)) (
--     PARTITION p202501 VALUES LESS THAN (UNIX_TIMESTAMP('2025-02-01')),
--     PARTITION p202502 VALUES LESS THAN (UNIX_TIMESTAMP('2025-03-01')),
--     PARTITION p202503 VALUES LESS THAN (UNIX_TIMESTAMP('2025-04-01'))
-- );

-- =====================================================
-- 7. Token计量表优化 (token_metering)
-- =====================================================

CREATE INDEX IF NOT EXISTS idx_metering_tenant_date
ON token_metering(tenant_id, date DESC)
COMMENT '租户每日用量 - tenant_id + date';

CREATE INDEX IF NOT EXISTS idx_metering_tenant_model_date
ON token_metering(tenant_id, model_id, date DESC)
COMMENT '租户模型用量 - tenant_id + model_id + date';

CREATE INDEX IF NOT EXISTS idx_metering_conversation
ON token_metering(conversation_id, created_at)
COMMENT '会话token记录 - conversation_id + created_at';

-- =====================================================
-- 8. 订阅和配额表优化 (subscriptions, quotas)
-- =====================================================

CREATE INDEX IF NOT EXISTS idx_subscription_tenant_status
ON subscriptions(tenant_id, status)
COMMENT '租户订阅查询 - tenant_id + status';

CREATE INDEX IF NOT EXISTS idx_subscription_tenant_type
ON subscriptions(tenant_id, subscription_type)
COMMENT '租户订阅类型 - tenant_id + subscription_type';

CREATE INDEX IF NOT EXISTS idx_quota_tenant_resource
ON quotas(tenant_id, resource_type)
COMMENT '租户配额查询 - tenant_id + resource_type';

-- =====================================================
-- 9. 工作流表优化 (workflows, workflow_executions)
-- =====================================================

CREATE INDEX IF NOT EXISTS idx_workflow_tenant_deleted
ON workflows(tenant_id, deleted_at)
COMMENT '租户工作流 - tenant_id + deleted_at';

CREATE INDEX IF NOT EXISTS idx_workflow_tenant_status
ON workflows(tenant_id, status, deleted_at)
COMMENT '租户工作流状态 - tenant_id + status + deleted_at';

CREATE INDEX IF NOT EXISTS idx_workflow_execution_created
ON workflow_executions(workflow_id, created_at DESC)
COMMENT '工作流执行历史 - workflow_id + created_at';

CREATE INDEX IF NOT EXISTS idx_workflow_execution_status
ON workflow_executions(workflow_id, status, created_at DESC)
COMMENT '工作流执行状态 - workflow_id + status + created_at';

-- =====================================================
-- 10. 通用优化建议
-- =====================================================

-- 10.1 分析表并更新索引统计
ANALYZE TABLE organizations;
ANALYZE TABLE departments;
ANALYZE TABLE roles;
ANALYZE TABLE bot_store_items;
ANALYZE TABLE audit_logs;

-- 10.2 检查慢查询日志
-- SET GLOBAL slow_query_log = 'ON';
-- SET GLOBAL long_query_time = 0.1; -- 记录超过100ms的查询

-- 10.3 启用查询缓存 (MySQL 5.7及以下版本)
-- SET GLOBAL query_cache_type = ON;
-- SET GLOBAL query_cache_size = 268435456; -- 256MB

-- =====================================================
-- 性能验证查询
-- =====================================================

-- 验证索引是否创建成功
SELECT
    TABLE_NAME,
    INDEX_NAME,
    COLUMN_NAME,
    SEQ_IN_INDEX,
    INDEX_TYPE
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME IN (
        'organizations', 'departments', 'roles',
        'bot_store_items', 'audit_logs', 'token_metering'
    )
ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX;

-- 查看表大小和索引使用情况
SELECT
    TABLE_NAME,
    ROUND(((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024), 2) AS 'Size (MB)',
    TABLE_ROWS
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
ORDER BY (DATA_LENGTH + INDEX_LENGTH) DESC;

-- =====================================================
-- 回滚脚本 (如需删除索引)
-- =====================================================

/*
DROP INDEX idx_org_tenant_deleted ON organizations;
DROP INDEX idx_org_tenant_status_deleted ON organizations;
DROP INDEX idx_dept_tenant_deleted ON departments;
DROP INDEX idx_role_tenant_deleted ON roles;
DROP INDEX idx_botstore_category_status_created ON bot_store_items;
-- ... 其他索引
*/
