# E2E测试快速检查清单

**用途**: 快速验证E2E测试的完整性和正确性

---

## ✅ 文件完整性检查

### 核心文件（必须存在）

- [ ] `test_application.go` - 测试应用封装
- [ ] `helpers/test_helpers.go` - 测试辅助函数
- [ ] `tenant_registration_journey_test.go` - 租户注册旅程
- [ ] `bot_creation_approval_journey_test.go` - Bot审核旅程
- [ ] `organization_permission_journey_test.go` - 组织权限旅程
- [ ] `conversation_memory_journey_test.go` - 对话记忆旅程

### 文档和脚本（必须存在）

- [ ] `README.md` - E2E测试文档
- [ ] `run_e2e_tests.sh` - 测试执行脚本（可执行权限）
- [ ] `verify_e2e_tests.go` - 测试验证器
- [ ] `E2E_TEST_REPORT.md` - 测试报告

---

## ✅ 代码质量检查

### 单一职责原则（SRP）

- [ ] 每个测试文件只负责一个用户旅程
- [ ] 每个测试函数只负责一个完整的业务流程
- [ ] 每个辅助函数只负责一个任务

### 命名规范

- [ ] 测试文件命名: `{journey}_journey_test.go`
- [ ] 测试函数命名: `TestE2E_{JourneyName}`
- [ ] 子测试命名: `Step{N}_{步骤描述}`
- [ ] 辅助函数命名: `{Action}{Entity}` (如CreateTestTenant)

### 注释完整性

- [ ] 每个测试文件有文件头注释
- [ ] 每个测试函数有功能描述注释
- [ ] 每个辅助函数有职责注释
- [ ] 复杂逻辑有行内注释

---

## ✅ 功能完整性检查

### 租户注册旅程

- [ ] Step1: 提交注册信息
- [ ] Step2: 验证租户已创建
- [ ] Step3: 验证默认角色已创建
- [ ] Step4: 验证默认配额已分配
- [ ] Step5: 验证订阅信息

### Bot审核旅程

- [ ] Step1: 创建Bot（低AI置信度）
- [ ] Step2: 验证Bot为草稿状态
- [ ] Step3: 验证审核任务已创建
- [ ] Step4: 审核人接取任务
- [ ] Step5: 提交审核结果（批准）
- [ ] Step6: 验证Bot已发布
- [ ] Step7: 验证审核历史已记录

### 组织权限旅程

- [ ] Step1: 创建部门
- [ ] Step2: 添加成员到部门
- [ ] Step3: 分配部门数据权限
- [ ] Step4: 验证权限缓存已更新
- [ ] Step5: 验证数据权限过滤生效

### 对话记忆旅程

- [ ] Step1: 第一次对话（用户自我介绍）
- [ ] Step2: 验证自动提取了实体记忆
- [ ] Step3: 第二次对话（验证记忆被注入）
- [ ] Step4: 验证回答风格符合用户偏好
- [ ] Step5: 验证记忆重要性评分
- [ ] Step6: 验证记忆访问次数

---

## ✅ 技术实现检查

### 测试应用封装

- [ ] TestApplication包含HTTP服务器
- [ ] TestApplication包含MySQL容器
- [ ] TestApplication包含GORM数据库
- [ ] TestApplication提供HTTP客户端封装
- [ ] TestApplication提供资源清理方法

### 测试辅助函数

- [ ] CreateTestTenant - 创建测试租户
- [ ] CreateTestUser - 创建测试用户
- [ ] CreateTestBot - 创建测试Bot
- [ ] CreateTestOrganization - 创建测试组织
- [ ] AssertTenantExists - 验证租户存在
- [ ] AssertBotStatus - 验证Bot状态
- [ ] AssertReviewTaskExists - 验证审核任务存在
- [ ] AssertMemoryExists - 验证记忆存在
- [ ] CleanupTenantData - 清理租户数据

### Mock Handlers

- [ ] handleTenantRegister - 租户注册
- [ ] handleCreateTenant - 创建租户
- [ ] handleGetTenant - 获取租户
- [ ] handleListRoles - 列出角色
- [ ] handleGetQuota - 获取配额
- [ ] handleCreateOrganization - 创建组织
- [ ] handleAssignDataPermission - 分配权限
- [ ] handleCreateBot - 创建Bot
- [ ] handleGetBot - 获取Bot
- [ ] handleListReviewTasks - 列出审核任务
- [ ] handleAssignReviewTask - 分配审核任务
- [ ] handleSubmitReview - 提交审核
- [ ] handleChat - 对话
- [ ] handleGetMemories - 获取记忆

---

## ✅ 执行验证检查

### 环境检查

- [ ] Go已安装（版本 >= 1.24）
- [ ] Docker已安装并运行
- [ ] testcontainers-go可访问
- [ ] 网络连接正常

### 测试执行

```bash
# 进入E2E目录
cd backend/tests/e2e

# 执行所有测试
go test -v -timeout 30m

# 预期结果: 16个测试全部通过
```

- [ ] 所有测试通过（无FAIL）
- [ ] 测试执行时间 < 10分钟
- [ ] 无内存泄漏
- [ ] 无资源泄漏

### 覆盖率检查

```bash
# 生成覆盖率
go test -cover -coverprofile=coverage.out
go tool cover -func=coverage.out
```

- [ ] 覆盖率 ≥ 70%
- [ ] 覆盖率报告已生成
- [ ] HTML覆盖率报告可访问

---

## ✅ 文档完整性检查

### README.md

- [ ] 概述和设计原则
- [ ] 测试套件结构
- [ ] 核心组件说明
- [ ] 测试场景详细说明
- [ ] 执行指南
- [ ] 测试覆盖率目标
- [ ] 故障排查指南
- [ ] 最佳实践

### E2E_TEST_REPORT.md

- [ ] 执行摘要
- [ ] 测试范围
- [ ] 架构设计
- [ ] 交付物清单
- [ ] 验证清单
- [ ] 执行指南
- [ ] 预期测试结果
- [ ] 技术栈
- [ ] 后续优化建议

### 代码注释

- [ ] 文件头注释（版权信息）
- [ ] 包注释（package说明）
- [ ] 函数注释（功能描述）
- [ ] 关键逻辑注释

---

## ✅ 集成检查

### 与现有测试集成

- [ ] 不与unit测试冲突
- [ ] 不与integration测试冲突
- [ ] 可独立运行
- [ ] 可在CI/CD中运行

### 与代码库集成

- [ ] 遵循项目代码规范
- [ ] 遵循项目目录结构
- [ ] 使用统一的错误处理
- [ ] 使用统一的日志格式

---

## 🚀 快速验证命令

### 一键验证（推荐）

```bash
# 方式1: 使用验证器
cd backend/tests/e2e
go run verify_e2e_tests.go

# 方式2: 使用脚本
./run_e2e_tests.sh --cover --report
```

### 分步验证

```bash
# 1. 检查文件
ls -la backend/tests/e2e/

# 2. 执行测试
cd backend/tests/e2e
go test -v -timeout 30m

# 3. 生成覆盖率
go test -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 📊 验收标准

### 必须满足（P0）

- ✅ 所有核心文件存在
- ✅ 4个用户旅程测试通过
- ✅ 测试覆盖率 ≥ 70%
- ✅ 测试执行时间 < 10分钟
- ✅ 文档完整

### 应该满足（P1）

- ✅ 测试代码符合规范
- ✅ 测试可独立运行
- ✅ 测试可重复执行
- ✅ 资源正确清理

### 最好满足（P2）

- ✅ 并发测试通过
- ✅ 边界测试完整
- ✅ 错误场景覆盖
- ✅ 性能基线建立

---

## 📝 问题跟踪

### 发现问题

- 问题描述:
- 严重程度:
- 负责人:
- 状态:

### 解决方案

- 解决方案:
- 实施时间:
- 验证结果:

---

**检查清单版本**: v1.0
**最后更新**: 2025-01-01
