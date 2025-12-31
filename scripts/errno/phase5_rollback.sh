#!/bin/bash
# Phase 5: 快速回滚到迁移前的状态
# 作者: ZKER架构团队
# 日期: 2025-01-03
# 用途: 企业级errno迁移的第五步 - 安全回滚

set -euo pipefail

# ========================================
# 配置
# ========================================
BACKUP_DIR="scripts/errno/backups"
BACKEND_DIR="backend"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
ROLLBACK_BACKUP="$BACKUP_DIR/before_rollback_$TIMESTAMP"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

# 检查备份目录
check_backup_dir() {
    log_info "检查备份目录..."
    echo ""

    if [ ! -d "$BACKUP_DIR" ]; then
        log_error "备份目录不存在: $BACKUP_DIR"
        echo ""
        echo "可能的原因:"
        echo "  1. 还没有执行过迁移"
        echo "  2. 备份目录被误删"
        echo ""
        echo "替代方案:"
        echo "  使用Git回滚: ./scripts/errno/phase5_rollback_git.sh"
        echo ""
        exit 1
    fi

    # 统计备份文件
    local backup_count=$(find "$BACKUP_DIR/backend" -type f 2>/dev/null | wc -l)

    if [ "$backup_count" -eq 0 ]; then
        log_error "备份目录为空: $BACKUP_DIR"
        echo ""
        exit 1
    fi

    log_success "找到 $backup_count 个备份文件"
    echo ""

    # 显示备份目录结构
    log_info "备份目录结构:"
    tree -L 2 "$BACKUP_DIR/backend" 2>/dev/null || \
        find "$BACKUP_DIR/backend" -type d | head -20
    echo ""
}

# 创建当前状态的备份
backup_current_state() {
    log_info "备份当前状态（以防万一）..."
    echo ""

    mkdir -p "$ROLLBACK_BACKUP"

    local backup_count=0

    # 获取所有修改过的文件
    local modified_files=$(git diff --name-only 2>/dev/null || echo "")

    if [ -z "$modified_files" ]; then
        log_warning "没有检测到修改的文件"
        log_info "将备份整个backend目录..."

        # 备份整个backend目录
        cp -r "$BACKEND_DIR" "$ROLLBACK_BACKUP/backend"
        backup_count=$(find "$ROLLBACK_BACKUP/backend" -type f | wc -l)
    else
        # 只备份修改过的文件
        while IFS= read -r file; do
            if [ -f "$file" ]; then
                local target_dir="$ROLLBACK_BACKUP/$(dirname "$file")"
                mkdir -p "$target_dir"
                cp "$file" "$target_dir/"
                backup_count=$((backup_count + 1))
            fi
        done <<< "$modified_files"
    fi

    log_success "已备份 $backup_count 个文件到: $ROLLBACK_BACKUP"
    echo ""
}

# 从备份恢复文件
restore_from_backup() {
    log_info "从备份恢复文件..."
    echo ""

    local restored=0
    local failed=0

    # 查找所有备份文件
    find "$BACKUP_DIR/backend" -type f | while read backup_file; do
        # 计算相对路径
        local rel_path="${backup_file#$BACKUP_DIR/}"
        local target_file="$rel_path"

        # 检查目标文件是否存在
        if [ ! -f "$target_file" ]; then
            log_warning "目标文件不存在，跳过: $target_file"
            continue
        fi

        # 恢复文件
        if cp "$backup_file" "$target_file"; then
            log_success "恢复: $target_file"
            restored=$((restored + 1))
        else
            log_error "恢复失败: $target_file"
            failed=$((failed + 1))
        fi
    done

    echo ""
    log_success "恢复完成"
    echo "  成功: $restored 个文件"
    echo "  失败: $failed 个文件"
    echo ""
}

# 验证恢复结果
verify_restore() {
    log_info "验证恢复结果..."
    echo ""

    # 1. 检查Go语法
    echo "  1. 检查Go语法..."
    if go fmt ./... > /dev/null 2>&1; then
        log_success "   语法检查通过"
    else
        log_error "   语法检查失败"
        return 1
    fi

    # 2. 尝试编译
    echo "  2. 尝试编译..."
    if go build ./... > /dev/null 2>&1; then
        log_success "   编译成功"
    else
        log_warning "   编译失败（可能是因为回滚到错误版本）"
    fi

    # 3. 检查errno引用
    echo "  3. 检查errno引用..."
    local direct_errno=$(grep -rn "return.*errno\." "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
        grep -v "_test.go" | wc -l)
    echo "     直接errno引用: $direct_errno"

    echo ""
    log_success "验证完成"
    echo ""
}

# 显示回滚后操作指引
show_post_rollback_instructions() {
    print_header "回滚后操作指引"

    echo "1. 检查Git状态:"
    echo "   git status"
    echo ""

    echo "2. 查看恢复的文件:"
    echo "   git diff"
    echo ""

    echo "3. 重新编译验证:"
    echo "   cd $BACKEND_DIR && go build ./..."
    echo ""

    echo "4. 运行测试:"
    echo "   cd $BACKEND_DIR && go test ./..."
    echo ""

    echo "5. 如果需要，可以使用Git进一步回滚:"
    echo "   git log --oneline -10"
    echo "   git reset --hard <commit-hash>"
    echo ""

    echo "6. 当前状态备份位置:"
    echo "   $ROLLBACK_BACKUP"
    echo ""
}

# ========================================
# 主函数
# ========================================

main() {
    print_header "errno迁移回滚"

    log_warning "⚠️  警告：此操作将覆盖当前文件！"
    echo ""
    echo "此操作将:"
    echo "  1. 备份当前状态到: $ROLLBACK_BACKUP"
    echo "  2. 从备份目录恢复文件: $BACKUP_DIR"
    echo "  3. 验证恢复结果"
    echo ""
    echo "备份目录: $BACKUP_DIR"
    echo ""

    # 检查备份目录
    check_backup_dir

    # 确认回滚
    echo "请确认回滚操作:"
    echo "  • 输入 'yes' 继续回滚"
    echo "  • 输入 'q'   退出"
    echo ""
    read -p "> " confirm

    if [ "$confirm" != "yes" ]; then
        echo ""
        log_info "已取消回滚"
        exit 0
    fi

    echo ""
    log_info "开始回滚..."
    echo ""

    # 1. 备份当前状态
    backup_current_state

    # 2. 从备份恢复
    restore_from_backup

    # 3. 验证恢复结果
    if verify_restore; then
        print_header "✅ 回滚成功！"
    else
        print_header "⚠️  回滚完成（验证有警告）"
    fi

    # 4. 显示后续操作指引
    show_post_rollback_instructions
}

# 执行主函数
main "$@"
