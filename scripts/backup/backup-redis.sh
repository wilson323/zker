#!/bin/bash
# =============================================================================
# Redis自动备份脚本
# 版本: v1.0.0
# 更新: 2025-01-03
# 功能: RDB快照 + AOF备份 + 自动清理 + 验证
# =============================================================================

set -euo pipefail

# ==================== 配置区域 ====================
BACKUP_ROOT="/backup/redis"
BACKUP_LOG="/var/log/redis-backup.log"
RETENTION_DAYS=7

# Redis连接配置
REDIS_HOST="${REDIS_HOST:-localhost}"
REDIS_PORT="${REDIS_PORT:-6379}"
REDIS_PASSWORD="${REDIS_PASSWORD:-}"
REDIS_DB="${REDIS_DB:-0}"

# 备份类型 (rdb/aof)
BACKUP_TYPE="${BACKUP_TYPE:-rdb}"

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

    # 检查redis-cli
    if ! command -v redis-cli &> /dev/null; then
        log "ERROR" "redis-cli命令未安装"
        exit 1
    fi

    # 测试Redis连接
    local auth_cmd=""
    if [ -n "${REDIS_PASSWORD}" ]; then
        auth_cmd="-a ${REDIS_PASSWORD}"
    fi

    if ! redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} PING 2>&1 | grep -q "PONG"; then
        log "ERROR" "Redis连接失败: ${REDIS_HOST}:${REDIS_PORT}"
        exit 1
    fi

    # 创建备份目录
    mkdir -p "${BACKUP_ROOT}/rdb"
    mkdir -p "${BACKUP_ROOT}/aof"
    mkdir -p "${BACKUP_ROOT}/logs"

    log "INFO" "前置检查完成"
}

# ==================== RDB快照备份 ====================
rdb_backup() {
    log "INFO" "开始RDB快照备份..."
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local backup_file="${BACKUP_ROOT}/rdb/redis-dump-${timestamp}.rdb"

    # Redis认证命令
    local auth_cmd=""
    if [ -n "${REDIS_PASSWORD}" ]; then
        auth_cmd="-a ${REDIS_PASSWORD}"
    fi

    # 触发BGSAVE
    log "INFO" "触发BGSAVE..."
    local bgsave_result=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} BGSAVE 2>&1)

    if echo "${bgsave_result}" | grep -q "Background saving started"; then
        log "INFO" "BGSAVE已启动"
    else
        log "ERROR" "BGSAVE启动失败: ${bgsave_result}"
        return 1
    fi

    # 等待BGSAVE完成
    log "INFO" "等待BGSAVE完成..."
    local max_wait=300  # 最多等待5分钟
    local waited=0

    while [ ${waited} -lt ${max_wait} ]; do
        local lastsave=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} LASTSAVE 2>&1)
        local current_time=$(date +%s)
        local lastsave_time=$(date -d "${lastsave}" +%s 2>/dev/null || echo 0)

        sleep 2
        ((waited+=2))

        # 检查是否在备份中
        local info=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} INFO persistence 2>&1)
        if echo "${info}" | grep -q "rdb_bgsave_in_progress:1"; then
            continue
        fi

        # 检查LASTSAVE是否更新
        local new_lastsave=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} LASTSAVE 2>&1)
        local new_lastsave_time=$(date -d "${new_lastsave}" +%s 2>/dev/null || echo 0)

        if [ ${new_lastsave_time} -gt ${lastsave_time} ]; then
            log "INFO" "BGSAVE已完成"
            break
        fi
    done

    if [ ${waited} -ge ${max_wait} ]; then
        log "ERROR" "BGSAVE超时"
        return 1
    fi

    # 获取RDB文件路径
    local config=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} CONFIG GET dir 2>&1)
    local redis_dir=$(echo "${config}" | tail -n 1)

    # 复制RDB文件
    local rdb_file="${redis_dir}/dump.rdb"
    if [ ! -f "${rdb_file}" ]; then
        log "ERROR" "RDB文件不存在: ${rdb_file}"
        return 1
    fi

    cp -f "${rdb_file}" "${backup_file}"

    if [ $? -eq 0 ]; then
        local backup_size=$(du -h "${backup_file}" | cut -f1)
        log "INFO" "RDB备份成功: ${backup_file} (${backup_size})"

        # 记录备份元数据
        echo "${timestamp}|${backup_file}|${backup_size}|rdb" >> "${BACKUP_ROOT}/logs/backups.log"

        return 0
    else
        log "ERROR" "RDB备份失败"
        return 1
    fi
}

# ==================== AOF备份 ====================
aof_backup() {
    log "INFO" "开始AOF备份..."
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local backup_file="${BACKUP_ROOT}/aof/redis-appendonly-${timestamp}.aof"

    # Redis认证命令
    local auth_cmd=""
    if [ -n "${REDIS_PASSWORD}" ]; then
        auth_cmd="-a ${REDIS_PASSWORD}"
    fi

    # 检查AOF是否启用
    local info=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} INFO persistence 2>&1)
    if ! echo "${info}" | grep -q "aof_enabled:1"; then
        log "WARN" "AOF未启用，无法执行AOF备份"
        return 1
    fi

    # 触发AOF重写
    log "INFO" "触发AOF重写..."
    redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} BGREWRITEAOF &> /dev/null || true

    # 等待AOF重写完成
    sleep 5

    # 获取AOF文件路径
    local config=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} CONFIG GET dir 2>&1)
    local redis_dir=$(echo "${config}" | tail -n 1)
    local appendonly_file=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} CONFIG GET appendfilename 2>&1 | tail -n 1)

    # 复制AOF文件
    local aof_file="${redis_dir}/${appendonly_file}"
    if [ ! -f "${aof_file}" ]; then
        log "ERROR" "AOF文件不存在: ${aof_file}"
        return 1
    fi

    cp -f "${aof_file}" "${backup_file}"

    if [ $? -eq 0 ]; then
        local backup_size=$(du -h "${backup_file}" | cut -f1)
        log "INFO" "AOF备份成功: ${backup_file} (${backup_size})"

        # 记录备份元数据
        echo "${timestamp}|${backup_file}|${backup_size}|aof" >> "${BACKUP_ROOT}/logs/backups.log"

        return 0
    else
        log "ERROR" "AOF备份失败"
        return 1
    fi
}

# ==================== 数据导出备份 ====================
export_backup() {
    log "INFO" "开始数据导出备份..."
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local backup_file="${BACKUP_ROOT}/redis-export-${timestamp}.json"

    # Redis认证命令
    local auth_cmd=""
    if [ -n "${REDIS_PASSWORD}" ]; then
        auth_cmd="-a ${REDIS_PASSWORD}"
    fi

    # 使用redis-dump工具（如果可用）
    if command -v redis-dump &> /dev/null; then
        log "INFO" "使用redis-dump导出数据..."
        redis-dump -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} > "${backup_file}"
    else
        # 使用redis-cli导出
        log "INFO" "使用redis-cli导出数据..."
        local keys=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} --scan --pattern "*" 2>&1)

        echo "{" > "${backup_file}"
        local first=true

        for key in ${keys}; do
            if [ "${first}" = true ]; then
                first=false
            else
                echo "," >> "${backup_file}"
            fi

            local key_type=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} TYPE "${key}" 2>&1)
            local ttl=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} TTL "${key}" 2>&1)

            echo "\"${key}\": {" >> "${backup_file}"
            echo "  \"type\": \"${key_type}\"," >> "${backup_file}"
            echo "  \"ttl\": ${ttl}," >> "${backup_file}"

            case "${key_type}" in
                string)
                    local value=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} GET "${key}" 2>&1)
                    echo "  \"value\": \"${value}\"" >> "${backup_file}"
                    ;;
                list)
                    local value=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} LRANGE "${key}" 0 -1 2>&1)
                    echo "  \"value\": ${value}" >> "${backup_file}"
                    ;;
                hash)
                    local value=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} HGETALL "${key}" 2>&1)
                    echo "  \"value\": ${value}" >> "${backup_file}"
                    ;;
                set)
                    local value=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} SMEMBERS "${key}" 2>&1)
                    echo "  \"value\": ${value}" >> "${backup_file}"
                    ;;
                zset)
                    local value=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} ZRANGE "${key}" 0 -1 WITHSCORES 2>&1)
                    echo "  \"value\": ${value}" >> "${backup_file}"
                    ;;
            esac

            echo "}" >> "${backup_file}"
        done

        echo "}" >> "${backup_file}"
    fi

    # 压缩备份文件
    gzip -f "${backup_file}"
    backup_file="${backup_file}.gz"

    local backup_size=$(du -h "${backup_file}" | cut -f1)
    log "INFO" "数据导出备份成功: ${backup_file} (${backup_size})"

    return 0
}

# ==================== 清理旧备份 ====================
cleanup_old_backups() {
    log "INFO" "清理${RETENTION_DAYS}天前的旧备份..."

    local deleted_count=0

    # 清理RDB备份
    while IFS= read -r -d '' old_file; do
        rm -f "${old_file}"
        log "INFO" "删除: ${old_file}"
        ((deleted_count++))
    done < <(find "${BACKUP_ROOT}/rdb" -name "redis-dump-*.rdb" -mtime +${RETENTION_DAYS} -print0)

    # 清理AOF备份
    while IFS= read -r -d '' old_file; do
        rm -f "${old_file}"
        log "INFO" "删除: ${old_file}"
        ((deleted_count++))
    done < <(find "${BACKUP_ROOT}/aof" -name "redis-appendonly-*.aof" -mtime +${RETENTION_DAYS} -print0)

    # 清理导出备份
    while IFS= read -r -d '' old_file; do
        rm -f "${old_file}"
        log "INFO" "删除: ${old_file}"
        ((deleted_count++))
    done < <(find "${BACKUP_ROOT}" -name "redis-export-*.json.gz" -mtime +${RETENTION_DAYS} -print0)

    log "INFO" "清理完成，删除${deleted_count}个文件"
}

# ==================== 验证备份 ====================
verify_backup() {
    local backup_file=$1
    local backup_type=$2

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

    # 根据类型验证
    case "${backup_type}" in
        rdb)
            # RDB文件验证（使用redis-check-rdb）
            if command -v redis-check-rdb &> /dev/null; then
                if ! redis-check-rdb "${backup_file}" 2>&1 | grep -q "Checksum OK"; then
                    log "ERROR" "RDB文件验证失败"
                    return 1
                fi
            fi
            ;;
        aof)
            # AOF文件验证（使用redis-check-aof）
            if command -v redis-check-aof &> /dev/null; then
                if ! redis-check-aof "${backup_file}" 2>&1 | grep -q "File is valid"; then
                    log "ERROR" "AOF文件验证失败"
                    return 1
                fi
            fi
            ;;
    esac

    log "INFO" "验证通过: ${backup_file} (${file_size} bytes)"
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
    log "INFO" "Redis备份任务开始"
    log "INFO" "备份类型: ${BACKUP_TYPE}"
    log "INFO" "=========================================="

    # 前置检查
    pre_check || exit 1

    # 执行备份
    local backup_status=0
    local backup_file=""

    case "${BACKUP_TYPE}" in
        rdb)
            rdb_backup || backup_status=1
            backup_file="${BACKUP_ROOT}/rdb/redis-dump-$(date +%Y%m%d-%H%M%S).rdb"
            ;;
        aof)
            aof_backup || backup_status=1
            backup_file="${BACKUP_ROOT}/aof/redis-appendonly-$(date +%Y%m%d-%H%M%S).aof"
            ;;
        export)
            export_backup || backup_status=1
            backup_file="${BACKUP_ROOT}/redis-export-$(date +%Y%m%d-%H%M%S).json.gz"
            ;;
        all)
            rdb_backup || backup_status=1
            aof_backup || backup_status=1
            export_backup || backup_status=1
            ;;
    esac

    # 验证备份
    if [ ${backup_status} -eq 0 ] && [ -n "${backup_file}" ]; then
        verify_backup "${backup_file}" "${BACKUP_TYPE}" || backup_status=1
    fi

    # 清理旧备份
    if [ ${backup_status} -eq 0 ]; then
        cleanup_old_backups
    fi

    # 发送通知
    if [ ${backup_status} -eq 0 ]; then
        log "INFO" "备份任务成功完成"
        send_notification "success" "Redis备份成功"
    else
        log "ERROR" "备份任务失败"
        send_notification "failure" "Redis备份失败"
        exit 1
    fi

    log "INFO" "=========================================="
}

# 执行主流程
main "$@"
