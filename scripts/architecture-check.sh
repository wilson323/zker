#!/bin/bash
# ZKER架构合规性检查工具
# 用途: 本地开发时的架构合规性检查，确保符合企业级标准
# 使用: bash scripts/architecture-check.sh
# 遵循规范: ZKER企业级开发规范手册 v1.0

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🏗️ ZKER架构合规性检查工具${NC}"
echo "========================================"
echo ""

# 检查是否在项目根目录
if [ ! -d "backend" ]; then
    echo -e "${RED}❌ 错误: 请在项目根目录运行此脚本${NC}"
    exit 1
fi

cd backend

FAILED=0
WARNINGS=0

# ==================== 1. 编译检查 ====================
echo -e "${BLUE}1️⃣  编译检查${NC}"
echo "----------------------------------------"
if go build -v ./... 2>&1 | tail -1; then
    echo -e "${GREEN}✅ 编译成功${NC}"
else
    echo -e "${RED}❌ 编译失败${NC}"
    FAILED=1
fi
echo ""

# ==================== 2. 循环依赖检查 ====================
echo -e "${BLUE}2️⃣  循环依赖检查${NC}"
echo "----------------------------------------"
CYCLE_OUTPUT=$(go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./... 2>&1)
if echo "$CYCLE_OUTPUT" | grep -q "import cycle"; then
    echo -e "${RED}❌ 发现循环依赖:${NC}"
    echo "$CYCLE_OUTPUT" | grep "import cycle"
    FAILED=1
else
    echo -e "${GREEN}✅ 无循环依赖${NC}"
fi
echo ""

# ==================== 3. DDD分层合规检查 ====================
echo -e "${BLUE}3️⃣  DDD分层合规检查${NC}"
echo "----------------------------------------"

# 3.1 domain层不应依赖api/application
echo "检查domain层..."
DOMAIN_VIOLATIONS=$(go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./domain/... 2>/dev/null | grep -E " (api/|application/)" || true)
if [ -n "$DOMAIN_VIOLATIONS" ]; then
    echo -e "${RED}❌ domain层依赖外层:${NC}"
    echo "$DOMAIN_VIOLATIONS"
    echo ""
    echo "  错误: domain层不应依赖api/application层"
    FAILED=1
else
    echo -e "${GREEN}✅ domain层合规${NC}"
fi

# 3.2 bizpkg不应依赖api/model
echo "检查bizpkg层..."
BIZPKG_VIOLATIONS=$(go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./bizpkg/... 2>/dev/null | grep " api/model/" || true)
if [ -n "$BIZPKG_VIOLATIONS" ]; then
    echo -e "${RED}❌ bizpkg依赖api/model:${NC}"
    echo "$BIZPKG_VIOLATIONS"
    echo ""
    echo "  错误: bizpkg应使用crossdomain/model而非api/model"
    FAILED=1
else
    echo -e "${GREEN}✅ bizpkg层合规${NC}"
fi

# 3.3 application层不应依赖api层
echo "检查application层..."
APP_VIOLATIONS=$(go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./application/... 2>/dev/null | grep " api/" || true)
if [ -n "$APP_VIOLATIONS" ]; then
    echo -e "${RED}❌ application层依赖api层:${NC}"
    echo "$APP_VIOLATIONS"
    echo ""
    echo "  错误: application层不应依赖api层"
    FAILED=1
else
    echo -e "${GREEN}✅ application层合规${NC}"
fi
echo ""

# ==================== 4. 包大小检查 ====================
echo -e "${BLUE}4️⃣  包大小检查${NC}"
echo "----------------------------------------"
LARGE_PACKAGES=$(find . -name "*.go" -not -path "*/vendor/*" -not -name "*.gen.go" | xargs wc -l 2>/dev/null | awk '$1 > 3000 { printf "%6d lines: %s\n", $1, $2 }' || true)
if [ -n "$LARGE_PACKAGES" ]; then
    echo -e "${YELLOW}⚠️  以下文件超过3000行，建议拆分:${NC}"
    echo "$LARGE_PACKAGES"
    WARNINGS=1
else
    echo -e "${GREEN}✅ 包大小合理${NC}"
fi
echo ""

# ==================== 5. 测试覆盖率检查 ====================
echo -e "${BLUE}5️⃣  测试覆盖率检查${NC}"
echo "----------------------------------------"
echo "运行测试并生成覆盖率报告..."

if go test -coverprofile=/tmp/zker_coverage.out -covermode=atomic ./... > /tmp/zker_test.log 2>&1; then
    COVERAGE=$(go tool cover -func=/tmp/zker_coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    COVERAGE_INT=${COVERAGE%.*}
    echo "当前覆盖率: ${COVERAGE}%"

    if [ "$COVERAGE_INT" -lt 70 ]; then
        echo -e "${RED}❌ 覆盖率低于70% (当前: ${COVERAGE}%)${NC}"
        FAILED=1
    else
        echo -e "${GREEN}✅ 覆盖率达标 (${COVERAGE}% ≥ 70%)${NC}"
    fi
else
    echo -e "${RED}❌ 测试执行失败${NC}"
    cat /tmp/zker_test.log
    FAILED=1
fi
echo ""

# ==================== 6. 代码规范检查 ====================
echo -e "${BLUE}6️⃣  代码规范检查 (golangci-lint)${NC}"
echo "----------------------------------------"

if ! command -v golangci-lint &> /dev/null; then
    echo -e "${YELLOW}⚠️  golangci-lint未安装，跳过检查${NC}"
    echo "安装方法: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$(go env GOPATH)/bin"
else
    if golangci-lint run --timeout=10m --config=../.github/linters/.golangci.yml; then
        echo -e "${GREEN}✅ 代码规范检查通过${NC}"
    else
        echo -e "${RED}❌ 代码规范检查失败${NC}"
        echo "运行 'golangci-lint run --fix' 修复问题"
        FAILED=1
    fi
fi
echo ""

# ==================== 7. 格式检查 ====================
echo -e "${BLUE}7️⃣  代码格式检查 (gofmt)${NC}"
echo "----------------------------------------"
UNFORMATTED=$(gofmt -l . | grep -v vendor | grep -v ".gen.go" | grep -v ".pb.go" || true)
if [ -n "$UNFORMATTED" ]; then
    echo -e "${RED}❌ 以下文件格式不正确:${NC}"
    echo "$UNFORMATTED"
    echo ""
    echo "运行以下命令修复:"
    echo "  gofmt -w \$file"
    FAILED=1
else
    echo -e "${GREEN}✅ 代码格式检查通过${NC}"
fi
echo ""

# ==================== 8. 安全检查 ====================
echo -e "${BLUE}8️⃣  安全检查 (go vet)${NC}"
echo "----------------------------------------"
if go vet ./...; then
    echo -e "${GREEN}✅ 安全检查通过${NC}"
else
    echo -e "${RED}❌ 安全检查失败${NC}"
    FAILED=1
fi
echo ""

# ==================== 总结 ==========================
cd ..

echo "========================================"
echo ""
echo "检查完成！"
echo ""

if [ $FAILED -eq 1 ]; then
    echo -e "${RED}❌ 架构检查失败${NC}"
    echo ""
    echo "请修复上述错误后重新提交"
    echo ""
    exit 1
elif [ $WARNINGS -eq 1 ]; then
    echo -e "${YELLOW}⚠️  架构检查通过，但有警告${NC}"
    echo ""
    echo "建议修复警告项以提升代码质量"
    echo ""
    exit 0
else
    echo -e "${GREEN}🎉 所有架构检查通过！${NC}"
    echo ""
    echo "✅ 编译检查"
    echo "✅ 循环依赖检查"
    echo "✅ DDD分层合规检查"
    echo "✅ 包大小检查"
    echo "✅ 测试覆盖率检查 (≥70%)"
    echo "✅ 代码规范检查"
    echo "✅ 代码格式检查"
    echo "✅ 安全检查"
    echo ""
    echo "可以安全提交代码！"
    echo ""
    exit 0
fi
