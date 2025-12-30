#!/bin/bash

# 组织中心集成测试运行脚本
# 用途: 简化测试执行，支持多种运行模式

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 切换到脚本所在目录
cd "$(dirname "$0")"
cd ../../

# 显示帮助信息
show_help() {
    echo "组织中心集成测试运行脚本"
    echo ""
    echo "用法: ./tests/integration/run-tests.sh [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help              显示帮助信息"
    echo "  -v, --verbose           详细输出模式"
    echo "  -c, --cover             生成覆盖率报告"
    echo "  -r, --run <name>        运行特定测试"
    echo "  -p, --parallel <num>    并行运行测试（默认4）"
    echo "  -s, --short             跳过集成测试"
    echo "  --clean                 清理测试缓存"
    echo ""
    echo "示例:"
    echo "  ./tests/integration/run-tests.sh                 # 运行所有测试"
    echo "  ./tests/integration/run-tests.sh -v              # 详细输出"
    echo "  ./tests/integration/run-tests.sh -c              # 生成覆盖率"
    echo "  ./tests/integration/run-tests.sh -r OrgCRUD      # 运行特定测试"
    echo "  ./tests/integration/run-tests.sh -p 8            # 8个并行"
    echo ""
}

# 默认参数
VERBOSE=""
COVER=""
RUN_TEST=""
PARALLEL="-parallel 4"
SHORT=""
CLEAN=""

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--verbose)
            VERBOSE="-v"
            shift
            ;;
        -c|--cover)
            COVER="-cover"
            shift
            ;;
        -r|--run)
            RUN_TEST="$2"
            shift 2
            ;;
        -p|--parallel)
            PARALLEL="-parallel $2"
            shift 2
            ;;
        -s|--short)
            SHORT="-short"
            shift
            ;;
        --clean)
            CLEAN="go clean -testcache"
            shift
            ;;
        *)
            echo -e "${RED}未知选项: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# 显示配置
echo -e "${GREEN}================================================${NC}"
echo -e "${GREEN}组织中心集成测试${NC}"
echo -e "${GREEN}================================================${NC}"
echo ""
echo "测试配置:"
echo "  详细输出: $VERBOSE"
echo "  覆盖率:   $COVER"
echo "  并行数:   $PARALLEL"
echo "  特定测试: $RUN_TEST"
echo "  跳过测试: $SHORT"
echo ""

# 清理缓存
if [ -n "$CLEAN" ]; then
    echo -e "${YELLOW}清理测试缓存...${NC}"
    go clean -testcache
    echo -e "${GREEN}✓ 缓存已清理${NC}"
    echo ""
fi

# 检查Docker
echo -e "${YELLOW}检查Docker环境...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗ Docker未安装${NC}"
    exit 1
fi

if ! docker ps &> /dev/null; then
    echo -e "${RED}✗ Docker未运行${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Docker环境正常${NC}"
echo ""

# 构建测试命令
TEST_CMD="go test $VERBOSE $COVER $PARALLEL $SHORT -tags=integration -timeout 10m ./tests/integration/..."

# 添加特定测试
if [ -n "$RUN_TEST" ]; then
    TEST_CMD="$TEST_CMD -run $RUN_TEST"
fi

# 运行测试
echo -e "${YELLOW}运行测试...${NC}"
echo ""
echo "命令: $TEST_CMD"
echo ""

# 执行测试
if eval $TEST_CMD; then
    echo ""
    echo -e "${GREEN}================================================${NC}"
    echo -e "${GREEN}✓ 所有测试通过！${NC}"
    echo -e "${GREEN}================================================${NC}"
    exit 0
else
    echo ""
    echo -e "${RED}================================================${NC}"
    echo -e "${RED}✗ 测试失败${NC}"
    echo -e "${RED}================================================${NC}"
    exit 1
fi
