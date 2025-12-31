#!/bin/bash
# 架构合规检查系统验证脚本
# 用途: 验证所有文件和配置是否正确安装
# 使用: bash scripts/verify-architecture-check.sh

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔍 ZKER架构合规检查系统验证${NC}"
echo "========================================"
echo ""

PASSED=0
FAILED=0
WARNINGS=0

# 检查文件是否存在
check_file() {
    local file=$1
    local description=$2

    if [ -f "$file" ]; then
        echo -e "${GREEN}✅ $description${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}❌ $description - 文件不存在${NC}"
        ((FAILED++))
        return 1
    fi
}

# 检查文件是否可执行
check_executable() {
    local file=$1
    local description=$2

    if [ -x "$file" ]; then
        echo -e "${GREEN}✅ $description (可执行)${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "${YELLOW}⚠️  $description (不可执行)${NC}"
        ((WARNINGS++))
        return 1
    fi
}

# 检查脚本语法
check_syntax() {
    local file=$1
    local description=$2

    if bash -n "$file" 2>/dev/null; then
        echo -e "${GREEN}✅ $description (语法正确)${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}❌ $description (语法错误)${NC}"
        ((FAILED++))
        return 1
    fi
}

echo "1️⃣  检查GitHub Actions工作流..."
echo "----------------------------------------"
check_file ".github/workflows/architecture-compliance.yml" "GitHub Actions工作流"
echo ""

echo "2️⃣  检查golangci配置..."
echo "----------------------------------------"
check_file ".github/linters/.golangci.yml" "golangci配置"
echo ""

echo "3️⃣  检查Git Hooks..."
echo "----------------------------------------"
check_file ".githooks/pre-commit" "pre-commit钩子"
check_executable ".githooks/pre-commit" "pre-commit钩子"
check_syntax ".githooks/pre-commit" "pre-commit钩子"
echo ""

echo "4️⃣  检查脚本文件..."
echo "----------------------------------------"
check_file "scripts/architecture-check.sh" "架构检查脚本"
check_executable "scripts/architecture-check.sh" "架构检查脚本"
check_syntax "scripts/architecture-check.sh" "架构检查脚本"

check_file "scripts/install-hooks.sh" "Hooks安装脚本"
check_executable "scripts/install-hooks.sh" "Hooks安装脚本"
check_syntax "scripts/install-hooks.sh" "Hooks安装脚本"
echo ""

echo "5️⃣  检查文档文件..."
echo "----------------------------------------"
check_file "scripts/README-ARCHITECTURE-CHECK.md" "脚本README"
check_file "docs/企业级功能完善与统一性设计方案/ZKER-CI-CD架构合规检查系统使用指南_v1.0.md" "完整使用指南"
echo ""

echo "6️⃣  检查Git Hooks安装状态..."
echo "----------------------------------------"
if [ -f ".git/hooks/pre-commit" ]; then
    echo -e "${GREEN}✅ pre-commit钩子已安装到.git/hooks/${NC}"
    ((PASSED++))

    if [ -x ".git/hooks/pre-commit" ]; then
        echo -e "${GREEN}✅ pre-commit钩子可执行${NC}"
        ((PASSED++))
    else
        echo -e "${YELLOW}⚠️  pre-commit钩子不可执行${NC}"
        echo "   运行: chmod +x .git/hooks/pre-commit"
        ((WARNINGS++))
    fi

    # 检查是否是我们安装的钩子
    if grep -q "ZKER架构合规性检查" .git/hooks/pre-commit 2>/dev/null; then
        echo -e "${GREEN}✅ pre-commit钩子版本正确${NC}"
        ((PASSED++))
    else
        echo -e "${YELLOW}⚠️  pre-commit钩子可能不是ZKER版本${NC}"
        ((WARNINGS++))
    fi
else
    echo -e "${YELLOW}⚠️  pre-commit钩子未安装${NC}"
    echo "   运行: bash scripts/install-hooks.sh"
    ((WARNINGS++))
fi
echo ""

echo "7️⃣  检查必需工具..."
echo "----------------------------------------"

# 检查Go
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✅ Go已安装 ($GO_VERSION)${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ Go未安装${NC}"
    ((FAILED++))
fi

# 检查golangci-lint
if command -v golangci-lint &> /dev/null; then
    LINT_VERSION=$(golangci-lint --version 2>/dev/null | head -1)
    echo -e "${GREEN}✅ golangci-lint已安装${NC}"
    echo "   $LINT_VERSION"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠️  golangci-lint未安装${NC}"
    echo "   安装方法: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$(go env GOPATH)/bin"
    ((WARNINGS++))
fi

# 检查gofmt
if command -v gofmt &> /dev/null; then
    echo -e "${GREEN}✅ gofmt已安装${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ gofmt未安装${NC}"
    ((FAILED++))
fi

# 检查goimports
if command -v goimports &> /dev/null; then
    echo -e "${GREEN}✅ goimports已安装${NC}"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠️  goimports未安装${NC}"
    echo "   安装方法: go install golang.org/x/tools/cmd/goimports@latest"
    ((WARNINGS++))
fi

# 检查gosec
if command -v gosec &> /dev/null; then
    echo -e "${GREEN}✅ gosec已安装${NC}"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠️  gosec未安装${NC}"
    echo "   安装方法: go install github.com/securego/gosec/v2/cmd/gosec@latest"
    ((WARNINGS++))
fi

echo ""
echo "========================================"
echo "验证完成！"
echo ""
echo "统计:"
echo "  ✅ 通过: $PASSED"
echo "  ⚠️  警告: $WARNINGS"
echo "  ❌ 失败: $FAILED"
echo ""

if [ $FAILED -eq 0 ]; then
    if [ $WARNINGS -eq 0 ]; then
        echo -e "${GREEN}🎉 所有检查通过！架构合规检查系统已正确安装${NC}"
        echo ""
        echo "下一步:"
        echo "  1. 运行架构检查: bash scripts/architecture-check.sh"
        echo "  2. 提交代码测试: git commit -m 'test: verify architecture check'"
        echo ""
        exit 0
    else
        echo -e "${YELLOW}⚠️  系统已安装，但有一些警告${NC}"
        echo ""
        echo "建议:"
        echo "  1. 安装缺失的工具（如golangci-lint）"
        echo "  2. 运行: bash scripts/install-hooks.sh"
        echo ""
        exit 0
    fi
else
    echo -e "${RED}❌ 验证失败，请检查上述错误${NC}"
    echo ""
    echo "常见问题:"
    echo "  1. 文件缺失 - 请检查文件是否存在"
    echo "  2. 语法错误 - 请检查脚本语法"
    echo "  3. 权限问题 - 运行: chmod +x scripts/*.sh .githooks/*"
    echo ""
    exit 1
fi
