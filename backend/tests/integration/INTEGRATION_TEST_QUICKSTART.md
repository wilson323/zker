# 集成测试快速开始指南

## 🚀 快速开始

### 一键运行（推荐）

#### Linux/Mac
```bash
chmod +x scripts/run-integration-tests.sh
./scripts/run-integration-tests.sh
```

#### Windows
```cmd
scripts\run-integration-tests.bat
```

## 📦 测试套件

### 后端集成测试
- 计费集成测试: `go test -v ./tests/integration/billing/...`
- 多租户集成测试: `go test -v ./tests/integration/tenant/...`
- API集成测试: `go test -v ./tests/integration/api/...`

### 前端E2E测试
- 计费管理: `npx playwright test billing-management.spec.ts`
- 知识管理: `npx playwright test knowledge-management.spec.ts`

## 📊 测试覆盖率目标
- 总体: ≥60%
- 计费系统: 70%+
- 多租户: 65%+
- API层: 60%+

## 📝 更多信息
参见: [INTEGRATION_TEST_IMPLEMENTATION_REPORT.md](../../INTEGRATION_TEST_IMPLEMENTATION_REPORT.md)
