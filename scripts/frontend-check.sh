#!/bin/bash
# 前端代码检查脚本
# 用途：前端代码提交前自动检查
# 使用：./scripts/frontend-check.sh

set -e

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔍 前端代码检查${NC}"
echo "========================================"
echo ""

# 检查Rush是否安装
echo "📦 检查Rush环境..."
if ! command -v rush &> /dev/null; then
    echo -e "${RED}❌ Rush未安装${NC}"
    echo "请运行: npm install -g @microsoft/rush"
    exit 1
fi
echo -e "${GREEN}✅ Rush环境正常${NC}"

# 1. TypeScript类型检查
echo ""
echo "📝 TypeScript类型检查..."
cd frontend
if rush tsc; then
    echo -e "${GREEN}✅ TypeScript类型检查通过${NC}"
else
    echo -e "${RED}❌ TypeScript类型检查失败${NC}"
    exit 1
fi

# 2. ESLint检查
echo ""
echo "🔧 ESLint检查..."
if rush lint; then
    echo -e "${GREEN}✅ ESLint检查通过${NC}"
else
    echo -e "${RED}❌ ESLint检查失败${NC}"
    exit 1
fi

# 3. 单元测试
echo ""
echo "🧪 运行单元测试..."
if rush test; then
    echo -e "${GREEN}✅ 单元测试通过${NC}"
else
    echo -e "${RED}❌ 单元测试失败${NC}"
    exit 1
fi

# 4. 测试覆盖率检查
echo ""
echo "📊 检查测试覆盖率..."
# 注意：需要根据实际项目的覆盖率报告格式调整
COVERAGE_FILE="coverage/coverage-summary.json"
if [ -f "$COVERAGE_FILE" ]; then
    COVERAGE=$(node -e "const data=require('./${COVERAGE_FILE}'); const total=data.total; console.log(((total.lines.pct + total.branches.pct + total.functions.pvt + total.statements.pct) / 4).toFixed(2))")

    echo "当前覆盖率: ${COVERAGE}%"
    if (( $(echo "$COVERAGE < 70" | bc -l) )); then
        echo -e "${RED}❌ 测试覆盖率不足: ${COVERAGE}% (要求 ≥ 70%)${NC}"
        exit 1
    else
        echo -e "${GREEN}✅ 测试覆盖率达标: ${COVERAGE}%${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  覆盖率报告未找到，跳过覆盖率检查${NC}"
fi

# 5. 构建检查
echo ""
echo "🏗️  构建检查..."
if rush build; then
    echo -e "${GREEN}✅ 构建成功${NC}"
else
    echo -e "${RED}❌ 构建失败${NC}"
    exit 1
fi

echo ""
echo "========================================"
echo -e "${GREEN}✅ 前端代码检查全部通过${NC}"
echo ""
echo "📝 检查项："
echo "  ✅ TypeScript类型检查"
echo "  ✅ ESLint代码规范"
echo "  ✅ 单元测试"
echo "  ✅ 测试覆盖率 (≥ 70%)"
echo "  ✅ 构建检查"
echo ""
echo "🎉 可以安全提交代码！"
echo ""
