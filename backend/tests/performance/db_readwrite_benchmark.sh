#!/bin/bash
# 数据库读写分离性能测试脚本

echo "=== 数据库读写分离性能测试 ==="

# 配置
PROXYSQL_HOST="localhost"
PROXYSQL_PORT=6033
MYSQL_USER="root"
MYSQL_PASSWORD="root"
DATABASE="zker"

# 并发数
CONCURRENCY=(10 50 100 200)

echo "测试配置:"
echo "ProxySQL: $PROXYSQL_HOST:$PROXYSQL_PORT"
echo "数据库: $DATABASE"
echo "并发级别: ${CONCURRENCY[@]}"
echo ""

# 测试1: 读性能测试
echo "=== 读性能测试 ==="
for concurrency in "${CONCURRENCY[@]}"; do
    echo "并发数: $concurrency"

    start_time=$(date +%s.%N)

    for i in $(seq 1 $concurrency); do
        mysql -h $PROXYSQL_HOST -P $PROXYSQL_PORT -u $MYSQL_USER -p$MYSQL_PASSWORD $DATABASE \
            -e "SELECT * FROM tenants LIMIT 10;" > /dev/null &
    done

    wait

    end_time=$(date +%s.%N)
    duration=$(echo "$end_time - $start_time" | bc)
    qps=$(echo "scale=2; $concurrency / $duration" | bc)

    echo "  耗时: ${duration}s"
    echo "  QPS: ${qps}"
done

echo ""

# 测试2: 写性能测试
echo "=== 写性能测试 (只测试Master) ==="
for concurrency in "${CONCURRENCY[@]}"; do
    echo "并发数: $concurrency"

    start_time=$(date +%s.%N)

    for i in $(seq 1 $concurrency); do
        mysql -h $PROXYSQL_HOST -P $PROXYSQL_PORT -u $MYSQL_USER -p$MYSQL_PASSWORD $DATABASE \
            -e "INSERT INTO test_table (id, data) VALUES ($i, 'test data');" > /dev/null &
    done

    wait

    end_time=$(date +%s.%N)
    duration=$(echo "$end_time - $start_time" | bc)
    qps=$(echo "scale=2; $concurrency / $duration" | bc)

    echo "  耗时: ${duration}s"
    echo "  QPS: ${qps}"
done

echo ""
echo "=== 性能测试完成 ==="

# 清理测试数据
echo "清理测试数据..."
mysql -h $PROXYSQL_HOST -P $PROXYSQL_PORT -u $MYSQL_USER -p$MYSQL_PASSWORD $DATABASE \
    -e "DROP TABLE IF EXISTS test_table;"
echo "✓ 清理完成"
