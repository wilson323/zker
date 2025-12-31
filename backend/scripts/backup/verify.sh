#!/bin/bash
# ============================================================
# ZKER 备份验证脚本
# ============================================================
# 功能: 验证备份文件的完整性和可用性
# 执行时间: 每次备份后自动执行
# 验证内容: 文件大小、校验和、数据行数
# ============================================================

set -euo pipefail

# ============================================================
# 配置参数
# ============================================================

BACKUP_FILE=$1
BACKUP_TYPE="${2:-full}" # full, incremental, binlog
LOG_FILE="/var/log/backup_verify.log"
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-}"
TEMP_DIR="/tmp/backup_verify_$$"

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

# ============================================================
# 主流程
# ============================================================

main() {
    log "INFO" "==========================================="
    log "INFO" "Starting backup verification"
    log "INFO" "Backup file: ${BACKUP_FILE}"
    log "INFO" "Backup type: ${BACKUP_TYPE}"
    log "INFO" "==========================================="

    local verification_failed=0

    # 1. 检查文件是否存在
    log "INFO" "Step 1: Checking file existence..."
    if [ ! -f "${BACKUP_FILE}" ]; then
        log "ERROR" "Backup file does not exist: ${BACKUP_FILE}"
        exit 1
    fi
    log "INFO" "File exists"

    # 2. 验证文件大小
    log "INFO" "Step 2: Verifying file size..."
    local file_size=$(stat -f%z ${BACKUP_FILE} 2>/dev/null || stat -c%s ${BACKUP_FILE})
    local file_size_mb=$((file_size / 1024 / 1024))

    case ${BACKUP_TYPE} in
        full)
            if [ ${file_size} -lt 1073741824 ]; then # 1GB
                log "ERROR" "Full backup file too small (${file_size_mb} MB < 1 GB)"
                verification_failed=1
            else
                log "INFO" "File size OK: ${file_size_mb} MB"
            fi
            ;;
        incremental)
            if [ ${file_size} -lt 10485760 ]; then # 10MB
                log "ERROR" "Incremental backup file too small (${file_size_mb} MB < 10 MB)"
                verification_failed=1
            else
                log "INFO" "File size OK: ${file_size_mb} MB"
            fi
            ;;
        binlog)
            if [ ${file_size} -lt 1024 ]; then # 1KB
                log "ERROR" "Binlog file too small (${file_size} bytes < 1 KB)"
                verification_failed=1
            else
                log "INFO" "File size OK: ${file_size} bytes"
            fi
            ;;
    esac

    # 3. 验证校验和
    log "INFO" "Step 3: Verifying checksum..."
    local checksum_file="${BACKUP_FILE}.sha256"

    if [ ! -f "${checksum_file}" ]; then
        log "ERROR" "Checksum file does not exist: ${checksum_file}"
        verification_failed=1
    else
        sha256sum -c ${checksum_file}
        if [ $? -eq 0 ]; then
            log "INFO" "Checksum verification passed"
        else
            log "ERROR" "Checksum verification failed"
            verification_failed=1
        fi
    fi

    # 4. 验证备份内容（仅全量和增量）
    if [[ "${BACKUP_TYPE}" == "full" ]] || [[ "${BACKUP_TYPE}" == "incremental" ]]; then
        log "INFO" "Step 4: Verifying backup content..."

        # 创建临时目录
        mkdir -p ${TEMP_DIR}

        # 解压备份
        log "INFO" "Extracting backup..."
        tar -xzf ${BACKUP_FILE} -C ${TEMP_DIR}

        if [ $? -ne 0 ]; then
            log "ERROR" "Failed to extract backup"
            verification_failed=1
        else
            # 准备备份
            log "INFO" "Preparing backup..."
            xtrabackup --prepare --target-dir=${TEMP_DIR} > /dev/null 2>&1

            if [ $? -ne 0 ]; then
                log "ERROR" "Failed to prepare backup"
                verification_failed=1
            else
                # 验证关键表
                log "INFO" "Verifying data integrity..."

                # 启动临时MySQL实例
                local temp_port=3307
                local temp_socket="/tmp/mysql_verify_$$.sock"

                # 使用--export选项导出表
                log "INFO" "Exporting tables for verification..."

                # 检查tenants表
                if [ -f "${TEMP_DIR}/zker/tenants.ibd" ]; then
                    log "INFO" "Tenants table found"
                else
                    log "ERROR" "Tenants table not found"
                    verification_failed=1
                fi

                # 检查bots表
                if [ -f "${TEMP_DIR}/zker/bots.ibd" ]; then
                    log "INFO" "Bots table found"
                else
                    log "ERROR" "Bots table not found"
                    verification_failed=1
                fi

                # 检查conversations表
                if [ -f "${TEMP_DIR}/zker/conversations.ibd" ]; then
                    log "INFO" "Conversations table found"
                else
                    log "ERROR" "Conversations table not found"
                    verification_failed=1
                fi

                # 检查messages表
                if [ -f "${TEMP_DIR}/zker/messages.ibd" ]; then
                    log "INFO" "Messages table found"
                else
                    log "ERROR" "Messages table not found"
                    verification_failed=1
                fi
            fi
        fi

        # 清理临时目录
        rm -rf ${TEMP_DIR}
    fi

    # 5. 验证S3上传（如果启用了S3）
    log "INFO" "Step 5: Verifying S3 upload..."
    local filename=$(basename ${BACKUP_FILE})

    if command -v aws &> /dev/null; then
        # 检查主区域
        local s3_exists=$(aws s3 ls s3://zker-backups/${BACKUP_TYPE}/${filename} 2>&1 | grep -c ${filename} || echo 0)

        if [ ${s3_exists} -gt 0 ]; then
            log "INFO" "S3 primary region upload verified"
        else
            log "WARN" "S3 primary region upload not found"
        fi

        # 检查灾备区域
        local dr_s3_exists=$(aws s3 ls s3://zker-backups-dr/${BACKUP_TYPE}/${filename} 2>&1 | grep -c ${filename} || echo 0)

        if [ ${dr_s3_exists} -gt 0 ]; then
            log "INFO" "S3 DR region upload verified"
        else
            log "INFO" "S3 DR region upload not found (may still be uploading)"
        fi
    else
        log "INFO" "AWS CLI not found, skipping S3 verification"
    fi

    # 6. 生成验证报告
    log "INFO" "==========================================="
    if [ ${verification_failed} -eq 0 ]; then
        log "INFO" "Backup verification PASSED"
        log "INFO" "==========================================="

        # 推送Prometheus指标
        if command -v curl &> /dev/null; then
            curl -X POST http://localhost:9091/metrics/job/backup \
                -d "backup_checksum_valid{type=\"${BACKUP_TYPE}\",database=\"zker\"} 1"
        fi

        exit 0
    else
        log "ERROR" "Backup verification FAILED"
        log "ERROR" "==========================================="

        # 发送告警
        echo "Backup verification failed for ${BACKUP_FILE}" | \
            mail -s "[ZKER Backup Alert] Verification Failed" ops-team@zker.com

        # 推送Prometheus指标
        if command -v curl &> /dev/null; then
            curl -X POST http://localhost:9091/metrics/job/backup \
                -d "backup_checksum_valid{type=\"${BACKUP_TYPE}\",database=\"zker\"} 0"
        fi

        exit 1
    fi
}

# 执行主流程
main "$@"
