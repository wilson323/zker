#!/bin/bash
# =====================================================
# 数据验证脚本: User和Space表tenant_id字段迁移验证
# 版本: v1.0.0
# 日期: 2025-01-01
# 说明: 验证tenant_id字段迁移是否成功
# =====================================================

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 数据库连接配置（从环境变量读取）
MYSQL_HOST="${MYSQL_HOST:-localhost}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-}"
MYSQL_DATABASE="${MYSQL_DATABASE:-opencoze}"

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# MySQL执行函数
execute_sql() {
    local sql="$1"
    mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" -e "$sql" 2>/dev/null
}

# 验证函数
verify_column_exists() {
    local table="$1"
    local column="$2"

    log_info "验证表 $table 是否有列 $column..."

    local count=$(execute_sql "
        SELECT COUNT(*)
        FROM INFORMATION_SCHEMA.COLUMNS
        WHERE TABLE_SCHEMA = '$MYSQL_DATABASE'
          AND TABLE_NAME = '$table'
          AND COLUMN_NAME = '$column';
    " | tail -n 1)

    if [ "$count" -eq 1 ]; then
        log_info "✅ 表 $table 的列 $column 存在"
        return 0
    else
        log_error "❌ 表 $table 的列 $column 不存在"
        return 1
    fi
}

verify_index_exists() {
    local table="$1"
    local index="$2"

    log_info "验证表 $table 是否有索引 $index..."

    local count=$(execute_sql "
        SELECT COUNT(*)
        FROM INFORMATION_SCHEMA.STATISTICS
        WHERE TABLE_SCHEMA = '$MYSQL_DATABASE'
          AND TABLE_NAME = '$table'
          AND INDEX_NAME = '$index';
    " | tail -n 1)

    if [ "$count" -gt 0 ]; then
        log_info "✅ 表 $table 的索引 $index 存在"
        return 0
    else
        log_warn "⚠️  表 $table 的索引 $index 不存在"
        return 1
    fi
}

verify_data_migration() {
    local table="$1"
    local column="$2"

    log_info "验证表 $table 的数据迁移..."

    local result=$(execute_sql "
        SELECT
            COUNT(*) AS total,
            COUNT(CASE WHEN $column IS NOT NULL AND $column != '' THEN 1 END) AS with_value,
            COUNT(CASE WHEN $column IS NULL OR $column = '' THEN 1 END) AS without_value
        FROM $table
        WHERE deleted_at IS NULL;
    " | tail -n 1)

    IFS=$'\t' read -r total with_value without_value <<< "$result"

    log_info "总记录数: $total"
    log_info "有值记录: $with_value"
    log_info "缺失记录: $without_value"

    if [ "$without_value" -eq 0 ]; then
        log_info "✅ 表 $table 的数据迁移成功"
        return 0
    else
        log_error "❌ 表 $table 有 $without_value 条记录缺失 $column 值"
        return 1
    fi
}

verify_consistency() {
    log_info "验证数据一致性..."

    # 验证space_user的tenant_id应该与space的tenant_id一致
    local result=$(execute_sql "
        SELECT COUNT(*) AS count
        FROM space_user su
        INNER JOIN space s ON su.space_id = s.id
        WHERE su.tenant_id != s.tenant_id;
    " | tail -n 1)

    if [ "$result" -eq 0 ]; then
        log_info "✅ space_user的tenant_id与space的tenant_id一致"
        return 0
    else
        log_error "❌ 有 $result 条space_user记录的tenant_id与space不一致"
        return 1
    fi
}

# =====================================================
# 主验证流程
# =====================================================

main() {
    log_info "=========================================="
    log_info "开始验证tenant_id字段迁移"
    log_info "数据库: $MYSQL_DATABASE"
    log_info "=========================================="

    local failed=0

    # 1. 验证user表
    echo ""
    log_info "【1/3】验证user表..."
    if ! verify_column_exists "user" "tenant_id"; then
        failed=1
    fi
    if ! verify_index_exists "user" "idx_tenant_id"; then
        failed=1
    fi
    if ! verify_index_exists "user" "idx_tenant_id_email"; then
        failed=1
    fi
    if ! verify_data_migration "user" "tenant_id"; then
        failed=1
    fi

    # 2. 验证space表
    echo ""
    log_info "【2/3】验证space表..."
    if ! verify_column_exists "space" "tenant_id"; then
        failed=1
    fi
    if ! verify_index_exists "space" "idx_tenant_id"; then
        failed=1
    fi
    if ! verify_index_exists "space" "idx_tenant_id_owner"; then
        failed=1
    fi
    if ! verify_data_migration "space" "tenant_id"; then
        failed=1
    fi

    # 3. 验证space_user表
    echo ""
    log_info "【3/3】验证space_user表..."
    if ! verify_column_exists "space_user" "tenant_id"; then
        failed=1
    fi
    if ! verify_index_exists "space_user" "idx_tenant_id"; then
        failed=1
    fi
    if ! verify_index_exists "space_user" "idx_tenant_id_user"; then
        failed=1
    fi
    if ! verify_data_migration "space_user" "tenant_id"; then
        failed=1
    fi

    # 4. 验证数据一致性
    echo ""
    log_info "【4/4】验证数据一致性..."
    if ! verify_consistency; then
        failed=1
    fi

    # 5. 统计数据
    echo ""
    log_info "【统计数据】"
    execute_sql "
        SELECT
            'Total Users' AS metric,
            COUNT(*) AS value
        FROM user
        WHERE deleted_at IS NULL

        UNION ALL

        SELECT
            'Users with Tenant ID' AS metric,
            COUNT(*) AS value
        FROM user
        WHERE tenant_id IS NOT NULL
          AND tenant_id != ''
          AND deleted_at IS NULL

        UNION ALL

        SELECT
            'Total Spaces' AS metric,
            COUNT(*) AS value
        FROM space
        WHERE deleted_at IS NULL

        UNION ALL

        SELECT
            'Spaces with Tenant ID' AS metric,
            COUNT(*) AS value
        FROM space
        WHERE tenant_id IS NOT NULL
          AND tenant_id != ''
          AND deleted_at IS NULL

        UNION ALL

        SELECT
            'Total Space Users' AS metric,
            COUNT(*) AS value
        FROM space_user

        UNION ALL

        SELECT
            'Space Users with Tenant ID' AS metric,
            COUNT(*) AS value
        FROM space_user
        WHERE tenant_id IS NOT NULL
          AND tenant_id != '';
    "

    # 最终结果
    echo ""
    log_info "=========================================="
    if [ $failed -eq 0 ]; then
        log_info "✅ 所有验证通过！迁移成功！"
        log_info "=========================================="
        exit 0
    else
        log_error "❌ 验证失败！请检查错误信息！"
        log_info "=========================================="
        exit 1
    fi
}

# 执行主函数
main "$@"
