# ZKER Pre-commit Hook 配置指南

**版本**: v1.0
**日期**: 2025-01-01
**作者**: DevOps团队
**状态**: ✅ 已实施

---

## 📋 概述

ZKER项目配置了自动化Pre-commit Hooks，确保每次提交的代码符合企业级规范。通过在代码提交前进行自动检查，防止不符合规范的代码进入仓库，有效避免代码退化。

### 🎯 核心目标

1. **API响应格式一致性**: 强制使用统一的API响应格式
2. **错误处理规范性**: 强制使用统一错误码系统
3. **安全问题检查**: 防止引入安全漏洞（panic、SSRF等）
4. **代码质量**: 自动格式化、静态分析、单元测试

---

## 🔧 安装步骤

### 方式一: 自动安装（推荐）

```bash
# 运行安装脚本
bash scripts/install-hooks.sh
```

**输出示例**:
```
🔧 安装ZKER Pre-commit Hooks...
✅ Pre-commit Hooks 安装完成!

📝 钩子已启用，将在每次提交前自动检查:
  - API响应格式一致性
  - 错误处理规范性
  - 安全问题检查
  - 代码格式化
  - go vet检查
  - 单元测试
```

### 方式二: 手动安装

```bash
# 1. 复制hook文件
cp .githooks/pre-commit .git/hooks/pre-commit

# 2. 设置执行权限
chmod +x .git/hooks/pre-commit

# 3. 验证安装
ls -la .git/hooks/pre-commit
```

---

## ✅ 检查项目详解

### 1. API响应格式一致性

**检查目标**: 防止直接使用`c.JSON()`，确保使用统一的响应格式

**违规示例**:
```go
// ❌ 错误: 直接使用c.JSON
c.JSON(http.StatusOK, resp)
c.JSON(http.StatusBadRequest, errResp)
```

**正确示例**:
```go
// ✅ 正确: 使用httputil
httputil.BuildSuccessResp(c, resp)
httputil.BuildErrorResp(c, errno.ErrBadRequestCode, errorx.KV("reason", "msg"))
```

**检查逻辑**:
```bash
# 检查暂存区中是否有c.JSON调用
BAD_JSON=$(git diff --cached | grep -E "^\+.*c\.JSON\(" || true)
```

---

### 2. 错误处理规范性

**检查目标**: 防止使用`errors.New()`和`fmt.Errorf()`，强制使用统一错误码

**违规示例**:
```go
// ❌ 错误: 未使用统一错误码
return errors.New("bot not found")
return fmt.Errorf("failed to create bot: %w", err)
```

**正确示例**:
```go
// ✅ 正确: 使用errorx和errno
return errorx.New(errno.ErrBotNotFoundCode, errorx.KV("bot_id", botID))
return errorx.Wrap(err, errno.ErrBotCreateFailedCode, errorx.KV("reason", "db error"))
```

**检查范围**:
- `backend/domain/*/service/*.go`
- `backend/application/**/*.go`

---

### 3. 安全问题检查

#### 3.1 panic检查

**检查目标**: 防止在生产代码中使用`panic()`

**违规示例**:
```go
// ❌ 错误: 生产代码中使用panic
if bot == nil {
    panic("bot is nil")
}
```

**正确示例**:
```go
// ✅ 正确: 返回错误
if bot == nil {
    return nil, fmt.Errorf("bot not found")
}
```

**例外**: `main.go`中的panic调用被允许（启动失败时）

#### 3.2 SSRF漏洞检查

**检查目标**: 确保所有HTTP请求都有URL验证

**检查文件**:
- `backend/pkg/urltobase64url/parser.go`
- `backend/domain/knowledge/**/*.go`

**验证要求**: 必须存在`validateURL()`函数

---

### 4. 代码格式化（gofmt）

**检查目标**: 确保代码格式统一

**自动修复**:
```bash
gofmt -w backend/
```

**检查逻辑**:
```bash
UNFORMATTED=$(gofmt -l backend/ 2>/dev/null || true)
```

---

### 5. 静态检查（go vet）

**检查目标**: 发现代码中的常见问题

**检查项目**:
- 未使用的变量
- 不可能的类型断言
- 错误的格式化字符串
- 死代码

**手动运行**:
```bash
go vet ./...
```

---

### 6. 单元测试

**检查目标**: 确保修改的代码通过所有测试

**测试范围**: 只运行修改文件相关的测试

**手动运行**:
```bash
# 运行所有短测试
go test -short ./...

# 运行特定包的测试
go test ./domain/tenant/...
```

---

## 🚨 常见问题

### 问题1: Hook被跳过

**场景**: 需要临时跳过检查（紧急修复）

**解决方案**:
```bash
# 使用--no-verify跳过pre-commit检查
git commit --no-verify -m "hotfix: urgent fix"

# 使用--no-verify跳过pre-push检查
git push --no-verify
```

⚠️ **注意**: 仅在紧急情况下使用，常规开发应修复问题而非跳过检查

---

### 问题2: 检查失败 - API响应格式

**错误信息**:
```
❌ 发现直接使用c.JSON的代码!
请使用httputil.BuildSuccessResp()或httutil.BuildErrorResp()
```

**解决方案**:
```go
// 1. 找到违规代码
git diff --cached | grep "c\.JSON"

// 2. 修改为正确格式
// 旧代码
c.JSON(http.StatusOK, data)

// 新代码
httputil.BuildSuccessResp(c, data)

// 3. 重新提交
git add .
git commit -m "fix: use httputil for API response"
```

---

### 问题3: 检查失败 - 错误处理

**错误信息**:
```
❌ 发现未使用统一错误码的错误处理!
请使用errorx.New(errno.ErrXxxCode, errorx.KV(...))
```

**解决方案**:
```go
// 1. 查看errno包中的错误码定义
cat backend/types/errno/errno.go

// 2. 找到合适的错误码
const ErrBotNotFoundCode = 200404 // Bot不存在

// 3. 修改代码
// 旧代码
return errors.New("bot not found")

// 新代码
return errorx.New(errno.ErrBotNotFoundCode)
```

---

### 问题4: 检查失败 - 代码格式

**错误信息**:
```
❌ 发现未格式化的Go文件:
backend/domain/tenant/service.go
```

**解决方案**:
```bash
# 自动格式化
gofmt -w backend/

# 重新提交
git add .
git commit -m "style: format code with gofmt"
```

---

### 问题5: 检查失败 - 单元测试

**错误信息**:
```
❌ 单元测试失败: backend/domain/tenant
```

**解决方案**:
```bash
# 1. 本地运行测试查看详细错误
cd backend/domain/tenant
go test -v ./...

# 2. 修复失败的测试
vim tenant_service_test.go

# 3. 重新运行测试确认通过
go test -v ./...

# 4. 重新提交
cd -
git add .
git commit -m "test: fix failing unit tests"
```

---

## 📊 CI/CD集成

### Pre-commit vs CI/CD

| 维度 | Pre-commit Hook | CI/CD |
|------|----------------|-------|
| **运行时机** | 本地提交前 | 云端Push/PR后 |
| **检查速度** | 快（只检查变更文件） | 慢（检查全量） |
| **阻止提交** | ✅ 是 | ✅ 是（合并前） |
| **覆盖率** | 基础检查 | 完整检查（性能测试、安全扫描） |

### CI/CD工作流

**文件**: `.github/workflows/global-quality-check.yml`

**触发条件**:
- PR到`main`或`develop`分支
- Push到`main`或`develop`分支
- 路径包含`backend/**`

**检查项目**:
1. **代码质量检查**
   - API响应格式
   - 错误处理规范
   - 安全扫描（gosec）
   - 代码质量（golangci-lint）
   - 单元测试
   - 覆盖率上传

2. **安全漏洞扫描**
   - Gosec安全扫描
   - SARIF上传到GitHub
   - SSRF漏洞检查

3. **性能基准测试**（仅PR）
   - K6性能测试
   - 性能基线对比

**查看结果**:
```bash
# GitHub Actions页面
https://github.com/coze-dev/coze-studio/actions

# PR检查状态
https://github.com/coze-dev/coze-studio/pull/123
```

---

## 🔧 配置文件

### .golangci.yml

**位置**: 项目根目录

**作用**: golangci-lint配置，定义启用哪些linter及其规则

**关键配置**:
```yaml
linters:
  enable:
    - gosec          # 安全检查
    - govet          # go vet
    - staticcheck    # 静态分析
    - errcheck       # 检查未处理的错误
    - gosimple       # 简化代码
    - gofumpt        # 代码格式化

linters-settings:
  govet:
    check-shadowing: true  # 检查变量遮蔽
  gocyclo:
    min-complexity: 15     # 圈复杂度阈值

run:
  timeout: 5m
  skip-dirs:
    - vendor
    - api/model  # 自动生成的代码
```

---

## 📝 最佳实践

### 1. 提交前检查清单

```bash
# 1. 手动运行Pre-commit检查
./.githooks/pre-commit

# 2. 运行完整测试套件
cd backend && go test ./...

# 3. 运行架构检查
bash scripts/architecture-check.sh

# 4. 提交代码
git add .
git commit -m "feat(tenant): add tenant isolation"
```

### 2. 提交信息规范

遵循**Conventional Commits**规范：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型**:
- `feat`: 新功能
- `fix`: Bug修复
- `docs`: 文档更新
- `test`: 测试相关
- `refactor`: 重构
- `perf`: 性能优化
- `style`: 代码格式
- `chore`: 杂项
- `ci`: CI/CD相关

**示例**:
```bash
git commit -m "feat(tenant): add tenant_id field to all tables"

git commit -m "fix(permission): resolve role check bug in data permission"

git commit -m "docs(api): update tenant management API documentation"
```

### 3. 分支策略

```
main         ← 生产环境，受保护
  ↑
develop      ← 开发环境，日常开发分支
  ↑
feature/*    ← 功能分支，从develop分出
bugfix/*     ← Bug修复分支
hotfix/*     ← 紧急修复分支
```

---

## 🔍 故障排查

### 问题: Hook无法执行

**症状**:
```
Permission denied: .git/hooks/pre-commit
```

**解决**:
```bash
chmod +x .git/hooks/pre-commit
```

### 问题: Hook路径错误

**症状**:
```
./githooks/pre-commit: No such file or directory
```

**解决**:
```bash
# 检查文件是否存在
ls -la .githooks/pre-commit

# 重新安装
bash scripts/install-hooks.sh
```

### 问题: CI/CD检查通过但本地Hook失败

**原因**: CI/CD环境可能与本地环境不同

**解决**:
```bash
# 同步环境
go mod download
go mod tidy

# 清理缓存
go clean -cache
go clean -testcache
```

---

## 📚 参考文档

- [ZKER企业级开发规范手册 v1.0](./ZKER-企业级开发规范手册_v1.0.md)
- [统一错误码定义规范](./ZKER-统一错误码定义规范.md)
- [全局一致性检查清单](./ZKER-全局一致性检查清单_v1.0.md)
- [API设计规范文档](./API设计规范文档.md)

---

## 🎉 总结

通过Pre-commit Hooks和CI/CD自动化检查，ZKER项目实现了：

✅ **代码质量保障**: 6项自动检查，防止代码退化
✅ **快速反馈**: 本地检查<1分钟，CI/CD<5分钟
✅ **规范执行**: 强制执行开发规范，无需人工审查
✅ **安全防护**: 自动检测安全漏洞和常见问题
✅ **测试覆盖**: 确保每次提交都通过测试

---

**需要帮助？**

- 查看故障排查手册
- 联系DevOps团队
- 提交Issue到GitHub

---

**最后更新**: 2025-01-01
**维护者**: DevOps团队
