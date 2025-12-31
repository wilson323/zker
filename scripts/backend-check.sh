#!/bin/bash
# 后端代码检查脚本
# 用途：后端代码提交前自动检查
# 使用：./scripts/backend-check.sh

set -e

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔍 后端代码检查${NC}"
echo "========================================"
echo ""

# 检查Go环境
echo "📦 检查Go环境..."
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go未安装${NC}"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}✅ Go版本: $GO_VERSION${NC}"

cd backend

# 1. gofmt格式检查
echo ""
echo "📝 代码格式检查 (gofmt)..."
UNFORMATTED=$(gofmt -l . | grep -v vendor | grep .go || true)
if [ -z "$UNFORMATTED" ]; then
    echo -e "${GREEN}✅ 代码格式检查通过${NC}"
else
    echo -e "${RED}❌ 以下文件格式未通过gofmt检查:${NC}"
    echo "$UNFORMATTED"
    echo "请运行: gofmt -w <file>"
    exit 1
fi

# 2. golangci-lint检查
echo ""
echo "🔧 golangci-lint检查..."
if ! command -v golangci-lint &> /dev/null; then
    echo -e "${YELLOW}⚠️  golangci-lint未安装，跳过检查${NC}"
    echo "安装方法: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$(go env GOPATH)/bin"
else
    if golangci-lint run; then
        echo -e "${GREEN}✅ golangci-lint检查通过${NC}"
    else
        echo -e "${RED}❌ golangci-lint检查失败${NC}"
        exit 1
    fi
fi

# 3. 单元测试
echo ""
echo "🧪 运行单元测试..."
TEST_OUTPUT=$(go test ./... -cover 2>&1)
TEST_EXIT_CODE=$?

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✅ 单元测试通过${NC}"
else
    echo -e "${RED}❌ 单元测试失败${NC}"
    echo "$TEST_OUTPUT"
    exit 1
fi

# 4. 测试覆盖率检查
echo ""
echo "📊 检查测试覆盖率..."
# 计算平均覆盖率
COVERAGE=$(echo "$TEST_OUTPUT" | grep "coverage:" | awk '{sum+=$3; count++} END {if(count>0) print sum/count; else print 0}')

if [ -n "$COVERAGE" ]; then
    COVERAGE_INT=${COVERAGE%.*}
    echo "当前覆盖率: ${COVERAGE}%"
    if [ "$COVERAGE_INT" -lt 80 ]; then
        echo -e "${RED}❌ 测试覆盖率不足: ${COVERAGE}% (要求 ≥ 80%)${NC}"
        exit 1
    else
        echo -e "${GREEN}✅ 测试覆盖率达标: ${COVERAGE}%${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  无法获取覆盖率数据${NC}"
fi

# 5. 竞态检查
echo ""
echo "🏃 竞态检查 (-race)..."
if go test ./... -race -short; then
    echo -e "${GREEN}✅ 竞态检查通过${NC}"
else
    echo -e "${RED}❌ 发现竞态条件${NC}"
    exit 1
fi

# 6. 构建检查
echo ""
echo "🏗️  构建检查..."
if go build -o /dev/null ./...; then
    echo -e "${GREEN}✅ 构建成功${NC}"
else
    echo -e "${RED}❌ 构建失败${NC}"
    exit 1
fi

# 7. 安全检查 (go vet)
echo ""
echo "🔒 安全检查 (go vet)..."
if go vet ./...; then
    echo -e "${GREEN}✅ 安全检查通过${NC}"
else
    echo -e "${RED}❌ 安全检查失败${NC}"
    exit 1
fi

cd ..

echo ""
echo "========================================"
echo -e "${GREEN}✅ 后端代码检查全部通过${NC}"
echo ""
echo "📝 检查项："
echo "  ✅ 代码格式 (gofmt)"
echo "  ✅ 代码规范 (golangci-lint)"
echo "  ✅ 单元测试"
echo "  ✅ 测试覆盖率 (≥ 80%)"
echo "  ✅ 竞态检查 (-race)"
echo "  ✅ 构建检查"
echo "  ✅ 安全检查 (go vet)"
echo ""
echo "🎉 可以安全提交代码！"
echo ""
