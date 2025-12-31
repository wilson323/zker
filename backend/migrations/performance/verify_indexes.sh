#!/bin/bash

# ============================================
# ZKER 性能优化 - 索引验证脚本
# 版本: v1.0.0
# 生成时间: 2025-01-03
# 用途: 验证性能索引的创建情况和效果
# ============================================

set -e

# 数据库配置
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASS="${DB_PASS:-}"
DB_NAME="${DB_NAME:-coze_studio}"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "==================================="
echo "ZKER 性能索引验证脚本"
echo "数据库: ${DB_NAME}"
echo "==================================="

# MySQL命令
MYSQL_CMD="mysql -h${DB_HOST} -P${DB_PORT} -u${DB_USER} -p${DB_PASS} ${DB_NAME}"

# 1. 检查索引是否创建成功
echo -e "\n${YELLOW}1. 检查索引创建情况...${NC}"

INDEXES=(
    "idx_knowledge_app_space_status_created:knowledge"
    "idx_document_knowledge_status_created:knowledge_document"
    "idx_message_conversation_status_created:message"
    "idx_conversation_tenant_space_updated:conversation"
    "idx_data_permission_role_resource:data_permission"
)

for index in "${INDEXES[@]}"; do
    INDEX_NAME=$(echo $index | cut -d: -f1)
    TABLE_NAME=$(echo $index | cut -d: -f2)

    COUNT=$(echo "SELECT COUNT(*) FROM information_schema.STATISTICS
                 WHERE TABLE_SCHEMA='${DB_NAME}'
                 AND TABLE_NAME='${TABLE_NAME}'
                 AND INDEX_NAME='${INDEX_NAME}'" | $MYSQL_CMD -N 2>/dev/null || echo "0")

    if [ "$COUNT" == "1" ]; then
        echo -e "  ${GREEN}✓${NC} ${INDEX_NAME} 已创建"
    else
        echo -e "  ${RED}✗${NC} ${INDEX_NAME} 未找到"
    fi
done

# 2. 检查索引选择性
echo -e "\n${YELLOW}2. 分析索引选择性...${NC}"

echo "SELECT
    TABLE_NAME,
    INDEX_NAME,
    ROUND(CARDINALITY / TABLE_ROWS * 100, 2) AS selectivity_pct
FROM information_schema.STATISTICS s
JOIN information_schema.TABLES t
  ON s.TABLE_SCHEMA = t.TABLE_SCHEMA
  AND s.TABLE_NAME = t.TABLE_NAME
WHERE s.TABLE_SCHEMA = '${DB_NAME}'
  AND s.INDEX_NAME IN (
    'idx_knowledge_app_space_status_created',
    'idx_document_knowledge_status_created',
    'idx_message_conversation_status_created',
    'idx_conversation_tenant_space_updated',
    'idx_data_permission_role_resource'
  )
  AND s.SEQ_IN_INDEX = 1
ORDER BY selectivity_pct DESC;" | $MYSQL_CMD 2>/dev/null

# 3. 执行EXPLAIN分析查询性能
echo -e "\n${YELLOW}3. 分析查询执行计划...${NC}"

# 知识库查询
echo -e "\n--- 知识库列表查询 ---"
echo "EXPLAIN SELECT *
FROM knowledge
WHERE app_id = 1
  AND space_id = 1
  AND status = 1
ORDER BY created_at DESC
LIMIT 20;" | $MYSQL_CMD 2>/dev/null | grep -E "type|Extra|key" || echo "查询执行"

# 文档查询
echo -e "\n--- 文档列表查询 ---"
echo "EXPLAIN SELECT id, name, status, created_at, size
FROM knowledge_document
WHERE knowledge_id = 1
  AND status = 1
ORDER BY created_at DESC
LIMIT 20;" | $MYSQL_CMD 2>/dev/null | grep -E "type|Extra|key" || echo "查询执行"

# 消息查询
echo -e "\n--- 消息列表查询 ---"
echo "EXPLAIN SELECT *
FROM message
WHERE conversation_id = 1
  AND status = 1
ORDER BY created_at DESC
LIMIT 50;" | $MYSQL_CMD 2>/dev/null | grep -E "type|Extra|key" || echo "查询执行"

# 权限查询
echo -e "\n--- 权限检查查询 ---"
echo "EXPLAIN SELECT *
FROM data_permission
WHERE role_id = 1
  AND resource_type = 'bot'
  AND resource_id = '123';" | $MYSQL_CMD 2>/dev/null | grep -E "type|Extra|key" || echo "查询执行"

# 4. 索引大小统计
echo -e "\n${YELLOW}4. 统计索引大小...${NC}"

echo "SELECT
    TABLE_NAME,
    INDEX_NAME,
    ROUND(STAT_VALUE * @@innodb_page_size / 1024 / 1024, 2) AS size_mb
FROM mysql.innodb_index_stats
WHERE database_name = '${DB_NAME}'
  AND stat_name = 'size'
  AND stat_description != 'Number of pages in the index'
  AND INDEX_NAME IN (
    'idx_knowledge_app_space_status_created',
    'idx_document_knowledge_status_created',
    'idx_message_conversation_status_created',
    'idx_conversation_tenant_space_updated',
    'idx_data_permission_role_resource'
  )
ORDER BY size_mb DESC;" | $MYSQL_CMD 2>/dev/null

# 5. 性能对比建议
echo -e "\n${YELLOW}5. 性能对比建议...${NC}"
echo "建议执行以下性能测试对比优化效果:"
echo "  1. 后端基准测试: cd backend/tests/performance && go test -bench=. -benchmem"
echo "  2. API压力测试: k6 run tests/performance/k6/api_test.js"
echo "  3. 数据库慢查询分析: SHOW ENGINE INNODB STATUS\\G"

echo -e "\n${GREEN}索引验证完成!${NC}"
echo "==================================="
