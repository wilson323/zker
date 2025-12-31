#!/bin/bash

#############################################
# ZKER 一致性评分计算工具
# 用途: 自动计算代码一致性评分
# 作者: ZKER开发团队
# 更新: 2025-01-03
#############################################

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}ZKER 一致性评分计算${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 初始化分数
total_score=0
max_score=100

# ========================================
# 1. API响应格式检查 (20分)
# ========================================
echo -e "${BLUE}1️⃣  检查API响应格式统一性...${NC}"

if [ -d "$BACKEND_DIR/api/handler" ]; then
    api_violations=$(grep -r "c.JSON(" "$BACKEND_DIR/api/handler" --include="*_handler.go" 2>/dev/null | \
        grep -v "BuildSuccessResp\|BuildErrorResp\|http.StatusOK\|http.Status" | \
        wc -l | tr -d ' ')

    # 检查自定义响应结构
    custom_response=$(grep -r "type APIResponse\|type.*Response struct" "$BACKEND_DIR/api/handler" --include="*.go" 2>/dev/null | \
        grep -v "api/model\|domain" | wc -l | tr -d ' ')

    total_api_issues=$((api_violations + custom_response))

    if [ $total_api_issues -eq 0 ]; then
        api_score=20
        echo -e "   ${GREEN}✅ 无违规，得分: 20/20${NC}"
    elif [ $total_api_issues -le 5 ]; then
        api_score=15
        echo -e "   ${YELLOW}⚠️  发现 $total_api_issues 处违规，得分: 15/20${NC}"
    elif [ $total_api_issues -le 20 ]; then
        api_score=10
        echo -e "   ${YELLOW}⚠️  发现 $total_api_issues 处违规，得分: 10/20${NC}"
    else
        api_score=5
        echo -e "   ${RED}❌ 发现 $total_api_issues 处违规，得分: 5/20${NC}"
    fi
else
    api_score=0
    echo -e "   ${RED}❌ 未找到handler目录，得分: 0/20${NC}"
fi

total_score=$((total_score + api_score))
echo ""

# ========================================
# 2. 错误处理检查 (20分)
# ========================================
echo -e "${BLUE}2️⃣  检查错误处理统一性...${NC}"

if [ -d "$BACKEND_DIR/domain" ]; then
    error_violations=$(grep -r "errors\.New\|fmt\.Errorf" "$BACKEND_DIR/domain" --include="*.go" 2>/dev/null | \
        grep -v "test.go" | wc -l | tr -d ' ')

    if [ $error_violations -eq 0 ]; then
        error_score=20
        echo -e "   ${GREEN}✅ 无违规，得分: 20/20${NC}"
    elif [ $error_violations -le 50 ]; then
        error_score=18
        echo -e "   ${GREEN}✅ 发现 $error_violations 处违规，得分: 18/20${NC}"
    elif [ $error_violations -le 200 ]; then
        error_score=15
        echo -e "   ${YELLOW}⚠️  发现 $error_violations 处违规，得分: 15/20${NC}"
    elif [ $error_violations -le 500 ]; then
        error_score=10
        echo -e "   ${YELLOW}⚠️  发现 $error_violations 处违规，得分: 10/20${NC}"
    elif [ $error_violations -le 1000 ]; then
        error_score=7
        echo -e "   ${YELLOW}⚠️  发现 $error_violations 处违规，得分: 7/20${NC}"
    else
        error_score=5
        echo -e "   ${RED}❌ 发现 $error_violations 处违规，得分: 5/20${NC}"
    fi
else
    error_score=0
    echo -e "   ${RED}❌ 未找到domain目录，得分: 0/20${NC}"
fi

total_score=$((total_score + error_score))
echo ""

# ========================================
# 3. 多租户隔离检查 (20分)
# ========================================
echo -e "${BLUE}3️⃣  检查多租户隔离...${NC}"

if [ -d "$BACKEND_DIR/domain" ]; then
    # 查找可能缺少租户隔离的查询
    tenant_violations=$(grep -r "db\.\(Find\|First\)" "$BACKEND_DIR/domain" --include="*.go" 2>/dev/null | \
        grep -v "tenant_id\|test.go\|deleted_at" | \
        wc -l | tr -d ' ')

    if [ $tenant_violations -eq 0 ]; then
        tenant_score=20
        echo -e "   ${GREEN}✅ 无违规，得分: 20/20${NC}"
    elif [ $tenant_violations -le 5 ]; then
        tenant_score=18
        echo -e "   ${GREEN}✅ 发现 $tenant_violations 处可能违规，得分: 18/20${NC}"
    elif [ $tenant_violations -le 20 ]; then
        tenant_score=15
        echo -e "   ${YELLOW}⚠️  发现 $tenant_violations 处可能违规，得分: 15/20${NC}"
    else
        tenant_score=10
        echo -e "   ${YELLOW}⚠️  发现 $tenant_violations 处可能违规，得分: 10/20${NC}"
    fi
else
    tenant_score=0
    echo -e "   ${RED}❌ 未找到domain目录，得分: 0/20${NC}"
fi

total_score=$((total_score + tenant_score))
echo ""

# ========================================
# 4. 架构分层检查 (15分)
# ========================================
echo -e "${BLUE}4️⃣  检查架构分层...${NC}"

if [ -d "$BACKEND_DIR/api/handler" ]; then
    arch_violations=$(grep -r "repository\." "$BACKEND_DIR/api/handler" --include="*.go" 2>/dev/null | \
        grep -v "service" | wc -l | tr -d ' ')

    if [ $arch_violations -eq 0 ]; then
        arch_score=15
        echo -e "   ${GREEN}✅ 无违规，得分: 15/15${NC}"
    else
        arch_score=$((15 - arch_violations))
        if [ $arch_score -lt 0 ]; then arch_score=0; fi
        echo -e "   ${RED}❌ 发现 $arch_violations 处违规，得分: $arch_score/15${NC}"
    fi
else
    arch_score=0
    echo -e "   ${RED}❌ 未找到handler目录，得分: 0/15${NC}"
fi

total_score=$((total_score + arch_score))
echo ""

# ========================================
# 5. 代码规范检查 (10分)
# ========================================
echo -e "${BLUE}5️⃣  检查代码规范...${NC}"

if [ -d "$BACKEND_DIR" ]; then
    # 检查包名规范
    naming_violations=$(find "$BACKEND_DIR" -name "*.go" -exec grep -l "^package [A-Z]" {} \; 2>/dev/null | wc -l | tr -d ' ')

    if [ $naming_violations -eq 0 ]; then
        code_score=10
        echo -e "   ${GREEN}✅ 包名规范，得分: 10/10${NC}"
    else
        code_score=$((10 - naming_violations))
        if [ $code_score -lt 0 ]; then code_score=0; fi
        echo -e "   ${YELLOW}⚠️  发现 $naming_violations 处命名违规，得分: $code_score/10${NC}"
    fi
else
    code_score=0
    echo -e "   ${RED}❌ 未找到backend目录，得分: 0/10${NC}"
fi

total_score=$((total_score + code_score))
echo ""

# ========================================
# 6. 测试覆盖率检查 (10分)
# ========================================
echo -e "${BLUE}6️⃣  检查测试覆盖率...${NC}"

if [ -f "$BACKEND_DIR/coverage.out" ]; then
    coverage=$(go tool cover -func="$BACKEND_DIR/coverage.out" 2>/dev/null | grep total | awk '{print $3}' | sed 's/%//' || echo "0")

    if [ ! -z "$coverage" ]; then
        coverage_int=${coverage%.*}
        if [ $coverage_int -ge 90 ]; then
            test_score=10
            echo -e "   ${GREEN}✅ 测试覆盖率: ${coverage}%，得分: 10/10${NC}"
        elif [ $coverage_int -ge 85 ]; then
            test_score=9
            echo -e "   ${GREEN}✅ 测试覆盖率: ${coverage}%，得分: 9/10${NC}"
        elif [ $coverage_int -ge 80 ]; then
            test_score=8
            echo -e "   ${GREEN}✅ 测试覆盖率: ${coverage}%，得分: 8/10${NC}"
        elif [ $coverage_int -ge 70 ]; then
            test_score=6
            echo -e "   ${YELLOW}⚠️  测试覆盖率: ${coverage}%，得分: 6/10${NC}"
        elif [ $coverage_int -ge 60 ]; then
            test_score=4
            echo -e "   ${YELLOW}⚠️  测试覆盖率: ${coverage}%，得分: 4/10${NC}"
        else
            test_score=2
            echo -e "   ${RED}❌ 测试覆盖率: ${coverage}%，得分: 2/10${NC}"
        fi
    else
        test_score=0
        echo -e "   ${RED}❌ 无法解析覆盖率，得分: 0/10${NC}"
    fi
else
    test_score=0
    echo -e "   ${YELLOW}⚠️  未找到覆盖率报告，得分: 0/10${NC}"
    echo -e "   ${YELLOW}   提示: 运行 'go test ./... -coverprofile=coverage.out' 生成报告${NC}"
fi

total_score=$((total_score + test_score))
echo ""

# ========================================
# 7. 安全规范检查 (5分)
# ========================================
echo -e "${BLUE}7️⃣  检查安全规范...${NC}"

if [ -d "$BACKEND_DIR/domain" ]; then
    # 检查SQL注入风险
    sql_injection=$(grep -r "db\.Exec\|db\.Raw" "$BACKEND_DIR/domain" --include="*.go" 2>/dev/null | \
        grep -v "test.go" | wc -l | tr -d ' ')

    if [ $sql_injection -eq 0 ]; then
        security_score=5
        echo -e "   ${GREEN}✅ 无SQL注入风险，得分: 5/5${NC}"
    else
        security_score=$((5 - sql_injection))
        if [ $security_score -lt 0 ]; then security_score=0; fi
        echo -e "   ${YELLOW}⚠️  发现 $sql_injection 处SQL注入风险，得分: $security_score/5${NC}"
    fi
else
    security_score=0
    echo -e "   ${RED}❌ 未找到domain目录，得分: 0/5${NC}"
fi

total_score=$((total_score + security_score))
echo ""

# ========================================
# 输出总分
# ========================================
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}📊 总分: ${total_score}/100${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 评级
if [ $total_score -ge 95 ]; then
    echo -e "${GREEN}🏆 等级: 卓越 (Enterprise-Grade)${NC}"
    echo -e "${GREEN}   达到国际一流水准，可分享经验${NC}"
elif [ $total_score -ge 90 ]; then
    echo -e "${GREEN}✅ 等级: 优秀 (Excellent)${NC}"
    echo -e "${GREEN}   达到企业级标准，可定期审计${NC}"
elif [ $total_score -ge 85 ]; then
    echo -e "${GREEN}✅ 等级: 良好 (Good)${NC}"
    echo -e "${YELLOW}   接近企业级标准，需优化细节${NC}"
elif [ $total_score -ge 75 ]; then
    echo -e "${YELLOW}⚠️  等级: 需改进 (Needs Improvement)${NC}"
    echo -e "${YELLOW}   基本可用，需要系统性优化${NC}"
else
    echo -e "${RED}❌ 等级: 不合格 (Unacceptable)${NC}"
    echo -e "${RED}   存在严重问题，必须立即修复${NC}"
fi

echo ""
echo -e "${BLUE}详细报告:${NC}"
echo -e "  📄 GLOBAL_CONSISTENCY_CHECK_REPORT.md"
echo -e "  📋 CONSISTENCY_FIX_SUGGESTIONS.md"
echo -e "  📊 CONSISTENCY_SCORE.md"
echo ""

# ========================================
# 生成JSON报告（可选）
# ========================================
if [ "$1" == "--json" ]; then
    cat <<EOF
{
  "total_score": $total_score,
  "max_score": 100,
  "details": {
    "api_response_format": {
      "score": $api_score,
      "max_score": 20,
      "violations": $total_api_issues
    },
    "error_handling": {
      "score": $error_score,
      "max_score": 20,
      "violations": $error_violations
    },
    "tenant_isolation": {
      "score": $tenant_score,
      "max_score": 20,
      "violations": $tenant_violations
    },
    "architecture_layers": {
      "score": $arch_score,
      "max_score": 15,
      "violations": $arch_violations
    },
    "code_standards": {
      "score": $code_score,
      "max_score": 10,
      "violations": $naming_violations
    },
    "test_coverage": {
      "score": $test_score,
      "max_score": 10,
      "coverage": "${coverage}%"
    },
    "security": {
      "score": $security_score,
      "max_score": 5,
      "violations": $sql_injection
    }
  },
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
}
EOF
fi

# 返回退出码
if [ $total_score -lt 75 ]; then
    exit 1
else
    exit 0
fi
