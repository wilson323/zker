#!/bin/bash
# deploy-all.sh - 监控系统一键部署脚本
# 职责: 协调调用各个组件的安装脚本，实现一键部署
# 版本: v1.0.0
# 作者: ZKER DevOps Team

set -e  # 遇到错误立即退出

# ==================== 配置变量 ====================
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

# ==================== 颜色输出 ====================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# ==================== 显示Banner ====================
show_banner() {
    echo ""
    echo "=========================================="
    echo "  ZKER 企业级监控系统 - 一键部署"
    echo "=========================================="
    echo "版本: v1.0.0"
    echo "日期: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "=========================================="
    echo ""
}

# ==================== 检查前置条件 ====================
check_prerequisites() {
    log_step "检查前置条件..."

    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，请先安装Docker"
        exit 1
    fi
    log_info "✅ Docker已安装"

    # 检查Docker Compose
    if ! command -v docker compose &> /dev/null; then
        log_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    log_info "✅ Docker Compose已安装"

    # 检查网络
    if ! docker network ls --format '{{.Name}}' | grep -q '^zker-network$'; then
        log_warn "zker-network网络不存在，将自动创建"
        docker network create zker-network 2>/dev/null || true
    fi
    log_info "✅ Docker网络已就绪"

    echo ""
}

# ==================== 安装Prometheus ====================
install_prometheus() {
    log_step "安装Prometheus..."
    echo ""

    if [ -f "$SCRIPT_DIR/install-prometheus.sh" ]; then
        bash "$SCRIPT_DIR/install-prometheus.sh"
        echo ""
    else
        log_error "install-prometheus.sh脚本未找到"
        exit 1
    fi
}

# ==================== 安装Grafana ====================
install_grafana() {
    log_step "安装Grafana..."
    echo ""

    if [ -f "$SCRIPT_DIR/install-grafana.sh" ]; then
        bash "$SCRIPT_DIR/install-grafana.sh"
        echo ""
    else
        log_error "install-grafana.sh脚本未找到"
        exit 1
    fi
}

# ==================== 安装AlertManager ====================
install_alertmanager() {
    log_step "安装AlertManager..."
    echo ""

    if [ -f "$SCRIPT_DIR/install-alertmanager.sh" ]; then
        bash "$SCRIPT_DIR/install-alertmanager.sh"
        echo ""
    else
        log_error "install-alertmanager.sh脚本未找到"
        exit 1
    fi
}

# ==================== 安装Exporters ====================
install_exporters() {
    log_step "安装Exporters..."
    echo ""

    if [ -f "$SCRIPT_DIR/install-exporters.sh" ]; then
        bash "$SCRIPT_DIR/install-exporters.sh" all
        echo ""
    else
        log_error "install-exporters.sh脚本未找到"
        exit 1
    fi
}

# ==================== 等待服务就绪 ====================
wait_for_services() {
    log_step "等待所有服务启动..."
    echo ""

    log_info "等待Prometheus启动..."
    for i in {1..30}; do
        if curl -f -s http://localhost:9090/-/healthy > /dev/null 2>&1; then
            log_info "✅ Prometheus已就绪"
            break
        fi
        if [ $i -eq 30 ]; then
            log_warn "⚠️  Prometheus启动超时"
        fi
        sleep 2
    done

    log_info "等待Grafana启动..."
    for i in {1..30}; do
        if curl -f -s http://localhost:3000/api/health > /dev/null 2>&1; then
            log_info "✅ Grafana已就绪"
            break
        fi
        if [ $i -eq 30 ]; then
            log_warn "⚠️  Grafana启动超时"
        fi
        sleep 2
    done

    log_info "等待AlertManager启动..."
    for i in {1..30}; do
        if curl -f -s http://localhost:9093/-/healthy > /dev/null 2>&1; then
            log_info "✅ AlertManager已就绪"
            break
        fi
        if [ $i -eq 30 ]; then
            log_warn "⚠️  AlertManager启动超时"
        fi
        sleep 2
    done

    echo ""
}

# ==================== 验证部署 ====================
verify_deployment() {
    log_step "验证部署..."
    echo ""

    # 检查所有容器状态
    log_info "容器状态:"
    docker ps --filter "name=zker-" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    echo ""

    # 检查服务端点
    log_info "服务端点检查:"
    if curl -f -s http://localhost:9090/-/healthy > /dev/null 2>&1; then
        log_info "✅ Prometheus: http://localhost:9090"
    else
        log_warn "⚠️  Prometheus: 未响应"
    fi

    if curl -f -s http://localhost:3000/api/health > /dev/null 2>&1; then
        log_info "✅ Grafana: http://localhost:3000"
    else
        log_warn "⚠️  Grafana: 未响应"
    fi

    if curl -f -s http://localhost:9093/-/healthy > /dev/null 2>&1; then
        log_info "✅ AlertManager: http://localhost:9093"
    else
        log_warn "⚠️  AlertManager: 未响应"
    fi

    echo ""
}

# ==================== 显示部署摘要 ====================
show_summary() {
    echo ""
    echo "=========================================="
    echo "  部署完成！"
    echo "=========================================="
    echo ""
    echo "📍 访问地址:"
    echo "  Prometheus: http://localhost:9090"
    echo "  Grafana:     http://localhost:3000 (admin/admin)"
    echo "  AlertManager: http://localhost:9093"
    echo ""
    echo "📊 监控指标端点:"
    echo "  Node Exporter:    http://localhost:9100/metrics"
    echo "  MySQL Exporter:   http://localhost:9104/metrics"
    echo "  Redis Exporter:   http://localhost:9121/metrics"
    echo ""
    echo "🔧 管理命令:"
    echo "  查看所有容器: docker ps -a | grep zker-"
    echo "  查看服务日志: docker logs -f <container_name>"
    echo "  停止所有服务: docker stop zker-prometheus zker-grafana zker-alertmanager zker-node-exporter zker-mysqld-exporter zker-redis-exporter"
    echo "  启动所有服务: docker start zker-prometheus zker-grafana zker-alertmanager zker-node-exporter zker-mysqld-exporter zker-redis-exporter"
    echo ""
    echo "📖 配置文件位置:"
    echo "  Prometheus配置: /etc/prometheus/prometheus.yml"
    echo "  Grafana配置: /etc/grafana/provisioning/"
    echo "  AlertManager配置: /etc/alertmanager/alertmanager.yml"
    echo ""
    echo "⚠️  后续步骤:"
    echo "  1. 登录Grafana (admin/admin) 并修改密码"
    echo "  2. 配置告警通知渠道 (Slack/PagerDuty/Email)"
    echo "  3. 导入自定义监控仪表盘"
    echo "  4. 根据需要调整告警阈值"
    echo ""
    echo "=========================================="
}

# ==================== 显示使用帮助 ====================
show_help() {
    cat << EOF
用法: $0 [选项]

选项:
    all         安装所有组件 (默认)
    prometheus  仅安装Prometheus
    grafana     仅安装Grafana
    alertmanager 仅安装AlertManager
    exporters   仅安装Exporters
    help        显示此帮助信息

示例:
    $0              # 安装所有组件
    $0 prometheus   # 仅安装Prometheus
    $0 help         # 显示帮助

环境变量:
    PROMETHEUS_PORT          Prometheus端口 (默认: 9090)
    GRAFANA_PORT             Grafana端口 (默认: 3000)
    ALERTMANAGER_PORT        AlertManager端口 (默认: 9093)
    GRAFANA_ADMIN_USER       Grafana管理员用户 (默认: admin)
    GRAFANA_ADMIN_PASSWORD   Grafana管理员密码 (默认: admin)
    MYSQL_EXPORTER_DSN       MySQL数据源 (默认: root:root@(zker-mysql:3306)/)

EOF
}

# ==================== 主函数 ====================
main() {
    show_banner

    # 解析参数
    case "${1:-all}" in
        all)
            check_prerequisites
            install_prometheus
            install_grafana
            install_alertmanager
            install_exporters
            wait_for_services
            verify_deployment
            show_summary
            ;;
        prometheus)
            check_prerequisites
            install_prometheus
            ;;
        grafana)
            check_prerequisites
            install_grafana
            ;;
        alertmanager)
            check_prerequisites
            install_alertmanager
            ;;
        exporters)
            check_prerequisites
            install_exporters
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
}

# ==================== 执行 ====================
main "$@"
