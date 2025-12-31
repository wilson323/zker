#!/bin/bash
# install-exporters.sh - Exporters安装脚本
# 职责: 只负责各类Exporters的安装、配置和启动
# 版本: v1.0.0
# 作者: ZKER DevOps Team

set -e  # 遇到错误立即退出

# ==================== 配置变量 ====================
NODE_EXPORTER_VERSION="v1.7.0"
MYSQL_EXPORTER_VERSION="v0.15.1"
REDIS_EXPORTER_VERSION="v1.55.0"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

# ==================== 颜色输出 ====================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# ==================== 检查函数 ====================
check_prerequisites() {
    log_info "检查前置条件..."

    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，请先安装Docker"
        exit 1
    fi

    log_info "前置条件检查通过"
}

# ==================== 安装Node Exporter ====================
install_node_exporter() {
    log_info "安装Node Exporter..."

    # 检查容器是否已存在
    if docker ps -a --format '{{.Names}}' | grep -q '^zker-node-exporter$'; then
        log_warn "Node Exporter容器已存在，正在删除..."
        docker stop zker-node-exporter 2>/dev/null || true
        docker rm zker-node-exporter 2>/dev/null || true
    fi

    # 启动Node Exporter容器
    docker run -d \
        --name zker-node-exporter \
        --restart unless-stopped \
        -p 9100:9100 \
        --privileged \
        -v /proc:/host/proc:ro \
        -v /sys:/host/sys:ro \
        -v /:/rootfs:ro \
        prom/node-exporter:"${NODE_EXPORTER_VERSION}" \
        --path.procfs=/host/proc \
        --path.sysfs=/host/sys \
        --collector.filesystem.mount-points-exclude=^/(sys|proc|dev|host|etc)($$|/)

    log_info "Node Exporter安装完成"
}

# ==================== 安装MySQL Exporter ====================
install_mysql_exporter() {
    log_info "安装MySQL Exporter..."

    # 检查MySQL连接配置
    MYSQL_DSN="${MYSQL_EXPORTER_DSN:-root:root@(zker-mysql:3306)/}"

    # 检查容器是否已存在
    if docker ps -a --format '{{.Names}}' | grep -q '^zker-mysqld-exporter$'; then
        log_warn "MySQL Exporter容器已存在，正在删除..."
        docker stop zker-mysqld-exporter 2>/dev/null || true
        docker rm zker-mysqld-exporter 2>/dev/null || true
    fi

    # 启动MySQL Exporter容器
    docker run -d \
        --name zker-mysqld-exporter \
        --restart unless-stopped \
        -p 9104:9104 \
        -e "DATA_SOURCE_NAME=${MYSQL_DSN}" \
        --network zker-network \
        prom/mysqld-exporter:"${MYSQL_EXPORTER_VERSION}"

    log_info "MySQL Exporter安装完成"
}

# ==================== 安装Redis Exporter ====================
install_redis_exporter() {
    log_info "安装Redis Exporter..."

    # 检查容器是否已存在
    if docker ps -a --format '{{.Names}}' | grep -q '^zker-redis-exporter$'; then
        log_warn "Redis Exporter容器已存在，正在删除..."
        docker stop zker-redis-exporter 2>/dev/null || true
        docker rm zker-redis-exporter 2>/dev/null || true
    fi

    # 启动Redis Exporter容器
    docker run -d \
        --name zker-redis-exporter \
        --restart unless-stopped \
        -p 9121:9121 \
        -e "REDIS_ADDR=redis://zker-redis:6379" \
        --network zker-network \
        oliver006/redis_exporter:"${REDIS_EXPORTER_VERSION}"

    log_info "Redis Exporter安装完成"
}

# ==================== 验证安装 ====================
verify_installation() {
    log_info "验证Exporters安装..."

    sleep 5

    # 检查Node Exporter
    if docker ps --format '{{.Names}}' | grep -q '^zker-node-exporter$'; then
        if curl -f -s http://localhost:9100/metrics > /dev/null 2>&1; then
            log_info "✅ Node Exporter运行正常 (http://localhost:9100/metrics)"
        else
            log_warn "⚠️  Node Exporter未响应"
        fi
    fi

    # 检查MySQL Exporter
    if docker ps --format '{{.Names}}' | grep -q '^zker-mysqld-exporter$'; then
        if curl -f -s http://localhost:9104/metrics > /dev/null 2>&1; then
            log_info "✅ MySQL Exporter运行正常 (http://localhost:9104/metrics)"
        else
            log_warn "⚠️  MySQL Exporter未响应"
        fi
    fi

    # 检查Redis Exporter
    if docker ps --format '{{.Names}}' | grep -q '^zker-redis-exporter$'; then
        if curl -f -s http://localhost:9121/metrics > /dev/null 2>&1; then
            log_info "✅ Redis Exporter运行正常 (http://localhost:9121/metrics)"
        else
            log_warn "⚠️  Redis Exporter未响应"
        fi
    fi

    log_info "Exporters安装验证完成"
}

# ==================== 显示访问信息 ====================
show_access_info() {
    echo ""
    echo "=========================================="
    echo "Exporters 安装完成！"
    echo "=========================================="
    echo "📍 Node Exporter: http://localhost:9100/metrics"
    echo "📍 MySQL Exporter: http://localhost:9104/metrics"
    echo "📍 Redis Exporter: http://localhost:9121/metrics"
    echo ""
    echo "常用命令:"
    echo "  查看所有Exporters: docker ps | grep exporter"
    echo "  停止所有Exporters: docker stop zker-node-exporter zker-mysqld-exporter zker-redis-exporter"
    echo "  删除所有Exporters: docker rm -f zker-node-exporter zker-mysqld-exporter zker-redis-exporter"
    echo ""
    echo "=========================================="
}

# ==================== 主函数 ====================
main() {
    echo "=========================================="
    echo "Exporters 安装脚本"
    echo "=========================================="
    echo ""

    check_prerequisites

    # 选择要安装的Exporters
    INSTALL_ALL="${1:-all}"

    if [ "$INSTALL_ALL" = "all" ] || [ "$INSTALL_ALL" = "node" ]; then
        install_node_exporter
    fi

    if [ "$INSTALL_ALL" = "all" ] || [ "$INSTALL_ALL" = "mysql" ]; then
        install_mysql_exporter
    fi

    if [ "$INSTALL_ALL" = "all" ] || [ "$INSTALL_ALL" = "redis" ]; then
        install_redis_exporter
    fi

    verify_installation
    show_access_info
}

# ==================== 执行 ====================
main "$@"
