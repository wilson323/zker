#!/bin/bash
# =============================================================================
# 数据库自动备份快速配置脚本
# 版本: v1.0.0
# 更新: 2025-01-03
# =============================================================================

set -euo pipefail

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查root权限
check_root() {
    if [ "$EUID" -ne 0 ]; then
        error "请使用root权限运行此脚本"
        echo "使用: sudo $0"
        exit 1
    fi
}

# 创建备份目录
create_backup_dirs() {
    log "创建备份目录..."

    mkdir -p /backup/mysql/full
    mkdir -p /backup/mysql/incremental
    mkdir -p /backup/mysql/logs
    mkdir -p /backup/redis/rdb
    mkdir -p /backup/redis/aof
    mkdir -p /backup/redis/logs
    mkdir -p /var/log/backup

    log "✅ 备份目录已创建"
}

# 复制脚本
copy_scripts() {
    log "复制备份脚本..."

    local script_dir="$(dirname "$0")/backup"

    if [ ! -d "${script_dir}" ]; then
        error "备份脚本目录不存在: ${script_dir}"
        exit 1
    fi

    # 设置执行权限
    chmod +x "${script_dir}/"*.sh

    # 复制到/usr/local/bin
    cp -f "${script_dir}/backup-mysql.sh" /usr/local/bin/zker-backup-mysql
    cp -f "${script_dir}/backup-redis.sh" /usr/local/bin/zker-backup-redis
    cp -f "${script_dir}/verify-backup.sh" /usr/local/bin/zker-verify-backup
    cp -f "${script_dir}/test-restore.sh" /usr/local/bin/zker-test-restore

    log "✅ 脚本已安装到 /usr/local/bin"
}

# 配置环境变量
configure_env() {
    log "配置环境变量..."

    local env_file="/etc/zker-backup.conf"

    if [ ! -f "${env_file}" ]; then
        cat > "${env_file}" << 'EOF'
# ZKER 数据库备份配置
# 更新日期: 2025-01-03

# MySQL配置
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=your_password_here
MYSQL_DATABASE=--all-databases

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# 备份配置
BACKUP_ROOT=/backup
RETENTION_DAYS_MYSQL=30
RETENTION_DAYS_REDIS=7

# 通知配置（可选）
NOTIFICATION_WEBHOOK=
EOF
        warn "请编辑配置文件: ${env_file}"
        warn "修改MySQL和Redis密码"
    fi

    log "✅ 配置文件已创建: ${env_file}"
}

# 配置定时任务
setup_cron() {
    log "配置定时任务..."

    local setup_script="$(dirname "$0")/backup/setup-backup-cron.sh"

    if [ -f "${setup_script}" ]; then
        bash "${setup_script}"
    else
        error "Cron配置脚本不存在: ${setup_script}"
        exit 1
    fi
}

# 测试备份
test_backup() {
    log "测试备份脚本..."

    log "测试MySQL备份（dry-run）..."
    /usr/local/bin/zker-backup-mysql 2>&1 | head -20 || true

    log "测试Redis备份（dry-run）..."
    /usr/local/bin/zker-backup-redis 2>&1 | head -20 || true

    log "✅ 备份脚本测试完成"
}

# 显示完成信息
show_info() {
    echo ""
    echo "=========================================="
    echo "🎉 数据库自动备份配置完成！"
    echo "=========================================="
    echo ""
    echo "备份脚本位置:"
    echo "  - /usr/local/bin/zker-backup-mysql"
    echo "  - /usr/local/bin/zker-backup-redis"
    echo "  - /usr/local/bin/zker-verify-backup"
    echo "  - /usr/local/bin/zker-test-restore"
    echo ""
    echo "配置文件:"
    echo "  - /etc/zker-backup.conf"
    echo ""
    echo "定时任务:"
    echo "  - MySQL全量备份: 每天 02:00"
    echo "  - MySQL增量备份: 每6小时"
    echo "  - Redis RDB备份: 每6小时"
    echo "  - 备份验证: 每天 04:00"
    echo ""
    echo "日志位置:"
    echo "  - /var/log/backup/mysql-backup.log"
    echo "  - /var/log/backup/redis-backup.log"
    echo "  - /var/log/backup/verify-backup.log"
    echo ""
    echo "下一步:"
    echo "  1. 编辑配置文件: vim /etc/zker-backup.conf"
    echo "  2. 手动测试备份: zker-backup-mysql BACKUP_TYPE=full"
    echo "  3. 查看定时任务: crontab -l"
    echo "  4. 验证备份: zker-verify-backup"
    echo ""
    echo "详细文档: scripts/backup/README.md"
    echo "=========================================="
}

# 主流程
main() {
    echo "=========================================="
    echo "ZKER 数据库备份配置工具 v1.0.0"
    echo "=========================================="
    echo ""

    check_root
    create_backup_dirs
    copy_scripts
    configure_env
    setup_cron
    test_backup
    show_info
}

main "$@"
