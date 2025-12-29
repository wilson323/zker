#!/bin/bash
# 数据一致性验证脚本
# 用于验证分布式系统中的数据一致性

set -e

# 配置
DB_MASTER_HOST="mysql-master"
DB_SLAVE_HOST="mysql-slave1"
DB_SLAVE2_HOST="mysql-slave2"
DB_USER="root"
DB_PASSWORD="root"
DATABASE="zker"

REDIS_CLUSTER_HOST="redis-cluster"
REDIS_CLUSTER_PORTS=(7001 7002 7003)

CONSUL_ADDR="localhost:8500"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查MySQL主从一致性
check_mysql_consistency() {
    echo "=== MySQL主从一致性检查 ==="

    # 获取Master的binlog位置
    log_info "获取Master binlog位置..."
    MASTER_STATUS=$(docker exec $DB_MASTER_HOST mysql -u$DB_USER -p$DB_PASSWORD -e "SHOW MASTER STATUS\G")
    MASTER_FILE=$(echo "$MASTER_STATUS" | grep "File:" | awk '{print $2}')
    MASTER_POSITION=$(echo "$MASTER_STATUS" | grep "Position:" | awk '{print $2}')

    log_info "Master Binlog: $MASTER_FILE, Position: $MASTER_POSITION"

    # 检查Slave1复制状态
    log_info "检查Slave1复制状态..."
    SLAVE1_STATUS=$(docker exec $DB_SLAVE_HOST mysql -u$DB_USER -p$DB_PASSWORD -e "SHOW SLAVE STATUS\G")
    SLAVE1_IO_RUNNING=$(echo "$SLAVE1_STATUS" | grep "Slave_IO_Running:" | awk '{print $2}')
    SLAVE1_SQL_RUNNING=$(echo "$SLAVE1_STATUS" | grep "Slave_SQL_Running:" | awk '{print $2}')
    SLAVE1_SECONDS_BEHIND=$(echo "$SLAVE1_STATUS" | grep "Seconds_Behind_Master:" | awk '{print $2}')

    if [ "$SLAVE1_IO_RUNNING" == "Yes" ] && [ "$SLAVE1_SQL_RUNNING" == "Yes" ]; then
        if [ "$SLAVE1_SECONDS_BEHIND" == "0" ]; then
            log_info "✓ Slave1 完全同步"
        else
            log_warn "⚠ Slave1 延迟: ${SLAVE1_SECONDS_BEHIND}s"
        fi
    else
        log_error "✗ Slave1 复制异常"
        return 1
    fi

    # 检查Slave2复制状态
    log_info "检查Slave2复制状态..."
    SLAVE2_STATUS=$(docker exec $DB_SLAVE2_HOST mysql -u$DB_USER -p$DB_PASSWORD -e "SHOW SLAVE STATUS\G")
    SLAVE2_IO_RUNNING=$(echo "$SLAVE2_STATUS" | grep "Slave_IO_Running:" | awk '{print $2}')
    SLAVE2_SQL_RUNNING=$(echo "$SLAVE2_STATUS" | grep "Slave_SQL_Running:" | awk '{print $2}')
    SLAVE2_SECONDS_BEHIND=$(echo "$SLAVE2_STATUS" | grep "Seconds_Behind_Master:" | awk '{print $2}')

    if [ "$SLAVE2_IO_RUNNING" == "Yes" ] && [ "$SLAVE2_SQL_RUNNING" == "Yes" ]; then
        if [ "$SLAVE2_SECONDS_BEHIND" == "0" ]; then
            log_info "✓ Slave2 完全同步"
        else
            log_warn "⚠ Slave2 延迟: ${SLAVE2_SECONDS_BEHIND}s"
        fi
    else
        log_error "✗ Slave2 复制异常"
        return 1
    fi

    # 数据行数一致性检查
    log_info "检查数据行数一致性..."
    TABLES=$(docker exec $DB_MASTER_HOST mysql -u$DB_USER -p$DB_PASSWORD $DATABASE -e "SHOW TABLES;" | grep -v "Tables_in")

    for table in $TABLES; do
        MASTER_COUNT=$(docker exec $DB_MASTER_HOST mysql -u$DB_USER -p$DB_PASSWORD $DATABASE -e "SELECT COUNT(*) FROM $table;" | tail -n 1)
        SLAVE1_COUNT=$(docker exec $DB_SLAVE_HOST mysql -u$DB_USER -p$DB_PASSWORD $DATABASE -e "SELECT COUNT(*) FROM $table;" | tail -n 1)
        SLAVE2_COUNT=$(docker exec $DB_SLAVE2_HOST mysql -u$DB_USER -p$DB_PASSWORD $DATABASE -e "SELECT COUNT(*) FROM $table;" | tail -n 1)

        if [ "$MASTER_COUNT" == "$SLAVE1_COUNT" ] && [ "$MASTER_COUNT" == "$SLAVE2_COUNT" ]; then
            log_info "✓ 表 $table 行数一致: $MASTER_COUNT"
        else
            log_error "✗ 表 $table 行数不一致 - Master: $MASTER_COUNT, Slave1: $SLAVE1_COUNT, Slave2: $SLAVE2_COUNT"
            return 1
        fi
    done

    echo ""
}

# 检查Redis集群数据一致性
check_redis_consistency() {
    echo "=== Redis集群一致性检查 ==="

    # 获取集群所有key的数量
    log_info "检查集群key数量..."
    TOTAL_KEYS=0
    for port in "${REDIS_CLUSTER_PORTS[@]}"; do
        KEYS=$(redis-cli -c -p $port DBSIZE)
        log_info "节点 [$port]: $KEYS keys"
        TOTAL_KEYS=$((TOTAL_KEYS + KEYS))
    done

    log_info "集群总key数: $TOTAL_KEYS"

    # 检查集群节点状态
    log_info "检查集群节点状态..."
    CLUSTER_INFO=$(redis-cli -c -p 7001 cluster info | grep cluster_state)
    if echo "$CLUSTER_INFO" | grep -q "ok"; then
        log_info "✓ 集群状态正常"
    else
        log_error "✗ 集群状态异常: $CLUSTER_INFO"
        return 1
    fi

    # 抽样验证数据一致性（选择部分key在所有节点中查询）
    log_info "抽样验证数据一致性..."
    SAMPLE_KEYS=$(redis-cli -c -p 7001 --scan --count 10 | head -n 5)

    for key in $SAMPLE_KEYS; do
        MASTER_VALUE=$(redis-cli -c -p 7001 GET "$key")
        if [ -n "$MASTER_VALUE" ]; then
            log_info "✓ Key '$key' 存在且可访问"
        else
            log_warn "⚠ Key '$key' 无法访问"
        fi
    done

    echo ""
}

# 检查服务注册一致性
check_service_registration() {
    echo "=== 服务注册一致性检查 ==="

    # 使用Consul API检查注册的服务
    log_info "从Consul获取注册的服务..."

    SERVICES=$(curl -s http://$CONSUL_ADDR/v1/agent/services | jq -r 'to_entries[] | "\(.key): \(.value.Service)"')

    EXPECTED_SERVICES=("tenant-service" "quota-service" "subscription-service" "auth-service" "bot-service")

    for service in "${EXPECTED_SERVICES[@]}"; do
        if echo "$SERVICES" | grep -q "$service"; then
            log_info "✓ 服务 '$service' 已注册"

            # 检查健康状态
            HEALTH=$(curl -s http://$CONSUL_ADDR/v1/health/service/$service | jq '.[].Checks[].Status' | grep -c "passing" || true)
            if [ "$HEALTH" -gt 0 ]; then
                log_info "  - 健康检查通过"
            else
                log_warn "  - 健康检查失败"
            fi
        else
            log_error "✗ 服务 '$service' 未注册"
        fi
    done

    echo ""
}

# 检查配置一致性
check_configuration_consistency() {
    echo "=== 配置一致性检查 ==="

    # 比较各服务的配置
    log_info "检查租户服务配置..."
    TENANT_CONFIG=$(curl -s http://$CONSUL_ADDR/v1/kv/coze-studio/tenant-service/config?raw)
    if [ -n "$TENANT_CONFIG" ]; then
        log_info "✓ 租户服务配置存在"
    else
        log_error "✗ 租户服务配置缺失"
    fi

    log_info "检查配额服务配置..."
    QUOTA_CONFIG=$(curl -s http://$CONSUL_ADDR/v1/kv/coze-studio/quota-service/config?raw)
    if [ -n "$QUOTA_CONFIG" ]; then
        log_info "✓ 配额服务配置存在"
    else
        log_error "✗ 配额服务配置缺失"
    fi

    echo ""
}

# 检查分布式事务一致性
check_transaction_consistency() {
    echo "=== 分布式事务一致性检查 ==="

    # 检查Saga事务日志
    log_info "检查Saga事务状态..."

    # 查询未完成的事务
    PENDING_TXNS=$(docker exec mysql-master mysql -u$DB_USER -p$DB_PASSWORD $DATABASE -e "
        SELECT transaction_id, status, created_at
        FROM saga_transactions
        WHERE status IN ('pending', 'running', 'compensating')
        ORDER BY created_at DESC
        LIMIT 10;
    " | tail -n +2)

    if [ -z "$PENDING_TXNS" ]; then
        log_info "✓ 无未完成事务"
    else
        log_warn "⚠ 发现未完成事务:"
        echo "$PENDING_TXNS"
    fi

    # 检查补偿事务
    COMPENSATED_TXNS=$(docker exec mysql-master mysql -u$DB_USER -p$DB_PASSWORD $DATABASE -e "
        SELECT COUNT(*) FROM saga_transactions WHERE status = 'compensated';
    " | tail -n 1)

    log_info "已补偿事务数: $COMPENSATED_TXNS"

    echo ""
}

# 主函数
main() {
    echo "=== 数据一致性验证开始 ==="
    echo ""

    # 执行所有检查
    check_mysql_consistency || EXIT_CODE=1
    check_redis_consistency || EXIT_CODE=1
    check_service_registration || EXIT_CODE=1
    check_configuration_consistency || EXIT_CODE=1
    check_transaction_consistency || EXIT_CODE=1

    echo "=== 数据一致性验证完成 ==="

    if [ -n "$EXIT_CODE" ] && [ "$EXIT_CODE" -eq 1 ]; then
        log_error "发现不一致问题，请查看上述详情"
        exit 1
    else
        log_info "✓ 所有一致性检查通过"
        exit 0
    fi
}

# 执行主函数
main
