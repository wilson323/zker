# golangci-lint配置说明

本目录包含ZKER项目的golangci-lint配置文件。

## 文件说明

- **`.golangci.yml`**: golangci-lint主配置文件
  - 定义了启用的linters
  - 配置了各linter的参数
  - 设置了排除规则
  - 定义了输出格式

## 使用方法

### 本地使用

```bash
# 运行检查（从项目根目录）
cd backend
golangci-lint run --config=../.github/linters/.golangci.yml

# 自动修复问题
golangci-lint run --fix --config=../.github/linters/.golangci.yml

# 仅运行特定linter
golangci-lint run --disable-all --enable=gofmt,goimports --config=../.github/linters/.golangci.yml
```

### CI/CD使用

配置文件会被GitHub Actions工作流自动使用：
- 文件: `.github/workflows/architecture-compliance.yml`
- 步骤: "5️⃣ 代码规范检查 (Lint Check)"

## 配置亮点

### 1. 启用的Linters

- **格式化**: gofmt, goimports
- **静态分析**: govet, staticcheck, structcheck, varcheck
- **错误处理**: errcheck, errorlint
- **代码质量**: unused, deadcode, ineffassign, unconvert
- **复杂度**: gocyclo (圈复杂度≤15), dupl (重复代码≤100 tokens)
- **性能**: prealloc (slice预分配优化)
- **安全**: gosec (安全漏洞扫描)
- **命名**: revive (替代golint，支持更多规则)
- **其他**: misspell, lll, nolintlint

### 2. 排除规则

自动排除以下文件：
- 测试文件 (`*_test.go`)
- 生成的文件 (`*.gen.go`, `*.pb.go`)
- vendor目录

测试文件放宽限制：
- 允许重复代码 (dupl)
- 允许高圈复杂度 (gocyclo)
- 允许未处理的错误 (errcheck)
- 允许安全问题 (gosec)
- 允许长行 (lll)

### 3. 严重性配置

- 默认严重性: `error`
- 警告级别（不阻止CI）:
  - dupl (重复代码)
  - goconst (重复常量)
  - misspell (拼写)
  - lll (行长度)

## 自定义配置

如需自定义配置，请编辑 `.golangci.yml` 文件：

```yaml
linters-settings:
  gocyclo:
    min-complexity: 20  # 修改圈复杂度阈值

  lll:
    line-length: 150    # 修改行长度阈值

linters:
  enable:
    - your-custom-linter  # 添加自定义linter
```

## 相关文档

- [golangci-lint官方文档](https://golangci-lint.run/)
- [支持的Linters](https://golangci-lint.run/usage/linters/)
- [配置示例](https://golangci-lint.run/usage/configuration/)

---

**维护者**: ZKER开发团队
