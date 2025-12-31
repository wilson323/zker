#!/bin/bash
# =============================================================================
# MySQL自动备份脚本
# 版本: v1.0.0
# 更新: 2025-01-03
# 功能: 全量备份 + 增量备份 + 自动清理 + 验证
# =============================================================================

set -euo pipefail

# ==================== 配置区域 ====================
BACKUP_ROOT="/backup/mysql"
BACKUP_LOG="/var/log/mysql-backup.log"
RETENTION_DAYS=30

# MySQL连接配置
MYSQL_HOST="${MYSQL_HOST:-localhost}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-root}"
MYSQL_DATABASE="${MYSQL_DATABASE:---all-databases}"

# 备份类型 (full/incremental)
BACKUP_TYPE="${BACKUP_TYPE:-full}"

# 通知配置
NOTIFICATION_WEBHOOK="${NOTIFICATION_WEBHOOK:-}"

# ==================== 日志函数 ====================
log() {
    local level=$1
    shift
    local message="$@"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo "[${timestamp}] [${level}] ${message}" | tee -a "${BACKUP_LOG}"
}

# ==================== 前置检查 ====================
pre_check() {
    log "INFO" "开始前置检查..."

    # 检查MySQL连接
    if ! command -v mysql &> /dev/null; then
        log "ERROR" "mysql命令未安装"
        exit 1
    fi

    # 检查mysqldump
    if ! command -v mysqldump &> /dev/null; then
        log "ERROR" "mysqldump命令未安装"
        exit 1
    fi

    # 测试MySQL连接
    if ! mysql -h"${MYSQL_HOST}" -P"${MYSQL_PORT}" -u"${MYSQL_USER}" -p"${MYSQL_PASSWORD}" -e "SELECT 1" &> /dev/null; then
        log "ERROR" "MySQL连接失败: ${MYSQL_HOST}:${MYSQL_PORT}"
        exit 1
    fi

    # 创建备份目录
    mkdir -p "${BACKUP_ROOT}/full"
    mkdir -p "${BACKUP_ROOT}/incremental"
    mkdir -p "${BACKUP_ROOT}/logs"

    log "INFO" "前置检查完成"
}

# ==================== 全量备份 ====================
full_backup() {
    log "INFO" "开始全量备份..."
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local backup_file="${BACKUP_ROOT}/full/mysql-full-${timestamp}.sql.gz"

    # 执行备份
    log "INFO" "备份文件: ${backup_file}"

    if mysqldump \
        -h"${MYSQL_HOST}" \
        -P"${MYSQL_PORT}" \
        -u"${MYSQL_USER}" \
        -p"${MYSQL_PASSWORD}" \
        ${MYSQL_DATABASE} \
        --single-transaction \
        --quick \
        --lock-tables=false \
        --routines \
        --triggers \
        --events \
        --hex-blob \
        --set-gtid-purged=OFF \
        --master-data=2 \
        --flush-logs \
        2>&1 | tee -a "${BACKUP_LOG}" \
        | gzip > "${backup_file}"; then

        local backup_size=$(du -h "${backup_file}" | cut -f1)
        log "INFO" "全量备份成功: ${backup_file} (${backup_size})"

        # 记录备份元数据
        echo "${timestamp}|${backup_file}|${backup_size}|full" >> "${BACKUP_ROOT}/logs/backups.log"

        return 0
    else
        log "ERROR" "全量备份失败"
        return 1
    fi
}

# ==================== 增量备份 (基于binlog) ====================
incremental_backup() {
    log "INFO" "开始增量备份..."
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local backup_file="${BACKUP_ROOT}/incremental/mysql-incremental-${timestamp}.sql.gz"

    # 获取当前binlog位置
    local binlog_info=$(mysql -h"${MYSQL_HOST}" -P"${MYSQL_PORT}" -u"${MYSQL_USER}" -p"${MYSQL_PASSWORD}" -e "SHOW MASTER STATUS\G" 2>&1)
    local current_log=$(echo "${binlog_info}" | grep "File:" | awk '{print $2}')
    local current_pos=$(echo "${binlog_info}" | grep "Position:" | awk '{print $2}')

    log "INFO" "当前Binlog: ${current_log}, Position: ${current_pos}"

    # 导出binlog
    if mysqlbinlog \
        --read-from-remote-server \
        --host="${MYSQL_HOST}" \
        --port="${MYSQL_PORT}" \
        --user="${MYSQL_USER}" \
        --password="${MYSQL_PASSWORD}" \
        --raw \
        --stop-never \
        --result-file="${BACKUP_ROOT}/incremental/" \
        ${current_log} 2>&1 | tee -a "${BACKUP_LOG}"; then

        log "INFO" "增量备份成功"

        # 压缩binlog文件
        cd "${BACKUP_ROOT}/incremental"
        for binlog in mysql-bin.*; do
            if [ -f "${binlog}" ]; then
                gzip "${binlog}"
                log "INFO" "压缩: ${binlog}.gz"
            fi
        done

        return 0
    else
        log "ERROR" "增量备份失败"
        return 1
    fi
}

# ==================== 清理旧备份 ====================
cleanup_old_backups() {
    log "INFO" "清理${RETENTION_DAYS}天前的旧备份..."

    local deleted_count=0

    # 清理全量备份
    while IFS= read -r -d '' old_file; do
        rm -f "${old_file}"
        log "INFO" "删除: ${old_file}"
        ((deleted_count++))
    done < <(find "${BACKUP_ROOT}/full" -name "mysql-full-*.sql.gz" -mtime +${RETENTION_DAYS} -print0)

    # 清理增量备份
    while IFS= read -r -d '' old_file; do
        rm -f "${old_file}"
        log "INFO" "删除: ${old_file}"
        ((deleted_count++))
    done < <(find "${BACKUP_ROOT}/incremental" -name "*.gz" -mtime +${RETENTION_DAYS} -print0)

    log "INFO" "清理完成，删除${deleted_count}个文件"
}

# ==================== 验证备份 ====================
verify_backup() {
    local backup_file=$1

    log "INFO" "验证备份: ${backup_file}"

    # 检查文件是否存在
    if [ ! -f "${backup_file}" ]; then
        log "ERROR" "备份文件不存在: ${backup_file}"
        return 1
    fi

    # 检查文件大小
    local file_size=$(stat -f%z "${backup_file}" 2>/dev/null || stat -c%s "${backup_file}" 2>/dev/null)
    if [ "${file_size}" -lt 1024 ]; then
        log "ERROR" "备份文件过小: ${file_size} bytes"
        return 1
    fi

    # 验证gzip完整性
    if ! gzip -t "${backup_file}" 2>/dev/null; then
        log "ERROR" "备份文件损坏"
        return 1
    fi

    # 尝试解压并检查SQL内容
    local table_count=$(gunzip -c "${backup_file}" | grep -c "CREATE TABLE" || true)
    if [ "${table_count}" -eq 0 ]; then
        log "WARN" "备份文件中未发现表定义"
    fi

    log "INFO" "验证通过: ${backup_file} (${file_size} bytes, ${table_count} tables)"
    return 0
}

# ==================== 发送通知 ====================
send_notification() {
    local status=$1
    local message=$2

    if [ -n "${NOTIFICATION_WEBHOOK}" ]; then
        local payload=$(cat <<EOF
{
    "status": "${status}",
    "message": "${message}",
    "timestamp": "$(date -Iseconds)",
    "host": "$(hostname)"
}
EOF
)
        curl -X POST "${NOTIFICATION_WEBHOOK}" \
            -H "Content-Type: application/json" \
            -d "${payload}" \
            --silent --show-error \
            2>&1 | tee -a "${BACKUP_LOG}"
    fi
}

# ==================== 主流程 ====================
main() {
    log "INFO" "=========================================="
    log "INFO" "MySQL备份任务开始"
    log "INFO" "备份类型: ${BACKUP_TYPE}"
    log "INFO" "=========================================="

    # 前置检查
    pre_check || exit 1

    # 执行备份
    local backup_status=0
    local backup_file=""

    if [ "${BACKUP_TYPE}" = "full" ]; then
        full_backup || backup_status=1
        backup_file="${BACKUP_ROOT}/full/mysql-full-$(date +%Y%m%d-%H%M%S).sql.gz"
    elif [ "${BACKUP_TYPE}" = "incremental" ]; then
        incremental_backup || backup_status=1
    fi

    # 验证备份
    if [ ${backup_status} -eq 0 ] && [ -n "${backup_file}" ]; then
        verify_backup "${backup_file}" || backup_status=1
    fi

    # 清理旧备份
    if [ ${backup_status} -eq 0 ]; then
        cleanup_old_backups
    fi

    # 发送通知
    if [ ${backup_status} -eq 0 ]; then
        log "INFO" "备份任务成功完成"
        send_notification "success" "MySQL备份成功"
    else
        log "ERROR" "备份任务失败"
        send_notification "failure" "MySQL备份失败"
        exit 1
    fi

    log "INFO" "=========================================="
}

# 执行主流程
main "$@"
