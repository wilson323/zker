# E2E测试快速启动指南

**目标**: 快速执行Coze Studio的E2E测试套件

---

## ⚡ 5分钟快速启动

### Step 1: 启动Docker (2分钟)

```bash
# Windows: 启动Docker Desktop应用程序
# 在开始菜单中找到Docker Desktop并启动

# 等待Docker完全启动(约1-2分钟)
# 验证Docker是否运行
docker ps

# 预期输出:
# CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
# (如果没有容器运行,这是正常的)
```

### Step 2: 安装依赖 (30秒)

```bash
cd D:/code/coze-studio/backend/tests/e2e
go mod tidy
```

### Step 3: 执行测试 (2-3分钟)

```bash
# 方式1: 使用脚本(推荐)
./run_e2e_tests.sh

# 方式2: 使用go test
go test -v -timeout 30m
```

### Step 4: 查看结果 (即时)

测试执行后会自动生成:
- `test_results.txt` - 测试结果日志
- `coverage.out` - 覆盖率数据(如果使用--cover)
- `coverage.html` - HTML覆盖率报告(如果使用--report)

---

## 🎯 测试场景概览

| 场景 | 测试文件 | 预期时间 | 关键验证点 |
|------|---------|---------|-----------|
| 租户注册 | `tenant_registration_journey_test.go` | 30-60秒 | 租户创建、RBAC初始化、配额分配 |
| Bot审核 | `bot_creation_approval_journey_test.go` | 40-80秒 | 创建、审核、批准/拒绝、发布 |
| 组织权限 | `organization_permission_journey_test.go` | 50-90秒 | 组织创建、权限分配、数据过滤 |
| 对话记忆 | `conversation_memory_journey_test.go` | 40-70秒 | 对话、记忆提取、记忆检索 |

---

## 📋 执行选项

### 选项1: 执行所有测试

```bash
./run_e2e_tests.sh
```

**预期时间**: 8-12分钟
**预期通过率**: ≥95%

### 选项2: 执行特定测试

```bash
# 只测试租户注册
./run_e2e_tests.sh --run TestE2E_TenantRegistrationJourney

# 只测试Bot审核
./run_e2e_tests.sh --run TestE2E_BotCreationWithApprovalJourney
```

**预期时间**: 1-2分钟/测试
**预期通过率**: 100%

### 选项3: 生成覆盖率报告

```bash
./run_e2e_tests.sh --cover --report
```

**输出文件**:
- `coverage.out` - 覆盖率数据
- `coverage.html` - HTML可视化报告

**预期覆盖率**: ≥70%

### 选项4: 并发执行

```bash
./run_e2e_tests.sh --parallel 4
```

**预期时间**: 3-5分钟
**注意**: 需要确保数据隔离

---

## ⚠️ 常见问题

### 问题1: Docker未启动

**错误信息**:
```
failed to connect to the docker API
```

**解决方案**:
1. 启动Docker Desktop应用程序
2. 等待Docker引擎完全启动
3. 验证: `docker ps`

### 问题2: MySQL容器启动失败

**错误信息**:
```
MySQL容器启动失败
```

**解决方案**:
```bash
# 检查Docker资源
docker stats

# 检查端口占用
netstat -an | grep 3306

# 重启Docker
# Windows: 右键Docker Desktop图标 -> Restart
```

### 问题3: 测试超时

**错误信息**:
```
test timed out after 30m
```

**解决方案**:
```bash
# 增加超时时间
./run_e2e_tests.sh -timeout 60m
```

### 问题4: 端口冲突

**错误信息**:
```
port 3306 already in use
```

**解决方案**:
```bash
# Windows: 检查端口占用
netstat -ano | findstr :3306

# 停止占用端口的进程
taskkill /PID <PID> /F

# 或者修改MySQL端口配置
```

---

## 📊 预期测试结果

### 成功标准

✅ 所有P0用户旅程测试通过(4/4)
✅ 端到端流程无阻塞
✅ 关键业务流程验证通过
✅ 测试覆盖率 ≥70%

### 测试通过率预期

| 类别 | 通过数 | 总数 | 通过率 |
|------|-------|------|-------|
| P0用户旅程 | 4 | 4 | 100% |
| P1验证测试 | 9 | 10 | ≥90% |
| P2并发测试 | 2 | 2 | ≥80% |
| **总计** | **15** | **16** | **≥95%** |

---

## 🔍 详细日志

### 查看测试日志

```bash
# 查看完整测试日志
cat test_results.txt

# 查看失败的测试
grep "FAIL:" test_results.txt

# 查看通过的测试
grep "PASS:" test_results.txt
```

### 查看覆盖率

```bash
# 查看覆盖率摘要
go tool cover -func=coverage.out | tail -1

# 在浏览器中查看HTML报告
# Windows:
start coverage.html

# macOS:
open coverage.html

# Linux:
xdg-open coverage.html
```

---

## 🚀 下一步

### 测试通过后

1. **生成完整报告**: `./run_e2e_tests.sh --cover --report`
2. **查看覆盖率**: `go tool cover -html=coverage.out`
3. **提交结果**: 将报告提交给团队

### 测试失败后

1. **查看失败日志**: `grep "FAIL:" test_results.txt`
2. **重新运行失败测试**: `./run_e2e_tests.sh --run <TestName>`
3. **修复问题**: 根据错误信息修复
4. **重新验证**: 运行完整测试套件

---

## 📞 获取帮助

### 文档

- [E2E测试README](./README.md) - 详细文档
- [E2E测试执行报告](./E2E_TEST_EXECUTION_REPORT.md) - 执行报告

### 命令行帮助

```bash
./run_e2e_tests.sh --help
```

### 常用命令

```bash
# 列出所有测试
go test -list=.

# 运行特定测试
go test -v -run TestE2E_TenantRegistrationJourney

# 跳过短测试
go test -v -short=false

# 详细输出
go test -v
```

---

**快速启动指南版本**: v1.0
**最后更新**: 2025-01-01
