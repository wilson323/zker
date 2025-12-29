#!/bin/bash

# =====================================================
# 租户隔离数据库迁移自动化脚本
# 版本: v1.0
# 日期: 2025-01-01
# 作者: 研发B（后端工程师）
# 说明: 自动执行数据库备份、迁移和验证
# =====================================================

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
MYSQL_CONTAINER="ioedream-mysql"
MYSQL_USER="coze"
MYSQL_PASSWORD="coze123"
MYSQL_DATABASE="opencoze"
BACKUP_DIR="./backup/migration_$(date +%Y%m%d_%H%M%S)"
MIGRATION_SCRIPT="./docker/atlas/migrations/20251230120000_add_user_tenant_isolation.sql"
ROLLBACK_SCRIPT="./docker/atlas/migrations/20251230120000_add_user_tenant_isolation_rollback.sql"

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

# 步骤1: 检查环境
check_environment() {
    log_info "Step 1: Checking environment..."

    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi

    # 检查MySQL容器
    if ! docker ps | grep -q $MYSQL_CONTAINER; then
        log_error "MySQL container ($MYSQL_CONTAINER) is not running"
        exit 1
    fi
    log_success "MySQL container is running"

    # 检查迁移脚本
    if [ ! -f "$MIGRATION_SCRIPT" ]; then
        log_error "Migration script not found: $MIGRATION_SCRIPT"
        exit 1
    fi
    log_success "Migration script found"

    # 检查回滚脚本
    if [ ! -f "$ROLLBACK_SCRIPT" ]; then
        log_error "Rollback script not found: $ROLLBACK_SCRIPT"
        exit 1
    fi
    log_success "Rollback script found"

    # 创建备份目录
    mkdir -p "$BACKUP_DIR"
    log_success "Backup directory created: $BACKUP_DIR"
}

# 步骤2: 备份数据库
backup_database() {
    log_info "Step 2: Backing up database..."

    local backup_file="$BACKUP_DIR/opencoze_backup_$(date +%Y%m%d_%H%M%S).sql"

    log_info "Creating backup: $backup_file"
    docker exec $MYSQL_CONTAINER mysqldump \
        -u$MYSQL_USER \
        -p$MYSQL_PASSWORD \
        $MYSQL_DATABASE \
        > "$backup_file"

    if [ $? -eq 0 ]; then
        local backup_size=$(du -h "$backup_file" | cut -f1)
        log_success "Backup completed (size: $backup_size)"
    else
        log_error "Backup failed"
        exit 1
    fi
}

# 步骤3: 执行迁移
execute_migration() {
    log_info "Step 3: Executing migration..."

    log_warning "This will modify the database structure. Press Ctrl+C to cancel..."
    sleep 3

    log_info "Executing migration script..."
    docker exec -i $MYSQL_CONTAINER mysql \
        -u$MYSQL_USER \
        -p$MYSQL_PASSWORD \
        $MYSQL_DATABASE < "$MIGRATION_SCRIPT"

    if [ $? -eq 0 ]; then
        log_success "Migration script executed successfully"
    else
        log_error "Migration failed! Please check the error messages above."
        log_info "You can restore from backup: $backup_file"
        exit 1
    fi
}

# 步骤4: 验证迁移结果
verify_migration() {
    log_info "Step 4: Verifying migration results..."

    # 验证1: users表应该有tenant_id列
    log_info "Checking if users table has tenant_id column..."
    local column_count=$(docker exec $MYSQL_CONTAINER mysql \
        -u$MYSQL_USER \
        -p$MYSQL_PASSWORD \
        $MYSQL_DATABASE \
        -se "SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA='$MYSQL_DATABASE' AND TABLE_NAME='users' AND COLUMN_NAME='tenant_id';")

    if [ "$column_count" -eq "1" ]; then
        log_success "✓ users.tenant_id column exists"
    else
        log_error "✗ users.tenant_id column not found!"
        return 1
    fi

    # 验证2: user_tenants表应该存在
    log_info "Checking if user_tenants table exists..."
    local table_count=$(docker exec $MYSQL_CONTAINER mysql \
        -u$MYSQL_USER \
        -p$MYSQL_PASSWORD \
        $MYSQL_DATABASE \
        -se "SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='$MYSQL_DATABASE' AND TABLE_NAME='user_tenants';")

    if [ "$table_count" -eq "1" ]; then
        log_success "✓ user_tenants table exists"
    else
        log_error "✗ user_tenants table not found!"
        return 1
    fi

    # 验证3: 所有users都应该有tenant_id
    log_info "Checking if all users have tenant_id..."
    local users_without_tenant=$(docker exec $MYSQL_CONTAINER mysql \
        -u$MYSQL_USER \
        -p$MYSQL_PASSWORD \
        $MYSQL_DATABASE \
        -se "SELECT COUNT(*) FROM users WHERE tenant_id IS NULL AND deleted_at IS NULL;")

    if [ "$users_without_tenant" -eq "0" ]; then
        log_success "✓ All users have tenant_id"
    else
        log_warning "✗ $users_without_tenant users without tenant_id"
    fi

    # 验证4: 统计数据
    log_info "Collecting statistics..."
    echo ""
    echo "=== Migration Statistics ==="

    local total_users=$(docker exec $MYSQL_CONTAINER mysql \
        -u$MYSQL_USER \
        -p$MYSQL_PASSWORD \
        $MYSQL_DATABASE \
        -se "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL;")

    local users_with_tenant=$(docker exec $MYSQL_CONTAINER mysql \
        -u$MYSQL_USER \
        -p$MYSQL_PASSWORD \
        $MYSQL_DATABASE \
        -se "SELECT COUNT(*) FROM users WHERE tenant_id IS NOT NULL AND deleted_at IS NULL;")

    local total_tenants=$(docker exec $MYSQL_CONTAINER mysql \
        -u$MYSQL_USER \
        -p$MYSQL_PASSWORD \
        $MYSQL_DATABASE \
        -se "SELECT COUNT(*) FROM tenants WHERE deleted_at IS NULL;")

    local individual_tenants=$(docker exec $MYSQL_CONTAINER mysql \
        -u$MYSQL_USER \
        -p$MYSQL_PASSWORD \
        $MYSQL_DATABASE \
        -se "SELECT COUNT(*) FROM tenants WHERE tenant_type='individual' AND deleted_at IS NULL;")

    local user_tenant_relations=$(docker exec $MYSQL_CONTAINER mysql \
        -u$MYSQL_USER \
        -p$MYSQL_PASSWORD \
        $MYSQL_DATABASE \
        -se "SELECT COUNT(*) FROM user_tenants WHERE deleted_at IS NULL;")

    echo "Total Users: $total_users"
    echo "Users with Tenant ID: $users_with_tenant"
    echo "Total Tenants: $total_tenants"
    echo "Individual Tenants: $individual_tenants"
    echo "User-Tenant Relations: $user_tenant_relations"
    echo "========================="
    echo ""
}

# 步骤5: 生成验证报告
generate_report() {
    log_info "Step 5: Generating verification report..."

    local report_file="$BACKUP_DIR/migration_report.txt"

    cat > "$report_file" << EOF
====================================================
Tenant Isolation Migration Report
====================================================
Date: $(date)
Database: $MYSQL_DATABASE
Container: $MYSQL_CONTAINER

====================================================
Migration Summary
====================================================

Migration Script: $MIGRATION_SCRIPT
Backup Location: $BACKUP_DIR

====================================================
Verification Results
====================================================

EOF

    # 添加详细统计
    docker exec $MYSQL_CONTAINER mysql \
        -u$MYSQL_USER \
        -p$MYSQL_PASSWORD \
        $MYSQL_DATABASE \
        -e "
SELECT
    'Total Users' AS metric,
    COUNT(*) AS value
FROM users
WHERE deleted_at IS NULL

UNION ALL

SELECT
    'Users with Tenant ID' AS metric,
    COUNT(*) AS value
FROM users
WHERE tenant_id IS NOT NULL
  AND deleted_at IS NULL

UNION ALL

SELECT
    'Total Tenants' AS metric,
    COUNT(*) AS value
FROM tenants
WHERE deleted_at IS NULL

UNION ALL

SELECT
    'Individual Tenants' AS metric,
    COUNT(*) AS value
FROM tenants
WHERE tenant_type = 'individual'
  AND deleted_at IS NULL

UNION ALL

SELECT
    'User-Tenant Relations' AS metric,
    COUNT(*) AS value
FROM user_tenants
WHERE deleted_at IS NULL;
" >> "$report_file"

    cat >> "$report_file" << EOF

====================================================
Next Steps
====================================================

1. Test the application with the new tenant isolation
2. Monitor for any errors or performance issues
3. If issues arise, rollback using:
   docker exec -i $MYSQL_CONTAINER mysql -u$MYSQL_USER -p$MYSQL_PASSWORD $MYSQL_DATABASE < $ROLLBACK_SCRIPT

====================================================
End of Report
====================================================
EOF

    log_success "Report generated: $report_file"
    echo ""
    cat "$report_file"
}

# 主函数
main() {
    echo ""
    echo "===================================================="
    echo "  Tenant Isolation Database Migration"
    echo "===================================================="
    echo ""

    # 检查环境
    check_environment
    echo ""

    # 备份数据库
    backup_database
    echo ""

    # 执行迁移
    execute_migration
    echo ""

    # 验证迁移
    verify_migration
    if [ $? -ne 0 ]; then
        log_error "Migration verification failed!"
        log_info "Please restore from backup or check the errors manually"
        exit 1
    fi
    echo ""

    # 生成报告
    generate_report
    echo ""

    log_success "===================================================="
    log_success "Migration completed successfully!"
    log_success "===================================================="
    log_info "Backup location: $BACKUP_DIR"
    log_info "If you need to rollback, use: $ROLLBACK_SCRIPT"
    echo ""
}

# 执行主函数
main "$@"
