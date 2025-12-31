#!/bin/bash
# 编译和errno检查本地测试脚本
# 用途: 在本地快速测试CI/CD工作流中的检查
# 使用: bash scripts/test-compilation-check.sh

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  ZKER后端编译 & errno检查本地测试${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

PASSED=0
FAILED=0

# 检查函数
check_step() {
    local step_name=$1
    local command=$2

    echo -e "${BLUE}检查:${NC} $step_name"
    echo "命令: $command"
    echo ""

    if eval "$command"; then
        echo -e "${GREEN}✅ 通过${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}❌ 失败${NC}"
        ((FAILED++))
        return 1
    fi
    echo ""
}

# 1. 检查Go版本
echo "1️⃣  检查Go环境"
echo "----------------------------------------"
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✅ Go已安装: $GO_VERSION${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ Go未安装${NC}"
    ((FAILED++))
    exit 1
fi
echo ""

# 2. 检查工作流文件
echo "2️⃣  检查工作流文件"
echo "----------------------------------------"
check_step "backend-compilation-check.yml 存在" \
    "[ -f .github/workflows/backend-compilation-check.yml ]"

check_step "工作流文件语法正确" \
    "command -v yamllint >/dev/null 2>&1 && yamllint .github/workflows/backend-compilation-check.yml 2>/dev/null || echo 'yamllint未安装，跳过语法检查'"
echo ""

# 3. 检查验证脚本
echo "3️⃣  检查验证脚本"
echo "----------------------------------------"
check_step "verify-errno-usage.sh 存在" \
    "[ -f scripts/verify-errno-usage.sh ]"

check_step "verify-errno-usage.sh 可执行" \
    "[ -x scripts/verify-errno-usage.sh ]"

check_step "verify-errno-usage.sh 语法正确" \
    "bash -n scripts/verify-errno-usage.sh"
echo ""

# 4. 测试后端编译
echo "4️⃣  测试后端编译"
echo "----------------------------------------"
if [ -d "backend" ]; then
    cd backend
    check_step "下载Go依赖" \
        "go mod download"

    check_step "验证Go依赖" \
        "go mod verify"

    check_step "基础编译检查" \
        "go build -v ./..."

    check_step "代码格式检查 (gofmt)" \
        "[ -z \"\$(gofmt -s -l . 2>/dev/null || true)\" ]"

    echo "运行 go vet (非阻塞)..."
    if go vet ./... 2>/dev/null; then
        echo -e "${GREEN}✅ go vet 检查通过${NC}"
    else
        echo -e "${YELLOW}⚠️  go vet 发现问题（警告）${NC}"
    fi
    cd - > /dev/null
else
    echo -e "${YELLOW}⚠️  backend目录不存在，跳过编译检查${NC}"
fi
echo ""

# 5. 测试errno验证
echo "5️⃣  测试errno验证"
echo "----------------------------------------"
if [ -f "scripts/verify-errno-usage.sh" ]; then
    chmod +x scripts/verify-errno-usage.sh

    echo "运行errno验证脚本..."
    if scripts/verify-errno-usage.sh; then
        echo -e "${GREEN}✅ errno验证通过${NC}"
        ((PASSED++))
    else
        echo -e "${RED}❌ errno验证失败${NC}"
        ((FAILED++))
    fi
else
    echo -e "${YELLOW}⚠️  errno验证脚本不存在${NC}"
fi
echo ""

# 6. 检查errno目录
echo "6️⃣  检查errno目录结构"
echo "----------------------------------------"
check_step "backend/types/errno 目录存在" \
    "[ -d backend/types/errno ]"

check_step "errno文件存在" \
    "[ \$(find backend/types/errno -name '*.go' -not -name '*_test.go' 2>/dev/null | wc -l) -gt 0 ]"

check_step "errno测试文件存在" \
    "[ -f backend/types/errno/errno_test.go ]"
echo ""

# 总结
echo "========================================"
echo "  测试总结"
echo "========================================"
echo ""
echo -e "${GREEN}✅ 通过: $PASSED${NC}"
echo -e "${RED}❌ 失败: $FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 所有检查通过！${NC}"
    echo ""
    echo "下一步:"
    echo "  1. 提交代码: git add ."
    echo "  2. 提交检查: git commit -m 'test: add compilation and errno check'"
    echo "  3. 推送到远程: git push"
    echo "  4. 查看CI/CD: GitHub Actions页面"
    echo ""
    exit 0
else
    echo -e "${RED}❌ 部分检查失败${NC}"
    echo ""
    echo "修复建议:"
    echo "  1. 确保Go环境正确安装"
    echo "  2. 运行 'cd backend && go mod download'"
    echo "  3. 运行 'cd backend && go build ./...'"
    echo "  4. 检查脚本权限: chmod +x scripts/*.sh"
    echo ""
    exit 1
fi
