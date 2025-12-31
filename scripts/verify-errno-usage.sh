#!/bin/bash
# errno使用规范验证脚本
# 用途: 在CI/CD中验证errno使用是否符合企业级规范
# 使用: bash scripts/verify-errno-usage.sh
# 作者: ZKER架构团队
# 日期: 2025-01-03

set -euo pipefail

# ========================================
# 配置
# ========================================
BACKEND_DIR="backend"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
RESULT_FILE="$SCRIPT_DIR/errno/errno_check_result.json"
LOG_FILE="$SCRIPT_DIR/errno/errno_check.log"

# 创建必要目录
mkdir -p "$(dirname "$RESULT_FILE")"
mkdir -p "$(dirname "$LOG_FILE")"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 统计变量
TOTAL_ISSUES=0
ERRORS=0
WARNINGS=0

# ========================================
# 函数定义
# ========================================

log_info() {
    echo -e "${BLUE}ℹ${NC} $1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [INFO] $1" >> "$LOG_FILE"
}

log_success() {
    echo -e "${GREEN}✓${NC} $1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [SUCCESS] $1" >> "$LOG_FILE"
}

log_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [WARNING] $1" >> "$LOG_FILE"
    ((WARNINGS++))
}

log_error() {
    echo -e "${RED}✗${NC} $1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [ERROR] $1" >> "$LOG_FILE"
    ((ERRORS++))
    ((TOTAL_ISSUES++))
}

print_header() {
    echo ""
    echo "=========================================="
    echo "  $1"
    echo "=========================================="
    echo ""
}

# ========================================
# 检查1: 直接errno引用检查
# ========================================
check_direct_errno_usage() {
    print_header "检查1: 直接errno引用检查"

    local direct_errno_count=0
    local violation_files=()

    log_info "扫描直接errno引用（排除测试文件）..."

    # 查找所有非测试Go文件中的直接errno引用
    while IFS= read -r file; do
        if [ ! -f "$file" ]; then
            continue
        fi

        # 检查文件中是否有直接errno引用
        local violations=$(grep -n "return.*errno\." "$file" 2>/dev/null | \
            grep -v "errorx\." | \
            grep -v "//" | \
            grep -v "_" || true)

        if [ -n "$violations" ]; then
            ((direct_errno_count += $(echo "$violations" | wc -l)))
            violation_files+=("$file")
            log_error "文件 $file 存在直接errno引用:"
            echo "$violations" | while read -r line; do
                echo "    $line"
            done
        fi
    done < <(find "$BACKEND_DIR" -name "*.go" -not -name "*_test.go" 2>/dev/null)

    echo ""
    if [ $direct_errno_count -eq 0 ]; then
        log_success "未发现直接errno引用（0处违规）"
        return 0
    else
        log_error "发现 $direct_errno_count 处直接errno引用"
        echo "  违规文件数: ${#violation_files[@]}"
        return 1
    fi
}

# ========================================
# 检查2: errorx.New/Wrap使用检查
# ========================================
check_errorx_usage() {
    print_header "检查2: errorx.New/Wrap使用检查"

    log_info "统计errorx使用情况..."

    # 统计errorx.New使用
    local errorx_new_count=$(grep -rn "errorx\.New(errno\." "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
        grep -v "_test.go" | wc -l)

    # 统计errorx.Wrap使用
    local errorx_wrap_count=$(grep -rn "errorx\.Wrap.*errno\." "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
        grep -v "_test.go" | wc -l)

    local total_errorx=$((errorx_new_count + errorx_wrap_count))

    echo "  errorx.New使用: $errorx_new_count 处"
    echo "  errorx.Wrap使用: $errorx_wrap_count 处"
    echo "  errorx总使用: $total_errorx 处"
    echo ""

    if [ $total_errorx -gt 0 ]; then
        log_success "发现errorx使用（符合规范）"
        return 0
    else
        log_warning "未发现errorx使用，可能未完成errno迁移"
        return 1
    fi
}

# ========================================
# 检查3: 错误处理覆盖率检查
# ========================================
check_error_coverage() {
    print_header "检查3: 错误处理覆盖率检查"

    log_info "分析错误处理覆盖率..."

    # 统计return语句总数
    local total_returns=$(grep -rn "return" "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
        grep -v "_test.go" | \
        grep -v "//" | wc -l)

    # 统计包含错误处理的return语句
    local error_returns=$(grep -rn "return.*error" "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
        grep -v "_test.go" | \
        grep -v "//" | wc -l)

    # 统计使用errorx的return语句
    local errorx_returns=$(grep -rn "return.*errorx" "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
        grep -v "_test.go" | wc -l)

    echo "  return语句总数: $total_returns"
    echo "  包含错误的return: $error_returns"
    echo "  使用errorx的return: $errorx_returns"
    echo ""

    if [ $error_returns -gt 0 ]; then
        local coverage=$((errorx_returns * 100 / error_returns))
        echo "  错误处理覆盖率: $coverage%"

        if [ $coverage -ge 95 ]; then
            log_success "错误处理覆盖率达标（≥95%）"
            return 0
        elif [ $coverage -ge 80 ]; then
            log_warning "错误处理覆盖率较低（$coverage% < 95%）"
            return 0
        else
            log_error "错误处理覆盖率严重不足（$coverage% < 80%）"
            return 1
        fi
    else
        log_warning "未发现错误处理语句"
        return 1
    fi
}

# ========================================
# 检查4: 硬编码错误检查
# ========================================
check_hardcoded_errors() {
    print_header "检查4: 硬编码错误检查"

    log_info "检查硬编码错误字符串..."

    local hardcoded_count=0

    # 常见的硬编码错误模式
    local patterns=(
        'errors\.New('
        'fmt\.Errorf('
        '"invalid"'
        '"not found"'
        '"already exists"'
        '"unauthorized"'
        '"forbidden"'
        '"internal error"'
    )

    for pattern in "${patterns[@]}"; do
        local count=$(grep -rn "$pattern" "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
            grep -v "_test.go" | \
            grep -v "//" | \
            grep -v "errorx" | \
            wc -l)

        if [ $count -gt 0 ]; then
            echo "  模式 '$pattern': $count 处"
            ((hardcoded_count += count))
        fi
    done

    echo ""

    if [ $hardcoded_count -eq 0 ]; then
        log_success "未发现硬编码错误"
        return 0
    else
        log_warning "发现 $hardcoded_count 处可能的硬编码错误"
        log_info "建议: 使用errno包中的错误定义"
        return 0
    fi
}

# ========================================
# 检查5: errno导入检查
# ========================================
check_errno_imports() {
    print_header "检查5: errno导入检查"

    log_info "检查errno包导入规范..."

    local files_without_errno=0
    local files_using_errorx=0

    # 统计导入errno的文件
    local errno_imports=$(grep -rn '"coze-studio/backend/types/errno"' "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
        grep -v "_test.go" | wc -l)

    # 统计导入errorx的文件
    local errorx_imports=$(grep -rn '"coze-studio/backend/pkg/errorx"' "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
        grep -v "_test.go" | wc -l)

    echo "  导入errno的文件: $errno_imports 个"
    echo "  导入errorx的文件: $errorx_imports 个"
    echo ""

    if [ $errno_imports -gt 0 ] && [ $errorx_imports -gt 0 ]; then
        log_success "errno和errorx包导入正确"
        return 0
    else
        log_warning "部分文件可能未正确导入errno/errorx包"
        return 0
    fi
}

# ========================================
# 检查6: 错误传播检查
# ========================================
check_error_propagation() {
    print_header "检查6: 错误传播检查"

    log_info "检查错误传播是否规范..."

    # 检查裸错误传播
    local bare_returns=$(grep -rn "return nil, err" "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
        grep -v "_test.go" | \
        grep -v "//" | wc -l)

    echo "  裸错误传播（return nil, err）: $bare_returns 处"
    echo ""

    if [ $bare_returns -gt 0 ]; then
        log_info "发现 $bare_returns 处裸错误传播"
        log_info "建议: 使用 errorx.Wrap(err, errno.ErrXXX) 添加上下文"
        return 0
    else
        log_success "未发现裸错误传播"
        return 0
    fi
}

# ========================================
# 生成检查结果
# ========================================
generate_report() {
    print_header "生成检查结果"

    # 生成JSON结果
    cat > "$RESULT_FILE" <<EOF
{
  "timestamp": "$(date -u '+%Y-%m-%dT%H:%M:%SZ')",
  "commit": "${GITHUB_SHA:-$(git rev-parse HEAD 2>/dev/null || echo 'unknown')}",
  "branch": "${GITHUB_REF_NAME:-$(git branch --show-current 2>/dev/null || echo 'unknown')}",
  "summary": {
    "total_issues": $TOTAL_ISSUES,
    "errors": $ERRORS,
    "warnings": $WARNINGS
  },
  "status": "$([ $ERRORS -eq 0 ] && echo 'passed' || echo 'failed')"
}
EOF

    echo "检查结果已保存到: $RESULT_FILE"
    cat "$RESULT_FILE" | jq '.' 2>/dev/null || cat "$RESULT_FILE"
    echo ""
}

# ========================================
# 主函数
# ========================================
main() {
    cd "$PROJECT_ROOT"

    print_header "errno使用规范验证"
    echo "开始时间: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "项目根目录: $PROJECT_ROOT"
    echo ""

    # 清空日志
    > "$LOG_FILE"

    # 执行检查
    local all_passed=true

    check_direct_errno_usage || all_passed=false
    check_errorx_usage || all_passed=false
    check_error_coverage || all_passed=false
    check_hardcoded_errors || all_passed=false
    check_errno_imports || all_passed=false
    check_error_propagation || all_passed=false

    # 生成报告
    generate_report

    # 总结
    print_header "验证总结"
    echo "结束时间: $(date '+%Y-%m-%d %H:%M:%S')"
    echo ""
    echo "统计信息:"
    echo "  总问题数: $TOTAL_ISSUES"
    echo "  错误数: $ERRORS"
    echo "  警告数: $WARNINGS"
    echo ""

    if [ $ERRORS -eq 0 ]; then
        print_header "✅ errno验证通过"
        echo ""
        echo "所有errno使用符合企业级规范！"
        echo ""
        echo "详细日志: $LOG_FILE"
        echo "结果文件: $RESULT_FILE"
        return 0
    else
        print_header "❌ errno验证失败"
        echo ""
        echo "发现 $ERRORS 个错误需要修复"
        echo ""
        echo "详细日志: $LOG_FILE"
        echo "结果文件: $RESULT_FILE"
        echo ""
        echo "修复建议:"
        echo "  1. 使用 errorx.New(errno.ErrXXX) 替换直接 errno 引用"
        echo "  2. 使用 errorx.Wrap(err, errno.ErrXXX) 添加错误上下文"
        echo "  3. 运行 ./scripts/errno/phase4_validate_all.sh 进行完整验证"
        return 1
    fi
}

# ========================================
# 执行
# ========================================
main "$@"
