#!/bin/bash
# install-grafana.sh - Grafana安装脚本
# 职责: 只负责Grafana的安装、配置和启动
# 版本: v1.0.0
# 作者: ZKER DevOps Team

set -e  # 遇到错误立即退出

# ==================== 配置变量 ====================
GRAFANA_VERSION="10.0.0"
GRAFANA_PORT="${GRAFANA_PORT:-3000}"
GRAFANA_DIR="/opt/grafana"
CONFIG_DIR="/etc/grafana"
DATA_DIR="/var/lib/grafana"
PROVISIONING_DIR="/etc/grafana/provisioning"
DASHBOARDS_DIR="/var/lib/grafana/dashboards"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

# 默认管理员凭据
ADMIN_USER="${GRAFANA_ADMIN_USER:-admin}"
ADMIN_PASSWORD="${GRAFANA_ADMIN_PASSWORD:-admin}"

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
    if netstat -tuln 2>/dev/null | grep -q ":${GRAFANA_PORT} "; then
        log_warn "端口${GRAFANA_PORT}已被占用，请检查是否有其他Grafana实例运行"
    fi

    log_info "前置条件检查通过"
}

# ==================== 创建目录 ====================
create_directories() {
    log_info "创建Grafana目录..."

    sudo mkdir -p "$DATA_DIR"
    sudo mkdir -p "$PROVISIONING_DIR/datasources"
    sudo mkdir -p "$PROVISIONING_DIR/dashboards"
    sudo mkdir -p "$DASHBOARDS_DIR"

    log_info "目录创建完成"
}

# ==================== 复制配置文件 ====================
copy_configs() {
    log_info "复制Grafana配置文件..."

    # 复制数据源配置
    if [ -f "$PROJECT_ROOT/docker/volumes/monitoring/grafana/provisioning/datasources/prometheus.yml" ]; then
        sudo cp "$PROJECT_ROOT/docker/volumes/monitoring/grafana/provisioning/datasources/prometheus.yml" \
            "$PROVISIONING_DIR/datasources/prometheus.yml"
        log_info "Prometheus数据源配置复制完成"
    else
        log_warn "Prometheus数据源配置未找到，跳过"
    fi

    # 复制仪表盘配置
    if [ -f "$PROJECT_ROOT/docker/volumes/monitoring/grafana/provisioning/dashboards/org-dashboard.yml" ]; then
        sudo cp "$PROJECT_ROOT/docker/volumes/monitoring/grafana/provisioning/dashboards/org-dashboard.yml" \
            "$PROVISIONING_DIR/dashboards/org-dashboard.yml"
        log_info "仪表盘配置复制完成"
    else
        log_warn "仪表盘配置未找到，跳过"
    fi

    # 复制仪表盘JSON文件
    if [ -d "$PROJECT_ROOT/deploy/monitoring/dashboards" ]; then
        sudo cp -r "$PROJECT_ROOT/deploy/monitoring/dashboards"/*.json "$DASHBOARDS_DIR/" 2>/dev/null || true
        log_info "仪表盘JSON文件复制完成"
    else
        log_warn "仪表盘JSON文件未找到，跳过"
    fi

    # 设置权限
    sudo chown -R 472:472 "$DATA_DIR"
    sudo chown -R root:root "$PROVISIONING_DIR"
    sudo chmod 644 "$PROVISIONING_DIR"/*/*.yml 2>/dev/null || true

    log_info "配置文件复制完成"
}

# ==================== 启动Grafana ====================
start_grafana() {
    log_info "启动Grafana容器..."

    # 检查容器是否已存在
    if docker ps -a --format '{{.Names}}' | grep -q '^zker-grafana$'; then
        log_warn "Grafana容器已存在，正在删除..."
        docker stop zker-grafana 2>/dev/null || true
        docker rm zker-grafana 2>/dev/null || true
    fi

    # 启动Grafana容器
    docker run -d \
        --name zker-grafana \
        --restart unless-stopped \
        -p "${GRAFANA_PORT}:3000" \
        -e "GF_SECURITY_ADMIN_USER=${ADMIN_USER}" \
        -e "GF_SECURITY_ADMIN_PASSWORD=${ADMIN_PASSWORD}" \
        -e "GF_INSTALL_PLUGINS=grafana-piechart-panel" \
        -e "GF_SERVER_ROOT_URL=http://localhost:${GRAFANA_PORT}" \
        -e "GF_USERS_ALLOW_SIGN_UP=false" \
        -v "$DATA_DIR:/var/lib/grafana" \
        -v "$PROVISIONING_DIR:/etc/grafana/provisioning:ro" \
        -v "$DASHBOARDS_DIR:/var/lib/grafana/dashboards:ro" \
        grafana/grafana:"${GRAFANA_VERSION}"

    log_info "Grafana容器启动成功"
}

# ==================== 验证安装 ====================
verify_installation() {
    log_info "验证Grafana安装..."

    # 等待Grafana启动
    sleep 10

    # 检查容器状态
    if docker ps --format '{{.Names}}' | grep -q '^zker-grafana$'; then
        log_info "✅ Grafana容器运行正常"
    else
        log_error "❌ Grafana容器未运行"
        docker logs zker-grafana
        exit 1
    fi

    # 检查HTTP端点
    if curl -f -s http://localhost:${GRAFANA_PORT}/api/health > /dev/null 2>&1; then
        log_info "✅ Grafana HTTP端点响应正常"
    else
        log_warn "⚠️  Grafana HTTP端点未响应，可能需要更多时间启动"
    fi

    log_info "Grafana安装验证完成"
}

# ==================== 显示访问信息 ====================
show_access_info() {
    echo ""
    echo "=========================================="
    echo "Grafana 安装完成！"
    echo "=========================================="
    echo "📍 访问地址: http://localhost:${GRAFANA_PORT}"
    echo "👤 管理员用户: ${ADMIN_USER}"
    echo "🔑 管理员密码: ${ADMIN_PASSWORD}"
    echo "📖 数据目录: $DATA_DIR"
    echo "📁 配置目录: $PROVISIONING_DIR"
    echo ""
    echo "⚠️  首次登录后请立即修改管理员密码！"
    echo ""
    echo "常用命令:"
    echo "  查看日志: docker logs -f zker-grafana"
    echo "  停止服务: docker stop zker-grafana"
    echo "  启动服务: docker start zker-grafana"
    echo "  重启服务: docker restart zker-grafana"
    echo "  删除容器: docker rm -f zker-grafana"
    echo ""
    echo "=========================================="
}

# ==================== 主函数 ====================
main() {
    echo "=========================================="
    echo "Grafana 安装脚本"
    echo "版本: ${GRAFANA_VERSION}"
    echo "=========================================="
    echo ""

    check_prerequisites
    create_directories
    copy_configs
    start_grafana
    verify_installation
    show_access_info
}

# ==================== 执行 ====================
main "$@"
