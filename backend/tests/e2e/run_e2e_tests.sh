#!/bin/bash

# ==============================================================================
# E2E测试执行脚本
#
# 职责: 提供统一的E2E测试执行入口
# 遵循单一职责原则：只负责E2E测试的执行和管理
# ==============================================================================

set -e  # 遇到错误立即退出
set -u  # 使用未定义变量时报错

# ==================== 颜色定义 ====================

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ==================== 变量定义 ====================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"
E2E_DIR="${BACKEND_DIR}/tests/e2e"
COVERAGE_FILE="${E2E_DIR}/coverage.out"
COVERAGE_HTML="${E2E_DIR}/coverage.html"
TEST_RESULTS_FILE="${E2E_DIR}/test_results.txt"

# ==================== 日志函数 ====================

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# ==================== 帮助信息 ====================

show_help() {
    cat << EOF
用法: $0 [选项] [Go测试参数]

E2E测试执行脚本 - Coze Studio端到端测试

选项:
    -h, --help          显示帮助信息
    -c, --cover         生成测试覆盖率报告
    -v, --verbose       详细输出（传递给go test）
    -r, --run TEST      运行特定测试（例如: TestE2E_TenantRegistrationJourney）
    -t, --timeout DURATION 设置超时时间（默认: 30m）
    -p, --parallel N     并发执行数（默认: 1）
    --clean             清理测试覆盖文件
    --report            生成HTML覆盖率报告

示例:
    # 执行所有E2E测试
    $0

    # 执行所有E2E测试并生成覆盖率
    $0 --cover

    # 执行特定测试
    $0 --run TestE2E_TenantRegistrationJourney

    # 并发执行（4个并发）
    $0 --parallel 4

    # 生成覆盖率报告
    $0 --cover --report

环境变量:
    GO_TEST_FLAGS       额外的go test参数
    TEST_TIMEOUT        测试超时时间（默认: 30m）

EOF
}

# ==================== 参数解析 ====================

VERBOSE=false
COVER=false
RUN_TEST=""
TIMEOUT="30m"
PARALLEL=1
CLEAN=false
GENERATE_REPORT=false

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--cover)
            COVER=true
            shift
            ;;
        -r|--run)
            RUN_TEST="$2"
            shift 2
            ;;
        -t|--timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        -p|--parallel)
            PARALLEL="$2"
            shift 2
            ;;
        --clean)
            CLEAN=true
            shift
            ;;
        --report)
            GENERATE_REPORT=true
            shift
            ;;
        *)
            # 其他参数传递给go test
            break
            ;;
    esac
done

# ==================== 环境检查 ====================

check_environment() {
    log_info "检查测试环境..."

    # 检查Go是否安装
    if ! command -v go &> /dev/null; then
        log_error "Go未安装，请先安装Go"
        exit 1
    fi

    # 检查Docker是否运行
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，testcontainers需要Docker"
        exit 1
    fi

    if ! docker info &> /dev/null; then
        log_error "Docker未运行，请启动Docker"
        exit 1
    fi

    # 检查E2E目录是否存在
    if [[ ! -d "$E2E_DIR" ]]; then
        log_error "E2E测试目录不存在: $E2E_DIR"
        exit 1
    fi

    log_success "环境检查通过"
}

# ==================== 清理函数 ====================

clean_coverage() {
    log_info "清理测试覆盖文件..."
    rm -f "$COVERAGE_FILE" "$COVERAGE_HTML" "$TEST_RESULTS_FILE"
    log_success "清理完成"
}

# ==================== 执行测试 ====================

run_tests() {
    log_info "开始执行E2E测试..."
    log_info "测试目录: $E2E_DIR"
    log_info "超时时间: $TIMEOUT"
    log_info "并发数: $PARALLEL"

    # 构建go test参数
    TEST_ARGS="-v"
    TEST_ARGS="$TEST_ARGS -timeout=$TIMEOUT"
    TEST_ARGS="$TEST_ARGS -parallel=$PARALLEL"

    # 添加覆盖率
    if [[ "$COVER" == true ]]; then
        TEST_ARGS="$TEST_ARGS -coverprofile=$COVERAGE_FILE -covermode=atomic"
    fi

    # 添加运行特定测试
    if [[ -n "$RUN_TEST" ]]; then
        TEST_ARGS="$TEST_ARGS -run=$RUN_TEST"
        log_info "运行特定测试: $RUN_TEST"
    fi

    # 传递额外的参数
    TEST_ARGS="$TEST_ARGS $@"

    log_info "go test参数: $TEST_ARGS"

    # 切换到E2E目录
    cd "$E2E_DIR"

    # 执行测试
    log_info "=========================================="
    if go test $TEST_ARGS 2>&1 | tee "$TEST_RESULTS_FILE"; then
        log_success "=========================================="
        log_success "E2E测试全部通过！"

        # 显示测试摘要
        if [[ -f "$TEST_RESULTS_FILE" ]]; then
            log_info "测试摘要:"
            grep -E "^(PASS|FAIL|ok|---)" "$TEST_RESULTS_FILE" || true
        fi
    else
        log_error "=========================================="
        log_error "E2E测试失败！"
        exit 1
    fi
}

# ==================== 生成覆盖率报告 ====================

generate_coverage_report() {
    if [[ ! -f "$COVERAGE_FILE" ]]; then
        log_warning "覆盖率文件不存在: $COVERAGE_FILE"
        log_warning "请使用 --cover 参数生成覆盖率数据"
        return
    fi

    log_info "生成覆盖率报告..."

    # 生成HTML报告
    if go tool cover -html="$COVERAGE_FILE" -o "$COVERAGE_HTML"; then
        log_success "HTML覆盖率报告已生成: $COVERAGE_HTML"

        # 显示覆盖率摘要
        log_info "覆盖率摘要:"
        go tool cover -func="$COVERAGE_FILE" | tail -1

        # 尝试在浏览器中打开（可选）
        if command -v open &> /dev/null; then
            log_info "在浏览器中打开覆盖率报告..."
            open "$COVERAGE_HTML" 2>/dev/null || true
        fi
    else
        log_error "生成覆盖率报告失败"
        exit 1
    fi
}

# ==================== 显示测试结果 ====================

show_test_results() {
    if [[ ! -f "$TEST_RESULTS_FILE" ]]; then
        return
    fi

    log_info "=========================================="
    log_info "测试结果摘要:"
    log_info "=========================================="

    # 统计测试数量
    TOTAL=$(grep -c "^=== RUN" "$TEST_RESULTS_FILE" 2>/dev/null || echo "0")
    PASSED=$(grep -c "^--- PASS:" "$TEST_RESULTS_FILE" 2>/dev/null || echo "0")
    FAILED=$(grep -c "^--- FAIL:" "$TEST_RESULTS_FILE" 2>/dev/null || echo "0")
    SKIPPED=$(grep -c "^--- SKIP:" "$TEST_RESULTS_FILE" 2>/dev/null || echo "0")

    echo "总测试数: $TOTAL"
    echo "通过: $PASSED"
    echo "失败: $FAILED"
    echo "跳过: $SKIPPED"

    if [[ "$FAILED" -gt 0 ]]; then
        echo ""
        log_error "失败的测试:"
        grep "^--- FAIL:" "$TEST_RESULTS_FILE" || true
    fi

    log_info "=========================================="
}

# ==================== 主函数 ====================

main() {
    echo "=========================================="
    echo "  Coze Studio E2E测试执行器"
    echo "=========================================="
    echo ""

    # 检查环境
    check_environment

    # 清理
    if [[ "$CLEAN" == true ]]; then
        clean_coverage
        exit 0
    fi

    # 执行测试
    run_tests

    # 显示结果
    show_test_results

    # 生成覆盖率报告
    if [[ "$GENERATE_REPORT" == true ]]; then
        generate_coverage_report
    elif [[ "$COVER" == true ]]; then
        # 如果启用了覆盖，显示简要信息
        log_info "覆盖率数据已保存到: $COVERAGE_FILE"
        log_info "使用 --report 生成HTML报告"
    fi

    log_success "E2E测试执行完成！"
}

# 执行主函数
main "$@"
