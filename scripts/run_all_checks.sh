#!/bin/bash

#############################################
# ZKER 全局一致性一键检查工具
# 用途: 运行所有一致性检查并生成报告
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

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPORT_DIR="$PROJECT_ROOT/reports"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# 创建报告目录
mkdir -p "$REPORT_DIR"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}ZKER 全局一致性自动检查${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${BLUE}检查时间: $(date '+%Y-%m-%d %H:%M:%S')${NC}"
echo -e "${BLUE}报告目录: $REPORT_DIR${NC}"
echo ""

# ========================================
# 1. API响应格式检查
# ========================================
echo -e "${BLUE}1️⃣  API响应格式检查...${NC}"
echo ""

BACKEND_DIR="$PROJECT_ROOT/backend"

if [ -d "$BACKEND_DIR/api/handler" ]; then
    echo "查找直接使用 c.JSON() 的handler:"
    direct_json=$(grep -rn "c.JSON(" "$BACKEND_DIR/api/handler" --include="*_handler.go" 2>/dev/null | \
        grep -v "BuildSuccessResp\|BuildErrorResp\|http.StatusOK\|http.Status" || true)

    if [ -z "$direct_json" ]; then
        echo -e "   ${GREEN}✅ 未发现直接使用 c.JSON() 的代码${NC}"
    else
        echo -e "   ${YELLOW}⚠️  发现以下违规:${NC}"
        echo "$direct_json" | head -20
        echo "$direct_json" | wc -l | xargs -I {} echo "   总计: {} 处"
    fi

    echo ""
    echo "查找自定义响应结构:"
    custom_response=$(grep -rn "type APIResponse\|type.*Response struct" "$BACKEND_DIR/api/handler" --include="*.go" 2>/dev/null | \
        grep -v "api/model\|domain" || true)

    if [ -z "$custom_response" ]; then
        echo -e "   ${GREEN}✅ 未发现自定义响应结构${NC}"
    else
        echo -e "   ${YELLOW}⚠️  发现以下自定义结构:${NC}"
        echo "$custom_response"
    fi
else
    echo -e "   ${RED}❌ 未找到 handler 目录${NC}"
fi

echo ""

# ========================================
# 2. 错误处理检查
# ========================================
echo -e "${BLUE}2️⃣  错误处理检查...${NC}"
echo ""

if [ -d "$BACKEND_DIR/domain" ]; then
    echo "查找硬编码错误 (errors.New, fmt.Errorf):"

    # 统计总违规次数
    total_errors=$(grep -r "errors\.New\|fmt\.Errorf" "$BACKEND_DIR/domain" --include="*.go" 2>/dev/null | \
        grep -v "test.go" | wc -l | tr -d ' ')

    echo "   总违规次数: $total_errors"

    # 按模块统计
    echo ""
    echo "按模块统计 (Top 10):"
    echo ""
    printf "%-50s %10s\n" "模块" "违规次数"
    printf "%-50s %10s\n" "----" "--------"

    for dir in "$BACKEND_DIR/domain"/*/; do
        if [ -d "$dir" ]; then
            count=$(grep -r "errors\.New\|fmt\.Errorf" "$dir" --include="*.go" 2>/dev/null | \
                    grep -v "test.go" | wc -l | tr -d ' ')

            if [ $count -gt 0 ]; then
                module_name=$(basename "$dir")
                printf "%-50s %10d\n" "$module_name" "$count"
            fi
        fi
    done | sort -k2 -rn | head -10

    echo ""
    echo "高风险文件 (违规 > 10次):"
    echo ""

    grep -r "errors\.New\|fmt\.Errorf" "$BACKEND_DIR/domain" --include="*.go" 2>/dev/null | \
        grep -v "test.go" | \
        cut -d: -f1 | \
        sort | uniq -c | \
        awk '$1 > 10 {printf "%5d %s\n", $1, $2}' | \
        sort -rn || echo "   无高风险文件"
else
    echo -e "   ${RED}❌ 未找到 domain 目录${NC}"
fi

echo ""

# ========================================
# 3. 架构分层检查
# ========================================
echo -e "${BLUE}3️⃣  架构分层检查...${NC}"
echo ""

if [ -d "$BACKEND_DIR/api/handler" ]; then
    echo "检查 Handler → Repository 违规:"

    repo_violations=$(grep -rn "repository\." "$BACKEND_DIR/api/handler" --include="*.go" 2>/dev/null | \
        grep -v "service" || true)

    if [ -z "$repo_violations" ]; then
        echo -e "   ${GREEN}✅ 未发现 Handler 直接调用 Repository${NC}"
    else
        echo -e "   ${RED}❌ 发现以下违规:${NC}"
        echo "$repo_violations"
    fi

    echo ""
    echo "检查 Handler → Domain 违规:"

    domain_violations=$(grep -rn "domain\." "$BACKEND_DIR/api/handler" --include="*.go" 2>/dev/null | \
        grep -v "entity" || true)

    if [ -z "$domain_violations" ]; then
        echo -e "   ${GREEN}✅ 未发现 Handler 直接调用 Domain${NC}"
    else
        echo -e "   ${YELLOW}⚠️  发现以下调用 (需确认是否合理):${NC}"
        echo "$domain_violations" | head -10
    fi
else
    echo -e "   ${RED}❌ 未找到 handler 目录${NC}"
fi

echo ""

# ========================================
# 4. 多租户隔离检查
# ========================================
echo -e "${BLUE}4️⃣  多租户隔离检查...${NC}"
echo ""

if [ -d "$BACKEND_DIR/domain" ]; then
    echo "查找可能缺少 tenant_id 的查询:"

    missing_tenant=$(grep -rn "db\.\(Find\|First\)" "$BACKEND_DIR/domain" --include="*.go" 2>/dev/null | \
        grep -v "tenant_id\|test.go\|deleted_at\|\.Where" | \
        head -20 || true)

    if [ -z "$missing_tenant" ]; then
        echo -e "   ${GREEN}✅ 未发现明显缺少租户隔离的查询${NC}"
    else
        echo -e "   ${YELLOW}⚠️  发现以下可能违规 (需人工审查):${NC}"
        echo "$missing_tenant"
    fi

    echo ""
    echo "优秀实践示例 (已包含 tenant_id):"

    good_examples=$(grep -rn "tenant_id.*AND deleted_at" "$BACKEND_DIR/domain/org" --include="*.go" 2>/dev/null | \
        head -5 || true)

    if [ -z "$good_examples" ]; then
        echo "   未找到示例"
    else
        echo "$good_examples"
    fi
else
    echo -e "   ${RED}❌ 未找到 domain 目录${NC}"
fi

echo ""

# ========================================
# 5. 代码规范检查
# ========================================
echo -e "${BLUE}5️⃣  代码规范检查...${NC}"
echo ""

if [ -d "$BACKEND_DIR" ]; then
    echo "检查包名规范 (应该小写):"

    bad_packages=$(find "$BACKEND_DIR" -name "*.go" -exec grep -l "^package [A-Z]" {} \; 2>/dev/null || true)

    if [ -z "$bad_packages" ]; then
        echo -e "   ${GREEN}✅ 所有包名符合规范${NC}"
    else
        echo -e "   ${RED}❌ 发现以下大写包名:${NC}"
        echo "$bad_packages"
    fi

    echo ""
    echo "检查接口命名 (应该大写开头):"

    bad_interfaces=$(grep -rn "type [a-z]" "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
        grep "interface" | head -10 || true)

    if [ -z "$bad_interfaces" ]; then
        echo -e "   ${GREEN}✅ 未发现小写接口名${NC}"
    else
        echo -e "   ${YELLOW}⚠️  发现以下小写接口 (需确认):${NC}"
        echo "$bad_interfaces"
    fi

    echo ""
    echo "检查常量命名 (应该大写):"

    bad_constants=$(grep -rn "const [a-z]" "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
        grep -v "test.go" | head -10 || true)

    if [ -z "$bad_constants" ]; then
        echo -e "   ${GREEN}✅ 未发现小写常量${NC}"
    else
        echo -e "   ${YELLOW}⚠️  发现以下小写常量 (需确认):${NC}"
        echo "$bad_constants"
    fi
else
    echo -e "   ${RED}❌ 未找到 backend 目录${NC}"
fi

echo ""

# ========================================
# 6. 测试覆盖率检查
# ========================================
echo -e "${BLUE}6️⃣  测试覆盖率检查...${NC}"
echo ""

if [ -f "$BACKEND_DIR/coverage.out" ]; then
    echo "解析覆盖率报告..."

    if command -v go &> /dev/null; then
        coverage=$(go tool cover -func="$BACKEND_DIR/coverage.out" 2>/dev/null | grep total || echo "")

        if [ ! -z "$coverage" ]; then
            echo "$coverage"

            coverage_percent=$(echo "$coverage" | awk '{print $3}' | sed 's/%//')
            coverage_int=${coverage_percent%.*}

            echo ""
            if [ $coverage_int -ge 85 ]; then
                echo -e "   ${GREEN}✅ 覆盖率达标 (≥85%)${NC}"
            elif [ $coverage_int -ge 75 ]; then
                echo -e "   ${YELLOW}⚠️  覆盖率接近目标${NC}"
            else
                echo -e "   ${RED}❌ 覆盖率未达标 (<75%)${NC}"
            fi
        else
            echo -e "   ${YELLOW}⚠️  无法解析覆盖率报告${NC}"
        fi
    else
        echo -e "   ${YELLOW}⚠️  未安装 go 工具${NC}"
    fi
else
    echo -e "   ${YELLOW}⚠️  未找到覆盖率报告${NC}"
    echo "   提示: 运行以下命令生成报告:"
    echo "   cd backend && go test ./... -coverprofile=coverage.out"
fi

echo ""

# ========================================
# 7. 安全规范检查
# ========================================
echo -e "${BLUE}7️⃣  安全规范检查...${NC}"
echo ""

if [ -d "$BACKEND_DIR/domain" ]; then
    echo "检查 SQL 注入风险 (db.Exec, db.Raw):"

    sql_risks=$(grep -rn "db\.Exec\|db\.Raw" "$BACKEND_DIR/domain" --include="*.go" 2>/dev/null | \
        grep -v "test.go" || true)

    if [ -z "$sql_risks" ]; then
        echo -e "   ${GREEN}✅ 未发现 SQL 注入风险${NC}"
    else
        echo -e "   ${YELLOW}⚠️  发现以下潜在风险 (需人工审查):${NC}"
        echo "$sql_risks" | head -10
    fi

    echo ""
    echo "检查硬编码密钥/密码:"

    secrets=$(grep -rn "password.*=.*['\"][^'\"]*['\"]" "$BACKEND_DIR" --include="*.go" 2>/dev/null | \
        grep -v "test.go\|config.go\|example" | head -10 || true)

    if [ -z "$secrets" ]; then
        echo -e "   ${GREEN}✅ 未发现硬编码密钥${NC}"
    else
        echo -e "   ${RED}❌ 发现以下硬编码密钥 (必须移除):${NC}"
        echo "$secrets"
    fi
else
    echo -e "   ${RED}❌ 未找到 domain 目录${NC}"
fi

echo ""

# ========================================
# 8. 生成汇总报告
# ========================================
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}✅ 检查完成！${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

echo -e "${BLUE}详细报告:${NC}"
echo -e "  📄 GLOBAL_CONSISTENCY_CHECK_REPORT.md"
echo -e "  📋 CONSISTENCY_FIX_SUGGESTIONS.md"
echo -e "  📊 CONSISTENCY_SCORE.md"
echo ""

echo -e "${BLUE}下一步操作:${NC}"
echo -e "  1. 查看评分: bash scripts/calculate_consistency_score.sh"
echo -e "  2. 修复P0问题: 参考 CONSISTENCY_FIX_SUGGESTIONS.md"
echo -e "  3. 生成覆盖率: cd backend && go test ./... -coverprofile=coverage.out"
echo ""

# 保存报告到文件
REPORT_FILE="$REPORT_DIR/consistency_check_$TIMESTAMP.txt"

{
    echo "ZKER 全局一致性检查报告"
    echo "生成时间: $(date '+%Y-%m-%d %H:%M:%S')"
    echo ""
    echo "=========================================="
    echo "检查结果汇总"
    echo "=========================================="
    echo ""
    echo "1. API响应格式: 见上述输出"
    echo "2. 错误处理: 见上述输出"
    echo "3. 架构分层: 见上述输出"
    echo "4. 多租户隔离: 见上述输出"
    echo "5. 代码规范: 见上述输出"
    echo "6. 测试覆盖率: 见上述输出"
    echo "7. 安全规范: 见上述输出"
    echo ""
    echo "完整检查结果见上述输出"
} > "$REPORT_FILE"

echo -e "${GREEN}✅ 报告已保存到: $REPORT_FILE${NC}"
echo ""

# 返回退出码
bash "$PROJECT_ROOT/scripts/calculate_consistency_score.sh"
