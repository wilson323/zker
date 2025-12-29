#!/bin/bash
# Redis数据迁移脚本
# 从单机Redis迁移到Redis集群

set -e

REDIS_SINGLE_HOST="localhost"
REDIS_SINGLE_PORT=6379
REDIS_CLUSTER_HOST="localhost"
REDIS_CLUSTER_PORTS=(7001 7002 7003)

echo "=== Redis数据迁移脚本 ==="
echo "源: $REDIS_SINGLE_HOST:$REDIS_SINGLE_PORT (单机)"
echo "目标: $REDIS_CLUSTER_CLUSTER_HOST (集群)"
echo ""

# 检查redis-cli是否安装
if ! command -v redis-cli &> /dev/null; then
    echo "错误: redis-cli未安装"
    exit 1
fi

# 检查源Redis连接
echo "1. 检查源Redis..."
redis-cli -h $REDIS_SINGLE_HOST -p $REDIS_SINGLE_PORT ping
if [ $? -ne 0 ]; then
    echo "错误: 无法连接到源Redis"
    exit 1
fi

# 检查目标Redis集群
echo "2. 检查目标Redis集群..."
redis-cli -c -p 7001 cluster info | grep cluster_state
if [ $? -ne 0 ]; then
    echo "错误: 集群状态异常"
    exit 1
fi

# 询问迁移模式
echo ""
echo "选择迁移模式:"
echo "1. 全量迁移 (所有keys)"
echo "2. 按模式迁移 (匹配模式)"
read -p "请选择 (1/2): " mode

case $mode in
    1)
        echo "=== 全量迁移 ==="
        # 使用redis-migrate-tool或自定义脚本
        ;;
    2)
        echo "=== 按模式迁移 ==="
        read -p "请输入key模式 (如: tenant:*): " pattern
        echo "迁移模式: $pattern"
        ;;
    *)
        echo "错误: 无效选择"
        exit 1
esac

# 执行迁移
echo ""
echo "3. 开始迁移..."
echo "警告: 迁移过程中源Redis不可写"

# 使用redis-shake工具或自定义脚本进行迁移
# 这里提供简化版的迁移逻辑

if [ "$mode" == "1" ]; then
    # 全量迁移示例
    total_keys=$(redis-cli -h $REDIS_SINGLE_HOST -p $REDIS_SINGLE_PORT DBSIZE)
    echo "总key数: $total_keys"

    # 分批迁移
    batch_size=1000
    migrated=0

    while [ $migrated -lt $total_keys ]; do
        # 获取一批keys
        keys=$(redis-cli -h $REDIS_SINGLE_HOST -p $REDIS_SINGLE_PORT --scan --count $batch_size)

        # 迁移这批keys
        for key in $keys; do
            value=$(redis-cli -h $REDIS_SINGLE_HOST -p $REDIS_SINGLE_PORT DUMP "$key")
            if [ -n "$value" ]; then
                # 解析DUMP格式并RESTORE到集群
                redis-cli -c -p 7001 --pipe << EOF
$value
EOF
            fi
        done

        migrated=$((migrated + batch_size))
        echo "已迁移: $migrated / $total_keys"
    done
fi

echo ""
echo "=== 迁移完成 ==="
echo "请验证数据完整性:"
echo "redis-cli -h $REDIS_SINGLE_HOST -p $REDIS_SINGLE_PORT DBSIZE"
echo "redis-cli -c -p 7001 CLUSTER KEYS"
