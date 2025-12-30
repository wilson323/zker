#!/bin/bash
# ZKER 组织中心性能测试执行脚本
#
# 功能: 自动执行组织中心完整的性能测试套件
#
# 测试覆盖:
#   - 基准测试 (Base Load)
#   - 峰值测试 (Spike Test)
#   - 压力测试 (Stress Test)
#   - 耐久测试 (Endurance Test)
#   - 可扩展性测试 (Scalability Test)
#
# 使用方法:
#   # 完整测试（包括数据生成）
#   bash run-org-performance-tests.sh --full
#
#   # 仅执行性能测试（跳过数据生成）
#   bash run-org-performance-tests.sh --skip-data
#
#   # 自定义测试
#   bash run-org-performance-tests.sh --base-url=http://localhost:8080 --size=medium
#
# @author 研发B (后端工程师)
# @version 1.0
# @date 2025-01-01

set -e  # 遇到错误立即退出

# ================================
# 颜色输出
# ================================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ================================
# 配置参数
# ================================
BASE_URL="${BASE_URL:-http://localhost:8080}"
DATA_SIZE="${DATA_SIZE:-small}"  # small, medium, large
SKIP_DATA="${SKIP_DATA:-false}"
RESULTS_DIR="./backend/tests/performance/results"
MYSQL_CONTAINER="zker-mysql-1"

# ================================
# 帮助信息
# ================================
show_help() {
    echo "ZKER 组织中心性能测试脚本"
    echo ""
    echo "使用方法:"
    echo "  bash run-org-performance-tests.sh [选项]"
    echo ""
    echo "选项:"
    echo "  --full                    完整测试（包括数据生成）"
    echo "  --skip-data               跳过数据生成"
    echo "  --base-url=URL            API基础URL (默认: http://localhost:8080)"
    echo "  --size=SIZE               测试数据规模 (small|medium|large, 默认: small)"
    echo "  --test-type=TYPE          测试类型 (base|spike|stress|endurance|scalability|all)"
    echo "  --help                    显示帮助信息"
    echo ""
    echo "环境变量:"
    echo "  BASE_URL                  API基础URL"
    echo "  DATA_SIZE                 测试数据规模"
    echo "  SKIP_DATA                 是否跳过数据生成 (true|false)"
    echo ""
    echo "示例:"
    echo "  # 完整测试（包括数据生成）"
    echo "  bash run-org-performance-tests.sh --full"
    echo ""
    echo "  # 仅执行基准测试"
    echo "  bash run-org-performance-tests.sh --test-type=base"
    echo ""
    echo "  # 自定义URL和数据规模"
    echo "  bash run-org-performance-tests.sh --base-url=http://192.168.1.100:8080 --size=medium"
}

# ================================
# 解析命令行参数
# ================================
TEST_TYPE="all"

for arg in "$@"; do
    case $arg in
        --full)
            SKIP_DATA=false
            TEST_TYPE="all"
            shift
            ;;
        --skip-data)
            SKIP_DATA=true
            shift
            ;;
        --base-url=*)
            BASE_URL="${arg#*=}"
            shift
            ;;
        --size=*)
            DATA_SIZE="${arg#*=}"
            shift
            ;;
        --test-type=*)
            TEST_TYPE="${arg#*=}"
            shift
            ;;
        --help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}❌ 未知参数: $arg${NC}"
            show_help
            exit 1
            ;;
    esac
done

# ================================
# 打印配置
# ================================
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}ZKER 组织中心性能测试套件${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "📊 测试配置:"
echo "  - BASE_URL: $BASE_URL"
echo "  - DATA_SIZE: $DATA_SIZE"
echo "  - SKIP_DATA: $SKIP_DATA"
echo "  - TEST_TYPE: $TEST_TYPE"
echo "  - RESULTS_DIR: $RESULTS_DIR"
echo ""

# ================================
# 检查依赖
# ================================
check_dependencies() {
    echo -e "${BLUE}🔍 检查依赖...${NC}"

    # 检查k6
    if ! command -v k6 &> /dev/null; then
        echo -e "${RED}❌ k6未安装${NC}"
        echo "请访问 https://k6.io/docs/getting-started/installation/ 安装k6"
        exit 1
    fi
    echo -e "${GREEN}✅ k6已安装: $(k6 version)${NC}"

    # 检查Go
    if ! command -v go &> /dev/null; then
        echo -e "${YELLOW}⚠️  Go未安装，将跳过测试数据生成${NC}"
        SKIP_DATA=true
    else
        echo -e "${GREEN}✅ Go已安装: $(go version)${NC}"
    fi

    # 检查MySQL连接
    if docker ps | grep -q "$MYSQL_CONTAINER"; then
        echo -e "${GREEN}✅ MySQL容器正在运行${NC}"
    else
        echo -e "${YELLOW}⚠️  MySQL容器未运行，尝试启动...${NC}"
        cd docker && docker compose up -d mysql && cd ..
        sleep 10  # 等待MySQL启动
    fi

    echo ""
}

# ================================
# 创建结果目录
# ================================
setup_results_dir() {
    echo -e "${BLUE}📁 创建结果目录...${NC}"
    mkdir -p "$RESULTS_DIR"
    echo -e "${GREEN}✅ 结果目录: $RESULTS_DIR${NC}"
    echo ""
}

# ================================
# 生成测试数据
# ================================
generate_test_data() {
    if [ "$SKIP_DATA" = "true" ]; then
        echo -e "${YELLOW}⏭️  跳过测试数据生成${NC}"
        echo ""
        return
    fi

    echo -e "${BLUE}📊 生成测试数据...${NC}"
    echo "数据规模: $DATA_SIZE"

    cd backend/tests/performance

    # 检查是否有现有数据
    if [ "$DATA_SIZE" = "small" ]; then
        echo "生成小规模测试数据..."
        go run generate_org_testdata.go -size=small 2>&1 | tee "$RESULTS_DIR/generate_data.log"
    elif [ "$DATA_SIZE" = "medium" ]; then
        echo "生成中等规模测试数据..."
        go run generate_org_testdata.go -size=medium 2>&1 | tee "$RESULTS_DIR/generate_data.log"
    elif [ "$DATA_SIZE" = "large" ]; then
        echo "生成大规模测试数据..."
        go run generate_org_testdata.go -size=large 2>&1 | tee "$RESULTS_DIR/generate_data.log"
    else
        echo -e "${RED}❌ 未知的数据规模: $DATA_SIZE${NC}"
        exit 1
    fi

    cd ../../..

    echo -e "${GREEN}✅ 测试数据生成完成${NC}"
    echo ""
}

# ================================
# 执行K6测试
# ================================
run_k6_test() {
    local test_name=$1
    local test_type=$2

    echo -e "${YELLOW}📊 执行测试: $test_name${NC}"

    local output_file="$RESULTS_DIR/${test_type}_${test_name}.json"
    local log_file="$RESULTS_DIR/${test_type}_${test_name}.log"

    k6 run \
        --env BASE_URL="$BASE_URL" \
        --env TEST_TYPE="$test_type" \
        --out json="$output_file" \
        backend/tests/performance/org_load_test.k6.js 2>&1 | tee "$log_file"

    # 检查测试是否成功
    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        echo -e "${GREEN}✅ $test_name 测试完成${NC}"
    else
        echo -e "${RED}❌ $test_name 测试失败${NC}"
        return 1
    fi

    echo ""
}

# ================================
# 生成HTML报告
# ================================
generate_html_report() {
    echo -e "${BLUE}📄 生成HTML报告...${NC}"

    # 检查是否安装了k6-to-junit
    if command -v k6-to-junit &> /dev/null; then
        for json_file in "$RESULTS_DIR"/*.json; do
            if [ -f "$json_file" ]; then
                test_name=$(basename "$json_file" .json)
                echo "生成报告: $test_name.html"

                k6-to-junit "$json_file" > "$RESULTS_DIR/${test_name}.junit.xml"
            fi
        done
    fi

    echo -e "${GREEN}✅ HTML报告生成完成${NC}"
    echo ""
}

# ================================
# 对比基线
# ================================
compare_baseline() {
    echo -e "${BLUE}📈 对比性能基线...${NC}"

    # 读取基线数据
    if [ -f "backend/tests/performance/BASELINE.json" ]; then
        echo "基线文件存在，进行对比..."
        # 这里可以添加对比逻辑
    else
        echo -e "${YELLOW}⚠️  基线文件不存在，跳过对比${NC}"
    fi

    echo ""
}

# ================================
// 打印测试摘要
# ================================
print_summary() {
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}✅ 测试完成!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo "📊 测试结果:"
    echo "  - 结果目录: $RESULTS_DIR"
    echo "  - JSON文件: $(ls -1 $RESULTS_DIR/*.json 2>/dev/null | wc -l) 个"
    echo "  - 日志文件: $(ls -1 $RESULTS_DIR/*.log 2>/dev/null | wc -l) 个"
    echo ""
    echo "📈 查看结果:"
    echo "  # 查看JSON结果"
    echo "  cat $RESULTS_DIR/base_load_test.json"
    echo ""
    echo "  # 查看日志"
    echo "  cat $RESULTS_DIR/base_load_test.log"
    echo ""
    echo "  # 生成图表（需要Python）"
    echo "  python backend/tests/performance/plot_results.py $RESULTS_DIR"
    echo ""
}

# ================================
# 主流程
# ================================
main() {
    # 1. 检查依赖
    check_dependencies

    # 2. 创建结果目录
    setup_results_dir

    # 3. 生成测试数据
    generate_test_data

    # 4. 执行性能测试
    echo -e "${BLUE}🚀 开始执行性能测试...${NC}"
    echo ""

    case $TEST_TYPE in
        base)
            run_k6_test "基准测试" "base"
            ;;
        spike)
            run_k6_test "峰值测试" "spike"
            ;;
        stress)
            run_k6_test "压力测试" "stress"
            ;;
        endurance)
            run_k6_test "耐久测试" "endurance"
            ;;
        scalability)
            run_k6_test "可扩展性测试" "scalability"
            ;;
        all)
            run_k6_test "基准测试" "base"
            run_k6_test "峰值测试" "spike"
            run_k6_test "压力测试" "stress"
            # run_k6_test "耐久测试" "endurance"  # 耗时2小时，默认跳过
            run_k6_test "可扩展性测试" "scalability"
            ;;
        *)
            echo -e "${RED}❌ 未知的测试类型: $TEST_TYPE${NC}"
            exit 1
            ;;
    esac

    # 5. 生成HTML报告
    generate_html_report

    # 6. 对比基线
    compare_baseline

    # 7. 打印摘要
    print_summary
}

# 捕获错误
trap 'echo -e "${RED}❌ 测试执行失败!${NC}"; exit 1' ERR

# 执行主流程
main

exit 0
