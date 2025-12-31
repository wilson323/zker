# ZKER架构合规检查系统

**版本**: v1.0
**状态**: 生产就绪
**最后更新**: 2025-01-01

---

## 🚀 快速开始

```bash
# 1. 安装Git Hooks（推荐）
bash scripts/install-hooks.sh

# 2. 手动运行架构检查
bash scripts/architecture-check.sh

# 3. 提交代码（自动触发检查）
git commit -m "feat(tenant): add tenant isolation"
```

---

## 📁 文件结构

```
coze-studio/
├── .github/
│   ├── workflows/
│   │   └── architecture-compliance.yml    # CI工作流
│   └── linters/
│       └── .golangci.yml                  # golangci配置
├── .githooks/
│   └── pre-commit                         # Pre-commit钩子
├── scripts/
│   ├── architecture-check.sh              # 架构检查脚本
│   ├── install-hooks.sh                   # Hooks安装脚本
│   └── README-ARCHITECTURE-CHECK.md       # 本文档
└── backend/
    └── ... (Go代码)
```

---

## ✅ 检查项

| 检查项 | 说明 | 阈值 | 检查时机 |
|--------|------|------|---------|
| **编译检查** | 确保代码可以编译 | 必须通过 | 本地 + CI |
| **循环依赖检测** | 检测Go包循环依赖 | 禁止 | 本地 + CI |
| **DDD分层合规检查** | 验证层次依赖关系 | 严格 | 本地 + CI |
| **测试覆盖率检查** | 单元测试覆盖率 | ≥70% | 本地 + CI |
| **代码规范检查** | golangci-lint检查 | 必须通过 | 本地 + CI |
| **安全扫描** | gosec安全漏洞扫描 | 高危禁止 | CI |
| **包大小检查** | 文件行数检查 | ≤3000行 | 本地 + CI |
| **格式检查** | gofmt格式验证 | 必须通过 | 本地 |

---

## 📋 使用指南

### 本地开发

#### 1. Pre-commit检查（自动）

每次 `git commit` 时自动触发。

```bash
# 正常提交
git commit -m "feat(tenant): add tenant_id field"
# ✅ 自动运行pre-commit检查

# 跳过检查（不推荐）
git commit --no-verify -m "message"
```

#### 2. 架构检查脚本（手动）

```bash
# 运行完整检查
bash scripts/architecture-check.sh

# 检查内容：
# 1. 编译检查
# 2. 循环依赖检查
# 3. DDD分层合规检查
# 4. 包大小检查
# 5. 测试覆盖率检查 (≥70%)
# 6. 代码规范检查 (golangci-lint)
# 7. 代码格式检查 (gofmt)
# 8. 安全检查 (go vet)
```

#### 3. 代码格式化

```bash
# 格式化代码
gofmt -w .

# 排序imports
goimports -w .
```

### CI/CD集成

#### GitHub Actions

PR到 `main` 或 `develop` 分支时自动触发。

**查看结果**:
1. 打开PR页面
2. 点击"Checks"标签
3. 查看"Architecture Compliance Check"任务

**触发条件**:
- 推送到 `main/develop` 分支
- 创建或更新PR
- 手动触发（workflow_dispatch）

---

## 🔧 DDD分层规则

```
domain/        ← 领域层（核心业务逻辑）
  ↑            ← 不依赖任何外层
  |
crossdomain/   ← 跨领域共享层
  ↑
bizpkg/        ← 业务包层
  ↑
application/   ← 应用层（用例编排）
  ↑
api/           ← API层（HTTP处理器）
  ↑
infra/         ← 基础设施层（技术实现）
```

**关键规则**:
- ❌ `domain/` 禁止依赖 `api/` 和 `application/`
- ❌ `bizpkg/` 禁止依赖 `api/model/`（应使用 `crossdomain/model/`）
- ❌ `application/` 禁止依赖 `api/`

---

## 🛠️ 常用命令

```bash
# 架构检查
bash scripts/architecture-check.sh

# 安装hooks
bash scripts/install-hooks.sh

# 格式化代码
gofmt -w .
goimports -w .

# 运行linter
golangci-lint run --fix

# 运行测试
go test -v ./...
go test -race ./...
go test -coverprofile=coverage.out ./...

# 查看覆盖率
go tool cover -html=coverage.out

# 跳过检查
git commit --no-verify -m "message"
git push --no-verify
```

---

## 🐛 故障排查

### Pre-commit钩子未触发

```bash
# 重新安装钩子
bash scripts/install-hooks.sh

# 检查权限
ls -la .git/hooks/pre-commit
```

### golangci-lint检查失败

```bash
# 自动修复
golangci-lint run --fix

# 查看详情
golangci-lint run --config=.github/linters/.golangci.yml
```

### 测试覆盖率不足

```bash
# 查看覆盖率详情
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 📚 相关文档

- [ZKER-CI-CD架构合规检查系统使用指南 v1.0](../docs/企业级功能完善与统一性设计方案/ZKER-CI-CD架构合规检查系统使用指南_v1.0.md) - 完整使用指南
- [ZKER企业级开发规范手册 v1.0](../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md) - 开发规范
- [ZKER全局一致性检查清单 v1.0](../docs/企业级功能完善与统一性设计方案/ZKER-全局一致性检查清单_v1.0.md) - 检查清单

---

## 🎯 最佳实践

1. **开发流程**:
   ```bash
   # 1. 本地检查
   bash scripts/architecture-check.sh

   # 2. 提交代码
   git commit -m "feat(tenant): add tenant isolation"

   # 3. 推送代码
   git push origin feature/add-tenant-isolation

   # 4. 创建PR
   # GitHub Actions自动运行检查

   # 5. 等待CI通过
   ```

2. **编写测试**: 目标覆盖率≥70%

3. **遵循DDD分层**: 严格遵守层次依赖规则

4. **统一错误码**: 使用 `types/errno` 包

5. **提交信息**: 遵循Conventional Commits规范

---

**维护者**: ZKER开发团队
**反馈**: GitHub Issues
