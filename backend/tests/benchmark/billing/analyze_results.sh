#!/bin/bash

###############################################################################
# 性能测试结果分析脚本
#
# 功能：
#   1. 解析基准测试结果
#   2. 提取关键性能指标
#   3. 生成可视化图表（使用gnuplot）
#   4. 对比基线指标
#   5. 提供优化建议
#
# 用法：
#   ./analyze_results.sh <result_file>
#
# @author 研发B
# @date 2025-12-30
###############################################################################

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
RESULTS_DIR="$PROJECT_ROOT/backend/tests/results/benchmark"

# 性能基线（基于ZKER性能基线文档）
BASELINE_RECORD_TOKEN_USAGE=100  # ms
BASELINE_BATCH_RECORD_1000=1000  # ms
BASELINE_GET_USAGE_STATS=50       # ms
BASELINE_CHECK_BUDGET=50          # ms
BASELINE_PRICING_ENGINE_10000=100 # ms

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

# 检查参数
if [ $# -lt 1 ]; then
    log_error "用法: $0 <result_file>"
    echo ""
    echo "示例:"
    echo "  $0 benchmark_20251230_120000.txt"
    echo "  $0 benchmark_latest.txt"
    exit 1
fi

RESULT_FILE="$1"

# 检查文件是否存在
if [ ! -f "$RESULT_FILE" ]; then
    # 尝试在结果目录中查找
    if [ -f "$RESULTS_DIR/$RESULT_FILE" ]; then
        RESULT_FILE="$RESULTS_DIR/$RESULT_FILE"
    else
        log_error "文件不存在: $RESULT_FILE"
        exit 1
    fi
fi

# 打印文件信息
print_file_info() {
    echo ""
    echo "========================================"
    echo "性能测试结果分析"
    echo "========================================"
    echo "文件: $RESULT_FILE"
    echo "大小: $(du -h "$RESULT_FILE" | cut -f1)"
    echo "修改时间: $(stat -c %y "$RESULT_FILE" 2>/dev/null || stat -f %Sm "$RESULT_FILE")"
    echo ""
}

# 提取基准测试结果
extract_benchmarks() {
    log_info "提取基准测试结果..."

    # 提取所有Benchmark行
    grep "^Benchmark" "$RESULT_FILE" | while read -r line; do
        # 解析测试名称
        test_name=$(echo "$line" | awk '{print $1}')

        # 解析ns/op
        ns_op=$(echo "$line" | awk '{print $3}')

        # 解析B/op (bytes per operation)
        b_op=$(echo "$line" | awk '{print $5}')

        # 解析allocs/op
        allocs_op=$(echo "$line" | awk '{print $7}')

        # 计算ms/op (更易读的单位)
        ms_op=$(echo "scale=2; $ns_op / 1000000" | bc 2>/dev/null || echo "N/A")

        echo "$test_name|$ms_op|$ns_op|$b_op|$allocs_op"
    done
}

# 对比基线
compare_baseline() {
    local test_name=$1
    local actual_ms=$2

    case "$test_name" in
        *RecordTokenUsage*)
            baseline=$BASELINE_RECORD_TOKEN_USAGE
            ;;
        *BatchRecordTokenUsage*1000*)
            baseline=$BASELINE_BATCH_RECORD_1000
            ;;
        *GetUsageStats*)
            baseline=$BASELINE_GET_USAGE_STATS
            ;;
        *CheckBudget*)
            baseline=$BASELINE_CHECK_BUDGET
            ;;
        *PricingEngine*10000*)
            baseline=$BASELINE_PRICING_ENGINE_10000
            ;;
        *)
            baseline="N/A"
            ;;
    esac

    if [ "$baseline" != "N/A" ]; then
        # 计算是否达标
        is_passed=$(echo "$actual_ms < $baseline" | bc 2>/dev/null || echo "0")

        if [ "$is_passed" -eq 1 ]; then
            echo -e " ${GREEN}✓${NC} 基线: < ${baseline}ms"
        else
            echo -e " ${RED}✗${NC} 基线: < ${baseline}ms (未达标)"
        fi
    else
        echo "  基线: N/A"
    fi
}

# 生成性能指标表格
generate_metrics_table() {
    log_info "生成性能指标表格..."

    echo ""
    echo "========================================"
    echo "性能指标汇总"
    echo "========================================"
    echo ""
    printf "%-50s %12s %12s %12s %10s\n" "测试名称" "耗时(ms/op)" "ns/op" "B/op" "allocs/op"
    echo "--------------------------------------------------------------------------------------------------------"

    extract_benchmarks "$RESULT_FILE" | while IFS='|' read -r test_name ms_op ns_op b_op allocs_op; do
        printf "%-50s %12s %12s %12s %10s\n" "$test_name" "$ms_op" "$ns_op" "$b_op" "$allocs_op"

        # 对比基线
        if [ "$ms_op" != "N/A" ]; then
            compare_baseline "$test_name" "$ms_op"
        fi
        echo ""
    done
}

# 生成性能基线对比报告
generate_baseline_report() {
    log_info "生成性能基线对比报告..."

    local report_file="${RESULT_FILE%.txt}_baseline_report.md"

    cat > "$report_file" << EOF
# 性能基线对比报告

**测试文件**: $(basename "$RESULT_FILE")
**分析时间**: $(date +"%Y-%m-%d %H:%M:%S")

---

## 📊 性能基线对比

| 测试名称 | 实际值 (ms/op) | 基线值 (ms/op) | 状态 |
|---------|----------------|----------------|------|
EOF

    # 填充表格
    extract_benchmarks "$RESULT_FILE" | while IFS='|' read -r test_name ms_op ns_op b_op allocs_op; do
        if [ "$ms_op" != "N/A" ]; then
            baseline="N/A"
            status="⚪️"

            case "$test_name" in
                *RecordTokenUsage*)
                    baseline=$BASELINE_RECORD_TOKEN_USAGE
                    is_passed=$(echo "$ms_op < $baseline" | bc 2>/dev/null || echo "0")
                    if [ "$is_passed" -eq 1 ]; then
                        status="✅"
                    else
                        status="❌"
                    fi
                    ;;
                *BatchRecordTokenUsage*1000*)
                    baseline=$BASELINE_BATCH_RECORD_1000
                    is_passed=$(echo "$ms_op < $baseline" | bc 2>/dev/null || echo "0")
                    if [ "$is_passed" -eq 1 ]; then
                        status="✅"
                    else
                        status="❌"
                    fi
                    ;;
                *GetUsageStats*)
                    baseline=$BASELINE_GET_USAGE_STATS
                    is_passed=$(echo "$ms_op < $baseline" | bc 2>/dev/null || echo "0")
                    if [ "$is_passed" -eq 1 ]; then
                        status="✅"
                    else
                        status="❌"
                    fi
                    ;;
                *CheckBudget*)
                    baseline=$BASELINE_CHECK_BUDGET
                    is_passed=$(echo "$ms_op < $baseline" | bc 2>/dev/null || echo "0")
                    if [ "$is_passed" -eq 1 ]; then
                        status="✅"
                    else
                        status="❌"
                    fi
                    ;;
                *PricingEngine*10000*)
                    baseline=$BASELINE_PRICING_ENGINE_10000
                    is_passed=$(echo "$ms_op < $baseline" | bc 2>/dev/null || echo "0")
                    if [ "$is_passed" -eq 1 ]; then
                        status="✅"
                    else
                        status="❌"
                    fi
                    ;;
            esac

            echo "| $test_name | $ms_op | $baseline | $status |" >> "$report_file"
        fi
    done

    cat >> "$report_file" << EOF

---

## 🔧 性能优化建议

基于当前测试结果，提供以下优化建议：

### 1. 数据库优化

**建议**:
- 检查慢查询日志，优化索引
- 使用 \`EXPLAIN ANALYZE\` 分析查询计划
- 考虑使用读写分离

**实施**:
\`\`\`sql
-- 检查缺失的索引
SELECT * FROM sys.schema_unused_indexes
WHERE object_schema = 'zker_production';

-- 分析慢查询
SELECT * FROM mysql.slow_log
ORDER BY query_time DESC LIMIT 100;
\`\`\`

### 2. 缓存优化

**建议**:
- 实现多级缓存（本地缓存 + Redis）
- 预热热点数据
- 设置合理的TTL

**实施**:
\`\`\`go
// L1: 本地缓存
// L2: Redis缓存
// L3: 数据库
func GetBot(ctx context.Context, botID string) (*Bot, error) {
    // L1: 本地缓存 (5分钟TTL)
    if bot, ok := localCache.Get(botID); ok {
        return bot, nil
    }

    // L2: Redis缓存 (30分钟TTL)
    val, err := redisClient.Get(ctx, "bot:"+botID).Result()
    if err == nil {
        var bot Bot
        json.Unmarshal([]byte(val), &bot)
        localCache.Set(botID, &bot, 5*time.Minute)
        return &bot, nil
    }

    // L3: 数据库
    bot, err := botRepo.FindByID(ctx, botID)
    // ...
}
\`\`\`

### 3. 并发优化

**建议**:
- 使用worker pool限制并发数
- 避免创建过多的goroutine
- 使用context控制超时

**实施**:
\`\`\`go
// 使用worker pool
sem := make(chan struct{}, 100) // 最多100个并发
var wg sync.WaitGroup

for _, id := range botIDs {
    wg.Add(1)
    sem <- struct{}{}
    go func(botID string) {
        defer wg.Done()
        defer func() { <-sem }()
        processBot(botID)
    }(id)
}
wg.Wait()
\`\`\`

### 4. 内存优化

**建议**:
- 使用对象池减少内存分配
- 避免频繁的字符串拼接
- 及时释放大对象

**实施**:
\`\`\`go
// 使用sync.Pool
var botPool = sync.Pool{
    New: func() interface{} {
        return &Bot{}
    },
}

func getBotFromPool() *Bot {
    return botPool.Get().(*Bot)
}

func putBotToPool(bot *Bot) {
    bot.Reset()
    botPool.Put(bot)
}
\`\`\`

---

## 📈 性能趋势分析

建议建立性能趋势图，监控以下指标：
- Token记录性能趋势
- 批量操作性能趋势
- 查询性能趋势
- 内存使用趋势

可以使用Grafana或其他监控工具实现。

---

**生成时间**: $(date +"%Y-%m-%d %H:%M:%S")
EOF

    log_success "性能基线对比报告生成完成"
    log_info "报告文件: $report_file"
}

# 提取关键性能指标
extract_key_metrics() {
    log_info "提取关键性能指标..."

    echo ""
    echo "========================================"
    echo "关键性能指标"
    echo "========================================"
    echo ""

    # Token记录性能
    echo "📝 Token记录性能:"
    grep "BenchmarkRecordTokenUsage" "$RESULT_FILE" | head -5 || echo "  无数据"
    echo ""

    # 批量记录性能
    echo "📦 批量记录性能:"
    grep "BenchmarkBatchRecordTokenUsage" "$RESULT_FILE" | head -5 || echo "  无数据"
    echo ""

    # 定价引擎性能
    echo "💰 定价引擎性能:"
    grep "BenchmarkPricingEngine" "$RESULT_FILE" | head -3 || echo "  无数据"
    echo ""

    # 预算检查性能
    echo "🎯 预算检查性能:"
    grep "BenchmarkCheckBudget" "$RESULT_FILE" | head -3 || echo "  无数据"
    echo ""

    # 并发性能
    echo "⚡ 并发性能:"
    grep "BenchmarkConcurrent" "$RESULT_FILE" | head -5 || echo "  无数据"
    echo ""
}

# 主函数
main() {
    print_file_info
    extract_key_metrics
    generate_metrics_table
    generate_baseline_report

    echo ""
    log_success "结果分析完成！"
    echo ""
    echo "生成的文件:"
    echo "  - 基线对比报告: ${RESULT_FILE%.txt}_baseline_report.md"
    echo ""
}

# 执行主函数
main
