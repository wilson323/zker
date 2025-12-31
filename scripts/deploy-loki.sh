#!/bin/bash
# =============================================================================
# Loki日志系统快速部署脚本
# 版本: v1.0.0
# 更新: 2025-01-03
# =============================================================================

set -euo pipefail

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Docker
check_docker() {
    log "检查Docker环境..."

    if ! command -v docker &> /dev/null; then
        error "Docker未安装，请先安装Docker"
        exit 1
    fi

    if ! command -v docker compose &> /dev/null && ! docker compose version &> /dev/null; then
        error "Docker Compose未安装"
        exit 1
    fi

    log "✅ Docker环境正常"
}

# 启动Loki
start_loki() {
    log "启动Loki日志系统..."

    cd "$(dirname "$0")/../docker"

    # 启动服务
    docker compose -f docker-compose-monitoring.yml up -d loki promtail

    log "等待服务启动..."
    sleep 10

    # 检查服务状态
    if docker compose -f docker-compose-monitoring.yml ps loki promtail | grep -q "Up"; then
        log "✅ Loki和Promtail启动成功"
    else
        error "服务启动失败"
        docker compose -f docker-compose-monitoring.yml logs loki promtail
        exit 1
    fi
}

# 验证部署
verify_deployment() {
    log "验证Loki部署..."

    # 检查Loki健康状态
    if curl -s http://localhost:3100/ready | grep -q "ready"; then
        log "✅ Loki健康检查通过"
    else
        error "Loki健康检查失败"
        return 1
    fi

    # 检查Promtail
    if curl -s http://localhost:9080/metrics | grep -q "promtail"; then
        log "✅ Promtail健康检查通过"
    else
        warn "Promtail健康检查失败"
    fi

    # 查询标签
    local labels=$(curl -s http://localhost:3100/loki/api/v1/labels | jq -r '.data[:5]' 2>/dev/null || echo "[]")
    log "可用标签: ${labels}"

    return 0
}

# 显示访问信息
show_info() {
    echo ""
    echo "=========================================="
    echo "🎉 Loki日志系统部署成功！"
    echo "=========================================="
    echo ""
    echo "服务访问地址:"
    echo "  - Loki API:   http://localhost:3100"
    echo "  - Grafana:    http://localhost:3000"
    echo "  - Promtail:   http://localhost:9080"
    echo ""
    echo "查看日志:"
    echo "  docker logs -f zker-loki"
    echo "  docker logs -f zker-promtail"
    echo ""
    echo "LogQL查询示例:"
    echo "  curl -G http://localhost:3100/loki/api/v1/query --data-urlencode 'query={job=\"backend\"}'"
    echo ""
    echo "配置文件:"
    echo "  - Loki:     docker/volumes/monitoring/loki/loki-config.yml"
    echo "  - Promtail: docker/volumes/monitoring/promtail/promtail-config.yml"
    echo ""
    echo "详细文档: scripts/backup/README.md"
    echo "=========================================="
}

# 主流程
main() {
    echo "=========================================="
    echo "ZKER Loki日志系统部署工具 v1.0.0"
    echo "=========================================="
    echo ""

    check_docker
    start_loki
    verify_deployment
    show_info
}

main "$@"
