#!/bin/bash
# 测试环境管理脚本
#
# 用途：管理独立测试环境的启动、停止、清理
#
# @author 研发D
# @date 2025-01-02

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Docker Compose配置文件
COMPOSE_FILE="$PROJECT_ROOT/docker/docker-compose.test.yml"
ENV_FILE="$PROJECT_ROOT/docker/.env.test"

# ==================== 函数定义 ====================

/**
 * 显示帮助信息
 */
show_help() {
    echo "测试环境管理脚本"
    echo ""
    echo "用法: $0 [command]"
    echo ""
    echo "命令:"
    echo "  start     启动测试环境"
    echo "  stop      停止测试环境"
    echo "  restart   重启测试环境"
    echo "  status    查看测试环境状态"
    echo "  logs      查看日志"
    echo "  clean     清理测试环境（删除所有数据）"
    echo "  help      显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 start           # 启动测试环境"
    echo "  $0 logs mysql      # 查看MySQL日志"
}

/**
 * 检查环境文件
 */
check_env_file() {
    if [ ! -f "$ENV_FILE" ]; then
        echo -e "${YELLOW}⚠️  环境文件不存在，从示例创建${NC}"
        cp "$PROJECT_ROOT/docker/.env.test.example" "$ENV_FILE"
        echo -e "${GREEN}✅ 已创建 $ENV_FILE${NC}"
    fi
}

/**
 * 启动测试环境
 */
start_test_env() {
    echo -e "${BLUE}🚀 启动测试环境...${NC}"
    check_env_file

    cd "$PROJECT_ROOT/docker"
    docker-compose -f docker-compose.test.yml up -d

    echo -e "${GREEN}✅ 测试环境启动完成${NC}"
    echo ""
    echo "服务地址:"
    echo "  MySQL:   localhost:3307"
    echo "  Redis:   localhost:6380"
    echo ""
    echo "数据库连接:"
    echo "  Host:     localhost"
    echo "  Port:     3307"
    echo "  User:     coze_test"
    echo "  Password: coze_test123"
    echo "  Database: coze_test"
}

/**
 * 停止测试环境
 */
stop_test_env() {
    echo -e "${YELLOW}⏸️  停止测试环境...${NC}"
    cd "$PROJECT_ROOT/docker"
    docker-compose -f docker-compose.test.yml down
    echo -e "${GREEN}✅ 测试环境已停止${NC}"
}

/**
 * 重启测试环境
 */
restart_test_env() {
    echo -e "${BLUE}🔄 重启测试环境...${NC}"
    stop_test_env
    sleep 2
    start_test_env
}

/**
 * 查看状态
 */
show_status() {
    echo -e "${BLUE}📊 测试环境状态${NC}"
    cd "$PROJECT_ROOT/docker"
    docker-compose -f docker-compose.test.yml ps
}

/**
 * 查看日志
 */
show_logs() {
    local service=$1
    echo -e "${BLUE}📄 查看日志${NC}"
    cd "$PROJECT_ROOT/docker"
    if [ -z "$service" ]; then
        docker-compose -f docker-compose.test.yml logs -f
    else
        docker-compose -f docker-compose.test.yml logs -f "$service"
    fi
}

/**
 * 清理测试环境
 */
clean_test_env() {
    echo -e "${RED}⚠️  警告：此操作将删除所有测试数据！${NC}"
    read -p "确定要清理测试环境吗？(yes/no): " confirm

    if [ "$confirm" != "yes" ]; then
        echo -e "${YELLOW}❌ 已取消${NC}"
        return 0
    fi

    echo -e "${RED}🗑️  清理测试环境...${NC}"
    cd "$PROJECT_ROOT/docker"

    # 停止并删除容器
    docker-compose -f docker-compose.test.yml down -v

    # 删除数据卷
    rm -rf "$PROJECT_ROOT/docker/data/mysql-test"
    rm -rf "$PROJECT_ROOT/docker/data/redis-test"

    echo -e "${GREEN}✅ 测试环境已清理${NC}"
}

// ==================== 主程序 ====================

/**
 * 主函数
 */
main() {
    case "${1:-help}" in
        start)
            start_test_env
            ;;
        stop)
            stop_test_env
            ;;
        restart)
            restart_test_env
            ;;
        status)
            show_status
            ;;
        logs)
            show_logs "$2"
            ;;
        clean)
            clean_test_env
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            echo -e "${RED}❌ 未知命令: $1${NC}"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

// 执行主函数
main "$@"
