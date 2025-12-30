#!/bin/bash

###############################################################################
# Token计量和预算管理系统 - 性能基准测试执行脚本
#
# 功能：
#   1. 执行所有billing性能基准测试
#   2. 生成CPU和Memory profile
#   3. 运行压力测试和内存泄漏检测
#   4. 生成测试报告
#
# 用法：
#   ./run_benchmark.sh              # 运行所有基准测试
#   ./run_benchmark.sh quick        # 快速测试（跳过压力测试）
#   ./run_benchmark.sh full         # 完整测试（包含压力测试）
#   ./run_benchmark.sh profile      # 只生成profile
#
# @author 研发B
# @date 2025-12-30
###############################################################################

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
BENCHMARK_DIR="$PROJECT_ROOT/backend/tests/benchmark/billing"
RESULTS_DIR="$PROJECT_ROOT/backend/tests/results/benchmark"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# 创建结果目录
mkdir -p "$RESULTS_DIR"

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

print_header() {
    echo ""
    echo "========================================"
    echo "$1"
    echo "========================================"
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."

    if ! command -v go &> /dev/null; then
        log_error "Go未安装"
        exit 1
    fi

    if ! command -v pprof &> /dev/null; then
        log_warning "pprof未安装，将无法分析profile"
        log_info "安装: go install github.com/google/pprof@latest"
    fi

    log_success "依赖检查完成"
}

# 执行基准测试
run_benchmarks() {
    print_header "执行性能基准测试"

    cd "$BENCHMARK_DIR"

    log_info "运行基准测试..."
    go test -bench=. -benchmem -run=^$ \
        -benchtime=1s \
        -timeout=30m \
        . > "$RESULTS_DIR/benchmark_${TIMESTAMP}.txt" 2>&1

    if [ $? -eq 0 ]; then
        log_success "基准测试完成"
        log_info "结果文件: $RESULTS_DIR/benchmark_${TIMESTAMP}.txt"
    else
        log_error "基准测试失败"
        exit 1
    fi
}

# 生成CPU Profile
generate_cpu_profile() {
    print_header "生成CPU Profile"

    cd "$BENCHMARK_DIR"

    log_info "生成CPU profile..."
    go test -cpuprofile="$RESULTS_DIR/cpu_${TIMESTAMP}.prof" \
        -bench=BenchmarkCPUProfile \
        -benchtime=10s \
        -timeout=10m \
        . > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        log_success "CPU profile生成完成"
        log_info "Profile文件: $RESULTS_DIR/cpu_${TIMESTAMP}.prof"

        # 如果安装了pprof，自动分析
        if command -v pprof &> /dev/null; then
            log_info "分析CPU profile..."
            pprof -top -nodecount=10 "$RESULTS_DIR/cpu_${TIMESTAMP}.prof" > "$RESULTS_DIR/cpu_top10_${TIMESTAMP}.txt"
            log_success "CPU profile分析完成"
        fi
    else
        log_error "CPU profile生成失败"
    fi
}

# 生成Memory Profile
generate_memory_profile() {
    print_header "生成Memory Profile"

    cd "$BENCHMARK_DIR"

    log_info "生成Memory profile..."
    go test -memprofile="$RESULTS_DIR/memory_${TIMESTAMP}.prof" \
        -bench=BenchmarkMemoryProfile \
        -benchtime=10s \
        -timeout=10m \
        . > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        log_success "Memory profile生成完成"
        log_info "Profile文件: $RESULTS_DIR/memory_${TIMESTAMP}.prof"

        # 如果安装了pprof，自动分析
        if command -v pprof &> /dev/null; then
            log_info "分析Memory profile..."
            pprof -top -nodecount=10 "$RESULTS_DIR/memory_${TIMESTAMP}.prof" > "$RESULTS_DIR/memory_top10_${TIMESTAMP}.txt"
            log_success "Memory profile分析完成"
        fi
    else
        log_error "Memory profile生成失败"
    fi
}

# 运行内存泄漏检测
run_memory_leak_test() {
    print_header "运行内存泄漏检测"

    cd "$BENCHMARK_DIR"

    log_info "运行内存泄漏检测..."
    go test -v -run=TestMemoryLeak \
        -timeout=10m \
        . > "$RESULTS_DIR/memory_leak_${TIMESTAMP}.txt" 2>&1

    if [ $? -eq 0 ]; then
        log_success "内存泄漏检测完成"
        log_info "结果文件: $RESULTS_DIR/memory_leak_${TIMESTAMP}.txt"
    else
        log_error "内存泄漏检测失败"
        exit 1
    fi
}

# 运行压力测试
run_stress_tests() {
    print_header "运行压力测试"

    cd "$BENCHMARK_DIR"

    log_warning "压力测试需要较长时间，请耐心等待..."
    log_info "运行持续高负载测试..."
    go test -v -run=TestStressTokenRecording \
        -timeout=2h \
        . > "$RESULTS_DIR/stress_test_${TIMESTAMP}.txt" 2>&1

    if [ $? -eq 0 ]; then
        log_success "持续高负载测试完成"
        log_info "结果文件: $RESULTS_DIR/stress_test_${TIMESTAMP}.txt"
    else
        log_error "持续高负载测试失败"
    fi

    log_info "运行峰值负载测试..."
    go test -v -run=TestPeakLoad \
        -timeout=30m \
        . > "$RESULTS_DIR/peak_load_${TIMESTAMP}.txt" 2>&1

    if [ $? -eq 0 ]; then
        log_success "峰值负载测试完成"
        log_info "结果文件: $RESULTS_DIR/peak_load_${TIMESTAMP}.txt"
    else
        log_error "峰值负载测试失败"
    fi
}

# 生成测试报告摘要
generate_summary() {
    print_header "生成测试报告摘要"

    local summary_file="$RESULTS_DIR/summary_${TIMESTAMP}.md"

    cat > "$summary_file" << EOF
# Token计量和预算管理系统 - 性能测试报告

**测试时间**: $(date +"%Y-%m-%d %H:%M:%S")
**测试环境**: $(go version)
**测试目录**: $BENCHMARK_DIR

---

## 📊 测试执行摘要

### 基准测试结果

EOF

    # 提取基准测试结果
    if [ -f "$RESULTS_DIR/benchmark_${TIMESTAMP}.txt" ]; then
        echo '```' >> "$summary_file"
        grep -E "^(Benchmark|PASS|FAIL)" "$RESULTS_DIR/benchmark_${TIMESTAMP}.txt" | head -100 >> "$summary_file"
        echo '```' >> "$summary_file"
    fi

    cat >> "$summary_file" << EOF

### Profile分析结果

#### CPU Profile Top 10

EOF

    if [ -f "$RESULTS_DIR/cpu_top10_${TIMESTAMP}.txt" ]; then
        echo '```' >> "$summary_file"
        cat "$RESULTS_DIR/cpu_top10_${TIMESTAMP}.txt" >> "$summary_file"
        echo '```' >> "$summary_file"
    fi

    cat >> "$summary_file" << EOF

#### Memory Profile Top 10

EOF

    if [ -f "$RESULTS_DIR/memory_top10_${TIMESTAMP}.txt" ]; then
        echo '```' >> "$summary_file"
        cat "$RESULTS_DIR/memory_top10_${TIMESTAMP}.txt" >> "$summary_file"
        echo '```' >> "$summary_file"
    fi

    cat >> "$summary_file" << EOF

### 内存泄漏检测结果

EOF

    if [ -f "$RESULTS_DIR/memory_leak_${TIMESTAMP}.txt" ]; then
        echo '```' >> "$summary_file"
        tail -50 "$RESULTS_DIR/memory_leak_${TIMESTAMP}.txt" >> "$summary_file"
        echo '```' >> "$summary_file"
    fi

    cat >> "$summary_file" << EOF

---

## 📁 详细结果文件

- **基准测试**: \`benchmark_${TIMESTAMP}.txt\`
- **CPU Profile**: \`cpu_${TIMESTAMP}.prof\`
- **Memory Profile**: \`memory_${TIMESTAMP}.prof\`
- **内存泄漏检测**: \`memory_leak_${TIMESTAMP}.txt\`
EOF

    if [ -f "$RESULTS_DIR/stress_test_${TIMESTAMP}.txt" ]; then
        echo "- **持续高负载测试**: \`stress_test_${TIMESTAMP}.txt\`" >> "$summary_file"
    fi

    if [ -f "$RESULTS_DIR/peak_load_${TIMESTAMP}.txt" ]; then
        echo "- **峰值负载测试**: \`peak_load_${TIMESTAMP}.txt\`" >> "$summary_file"
    fi

    cat >> "$summary_file" << EOF

---

## 🔍 Profile分析命令

### CPU Profile分析
\`\`\`bash
# 查看Top 10
pprof -top -nodecount=10 $RESULTS_DIR/cpu_${TIMESTAMP}.prof

# 查看函数调用图
pprof -web $RESULTS_DIR/cpu_${TIMESTAMP}.prof

# 查看特定函数
pprof -list BenchmarkCPUProfile $RESULTS_DIR/cpu_${TIMESTAMP}.prof
\`\`\`

### Memory Profile分析
\`\`\`bash
# 查看Top 10
pprof -top -nodecount=10 $RESULTS_DIR/memory_${TIMESTAMP}.prof

# 查看内存分配
pprof -web $RESULTS_DIR/memory_${TIMESTAMP}.prof

# 查看特定函数
pprof -list BenchmarkMemoryProfile $RESULTS_DIR/memory_${TIMESTAMP}.prof
\`\`\`

---

**生成时间**: $(date +"%Y-%m-%d %H:%M:%S")
EOF

    log_success "测试报告摘要生成完成"
    log_info "报告文件: $summary_file"
}

# 清理旧结果
cleanup_old_results() {
    log_info "清理7天前的旧测试结果..."
    find "$RESULTS_DIR" -name "*.txt" -mtime +7 -delete 2>/dev/null || true
    find "$RESULTS_DIR" -name "*.prof" -mtime +7 -delete 2>/dev/null || true
    log_success "旧结果清理完成"
}

# 主函数
main() {
    print_header "Token计量和预算管理系统 - 性能基准测试"

    check_dependencies

    # 根据参数选择测试模式
    case "${1:-all}" in
        quick)
            log_info "快速测试模式（跳过压力测试）"
            run_benchmarks
            generate_cpu_profile
            generate_memory_profile
            run_memory_leak_test
            ;;

        full)
            log_info "完整测试模式（包含压力测试）"
            run_benchmarks
            generate_cpu_profile
            generate_memory_profile
            run_memory_leak_test
            run_stress_tests
            ;;

        profile)
            log_info "Profile生成模式"
            generate_cpu_profile
            generate_memory_profile
            ;;

        stress)
            log_info "压力测试模式"
            run_stress_tests
            ;;

        all|*)
            log_info "完整测试模式"
            run_benchmarks
            generate_cpu_profile
            generate_memory_profile
            run_memory_leak_test
            # run_stress_tests  # 默认不运行压力测试，因为耗时较长
            ;;
    esac

    generate_summary
    cleanup_old_results

    print_header "测试完成"
    log_success "所有测试已完成！"
    log_info "结果目录: $RESULTS_DIR"
}

# 执行主函数
main "$@"
