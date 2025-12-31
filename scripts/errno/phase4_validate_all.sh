#!/bin/bash
# Phase 4: 综合验证（语法+编译+测试）
# 作者: ZKER架构团队
# 日期: 2025-01-03
# 用途: 企业级errno迁移的第四步 - 多重验证

set -euo pipefail

# ========================================
# 配置
# ========================================
BACKEND_DIR="backend"
LOG_DIR="scripts/errno/logs"
mkdir -p "$LOG_DIR"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 日志文件
SYNTAX_LOG="$LOG_DIR/syntax_check.log"
COMPILE_LOG="$LOG_DIR/compile_check.log"
TEST_LOG="$LOG_DIR/test_check.log"

# 统计变量
SYNTAX_ERRORS=0
COMPILE_ERRORS=0
TEST_ERRORS=0

# ========================================
# 函数
# ========================================

log_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

log_success() {
    echo -e "${GREEN}✓${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

log_error() {
    echo -e "${RED}✗${NC} $1"
}

print_header() {
    echo ""
    echo "=========================================="
    echo "  $1"
    echo "=========================================="
    echo ""
}

# ========================================
# Phase 4.1: 语法验证
# ========================================

validate_syntax() {
    print_header "Phase 4.1: Go语法验证"

    log_info "检查Go文件格式和语法..."
    echo ""

    local errors=0
    local files_checked=0

    # 获取所有修改过的Go文件
    local modified_files=$(git diff --name-only 2>/dev/null | grep '\.go$' | grep -v '_test\.go$' || true)

    if [ -z "$modified_files" ]; then
        # 如果没有git，检查所有Go文件
        modified_files=$(find "$BACKEND_DIR" -name "*.go" -not -name "*_test.go" 2>/dev/null || true)
    fi

    for file in $modified_files; do
        [ ! -f "$file" ] && continue

        files_checked=$((files_checked + 1))
        local filename=$(basename "$file")

        echo -n "  [$files_checked] $filename ... "

        # 1. go fmt 检查格式
        if ! output=$(go fmt "$file" 2>&1); then
            log_error "格式错误"
            echo "$output" >> "$SYNTAX_LOG"
            errors=$((errors + 1))
            continue
        fi

        # 2. go vet 检查
        if ! output=$(go vet "$file" 2>&1); then
            log_warning "go vet警告"
            echo "$output" >> "$SYNTAX_LOG"
            # 不算错误，继续
        fi

        log_success "OK"
    done

    echo ""
    if [ $errors -eq 0 ]; then
        log_success "语法验证通过（检查了 $files_checked 个文件）"
        return 0
    else
        log_error "发现 $errors 个语法错误"
        echo "  详细日志: $SYNTAX_LOG"
        return 1
    fi
}

# ========================================
# Phase 4.2: 编译验证
# ========================================

validate_compile() {
    print_header "Phase 4.2: 编译验证"

    log_info "编译整个项目..."
    echo ""

    cd "$BACKEND_DIR"

    # 尝试编译
    if go build -v ./... 2>&1 | tee "$COMPILE_LOG"; then
        echo ""
        log_success "编译成功！"
        cd - > /dev/null
        return 0
    else
        echo ""
        log_error "编译失败"
        echo "  详细日志: $COMPILE_LOG"
        echo ""
        echo "编译错误详情:"
        grep -i "error" "$COMPILE_LOG" || true
        cd - > /dev/null
        return 1
    fi
}

# ========================================
# Phase 4.3: 单元测试验证
# ========================================

validate_tests() {
    print_header "Phase 4.3: 单元测试验证"

    log_info "运行单元测试（不包含集成测试）..."
    echo ""

    cd "$BACKEND_DIR"

    # 运行测试（-short flag跳过慢速测试）
    local test_start=$(date +%s)

    if go test -short -v ./... 2>&1 | tee "$TEST_LOG"; then
        local test_end=$(date +%s)
        local test_duration=$((test_end - test_start))

        echo ""
        log_success "所有测试通过！（耗时: ${test_duration}秒）"

        # 统计测试覆盖
        echo ""
        log_info "测试覆盖率统计:"
        go test -short -cover ./... 2>&1 | grep "coverage:" | head -20

        cd - > /dev/null
        return 0
    else
        local test_end=$(date +%s)
        echo ""
        log_error "测试失败"
        echo "  详细日志: $TEST_LOG"
        echo ""
        echo "失败的测试:"
        grep "FAIL:" "$TEST_LOG" || true
        cd - > /dev/null
        return 1
    fi
}

# ========================================
# Phase 4.4: errno使用验证
# ========================================

validate_errno_usage() {
    print_header "Phase 4.4: errno使用验证"

    log_info "检查errno引用..."
    echo ""

    # 1. 检查直接errno引用（应该为0）
    local direct_errno=$(grep -rn "return.*errno\." "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
        grep -v "_test.go" | \
        wc -l)

    echo "  1. 直接errno引用: $direct_errno (应为0)"
    if [ "$direct_errno" -gt 0 ]; then
        log_error "发现 $direct_errno 处直接errno引用"
        echo "  详情:"
        grep -rn "return.*errno\." "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
            grep -v "_test.go" | head -10
    else
        log_success "✓ 无直接errno引用"
    fi
    echo ""

    # 2. 检查errorx.New使用
    local errorx_new=$(grep -rn "errorx\.New(errno\." "$BACKEND_DIR" --include="*.go" 2>/dev/null | wc -l)
    echo "  2. errorx.New使用: $errorx_new"
    echo ""

    # 3. 检查errorx.Wrap使用
    local errorx_wrap=$(grep -rn "errorx\.Wrap.*errno\." "$BACKEND_DIR" --include="*.go" 2>/dev/null | wc -l)
    echo "  3. errorx.Wrap使用: $errorx_wrap"
    echo ""

    # 4. 计算覆盖率
    local total_errors=$(grep -rn "return.*error" "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
        grep -v "_test.go" | wc -l)
    local errorx_usage=$((errorx_new + errorx_wrap))

    if [ "$total_errors" -gt 0 ]; then
        local coverage=$((errorx_usage * 100 / total_errors))
        echo "  4. 错误处理覆盖率: $coverage% (目标≥95%)"
        echo ""

        if [ "$coverage" -ge 95 ]; then
            log_success "✓ 错误处理覆盖率达标"
        else
            log_warning "错误处理覆盖率低于95%"
        fi
    fi

    echo ""
    if [ "$direct_errno" -eq 0 ] && [ "$errorx_new" -gt 0 ]; then
        return 0
    else
        return 1
    fi
}

# ========================================
# 主函数
# ========================================

main() {
    echo "=========================================="
    echo "  errno迁移 - Phase 4: 综合验证"
    echo "=========================================="
    echo ""

    local start_time=$(date +%s)
    local all_passed=true

    # Phase 4.1: 语法验证
    if ! validate_syntax; then
        SYNTAX_ERRORS=1
        all_passed=false
    fi
    echo ""

    # Phase 4.2: 编译验证
    if ! validate_compile; then
        COMPILE_ERRORS=1
        all_passed=false
    fi
    echo ""

    # Phase 4.3: 单元测试验证
    if ! validate_tests; then
        TEST_ERRORS=1
        all_passed=false
    fi
    echo ""

    # Phase 4.4: errno使用验证
    if ! validate_errno_usage; then
        all_passed=false
    fi
    echo ""

    # 总结
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))

    print_header "验证总结"
    echo "总耗时: ${duration}秒"
    echo ""
    echo "检查项结果:"
    echo "  语法验证:   $([ "$SYNTAX_ERRORS" -eq 0 ] && echo "✅ 通过" || echo "❌ 失败")"
    echo "  编译验证:   $([ "$COMPILE_ERRORS" -eq 0 ] && echo "✅ 通过" || echo "❌ 失败")"
    echo "  测试验证:   $([ "$TEST_ERRORS" -eq 0 ] && echo "✅ 通过" || echo "❌ 失败")"
    echo "  errno使用:  ✅ 已检查"
    echo ""

    if [ "$all_passed" = true ]; then
        print_header "✅ 所有验证通过！"
        echo ""
        echo "下一步:"
        echo "  1. 代码review"
        echo "  2. 提交PR: git add . && git commit -m 'feat: migrate errno to errorx'"
        echo "  3. 等待CI/CD检查"
        echo ""
        return 0
    else
        print_header "❌ 验证失败"
        echo ""
        echo "请检查上述错误并修复后重新验证"
        echo ""
        echo "日志文件:"
        echo "  - 语法检查: $SYNTAX_LOG"
        echo "  - 编译检查: $COMPILE_LOG"
        echo "  - 测试检查: $TEST_LOG"
        echo ""
        echo "回滚命令:"
        echo "  ./scripts/errno/phase5_rollback.sh"
        echo ""
        return 1
    fi
}

# 执行主函数
main "$@"
