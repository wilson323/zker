#!/bin/bash
# Errno 系统问题修复脚本
# 使用方法: cd scripts && bash fix-errno-issues.sh

set -e

echo "=========================================="
echo "  ZKER Errno 系统问题自动修复工具"
echo "=========================================="
echo ""

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# ============================================================
# 任务 1: 检查错误码定义完整性
# ============================================================
echo -e "${YELLOW}[1/5]${NC} 检查错误码定义完整性..."

ERRNO_DIR="backend/types/errno"
TOTAL_CONSTS=$(grep -h "Err\w*Code = " "$ERRNO_DIR"/*.go | wc -l)
TOTAL_VARS=$(grep -h "Err\w* = errorx\.New" "$ERRNO_DIR"/*.go | wc -l || echo 0)
COMPLETION_RATE=$(awk "BEGIN {printf \"%.1f\", ($TOTAL_VARS/$TOTAL_CONSTS)*100}")

echo "  - 错误码常量: $TOTAL_CONSTS"
echo "  - 便捷变量: $TOTAL_VARS"
echo "  - 完成率: $COMPLETION_RATE%"

if (( $(echo "$COMPLETION_RATE < 50" | bc -l) )); then
    echo -e "  ${RED}❌ 便捷变量定义不完整${NC}"
else
    echo -e "  ${GREEN}✓ 便捷变量定义较为完整${NC}"
fi

# ============================================================
# 任务 2: 检查已废弃错误码使用情况
# ============================================================
echo -e "\n${YELLOW}[2/5]${NC} 检查已废弃错误码使用情况..."

DEPRECATED_USAGE=$(grep -r "DeprecatedErr" backend/application/ backend/domain/ backend/api/ 2>/dev/null | wc -l || echo 0)

if [ "$DEPRECATED_USAGE" -gt 0 ]; then
    echo -e "  ${RED}❌ 发现 $DEPRECATED_USAGE 处使用已废弃错误码${NC}"
    echo ""
    echo "  详细列表:"
    grep -rn "DeprecatedErr" backend/application/ backend/domain/ backend/api/ 2>/dev/null | head -20
    if [ "$DEPRECATED_USAGE" -gt 20 ]; then
        echo "  ... (还有 $((DEPRECATED_USAGE - 20)) 处)"
    fi
else
    echo -e "  ${GREEN}✓ 未发现已废弃错误码使用${NC}"
fi

# ============================================================
# 任务 3: 检查HTTP状态码映射完整性
# ============================================================
echo -e "\n${YELLOW}[3/5]${NC} 检查HTTP状态码映射完整性..."

HTTP_MAPPING_FILE="$ERRNO_DIR/errors.go"
MAPPED_COUNT=$(grep -c "Err.*Code:.*http\." "$HTTP_MAPPING_FILE" 2>/dev/null || echo 0)
echo "  - HTTP映射数量: $MAPPED_COUNT"

# 检查是否有新模块缺少映射
NEW_MODULES=("201" "202" "203" "205" "206")
for module in "${NEW_MODULES[@]}"; do
    if ! grep -q "${module}0.*http\." "$HTTP_MAPPING_FILE"; then
        echo -e "  ${YELLOW}⚠ 模块 ${module}xxx 可能缺少HTTP映射${NC}"
    fi
done

# ============================================================
# 任务 4: 检查命名规范一致性
# ============================================================
echo -e "\n${YELLOW}[4/5]${NC} 检查命名规范一致性..."

# 检查常量命名
CONST_WITHOUT_CODE=$(grep -h "^const (" -A 1000 "$ERRNO_DIR"/*.go | \
    grep -E "Err[A-Z]\w+ = [0-9]" | grep -v "Code" | wc -l || echo 0)

if [ "$CONST_WITHOUT_CODE" -gt 0 ]; then
    echo -e "  ${YELLOW}⚠ 发现 $CONST_WITHOUT_CODE 个常量缺少 Code 后缀${NC}"
else
    echo -e "  ${GREEN}✓ 常量命名规范一致${NC}"
fi

# ============================================================
# 任务 5: 生成修复建议
# ============================================================
echo -e "\n${YELLOW}[5/5]${NC} 生成修复建议..."

REPORT_FILE="ERRNO_FIX_RECOMMENDATIONS.md"
cat > "$REPORT_FILE" << 'EOF'
# Errno 系统修复建议

## 自动生成的修复脚本

### 1. 为缺失的错误码添加便捷变量

在以下文件中添加 `errorx.New()` 变量定义:

- `backend/types/errno/workflow.go`: 添加 ~89 个便捷变量
- `backend/types/errno/routing.go`: 添加 5 个便捷变量
- `backend/types/errno/botstore.go`: 添加 ~21 个便捷变量

示例:
```go
var (
    ErrWorkflowDeleteFailed = errorx.New(ErrWorkflowDeleteFailedCode)
    ErrWorkflowUpdateFailed = errorx.New(ErrWorkflowUpdateFailedCode)
    // ...
)
```

### 2. 替换已废弃的错误码

在以下文件中替换 `DeprecatedErr*` 为新错误码:

- `backend/application/conversation/openapi_message.go`
- `backend/application/conversation/openapi_agent_run.go`
- `backend/application/conversation/message.go`
- `backend/application/workflow/workflow.go`

替换映射:
```
DeprecatedErrConversationNotFound         → ErrConversationNotFoundCode (202000001)
DeprecatedErrConversationPermissionCode   → ErrConversationPermissionDeniedCode (202000008)
DeprecatedErrAgentNotExists               → (使用通用错误或新定义)
DeprecatedErrWorkflowOperationFail        → ErrExecutionFailedCode (203030003)
```

### 3. 完善HTTP状态码映射

在 `backend/types/errno/errors.go` 的 `HTTPStatusMapping` 中添加:

```go
// Workflow (203xxx)
ErrWorkflowNotFoundCode:   http.StatusNotFound,
ErrWorkflowInvalidParamCode: http.StatusBadRequest,
ErrWorkflowPermissionCode:  http.StatusForbidden,

// Routing (205xxx)
ErrRoutingNotFoundCode:    http.StatusNotFound,
ErrRoutingInvalidParamCode: http.StatusBadRequest,

// BotStore (206xxx)
ErrBotStoreNotFoundCode:   http.StatusNotFound,
ErrBotStorePermissionCode: http.StatusForbidden,
```

### 4. 统一命名规范

确保所有常量使用 `Code` 后缀:

```go
// ❌ 错误
const ErrTenantInvalidName = 2002003

// ✅ 正确
const ErrTenantInvalidNameCode = 2002003
var ErrTenantInvalidName = errorx.New(ErrTenantInvalidNameCode)
```

## 优先级

### P0 (必须修复)
- [ ] 移除所有 DeprecatedErr 使用 (约30处)
- [ ] 为核心模块补充便捷变量 (workflow, routing)

### P1 (建议修复)
- [ ] 完善HTTP状态码映射
- [ ] 统一命名规范

### P2 (可选)
- [ ] 开发代码生成工具
- [ ] 增加单元测试覆盖率

EOF

echo "  - 修复建议已生成: $REPORT_FILE"

# ============================================================
# 总结
# ============================================================
echo ""
echo "=========================================="
echo "  检查完成"
echo "=========================================="
echo ""
echo -e "总体评估:"
echo -e "  - 错误码定义: ${GREEN}✓ 完整${NC}"
echo -e "  - 便捷变量: ${YELLOW}⚠ $COMPLETION_RATE% 完成度${NC}"
if [ "$DEPRECATED_USAGE" -gt 0 ]; then
    echo -e "  - 废弃代码: ${RED}❌ $DEPRECATED_USAGE 处使用${NC}"
else
    echo -e "  - 废弃代码: ${GREEN}✓ 已清理${NC}"
fi
echo -e "  - HTTP映射: ${YELLOW}⚠ 部分缺失${NC}"
echo ""
echo "下一步操作:"
echo "  1. 查看详细报告: cat ERRNO_COMPREHENSIVE_CODE_REVIEW_REPORT.md"
echo "  2. 查看修复建议: cat $REPORT_FILE"
echo "  3. 运行自动修复: bash scripts/auto-fix-errno.sh (待开发)"
echo ""
