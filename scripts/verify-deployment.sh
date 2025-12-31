#!/bin/bash
# 验证Loki和备份系统部署完整性

echo "=========================================="
echo "ZKER 部署验证工具 v1.0.0"
echo "=========================================="
echo ""

# 颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

check_file() {
    if [ -f "$1" ]; then
        echo -e "${GREEN}✅${NC} $1"
        return 0
    else
        echo -e "${RED}❌${NC} $1"
        return 1
    fi
}

echo "1. Loki日志系统配置"
echo "-------------------"
check_file "docker/volumes/monitoring/loki/loki-config.yml"
check_file "docker/volumes/monitoring/promtail/promtail-config.yml"
check_file "docker/volumes/monitoring/grafana/provisioning/datasources/loki.yml"
echo ""

echo "2. 数据库备份脚本"
echo "-------------------"
check_file "scripts/backup/backup-mysql.sh"
check_file "scripts/backup/backup-redis.sh"
check_file "scripts/backup/setup-backup-cron.sh"
check_file "scripts/backup/verify-backup.sh"
check_file "scripts/backup/test-restore.sh"
echo ""

echo "3. Kubernetes部署文件"
echo "-------------------"
check_file "deploy/k8s/cronjob/mysql-backup.yaml"
check_file "deploy/k8s/cronjob/redis-backup.yaml"
echo ""

echo "4. 快速部署脚本"
echo "-------------------"
check_file "scripts/deploy-loki.sh"
check_file "scripts/setup-backup.sh"
echo ""

echo "5. 文档"
echo "-------------------"
check_file "scripts/backup/README.md"
check_file "LOKI_AND_BACKUP_DEPLOYMENT_SUMMARY.md"
check_file "LOKI_BACKUP_QUICK_REFERENCE.md"
echo ""

echo "=========================================="
echo "验证完成！"
echo "=========================================="
