#!/bin/bash
# 计费系统监控和日志集成验证脚本
# 版本: v1.0.0

set -e

echo "=========================================="
echo "计费系统监控和日志集成验证"
echo "=========================================="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 统计变量
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0

# 检查函数
check_service() {
    local service_name=$1
    local service_url=$2
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))

    echo -n "检查 $service_name... "
    if curl -s -f "$service_url" > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        return 0
    else
        echo -e "${RED}✗${NC}"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        return 1
    fi
}

check_file() {
    local file_path=$1
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))

    echo -n "检查文件 $file_path... "
    if [ -f "$file_path" ]; then
        echo -e "${GREEN}✓${NC}"
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        return 0
    else
        echo -e "${RED}✗${NC}"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        return 1
    fi
}

check_metrics_endpoint() {
    local metrics_url=$1
    local metric_pattern=$2
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))

    echo -n "检查指标端点 $metrics_url... "
    if curl -s "$metrics_url" | grep -q "$metric_pattern"; then
        echo -e "${GREEN}✓${NC}"
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        return 0
    else
        echo -e "${RED}✗${NC}"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        return 1
    fi
}

check_prometheus_rule() {
    local rule_file=$1
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))

    echo -n "验证Prometheus告警规则 $rule_file... "
    if promtool check rules "$rule_file" 2>&1 | grep -q "SUCCESS"; then
        echo -e "${GREEN}✓${NC}"
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        return 0
    else
        echo -e "${RED}✗${NC}"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        return 1
    fi
}

echo "1. 检查源代码文件"
echo "----------------------------------------"
check_file "backend/domain/billing/service/metrics.go"
check_file "backend/domain/billing/service/logging.go"
check_file "backend/domain/billing/service/tracing.go"
check_file "backend/domain/billing/service/monitoring_integration.go"
echo ""

echo "2. 检查配置文件"
echo "----------------------------------------"
check_file "deploy/monitoring/dashboards/billing-dashboard.json"
check_file "deploy/monitoring/alerts/billing-alerts.yml"
check_file "deploy/logging/filebeat-billing.yml"
echo ""

echo "3. 检查文档"
echo "----------------------------------------"
check_file "docs/企业级功能完善与统一性设计方案/BILLING_MONITORING_GUIDE.md"
echo ""

echo "4. 检查服务状态（假设本地运行）"
echo "----------------------------------------"
check_service "Prometheus" "http://localhost:9090/-/healthy"
check_service "Grafana" "http://localhost:3000/api/health"
check_service "Elasticsearch" "http://localhost:9200/_cluster/health"
check_service "Jaeger" "http://localhost:16686/api/services"
echo ""

echo "5. 检查Prometheus指标端点"
echo "----------------------------------------"
check_metrics_endpoint "http://localhost:9090/api/v1/query?query=up" "data"
echo ""

echo "6. 验证告警规则语法"
echo "----------------------------------------"
if command -v promtool &> /dev/null; then
    check_prometheus_rule "deploy/monitoring/alerts/billing-alerts.yml"
else
    echo -e "${YELLOW}⚠ promtool未安装，跳过告警规则验证${NC}"
fi
echo ""

echo "7. 检查指标定义完整性"
echo "----------------------------------------"
TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
echo -n "检查metrics.go中的指标数量... "
METRIC_COUNT=$(grep -c "promauto\." backend/domain/billing/service/metrics.go || true)
if [ "$METRIC_COUNT" -ge 20 ]; then
    echo -e "${GREEN}✓ (找到 $METRIC_COUNT 个指标)${NC}"
    PASSED_CHECKS=$((PASSED_CHECKS + 1))
else
    echo -e "${RED}✗ (仅找到 $METRIC_COUNT 个指标，期望至少20个)${NC}"
    FAILED_CHECKS=$((FAILED_CHECKS + 1))
fi
echo ""

echo "8. 检查日志函数完整性"
echo "----------------------------------------"
TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
echo -n "检查logging.go中的日志函数数量... "
LOG_FUNC_COUNT=$(grep -c "^func Log" backend/domain/billing/service/logging.go || true)
if [ "$LOG_FUNC_COUNT" -ge 15 ]; then
    echo -e "${GREEN}✓ (找到 $LOG_FUNC_COUNT 个函数)${NC}"
    PASSED_CHECKS=$((PASSED_CHECKS + 1))
else
    echo -e "${RED}✗ (仅找到 $LOG_FUNC_COUNT 个函数，期望至少15个)${NC}"
    FAILED_CHECKS=$((FAILED_CHECKS + 1))
fi
echo ""

echo "9. 检查追踪函数完整性"
echo "----------------------------------------"
TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
echo -n "检查tracing.go中的追踪函数数量... "
TRACE_FUNC_COUNT=$(grep -c "^func Start" backend/domain/billing/service/tracing.go || true)
if [ "$TRACE_FUNC_COUNT" -ge 10 ]; then
    echo -e "${GREEN}✓ (找到 $TRACE_FUNC_COUNT 个函数)${NC}"
    PASSED_CHECKS=$((PASSED_CHECKS + 1))
else
    echo -e "${RED}✗ (仅找到 $TRACE_FUNC_COUNT 个函数，期望至少10个)${NC}"
    FAILED_CHECKS=$((FAILED_CHECKS + 1))
fi
echo ""

echo "10. 检查Grafana Dashboard配置"
echo "----------------------------------------"
TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
echo -n "验证Dashboard JSON格式... "
if python3 -m json.tool "deploy/monitoring/dashboards/billing-dashboard.json" > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC}"
    PASSED_CHECKS=$((PASSED_CHECKS + 1))
else
    echo -e "${RED}✗ (JSON格式无效)${NC}"
    FAILED_CHECKS=$((FAILED_CHECKS + 1))
fi
echo ""

echo "=========================================="
echo "验证总结"
echo "=========================================="
echo -e "总检查数: $TOTAL_CHECKS"
echo -e "${GREEN}通过: $PASSED_CHECKS${NC}"
echo -e "${RED}失败: $FAILED_CHECKS${NC}"
echo ""

if [ $FAILED_CHECKS -eq 0 ]; then
    echo -e "${GREEN}所有检查通过！✓${NC}"
    echo ""
    echo "下一步："
    echo "1. 启动Prometheus: prometheus --config.file=/etc/prometheus/prometheus.yml"
    echo "2. 启动Grafana: grafana-server"
    echo "3. 导入Dashboard: 复制 billing-dashboard.json 到 Grafana"
    echo "4. 配置告警: 复制 billing-alerts.yml 到 Prometheus"
    echo "5. 启动Filebeat: filebeat -c filebeat-billing.yml"
    exit 0
else
    echo -e "${RED}部分检查失败，请修复上述问题后重试${NC}"
    exit 1
fi
