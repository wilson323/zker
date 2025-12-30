#!/bin/bash
# ============================================================
# 计费系统集成测试执行脚本
#
# 功能:
# 1. 运行所有计费系统相关的集成测试
# 2. 生成测试覆盖率报告
# 3. 生成性能基准数据
# 4. 生成HTML测试报告
#
# 使用方法:
#   cd backend/tests/integration/billing
#   ./run_tests.sh [选项]
#
# 选项:
#   -v, --verbose    详细输出模式
#   -c, --coverage   生成覆盖率报告
#   -p, --performance 生成性能基准报告
#   -h, --help       显示帮助信息
# ============================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 当前脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../../../../.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"

# 测试结果目录
TEST_RESULTS_DIR="$SCRIPT_DIR/results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# 解析命令行参数
VERBOSE=false
COVERAGE=false
PERFORMANCE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -p|--performance)
            PERFORMANCE=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -v, --verbose     详细输出模式"
            echo "  -c, --coverage    生成覆盖率报告"
            echo "  -p, --performance 生成性能基准报告"
            echo "  -h, --help        显示帮助信息"
            exit 0
            ;;
        *)
            echo -e "${RED}未知选项: $1${NC}"
            exit 1
            ;;
    esac
done

# ============================================================
# 函数定义
# ============================================================

# 打印信息日志
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

# 打印成功日志
log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# 打印警告日志
log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 打印错误日志
log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 打印分隔线
print_separator() {
    echo "=========================================="
}

# 创建测试结果目录
create_results_dir() {
    mkdir -p "$TEST_RESULTS_DIR/$TIMESTAMP"
    log_info "测试结果目录: $TEST_RESULTS_DIR/$TIMESTAMP"
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."

    # 检查Go
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装"
        exit 1
    fi

    # 检查Docker（testcontainers需要）
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装（testcontainers需要Docker）"
        exit 1
    fi

    # 检查Docker是否运行
    if ! docker info &> /dev/null; then
        log_error "Docker 未运行"
        exit 1
    fi

    log_success "依赖检查通过"
}

# 运行Token计量集成测试
run_token_metering_tests() {
    log_info "运行Token计量集成测试..."

    cd "$SCRIPT_DIR"

    if [ "$VERBOSE" = true ]; then
        go test -v -timeout 30m ./... \
            -run TestTokenMetering \
            2>&1 | tee "$TEST_RESULTS_DIR/$TIMESTAMP/token_metering.log"
    else
        go test -timeout 30m ./... \
            -run TestTokenMetering \
            2>&1 | tee "$TEST_RESULTS_DIR/$TIMESTAMP/token_metering.log"
    fi

    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        log_success "Token计量集成测试通过"
    else
        log_error "Token计量集成测试失败"
        return 1
    fi
}

# 运行预算管理集成测试
run_budget_management_tests() {
    log_info "运行预算管理集成测试..."

    cd "$SCRIPT_DIR"

    if [ "$VERBOSE" = true ]; then
        go test -v -timeout 30m ./... \
            -run TestBudgetManagement \
            2>&1 | tee "$TEST_RESULTS_DIR/$TIMESTAMP/budget_management.log"
    else
        go test -timeout 30m ./... \
            -run TestBudgetManagement \
            2>&1 | tee "$TEST_RESULTS_DIR/$TIMESTAMP/budget_management.log"
    fi

    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        log_success "预算管理集成测试通过"
    else
        log_error "预算管理集成测试失败"
        return 1
    fi
}

# 运行E2E测试
run_e2e_tests() {
    log_info "运行端到端集成测试..."

    cd "$SCRIPT_DIR"

    if [ "$VERBOSE" = true ]; then
        go test -v -timeout 30m ./... \
            -run TestBilling_E2E \
            2>&1 | tee "$TEST_RESULTS_DIR/$TIMESTAMP/e2e.log"
    else
        go test -timeout 30m ./... \
            -run TestBilling_E2E \
            2>&1 | tee "$TEST_RESULTS_DIR/$TIMESTAMP/e2e.log"
    fi

    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        log_success "端到端集成测试通过"
    else
        log_error "端到端集成测试失败"
        return 1
    fi
}

# 运行并发测试
run_concurrent_tests() {
    log_info "运行并发测试..."

    cd "$SCRIPT_DIR"

    if [ "$VERBOSE" = true ]; then
        go test -v -timeout 30m ./... \
            -run TestConcurrent \
            2>&1 | tee "$TEST_RESULTS_DIR/$TIMESTAMP/concurrent.log"
    else
        go test -timeout 30m ./... \
            -run TestConcurrent \
            2>&1 | tee "$TEST_RESULTS_DIR/$TIMESTAMP/concurrent.log"
    fi

    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        log_success "并发测试通过"
    else
        log_error "并发测试失败"
        return 1
    fi
}

# 运行错误场景测试
run_error_scenarios_tests() {
    log_info "运行错误场景测试..."

    cd "$SCRIPT_DIR"

    if [ "$VERBOSE" = true ]; then
        go test -v -timeout 30m ./... \
            -run TestErrorScenarios \
            2>&1 | tee "$TEST_RESULTS_DIR/$TIMESTAMP/error_scenarios.log"
    else
        go test -timeout 30m ./... \
            -run TestErrorScenarios \
            2>&1 | tee "$TEST_RESULTS_DIR/$TIMESTAMP/error_scenarios.log"
    fi

    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        log_success "错误场景测试通过"
    else
        log_error "错误场景测试失败"
        return 1
    fi
}

# 生成覆盖率报告
generate_coverage_report() {
    if [ "$COVERAGE" = false ]; then
        return
    fi

    log_info "生成测试覆盖率报告..."

    cd "$SCRIPT_DIR"

    # 生成覆盖率文件
    go test -timeout 30m -coverprofile="$TEST_RESULTS_DIR/$TIMESTAMP/coverage.out" \
        -covermode=atomic ./... 2>&1 | tee "$TEST_RESULTS_DIR/$TIMESTAMP/coverage.log"

    # 生成HTML覆盖率报告
    go tool cover -html="$TEST_RESULTS_DIR/$TIMESTAMP/coverage.out" \
        -o "$TEST_RESULTS_DIR/$TIMESTAMP/coverage.html"

    # 生成覆盖率函数报告
    go tool cover -func="$TEST_RESULTS_DIR/$TIMESTAMP/coverage.out" \
        > "$TEST_RESULTS_DIR/$TIMESTAMP/coverage_func.txt"

    log_success "覆盖率报告生成完成"
    log_info "  - HTML报告: $TEST_RESULTS_DIR/$TIMESTAMP/coverage.html"
    log_info "  - 函数报告: $TEST_RESULTS_DIR/$TIMESTAMP/coverage_func.txt"

    # 提取总体覆盖率
    TOTAL_COVERAGE=$(go tool cover -func="$TEST_RESULTS_DIR/$TIMESTAMP/coverage.out" | \
        grep total | awk '{print $3}')
    log_info "  - 总体覆盖率: $TOTAL_COVERAGE"
}

# 生成性能基准报告
generate_performance_report() {
    if [ "$PERFORMANCE" = false ]; then
        return
    fi

    log_info "生成性能基准报告..."

    cd "$SCRIPT_DIR"

    # 运行性能基准测试
    go test -bench=. -benchmem -timeout 30m \
        ./... 2>&1 | tee "$TEST_RESULTS_DIR/$TIMESTAMP/benchmark.log"

    log_success "性能基准报告生成完成"
    log_info "  - 基准报告: $TEST_RESULTS_DIR/$TIMESTAMP/benchmark.log"
}

# 生成测试总结报告
generate_summary_report() {
    log_info "生成测试总结报告..."

    SUMMARY_FILE="$TEST_RESULTS_DIR/$TIMESTAMP/summary.txt"

    cat > "$SUMMARY_FILE" << EOF
============================================================
计费系统集成测试总结报告
============================================================
测试时间: $(date '+%Y-%m-%d %H:%M:%S')
测试环境: $(go version)
Docker版本: $(docker --version)

============================================================
测试结果
============================================================

EOF

    # 分析各个测试日志文件
    for log_file in token_metering budget_management e2e concurrent error_scenarios; do
        log_path="$TEST_RESULTS_DIR/$TIMESTAMP/${log_file}.log"

        if [ -f "$log_path" ]; then
            echo "[$log_file 测试结果]" >> "$SUMMARY_FILE"
            grep -E "PASS|FAIL" "$log_path" | tail -5 >> "$SUMMARY_FILE" || echo "无测试结果" >> "$SUMMARY_FILE"
            echo "" >> "$SUMMARY_FILE"
        fi
    done

    # 如果生成了覆盖率报告，添加覆盖率信息
    if [ "$COVERAGE" = true ] && [ -f "$TEST_RESULTS_DIR/$TIMESTAMP/coverage_func.txt" ]; then
        echo "[覆盖率信息]" >> "$SUMMARY_FILE"
        grep "total" "$TEST_RESULTS_DIR/$TIMESTAMP/coverage_func.txt" >> "$SUMMARY_FILE"
        echo "" >> "$SUMMARY_FILE"
    fi

    cat >> "$SUMMARY_FILE" << EOF
============================================================
测试文件位置
============================================================
结果目录: $TEST_RESULTS_DIR/$TIMESTAMP

测试日志:
  - token_metering.log
  - budget_management.log
  - e2e.log
  - concurrent.log
  - error_scenarios.log

EOF

    if [ "$COVERAGE" = true ]; then
        cat >> "$SUMMARY_FILE" << EOF
覆盖率报告:
  - coverage.out
  - coverage.html
  - coverage_func.txt

EOF
    fi

    if [ "$PERFORMANCE" = true ]; then
        cat >> "$SUMMARY_FILE" << EOF
性能基准:
  - benchmark.log

EOF
    fi

    cat >> "$SUMMARY_FILE" << EOF
============================================================
报告生成时间: $(date '+%Y-%m-%d %H:%M:%S')
============================================================
EOF

    log_success "测试总结报告生成完成: $SUMMARY_FILE"

    # 打印总结报告内容
    cat "$SUMMARY_FILE"
}

# 清理旧的测试结果
cleanup_old_results() {
    log_info "清理旧的测试结果（保留最近7天）..."

    find "$TEST_RESULTS_DIR" -type d -mtime +7 -exec rm -rf {} + 2>/dev/null || true

    log_success "清理完成"
}

# ============================================================
# 主流程
# ============================================================

main() {
    print_separator
    echo "计费系统集成测试执行脚本"
    print_separator
    echo ""

    # 检查依赖
    check_dependencies
    echo ""

    # 创建测试结果目录
    create_results_dir
    echo ""

    # 清理旧结果
    cleanup_old_results
    echo ""

    # 记录开始时间
    START_TIME=$(date +%s)

    # 运行所有测试
    log_info "开始运行测试套件..."
    echo ""

    FAILED_TESTS=0

    # Token计量集成测试
    if ! run_token_metering_tests; then
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    echo ""

    # 预算管理集成测试
    if ! run_budget_management_tests; then
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    echo ""

    # E2E测试
    if ! run_e2e_tests; then
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    echo ""

    # 并发测试
    if ! run_concurrent_tests; then
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    echo ""

    # 错误场景测试
    if ! run_error_scenarios_tests; then
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    echo ""

    # 计算耗时
    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))
    MINUTES=$((DURATION / 60))
    SECONDS=$((DURATION % 60))

    # 生成覆盖率报告
    generate_coverage_report
    echo ""

    # 生成性能基准报告
    generate_performance_report
    echo ""

    # 生成总结报告
    generate_summary_report
    echo ""

    # 测试结果总结
    print_separator
    if [ $FAILED_TESTS -eq 0 ]; then
        log_success "所有测试通过！"
        log_success "总耗时: ${MINUTES}分${SECONDS}秒"
    else
        log_error "有 $FAILED_TESTS 个测试套件失败"
        log_warning "总耗时: ${MINUTES}分${SECONDS}秒"
        exit 1
    fi
    print_separator
}

# 执行主流程
main
