#!/bin/bash
# scripts/migrate.sh
# 数据库迁移自动化脚本
# 用途: 执行数据库schema迁移、版本管理、回滚操作

set -e  # 遇到错误立即退出

# ========== 颜色输出定义 ==========
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ========== 配置项 ==========
MIGRATIONS_DIR="${MIGRATIONS_DIR:-./docker/atlas/migrations}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASSWORD="${DB_PASSWORD:-root}"
DB_NAME="${DB_NAME:-zker}"
ATLAS_DIR="${ATLAS_DIR:-./docker/atlas}"

# ========== 日志函数 ==========
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# ========== 检查函数 ==========
check_atlas_installed() {
    if ! command -v atlas &> /dev/null; then
        log_warn "Atlas CLI not found. Installing..."
        curl -sSf https://atlasgo.sh | sh -s -- -y --community
        if [ $? -eq 0 ]; then
            log_info "Atlas CLI installed successfully"
        else
            log_error "Failed to install Atlas CLI"
            exit 1
        fi
    else
        log_info "Atlas CLI already installed: $(atlas version)"
    fi
}

check_mysql_connection() {
    log_step "检查MySQL连接..."
    if mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1" &> /dev/null; then
        log_info "MySQL连接正常"
        return 0
    else
        log_error "无法连接到MySQL: $DB_HOST:$DB_PORT"
        exit 1
    fi
}

check_migrations_dir() {
    if [ ! -d "$MIGRATIONS_DIR" ]; then
        log_warn "Migrations目录不存在: $MIGRATIONS_DIR"
        log_info "将创建目录: $MIGRATIONS_DIR"
        mkdir -p "$MIGRATIONS_DIR"
    fi
    log_info "Migrations目录: $MIGRATIONS_DIR"
}

# ========== 迁移函数 ==========
run_atlas_apply() {
    local schema_file=$1

    if [ ! -f "$schema_file" ]; then
        log_error "Schema文件不存在: $schema_file"
        exit 1
    fi

    log_step "执行Atlas迁移..."
    ATLAS_URL="mysql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME"

    atlas schema apply \
        -u "$ATLAS_URL" \
        --to "file://$schema_file" \
        --exclude "atlas_schema_revisions,table_*" \
        --dev-url "mysql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/${DB_NAME}_dev" \
        --format "{{ sql }}\n"

    if [ $? -eq 0 ]; then
        log_info "Atlas迁移成功"
    else
        log_error "Atlas迁移失败"
        exit 1
    fi
}

run_sql_migration() {
    local migration_file=$1
    local migration_name=$(basename "$migration_file" .sql)

    log_step "执行SQL迁移: $migration_name"

    # 执行SQL文件
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$migration_file"

    if [ $? -eq 0 ]; then
        log_info "SQL迁移成功: $migration_name"
    else
        log_error "SQL迁移失败: $migration_name"
        exit 1
    fi
}

run_all_sql_migrations() {
    log_step "执行所有SQL迁移文件..."

    local count=0
    for migration in $(ls -v "$MIGRATIONS_DIR"/*.sql 2>/dev/null || true); do
        if [ -f "$migration" ]; then
            run_sql_migration "$migration"
            count=$((count + 1))
        fi
    done

    if [ $count -eq 0 ]; then
        log_warn "未找到SQL迁移文件"
    else
        log_info "共执行了 $count 个SQL迁移"
    fi
}

# ========== 状态查询函数 ==========
show_migration_status() {
    log_step "查询迁移状态..."

    # 检查是否存在schema_migrations表
    local table_exists=$(mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -sN -e \
        "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = '$DB_NAME' AND table_name = 'schema_migrations';" 2>/dev/null || echo "0")

    if [ "$table_exists" = "0" ]; then
        log_warn "schema_migrations表不存在,数据库可能未初始化"
        return
    fi

    # 查询迁移记录
    echo ""
    echo "========================================"
    echo "           迁移历史记录"
    echo "========================================"
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -t -e \
        "SELECT * FROM schema_migrations ORDER BY version DESC LIMIT 20;"
    echo "========================================"
}

show_db_info() {
    log_step "数据库信息..."

    echo ""
    echo "========================================"
    echo "           数据库概览"
    echo "========================================"
    echo "主机: $DB_HOST:$DB_PORT"
    echo "数据库: $DB_NAME"
    echo "用户: $DB_USER"
    echo ""

    # 查询表数量
    local table_count=$(mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -sN -e \
        "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = '$DB_NAME';" 2>/dev/null || echo "0")
    echo "表数量: $table_count"
    echo ""

    # 查询数据库大小
    local db_size=$(mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -sN -e \
        "SELECT ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) FROM information_schema.tables WHERE table_schema = '$DB_NAME';" 2>/dev/null || echo "0")
    echo "数据库大小: ${db_size} MB"
    echo "========================================"
}

# ========== 创建迁移模板函数 ==========
create_migration_template() {
    local migration_name=$1
    local timestamp=$(date +%Y%m%d%H%M%S)
    local migration_file="$MIGRATIONS_DIR/${timestamp}_${migration_name}.sql"

    if [ -z "$migration_name" ]; then
        log_error "请提供迁移名称"
        echo "用法: $0 create <migration_name>"
        exit 1
    fi

    log_step "创建迁移模板: $migration_file"

    cat > "$migration_file" <<EOF
-- Migration: ${migration_name}
-- Created: $(date '+%Y-%m-%d %H:%M:%S')
-- Description: 请在此处添加迁移描述

-- 1. 创建表 (示例)
-- CREATE TABLE IF NOT EXISTS example_table (
--     id VARCHAR(36) PRIMARY KEY COMMENT '主键ID',
--     name VARCHAR(200) NOT NULL COMMENT '名称',
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
--     deleted_at TIMESTAMP NULL,
--     INDEX idx_name (name),
--     INDEX idx_created_at (created_at)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='示例表';

-- 2. 修改表 (示例)
-- ALTER TABLE existing_table ADD COLUMN new_column VARCHAR(100);

-- 3. 添加索引 (示例)
-- CREATE INDEX idx_new_column ON existing_table(new_column);

-- 4. 记录迁移版本
INSERT INTO schema_migrations (version, name, applied_at)
VALUES (${timestamp}, '${timestamp}_${migration_name}', NOW())
ON DUPLICATE KEY UPDATE applied_at = NOW();
EOF

    log_info "迁移模板已创建: $migration_file"
    log_info "请编辑文件并添加迁移内容"
}

# ========== 清理函数 ==========
cleanup_atlas_revisions() {
    log_step "清理Atlas版本记录..."

    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e \
        "DELETE FROM atlas_schema_revisions WHERE created_at < DATE_SUB(NOW(), INTERVAL 30 DAY);"

    if [ $? -eq 0 ]; then
        log_info "清理完成"
    else
        log_warn "清理失败或表不存在"
    fi
}

# ========== 帮助函数 ==========
show_usage() {
    cat << EOF
数据库迁移自动化脚本 v1.0

用法: $0 <command> [options]

命令:
  up [schema_file]        执行迁移 (默认使用Atlas schema)
  down                    回滚最后一次迁移 (需要手动编写down SQL)
  status                  查看迁移状态
  info                    显示数据库信息
  create <name>           创建新的迁移模板
  cleanup                 清理旧的Atlas版本记录
  help                    显示此帮助信息

环境变量:
  DB_HOST                 MySQL主机 (默认: localhost)
  DB_PORT                 MySQL端口 (默认: 3306)
  DB_USER                 MySQL用户 (默认: root)
  DB_PASSWORD             MySQL密码 (默认: root)
  DB_NAME                 数据库名称 (默认: zker)
  MIGRATIONS_DIR          迁移文件目录 (默认: ./docker/atlas/migrations)

示例:
  # 执行Atlas迁移
  $0 up

  # 执行特定schema文件
  $0 up ./docker/atlas/zker_schema.hcl

  # 查看迁移状态
  $0 status

  # 创建新迁移
  $0 create add_tenant_id_to_bots

  # 清理旧版本
  $0 cleanup

EOF
}

# ========== 主函数 ==========
main() {
    local command=${1:-help}
    shift || true

    case $command in
        up)
            log_info "开始数据库迁移..."
            check_atlas_installed
            check_mysql_connection
            check_migrations_dir

            local schema_file="${1:-$ATLAS_DIR/zker_schema.hcl}"

            # 先执行SQL迁移
            if [ -d "$MIGRATIONS_DIR" ] && [ "$(ls -A $MIGRATIONS_DIR/*.sql 2>/dev/null)" ]; then
                run_all_sql_migrations
            fi

            # 再执行Atlas迁移
            if [ -f "$schema_file" ]; then
                run_atlas_apply "$schema_file"
            else
                log_warn "未找到Atlas schema文件: $schema_file"
            fi

            log_info "数据库迁移完成!"
            show_migration_status
            ;;

        down)
            log_warn "回滚功能需要手动编写down SQL脚本"
            log_warn "请手动执行回滚SQL语句"
            exit 1
            ;;

        status)
            check_mysql_connection
            show_migration_status
            show_db_info
            ;;

        info)
            check_mysql_connection
            show_db_info
            ;;

        create)
            create_migration_template "$1"
            ;;

        cleanup)
            check_mysql_connection
            cleanup_atlas_revisions
            ;;

        help|--help|-h)
            show_usage
            ;;

        *)
            log_error "未知命令: $command"
            echo ""
            show_usage
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
