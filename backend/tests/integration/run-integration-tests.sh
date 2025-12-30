#!/bin/bash

# 集成测试执行脚本
# 职责: 执行所有集成测试并生成覆盖率报告

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
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

# 获取脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"
cd "$BACKEND_DIR"

log_info "当前工作目录: $BACKEND_DIR"

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."

    # 检查Go
    if ! command -v go &> /dev/null; then
        log_error "Go未安装"
        exit 1
    fi

    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装（testcontainers需要）"
        exit 1
    fi

    # 检查Docker是否运行
    if ! docker info &> /dev/null; then
        log_error "Docker未运行"
        exit 1
    fi

    log_success "依赖检查通过"
}

# 设置测试环境
setup_test_env() {
    log_info "设置测试环境..."

    # 设置Go环境变量
    export GO111MODULE=on
    export CGO_ENABLED=1

    # 设置测试环境变量
    export TEST_ENV=integration
    export LOG_LEVEL=debug

    log_success "测试环境设置完成"
}

# 运行集成测试
run_integration_tests() {
    log_info "开始运行集成测试..."

    local test_pattern="${1:-./tests/integration/...}"
    local verbose="${2:-false}"

    # 构建测试命令
    local go_test_cmd="go test -v"

    if [ "$verbose" == "true" ]; then
        go_test_cmd="$go_test_cmd -v"
    else
        go_test_cmd="$go_test_cmd"
    fi

    # 启用竞态检测
    go_test_cmd="$go_test_cmd -race"

    # 启用覆盖率
    go_test_cmd="$go_test_cmd -coverprofile=coverage.out -covermode=atomic"

    # 设置超时
    go_test_cmd="$go_test_cmd -timeout 10m"

    # 运行测试
    log_info "执行命令: $go_test_cmd $test_pattern"

    if eval "$go_test_cmd $test_pattern"; then
        log_success "集成测试全部通过"
        return 0
    else
        log_error "集成测试失败"
        return 1
    fi
}

# 生成覆盖率报告
generate_coverage_report() {
    log_info "生成覆盖率报告..."

    # 检查覆盖率文件是否存在
    if [ ! -f coverage.out ]; then
        log_warning "未找到覆盖率文件 coverage.out"
        return 1
    fi

    # 生成HTML覆盖率报告
    log_info "生成HTML覆盖率报告..."
    go tool cover -html=coverage.out -o coverage.html

    # 生成函数覆盖率报告
    log_info "生成函数覆盖率报告..."
    go tool cover -func=coverage.out -o coverage.txt

    # 计算总体覆盖率
    local coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')

    log_success "覆盖率报告生成完成"
    log_info "总体覆盖率: $coverage"
    log_info "HTML报告: $BACKEND_DIR/coverage.html"
    log_info "文本报告: $BACKEND_DIR/coverage.txt"

    # 检查覆盖率是否达标
    local coverage_percent=$(echo $coverage | sed 's/%//')
    local required_coverage=80

    if (( $(echo "$coverage_percent < $required_coverage" | bc -l) )); then
        log_warning "覆盖率 ${coverage}% 低于要求 ${required_coverage}%"
        return 1
    else
        log_success "覆盖率 ${coverage}% 达标"
        return 0
    fi
}

# 清理测试资源
cleanup_test_resources() {
    log_info "清理测试资源..."

    # 停止并删除testcontainers创建的容器
    docker ps -a --filter "label=org.testcontainers" --format "{{.ID}}" | while read -r container_id; do
        if [ -n "$container_id" ]; then
            log_info "删除测试容器: $container_id"
            docker rm -f "$container_id" &> /dev/null || true
        fi
    done

    log_success "测试资源清理完成"
}

# 运行特定测试套件
run_test_suite() {
    local suite_name=$1

    log_info "运行测试套件: $suite_name"

    case "$suite_name" in
        "tenant-permission")
            run_integration_tests "./tests/integration/ -run TestTenantPermissionIntegrationSuite"
            ;;
        "saga-business")
            run_integration_tests "./tests/integration/ -run TestSagaBusinessIntegrationSuite"
            ;;
        "humaninloop-bot")
            run_integration_tests "./tests/integration/ -run TestHumanInLoopBotIntegrationSuite"
            ;;
        "memory-conversation")
            run_integration_tests "./tests/integration/ -run TestMemoryConversationIntegrationSuite"
            ;;
        *)
            log_error "未知的测试套件: $suite_name"
            log_info "可用的测试套件:"
            log_info "  - tenant-permission"
            log_info "  - saga-business"
            log_info "  - humaninloop-bot"
            log_info "  - memory-conversation"
            exit 1
            ;;
    esac
}

# 打印使用说明
print_usage() {
    cat << EOF
集成测试执行脚本

用法: $0 [选项] [测试套件]

选项:
    -h, --help          显示此帮助信息
    -c, --cleanup-only  仅清理测试资源
    -v, --verbose       详细输出
    -s, --suite <name>  运行特定测试套件
    --no-coverage       跳过覆盖率报告生成

测试套件:
    tenant-permission      租户+权限系统集成测试
    saga-business          Saga+业务系统集成测试
    humaninloop-bot        人机协同+Bot系统集成测试
    memory-conversation    记忆+对话系统集成测试

示例:
    # 运行所有集成测试
    $0

    # 运行特定测试套件
    $0 -s tenant-permission

    # 详细输出模式
    $0 -v

    # 仅清理测试资源
    $0 -c

EOF
}

# 主函数
main() {
    local cleanup_only=false
    local verbose=false
    local suite_name=""
    local generate_coverage=true

    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                print_usage
                exit 0
                ;;
            -c|--cleanup-only)
                cleanup_only=true
                shift
                ;;
            -v|--verbose)
                verbose=true
                shift
                ;;
            -s|--suite)
                suite_name="$2"
                shift 2
                ;;
            --no-coverage)
                generate_coverage=false
                shift
                ;;
            *)
                log_error "未知参数: $1"
                print_usage
                exit 1
                ;;
        esac
    done

    # 设置退出时清理
    trap cleanup_test_resources EXIT

    # 仅清理模式
    if [ "$cleanup_only" == "true" ]; then
        cleanup_test_resources
        exit 0
    fi

    # 执行测试流程
    log_info "==================== 集成测试开始 ===================="
    echo ""

    check_dependencies
    setup_test_env

    if [ -n "$suite_name" ]; then
        run_test_suite "$suite_name"
    else
        run_integration_tests "" "$verbose"
    fi

    local test_result=$?

    echo ""

    if [ "$generate_coverage" == "true" ]; then
        generate_coverage_report
        local coverage_result=$?
    fi

    echo ""
    log_info "==================== 集成测试结束 ===================="

    # 返回测试结果
    if [ "$test_result" -ne 0 ]; then
        exit $test_result
    fi

    if [ "$generate_coverage" == "true" ] && [ "$coverage_result" -ne 0 ]; then
        exit $coverage_result
    fi

    exit 0
}

# 执行主函数
main "$@"
