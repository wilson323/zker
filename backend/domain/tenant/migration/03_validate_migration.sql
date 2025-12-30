-- ================================================================================
-- ZKER 多租户迁移 - 数据一致性验证脚本
-- ================================================================================
-- 版本: v1.0
-- 创建日期: 2025-01-01
-- 说明: 验证 tenant_id 迁移的完整性和数据一致性
--       必须在数据回填完成后执行
-- ================================================================================

-- ================================================================================
-- 第一部分: 字段存在性验证
-- ================================================================================

-- 1.1 验证所有业务表都已添加 tenant_id 字段
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    DATA_TYPE,
    IS_NULLABLE,
    COLUMN_COMMENT
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'tenant_id'
  AND TABLE_NAME IN (
    'app_draft',
    'app_release_record',
    'app_conversation_template_draft',
    'conversation',
    'message',
    'knowledge',
    'knowledge_document',
    'knowledge_document_slice',
    'files',
    'agent_tool_draft',
    'agent_tool_version',
    'agent_to_database',
    'chat_flow_role_config',
    'connector_workflow_version',
    'online_database_info',
    'draft_database_info',
    'api_key',
    'app_connector_release_ref',
    'app_dynamic_conversation_draft',
    'app_dynamic_conversation_online',
    'app_static_conversation_draft',
    'app_static_conversation_online',
    'app_conversation_template_online',
    'data_copy_task',
    'node_execution',
    'knowledge_document_review'
  )
ORDER BY TABLE_NAME;

-- ================================================================================
-- 第二部分: 数据完整性验证
-- ================================================================================

-- 2.1 检查每张表的 tenant_id NULL 值情况
CREATE TEMPORARY TABLE temp_validation_results (
    table_name VARCHAR(100),
    total_records BIGINT,
    null_tenant_id BIGINT,
    null_percentage DECIMAL(5,2)
);

INSERT INTO temp_validation_results
SELECT
    'app_draft' as table_name,
    COUNT(*) as total_records,
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) as null_tenant_id,
    ROUND(SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2) as null_percentage
FROM app_draft
UNION ALL
SELECT
    'app_release_record',
    COUNT(*),
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END),
    ROUND(SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2)
FROM app_release_record
UNION ALL
SELECT
    'conversation',
    COUNT(*),
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END),
    ROUND(SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2)
FROM conversation
UNION ALL
SELECT
    'message',
    COUNT(*),
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END),
    ROUND(SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2)
FROM message
UNION ALL
SELECT
    'knowledge',
    COUNT(*),
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END),
    ROUND(SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2)
FROM knowledge
UNION ALL
SELECT
    'knowledge_document',
    COUNT(*),
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END),
    ROUND(SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2)
FROM knowledge_document
UNION ALL
SELECT
    'files',
    COUNT(*),
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END),
    ROUND(SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2)
FROM files;

-- 显示验证结果
SELECT
    table_name,
    total_records,
    null_tenant_id,
    null_percentage,
    CASE
        WHEN null_tenant_id = 0 THEN '✓ PASS'
        WHEN null_percentage < 1 THEN '⚠ WARNING (少量NULL值)'
        ELSE '✗ FAIL'
    END as validation_status
FROM temp_validation_results
ORDER BY null_percentage DESC;

-- ================================================================================
-- 第三部分: 关联关系验证
-- ================================================================================

-- 3.1 验证 app_draft 和 app_release_record 的 tenant_id 一致性
SELECT
    COUNT(*) as total_records,
    SUM(CASE WHEN ad.tenant_id <> arr.tenant_id THEN 1 ELSE 0 END) as mismatch_count,
    CASE
        WHEN SUM(CASE WHEN ad.tenant_id <> arr.tenant_id THEN 1 ELSE 0 END) = 0 THEN '✓ PASS'
        ELSE '✗ FAIL'
    END as validation_status
FROM app_draft ad
INNER JOIN app_release_record arr ON ad.id = arr.app_id;

-- 3.2 验证 conversation 和 message 的 tenant_id 一致性
SELECT
    COUNT(*) as total_records,
    SUM(CASE WHEN c.tenant_id <> m.tenant_id THEN 1 ELSE 0 END) as mismatch_count,
    CASE
        WHEN SUM(CASE WHEN c.tenant_id <> m.tenant_id THEN 1 ELSE 0 END) = 0 THEN '✓ PASS'
        ELSE '✗ FAIL'
    END as validation_status
FROM conversation c
INNER JOIN message m ON c.id = m.conversation_id;

-- 3.3 验证 knowledge 和 knowledge_document 的 tenant_id 一致性
SELECT
    COUNT(*) as total_records,
    SUM(CASE WHEN k.tenant_id <> kd.tenant_id THEN 1 ELSE 0 END) as mismatch_count,
    CASE
        WHEN SUM(CASE WHEN k.tenant_id <> kd.tenant_id THEN 1 ELSE 0 END) = 0 THEN '✓ PASS'
        ELSE '✗ FAIL'
    END as validation_status
FROM knowledge k
INNER JOIN knowledge_document kd ON k.id = kd.knowledge_id;

-- 3.4 验证 knowledge_document 和 knowledge_document_slice 的 tenant_id 一致性
SELECT
    COUNT(*) as total_records,
    SUM(CASE WHEN kd.tenant_id <> kds.tenant_id THEN 1 ELSE 0 END) as mismatch_count,
    CASE
        WHEN SUM(CASE WHEN kd.tenant_id <> kds.tenant_id THEN 1 ELSE 0 END) = 0 THEN '✓ PASS'
        ELSE '✗ FAIL'
    END as validation_status
FROM knowledge_document kd
INNER JOIN knowledge_document_slice kds ON kd.id = kds.document_id;

-- ================================================================================
-- 第四部分: 租户分布统计
-- ================================================================================

-- 4.1 各表的租户分布
SELECT
    'app_draft' as table_name,
    tenant_id,
    COUNT(*) as record_count,
    ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM app_draft), 2) as percentage
FROM app_draft
GROUP BY tenant_id
ORDER BY record_count DESC
LIMIT 10;

SELECT
    'conversation' as table_name,
    tenant_id,
    COUNT(*) as record_count,
    ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM conversation), 2) as percentage
FROM conversation
GROUP BY tenant_id
ORDER BY record_count DESC
LIMIT 10;

SELECT
    'message' as table_name,
    tenant_id,
    COUNT(*) as record_count,
    ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM message), 2) as percentage
FROM message
GROUP BY tenant_id
ORDER BY record_count DESC
LIMIT 10;

SELECT
    'knowledge' as table_name,
    tenant_id,
    COUNT(*) as record_count,
    ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM knowledge), 2) as percentage
FROM knowledge
GROUP BY tenant_id
ORDER BY record_count DESC
LIMIT 10;

-- ================================================================================
-- 第五部分: 数据抽样验证
-- ================================================================================

-- 5.1 抽样检查 app_draft 的租户分配
SELECT
    ad.id,
    ad.space_id,
    ad.tenant_id,
    ad.name,
    ad.created_at
FROM app_draft ad
WHERE ad.tenant_id IS NOT NULL
ORDER BY ad.created_at DESC
LIMIT 10;

-- 5.2 抽样检查 conversation 的租户分配
SELECT
    c.id,
    c.connector_id,
    c.tenant_id,
    c.name,
    c.created_at
FROM conversation c
WHERE c.tenant_id IS NOT NULL
ORDER BY c.created_at DESC
LIMIT 10;

-- 5.3 抽样检查 message 的租户分配
SELECT
    m.id,
    m.conversation_id,
    m.tenant_id,
    m.role,
    m.created_at
FROM message m
WHERE m.tenant_id IS NOT NULL
ORDER BY m.created_at DESC
LIMIT 10;

-- ================================================================================
-- 第六部分: 性能验证
-- ================================================================================

-- 6.1 检查索引是否正确创建
SELECT
    TABLE_NAME,
    INDEX_NAME,
    COLUMN_NAME,
    SEQ_IN_INDEX
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'tenant_id'
  AND TABLE_NAME IN (
    'app_draft',
    'app_release_record',
    'conversation',
    'message',
    'knowledge'
  )
ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX;

-- 6.2 测试查询性能(使用 tenant_id 过滤)
-- 添加 EXPLAIN 查看执行计划
EXPLAIN SELECT * FROM app_draft WHERE tenant_id = 'tenant_12345' LIMIT 100;

-- ================================================================================
-- 第七部分: 孤立记录检查
-- ================================================================================

-- 7.1 检查是否有孤立的消息(没有对应的对话)
SELECT
    COUNT(*) as orphan_messages
FROM message m
LEFT JOIN conversation c ON m.conversation_id = c.id
WHERE c.id IS NULL;

-- 7.2 检查是否有孤立的文档(没有对应的知识库)
SELECT
    COUNT(*) as orphan_documents
FROM knowledge_document kd
LEFT JOIN knowledge k ON kd.knowledge_id = k.id
WHERE k.id IS NULL;

-- 7.3 检查是否有孤立的文档分块(没有对应的文档)
SELECT
    COUNT(*) as orphan_slices
FROM knowledge_document_slice kds
LEFT JOIN knowledge_document kd ON kds.document_id = kd.id
WHERE kd.id IS NULL;

-- ================================================================================
-- 第八部分: 迁移摘要报告
-- ================================================================================

-- 8.1 总体迁移摘要
SELECT
    '========== 迁移验证摘要报告 ==========' as report_section
UNION ALL
SELECT
    CONCAT('1. 业务表总数: ', COUNT(*)) as report_section
FROM (
    SELECT 'app_draft' as tn
    UNION SELECT 'app_release_record'
    UNION SELECT 'conversation'
    UNION SELECT 'message'
    UNION SELECT 'knowledge'
    UNION SELECT 'knowledge_document'
    UNION SELECT 'files'
) t
UNION ALL
SELECT
    CONCAT('2. 已添加 tenant_id 字段的表: ', COUNT(*))
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'tenant_id'
  AND TABLE_NAME IN (
    'app_draft', 'app_release_record', 'conversation', 'message',
    'knowledge', 'knowledge_document', 'files'
  )
UNION ALL
SELECT
    CONCAT('3. tenant_id 非NULL的记录数: ',
        (SELECT COUNT(*) FROM app_draft WHERE tenant_id IS NOT NULL) +
        (SELECT COUNT(*) FROM conversation WHERE tenant_id IS NOT NULL) +
        (SELECT COUNT(*) FROM message WHERE tenant_id IS NOT NULL) +
        (SELECT COUNT(*) FROM knowledge WHERE tenant_id IS NOT NULL))
UNION ALL
SELECT
    CONCAT('4. 租户总数(基于space_id): ', COUNT(DISTINCT tenant_id))
FROM app_draft
WHERE tenant_id IS NOT NULL
UNION ALL
SELECT
    '========================================'
ORDER BY report_section;

-- 清理临时表
DROP TEMPORARY TABLE IF EXISTS temp_validation_results;

-- ================================================================================
-- 执行说明
-- ================================================================================
-- 1. 执行前确保已完成数据回填
-- 2. 执行本脚本: mysql -u root -p database_name < 03_validate_migration.sql
-- 3. 检查所有验证结果:
--    - 所有表的 NULL 值应为 0
--    - 关联关系验证应全部为 PASS
--    - 孤立记录数量应为 0
-- 4. 如果发现验证失败:
--    - 检查对应表的回填逻辑
--    - 补充缺失的数据
--    - 重新执行验证
-- 5. 验证通过后,可以开始启用双写适配器
-- ================================================================================

-- 迁移完成标记
-- 验证执行时间: ____________________
-- 验证执行人: ____________________
-- 验证结果: [] 全部通过  [ ] 存在问题需要修复
-- 审核人: ____________________
