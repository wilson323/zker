#!/bin/bash
# ZKER 模拟业务指标生成器
# 版本: v1.0.0
# 用途: 生成模拟的业务指标数据，用于测试监控系统

set -e

# 配置
API_URL="${API_URL:-http://localhost:8888}"

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

# 模拟API请求（生成HTTP指标）
simulate_api_requests() {
    log_info "模拟API请求..."

    local endpoints=(
        "/api/v1/tenant"
        "/api/v1/bot"
        "/api/v1/conversation"
        "/api/v1/workflow"
        "/api/v1/quota/check"
    )

    local duration=60
    local qps=100

    log_info "启动 ${qps} QPS 模拟，持续 ${duration} 秒..."

    for i in $(seq 1 $qps); do
        (
            while true; do
                endpoint=${endpoints[$RANDOM % ${#endpoints[@]}]}
                curl -s -w "\n" "$API_URL$endpoint" > /dev/null 2>&1
                sleep 0.$((RANDOM % 10))
            done
        ) &
    done

    sleep $duration

    log_success "API请求模拟完成"
}

# 模拟Bot调用（生成Bot指标）
simulate_bot_invocations() {
    log_info "模拟Bot调用..."

    local bot_ids=("bot-001" "bot-002" "bot-003")
    local duration=60
    local qps=50

    log_info "启动 ${qps} QPS Bot调用模拟，持续 ${duration} 秒..."

    for i in $(seq 1 $qps); do
        (
            while true; do
                bot_id=${bot_ids[$RANDOM % ${#bot_ids[@]}]}
                # 模拟成功和失败
                if [ $((RANDOM % 10)) -lt 9 ]; then
                    status="success"
                else
                    status="failed"
                fi
                curl -s "$API_URL/api/v1/bot/$bot_id/invoke" -X POST -d "{\"status\":\"$status\"}" > /dev/null 2>&1
                sleep 0.$((RANDOM % 10))
            done
        ) &
    done

    sleep $duration

    log_success "Bot调用模拟完成"
}

# 模拟数据库查询（生成数据库指标）
simulate_db_queries() {
    log_info "模拟数据库查询..."

    local duration=60
    local qps=200

    log_info "启动 ${qps} QPS 数据库查询模拟，持续 ${duration} 秒..."

    for i in $(seq 1 $qps); do
        (
            while true; do
                # 通过API间接触发数据库查询
                curl -s "$API_URL/api/v1/tenant/list" > /dev/null 2>&1
                sleep 0.$((RANDOM % 10))
            done
        ) &
    done

    sleep $duration

    log_success "数据库查询模拟完成"
}

# 模拟缓存操作（生成缓存指标）
simulate_cache_operations() {
    log_info "模拟缓存操作..."

    local duration=60
    local qps=300

    log_info "启动 ${qps} QPS 缓存操作模拟，持续 ${duration} 秒..."

    for i in $(seq 1 $qps); do
        (
            while true; do
                # 模拟缓存命中和未命中
                if [ $((RANDOM % 10)) -lt 7 ]; then
                    # 命中
                    curl -s "$API_URL/api/v1/cache/get/hit" > /dev/null 2>&1
                else
                    # 未命中
                    curl -s "$API_URL/api/v1/cache/get/miss" > /dev/null 2>&1
                fi
                sleep 0.$((RANDOM % 10))
            done
        ) &
    done

    sleep $duration

    log_success "缓存操作模拟完成"
}

# 模拟配额检查（生成配额指标）
simulate_quota_checks() {
    log_info "模拟配额检查..."

    local tenant_ids=("tenant-001" "tenant-002" "tenant-003")
    local resource_types=("bots" "conversations" "storage")
    local duration=60
    local qps=20

    log_info "启动 ${qps} QPS 配额检查模拟，持续 ${duration} 秒..."

    for i in $(seq 1 $qps); do
        (
            while true; do
                tenant_id=${tenant_ids[$RANDOM % ${#tenant_ids[@]}]}
                resource_type=${resource_types[$RANDOM % ${#resource_types[@]}]}
                curl -s "$API_URL/api/v1/quota/check" \
                    -X POST \
                    -d "{\"tenant_id\":\"$tenant_id\",\"resource_type\":\"$resource_type\"}" \
                    > /dev/null 2>&1
                sleep 0.$((RANDOM % 10))
            done
        ) &
    done

    sleep $duration

    log_success "配额检查模拟完成"
}

# 主函数
main() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  ZKER 业务指标生成器${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""

    # 检查参数
    MODE=${1:-all}

    case $MODE in
        api)
            simulate_api_requests
            ;;
        bot)
            simulate_bot_invocations
            ;;
        db)
            simulate_db_queries
            ;;
        cache)
            simulate_cache_operations
            ;;
        quota)
            simulate_quota_checks
            ;;
        all)
            log_info "启动所有模拟器..."
            simulate_api_requests &
            simulate_bot_invocations &
            simulate_db_queries &
            simulate_cache_operations &
            simulate_quota_checks &
            wait
            ;;
        *)
            echo "用法: $0 [api|bot|db|cache|quota|all]"
            exit 1
            ;;
    esac

    echo ""
    log_success "所有模拟完成！"
    log_info "请检查 Prometheus 和 Grafana 查看指标数据"
}

# 执行主函数
main "$@"
