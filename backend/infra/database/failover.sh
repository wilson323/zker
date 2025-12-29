#!/bin/bash
# MySQL主从故障切换脚本
# 用于在Master宕机时手动切换到Slave

set -e

# 配置
MYSQL_MASTER="mysql-master"
MYSQL_SLAVE="mysql-slave1"
MYSQL_USER="root"
MYSQL_PASSWORD="root"

echo "=== MySQL主从故障切换脚本 ==="
echo "当前Master: $MYSQL_MASTER"
echo "切换目标Slave: $MYSQL_SLAVE"
echo ""

# 确认切换
read -p "确认要执行故障切换吗？(yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    echo "取消切换"
    exit 0
fi

echo "=== 开始故障切换流程 ==="

# 步骤1: 停止Slave复制
echo "1. 停止Slave复制..."
docker exec $MYSQL_SLAVE mysql -u$MYSQL_USER -p$MYSQL_PASSWORD -e "STOP SLAVE;"

# 步骤2: 等待Slave复制完成
echo "2. 等待Slave复制完成..."
docker exec $MYSQL_SLAVE mysql -u$MYSQL_USER -p$MYSQL_PASSWORD -e "SHOW SLAVE STATUS\G" | grep "Seconds_Behind_Master: 0"
if [ $? -eq 0 ]; then
    echo "✓ Slave已完全同步"
else
    echo "✗ Slave仍有复制延迟，请手动检查"
    exit 1
fi

# 步骤3: 提升Slave为新Master
echo "3. 提升Slave为新Master..."
docker exec $MYSQL_SLAVE mysql -u$MYSQL_USER -p$MYSQL_PASSWORD -e "RESET SLAVE; SET GLOBAL read_only = OFF;"

# 步骤4: 更新ProxySQL配置
echo "4. 更新ProxySQL路由规则..."
docker exec zker-proxysql mysql -h 127.0.0.1 -P 6032 -u admin -padmin_password_2025 << EOF
UPDATE mysql_servers SET hostgroup_id = 10 WHERE hostname = '$MYSQL_SLAVE';
UPDATE mysql_servers SET hostgroup_id = 20 WHERE hostname = '$MYSQL_MASTER';
LOAD MYSQL SERVERS TO RUNTIME;
EOF

echo "✓ ProxySQL路由已更新"

# 步骤5: 验证切换
echo "5. 验证切换结果..."
docker exec zker-proxysql mysql -h 127.0.0.1 -P 6032 -u root -p -e "SELECT * FROM tenants LIMIT 1;"
if [ $? -eq 0 ]; then
    echo "✓ 切换成功，数据库读写正常"
else
    echo "✗ 切换失败，请检查"
    exit 1
fi

echo ""
echo "=== 故障切换完成 ==="
echo "新Master: $MYSQL_SLAVE"
echo "下一步: 修复原Master并配置为新Slave"
