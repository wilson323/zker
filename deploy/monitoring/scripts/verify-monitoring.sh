#!/bin/bash
# ZKER 监控系统验证脚本
# 版本: v1.0.0
# 用途: 验证 Prometheus、Grafana、AlertManager 和 Exporters 的健康状态

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
PROMETHEUS_URL="${PROMETHEUS_URL:-http://localhost:9090}"
GRAFANA_URL="${GRAFANA_URL:-http://localhost:3000}"
ALERTMANAGER_URL="${ALERTMANAGER_URL:-http://localhost:9093}"
API_METRICS_URL="${API_METRICS_URL:-http://localhost:8888/metrics}"

# 测试计数器
PASS=0
FAIL=0
WARN=0

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
    ((PASS++))
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
    ((FAIL++))
}

log_warn() {
    echo -e "${YELLOW}[!]${NC} $1"
    ((WARN++))
}

# HTTP健康检查函数
check_http() {
    local url=$1
    local name=$2
    local expected_code=${3:-200}

    log_info "检查 $name ($url)..."

    if curl -s -f -o /dev/null -w "%{http_code}" "$url" | grep -q "$expected_code"; then
        log_success "$name 健康检查通过"
        return 0
    else
        log_error "$name 健康检查失败"
        return 1
    fi
}

# Prometheus指标检查
check_prometheus_metrics() {
    log_info "检查 Prometheus 指标采集..."

    local metrics=(
        "http_requests_total"
        "http_request_duration_seconds"
        "quota_usage"
        "bot_invocation_total"
        "db_query_duration_seconds"
    )

    for metric in "${metrics[@]}"; do
        local count=$(curl -s "${PROMETHEUS_URL}/api/v1/query?query=${metric}" | jq -r '.data.result | length')

        if [ "$count" -gt 0 ]; then
            log_success "指标 ${metric} 存在 (${count} 条数据)"
        else
            log_warn "指标 ${metric} 不存在或无数据"
        fi
    done
}

# Prometheus Targets检查
check_prometheus_targets() {
    log_info "检查 Prometheus Targets 状态..."

    local targets_output=$(curl -s "${PROMETHEUS_URL}/api/v1/targets")
    local total=$(echo "$targets_output" | jq -r '.data.activeTargets | length')
    local up=$(echo "$targets_output" | jq -r '[.data.activeTargets[] | select(.health=="up")] | length')

    log_info "Targets: $up/$total 在线"

    echo "$targets_output" | jq -r '.data.activeTargets[] | "\(.labels.job) - \(.health)"' | while read -r line; do
        if echo "$line" | grep -q "up"; then
            log_success "$line"
        else
            log_error "$line"
        fi
    done
}

# Grafana仪表板检查
check_grafana_dashboards() {
    log_info "检查 Grafana 仪表板..."

    local dashboards=$(curl -s "${GRAFANA_URL}/api/search?query=zker" \
        -u "admin:admin" | jq -r '.[] | .uid')

    if [ -z "$dashboards" ]; then
        log_warn "未找到 ZKER 相关仪表板"
    else
        echo "$dashboards" | while read -r uid; do
            log_success "仪表板 UID: $uid"
        done
    fi
}

# AlertManager告警检查
check_alertmanager_alerts() {
    log_info "检查 AlertManager 告警状态..."

    local alerts=$(curl -s "${ALERTMANAGER_URL}/api/v2/alerts")
    local count=$(echo "$alerts" | jq 'length')

    if [ "$count" -eq 0 ]; then
        log_success "当前没有活跃告警"
    else
        log_warn "当前有 $count 个活跃告警"
        echo "$alerts" | jq -r '.[] | "\(.labels.alertname) - \(.labels.severity)"'
    fi
}

# API /metrics端点检查
check_api_metrics() {
    log_info "检查 API /metrics 端点..."

    local metrics_content=$(curl -s "$API_METRICS_URL")

    if echo "$metrics_content" | grep -q "http_requests_total"; then
        log_success "http_requests_total 指标已暴露"
    else
        log_error "http_requests_total 指标未找到"
    fi

    if echo "$metrics_content" | grep -q "http_request_duration_seconds"; then
        log_success "http_request_duration_seconds 指标已暴露"
    else
        log_error "http_request_duration_seconds 指标未找到"
    fi
}

# 性能测试（模拟高QPS）
performance_test() {
    log_info "执行性能测试（模拟高QPS）..."

    local duration=10
    local concurrent=100

    log_info "启动 ${concurrent} 并发请求，持续 ${duration} 秒..."

    for i in $(seq 1 $concurrent); do
        (
            for j in $(seq 1 $((duration * 10))); do
                curl -s "$API_METRICS_URL" > /dev/null
                sleep 0.1
            done
        ) &
    done

    wait

    log_success "性能测试完成，检查 Prometheus 采集的 QPS..."

    sleep 5

    local qps=$(curl -s "${PROMETHEUS_URL}/api/v1/query?query=rate(http_requests_total[5m])" \
        | jq -r '.data.result[0].value[1] // "0"')

    log_info "当前 QPS: ${qps}"
}

# 主函数
main() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  ZKER 监控系统验证脚本${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""

    # 基础健康检查
    log_info "=== 基础健康检查 ==="
    check_http "$PROMETHEUS_URL/-/healthy" "Prometheus"
    check_http "$GRAFANA_URL/api/health" "Grafana"
    check_http "$ALERTMANAGER_URL/-/healthy" "AlertManager"
    echo ""

    # Prometheus检查
    log_info "=== Prometheus 检查 ==="
    check_prometheus_targets
    check_prometheus_metrics
    echo ""

    # Grafana检查
    log_info "=== Grafana 检查 ==="
    check_grafana_dashboards
    echo ""

    # AlertManager检查
    log_info "=== AlertManager 检查 ==="
    check_alertmanager_alerts
    echo ""

    # API Metrics检查
    log_info "=== API Metrics 检查 ==="
    check_api_metrics
    echo ""

    # 性能测试（可选）
    if [ "${RUN_PERFORMANCE_TEST:-false}" = "true" ]; then
        log_info "=== 性能测试 ==="
        performance_test
        echo ""
    fi

    # 汇总
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  验证结果汇总${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo -e "${GREEN}通过: $PASS${NC}"
    echo -e "${YELLOW}警告: $WARN${NC}"
    echo -e "${RED}失败: $FAIL${NC}"
    echo ""

    if [ $FAIL -eq 0 ]; then
        log_success "所有核心检查通过！"
        exit 0
    else
        log_error "有 $FAIL 项检查失败，请查看日志"
        exit 1
    fi
}

# 执行主函数
main "$@"
