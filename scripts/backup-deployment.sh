#!/bin/bash
# scripts/backup-deployment.sh - 部署备份脚本
# 版本: v1.0.0
# 说明: 在部署前自动备份Kubernetes资源配置

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 配置
NAMESPACE="${1:-coze-studio-prod}"
DEPLOYMENT="${2:-backend}"
BACKUP_DIR="${BACKUP_DIR:-./deploy-backups}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/${DEPLOYMENT}_${TIMESTAMP}.yaml"

# 创建备份目录
mkdir -p "$BACKUP_DIR"

log_info "=========================================="
log_info "  部署备份脚本"
log_info "=========================================="
log_info "命名空间: $NAMESPACE"
log_info "部署: $DEPLOYMENT"
log_info "备份目录: $BACKUP_DIR"
log_info ""

# 检查kubectl
if ! command -v kubectl &> /dev/null; then
    log_error "kubectl未安装或不在PATH中"
    exit 1
fi

# 检查Deployment是否存在
if ! kubectl get deployment "$DEPLOYMENT" -n "$NAMESPACE" &> /dev/null; then
    log_error "Deployment $DEPLOYMENT 在命名空间 $NAMESPACE 中不存在"
    exit 1
fi

# 备份Deployment配置
log_info "备份Deployment配置..."
kubectl get deployment "$DEPLOYMENT" -n "$NAMESPACE" -o yaml > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    log_success "Deployment配置已备份到: $BACKUP_FILE"

    # 显示备份文件大小
    SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    log_info "文件大小: $SIZE"

    # 创建最新备份的软链接
    LATEST_LINK="${BACKUP_DIR}/${DEPLOYMENT}_latest.yaml"
    ln -sf "$(basename "$BACKUP_FILE")" "$LATEST_LINK"
    log_info "最新备份链接: $LATEST_LINK"

    # 清理旧备份（保留最近10个）
    log_info "清理旧备份（保留最近10个）..."
    ls -t "${BACKUP_DIR}/${DEPLOYMENT}"_*.yaml 2>/dev/null | tail -n +11 | xargs rm -f 2>/dev/null || true

    # 显示当前备份列表
    log_info ""
    log_info "当前备份列表:"
    ls -lh "${BACKUP_DIR}/${DEPLOYMENT}"_*.yaml 2>/dev/null | tail -5 || log_warning "未找到备份文件"

else
    log_error "备份失败！"
    exit 1
fi

log_success ""
log_success "=========================================="
log_success "  备份完成！"
log_success "=========================================="
log_info ""
log_info "使用方法:"
log_info "  # 恢复最新备份"
log_info "  ./scripts/rollback.sh $NAMESPACE $DEPLOYMENT"
log_info ""
log_info "  # 恢复指定备份"
log_info "  kubectl apply -f $BACKUP_FILE"
log_info ""
