#!/bin/bash

# ====================================================================
# 租户隔离机制部署脚本
# ====================================================================
#
# 功能：自动化部署租户隔离机制
# 使用：bash scripts/deploy/tenant_isolation_setup.sh
#
# ====================================================================

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查前置条件
check_prerequisites() {
    log_info "检查前置条件..."

    # 检查MySQL
    if ! command -v mysql &> /dev/null; then
        log_error "MySQL未安装，请先安装MySQL 8.4.5+"
        exit 1
    fi

    # 检查Redis
    if ! command -v redis-cli &> /dev/null; then
        log_error "Redis未安装，请先安装Redis 8.0+"
        exit 1
    fi

    # 检查Go
    if ! command -v go &> /dev/null; then
        log_error "Go未安装，请先安装Go 1.24.0+"
        exit 1
    fi

    log_info "✅ 前置条件检查通过"
}

# 部署MySQL触发器
deploy_mysql_triggers() {
    log_info "部署MySQL触发器..."

    # 读取配置
    if [ -f .env ]; then
        source .env
    else
        log_error ".env文件不存在，请先创建配置文件"
        exit 1
    fi

    # 执行触发器脚本
    mysql -h "${DB_HOST}" -u "${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" \
        < backend/database/triggers/tenant_isolation_triggers.sql

    # 验证触发器
    TRIGGER_COUNT=$(mysql -h "${DB_HOST}" -u "${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" -N \
        -e "SELECT COUNT(*) FROM INFORMATION_SCHEMA.TRIGGERS WHERE TRIGGER_NAME LIKE 'trg_%_tenant%'")

    log_info "✅ 已创建 ${TRIGGER_COUNT} 个租户隔离触发器"
}

# 运行测试
run_tests() {
    log_info "运行测试..."

    # 单元测试
    log_info "运行单元测试..."
    go test ./backend/api/middleware/... -v -run TestTenant
    go test ./backend/infra/cache/... -v -run TestTenant
    go test ./backend/infra/search/... -v -run TestTenant
    go test ./backend/infra/storage/... -v -run TestTenant

    # 安全测试
    log_info "运行安全测试..."
    go test ./backend/tests/security/... -v

    # 性能测试
    log_info "运行性能测试..."
    go test ./backend/api/middleware/... -bench=BenchmarkTenant -benchmem

    log_info "✅ 所有测试通过"
}

# 验证部署
verify_deployment() {
    log_info "验证部署..."

    # 1. 检查触发器
    log_info "检查MySQL触发器..."
    mysql -h "${DB_HOST}" -u "${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" \
        -e "SELECT TRIGGER_NAME, EVENT_OBJECT_TABLE FROM INFORMATION_SCHEMA.TRIGGERS WHERE TRIGGER_NAME LIKE 'trg_%_tenant%' ORDER BY EVENT_OBJECT_TABLE;"

    # 2. 检查Redis连接
    log_info "检查Redis连接..."
    redis-cli -h "${REDIS_HOST}" -a "${REDIS_PASSWORD}" PING

    # 3. 检查代码编译
    log_info "检查代码编译..."
    cd backend
    go build ./...
    cd ..

    log_info "✅ 部署验证通过"
}

# 清理函数（部署失败时回滚）
cleanup() {
    log_warn "部署失败，执行清理..."

    # 删除触发器
    mysql -h "${DB_HOST}" -u "${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" \
        -e "
        DROP PROCEDURE IF EXISTS add_tenant_isolation_triggers;
        DROP PROCEDURE IF EXISTS remove_tenant_isolation_triggers;
        "

    log_info "清理完成"
}

# 主流程
main() {
    log_info "开始部署租户隔离机制..."
    echo ""

    # 检查前置条件
    check_prerequisites
    echo ""

    # 部署MySQL触发器
    deploy_mysql_triggers
    echo ""

    # 运行测试
    run_tests
    echo ""

    # 验证部署
    verify_deployment
    echo ""

    log_info "=========================================="
    log_info "🎉 租户隔离机制部署成功！"
    log_info "=========================================="
    echo ""
    log_info "下一步操作："
    log_info "1. 重启应用服务"
    log_info "2. 监控运行状态"
    log_info "3. 查看文档: docs/企业级功能完善与统一性设计方案/租户隔离机制实现总结_v1.0.md"
}

# 捕获错误并清理
trap cleanup ERR

# 执行主流程
main "$@"
