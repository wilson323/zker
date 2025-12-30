-- ================================================================================
-- ZKER 多租户迁移 - tenant_id 数据回填脚本（简化版）
-- ================================================================================
-- 版本: v1.0-simple
-- 创建日期: 2025-01-01
-- 说明: 直接更新所有记录，适用于中小型数据集
--       必须在 01_add_tenant_id_to_business_tables.sql 之后执行
-- ================================================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ================================================================================
-- 第一步: 创建租户映射关系
-- ================================================================================

-- 1.1 创建临时租户映射表
DROP TEMPORARY TABLE IF EXISTS temp_tenant_mapping;
CREATE TEMPORARY TABLE temp_tenant_mapping (
    space_id BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    tenant_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=MEMORY;

-- 1.2 从现有的 space_id 生成租户映射
-- 逻辑: 每个 space_id 对应一个租户,租户ID格式为 'tenant_{space_id}'
INSERT INTO temp_tenant_mapping (space_id, tenant_id, tenant_name)
SELECT DISTINCT
    space_id,
    CONCAT('tenant_', space_id) as tenant_id,
    CONCAT('Space ', space_id, ' Tenant') as tenant_name
FROM app_draft
WHERE space_id IS NOT NULL AND space_id > 0;

-- 1.3 为没有 space_id 的数据创建默认租户
INSERT INTO temp_tenant_mapping (space_id, tenant_id, tenant_name)
VALUES (0, 'tenant_default', 'Default Tenant')
ON DUPLICATE KEY UPDATE tenant_id = tenant_id;

-- ================================================================================
-- 第二步: 数据回填 - 直接更新所有记录
-- ================================================================================

-- 2.1 应用草稿表 (app_draft)
UPDATE app_draft ad
INNER JOIN temp_tenant_mapping tm ON tm.space_id = IFNULL(ad.space_id, 0)
SET ad.tenant_id = tm.tenant_id
WHERE ad.tenant_id IS NULL;

SELECT CONCAT('✅ app_draft: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.2 应用发布记录表 (app_release_record)
UPDATE app_release_record arr
INNER JOIN temp_tenant_mapping tm ON tm.space_id = IFNULL(arr.space_id, 0)
SET arr.tenant_id = tm.tenant_id
WHERE arr.tenant_id IS NULL;

SELECT CONCAT('✅ app_release_record: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.3 应用对话模板草稿 (app_conversation_template_draft)
UPDATE app_conversation_template_draft actd
INNER JOIN temp_tenant_mapping tm ON tm.space_id = IFNULL(actd.space_id, 0)
SET actd.tenant_id = tm.tenant_id
WHERE actd.tenant_id IS NULL;

SELECT CONCAT('✅ app_conversation_template_draft: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.4 对话表 (conversation)
-- 对话表需要通过 connector_id 关联到 app,再找到 space_id
UPDATE conversation c
INNER JOIN app_draft ad ON c.connector_id = ad.id
LEFT JOIN temp_tenant_mapping tm ON tm.space_id = IFNULL(ad.space_id, 0)
SET c.tenant_id = IFNULL(tm.tenant_id, 'tenant_default')
WHERE c.tenant_id IS NULL;

SELECT CONCAT('✅ conversation: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.5 消息表 (message)
-- 消息表通过 conversation_id 关联
UPDATE message m
INNER JOIN conversation c ON m.conversation_id = c.id
SET m.tenant_id = c.tenant_id
WHERE m.tenant_id IS NULL AND c.tenant_id IS NOT NULL;

SELECT CONCAT('✅ message: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.6 知识库表 (knowledge)
UPDATE knowledge k
INNER JOIN temp_tenant_mapping tm ON tm.space_id = IFNULL(k.space_id, 0)
SET k.tenant_id = tm.tenant_id
WHERE k.tenant_id IS NULL;

SELECT CONCAT('✅ knowledge: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.7 知识库文档表 (knowledge_document)
UPDATE knowledge_document kd
INNER JOIN knowledge k ON kd.knowledge_id = k.id
SET kd.tenant_id = k.tenant_id
WHERE kd.tenant_id IS NULL AND k.tenant_id IS NOT NULL;

SELECT CONCAT('✅ knowledge_document: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.8 知识库文档分块表 (knowledge_document_slice)
UPDATE knowledge_document_slice kds
INNER JOIN knowledge_document kd ON kds.document_id = kd.id
SET kds.tenant_id = kd.tenant_id
WHERE kds.tenant_id IS NULL AND kd.tenant_id IS NOT NULL;

SELECT CONCAT('✅ knowledge_document_slice: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.9 文件表 (files)
-- 对于没有明确关联的文件,分配到默认租户
UPDATE files
SET tenant_id = 'tenant_default'
WHERE tenant_id IS NULL;

SELECT CONCAT('✅ files: 更新了 ', ROW_COUNT(), ' 条记录到默认租户') AS progress;

-- 2.10 Agent 工具草稿表 (agent_tool_draft)
UPDATE agent_tool_draft atd
INNER JOIN app_draft ad ON atd.agent_id = ad.id
SET atd.tenant_id = IFNULL(ad.tenant_id, 'tenant_default')
WHERE atd.tenant_id IS NULL;

SELECT CONCAT('✅ agent_tool_draft: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.11 Agent 工具版本表 (agent_tool_version)
UPDATE agent_tool_version atv
INNER JOIN app_draft ad ON atv.agent_id = ad.id
SET atv.tenant_id = IFNULL(ad.tenant_id, 'tenant_default')
WHERE atv.tenant_id IS NULL;

SELECT CONCAT('✅ agent_tool_version: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.12 Agent 数据库关联表 (agent_to_database)
UPDATE agent_to_database atd
INNER JOIN app_draft ad ON atd.agent_id = ad.id
SET atd.tenant_id = IFNULL(ad.tenant_id, 'tenant_default')
WHERE atd.tenant_id IS NULL;

SELECT CONCAT('✅ agent_to_database: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.13 对话流角色配置表 (chat_flow_role_config)
UPDATE chat_flow_role_config cfc
INNER JOIN app_draft ad ON cfc.workflow_id = ad.id
SET cfc.tenant_id = IFNULL(ad.tenant_id, 'tenant_default')
WHERE cfc.tenant_id IS NULL;

SELECT CONCAT('✅ chat_flow_role_config: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.14 连接器工作流版本表 (connector_workflow_version)
UPDATE connector_workflow_version cwv
INNER JOIN app_draft ad ON cwv.app_id = ad.id
SET cwv.tenant_id = IFNULL(ad.tenant_id, 'tenant_default')
WHERE cwv.tenant_id IS NULL;

SELECT CONCAT('✅ connector_workflow_version: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.15 在线数据库信息表 (online_database_info)
UPDATE online_database_info odi
INNER JOIN app_draft ad ON odi.app_id = ad.id
SET odi.tenant_id = IFNULL(ad.tenant_id, 'tenant_default')
WHERE odi.tenant_id IS NULL;

SELECT CONCAT('✅ online_database_info: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.16 草稿数据库信息表 (draft_database_info)
UPDATE draft_database_info ddi
INNER JOIN temp_tenant_mapping tm ON tm.space_id = IFNULL(ddi.space_id, 0)
SET ddi.tenant_id = tm.tenant_id
WHERE ddi.tenant_id IS NULL;

SELECT CONCAT('✅ draft_database_info: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.17 API 密钥表 (api_key)
-- API Key 分配到默认租户
UPDATE api_key
SET tenant_id = 'tenant_default'
WHERE tenant_id IS NULL;

SELECT CONCAT('✅ api_key: 更新了 ', ROW_COUNT(), ' 条记录到默认租户') AS progress;

-- 2.18 应用连接器发布关联表 (app_connector_release_ref)
UPDATE app_connector_release_ref acrr
INNER JOIN app_release_record arr ON acrr.record_id = arr.id
SET acrr.tenant_id = IFNULL(arr.tenant_id, 'tenant_default')
WHERE acrr.tenant_id IS NULL;

SELECT CONCAT('✅ app_connector_release_ref: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

-- 2.19 其他小表 - 直接分配默认租户
UPDATE app_dynamic_conversation_draft SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;
SELECT CONCAT('✅ app_dynamic_conversation_draft: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

UPDATE app_dynamic_conversation_online SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;
SELECT CONCAT('✅ app_dynamic_conversation_online: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

UPDATE app_static_conversation_draft SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;
SELECT CONCAT('✅ app_static_conversation_draft: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

UPDATE app_static_conversation_online SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;
SELECT CONCAT('✅ app_static_conversation_online: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

UPDATE app_conversation_template_online SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;
SELECT CONCAT('✅ app_conversation_template_online: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

UPDATE data_copy_task SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;
SELECT CONCAT('✅ data_copy_task: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

UPDATE node_execution SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;
SELECT CONCAT('✅ node_execution: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

UPDATE knowledge_document_review SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;
SELECT CONCAT('✅ knowledge_document_review: 更新了 ', ROW_COUNT(), ' 条记录') AS progress;

SET FOREIGN_KEY_CHECKS = 1;

-- ================================================================================
-- 第三步: 数据验证
-- ================================================================================

-- 检查是否还有 NULL 的 tenant_id
SELECT '========== 数据验证报告 ==========' AS '';

SELECT
    'app_draft' as table_name,
    COUNT(*) as total_records,
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) as null_tenant_id_count
FROM app_draft
UNION ALL
SELECT
    'app_release_record',
    COUNT(*),
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END)
FROM app_release_record
UNION ALL
SELECT
    'conversation',
    COUNT(*),
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END)
FROM conversation
UNION ALL
SELECT
    'message',
    COUNT(*),
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END)
FROM message
UNION ALL
SELECT
    'knowledge',
    COUNT(*),
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END)
FROM knowledge
UNION ALL
SELECT
    'knowledge_document',
    COUNT(*),
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END)
FROM knowledge_document;

-- 显示租户分布统计
SELECT '========== 租户分布统计 ==========' AS '';
SELECT
    tenant_id,
    COUNT(*) as record_count,
    MIN(created_at) as earliest_record,
    MAX(created_at) as latest_record
FROM app_draft
GROUP BY tenant_id
ORDER BY record_count DESC;

SELECT '========== 迁移完成 ==========' AS '';
