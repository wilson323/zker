// Package saga 提供Saga分布式事务编排框架
//
// Saga模式是一种长活事务(Long Lived Transaction)模式,将分布式事务拆分为多个本地事务,
// 每个本地事务都有对应的补偿事务,保证最终一致性。
//
// 核心概念:
//   - Saga: 完整的分布式事务定义
//   - SagaStep: 单个执行步骤
//   - CompensationStep: 补偿步骤
//   - SagaExecution: Saga执行记录
//   - Orchestrator: Saga协调器
package saga

import (
	"context"
	"time"
)

// SagaStatus Saga整体状态
type SagaStatus string

const (
	SagaStatusPending      SagaStatus = "pending"       // 待执行
	SagaStatusRunning      SagaStatus = "running"       // 执行中
	SagaStatusCompleted    SagaStatus = "completed"     // 已完成
	SagaStatusFailed       SagaStatus = "failed"        // 失败
	SagaStatusCompensating SagaStatus = "compensating"  // 补偿中
	SagaStatusCompensated  SagaStatus = "compensated"   // 已补偿
)

// StepStatus 步骤状态
type StepStatus string

const (
	StepStatusPending      StepStatus = "pending"       // 待执行
	StepStatusRunning      StepStatus = "running"       // 执行中
	StepStatusCompleted    StepStatus = "completed"     // 已完成
	StepStatusFailed       StepStatus = "failed"        // 执行失败
	StepStatusCompensating StepStatus = "compensating"  // 补偿中
	StepStatusCompensated  StepStatus = "compensated"   // 已补偿
)

// Saga Saga定义
type Saga struct {
	ID            string            // Saga唯一标识
	Name          string            // Saga名称
	Description   string            // Saga描述
	Steps         []SagaStep        // 执行步骤
	Compensations []CompensationStep // 补偿步骤
	Timeout       time.Duration     // 超时时间
	RetryPolicy   *RetryPolicy      // 重试策略
}

// RetryPolicy 重试策略
type RetryPolicy struct {
	MaxAttempts     int           // 最大重试次数
	InitialInterval time.Duration // 初始重试间隔
	MaxInterval     time.Duration // 最大重试间隔
	Multiplier      float64       // 退避倍数
}

// DefaultRetryPolicy 默认重试策略
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: 1 * time.Second,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
	}
}

// SagaStep Saga步骤接口
type SagaStep interface {
	// Execute 执行步骤
	Execute(ctx context.Context, data interface{}) (interface{}, error)

	// Name 步骤名称
	Name() string

	// Timeout 超时时间
	Timeout() time.Duration
}

// CompensationStep 补偿步骤接口
type CompensationStep interface {
	// Compensate 执行补偿
	Compensate(ctx context.Context, data interface{}) error

	// Name 补偿步骤名称
	Name() string

	// Timeout 超时时间
	Timeout() time.Duration
}

// SagaExecution Saga执行记录
type SagaExecution struct {
	ID             string           // 执行记录ID
	SagaID         string           // Saga定义ID
	Status         SagaStatus       // 执行状态
	CurrentStep    int              // 当前步骤索引
	InputData      interface{}      // 输入数据
	OutputData     interface{}      // 输出数据
	Error          error            // 错误信息
	StartedAt      time.Time        // 开始时间
	CompletedAt    *time.Time       // 完成时间
	StepExecutions []StepExecution  // 步骤执行记录
}

// StepExecution 步骤执行记录
type StepExecution struct {
	StepName    string       // 步骤名称
	Status      StepStatus   // 执行状态
	Input       interface{}  // 输入数据
	Output      interface{}  // 输出数据
	Error       error        // 错误信息
	StartedAt   time.Time    // 开始时间
	CompletedAt *time.Time   // 完成时间
}

// SagaError Saga错误
type SagaError struct {
	Code        string // 错误码
	Message     string // 错误消息
	SagaID      string // Saga ID
	StepName    string // 步骤名称
	Cause       error  // 原始错误
	Recoverable bool   // 是否可恢复
}

func (e *SagaError) Error() string {
	return e.Message
}

func (e *SagaError) Unwrap() error {
	return e.Cause
}

// 预定义错误
var (
	ErrSagaNotFound         = &SagaError{Code: "SAGA_NOT_FOUND", Message: "Saga定义不存在", Recoverable: false}
	ErrSagaTimeout          = &SagaError{Code: "SAGA_TIMEOUT", Message: "Saga执行超时", Recoverable: true}
	ErrStepFailed           = &SagaError{Code: "STEP_FAILED", Message: "步骤执行失败", Recoverable: true}
	ErrCompensationFailed   = &SagaError{Code: "COMPENSATION_FAILED", Message: "补偿执行失败", Recoverable: false}
	ErrInvalidSagaDefinition = &SagaError{Code: "INVALID_SAGA_DEFINITION", Message: "无效的Saga定义", Recoverable: false}
)
