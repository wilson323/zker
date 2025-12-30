-- ================================================================================
-- ZKER 多租户迁移 - tenant_id 数据回填脚本
-- ================================================================================
-- 版本: v1.0
-- 创建日期: 2025-01-01
-- 说明: 为现有数据回填 tenant_id,基于 space_id 创建租户映射
--       必须在 01_add_tenant_id_to_business_tables.sql 之后执行
--       执行前请务必备份数据库!
-- ================================================================================

SET FOREIGN_KEY_CHECKS = 0;
SET AUTOCOMMIT = 0;

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
-- 第二步: 数据回填 - 分批处理避免长事务
-- ================================================================================

-- 2.1 应用草稿表 (app_draft)
-- 分批处理,每批 10000 条
SET @batch_size = 10000;
SET @affected_rows = 1;

WHILE @affected_rows > 0 DO
    UPDATE app_draft
    SET tenant_id = (
        SELECT tm.tenant_id
        FROM temp_tenant_mapping tm
        WHERE tm.space_id = app_draft.space_id
        LIMIT 1
    )
    WHERE tenant_id IS NULL
    LIMIT @batch_size;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('app_draft: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.2 应用发布记录表 (app_release_record)
SET @affected_rows = 1;
WHILE @affected_rows > 0 DO
    UPDATE app_release_record
    SET tenant_id = (
        SELECT tm.tenant_id
        FROM temp_tenant_mapping tm
        WHERE tm.space_id = app_release_record.space_id
        LIMIT 1
    )
    WHERE tenant_id IS NULL
    LIMIT @batch_size;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('app_release_record: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.3 应用对话模板草稿 (app_conversation_template_draft)
SET @affected_rows = 1;
WHILE @affected_rows > 0 DO
    UPDATE app_conversation_template_draft
    SET tenant_id = (
        SELECT tm.tenant_id
        FROM temp_tenant_mapping tm
        WHERE tm.space_id = app_conversation_template_draft.space_id
        LIMIT 1
    )
    WHERE tenant_id IS NULL
    LIMIT @batch_size;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('app_conversation_template_draft: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.4 对话表 (conversation)
-- 对话表需要通过 connector_id 关联到 app,再找到 space_id
SET @affected_rows = 1;
WHILE @affected_rows > 0 DO
    UPDATE conversation c
    INNER JOIN (
        SELECT DISTINCT
            c.id as conv_id,
            tm.tenant_id
        FROM conversation c
        INNER JOIN app_draft ad ON c.connector_id = ad.id
        INNER JOIN temp_tenant_mapping tm ON tm.space_id = ad.space_id
        WHERE c.tenant_id IS NULL
        LIMIT @batch_size
    ) t ON c.id = t.conv_id
    SET c.tenant_id = t.tenant_id;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('conversation: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.5 消息表 (message)
-- 消息表通过 conversation_id 关联
SET @affected_rows = 1;
WHILE @affected_rows > 0 DO
    UPDATE message m
    INNER JOIN conversation c ON m.conversation_id = c.id
    SET m.tenant_id = c.tenant_id
    WHERE m.tenant_id IS NULL AND c.tenant_id IS NOT NULL
    LIMIT @batch_size;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('message: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.6 知识库表 (knowledge)
SET @affected_rows = 1;
WHILE @affected_rows > 0 DO
    UPDATE knowledge
    SET tenant_id = (
        SELECT tm.tenant_id
        FROM temp_tenant_mapping tm
        WHERE tm.space_id = knowledge.space_id
        LIMIT 1
    )
    WHERE tenant_id IS NULL
    LIMIT @batch_size;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('knowledge: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.7 知识库文档表 (knowledge_document)
SET @affected_rows = 1;
WHILE @affected_rows > 0 DO
    UPDATE knowledge_document kd
    INNER JOIN knowledge k ON kd.knowledge_id = k.id
    SET kd.tenant_id = k.tenant_id
    WHERE kd.tenant_id IS NULL AND k.tenant_id IS NOT NULL
    LIMIT @batch_size;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('knowledge_document: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.8 知识库文档分块表 (knowledge_document_slice)
SET @affected_rows = 1;
WHILE @affected_rows > 0 DO
    UPDATE knowledge_document_slice kds
    INNER JOIN knowledge_document kd ON kds.document_id = kd.id
    SET kds.tenant_id = kd.tenant_id
    WHERE kds.tenant_id IS NULL AND kd.tenant_id IS NOT NULL
    LIMIT @batch_size;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('knowledge_document_slice: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.9 文件表 (files)
-- 文件表使用 creator_id,需要查找关联的租户
-- 对于没有明确关联的文件,分配到默认租户
UPDATE files
SET tenant_id = 'tenant_default'
WHERE tenant_id IS NULL;

SELECT CONCAT('files: 更新了 ', ROW_COUNT(), ' 条记录到默认租户') AS progress;

-- 2.10 Agent 工具草稿表 (agent_tool_draft)
-- 通过 agent_id 关联到 app_draft
SET @affected_rows = 1;
WHILE @affected_rows > 0 DO
    UPDATE agent_tool_draft atd
    INNER JOIN app_draft ad ON atd.agent_id = ad.id
    SET atd.tenant_id = ad.tenant_id
    WHERE atd.tenant_id IS NULL AND ad.tenant_id IS NOT NULL
    LIMIT @batch_size;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('agent_tool_draft: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.11 Agent 工具版本表 (agent_tool_version)
SET @affected_rows = 1;
WHILE @affected_rows > 0 DO
    UPDATE agent_tool_version atv
    INNER JOIN app_draft ad ON atv.agent_id = ad.id
    SET atv.tenant_id = ad.tenant_id
    WHERE atv.tenant_id IS NULL AND ad.tenant_id IS NOT NULL
    LIMIT @batch_size;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('agent_tool_version: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.12 Agent 数据库关联表 (agent_to_database)
SET @affected_rows = 1;
WHILE @affected_rows > 0 DO
    UPDATE agent_to_database atd
    INNER JOIN app_draft ad ON atd.agent_id = ad.id
    SET atd.tenant_id = ad.tenant_id
    WHERE atd.tenant_id IS NULL AND ad.tenant_id IS NOT NULL
    LIMIT @batch_size;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('agent_to_database: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.13 对话流角色配置表 (chat_flow_role_config)
-- 通过 workflow_id 或 connector_id 关联
SET @affected_rows = 1;
WHILE @affected_rows > 0 DO
    UPDATE chat_flow_role_config cfc
    INNER JOIN app_draft ad ON cfc.workflow_id = ad.id
    SET cfc.tenant_id = ad.tenant_id
    WHERE cfc.tenant_id IS NULL AND ad.tenant_id IS NOT NULL
    LIMIT @batch_size;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('chat_flow_role_config: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.14 连接器工作流版本表 (connector_workflow_version)
SET @affected_rows = 1;
WHILE @affected_rows > 0 DO
    UPDATE connector_workflow_version cwv
    INNER JOIN app_draft ad ON cwv.app_id = ad.id
    SET cwv.tenant_id = ad.tenant_id
    WHERE cwv.tenant_id IS NULL AND ad.tenant_id IS NOT NULL
    LIMIT @batch_size;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('connector_workflow_version: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.15 在线数据库信息表 (online_database_info)
SET @affected_rows = 1;
WHILE @affected_rows > 0 DO
    UPDATE online_database_info odi
    INNER JOIN app_draft ad ON odi.app_id = ad.id
    SET odi.tenant_id = ad.tenant_id
    WHERE odi.tenant_id IS NULL AND ad.tenant_id IS NOT NULL
    LIMIT @batch_size;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('online_database_info: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.16 草稿数据库信息表 (draft_database_info)
SET @affected_rows = 1;
WHILE @affected_rows > 0 DO
    UPDATE draft_database_info ddi
    SET tenant_id = (
        SELECT tm.tenant_id
        FROM temp_tenant_mapping tm
        WHERE tm.space_id = ddi.space_id
        LIMIT 1
    )
    WHERE tenant_id IS NULL
    LIMIT @batch_size;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('draft_database_info: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.17 API 密钥表 (api_key)
-- API Key 分配到默认租户
UPDATE api_key
SET tenant_id = 'tenant_default'
WHERE tenant_id IS NULL;

SELECT CONCAT('api_key: 更新了 ', ROW_COUNT(), ' 条记录到默认租户') AS progress;

-- 2.18 应用连接器发布关联表 (app_connector_release_ref)
SET @affected_rows = 1;
WHILE @affected_rows > 0 DO
    UPDATE app_connector_release_ref acrr
    INNER JOIN app_release_record arr ON acrr.record_id = arr.id
    SET acrr.tenant_id = arr.tenant_id
    WHERE acrr.tenant_id IS NULL AND arr.tenant_id IS NOT NULL
    LIMIT @batch_size;

    SET @affected_rows = ROW_COUNT();
    COMMIT;
    SELECT CONCAT('app_connector_release_ref: 更新了 ', @affected_rows, ' 条记录') AS progress;
END WHILE;

-- 2.19 其他小表 - 直接分配默认租户
UPDATE app_dynamic_conversation_draft SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;
UPDATE app_dynamic_conversation_online SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;
UPDATE app_static_conversation_draft SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;
UPDATE app_static_conversation_online SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;
UPDATE app_conversation_template_online SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;
UPDATE data_copy_task SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;
UPDATE node_execution SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;
UPDATE knowledge_document_review SET tenant_id = 'tenant_default' WHERE tenant_id IS NULL;

SELECT CONCAT('其他小表: 更新完成') AS progress;

COMMIT;
SET AUTOCOMMIT = 1;

-- ================================================================================
-- 第三步: 数据验证
-- ================================================================================

-- 3.1 检查是否还有 NULL 的 tenant_id
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

-- 3.2 显示租户分布统计
SELECT
    tenant_id,
    COUNT(*) as record_count,
    MIN(created_at) as earliest_record,
    MAX(created_at) as latest_record
FROM app_draft
GROUP BY tenant_id
ORDER BY record_count DESC;

SET FOREIGN_KEY_CHECKS = 1;

-- ================================================================================
-- 执行说明
-- ================================================================================
-- 1. 执行前确保已完成 01_add_tenant_id_to_business_tables.sql
-- 2. 执行前备份: mysqldump -u root -p database_name > backup_before_backfill.sql
-- 3. 执行本脚本: mysql -u root -p database_name < 02_backfill_tenant_id_data.sql
-- 4. 监控执行进度,确保所有分批更新都完成
-- 5. 验证第三步的数据验证查询,确保所有记录都已填充 tenant_id
-- 6. 如果有 NULL 值,需要手动处理或调整回填逻辑
-- ================================================================================

-- 迁移完成标记
-- 请记录执行完成时间: ____________________
-- 执行人: ____________________
-- 审核人: ____________________
-- 验证结果: [] 所有记录已填充 tenant_id  [ ] 存在问题需要处理
