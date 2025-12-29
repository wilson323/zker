#!/bin/bash
# Redis集群初始化脚本

echo "=== 开始初始化Redis集群 ==="

# 等待所有Redis节点启动
echo "等待Redis节点启动..."
sleep 10

# 创建集群 (3主3从)
# 节点分配:
# - 7001 (Master 1) -> 负责 slot 0-5460
# - 7002 (Master 2) -> 负责 slot 5461-10922
# - 7003 (Master 3) -> 负责 slot 10923-16383
# - 7004 (Slave 1) -> 复制 Master 1
# - 7005 (Slave 2) -> 复制 Master 2
# - 7006 (Slave 3) -> 复制 Master 3

echo "创建Redis集群..."
redis-cli --cluster create \
  127.0.0.1:7001 \
  127.0.0.1:7002 \
  127.0.0.1:7003 \
  127.0.0.1:7004 \
  127.0.0.1:7005 \
  127.0.0.1:7006 \
  --cluster-replicas 1 \
  --cluster-yes

echo ""
echo "=== 验证集群状态 ==="

# 检查集群信息
echo "集群信息:"
redis-cli -c -p 7001 cluster info | grep -E "cluster_state|cluster_slots_assigned"

# 检查节点状态
echo ""
echo "集群节点:"
redis-cli -c -p 7001 cluster nodes

# 检查槽位分配
echo ""
echo "槽位分配:"
redis-cli -c -p 7001 cluster slots | head -10

echo ""
echo "=== 集群初始化完成 ==="
echo "使用 redis-cli -c -p 7001 连接到集群"
echo "访问 Redis Cluster 监控: docker exec zker-redis-exporter-1 curl http://localhost:9121/metrics"
