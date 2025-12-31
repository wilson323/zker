#!/bin/bash
# Pre-commit Hook: 错误处理规范性检查
#
# 安装方法：
#   cp scripts/pre-commit-error-check.sh .git/hooks/pre-commit
#   chmod +x .git/hooks/pre-commit
#
# 功能：检查Go代码的错误处理规范性，阻止不符合规范的代码提交

set -e

echo "🔍 检查错误处理规范性..."

# 获取暂存的Go文件
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

if [ -z "$STAGED_GO_FILES" ]; then
    echo "⚠️  没有暂存的Go文件"
    exit 0
fi

echo "📄 检查以下文件："
echo "$STAGED_GO_FILES"
echo ""

# 临时保存当前工作目录
CURRENT_DIR=$(pwd)

# 创建临时目录存放待检查的文件
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# 复制暂存的文件到临时目录
for FILE in $STAGED_GO_FILES; do
    mkdir -p "$TEMP_DIR/$(dirname $FILE)"
    git show ":$FILE" > "$TEMP_DIR/$FILE"
done

# 切换到临时目录进行扫描
cd $TEMP_DIR

# 运行错误处理扫描器
python $CURRENT_DIR/tools/error_handler_migrator.py backend 2>&1 | tee /tmp/error_check_report.txt

# 提取违规数量
VIOLATIONS=$(grep "违规总数" /tmp/error_check_report.txt | awk '{print $3}' || echo "0")

echo ""
echo "========================================="
echo "📊 扫描结果"
echo "========================================="
echo "违规总数: $VIOLATIONS 处"

if [ "$VIOLATIONS" -gt "0" ]; then
    echo ""
    echo "❌ 发现 $VIOLATIONS 处错误处理违规！"
    echo ""
    echo "违规类型："
    grep -A 10 "违规类型统计" /tmp/error_check_report.txt | tail -5
    echo ""
    echo "请使用以下工具修复："
    echo "  - python tools/error_handler_migrator.py backend/domain/<模块名>"
    echo "  - 参考: docs/企业级功能完善与统一性设计方案/ZKER-P0模块错误处理标准化完成报告_v1.0.md"
    echo ""
    exit 1
else
    echo ""
    echo "✅ 错误处理检查通过！"
    echo ""
    exit 0
fi
