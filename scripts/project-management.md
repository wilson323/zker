# 项目管理和进度跟踪脚本使用手册

> **最后更新**: 2025-01-01
> **维护人**: 项目团队

---

## 📖 概述

本目录包含用于项目管理和进度跟踪的实用脚本，辅助4人研发团队高效协作。

**脚本列表**:
1. `daily-check.sh` - 每日环境检查
2. `weekly-progress-report.sh` - 周进度报告生成
3. `frontend-check.sh` - 前端代码检查
4. `backend-check.sh` - 后端代码检查

---

## 🚀 快速开始

### 1. daily-check.sh - 每日环境检查

**用途**: 检查开发环境的服务状态

**使用方法**:
```bash
chmod +x scripts/daily-check.sh
./scripts/daily-check.sh
```

**检查项**:
- ✅ Docker服务状态
- ✅ MySQL数据库连接
- ✅ Redis缓存状态
- ✅ Elasticsearch服务
- ✅ 后端服务健康检查
- ✅ 前端服务健康检查
- ✅ 磁盘空间使用
- ✅ Docker资源统计

**输出示例**:
```
🔍 每日环境检查 - 2025-01-01 09:00:00
========================================

📦 检查Docker服务...
✅ Docker服务正常

📦 检查中间件容器...
NAMES              STATUS
coze-mysql-1       Up 2 hours
coze-redis-1       Up 2 hours

💾 检查MySQL数据库...
✅ MySQL数据库正常
📊 数据库列表:
information_schema
coze
coze_test
```

**推荐使用时机**: 每天开始工作前

---

### 2. weekly-progress-report.sh - 周进度报告生成

**用途**: 自动生成每周进度报告模板

**使用方法**:
```bash
chmod +x scripts/weekly-progress-report.sh
./scripts/weekly-progress-report.sh
```

**生成的报告包含**:
- 📋 本周完成任务清单
- 📊 进度统计表
- ⚠️ 风险和问题列表
- 🎯 下周计划
- 📈 关键指标（代码质量、性能、监控）

**输出文件**: `docs/weekly-reports/week-{N}-{YEAR}.md`

**工作流程**:
1. 运行脚本生成报告模板
2. 编辑报告，填写实际进度
3. 提交到代码仓库
4. 在周会上分享

**推荐使用时机**: 每周五下午

---

### 3. frontend-check.sh - 前端代码检查

**用途**: 前端代码提交前自动检查

**使用方法**:
```bash
chmod +x scripts/frontend-check.sh
./scripts/frontend-check.sh
```

**检查项**:
- ✅ TypeScript类型检查
- ✅ ESLint代码规范
- ✅ 单元测试
- ✅ 测试覆盖率 (≥ 70%)
- ✅ 构建检查

**推荐使用时机**:
- Git commit前
- Pull Request创建前
- CI/CD pipeline中

**Git Hook集成**:
```bash
# 在 .git/hooks/pre-commit 中添加
./scripts/frontend-check.sh || exit 1
```

---

### 4. backend-check.sh - 后端代码检查

**用途**: 后端代码提交前自动检查

**使用方法**:
```bash
chmod +x scripts/backend-check.sh
./scripts/backend-check.sh
```

**检查项**:
- ✅ 代码格式 (gofmt)
- ✅ 代码规范 (golangci-lint)
- ✅ 单元测试
- ✅ 测试覆盖率 (≥ 80%)
- ✅ 竞态检查 (-race)
- ✅ 构建检查
- ✅ 安全检查 (go vet)

**推荐使用时机**:
- Git commit前
- Pull Request创建前
- CI/CD pipeline中

**Git Hook集成**:
```bash
# 在 .git/hooks/pre-commit 中添加
./scripts/backend-check.sh || exit 1
```

---

## 📅 每日工作流

### 开始工作前

```bash
# 1. 检查环境状态
./scripts/daily-check.sh

# 2. 拉取最新代码
git pull origin develop

# 3. 查看今日任务
# （在项目管理工具中查看）
```

### 代码开发

```bash
# 1. 创建功能分支
git checkout -b feature/your-feature

# 2. 开发代码
# ... 编写代码 ...

# 3. 代码检查（提交前）
./scripts/frontend-check.sh  # 前端项目
./scripts/backend-check.sh   # 后端项目

# 4. 提交代码
git add .
git commit -m "feat: add your feature"

# 5. 推送并创建PR
git push origin feature/your-feature
```

### 结束工作时

```bash
# 1. 更新任务状态（在项目管理工具中）
# 2. 准备明天的任务列表
# 3. 提交未完成的工作
```

---

## 📊 周度工作流

### 周一

- [ ] 参加每日站会（09:00-09:15）
- [ ] 查看本周任务计划
- [ ] 更新进度报告

### 周二-周四

- [ ] 正常开发工作
- [ ] 代码审查（14:00-15:00）
- [ ] 技术讨论（按需）

### 周五

- [ ] 生成周进度报告
  ```bash
  ./scripts/weekly-progress-report.sh
  ```
- [ ] 填写实际完成情况
- [ ] 记录风险和问题
- [ ] 提交报告并分享
- [ ] 参加周会（16:00-17:00）

---

## 🛠️ 安装依赖

### 前端工具

```bash
# 安装Rush
npm install -g @microsoft/rush

# 安装依赖
cd frontend
rush update
```

### 后端工具

```bash
# 安装Go工具
go install golang.org/x/tools/cmd/...@latest

# 安装golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# 验证安装
go version
golangci-lint version
```

---

## 🔧 配置Git Hook（可选）

### 自动代码检查

创建 `.git/hooks/pre-commit`:

```bash
#!/bin/bash
# Git Pre-Commit Hook

# 检查修改的文件
FRONTEND_CHANGED=$(git diff --cached --name-only | grep "^frontend/" || true)
BACKEND_CHANGED=$(git diff --cached --name-only | grep "^backend/" || true)

# 运行前端检查
if [ -n "$FRONTEND_CHANGED" ]; then
    echo "🔍 检查前端代码..."
    ./scripts/frontend-check.sh || exit 1
fi

# 运行后端检查
if [ -n "$BACKEND_CHANGED" ]; then
    echo "🔍 检查后端代码..."
    ./scripts/backend-check.sh || exit 1
fi

echo "✅ 代码检查通过"
```

添加执行权限:
```bash
chmod +x .git/hooks/pre-commit
```

---

## 📊 进度报告模板

### 每周进度报告结构

```markdown
# 第N周进度报告

**报告日期**: YYYY-MM-DD
**周次**: 第N周 (YYYY)
**日期范围**: YYYY-MM-DD - YYYY-MM-DD

---

## 📋 本周完成任务

### UI组件实现（研发C）
- [x] 完成租户管理组件（3/3天）
- [x] 完成RBAC权限组件（4/4天）
- [ ] 订阅配额组件（2/3天，进行中）

### 集成测试覆盖（研发B）
- [x] 完成后端集成测试框架（2/2天）
- [ ] 租户管理API测试（1/3天）

---

## 📊 进度统计

| 模块 | 计划工期 | 已用工期 | 完成度 | 状态 |
|------|---------|---------|--------|------|
| UI组件实现 | 10天 | 9天 | 90% | 🟢 正常 |
| 集成测试 | 6天 | 3天 | 50% | 🟢 正常 |

**图例**:
- 🟢 正常（按计划进行）
- 🟡 预警（轻微延期）
- 🔴 严重（严重延期）

---

## ⚠️ 风险和问题

| 风险/问题 | 影响 | 缓解措施 | 负责人 | 状态 |
|----------|------|---------|--------|------|
| 权限树性能问题 | 中 | 使用虚拟滚动 | 研发C | 🔄 |

---

## 🎯 下周计划

**UI组件实现（研发C）**:
- [ ] 完成订阅配额组件
- [ ] 开始Playwright E2E测试

**集成测试覆盖（研发B）**:
- [ ] 完成后端集成测试
- [ ] 搭建E2E测试框架
```

---

## 🔍 故障排查

### 问题1: 脚本无法执行

**症状**: `Permission denied`

**解决方案**:
```bash
chmod +x scripts/*.sh
```

### 问题2: 环境检查失败

**症状**: 服务连接失败

**解决方案**:
```bash
# 启动中间件服务
cd docker
docker-compose up -d

# 等待服务就绪
sleep 30

# 重新检查
./scripts/daily-check.sh
```

### 问题3: 代码检查不通过

**症状**: 测试覆盖率不足

**解决方案**:
```bash
# 前端：添加测试用例
cd frontend
rush test

# 后端：添加测试用例
cd backend
go test ./...
```

---

## 📚 相关文档

- [实施计划-5任务执行指南](../docs/企业级功能完善与统一性设计方案/实施计划-5任务执行指南_v1.0.md)
- [企业级开发规范手册](../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [全局一致性检查清单](../docs/企业级功能完善与统一性设计方案/ZKER-全局一致性检查清单_v1.0.md)

---

## 📞 获取帮助

**问题反馈**: 创建GitHub Issue
**功能建议**: 创建GitHub Issue或联系项目团队

**团队联系**:
- @研发A - 后端架构师
- @研发B - 后端工程师
- @研发C - 前端工程师
- @研发D - DevOps工程师

---

**文档版本**: v1.0
**最后更新**: 2025-01-01
