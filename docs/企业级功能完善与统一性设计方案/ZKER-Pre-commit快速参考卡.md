# ZKER Pre-commit Hook 快速参考卡

## 🚀 快速安装

```bash
bash scripts/install-hooks.sh
```

## ✅ 6项自动检查

| # | 检查项 | 说明 | 修复命令 |
|---|--------|------|----------|
| 1 | API响应格式 | 禁止`c.JSON()`，使用`httputil` | 手动修改代码 |
| 2 | 错误处理 | 禁止`errors.New()`，使用`errorx` | 手动修改代码 |
| 3 | 安全问题 | 禁止`panic()` | 改用`return error` |
| 4 | 代码格式 | `gofmt`格式检查 | `gofmt -w backend/` |
| 5 | 静态检查 | `go vet`检查 | 手动修复问题 |
| 6 | 单元测试 | 运行修改文件的测试 | `go test -short ./...` |

## 📝 常用命令

```bash
# 手动运行pre-commit检查
./.githooks/pre-commit

# 跳过检查（不推荐）
git commit --no-verify -m "message"

# 运行完整测试
cd backend && go test ./...

# 格式化代码
gofmt -w backend/

# 静态检查
go vet ./...

# 运行linter
golangci-lint run
```

## 🔧 代码修复示例

### API响应格式
```go
// ❌ 错误
c.JSON(http.StatusOK, data)

// ✅ 正确
httputil.BuildSuccessResp(c, data)
```

### 错误处理
```go
// ❌ 错误
return errors.New("not found")

// ✅ 正确
return errorx.New(errno.ErrNotFoundCode)
```

### 安全问题
```go
// ❌ 错误
if bot == nil {
    panic("bot is nil")
}

// ✅ 正确
if bot == nil {
    return nil, fmt.Errorf("bot not found")
}
```

## 📚 详细文档

查看完整配置指南：
```bash
cat docs/企业级功能完善与统一性设计方案/ZKER-Pre-commit-Hook配置指南_v1.0.md
```
