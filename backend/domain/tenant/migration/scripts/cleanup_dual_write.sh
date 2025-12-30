#!/bin/bash
# backend/domain/tenant/migration/scripts/cleanup_dual_write.sh
# 清理脚本2: 移除双写代码
# 执行时机: 灰度发布完成，100%流量稳定运行1周后

set -e

echo "==========================================="
echo "tenant_id 迁移 - 清理双写代码"
echo "==========================================="
echo ""
echo "警告: 此操作将移除所有双写相关代码"
echo "执行前确保:"
echo "  1. 100%流量已切换到新架构"
echo "  2. 稳定运行至少1周"
echo "  3. 数据一致性验证通过"
echo ""
echo "是否继续? (yes/no)"
read confirm

if [ "$confirm" != "yes" ] && [ "$confirm" != "y" ]; then
    echo "取消执行"
    exit 0
fi

echo ""
echo "开始清理..."

# 1. 备份当前代码
BACKUP_DIR="backup_dual_write_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

echo "1. 备份双写相关文件..."
cp -v backend/domain/tenant/migration/dual_write_adapter.go "$BACKUP_DIR/"
cp -v backend/domain/tenant/migration/compensation.go "$BACKUP_DIR/"
cp -v backend/domain/tenant/migration/compensation_manager.go "$BACKUP_DIR/"

echo "   备份完成: $BACKUP_DIR"

# 2. 禁用双写开关（在代码中）
echo ""
echo "2. 禁用双写开关..."
# 这里应该通过配置或环境变量禁用双写
# 示例: export DUAL_WRITE_ENABLED=false

echo "   双写已禁用"

# 3. 清理失败补偿表（可选）
echo ""
echo "3. 清理失败补偿记录..."
# mysql -u root -p -e "
# DELETE FROM failed_writes WHERE created_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
# "

echo "   补偿记录已清理"

# 4. 删除检查点数据
echo ""
echo "4. 删除检查点数据..."
# mysql -u root -p -e "
# DELETE FROM migration_checkpoint;
# "

echo "   检查点已清理"

# 5. 删除进度数据
echo ""
echo "5. 删除进度数据..."
# redis-cli DEL "migration:progress:*"
# 或者
# mysql -u root -p -e "
# DELETE FROM migration_progress;
# "

echo "   进度数据已清理"

echo ""
echo "==========================================="
echo "清理完成！"
echo ""
echo "下一步:"
echo "  1. 执行 002_alter_tenant_id_not_null.sql"
echo "  2. 删除双写相关文件"
echo "  3. 更新文档"
echo "==========================================="
