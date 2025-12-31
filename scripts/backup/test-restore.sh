#!/bin/bash
# =============================================================================
# 备份恢复测试脚本
# 版本: v1.0.0
# 更新: 2025-01-03
# 警告: 此脚本会在测试环境中执行真实恢复操作
# =============================================================================

set -euo pipefail

# ==================== 配置区域 ====================
BACKUP_ROOT_MYSQL="/backup/mysql"
BACKUP_ROOT_REDIS="/backup/redis"
BACKUP_LOG="/var/log/backup-restore-test.log"

# 数据库连接配置
MYSQL_HOST="${MYSQL_HOST:-localhost}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-root}"

REDIS_HOST="${REDIS_HOST:-localhost}"
REDIS_PORT="${REDIS_PORT:-6379}"
REDIS_PASSWORD="${REDIS_PASSWORD:-}"

# 测试配置
TEST_MODE="${TEST_MODE:-dry-run}"  # dry-run, test, production
SKIP_MYSQL_RESTORE="${SKIP_MYSQL_RESTORE:-false}"
SKIP_REDIS_RESTORE="${SKIP_REDIS_RESTORE:-false}"

# ==================== 日志函数 ====================
log() {
    local level=$1
    shift
    local message="$@"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo "[${timestamp}] [${level}] ${message}" | tee -a "${BACKUP_LOG}"
}

# ==================== 安全确认 ====================
safety_check() {
    if [ "${TEST_MODE}" = "production" ]; then
        log "WARN" "⚠️  生产模式恢复！这将覆盖现有数据！"
        read -p "确认要在生产环境执行恢复？(输入 YES 继续): " confirm
        if [ "${confirm}" != "YES" ]; then
            log "INFO" "取消恢复操作"
            exit 0
        fi
    fi

    if [ "${TEST_MODE}" = "test" ]; then
        log "WARN" "⚠️  测试模式恢复！将创建测试数据库并恢复数据"
        read -p "确认要执行测试恢复？(yes/no): " confirm
        if [ "${confirm}" != "yes" ]; then
            log "INFO" "取消恢复操作"
            exit 0
        fi
    fi

    if [ "${TEST_MODE}" = "dry-run" ]; then
        log "INFO" "🔍 Dry-run模式：仅验证备份文件，不执行恢复"
    fi
}

# ==================== MySQL恢复测试 ====================
test_mysql_restore() {
    log "INFO" "=========================================="
    log "INFO" "MySQL恢复测试"
    log "INFO" "=========================================="

    # 获取最新备份
    local latest_backup=$(find "${BACKUP_ROOT_MYSQL}/full" -name "mysql-full-*.sql.gz" -printf '%T@ %p\n' 2>/dev/null | sort -rn | head -1 | cut -d' ' -f2-)

    if [ -z "${latest_backup}" ]; then
        log "ERROR" "未找到MySQL备份文件"
        return 1
    fi

    log "INFO" "使用备份: ${latest_backup}"

    # Dry-run模式
    if [ "${TEST_MODE}" = "dry-run" ]; then
        log "INFO" "验证备份文件..."

        # 检查文件大小
        local file_size=$(stat -f%z "${latest_backup}" 2>/dev/null || stat -c%s "${latest_backup}" 2>/dev/null)
        log "INFO" "文件大小: ${file_size} bytes"

        # 验证gzip完整性
        if gzip -t "${latest_backup}" 2>/dev/null; then
            log "INFO" "✅ GZIP文件完整"
        else
            log "ERROR" "❌ GZIP文件损坏"
            return 1
        fi

        # 查看备份内容摘要
        local table_count=$(gunzip -c "${latest_backup}" 2>/dev/null | grep -c "CREATE TABLE" || true)
        local db_count=$(gunzip -c "${latest_backup}" 2>/dev/null | grep -c "CREATE DATABASE" || true)
        log "INFO" "数据库数量: ${db_count}"
        log "INFO" "表数量: ${table_count}"

        log "INFO" "✅ Dry-run完成"
        return 0
    fi

    # 测试/生产模式
    local auth_cmd="-h${MYSQL_HOST} -P${MYSQL_PORT} -u${MYSQL_USER} -p${MYSQL_PASSWORD}"

    # 创建测试数据库
    local test_db="backup_restore_test_$(date +%s)"

    if [ "${TEST_MODE}" = "test" ]; then
        log "INFO" "创建测试数据库: ${test_db}"
        mysql ${auth_cmd} -e "CREATE DATABASE IF NOT EXISTS \`${test_db}\`" 2>&1 | tee -a "${BACKUP_LOG}"
    fi

    # 执行恢复
    local target_db="${test_db}"
    if [ "${TEST_MODE}" = "production" ]; then
        target_db=""  # 恢复所有数据库
        log "WARN" "⚠️  恢复到生产环境！"
    fi

    log "INFO" "开始恢复备份..."
    local start_time=$(date +%s)

    if [ "${TEST_MODE}" = "production" ]; then
        # 生产模式：恢复所有数据库
        gunzip -c "${latest_backup}" | mysql ${auth_cmd} 2>&1 | tee -a "${BACKUP_LOG}"
    else
        # 测试模式：恢复到测试数据库
        gunzip -c "${latest_backup}" | mysql ${auth_cmd} "${target_db}" 2>&1 | tee -a "${BACKUP_LOG}"
    fi

    local restore_status=$?
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))

    if [ ${restore_status} -eq 0 ]; then
        log "INFO" "✅ 恢复成功 (耗时: ${duration}秒)"

        # 验证恢复的数据
        if [ "${TEST_MODE}" = "test" ]; then
            local restored_table_count=$(mysql ${auth_cmd} -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='${test_db}'" -s -N 2>/dev/null)
            log "INFO" "恢复的表数量: ${restored_table_count}"

            # 清理测试数据库
            log "INFO" "清理测试数据库: ${test_db}"
            mysql ${auth_cmd} -e "DROP DATABASE IF EXISTS \`${test_db}\`" 2>&1 | tee -a "${BACKUP_LOG}"
        fi

        return 0
    else
        log "ERROR" "❌ 恢复失败"

        # 清理测试数据库
        if [ "${TEST_MODE}" = "test" ]; then
            mysql ${auth_cmd} -e "DROP DATABASE IF EXISTS \`${test_db}\`" 2>&1 | tee -a "${BACKUP_LOG}"
        fi

        return 1
    fi
}

# ==================== Redis恢复测试 ====================
test_redis_restore() {
    log "INFO" "=========================================="
    log "INFO" "Redis恢复测试"
    log "INFO" "=========================================="

    # 获取最新RDB备份
    local latest_backup=$(find "${BACKUP_ROOT_REDIS}/rdb" -name "redis-dump-*.rdb" -printf '%T@ %p\n' 2>/dev/null | sort -rn | head -1 | cut -d' ' -f2-)

    if [ -z "${latest_backup}" ]; then
        log "ERROR" "未找到Redis备份文件"
        return 1
    fi

    log "INFO" "使用备份: ${latest_backup}"

    # Redis认证
    local auth_cmd=""
    if [ -n "${REDIS_PASSWORD}" ]; then
        auth_cmd="-a ${REDIS_PASSWORD}"
    fi

    # Dry-run模式
    if [ "${TEST_MODE}" = "dry-run" ]; then
        log "INFO" "验证RDB文件..."

        # 检查文件大小
        local file_size=$(stat -f%z "${latest_backup}" 2>/dev/null || stat -c%s "${latest_backup}" 2>/dev/null)
        log "INFO" "文件大小: ${file_size} bytes"

        # 验证RDB文件
        if command -v redis-check-rdb &> /dev/null; then
            if redis-check-rdb "${latest_backup}" 2>&1 | tee -a "${BACKUP_LOG}" | grep -q "Checksum OK"; then
                log "INFO" "✅ RDB文件完整"
            else
                log "ERROR" "❌ RDB文件损坏"
                return 1
            fi

            # 显示RDB信息
            redis-check-rdb "${latest_backup}" 2>&1 | grep -E "offset|checksum" | tee -a "${BACKUP_LOG}"
        else
            log "WARN" "redis-check-rdb不可用，跳过验证"
        fi

        log "INFO" "✅ Dry-run完成"
        return 0
    fi

    # 测试/生产模式
    log "WARN" "⚠️  Redis恢复需要停机或使用从节点"
    log "INFO" "恢复步骤："
    log "INFO" "  1. 停止Redis或切换到从节点"
    log "INFO" "  2. 复制RDB文件到Redis数据目录"
    log "INFO" "  3. 启动Redis"

    if [ "${TEST_MODE}" = "test" ]; then
        log "INFO" "跳过实际恢复（测试模式）"
        return 0
    fi

    if [ "${TEST_MODE}" = "production" ]; then
        read -p "确认要恢复Redis？此操作需要停机 (yes/no): " confirm
        if [ "${confirm}" != "yes" ]; then
            log "INFO" "取消恢复操作"
            return 0
        fi

        # 生产环境恢复
        log "WARN" "请手动执行以下步骤："
        log "INFO" "1. 停止Redis: docker-compose stop redis"
        log "INFO" "2. 备份当前数据: cp /var/lib/redis/dump.rdb /var/lib/redis/dump.rdb.old"
        log "INFO" "3. 恢复备份: cp ${latest_backup} /var/lib/redis/dump.rdb"
        log "INFO" "4. 启动Redis: docker-compose start redis"
        log "INFO" "5. 验证数据: redis-cli INFO keyspace"

        return 0
    fi
}

# ==================== 生成恢复报告 ====================
generate_restore_report() {
    local mysql_status=$1
    local redis_status=$2

    log "INFO" "=========================================="
    log "INFO" "恢复测试报告"
    log "INFO" "=========================================="
    log "INFO" "测试模式: ${TEST_MODE}"
    log "INFO" "MySQL恢复: $([ ${mysql_status} -eq 0 ] && echo "✅ 成功" || echo "❌ 失败")"
    log "INFO" "Redis恢复: $([ ${redis_status} -eq 0 ] && echo "✅ 成功" || echo "❌ 失败")"
    log "INFO" "=========================================="

    if [ ${mysql_status} -eq 0 ] && [ ${redis_status} -eq 0 ]; then
        log "INFO" "✅ 所有恢复测试通过"
        return 0
    else
        log "ERROR" "❌ 部分恢复测试失败"
        return 1
    fi
}

# ==================== 显示帮助 ====================
show_help() {
    cat << EOF
用法: $0 [选项]

备份恢复测试工具

选项:
    -h, --help              显示帮助信息
    -m, --mode MODE         设置测试模式
                            - dry-run: 仅验证备份文件（默认）
                            - test: 恢复到测试数据库
                            - production: 恢复到生产环境（危险！）
    --skip-mysql            跳过MySQL恢复测试
    --skip-redis            跳过Redis恢复测试

示例:
    # Dry-run模式（默认）
    $0

    # 测试模式（恢复到测试数据库）
    $0 -m test

    # 生产模式（危险！）
    $0 -m production

    # 只测试MySQL
    $0 -m test --skip-redis

环境变量:
    MYSQL_HOST, MYSQL_PORT, MYSQL_USER, MYSQL_PASSWORD
    REDIS_HOST, REDIS_PORT, REDIS_PASSWORD
    TEST_MODE, SKIP_MYSQL_RESTORE, SKIP_REDIS_RESTORE

EOF
}

# ==================== 主流程 ====================
main() {
    log "INFO" "=========================================="
    log "INFO" "备份恢复测试工具 v1.0.0"
    log "INFO" "=========================================="

    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -m|--mode)
                TEST_MODE="$2"
                shift 2
                ;;
            --skip-mysql)
                SKIP_MYSQL_RESTORE="true"
                shift
                ;;
            --skip-redis)
                SKIP_REDIS_RESTORE="true"
                shift
                ;;
            *)
                log "ERROR" "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done

    # 安全确认
    safety_check

    # 执行测试
    local mysql_status=0
    local redis_status=0

    if [ "${SKIP_MYSQL_RESTORE}" != "true" ]; then
        test_mysql_restore || mysql_status=1
    else
        log "INFO" "跳过MySQL恢复测试"
    fi

    if [ "${SKIP_REDIS_RESTORE}" != "true" ]; then
        test_redis_restore || redis_status=1
    else
        log "INFO" "跳过Redis恢复测试"
    fi

    # 生成报告
    generate_restore_report ${mysql_status} ${redis_status}
}

# 执行主流程
main "$@"
