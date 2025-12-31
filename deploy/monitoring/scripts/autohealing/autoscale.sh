#!/bin/bash
# 自动扩缩容脚本
# ZKER 企业级自治愈系统
# 版本: v1.0.0

set -euo pipefail

# ========================================
# 配置
# ========================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="/var/log/zker/autohealing/autoscale.log"
METRICS_FILE="/var/node_exporter/textfile_collector/zker_autoscale.prom"

# 扩容阈值
CPU_SCALE_UP_THRESHOLD=80
CPU_SCALE_DOWN_THRESHOLD=20
MEMORY_SCALE_UP_THRESHOLD=85
MEMORY_SCALE_DOWN_THRESHOLD=30
QPS_SCALE_UP_THRESHOLD=5000
QPS_SCALE_DOWN_THRESHOLD=1000

# 扩容配置
MIN_REPLICAS=2
MAX_REPLICAS=10
SCALE_UP_STEP=2
SCALE_DOWN_STEP=1
COOLDOWN_PERIOD=300  # 5分钟

# Kubernetes配置
KUBECONFIG="${KUBECONFIG:-/root/.kube/config}"
NAMESPACE="${NAMESPACE:-zker-production}"

# 部署列表
declare -A DEPLOYMENTS=(
    ["zker-backend"]="zker-backend"
    ["zker-workflow"]="zker-workflow-engine"
    ["zker-frontend"]="zker-frontend"
)

# 日志函数
log() {
    local level="$1"
    shift
    local message="[$(date +'%Y-%m-%d %H:%M:%S')] [$level] $*"
    echo "$message" | tee -a "$LOG_FILE"
}

# 记录指标
record_metric() {
    local metric_name="$1"
    local metric_value="$2"
    local labels="${3:-}"

    mkdir -p "$(dirname "$METRICS_FILE")"
    echo "${metric_name}${labels} ${metric_value} $(date +%s)" >> "$METRICS_FILE"
}

# 获取部署当前副本数
get_replicas() {
    local deployment="$1"

    kubectl get deployment "$deployment" \
        -n "$NAMESPACE" \
        -o jsonpath='{.spec.replicas}' 2>/dev/null || echo "0"
}

# 获取部署期望副本数
get_desired_replicas() {
    local deployment="$1"

    kubectl get deployment "$deployment" \
        -n "$NAMESPACE" \
        -o jsonpath='{.status.replicas}' 2>/dev/null || echo "0"
}

# 获取Pod平均CPU使用率
get_avg_cpu_usage() {
    local deployment="$1"

    kubectl top pods -n "$NAMESPACE" \
        -l app="$deployment" \
        --no-headers | \
        awk '{sum+=$2; count++} END {print sum/count}' || echo "0"
}

# 获取Pod平均内存使用率
get_avg_memory_usage() {
    local deployment="$1"

    kubectl top pods -n "$NAMESPACE" \
        -l app="$deployment" \
        --no-headers | \
        awk '{sum+=$3; count++} END {print sum/count}' || echo "0"
}

# 获取部署QPS
get_deployment_qps() {
    local deployment="$1"

    # 从Prometheus查询QPS
    local prometheus_url="http://prometheus:9090"
    local query="sum(rate(zker_api_requests_total{job=\"$deployment\"}[5m]))"

    curl -s -G "$prometheus_url/api/v1/query" \
        --data-urlencode "query=$query" | \
        jq -r '.data.result[0].value[1] // "0"' || echo "0"
}

# 计算目标副本数
calculate_target_replicas() {
    local deployment="$1"
    local current_replicas="$2"
    local cpu_usage="$3"
    local memory_usage="$4"
    local qps="$5"

    local target_replicas="$current_replicas"
    local scale_reason=()

    # CPU扩容判断
    if (( $(echo "$cpu_usage > $CPU_SCALE_UP_THRESHOLD" | bc -l) )); then
        local cpu_replicas
        cpu_replicas=$(awk "BEGIN {print int($current_replicas * ($cpu_usage / $CPU_SCALE_UP_THRESHOLD))}")
        target_replicas=$((cpu_replicas > target_replicas ? cpu_replicas : target_replicas))
        scale_reason+=("CPU高: ${cpu_usage}%")
    fi

    # CPU缩容判断
    if (( $(echo "$cpu_usage < $CPU_SCALE_DOWN_THRESHOLD" | bc -l) )); then
        local cpu_replicas
        cpu_replicas=$(awk "BEGIN {print int($current_replicas * ($cpu_usage / $CPU_SCALE_DOWN_THRESHOLD))}")
        if (( cpu_replicas < target_replicas )); then
            target_replicas=$cpu_replicas
        fi
        scale_reason+=("CPU低: ${cpu_usage}%")
    fi

    # 内存扩容判断
    if (( $(echo "$memory_usage > $MEMORY_SCALE_UP_THRESHOLD" | bc -l) )); then
        local mem_replicas
        mem_replicas=$(awk "BEGIN {print int($current_replicas * ($memory_usage / $MEMORY_SCALE_UP_THRESHOLD))}")
        target_replicas=$((mem_replicas > target_replicas ? mem_replicas : target_replicas))
        scale_reason+=("内存高: ${memory_usage}%")
    fi

    # QPS扩容判断
    if (( $(echo "$qps > $QPS_SCALE_UP_THRESHOLD" | bc -l) )); then
        local qps_replicas
        qps_replicas=$(awk "BEGIN {print int($current_replicas * ($qps / $QPS_SCALE_UP_THRESHOLD))}")
        target_replicas=$((qps_replicas > target_replicas ? qps_replicas : target_replicas))
        scale_reason+=("QPS高: ${qps}")
    fi

    # 限制副本数范围
    if (( target_replicas < MIN_REPLICAS )); then
        target_replicas=$MIN_REPLICAS
    fi

    if (( target_replicas > MAX_REPLICAS )); then
        target_replicas=$MAX_REPLICAS
    fi

    echo "$target_replicas|${scale_reason[*]}"
}

# 执行扩缩容
scale_deployment() {
    local deployment="$1"
    local current_replicas="$2"
    local target_replicas="$3"
    shift 3
    local reasons=("$@")

    if (( target_replicas == current_replicas )); then
        log "INFO" "副本数无需调整: $deployment ($current_replicas)"
        return 0
    fi

    log "INFO" "开始扩缩容: $deployment"
    log "INFO" "  当前副本数: $current_replicas"
    log "INFO" "  目标副本数: $target_replicas"
    log "INFO" "  扩缩容原因: ${reasons[*]}"

    # 记录指标
    if (( target_replicas > current_replicas )); then
        record_metric "zker_autoscale_scale_up_total" "1" "deployment=\"$deployment\",from=$current_replicas,to=$target_replicas"
    else
        record_metric "zker_autoscale_scale_down_total" "1" "deployment=\"$deployment\",from=$current_replicas,to=$target_replicas"
    fi

    # 执行扩缩容
    if kubectl scale deployment "$deployment" \
        -n "$NAMESPACE" \
        --replicas="$target_replicas"; then
        log "INFO" "扩缩容命令已执行: $deployment -> $target_replicas replicas"

        # 等待扩缩容完成
        local timeout=300
        local elapsed=0
        while (( elapsed < timeout )); do
            local ready_replicas
            ready_replicas=$(kubectl get deployment "$deployment" \
                -n "$NAMESPACE" \
                -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")

            if (( ready_replicas == target_replicas )); then
                log "INFO" "扩缩容完成: $deployment ($ready_replicas/$target_replicas ready)"
                record_metric "zker_autoscale_success_total" "1" "deployment=\"$deployment\""
                return 0
            fi

            sleep 5
            ((elapsed += 5))
        done

        log "ERROR" "扩缩容超时: $deployment (等待${timeout}秒后仍为${ready_replicas}/${target_replicas} ready)"
        record_metric "zker_autoscale_timeout_total" "1" "deployment=\"$deployment\""
        return 1
    else
        log "ERROR" "扩缩容失败: $deployment"
        record_metric "zker_autoscale_failure_total" "1" "deployment=\"$deployment\""
        return 1
    fi
}

# 检查冷却期
check_cooldown() {
    local deployment="$1"

    local last_scale_file="/var/run/zker/autoscale_${deployment}_last_scale"

    if [[ -f "$last_scale_file" ]]; then
        local last_scale_time
        last_scale_time=$(cat "$last_scale_file")
        local current_time
        current_time=$(date +%s)
        local elapsed=$((current_time - last_scale_time))

        if (( elapsed < COOLDOWN_PERIOD )); then
            log "INFO" "冷却期未过: $deployment (已过${elapsed}秒，需${COOLDOWN_PERIOD}秒)"
            return 1
        fi
    fi

    return 0
}

# 更新扩缩容时间
update_scale_time() {
    local deployment="$1"

    mkdir -p "$(dirname "$last_scale_file")"
    date +%s > "/var/run/zker/autoscale_${deployment}_last_scale"
}

# 处理单个部署
process_deployment() {
    local deployment="$1"

    log "INFO" "检查部署扩缩容: $deployment"

    # 检查冷却期
    if ! check_cooldown "$deployment"; then
        return 0
    fi

    # 获取当前副本数
    local current_replicas
    current_replicas=$(get_replicas "$deployment")

    if (( current_replicas == 0 )); then
        log "WARN" "部署不存在或副本数为0: $deployment"
        return 1
    fi

    # 获取指标
    local cpu_usage memory_usage qps
    cpu_usage=$(get_avg_cpu_usage "$deployment")
    memory_usage=$(get_avg_memory_usage "$deployment")
    qps=$(get_deployment_qps "$deployment")

    log "INFO" "  CPU使用率: ${cpu_usage}%"
    log "INFO" "  内存使用率: ${memory_usage}%"
    log "INFO" "  QPS: ${qps}"

    # 计算目标副本数
    local result
    result=$(calculate_target_replicas "$deployment" "$current_replicas" "$cpu_usage" "$memory_usage" "$qps")
    local target_replicas="${result%%|*}"
    local reasons="${result#*|}"

    # 执行扩缩容
    if scale_deployment "$deployment" "$current_replicas" "$target_replicas" "$reasons"; then
        update_scale_time "$deployment"
    fi
}

# 主函数
main() {
    # 创建日志目录
    mkdir -p "$(dirname "$LOG_FILE")"

    log "INFO" "=========================================="
    log "INFO" "ZKER 自动扩缩容守护进程"
    log "INFO" "版本: v1.0.0"
    log "INFO" "命名空间: $NAMESPACE"
    log "INFO" "=========================================="

    # 检查kubectl
    if ! command -v kubectl &>/dev/null; then
        log "ERROR" "kubectl未安装"
        exit 1
    fi

    # 检查集群连接
    if ! kubectl cluster-info &>/dev/null; then
        log "ERROR" "无法连接Kubernetes集群"
        exit 1
    fi

    # 处理所有部署
    for deployment in "${!DEPLOYMENTS[@]}"; do
        process_deployment "$deployment"
    done
}

# 如果直接运行脚本
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
