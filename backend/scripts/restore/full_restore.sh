#!/bin/bash
# ============================================================
# ZKER 全量恢复脚本
# ============================================================
# 功能: 从全量备份恢复MySQL数据库
# 用途: 数据库宕机、数据损坏、灾难恢复
# 预期RTO: 30-45分钟
# ============================================================

set -euo pipefail

# ============================================================
# 配置参数
# ============================================================

BACKUP_FILE=$1
TARGET_DIR="${2:-/var/lib/mysql}"
BACKUP_DIR="/data/backups/full"
LOG_FILE="/var/log/restore.log"
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-}"
MYSQL_PORT="${MYSQL_PORT:-3306}"

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

confirm() {
    local prompt="$1"
    local response

    echo -n "${prompt} (yes/no): "
    read response

    if [ "${response}" != "yes" ]; then
        log "INFO" "Restore cancelled by user"
        exit 0
    fi
}

# ============================================================
# 主流程
# ============================================================

main() {
    log "INFO" "==========================================="
    log "INFO" "Starting full restore process"
    log "INFO" "Backup file: ${BACKUP_FILE}"
    log "INFO" "Target directory: ${TARGET_DIR}"
    log "INFO" "==========================================="

    local start_time=$(date +%s)

    # 1. 确认操作
    confirm "This will replace all data in ${TARGET_DIR}. Are you sure?"

    # 2. 检查备份文件
    log "INFO" "Step 1: Checking backup file..."
    if [ ! -f "${BACKUP_FILE}" ]; then
        log "ERROR" "Backup file does not exist: ${BACKUP_FILE}"
        exit 1
    fi

    # 验证校验和
    if [ -f "${BACKUP_FILE}.sha256" ]; then
        log "INFO" "Verifying checksum..."
        sha256sum -c ${BACKUP_FILE}.sha256
        if [ $? -ne 0 ]; then
            log "ERROR" "Checksum verification failed"
            exit 1
        fi
        log "INFO" "Checksum OK"
    else
        log "WARN" "Checksum file not found, skipping verification"
    fi

    # 3. 停止MySQL服务
    log "INFO" "Step 2: Stopping MySQL service..."
    systemctl stop mysql

    if [ $? -ne 0 ]; then
        log "ERROR" "Failed to stop MySQL service"
        exit 1
    fi

    # 等待MySQL完全停止
    sleep 5

    # 4. 备份当前数据目录
    log "INFO" "Step 3: Backing up current data directory..."
    local backup_timestamp=$(date +%Y%m%d_%H%M%S)
    local current_backup="${TARGET_DIR}.backup_${backup_timestamp}"

    if [ -d "${TARGET_DIR}" ]; then
        mv ${TARGET_DIR} ${current_backup}
        log "INFO" "Current data backed up to: ${current_backup}"
    else
        log "INFO" "Target directory does not exist, creating..."
        mkdir -p ${TARGET_DIR}
    fi

    # 5. 创建临时目录并解压备份
    log "INFO" "Step 4: Extracting backup file..."
    local temp_dir="/tmp/restore_$$"
    mkdir -p ${temp_dir}

    cd ${temp_dir}
    tar -xzf ${BACKUP_FILE}

    if [ $? -ne 0 ]; then
        log "ERROR" "Failed to extract backup file"
        # 恢复原数据目录
        mv ${current_backup} ${TARGET_DIR}
        systemctl start mysql
        exit 1
    fi

    # 6. 准备备份
    log "INFO" "Step 5: Preparing backup..."
    xtrabackup --prepare --target-dir=${temp_dir} 2>&1 | tee -a ${LOG_FILE}

    if [ $? -ne 0 ]; then
        log "ERROR" "Failed to prepare backup"
        # 恢复原数据目录
        rm -rf ${temp_dir}
        mv ${current_backup} ${TARGET_DIR}
        systemctl start mysql
        exit 1
    fi

    # 7. 恢复数据
    log "INFO" "Step 6: Restoring data..."
    xtrabackup --copy-back --target-dir=${temp_dir} 2>&1 | tee -a ${LOG_FILE}

    if [ $? -ne 0 ]; then
        log "ERROR" "Failed to restore backup"
        # 恢复原数据目录
        rm -rf ${temp_dir}
        mv ${current_backup} ${TARGET_DIR}
        systemctl start mysql
        exit 1
    fi

    # 8. 恢复权限
    log "INFO" "Step 7: Restoring permissions..."
    chown -R mysql:mysql ${TARGET_DIR}
    chmod -R 750 ${TARGET_DIR}

    # 9. 清理临时文件
    log "INFO" "Step 8: Cleaning up temporary files..."
    rm -rf ${temp_dir}

    # 10. 启动MySQL服务
    log "INFO" "Step 9: Starting MySQL service..."
    systemctl start mysql

    if [ $? -ne 0 ]; then
        log "ERROR" "Failed to start MySQL service"
        log "INFO" "Check MySQL error log for details"
        exit 1
    fi

    # 等待MySQL启动
    log "INFO" "Waiting for MySQL to be ready..."
    local max_attempts=60
    local attempt=0

    while [ ${attempt} -lt ${max_attempts} ]; do
        if mysql -u root -p${MYSQL_ROOT_PASSWORD} -e "SELECT 1" > /dev/null 2>&1; then
            log "INFO" "MySQL is ready"
            break
        fi

        attempt=$((attempt + 1))
        sleep 1
    done

    if [ ${attempt} -eq ${max_attempts} ]; then
        log "ERROR" "MySQL failed to start within ${max_attempts} seconds"
        exit 1
    fi

    # 11. 验证数据
    log "INFO" "Step 10: Verifying data integrity..."

    # 检查表数量
    local table_count=$(mysql -u root -p${MYSQL_ROOT_PASSWORD} \
        -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'zker';" \
        --skip-column-names)

    log "INFO" "Tables in zker database: ${table_count}"

    if [ ${table_count} -lt 10 ]; then
        log "ERROR" "Table count too low, restore may have failed"
        exit 1
    fi

    # 检查关键表
    log "INFO" "Checking critical tables..."

    mysql -u root -p${MYSQL_ROOT_PASSWORD} -e "
        SELECT
            (SELECT COUNT(*) FROM zker.tenants) AS tenant_count,
            (SELECT COUNT(*) FROM zker.bots) AS bot_count,
            (SELECT COUNT(*) FROM zker.conversations) AS conversation_count,
            (SELECT COUNT(*) FROM zker.messages) AS message_count;
    " | tee -a ${LOG_FILE}

    if [ $? -ne 0 ]; then
        log "ERROR" "Failed to verify data"
        exit 1
    fi

    # 12. 计算恢复时间
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    local duration_minutes=$((duration / 60))

    log "INFO" "==========================================="
    log "INFO" "Full restore completed successfully"
    log "INFO" "Total duration: ${duration} seconds (${duration_minutes} minutes)"
    log "INFO" "==========================================="

    # 13. 发送通知
    echo "Full restore completed successfully in ${duration_minutes} minutes" | \
        mail -s "[ZKER Restore] Success" ops-team@zker.com

    # 14. 推送Prometheus指标
    if command -v curl &> /dev/null; then
        curl -X POST http://localhost:9091/metrics/job/restore \
            -d "restore_success{type=\"full\",database=\"zker\"} 1" \
            -d "restore_duration_seconds{type=\"full\",database=\"zker\"} ${duration}"
    fi

    return 0
}

# 执行主流程
main "$@"
