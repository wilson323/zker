#!/bin/bash

# =====================================================
# ioedream数据库租户隔离迁移
# 版本: v1.0
# =====================================================

set -e

MYSQL_CONTAINER="ioedream-mysql"
MYSQL_ROOT_PASSWORD="123456"
MYSQL_DATABASE="ioedream"
MIGRATION_SCRIPT="./docker/atlas/migrations/20251230140000_ioedream_tenant_isolation.sql"

echo "===================================================="
echo "  ioedream Tenant Isolation Migration"
echo "===================================================="
echo "Database: $MYSQL_DATABASE"
echo "Container: $MYSQL_CONTAINER"
echo ""

# 步骤1: 创建备份
echo "Step 1: Creating backup..."
BACKUP_DIR="./backup/ioedream_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"
BACKUP_FILE="$BACKUP_DIR/ioedream_backup.sql"

docker exec $MYSQL_CONTAINER mysqldump -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE > "$BACKUP_FILE" 2>/dev/null
if [ $? -eq 0 ]; then
    echo "✓ Backup completed: $BACKUP_FILE"
else
    echo "✗ Backup failed!"
    exit 1
fi
echo ""

# 步骤2: 执行迁移
echo "Step 2: Executing migration..."
sleep 2

docker exec -i $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE < "$MIGRATION_SCRIPT" 2>&1 | grep -v "Warning"

if [ ${PIPESTATUS[0]} -eq 0 ]; then
    echo "✓ Migration executed successfully"
else
    echo "✗ Migration failed!"
    echo "You can restore from backup: $BACKUP_FILE"
    exit 1
fi
echo ""

# 步骤3: 验证结果
echo "Step 3: Verifying migration..."
echo ""

# 检查tenants表
TENANTS_TABLE=$(docker exec $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE -se "SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='$MYSQL_DATABASE' AND TABLE_NAME='tenants';" 2>/dev/null)
if [ "$TENANTS_TABLE" -eq "1" ]; then
    echo "✓ tenants table exists"
else
    echo "✗ tenants table not found!"
    exit 1
fi

# 检查user_tenants表
USER_TENANTS_TABLE=$(docker exec $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE -se "SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='$MYSQL_DATABASE' AND TABLE_NAME='user_tenants';" 2>/dev/null)
if [ "$USER_TENANTS_TABLE" -eq "1" ]; then
    echo "✓ user_tenants table exists"
else
    echo "✗ user_tenants table not found!"
    exit 1
fi

# 检查tenant_id列
TENANT_COLUMN=$(docker exec $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE -se "SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA='$MYSQL_DATABASE' AND TABLE_NAME='t_common_user' AND COLUMN_NAME='tenant_id';" 2>/dev/null)
if [ "$TENANT_COLUMN" -eq "1" ]; then
    echo "✓ t_common_user.tenant_id column exists"
else
    echo "✗ t_common_user.tenant_id column not found!"
    exit 1
fi

# 统计数据
echo ""
echo "=== Statistics ==="
TOTAL_USERS=$(docker exec $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE -se "SELECT COUNT(*) FROM t_common_user WHERE deleted_flag = 0;" 2>/dev/null)
USERS_WITH_TENANT=$(docker exec $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE -se "SELECT COUNT(*) FROM t_common_user WHERE tenant_id IS NOT NULL AND deleted_flag = 0;" 2>/dev/null)
TOTAL_TENANTS=$(docker exec $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE -se "SELECT COUNT(*) FROM tenants;" 2>/dev/null)
INDIVIDUAL_TENANTS=$(docker exec $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE -se "SELECT COUNT(*) FROM tenants WHERE tenant_type='individual';" 2>/dev/null)
USER_TENANT_RELATIONS=$(docker exec $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE -se "SELECT COUNT(*) FROM user_tenants;" 2>/dev/null)

echo "Total Users: $TOTAL_USERS"
echo "Users with Tenant ID: $USERS_WITH_TENANT"
echo "Total Tenants: $TOTAL_TENANTS"
echo "Individual Tenants: $INDIVIDUAL_TENANTS"
echo "User-Tenant Relations: $USER_TENANT_RELATIONS"
echo ""

echo "===================================================="
echo "✓ Migration completed successfully!"
echo "===================================================="
echo "Backup location: $BACKUP_DIR"
echo ""

# 显示租户详情
echo "=== Tenant Details ==="
docker exec $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE -e "
SELECT
    t.tenant_id,
    t.tenant_name,
    t.tenant_type,
    t.status,
    u.username as owner_username
FROM tenants t
LEFT JOIN user_tenants ut ON t.tenant_id = ut.tenant_id AND ut.role = 'owner'
LEFT JOIN t_common_user u ON ut.user_id = u.user_id
ORDER BY t.create_time;
" 2>&1 | grep -v "Warning"
echo ""
