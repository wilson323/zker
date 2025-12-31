#!/bin/bash
# install-prometheus.sh - Prometheus安装脚本
# 职责: 只负责Prometheus的安装、配置和启动
# 版本: v1.0.0
# 作者: ZKER DevOps Team

set -e  # 遇到错误立即退出

# ==================== 配置变量 ====================
PROMETHEUS_VERSION="v2.45.0"
PROMETHEUS_PORT="${PROMETHEUS_PORT:-9090}"
PROMETHEUS_DIR="/opt/prometheus"
CONFIG_DIR="/etc/prometheus"
DATA_DIR="/var/lib/prometheus"
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

    # 检查Docker Compose
    if ! command -v docker compose &> /dev/null; then
        log_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi

    # 检查端口是否被占用
    if netstat -tuln 2>/dev/null | grep -q ":${PROMETHEUS_PORT} "; then
        log_warn "端口${PROMETHEUS_PORT}已被占用，请检查是否有其他Prometheus实例运行"
    fi

    log_info "前置条件检查通过"
}

# ==================== 创建目录 ====================
create_directories() {
    log_info "创建Prometheus目录..."

    sudo mkdir -p "$CONFIG_DIR"
    sudo mkdir -p "$DATA_DIR"
    sudo mkdir -p "$PROMETHEUS_DIR/rules"

    log_info "目录创建完成"
}

# ==================== 复制配置文件 ====================
copy_configs() {
    log_info "复制Prometheus配置文件..."

    # 复制prometheus.yml
    if [ -f "$PROJECT_ROOT/docker/volumes/monitoring/prometheus/prometheus.yml" ]; then
        sudo cp "$PROJECT_ROOT/docker/volumes/monitoring/prometheus/prometheus.yml" "$CONFIG_DIR/prometheus.yml"
        log_info "prometheus.yml复制完成"
    else
        log_warn "prometheus.yml未找到，使用默认配置"
    fi

    # 复制alerts.yml
    if [ -f "$PROJECT_ROOT/docker/volumes/monitoring/prometheus/alerts.yml" ]; then
        sudo cp "$PROJECT_ROOT/docker/volumes/monitoring/prometheus/alerts.yml" "$CONFIG_DIR/alerts.yml"
        log_info "alerts.yml复制完成"
    else
        log_warn "alerts.yml未找到，跳过"
    fi

    # 设置权限
    sudo chown -R root:root "$CONFIG_DIR"
    sudo chmod 644 "$CONFIG_DIR"/*.yml 2>/dev/null || true

    log_info "配置文件复制完成"
}

# ==================== 启动Prometheus ====================
start_prometheus() {
    log_info "启动Prometheus容器..."

    # 检查容器是否已存在
    if docker ps -a --format '{{.Names}}' | grep -q '^zker-prometheus$'; then
        log_warn "Prometheus容器已存在，正在删除..."
        docker stop zker-prometheus 2>/dev/null || true
        docker rm zker-prometheus 2>/dev/null || true
    fi

    # 启动Prometheus容器
    docker run -d \
        --name zker-prometheus \
        --restart unless-stopped \
        -p "${PROMETHEUS_PORT}:9090" \
        -v "$CONFIG_DIR/prometheus.yml:/etc/prometheus/prometheus.yml:ro" \
        -v "$CONFIG_DIR/alerts.yml:/etc/prometheus/alerts.yml:ro" \
        -v "$DATA_DIR:/prometheus" \
        prom/prometheus:"${PROMETHEUS_VERSION}" \
        --config.file=/etc/prometheus/prometheus.yml \
        --storage.tsdb.path=/prometheus \
        --storage.tsdb.retention.time=30d \
        --web.console.libraries=/usr/share/prometheus/console_libraries \
        --web.console.templates=/usr/share/prometheus/consoles \
        --web.enable-lifecycle \
        --web.enable-admin-api

    log_info "Prometheus容器启动成功"
}

# ==================== 验证安装 ====================
verify_installation() {
    log_info "验证Prometheus安装..."

    # 等待Prometheus启动
    sleep 5

    # 检查容器状态
    if docker ps --format '{{.Names}}' | grep -q '^zker-prometheus$'; then
        log_info "✅ Prometheus容器运行正常"
    else
        log_error "❌ Prometheus容器未运行"
        docker logs zker-prometheus
        exit 1
    fi

    # 检查HTTP端点
    if curl -f -s http://localhost:${PROMETHEUS_PORT}/-/healthy > /dev/null 2>&1; then
        log_info "✅ Prometheus HTTP端点响应正常"
    else
        log_warn "⚠️  Prometheus HTTP端点未响应，可能需要更多时间启动"
    fi

    # 检查配置文件
    if docker exec zker-prometheus promtool check config /etc/prometheus/prometheus.yml > /dev/null 2>&1; then
        log_info "✅ Prometheus配置文件语法正确"
    else
        log_warn "⚠️  Prometheus配置文件可能存在问题"
    fi

    log_info "Prometheus安装验证完成"
}

# ==================== 显示访问信息 ====================
show_access_info() {
    echo ""
    echo "=========================================="
    echo "Prometheus 安装完成！"
    echo "=========================================="
    echo "📍 访问地址: http://localhost:${PROMETHEUS_PORT}"
    echo "📊 Web UI: http://localhost:${PROMETHEUS_PORT}/graph"
    echo "📖 配置文件: $CONFIG_DIR/prometheus.yml"
    echo "💾 数据目录: $DATA_DIR"
    echo ""
    echo "常用命令:"
    echo "  查看日志: docker logs -f zker-prometheus"
    echo "  停止服务: docker stop zker-prometheus"
    echo "  启动服务: docker start zker-prometheus"
    echo "  重启服务: docker restart zker-prometheus"
    echo "  删除容器: docker rm -f zker-prometheus"
    echo ""
    echo "配置热重载:"
    echo "  docker exec zker-prometheus kill -HUP 1"
    echo ""
    echo "=========================================="
}

# ==================== 主函数 ====================
main() {
    echo "=========================================="
    echo "Prometheus 安装脚本"
    echo "版本: ${PROMETHEUS_VERSION}"
    echo "=========================================="
    echo ""

    check_prerequisites
    create_directories
    copy_configs
    start_prometheus
    verify_installation
    show_access_info
}

# ==================== 执行 ====================
main "$@"
