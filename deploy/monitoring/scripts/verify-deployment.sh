#!/bin/bash
# verify-deployment.sh - 监控系统部署验证脚本
# 职责: 验证所有监控组件的部署状态和健康度
# 版本: v1.0.0
# 作者: ZKER DevOps Team

set -e  # 遇到错误立即退出

# ==================== 颜色输出 ====================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[✓]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[!]${NC} $1"
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# ==================== 检查容器状态 ====================
check_containers() {
    log_step "检查容器状态..."
    echo ""

    local all_running=true

    # 检查Prometheus
    if docker ps --format '{{.Names}}' | grep -q '^zker-prometheus$'; then
        log_info "Prometheus容器运行正常"
    else
        log_error "Prometheus容器未运行"
        all_running=false
    fi

    # 检查Grafana
    if docker ps --format '{{.Names}}' | grep -q '^zker-grafana$'; then
        log_info "Grafana容器运行正常"
    else
        log_error "Grafana容器未运行"
        all_running=false
    fi

    # 检查AlertManager
    if docker ps --format '{{.Names}}' | grep -q '^zker-alertmanager$'; then
        log_info "AlertManager容器运行正常"
    else
        log_error "AlertManager容器未运行"
        all_running=false
    fi

    # 检查Exporters
    if docker ps --format '{{.Names}}' | grep -q '^zker-node-exporter$'; then
        log_info "Node Exporter容器运行正常"
    else
        log_warn "Node Exporter容器未运行"
    fi

    if docker ps --format '{{.Names}}' | grep -q '^zker-mysqld-exporter$'; then
        log_info "MySQL Exporter容器运行正常"
    else
        log_warn "MySQL Exporter容器未运行"
    fi

    if docker ps --format '{{.Names}}' | grep -q '^zker-redis-exporter$'; then
        log_info "Redis Exporter容器运行正常"
    else
        log_warn "Redis Exporter容器未运行"
    fi

    echo ""

    if [ "$all_running" = false ]; then
        return 1
    fi
}

# ==================== 检查HTTP端点 ====================
check_endpoints() {
    log_step "检查HTTP端点..."
    echo ""

    # 检查Prometheus
    if curl -f -s http://localhost:9090/-/healthy > /dev/null 2>&1; then
        log_info "Prometheus HTTP端点响应正常 (http://localhost:9090)"
    else
        log_error "Prometheus HTTP端点未响应"
    fi

    # 检查Grafana
    if curl -f -s http://localhost:3000/api/health > /dev/null 2>&1; then
        log_info "Grafana HTTP端点响应正常 (http://localhost:3000)"
    else
        log_error "Grafana HTTP端点未响应"
    fi

    # 检查AlertManager
    if curl -f -s http://localhost:9093/-/healthy > /dev/null 2>&1; then
        log_info "AlertManager HTTP端点响应正常 (http://localhost:9093)"
    else
        log_error "AlertManager HTTP端点未响应"
    fi

    # 检查Node Exporter
    if curl -f -s http://localhost:9100/metrics > /dev/null 2>&1; then
        log_info "Node Exporter HTTP端点响应正常 (http://localhost:9100/metrics)"
    else
        log_warn "Node Exporter HTTP端点未响应"
    fi

    # 检查MySQL Exporter
    if curl -f -s http://localhost:9104/metrics > /dev/null 2>&1; then
        log_info "MySQL Exporter HTTP端点响应正常 (http://localhost:9104/metrics)"
    else
        log_warn "MySQL Exporter HTTP端点未响应"
    fi

    # 检查Redis Exporter
    if curl -f -s http://localhost:9121/metrics > /dev/null 2>&1; then
        log_info "Redis Exporter HTTP端点响应正常 (http://localhost:9121/metrics)"
    else
        log_warn "Redis Exporter HTTP端点未响应"
    fi

    echo ""
}

# ==================== 检查配置文件 ====================
check_configs() {
    log_step "检查配置文件..."
    echo ""

    # 检查Prometheus配置
    if docker exec zker-prometheus promtool check config /etc/prometheus/prometheus.yml > /dev/null 2>&1; then
        log_info "Prometheus配置文件语法正确"
    else
        log_warn "Prometheus配置文件可能存在问题"
    fi

    # 检查AlertManager配置
    if docker exec zker-alertmanager amtool check-config /etc/alertmanager/alertmanager.yml > /dev/null 2>&1; then
        log_info "AlertManager配置文件语法正确"
    else
        log_warn "AlertManager配置文件可能存在问题"
    fi

    echo ""
}

# ==================== 检查Prometheus targets ====================
check_prometheus_targets() {
    log_step "检查Prometheus抓取目标..."
    echo ""

    local targets_json
    targets_json=$(curl -s http://localhost:9090/api/v1/targets 2>/dev/null)

    if [ -z "$targets_json" ]; then
        log_error "无法获取Prometheus targets信息"
        echo ""
        return 1
    fi

    # 解析targets状态
    local up_count
    local down_count
    up_count=$(echo "$targets_json" | grep -o '"health":"up"' | wc -l)
    down_count=$(echo "$targets_json" | grep -o '"health":"down"' | wc -l)

    log_info "Prometheus抓取目标状态: UP=${up_count}, DOWN=${down_count}"

    if [ "$down_count" -gt 0 ]; then
        log_warn "有${down_count}个目标未响应，请检查"
    fi

    echo ""
}

# ==================== 检查Grafana数据源 ====================
check_grafana_datasources() {
    log_step "检查Grafana数据源..."
    echo ""

    if ! docker exec zker-grafana grafana-cli admin list-datasources > /dev/null 2>&1; then
        log_warn "无法列出Grafana数据源"
    else
        log_info "Grafana数据源已配置"
    fi

    echo ""
}

# ==================== 检查告警规则 ====================
check_alert_rules() {
    log_step "检查告警规则..."
    echo ""

    local rules_json
    rules_json=$(curl -s http://localhost:9090/api/v1/rules 2>/dev/null)

    if [ -z "$rules_json" ]; then
        log_error "无法获取Prometheus告警规则"
        echo ""
        return 1
    fi

    # 统计规则数量
    local rule_count
    rule_count=$(echo "$rules_json" | grep -o '"name":' | wc -l)

    log_info "已加载${rule_count}条告警规则"

    # 检查活跃告警
    local alerts_json
    alerts_json=$(curl -s http://localhost:9090/api/v1/alerts 2>/dev/null)

    if [ -n "$alerts_json" ]; then
        local active_alerts
        active_alerts=$(echo "$alerts_json" | grep -o '"state":"firing"' | wc -l)

        if [ "$active_alerts" -gt 0 ]; then
            log_warn "当前有${active_alerts}个活跃告警"
        else
            log_info "当前无活跃告警"
        fi
    fi

    echo ""
}

# ==================== 显示验证摘要 ====================
show_summary() {
    echo ""
    echo "=========================================="
    echo "  验证完成"
    echo "=========================================="
    echo ""
    echo "访问地址:"
    echo "  Prometheus:    http://localhost:9090"
    echo "  Grafana:       http://localhost:3000 (admin/admin)"
    echo "  AlertManager:  http://localhost:9093"
    echo ""
    echo "监控指标端点:"
    echo "  Node Exporter:    http://localhost:9100/metrics"
    echo "  MySQL Exporter:   http://localhost:9104/metrics"
    echo "  Redis Exporter:   http://localhost:9121/metrics"
    echo ""
    echo "=========================================="
}

# ==================== 主函数 ====================
main() {
    echo ""
    echo "=========================================="
    echo "  ZKER 监控系统部署验证"
    echo "=========================================="
    echo ""

    local exit_code=0

    # 执行各项检查
    check_containers || exit_code=$?
    check_endpoints || exit_code=$?
    check_configs || exit_code=$?
    check_prometheus_targets || exit_code=$?
    check_grafana_datasources || exit_code=$?
    check_alert_rules || exit_code=$?

    show_summary

    exit $exit_code
}

# ==================== 执行 ====================
main "$@"
