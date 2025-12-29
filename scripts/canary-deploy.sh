#!/bin/bash
# scripts/canary-deploy.sh
# 灰度发布脚本 - 基于 Istio 的金丝雀部署
# 用途: 逐步将流量从旧版本切换到新版本,降低发布风险
# 遵循规范: 灰度发布最佳实践 + 自动化验证

set -e  # 遇到错误立即退出

# ========== 颜色输出定义 ==========
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# ========== 配置项 ==========
NAMESPACE="${NAMESPACE:-zker-prod}"
SERVICE="${SERVICE:-zker-api}"
NEW_VERSION="${1:-}"
CANARY_PERCENT="${2:-10}"  # 默认10%流量
MAX_CANARY_PERCENT=100
INCREMENT_STEP=10         # 每次增加10%
WAIT_TIME=300             # 每次增加流量后等待5分钟

# Prometheus 配置
PROMETHEUS_URL="${PROMETHEUS_URL:-http://prometheus:9090}"
ERROR_RATE_THRESHOLD=0.05  # 错误率阈值5%
LATENCY_THRESHOLD_P95=1000  # P95延迟阈值1000ms

# ========== 日志函数 ==========
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

# ========== 检查函数 ==========
check_args() {
    if [ -z "$NEW_VERSION" ]; then
        log_error "请提供新版本号"
        echo "用法: $0 <new_version> [canary_percent]"
        echo ""
        echo "示例:"
        echo "  $0 v1.0.1        # 使用默认10%流量"
        echo "  $0 v1.0.1 20     # 使用20%流量"
        echo "  $0 v1.0.1 auto   # 自动逐步增加流量"
        exit 1
    fi

    if [ "$CANARY_PERCENT" = "auto" ]; then
        AUTO_MODE=true
        CANARY_PERCENT=10
    else
        AUTO_MODE=false
        # 验证百分比
        if ! [[ "$CANARY_PERCENT" =~ ^[0-9]+$ ]] || [ "$CANARY_PERCENT" -lt 1 ] || [ "$CANARY_PERCENT" -gt "$MAX_CANARY_PERCENT" ]; then
            log_error "金丝雀流量百分比必须在 1-$MAX_CANARY_PERCENT 之间"
            exit 1
        fi
    fi

    log_info "灰度发布配置:"
    echo "  命名空间: $NAMESPACE"
    echo "  服务: $SERVICE"
    echo "  新版本: $NEW_VERSION"
    echo "  初始流量: $CANARY_PERCENT%"
    echo "  自动模式: $AUTO_MODE"
}

check_istio_installed() {
    log_step "检查 Istio 是否安装..."

    if ! kubectl get namespace istio-system &> /dev/null; then
        log_error "Istio 未安装,请先安装 Istio"
        exit 1
    fi

    log_info "Istio 已安装"
}

check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl 未安装"
        exit 1
    fi
}

# ========== Istio 配置函数 ==========
create_virtualservice() {
    log_step "创建/更新 VirtualService..."

    cat <<EOF | kubectl apply -f -
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: $SERVICE
  namespace: $NAMESPACE
spec:
  hosts:
  - $SERVICE
  http:
  - match:
    - headers:
        x-canary:
          exact: "true"
    route:
    - destination:
        host: $SERVICE
        subset: v2
      weight: 100
  - route:
    - destination:
        host: $SERVICE
        subset: v1
      weight: $((100 - CANARY_PERCENT))
    - destination:
        host: $SERVICE
        subset: v2
      weight: $CANARY_PERCENT
EOF

    log_info "VirtualService 已创建/更新"
}

create_destinationrule() {
    log_step "创建/更新 DestinationRule..."

    cat <<EOF | kubectl apply -f -
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: $SERVICE
  namespace: $NAMESPACE
spec:
  host: $SERVICE
  subsets:
  - name: v1
    labels:
      version: current
  - name: v2
    labels:
      version: $NEW_VERSION
EOF

    log_info "DestinationRule 已创建/更新"
}

# ========== 部署金丝雀版本 ==========
deploy_canary() {
    log_step "部署金丝雀版本 $NEW_VERSION..."

    # 创建金丝雀 Deployment
    cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: $SERVICE-canary
  namespace: $NAMESPACE
  labels:
    app: $SERVICE
    version: $NEW_VERSION
spec:
  replicas: 1
  selector:
    matchLabels:
      app: $SERVICE
      version: $NEW_VERSION
  template:
    metadata:
      labels:
        app: $SERVICE
        version: $NEW_VERSION
    spec:
      containers:
      - name: $SERVICE
        image: zker/api:$NEW_VERSION
        # ... 其他配置从主 Deployment 复制
EOF

    log_info "等待金丝雀 Pod 就绪..."
    kubectl wait --for=condition=available --timeout=5m \
      deployment/$SERVICE-canary -n $NAMESPACE

    log_info "✅ 金丝雀版本已部署"
}

# ========== 更新流量比例 ==========
update_traffic_weight() {
    local weight=$1
    local version_label=$2

    log_step "更新流量比例: ${weight}% -> 金丝雀版本"

    cat <<EOF | kubectl apply -f -
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: $SERVICE
  namespace: $NAMESPACE
spec:
  hosts:
  - $SERVICE
  http:
  - match:
    - headers:
        x-canary:
          exact: "true"
    route:
    - destination:
        host: $SERVICE
        subset: v2
      weight: 100
  - route:
    - destination:
        host: $SERVICE
        subset: v1
      weight: $((100 - weight))
    - destination:
        host: $SERVICE
        subset: v2
      weight: $weight
EOF

    log_info "流量已更新: 旧版本 $((100 - weight))%, 金丝雀版本 ${weight}%"
}

# ========== 监控和验证函数 ==========
wait_for_pods_ready() {
    log_step "等待所有 Pod 就绪..."

    kubectl rollout status deployment/$SERVICE-canary -n $NAMESPACE --timeout=5m

    log_info "✅ 所有 Pod 已就绪"
}

check_metrics() {
    log_step "检查指标..."

    # 查询错误率
    local error_rate=$(curl -s "$PROMETHEUS_URL/api/v1/query?query=rate(http_requests_total{namespace=\"$NAMESPACE\",service=\"$SERVICE\",status=~\"5..\"}[5m])" | \
      jq -r '.data.result[0].value[1] // "0"')

    # 查询P95延迟
    local latency_p95=$(curl -s "$PROMETHEUS_URL/api/v1/query?query=histogram_quantile(0.95, rate(http_request_duration_seconds_bucket{namespace=\"$NAMESPACE\",service=\"$SERVICE\"}[5m]))" | \
      jq -r '.data.result[0].value[1] // "0"')

    log_info "当前指标:"
    echo "  错误率: $(echo "$error_rate * 100" | bc)%"
    echo "  P95延迟: $(echo "$latency_p95 * 1000" | bc)ms"

    # 检查阈值
    if (( $(echo "$error_rate > $ERROR_RATE_THRESHOLD" | bc -l) )); then
        log_error "❌ 错误率超过阈值: $(echo "$error_rate * 100" | bc)% > $(echo "$ERROR_RATE_THRESHOLD * 100" | bc)%"
        return 1
    fi

    if (( $(echo "$latency_p95 * 1000 > $LATENCY_THRESHOLD_P95" | bc -l) )); then
        log_error "❌ P95延迟超过阈值: $(echo "$latency_p95 * 1000" | bc)ms > ${LATENCY_THRESHOLD_P95}ms"
        return 1
    fi

    log_info "✅ 指标检查通过"
    return 0
}

# ========== 逐步增加流量 ==========
increment_traffic() {
    log_step "自动逐步增加流量模式..."

    local current_percent=$CANARY_PERCENT

    while [ $current_percent -lt $MAX_CANARY_PERCENT ]; do
        log_info "当前流量: ${current_percent}%"

        # 更新流量
        update_traffic_weight $current_percent $NEW_VERSION

        # 等待观察
        log_info "等待 ${WAIT_TIME}s 观察期..."
        sleep $WAIT_TIME

        # 检查指标
        if ! check_metrics; then
            log_error "指标检查失败,停止增加流量!"
            return 1
        fi

        # 增加流量
        current_percent=$((current_percent + INCREMENT_STEP))
        if [ $current_percent -gt $MAX_CANARY_PERCENT ]; then
            current_percent=$MAX_CANARY_PERCENT
        fi
    done

    log_info "✅ 已达到100%流量切换"
    return 0
}

# ========== 清理金丝雀资源 ==========
cleanup_canary() {
    log_step "清理金丝雀资源..."

    # 删除金丝雀 Deployment
    kubectl delete deployment $SERVICE-canary -n $NAMESPACE --ignore-not-found=true

    # 更新主 Deployment 标签
    kubectl patch deployment $SERVICE -n $NAMESPACE -p '{"spec":{"template":{"metadata":{"labels":{"version":"'$NEW_VERSION'"}}}}}'

    log_info "✅ 金丝雀资源已清理"
}

# ========== 回滚函数 ==========
rollback_canary() {
    log_error "执行回滚..."

    # 将流量切回旧版本
    cat <<EOF | kubectl apply -f -
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: $SERVICE
  namespace: $NAMESPACE
spec:
  hosts:
  - $SERVICE
  http:
  - route:
    - destination:
        host: $SERVICE
        subset: v1
      weight: 100
EOF

    # 删除金丝雀 Deployment
    kubectl delete deployment $SERVICE-canary -n $NAMESPACE --ignore-not-found=true

    log_info "✅ 已回滚到旧版本"
}

# ========== 主函数 ==========
main() {
    log_info "========================================="
    log_info "     灰度发布 - ZKER $SERVICE"
    log_info "========================================="
    echo ""

    # 检查
    check_kubectl
    check_args
    check_istio_installed

    # 创建 Istio 配置
    create_destinationrule
    create_virtualservice

    # 部署金丝雀
    deploy_canary
    wait_for_pods_ready

    # 更新流量
    update_traffic_weight $CANARY_PERCENT $NEW_VERSION

    # 自动模式
    if [ "$AUTO_MODE" = true ]; then
        if ! increment_traffic; then
            rollback_canary
            exit 1
        fi
    else
        log_info "手动模式: 请验证金丝雀版本"
        log_info "查看金丝雀日志: kubectl logs -f -n $NAMESPACE -l app=$SERVICE,version=$NEW_VERSION"
        log_info "测试金丝雀: curl -H 'x-canary: true' http://$SERVICE/"
    fi

    # 清理
    if [ "$AUTO_MODE" = true ]; then
        cleanup_canary
    fi

    log_info "========================================="
    log_info "     灰度发布完成!"
    log_info "========================================="
}

# 执行主函数
main "$@"
