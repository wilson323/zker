#!/bin/bash

#############################################
# API响应格式统一修复脚本
# 用途: 批量修复API响应格式
# 作者: ZKER开发团队
# 更新: 2025-01-03
#############################################

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}API响应格式统一修复工具${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 函数：修复单个文件
fix_file() {
    local file=$1
    local temp_file="${file}.tmp"

    echo -e "${YELLOW}修复文件: $file${NC}"

    # 备份原文件
    cp "$file" "${file}.backup"

    # 应用修复规则
    sed -E '
        # 删除自定义APIResponse结构
        /^type APIResponse struct/,/^}$/d;

        # 替换 c.JSON(http.StatusOK, &XXXResponse{Code:, Message:, Data:}) 为 BuildSuccessResp
        s/c\.JSON\(http\.StatusOK, \&[^{]+\{Code:[[:space:]]*0,[[:space:]]*Message:[[:space:]]*"[^"]*",[[:space:]]*Data:[[:space:]]*([^}]*)\}\)/httputil.BuildSuccessResp(c, \1)/g;

        # 替换 c.JSON(http.StatusOK, XXXResponse{Code:, Message:, Data:}) 为 BuildSuccessResp
        s/c\.JSON\(http\.StatusOK, [^{]+\{Code:[[:space:]]*0,[[:space:]]*Message:[[:space:]]*"[^"]*",[[:space:]]*Data:[[:space:]]*([^}]*)\}\)/httputil.BuildSuccessResp(c, \1)/g;

        # 替换 c.JSON(http.StatusInternalServerError, &XXXResponse{Code:, Message:}) 为 BuildErrorResp
        s/c\.JSON\(http\.StatusInternalServerError, \&[^{]+\{Code:[[:space:]]*[0-9]+,[[:space:]]*Message:[[:space:]]*"([^"]*)"\}\)/httputil.BuildErrorResp(c, errno.ErrInternalServer, "\1", "", nil)/g;

        # 替换 c.JSON(http.StatusBadRequest, &XXXResponse{Code:, Message:}) 为 BuildErrorResp
        s/c\.JSON\(http\.StatusBadRequest, \&[^{]+\{Code:[[:space:]]*[0-9]+,[[:space:]]*Message:[[:space:]]*"([^"]*)"\}\)/httputil.BuildErrorResp(c, errno.ErrInvalidParam, "\1", "", nil)/g;

        # 删除删除的自定义响应结构引用
        /APIResponse{/d;
    ' "$file" > "$temp_file"

    # 替换原文件
    mv "$temp_file" "$file"

    echo -e "${GREEN}✅ 修复完成${NC}"
    echo ""
}

# 函数：检查文件是否需要修复
needs_fix() {
    local file=$1

    # 检查是否包含自定义响应结构
    if grep -q "type APIResponse\|type.*Response struct" "$file" 2>/dev/null; then
        return 0
    fi

    # 检查是否直接使用c.JSON
    if grep -q "c.JSON(http.StatusOK\|c.JSON(http.StatusInternalServerError\|c.JSON(http.StatusBadRequest" "$file" 2>/dev/null; then
        return 0
    fi

    return 1
}

# 主逻辑
if [ -z "$1" ]; then
    echo -e "${YELLOW}用法: $0 [file|directory]${NC}"
    echo ""
    echo "示例:"
    echo "  $0 backend/api/handler/coze/monitoring/monitoring_handler.go"
    echo "  $0 backend/api/handler/coze/monitoring/"
    echo "  $0 backend/api/handler/coze/"
    echo ""
    echo "注意: 修复前会自动备份文件 (.backup)"
    exit 1
fi

TARGET="$1"

# 统计
total_files=0
fixed_files=0
skipped_files=0

if [ -f "$TARGET" ]; then
    # 单个文件
    total_files=1
    if needs_fix "$TARGET"; then
        fix_file "$TARGET"
        fixed_files=1
    else
        echo -e "${GREEN}✅ 文件无需修复: $TARGET${NC}"
        skipped_files=1
    fi
elif [ -d "$TARGET" ]; then
    # 目录
    echo -e "${BLUE}扫描目录: $TARGET${NC}"
    echo ""

    # 查找所有需要修复的文件
    while IFS= read -r -d '' file; do
        total_files=$((total_files + 1))

        if needs_fix "$file"; then
            fix_file "$file"
            fixed_files=$((fixed_files + 1))
        else
            skipped_files=$((skipped_files + 1))
        fi
    done < <(find "$TARGET" -name "*_handler.go" -type f -print0)
else
    echo -e "${RED}❌ 错误: $TARGET 不是文件或目录${NC}"
    exit 1
fi

# 汇总
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}修复完成${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "扫描文件: $total_files"
echo -e "${GREEN}已修复:   $fixed_files${NC}"
echo -e "${YELLOW}跳过:     $skipped_files${NC}"
echo ""

if [ $fixed_files -gt 0 ]; then
    echo -e "${YELLOW}⚠️  备份文件已创建 (.backup 后缀)${NC}"
    echo -e "${YELLOW}   请确认修复无误后删除备份文件${NC}"
    echo ""
    echo -e "删除备份命令:"
    echo -e "  find \"$TARGET\" -name '*.backup' -delete"
    echo ""
fi

echo -e "${BLUE}下一步:${NC}"
echo -e "  1. 运行测试: cd backend && go test ./..."
echo -e "  2. 代码审查: git diff"
echo -e "  3. 提交修复: git commit -am 'fix: 统一API响应格式'"
echo ""
