#!/bin/bash
# 服务自动重启脚本
# ZKER 企业级自治愈系统
# 版本: v1.0.0

set -euo pipefail

# ========================================
# 配置
# ========================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="/var/log/zker/autohealing/service-restart.log"
METRICS_FILE="/var/node_exporter/textfile_collector/zker_autohealing.prom"
ALERTMANAGER_URL="${ALERTMANAGER_URL:-http://localhost:9093}"
SLACK_WEBHOOK="${SLACK_WEBHOOK:-}"

# 服务列表
declare -A SERVICES=(
    ["zker-backend"]="coze-studio"
    ["zker-workflow"]="zker-workflow-engine"
    ["zker-frontend"]="zker-frontend"
    ["mysql-master"]="mysql-master"
    ["redis-cluster-node-1"]="redis-node-1"
)

# 重试配置
MAX_RESTART_ATTEMPTS=3
RESTART_DELAY=10
HEALTH_CHECK_TIMEOUT=60

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

# 检查服务健康状态
check_service_health() {
    local service_name="$1"
    local container_name="${SERVICES[$service_name]}"

    log "INFO" "检查服务健康状态: $service_name"

    # 检查容器是否存在
    if ! docker inspect "$container_name" &>/dev/null; then
        log "ERROR" "容器不存在: $container_name"
        return 2
    fi

    # 检查容器是否运行
    if ! docker inspect -f '{{.State.Running}}' "$container_name" | grep -q 'true'; then
        log "ERROR" "容器未运行: $container_name"
        return 3
    fi

    # 检查健康状态（如果配置了健康检查）
    local health_status
    health_status=$(docker inspect -f '{{.State.Health.Status}}' "$container_name" 2>/dev/null || echo "unknown")

    case "$health_status" in
        healthy)
            log "INFO" "服务健康: $service_name"
            return 0
            ;;
        unhealthy)
            log "WARN" "服务不健康: $service_name"
            return 1
            ;;
        unknown)
            # 容器没有健康检查，检查进程
            if docker exec "$container_name" pgrep -f "$service_name" &>/dev/null; then
                log "INFO" "服务进程运行正常: $service_name"
                return 0
            else
                log "WARN" "服务进程不存在: $service_name"
                return 1
            fi
            ;;
        starting)
            log "INFO" "服务启动中: $service_name"
            return 1
            ;;
    esac
}

# 重启服务
restart_service() {
    local service_name="$1"
    local attempt="${2:-1}"
    local container_name="${SERVICES[$service_name]}"

    log "WARN" "尝试重启服务: $service_name (第 $attempt/$MAX_RESTART_ATTEMPTS 次)"

    # 记录重启指标
    record_metric "zker_autohealing_restart_total" "1" "service=\"$service_name\""

    # 保存容器日志
    local log_backup="/var/log/zker/autohealing/${container_name}_$(date +%Y%m%d_%H%M%S).log"
    docker logs "$container_name" > "$log_backup" 2>&1
    log "INFO" "容器日志已保存: $log_backup"

    # 重启容器
    if docker restart "$container_name"; then
        log "INFO" "容器重启命令已执行: $container_name"

        # 等待容器启动
        sleep "$RESTART_DELAY"

        # 检查服务是否恢复
        if check_service_health "$service_name"; then
            log "INFO" "服务恢复成功: $service_name"
            record_metric "zker_autohealing_recovery_success_total" "1" "service=\"$service_name\""

            # 发送恢复通知
            send_notification "$service_name" "recovered"

            return 0
        else
            log "ERROR" "服务重启后仍然不健康: $service_name"

            # 重试
            if (( attempt < MAX_RESTART_ATTEMPTS )); then
                log "INFO" "等待 $RESTART_DELAY 秒后重试..."
                sleep "$RESTART_DELAY"
                restart_service "$service_name" $((attempt + 1))
                return $?
            else
                log "ERROR" "达到最大重试次数: $service_name"
                record_metric "zker_autohealing_failure_total" "1" "service=\"$service_name\",reason=\"max_retries\""

                # 发送失败通知
                send_notification "$service_name" "failed"

                return 1
            fi
        fi
    else
        log "ERROR" "容器重启失败: $container_name"
        record_metric "zker_autohealing_failure_total" "1" "service=\"$service_name\",reason=\"restart_failed\""

        # 发送失败通知
        send_notification "$service_name" "restart_failed"

        return 1
    fi
}

# 发送通知
send_notification() {
    local service_name="$1"
    local status="$2"

    if [[ -n "$SLACK_WEBHOOK" ]]; then
        local message=""
        local color=""

        case "$status" in
            recovered)
                message="✅ 服务已恢复: $service_name"
                color="good"
                ;;
            failed|restart_failed)
                message="🚨 服务恢复失败: $service_name"
                color="danger"
                ;;
            restarting)
                message="⚠️ 正在重启服务: $service_name"
                color="warning"
                ;;
        esac

        curl -X POST "$SLACK_WEBHOOK" \
            -H 'Content-Type: application/json' \
            -d "{
                \"text\": \"$message\",
                \"attachments\": [{
                    \"color\": \"$color\",
                    \"fields\": [
                        {\"title\": \"服务\", \"value\": \"$service_name\", \"short\": true},
                        {\"title\": \"时间\", \"value\": \"$(date)\", \"short\": true},
                        {\"title\": \"主机\", \"value\": \"$(hostname)\", \"short\": true}
                    ]
                }]
            }" &>/dev/null || true
    fi
}

# 处理告警
process_alert() {
    local alert_json="$1"

    # 解析告警信息
    local alert_name
    local service_name
    local severity

    alert_name=$(echo "$alert_json" | jq -r '.alert.labels.alert // empty')
    severity=$(echo "$alert_json" | jq -r '.alert.labels.severity // empty')

    # 从告警标签中提取服务名
    service_name=$(echo "$alert_json" | jq -r '.alert.labels.job // .alert.labels.service // empty')

    if [[ -z "$service_name" ]]; then
        log "WARN" "无法从告警中提取服务名: $alert_name"
        return 1
    fi

    # 只处理P0级别的服务宕机告警
    if [[ "$severity" != "P0" ]] && [[ "$alert_name" != "ServiceDown" ]]; then
        log "INFO" "跳过非P0告警: $alert_name ($severity)"
        return 0
    fi

    log "WARN" "收到P0告警: $alert_name - 服务: $service_name"

    # 检查服务是否真的宕机
    if check_service_health "$service_name"; then
        log "INFO" "服务健康，无需重启: $service_name"
        return 0
    fi

    # 发送重启通知
    send_notification "$service_name" "restarting"

    # 重启服务
    restart_service "$service_name"
}

# 监听Alertmanager Webhook
listen_alerts() {
    log "INFO" "启动告警监听服务..."

    # 创建临时管道
    local tmp_pipe
    tmp_pipe=$(mktemp -u)
    mkfifo -m 600 "$tmp_pipe"

    # 启动nc监听
    while true; do
        {
            read -r request
            # 读取HTTP body
            read -r content_length
            content_length=$(echo "$content_length" | grep -i 'Content-Length:' | awk '{print $2}')

            if [[ -n "$content_length" ]]; then
                read -r -n "$content_length" body
                echo "$body"

                # 返回HTTP响应
                echo "HTTP/1.1 200 OK"
                echo "Content-Type: text/plain"
                echo ""
                echo "OK"
            fi
        } < "$tmp_pipe" | nc -l -p 9099 | while read -r line; do
            # 处理每个告警
            echo "$line" | jq -c '.alerts[]?' 2>/dev/null | while read -r alert; do
                process_alert "$alert" &
            done
        done
    done

    rm -f "$tmp_pipe"
}

# 主函数
main() {
    # 创建日志目录
    mkdir -p "$(dirname "$LOG_FILE")"
    mkdir -p "/var/log/zker/autohealing"

    log "INFO" "=========================================="
    log "INFO" "ZKER 服务自动重启守护进程"
    log "INFO" "版本: v1.0.0"
    log "INFO" "主机: $(hostname)"
    log "INFO" "=========================================="

    # 检查依赖
    for cmd in docker jq curl nc; do
        if ! command -v "$cmd" &>/dev/null; then
            log "ERROR" "缺少依赖命令: $cmd"
            exit 1
        fi
    done

    # 启动监听
    listen_alerts
}

# 如果直接运行脚本
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
