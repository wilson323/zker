#!/bin/bash
# ============================================================
# ZKER Binlog备份脚本
# ============================================================
# 功能: 实时备份MySQL Binlog（用于PITR）
# 执行时间: 每5分钟
# 保留策略: 本地3天，异地7天
# ============================================================

set -euo pipefail

# ============================================================
# 配置参数
# ============================================================

BINLOG_DIR="/var/lib/mysql"
BACKUP_DIR="/data/backups/binlog"
LOG_FILE="/var/log/backup.log"
LOCK_FILE="/var/lock/binlog_backup.lock"
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-}"
S3_BUCKET="s3://zker-backups/binlog"
RETENTION_DAYS=3

# ============================================================
# 工具函数
# ============================================================

log() {
    local level=$1
    shift
    local message="$@"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo "[${timestamp}] [${level}] ${message}" | tee -a ${LOG_FILE}
}

check_lock() {
    if [ -f "${LOCK_FILE}" ]; then
        local pid=$(cat ${LOCK_FILE})
        if ps -p ${pid} > /dev/null 2>&1; then
            log "WARN" "Another backup process is running (PID: ${pid}), skipping"
            exit 0
        else
            log "WARN" "Removing stale lock file"
            rm -f ${LOCK_FILE}
        fi
    fi
}

create_lock() {
    echo $$ > ${LOCK_FILE}
}

remove_lock() {
    rm -f ${LOCK_FILE}
}

cleanup() {
    local exit_code=$?
    if [ ${exit_code} -ne 0 ]; then
        log "ERROR" "Binlog backup failed with exit code ${exit_code}"
        send_alert "Binlog backup failed" "Exit code: ${exit_code}"
    fi
    remove_lock
    exit ${exit_code}
}

send_alert() {
    local subject=$1
    local message=$2
    echo "${message}" | mail -s "[ZKER Backup Alert] ${subject}" ops-team@zker.com
}

# ============================================================
# 主流程
# ============================================================

main() {
    trap cleanup EXIT INT TERM

    log "DEBUG" "Starting binlog backup process"

    # 检查锁文件
    check_lock
    create_lock

    # 1. 刷新binlog，切换到新文件
    log "DEBUG" "Flushing binary logs..."
    mysql -u root -p${MYSQL_ROOT_PASSWORD} -e "FLUSH BINARY LOGS;" 2>&1 | tee -a ${LOG_FILE}

    if [ $? -ne 0 ]; then
        log "ERROR" "Failed to flush binary logs"
        exit 1
    fi

    # 2. 查找最新的binlog文件
    log "DEBUG" "Finding latest binlog file..."
    local latest_binlog=$(ls -t ${BINLOG_DIR}/mysql-bin.0* 2>/dev/null | head -1)

    if [ -z "${latest_binlog}" ]; then
        log "ERROR" "No binlog files found"
        exit 1
    fi

    local binlog_name=$(basename ${latest_binlog})
    log "DEBUG" "Latest binlog: ${binlog_name}"

    # 3. 检查是否已备份
    if [ -f "${BACKUP_DIR}/${binlog_name}.gz" ]; then
        log "DEBUG" "Binlog ${binlog_name} already backed up, skipping"
        exit 0
    fi

    # 4. 复制binlog文件
    log "DEBUG" "Copying binlog file..."
    cp ${latest_binlog} ${BACKUP_DIR}/

    if [ $? -ne 0 ]; then
        log "ERROR" "Failed to copy binlog file"
        exit 1
    fi

    # 5. 压缩binlog
    log "DEBUG" "Compressing binlog..."
    gzip ${BACKUP_DIR}/${binlog_name}

    if [ $? -ne 0 ]; then
        log "ERROR" "Failed to compress binlog"
        exit 1
    fi

    # 6. 计算校验和
    log "DEBUG" "Calculating checksum..."
    sha256sum ${BACKUP_DIR}/${binlog_name}.gz > ${BACKUP_DIR}/${binlog_name}.gz.sha256

    if [ $? -ne 0 ]; then
        log "ERROR" "Failed to calculate checksum"
        exit 1
    fi

    # 7. 获取文件大小
    local file_size=$(stat -f%z ${BACKUP_DIR}/${binlog_name}.gz 2>/dev/null || stat -c%s ${BACKUP_DIR}/${binlog_name}.gz)
    local file_size_kb=$((file_size / 1024))
    log "DEBUG" "Binlog file size: ${file_size_kb} KB"

    # 8. 上传到S3
    log "DEBUG" "Uploading to S3..."
    aws s3 cp ${BACKUP_DIR}/${binlog_name}.gz \
        ${S3_BUCKET}/ \
        --storage-class STANDARD \
        2>&1 | tee -a ${LOG_FILE}

    if [ $? -ne 0 ]; then
        log "ERROR" "Failed to upload to S3"
        exit 1
    fi

    # 9. 上传校验和文件
    aws s3 cp ${BACKUP_DIR}/${binlog_name}.gz.sha256 ${S3_BUCKET}/

    # 10. 清理旧binlog（保留3天）
    log "DEBUG" "Cleaning up old binlogs (older than ${RETENTION_DAYS} days)..."
    find ${BACKUP_DIR} -name "*.gz" -mtime +${RETENTION_DAYS} -exec rm -f {} \;
    find ${BACKUP_DIR} -name "*.gz.sha256" -mtime +${RETENTION_DAYS} -exec rm -f {} \;

    # 11. 清理S3旧binlog
    log "DEBUG" "Cleaning up S3 old binlogs..."
    aws s3 ls ${S3_BUCKET}/ | while read -r line; do
        local file_date=$(echo ${line} | awk '{print $4}')
        local file_time=$(echo ${line} | awk '{print $1" "$2}')
        local file_timestamp=$(date -d "${file_time}" +%s 2>/dev/null || echo 0)
        local current_timestamp=$(date +%s)
        local file_age=$(( (current_timestamp - file_timestamp) / 86400 ))

        if [ ${file_age} -gt ${RETENTION_DAYS} ]; then
            log "DEBUG" "Deleting old S3 binlog: ${file_date}"
            aws s3 rm ${S3_BUCKET}/${file_date}
        fi
    done

    log "DEBUG" "Binlog backup completed successfully: ${binlog_name}.gz (${file_size_kb} KB)"

    # 12. 推送Prometheus指标
    if command -v curl &> /dev/null; then
        curl -X POST http://localhost:9091/metrics/job/backup \
            -d "backup_success{type=\"binlog\",database=\"zker\"} 1" \
            -d "backup_size_bytes{type=\"binlog\",database=\"zker\"} ${file_size}"
    fi

    return 0
}

# 执行主流程
main "$@"
