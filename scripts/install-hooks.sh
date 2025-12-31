#!/bin/bash
# ZKER Pre-commit Hook 安装脚本

echo "🔧 安装ZKER Pre-commit Hooks..."

# 复制hooks到.git/hooks/
cp -f .githooks/pre-commit .git/hooks/pre-commit
cp -f .githooks/pre-commit .git/hooks/pre-commit.sample 2>/dev/null || true

# 设置执行权限
chmod +x .git/hooks/pre-commit

echo "✅ Pre-commit Hooks 安装完成!"
echo ""
echo "📝 钩子已启用，将在每次提交前自动检查:"
echo "  - API响应格式一致性"
echo "  - 错误处理规范性"
echo "  - 安全问题检查"
echo "  - 代码格式化"
echo "  - go vet检查"
echo "  - 单元测试"
echo ""
