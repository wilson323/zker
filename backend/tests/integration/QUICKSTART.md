# 组织中心集成测试 - 快速参考指南

## 一、测试文件结构

```
backend/tests/integration/
├── org_integration_test.go    # 主测试文件
├── README.md                  # 详细文档
├── QUICKSTART.md              # 本文件
├── run-tests.sh               # Linux/Mac运行脚本
└── run-tests.bat              # Windows运行脚本
```

## 二、快速开始

### 1. 前置条件检查

```bash
# 检查Docker
docker --version
docker ps

# 检查Go
go version
```

### 2. 运行所有测试

```bash
# Windows
cd backend
.\tests\integration\run-tests.bat

# Linux/Mac
cd backend
chmod +x ./tests/integration/run-tests.sh
./tests/integration/run-tests.sh
```

### 3. 运行特定测试

```bash
# 只测试组织管理
go test -v -tags=integration ./tests/integration/... -run TestOrgIntegration/组织管理

# 只测试多租户隔离
go test -v -tags=integration ./tests/integration/... -run TestOrgIntegration/多租户
```

## 三、测试用例清单

### 1. 组织管理完整CRUD流程 (7个测试)
- ✅ 创建公司组织
- ✅ 创建分公司组织
- ✅ 创建项目组组织
- ✅ 查询组织树
- ✅ 更新组织信息
- ✅ 查询后代组织
- ✅ 删除项目组组织

### 2. 部门管理完整CRUD流程 (2个测试)
- ✅ 创建根部门
- ✅ 查询部门树

### 3. 岗位管理完整CRUD流程 (2个测试)
- ✅ 创建岗位
- ✅ 查询岗位列表

### 4. 员工管理完整CRUD流程 (2个测试)
- ✅ 创建员工
- ✅ 员工转正

### 5. 通讯录查询流程 (1个测试)
- ✅ 查询组织通讯录

### 6. HR生命周期流程 (3个测试)
- ✅ 创建员工合同
- ✅ 员工调岗
- ✅ 员工离职

### 7. 多租户隔离验证 (2个测试)
- ✅ 租户数据隔离
- ✅ 租户完全隔离验证

### 8. 事务测试 (2个测试)
- ✅ 事务提交成功
- ✅ 事务回滚失败操作

### 9. 复杂业务场景 (2个测试)
- ✅ 组织树结构查询
- ✅ 员工状态转换

**总计: 25个测试用例**

## 四、常用命令

### 基础测试命令

```bash
# 运行所有集成测试（详细输出）
go test -v -tags=integration ./tests/integration/...

# 运行并生成覆盖率报告
go test -v -cover -tags=integration ./tests/integration/...

# 跳过集成测试（快速测试）
go test -v -short ./tests/integration/...

# 并行运行测试（加速）
go test -v -parallel 8 -tags=integration ./tests/integration/...
```

### 调试命令

```bash
# 只运行特定测试
go test -v -tags=integration ./tests/integration/... -run TestOrgIntegration/testOrgCRUD

# 显示详细日志
go test -v -tags=integration ./tests/integration/... -run TestOrgIntegration 2>&1 | tee test.log

# CPU性能分析
go test -v -cpuprofile=cpu.prof -tags=integration ./tests/integration/...
go tool pprof cpu.prof

# 内存分析
go test -v -memprofile=mem.prof -tags=integration ./tests/integration/...
go tool pprof mem.prof
```

### 容器管理

```bash
# 查看运行中的容器
docker ps

# 查看容器日志
docker logs <container_id>

# 停止所有测试容器
docker ps -q --filter "ancestor=mysql:8.4.5" | xargs docker stop

# 清理容器
docker system prune -f
```

## 五、故障排查

### 问题1: Docker启动失败

**症状**: `启动MySQL容器失败`

**解决方案**:
```bash
# 检查Docker是否运行
docker ps

# 检查端口占用
netstat -ano | findstr :3306

# 重启Docker Desktop
```

### 问题2: 数据库连接失败

**症状**: `连接数据库失败`

**解决方案**:
```bash
# 检查容器日志
docker logs <container_id>

# 增加重试次数（代码中已实现20次重试）
# 每次等待2秒，最多等待40秒
```

### 问题3: 测试超时

**症状**: 测试运行超过10分钟

**解决方案**:
```bash
# 增加超时时间
go test -v -timeout 20m -tags=integration ./tests/integration/...

# 或跳过慢速测试
go test -v -short -tags=integration ./tests/integration/...
```

### 问题4: 端口冲突

**症状**: `bind: address already in use`

**解决方案**:
```bash
# 查找占用端口的进程
netstat -ano | findstr :3306

# 关闭占用端口的进程
taskkill /PID <pid> /F

# 或修改Docker端口映射
```

## 六、测试数据说明

### 租户数据
- tenant-001: 主租户（用于大部分测试）
- tenant-002: 用于多租户隔离测试

### 组织编码规则
- COMP001, COMP002...: 公司组织
- DIV001, DIV002...: 分公司组织
- PROJ001, PROJ002...: 项目组组织

### 员工编码规则
- EMP001, EMP002...: 员工工号

### 部门编码规则
- TECH, TECH3, TECH4...: 技术部
- PRODUCT, PRODUCT4...: 产品部

### 岗位编码规则
- SENIOR_SE: 高级软件工程师
- SE001, SE003...: 软件工程师
- PM004...: 产品经理

## 七、测试执行时间（参考）

| 测试组 | 用例数 | 执行时间 |
|--------|--------|----------|
| 组织管理CRUD | 7 | ~2s |
| 部门管理CRUD | 2 | ~1.5s |
| 岗位管理CRUD | 2 | ~1s |
| 员工管理CRUD | 2 | ~2.5s |
| 通讯录查询 | 1 | ~0.5s |
| HR生命周期 | 3 | ~2s |
| 多租户隔离 | 2 | ~1s |
| 事务测试 | 2 | ~1.5s |
| 复杂业务场景 | 2 | ~3s |
| **总计** | **25** | **~15s** |

## 八、扩展测试

### 添加新测试用例

```go
t.Run("新测试场景", func(t *testing.T) {
    ctx := context.Background()

    // 1. 准备测试数据
    org, _ := setup.OrgSvc.CreateOrganization(ctx, &service.CreateOrganizationRequest{
        TenantID:  setup.TenantID,
        OrgName:   "测试组织",
        OrgType:   entity.OrgTypeCompany,
        OrgCode:   "TEST_ORG",
        SortOrder: 1,
    })

    // 2. 执行测试操作
    result, err := setup.Service.DoSomething(ctx, params)

    // 3. 验证结果
    require.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, expected, result.Field)
})
```

### 添加新Service测试

1. 在 `setupTestEnvironment` 中添加新Service
2. 在 `TestSetup` 结构中添加字段
3. 创建新的测试函数 `testXXXCRUD`
4. 在 `TestOrgIntegration` 中注册测试

## 九、CI/CD集成

### GitHub Actions示例

```yaml
name: Org Integration Tests

on:
  push:
    paths:
      - 'backend/domain/org/**'
      - 'backend/tests/integration/**'
  pull_request:
    paths:
      - 'backend/domain/org/**'
      - 'backend/tests/integration/**'

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24.0'

      - name: Run integration tests
        run: |
          cd backend
          go test -v -cover -tags=integration ./tests/integration/...
        timeout-minutes: 15
```

## 十、相关文档

- [详细文档](./README.md)
- [企业级开发规范手册](../../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [全局一致性检查清单](../../../docs/企业级功能完善与统一性设计方案/ZKER-全局一致性检查清单_v1.0.md)
- [统一错误码定义规范](../../../docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)

## 十一、联系和支持

如有问题或建议，请联系：
- 研发A（后端架构师）：负责组织中心架构设计和测试框架
- 研发B（后端工程师）：负责测试性能优化和CI/CD集成

---

**最后更新**: 2025-01-01
**版本**: 1.0.0
