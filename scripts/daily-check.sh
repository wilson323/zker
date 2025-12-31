#!/bin/bash
# 每日环境检查脚本
# 用途：检查开发环境的服务状态
# 使用：./scripts/daily-check.sh

set -e

echo "🔍 每日环境检查 - $(date '+%Y-%m-%d %H:%M:%S')"
echo "========================================"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查函数
check_service() {
    local service_name=$1
    local check_command=$2

    echo -n "检查 $service_name ... "
    if eval "$check_command" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 正常${NC}"
        return 0
    else
        echo -e "${RED}❌ 异常${NC}"
        return 1
    fi
}

# 1. 检查Docker服务
echo ""
echo "📦 检查Docker服务..."
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker服务未运行${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Docker服务正常${NC}"

# 2. 检查中间件容器
echo ""
echo "📦 检查中间件容器..."
docker ps --format "table {{.Names}}\t{{.Status}}" | grep -E "mysql|redis|es|minio|etcd" || true

# 3. 检查MySQL数据库
echo ""
echo "💾 检查MySQL数据库..."
MYSQL_CONTAINER=$(docker ps --filter "name=mysql" --format "{{.Names}}" | head -n 1)
if [ -n "$MYSQL_CONTAINER" ]; then
    if docker exec "$MYSQL_CONTAINER" mysql -uroot -p123456 -e "SELECT 1" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ MySQL数据库正常${NC}"

        # 显示数据库列表
        echo "📊 数据库列表:"
        docker exec "$MYSQL_CONTAINER" mysql -uroot -p123456 -e "SHOW DATABASES" | tail -n +2
    else
        echo -e "${RED}❌ MySQL数据库连接失败${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  MySQL容器未运行${NC}"
fi

# 4. 检查Redis缓存
echo ""
echo "🚀 检查Redis缓存..."
REDIS_CONTAINER=$(docker ps --filter "name=redis" --format "{{.Names}}" | head -n 1)
if [ -n "$REDIS_CONTAINER" ]; then
    if docker exec "$REDIS_CONTAINER" redis-cli PING > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Redis缓存正常${NC}"
    else
        echo -e "${RED}❌ Redis缓存连接失败${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Redis容器未运行${NC}"
fi

# 5. 检查Elasticsearch
echo ""
echo "🔍 检查Elasticsearch..."
ES_CONTAINER=$(docker ps --filter "name=elasticsearch" --format "{{.Names}}" | head -n 1)
if [ -n "$ES_CONTAINER" ]; then
    ES_URL="http://localhost:9200"
    if curl -s "$ES_URL/_cluster/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Elasticsearch正常${NC}"
        ES_HEALTH=$(curl -s "$ES_URL/_cluster/health" | grep -o '"status":"[^"]*"' | cut -d'"' -f4)
        echo "   集群状态: $ES_HEALTH"
    else
        echo -e "${RED}❌ Elasticsearch连接失败${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Elasticsearch容器未运行${NC}"
fi

# 6. 检查后端服务
echo ""
echo "🔧 检查后端服务..."
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ 后端服务正常${NC}"
else
    echo -e "${YELLOW}⚠️  后端服务未响应（可能未启动）${NC}"
fi

# 7. 检查前端服务
echo ""
echo "🎨 检查前端服务..."
if curl -s http://localhost:8888 > /dev/null 2>&1; then
    echo -e "${GREEN}✅ 前端服务正常${NC}"
else
    echo -e "${YELLOW}⚠️  前端服务未响应（可能未启动）${NC}"
fi

# 8. 检查磁盘空间
echo ""
echo "💿 检查磁盘空间..."
df -h | grep -E "/$|/data|/var/lib/docker" || true

# 9. 检查Docker资源使用
echo ""
echo "📊 Docker资源统计..."
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}" || true

echo ""
echo "========================================"
echo -e "${GREEN}✅ 环境检查完成${NC}"
echo ""
