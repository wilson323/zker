#!/bin/bash
# install-alertmanager.sh - AlertManager安装脚本
# 职责: 只负责AlertManager的安装、配置和启动
# 版本: v1.0.0
# 作者: ZKER DevOps Team

set -e  # 遇到错误立即退出

# ==================== 配置变量 ====================
ALERTMANAGER_VERSION="v0.26.0"
ALERTMANAGER_PORT="${ALERTMANAGER_PORT:-9093}"
ALERTMANAGER_DIR="/opt/alertmanager"
CONFIG_DIR="/etc/alertmanager"
DATA_DIR="/var/lib/alertmanager"
TEMPLATES_DIR="/etc/alertmanager/templates"
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
    if netstat -tuln 2>/dev/null | grep -q ":${ALERTMANAGER_PORT} "; then
        log_warn "端口${ALERTMANAGER_PORT}已被占用，请检查是否有其他AlertManager实例运行"
    fi

    log_info "前置条件检查通过"
}

# ==================== 创建目录 ====================
create_directories() {
    log_info "创建AlertManager目录..."

    sudo mkdir -p "$CONFIG_DIR"
    sudo mkdir -p "$DATA_DIR"
    sudo mkdir -p "$TEMPLATES_DIR"

    log_info "目录创建完成"
}

# ==================== 复制配置文件 ====================
copy_configs() {
    log_info "复制AlertManager配置文件..."

    # 复制alertmanager.yml
    if [ -f "$PROJECT_ROOT/deploy/monitoring/alertmanager/config.yml" ]; then
        sudo cp "$PROJECT_ROOT/deploy/monitoring/alertmanager/config.yml" "$CONFIG_DIR/alertmanager.yml"
        log_info "alertmanager.yml复制完成"
    elif [ -f "$PROJECT_ROOT/docker/volumes/monitoring/alertmanager/alertmanager.yml" ]; then
        sudo cp "$PROJECT_ROOT/docker/volumes/monitoring/alertmanager/alertmanager.yml" "$CONFIG_DIR/alertmanager.yml"
        log_info "alertmanager.yml复制完成"
    else
        log_warn "alertmanager.yml未找到，使用默认配置"
    fi

    # 设置权限
    sudo chown -R root:root "$CONFIG_DIR"
    sudo chmod 644 "$CONFIG_DIR"/*.yml 2>/dev/null || true

    log_info "配置文件复制完成"
}

# ==================== 启动AlertManager ====================
start_alertmanager() {
    log_info "启动AlertManager容器..."

    # 检查容器是否已存在
    if docker ps -a --format '{{.Names}}' | grep -q '^zker-alertmanager$'; then
        log_warn "AlertManager容器已存在，正在删除..."
        docker stop zker-alertmanager 2>/dev/null || true
        docker rm zker-alertmanager 2>/dev/null || true
    fi

    # 启动AlertManager容器
    docker run -d \
        --name zker-alertmanager \
        --restart unless-stopped \
        -p "${ALERTMANAGER_PORT}:9093" \
        -v "$CONFIG_DIR/alertmanager.yml:/etc/alertmanager/alertmanager.yml:ro" \
        -v "$TEMPLATES_DIR:/etc/alertmanager/templates:ro" \
        -v "$DATA_DIR:/alertmanager" \
        prom/alertmanager:"${ALERTMANAGER_VERSION}" \
        --config.file=/etc/alertmanager/alertmanager.yml \
        --storage.path=/alertmanager \
        --web.external-url=http://localhost:${ALERTMANAGER_PORT}

    log_info "AlertManager容器启动成功"
}

# ==================== 验证安装 ====================
verify_installation() {
    log_info "验证AlertManager安装..."

    # 等待AlertManager启动
    sleep 5

    # 检查容器状态
    if docker ps --format '{{.Names}}' | grep -q '^zker-alertmanager$'; then
        log_info "✅ AlertManager容器运行正常"
    else
        log_error "❌ AlertManager容器未运行"
        docker logs zker-alertmanager
        exit 1
    fi

    # 检查HTTP端点
    if curl -f -s http://localhost:${ALERTMANAGER_PORT}/-/healthy > /dev/null 2>&1; then
        log_info "✅ AlertManager HTTP端点响应正常"
    else
        log_warn "⚠️  AlertManager HTTP端点未响应，可能需要更多时间启动"
    fi

    # 检查配置文件
    if docker exec zker-alertmanager amtool check-config /etc/alertmanager/alertmanager.yml > /dev/null 2>&1; then
        log_info "✅ AlertManager配置文件语法正确"
    else
        log_warn "⚠️  AlertManager配置文件可能存在问题"
    fi

    log_info "AlertManager安装验证完成"
}

# ==================== 显示访问信息 ====================
show_access_info() {
    echo ""
    echo "=========================================="
    echo "AlertManager 安装完成！"
    echo "=========================================="
    echo "📍 访问地址: http://localhost:${ALERTMANAGER_PORT}"
    echo "📖 配置文件: $CONFIG_DIR/alertmanager.yml"
    echo "💾 数据目录: $DATA_DIR"
    echo ""
    echo "常用命令:"
    echo "  查看日志: docker logs -f zker-alertmanager"
    echo "  停止服务: docker stop zker-alertmanager"
    echo "  启动服务: docker start zker-alertmanager"
    echo "  重启服务: docker restart zker-alertmanager"
    echo "  删除容器: docker rm -f zker-alertmanager"
    echo ""
    echo "配置验证:"
    echo "  docker exec zker-alertmanager amtool check-config /etc/alertmanager/alertmanager.yml"
    echo ""
    echo "=========================================="
}

# ==================== 主函数 ====================
main() {
    echo "=========================================="
    echo "AlertManager 安装脚本"
    echo "版本: ${ALERTMANAGER_VERSION}"
    echo "=========================================="
    echo ""

    check_prerequisites
    create_directories
    copy_configs
    start_alertmanager
    verify_installation
    show_access_info
}

# ==================== 执行 ====================
main "$@"
