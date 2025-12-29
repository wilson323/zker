#!/bin/bash
# scripts/rollback.sh
# 一键回滚脚本 - 快速回滚到上一个稳定版本
# 用途: 当部署出现问题时,5分钟内快速回滚
# 遵循规范: 最小化故障影响 + 自动化回滚

set -e  # 遇到错误立即退出

# ========== 颜色输出定义 ==========
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# ========== 配置项 ==========
NAMESPACE="${NAMESPACE:-zker-prod}"
DEPLOYMENT="${DEPLOYMENT:-zker-api}"
TARGET_VERSION="${1:-}"
FORCE="${FORCE:-false}"
BACKUP_ENABLED="${BACKUP_ENABLED:-true}"

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
check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl 未安装"
        exit 1
    fi
}

check_namespace() {
    if ! kubectl get namespace $NAMESPACE &> /dev/null; then
        log_error "命名空间不存在: $NAMESPACE"
        exit 1
    fi
}

check_deployment() {
    if ! kubectl get deployment $DEPLOYMENT -n $NAMESPACE &> /dev/null; then
        log_error "Deployment 不存在: $DEPLOYMENT"
        exit 1
    fi
}

# ========== 版本管理函数 ==========
get_current_version() {
    kubectl get deployment $DEPLOYMENT -n $NAMESPACE \
      -o jsonpath='{.spec.template.metadata.labels.version}' || echo "unknown"
}

get_previous_version() {
    # 从 ReplicaSets 获取历史版本
    local prev_version=$(kubectl get replicasets -n $NAMESPACE \
      -l app=$DEPLOYMENT \
      --sort-by=.metadata.creationTimestamp \
      -o jsonpath='{.items[-2].spec.template.metadata.labels.version}' 2>/dev/null || echo "")

    if [ -z "$prev_version" ]; then
        log_warn "未找到上一个版本"
    fi

    echo "$prev_version"
}

list_available_versions() {
    log_step "可用的版本列表:"

    echo ""
    echo "📌 已部署的 ReplicaSets:"
    kubectl get replicasets -n $NAMESPACE \
      -l app=$DEPLOYMENT \
      --sort-by=.metadata.creationTimestamp \
      -o custom-columns="VERSION:.spec.template.metadata.labels.version,NAME:.metadata.name,AGE:.metadata.creationTimestamp,READY:.status.readyReplicas" | head -10

    echo ""
    echo "📌 Deployment 历史记录:"
    kubectl rollout history deployment/$DEPLOYMENT -n $NAMESPACE | head -10
}

# ========== 备份函数 ==========
backup_current_config() {
    if [ "$BACKUP_ENABLED" != "true" ]; then
        return
    fi

    log_step "备份当前配置..."

    local backup_dir="./rollback-backups"
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local backup_file="$backup_dir/${DEPLOYMENT}_${timestamp}.yaml"

    mkdir -p $backup_dir

    # 导出当前 Deployment 配置
    kubectl get deployment $DEPLOYMENT -n $NAMESPACE -o yaml > $backup_file

    log_info "✅ 配置已备份到: $backup_file"

    # 保留最近10个备份
    ls -t $backup_dir/${DEPLOYMENT}_*.yaml | tail -n +11 | xargs rm -f 2>/dev/null || true
}

# ========== 回滚函数 ==========
rollback_to_previous() {
    log_step "回滚到上一个版本..."

    # 使用 kubectl rollout undo
    kubectl rollout undo deployment/$DEPLOYMENT -n $NAMESPACE

    log_info "等待回滚完成..."
    kubectl rollout status deployment/$DEPLOYMENT -n $NAMESPACE --timeout=5m

    local new_version=$(get_current_version)
    log_info "✅ 已回滚到版本: $new_version"
}

rollback_to_specific_version() {
    local target_version=$1

    log_step "回滚到指定版本: $target_version"

    # 查找对应版本的 ReplicaSet
    local replicaset_name=$(kubectl get replicasets -n $NAMESPACE \
      -l app=$DEPLOYMENT,version=$target_version \
      -o jsonpath='{.items[0].metadata.name}')

    if [ -z "$replicaset_name" ]; then
        log_error "未找到版本 $target_version 对应的 ReplicaSet"
        return 1
    fi

    log_info "找到 ReplicaSet: $replicaset_name"

    # 恢复到指定 ReplicaSet
    kubectl rollout undo deployment/$DEPLOYMENT \
      --to-revision=$(kubectl get replicasets -n $NAMESPACE \
        -l app=$DEPLOYMENT,version=$target_version \
        -o jsonpath='{.items[0].metadata.annotations.deployment\.kubernetes\.io/revision}') \
      -n $NAMESPACE

    log_info "等待回滚完成..."
    kubectl rollout status deployment/$DEPLOYMENT -n $NAMESPACE --timeout=5m

    log_info "✅ 已回滚到版本: $target_version"
}

rollback_from_backup() {
    local backup_file=$1

    if [ ! -f "$backup_file" ]; then
        log_error "备份文件不存在: $backup_file"
        return 1
    fi

    log_step "从备份恢复: $backup_file"

    kubectl apply -f $backup_file

    log_info "等待恢复完成..."
    kubectl rollout status deployment/$DEPLOYMENT -n $NAMESPACE --timeout=5m

    log_info "✅ 已从备份恢复"
}

# ========== 验证函数 ==========
verify_rollback() {
    log_step "验证回滚结果..."

    # 检查 Pod 状态
    local ready_pods=$(kubectl get deployment $DEPLOYMENT -n $NAMESPACE \
      -o jsonpath='{.status.readyRepresents}')

    local desired_pods=$(kubectl get deployment $DEPLOYMENT -n $NAMESPACE \
      -o jsonpath='{.spec.replicas}')

    if [ "$ready_pods" -lt "$desired_pods" ]; then
        log_error "Pod 未就绪: $ready_pods/$desired_pods"
        return 1
    fi

    # 检查版本
    local current_version=$(get_current_version)
    log_info "当前版本: $current_version"

    # 显示 Pod 状态
    echo ""
    log_info "Pod 状态:"
    kubectl get pods -n $NAMESPACE -l app=$DEPLOYMENT

    log_info "✅ 回滚验证通过"
}

# ========== Istio 回滚 (如果使用灰度发布) ==========
rollback_istio_canary() {
    log_step "回滚 Istio 灰度发布..."

    # 检查是否存在 VirtualService
    if kubectl get virtualservice $DEPLOYMENT -n $NAMESPACE &> /dev/null; then
        # 将流量切回旧版本
        cat <<EOF | kubectl apply -f -
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: $DEPLOYMENT
  namespace: $NAMESPACE
spec:
  hosts:
  - $DEPLOYMENT
  http:
  - route:
    - destination:
        host: $DEPLOYMENT
        subset: v1
      weight: 100
EOF

        log_info "✅ Istio 流量已回滚"
    fi

    # 删除金丝雀 Deployment
    if kubectl get deployment ${DEPLOYMENT}-canary -n $NAMESPACE &> /dev/null; then
        kubectl delete deployment ${DEPLOYMENT}-canary -n $NAMESPACE
        log_info "✅ 金丝雀 Deployment 已删除"
    fi
}

# ========== 数据库回滚 (可选) ==========
rollback_database() {
    local db_version=$1

    log_warn "数据库回滚需要手动执行:"
    echo ""
    echo "1. 检查数据库迁移状态:"
    echo "   ./scripts/migrate.sh status"
    echo ""
    echo "2. 回滚数据库迁移 (如果需要):"
    echo "   ./scripts/migrate.sh down $db_version"
    echo ""
    echo "⚠️  注意: 数据库回滚可能导致数据丢失,请谨慎操作!"

    read -p "是否需要回滚数据库? (yes/no): " confirm

    if [ "$confirm" = "yes" ]; then
        log_step "执行数据库回滚..."
        # TODO: 执行数据库回滚逻辑
        log_warn "数据库回滚功能待实现"
    fi
}

# ========== 通知函数 ==========
send_notification() {
    local status=$1
    local version=$(get_current_version)

    log_info "========================================="
    log_info "     回滚完成!"
    log_info "========================================="
    echo ""
    echo "状态: $status"
    echo "当前版本: $version"
    echo "时间: $(date '+%Y-%m-%d %H:%M:%S')"
    echo ""

    # TODO: 发送通知到 Slack/DingTalk/Email
    # curl -X POST $SLACK_WEBHOOK \
    #   -H 'Content-Type: application/json' \
    #   -d "{\"text\":\"🔄 $DEPLOYMENT 已回滚到版本 $version\"}"
}

# ========== 帮助函数 ==========
show_usage() {
    cat << EOF
一键回滚脚本 v1.0 - ZKER Deployment

用法: $0 [options] [version]

选项:
  -n, --namespace <name>      指定命名空间 (默认: zker-prod)
  -d, --deployment <name>     指定 Deployment (默认: zker-api)
  -f, --force                 强制回滚,不询问确认
  --no-backup                 不创建备份
  -l, --list                  列出可用版本
  -h, --help                  显示此帮助信息

参数:
  version                     目标版本号 (可选)

示例:
  # 回滚到上一个版本
  $0

  # 回滚到指定版本
  $0 v1.0.0

  # 列出可用版本
  $0 --list

  # 指定命名空间和 Deployment
  $0 -n zker-staging -d zker-api v1.0.1

  # 强制回滚 (不确认)
  $0 --force v1.0.0

环境变量:
  NAMESPACE                   命名空间
  DEPLOYMENT                  Deployment 名称
  FORCE                       强制执行
  BACKUP_ENABLED              启用备份

EOF
}

# ========== 主函数 ==========
main() {
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -n|--namespace)
                NAMESPACE="$2"
                shift 2
                ;;
            -d|--deployment)
                DEPLOYMENT="$2"
                shift 2
                ;;
            -f|--force)
                FORCE=true
                shift
                ;;
            --no-backup)
                BACKUP_ENABLED=false
                shift
                ;;
            -l|--list)
                check_kubectl
                check_namespace
                list_available_versions
                exit 0
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            -*)
                log_error "未知选项: $1"
                show_usage
                exit 1
                ;;
            *)
                TARGET_VERSION="$1"
                shift
                ;;
        esac
    done

    # 显示标题
    log_info "========================================="
    log_info "     一键回滚 - ZKER $DEPLOYMENT"
    log_info "========================================="
    echo ""

    # 检查
    check_kubectl
    check_namespace
    check_deployment

    # 显示当前状态
    local current_version=$(get_current_version)
    log_info "当前版本: $current_version"

    # 如果没有指定版本,显示上一个版本
    if [ -z "$TARGET_VERSION" ]; then
        local prev_version=$(get_previous_version)
        if [ -n "$prev_version" ]; then
            log_info "上一个版本: $prev_version"
        fi
    fi

    # 确认回滚
    if [ "$FORCE" != "true" ]; then
        echo ""
        read -p "确认要回滚 $DEPLOYMENT? (yes/no): " confirm

        if [ "$confirm" != "yes" ]; then
            log_info "回滚已取消"
            exit 0
        fi
    fi

    # 备份当前配置
    backup_current_config

    # 执行回滚
    if [ -n "$TARGET_VERSION" ]; then
        rollback_to_specific_version "$TARGET_VERSION"
    else
        rollback_to_previous
    fi

    # 如果使用了 Istio 灰度,也需要回滚
    rollback_istio_canary

    # 验证回滚
    if verify_rollback; then
        send_notification "success"
    else
        log_error "回滚验证失败!"
        send_notification "failed"
        exit 1
    fi

    # 提示数据库回滚
    if [ -n "$TARGET_VERSION" ]; then
        rollback_database "$TARGET_VERSION"
    fi

    echo ""
    log_info "========================================="
    log_info "     回滚成功完成!"
    log_info "========================================="
}

# 执行主函数
main "$@"
