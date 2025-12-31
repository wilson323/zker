# ZKER CI/CD架构合规检查系统使用指南

**版本**: v1.0
**创建时间**: 2025-01-01
**适用阶段**: 企业级功能完善阶段
**遵循规范**: ZKER企业级开发规范手册 v1.0

---

## 📖 目录

1. [概述](#概述)
2. [功能特性](#功能特性)
3. [快速开始](#快速开始)
4. [本地使用](#本地使用)
5. [CI/CD集成](#cicd集成)
6. [检查项详解](#检查项详解)
7. [故障排查](#故障排查)
8. [最佳实践](#最佳实践)
9. [附录](#附录)

---

## 概述

ZKER CI/CD架构合规检查系统是一套完整的自动化质量门禁，用于确保代码符合企业级架构标准。系统通过Git hooks、本地脚本和GitHub Actions三层检查，在代码提交前、提交后和PR阶段全面验证代码质量。

### 核心目标

- ✅ **自动化优先**: 所有检查自动化执行，无需人工干预
- ✅ **快速反馈**: 本地检查<2分钟，CI检查<5分钟
- ✅ **清晰输出**: 明确指出问题位置和修复建议
- ✅ **可配置性**: 通过配置文件灵活调整规则
- ✅ **CI/CD集成**: 无缝集成到GitHub Actions流水线

---

## 功能特性

### 1. 架构合规性检查

| 检查项 | 说明 | 阈值 |
|--------|------|------|
| **编译检查** | 确保代码可以正常编译 | 必须通过 |
| **循环依赖检测** | 检测Go包循环依赖 | 禁止 |
| **DDD分层合规检查** | 验证层次依赖关系 | 严格 |
| **测试覆盖率检查** | 单元测试覆盖率 | ≥70% |
| **代码规范检查** | golangci-lint检查 | 必须通过 |
| **安全扫描** | gosec安全漏洞扫描 | 高危禁止 |
| **包大小检查** | 文件行数检查 | ≤3000行 |
| **格式检查** | gofmt格式验证 | 必须通过 |

### 2. DDD分层依赖规则

```
domain/        ← 领域层，核心业务逻辑
  ↑            ← 不依赖任何外层
  |
crossdomain/   ← 跨领域共享层
  ↑
bizpkg/        ← 业务包层
  ↑
application/   ← 应用层，用例编排
  ↑
api/           ← API层，HTTP处理器
  ↑
infra/         ← 基础设施层，技术实现
```

**关键规则**:
- ❌ `domain/` 禁止依赖 `api/` 和 `application/`
- ❌ `bizpkg/` 禁止依赖 `api/model/`（应使用 `crossdomain/model/`）
- ❌ `application/` 禁止依赖 `api/`

---

## 快速开始

### 安装Git Hooks（推荐）

```bash
# 1. 进入项目根目录
cd D:/code/coze-studio

# 2. 运行安装脚本
bash scripts/install-hooks.sh

# 3. 验证安装
ls -la .git/hooks/
```

安装成功后，每次 `git commit` 都会自动触发pre-commit检查。

### 手动运行架构检查

```bash
# 运行完整架构检查
bash scripts/architecture-check.sh
```

---

## 本地使用

### 1. Pre-commit检查（自动）

每次 `git commit` 时自动触发，检查时间约30-60秒。

**检查内容**:
- ✅ 代码格式（gofmt）
- ✅ Import排序（goimports）
- ✅ 静态检查（go vet）
- ✅ 循环依赖检查
- ✅ DDD分层合规检查
- ✅ 单元测试（短测试）

**示例输出**:

```bash
$ git commit -m "feat(tenant): add tenant_id field"

🔍 运行Pre-commit架构合规性检查
========================================

1️⃣  代码格式检查 (gofmt)...
✅ 代码格式检查通过

2️⃣  Import排序检查 (goimports)...
✅ Import排序检查通过

3️⃣  静态检查 (go vet)...
✅ 静态检查通过

4️⃣  循环依赖检查...
✅ 无循环依赖

5️⃣  DDD分层合规检查...
✅ domain层合规
✅ bizpkg层合规

6️⃣  运行单元测试 (短测试)...
✅ 单元测试通过

========================================
✅ 所有Pre-commit检查通过

🎉 可以安全提交代码！
```

### 2. 架构检查脚本（手动）

运行完整的架构检查套件，检查时间约2-5分钟。

```bash
bash scripts/architecture-check.sh
```

**检查内容**:
1. 编译检查
2. 循环依赖检查
3. DDD分层合规检查
4. 包大小检查
5. 测试覆盖率检查（≥70%）
6. 代码规范检查（golangci-lint）
7. 代码格式检查
8. 安全检查（go vet）

### 3. 跳过检查（不推荐）

**跳过pre-commit检查**:
```bash
git commit --no-verify -m "message"
```

**跳过pre-push检查**:
```bash
git push --no-verify
```

⚠️ **警告**: 跳过检查可能导致代码质量下降，仅在紧急情况下使用。

---

## CI/CD集成

### GitHub Actions工作流

架构合规检查会在以下情况自动触发：
- 推送代码到 `main` 或 `develop` 分支
- 创建或更新Pull Request到 `main` 或 `develop`
- 手动触发（workflow_dispatch）

### 工作流文件

**文件位置**: `.github/workflows/architecture-compliance.yml`

**检查流程**:

```yaml
1️⃣  编译检查          → 必须通过
2️⃣  循环依赖检测       → 必须通过
3️⃣  跨层依赖检测       → 必须通过
4️⃣  测试覆盖率检查     → ≥70%
5️⃣  代码规范检查       → 必须通过
6️⃣  安全扫描          → 高危禁止
7️⃣  包大小检查        → 警告
8️⃣  生成检查报告      → 始终执行
```

### 查看CI结果

**GitHub界面**:
1. 进入Pull Request页面
2. 点击"Checks"标签
3. 查看"Architecture Compliance Check"任务

**命令行**:
```bash
# 查看CI状态
gh run list --workflow=architecture-compliance.yml

# 查看特定运行结果
gh run view <run-id>
```

### CI检查失败处理

1. **查看失败详情**:
   - 点击失败的检查步骤
   - 查看详细日志和错误信息

2. **本地复现问题**:
   ```bash
   bash scripts/architecture-check.sh
   ```

3. **修复问题**:
   - 根据错误信息修复代码
   - 运行本地检查验证修复
   - 提交修复并推送

4. **重新触发CI**:
   - CI会自动重新运行
   - 或手动触发：`gh workflow run architecture-compliance.yml`

---

## 检查项详解

### 1. 编译检查

**目的**: 确保代码可以正常编译

**命令**:
```bash
go build -v ./...
```

**失败原因**:
- 语法错误
- 类型错误
- 缺少依赖

**修复方法**:
```bash
# 查看编译错误详情
go build -v ./...

# 更新依赖
go mod tidy

# 修复语法或类型错误
```

### 2. 循环依赖检测

**目的**: 检测Go包循环依赖

**命令**:
```bash
go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./...
```

**示例错误**:
```
❌ 发现循环依赖:
github.com/coze-dev/coze-studio/backend/domain/tenant
    imports github.com/coze-dev/coze-studio/backend/api/model
    imports github.com/coze-dev/coze-studio/backend/domain/tenant
```

**修复方法**:
- 提取共享代码到独立的包
- 使用依赖注入解除循环依赖
- 重新设计包结构

### 3. DDD分层合规检查

**目的**: 确保符合DDD分层架构

**检查规则**:

| 层次 | 禁止依赖 | 允许依赖 |
|------|---------|---------|
| domain/ | api/, application/ | 无（核心层） |
| crossdomain/ | api/, application/ | domain/ |
| bizpkg/ | api/model/ | domain/, crossdomain/ |
| application/ | api/ | domain/, bizpkg/, crossdomain/, infra/ |
| api/ | - | domain/, application/, bizpkg/, infra/ |
| infra/ | - | 所有层 |

**示例错误**:
```
❌ domain层依赖外层:
domain/tenant/service.go: api/model/common.go
```

**修复方法**:
```go
// ❌ 错误：domain层依赖api/model
import "github.com/coze-dev/coze-studio/backend/api/model"

// ✅ 正确：使用crossdomain/model
import "github.com/coze-dev/coze-studio/backend/crossdomain/model"
```

### 4. 测试覆盖率检查

**目的**: 确保代码测试覆盖率≥70%

**命令**:
```bash
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out | grep total
```

**覆盖率阈值**: 70%

**查看覆盖率详情**:
```bash
# 生成HTML报告
go tool cover -html=coverage.out -o coverage.html

# 按包查看覆盖率
go tool cover -func=coverage.out
```

**提高覆盖率**:
```bash
# 查看未覆盖的代码
go tool cover -html=coverage.out

# 添加测试用例
# 在*_test.go文件中添加测试函数
```

### 5. 代码规范检查（golangci-lint）

**目的**: 统一代码风格和最佳实践

**配置文件**: `.github/linters/.golangci.yml`

**常用命令**:
```bash
# 运行检查
golangci-lint run --config=.github/linters/.golangci.yml

# 自动修复
golangci-lint run --fix --config=.github/linters/.golangci.yml

# 仅运行特定linter
golangci-lint run --disable-all --enable=gofmt,goimports
```

**常见问题**:

| 问题 | 说明 | 修复方法 |
|------|------|---------|
| `expected ... have ...` | 格式错误 | 运行 `gofmt -w file.go` |
| `SA9003: ...` | 未使用的分支 | 删除或修复代码 |
| `G101: ...` | 硬编码凭证 | 使用环境变量 |
| `G401: ...` | 弱加密算法 | 使用强加密算法 |

### 6. 安全扫描（gosec）

**目的**: 检测安全漏洞

**命令**:
```bash
gosec -quiet -fmt json ./...
```

**常见安全问题**:

| 编号 | 说明 | 修复方法 |
|------|------|---------|
| G101 | 硬编码凭证 | 使用环境变量或配置文件 |
| G201 | SQL注入 | 使用参数化查询 |
| G204 | 命令注入 | 避免直接拼接命令 |
| G301/302 | 文件权限 | 使用正确的文件权限 |
| G401 | 弱加密算法 | 使用AES/RSA等强算法 |

### 7. 包大小检查

**目的**: 保持文件简洁，便于维护

**阈值**: 3000行

**检查命令**:
```bash
find . -name "*.go" -not -path "*/vendor/*" | xargs wc -l | awk '$1 > 3000'
```

**建议**:
- 拆分大型文件
- 按功能组织代码
- 保持单一职责原则

### 8. 格式检查（gofmt）

**目的**: 统一代码格式

**命令**:
```bash
gofmt -l . | grep -v vendor
```

**自动修复**:
```bash
gofmt -w .
```

---

## 故障排查

### 问题1: pre-commit钩子未触发

**症状**: `git commit` 没有运行检查

**原因**:
1. 钩子未安装
2. 钩子没有执行权限
3. 使用了 `--no-verify` 参数

**解决方法**:
```bash
# 重新安装钩子
bash scripts/install-hooks.sh

# 检查权限
ls -la .git/hooks/pre-commit

# 手动触发钩子
.git/hooks/pre-commit
```

### 问题2: golangci-lint检查失败

**症状**: 代码规范检查失败，但不确定如何修复

**解决方法**:
```bash
# 1. 查看详细错误
golangci-lint run --config=.github/linters/.golangci.yml

# 2. 尝试自动修复
golangci-lint run --fix --config=.github/linters/.golangci.yml

# 3. 查看特定linter的帮助
golangci-lint run --disable-all --enable=<linter-name> --help
```

### 问题3: 测试覆盖率不足

**症状**: 覆盖率低于70%

**解决方法**:
```bash
# 1. 查看覆盖率详情
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -html=coverage.out

# 2. 识别未覆盖的代码
# 在浏览器中打开coverage.html

# 3. 添加测试用例
# 在*_test.go文件中添加测试函数
```

### 问题4: CI检查超时

**症状**: GitHub Actions检查超时失败

**原因**:
- 测试执行时间过长
- 网络问题导致依赖下载慢
- 资源限制

**解决方法**:
```yaml
# 增加timeout
timeout-minutes: 30

# 使用缓存
- uses: actions/setup-go@v5
  with:
    cache: true

# 跳过慢速测试
go test -short ./...
```

### 问题5: 循环依赖误报

**症状**: 确定没有循环依赖，但检查失败

**原因**:
- vendor目录干扰
- 生成的代码问题

**解决方法**:
```bash
# 清理vendor
go mod vendor
rm -rf vendor/

# 重新生成
go mod tidy

# 检查生成的代码
# 在.golangci.yml中排除
exclude-files:
  - ".*\\.gen\\.go"
```

---

## 最佳实践

### 1. 开发工作流

**推荐流程**:
```bash
# 1. 创建功能分支
git checkout -b feature/add-tenant-isolation

# 2. 开发代码
# ... 编写代码 ...

# 3. 本地检查
bash scripts/architecture-check.sh

# 4. 提交代码（自动触发pre-commit）
git add .
git commit -m "feat(tenant): add tenant isolation middleware"

# 5. 推送代码（自动触发pre-push）
git push origin feature/add-tenant-isolation

# 6. 创建PR
# GitHub Actions自动运行完整检查

# 7. 等待CI通过
# 修复失败的问题

# 8. 合并代码
```

### 2. 编写测试

**目标**: 覆盖率≥70%

**示例**:
```go
// service_test.go
package tenant

import (
    "testing"
)

func TestTenantService_Create(t *testing.T) {
    tests := []struct {
        name    string
        req     *CreateTenantRequest
        want    *Tenant
        wantErr bool
    }{
        {
            name: "正常创建",
            req:  &CreateTenantRequest{Name: "test"},
            want: &Tenant{ID: "1", Name: "test"},
            wantErr: false,
        },
        {
            name: "名称为空",
            req:  &CreateTenantRequest{Name: ""},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := service.Create(tt.req)
            if (err != nil) != tt.wantErr {
                t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Create() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 3. 代码组织

**DDD分层**:
```
backend/domain/tenant/
├── entity/              # 实体
│   └── tenant.go
├── repository/          # 仓储接口
│   └── tenant_repository.go
├── service/             # 领域服务
│   └── tenant_service.go
└── internal/            # 内部实现（不对外暴露）
    ├── dal/             # 数据访问层
    └── convert/         # 数据转换
```

**依赖规则**:
- ✅ `domain/entity` 不依赖任何外层
- ✅ `domain/repository` 仅定义接口
- ✅ `domain/service` 实现业务逻辑
- ❌ `domain/` 不依赖 `api/` 和 `application/`

### 4. 错误处理

**统一错误码**:
```go
import "github.com/coze-dev/coze-studio/backend/types/errno"

// 使用统一错误码
if req.TenantID == "" {
    return nil, errno.ErrTenantNotFound
}

// 包装错误
if err != nil {
    return nil, fmt.Errorf("failed to create tenant: %w", err)
}
```

### 5. 提交信息

**Conventional Commits**:
```bash
# 格式: <type>(<scope>): <subject>

# type: feat, fix, docs, style, refactor, perf, test, chore, ci, build
# scope: 模块名（tenant, permission, routing等）
# subject: 简短描述（≤50字符）

# 示例
feat(tenant): add tenant_id field to all tables
fix(permission): resolve role check bug in data permission
docs(api): update tenant management API documentation
refactor(routing): simplify routing engine logic
test(knowledge): add unit tests for knowledge service
```

---

## 附录

### A. 配置文件

| 文件 | 位置 | 说明 |
|------|------|------|
| **golangci配置** | `.github/linters/.golangci.yml` | golangci-lint规则配置 |
| **pre-commit钩子** | `.githooks/pre-commit` | 提交前检查脚本 |
| **架构检查脚本** | `scripts/architecture-check.sh` | 完整架构检查 |
| **hooks安装脚本** | `scripts/install-hooks.sh` | 安装Git hooks |
| **CI工作流** | `.github/workflows/architecture-compliance.yml` | CI检查配置 |

### B. 命令速查

```bash
# 安装hooks
bash scripts/install-hooks.sh

# 运行架构检查
bash scripts/architecture-check.sh

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

# 检查循环依赖
go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./...

# 跳过检查
git commit --no-verify -m "message"
git push --no-verify
```

### C. 相关文档

- [ZKER企业级开发规范手册 v1.0](./ZKER-企业级开发规范手册_v1.0.md)
- [ZKER全局一致性检查清单 v1.0](./ZKER-全局一致性检查清单_v1.0.md)
- [ZKER统一错误码定义规范](./ZKER-统一错误码定义规范.md)
- [ZKER实现差距分析与研发计划](./ZKER-实现差距分析与研发计划_v1.0.md)

### D. 参考资料

- [Go官方代码审查检查清单](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go)
- [golangci-lint文档](https://golangci-lint.run/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [DDD分层架构](https://herbertograca.com/2017/09/14/ddd-the-clear-architecture/)

---

**更新日志**:

| 版本 | 日期 | 更新内容 |
|------|------|---------|
| v1.0 | 2025-01-01 | 初始版本，完整CI/CD架构合规检查系统 |

---

**维护者**: ZKER开发团队
**反馈**: 请通过GitHub Issues提交问题和建议
