#!/bin/bash

###############################################################################
# ZKER API规范自动修复工具
#
# 用途: 自动修复API规范违规问题
# 作者: API规范专家 (Claude Code)
# 日期: 2025-01-01
# 版本: v1.0
#
# 使用方法:
#   ./scripts/api-linter/fix-api-violations.sh [选项]
#
# 选项:
#   --check-only    仅检查,不修复
#   --fix-response  修复响应格式
#   --fix-methods   修复HTTP方法
#   --fix-routes    修复路由命名
#   --all           修复所有问题 (默认)
#
###############################################################################

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
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

# 检查项目目录
check_project_dir() {
    if [ ! -d "$BACKEND_DIR" ]; then
        log_error "后端目录不存在: $BACKEND_DIR"
        exit 1
    fi
    log_success "项目目录检查通过: $PROJECT_ROOT"
}

# 检查违规情况
check_violations() {
    log_info "开始检查API规范违规..."

    local violation_count=0

    echo ""
    echo "========================================"
    echo "📊 API规范违规检查报告"
    echo "========================================"
    echo ""

    # 1. 检查响应格式
    log_info "检查响应格式..."
    local response_violations=$(grep -rn "c\.JSON(consts\.StatusOK, resp)" "$BACKEND_DIR/api/handler" --include="*.go" 2>/dev/null | wc -l || echo "0")
    echo "  响应格式违规: $response_violations 处"
    ((violation_count+=response_violations))

    # 2. 检查POST用于查询
    log_info "检查POST用于查询..."
    local post_get_violations=$(grep -rn "\.POST.*\/get\|\.POST.*\/list\|\.POST.*\/detail" "$BACKEND_DIR/api/router" --include="*.go" 2>/dev/null | wc -l || echo "0")
    echo "  POST用于查询: $post_get_violations 处"
    ((violation_count+=post_get_violations))

    # 3. 检查POST用于删除
    log_info "检查POST用于删除..."
    local post_delete_violations=$(grep -rn "\.POST.*\/delete" "$BACKEND_DIR/api/router" --include="*.go" 2>/dev/null | wc -l || echo "0")
    echo "  POST用于删除: $post_delete_violations 处"
    ((violation_count+=post_delete_violations))

    # 4. 检查单数命名
    log_info "检查单数命名..."
    local singular_violations=$(grep -rn '\.Group("/[^/]*",[^)]*)' "$BACKEND_DIR/api/router" --include="*.go" | grep -v "s\|ies\|data\|info\|config" | wc -l || echo "0")
    echo "  单数命名: $singular_violations 处"

    # 5. 检查Swagger文档
    log_info "检查Swagger文档..."
    local total_files=$(find "$BACKEND_DIR/api/handler" -name "*.go" -type f | wc -l)
    local swagger_files=$(find "$BACKEND_DIR/api/handler" -name "*.go" -type f -exec grep -l "@Router" {} \; | wc -l)
    local swagger_coverage=$((swagger_files * 100 / total_files))
    echo "  Swagger文档覆盖率: $swagger_coverage% ($swagger_files/$total_files 文件)"

    echo ""
    echo "========================================"
    echo "总违规数: $violation_count 处"
    echo "========================================"
    echo ""

    if [ "$violation_count" -eq 0 ]; then
        log_success "未发现API规范违规!"
        return 0
    else
        log_warning "发现 $violation_count 处API规范违规"
        return 1
    fi
}

# 修复响应格式
fix_response_format() {
    log_info "开始修复响应格式..."

    local fixed_count=0
    local handler_dir="$BACKEND_DIR/api/handler"

    # 查找所有需要修复的文件
    local files=$(grep -rl "c\.JSON(consts\.StatusOK, resp)" "$handler_dir" --include="*.go" 2>/dev/null || true)

    if [ -z "$files" ]; then
        log_success "响应格式已全部符合规范"
        return 0
    fi

    echo "$files" | while read -r file; do
        log_info "处理文件: $file"

        # 检查是否已导入httputil
        if ! grep -q "backend/api/internal/httputil" "$file"; then
            log_warning "文件缺少httputil导入,请手动添加: $file"
            continue
        fi

        # 备份文件
        cp "$file" "$file.bak"

        # 替换响应调用
        sed -i 's/c\.JSON(consts\.StatusOK, resp)/httputil.BuildSuccessResp(c, resp)/g' "$file"
        sed -i 's/c\.JSON(consts\.StatusOK, \&resp)/httputil.BuildSuccessResp(c, \&resp)/g' "$file"

        # 检查修改是否成功
        if [ $? -eq 0 ]; then
            log_success "修复成功: $file"
            rm "$file.bak"
            ((fixed_count++))
        else
            log_error "修复失败: $file,恢复备份"
            mv "$file.bak" "$file"
        fi
    done

    log_success "响应格式修复完成,共修复 $fixed_count 个文件"
}

# 修复HTTP方法
fix_http_methods() {
    log_info "开始修复HTTP方法..."
    log_warning "HTTP方法修复需要手动操作,以下是需要修改的文件:"

    local violations_file="/tmp/http_method_violations.txt"

    # 查找所有违规
    {
        echo "========================================="
        echo "POST用于查询 (应改为GET)"
        echo "========================================="
        grep -rn "\.POST.*\/get\|\.POST.*\/list\|\.POST.*\/detail" "$BACKEND_DIR/api/router" --include="*.go" 2>/dev/null || true

        echo ""
        echo "========================================="
        echo "POST用于删除 (应改为DELETE)"
        echo "========================================="
        grep -rn "\.POST.*\/delete" "$BACKEND_DIR/api/router" --include="*.go" 2>/dev/null || true

        echo ""
        echo "========================================="
        echo "POST用于更新 (应改为PUT/PATCH)"
        echo "========================================="
        grep -rn "\.POST.*\/update" "$BACKEND_DIR/api/router" --include="*.go" 2>/dev/null || true
    } > "$violations_file"

    cat "$violations_file"

    log_info "详细违规清单已保存到: $violations_file"
    log_warning "请根据上述清单手动修改路由文件"
}

# 修复路由命名
fix_route_names() {
    log_info "开始检查路由命名..."
    log_warning "路由命名修复需要手动操作,以下是需要修改的地方:"

    local violations_file="/tmp/route_name_violations.txt"

    # 查找所有违规
    {
        echo "========================================="
        echo "单数命名 (应改为复数)"
        echo "========================================="
        grep -rn '\.Group("/[^"'\'']*s' "$BACKEND_DIR/api/router" --include="*.go" | grep -v "s\|ies\|data\|info\|config" || true

        echo ""
        echo "========================================="
        echo "驼峰命名 (应改为连字符)"
        echo "========================================="
        grep -rn "\.POST\|\.GET\|\.PUT\|\.DELETE" "$BACKEND_DIR/api/router" --include="*.go" | grep -E "[A-Z]" || true
    } > "$violations_file"

    cat "$violations_file"

    log_info "详细违规清单已保存到: $violations_file"
    log_warning "请根据上述清单手动修改路由命名"
}

# 生成Swagger文档
generate_swagger() {
    log_info "开始生成Swagger文档..."

    # 检查swag工具
    if ! command -v swag &> /dev/null; then
        log_warning "swag工具未安装,正在安装..."
        go install github.com/swaggo/swag/cmd/swag@latest
    fi

    cd "$BACKEND_DIR"

    # 生成文档
    log_info "正在生成Swagger文档..."
    if swag init -g api/main.go -o docs --parseInternal --parseDepth 1 2>/dev/null; then
        log_success "Swagger文档生成成功!"
        log_info "文档路径: $BACKEND_DIR/docs"
        log_info "访问地址: http://localhost:8080/swagger/index.html"
    else
        log_error "Swagger文档生成失败"
        log_warning "请检查main.go是否存在,或手动运行: swag init -g api/main.go"
    fi
}

# 主函数
main() {
    echo ""
    echo "========================================"
    echo "  ZKER API规范自动修复工具 v1.0"
    echo "========================================"
    echo ""

    check_project_dir

    # 解析命令行参数
    local check_only=false
    local fix_response=false
    local fix_methods=false
    local fix_routes=false
    local gen_swagger=false

    while [[ $# -gt 0 ]]; do
        case $1 in
            --check-only)
                check_only=true
                shift
                ;;
            --fix-response)
                fix_response=true
                shift
                ;;
            --fix-methods)
                fix_methods=true
                shift
                ;;
            --fix-routes)
                fix_routes=true
                shift
                ;;
            --gen-swagger)
                gen_swagger=true
                shift
                ;;
            --all)
                fix_response=true
                fix_methods=true
                fix_routes=true
                gen_swagger=true
                shift
                ;;
            *)
                log_error "未知选项: $1"
                echo "使用方法: $0 [选项]"
                echo "选项:"
                echo "  --check-only    仅检查,不修复"
                echo "  --fix-response  修复响应格式"
                echo "  --fix-methods   修复HTTP方法"
                echo "  --fix-routes    修复路由命名"
                echo "  --gen-swagger   生成Swagger文档"
                echo "  --all           修复所有问题 (默认)"
                exit 1
                ;;
        esac
    done

    # 如果没有指定选项,默认执行所有检查
    if [ "$check_only" = false ] && [ "$fix_response" = false ] && [ "$fix_methods" = false ] && [ "$fix_routes" = false ] && [ "$gen_swagger" = false ]; then
        check_only=true
    fi

    # 仅检查模式
    if [ "$check_only" = true ]; then
        check_violations
        exit $?
    fi

    # 修复模式
    echo ""
    log_info "开始修复API规范违规..."
    echo ""

    if [ "$fix_response" = true ]; then
        fix_response_format
        echo ""
    fi

    if [ "$fix_methods" = true ]; then
        fix_http_methods
        echo ""
    fi

    if [ "$fix_routes" = true ]; then
        fix_route_names
        echo ""
    fi

    if [ "$gen_swagger" = true ]; then
        generate_swagger
        echo ""
    fi

    # 最终检查
    log_info "最终检查..."
    check_violations

    echo ""
    log_success "修复工具执行完成!"
    echo ""
    echo "📋 下一步操作:"
    echo "  1. 查看修改后的代码"
    echo "  2. 运行测试: cd $BACKEND_DIR && go test ./..."
    echo "  3. 提交代码: git add . && git commit -m 'fix(api): 修复API规范违规'"
    echo ""
}

# 执行主函数
main "$@"
