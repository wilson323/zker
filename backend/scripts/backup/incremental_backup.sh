#!/bin/bash
# ============================================================
# ZKER 增量备份脚本
# ============================================================
# 功能: 执行MySQL增量备份（Percona XtraBackup）
# 执行时间: 每日 02:00
# 保留策略: 本地7天，异地30天
# ============================================================

set -euo pipefail

# ============================================================
# 配置参数
# ============================================================

BACKUP_DIR="/data/backups/incremental"
FULL_BACKUP_DIR="/data/backups/full"
LOG_FILE="/var/log/backup.log"
LOCK_FILE="/var/lock/incremental_backup.lock"
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-}"
S3_BUCKET="s3://zker-backups/incremental"
DR_S3_BUCKET="s3://zker-backups-dr/incremental"
RETENTION_DAYS=7

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
            log "ERROR" "Another backup process is running (PID: ${pid})"
            exit 1
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
        log "ERROR" "Incremental backup failed with exit code ${exit_code}"
        send_alert "Incremental backup failed" "Exit code: ${exit_code}"
    fi
    remove_lock
    exit ${exit_code}
}

send_alert() {
    local subject=$1
    local message=$2
    echo "${message}" | mail -s "[ZKER Backup Alert] ${subject}" ops-team@zker.com
    curl -X POST -H 'Content-type: application/json' \
        --data "{\"text\":\"${subject}: ${message}\"}" \
        ${SLACK_WEBHOOK_URL}
}

# ============================================================
# 主流程
# ============================================================

main() {
    trap cleanup EXIT INT TERM

    log "INFO" "==========================================="
    log "INFO" "Starting incremental backup process"
    log "INFO" "==========================================="

    # 检查锁文件
    check_lock
    create_lock

    # 获取日期信息
    local date=$(date +%Y%m%d)
    local day_of_week=$(date +%u) # 1=Monday, 7=Sunday
    local timestamp=$(date +%s)
    local backup_file="incr_backup_${date}.tar.gz"
    local backup_path="${BACKUP_DIR}/${date}"
    local temp_backup_path="${BACKUP_DIR}/temp_${date}"

    # 1. 查找最新的全量备份
    log "INFO" "Finding latest full backup..."
    local latest_full=$(ls -t ${FULL_BACKUP_DIR}/full_backup_*.tar.gz 2>/dev/null | head -1)

    if [ -z "${latest_full}" ]; then
        log "ERROR" "No full backup found. Please run full backup first."
        exit 1
    fi

    log "INFO" "Base full backup: ${latest_full}"

    # 2. 解压全量备份作为基础
    log "INFO" "Extracting base full backup..."
    local base_backup_path="${BACKUP_DIR}/base_full"
    rm -rf ${base_backup_path}
    mkdir -p ${base_backup_path}

    tar -xzf ${latest_full} -C ${base_backup_path}

    if [ $? -ne 0 ]; then
        log "ERROR" "Failed to extract full backup"
        exit 1
    fi

    # 3. 执行增量备份
    log "INFO" "Executing XtraBackup incremental backup..."
    local start_time=$(date +%s)

    xtrabackup --backup \
        --target-dir=${temp_backup_path} \
        --incremental-basedir=${base_backup_path} \
        --user=root \
        --password=${MYSQL_ROOT_PASSWORD} \
        --parallel=4 \
        --compress \
        --compress-threads=4 \
        2>&1 | tee -a ${LOG_FILE}

    if [ $? -ne 0 ]; then
        log "ERROR" "XtraBackup incremental backup failed"
        exit 1
    fi

    local backup_end_time=$(date +%s)
    local backup_duration=$((backup_end_time - start_time))
    log "INFO" "XtraBackup incremental completed in ${backup_duration} seconds"

    # 4. 准备增量备份
    log "INFO" "Preparing incremental backup..."
    xtrabackup --prepare \
        --target-dir=${base_backup_path} \
        --apply-log-only \
        2>&1 | tee -a ${LOG_FILE}

    xtrabackup --prepare \
        --target-dir=${base_backup_path} \
        --incremental-dir=${temp_backup_path} \
        2>&1 | tee -a ${LOG_FILE}

    if [ $? -ne 0 ]; then
        log "ERROR" "XtraBackup prepare failed"
        exit 1
    fi

    # 5. 打包增量备份
    log "INFO" "Compressing incremental backup to ${backup_file}..."
    cd ${temp_backup_path}
    tar -czf ${BACKUP_DIR}/${backup_file} .
    cd -

    if [ $? -ne 0 ]; then
        log "ERROR" "Incremental backup compression failed"
        exit 1
    fi

    # 6. 计算校验和
    log "INFO" "Calculating SHA256 checksum..."
    sha256sum ${BACKUP_DIR}/${backup_file} > ${BACKUP_DIR}/${backup_file}.sha256

    if [ $? -ne 0 ]; then
        log "ERROR" "Checksum calculation failed"
        exit 1
    fi

    local checksum=$(cat ${BACKUP_DIR}/${backup_file}.sha256 | awk '{print $1}')
    log "INFO" "SHA256 checksum: ${checksum}"

    # 7. 获取文件大小
    local file_size=$(stat -f%z ${BACKUP_DIR}/${backup_file} 2>/dev/null || stat -c%s ${BACKUP_DIR}/${backup_file})
    local file_size_mb=$((file_size / 1024 / 1024))
    log "INFO" "Incremental backup file size: ${file_size_mb} MB"

    # 8. 计算压缩率
    local base_size=$(tar -tzf ${latest_full} | wc -c)
    local compression_ratio=$(echo "scale=2; ${base_size} / ${file_size}" | bc)
    log "INFO" "Compression ratio: ${compression_ratio}x"

    # 9. 上传到S3主区域
    log "INFO" "Uploading to S3 primary region..."
    aws s3 cp ${BACKUP_DIR}/${backup_file} \
        ${S3_BUCKET}/ \
        --storage-class STANDARD_IA \
        --metadata "checksum=${checksum},timestamp=${timestamp}" \
        2>&1 | tee -a ${LOG_FILE}

    if [ $? -ne 0 ]; then
        log "ERROR" "S3 upload to primary region failed"
        exit 1
    fi

    # 10. 上传校验和文件
    aws s3 cp ${BACKUP_DIR}/${backup_file}.sha256 ${S3_BUCKET}/

    # 11. 上传到S3灾备区域
    log "INFO" "Uploading to S3 DR region..."
    aws s3 cp ${BACKUP_DIR}/${backup_file} \
        ${DR_S3_BUCKET}/ \
        --storage-class STANDARD_IA \
        --metadata "checksum=${checksum},timestamp=${timestamp}" \
        2>&1 | tee -a ${LOG_FILE}

    if [ $? -ne 0 ]; then
        log "WARN" "S3 upload to DR region failed (non-fatal)"
    fi

    # 12. 清理临时文件
    log "INFO" "Cleaning up temporary files..."
    rm -rf ${temp_backup_path}
    rm -rf ${base_backup_path}

    # 13. 清理旧备份（保留7天）
    log "INFO" "Cleaning up old backups (older than ${RETENTION_DAYS} days)..."
    find ${BACKUP_DIR} -name "incr_backup_*.tar.gz" -mtime +${RETENTION_DAYS} -exec rm -f {} \;
    find ${BACKUP_DIR} -name "incr_backup_*.sha256" -mtime +${RETENTION_DAYS} -exec rm -f {} \;

    # 14. 清理S3旧备份
    log "INFO" "Cleaning up S3 old backups..."
    aws s3 ls ${S3_BUCKET}/ | while read -r line; do
        local file_date=$(echo ${line} | awk '{print $4}')
        local file_time=$(echo ${line} | awk '{print $1" "$2}')
        local file_timestamp=$(date -d "${file_time}" +%s 2>/dev/null || echo 0)
        local current_timestamp=$(date +%s)
        local file_age=$(( (current_timestamp - file_timestamp) / 86400 ))

        if [ ${file_age} -gt ${RETENTION_DAYS} ]; then
            log "INFO" "Deleting old S3 backup: ${file_date}"
            aws s3 rm ${S3_BUCKET}/${file_date}
        fi
    done

    # 15. 记录备份成功
    local end_time=$(date +%s)
    local total_duration=$((end_time - start_time))
    log "INFO" "==========================================="
    log "INFO" "Incremental backup completed successfully"
    log "INFO" "Backup file: ${backup_file}"
    log "INFO" "File size: ${file_size_mb} MB"
    log "INFO" "Checksum: ${checksum}"
    log "INFO" "Total duration: ${total_duration} seconds"
    log "INFO" "==========================================="

    # 16. 推送Prometheus指标
    if command -v curl &> /dev/null; then
        curl -X POST http://localhost:9091/metrics/job/backup \
            -d "backup_success{type=\"incremental\",database=\"zker\"} 1" \
            -d "backup_size_bytes{type=\"incremental\",database=\"zker\"} ${file_size}" \
            -d "backup_duration_seconds{type=\"incremental\",database=\"zker\"} ${total_duration}"
    fi

    return 0
}

# 执行主流程
main "$@"
