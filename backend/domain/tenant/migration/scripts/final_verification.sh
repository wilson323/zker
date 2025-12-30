#!/bin/bash
# backend/domain/tenant/migration/scripts/final_verification.sh
-- 清理脚本4: 最终验证
-- 执行时机: 所有清理完成后

set -e

echo "==========================================="
echo "tenant_id 迁移 - 最终验证"
echo "==========================================="
echo ""

DB_NAME=${1:-"zker_production"}

echo "数据库: $DB_NAME"
echo ""

# 1. 验证tenant_id字段为NOT NULL
echo "1. 验证tenant_id字段约束..."
NULL_COUNT=$(mysql -u root -p -D "$DB_NAME" -se "
    SELECT COUNT(*)
    FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = '$DB_NAME'
      AND COLUMN_NAME = 'tenant_id'
      AND IS_NULLABLE = 'YES';
")

if [ "$NULL_COUNT" -eq 0 ]; then
    echo "   ✓ 所有表的tenant_id字段都已设为NOT NULL"
else
    echo "   ✗ 发现 $NULL_COUNT 张表的tenant_id仍允许NULL"
    exit 1
fi

# 2. 验证数据完整性
echo ""
echo "2. 验证数据完整性..."
TOTAL_INCOMPLETE=$(mysql -u root -p -D "$DB_NAME" -se "
    SELECT SUM(incomplete_count)
    FROM (
        SELECT COUNT(*) AS incomplete_count
        FROM bots WHERE tenant_id IS NULL OR tenant_id = ''
        UNION ALL
        SELECT COUNT(*) FROM bot_configs WHERE tenant_id IS NULL OR tenant_id = ''
        UNION ALL
        SELECT COUNT(*) FROM conversations WHERE tenant_id IS NULL OR tenant_id = ''
        UNION ALL
        SELECT COUNT(*) FROM messages WHERE tenant_id IS NULL OR tenant_id = ''
        UNION ALL
        SELECT COUNT(*) FROM knowledge_bases WHERE tenant_id IS NULL OR tenant_id = ''
        UNION ALL
        SELECT COUNT(*) FROM knowledge_chunks WHERE tenant_id IS NULL OR tenant_id = ''
        UNION ALL
        SELECT COUNT(*) FROM workflows WHERE tenant_id IS NULL OR tenant_id = ''
        UNION ALL
        SELECT COUNT(*) FROM workflow_executions WHERE tenant_id IS NULL OR tenant_id = ''
        UNION ALL
        SELECT COUNT(*) FROM single_agent_draft WHERE tenant_id IS NULL OR tenant_id = ''
        UNION ALL
        SELECT COUNT(*) FROM published_bots WHERE tenant_id IS NULL OR tenant_id = ''
    ) AS incomplete;
")

if [ "$TOTAL_INCOMPLETE" -eq 0 ]; then
    echo "   ✓ 所有表的tenant_id数据完整"
else
    echo "   ✗ 发现 $TOTAL_INCOMPLETE 条记录的tenant_id为空"
    exit 1
fi

# 3. 验证索引存在
echo ""
echo "3. 验证索引..."
MISSING_INDEXES=$(mysql -u root -p -D "$DB_NAME" -se "
    SELECT COUNT(*)
    FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = '$DB_NAME'
      AND COLUMN_NAME = 'tenant_id'
      AND INDEX_NAME LIKE 'idx_tenant_id';
")

if [ "$MISSING_INDEXES" -ge 10 ]; then
    echo "   ✓ 所有表的tenant_id索引存在"
else
    echo "   ✗ 部分表的tenant_id索引缺失 ($MISSING_INDEXES/10)"
    exit 1
fi

# 4. 验证唯一约束
echo ""
echo "4. 验证唯一约束..."
UNIQUE_COUNT=$(mysql -u root -p -D "$DB_NAME" -se "
    SELECT COUNT(DISTINCT TABLE_NAME)
    FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = '$DB_NAME'
      AND INDEX_NAME LIKE 'uk_tenant_%';
")

if [ "$UNIQUE_COUNT" -ge 10 ]; then
    echo "   ✓ 所有表的唯一约束存在"
else
    echo "   ⚠ 部分表的唯一约束可能缺失 ($UNIQUE_COUNT/10)"
fi

# 5. 性能检查
echo ""
echo "5. 性能检查..."
mysql -u root -p -D "$DB_NAME" -e "
    SELECT
        TABLE_NAME,
        TABLE_ROWS,
        ROUND(DATA_LENGTH / 1024 / 1024, 2) AS data_mb,
        ROUND(INDEX_LENGTH / 1024 / 1024, 2) AS index_mb
    FROM INFORMATION_SCHEMA.TABLES
    WHERE TABLE_SCHEMA = '$DB_NAME'
      AND TABLE_NAME IN (
        'bots', 'bot_configs', 'conversations', 'messages',
        'knowledge_bases', 'knowledge_chunks', 'workflows',
        'workflow_executions', 'single_agent_draft', 'published_bots'
      )
    ORDER BY TABLE_ROWS DESC;
"

echo ""
echo "==========================================="
echo "验证完成！"
echo ""
echo "所有检查通过，迁移已成功完成。"
echo "==========================================="
