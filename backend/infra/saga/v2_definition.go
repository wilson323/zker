// Package saga 提供Saga分布式事务编排框架 v2.0
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

// V2SagaStatus Saga整体状态 (v2.0)
type V2SagaStatus string

const (
	V2SagaStatusPending      V2SagaStatus = "pending"       // 待执行
	V2SagaStatusRunning      V2SagaStatus = "running"       // 执行中
	V2SagaStatusCompleted    V2SagaStatus = "completed"     // 已完成
	V2SagaStatusFailed       V2SagaStatus = "failed"        // 失败
	V2SagaStatusCompensating V2SagaStatus = "compensating"  // 补偿中
	V2SagaStatusCompensated  V2SagaStatus = "compensated"   // 已补偿
)

// V2StepStatus 步骤状态 (v2.0)
type V2StepStatus string

const (
	V2StepStatusPending     V2StepStatus = "pending"     // 待执行
	V2StepStatusCompleted   V2StepStatus = "completed"   // 已完成
	V2StepStatusFailed      V2StepStatus = "failed"      // 执行失败
	V2StepStatusCompensated V2StepStatus = "compensated" // 已补偿
)

// V2Saga Saga定义 (v2.0)
type V2Saga struct {
	ID            string            // Saga唯一标识
	Name          string            // Saga名称
	Description   string            // Saga描述
	Steps         []V2SagaStep      // 执行步骤
	Compensations []V2CompensationStep // 补偿步骤
	Timeout       time.Duration     // 超时时间
	RetryPolicy   *V2RetryPolicy    // 重试策略
}

// V2RetryPolicy 重试策略 (v2.0)
type V2RetryPolicy struct {
	MaxAttempts     int           // 最大重试次数
	InitialInterval time.Duration // 初始重试间隔
	MaxInterval     time.Duration // 最大重试间隔
	Multiplier      float64       // 退避倍数
}

// DefaultV2RetryPolicy 默认重试策略
func DefaultV2RetryPolicy() *V2RetryPolicy {
	return &V2RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: 1 * time.Second,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
	}
}

// V2SagaStep Saga步骤接口 (v2.0)
type V2SagaStep interface {
	// Execute 执行步骤
	Execute(ctx context.Context, data interface{}) (interface{}, error)

	// Name 步骤名称
	Name() string

	// Timeout 超时时间
	Timeout() time.Duration
}

// V2CompensationStep 补偿步骤接口 (v2.0)
type V2CompensationStep interface {
	// Compensate 执行补偿
	Compensate(ctx context.Context, data interface{}) error

	// Name 补偿步骤名称
	Name() string

	// Timeout 超时时间
	Timeout() time.Duration
}

// V2SagaExecution Saga执行记录 (v2.0)
type V2SagaExecution struct {
	ID             string              // 执行记录ID
	SagaID         string              // Saga定义ID
	Status         V2SagaStatus        // 执行状态
	CurrentStep    int                 // 当前步骤索引
	InputData      interface{}         // 输入数据
	OutputData     interface{}         // 输出数据
	Error          error               // 错误信息
	StartedAt      time.Time           // 开始时间
	CompletedAt    *time.Time          // 完成时间
	StepExecutions []V2StepExecution   // 步骤执行记录
}

// V2StepExecution 步骤执行记录 (v2.0)
type V2StepExecution struct {
	StepName    string         // 步骤名称
	Status      V2StepStatus   // 执行状态
	Input       interface{}    // 输入数据
	Output      interface{}    // 输出数据
	Error       error          // 错误信息
	StartedAt   time.Time      // 开始时间
	CompletedAt *time.Time     // 完成时间
}

// V2SagaError Saga错误 (v2.0)
type V2SagaError struct {
	Code        string // 错误码
	Message     string // 错误消息
	SagaID      string // Saga ID
	StepName    string // 步骤名称
	Cause       error  // 原始错误
	Recoverable bool   // 是否可恢复
}

func (e *V2SagaError) Error() string {
	return e.Message
}

func (e *V2SagaError) Unwrap() error {
	return e.Cause
}

// 预定义错误
var (
	ErrV2SagaNotFound         = &V2SagaError{Code: "SAGA_NOT_FOUND", Message: "Saga定义不存在", Recoverable: false}
	ErrV2SagaTimeout          = &V2SagaError{Code: "SAGA_TIMEOUT", Message: "Saga执行超时", Recoverable: true}
	ErrV2StepFailed           = &V2SagaError{Code: "STEP_FAILED", Message: "步骤执行失败", Recoverable: true}
	ErrV2CompensationFailed   = &V2SagaError{Code: "COMPENSATION_FAILED", Message: "补偿执行失败", Recoverable: false}
	ErrV2InvalidSagaDefinition = &V2SagaError{Code: "INVALID_SAGA_DEFINITION", Message: "无效的Saga定义", Recoverable: false}
)
