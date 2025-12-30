# 组织中心集成测试

## 概述

本测试套件提供组织中心的完整集成测试，使用 testcontainers + MySQL 8.4.5 进行真实环境测试。

## 测试覆盖范围

### 1. 核心业务流程测试

#### 组织管理完整CRUD流程
- ✅ 创建公司组织
- ✅ 创建分公司组织
- ✅ 创建项目组组织
- ✅ 查询组织详情
- ✅ 查询组织树
- ✅ 更新组织信息
- ✅ 查询子组织
- ✅ 查询祖先组织
- ✅ 查询后代组织
- ✅ 分页查询组织列表
- ✅ 删除组织

#### 部门管理完整CRUD流程
- ✅ 创建根部门
- ✅ 创建子部门
- ✅ 查询部门树
- ✅ 更新部门信息
- ✅ 删除部门

#### 岗位管理完整CRUD流程
- ✅ 创建岗位
- ✅ 查询岗位列表
- ✅ 更新岗位信息
- ✅ 删除岗位

#### 员工管理完整CRUD流程
- ✅ 创建员工
- ✅ 查询员工详情
- ✅ 查询员工列表
- ✅ 更新员工信息
- ✅ 员工转正
- ✅ 删除员工

#### 通讯录查询流程
- ✅ 查询组织通讯录
- ✅ 查询部门通讯录
- ✅ 搜索员工

#### HR生命周期流程
- ✅ 创建员工合同
- ✅ 签署合同
- ✅ 员工调岗
- ✅ 员工离职
- ✅ 审批离职

### 2. 多租户隔离测试

#### 租户数据隔离
- ✅ 不同租户的数据完全隔离
- ✅ 租户ID过滤验证
- ✅ 跨租户访问拒绝验证

#### 租户完全隔离验证
- ✅ 所有表的数据都按tenant_id隔离
- ✅ organizations表隔离验证
- ✅ employees表隔离验证
- ✅ departments表隔离验证

### 3. 事务测试

#### 事务提交成功
- ✅ 多表操作事务提交
- ✅ 数据一致性验证

#### 事务回滚失败操作
- ✅ 错误触发事务回滚
- ✅ 回滚后数据未变更

#### 并发安全
- ✅ 并发创建组织（编码唯一性约束）
- ✅ 乐观锁验证
- ✅ 数据竞争检测

### 4. 复杂业务场景测试

#### 组织树结构查询
- ✅ 创建多层组织树（3级以上）
- ✅ 查询祖先链验证
- ✅ 查询后代树验证
- ✅ 闭包表路径验证

#### 员工状态转换
- ✅ 试用期 → 在职状态转换
- ✅ 离职状态转换
- ✅ 状态流转规则验证

#### 合同状态流转
- ✅ 草稿 → 生效状态转换
- ✅ 签署流程验证
- ✅ 合同到期处理

#### 离职审批流程
- ✅ 提交离职申请（pending状态）
- ✅ 审批通过（approved状态）
- ✅ 审批拒绝（rejected状态）
- ✅ 工作交接流程

## 技术架构

### 测试环境
- **数据库**: MySQL 8.4.5 (Docker容器)
- **ORM**: GORM v1.25.11
- **测试框架**: testify + testcontainers-go
- **字符集**: utf8mb4_unicode_ci

### 测试模式
- ✅ 真实数据库测试（不使用Mock）
- ✅ 完整的端到端测试
- ✅ 表驱动测试
- ✅ 子测试组织（t.Run）
- ✅ Setup/Teardown模式
- ✅ 独立测试数据

## 运行测试

### 前置条件

1. **安装Docker Desktop**
   ```bash
   # 验证Docker安装
   docker --version
   docker ps
   ```

2. **确保端口未被占用**
   - MySQL端口: 3306
   - MySQL X协议端口: 33060

### 运行所有测试

```bash
cd backend

# 运行所有集成测试
go test -v ./tests/integration/...

# 运行并显示覆盖率
go test -v -cover ./tests/integration/...

# 运行特定测试
go test -v ./tests/integration/... -run TestOrgIntegration
```

### 运行特定测试用例

```bash
# 只测试组织管理
go test -v ./tests/integration/... -run TestOrgIntegration/组织管理

# 只测试多租户隔离
go test -v ./tests/integration/... -run TestOrgIntegration/多租户隔离

# 只测试事务
go test -v ./tests/integration/... -run TestOrgIntegration/事务测试
```

### 并行测试

```bash
# 并行运行测试（加速）
go test -v -parallel 4 ./tests/integration/...
```

### 性能分析

```bash
# CPU性能分析
go test -v -cpuprofile=cpu.prof ./tests/integration/...
go tool pprof cpu.prof

# 内存分析
go test -v -memprofile=mem.prof ./tests/integration/...
go tool pprof mem.prof
```

## 测试数据准备

### Fixtures

测试使用以下fixture数据：

1. **租户数据**
   - tenant-001: 主租户
   - tenant-002: 用于多租户隔离测试

2. **组织数据**
   - 自动生成唯一编码（使用时间戳）
   - 每个测试独立的数据

3. **员工数据**
   - 基本信息完整
   - 包含试用期、在职、离职等状态

4. **合同数据**
   - 草稿、生效、过期、终止等状态

### 数据清理

- 每个测试独立的数据
- 测试结束后自动清理Docker容器
- 软删除数据保留（deleted_at标记）

## 企业级规范遵循

### 代码规范

1. **函数长度 < 50行**
   - ✅ 所有测试函数遵循此规范
   - ✅ 复杂逻辑拆分为辅助函数

2. **中文注释完整**
   - ✅ 每个测试函数都有清晰的注释
   - ✅ 复杂逻辑都有说明

3. **表驱动测试**
   - ✅ 使用t.Run组织子测试
   - ✅ 测试用例清晰分离

4. **Setup/Teardown模式**
   - ✅ setupTestEnvironment初始化环境
   - ✅ defer setup.Cleanup()清理资源

### 测试规范

1. **断言清晰**
   - ✅ 使用testify的require/assert
   - ✅ 错误信息描述准确

2. **测试独立性**
   - ✅ 每个测试独立运行
   - ✅ 无依赖关系

3. **数据隔离**
   - ✅ 使用随机ID避免冲突
   - ✅ 每个租户独立数据

## 测试覆盖率目标

- ✅ Service层: ≥ 80%
- ✅ Repository层: ≥ 70%
- ✅ 核心业务流程: 100%

## 故障排查

### 常见问题

#### 1. Docker容器启动失败
```bash
# 检查Docker是否运行
docker ps

# 检查端口占用
netstat -ano | findstr :3306
```

#### 2. 数据库连接失败
```bash
# 检查容器日志
docker logs <container_id>

# 增加重试次数（代码中已实现）
for i := 0; i < 10; i++ {
    // 连接重试逻辑
}
```

#### 3. 测试超时
```bash
# 增加超时时间
go test -v -timeout 10m ./tests/integration/...
```

#### 4. 端口冲突
```bash
# 修改docker-compose端口映射
# 或关闭占用端口的进程
```

## 扩展测试

### 添加新测试用例

```go
t.Run("新测试场景", func(t *testing.T) {
    // 1. 准备测试数据
    ctx := context.Background()

    // 2. 执行测试操作
    result, err := setup.Service.DoSomething(ctx, params)

    // 3. 验证结果
    require.NoError(t, err)
    assert.NotNil(t, result)
})
```

### 添加新Service测试

```go
// 1. 在setup中添加新Service
NewService := service.NewXXXService(repo, db)

// 2. 在TestSetup结构中添加字段
XXXSvc *service.XXXService

// 3. 添加测试函数
func testXXXCRUD(t *testing.T, setup *TestSetup) {
    // 测试代码
}
```

## CI/CD集成

### GitHub Actions示例

```yaml
name: Integration Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      docker:
        image: docker:latest

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24.0'

      - name: Run integration tests
        run: |
          cd backend
          go test -v -cover ./tests/integration/...
```

## 性能基线

### 测试执行时间（参考）

- 组织管理CRUD: ~2s
- 部门管理CRUD: ~1.5s
- 员工管理CRUD: ~2.5s
- 多租户隔离: ~1s
- 事务测试: ~1.5s
- 复杂业务场景: ~3s

**总计**: ~11.5s

### 优化建议

1. 并行运行测试（-parallel 4）
2. 跳过某些测试（-skip）
3. 使用测试缓存

## 相关文档

- [企业级开发规范手册](../../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [全局一致性检查清单](../../../docs/企业级功能完善与统一性设计方案/ZKER-全局一致性检查清单_v1.0.md)
- [统一错误码定义规范](../../../docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)
- [数据迁移方案](../../../docs/企业级功能完善与统一性设计方案/ZKER-数据迁移方案_v1.0.md)

## 联系方式

如有问题，请联系：
- 研发A（后端架构师）: 负责组织中心架构设计
- 研发B（后端工程师）: 负责测试框架和性能测试
