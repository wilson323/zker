#!/bin/bash
# MySQL主从复制监控脚本
# 定期检查复制状态并上报到Pushgateway

PUSHGATEWAY_URL="http://localhost:9091"
MYSQL_MASTER_HOST="mysql-master"
MYSQL_SLAVE_HOSTS=("mysql-slave1" "mysql-slave2")
MYSQL_USER="root"
MYSQL_PASSWORD="root"

# 检查Slave复制状态
check_replication_status() {
    local master=$1
    local slave=$2
    local slave_host=$3

    # 连接到Slave查询复制状态
    status=$(mysql -h"$slave_host" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" -e "SHOW SLAVE STATUS\G" 2>/dev/null)

    # 提取复制延迟 (Seconds_Behind_Master)
    lag=$(echo "$status" | grep "Seconds_Behind_Master:" | awk '{print $2}')

    # 提取复制状态 (Slave_IO_Running 和 Slave_SQL_Running)
    io_running=$(echo "$status" | grep "Slave_IO_Running:" | awk '{print $2}')
    sql_running=$(echo "$status" | grep "Slave_SQL_Running:" | awk '{print $2}')

    # 判断复制状态 (1=running, 0=stopped)
    if [ "$io_running" = "Yes" ] && [ "$sql_running" = "Yes" ]; then
        replication_status=1
    else
        replication_status=0
    fi

    # 上报到Pushgateway
    cat <<EOF | curl --data-binary @- "${PUSHGATEWAY_URL}/metrics/job/mysql_replication/instance/${slave}"
# HELP db_replication_lag_seconds Database replication lag in seconds
# TYPE db_replication_lag_seconds gauge
db_replication_lag_seconds{master="${master}",slave="${slave}"} ${lag:-0}
# HELP db_replication_status Database replication status (1=running, 0=stopped)
# TYPE db_replication_status gauge
db_replication_status{master="${master}",slave="${slave}"} ${replication_status}
EOF

    echo "[$slave] Replication Lag: ${lag}s, Status: ${replication_status}"
}

# 检查连接池使用率
check_connection_pool() {
    local host=$1
    local type=$2  # write or read

    # 查询连接数
    connections=$(mysql -h"$host" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" -e "SHOW STATUS LIKE 'Threads_connected'" | tail -n 1 | awk '{print $2}')

    # 查询最大连接数
    max_connections=$(mysql -h"$host" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" -e "SHOW VARIABLES LIKE 'max_connections'" | tail -n 1 | awk '{print $2}')

    # 计算使用率
    usage=$(awk "BEGIN {printf \"%.2f\", ${connections}/${max_connections}}")

    # 上报到Pushgateway
    cat <<EOF | curl --data-binary @- "${PUSHGATEWAY_URL}/metrics/job/mysql_connection_pool/instance/${host}"
# HELP db_connection_pool_usage Database connection pool usage (0-1)
# TYPE db_connection_pool_usage gauge
db_connection_pool_usage{database="zker",type="${type}"} ${usage}
EOF

    echo "[$host] Connection Pool: ${connections}/${max_connections} (${usage})"
}

# 主循环
while true; do
    echo "=== $(date '+%Y-%m-%d %H:%M:%S') ==="

    # 检查所有Slave的复制状态
    for slave_host in "${MYSQL_SLAVE_HOSTS[@]}"; do
        slave_name=${slave_host}
        check_replication_status "$MYSQL_MASTER_HOST" "$slave_name" "$slave_host"
    done

    # 检查Master连接池 (写库)
    check_connection_pool "$MYSQL_MASTER_HOST" "write"

    # 检查Slave连接池 (读库)
    for slave_host in "${MYSQL_SLAVE_HOSTS[@]}"; do
        check_connection_pool "$slave_host" "read"
    done

    echo ""
    sleep 60  # 每分钟检查一次
done
