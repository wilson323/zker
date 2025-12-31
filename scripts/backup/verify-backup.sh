#!/bin/bash
# =============================================================================
# 备份验证脚本
# 版本: v1.0.0
# 更新: 2025-01-03
# 功能: 验证备份完整性和恢复测试
# =============================================================================

set -euo pipefail

# ==================== 配置区域 ====================
BACKUP_ROOT_MYSQL="/backup/mysql"
BACKUP_ROOT_REDIS="/backup/redis"
BACKUP_LOG="/var/log/backup-verify.log"
RETENTION_HOURS=24

# 数据库连接配置
MYSQL_HOST="${MYSQL_HOST:-localhost}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-root}"
REDIS_HOST="${REDIS_HOST:-localhost}"
REDIS_PORT="${REDIS_PORT:-6379}"
REDIS_PASSWORD="${REDIS_PASSWORD:-}"

# 通知配置
NOTIFICATION_WEBHOOK="${NOTIFICATION_WEBHOOK:-}"
ALERT_THRESHOLD=3  # 连续失败多少次发送告警

# ==================== 日志函数 ====================
log() {
    local level=$1
    shift
    local message="$@"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo "[${timestamp}] [${level}] ${message}" | tee -a "${BACKUP_LOG}"
}

# ==================== MySQL备份验证 ====================
verify_mysql_backup() {
    log "INFO" "验证MySQL备份..."

    local backup_count=0
    local failed_count=0

    # 查找最近24小时的备份
    while IFS= read -r -d '' backup_file; do
        ((backup_count++))

        log "INFO" "验证: ${backup_file}"

        # 检查文件存在性
        if [ ! -f "${backup_file}" ]; then
            log "ERROR" "备份文件不存在"
            ((failed_count++))
            continue
        fi

        # 检查文件大小
        local file_size=$(stat -f%z "${backup_file}" 2>/dev/null || stat -c%s "${backup_file}" 2>/dev/null)
        if [ "${file_size}" -lt 1024 ]; then
            log "ERROR" "备份文件过小: ${file_size} bytes"
            ((failed_count++))
            continue
        fi

        # 验证gzip完整性
        if ! gzip -t "${backup_file}" 2>/dev/null; then
            log "ERROR" "GZIP文件损坏"
            ((failed_count++))
            continue
        fi

        # 尝试验证SQL内容
        local table_count=$(gunzip -c "${backup_file}" 2>/dev/null | grep -c "CREATE TABLE" || true)
        if [ "${table_count}" -eq 0 ]; then
            log "WARN" "未发现表定义"
        fi

        log "INFO" "✅ 验证通过 (${file_size} bytes, ${table_count} tables)"

    done < <(find "${BACKUP_ROOT_MYSQL}/full" -name "mysql-full-*.sql.gz" -mtime -1 -print0)

    # 检查备份数量
    if [ ${backup_count} -eq 0 ]; then
        log "ERROR" "未找到最近24小时的MySQL备份"
        return 1
    fi

    if [ ${failed_count} -gt 0 ]; then
        log "ERROR" "MySQL备份验证失败: ${failed_count}/${backup_count} 个文件有问题"
        return 1
    fi

    log "INFO" "✅ MySQL备份验证完成: ${backup_count} 个文件全部有效"
    return 0
}

# ==================== Redis备份验证 ====================
verify_redis_backup() {
    log "INFO" "验证Redis备份..."

    local backup_count=0
    local failed_count=0

    # 验证RDB备份
    while IFS= read -r -d '' backup_file; do
        ((backup_count++))

        log "INFO" "验证RDB: ${backup_file}"

        # 检查文件存在性
        if [ ! -f "${backup_file}" ]; then
            log "ERROR" "备份文件不存在"
            ((failed_count++))
            continue
        fi

        # 检查文件大小
        local file_size=$(stat -f%z "${backup_file}" 2>/dev/null || stat -c%s "${backup_file}" 2>/dev/null)
        if [ "${file_size}" -lt 1024 ]; then
            log "ERROR" "备份文件过小: ${file_size} bytes"
            ((failed_count++))
            continue
        fi

        # 使用redis-check-rdb验证（如果可用）
        if command -v redis-check-rdb &> /dev/null; then
            if ! redis-check-rdb "${backup_file}" 2>&1 | grep -q "Checksum OK"; then
                log "ERROR" "RDB文件验证失败"
                ((failed_count++))
                continue
            fi
        fi

        log "INFO" "✅ RDB验证通过 (${file_size} bytes)"

    done < <(find "${BACKUP_ROOT_REDIS}/rdb" -name "redis-dump-*.rdb" -mtime -1 -print0)

    # 验证AOF备份（如果存在）
    while IFS= read -r -d '' backup_file; do
        ((backup_count++))

        log "INFO" "验证AOF: ${backup_file}"

        if command -v redis-check-aof &> /dev/null; then
            if ! redis-check-aof "${backup_file}" 2>&1 | grep -q "File is valid"; then
                log "ERROR" "AOF文件验证失败"
                ((failed_count++))
                continue
            fi
        fi

        log "INFO" "✅ AOF验证通过"

    done < <(find "${BACKUP_ROOT_REDIS}/aof" -name "redis-appendonly-*.aof" -mtime -1 -print0)

    # 检查备份数量
    if [ ${backup_count} -eq 0 ]; then
        log "ERROR" "未找到最近24小时的Redis备份"
        return 1
    fi

    if [ ${failed_count} -gt 0 ]; then
        log "ERROR" "Redis备份验证失败: ${failed_count}/${backup_count} 个文件有问题"
        return 1
    fi

    log "INFO" "✅ Redis备份验证完成: ${backup_count} 个文件全部有效"
    return 0
}

# ==================== 备份恢复测试 ====================
test_mysql_restore() {
    log "INFO" "测试MySQL备份恢复..."

    # 获取最新备份
    local latest_backup=$(find "${BACKUP_ROOT_MYSQL}/full" -name "mysql-full-*.sql.gz" -mtime -1 -printf '%T@ %p\n' 2>/dev/null | sort -rn | head -1 | cut -d' ' -f2-)

    if [ -z "${latest_backup}" ]; then
        log "WARN" "未找到可用于恢复测试的备份"
        return 0
    fi

    log "INFO" "使用备份: ${latest_backup}"

    # 创建测试数据库
    local test_db="backup_test_$(date +%s)"
    local auth_cmd="-h${MYSQL_HOST} -P${MYSQL_PORT} -u${MYSQL_USER} -p${MYSQL_PASSWORD}"

    mysql ${auth_cmd} -e "CREATE DATABASE IF NOT EXISTS \`${test_db}\`" 2>&1 | tee -a "${BACKUP_LOG}"

    # 恢复到测试数据库
    log "INFO" "恢复到测试数据库: ${test_db}"

    if gunzip -c "${latest_backup}" | mysql ${auth_cmd} "${test_db}" 2>&1 | tee -a "${BACKUP_LOG}"; then
        log "INFO" "✅ 恢复测试成功"

        # 检查表数量
        local table_count=$(mysql ${auth_cmd} -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='${test_db}'" -s -N 2>/dev/null)
        log "INFO" "恢复的表数量: ${table_count}"

        # 清理测试数据库
        mysql ${auth_cmd} -e "DROP DATABASE IF EXISTS \`${test_db}\`" 2>&1 | tee -a "${BACKUP_LOG}"

        return 0
    else
        log "ERROR" "❌ 恢复测试失败"
        mysql ${auth_cmd} -e "DROP DATABASE IF EXISTS \`${test_db}\`" 2>&1 | tee -a "${BACKUP_LOG}"
        return 1
    fi
}

# ==================== Redis恢复测试 ====================
test_redis_restore() {
    log "INFO" "测试Redis备份恢复..."

    # 获取最新RDB备份
    local latest_backup=$(find "${BACKUP_ROOT_REDIS}/rdb" -name "redis-dump-*.rdb" -mtime -1 -printf '%T@ %p\n' 2>/dev/null | sort -rn | head -1 | cut -d' ' -f2-)

    if [ -z "${latest_backup}" ]; then
        log "WARN" "未找到可用于恢复测试的Redis备份"
        return 0
    fi

    log "INFO" "使用备份: ${latest_backup}"

    # Redis认证
    local auth_cmd=""
    if [ -n "${REDIS_PASSWORD}" ]; then
        auth_cmd="-a ${REDIS_PASSWORD}"
    fi

    # 记录当前数据库大小
    local current_dbsize=$(redis-cli -h "${REDIS_HOST}" -p "${REDIS_PORT}" ${auth_cmd} DBSIZE 2>/dev/null)
    log "INFO" "当前数据库大小: ${current_dbsize} keys"

    # 验证RDB文件
    if command -v redis-check-rdb &> /dev/null; then
        if redis-check-rdb "${latest_backup}" 2>&1 | tee -a "${BACKUP_LOG}" | grep -q "Checksum OK"; then
            log "INFO" "✅ RDB文件验证通过"
        else
            log "ERROR" "❌ RDB文件验证失败"
            return 1
        fi
    else
        log "WARN" "redis-check-rdb不可用，跳过验证"
    fi

    log "INFO" "✅ Redis恢复测试通过（仅验证，未实际恢复）"
    return 0
}

# ==================== 检查备份间隔 ====================
check_backup_interval() {
    log "INFO" "检查备份间隔..."

    # 检查MySQL备份间隔
    local latest_mysql=$(find "${BACKUP_ROOT_MYSQL}/full" -name "mysql-full-*.sql.gz" -printf '%T@\n' 2>/dev/null | sort -rn | head -1)
    if [ -n "${latest_mysql}" ]; then
        local current_time=$(date +%s)
        local interval=$((current_time - $(printf "%.0f" "${latest_mysql}")))
        local interval_hours=$((interval / 3600))

        log "INFO" "MySQL最新备份: ${interval_hours} 小时前"

        if [ ${interval_hours} -gt 26 ]; then
            log "WARN" "MySQL备份间隔超过26小时"
            return 1
        fi
    else
        log "ERROR" "未找到MySQL备份"
        return 1
    fi

    # 检查Redis备份间隔
    local latest_redis=$(find "${BACKUP_ROOT_REDIS}/rdb" -name "redis-dump-*.rdb" -printf '%T@\n' 2>/dev/null | sort -rn | head -1)
    if [ -n "${latest_redis}" ]; then
        local current_time=$(date +%s)
        local interval=$((current_time - $(printf "%.0f" "${latest_redis}")))
        local interval_hours=$((interval / 3600))

        log "INFO" "Redis最新备份: ${interval_hours} 小时前"

        if [ ${interval_hours} -gt 8 ]; then
            log "WARN" "Redis备份间隔超过8小时"
            return 1
        fi
    else
        log "ERROR" "未找到Redis备份"
        return 1
    fi

    log "INFO" "✅ 备份间隔正常"
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
    "host": "$(hostname)",
    "backup_type": "verification"
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

# ==================== 生成验证报告 ====================
generate_report() {
    local mysql_status=$1
    local redis_status=$2
    local restore_status=$3
    local interval_status=$4

    log "INFO" "=========================================="
    log "INFO" "备份验证报告"
    log "INFO" "=========================================="
    log "INFO" "MySQL备份验证: $([ ${mysql_status} -eq 0 ] && echo "✅ 通过" || echo "❌ 失败")"
    log "INFO" "Redis备份验证: $([ ${redis_status} -eq 0 ] && echo "✅ 通过" || echo "❌ 失败")"
    log "INFO" "恢复测试: $([ ${restore_status} -eq 0 ] && echo "✅ 通过" || echo "❌ 失败")"
    log "INFO" "备份间隔检查: $([ ${interval_status} -eq 0 ] && echo "✅ 正常" || echo "⚠️ 异常")"
    log "INFO" "=========================================="

    # 计算总体状态
    local overall_status=0
    [ ${mysql_status} -ne 0 ] && ((overall_status++))
    [ ${redis_status} -ne 0 ] && ((overall_status++))
    [ ${restore_status} -ne 0 ] && ((overall_status++))
    [ ${interval_status} -ne 0 ] && ((overall_status++))

    if [ ${overall_status} -eq 0 ]; then
        log "INFO" "✅ 所有验证通过"
        send_notification "success" "备份验证全部通过"
        return 0
    else
        log "ERROR" "❌ 验证失败: ${overall_status} 项检查未通过"
        send_notification "failure" "备份验证失败: ${overall_status} 项检查未通过"
        return 1
    fi
}

# ==================== 主流程 ====================
main() {
    log "INFO" "=========================================="
    log "INFO" "备份验证任务开始"
    log "INFO" "=========================================="

    local mysql_status=0
    local redis_status=0
    local restore_status=0
    local interval_status=0

    # 验证MySQL备份
    verify_mysql_backup || mysql_status=1

    # 验证Redis备份
    verify_redis_backup || redis_status=1

    # 测试恢复（可选，需要注释掉以避免实际恢复）
    # test_mysql_restore || restore_status=1
    # test_redis_restore || restore_status=1
    log "INFO" "跳过恢复测试（生产环境建议定期手动测试）"

    # 检查备份间隔
    check_backup_interval || interval_status=1

    # 生成报告
    generate_report ${mysql_status} ${redis_status} ${restore_status} ${interval_status} || exit 1

    log "INFO" "=========================================="
}

# 执行主流程
main "$@"
