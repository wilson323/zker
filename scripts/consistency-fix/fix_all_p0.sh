#!/bin/bash

#############################################
# ZKER P0问题一键修复工具
# 用途: 自动修复所有P0优先级问题
# 作者: ZKER开发团队
# 更新: 2025-01-03
#############################################

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}ZKER P0问题一键修复工具${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${YELLOW}⚠️  此脚本将自动修复所有P0问题${NC}"
echo -e "${YELLOW}   修复前会自动备份文件${NC}"
echo ""
echo -e "${BLUE}P0问题列表:${NC}"
echo "  1. 统一监控模块API响应格式"
echo "  2. 统一路由模块API响应格式"
echo "  3. 修复billing_engine错误处理"
echo "  4. 修复workflow核心错误处理"
echo "  5. 修复payment_service错误处理"
echo ""

read -p "是否继续? (y/N) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}已取消${NC}"
    exit 0
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}开始修复...${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# ========================================
# 1. 统一监控模块API响应格式
# ========================================
echo -e "${BLUE}[1/5] 统一监控模块API响应格式...${NC}"
echo ""

MONITORING_HANDLER="$PROJECT_ROOT/backend/api/handler/coze/monitoring/monitoring_handler.go"

if [ -f "$MONITORING_HANDLER" ]; then
    bash "$PROJECT_ROOT/scripts/consistency-fix/fix_api_response.sh" "$MONITORING_HANDLER"
    echo -e "${GREEN}✅ 监控模块API响应格式修复完成${NC}"
else
    echo -e "${YELLOW}⚠️  监控模块文件不存在，跳过${NC}"
fi

echo ""

# ========================================
# 2. 统一路由模块API响应格式
# ========================================
echo -e "${BLUE}[2/5] 统一路由模块API响应格式...${NC}"
echo ""

ROUTING_DIR="$PROJECT_ROOT/backend/api/handler/coze/routing"

if [ -d "$ROUTING_DIR" ]; then
    bash "$PROJECT_ROOT/scripts/consistency-fix/fix_api_response.sh" "$ROUTING_DIR"
    echo -e "${GREEN}✅ 路由模块API响应格式修复完成${NC}"
else
    echo -e "${YELLOW}⚠️  路由模块目录不存在，跳过${NC}"
fi

echo ""

# ========================================
# 3. 修复billing_engine错误处理
# ========================================
echo -e "${BLUE}[3/5] 修复billing_engine错误处理...${NC}"
echo ""

BILLING_ENGINE="$PROJECT_ROOT/backend/domain/billing/service/billing_engine.go"

if [ -f "$BILLING_ENGINE" ]; then
    python3 "$PROJECT_ROOT/scripts/consistency-fix/fix_error_handling.py" "$BILLING_ENGINE"
    echo -e "${GREEN}✅ billing_engine错误处理修复完成${NC}"
else
    echo -e "${YELLOW}⚠️  billing_engine文件不存在，跳过${NC}"
fi

echo ""

# ========================================
# 4. 修复workflow核心错误处理
# ========================================
echo -e "${BLUE}[4/5] 修复workflow核心错误处理...${NC}"
echo ""

WORKFLOW_DIR="$PROJECT_ROOT/backend/domain/workflow/internal/compose"

if [ -d "$WORKFLOW_DIR" ]; then
    python3 "$PROJECT_ROOT/scripts/consistency-fix/fix_error_handling.py" "$WORKFLOW_DIR"
    echo -e "${GREEN}✅ workflow核心错误处理修复完成${NC}"
else
    echo -e "${YELLOW}⚠️  workflow目录不存在，跳过${NC}"
fi

echo ""

# ========================================
# 5. 修复payment_service错误处理
# ========================================
echo -e "${BLUE}[5/5] 修复payment_service错误处理...${NC}"
echo ""

PAYMENT_SERVICE="$PROJECT_ROOT/backend/domain/billing/service/payment_service.go"

if [ -f "$PAYMENT_SERVICE" ]; then
    python3 "$PROJECT_ROOT/scripts/consistency-fix/fix_error_handling.py" "$PAYMENT_SERVICE"
    echo -e "${GREEN}✅ payment_service错误处理修复完成${NC}"
else
    echo -e "${YELLOW}⚠️  payment_service文件不存在，跳过${NC}"
fi

echo ""

# ========================================
# 汇总
# ========================================
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}✅ P0问题修复完成！${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

echo -e "${BLUE}下一步操作:${NC}"
echo "  1. 运行测试: cd backend && go test ./..."
echo "  2. 查看变更: git diff"
echo "  3. 提交修复: git add . && git commit -am 'fix: 修复P0问题'"
echo ""

echo -e "${YELLOW}⚠️  提示:${NC}"
echo "  - 所有备份文件使用 .backup 后缀"
echo "  - 请确认修复无误后删除备份文件"
echo "  - 如有问题可从备份文件恢复"
echo ""

read -p "是否删除备份文件? (y/N) " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}删除备份文件...${NC}"
    find "$PROJECT_ROOT/backend" -name "*.backup" -delete
    echo -e "${GREEN}✅ 备份文件已删除${NC}"
else
    echo -e "${YELLOW}保留备份文件${NC}"
    echo -e "${YELLOW}删除命令: find backend -name '*.backup' -delete${NC}"
fi

echo ""
echo -e "${GREEN}✅ 修复流程完成！${NC}"
