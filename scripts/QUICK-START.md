# 🚀 CI/CD架构合规检查系统 - 快速开始

## 一分钟快速开始

```bash
# 1. 安装Git Hooks（推荐）
bash scripts/install-hooks.sh

# 2. 验证安装
bash scripts/verify-architecture-check.sh

# 3. 测试pre-commit钩子
git commit --allow-empty -m "test: verify architecture check"
```

## ✅ 系统已就绪！

所有组件已成功部署：

| 组件 | 状态 | 文件 |
|------|------|------|
| **GitHub Actions** | ✅ | `.github/workflows/architecture-compliance.yml` |
| **golangci配置** | ✅ | `.github/linters/.golangci.yml` |
| **Pre-commit钩子** | ✅ | `.githooks/pre-commit` |
| **架构检查脚本** | ✅ | `scripts/architecture-check.sh` |
| **Hooks安装脚本** | ✅ | `scripts/install-hooks.sh` |
| **验证脚本** | ✅ | `scripts/verify-architecture-check.sh` |
| **使用文档** | ✅ | `scripts/README-ARCHITECTURE-CHECK.md` |
| **完整指南** | ✅ | `docs/企业级功能完善与统一性设计方案/ZKER-CI-CD架构合规检查系统使用指南_v1.0.md` |

## 📋 8项检查

1. ✅ 编译检查
2. ✅ 循环依赖检测
3. ✅ DDD分层合规检查
4. ✅ 包大小检查
5. ✅ 测试覆盖率检查 (≥70%)
6. ✅ 代码规范检查 (golangci-lint)
7. ✅ 代码格式检查 (gofmt)
8. ✅ 安全检查 (go vet)

## 🔧 常用命令

```bash
# 安装hooks
bash scripts/install-hooks.sh

# 运行完整检查
bash scripts/architecture-check.sh

# 验证安装
bash scripts/verify-architecture-check.sh

# 格式化代码
gofmt -w .
goimports -w .

# 运行linter
golangci-lint run --fix

# 运行测试
go test -v ./...
go test -coverprofile=coverage.out ./...

# 查看覆盖率
go tool cover -html=coverage.out
```

## 📚 文档

- **快速参考**: `scripts/README-ARCHITECTURE-CHECK.md`
- **完整指南**: `docs/企业级功能完善与统一性设计方案/ZKER-CI-CD架构合规检查系统使用指南_v1.0.md`
- **部署总结**: `scripts/ARCHITECTURE-CHECK-SUMMARY.md`

## 🎯 下一步

1. ✅ 安装Git hooks: `bash scripts/install-hooks.sh`
2. ✅ 运行验证: `bash scripts/verify-architecture-check.sh`
3. ✅ 测试检查: `bash scripts/architecture-check.sh`
4. ✅ 提交代码: `git commit -m "feat(tenant): add feature"`（自动触发pre-commit）

## 🎉 完成！

系统已就绪，开始使用吧！

---

**文档**: `docs/企业级功能完善与统一性设计方案/ZKER-CI-CD架构合规检查系统使用指南_v1.0.md`
**维护**: ZKER开发团队
