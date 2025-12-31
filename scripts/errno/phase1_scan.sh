#!/bin/bash
# Phase 1: 扫描所有errno引用，生成详细清单
# 作者: ZKER架构团队
# 日期: 2025-01-03
# 用途: 企业级errno迁移的第一步 - 静态分析

set -euo pipefail

# ========================================
# 配置
# ========================================
BACKEND_DIR="backend"
OUTPUT_DIR="scripts/errno/analysis"
REPORT_FILE="$OUTPUT_DIR/errno_inventory_report.csv"
SUMMARY_FILE="$OUTPUT_DIR/errno_summary.txt"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# 初始化输出目录
init_output_dir() {
    mkdir -p "$OUTPUT_DIR"
    log_info "创建输出目录: $OUTPUT_DIR"

    # 清空报告文件
    > "$REPORT_FILE"
    echo "file_path,line_number,context,pattern,recommended_action,reason" > "$REPORT_FILE"
}

# 扫描单个文件
scan_file() {
    local file="$1"
    local rel_path="${file#$BACKEND_DIR/}"

    local file_refs=0
    local file_should_migrate=0
    local file_should_skip=0

    # 使用grep获取所有包含errno的行
    while IFS= read -r line; do
        # 提取行号
        local line_num=$(echo "$line" | cut -d: -f1)

        # 提取行内容
        local line_content=$(echo "$line" | cut -d: -f2-)

        # 提取上下文（前后2行）
        local start_line=$((line_num - 2))
        local end_line=$((line_num + 2))
        local context=$(sed -n "${start_line},${end_line}p" "$file" | sed 's/,/;/g' | tr '\n' '|')

        # 提取errno模式
        local pattern=$(echo "$line_content" | grep -oP 'errno\.[A-Z][a-zA-Z0-9_]+' || echo "")

        # 判断是否需要迁移
        local action="MIGRATE"
        local reason=""

        # 跳过的情况
        if echo "$line_content" | grep -qP '^\s*//'; then
            action="SKIP_COMMENT"
            reason="注释中的errno引用"
            file_should_skip=$((file_should_skip + 1))
        elif echo "$line_content" | grep -qP '(assert|require)\.(Error|Equal).*errno\.'; then
            action="SKIP_TEST_ASSERT"
            reason="测试断言中的errno引用"
            file_should_skip=$((file_should_skip + 1))
        elif echo "$line_content" | grep -qP '"[^"]*errno\.'; then
            action="SKIP_STRING_LITERAL"
            reason="字符串字面量中的errno"
            file_should_skip=$((file_should_skip + 1))
        elif echo "$line_content" | grep -qP 'errorx\.(New|Wrap|WrapByCode).*errno\.'; then
            action="SKIP_ALREADY_MIGRATED"
            reason="已经使用errorx包装"
            file_should_skip=$((file_should_skip + 1))
        elif echo "$line_content" | grep -qP '\s+//\s+TODO'; then
            action="SKIP_TODO"
            reason="TODO注释，暂时跳过"
            file_should_skip=$((file_should_skip + 1))
        else
            file_should_migrate=$((file_should_migrate + 1))
        fi

        # 写入CSV（转义逗号和引号）
        local escaped_context=$(echo "$context" | sed 's/"/""/g')
        echo "\"$rel_path\",\"$line_num\",\"$escaped_context\",\"$pattern\",\"$action\",\"$reason\"" >> "$REPORT_FILE"

        file_refs=$((file_refs + 1))
    done < <(grep -n "errno\." "$file" || true)

    # 返回统计信息
    echo "$file_refs|$file_should_migrate|$file_should_skip"
}

# 扫描所有文件
scan_all_files() {
    log_info "开始扫描errno引用..."
    echo ""

    local total_files=0
    local total_refs=0
    local total_should_migrate=0
    local total_should_skip=0

    # 扫描所有Go文件（排除测试文件和生成文件）
    while IFS= read -r -d '' file; do
        # 跳过测试文件
        [[ "$file" == *"_test.go" ]] && continue
        # 跳过生成文件
        [[ "$file" == *"/gen/"* ]] && continue
        [[ "$file" == *"/gen.go" ]] && continue

        local rel_path="${file#$BACKEND_DIR/}"
        log_info "扫描: $rel_path"

        # 扫描文件
        local stats=$(scan_file "$file")
        IFS='|' read -r file_refs file_should_migrate file_should_skip <<< "$stats"

        if [ "$file_refs" -gt 0 ]; then
            log_success "$rel_path: $file_refs处引用 ($file_should_migrate需迁移, $file_should_skip应跳过)"
            total_files=$((total_files + 1))
            total_refs=$((total_refs + file_refs))
            total_should_migrate=$((total_should_migrate + file_should_migrate))
            total_should_skip=$((total_should_skip + file_should_skip))
        fi
    done < <(find "$BACKEND_DIR" -name "*.go" -print0)

    # 输出统计
    echo ""
    log_success "扫描完成！"
    echo ""
    echo "📊 统计结果:"
    echo "  总文件数: $total_files"
    echo "  总引用数: $total_refs"
    echo "  需迁移: $total_should_migrate"
    echo "  应跳过: $total_should_skip"
    echo ""

    # 保存汇总信息
    cat > "$SUMMARY_FILE" << EOF
errno扫描汇总
==================
扫描时间: $(date '+%Y-%m-%d %H:%M:%S')

统计结果:
  总文件数: $total_files
  总引用数: $total_refs
  需迁移: $total_should_migrate
  应跳过: $total_should_skip

详细报告: $REPORT_FILE
EOF

    log_info "汇总信息已保存到: $SUMMARY_FILE"
}

# ========================================
# 主函数
# ========================================
main() {
    echo "=========================================="
    echo "  errno迁移 - Phase 1: 静态扫描"
    echo "=========================================="
    echo ""

    # 初始化
    init_output_dir

    # 扫描
    scan_all_files

    echo ""
    log_success "Phase 1 完成！"
    echo ""
    echo "下一步:"
    echo "  1. 查看详细报告: cat $REPORT_FILE | column -t -s, -o '|'"
    echo "  2. 生成Markdown报告: ./scripts/errno/phase1_generate_report.sh"
    echo ""
}

# 执行主函数
main "$@"
