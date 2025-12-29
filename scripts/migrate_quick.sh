#!/bin/bash

# =====================================================
# 租户隔离数据库迁移 - 快速执行版本
# 版本: v1.0
# =====================================================

set -e

MYSQL_CONTAINER="ioedream-mysql"
MYSQL_ROOT_PASSWORD="123456"
MYSQL_DATABASE="ioedream"
MIGRATION_SCRIPT="./docker/atlas/migrations/20251230120000_add_user_tenant_isolation.sql"

echo "===================================================="
echo "  Tenant Isolation Database Migration"
echo "===================================================="
echo "Database: $MYSQL_DATABASE"
echo "Container: $MYSQL_CONTAINER"
echo ""

# 步骤1: 创建备份
echo "Step 1: Creating backup..."
mkdir -p ./backup/migration_$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="./backup/migration_$(date +%Y%m%d_%H%M%S)/ioedream_backup.sql"

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

docker exec -i $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE < "$MIGRATION_SCRIPT" 2>/dev/null

if [ $? -eq 0 ]; then
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

# 检查tenant_id列
TENANT_COLUMN=$(docker exec $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE -se "SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA='$MYSQL_DATABASE' AND TABLE_NAME='users' AND COLUMN_NAME='tenant_id';" 2>/dev/null)
if [ "$TENANT_COLUMN" -eq "1" ]; then
    echo "✓ users.tenant_id column exists"
else
    echo "✗ users.tenant_id column not found!"
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

# 统计数据
echo ""
echo "=== Statistics ==="
TOTAL_USERS=$(docker exec $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE -se "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL;" 2>/dev/null)
USERS_WITH_TENANT=$(docker exec $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE -se "SELECT COUNT(*) FROM users WHERE tenant_id IS NOT NULL AND deleted_at IS NULL;" 2>/dev/null)
TOTAL_TENANTS=$(docker exec $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE -se "SELECT COUNT(*) FROM tenants WHERE deleted_at IS NULL;" 2>/dev/null)
INDIVIDUAL_TENANTS=$(docker exec $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE -se "SELECT COUNT(*) FROM tenants WHERE tenant_type='individual' AND deleted_at IS NULL;" 2>/dev/null)
USER_TENANT_RELATIONS=$(docker exec $MYSQL_CONTAINER mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE -se "SELECT COUNT(*) FROM user_tenants WHERE deleted_at IS NULL;" 2>/dev/null)

echo "Total Users: $TOTAL_USERS"
echo "Users with Tenant ID: $USERS_WITH_TENANT"
echo "Total Tenants: $TOTAL_TENANTS"
echo "Individual Tenants: $INDIVIDUAL_TENANTS"
echo "User-Tenant Relations: $USER_TENANT_RELATIONS"
echo ""

echo "===================================================="
echo "✓ Migration completed successfully!"
echo "===================================================="
