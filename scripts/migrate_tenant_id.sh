#!/bin/bash

################################################################################
# ZKER tenant_id 迁移执行脚本
################################################################################
# 版本: v1.0
# 创建日期: 2025-01-01
# 说明: 安全执行 tenant_id 迁移，包含备份、执行、验证、回滚功能
################################################################################

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 数据库配置
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_NAME="${DB_NAME:-coze_studio}"
BACKUP_DIR="./backups"

################################################################################
# 辅助函数
################################################################################

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

################################################################################
# 功能函数
################################################################################

# 检查 MySQL 连接
check_mysql_connection() {
    log_info "检查 MySQL 连接..."
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -e "USE $DB_NAME; SELECT 1;" &>/dev/null; then
        log_success "MySQL 连接正常"
        return 0
    else
        log_error "MySQL 连接失败"
        log_error "请检查数据库配置和环境变量"
        exit 1
    fi
}

# 创建备份目录
create_backup_dir() {
    log_info "创建备份目录..."
    mkdir -p "$BACKUP_DIR"
    log_success "备份目录: $BACKUP_DIR"
}

# 备份数据库
backup_database() {
    local backup_file="$BACKUP_DIR/backup_before_tenant_id_$(date +%Y%m%d_%H%M%S).sql"

    log_info "备份数据库..."
    log_info "备份文件: $backup_file"

    if mysqldump -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" \
        --single-transaction \
        --routines \
        --triggers \
        --events \
        "$DB_NAME" > "$backup_file"; then
        log_success "数据库备份完成 ($(du -h "$backup_file" | cut -f1))"
    else
        log_error "数据库备份失败"
        exit 1
    fi
}

# 执行 SQL 脚本
execute_sql_script() {
    local sql_file="$1"
    local step_name="$2"

    log_info "$step_name"
    log_info "执行脚本: $sql_file"

    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" "$DB_NAME" < "$sql_file"; then
        log_success "$step_name 完成"
    else
        log_error "$step_name 失败"
        log_error "脚本: $sql_file"
        exit 1
    fi
}

# 验证迁移结果
validate_migration() {
    log_info "验证迁移结果..."

    # 检查是否有 tenant_id 字段为 NULL 的记录
    local query="
        SELECT
            TABLE_NAME,
            COLUMN_NAME,
            TABLE_ROWS
        FROM information_schema.TABLES t
        JOIN information_schema.COLUMNS c
            ON t.TABLE_NAME = c.TABLE_NAME
        WHERE t.TABLE_SCHEMA = '$DB_NAME'
          AND c.COLUMN_NAME = 'tenant_id'
        ORDER BY t.TABLE_NAME;
    "

    log_info "检查 tenant_id 字段添加情况..."
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -e "$query"

    # 检查是否还有 NULL 值
    local null_check="
        SELECT
            TABLE_NAME,
            COUNT(*) as total_rows,
            SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) as null_rows
        FROM information_schema.TABLES t
        WHERE t.TABLE_SCHEMA = '$DB_NAME'
          AND EXISTS (
              SELECT 1 FROM information_schema.COLUMNS c
              WHERE c.TABLE_NAME = t.TABLE_NAME
                AND c.COLUMN_NAME = 'tenant_id'
          )
        GROUP BY TABLE_NAME
        HAVING null_rows > 0;
    "

    log_info "检查 tenant_id NULL 值情况..."
    local null_count=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -N -e "$null_check" | wc -l)

    if [ "$null_count" -eq 0 ]; then
        log_success "所有表的 tenant_id 字段都已正确填充"
    else
        log_warning "以下表的 tenant_id 字段仍存在 NULL 值："
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -e "$null_check"
    fi
}

# 回滚迁移
rollback_migration() {
    log_warning "开始回滚迁移..."
    log_warning "这将删除所有 tenant_id 字段！"

    read -p "确认回滚? (输入 'YES' 继续): " confirm
    if [ "$confirm" != "YES" ]; then
        log_info "回滚已取消"
        return 0
    fi

    # 执行回滚 SQL
    local rollback_sql="$BACKUP_DIR/rollback_tenant_id.sql"

    log_info "生成回滚脚本..."
    cat > "$rollback_sql" << 'EOF'
SET FOREIGN_KEY_CHECKS = 0;

-- 删除所有 tenant_id 字段
ALTER TABLE app_draft DROP COLUMN tenant_id;
ALTER TABLE app_release_record DROP COLUMN tenant_id;
ALTER TABLE app_conversation_template_draft DROP COLUMN tenant_id;
ALTER TABLE conversation DROP COLUMN tenant_id;
ALTER TABLE message DROP COLUMN tenant_id;
ALTER TABLE knowledge DROP COLUMN tenant_id;
ALTER TABLE knowledge_document DROP COLUMN tenant_id;
ALTER TABLE knowledge_document_slice DROP COLUMN tenant_id;
ALTER TABLE files DROP COLUMN tenant_id;
ALTER TABLE agent_tool_draft DROP COLUMN tenant_id;
ALTER TABLE agent_tool_version DROP COLUMN tenant_id;
ALTER TABLE agent_to_database DROP COLUMN tenant_id;
ALTER TABLE chat_flow_role_config DROP COLUMN tenant_id;
ALTER TABLE connector_workflow_version DROP COLUMN tenant_id;
ALTER TABLE online_database_info DROP COLUMN tenant_id;
ALTER TABLE draft_database_info DROP COLUMN tenant_id;
ALTER TABLE api_key DROP COLUMN tenant_id;
ALTER TABLE app_connector_release_ref DROP COLUMN tenant_id;
ALTER TABLE app_dynamic_conversation_draft DROP COLUMN tenant_id;
ALTER TABLE app_dynamic_conversation_online DROP COLUMN tenant_id;
ALTER TABLE app_static_conversation_draft DROP COLUMN tenant_id;
ALTER TABLE app_static_conversation_online DROP COLUMN tenant_id;
ALTER TABLE app_conversation_template_online DROP COLUMN tenant_id;
ALTER TABLE data_copy_task DROP COLUMN tenant_id;
ALTER TABLE node_execution DROP COLUMN tenant_id;
ALTER TABLE knowledge_document_review DROP COLUMN tenant_id;

SET FOREIGN_KEY_CHECKS = 1;
EOF

    log_info "执行回滚..."
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" "$DB_NAME" < "$rollback_sql"; then
        log_success "回滚完成"
    else
        log_error "回滚失败"
        exit 1
    fi
}

################################################################################
# 主流程
################################################################################

show_usage() {
    cat << EOF
用法: $0 [COMMAND] [OPTIONS]

命令:
    migrate       执行完整的迁移流程（备份 + 添加字段 + 回填数据 + 验证）
    rollback      回滚迁移（删除所有 tenant_id 字段）
    validate       仅验证迁移结果
    backup        仅执行数据库备份

环境变量:
    DB_HOST       MySQL 主机地址 (默认: localhost)
    DB_PORT       MySQL 端口 (默认: 3306)
    DB_USER       MySQL 用户名 (默认: root)
    DB_NAME       数据库名称 (默认: coze_studio)

示例:
    # 完整迁移流程
    $0 migrate

    # 仅验证
    $0 validate

    # 回滚
    $0 rollback

    # 指定数据库
    DB_NAME=coze_production $0 migrate

EOF
}

# 主函数
main() {
    local command="${1:-migrate}"

    # 显示标题
    echo -e "${BLUE}============================================${NC}"
    echo -e "${BLUE}    ZKER tenant_id 迁移工具 v1.0${NC}"
    echo -e "${BLUE}============================================${NC}"
    echo ""

    # 检查命令
    case "$command" in
        migrate)
            log_info "开始完整迁移流程..."
            check_mysql_connection
            create_backup_dir
            backup_database
            execute_sql_script "backend/domain/tenant/migration/01_add_tenant_id_to_business_tables.sql" "步骤01: 添加 tenant_id 字段"
            execute_sql_script "backend/domain/tenant/migration/02_backfill_tenant_id_data.sql" "步骤02: 回填 tenant_id 数据"
            validate_migration
            log_success "🎉 迁移完成！"
            ;;
        rollback)
            rollback_migration
            ;;
        validate)
            check_mysql_connection
            validate_migration
            ;;
        backup)
            check_mysql_connection
            create_backup_dir
            backup_database
            ;;
        help|--help|-h)
            show_usage
            ;;
        *)
            log_error "未知命令: $command"
            show_usage
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
