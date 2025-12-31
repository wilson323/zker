# 后端编译和errno检查 CI/CD 集成指南

> **版本**: v1.0
> **更新**: 2025-01-03
> **维护者**: ZKER架构团队

---

## 📋 目录

- [概述](#概述)
- [文件清单](#文件清单)
- [工作流说明](#工作流说明)
- [errno验证脚本](#errno验证脚本)
- [本地测试](#本地测试)
- [CI/CD集成](#cicd集成)
- [故障排查](#故障排查)

---

## 概述

本文档描述了 ZKER 项目中新增的**后端编译检查**和**errno使用规范验证** CI/CD 工作流。

### 目标

- ✅ 在每次代码提交时自动检查后端代码编译
- ✅ 验证 errno 使用符合企业级规范
- ✅ 在 PR 阶段拦截不符合规范的代码
- ✅ 提供清晰的错误报告和修复建议

### 检查内容

| 检查项 | 说明 | 优先级 |
|--------|------|--------|
| **编译检查** | 确保所有Go代码可以编译 | P0 |
| **errno规范** | 验证errorx.New/Wrap使用 | P0 |
| **代码格式** | gofmt格式检查 | P1 |
| **静态分析** | go vet静态分析 | P1 |
| **依赖安全** | gosec安全扫描 | P1 |

---

## 文件清单

### 新增文件

```
.github/workflows/
└── backend-compilation-check.yml    # CI/CD工作流定义

scripts/
├── verify-errno-usage.sh           # errno验证脚本
└── test-compilation-check.sh       # 本地测试脚本

docs/企业级功能完善与统一性设计方案/
└── BACKEND_COMPILATION_ERRNO_CHECK_GUIDE.md  # 本文档
```

### 文件说明

#### 1. backend-compilation-check.yml

GitHub Actions工作流，包含以下任务：

- **compilation-check**: 后端代码编译检查
- **errno-usage-check**: errno使用规范验证
- **code-quality-check**: 代码质量检查（gofmt, go vet）
- **dependency-check**: 依赖安全检查（gosec）
- **summary**: 综合总结报告

#### 2. verify-errno-usage.sh

errno使用规范验证脚本，执行以下检查：

1. **直接errno引用检查**: 检查是否有 `return errno.ErrXXX`（违规）
2. **errorx使用检查**: 统计 `errorx.New(errno.ErrXXX)` 使用
3. **错误处理覆盖率**: 计算错误处理覆盖率（目标≥95%）
4. **硬编码错误检查**: 检查硬编码错误字符串
5. **errno导入检查**: 验证errno/errorx包导入
6. **错误传播检查**: 检查裸错误传播

#### 3. test-compilation-check.sh

本地快速测试脚本，用于在提交前验证检查项。

---

## 工作流说明

### 触发条件

工作流在以下情况自动触发：

- **Push**: 推送到任何分支（main, develop, feature/*, bugfix/*）
- **Pull Request**: 创建或更新 PR
- **手动触发**: 在 GitHub Actions 页面手动运行

### 执行流程

```
触发工作流
    │
    ├─→ 1. 编译检查
    │   └─→ go build ./...
    │
    ├─→ 2. errno验证（依赖编译检查）
    │   └─→ scripts/verify-errno-usage.sh
    │
    ├─→ 3. 代码质量检查（依赖编译检查）
    │   ├─→ gofmt
    │   └─→ go vet
    │
    ├─→ 4. 依赖安全检查
    │   ├─→ go mod verify
    │   └─→ gosec
    │
    └─→ 5. 总结报告
        ├─→ 生成 GitHub Summary
        └─→ PR评论（失败时）
```

### 失败处理

当检查失败时：

1. **GitHub Summary**: 显示详细的错误信息
2. **PR 评论**: 自动在 PR 中添加失败评论
3. **日志文件**: 保存完整的检查日志
4. **阻止合并**: PR 无法合并（如果设置branch protection）

---

## errno验证脚本

### 使用方法

```bash
# 直接运行
bash scripts/verify-errno-usage.sh

# 或使用可执行权限
chmod +x scripts/verify-errno-usage.sh
./scripts/verify-errno-usage.sh
```

### 输出说明

#### 成功输出

```
==========================================
  errno使用规范验证
==========================================

开始时间: 2025-01-03 10:30:00

==========================================
  检查1: 直接errno引用检查
==========================================

ℹ 扫描直接errno引用（排除测试文件）...
✓ 未发现直接errno引用（0处违规）

==========================================
  检查2: errorx.New/Wrap使用检查
==========================================

ℹ 统计errorx使用情况...
  errorx.New使用: 375 处
  errorx.Wrap使用: 42 处
  errorx总使用: 417 处

✓ 发现errorx使用（符合规范）

...

==========================================
✅ errno验证通过
==========================================

所有errno使用符合企业级规范！

详细日志: scripts/errno/errno_check.log
结果文件: scripts/errno/errno_check_result.json
```

#### 失败输出

```
==========================================
❌ errno验证失败
==========================================

发现 25 个错误需要修复

详细日志: scripts/errno/errno_check.log
结果文件: scripts/errno/errno_check_result.json

修复建议:
  1. 使用 errorx.New(errno.ErrXXX) 替换直接 errno 引用
  2. 使用 errorx.Wrap(err, errno.ErrXXX) 添加错误上下文
  3. 运行 ./scripts/errno/phase4_validate_all.sh 进行完整验证
```

### 检查项详解

#### 检查1: 直接errno引用检查

**违规示例**:
```go
// ❌ 错误
if err != nil {
    return nil, errno.ErrBotNotFound
}
```

**正确示例**:
```go
// ✅ 正确
if err != nil {
    return nil, errorx.New(errno.ErrBotNotFound)
}
```

#### 检查2: errorx.New/Wrap使用检查

统计 `errorx.New(errno.ErrXXX)` 和 `errorx.Wrap(err, errno.ErrXXX)` 的使用数量。

**目标**: 使用数量 > 0（表示已完成errno迁移）

#### 检查3: 错误处理覆盖率

计算错误处理覆盖率：

```
覆盖率 = (使用errorx的return语句数 / 包含错误的return语句数) × 100%
```

**目标**: ≥ 95%

#### 检查4: 硬编码错误检查

检查常见的硬编码错误模式：

- `errors.New("...")`
- `fmt.Errorf("...")`
- "invalid", "not found", "already exists" 等

**建议**: 使用 errno 包中的错误定义

#### 检查5: errno导入检查

验证以下导入：

```go
import (
    "coze-studio/backend/types/errno"
    "coze-studio/backend/pkg/errorx"
)
```

#### 检查6: 错误传播检查

检查裸错误传播：

```go
// ⚠️ 不推荐
return nil, err

// ✅ 推荐
return nil, errorx.Wrap(err, errno.ErrBotNotFound)
```

---

## 本地测试

### 快速测试

在提交代码前，运行本地测试脚本：

```bash
bash scripts/test-compilation-check.sh
```

### 手动测试

#### 1. 编译检查

```bash
cd backend
go mod download
go build ./...
```

#### 2. 代码格式检查

```bash
cd backend
gofmt -s -l .
```

#### 3. 静态分析

```bash
cd backend
go vet ./...
```

#### 4. errno验证

```bash
bash scripts/verify-errno-usage.sh
```

### 预提交检查

建议添加到 pre-commit hook：

```bash
# .githooks/pre-commit
# 编译检查
echo "检查后端编译..."
cd backend && go build ./... || {
    echo "❌ 后端编译失败"
    exit 1
}

# errno验证
echo "检查errno使用..."
bash scripts/verify-errno-usage.sh || {
    echo "❌ errno验证失败"
    exit 1
}
```

---

## CI/CD集成

### 配置Secrets

如果需要失败通知，配置以下 Secrets：

| Secret | 说明 | 示例 |
|--------|------|------|
| `SLACK_WEBHOOK` | Slack Webhook URL | `https://hooks.slack.com/...` |
| `EMAIL_TO` | 收件人邮箱 | `team@example.com` |
| `EMAIL_FROM` | 发件人邮箱 | `ci@example.com` |
| `EMAIL_PASSWORD` | 邮箱密码 | `password_or_app_token` |

### Branch Protection

建议设置分支保护规则：

1. 进入 GitHub Settings → Branches
2. 添加规则（如 `main`, `develop`）
3. 勾选以下选项：
   - ✅ Require status checks to pass
   - ✅ Require branches to be up to date
   - 选择必选检查：
     - `Backend Compilation & Errno Check / compilation-check`
     - `Backend Compilation & Errno Check / errno-usage-check`

### 工作流徽章

在 README.md 中添加状态徽章：

```markdown
![Backend Check](https://github.com/coze-dev/coze-studio/actions/workflows/backend-compilation-check.yml/badge.svg)
```

---

## 故障排查

### 问题1: 编译检查失败

**症状**:

```
❌ 编译检查失败
# backend/domain/agent/service.go:123: undefined: errorx
```

**原因**:

缺少 `errorx` 包或导入路径错误

**解决**:

1. 检查 `backend/go.mod` 中是否有 `errorx` 依赖
2. 检查导入路径是否正确
3. 运行 `go mod tidy` 更新依赖

### 问题2: errno验证失败

**症状**:

```
❌ 发现 25 处直接errno引用
```

**原因**:

代码中存在直接返回 `errno.ErrXXX` 而未使用 `errorx.New()` 包装

**解决**:

1. 查看详细日志：`cat scripts/errno/errno_check.log`
2. 逐个修复违规代码
3. 使用 errno 迁移脚本：`./scripts/errno/phase3_migrate_batch.sh`

### 问题3: 工作流未触发

**症状**:

推送代码后，GitHub Actions 未运行工作流

**原因**:

1. 工作流文件路径不正确
2. 工作流文件语法错误
3. GitHub Actions 未启用

**解决**:

1. 检查文件路径：`.github/workflows/backend-compilation-check.yml`
2. 验证 YAML 语法：使用在线工具
3. 检查 Actions 设置：Settings → Actions → General

### 问题4: false positives

**症状**:

errno验证报告了错误，但代码实际上是正确的

**原因**:

1. 注释中的代码被误判
2. 生成的代码（如 `*_gen.go`）未排除
3. 测试代码未正确排除

**解决**:

1. 在违规代码行添加注释：`// errno:direct`
2. 将生成代码移到 `gen/` 目录
3. 确保测试文件以 `_test.go` 结尾

---

## 最佳实践

### 开发流程

1. **编写代码**: 遵循 errno 使用规范
2. **本地测试**: 运行 `bash scripts/test-compilation-check.sh`
3. **提交代码**: `git commit -m "feat: ..."`
4. **推送代码**: `git push`
5. **查看CI**: GitHub Actions 页面查看结果
6. **修复问题**: 如失败，查看日志并修复

### errno使用规范

#### ✅ 推荐做法

```go
// 1. 创建新错误
if err != nil {
    return nil, errorx.New(errno.ErrBotNotFound)
}

// 2. 包装错误（添加上下文）
if err != nil {
    return nil, errorx.Wrap(err, errno.ErrDatabaseFailed)
}

// 3. 错误判断
if errors.Is(err, errno.ErrBotNotFound) {
    // 处理Bot未找到
}
```

#### ❌ 避免做法

```go
// 1. 直接返回errno
if err != nil {
    return nil, errno.ErrBotNotFound  // ❌
}

// 2. 硬编码错误
if err != nil {
    return nil, errors.New("bot not found")  // ❌
}

// 3. 裸错误传播
if err != nil {
    return nil, err  // ⚠️ 不推荐
}
```

---

## 相关文档

- [errno迁移方案](./ZKER-企业级错误码系统全局一致性增强_v1.0.md)
- [errno迁移脚本使用指南](../../scripts/errno/README.md)
- [代码质量改进设计文档](./代码质量改进设计文档_v1.0.md)
- [CI/CD架构合规检查系统使用指南](./ZKER-CI-CD架构合规检查系统使用指南_v1.0.md)

---

## 更新日志

| 版本 | 日期 | 变更内容 |
|------|------|---------|
| v1.0 | 2025-01-03 | 初始版本，添加编译检查和errno验证 |

---

## 支持

如有问题，请联系：

- **技术负责人**: ZKER架构团队
- **文档维护**: DevOps团队
- **问题反馈**: [GitHub Issues](https://github.com/coze-dev/coze-studio/issues)

---

**文档版本**: v1.0
**最后更新**: 2025-01-03
**文档状态**: ✅ 已发布
