# ZKER 自动化质量检查系统

## 📋 概述

ZKER项目配置了完整的自动化质量检查系统，包括本地Pre-commit Hooks和云端CI/CD工作流，确保代码质量和全局一致性。

## 🎯 检查覆盖

### 本地检查 (Pre-commit Hooks)
- ✅ API响应格式一致性
- ✅ 错误处理规范性
- ✅ 安全问题检查
- ✅ 代码格式化 (gofmt)
- ✅ 静态检查 (go vet)
- ✅ 单元测试

### 云端检查 (CI/CD)
- ✅ 完整代码质量检查 (golangci-lint)
- ✅ 安全漏洞扫描 (gosec)
- ✅ 性能基准测试 (K6)
- ✅ 覆盖率报告 (codecov)

## 🚀 快速开始

### 1. 安装本地Hooks

```bash
bash scripts/install-hooks.sh
```

### 2. 验证安装

```bash
# 手动运行检查
./.githooks/pre-commit

# 查看已安装的hooks
ls -la .git/hooks/
```

### 3. 正常提交

```bash
git add .
git commit -m "feat(tenant): add tenant isolation"
# 自动运行pre-commit检查
```

## 📚 文档

- [完整配置指南](../docs/企业级功能完善与统一性设计方案/ZKER-Pre-commit-Hook配置指南_v1.0.md)
- [快速参考卡](../docs/企业级功能完善与统一性设计方案/ZKER-Pre-commit快速参考卡.md)
- [企业级开发规范手册](../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)

## 🔧 配置文件

- `.githooks/pre-commit` - Pre-commit hook脚本
- `.github/workflows/global-quality-check.yml` - CI/CD工作流
- `.golangci.yml` - golangci-lint配置

## 🆘 获取帮助

遇到问题？查看：
1. [常见问题](../docs/企业级功能完善与统一性设计方案/ZKER-Pre-commit-Hook配置指南_v1.0.md#-常见问题)
2. [故障排查](../docs/企业级功能完善与统一性设计方案/ZKER-Pre-commit-Hook配置指南_v1.0.md#-故障排查)
3. [DevOps团队](mailto:devops@zker.com)
