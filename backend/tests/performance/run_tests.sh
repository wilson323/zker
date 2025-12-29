#!/bin/bash
# ZKER Performance Test Runner
# 运行所有性能测试脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
BASE_URL=${BASE_URL:-"http://localhost:8080"}
TENANT_ID=${TENANT_ID:-"test-tenant"}
RESULTS_DIR="./tests/performance/results"

# 创建结果目录
mkdir -p "$RESULTS_DIR"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}ZKER 性能测试套件${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "测试配置:"
echo "  BASE_URL: $BASE_URL"
echo "  TENANT_ID: $TENANT_ID"
echo "  RESULTS_DIR: $RESULTS_DIR"
echo ""

# 检查k6是否安装
if ! command -v k6 &> /dev/null; then
    echo -e "${RED}❌ k6未安装${NC}"
    echo "请访问 https://k6.io/docs/getting-started/installation/ 安装k6"
    exit 1
fi

echo -e "${GREEN}✅ k6已安装${NC}"
echo ""

# 测试1: 租户管理API负载测试
echo -e "${YELLOW}📊 测试1: 租户管理API负载测试${NC}"
k6 run \
    --env BASE_URL="$BASE_URL" \
    --out json="$RESULTS_DIR/tenant_load.json" \
    ./tests/performance/tenant_load_test.js
echo ""

# 测试2: 配额检查性能测试
echo -e "${YELLOW}📊 测试2: 配额检查性能测试${NC}"
k6 run \
    --env BASE_URL="$BASE_URL" \
    --env TENANT_ID="$TENANT_ID" \
    --out json="$RESULTS_DIR/quota_load.json" \
    ./tests/performance/quota_load_test.js
echo ""

# 测试3: 压力测试
echo -e "${YELLOW}📊 测试3: API压力测试${NC}"
k6 run \
    --env BASE_URL="$BASE_URL" \
    --out json="$RESULTS_DIR/stress_test.json" \
    ./tests/performance/stress_test.js
echo ""

# 生成HTML报告
echo -e "${YELLOW}📄 生成测试报告${NC}"

# 检查是否安装了k6-reporter
if command -v k6-reporter &> /dev/null; then
    # 生成HTML报告
    for json_file in "$RESULTS_DIR"/*.json; do
        test_name=$(basename "$json_file" .json)
        echo "生成报告: $test_name.html"

        k6 run \
            --env BASE_URL="$BASE_URL" \
            --out json="$json_file" \
            "$(dirname "$json_file")/${test_name}.js" \
        2>&1 | k6-reporter html "$RESULTS_DIR/${test_name}.html" || true
    done
else
    echo -e "${YELLOW}⚠️  k6-reporter未安装,跳过HTML报告生成${NC}"
    echo "安装命令: go install github.com/k6reporter/k6reporter/cmd/k6reporter@latest"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✅ 所有测试完成!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "测试结果保存在: $RESULTS_DIR"
echo ""
echo "查看JSON结果:"
echo "  cat $RESULTS_DIR/tenant_load.json"
echo ""
echo "查看HTML报告(如果生成):"
echo "  open $RESULTS_DIR/tenant_load.html"
echo ""
