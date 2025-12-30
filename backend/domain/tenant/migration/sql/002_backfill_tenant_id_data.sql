-- ================================================================================
-- 002_backfill_tenant_id_data.sql
-- 回填 tenant_id 数据
-- ================================================================================
-- 用途: 为历史数据分配正确的 tenant_id
-- 策略: 根据用户关联关系分配租户ID
-- ================================================================================

-- 设置批处理大小
SET @batch_size = 1000;

-- ================================================================================
-- 1. 为 users 表回填 tenant_id
-- ================================================================================
-- 策略:
--   - 管理员用户 -> system_tenant
--   - 普通用户 -> 创建个人租户
-- ================================================================================

-- Step 1: 创建系统租户（如果不存在）
INSERT INTO tenants (tenant_id, tenant_name, tenant_type, status, subscription_tier, created_at, updated_at)
SELECT 'system_tenant', 'System Tenant', 'system', 'active', 'enterprise', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM tenants WHERE tenant_id = 'system_tenant');

-- Step 2: 为普通用户创建个人租户
INSERT INTO tenants (tenant_id, tenant_name, tenant_type, status, subscription_tier, created_at, updated_at)
SELECT DISTINCT
    CONCAT('tenant_', u.id) as tenant_id,
    CONCAT(u.nickname, ' 的个人空间') as tenant_name,
    'individual' as tenant_type,
    'active' as status,
    'free' as subscription_tier,
    NOW() as created_at,
    NOW() as updated_at
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM tenants t WHERE t.tenant_id = CONCAT('tenant_', u.id)
)
AND u.id NOT IN (SELECT user_id FROM user_roles WHERE role_id = 'admin');

-- Step 3: 更新 users 表的 tenant_id
UPDATE users u
LEFT JOIN user_roles ur ON u.id = ur.user_id AND ur.role_id = 'admin'
SET u.tenant_id = CASE
    WHEN ur.user_id IS NOT NULL THEN 'system_tenant'
    ELSE CONCAT('tenant_', u.id)
END
WHERE u.tenant_id = 'system_tenant' OR u.tenant_id IS NULL;

-- ================================================================================
-- 2. 为 bot_draft 表回填 tenant_id
-- ================================================================================

-- 分批更新（避免锁表）
SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE bot_draft bd
    INNER JOIN (
        SELECT id, creator_id
        FROM bot_draft
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON bd.id = batch.id
    INNER JOIN users u ON bd.creator_id = u.id
    SET bd.tenant_id = u.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- ================================================================================
-- 3. 为 bot_published 表回填 tenant_id
-- ================================================================================

SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE bot_published bp
    INNER JOIN (
        SELECT id, creator_id
        FROM bot_published
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON bp.id = batch.id
    INNER JOIN users u ON bp.creator_id = u.id
    SET bp.tenant_id = u.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- ================================================================================
-- 4. 为 conversation 表回填 tenant_id
-- ================================================================================

SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE conversation c
    INNER JOIN (
        SELECT id, creator_id
        FROM conversation
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON c.id = batch.id
    INNER JOIN users u ON c.creator_id = u.id
    SET c.tenant_id = u.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- ================================================================================
-- 5. 为 message 表回填 tenant_id
-- ================================================================================

-- 方式1: 通过 conversation 关联回填
SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE message m
    INNER JOIN (
        SELECT id, conversation_id
        FROM message
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON m.id = batch.id
    INNER JOIN conversation c ON m.conversation_id = c.id
    SET m.tenant_id = c.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- ================================================================================
-- 6. 为 knowledge 相关表回填 tenant_id
-- ================================================================================

-- knowledge 表
SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE knowledge k
    INNER JOIN (
        SELECT id, creator_id
        FROM knowledge
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON k.id = batch.id
    INNER JOIN users u ON k.creator_id = u.id
    SET k.tenant_id = u.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- knowledge_document 表
SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE knowledge_document kd
    INNER JOIN (
        SELECT id, knowledge_id
        FROM knowledge_document
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON kd.id = batch.id
    INNER JOIN knowledge k ON kd.knowledge_id = k.id
    SET kd.tenant_id = k.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- knowledge_document_slice 表
SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE knowledge_document_slice kds
    INNER JOIN (
        SELECT id, document_id
        FROM knowledge_document_slice
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON kds.id = batch.id
    INNER JOIN knowledge_document kd ON kds.document_id = kd.id
    SET kds.tenant_id = kd.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- ================================================================================
-- 7. 为 workflow 表回填 tenant_id
-- ================================================================================

SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE workflow w
    INNER JOIN (
        SELECT id, creator_id
        FROM workflow
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON w.id = batch.id
    INNER JOIN users u ON w.creator_id = u.id
    SET w.tenant_id = u.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- ================================================================================
-- 8. 为 files 表回填 tenant_id
-- ================================================================================

SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE files f
    INNER JOIN (
        SELECT id, creator_id
        FROM files
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON f.id = batch.id
    INNER JOIN users u ON f.creator_id = u.id
    SET f.tenant_id = u.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- ================================================================================
-- 9. 为 plugin 相关表回填 tenant_id
-- ================================================================================

-- plugin_draft 表
SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE plugin_draft pd
    INNER JOIN (
        SELECT id, creator_id
        FROM plugin_draft
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON pd.id = batch.id
    INNER JOIN users u ON pd.creator_id = u.id
    SET pd.tenant_id = u.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- plugin_published 表
SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE plugin_published pp
    INNER JOIN (
        SELECT id, creator_id
        FROM plugin_published
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON pp.id = batch.id
    INNER JOIN users u ON pp.creator_id = u.id
    SET pp.tenant_id = u.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- ================================================================================
-- 10. 为 database 相关表回填 tenant_id
-- ================================================================================

-- online_database_info 表
SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE online_database_info odi
    INNER JOIN (
        SELECT id, creator_id
        FROM online_database_info
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON odi.id = batch.id
    INNER JOIN users u ON odi.creator_id = u.id
    SET odi.tenant_id = u.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- draft_database_info 表
SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE draft_database_info ddi
    INNER JOIN (
        SELECT id, creator_id
        FROM draft_database_info
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON ddi.id = batch.id
    INNER JOIN users u ON ddi.creator_id = u.id
    SET ddi.tenant_id = u.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- ================================================================================
-- 11. 为 api_key 表回填 tenant_id
-- ================================================================================

SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE api_key ak
    INNER JOIN (
        SELECT id, creator_id
        FROM api_key
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON ak.id = batch.id
    INNER JOIN users u ON ak.creator_id = u.id
    SET ak.tenant_id = u.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- ================================================================================
-- 12. 为其他表回填 tenant_id
-- ================================================================================

-- agent_to_database 表
SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE agent_to_database atd
    INNER JOIN (
        SELECT id, agent_id
        FROM agent_to_database
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON atd.id = batch.id
    INNER JOIN bot_draft bd ON atd.agent_id = bd.id
    SET atd.tenant_id = bd.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- data_copy_task 表
SET @rows_affected = 1;
WHILE @rows_affected > 0 DO
    UPDATE data_copy_task dct
    INNER JOIN (
        SELECT id, creator_id
        FROM data_copy_task
        WHERE tenant_id = 'system_tenant'
        LIMIT @batch_size
    ) AS batch ON dct.id = batch.id
    INNER JOIN users u ON dct.creator_id = u.id
    SET dct.tenant_id = u.tenant_id;

    SET @rows_affected = ROW_COUNT();
END WHILE;

-- ================================================================================
-- 验证脚本
-- ================================================================================

-- 检查是否还有未分配tenant_id的记录
SELECT
    'users' as table_name,
    COUNT(*) as total,
    SUM(CASE WHEN tenant_id = 'system_tenant' THEN 1 ELSE 0 END) as system_tenant_count,
    SUM(CASE WHEN tenant_id != 'system_tenant' THEN 1 ELSE 0 END) as assigned_count
FROM users
UNION ALL
SELECT
    'bot_draft',
    COUNT(*),
    SUM(CASE WHEN tenant_id = 'system_tenant' THEN 1 ELSE 0 END),
    SUM(CASE WHEN tenant_id != 'system_tenant' THEN 1 ELSE 0 END)
FROM bot_draft
UNION ALL
SELECT
    'conversation',
    COUNT(*),
    SUM(CASE WHEN tenant_id = 'system_tenant' THEN 1 ELSE 0 END),
    SUM(CASE WHEN tenant_id != 'system_tenant' THEN 1 ELSE 0 END)
FROM conversation;

-- 验证租户分配是否正确（应该没有使用默认租户的记录）
SELECT
    t.table_name,
    t.default_tenant_count,
    t.total_count,
    CONCAT(ROUND(t.default_tenant_count / t.total_count * 100, 2), '%') as default_percentage
FROM (
    SELECT
        'users' as table_name,
        SUM(CASE WHEN tenant_id = 'system_tenant' THEN 1 ELSE 0 END) as default_tenant_count,
        COUNT(*) as total_count
    FROM users
) t
WHERE t.default_tenant_count > 0;
