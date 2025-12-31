#!/bin/bash
# =============================================================================
# 配置自动备份定时任务
# 版本: v1.0.0
# 更新: 2025-01-03
# =============================================================================

set -euo pipefail

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKUP_DIR="${SCRIPT_DIR}"
CRON_FILE="/etc/cron.d/zker-backup"
LOG_DIR="/var/log/backup"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[$(date '+%Y-%m-%d %H:%M:%S')]${NC} $1"
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
        exit 1
    fi
}

# 创建日志目录
create_log_dir() {
    log "创建日志目录: ${LOG_DIR}"
    mkdir -p "${LOG_DIR}"
    chmod 755 "${LOG_DIR}"
}

# 设置MySQL备份任务
setup_mysql_cron() {
    log "配置MySQL自动备份..."

    local mysql_backup_script="${BACKUP_DIR}/backup-mysql.sh"

    if [ ! -f "${mysql_backup_script}" ]; then
        error "MySQL备份脚本不存在: ${mysql_backup_script}"
        return 1
    fi

    chmod +x "${mysql_backup_script}"

    # 每天凌晨2点执行全量备份
    echo "0 2 * * * root ${mysql_backup_script} BACKUP_TYPE=full >> ${LOG_DIR}/mysql-backup.log 2>&1" >> "${CRON_FILE}"

    # 每6小时执行增量备份
    echo "0 */6 * * * root ${mysql_backup_script} BACKUP_TYPE=incremental >> ${LOG_DIR}/mysql-backup.log 2>&1" >> "${CRON_FILE}"

    log "MySQL备份任务已配置"
    log "  - 全量备份: 每天 02:00"
    log "  - 增量备份: 每6小时 (00:00, 06:00, 12:00, 18:00)"
}

# 设置Redis备份任务
setup_redis_cron() {
    log "配置Redis自动备份..."

    local redis_backup_script="${BACKUP_DIR}/backup-redis.sh"

    if [ ! -f "${redis_backup_script}" ]; then
        error "Redis备份脚本不存在: ${redis_backup_script}"
        return 1
    fi

    chmod +x "${redis_backup_script}"

    # 每6小时执行RDB备份
    echo "0 */6 * * * root ${redis_backup_script} BACKUP_TYPE=rdb >> ${LOG_DIR}/redis-backup.log 2>&1" >> "${CRON_FILE}"

    # 每天凌晨3点执行数据导出
    echo "0 3 * * * root ${redis_backup_script} BACKUP_TYPE=export >> ${LOG_DIR}/redis-backup.log 2>&1" >> "${CRON_FILE}"

    log "Redis备份任务已配置"
    log "  - RDB备份: 每6小时"
    log "  - 数据导出: 每天 03:00"
}

# 设置备份验证任务
setup_verify_cron() {
    log "配置备份验证任务..."

    local verify_script="${BACKUP_DIR}/verify-backup.sh"

    if [ ! -f "${verify_script}" ]; then
        warn "备份验证脚本不存在，跳过"
        return 0
    fi

    chmod +x "${verify_script}"

    # 每天凌晨4点验证备份
    echo "0 4 * * * root ${verify_script} >> ${LOG_DIR}/verify-backup.log 2>&1" >> "${CRON_FILE}"

    log "备份验证任务已配置: 每天 04:00"
}

# 创建Crontab
create_cron_file() {
    log "创建Crontab文件: ${CRON_FILE}"

    cat > "${CRON_FILE}" << 'CRON_HEADER'
# =============================================================================
# ZKER 自动备份定时任务
# =============================================================================
# 说明:
#   - MySQL全量备份: 每天 02:00
#   - MySQL增量备份: 每6小时
#   - Redis RDB备份: 每6小时
#   - Redis数据导出: 每天 03:00
#   - 备份验证: 每天 04:00
#
# 日志位置: /var/log/backup/
# =============================================================================

CRON_HEADER

    # 添加环境变量
    echo "" >> "${CRON_FILE}"
    echo "# 环境变量" >> "${CRON_FILE}"
    echo "SHELL=/bin/bash" >> "${CRON_FILE}"
    echo "PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin" >> "${CRON_FILE}"
    echo "" >> "${CRON_FILE}"
}

# 重启Cron服务
restart_cron() {
    log "重启Cron服务..."

    if command -v systemctl &> /dev/null; then
        systemctl restart cron || systemctl restart crond
        systemctl enable cron || systemctl enable crond
    elif command -v service &> /dev/null; then
        service cron restart || service crond restart
    else
        warn "无法重启Cron服务，请手动重启"
    fi
}

# 显示Crontab
show_crontab() {
    log "当前Crontab配置:"
    echo ""
    cat "${CRON_FILE}"
    echo ""
}

# 测试备份脚本
test_backup_scripts() {
    log "测试备份脚本..."

    local mysql_script="${BACKUP_DIR}/backup-mysql.sh"
    local redis_script="${BACKUP_DIR}/backup-redis.sh"

    if [ -f "${mysql_script}" ]; then
        log "测试MySQL备份脚本..."
        bash -n "${mysql_script}" && log "✅ MySQL脚本语法正确" || error "❌ MySQL脚本语法错误"
    fi

    if [ -f "${redis_script}" ]; then
        log "测试Redis备份脚本..."
        bash -n "${redis_script}" && log "✅ Redis脚本语法正确" || error "❌ Redis脚本语法错误"
    fi
}

# 显示使用帮助
show_help() {
    cat << EOF
用法: $0 [选项]

选项:
    -h, --help          显示帮助信息
    -s, --show          显示当前Crontab配置
    -t, --test          测试备份脚本
    -r, --remove        移除所有备份定时任务

示例:
    $0                  # 配置所有备份定时任务
    $0 --show           # 显示当前配置
    $0 --test           # 测试脚本
    $0 --remove         # 移除配置

EOF
}

# 移除Cron配置
remove_cron() {
    log "移除备份定时任务..."

    if [ -f "${CRON_FILE}" ]; then
        rm -f "${CRON_FILE}"
        log "已删除: ${CRON_FILE}"
    else
        log "Cron文件不存在"
    fi

    restart_cron
    log "定时任务已移除"
}

# 主函数
main() {
    log "=========================================="
    log "ZKER 自动备份配置工具 v1.0.0"
    log "=========================================="

    # 解析参数
    case "${1:-}" in
        -h|--help)
            show_help
            exit 0
            ;;
        -s|--show)
            show_crontab
            exit 0
            ;;
        -t|--test)
            test_backup_scripts
            exit 0
            ;;
        -r|--remove)
            check_root
            remove_cron
            exit 0
            ;;
        *)
            # 正常配置流程
            ;;
    esac

    # 检查root权限
    check_root

    # 创建日志目录
    create_log_dir

    # 创建Crontab文件
    create_cron_file

    # 配置各项备份任务
    setup_mysql_cron
    setup_redis_cron
    setup_verify_cron

    # 测试脚本
    test_backup_scripts

    # 重启Cron服务
    restart_cron

    # 显示配置
    show_crontab

    log "=========================================="
    log "✅ 自动备份配置完成！"
    log ""
    log "日志目录: ${LOG_DIR}"
    log "Crontab文件: ${CRON_FILE}"
    log ""
    log "查看备份日志:"
    log "  tail -f ${LOG_DIR}/mysql-backup.log"
    log "  tail -f ${LOG_DIR}/redis-backup.log"
    log "=========================================="
}

# 执行主函数
main "$@"
