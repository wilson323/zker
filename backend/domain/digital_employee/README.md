# 数字员工管理模块

## 概述

数字员工管理模块提供了完整的AI数字员工画像管理、任务分配和绩效统计功能。

## 功能特性

### 1. 员工画像管理
- ✅ 创建员工画像（名称、角色、技能、性格特征）
- ✅ 更新员工画像
- ✅ 获取员工画像详情
- ✅ 列出员工画像（支持分页、角色、状态筛选）
- ✅ 删除员工画像（软删除）
- ✅ 关联Bot和知识库

### 2. 任务分配系统
- ✅ 手动分配任务给员工
- ✅ 基于技能自动分配任务
- ✅ 完成任务标记
- ✅ 失败任务标记
- ✅ 查看员工任务列表
- ✅ 任务状态管理（assigned、in_progress、completed、failed）

### 3. 绩效统计
- ✅ 员工绩效查询（日/周/月）
- ✅ 团队绩效统计
- ✅ 完成率自动计算
- ✅ 平均响应时间统计
- ✅ 客户评分支持

## 架构设计

### 分层架构

```
backend/domain/digital_employee/
├── entity/                   # 实体层
│   ├── profile.go           # 员工画像实体
│   ├── task_assignment.go   # 任务分配实体
│   └── performance.go       # 绩效统计实体
├── repository/              # 仓储接口层
│   ├── profile_repository.go
│   ├── task_repository.go
│   └── performance_repository.go
├── internal/dal/           # 数据访问层
│   └── dao.go             # DAL实现
└── service/               # 服务层
    ├── profile_service.go
    ├── task_service.go
    ├── performance_service.go
    └── service.go         # 服务容器
```

### DDD设计原则

- **Entity**: 纯业务实体，无技术依赖
- **Repository**: 数据访问接口抽象
- **DAL**: Repository的具体实现（GORM）
- **Service**: 业务逻辑编排

## API接口

### 员工画像管理

```
POST   /api/v1/digital-employees              # 创建员工画像
PUT    /api/v1/digital-employees/:employee_id # 更新员工画像
GET    /api/v1/digital-employees/:employee_id # 获取员工画像
GET    /api/v1/digital-employees              # 列出员工画像
DELETE /api/v1/digital-employees/:employee_id # 删除员工画像
```

### 任务分配

```
POST /api/v1/digital-employees/tasks/assign                    # 分配任务
POST /api/v1/digital-employees/tasks/auto-assign              # 自动分配
PUT  /api/v1/digital-employees/tasks/:assignment_id/complete   # 完成任务
PUT  /api/v1/digital-employees/tasks/:assignment_id/fail       # 标记失败
GET  /api/v1/digital-employees/:employee_id/tasks             # 任务列表
```

### 绩效统计

```
GET /api/v1/digital-employees/:employee_id/performance  # 员工绩效
GET /api/v1/digital-employees/performance/team          # 团队绩效
```

## 数据库设计

### digital_employee_profiles (员工画像表)

| 字段 | 类型 | 说明 |
|------|------|------|
| employee_id | VARCHAR(36) | 主键 |
| tenant_id | VARCHAR(36) | 租户ID |
| name | VARCHAR(100) | 员工名称 |
| avatar | VARCHAR(255) | 头像URL |
| role | ENUM | 角色：customer_service、sales等 |
| skills | JSON | 技能标签数组 |
| personality | VARCHAR(100) | 性格特征 |
| knowledge_base_id | VARCHAR(36) | 知识库ID |
| bot_id | VARCHAR(36) | 关联Bot ID |
| status | ENUM | 状态：active、inactive、deleted |

### digital_employee_task_assignments (任务分配表)

| 字段 | 类型 | 说明 |
|------|------|------|
| assignment_id | VARCHAR(36) | 主键 |
| task_id | VARCHAR(36) | 任务ID |
| employee_id | VARCHAR(36) | 员工ID |
| tenant_id | VARCHAR(36) | 租户ID |
| task_type | ENUM | 任务类型 |
| priority | ENUM | 优先级：high、medium、low |
| status | ENUM | 状态：assigned、in_progress、completed、failed |
| assigned_at | BIGINT | 分配时间（毫秒） |
| completed_at | BIGINT | 完成时间（毫秒） |
| result | TEXT | 任务结果 |
| failure_reason | VARCHAR(255) | 失败原因 |

### digital_employee_performance (绩效统计表)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键（自增） |
| employee_id | VARCHAR(36) | 员工ID |
| tenant_id | VARCHAR(36) | 租户ID |
| total_tasks | INT | 总任务数 |
| completed_tasks | INT | 完成任务数 |
| failed_tasks | INT | 失败任务数 |
| completion_rate | DECIMAL(5,2) | 完成率（%） |
| avg_response_time | INT | 平均响应时间（秒） |
| customer_rating | DECIMAL(3,2) | 客户评分（1-5分） |
| period | ENUM | 统计周期：daily、weekly、monthly |
| date | DATE | 统计日期 |

## 错误码定义

| 错误码 | 说明 |
|--------|------|
| DIGITAL_EMPLOYEE_30001 | 员工不存在 |
| DIGITAL_EMPLOYEE_30002 | 员工名称已存在 |
| DIGITAL_EMPLOYEE_30003 | 无效的员工角色 |
| DIGITAL_EMPLOYEE_30004 | 员工未激活 |
| DIGITAL_EMPLOYEE_30005 | 缺少必需技能 |
| DIGITAL_EMPLOYEE_31001 | 任务分配不存在 |
| DIGITAL_EMPLOYEE_31002 | 任务已分配 |
| DIGITAL_EMPLOYEE_31003 | 无效的任务状态 |
| DIGITAL_EMPLOYEE_31004 | 未找到可用员工 |
| DIGITAL_EMPLOYEE_31005 | 任务已完成 |
| DIGITAL_EMPLOYEE_32001 | 绩效记录不存在 |
| DIGITAL_EMPLOYEE_32002 | 无效的统计周期 |
| DIGITAL_EMPLOYEE_32003 | 绩效计算失败 |

## 使用示例

### 创建员工画像

```go
req := &entity.CreateProfileRequest{
    TenantID:    "tenant-123",
    Name:        "智能客服小助手",
    Avatar:      "https://example.com/avatar.png",
    Role:        entity.EmployeeRoleCustomerService,
    Skills:      []string{"客户咨询", "问题解答", "订单查询"},
    Personality: "热情友好，专业高效",
    BotID:       &botID,
}

profile, err := service.CreateProfile(ctx, req)
```

### 分配任务

```go
req := &entity.AssignTaskRequest{
    TaskID:     "task-456",
    EmployeeID: "employee-123",
    TaskType:   entity.TaskTypeCustomerService,
    Priority:   entity.TaskPriorityHigh,
}

assignment, err := service.AssignTask(ctx, req)
```

### 自动分配任务

```go
req := &entity.AutoAssignTaskRequest{
    TaskID:         "task-789",
    TaskType:       entity.TaskTypeTechSupport,
    Priority:       entity.TaskPriorityMedium,
    RequiredSkills: []string{"技术支持", "故障排查"},
}

assignment, err := service.AutoAssignTask(ctx, req)
```

### 查询绩效

```go
req := &entity.GetPerformanceRequest{
    EmployeeID: "employee-123",
    Period:     entity.PerformancePeriodWeekly,
}

performance, err := service.GetPerformance(ctx, req)
fmt.Printf("完成率: %.2f%%\n", performance.CompletionRate)
```

## 测试

运行单元测试：

```bash
cd backend/domain/digital_employee
go test ./... -v -cover
```

运行覆盖率测试：

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 待完成功能

- [ ] 应用层（Application）集成
- [ ] 完整的错误处理和日志
- [ ] 性能优化（缓存、批量查询）
- [ ] 定时任务：自动更新绩效统计
- [ ] WebSocket推送任务状态变更
- [ ] 更多技能匹配算法
- [ ] 负载均衡策略优化

## 贡献指南

1. 遵循企业级开发规范
2. 所有代码必须有单元测试（覆盖率 ≥ 80%）
3. 使用统一错误码
4. API文档完整
5. 代码审查通过

## License

Copyright 2025 coze-dev Authors
