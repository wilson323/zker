# 分布式事务方案 - Saga模式

> **文档编号**: ARCH-DIST-2025-001
> **文档类型**: 架构设计方案
> **创建日期**: 2025-12-30
> **版本**: v1.0
> **密级**: 内部公开

---

## 📋 目录

- [1. 方案概述](#1-方案概述)
- [2. Saga模式设计](#2-saga模式设计)
- [3. 核心业务流程Saga实现](#3-核心业务流程saga实现)
- [4. 异常处理与补偿](#4-异常处理与补偿)
- [5. 实施指南](#5-实施指南)

---

## 1. 方案概述

### 1.1 背景与问题

在ZKER微服务架构中，多个业务流程涉及跨服务事务：

**典型场景**：
- **租户注册流程** - 涉及租户服务、组织服务、用户服务、通知服务
- **Bot创建流程** - 涉及Bot服务、知识库服务、权限服务
- **工作流执行** - 涉及对话服务、LLM服务、插件服务、知识库服务

**挑战**：
- 无法使用传统的两阶段提交（2PC）
- 需要保证最终一致性
- 需要支持补偿回滚

### 1.2 Saga模式简介

**Saga模式**是一种长活事务（Long Lived Transaction）模式，将一个分布式事务拆分为多个本地事务，每个本地事务都有对应的补偿事务。

**核心思想**：
1. 将分布式事务拆分为多个本地事务
2. 按顺序执行每个本地事务
3. 如果某个本地事务失败，执行前面已完成的本地事务的补偿操作
4. 最终达到一致状态

**优势**：
- ✅ 无需分布式锁，性能高
- ✅ 支持长事务（秒级到分钟级）
- ✅ 最终一致性保证
- ✅ 易于理解和实施

### 1.3 Saga模式类型

**1. 编排式（Choreography）**
- 通过事件驱动，服务间发布/订阅事件
- 去中心化，每个服务监听事件并执行本地事务
- 适合简单场景

**2. 协调式（Orchestration）**
- 由中央协调器（Saga协调器）控制整个流程
- 集中管理，易于监控和调试
- 适合复杂场景

**ZKER选择**：**协调式Saga**（Orchestration）

---

## 2. Saga模式设计

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Saga协调器架构                            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  【业务服务】                                                 │
│  ├── 租户服务 (tenant-service)                              │
│  ├── 组织服务 (organization-service)                        │
│  ├── 用户服务 (user-service)                                │
│  ├── Bot服务 (bot-service)                                  │
│  └── 通知服务 (notification-service)                        │
│                                                              │
│  【Saga协调器层】                                             │
│  ├── Saga定义注册                                            │
│  ├── Saga执行引擎                                            │
│  ├── 补偿事务管理                                            │
│  ├── 状态持久化                                              │
│  └── 事件发布/订阅                                           │
│                                                              │
│  【基础设施层】                                               │
│  ├── MySQL (Saga状态存储)                                    │
│  ├── Redis (缓存)                                           │
│  └── NSQ (事件总线)                                         │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 核心接口定义

```go
package saga

import (
    "context"
    "time"
)

// Saga定义
type Saga struct {
    ID              string
    Name            string
    Description     string
    Steps           []SagaStep
    Compensations   []CompensationStep
    Timeout         time.Duration
    RetryPolicy     *RetryPolicy
}

// Saga步骤
type SagaStep interface {
    // 执行步骤
    Execute(ctx context.Context, data interface{}) (interface{}, error)
    // 步骤名称
    Name() string
    // 超时时间
    Timeout() time.Duration
}

// 补偿步骤
type CompensationStep interface {
    // 执行补偿
    Compensate(ctx context.Context, data interface{}) error
    // 补偿步骤名称
    Name() string
    // 超时时间
    Timeout() time.Duration
}

// 重试策略
type RetryPolicy struct {
    MaxAttempts int           // 最大重试次数
    InitialInterval time.Duration  // 初始重试间隔
    MaxInterval     time.Duration  // 最大重试间隔
    Multiplier      float64        // 退避倍数
}

// Saga状态
type SagaStatus string

const (
    SagaStatusPending    SagaStatus = "pending"     // 待执行
    SagaStatusRunning    SagaStatus = "running"     // 执行中
    SagaStatusCompleted  SagaStatus = "completed"   // 已完成
    SagaStatusFailed     SagaStatus = "failed"      // 失败
    SagaStatusCompensating SagaStatus = "compensating" // 补偿中
    SagaStatusCompensated SagaStatus = "compensated" // 已补偿
)

// Saga执行记录
type SagaExecution struct {
    ID              string
    SagaID          string
    Status          SagaStatus
    CurrentStep     int
    InputData       interface{}
    OutputData      interface{}
    Error           error
    StartedAt       time.Time
    CompletedAt     *time.Time
    StepExecutions  []StepExecution
}

// 步骤执行记录
type StepExecution struct {
    StepName        string
    Status          string
    Input           interface{}
    Output          interface{}
    Error           error
    StartedAt       time.Time
    CompletedAt     *time.Time
}
```

### 2.3 Saga协调器实现

```go
package saga

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Saga协调器
type SagaCoordinator struct {
    sagaRepo    SagaRepository
    eventBus    EventBus
    logger      Logger
}

// 创建Saga协调器
func NewSagaCoordinator(
    sagaRepo SagaRepository,
    eventBus EventBus,
    logger Logger,
) *SagaCoordinator {
    return &SagaCoordinator{
        sagaRepo: sagaRepo,
        eventBus: eventBus,
        logger:   logger,
    }
}

// 定义Saga
func (c *SagaCoordinator) DefineSaga(saga *Saga) error {
    return c.sagaRepo.Save(context.Background(), saga)
}

// 执行Saga
func (c *SagaCoordinator) ExecuteSaga(
    ctx context.Context,
    sagaName string,
    input interface{},
) (*SagaExecution, error) {
    // 1. 加载Saga定义
    saga, err := c.sagaRepo.FindByName(ctx, sagaName)
    if err != nil {
        return nil, fmt.Errorf("saga not found: %w", err)
    }

    // 2. 创建Saga执行记录
    execution := &SagaExecution{
        ID:           generateID(),
        SagaID:       saga.ID,
        Status:       SagaStatusRunning,
        CurrentStep:  0,
        InputData:    input,
        StartedAt:    time.Now(),
        StepExecutions: make([]StepExecution, 0, len(saga.Steps)),
    }

    // 3. 保存执行记录
    if err := c.sagaRepo.SaveExecution(ctx, execution); err != nil {
        return nil, err
    }

    // 4. 依次执行每个步骤
    var lastOutput interface{}
    var compensationData []interface{}

    for i, step := range saga.Steps {
        c.logger.Infof("执行Saga步骤: %s", step.Name())

        // 准备输入数据
        stepInput := input
        if i > 0 {
            stepInput = lastOutput
        }

        // 执行步骤
        stepCtx, cancel := context.WithTimeout(ctx, step.Timeout())
        defer cancel()

        output, err := c.executeStepWithRetry(
            stepCtx,
            step,
            stepInput,
            saga.RetryPolicy,
        )

        // 记录步骤执行
        stepExec := StepExecution{
            StepName:  step.Name(),
            Status:    "completed",
            Input:     stepInput,
            Output:    output,
            StartedAt: time.Now(),
        }

        if err != nil {
            stepExec.Status = "failed"
            stepExec.Error = err
            completedAt := time.Now()
            stepExec.CompletedAt = &completedAt
            execution.StepExecutions = append(execution.StepExecutions, stepExec)

            // 步骤执行失败，开始补偿
            c.logger.Errorf("Saga步骤执行失败: %s, 错误: %v", step.Name(), err)
            execution.Status = SagaStatusCompensating
            execution.Error = err

            if err := c.sagaRepo.UpdateExecution(ctx, execution); err != nil {
                c.logger.Errorf("更新Saga执行记录失败: %v", err)
            }

            // 执行补偿
            if err := c.compensate(ctx, execution, compensationData); err != nil {
                c.logger.Errorf("Saga补偿失败: %v", err)
                execution.Status = SagaStatusFailed
            } else {
                execution.Status = SagaStatusCompensated
            }

            return execution, fmt.Errorf("saga执行失败: %w", err)
        }

        completedAt := time.Now()
        stepExec.CompletedAt = &completedAt
        execution.StepExecutions = append(execution.StepExecutions, stepExec)

        // 保存补偿数据
        compensationData = append(compensationData, output)

        // 更新最后输出
        lastOutput = output
        execution.CurrentStep = i + 1

        // 更新执行记录
        if err := c.sagaRepo.UpdateExecution(ctx, execution); err != nil {
            c.logger.Errorf("更新Saga执行记录失败: %v", err)
        }
    }

    // 所有步骤执行成功
    execution.Status = SagaStatusCompleted
    execution.OutputData = lastOutput
    completedAt := time.Now()
    execution.CompletedAt = &completedAt

    if err := c.sagaRepo.UpdateExecution(ctx, execution); err != nil {
        c.logger.Errorf("更新Saga执行记录失败: %v", err)
    }

    c.logger.Infof("Saga执行成功: %s", sagaName)

    return execution, nil
}

// 带重试的步骤执行
func (c *SagaCoordinator) executeStepWithRetry(
    ctx context.Context,
    step SagaStep,
    input interface{},
    retryPolicy *RetryPolicy,
) (interface{}, error) {
    var lastErr error

    for attempt := 0; attempt <= retryPolicy.MaxAttempts; attempt++ {
        if attempt > 0 {
            c.logger.Infof("重试步骤: %s, 第%d次重试", step.Name(), attempt)

            // 计算退避时间
            backoff := time.Duration(float64(retryPolicy.InitialInterval) *
                float64(1<<uint(attempt-1)) * retryPolicy.Multiplier)
            if backoff > retryPolicy.MaxInterval {
                backoff = retryPolicy.MaxInterval
            }

            select {
            case <-time.After(backoff):
            case <-ctx.Done():
                return nil, ctx.Err()
            }
        }

        output, err := step.Execute(ctx, input)
        if err == nil {
            return output, nil
        }

        lastErr = err
        c.logger.Errorf("步骤执行失败: %s, 错误: %v", step.Name(), err)
    }

    return nil, fmt.Errorf("步骤执行失败，已重试%d次: %w", retryPolicy.MaxAttempts, lastErr)
}

// 补偿执行
func (c *SagaCoordinator) compensate(
    ctx context.Context,
    execution *SagaExecution,
    compensationData []interface{},
) error {
    c.logger.Infof("开始Saga补偿: %s", execution.ID)

    // 加载Saga定义
    saga, err := c.sagaRepo.FindByID(ctx, execution.SagaID)
    if err != nil {
        return err
    }

    // 反向执行补偿步骤
    for i := execution.CurrentStep - 1; i >= 0; i-- {
        compensation := saga.Compensations[i]
        data := compensationData[i]

        c.logger.Infof("执行补偿步骤: %s", compensation.Name())

        compCtx, cancel := context.WithTimeout(ctx, compensation.Timeout())
        defer cancel()

        if err := compensation.Compensate(compCtx, data); err != nil {
            c.logger.Errorf("补偿步骤执行失败: %s, 错误: %v", compensation.Name(), err)
            // 继续执行后续补偿
        } else {
            c.logger.Infof("补偿步骤执行成功: %s", compensation.Name())
        }
    }

    return nil
}
```

### 2.4 Saga仓储实现

```go
package repository

import (
    "context"
    "encoding/json"
    "time"
)

// Saga仓储接口
type SagaRepository interface {
    // Saga定义管理
    Save(ctx context.Context, saga *saga.Saga) error
    FindByID(ctx context.Context, id string) (*saga.Saga, error)
    FindByName(ctx context.Context, name string) (*saga.Saga, error)
    List(ctx context.Context) ([]*saga.Saga, error)

    // Saga执行管理
    SaveExecution(ctx context.Context, execution *saga.SagaExecution) error
    UpdateExecution(ctx context.Context, execution *saga.SagaExecution) error
    FindExecutionByID(ctx context.Context, id string) (*saga.SagaExecution, error)
    ListExecutionsBySaga(ctx context.Context, sagaID string) ([]*saga.SagaExecution, error)
}

// MySQL实现
type MySQLSagaRepository struct {
    db *gorm.DB
}

func NewMySQLSagaRepository(db *gorm.DB) *MySQLSagaRepository {
    return &MySQLSagaRepository{db: db}
}

func (r *MySQLSagaRepository) Save(ctx context.Context, saga *saga.Saga) error {
    // 序列化步骤和补偿
    stepsJSON, _ := json.Marshal(saga.Steps)
    compensationsJSON, _ := json.Marshal(saga.Compensations)

    record := &SagaRecord{
        ID:            saga.ID,
        Name:          saga.Name,
        Description:   saga.Description,
        Steps:         string(stepsJSON),
        Compensations: string(compensationsJSON),
        Timeout:       int64(saga.Timeout),
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    return r.db.WithContext(ctx).Create(record).Error
}

func (r *MySQLSagaRepository) FindByName(ctx context.Context, name string) (*saga.Saga, error) {
    var record SagaRecord
    err := r.db.WithContext(ctx).
        Where("name = ?", name).
        First(&record).
        Error

    if err != nil {
        return nil, err
    }

    return record.ToSaga()
}

func (r *MySQLSagaRepository) SaveExecution(
    ctx context.Context,
    execution *saga.SagaExecution,
) error {
    inputJSON, _ := json.Marshal(execution.InputData)
    stepsJSON, _ := json.Marshal(execution.StepExecutions)

    record := &SagaExecutionRecord{
        ID:           execution.ID,
        SagaID:       execution.SagaID,
        Status:       string(execution.Status),
        CurrentStep:  execution.CurrentStep,
        InputData:    string(inputJSON),
        StepExecutions: string(stepsJSON),
        StartedAt:    execution.StartedAt,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    if execution.CompletedAt != nil {
        record.CompletedAt = *execution.CompletedAt
    }

    if execution.Error != nil {
        record.Error = execution.Error.Error()
    }

    return r.db.WithContext(ctx).Create(record).Error
}

func (r *MySQLSagaRepository) UpdateExecution(
    ctx context.Context,
    execution *saga.SagaExecution,
) error {
    inputJSON, _ := json.Marshal(execution.InputData)
    stepsJSON, _ := json.Marshal(execution.StepExecutions)

    updates := map[string]interface{}{
        "status":          string(execution.Status),
        "current_step":    execution.CurrentStep,
        "input_data":      string(inputJSON),
        "step_executions": string(stepsJSON),
        "updated_at":      time.Now(),
    }

    if execution.CompletedAt != nil {
        updates["completed_at"] = *execution.CompletedAt
    }

    if execution.Error != nil {
        updates["error"] = execution.Error.Error()
    }

    return r.db.WithContext(ctx).
        Model(&SagaExecutionRecord{}).
        Where("id = ?", execution.ID).
        Updates(updates).
        Error
}

// 数据库表定义
type SagaRecord struct {
    ID            string `gorm:"primaryKey"`
    Name          string `gorm:"unique;not null"`
    Description   string
    Steps         string `gorm:"type:text"`
    Compensations string `gorm:"type:text"`
    Timeout       int64
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

type SagaExecutionRecord struct {
    ID             string `gorm:"primaryKey"`
    SagaID         string `gorm:"not null;index"`
    Status         string `gorm:"not null;index"`
    CurrentStep    int
    InputData      string `gorm:"type:text"`
    OutputData     string `gorm:"type:text"`
    Error          string `gorm:"type:text"`
    StepExecutions string `gorm:"type:text"`
    StartedAt      time.Time `gorm:"index"`
    CompletedAt    *time.Time
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

---

## 3. 核心业务流程Saga实现

### 3.1 租户注册Saga

**业务场景**：企业自助注册ZKER平台

**涉及服务**：
1. 租户服务 - 创建租户记录
2. 组织服务 - 创建默认组织
3. 用户服务 - 创建管理员账户
4. 通知服务 - 发送欢迎邮件

```go
package saga

import (
    "context"
    "time"
)

// 租户注册Saga定义
func NewTenantRegistrationSaga() *Saga {
    return &Saga{
        ID:          "saga-tenant-registration",
        Name:        "租户注册Saga",
        Description: "处理企业租户注册流程，包括创建租户、组织、管理员账户和发送欢迎邮件",
        Timeout:     5 * time.Minute,
        Steps: []SagaStep{
            &CreateTenantStep{},
            &CreateDefaultOrganizationStep{},
            &CreateAdminUserStep{},
            &SendWelcomeEmailStep{},
        },
        Compensations: []CompensationStep{
            &DeleteAdminUserStep{},
            &DeleteDefaultOrganizationStep{},
            &DeleteTenantStep{},
            &SendCancellationEmailStep{},
        },
        RetryPolicy: &RetryPolicy{
            MaxAttempts:     3,
            InitialInterval: 1 * time.Second,
            MaxInterval:     10 * time.Second,
            Multiplier:      2.0,
        },
    }
}

// ============ 步骤实现 ============

// 步骤1: 创建租户
type CreateTenantStep struct{}

func (s *CreateTenantStep) Name() string {
    return "创建租户"
}

func (s *CreateTenantStep) Timeout() time.Duration {
    return 30 * time.Second
}

func (s *CreateTenantStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    cmd := data.(*RegisterTenantCommand)

    tenant := &Tenant{
        ID:             generateTenantID(),
        Name:           cmd.EnterpriseName,
        Industry:       cmd.Industry,
        Scale:          cmd.Scale,
        Subdomain:      cmd.Subdomain,
        Status:         TenantStatusTrial,
        TrialEndsAt:    time.Now().Add(14 * 24 * time.Hour),
        ContactName:    cmd.ContactName,
        ContactEmail:   cmd.ContactEmail,
        ContactPhone:   cmd.ContactPhone,
    }

    if err := tenantRepo.Save(ctx, tenant); err != nil {
        return nil, fmt.Errorf("创建租户失败: %w", err)
    }

    return tenant, nil
}

// 补偿1: 删除租户
type DeleteTenantStep struct{}

func (s *DeleteTenantStep) Name() string {
    return "删除租户"
}

func (s *DeleteTenantStep) Timeout() time.Duration {
    return 30 * time.Second
}

func (s *DeleteTenantStep) Compensate(ctx context.Context, data interface{}) error {
    tenant := data.(*Tenant)

    if err := tenantRepo.Delete(ctx, tenant.ID); err != nil {
        return fmt.Errorf("删除租户失败: %w", err)
    }

    return nil
}

// 步骤2: 创建默认组织
type CreateDefaultOrganizationStep struct{}

func (s *CreateDefaultOrganizationStep) Name() string {
    return "创建默认组织"
}

func (s *CreateDefaultOrganizationStep) Timeout() time.Duration {
    return 30 * time.Second
}

func (s *CreateDefaultOrganizationStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    tenant := data.(*Tenant)

    org := &Organization{
        ID:       generateOrgID(),
        TenantID: tenant.ID,
        Name:     tenant.Name + "-默认组织",
        Type:     OrgTypeCompany,
        ParentID: nil,
        Path:     "/" + generateOrgID(),
        Level:    1,
        Status:   OrgStatusActive,
    }

    if err := orgRepo.Save(ctx, org); err != nil {
        return nil, fmt.Errorf("创建默认组织失败: %w", err)
    }

    return org, nil
}

// 补偿2: 删除默认组织
type DeleteDefaultOrganizationStep struct{}

func (s *DeleteDefaultOrganizationStep) Name() string {
    return "删除默认组织"
}

func (s *DeleteDefaultOrganizationStep) Timeout() time.Duration {
    return 30 * time.Second
}

func (s *DeleteDefaultOrganizationStep) Compensate(ctx context.Context, data interface{}) error {
    org := data.(*Organization)

    if err := orgRepo.Delete(ctx, org.ID); err != nil {
        return fmt.Errorf("删除默认组织失败: %w", err)
    }

    return nil
}

// 步骤3: 创建管理员账户
type CreateAdminUserStep struct{}

func (s *CreateAdminUserStep) Name() string {
    return "创建管理员账户"
}

func (s *CreateAdminUserStep) Timeout() time.Duration {
    return 30 * time.Second
}

func (s *CreateAdminUserStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    org := data.(*Organization)

    user := &User{
        ID:           generateUserID(),
        TenantID:     org.TenantID,
        Name:         org.Tenant.ContactName,
        Email:        org.Tenant.ContactEmail,
        Phone:        org.Tenant.ContactPhone,
        PasswordHash: hashPassword(generateRandomPassword()),
        OrganizationID: org.ID,
        Role:         UserRoleAdmin,
        Status:       UserStatusActive,
    }

    if err := userRepo.Save(ctx, user); err != nil {
        return nil, fmt.Errorf("创建管理员账户失败: %w", err)
    }

    return user, nil
}

// 补偿3: 删除管理员账户
type DeleteAdminUserStep struct{}

func (s *DeleteAdminUserStep) Name() string {
    return "删除管理员账户"
}

func (s *DeleteAdminUserStep) Timeout() time.Duration {
    return 30 * time.Second
}

func (s *DeleteAdminUserStep) Compensate(ctx context.Context, data interface{}) error {
    user := data.(*User)

    if err := userRepo.Delete(ctx, user.ID); err != nil {
        return fmt.Errorf("删除管理员账户失败: %w", err)
    }

    return nil
}

// 步骤4: 发送欢迎邮件
type SendWelcomeEmailStep struct{}

func (s *SendWelcomeEmailStep) Name() string {
    return "发送欢迎邮件"
}

func (s *SendWelcomeEmailStep) Timeout() time.Duration {
    return 30 * time.Second
}

func (s *SendWelcomeEmailStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    user := data.(*User)

    email := &Email{
        To:       user.Email,
        Subject:  "欢迎使用ZKER企业级AI平台",
        Template: "welcome",
        Data: map[string]interface{}{
            "UserName":  user.Name,
            "TenantName": user.Tenant.Name,
            "LoginURL":  fmt.Sprintf("https://%s.zker.com", user.Tenant.Subdomain),
        },
    }

    if err := emailService.Send(ctx, email); err != nil {
        return nil, fmt.Errorf("发送欢迎邮件失败: %w", err)
    }

    return nil, nil
}

// 补偿4: 发送取消通知邮件
type SendCancellationEmailStep struct{}

func (s *SendCancellationEmailStep) Name() string {
    return "发送取消通知邮件"
}

func (s *SendCancellationEmailStep) Timeout() time.Duration {
    return 30 * time.Second
}

func (s *SendCancellationEmailStep) Compensate(ctx context.Context, data interface{}) error {
    user := data.(*User)

    email := &Email{
        To:       user.Email,
        Subject:  "ZKER注册失败通知",
        Template: "registration_failed",
        Data: map[string]interface{}{
            "UserName": user.Name,
            "Reason":  "系统错误，注册未完成",
        },
    }

    // 发送失败邮件，不阻塞补偿流程
    _ = emailService.Send(ctx, email)

    return nil
}
```

### 3.2 Bot创建Saga

**业务场景**：用户创建新的Bot

**涉及服务**：
1. Bot服务 - 创建Bot记录
2. 知识库服务 - 创建默认知识库
3. 权限服务 - 分配Bot权限

```go
package saga

// Bot创建Saga定义
func NewBotCreationSaga() *Saga {
    return &Saga{
        ID:          "saga-bot-creation",
        Name:        "Bot创建Saga",
        Description: "处理Bot创建流程，包括创建Bot、默认知识库和权限分配",
        Timeout:     3 * time.Minute,
        Steps: []SagaStep{
            &CreateBotStep{},
            &CreateDefaultKnowledgeBaseStep{},
            &AssignBotPermissionsStep{},
        },
        Compensations: []CompensationStep{
            &RevokeBotPermissionsStep{},
            &DeleteDefaultKnowledgeBaseStep{},
            &DeleteBotStep{},
        },
        RetryPolicy: &RetryPolicy{
            MaxAttempts:     3,
            InitialInterval: 1 * time.Second,
            MaxInterval:     10 * time.Second,
            Multiplier:      2.0,
        },
    }
}

// 步骤1: 创建Bot
type CreateBotStep struct{}

func (s *CreateBotStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    cmd := data.(*CreateBotCommand)

    bot := &Bot{
        ID:          generateBotID(),
        TenantID:    cmd.TenantID,
        Name:        cmd.Name,
        Description: cmd.Description,
        Avatar:      cmd.Avatar,
        Type:        cmd.Type,
        Status:      BotStatusDraft,
        CreatedBy:   cmd.CreatorID,
    }

    if err := botRepo.Save(ctx, bot); err != nil {
        return nil, fmt.Errorf("创建Bot失败: %w", err)
    }

    return bot, nil
}

// 步骤2: 创建默认知识库
type CreateDefaultKnowledgeBaseStep struct{}

func (s *CreateDefaultKnowledgeBaseStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    bot := data.(*Bot)

    kb := &KnowledgeBase{
        ID:          generateKbID(),
        TenantID:    bot.TenantID,
        BotID:       bot.ID,
        Name:        bot.Name + "-默认知识库",
        Type:        KBTypePrivate,
        Status:      KBStatusActive,
        CreatedBy:   bot.CreatedBy,
    }

    if err := kbRepo.Save(ctx, kb); err != nil {
        return nil, fmt.Errorf("创建默认知识库失败: %w", err)
    }

    return kb, nil
}

// 步骤3: 分配Bot权限
type AssignBotPermissionsStep struct{}

func (s *AssignBotPermissionsStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    kb := data.(*KnowledgeBase)

    permissions := []Permission{
        {
            ResourceID: kb.BotID,
            ResourceType: "bot",
            UserID:      kb.CreatedBy,
            Role:        "owner",
            Actions:     []string{"read", "write", "delete", "share"},
        },
    }

    if err := permissionRepo.BatchGrant(ctx, permissions); err != nil {
        return nil, fmt.Errorf("分配Bot权限失败: %w", err)
    }

    return nil, nil
}

// 补偿步骤实现省略...
```

---

## 4. 异常处理与补偿

### 4.1 异常类型

```go
package saga

// Saga异常类型
type SagaError struct {
    Code       string
    Message    string
    SagaID     string
    StepName   string
    Cause      error
    Recoverable bool
}

func (e *SagaError) Error() string {
    return fmt.Sprintf("[%s] %s: %s", e.Code, e.StepName, e.Message)
}

// 预定义错误
var (
    ErrSagaNotFound      = &SagaError{Code: "SAGA_NOT_FOUND", Recoverable: false}
    ErrSagaTimeout       = &SagaError{Code: "SAGA_TIMEOUT", Recoverable: true}
    ErrStepFailed        = &SagaError{Code: "STEP_FAILED", Recoverable: true}
    ErrCompensationFailed = &SagaError{Code: "COMPENSATION_FAILED", Recoverable: false}
)
```

### 4.2 补偿策略

```go
package saga

// 补偿策略
type CompensationStrategy int

const (
    // 自动补偿：失败时自动执行所有补偿步骤
    CompensationAuto CompensationStrategy = iota

    // 手动补偿：失败时记录状态，等待人工介入
    CompensationManual

    // 部分补偿：某些步骤失败时跳过补偿
    CompensationPartial
)

// 补偿配置
type CompensationConfig struct {
    Strategy      CompensationStrategy
    MaxRetries    int
    RetryInterval time.Duration
    OnFailure     func(*SagaExecution) error
}

// 应用补偿策略
func (c *SagaCoordinator) applyCompensationStrategy(
    ctx context.Context,
    execution *SagaExecution,
    config *CompensationConfig,
) error {
    switch config.Strategy {
    case CompensationAuto:
        return c.compensate(ctx, execution, nil)

    case CompensationManual:
        execution.Status = SagaStatusCompensating
        c.sagaRepo.UpdateExecution(ctx, execution)
        // 发送告警通知
        c.alertService.SendAlert(ctx, &Alert{
            Type:    "saga_compensation_required",
            SagaID:  execution.SagaID,
            Message: "Saga执行失败，需要手动补偿",
        })
        return nil

    case CompensationPartial:
        return c.partialCompensate(ctx, execution)

    default:
        return fmt.Errorf("未知的补偿策略: %v", config.Strategy)
    }
}

// 部分补偿（跳过可忽略的步骤）
func (c *SagaCoordinator) partialCompensate(
    ctx context.Context,
    execution *SagaExecution,
) error {
    saga, _ := c.sagaRepo.FindByID(ctx, execution.SagaID)

    for i := execution.CurrentStep - 1; i >= 0; i-- {
        compensation := saga.Compensations[i]

        // 检查是否可以跳过补偿
        if c.canSkipCompensation(compensation) {
            c.logger.Infof("跳过补偿步骤: %s", compensation.Name())
            continue
        }

        // 执行补偿
        data := execution.StepExecutions[i].Input
        compCtx, cancel := context.WithTimeout(ctx, compensation.Timeout())
        defer cancel()

        if err := compensation.Compensate(compCtx, data); err != nil {
            c.logger.Errorf("补偿步骤执行失败: %s", compensation.Name())
        }
    }

    return nil
}

func (c *SagaCoordinator) canSkipCompensation(compensation CompensationStep) bool {
    // 某些步骤（如发送邮件）可以跳过补偿
    switch compensation.Name() {
    case "发送欢迎邮件":
        return true
    default:
        return false
    }
}
```

### 4.3 Saga监控

```go
package saga

// Saga监控指标
type SagaMetrics struct {
    TotalExecutions    int64
    SuccessfulExecutions int64
    FailedExecutions   int64
    CompensatedExecutions int64
    AverageDuration    time.Duration
    StepFailureRates   map[string]float64
}

// 监控服务
type SagaMonitor struct {
    metricsRepo SagaMetricsRepository
    prometheus  PrometheusClient
}

// 记录Saga执行完成
func (m *SagaMonitor) RecordExecution(execution *SagaExecution) {
    metrics := &SagaMetrics{
        TotalExecutions: 1,
    }

    switch execution.Status {
    case SagaStatusCompleted:
        metrics.SuccessfulExecutions = 1
    case SagaStatusFailed:
        metrics.FailedExecutions = 1
    case SagaStatusCompensated:
        metrics.CompensatedExecutions = 1
    }

    duration := execution.CompletedAt.Sub(execution.StartedAt)
    metrics.AverageDuration = duration

    // 发送到Prometheus
    m.prometheus.RecordSagaExecution(metrics)

    // 保存到数据库
    m.metricsRepo.Save(metrics)
}
```

---

## 5. 实施指南

### 5.1 使用Saga协调器

```go
package main

import (
    "context"
    "log"
)

func main() {
    // 1. 初始化依赖
    db := initMySQL()
    redis := initRedis()
    nsq := initNSQ()

    sagaRepo := repository.NewMySQLSagaRepository(db)
    eventBus := messaging.NewNSQEventBus(nsq)
    logger := logger.NewZapLogger()

    // 2. 创建Saga协调器
    coordinator := saga.NewSagaCoordinator(sagaRepo, eventBus, logger)

    // 3. 注册Saga定义
    coordinator.DefineSaga(saga.NewTenantRegistrationSaga())
    coordinator.DefineSaga(saga.NewBotCreationSaga())

    // 4. 执行Saga
    cmd := &RegisterTenantCommand{
        EnterpriseName: "示例公司",
        ContactEmail:   "admin@example.com",
        Subdomain:      "example-company",
        // ...
    }

    execution, err := coordinator.ExecuteSaga(context.Background(), "租户注册Saga", cmd)
    if err != nil {
        log.Printf("Saga执行失败: %v", err)
        return
    }

    log.Printf("Saga执行成功: %s", execution.ID)
}
```

### 5.2 测试Saga

```go
package saga_test

import (
    "context"
    "testing"
    "time"
)

func TestTenantRegistrationSaga_Success(t *testing.T) {
    // 准备测试环境
    testDB := setupTestDB()
    defer cleanupTestDB(testDB)

    // 创建Mock服务
    mockTenantRepo := &MockTenantRepository{}
    mockOrgRepo := &MockOrganizationRepository{}
    mockUserRepo := &MockUserRepository{}
    mockEmailService := &MockEmailService{}

    // 创建Saga协调器
    coordinator := saga.NewSagaCoordinator(
        mockTenantRepo,
        mockOrgRepo,
        mockUserRepo,
        mockEmailService,
    )

    // 执行Saga
    cmd := &RegisterTenantCommand{
        EnterpriseName: "测试公司",
        ContactEmail:   "test@example.com",
        Subdomain:      "test-company",
    }

    execution, err := coordinator.ExecuteSaga(context.Background(), "租户注册Saga", cmd)

    // 验证结果
    assert.NoError(t, err)
    assert.Equal(t, SagaStatusCompleted, execution.Status)
    assert.Equal(t, 4, len(execution.StepExecutions))

    // 验证每个步骤都成功
    for _, stepExec := range execution.StepExecutions {
        assert.Equal(t, "completed", stepExec.Status)
        assert.NoError(t, stepExec.Error)
    }
}

func TestTenantRegistrationSaga_Compensation(t *testing.T) {
    // 准备测试环境
    testDB := setupTestDB()
    defer cleanupTestDB(testDB)

    // 创建Mock服务（用户服务会失败）
    mockTenantRepo := &MockTenantRepository{}
    mockOrgRepo := &MockOrganizationRepository{}
    mockUserRepo := &MockUserRepository{FailOnCreate: true}
    mockEmailService := &MockEmailService{}

    coordinator := saga.NewSagaCoordinator(
        mockTenantRepo,
        mockOrgRepo,
        mockUserRepo,
        mockEmailService,
    )

    // 执行Saga
    cmd := &RegisterTenantCommand{
        EnterpriseName: "测试公司",
        ContactEmail:   "test@example.com",
        Subdomain:      "test-company",
    }

    execution, err := coordinator.ExecuteSaga(context.Background(), "租户注册Saga", cmd)

    // 验证补偿
    assert.Error(t, err)
    assert.Equal(t, SagaStatusCompensated, execution.Status)
    assert.Equal(t, 2, len(execution.StepExecutions)) // 前2步成功，第3步失败

    // 验证补偿步骤执行
    assert.True(t, mockTenantRepo.Deleted)
    assert.True(t, mockOrgRepo.Deleted)
}
```

### 5.3 部署配置

```yaml
# docker-compose.yml

version: '3.8'

services:
  saga-coordinator:
    image: zker/saga-coordinator:latest
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - REDIS_HOST=redis
      - NSQ_NSQLOOKUPD=nsqlookupd:4160
    depends_on:
      - mysql
      - redis
      - nsqlookupd

  mysql:
    image: mysql:8.4.5
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: zker_saga
    volumes:
      - mysql_data:/var/lib/mysql

  redis:
    image: redis:8.0
    volumes:
      - redis_data:/data

  nsqlookupd:
    image: nsqio/nsq
    command: /nsqlookupd
    ports:
      - "4160:4160"
      - "4161:4161"

  nsqd:
    image: nsqio/nsq
    command: /nsqd --lookupd-tcp-address=nsqlookupd:4160
    ports:
      - "4150:4150"
      - "4151:4151"

volumes:
  mysql_data:
  redis_data:
```

---

## 6. 总结

### 6.1 核心要点

1. **Saga模式选择** - ZKER采用协调式Saga（Orchestration）
2. **补偿事务** - 每个步骤都有对应的补偿操作
3. **最终一致性** - 保证所有服务最终达到一致状态
4. **监控告警** - 完整的监控指标和告警机制

### 6.2 适用场景

✅ **适合使用Saga的场景**：
- 跨服务业务流程（租户注册、Bot创建）
- 长事务（秒级到分钟级）
- 可以接受最终一致性
- 需要补偿回滚

❌ **不适合使用Saga的场景**：
- 强一致性要求（如支付）
- 短事务（毫秒级）
- 无法补偿的操作（如发送短信后无法回收）

### 6.3 最佳实践

1. **幂等性** - 所有步骤和补偿操作都必须是幂等的
2. **超时控制** - 每个步骤都要设置合理的超时时间
3. **重试机制** - 失败步骤自动重试，避免瞬时故障
4. **监控告警** - 实时监控Saga执行状态，及时发现问题
5. **人工介入** - 复杂补偿流程支持人工介入处理

---

**文档版本**: v1.0
**最后更新**: 2025-12-30
**作者**: ZKER架构团队
