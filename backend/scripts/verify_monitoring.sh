#!/bin/bash
# 🔧 P0修复：监控系统验证脚本
# 用于验证监控日志系统是否正确初始化和运行

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
API_BASE_URL="${API_BASE_URL:-http://localhost:8888}"
METRICS_ENDPOINT="${API_BASE_URL}/metrics"
HEALTH_ENDPOINT="${API_BASE_URL}/health"
JAEGER_UI_URL="${JAEGER_UI_URL:-http://localhost:16686}"
PROMETHEUS_URL="${PROMETHEUS_URL:-http://localhost:9090}"
GRAFANA_URL="${GRAFANA_URL:-http://localhost:3001}"

echo "🔍 开始验证监控系统..."
echo ""

# 1. 检查Prometheus /metrics端点
echo "📊 检查 1: Prometheus /metrics 端点"
METRICS_RESPONSE=$(curl -s "${METRICS_ENDPOINT}" || echo "")
if [[ $METRICS_RESPONSE == *"go_memstats"* ]] || [[ $METRICS_RESPONSE == *"http_requests_total"* ]]; then
    echo -e "${GREEN}✅ Prometheus /metrics 端点正常${NC}"
    echo "   找到指标示例:"
    echo "$METRICS_RESPONSE" | head -n 5
    echo ""
else
    echo -e "${RED}❌ Prometheus /metrics 端点异常或无响应${NC}"
    echo "   响应: $METRICS_RESPONSE"
    exit 1
fi

# 2. 检查/health端点
echo "💓 检查 2: 健康检查 /health 端点"
HEALTH_RESPONSE=$(curl -s "${HEALTH_ENDPOINT}" || echo "")
if [[ $HEALTH_RESPONSE == *"healthy"* ]] || [[ $HEALTH_RESPONSE == *"status"* ]]; then
    echo -e "${GREEN}✅ 健康检查端点正常${NC}"
    echo "   响应: $HEALTH_RESPONSE"
    echo ""
else
    echo -e "${YELLOW}⚠️  健康检查端点可能未配置（非必需）${NC}"
    echo "   响应: $HEALTH_RESPONSE"
    echo ""
fi

# 3. 检查关键指标是否存在
echo "📈 检查 3: 关键业务指标"
CRITICAL_METRICS=(
    "http_requests_total"
    "http_request_duration_seconds"
    "tenant_total"
    "db_query_duration_seconds"
)

ALL_METRICS_FOUND=true
for metric in "${CRITICAL_METRICS[@]}"; do
    if curl -s "${METRICS_ENDPOINT}" | grep -q "$metric"; then
        echo -e "${GREEN}✅ 找到指标: $metric${NC}"
    else
        echo -e "${YELLOW}⚠️  未找到指标: $metric${NC}"
        ALL_METRICS_FOUND=false
    fi
done
echo ""

# 4. 检查Jaeger连接
echo "🔍 检查 4: Jaeger分布式追踪"
JAEGER_STATUS=$(curl -s "${JAEGER_UI_URL}/api/status" || echo "")
if [[ $JAEGER_STATUS == *"jaeger"* ]]; then
    echo -e "${GREEN}✅ Jaeger UI可访问: ${JAEGER_UI_URL}${NC}"
else
    echo -e "${YELLOW}⚠️  Jaeger UI可能未启动（非必需）${NC}"
    echo "   URL: ${JAEGER_UI_URL}"
fi
echo ""

# 5. 检查Prometheus服务
echo "🔍 检查 5: Prometheus服务"
PROMETHEUS_HEALTH=$(curl -s "${PROMETHEUS_URL}/-/healthy" || echo "")
if [[ $PROMETHEUS_HEALTH == *"Prometheus is Healthy"* ]]; then
    echo -e "${GREEN}✅ Prometheus服务正常${NC}"
    echo "   URL: ${PROMETHEUS_URL}"
else
    echo -e "${YELLOW}⚠️  Prometheus服务可能未启动（非必需）${NC}"
    echo "   URL: ${PROMETHEUS_URL}"
fi
echo ""

# 6. 检查Grafana服务
echo "🔍 检查 6: Grafana可视化服务"
GRAFANA_HEALTH=$(curl -s "${GRAFANA_URL}/api/health" || echo "")
if [[ $GRAFANA_HEALTH == *"commit"* ]] || [[ $GRAFANA_HEALTH == *"database"* ]]; then
    echo -e "${GREEN}✅ Grafana服务正常${NC}"
    echo "   URL: ${GRAFANA_URL}"
else
    echo -e "${YELLOW}⚠️  Grafana服务可能未启动（非必需）${NC}"
    echo "   URL: ${GRAFANA_URL}"
fi
echo ""

# 7. 测试指标生成（发送HTTP请求）
echo "🧪 检查 7: 测试指标生成"
echo "   发送测试请求到 ${API_BASE_URL}..."
curl -s "${API_BASE_URL}" > /dev/null 2>&1 || true
sleep 2 # 等待指标更新

NEW_REQUESTS=$(curl -s "${METRICS_ENDPOINT}" | grep "^http_requests_total" | wc -l)
if [[ $NEW_REQUESTS -gt 0 ]]; then
    echo -e "${GREEN}✅ 指标正在生成（找到 ${NEW_REQUESTS} 个http_requests_total指标）${NC}"
else
    echo -e "${YELLOW}⚠️  未检测到新的HTTP请求指标${NC}"
fi
echo ""

# 总结
echo "=========================================="
echo -e "${GREEN}✨ 监控系统验证完成！${NC}"
echo ""
echo "📋 访问地址："
echo "   - API服务:     ${API_BASE_URL}"
echo "   - Prometheus:  ${PROMETHEUS_URL}"
echo "   - Grafana:     ${GRAFANA_URL} (admin/admin)"
echo "   - Jaeger UI:   ${JAEGER_UI_URL}"
echo ""
echo "💡 提示："
echo "   1. 检查/metrics端点是否暴露所有预期指标"
echo "   2. 在Prometheus中查询: up{job=\"coze-api\"}"
echo "   3. 在Grafana中创建仪表盘可视化"
echo "   4. 在Jaeger中查看分布式追踪"
echo "=========================================="

# 返回状态
if $ALL_METRICS_FOUND; then
    exit 0
else
    echo -e "${YELLOW}⚠️  部分关键指标未找到，请检查配置${NC}"
    exit 1
fi
