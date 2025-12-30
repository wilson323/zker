#!/bin/bash

# ============================================================================
# 租户表DDL执行脚本
# 版本: v1.0.0
# 日期: 2025-01-01
# 目的: 自动化执行租户表DDL脚本
# ============================================================================

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
MYSQL_HOST="${MYSQL_HOST:-localhost}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-123456}"
MYSQL_DATABASE="opencoze"
MIGRATIONS_DIR="$(cd "$(dirname "$0")/../docker/atlas/migrations" && pwd)"
VERIFY_SCRIPT="$(cd "$(dirname "$0")" && pwd)/verify_tenant_ddl.sql"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}租户表DDL部署脚本${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# 检查MySQL连接
echo -e "${YELLOW}检查MySQL连接...${NC}"
mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" -e "SELECT 1;" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ MySQL连接成功${NC}"
else
    echo -e "${RED}✗ MySQL连接失败${NC}"
    exit 1
fi
echo ""

# 执行DDL脚本
echo -e "${YELLOW}创建租户系统表...${NC}"
mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" < "$MIGRATIONS_DIR/20251230025000_create_tenant_tables.sql"
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 租户表创建成功${NC}"
else
    echo -e "${RED}✗ 租户表创建失败${NC}"
    exit 1
fi
echo ""

# 创建触发器
echo -e "${YELLOW}创建触发器...${NC}"
mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" < "$MIGRATIONS_DIR/20251230025001_tenant_status_trigger.sql"
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 触发器创建成功${NC}"
else
    echo -e "${RED}✗ 触发器创建失败${NC}"
    exit 1
fi
echo ""

# 创建视图
echo -e "${YELLOW}创建统计视图...${NC}"
mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" < "$MIGRATIONS_DIR/20251230025002_tenant_statistics_view.sql"
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 视图创建成功${NC}"
else
    echo -e "${RED}✗ 视图创建失败${NC}"
    exit 1
fi
echo ""

# 初始化数据
echo -e "${YELLOW}初始化数据...${NC}"
mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" < "$MIGRATIONS_DIR/20251230025003_init_tenant_data.sql"
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 数据初始化成功${NC}"
else
    echo -e "${RED}✗ 数据初始化失败${NC}"
    exit 1
fi
echo ""

# 验证部署
echo -e "${YELLOW}验证部署结果...${NC}"
if [ -f "$VERIFY_SCRIPT" ]; then
    mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" < "$VERIFY_SCRIPT"
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ 验证通过${NC}"
    else
        echo -e "${YELLOW}⚠ 验证脚本执行完成，请检查输出${NC}"
    fi
else
    echo -e "${YELLOW}⚠ 验证脚本不存在: $VERIFY_SCRIPT${NC}"
fi
echo ""

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}部署完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "后续步骤："
echo "1. 检查表结构: mysql -u$MYSQL_USER -p$MYSQL_PASSWORD $MYSQL_DATABASE -e 'SHOW TABLES LIKE \"%tenant%\";'"
echo "2. 查看系统租户: mysql -u$MYSQL_USER -p$MYSQL_PASSWORD $MYSQL_DATABASE -e 'SELECT * FROM tenants WHERE tenant_id=\"system-default\";'"
echo "3. 回滚命令: mysql -u$MYSQL_USER -p$MYSQL_PASSWORD $MYSQL_DATABASE < $MIGRATIONS_DIR/rollback/20251230025000_tenant_tables_rollback.sql"
